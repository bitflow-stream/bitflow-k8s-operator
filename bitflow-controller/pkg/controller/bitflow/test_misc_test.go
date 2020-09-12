package bitflow

import (
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
)

type MiscTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestMisc() {
	s.SubTestSuite(new(MiscTestSuite))
}

func (s *MiscTestSuite) TestHashFunction() {
	id := common.HashName("pre", "step-1", "source-1")
	s.Equal(common.HashSuffixLength, len(id)-len("pre"), "Length of identifier does not have the correct length")
	id2 := common.HashName("pre", "step-1", "source-1")
	s.Equal(id, id2, "Hash must be consistent")
}

func (s *MiscTestSuite) TestBuildDataSource() {
	urlString := buildDataSource("podname", "127.0.0.1", 8888)
	s.NotEmpty(urlString, "Expected a valid url but it is empty")
}

func (s *MiscTestSuite) TestOwnerRefs() {
	pod := s.Pod("name")
	source := s.Source("source1", nil)
	source.UID = "456"

	setOwnerReferenceForPod(pod, source)
	s.NotEmpty(pod.OwnerReferences)
}
