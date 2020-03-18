package bitflow

import (
	"context"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"k8s.io/apimachinery/pkg/types"
)

type ValidationTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestValidation() {
	s.SubTestSuite(new(ValidationTestSuite))
}

func (s *ValidationTestSuite) TestValidation() {
	stepName := "bitflow-step-1"
	step := s.Step(stepName, bitflowv1.StepTypeOneToOne, "case", "one")

	r := s.initReconciler(step, s.Node("node1"))
	s.testReconcile(r, stepName)

	var found bitflowv1.BitflowStep
	s.NoError(r.client.Get(context.TODO(), types.NamespacedName{Name: stepName, Namespace: common.TestNamespace}, &found))
	s.Empty(found.Status.ValidationError)
	for _, out := range found.Spec.Outputs {
		s.NotNil(out.GetUrl())
	}
}

func (s *ValidationTestSuite) TestValidationCorruptOutput() {
	stepName := "bitflow-step-1"
	step := s.Step(stepName, "", "case", "one")
	s.AddStepOutput(step, "out", make(map[string]string))

	// corrupt URL
	step.Spec.Outputs[0].URL = "++fail://"

	r := s.initReconciler(step, s.Node("node1"))
	s.testReconcile(r, stepName)

	var found bitflowv1.BitflowStep
	s.NoError(r.client.Get(context.TODO(), types.NamespacedName{Name: stepName, Namespace: common.TestNamespace}, &found))
	s.Contains(found.Status.ValidationError, "URL")
	s.Contains(found.Status.ValidationError, "label")
}

func (s *ValidationTestSuite) TestValidationSingletonWithIngest() {
	stepName := "bitflow-step-1"
	step := s.Step(stepName, bitflowv1.StepTypeSingleton, "case", "one")
	source := s.Source("source-1", map[string]string{"case": "one"})

	r := s.initReconciler(step, source, s.Node("node1"))
	s.testReconcile(r, stepName)

	var found bitflowv1.BitflowStep
	s.NoError(r.client.Get(context.TODO(), types.NamespacedName{Name: stepName, Namespace: common.TestNamespace}, &found))
	s.Contains(found.Status.ValidationError, "empty .spec.ingest")
	s.assertNoPodsExist(r.client)
}
