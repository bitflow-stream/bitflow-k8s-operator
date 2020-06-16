package skd

import (
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	"math"
	"testing"
)

type CalculatePenaltyTestSuite struct {
	common.AbstractTestSuite
}

func TestCalculatePenalty(t *testing.T) {
	suite.Run(t, new(CalculatePenaltyTestSuite))
}

const float64EqualityThreshold = 1e-4

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func someCurve() Curve {
	return Curve{
		a: 6.71881241016441,
		b: 0.0486498280492762,
		c: 2.0417306475862214,
		d: 15.899403720950454,
	}
}

func (s *CalculatePenaltyTestSuite) assertAlmostEqual(a float64, b float64) {
	s.True(almostEqual(a, b), "%v is not almost equal to %v", a, b)
}

func (s *CalculatePenaltyTestSuite) Test_shouldCalculateExecutionTime() {
	executionTime := CalculateExecutionTime(
		0.025,
		Curve{
			a: 6.71881241016441,
			b: 0.0486498280492762,
			c: 2.0417306475862214,
			d: 15.899403720950454,
		})
	s.assertAlmostEqual(1396.99168564, executionTime)
}

func (s *CalculatePenaltyTestSuite) Test_shouldGetNumberOfPodSlots_initial2Scaling2_0Pods1ToAdd() {
	simulatedNode := SimulatedNode{
		nodeData: &NodeData{
			node:                    s.Node("Some Node"),
			curve:                   someCurve(),
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		pods: []*corev1.Pod{},
	}

	numberOfSlots, err := GetNumberOfPodSlotsAfterAddingPods(simulatedNode, 1)

	s.Nil(err)
	s.Equal(int64(2), numberOfSlots)
}

func (s *CalculatePenaltyTestSuite) Test_shouldGetNumberOfPodSlots_initial2Scaling2_1Pods1ToAdd() {
	simulatedNode := SimulatedNode{
		nodeData: &NodeData{
			node:                    s.Node("Some Node"),
			curve:                   someCurve(),
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		pods: []*corev1.Pod{
			s.Pod("pod1"),
		},
	}

	numberOfSlots, err := GetNumberOfPodSlotsAfterAddingPods(simulatedNode, 1)

	s.Nil(err)
	s.Equal(int64(2), numberOfSlots)
}

func (s *CalculatePenaltyTestSuite) Test_shouldGetNumberOfPodSlots_initial2Scaling2_2Pods1ToAdd() {
	simulatedNode := SimulatedNode{
		nodeData: &NodeData{
			node:                    s.Node("Some Node"),
			curve:                   someCurve(),
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		pods: []*corev1.Pod{
			s.Pod("pod1"),
			s.Pod("pod2"),
		},
	}

	numberOfSlots, err := GetNumberOfPodSlotsAfterAddingPods(simulatedNode, 1)

	s.Nil(err)
	s.Equal(int64(4), numberOfSlots)
}

func (s *CalculatePenaltyTestSuite) Test_shouldGetNumberOfPodSlots_initial4Scaling2_3Pods1ToAdd() {
	simulatedNode := SimulatedNode{
		nodeData: &NodeData{
			node:                    s.Node("Some Node"),
			curve:                   someCurve(),
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		pods: []*corev1.Pod{
			s.Pod("pod1"),
			s.Pod("pod2"),
			s.Pod("pod3"),
		},
	}

	numberOfSlots, err := GetNumberOfPodSlotsAfterAddingPods(simulatedNode, 1)

	s.Nil(err)
	s.Equal(int64(4), numberOfSlots)
}

func (s *CalculatePenaltyTestSuite) Test_shouldGetNumberOfPodSlots_initial4Scaling2_4Pods1ToAdd() {
	simulatedNode := SimulatedNode{
		nodeData: &NodeData{
			node:                    s.Node("Some Node"),
			curve:                   someCurve(),
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		pods: []*corev1.Pod{
			s.Pod("pod1"),
			s.Pod("pod2"),
			s.Pod("pod3"),
			s.Pod("pod4"),
		},
	}

	numberOfSlots, err := GetNumberOfPodSlotsAfterAddingPods(simulatedNode, 1)

	s.Nil(err)
	s.Equal(int64(8), numberOfSlots)
}

func (s *CalculatePenaltyTestSuite) Test_shouldGetNumberOfPodSlots_initial2Scaling3_18Pods0ToAdd() {
	simulatedNode := SimulatedNode{
		nodeData: &NodeData{
			node:                    s.Node("Some Node"),
			curve:                   someCurve(),
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    3,
			resourceLimit:           0.1,
		},
		pods: []*corev1.Pod{
			s.Pod("pod1"),
			s.Pod("pod2"),
			s.Pod("pod3"),
			s.Pod("pod4"),
			s.Pod("pod5"),
			s.Pod("pod6"),
			s.Pod("pod7"),
			s.Pod("pod8"),
			s.Pod("pod9"),
			s.Pod("pod10"),
			s.Pod("pod11"),
			s.Pod("pod12"),
			s.Pod("pod13"),
			s.Pod("pod14"),
			s.Pod("pod15"),
			s.Pod("pod16"),
			s.Pod("pod17"),
			s.Pod("pod18"),
		},
	}

	numberOfSlots, err := GetNumberOfPodSlotsAfterAddingPods(simulatedNode, 0)

	s.Nil(err)
	s.Equal(int64(18), numberOfSlots)
}

func (s *CalculatePenaltyTestSuite) Test_shouldGetNumberOfPodSlots_initial2Scaling3_18Pods1ToAdd() {
	simulatedNode := SimulatedNode{
		nodeData: &NodeData{
			node:                    s.Node("Some Node"),
			curve:                   someCurve(),
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    3,
			resourceLimit:           0.1,
		},
		pods: []*corev1.Pod{
			s.Pod("pod1"),
			s.Pod("pod2"),
			s.Pod("pod3"),
			s.Pod("pod4"),
			s.Pod("pod5"),
			s.Pod("pod6"),
			s.Pod("pod7"),
			s.Pod("pod8"),
			s.Pod("pod9"),
			s.Pod("pod10"),
			s.Pod("pod11"),
			s.Pod("pod12"),
			s.Pod("pod13"),
			s.Pod("pod14"),
			s.Pod("pod15"),
			s.Pod("pod16"),
			s.Pod("pod17"),
			s.Pod("pod18"),
		},
	}

	numberOfSlots, err := GetNumberOfPodSlotsAfterAddingPods(simulatedNode, 1)

	s.Nil(err)
	s.Equal(int64(54), numberOfSlots)
}

func (s *CalculatePenaltyTestSuite) Test_shouldCalculatePenaltyForNodeAfterAddingPods() {
	simulatedNode := SimulatedNode{
		nodeData: &NodeData{
			node: s.NodeWithCpuAndMemory("Some Node", 2000, 4, 0.1),
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		pods: []*corev1.Pod{
			s.Pod("pod1"),
			s.Pod("pod2"),
			s.Pod("pod3"),
			s.Pod("pod4"),
		},
	}

	actualPenalty := CalculatePenaltyForNodeAfterAddingPods(simulatedNode, 1)

	s.Equal(15.908765355793092, actualPenalty)
}

func (s *CalculatePenaltyTestSuite) Test_shouldGetLowestPenaltyNode_differenceInCpu() {
	simulatedNodes := []*SimulatedNode{
		{
			nodeData: &NodeData{
				node: s.NodeWithCpuAndMemory("notLowestPenaltyNode", 2000, 4, 0.1),
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
			pods: []*corev1.Pod{
				s.Pod("pod1"),
				s.Pod("pod2"),
				s.Pod("pod3"),
				s.Pod("pod4"),
			},
		},
		{
			nodeData: &NodeData{
				node: s.NodeWithCpuAndMemory("lowestPenaltyNode", 3000, 4, 0.1),
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
			pods: []*corev1.Pod{
				s.Pod("pod1"),
				s.Pod("pod2"),
				s.Pod("pod3"),
				s.Pod("pod4"),
			},
		},
	}

	lowestPenaltyNode, err := GetLowestPenaltyNode(simulatedNodes)

	s.Nil(err)
	s.Equal("lowestPenaltyNode", lowestPenaltyNode.nodeData.node.Name)
}

func (s *CalculatePenaltyTestSuite) Test_shouldGetLowestPenaltyNode_differenceInPods() {
	simulatedNodes := []*SimulatedNode{
		{
			nodeData: &NodeData{
				node: s.NodeWithCpuAndMemory("notLowestPenaltyNode", 2000, 4, 0.1),
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
			pods: []*corev1.Pod{
				s.Pod("pod1"),
				s.Pod("pod2"),
				s.Pod("pod3"),
				s.Pod("pod4"),
			},
		},
		{
			nodeData: &NodeData{
				node: s.NodeWithCpuAndMemory("lowestPenaltyNode", 2000, 4, 0.1),
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
			pods: []*corev1.Pod{
				s.Pod("pod1"),
				s.Pod("pod2"),
				s.Pod("pod3"),
			},
		},
	}

	lowestPenaltyNode, err := GetLowestPenaltyNode(simulatedNodes)

	s.Nil(err)
	s.Equal("lowestPenaltyNode", lowestPenaltyNode.nodeData.node.Name)
}
