package scheduler

import (
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

func (s *CalculatePenaltyTestSuite) TestCalculateExecutionTimeOneCpu() {
	executionTime := CalculateExecutionTime(1)
	s.assertAlmostEqual(21.99716231, executionTime)
}

func (s *CalculatePenaltyTestSuite) TestCalculateExecutionTimeTwoCpus() {
	executionTime := CalculateExecutionTime(2)
	s.assertAlmostEqual(17.45302995, executionTime)
}

func (s *CalculatePenaltyTestSuite) TestCalculateExecutionTimeHalfCpu() {
	executionTime := CalculateExecutionTime(0.5)
	s.assertAlmostEqual(38.78601002, executionTime)
}

func (s *CalculatePenaltyTestSuite) TestCalculateExecutionTimeSmallR() {
	executionTime := CalculateExecutionTime(0.025)
	s.assertAlmostEqual(1396.99168564, executionTime)
}

// TODO test for very small cpu count + very high cpu count

func (s *CalculatePenaltyTestSuite) TestGetAllocatableCpu() {
	cpu := getAllocatableCpu(*s.Node("node"))
	s.Equal(2000.0, cpu)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorTwoValueOne() {
	next, _ := getNextHigherNumberOfPodSlots(2, 1)
	s.Equal(2.0, next)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorTwoValueOneAndAHalf() {
	next, _ := getNextHigherNumberOfPodSlots(2, 1.5)
	s.Equal(2.0, next)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorTwoValueTwo() {
	next, _ := getNextHigherNumberOfPodSlots(2, 2)
	s.Equal(2.0, next)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorTwoValueTwoAndAHalf() {
	next, _ := getNextHigherNumberOfPodSlots(2, 2.5)
	s.Equal(4.0, next)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorTwoValueFour() {
	next, _ := getNextHigherNumberOfPodSlots(2, 4)
	s.Equal(4.0, next)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorTwoValueSeven() {
	next, _ := getNextHigherNumberOfPodSlots(2, 7)
	s.Equal(8.0, next)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorTwoValueOneHundredTwentyNine() {
	next, _ := getNextHigherNumberOfPodSlots(2, 129)
	s.Equal(256.0, next)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorThreeValueOne() {
	next, _ := getNextHigherNumberOfPodSlots(3, 1)
	s.Equal(3.0, next)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorThreeValueThree() {
	next, _ := getNextHigherNumberOfPodSlots(3, 3)
	s.Equal(3.0, next)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorThreeValueThreeAndAHalf() {
	next, _ := getNextHigherNumberOfPodSlots(3, 3.5)
	s.Equal(9.0, next)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorThreeValueEightyAndAHalf() {
	next, _ := getNextHigherNumberOfPodSlots(3, 80.5)
	s.Equal(81.0, next)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorThreeValueEightyOne() {
	next, _ := getNextHigherNumberOfPodSlots(3, 81)
	s.Equal(81.0, next)
}

func (s *CalculatePenaltyTestSuite) TestGetNextHigherNumberOfPodSlotsFactorThreeValueNinetySixAndAHalf() {
	next, _ := getNextHigherNumberOfPodSlots(3, 96.5)
	s.Equal(243.0, next)
}
