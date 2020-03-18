package bitflow

import (
	"context"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
)

type ReconcileAllToOneTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestReconcileAllToOne() {
	s.SubTestSuite(new(ReconcileAllToOneTestSuite))
}

func (s *ReconcileAllToOneTestSuite) TestAllToOneStepNoSources() {
	name := "bitflow-step-1"
	r := s.initReconciler(
		s.Node("node1"),
		s.Step(name, bitflowv1.StepTypeAllToOne, "x", "y"))

	s.testReconcile(r, name)
	s.assertNoPodsExist(r.client)
}

func (s *ReconcileAllToOneTestSuite) TestAllToOneStepNoMatchingSources() {
	name := "bitflow-step-1"
	r := s.initReconciler(
		s.Node("node1"),
		s.Source("source1", map[string]string{"a": "b"}), s.Source("source2", map[string]string{"c": "d"}),
		s.Step(name, bitflowv1.StepTypeAllToOne, "x", "y"))

	s.testReconcile(r, name)
	s.assertNoPodsExist(r.client)
}

func (s *ReconcileAllToOneTestSuite) TestAllToOneStepTwoSources() {
	name := "bitflow-step-1"
	r := s.initReconciler(
		s.Node("node1"),
		s.Source("source1", map[string]string{"x": "y"}), s.Source("source2", map[string]string{"x": "y"}),
		s.Source("other-source", map[string]string{"HELLO": "WORLD"}),
		s.Step(name, bitflowv1.StepTypeAllToOne, "x", "y"))

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)
}

func (s *ReconcileAllToOneTestSuite) TestAllToOneStepTwoInvalidSources() {
	name := "bitflow-step-1"
	source1 := s.Source("source1", map[string]string{"x": "y"})
	source2 := s.Source("source2", map[string]string{"x": "y"})
	r := s.initReconciler(
		s.Node("node1"),
		source1, source2, s.Step(name, bitflowv1.StepTypeAllToOne, "x", "y"))

	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)

	// Invalidate first source -> pod should remain
	source1.Spec.URL = "++fail://"
	source1.Validate()
	s.NoError(r.client.Update(context.TODO(), source1))
	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)

	// Invalidate second source -> pod should be deleted
	source2.Spec.URL = "++fail://"
	source2.Validate()
	s.NoError(r.client.Update(context.TODO(), source2))
	s.testReconcile(r, name)
	s.assertNoPodsExist(r.client)
}
