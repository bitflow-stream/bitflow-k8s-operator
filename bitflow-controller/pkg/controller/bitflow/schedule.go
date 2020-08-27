package bitflow

import (
	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	log "github.com/sirupsen/logrus"
)

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
