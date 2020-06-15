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

func getNextHigherNumberOfPodSlots(incrementFactor int64, value int64) (int64, error) {
	if value < incrementFactor {
		return incrementFactor, nil
	}
	count := incrementFactor
	for true {
		if count >= value {
			return count, nil
		}
		count *= incrementFactor
	}
	return -1, errors.New("should never happen")
}

func GetNumberOfPodSlotsAllocatedForNodeAfterAddingPods(node SimulatedNode, numberOfPodsToAdd int64) int64 {
	incrementFactor := node.nodeData.podSlotScalingFactor
	numberOfPodsOnNode := int64(len(node.pods))
	slots, _ := getNextHigherNumberOfPodSlots(incrementFactor, numberOfPodsOnNode+numberOfPodsToAdd)
	return slots
}

// lower is better
func CalculatePenaltyForNodeAfterAddingPods(simulatedNode SimulatedNode, numberOfPodsToAdd int64) float64 {
	R := getAllocatableCpu(*simulatedNode.nodeData.node) * simulatedNode.nodeData.resourceLimit / float64(GetNumberOfPodSlotsAllocatedForNodeAfterAddingPods(simulatedNode, numberOfPodsToAdd))

	return CalculateExecutionTime(R, simulatedNode.nodeData.curve)
}

func getLowestPenaltyNode(simulatedNodes []*SimulatedNode) *SimulatedNode {
	minPenaltyIndex := -1
	minPenalty := -1.0
	for i, simulatedNode := range simulatedNodes {
		penalty := CalculatePenaltyForNodeAfterAddingPods(*simulatedNode, 1)
		if minPenaltyIndex < 0 || penalty < minPenalty {
			minPenaltyIndex = i
			minPenalty = penalty
		}
	}
	return simulatedNodes[minPenaltyIndex]
}
