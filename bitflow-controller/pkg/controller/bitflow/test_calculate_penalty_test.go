package bitflow

import "math"

type CalculatePenaltyTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestCalculatePenalty() {
	s.SubTestSuite(new(CalculatePenaltyTestSuite))
}

const float64EqualityThreshold = 1e-5

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func (s *CalculatePenaltyTestSuite) TestCalculatePenaltyForNode() {
	numberOfPods := 4

	labels := map[string]string{"hello": "world"}
	r := s.initReconciler(
		s.Node("node1"),
		s.Source("source1", labels), s.Source("source2", labels),
		s.Source("source3", labels), s.Source("source4", labels),
		s.Step("step1", "", "hello", "world"))
	s.testReconcile(r, "step1")

	penalty, err := calculatePenaltyForNode(r.client, "node1")

	s.NoError(err)
	s.Equal(numberOfPods*getPodPenalty(), penalty)
}

func (s *CalculatePenaltyTestSuite) TestCalculatePenaltyForNonExistingNode() {
	r := s.initReconciler()
	penalty, err := calculatePenaltyForNode(r.client, "node1")

	s.NoError(err)
	s.Equal(0, penalty)
}

func (s *CalculatePenaltyTestSuite) TestCalculateExecutionTimeOneCpu() {
	executionTime := calculateExecutionTime(1)
	s.Equal(2.0, executionTime)
}

func (s *CalculatePenaltyTestSuite) TestCalculateExecutionTimeTwoCpus() {
	executionTime := calculateExecutionTime(2)
	s.True(almostEqual(1.641133, executionTime))
}

func (s *CalculatePenaltyTestSuite) TestCalculateExecutionTimeHalfCpu() {
	executionTime := calculateExecutionTime(0.5)
	s.True(almostEqual(2.33484, executionTime))
}

// TODO test for very small cpu count + very high cpu count once actual curve parameters are defined
