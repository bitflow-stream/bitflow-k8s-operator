package bitflow

import (
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const podPenalty = 5

func getPodPenalty() int {
	return podPenalty
}

// lower is better
func calculatePenaltyForNode(cli client.Client, nodeName string) (int, error) {
	count, err := common.GetNumberOfPodsForNode(cli, nodeName)

	penalty := count * podPenalty

	return penalty, err
}
