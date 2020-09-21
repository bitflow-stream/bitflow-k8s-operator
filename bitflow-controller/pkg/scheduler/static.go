package scheduler

import (
	"sort"
)

type StaticScheduler struct {
	// Maps pod Names to Node Names. Empty node name means the pod is not yet scheduled.
	CurrentState map[string]string

	// Maps Pod names to nodes that contain the data sources for the pod. If a pod is contained in this map,
	// it is scheduled on one of the data-source-nodes. Otherwise, it is scheduled on a random node.
	DataSourceNodes map[string][]string

	// AvailableNodes contains a list of all available (ready) nodes where each pod can be placed. For most pods, the list of available
	// nodes will be the same, but some pods may have restrictions regarding target nodes.
	AvailableNodes map[string][]string
}

type nodeSelectionFunc func(nodes []string, podsPerNode map[string]int) string

func (s *StaticScheduler) schedule(selectNode nodeSelectionFunc) (bool, map[string]string, error) {
	// Construct the current schedule and count number of pods on each node
	schedule := make(map[string]string)
	podsPerNode := make(map[string]int)
	var newPods []string
	for pod, node := range s.CurrentState {
		if node != "" {
			schedule[pod] = node
			podsPerNode[node] = podsPerNode[node] + 1
		} else {
			newPods = append(newPods, pod)
		}
	}

	// Make scheduling deterministic
	sort.Strings(newPods)

	// Select nodes for all new pods
	changed := false
	for _, pod := range newPods {
		availableNodes := s.DataSourceNodes[pod]
		sort.Strings(availableNodes) // Make scheduling deterministic
		node := selectNode(availableNodes, podsPerNode)
		if node == "" {
			availableNodes = s.AvailableNodes[pod]
			sort.Strings(availableNodes)
			node = selectNode(availableNodes, podsPerNode)
		}
		if node != "" {
			schedule[pod] = node
			podsPerNode[node] = podsPerNode[node] + 1
			changed = true
		}
	}
	return changed, schedule, nil
}

// RandomStaticScheduler implements the Scheduler interface without changing already existing pod placements.
// Newly created pods are placed on a random node, giving priority to data source nodes.
type RandomStaticScheduler struct {
	StaticScheduler
}

func (s *RandomStaticScheduler) Schedule() (bool, map[string]string, error) {
	return s.schedule(s.selectNode)
}

func (s *RandomStaticScheduler) selectNode(nodes []string, _ map[string]int) string {
	return selectRandomNode(nodes)
}

// LeastOccupiedStaticScheduler places new pods on the node with the least pods already scheduled on it. Pods are not
// moved after being scheduled once, meaning the load is not dynamically re-balanced.
// For meaning of the fields, see RandomStaticScheduler.
type LeastOccupiedStaticScheduler struct {
	StaticScheduler
	AvailableNodes map[string][]string
}

func (eds LeastOccupiedStaticScheduler) Schedule() (bool, map[string]string, error) {
	return eds.schedule(eds.selectNode)
}

func (eds LeastOccupiedStaticScheduler) selectNode(nodes []string, podsPerNode map[string]int) string {
	if len(nodes) == 0 {
		return ""
	}

	node := nodes[0]
	minPodsOnNode := podsPerNode[node]
	for _, availableNode := range nodes {
		numPods := podsPerNode[availableNode]
		if numPods < minPodsOnNode {
			minPodsOnNode = numPods
			node = availableNode
		}
	}
	return node
}
