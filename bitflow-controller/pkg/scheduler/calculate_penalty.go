package scheduler

import (
	"errors"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/config"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"math"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func getStepCurveParameters(stepName string) (float64, float64, float64, float64) {
	// dummy data until steps contain curve parameters
	a := 6.71881241016441
	b := 0.0486498280492762
	c := 2.0417306475862214
	d := 15.899403720950454
	return a, b, c, d
}

func getAllocatableCpu(node corev1.Node) float64 {
	// TODO MilliValue() is correct, for memory use Value()
	return float64(node.Status.Allocatable.Cpu().MilliValue())
}

func GetTotalResourceLimit(node corev1.Node, config *config.Config) float64 {
	return resources.RequestBitflowResourceLimitByNode(&node, config)
}

func getNumberOfPodsForNode(client client.Client, nodeName string) int {
	count, err := common.GetNumberOfPodsForNode(client, nodeName)
	if err != nil {
		return 0
	}
	return count
}

func getNextHigherNumberOfPodSlots(bufferInitSize float64, incrementFactor float64, value float64) (float64, error) {
	println(value, "<", bufferInitSize)
	if value < bufferInitSize {
		return bufferInitSize, nil
	}
	count := incrementFactor
	for true {
		println(count, ">=", value)
		if count >= value {
			return count, nil
		}
		count *= incrementFactor
	}
	return -1, errors.New("Should never happen")
}

func GetNumberOfPodSlotsAllocatedForNodeAfterAddingPods(client client.Client, config *config.Config, nodeName string, numberOfPodsToAdd float64) float64 {
	bufferInitSize := float64(config.GetInitialResourceBufferSize())
	incrementFactor := config.GetResourceBufferIncrementFactor()
	numberOfPodsOnNode := float64(getNumberOfPodsForNode(client, nodeName))
	println(numberOfPodsOnNode + numberOfPodsToAdd) // TODO remove
	slots, _ := getNextHigherNumberOfPodSlots(bufferInitSize, incrementFactor, numberOfPodsOnNode+numberOfPodsToAdd)
	return slots
}

func CalculateExecutionTime(cpus float64) float64 {
	a, b, c, d := getStepCurveParameters("some-step-name")
	return a*math.Pow(cpus+b, -c) + d
}

// lower is better
func CalculatePenaltyForNode(client client.Client, config *config.Config, node corev1.Node) (float64, error) {
	return CalculatePenaltyForNodeAfterAddingPods(client, config, node, 0)
}

// lower is better
func CalculatePenaltyForNodeAfterAddingPods(client client.Client, config *config.Config, node corev1.Node, numberOfPodsToAdd float64) (float64, error) {
	R := getAllocatableCpu(node) * GetTotalResourceLimit(node, config) / GetNumberOfPodSlotsAllocatedForNodeAfterAddingPods(client, config, node.Name, numberOfPodsToAdd)

	// TODO error in return necessary?
	return CalculateExecutionTime(R), nil
}
