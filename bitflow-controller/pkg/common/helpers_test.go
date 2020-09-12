package common

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type HelpersTestSuite struct {
	AbstractTestSuite
}

func TestHelpersSuite(t *testing.T) {
	suite.Run(t, new(HelpersTestSuite))
}

func (s *HelpersTestSuite) TestGetNodeName() {
	pod := s.Pod("pod")
	s.Empty(GetTargetNode(pod), "Expected no node to be found for new pod")

	expectedNodeName := "node-name"
	SetTargetNodeName(pod, expectedNodeName)
	s.Equal(expectedNodeName, GetTargetNode(pod), "Expected different node to be found")
}

func (s *HelpersTestSuite) TestNodePatchRequired() {
	node := s.SchedulerNode()
	pod := s.Pod("pod1")
	SetTargetNode(pod, node)
	nodeVal := pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values[0]
	s.Equal(node.Name, nodeVal)
}
