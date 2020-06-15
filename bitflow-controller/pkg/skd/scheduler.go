package skd

import (
	"errors"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	corev1 "k8s.io/api/core/v1"
)

type Scheduler struct {
	nodes []*NodeData
	pods  []*corev1.Pod
}

type NodeData struct {
	node                    *corev1.Node
	curve                   Curve
	initialNumberOfPodSlots int64
	podSlotScalingFactor    int64
	resourceLimit           float64
}

type Curve struct {
	a, b, c, d float64
}

type SimulatedNode struct {
	nodeData *NodeData
	pods     []*corev1.Pod
}

func (s *Scheduler) setNodeAffinityForPods() error {
	err := validateScheduler(*s)
	if err != nil {
		return err
	}
	simulate(*s)
	return nil
}

func validateScheduler(scheduler Scheduler) error {
	if len(scheduler.nodes) == 0 {
		return errors.New("no node data in scheduler")
	}
	if len(scheduler.pods) == 0 {
		return errors.New("no pods in scheduler")
	}
	for _, nodeData := range scheduler.nodes {
		if nodeData.node == nil {
			return errors.New("no node in NodeData")
		}
		if nodeData.curve == (Curve{}) {
			return errors.New("empty curve")
		}
		if nodeData.resourceLimit == 0 {
			return errors.New("resourceLimit is 0")
		}
		if nodeData.initialNumberOfPodSlots == 0 {
			return errors.New("initialNumberOfPodSlots is 0")
		}
		if nodeData.podSlotScalingFactor == 0 {
			return errors.New("podSlotScalingFactor is 0")
		}
	}
	return nil
}

func simulate(scheduler Scheduler) {
	var simulatedNodes []*SimulatedNode
	for _, nodeData := range scheduler.nodes {
		simulatedNodes = append(simulatedNodes, &SimulatedNode{
			nodeData: nodeData,
			pods:     nil,
		})
	}

	for _, pod := range scheduler.pods {
		lowestPenaltyNode := getLowestPenaltyNode(simulatedNodes)
		lowestPenaltyNode.pods = append(lowestPenaltyNode.pods, pod)
	}

	for _, simulatedNode := range simulatedNodes {
		for _, pod := range simulatedNode.pods {
			common.SetTargetNode(simulatedNode.nodeData.node, pod)
		}
	}
}
