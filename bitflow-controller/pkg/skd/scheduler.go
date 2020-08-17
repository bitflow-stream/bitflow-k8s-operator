package skd

import (
	"errors"
)

type Scheduler interface {
	Schedule() (bool, map[string]string, error)
}

type EqualDistributionScheduler struct {
	nodeNames []string
	podNames  []string
}

func (eds EqualDistributionScheduler) Schedule() (bool, map[string]string, error) {
	if err := validateEqualDistributionScheduler(eds); err != nil {
		return false, nil, err
	}

	m := make(map[string]string)

	if len(eds.podNames) == 0 {
		return true, m, nil
	}

	var nodeIndex = 0
	for _, podName := range eds.podNames {
		if nodeIndex >= len(eds.nodeNames) {
			nodeIndex = 0
		}

		m[podName] = eds.nodeNames[nodeIndex]

		nodeIndex++
	}

	return true, m, nil
}

type AdvancedScheduler struct {
	nodes              []*NodeData
	pods               []*PodData
	networkPenalty     float64
	memoryPenalty      float64
	thresholdPercent   float64
	previousScheduling map[string]string
}

func (as AdvancedScheduler) findBestSchedulingCheckingAllPermutations(state SystemState, podsLeft []*PodData) (SystemState, float64, error) {
	if len(podsLeft) == 0 {
		penalty, err := CalculatePenalty(state, as.networkPenalty, as.memoryPenalty)
		return state, penalty, err
	}

	currentPod := podsLeft[0]

	var lowestPenalty float64 = -1
	var lowestPenaltySystemState SystemState

	for i, nodeState := range state.nodes {
		nodeState.pods = append(nodeState.pods, currentPod)
		state.nodes[i] = nodeState
		newSystemState, currentPenalty, err := as.findBestSchedulingCheckingAllPermutations(state, podsLeft[1:])
		if err != nil {
			continue
		}
		if lowestPenalty == -1 || currentPenalty < lowestPenalty {
			lowestPenalty = currentPenalty

			// copying "manually" to prevent lowestPenaltySystemState and newSystemState from having the same memory address, which leads to problems
			lowestPenaltySystemState = SystemState{}
			for _, newSystemStateNodeState := range newSystemState.nodes {
				lowestPenaltySystemState.nodes = append(lowestPenaltySystemState.nodes, newSystemStateNodeState)
			}
		}

		// deleting previously added pod in preparation for next iteration
		// copying "manually" to prevent errors -  state.nodes[i].pods = state.nodes[i].pods[:len(state.nodes[i].pods)-1] does NOT work
		var tempPods []*PodData
		for j, pod := range state.nodes[i].pods {
			if j == len(state.nodes[i].pods)-1 {
				break
			}
			tempPods = append(tempPods, pod)
		}
		state.nodes[i].pods = tempPods
	}

	if lowestPenalty == -1 {
		return SystemState{}, -1, errors.New("pod " + currentPod.name + " could not be scheduled onto any node")
	}

	return lowestPenaltySystemState, lowestPenalty, nil
}

func NewDistributionPenaltyLowerConsideringThreshold(previousPenalty float64, newPenalty float64, thresholdPercent float64) bool {
	var previousPenaltyMinusThreshold = previousPenalty * ((100 - thresholdPercent) / 100)
	if newPenalty <= previousPenaltyMinusThreshold {
		return true
	}
	return false
}

func (as AdvancedScheduler) getPreviousSystemState() SystemState {
	systemState := SystemState{[]NodeState{}}

	for _, node := range as.nodes {
		nodeState := NodeState{
			node: node,
			pods: []*PodData{},
		}
		for _, pod := range as.pods {
			if as.previousScheduling[pod.name] == node.name {
				nodeState.pods = append(nodeState.pods, pod)
			}
		}
		systemState.nodes = append(systemState.nodes, nodeState)
	}

	return systemState
}

func (as AdvancedScheduler) Schedule() (bool, map[string]string, error) {
	if err := validateAdvancedScheduler(as); err != nil {
		return false, nil, err
	}

	systemState := SystemState{[]NodeState{}}
	for _, node := range as.nodes {
		systemState.nodes = append(systemState.nodes, NodeState{
			node: node,
			pods: []*PodData{},
		})
	}

	bestDistributionState, bestDistributionPenalty, err := as.findBestSchedulingCheckingAllPermutations(systemState, as.pods)

	if as.previousScheduling != nil {
		previousPenalty, err := CalculatePenalty(as.getPreviousSystemState(), as.networkPenalty, as.memoryPenalty)
		if err == nil && !NewDistributionPenaltyLowerConsideringThreshold(previousPenalty, bestDistributionPenalty, as.thresholdPercent) {
			return false, nil, nil
		}
	}

	if err != nil {
		return false, nil, err
	}

	m := make(map[string]string)
	for _, nodeState := range bestDistributionState.nodes {
		nodeName := nodeState.node.name
		for _, pod := range nodeState.pods {
			m[pod.name] = nodeName
		}
	}

	return true, m, nil
}

type NodeData struct {
	name                    string
	allocatableCpu          float64 // 1000 == 1 CPU core
	memory                  float64 // memory in MB
	initialNumberOfPodSlots int
	podSlotScalingFactor    int
	resourceLimit           float64
}

type PodData struct {
	name             string
	receivesDataFrom []string // list of pod names
	curve            Curve
	minimumMemory    float64 // memory in MB
}

type Curve struct {
	a, b, c, d float64
}

type NodeState struct {
	node *NodeData
	pods []*PodData
}

func (state SystemState) toString() string {
	var str = "("

	for _, nodeState := range state.nodes {
		str += nodeState.node.name + "["

		for _, pod := range nodeState.pods {
			str += pod.name + " "
		}

		str += "] "
	}
	str += ")"
	return str
}

type SystemState struct {
	nodes []NodeState
}

func validateEqualDistributionScheduler(scheduler EqualDistributionScheduler) error {
	if len(scheduler.nodeNames) == 0 {
		return errors.New("no nodes in scheduler")
	}
	for _, name := range scheduler.nodeNames {
		if name == "" {
			return errors.New("empty name in nodeNames")
		}
	}
	if len(scheduler.podNames) == 0 {
		return errors.New("no pods in scheduler")
	}
	for _, name := range scheduler.podNames {
		if name == "" {
			return errors.New("empty name in podNames")
		}
	}
	return nil
}

func validateAdvancedScheduler(scheduler AdvancedScheduler) error {
	if len(scheduler.nodes) == 0 {
		return errors.New("no node data in scheduler")
	}
	for _, nodeData := range scheduler.nodes {
		if nodeData.name == "" {
			return errors.New("empty name in NodeData")
		}
		if nodeData.memory <= 0 {
			return errors.New("memory is <= 0")
		}
		if nodeData.initialNumberOfPodSlots <= 0 {
			return errors.New("initialNumberOfPodSlots is <= 0")
		}
		if nodeData.podSlotScalingFactor <= 0 {
			return errors.New("podSlotScalingFactor is <= 0")
		}
		if nodeData.resourceLimit <= 0 {
			return errors.New("resourceLimit is <= 0")
		}
	}
	for _, podData := range scheduler.pods {
		if podData.name == "" {
			return errors.New("empty name in PodData")
		}
		if podData.receivesDataFrom == nil {
			return errors.New("receivesDataFrom is nil")
		}
		if podData.curve == (Curve{}) {
			return errors.New("empty curve")
		}
		if podData.minimumMemory <= 0 {
			return errors.New("minimumMemory is <= 0")
		}
	}
	if scheduler.thresholdPercent < 0 || scheduler.thresholdPercent > 100 {
		return errors.New("thresholdPercent needs to be >= 0 and <= 100")
	}
	return nil
}
