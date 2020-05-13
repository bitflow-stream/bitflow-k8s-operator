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
	s.Equal(2.0, cpu)
}
