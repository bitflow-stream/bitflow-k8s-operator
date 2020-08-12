package skd

import (
	"errors"
)

type Scheduler interface {
	Schedule() (map[string]string, error)
}

type EqualDistributionScheduler struct {
	nodeNames []string
	podNames  []string
}

func (eds EqualDistributionScheduler) Schedule() (map[string]string, error) {
	if err := validateEqualDistributionScheduler(eds); err != nil {
		return nil, err
	}

	m := make(map[string]string)

	if len(eds.podNames) == 0 {
		return m, nil
	}

	var nodeIndex = 0
	for _, podName := range eds.podNames {
		if nodeIndex >= len(eds.nodeNames) {
			nodeIndex = 0
		}

		m[podName] = eds.nodeNames[nodeIndex]

		nodeIndex++
	}

	return m, nil
}

type AdvancedScheduler struct {
	nodes []*NodeData
	pods  []*PodData
}

func (as AdvancedScheduler) Schedule() (map[string]string, error) {
	if err := validateAdvancedScheduler(as); err != nil {
		return nil, err
	}

	// TODO implement

	return nil, nil
}

type NodeData struct {
	name                    string
	allocatableCpu          int // TODO which unit? Data is found here: node.Status.Allocatable.Cpu().MilliValue()
	memory                  int // memory in MB
	initialNumberOfPodSlots int
	podSlotScalingFactor    int
	resourceLimit           float64
}

type PodData struct {
	name             string
	receivesDataFrom []string // list of pod names
	curve            Curve
	minimumMemory    int // memory in MB
}

type Curve struct {
	a, b, c, d float64
}

type NodeState struct {
	node *NodeData
	pods []*PodData
}

type SystemState struct {
	nodes []*NodeState
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
			return errors.New("resourceLimit is 0")
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
	return nil
}
