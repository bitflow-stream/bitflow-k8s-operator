package skd

import (
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	corev1 "k8s.io/api/core/v1"
)

type Scheduler struct {
	nodes []NodeData
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

func (s *Scheduler) setNodeAffinityForPods() {
	for _, pod := range s.pods {
		common.SetTargetNode(s.nodes[0].node, pod)
	}
}
