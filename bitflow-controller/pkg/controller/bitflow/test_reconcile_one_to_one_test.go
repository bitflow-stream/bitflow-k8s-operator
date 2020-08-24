package bitflow

import (
	"context"
	"strconv"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ReconcileOneToOneTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestReconcileOneToOne() {
	s.SubTestSuite(new(ReconcileOneToOneTestSuite))
}

func (s *ReconcileOneToOneTestSuite) TestOnlyStep() {
	name := "bitflow-step-1"
	r := s.initReconciler(s.Node("node1"),
		s.StepWithOutput(name, "", "out", map[string]string{"a": "b"}, "x", "y"))

	s.testReconcile(r, name)
	s.assertNoPodsExist(r.client)
	s.assertNoSourceExists(r.client)
}

func (s *ReconcileOneToOneTestSuite) TestOneSourceOneStep() {
	name := "bitflow-step-1"
	r := s.initReconciler(s.Node("node1"),
		s.StepWithOutput(name, "", "out", map[string]string{"a": "b"}, "x", "y"),
		s.Source("source1", map[string]string{"x": "y"}))

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)

	s.assignIPToPods(r.client)
	s.testReconcile(r, name)
	s.assertOutputSources(r.client, 1)
}

func (s *ReconcileOneToOneTestSuite) TestOneSourceOneStepDeletePod() {
	name := "bitflow-step-1"
	r := s.initReconciler(s.Node("node1"),
		s.StepWithOutput(name, "", "out", map[string]string{"a": "b"}, "x", "y"),
		s.Source("source1", map[string]string{"x": "y"}))

	s.testReconcile(r, name)
	s.assignIPToPods(r.client)
	s.testReconcile(r, name)

	// Delete all pods
	var list corev1.PodList
	err := r.client.List(context.TODO(), &client.ListOptions{}, &list)
	s.NoError(err)
	for _, pod := range list.Items {
		s.NoError(r.client.Delete(context.TODO(), &pod))
	}

	// Re-do the reconcile
	s.testReconcile(r, name)
	s.assignIPToPods(r.client)
	s.testReconcile(r, name)

	s.assertPodsForStep(r.client, name, 1)
	s.assertOutputSources(r.client, 1)
}

func (s *ReconcileOneToOneTestSuite) makeReconcilerWithSources(stepName string, numSources int) *BitflowReconciler {
	return s.initReconciler(
		s.addSources("source", numSources, map[string]string{"x": "y"},
			s.Node("node1"),
			s.StepWithOutput(stepName, "", "out", map[string]string{"a": "b"}, "x", "y"))...)
}

func (s *ReconcileOneToOneTestSuite) TestOneToOneStepKeepOutput() {
	name := "bitflow-step-1"
	r := s.makeReconcilerWithSources(name, common.TestNodeBufferInitSize)

	s.testReconcile(r, name)
	s.assignIPToPods(r.client)
	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, common.TestNodeBufferInitSize)
	s.assertOutputSources(r.client, common.TestNodeBufferInitSize)

	// Delete the IP from all pods
	var list corev1.PodList
	s.NoError(r.client.List(context.TODO(), &client.ListOptions{}, &list))
	for _, pod := range list.Items {
		pod.Status.PodIP = ""
		s.NoError(r.client.Update(context.TODO(), &pod))
	}

	s.testReconcile(r, name)
	s.assertOutputSources(r.client, common.TestNodeBufferInitSize)
}

func (s *ReconcileOneToOneTestSuite) TestRespawningOneToOne() {
	name := "bitflow-step-1"
	r := s.makeReconcilerWithSources(name, common.TestNodeBufferInitSize+1)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)
	s.assertRespawningPods(r, common.TestNodeBufferInitSize)

	// TODO this is actually weird behavior in the resource assigner, which should be fixed
	// currently, newly created pods are not taken into account when computing resources

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 2)
	s.assertRespawningPods(r, common.TestNodeBufferInitSize+1-2)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 3)
	s.assertRespawningPods(r, common.TestNodeBufferInitSize+1-3)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 4)
	s.assertRespawningPods(r, common.TestNodeBufferInitSize+1-4)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, common.TestNodeBufferInitSize+1)
	s.assertRespawningPods(r, 0)
}

func (s *ReconcileOneToOneTestSuite) TestOneToOneDeleteSourceCheckRespawning1() {
	name := "bitflow-step-1"
	r := s.makeReconcilerWithSources(name, common.TestNodeBufferInitSize+1)

	// 1 pod started with correct resources, 4 pods pods due to wrong resources
	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)
	s.assertRespawningPods(r, common.TestNodeBufferInitSize)

	s.NoError(r.client.Delete(context.TODO(), s.Source("source0", nil)))

	// the pod that was originally started with correct resources is now too small, therefore pods.
	// the other 3 pods were started correctly.
	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, common.TestNodeBufferInitSize-1)
	s.assertRespawningPods(r, 1)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, common.TestNodeBufferInitSize)
	s.assertRespawningPods(r, 0)
}

func (s *ReconcileOneToOneTestSuite) TestOneToOneDeleteSourceCheckRespawning2() {
	name := "bitflow-step-1"
	r := s.makeReconcilerWithSources(name, common.TestNodeBufferInitSize+1)

	// 1 pod started with correct resources, 4 pods pods due to wrong resources
	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)
	s.assertRespawningPods(r, common.TestNodeBufferInitSize)

	// Delete the source for the pod that was started correctly (the last one)
	sourceName := "source" + strconv.Itoa(common.TestNodeBufferInitSize)
	s.NoError(r.client.Delete(context.TODO(), s.Source(sourceName, nil)))
	s.deletePodForSource(r.client, sourceName)

	// Now all pods should immediately be started correctly
	s.testReconcile(r, name)
	s.assertRespawningPods(r, 0)
	s.assertPodsForStep(r.client, name, common.TestNodeBufferInitSize)
}

func (s *ReconcileOneToOneTestSuite) TestOneToOneDeleteSourceCheckRespawning3() {
	name := "bitflow-step-1"
	num := common.TestNodeBufferInitSize + 2
	r := s.makeReconcilerWithSources(name, num)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 2)
	s.assertRespawningPods(r, common.TestNodeBufferInitSize)

	// Delete one of the correctly running pods and its source
	sourceName := "source" + strconv.Itoa(num-1)
	s.NoError(r.client.Delete(context.TODO(), s.Source(sourceName, nil)))
	s.deletePodForSource(r.client, sourceName)

	// TODO see TestRespawningOneToOne, this behavior should be fixed

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 2)
	s.assertRespawningPods(r, num-3)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 3)
	s.assertRespawningPods(r, num-4)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 4)
	s.assertRespawningPods(r, num-5)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 5)
	s.assertRespawningPods(r, 0)
}
