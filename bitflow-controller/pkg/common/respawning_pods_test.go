package common

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RespawningPodsTestSuite struct {
	AbstractTestSuite
}

func TestRespawningPods(t *testing.T) {
	suite.Run(t, new(RespawningPodsTestSuite))
}

func (s *RespawningPodsTestSuite) TestAdd() {
	name1 := "name1"
	name2 := "name2"
	pod1 := s.Pod(name1)
	pod2 := s.Pod(name2)
	pod1.Status.PodIP = "1.2.3.4"
	pod2.Status.PodIP = "4.5.6.7"

	pods := NewRespawningPods()
	pods.Add(pod1)
	pods.Add(pod2)

	_, restarting1 := pods.IsPodRestarting(name1)
	s.True(restarting1, "Pod should be restarting")

	pods.Delete(name2)
	_, restarting2 := pods.IsPodRestarting(name2)
	s.False(restarting2, "Pod should be deleted")
}

func (s *RespawningPodsTestSuite) TestDeleteWithLabels() {
	name1 := "helloStep1-pod-1"
	name2 := "helloStep1-pod-2"
	name3 := "helloStep1-pod-3"
	name4 := "helloStep2-pod-1"
	name5 := "helloStep2-pod-2"
	pod1 := s.PodLabels(name1, map[string]string{"bitflow-step-name": "helloStep1"})
	pod2 := s.PodLabels(name2, map[string]string{"bitflow-step-name": "helloStep1"})
	pod3 := s.PodLabels(name3, map[string]string{"bitflow-step-name": "helloStep1"})
	pod4 := s.PodLabels(name4, map[string]string{"bitflow-step-name": "helloStep2"})
	pod5 := s.PodLabels(name5, map[string]string{"bitflow-step-name": "helloStep2"})
	pod1.Status.PodIP = "4.5.6.7"
	pod2.Status.PodIP = "4.5.6.7"
	pod3.Status.PodIP = "4.5.6.7"
	pod5.Status.PodIP = "4.5.6.7"
	pod4.Status.PodIP = "4.5.6.7"

	pods := NewRespawningPods()
	pods.Add(pod1)
	pods.Add(pod2)
	pods.Add(pod3)
	pods.Add(pod4)
	pods.Add(pod5)
	s.Len(pods.ListPods(), 5)

	pods.DeletePodsWithLabel("bitflow-step-name", "helloStep1")
	s.Len(pods.ListPods(), 2)

	_, ok1 := pods.IsPodRestarting(name1)
	s.False(ok1)
	_, ok2 := pods.IsPodRestarting(name2)
	s.False(ok2)
	_, ok3 := pods.IsPodRestarting(name3)
	s.False(ok3)
	_, ok4 := pods.IsPodRestarting(name4)
	s.True(ok4)
	_, ok5 := pods.IsPodRestarting(name5)
	s.True(ok5)
}

func (s *RespawningPodsTestSuite) TestDeleteExcept() {
	name1 := "helloStep1-pod-1"
	name2 := "helloStep1-pod-2"
	name3 := "helloStep1-pod-3"
	name4 := "helloStep2-pod-1"
	name5 := "helloStep2-pod-2"
	pod1 := s.PodLabels(name1, map[string]string{"bitflow-step-name": "helloStep1"})
	pod2 := s.PodLabels(name2, map[string]string{"bitflow-step-name": "helloStep1"})
	pod3 := s.PodLabels(name3, map[string]string{"bitflow-step-name": "helloStep1"})
	pod4 := s.PodLabels(name4, map[string]string{"bitflow-step-name": "helloStep2"})
	pod5 := s.PodLabels(name5, map[string]string{"bitflow-step-name": "helloStep2"})
	pod1.Status.PodIP = "1.2.3"
	pod2.Status.PodIP = "4.5.6"
	pod3.Status.PodIP = "4.5.6"
	pod4.Status.PodIP = "4.5.6"
	pod5.Status.PodIP = "4.5.6"

	pods := NewRespawningPods()
	pods.Add(pod1)
	pods.Add(pod2)
	pods.Add(pod3)
	pods.Add(pod4)
	pods.Add(pod5)
	s.Len(pods.ListPods(), 5)

	pods.DeletePodsWithLabelExcept("bitflow-step-name", "helloStep1", []string{name2})
	s.Len(pods.ListPods(), 3)

	_, ok1 := pods.IsPodRestarting(name1)
	s.False(ok1)
	_, ok2 := pods.IsPodRestarting(name2)
	s.True(ok2)
	_, ok3 := pods.IsPodRestarting(name3)
	s.False(ok3)
	_, ok4 := pods.IsPodRestarting(name4)
	s.True(ok4)
	_, ok5 := pods.IsPodRestarting(name5)
	s.True(ok5)
}

func (s *RespawningPodsTestSuite) TestRespawningOnNode() {
	name := "pod-1"
	pod := s.Pod(name)
	pod.Status.PodIP = "1.2.3.4"
	node := "restartNode"
	pod.Spec.NodeName = node

	pods := NewRespawningPods()
	pods.Add(pod)

	_, ok := pods.IsPodRestartingOnNode("xxx", node)
	s.False(ok)
	_, ok = pods.IsPodRestartingOnNode(name, "xxx")
	s.False(ok)

	podStatus, ok := pods.IsPodRestartingOnNode(name, node)
	s.True(ok)
	s.Equal(node, podStatus.Pod.Spec.NodeName)
}
