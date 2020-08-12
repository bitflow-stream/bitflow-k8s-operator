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

// TODO test CalculateExecutionTime
// TODO test CalculatePenalty
