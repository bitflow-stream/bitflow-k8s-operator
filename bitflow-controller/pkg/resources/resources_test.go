package resources

import (
	"testing"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type ResourcesTestSuite struct {
	common.AbstractTestSuite
}

func TestResources(t *testing.T) {
	new(ResourcesTestSuite).Run(t)
}

func (s *ResourcesTestSuite) resourceList(cpu, mem int64) *corev1.ResourceList {
	return &corev1.ResourceList{
		corev1.ResourceCPU:    *resource.NewMilliQuantity(cpu, resource.DecimalSI),
		corev1.ResourceMemory: *resource.NewQuantity(mem, resource.BinarySI),
	}
}

func (s *ResourcesTestSuite) assignResources(existingContainers int, totalLimit float64, cpu, memory int64, initSize, respawning int, factor float64) *corev1.ResourceList {
	nodeInfo := NodeInfo{
		NumberOfBitflowContainers: existingContainers,
		TotalResourceLimit:        totalLimit,
		AllocatableResources:      *s.resourceList(cpu, memory),
	}
	return nodeInfo.GetCurrentResourceList(initSize, respawning, factor)
}

func (s *ResourcesTestSuite) TestResourceAssignment1() {
	res := s.assignResources(2, 0.1, 100, 512*MBytes, 2, 0, 2.0)
	s.Equal(int64(5), res.Cpu().MilliValue())
}

func (s *ResourcesTestSuite) TestResourceAssignment2() {
	res := s.assignResources(2, 0.1, 100, 512*MBytes, 2, 1, 2.0)
	s.Equal(int64(3), res.Cpu().MilliValue())
}

func (s *ResourcesTestSuite) TestResourceAssignment3() {
	res := s.assignResources(2, 0.1, 100, 512*MBytes, 1, 0, 4.0)
	s.Equal(int64(3), res.Cpu().MilliValue())
}

func (s *ResourcesTestSuite) TestResourceAssignment4() {
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
	PatchPodResourceLimitList(pod, s.resourceList(cpu, mem))

	res := pod.Spec.Containers[0].Resources
	s.Zero(res.Requests.Cpu().MilliValue())
	s.Zero(res.Requests.Memory().Value())
	s.Equal(mem, res.Limits.Memory().Value())
	s.Equal(cpu, res.Limits.Cpu().MilliValue())
}
