package skd

import (
	"errors"
	"math"
)

func CalculateExecutionTime(cpus float64, curve Curve) float64 {
	return curve.a*math.Pow(cpus+curve.b, -curve.c) + curve.d
}

func GetNumberOfPodSlots(nodeData *NodeData, numberOfPods int) (int, error) {
	initialNumberOfPodSlots := nodeData.initialNumberOfPodSlots
	if numberOfPods < initialNumberOfPodSlots {
		return initialNumberOfPodSlots, nil
	}
	updatedNumberOfPodSlots := initialNumberOfPodSlots
	for true {
		if updatedNumberOfPodSlots >= numberOfPods {
			return updatedNumberOfPodSlots, nil
		}
		updatedNumberOfPodSlots *= nodeData.podSlotScalingFactor
	}
	return -1, errors.New("should never happen")
}

func CalculatePenalty(state SystemState) (float64, error) {
	var penalty = 0.0

	for _, nodeState := range state.nodes {
		nodeData := nodeState.node

		numberOfPodSlots, err := GetNumberOfPodSlots(nodeData, len(nodeState.pods))
		if err != nil {
			return -1, err
		}

		availableCpus := float64(nodeData.allocatableCpu) * nodeData.resourceLimit / float64(numberOfPodSlots)

		for _, podData := range nodeState.pods {
			penalty += CalculateExecutionTime(availableCpus, podData.curve)
		}
	}

	return penalty, nil
}
