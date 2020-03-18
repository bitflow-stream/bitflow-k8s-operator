package bitflow

import (
	"fmt"
)

type MiscTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestMisc() {
	s.SubTestSuite(new(MiscTestSuite))
}

func (s *MiscTestSuite) TestHashFunction() {
	id := HashName("pre", "step-1", "source-1")
	s.Equal(HASH_SUFFIX_LENGTH, len(id)-len("pre"), "Length of identifier does not have the correct length")
	id2 := HashName("pre", "step-1", "source-1")
	s.Equal(id, id2, "Hash must be consistent")
}

func (s *MiscTestSuite) TestBuildDataSource() {
	urlString := buildDataSource("podname", "127.0.0.1", 8888)
	s.NotEmpty(urlString, "Expected a valid url but it is empty")
}

func (s *MiscTestSuite) TestRequeueError() {
	err := fmt.Errorf("Random error")
	requeue := NewRequeueError(err)
	s.False(isRequeueError(err))
	s.True(isRequeueError(requeue))
	s.False(isRequeueError(nil))
}

func (s *MiscTestSuite) TestOwnerRefs() {
	pod := s.Pod("name")
	source := s.Source("source1", nil)
	source.UID = "456"

	setOwnerReferenceForPod(pod, source)
	s.NotEmpty(pod.OwnerReferences)
}
