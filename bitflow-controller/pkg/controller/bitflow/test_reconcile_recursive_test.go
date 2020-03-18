package bitflow

import (
	"context"
	"strconv"
	"strings"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ReconcileRecursiveTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestReconcileRecursive() {
	s.SubTestSuite(new(ReconcileRecursiveTestSuite))
}

func (s *ReconcileRecursiveTestSuite) TestOutputSourceCreated() {
	name := "bitflow-step-1"
	r := s.initReconciler(
		s.Node("node1"),
		s.Source("source1", map[string]string{"hello": "world"}),
		s.StepWithOutput(name, "", "out", map[string]string{"x": "y"}, "hello", "world"))

	s.testReconcile(r, name)
	s.assignIPToPods(r.client)
	s.testReconcile(r, name)

	s.assertPodsForStep(r.client, name, 1)
	s.assertOutputSources(r.client, 1)
	s.assertOutputDepth(r.client, name, 1)
}

func (s *ReconcileRecursiveTestSuite) assertOutputDepth(cl client.Client, step string, depth int) {
	var sourceList bitflowv1.BitflowSourceList
	s.NoError(cl.List(context.TODO(), &client.ListOptions{}, &sourceList))

	for _, source := range sourceList.Items {
		if strings.Contains(source.Name, STEP_OUTPUT_PREFIX) && strings.Contains(source.Name, step+"out") {
			s.Equal(depth, source.Labels[bitflowv1.PipelineDepthLabel], "Wrong pipeline depth")
			s.Equal(step, source.Labels[bitflowv1.PipelinePathLabelPrefix+strconv.Itoa(depth)], "Wrong step name in output source")
		}
	}
}

func (s *ReconcileRecursiveTestSuite) TestRecursivePipeline() {
	name := "bitflow-step-1"
	r := s.initReconciler(
		s.Node("node1"),
		s.Source("source1", map[string]string{"hello": "world"}),
		s.StepWithOutput(name, "", "out", map[string]string{"x": "y"}, "hello", "world"))

	s.testReconcile(r, name)
	s.assignIPToPods(r.client)
	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)
	s.assertOutputSources(r.client, 1)

	name2 := "bitflow-step-2"
	step := s.StepWithOutput(name2, "", "out", map[string]string{"case": "two"}, "x", "y")
	s.NoError(r.client.Create(context.TODO(), step))

	s.testReconcile(r, name2)
	s.assignIPToPods(r.client)
	s.testReconcile(r, name2)

	s.assertPodsForStep(r.client, name, 1)
	s.assertPodsForStep(r.client, name2, 1)
	s.assertOutputSources(r.client, 2)
	s.assertOutputDepth(r.client, name, 1)
	s.assertOutputDepth(r.client, name2, 2)
}

func (s *ReconcileRecursiveTestSuite) TestRecursivePipelineLoop() {
	name := "bitflow-step-1"
	r := s.initReconciler(
		s.Node("node1"),
		s.Source("source1", map[string]string{"hello": "world"}),
		s.StepWithOutput(name, "", "out", map[string]string{"x": "y"}, "hello", "world"))

	s.testReconcile(r, name)
	s.assignIPToPods(r.client)
	s.testReconcile(r, name)
	s.assertPodsForStep(r.client, name, 1)
	s.assertOutputSources(r.client, 1)

	name2 := "bitflow-step-2"
	step := s.StepWithOutput(name2, "", "out", map[string]string{"hello": "world"}, "x", "y")
	s.NoError(r.client.Create(context.TODO(), step))

	s.testReconcile(r, name2)
	s.assignIPToPods(r.client)
	s.testReconcile(r, name2)
	s.testReconcile(r, name)

	s.assertOutputSources(r.client, 2)
	s.assertOutputDepth(r.client, name, 1)
	s.assertOutputDepth(r.client, name2, 2)
	s.assertPodsForStep(r.client, name, 1)
}

func (s *ReconcileRecursiveTestSuite) TestRecursivePipelineInverseCreation() {
	name := "bitflow-step-1"
	r := s.initReconciler(
		s.Node("node1"),
		s.Source("source1", map[string]string{"hello": "world"}),
		s.StepWithOutput(name, "", "out", map[string]string{"case": "two"}, "x", "y"))

	s.testReconcile(r, name)
	s.assertNoPodsExist(r.client)

	name2 := "bitflow-step-2"
	step := s.StepWithOutput(name2, "", "out", map[string]string{"x": "y"}, "hello", "world")
	s.NoError(r.client.Create(context.TODO(), step))

	s.testReconcile(r, name2)
	s.assignIPToPods(r.client)
	s.testReconcile(r, name2)
	s.testReconcile(r, name)
	s.assignIPToPods(r.client)
	s.testReconcile(r, name)

	s.assertOutputSources(r.client, 2)
	s.assertOutputDepth(r.client, name, 1)
	s.assertOutputDepth(r.client, name2, 2)
	s.assertPodsForStep(r.client, name, 1)
}
