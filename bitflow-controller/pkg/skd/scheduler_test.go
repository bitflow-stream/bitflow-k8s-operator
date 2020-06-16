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

func getNodeAffinityForPod(pod *corev1.Pod) string {
	return pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values[0]
}

func (s *SkdTestSuite) Test_shouldSetAffinityOnAllPods() {
	node1 := s.Node("node1")
	node2 := s.Node("node2")
	pod1 := s.Pod("pod1")
	pod2 := s.Pod("pod2")
	pod3 := s.Pod("pod3")
	pod4 := s.Pod("pod4")
	pod5 := s.Pod("pod5")

	scheduler := Scheduler{
		nodes: []*NodeData{
			{
				node: node1,
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
			{
				node: node2,
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*corev1.Pod{pod1, pod2, pod3, pod4, pod5},
	}

	err := scheduler.setNodeAffinityForPods()

	s.Nil(err)
	s.Equal(node1.Name, getNodeAffinityForPod(scheduler.pods[0]))
	s.Equal(node1.Name, getNodeAffinityForPod(scheduler.pods[1]))
	s.Equal(node2.Name, getNodeAffinityForPod(scheduler.pods[2]))
	s.Equal(node2.Name, getNodeAffinityForPod(scheduler.pods[3]))
	s.Equal(node1.Name, getNodeAffinityForPod(scheduler.pods[4]))
}

func (s *SkdTestSuite) Test_shouldSetAffinityOnAllPods_differenceInInitialNumberOfPodSlots() {
	node1 := s.Node("node1")
	node2 := s.Node("node2")
	pod1 := s.Pod("pod1")
	pod2 := s.Pod("pod2")
	pod3 := s.Pod("pod3")
	pod4 := s.Pod("pod4")
	pod5 := s.Pod("pod5")
	pod6 := s.Pod("pod6")
	pod7 := s.Pod("pod7")
	pod8 := s.Pod("pod8")

	scheduler := Scheduler{
		nodes: []*NodeData{
			{
				node: node1,
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				initialNumberOfPodSlots: 1,
				podSlotScalingFactor:    3,
				resourceLimit:           0.1,
			},
			{
				node: node2,
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				initialNumberOfPodSlots: 1,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*corev1.Pod{pod1, pod2, pod3, pod4, pod5, pod6, pod7, pod8},
	}

	err := scheduler.setNodeAffinityForPods()

	s.Nil(err)
	s.Equal(node1.Name, getNodeAffinityForPod(scheduler.pods[0]))
	s.Equal(node2.Name, getNodeAffinityForPod(scheduler.pods[1]))
	s.Equal(node2.Name, getNodeAffinityForPod(scheduler.pods[2]))
	s.Equal(node1.Name, getNodeAffinityForPod(scheduler.pods[3]))
	s.Equal(node1.Name, getNodeAffinityForPod(scheduler.pods[4]))
	s.Equal(node2.Name, getNodeAffinityForPod(scheduler.pods[5]))
	s.Equal(node2.Name, getNodeAffinityForPod(scheduler.pods[6]))
	s.Equal(node2.Name, getNodeAffinityForPod(scheduler.pods[7]))
}
