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

func (s *ReconcileOneToOneTestSuite) TestOneToOneDeleteSourceCheckRespawning1() {
	name := "bitflow-step-1"
	r := s.makeReconcilerWithSources(name, common.TestNodeBufferInitSize+1)

	// all pods started with halved resources
	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, common.TestNodeBufferInitSize+1)
	s.assertRespawningPods(r, 0)

	s.NoError(r.client.Delete(context.TODO(), s.Source("source0", nil)))

	// pods restarted with more resources. First reconcile: delete pods. Second reconcile: restart pods.
	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 0)
	s.assertRespawningPods(r, common.TestNodeBufferInitSize)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, common.TestNodeBufferInitSize)
	s.assertRespawningPods(r, 0)
}

func (s *ReconcileOneToOneTestSuite) TestOneToOneDeleteSourceNoRestart() {
	name := "bitflow-step-1"
	num := common.TestNodeBufferInitSize + 2
	r := s.makeReconcilerWithSources(name, num)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, common.TestNodeBufferInitSize+2)
	s.assertRespawningPods(r, 0)

	// Delete one of the running pods and its source
	sourceName := "source" + strconv.Itoa(num-1)
	s.NoError(r.client.Delete(context.TODO(), s.Source(sourceName, nil)))
	s.deletePodForSource(r.client, sourceName)

	// Pods should not be respawned - resources remain the same
	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, common.TestNodeBufferInitSize+1)
	s.assertRespawningPods(r, 0)
}

func (s *ReconcileOneToOneTestSuite) TestOneToOneCreateSourceCheckRespawning() {
	name := "bitflow-step-1"
	r := s.makeReconcilerWithSources(name, common.TestNodeBufferInitSize)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, common.TestNodeBufferInitSize)
	s.assertRespawningPods(r, 0)

	// Add another source - old pods should be restarted, new pod added
	sourceName := "extra-source"
	s.NoError(r.client.Create(context.TODO(), s.Source(sourceName, map[string]string{"x": "y"})))

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)
	s.assertRespawningPods(r, common.TestNodeBufferInitSize)

	// Second reconcile should create all respawned pods
	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, common.TestNodeBufferInitSize+1)
	s.assertRespawningPods(r, 0)
}
