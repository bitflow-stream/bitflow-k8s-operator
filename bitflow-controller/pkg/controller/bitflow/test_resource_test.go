package bitflow

import (
	"context"
	"math"
	"strings"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const resourceShare = common.TestNodeResourceLimit / float64(common.TestNodeBufferInitSize)

type ResourcesTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestResources() {
	s.SubTestSuite(new(ResourcesTestSuite))
}

func (s *ResourcesTestSuite) TestResourceAssignment1() {
	stepName := "bitflow-step-1"
	node := "node1"
	sourceLabels := map[string]string{"nodename": node, "x": "y"}
	r := s.initReconciler(
		s.Node(node), s.DefaultSchedulersStep(stepName, "", "x", "y"), s.Source("source1", sourceLabels))

	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, 1)
	s.assertPodResourceLimit(r.client, stepName, resourceShare)
}

func (s *ResourcesTestSuite) assertPodResourceLimit(cl client.Client, stepName string, share float64) {
	expectedCpu := int64(math.Round(float64(common.TestNodeCpu) * share))
	expectedMem := int64(math.Round(float64(common.TestNodeMem) * share))

	var list corev1.PodList
	s.NoError(cl.List(context.TODO(), &client.ListOptions{}, &list))
	for _, pod := range list.Items {
		if strings.Index(pod.Name, stepName) != 0 || common.IsBeingDeleted(&pod) {
			continue
		}
		for _, container := range pod.Spec.Containers {
			s.Equal(expectedCpu, container.Resources.Limits.Cpu().MilliValue(), "Wrong CPU limit")
			s.Equal(expectedMem, container.Resources.Limits.Memory().Value(), "Wrong mem limit")
		}
	}
}

func (s *ResourcesTestSuite) assignNodeNameToPods(cl client.Client, node string) {
	var list corev1.PodList
	s.NoError(cl.List(context.TODO(), &client.ListOptions{}, &list))
	for _, pod := range list.Items {
		pod.Spec.NodeName = node
		s.NoError(cl.Update(context.TODO(), &pod))
	}
}

/*
	Wenn wir mehrere Pods neu starten, muessen wir die Pods, die noch in der Pipe liegen beachten
	Sonst wird das Limit wieder falsch gesetzt un das ganze pendelt immer hin und her
	Unbedingt die Respawning Map benutzen und pflegen
*/
func (s *ResourcesTestSuite) TestResourceAssignment2() {
	stepName := "bitflow-step-1"
	node := "node1"
	sourceLabels := map[string]string{"nodename": node, "x": "y"}
	objects := s.addSources("source", common.TestNodeBufferInitSize, sourceLabels,
		s.Node(node),
		s.StepWithOutput(stepName, "", "out", map[string]string{"hello": "world"}, "x", "y"))

	r := s.initReconciler(objects...)
	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, common.TestNodeBufferInitSize)

	s.assignIPToPods(r.client)
	s.testReconcile(r, stepName)
	s.assertOutputSources(r.client, common.TestNodeBufferInitSize)
	s.assertPodResourceLimit(r.client, stepName, resourceShare)

	s.assignNodeNameToPods(r.client, node)

	addPods := int(float64(common.TestNodeBufferInitSize)*common.TestNodeBufferFactor - float64(common.TestNodeBufferInitSize))
	for _, source := range s.addSources("dynamicSource", addPods, sourceLabels) {
		s.NoError(r.client.Create(context.TODO(), source))
	}

	share := resourceShare / common.TestNodeBufferFactor
	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, addPods)
	s.assertPodResourceLimit(r.client, stepName, resourceShare/common.TestNodeBufferFactor)

	s.assertOutputSources(r.client, common.TestNodeBufferInitSize)
	s.assignNodeNameToPods(r.client, node)
	s.assignIPToPods(r.client)

	s.NoError(r.client.Create(context.TODO(), s.Source("dynamicSourceLast", sourceLabels)))

	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, 1)
	s.assertOutputSources(r.client, common.TestNodeBufferInitSize+addPods)

	share = share / common.TestNodeBufferFactor
	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, addPods+common.TestNodeBufferInitSize+1)
	s.assertPodResourceLimit(r.client, stepName, share)
}

func (s *ResourcesTestSuite) TestResourceAssignment3() {
	stepName := "bitflow-step-1"
	node := "node1"
	sourceLabels := map[string]string{"nodename": node, "x": "y"}

	r := s.initReconciler(
		s.addSources("source", common.TestNodeBufferInitSize, sourceLabels,
			s.Node(node),
			s.StepWithOutput(stepName, "", "out", map[string]string{"hello": "world"}, "x", "y"))...)

	s.testReconcile(r, stepName)
	s.assertPodResourceLimit(r.client, stepName, resourceShare)

	s.assignNodeNameToPods(r.client, node)

	source := s.Source("dynamicSource", sourceLabels)
	s.NoError(r.client.Create(context.TODO(), source))

	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, 1)

	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, common.TestNodeBufferInitSize+1)
	s.assertPodResourceLimit(r.client, stepName, resourceShare/common.TestNodeBufferFactor)

	s.assignNodeNameToPods(r.client, node)
	s.NoError(r.client.Delete(context.TODO(), source))
	s.deletePodForSource(r.client, source.Name)

	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, 0)

	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, common.TestNodeBufferInitSize)
	s.assertPodResourceLimit(r.client, stepName, resourceShare)
}

func (s *ResourcesTestSuite) TestResourceAssignment4() {
	stepName := "bitflow-step-1"
	node := "node1"
	sourceLabels := map[string]string{"nodename": node, "x": "y"}

	r := s.initReconciler(
		s.addSources("source", common.TestNodeBufferInitSize, sourceLabels,
			s.Node(node),
			s.StepWithOutput(stepName, "", "out", map[string]string{"hello": "world"}, "x", "y"))...)

	s.testReconcile(r, stepName)
	s.assertPodResourceLimit(r.client, stepName, resourceShare)

	s.assignNodeNameToPods(r.client, node)
	source := s.Source("dynamicSource", sourceLabels)
	s.NoError(r.client.Create(context.TODO(), source))

	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, 1)
	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, common.TestNodeBufferInitSize+1)
	s.assertPodResourceLimit(r.client, stepName, resourceShare/common.TestNodeBufferFactor)

	s.assignNodeNameToPods(r.client, node)

	// Mark pods as deleting
	var list corev1.PodList
	s.NoError(r.client.List(context.TODO(), &client.ListOptions{}, &list))
	timestamp := metav1.Now()
	for _, pod := range list.Items {
		if pod.Labels[bitflowv1.PodLabelOneToOneSourceName] == "dynamicSource" {
			pod.DeletionTimestamp = &timestamp
			s.NoError(r.client.Update(context.TODO(), &pod))
		}
	}

	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, 0)
	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, common.TestNodeBufferInitSize)
	s.assertPodResourceLimit(r.client, stepName, resourceShare)
}

func (s *ResourcesTestSuite) TestResourceAssignment5() {
	stepName := "bitflow-step-1"
	node := "node1"
	sourceLabels := map[string]string{"nodename": node, "x": "y"}

	r := s.initReconciler(
		s.addSources("source", common.TestNodeBufferInitSize, sourceLabels,
			s.Node(node),
			s.StepWithOutput(stepName, "", "out", map[string]string{"hello": "world"}, "x", "y"))...)

	s.testReconcile(r, stepName)
	source := s.Source("dynamicSource", sourceLabels)

	s.NoError(r.client.Create(context.TODO(), source))
	s.testReconcile(r, stepName)
	s.assertRespawningPods(r, common.TestNodeBufferInitSize)

	s.NoError(r.client.Delete(context.TODO(), s.DefaultSchedulersStep(stepName, "")))
	s.testReconcile(r, stepName)
	s.assertRespawningPods(r, 0)
	s.assertPodsForStep(r.client, stepName, 0)
}
