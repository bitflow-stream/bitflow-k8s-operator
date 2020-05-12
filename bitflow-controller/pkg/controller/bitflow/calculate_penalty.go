package bitflow

import (
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"math"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const podPenalty = 5

const a = -5
const b = 0
const c = -0.1
const d = 7

func getPodPenalty() int {
	return podPenalty
}

func calculateExecutionTime(cpus float64) float64 {
	return a*math.Pow(cpus+b, -c) + d
}

// lower is better
func calculatePenaltyForNode(cli client.Client, nodeName string) (int, error) {
	count, err := common.GetNumberOfPodsForNode(cli, nodeName)

	penalty := count * podPenalty

	return penalty, err
}
