package bitflow

import (
	"context"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ReconcileSingletonTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestReconcileSingleton() {
	s.SubTestSuite(new(ReconcileSingletonTestSuite))
}

func (s *ReconcileSingletonTestSuite) TestSingletonStep() {
	name := "bitflow-step-1"
	labels := map[string]string{}
	r := s.initReconciler(s.Node("node1"),
		s.StepWithOutput(name, bitflowv1.StepTypeSingleton, "out", map[string]string{"x": "y"}),
		s.Source("source1", labels), s.Source("source2", labels))

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)

	s.assignIPToPods(r.client)
	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)
	s.assertOutputSources(r.client, 1)
}

func (s *ReconcileSingletonTestSuite) TestSingletonStepWithWrongStepName() {
	name := "bitflow-step-1"
	r := s.initReconciler(s.Node("node1"),
		s.StepWithOutput(name, bitflowv1.StepTypeSingleton, "out", map[string]string{"x": "y"}))

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)

	s.assignIPToPods(r.client)
	s.changeStepName(r.client)
	s.testReconcile(r, name)
	s.assertNoPodsExist(r.client)

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)
}

func (s *ReconcileSingletonTestSuite) changeStepName(cl client.Client) {
	var list corev1.PodList
	s.NoError(cl.List(context.TODO(), &client.ListOptions{}, &list))
	for _, pod := range list.Items {
		pod.Labels[bitflowv1.LabelStepName] = "xx-random-xx"
		s.NoError(cl.Update(context.TODO(), &pod))
	}
}
