package bitflow

type CalculatePenaltyTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestCalculatePenalty() {
	s.SubTestSuite(new(CalculatePenaltyTestSuite))
}

func (s *CalculatePenaltyTestSuite) TestCalculatePenaltyForNode() {
	numberOfPods := 4

	labels := map[string]string{"hello": "world"}
	r := s.initReconciler(
		s.Node("node1"),
		s.Source("source1", labels), s.Source("source2", labels),
		s.Source("source3", labels), s.Source("source4", labels),
		s.DefaultSchedulersStep("step1", "", "hello", "world"))
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
