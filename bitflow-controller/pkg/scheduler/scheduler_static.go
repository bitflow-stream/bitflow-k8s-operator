package scheduler

import (
	"errors"
	"math/rand"
	"time"
)

var _staticSchedulerRNG = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))

// StaticScheduler implements the Scheduler interface without changing already existing pod placements.
// Newly created pods are placed on a random node.
type StaticScheduler struct {
	// Maps pod Names to Node Names. Empty node name means the pod is not yet scheduled.
	CurrentState map[string]string

	// Maps Pod names to nodes that contain the data sources for the pod. If a pod is contained in this map,
	// it is scheduled on one of the data-source-nodes. Otherwise, it is scheduled on a random node.
	DataSourceNodes map[string][]string

	// Nodes contains a list of all available (ready) nodes where pods can be placed.
	Nodes []string
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
	return nil
}

func (s *StaticScheduler) selectNode(pod string) string {
	node := s.selectNodeWithSourceAffinity(pod)
	if node == "" {
		// Fallback: pick random node
		node = s.selectRandomNode(s.Nodes)
	}
	return node
}

func (s *StaticScheduler) selectRandomNode(nodes []string) string {
	index := _staticSchedulerRNG.Int31n(int32(len(nodes)))
	return nodes[index]
}

func (s *StaticScheduler) selectNodeWithSourceAffinity(pod string) string {
	sourceNodes := s.DataSourceNodes[pod]
	if len(sourceNodes) > 0 {
		return s.selectRandomNode(sourceNodes)
	}
	return ""
}
