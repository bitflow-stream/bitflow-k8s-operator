package scheduler

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/config"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _staticSchedulerRNG = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

// StaticScheduler implements the Scheduler interface without changing already existing pod placements.
// Newly created pods are placed on a random node.
type StaticScheduler struct {
	// Maps pod Names to Node Names. Empty node name means the pod is not yet scheduled.
	CurrentState map[string]string

	// Nodes contains a list of all available (ready) nodes where pods can be placed.
	Nodes []string

	// SourceAffinityClient can be set to a non-nil value to enable placing pods as close as possible to their data source.
	// If the source cannot be determined, or if SourceAffinityClient is nil, new pods are placed on a random node.
	SourceAffinityClient client.Client

	// If SourceAffinityClient is set, the below variables must be set as well
	Config    *config.Config
	Namespace string
}

func (s *StaticScheduler) Schedule() (bool, map[string]string, error) {
	if err := s.validate(); err != nil {
		return false, nil, err
	}

	changed := false
	result := make(map[string]string)
	for pod, node := range s.CurrentState {
		if node == "" {
			node = s.selectNode(pod)
			changed = true
		}
		result[pod] = node
	}
	return changed, result, nil
}

func (s *StaticScheduler) validate() error {
	if len(s.Nodes) == 0 {
		return errors.New("Need at least one node to schedule on")
	}
	if s.SourceAffinityClient != nil && (s.Namespace == "" || s.Config == nil) {
		return errors.New("StaticScheduler with valid SourceAffinityClient also needs Namespace and Config variables")
	}
	return nil
}

func (s *StaticScheduler) selectNode(pod string) string {
	node, err := s.selectNodeWithSourceAffinity(pod)
	if err != nil {
		log.WithField("pod", pod).Warnln("Failed to select node based on source affinity:", err)
	}
	if node == "" {
		// Fallback: pick random node
		node = s.selectRandomNode()
	}
	return node
}

func (s *StaticScheduler) selectRandomNode() string {
	index := _staticSchedulerRNG.Int31n(int32(len(s.Nodes)))
	return s.Nodes[index]
}

func (s *StaticScheduler) selectNodeWithSourceAffinity(pod string) (string, error) {
	if s.SourceAffinityClient == nil {
		return "", nil
	}

	sources, err := s.getSourcesForPod(pod)
	if err != nil {
		return "", err
	}

	logger := log.WithField("pod", pod)
	switch len(sources) {
	case 0:
		return "", nil
	case 1:
		return s.findNodeForDataSource(sources[0], logger)
	default:
		return s.findNodeForDataSources(sources, logger)
	}
}

func (s *StaticScheduler) getSourcesForPod(pod string) ([]*bitflowv1.BitflowSource, error) {
	// TODO first get the step, then all currently matched sources.
	// Until this is implemented, treat all steps as if they have no data sources
	return nil, nil
}

func (s *StaticScheduler) findNodeForDataSource(source *bitflowv1.BitflowSource, logger *log.Entry) (string, error) {
	nodeLabel := s.Config.GetStandaloneSourceLabel()
	if nodeName, ok := source.Labels[nodeLabel]; ok {
		node, err := common.RequestReadyNode(s.SourceAffinityClient, nodeName)
		if err == nil {
			logger.Debugf("%v has label %v=%v, scheduling on node %v", source, nodeLabel, nodeName, node.Name)
			return node.Name, nil
		} else {
			return "", fmt.Errorf("Error requesting node '%v' based on label %v=%v of %v: %v", nodeName, nodeLabel, nodeName, source, err)
		}
	}
	return s.findNodeForDataSourcePod(source, logger)
}

func (s *StaticScheduler) findNodeForDataSourcePod(source *bitflowv1.BitflowSource, logger *log.Entry) (string, error) {
	podName := source.Labels[bitflowv1.SourceLabelPodName]
	if podName == "" {
		logger.Debugf("%v has no label %v, cannot use it for scheduling", source, bitflowv1.SourceLabelPodName)
		return "", nil
	}
	pod, err := common.RequestPod(s.SourceAffinityClient, podName, s.Namespace)
	if err != nil {
		return "", fmt.Errorf("Error requesting pod '%v' based on label %v=%v of %v: %v",
			podName, bitflowv1.SourceLabelPodName, podName, source, err)
	}
	node, err := common.RequestReadyNode(s.SourceAffinityClient, common.GetNodeName(pod))
	if err == nil {
		return node.Name, nil
	}
	return "", err
}

func (s *StaticScheduler) findNodeForDataSources(sources []*bitflowv1.BitflowSource, logger *log.Entry) (string, error) {
	sourcesOnNodes := make(map[string]int)

	var i int
	for _, source := range sources {
		node, err := s.findNodeForDataSource(source, logger)
		if err != nil || node == "" {
			if err != nil {
				logger.Warnln(err)
			}
			continue
		}
		sourcesOnNodes[node] = sourcesOnNodes[node] + 1
	}
	i = 0
	var maxNode string
	for node, numSources := range sourcesOnNodes {
		if numSources > i {
			i = numSources
			maxNode = node
		}
	}
	if i == 0 {
		return "", nil
	}
	return maxNode, nil
}
