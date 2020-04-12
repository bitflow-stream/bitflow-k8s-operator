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
	s.Empty(GetNodeName(pod), "Expected no node to be found for new pod")

	expectedNodeName := "node-name"
	pod.Spec.NodeName = expectedNodeName
	s.Equal(expectedNodeName, GetNodeName(pod), "Expected different node to be found")

	// TODO commented out due to dependency to 'scheduler' package... The node name of the pod SHOULD always be correct in pod.spec.node
	// labels := map[string]string{scheduler.HostnameLabel: expectedNodeName}
	// node := s.Node2(expectedNodeName, labels, nil)
	// scheduler.PatchPodNodeAffinityRequired(node, pod)
	// pod.Spec.NodeName = ""
	// s.Equal(expectedNodeName, GetNodeName(pod), "Expected different node to be found")
}
