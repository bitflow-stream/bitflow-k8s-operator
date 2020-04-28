package scheduler

import (
	"testing"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/config"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type SchedulerTestSuite struct {
	common.AbstractTestSuite
}

func TestScheduler(t *testing.T) {
	suite.Run(t, new(SchedulerTestSuite))
}

func (s *SchedulerTestSuite) getSchedulerNode() *corev1.Node {
	return s.Node2("node1",
		map[string]string{"test-node": "yes", HostnameLabel: "node1"},
		map[string]string{"bitflow-resource-limit": "0.1"})
}

func (s *SchedulerTestSuite) getNodeWithResources(name string, cpu int64, memory int64) *corev1.Node {
	return s.NodeWithResources(name,
		map[string]string{"test-node": "yes", HostnameLabel: name},
		map[string]string{"bitflow-resource-limit": "0.1"},
		cpu,
		memory*1024*1024)
}

func (s *SchedulerTestSuite) getScheduler(objects ...runtime.Object) *Scheduler {
	configMap := s.ConfigMap("bitflow-config")
	configMap.Data["schedulers"] = "" // Make sure there are not default schedulers
	objects = append(objects, configMap)
	cl := s.MakeFakeClient(objects...)
	conf := config.NewConfig(cl, common.TestNamespace, "bitflow-config")
	return &Scheduler{cl, conf, common.TestNamespace, map[string]string{"bitflow": "true"}}
}

func (s *SchedulerTestSuite) testSchedulerMultipleNodes(testName string, expectedSuccessfulScheduler string, schedulerList string, expectedNode *corev1.Node, nodes ...*corev1.Node) {
	s.SubTest(testName, func() {
		// Setup
		sources := []*bitflowv1.BitflowSource(nil)
		scheduledPod := s.Pod("scheduled-pod")
		step := s.Step("test-step", bitflowv1.StepTypeOneToOne)

		var runtimeObjectNodes = make([]runtime.Object, len(nodes))
		for i, node := range nodes {
			runtimeObjectNodes[i] = node
		}
		scheduler := s.getScheduler(runtimeObjectNodes...)
		step.Spec.Scheduler = schedulerList

		//Execution
		nodeScheduled, successfulScheduler := scheduler.SchedulePod(scheduledPod, step, sources)

		// Assertions
		s.Equal(expectedSuccessfulScheduler, successfulScheduler)
		if expectedNode != nil {
			s.NotNil(nodeScheduled)
			s.Equal(expectedNode.Name, nodeScheduled.Name)
		}
	})
}

func (s *SchedulerTestSuite) TestSchedulersMultipleNodes() {
	firstNode := s.getNodeWithResources("firstNode", 10, 1)
	mostCpuNode := s.getNodeWithResources("mostCpuNode", 1000, 1)
	modeMemNode := s.getNodeWithResources("modeMemNode", 100, 100)

	otherNode1 := s.getNodeWithResources("otherNode1", 564, 99)
	otherNode2 := s.getNodeWithResources("otherNode2", 239, 12)
	otherNode3 := s.getNodeWithResources("otherNode3", 0, 2)
	otherNode4 := s.getNodeWithResources("otherNode4", 999, 0)
	otherNode5 := s.getNodeWithResources("otherNode5", 388, 3)

	s.testSchedulerMultipleNodes("mostCPU",
		"mostCPU", "mostCPU,mostMem,first,random", mostCpuNode,
		mostCpuNode, modeMemNode, otherNode1, otherNode2, otherNode3, otherNode4, otherNode5)

	s.testSchedulerMultipleNodes("mostMem",
		"mostMem", "mostMem,mostCPU,first,random", modeMemNode,
		mostCpuNode, modeMemNode, otherNode1, otherNode2, otherNode3, otherNode4, otherNode5)

	s.testSchedulerMultipleNodes("first",
		"first", "first,mostCPU,mostMem,random", firstNode,
		firstNode, mostCpuNode, modeMemNode, otherNode1, otherNode2, otherNode3, otherNode4)

	s.testSchedulerMultipleNodes("randomOnlyOneNode",
		"random", "random,first", otherNode1,
		otherNode1)

	s.testSchedulerMultipleNodes("noNodesShouldNotFindScheduler",
		"", "random", nil)

	// TODO test leastContainers scheduler
	// TODO test sourceAffinity scheduler
}

func (s *SchedulerTestSuite) testSimpleScheduler(scheduler *Scheduler, schedulerList string, sources []*bitflowv1.BitflowSource, expectedSuccessfulScheduler string, expectedNode *corev1.Node) {
	s.SubTest(schedulerList, func() {
		scheduledPod := s.Pod("scheduled-pod")
		step := s.Step("test-step", bitflowv1.StepTypeOneToOne)
		step.Spec.Scheduler = schedulerList
		node, successfulScheduler := scheduler.SchedulePod(scheduledPod, step, sources)
		s.Equal(expectedSuccessfulScheduler, successfulScheduler)
		if expectedNode == nil {
			s.Nil(node)
		} else {
			s.NotNil(node)
			s.Equal(expectedNode.Name, node.Name)
		}
	})
}

func (s *SchedulerTestSuite) TestSimpleSchedulers() {
	node := s.getSchedulerNode()
	scheduler := s.getScheduler(node)
	sources := []*bitflowv1.BitflowSource(nil)

	s.testSimpleScheduler(scheduler, "first", sources, "first", node)
	s.testSimpleScheduler(scheduler, "random", sources, "random", node)
	s.testSimpleScheduler(scheduler, "leastContainers", sources, "leastContainers", node)
	s.testSimpleScheduler(scheduler, "mostCPU", sources, "mostCPU", node)
	s.testSimpleScheduler(scheduler, "mostMem", sources, "mostMem", node)

	s.testSimpleScheduler(scheduler, "WRONG_SCHEDULER_NAME", sources, "", nil)

	s.testSimpleScheduler(scheduler, "first,random", sources, "first", node)
	s.testSimpleScheduler(scheduler, "random,first", sources, "random", node)
	s.testSimpleScheduler(scheduler, "first,random,random", sources, "first", node)
	s.testSimpleScheduler(scheduler, "first,random,WRONG,random", sources, "first", node)
	s.testSimpleScheduler(scheduler, "first,WRONG,random", sources, "first", node)

	s.testSimpleScheduler(scheduler, "WRONG,random,first", sources, "random", node)
	s.testSimpleScheduler(scheduler, "WRONG1,WRONG2,first,random", sources, "first", node)

	// sourceAffinity always fails here, since there are no data sources defined
	s.testSimpleScheduler(scheduler, "sourceAffinity", sources, "", nil)
	s.testSimpleScheduler(scheduler, "sourceAffinity,first", sources, "first", node)
	s.testSimpleScheduler(scheduler, "sourceAffinity,WRONG,first", sources, "first", node)
	s.testSimpleScheduler(scheduler, "sourceAffinity,WRONG,WRONG", sources, "", nil)
	s.testSimpleScheduler(scheduler, "WRONG,sourceAffinity,WRONG,random,first", sources, "random", node)
}

func (s *SchedulerTestSuite) TestNodePatchPreferred() {
	node := s.getSchedulerNode()
	pod := s.Pod("pod1")
	SetPodNodeAffinityPreferred(node, pod)
	nodeVal := pod.Spec.Affinity.NodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution[0].Preference.MatchExpressions[0].Values[0]
	s.Equal(node.Name, nodeVal)
}

func (s *SchedulerTestSuite) TestNodePatchRequired() {
	node := s.getSchedulerNode()
	pod := s.Pod("pod1")
	SetPodNodeAffinityRequired(node, pod)
	nodeVal := pod.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution.NodeSelectorTerms[0].MatchExpressions[0].Values[0]
	s.Equal(node.Name, nodeVal)
}
