package bitflow

import (
	"strconv"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/scheduler"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	corev1 "k8s.io/api/core/v1"
)

type SchedulerTestSuite struct {
	BitflowControllerTestHelpers
}

func (s *BitflowControllerTestSuite) TestScheduler() {
	s.SubTestSuite(new(SchedulerTestSuite))
}

func (s *SchedulerTestSuite) TestGetAllNodes() {
	// Mark Node 2 as non-ready
	node2 := s.Node("node2")
	node2.Status.Conditions = []corev1.NodeCondition{{Type: corev1.NodeReady, Status: corev1.ConditionFalse}}

	r := s.initReconciler(s.Node("node1"), node2, s.Node("node3"), s.Node("node4"), s.Node("node5"))

	nodeList, err := common.RequestReadyNodes(r.client)
	s.NoError(err)
	s.Len(nodeList.Items, 4)
}

func (s *SchedulerTestSuite) testScheduler(schedulerName string, numReconciles int, expectedPodsOnNodes []int) {
	s.SubTest(schedulerName, func() {
		scheduler.Seed(42)     // Make tests reproducible
		totalExpectedPods := 8 // Based on the steps defined below

		numNodes := len(expectedPodsOnNodes)
		objects := []runtime.Object{
			s.Source("source1", map[string]string{"type": "input", "s": "a", "bitflow-nodename": "node" + strconv.Itoa((numNodes+0)%numNodes)}),
			s.Source("source2", map[string]string{"type": "input", "s": "b", "bitflow-nodename": "node" + strconv.Itoa((numNodes+1)%numNodes)}),
			s.Source("source3", map[string]string{"type": "input", "s": "c", "bitflow-nodename": "node" + strconv.Itoa((numNodes+2)%numNodes)}),
			s.Source("source4", map[string]string{"type": "input", "s": "d", "bitflow-nodename": "node" + strconv.Itoa((numNodes+3)%numNodes)}),
			s.StepWithOutput("step-one", v1.StepTypeOneToOne,
				"out", map[string]string{"type": "processed"}, "type", "input"),
			s.StepWithOutput("step-all", v1.StepTypeAllToOne,
				"out", map[string]string{"type": "processed"}, "type", "input"),
			s.StepWithOutput("step-filter", v1.StepTypeOneToOne,
				"out", map[string]string{"type": "filtered"}, "type", "input", "s", "a"),
			s.StepWithOutput("step-filter-processed", v1.StepTypeAllToOne,
				"out", map[string]string{"type": "processed"}, "type", "filtered"),
			s.Step("step-processed", v1.StepTypeAllToOne, "type", "processed"),
		}
		for i := range expectedPodsOnNodes {
			objects = append(objects, s.Node("node"+strconv.Itoa(i)))
		}
		r := s.initReconciler(objects...)

		// Configure the tested scheduler
		s.SetConfigValue(r.client, "bitflow-config", "scheduler", schedulerName)

		// Make sure all pods and output sources are created
		for i := 0; i < numReconciles; i++ {
			// This fake step makes the controller consider all steps at once, instead of checking the one-by-one
			_, err := s.performReconcile(r, ReconcileLoopFakeStepName)
			s.NoError(err)
			s.assignIPToPods(r.client)
		}

		podsOnNodes := make([]int, len(expectedPodsOnNodes))
		sumExpectedPods := 0
		sumEncounteredPods := 0
		for i, expectedPodsOnNode := range expectedPodsOnNodes {
			nodeName := "node" + strconv.Itoa(i)
			sumExpectedPods += expectedPodsOnNode
			count, err := s.getNumberOfPodsForNode(r.client, nodeName)
			s.NoError(err)
			podsOnNodes[i] = count
			sumEncounteredPods += count
		}
		s.Equal(totalExpectedPods, sumExpectedPods, "Wrong number of expected pods provided")
		s.Equal(totalExpectedPods, sumEncounteredPods, "Wrong number of scheduled pods")

		s.Equal(expectedPodsOnNodes, podsOnNodes, "Wrong schedule")
	})
}

func (s *SchedulerTestSuite) testSchedulers(numReconciles int, schedulers map[string][]int) {
	schedulers["WRONG-SCHEDULER"] = schedulers[SchedulerNameDefault]
	for schedulerName, distribution := range schedulers {
		s.testScheduler(schedulerName, numReconciles, distribution)
	}
}

func (s *SchedulerTestSuite) TestSchedule_1Node() {
	s.testSchedulers(3, map[string][]int{
		SchedulerNameLeastOccupied: {8},
		SchedulerNameRandom:        {8},
	})
}

func (s *SchedulerTestSuite) TestSchedule_2Nodes() {
	s.testSchedulers(4, map[string][]int{
		SchedulerNameLeastOccupied: {5, 3}, // Not 4,4 due to order of events and limitations of the naive scheduler
		SchedulerNameRandom:        {6, 2},
	})
}

func (s *SchedulerTestSuite) TestSchedule_3Nodes() {
	s.testSchedulers(6, map[string][]int{
		SchedulerNameLeastOccupied: {5, 2, 1}, // See comment above
		SchedulerNameRandom:        {6, 1, 1},
	})
}

func (s *SchedulerTestSuite) TestSchedule_4Nodes() {
	s.testSchedulers(10, map[string][]int{
		SchedulerNameLeastOccupied: {4, 2, 1, 1}, // See comment above
		SchedulerNameRandom:        {3, 3, 1, 1},
	})
}

func (s *SchedulerTestSuite) TestSchedule_5Nodes() {
	s.testSchedulers(10, map[string][]int{
		SchedulerNameLeastOccupied: {4, 2, 1, 1, 0}, // See comment above
		SchedulerNameRandom:        {3, 3, 1, 1, 0},
	})
}
