package bitflow

import (
	"context"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type SchedulerTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestScheduler() {
	s.SubTestSuite(new(SchedulerTestSuite))
}

func (s *SchedulerTestSuite) TestGetAllNodes() {
	r := s.initReconciler(
		s.Node("node1"),
		s.Node("node2"),
		s.Node("node3"),
		s.Node("node4"),
		s.Node("node5"))

	var nodeList *corev1.NodeList
	var err error
	nodeList, err = common.RequestReadyNodes(r.client)
	s.NoError(err)
	s.Equal(5, len(nodeList.Items))
}

func (s *SchedulerTestSuite) TestGetNumberOfPodsForNode() {
	labels := map[string]string{"hello": "world"}
	r := s.initReconciler(
		s.Node("node1"),
		s.Source("source1", labels), s.Source("source2", labels),
		s.Source("source3", labels), s.Source("source4", labels),
		s.Step("step1", "", "hello", "world"))
	s.testReconcile(r, "step1")

	count, err := common.GetNumberOfPodsForNode(r.client, "node1")

	s.NoError(err)
	s.Equal(4, count)
}

func (s *SchedulerTestSuite) TestLeastContainersScheduler() {
	labels := map[string]string{"hello": "world"}
	r := s.initReconciler(
		s.Node("node1"), s.Node("node2"),
		s.Source("source1", labels), s.Source("source2", labels),
		s.Source("source3", labels), s.Source("source4", labels),
		s.StepCustomSchedulers("step1", "", "leastContainers", "hello", "world"))
	s.testReconcile(r, "step1")

	s.assertNumberOfPodsForNode(r.client, "node1", 2)
	s.assertNumberOfPodsForNode(r.client, "node2", 2)
}

func (s *SchedulerTestSuite) TestScheduling2StandaloneSources() {
	doTest := func(sourceNode string) {
		s.SubTest(sourceNode, func() {
			labels := map[string]string{"nodename": sourceNode, "hello": "world"}
			r := s.initReconciler(
				s.Node("node1"), s.Node("node2"), s.Node("node3"),
				s.Source("source1", labels), s.Source("source2", labels), s.Source("source3", labels),
				s.Step("step1", "", "hello", "world"))
			s.testReconcile(r, "step1")

			s.assertPodsForStep(r.client, "step1", 3)
			s.assertPodNodeAffinity(r.client, "step1", sourceNode)
		})
	}

	doTest("node1")
	doTest("node2")
}

func (s *SchedulerTestSuite) assertPodNodeAffinity(cl client.Client, stepName string, nodeName string) {
	var list corev1.PodList
	s.NoError(cl.List(context.TODO(), &client.ListOptions{}, &list))

	found := false
	for _, pod := range list.Items {
		if pod.Labels[v1.LabelStepName] != stepName {
			continue
		}
		found = true

		s.NotNil(pod.Spec.Affinity)
		s.Equal(nodeName,
			pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values[0])
	}
	s.True(found, "No pod found for step %v", stepName)
}

func (s *SchedulerTestSuite) TestSchedulingOutputSource() {
	stepName := "step-1"

	// create random Pod
	pod := s.PodLabels("pod1", map[string]string{v1.LabelStepName: stepName})
	pod.Spec.NodeName = "node2"
	// create fake output source on that pod
	source := s.Source(ConstructSourceName("pod1", "randomSource"), map[string]string{
		v1.LabelStepName:      stepName,
		v1.SourceLabelPodName: pod.Name,
		"hello":               "world",
	})

	r := s.initReconciler(
		pod, source,
		s.Node("node1"), s.Node("node2"), s.Node("node3"),
		s.Step(stepName, "", "hello", "world"))

	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, 1)
	s.assertPodNodeAffinity(r.client, stepName, "node2")
}
