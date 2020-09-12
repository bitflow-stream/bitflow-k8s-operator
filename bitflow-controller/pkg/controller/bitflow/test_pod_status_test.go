package bitflow

import (
	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
)

type PodStatusTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestPodStatus() {
	s.SubTestSuite(new(PodStatusTestSuite))
}

func (s *PodStatusTestSuite) TestEmptyManagedPods() {
	pods := NewManagedPods()
	pod := s.Pod("xx")

	s.NotPanics(func() {
		pods.CleanupStep("xx", map[string]bool{})
		pods.UpdateExistingPod(pod)
		pods.MarkRespawning(pod, true)
	})
	s.Zero(pods.Len())
	s.Empty(pods.ListRespawningPods())
}

func (s *PodStatusTestSuite) TestManagedPods() {
	pod1, pod2, pod3 := s.Pod("x1"), s.Pod("x2"), s.Pod("x3")
	common.SetTargetNode(pod1, s.Node("nodeA"))
	common.SetTargetNode(pod2, s.Node("nodeA"))
	common.SetTargetNode(pod3, s.Node("nodeB"))
	missingPod := s.Pod("missing")
	step1, step2, step3 := s.Step("s1", bitflowv1.StepTypeOneToOne), s.Step("s2", bitflowv1.StepTypeOneToOne), s.Step("s3", bitflowv1.StepTypeOneToOne)
	source1, source2 := s.Source("source1", map[string]string{"a": "b"}), s.Source("source2", map[string]string{"x": "y"})

	pods := NewManagedPods()
	pods.Put(pod1, step1, []*bitflowv1.BitflowSource{source1, source2})
	pods.Put(pod2, step2, []*bitflowv1.BitflowSource{source1})
	pods.Put(pod3, step3, []*bitflowv1.BitflowSource{})

	// Initial state
	s.NotPanics(func() {
		pods.CleanupStep(missingPod.Name, map[string]bool{})
		pods.UpdateExistingPod(missingPod)
		pods.MarkRespawning(missingPod, true)
	})
	s.Equal(3, pods.Len())
	s.Empty(pods.ListRespawningPods())

	// Modifications: replace pod1, add pod4, update pod3, mark pod2 and pod3 respawning
	pod4 := s.Pod("x4")
	pod1Replacement := pod1.DeepCopy()
	common.SetTargetNode(pod4, s.Node("nodeB"))
	common.SetTargetNode(pod1Replacement, s.Node("nodeB"))
	pods.Put(pod4, step3, []*bitflowv1.BitflowSource{})
	pods.Put(pod1Replacement, step2, []*bitflowv1.BitflowSource{source1})
	pods.MarkRespawning(pod2, true)
	pods.MarkRespawning(pod3, true)
	pod3Updated := pod3.DeepCopy()
	common.SetTargetNode(pod3Updated, s.Node("nodeC"))
	pods.UpdateExistingPod(pod3Updated)

	// Updated state
	s.Equal(4, pods.Len())
	s.True(pods.pods[pod2.Name].respawning)
	s.True(pods.pods[pod3.Name].respawning)
	s.Len(pods.ListRespawningPods(), 2)

	// Remove respawning flag
	pods.MarkRespawning(pod3, false)
	s.Len(pods.ListRespawningPods(), 1)
	pods.MarkRespawning(pod2, false)
	s.Empty(pods.ListRespawningPods())
}
