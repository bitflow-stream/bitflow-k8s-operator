package bitflow

import (
	"context"
	"math"
	"strings"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	MBytes        = 1024 * 1024
	resourceShare = common.TestNodeResourceLimit / float64(common.TestNodeBufferInitSize)
)

type ResourcesTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestResources() {
	s.SubTestSuite(new(ResourcesTestSuite))
}

func (s *ResourcesTestSuite) resourceList(cpu, mem int64) *corev1.ResourceList {
	return &corev1.ResourceList{
		corev1.ResourceCPU:    *resource.NewMilliQuantity(cpu, resource.DecimalSI),
		corev1.ResourceMemory: *resource.NewQuantity(mem, resource.BinarySI),
	}
}

func (s *ResourcesTestSuite) assignResources(existingContainers int, totalLimit float64, cpu, memory int64, initSize, respawning int, factor float64) *corev1.ResourceList {
	return buildPodResourceList(initSize, factor, totalLimit, existingContainers+respawning, *s.resourceList(cpu, memory))
}

func (s *ResourcesTestSuite) TestSimpleResourceAssignment1() {
	res := s.assignResources(2, 0.1, 100, 512*MBytes, 2, 0, 2.0)
	s.Equal(int64(5), res.Cpu().MilliValue())
}

func (s *ResourcesTestSuite) TestSimpleResourceAssignment2() {
	res := s.assignResources(2, 0.1, 100, 512*MBytes, 2, 1, 2.0)
	s.Equal(int64(3), res.Cpu().MilliValue())
}

func (s *ResourcesTestSuite) TestSimpleResourceAssignment3() {
	res := s.assignResources(2, 0.1, 100, 512*MBytes, 1, 0, 4.0)
	s.Equal(int64(3), res.Cpu().MilliValue())
}

func (s *ResourcesTestSuite) TestSimpleResourceAssignment4() {
	res := s.assignResources(1, 0.1, 100, 512*MBytes, 1, 0, 4.0)
	s.Equal(int64(10), res.Cpu().MilliValue())
}

func (s *ResourcesTestSuite) TestResourceAssignmentFactorOne() {
	res := s.assignResources(1, 0.1, 100, 512*MBytes, 2, 0, 1.0)
	s.Equal(int64(5), res.Cpu().MilliValue())
}

func (s *ResourcesTestSuite) TestResourceAssignmentFactorHalf() {
	res := s.assignResources(2, 0.1, 100, 512*MBytes, 2, 1, 0.5)
	s.Equal(int64(3), res.Cpu().MilliValue())
}

func (s *ResourcesTestSuite) TestResourceAssignmentNoLimit1() {
	res := s.assignResources(2, 1.1, 100, 512*MBytes, 2, 1, 0.5)
	s.Nil(res)
}

func (s *ResourcesTestSuite) TestResourceAssignmentNoLimit2() {
	res := s.assignResources(2, 0, 100, 512*MBytes, 2, 1, 0.5)
	s.Nil(res)
}

func (s *ResourcesTestSuite) TestPatchResourceLimitList() {
	pod := s.Pod("pod1")
	var mem int64 = 512 * MBytes
	var cpu int64 = 100
	patchPodResourceLimits(pod, s.resourceList(cpu, mem))

	res := pod.Spec.Containers[0].Resources
	s.Zero(res.Requests.Cpu().MilliValue())
	s.Zero(res.Requests.Memory().Value())
	s.Equal(mem, res.Limits.Memory().Value())
	s.Equal(cpu, res.Limits.Cpu().MilliValue())
}

func (s *ResourcesTestSuite) TestResourceAssignment1() {
	stepName := "bitflow-step-1"
	node := "node1"
	sourceLabels := map[string]string{"bitflow-nodename": node, "x": "y"}
	r := s.initReconciler(
		s.Node(node), s.Step(stepName, "", "x", "y"), s.Source("source1", sourceLabels))

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

func (s *ResourcesTestSuite) TestResourceAssignment2() {
	stepName := "bitflow-step-1"
	node := "node1"
	sourceLabels := map[string]string{"bitflow-nodename": node, "x": "y"}
	objects := s.addSources("source", common.TestNodeBufferInitSize, sourceLabels,
		s.Node(node),
		s.StepWithOutput(stepName, "", "out", map[string]string{"hello": "world"}, "x", "y"))

	r := s.initReconciler(objects...)
	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, common.TestNodeBufferInitSize)
	s.assertPodResourceLimit(r.client, stepName, resourceShare)
	s.assertOutputSources(r.client, 0)

	s.assignIPToPods(r.client)
	s.testReconcile(r, stepName)
	s.assertOutputSources(r.client, common.TestNodeBufferInitSize)

	// Add enough pods to break the resource limit, making the old pods restart
	addPods := int(float64(common.TestNodeBufferInitSize) * (common.TestNodeBufferFactor - 1.0))
	for _, source := range s.addSources("dynamicSource", addPods, sourceLabels) {
		s.NoError(r.client.Create(context.TODO(), source))
	}

	// Only the new pods should be running, the old pods are respawning now
	share := resourceShare / common.TestNodeBufferFactor
	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, addPods)
	s.assertRespawningPods(r, common.TestNodeBufferInitSize)
	s.assertPodResourceLimit(r.client, stepName, resourceShare/common.TestNodeBufferFactor)

	// Existing output sources should not be deleted
	s.assertOutputSources(r.client, common.TestNodeBufferInitSize)
	s.assignIPToPods(r.client)

	// Add another source, break the resource limit again
	s.NoError(r.client.Create(context.TODO(), s.Source("dynamicSourceLast", sourceLabels)))

	// The old pods, and the new extra pod, should be running now. The addPods pods are being restarted, resources adjusted.
	// The output sources for the addPods pods are not yet created, because IPs have not been assigned in time.
	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, common.TestNodeBufferInitSize+1)
	s.assertRespawningPods(r, addPods)
	s.assertOutputSources(r.client, common.TestNodeBufferInitSize)

	share = share / common.TestNodeBufferFactor
	s.assignIPToPods(r.client)
	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, common.TestNodeBufferInitSize+addPods+1)
	s.assertRespawningPods(r, 0)
	s.assertOutputSources(r.client, common.TestNodeBufferInitSize+1)
	s.assertPodResourceLimit(r.client, stepName, share)

	s.assignIPToPods(r.client)
	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, common.TestNodeBufferInitSize+addPods+1)
	s.assertOutputSources(r.client, common.TestNodeBufferInitSize+addPods+1)
}

func (s *ResourcesTestSuite) TestResourceAssignment3() {
	stepName := "bitflow-step-1"
	node := "node1"
	sourceLabels := map[string]string{"bitflow-nodename": node, "x": "y"}

	r := s.initReconciler(
		s.addSources("source", common.TestNodeBufferInitSize, sourceLabels,
			s.Node(node),
			s.StepWithOutput(stepName, "", "out", map[string]string{"hello": "world"}, "x", "y"))...)

	s.testReconcile(r, stepName)
	s.assertPodResourceLimit(r.client, stepName, resourceShare)

	source := s.Source("dynamicSource", sourceLabels)
	s.NoError(r.client.Create(context.TODO(), source))

	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, 1)

	s.testReconcile(r, stepName)
	s.assertPodsForStep(r.client, stepName, common.TestNodeBufferInitSize+1)
	s.assertPodResourceLimit(r.client, stepName, resourceShare/common.TestNodeBufferFactor)

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
	sourceLabels := map[string]string{"bitflow-nodename": node, "x": "y"}

	r := s.initReconciler(
		s.addSources("source", common.TestNodeBufferInitSize, sourceLabels,
			s.Node(node),
			s.StepWithOutput(stepName, "", "out", map[string]string{"hello": "world"}, "x", "y"))...)

	s.testReconcile(r, stepName)
	source := s.Source("dynamicSource", sourceLabels)

	s.NoError(r.client.Create(context.TODO(), source))
	s.testReconcile(r, stepName)
	s.assertRespawningPods(r, common.TestNodeBufferInitSize)

	s.NoError(r.client.Delete(context.TODO(), s.Step(stepName, "")))
	s.testReconcile(r, stepName)
	s.assertRespawningPods(r, 0)
	s.assertPodsForStep(r.client, stepName, 0)
}
