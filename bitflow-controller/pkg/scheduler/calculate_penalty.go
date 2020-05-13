package scheduler

import (
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
	// TODO .Value() oder .MilliValue()?
	return float64(node.Status.Allocatable.Cpu().Value())
}

func getTotalResourceLimit(node corev1.Node, config *config.Config) float64 {
	return resources.RequestBitflowResourceLimitByNode(&node, config)
}

func getMaxPods(node corev1.Node) float64 {
	return 8.0 // TODO get actual value
}

func CalculateExecutionTime(cpus float64) float64 {
	a, b, c, d := getStepCurveParameters("some-step-name")
	return a*math.Pow(cpus+b, -c) + d
}

// lower is better
func CalculatePenaltyForNode(cli client.Client, config *config.Config, node corev1.Node) (float64, error) {
	R := getAllocatableCpu(node) * getTotalResourceLimit(node, config) / getMaxPods(node)

	// TODO error in return necessary?
	return CalculateExecutionTime(R), nil
}
