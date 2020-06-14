package skd

import (
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	"testing"
)

type SkdTestSuite struct {
	common.AbstractTestSuite
}

func TestSkd(t *testing.T) {
	suite.Run(t, new(SkdTestSuite))
}

func (s *SkdTestSuite) Test_shouldSetAffinityOnAllPods() {
	node1 := s.Node("Some Node Name")
	pod1 := s.Pod("Some Pod Name")
	pod2 := s.Pod("Some Other Pod Name")

	scheduler := Scheduler{
		nodes: []NodeData{
			{
				node:                    node1,
				curve:                   Curve{},
				initialNumberOfPodSlots: 0,
				podSlotScalingFactor:    0,
				resourceLimit:           0,
			},
		},
		pods: []*corev1.Pod{pod1, pod2},
	}

	scheduler.setNodeAffinityForPods()

	actualNodeAffinityPod1 := scheduler.pods[0].Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values[0]
	actualNodeAffinityPod2 := scheduler.pods[1].Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values[0]
	s.Equal(node1.Name, actualNodeAffinityPod1)
	s.Equal(node1.Name, actualNodeAffinityPod2)
}
