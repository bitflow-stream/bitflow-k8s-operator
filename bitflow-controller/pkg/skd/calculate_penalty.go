package skd

import (
	"errors"
	corev1 "k8s.io/api/core/v1"
	"math"
)

func getAllocatableCpu(node corev1.Node) float64 {
	// TODO MilliValue() is correct, for memory use Value()
	return float64(node.Status.Allocatable.Cpu().MilliValue())
}

func CalculateExecutionTime(cpus float64, curve Curve) float64 {
	return curve.a*math.Pow(cpus+curve.b, -curve.c) + curve.d
}

func GetNumberOfPodSlotsAfterAddingPods(node SimulatedNode, numberOfPodsToAdd int64) (int64, error) {
	numberOfPods := int64(len(node.pods)) + numberOfPodsToAdd
	initialNumberOfPodSlots := node.nodeData.initialNumberOfPodSlots
	if numberOfPods < initialNumberOfPodSlots {
		return initialNumberOfPodSlots, nil
	}
	updatedNumberOfPodSlots := initialNumberOfPodSlots
	for true {
		if updatedNumberOfPodSlots >= numberOfPods {
			return updatedNumberOfPodSlots, nil
		}
		updatedNumberOfPodSlots *= node.nodeData.podSlotScalingFactor
	}
	return -1, errors.New("should never happen")
}

// lower is better
func CalculatePenaltyForNodeAfterAddingPods(simulatedNode SimulatedNode, numberOfPodsToAdd int64) float64 {
	numberOfPodSlots, _ := GetNumberOfPodSlotsAfterAddingPods(simulatedNode, numberOfPodsToAdd) // TODO handle error or remove
	R := getAllocatableCpu(*simulatedNode.nodeData.node) * simulatedNode.nodeData.resourceLimit / float64(numberOfPodSlots)

	return CalculateExecutionTime(R, simulatedNode.nodeData.curve)
}

func GetLowestPenaltyNode(simulatedNodes []*SimulatedNode) (*SimulatedNode, error) {
	if len(simulatedNodes) == 0 {
		return nil, errors.New("simulatedNodes are empty")
	}
	minPenaltyIndex := 0
	minPenalty := CalculatePenaltyForNodeAfterAddingPods(*simulatedNodes[0], 1)
	for i, simulatedNode := range simulatedNodes {
		if i == 0 {
			continue
		}
		penalty := CalculatePenaltyForNodeAfterAddingPods(*simulatedNode, 1)
		if penalty < minPenalty {
			minPenaltyIndex = i
			minPenalty = penalty
		}
	}
	return simulatedNodes[minPenaltyIndex], nil
}
