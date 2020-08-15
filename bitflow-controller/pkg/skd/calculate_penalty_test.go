package skd

import (
	"fmt"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/stretchr/testify/suite"
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

func (s *CalculatePenaltyTestSuite) assertAlmostEqual(a float64, b float64) {
	s.True(almostEqual(a, b), "%v is not almost equal to %v", a, b)
}

func (s *CalculatePenaltyTestSuite) testGetNumberOfPodSlots(initialNumberOfPodSlots int, podSlotScalingFactor int, numberOfPods int, expectedNumberOfSlots int) {
	s.SubTest(fmt.Sprintf("init%d:scale%d:pods%d->%d", initialNumberOfPodSlots, podSlotScalingFactor, numberOfPods, expectedNumberOfSlots), func() {
		nodeData := &NodeData{
			initialNumberOfPodSlots: initialNumberOfPodSlots,
			podSlotScalingFactor:    podSlotScalingFactor,
		}

		actualNumberOfSlots, err := GetNumberOfPodSlots(nodeData, numberOfPods)

		s.Nil(err)
		s.Equal(expectedNumberOfSlots, actualNumberOfSlots)
	})
}

func (s *CalculatePenaltyTestSuite) Test_shouldGetNumberOfPodSlotsWithDifferentCombinations() {
	s.testGetNumberOfPodSlots(2, 2, 0, 2)
	s.testGetNumberOfPodSlots(2, 2, 1, 2)
	s.testGetNumberOfPodSlots(2, 2, 2, 2)
	s.testGetNumberOfPodSlots(2, 2, 3, 4)
	s.testGetNumberOfPodSlots(2, 2, 4, 4)
	s.testGetNumberOfPodSlots(2, 2, 5, 8)
	s.testGetNumberOfPodSlots(2, 2, 8, 8)
	s.testGetNumberOfPodSlots(2, 2, 9, 16)
	s.testGetNumberOfPodSlots(2, 2, 64, 64)
	s.testGetNumberOfPodSlots(2, 2, 65, 128)

	s.testGetNumberOfPodSlots(2, 3, 3, 6)
	s.testGetNumberOfPodSlots(2, 3, 6, 6)
	s.testGetNumberOfPodSlots(2, 3, 7, 18)

	s.testGetNumberOfPodSlots(3, 4, 0, 3)
	s.testGetNumberOfPodSlots(3, 4, 3, 3)
	s.testGetNumberOfPodSlots(3, 4, 4, 12)
	s.testGetNumberOfPodSlots(3, 4, 12, 12)
	s.testGetNumberOfPodSlots(3, 4, 13, 48)
}

func (s *CalculatePenaltyTestSuite) testCalculateExecutionTime(cpuMillis float64, curveA float64, curveB float64, curveC float64, curveD float64, expectedExecutionTime float64) {
	s.SubTest(fmt.Sprintf("cpuMillis%f:a%f:b%f:c%f:d%f->%f", cpuMillis, curveA, curveB, curveC, curveD, expectedExecutionTime), func() {
		actualExecutionTime := CalculateExecutionTime(cpuMillis, Curve{
			a: curveA,
			b: curveB,
			c: curveC,
			d: curveD,
		})

		s.assertAlmostEqual(actualExecutionTime, expectedExecutionTime)
	})
}

func (s *CalculatePenaltyTestSuite) Test_shouldCalculateExecutionTimeForDifferentCombinationsOfCpusAndCurves() {
	s.testCalculateExecutionTime(
		16000,
		6.71881241016441,
		0.0486498280492762,
		2.0417306475862214,
		15.899403720950454,
		15.9226)
	s.testCalculateExecutionTime(
		4000,
		6.71881241016441,
		0.0486498280492762,
		2.0417306475862214,
		15.899403720950454,
		16.2861)
	s.testCalculateExecutionTime(
		2000,
		6.71881241016441,
		0.0486498280492762,
		2.0417306475862214,
		15.899403720950454,
		17.4530)
	s.testCalculateExecutionTime(
		1000,
		6.71881241016441,
		0.0486498280492762,
		2.0417306475862214,
		15.899403720950454,
		21.9972)
	s.testCalculateExecutionTime(
		500,
		6.71881241016441,
		0.0486498280492762,
		2.0417306475862214,
		15.899403720950454,
		38.7860)
	s.testCalculateExecutionTime(
		250,
		6.71881241016441,
		0.0486498280492762,
		2.0417306475862214,
		15.899403720950454,
		95.1258)
	s.testCalculateExecutionTime(
		50,
		6.71881241016441,
		0.0486498280492762,
		2.0417306475862214,
		15.899403720950454,
		776.3603)
}

func (s *CalculatePenaltyTestSuite) testCalculatePenalty(testName string, state SystemState, expectedPenalty float64) {
	s.SubTest(testName, func() {
		actualPenalty, err := CalculatePenalty(state)

		s.Nil(err)
		s.assertAlmostEqual(actualPenalty, expectedPenalty)
	})
}

func (s *CalculatePenaltyTestSuite) Test_shouldCalculatePenaltyForDifferentStates() {
	s.testCalculatePenalty(
		"",
		SystemState{
			[]*NodeState{
				{
					node: &NodeData{
						name:                    "n1",
						allocatableCpu:          16000,
						memory:                  2000,
						initialNumberOfPodSlots: 8,
						podSlotScalingFactor:    2,
						resourceLimit:           0.5, // TODO bezieht sich das resourceLimit nur auf CPU oder auch auf memory?
					},
					pods: []*PodData{
						{
							name:             "p1",
							receivesDataFrom: []string{},
							curve: Curve{
								a: 6.71881241016441,
								b: 0.0486498280492762,
								c: 2.0417306475862214,
								d: 15.899403720950454,
							},
							minimumMemory: 64,
						}},
				}}},
		21.9972)
}
