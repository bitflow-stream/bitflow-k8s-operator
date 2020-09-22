package skd

import (
	"fmt"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SkdTestSuite struct {
	common.AbstractTestSuite
}

func TestSkd(t *testing.T) {
	suite.Run(t, new(SkdTestSuite))
}

func (s *SkdTestSuite) Test_EqualDistributionScheduler_shouldReturnCorrectMap() {
	var scheduler Scheduler
	scheduler = EqualDistributionScheduler{
		nodeNames: []string{"n1", "n2", "n3"},
		podNames:  []string{"p1", "p2", "p3", "p4", "p5", "p6", "p7", "p8", "p9", "p10"},
	}

	schedulingChanged, scheduledMap, err := scheduler.Schedule()

	s.Nil(err)
	s.True(schedulingChanged)
	s.Equal("n1", scheduledMap["p1"])
	s.Equal("n2", scheduledMap["p2"])
	s.Equal("n3", scheduledMap["p3"])
	s.Equal("n1", scheduledMap["p4"])
	s.Equal("n2", scheduledMap["p5"])
	s.Equal("n3", scheduledMap["p6"])
	s.Equal("n1", scheduledMap["p7"])
	s.Equal("n2", scheduledMap["p8"])
	s.Equal("n3", scheduledMap["p9"])
	s.Equal("n1", scheduledMap["p10"])
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldMapPodsCorrectlyWithoutNetworkPenalty() {
	var scheduler AdvancedScheduler
	scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				memory:                  64,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
			{
				name:                    "n2",
				allocatableCpu:          4000,
				memory:                  64,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
			{
				name:                    "n3",
				allocatableCpu:          4000,
				memory:                  64,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				receivesDataFrom: []string{},
				sendsDataTo:      []string{"p7", "p8"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
			{
				name:             "p2",
				receivesDataFrom: []string{"p10"},
				sendsDataTo:      []string{"p9"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
			{
				name:             "p3",
				receivesDataFrom: []string{},
				sendsDataTo:      []string{"p9"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
			{
				name:             "p4",
				receivesDataFrom: []string{},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
			{
				name:             "p5",
				receivesDataFrom: []string{},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
			{
				name:             "p6",
				receivesDataFrom: []string{},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
			{
				name:             "p7",
				receivesDataFrom: []string{"p1"},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
			{
				name:             "p8",
				receivesDataFrom: []string{"p1"},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
			{
				name:             "p9",
				receivesDataFrom: []string{"p2", "p3"},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
			{
				name:             "p10",
				receivesDataFrom: []string{},
				sendsDataTo:      []string{"p2"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
			{
				name:             "p11",
				receivesDataFrom: []string{},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
			{
				name:             "p12",
				receivesDataFrom: []string{},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
			{
				name:             "p13",
				receivesDataFrom: []string{},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
		},
		networkPenalty:   0,
		memoryPenalty:    10,
		thresholdPercent: 10,
	}

	schedulingChanged, scheduledMap, err := scheduler.ScheduleCheckingAllPermutations()

	s.Nil(err)
	s.True(schedulingChanged)
	s.Equal("n1", scheduledMap["p1"])
	s.Equal("n1", scheduledMap["p2"])
	s.Equal("n1", scheduledMap["p3"])
	s.Equal("n1", scheduledMap["p4"])
	s.Equal("n2", scheduledMap["p5"])
	s.Equal("n2", scheduledMap["p6"])
	s.Equal("n2", scheduledMap["p7"])
	s.Equal("n2", scheduledMap["p8"])
	s.Equal("n3", scheduledMap["p9"])
	s.Equal("n3", scheduledMap["p10"])
	s.Equal("n3", scheduledMap["p11"])
	s.Equal("n3", scheduledMap["p12"])
	s.Equal("n3", scheduledMap["p13"])
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldMapPodsCorrectlyWithNetworkPenalty() {
	var scheduler AdvancedScheduler
	nodes := []*NodeData{
		{
			name:                    "n1",
			allocatableCpu:          4000,
			memory:                  64,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		{
			name:                    "n2",
			allocatableCpu:          4000,
			memory:                  64,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		{
			name:                    "n3",
			allocatableCpu:          4000,
			memory:                  64,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
	}
	pods := []*PodData{
		{
			name:             "p1",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p2",
			receivesDataFrom: []string{"p10"},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p3",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p4",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p5",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p6",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p7",
			receivesDataFrom: []string{"p1"},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p8",
			receivesDataFrom: []string{"p1"},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p9",
			receivesDataFrom: []string{"p3", "p4", "p6", "p7", "p8"},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p10",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p11",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p12",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p13",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
	}
	scheduler = AdvancedScheduler{
		nodes:            nodes,
		pods:             pods,
		networkPenalty:   1_000_000,
		memoryPenalty:    0,
		thresholdPercent: 10,
	}

	schedulingChanged, scheduledMap, err := scheduler.ScheduleCheckingAllPermutations()

	s.Nil(err)
	s.True(schedulingChanged)
	s.Equal(scheduledMap["p2"], scheduledMap["p10"])
	s.Equal(scheduledMap["p7"], scheduledMap["p1"])
	s.Equal(scheduledMap["p8"], scheduledMap["p1"])
	s.Equal(scheduledMap["p9"], scheduledMap["p3"])
	s.Equal(scheduledMap["p9"], scheduledMap["p4"])
	s.Equal(scheduledMap["p9"], scheduledMap["p6"])
	s.Equal(scheduledMap["p9"], scheduledMap["p7"])
	s.Equal(scheduledMap["p9"], scheduledMap["p8"])
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldRecognizeSchedulingHasNotChangedWithThreshold() {
	var scheduler AdvancedScheduler
	nodes := []*NodeData{
		{
			name:                    "n1",
			allocatableCpu:          4000,
			memory:                  64,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		{
			name:                    "n2",
			allocatableCpu:          4000,
			memory:                  64,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		{
			name:                    "n3",
			allocatableCpu:          4000,
			memory:                  64,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
	}
	pods := []*PodData{
		{
			name:             "p1",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p2",
			receivesDataFrom: []string{"p10"},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p3",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p4",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p5",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p6",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p7",
			receivesDataFrom: []string{"p1"},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
	}
	scheduler = AdvancedScheduler{
		nodes:            nodes,
		pods:             pods,
		thresholdPercent: 10,
	}

	schedulingChanged1, scheduledMap, err1 := scheduler.ScheduleCheckingAllPermutations()

	s.Nil(err1)
	s.True(schedulingChanged1)

	scheduler = AdvancedScheduler{
		nodes:              nodes,
		pods:               pods,
		thresholdPercent:   10,
		previousScheduling: scheduledMap,
	}

	schedulingChanged2, _, err2 := scheduler.ScheduleCheckingAllPermutations()

	s.Nil(err2)
	s.False(schedulingChanged2)
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldRecognizeSchedulingHasChangedWithoutThreshold() {
	var scheduler AdvancedScheduler
	nodes := []*NodeData{
		{
			name:                    "n1",
			allocatableCpu:          4000,
			memory:                  64,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		{
			name:                    "n2",
			allocatableCpu:          4000,
			memory:                  64,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		{
			name:                    "n3",
			allocatableCpu:          4000,
			memory:                  64,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
	}
	pods := []*PodData{
		{
			name:             "p1",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p2",
			receivesDataFrom: []string{"p10"},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p3",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p4",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p5",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p6",
			receivesDataFrom: []string{},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
		{
			name:             "p7",
			receivesDataFrom: []string{"p1"},
			curve: Curve{
				a: 6.71881241016441,
				b: 0.0486498280492762,
				c: 2.0417306475862214,
				d: 15.899403720950454,
			},
			minimumMemory: 16,
		},
	}
	scheduler = AdvancedScheduler{
		nodes: nodes,
		pods:  pods,
	}

	schedulingChanged1, scheduledMap, err1 := scheduler.ScheduleCheckingAllPermutations()

	s.Nil(err1)
	s.True(schedulingChanged1)

	scheduler = AdvancedScheduler{
		nodes:              nodes,
		pods:               pods,
		previousScheduling: scheduledMap,
	}

	schedulingChanged2, _, err2 := scheduler.ScheduleCheckingAllPermutations()

	s.Nil(err2)
	s.True(schedulingChanged2)
}

func (s *SkdTestSuite) testNewDistributionPenaltyLowerConsideringThreshold(previousPenalty float64, newPenalty float64, thresholdPercent float64, expectedOutcome bool) {
	s.SubTest(fmt.Sprintf("previous%f:new%f:threshold%f->%v", previousPenalty, newPenalty, thresholdPercent, expectedOutcome), func() {
		actualOutcome := NewDistributionPenaltyIsLowerConsideringThreshold(previousPenalty, newPenalty, thresholdPercent)
		s.Equal(expectedOutcome, actualOutcome)
	})
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldDetermineIfNewDistributionPenaltyIsLowerConsideringThreshold() {
	s.testNewDistributionPenaltyLowerConsideringThreshold(20, 10, 0, true)
	s.testNewDistributionPenaltyLowerConsideringThreshold(10, 20, 0, false)
	s.testNewDistributionPenaltyLowerConsideringThreshold(1000, 500, 50, true)
	s.testNewDistributionPenaltyLowerConsideringThreshold(1000, 500.000001, 50, false)
	s.testNewDistributionPenaltyLowerConsideringThreshold(100, 90, 10, true)
	s.testNewDistributionPenaltyLowerConsideringThreshold(100, 90.000001, 10, false)
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldNotThrowErrorWhenAllNecessaryFieldsAreSet() {
	var scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				memory:                  64,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				receivesDataFrom: []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
		},
	}

	_, _, err := scheduler.ScheduleCheckingAllPermutations()

	s.Nil(err)
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldThrowErrorWhenNodeMemoryIsMissing() {
	var scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				receivesDataFrom: []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
		},
	}

	_, _, err := scheduler.ScheduleCheckingAllPermutations()

	s.NotNil(err)
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldThrowErrorWhenNodeMemoryIs0() {
	var scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				memory:                  0,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				receivesDataFrom: []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
		},
	}

	_, _, err := scheduler.ScheduleCheckingAllPermutations()

	s.NotNil(err)
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldThrowErrorWhenNodeResourceLimitIsMissing() {
	var scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				memory:                  64,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				receivesDataFrom: []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
		},
	}

	_, _, err := scheduler.ScheduleCheckingAllPermutations()

	s.NotNil(err)
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldThrowErrorWhenNodeResourceLimitIs0() {
	var scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				memory:                  64,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				receivesDataFrom: []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
		},
	}

	_, _, err := scheduler.ScheduleCheckingAllPermutations()

	s.NotNil(err)
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldThrowErrorWhenPodMinimumMemoryIsMissing() {
	var scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				memory:                  64,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				receivesDataFrom: []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
			},
		},
	}

	_, _, err := scheduler.ScheduleCheckingAllPermutations()

	s.NotNil(err)
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldThrowErrorWhenPodMinimumMemoryIs0() {
	var scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				memory:                  64,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				receivesDataFrom: []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 0,
			},
		},
	}

	_, _, err := scheduler.ScheduleCheckingAllPermutations()

	s.NotNil(err)
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldThrowErrorWhenThresholdIsLessThan0() {
	var scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				memory:                  64,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				receivesDataFrom: []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
		},
		thresholdPercent: -1,
	}

	_, _, err := scheduler.ScheduleCheckingAllPermutations()

	s.NotNil(err)
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldThrowErrorWhenThresholdIsGreaterThan100() {
	var scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				memory:                  64,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				receivesDataFrom: []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory: 16,
			},
		},
		thresholdPercent: 101,
	}

	_, _, err := scheduler.ScheduleCheckingAllPermutations()

	s.NotNil(err)
}

func (s *SkdTestSuite) Test_shouldSortPodsTopologically() {
	sortedPods, err := sortPodsUsingKahnsAlgorithm([]*PodData{
		{
			name:             "011(also-sends-to-010)",
			receivesDataFrom: []string{"01"},
			sendsDataTo:      []string{"010(also-sends-to-001)"},
		},
		{
			name:             "010(also-sends-to-001)",
			receivesDataFrom: []string{"01", "011(also-sends-to-010)"},
			sendsDataTo:      []string{"001"},
		},
		{
			name:             "001",
			receivesDataFrom: []string{"00(also-sends-to-010)", "010(also-sends-to-001)"},
			sendsDataTo:      []string{},
		},
		{
			name:             "003",
			receivesDataFrom: []string{"00(also-sends-to-010)"},
			sendsDataTo:      []string{},
		},
		{
			name:             "002",
			receivesDataFrom: []string{"00(also-sends-to-010)"},
			sendsDataTo:      []string{},
		},
		{
			name:             "00(also-sends-to-010)",
			receivesDataFrom: []string{"0"},
			sendsDataTo:      []string{"001", "002", "003", "010(also-sends-to-001)"},
		},
		{
			name:             "0",
			receivesDataFrom: []string{},
			sendsDataTo:      []string{"00(also-sends-to-010)", "01"},
		},
		{
			name:             "01",
			receivesDataFrom: []string{"0"},
			sendsDataTo:      []string{"010(also-sends-to-001)", "011(also-sends-to-010)"},
			dataSourceNodes:  []string{"test"},
			curve: Curve{
				a: 1,
				b: 2,
				c: 3,
				d: 4,
			},
			minimumMemory:        22,
			maximumExecutionTime: 23,
		},
	})

	s.Nil(err)
	s.Equal("0", sortedPods[0].name)
	s.Equal("00(also-sends-to-010)", sortedPods[1].name)
	s.Equal("01", sortedPods[2].name)
	s.Equal("002", sortedPods[3].name)
	s.Equal("003", sortedPods[4].name)
	s.Equal("011(also-sends-to-010)", sortedPods[5].name)
	s.Equal("010(also-sends-to-001)", sortedPods[6].name)
	s.Equal("001", sortedPods[7].name)

	// making sure other fields are copied properly
	s.Equal("test", sortedPods[2].dataSourceNodes[0])
	s.Equal(1.0, sortedPods[2].curve.a)
	s.Equal(2.0, sortedPods[2].curve.b)
	s.Equal(3.0, sortedPods[2].curve.c)
	s.Equal(4.0, sortedPods[2].curve.d)
	s.Equal(22.0, sortedPods[2].minimumMemory)
	s.Equal(23.0, sortedPods[2].maximumExecutionTime)
}

func (s *SkdTestSuite) Test_shouldReturnIndicesSortedByNumberOfPods() {
	nodeStates := []NodeState{
		{
			node: &NodeData{
				name: "n1",
			},
			pods: []*PodData{
				{
					name: "p1",
				},
				{
					name: "p2",
				},
				{
					name: "p3",
				},
				{
					name: "p4",
				},
			},
		},
		{
			node: &NodeData{
				name: "n2",
			},
			pods: []*PodData{},
		},
		{
			node: &NodeData{
				name: "n3",
			},
			pods: []*PodData{
				{
					name: "p5",
				},
				{
					name: "p6",
				},
				{
					name: "p7",
				},
			},
		},
		{
			node: &NodeData{
				name: "n4",
			},
			pods: []*PodData{
				{
					name: "p8",
				},
			},
		},
		{
			node: &NodeData{
				name: "n5",
			},
			pods: []*PodData{
				{
					name: "p9",
				},
				{
					name: "p10",
				},
			},
		},
	}
	indexSlice := getIndexSliceSortedByNumberOfPods(nodeStates)

	s.Equal(5, len(indexSlice))
	s.Equal("n2", nodeStates[indexSlice[0]].node.name)
	s.Equal("n4", nodeStates[indexSlice[1]].node.name)
	s.Equal("n5", nodeStates[indexSlice[2]].node.name)
	s.Equal("n3", nodeStates[indexSlice[3]].node.name)
	s.Equal("n1", nodeStates[indexSlice[4]].node.name)
}

func (s *SkdTestSuite) Test_shouldGetNodeIndexOfNodeStateContainingPod() {
	systemState := SystemState{
		nodes: []NodeState{
			{
				node: &NodeData{
					name: "n1",
				},
				pods: []*PodData{
					{
						name: "p1",
					},
					{
						name: "p2",
					},
					{
						name: "p3",
					},
					{
						name: "p4",
					},
				},
			},
			{
				node: &NodeData{
					name: "n2",
				},
				pods: []*PodData{},
			},
			{
				node: &NodeData{
					name: "n3",
				},
				pods: []*PodData{
					{
						name: "p5",
					},
					{
						name: "p6",
					},
					{
						name: "p7",
					},
				},
			},
			{
				node: &NodeData{
					name: "n4",
				},
				pods: []*PodData{
					{
						name: "p8",
					},
				},
			},
			{
				node: &NodeData{
					name: "n5",
				},
				pods: []*PodData{
					{
						name: "p9",
					},
					{
						name: "p10",
					},
				},
			},
		},
	}

	index, err := getNodeIndexOfNodeStateContainingPod("p6", systemState)

	s.Nil(err)
	s.Equal("n3", systemState.nodes[index].node.name)

	_, err = getNodeIndexOfNodeStateContainingPod("pDoesNotExist", systemState)

	s.NotNil(err)
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldFindGoodScheduling() {
	var scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				memory:                  640,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
			{
				name:                    "n2",
				allocatableCpu:          4000,
				memory:                  640,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				dataSourceNodes:  []string{"n1"},
				receivesDataFrom: []string{},
				sendsDataTo:      []string{"p2"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
			{
				name:             "p2",
				dataSourceNodes:  []string{},
				receivesDataFrom: []string{"p1"},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
			{
				name:             "p3",
				dataSourceNodes:  []string{"n2"},
				receivesDataFrom: []string{},
				sendsDataTo:      []string{"p4"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
			{
				name:             "p4",
				dataSourceNodes:  []string{},
				receivesDataFrom: []string{"p3"},
				sendsDataTo:      []string{"p5"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
			{
				name:             "p5",
				dataSourceNodes:  []string{},
				receivesDataFrom: []string{"p4"},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
		},
		networkPenalty:   10,
		memoryPenalty:    50,
		thresholdPercent: 5,
	}

	_, schedulingMap, err := scheduler.Schedule()

	s.Nil(err)
	s.Equal("n1", schedulingMap["p1"])
	s.Equal("n1", schedulingMap["p2"])
	s.Equal("n2", schedulingMap["p3"])
	s.Equal("n2", schedulingMap["p4"])
	s.Equal("n2", schedulingMap["p5"])
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldFindGoodSchedulingIfNetworkPenaltyPlaysARole() {
	var scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				memory:                  160,
				initialNumberOfPodSlots: 1,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
			{
				name:                    "n2",
				allocatableCpu:          40000,
				memory:                  6400,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				dataSourceNodes:  []string{"n1"},
				receivesDataFrom: []string{},
				sendsDataTo:      []string{"p2"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
			{
				name:             "p2",
				dataSourceNodes:  []string{},
				receivesDataFrom: []string{"p1"},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
			{
				name:             "p3",
				dataSourceNodes:  []string{"n2"},
				receivesDataFrom: []string{},
				sendsDataTo:      []string{"p4"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
			{
				name:             "p4",
				dataSourceNodes:  []string{},
				receivesDataFrom: []string{"p3"},
				sendsDataTo:      []string{"p5"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
			{
				name:             "p5",
				dataSourceNodes:  []string{},
				receivesDataFrom: []string{"p4"},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
		},
		networkPenalty:   10,
		memoryPenalty:    50,
		thresholdPercent: 5,
	}

	_, schedulingMap, err := scheduler.Schedule()

	s.Nil(err)
	s.Equal("n1", schedulingMap["p1"])
	s.Equal("n2", schedulingMap["p2"])
	s.Equal("n2", schedulingMap["p3"])
	s.Equal("n2", schedulingMap["p4"])
	s.Equal("n2", schedulingMap["p5"])
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldFindGoodSchedulingAndPreferExecutionTimeOverrunOverMemoryOverrun() {
	var scheduler = AdvancedScheduler{
		nodes: []*NodeData{
			{
				name:                    "n1",
				allocatableCpu:          4000,
				memory:                  160,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
			{
				name:                    "n2",
				allocatableCpu:          50,
				memory:                  6400,
				initialNumberOfPodSlots: 2,
				podSlotScalingFactor:    2,
				resourceLimit:           0.1,
			},
		},
		pods: []*PodData{
			{
				name:             "p1",
				dataSourceNodes:  []string{"n1"},
				receivesDataFrom: []string{},
				sendsDataTo:      []string{"p2"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
			{
				name:             "p2",
				dataSourceNodes:  []string{},
				receivesDataFrom: []string{"p1"},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
			{
				name:             "p3",
				dataSourceNodes:  []string{"n2"},
				receivesDataFrom: []string{},
				sendsDataTo:      []string{"p4"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
			{
				name:             "p4",
				dataSourceNodes:  []string{},
				receivesDataFrom: []string{"p3"},
				sendsDataTo:      []string{"p5"},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
			{
				name:             "p5",
				dataSourceNodes:  []string{},
				receivesDataFrom: []string{"p4"},
				sendsDataTo:      []string{},
				curve: Curve{
					a: 6.71881241016441,
					b: 0.0486498280492762,
					c: 2.0417306475862214,
					d: 15.899403720950454,
				},
				minimumMemory:        16,
				maximumExecutionTime: 3000,
			},
		},
		networkPenalty:   50,
		memoryPenalty:    1000,
		thresholdPercent: 5,
	}

	_, schedulingMap, err := scheduler.Schedule()

	s.Nil(err)
	s.Equal("n2", schedulingMap["p1"])
	s.Equal("n2", schedulingMap["p2"])
	s.Equal("n2", schedulingMap["p3"])
	s.Equal("n2", schedulingMap["p4"])
	s.Equal("n2", schedulingMap["p5"])
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldFindGoodSchedulingInRealisticScenario() {
	var scheduler AdvancedScheduler
	nodes := []*NodeData{
		{
			name:                    "n1",
			allocatableCpu:          80000,
			memory:                  2560,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		{
			name:                    "n2",
			allocatableCpu:          80000,
			memory:                  2560,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		{
			name:                    "n3",
			allocatableCpu:          80000,
			memory:                  2560,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		{
			name:                    "n4",
			allocatableCpu:          80000,
			memory:                  2560,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
		{
			name:                    "n5",
			allocatableCpu:          80000,
			memory:                  2560,
			initialNumberOfPodSlots: 2,
			podSlotScalingFactor:    2,
			resourceLimit:           0.1,
		},
	}
	curve := Curve{
		a: 6.71881241016441,
		b: 0.0486498280492762,
		c: 2.0417306475862214,
		d: 15.899403720950454,
	}
	pods := []*PodData{
		{
			name:                 "p1_1",
			dataSourceNodes:      []string{"n1"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p1_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p2_1",
			dataSourceNodes:      []string{"n1"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p2_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p3_1",
			dataSourceNodes:      []string{"n1"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p3_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p4_1",
			dataSourceNodes:      []string{"n1"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p4_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p5_1",
			dataSourceNodes:      []string{"n1"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p5_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p6_1",
			dataSourceNodes:      []string{"n2"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p6_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p7_1",
			dataSourceNodes:      []string{"n2"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p7_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p8_1",
			dataSourceNodes:      []string{"n2"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p8_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p9_1",
			dataSourceNodes:      []string{"n2"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p9_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p10_1",
			dataSourceNodes:      []string{"n2"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p10_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p11_1",
			dataSourceNodes:      []string{"n3"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p11_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p12_1",
			dataSourceNodes:      []string{"n3"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p12_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p13_1",
			dataSourceNodes:      []string{"n3"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p13_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p14_1",
			dataSourceNodes:      []string{"n3"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p14_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p15_1",
			dataSourceNodes:      []string{"n3"},
			receivesDataFrom:     []string{},
			sendsDataTo:          []string{"p15_2"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p1_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p1_1"},
			sendsDataTo:          []string{"p1_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p2_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p2_1"},
			sendsDataTo:          []string{"p2_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p3_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p3_1"},
			sendsDataTo:          []string{"p3_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p4_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p4_1"},
			sendsDataTo:          []string{"p4_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p5_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p5_1"},
			sendsDataTo:          []string{"p5_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p6_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p6_1"},
			sendsDataTo:          []string{"p6_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p7_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p7_1"},
			sendsDataTo:          []string{"p7_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p8_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p8_1"},
			sendsDataTo:          []string{"p8_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p9_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p9_1"},
			sendsDataTo:          []string{"p9_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p10_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p10_1"},
			sendsDataTo:          []string{"p10_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p11_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p11_1"},
			sendsDataTo:          []string{"p11_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p12_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p12_1"},
			sendsDataTo:          []string{"p12_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p13_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p13_1"},
			sendsDataTo:          []string{"p13_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p14_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p14_1"},
			sendsDataTo:          []string{"p14_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p15_2",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p15_1"},
			sendsDataTo:          []string{"p15_3"},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p1_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p1_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p2_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p2_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p3_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p3_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p4_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p4_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p5_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p5_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p6_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p6_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p7_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p7_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p8_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p8_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p9_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p9_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p10_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p10_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p11_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p11_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p12_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p12_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p13_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p13_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p14_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p14_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
		{
			name:                 "p15_3",
			dataSourceNodes:      []string{},
			receivesDataFrom:     []string{"p15_2"},
			sendsDataTo:          []string{},
			curve:                curve,
			minimumMemory:        16,
			maximumExecutionTime: 200,
		},
	}
	scheduler = AdvancedScheduler{
		nodes:            nodes,
		pods:             pods,
		networkPenalty:   500,
		memoryPenalty:    100,
		thresholdPercent: 10,
	}

	schedulingChanged, scheduledMap, err := scheduler.Schedule()

	s.Nil(err)
	s.True(schedulingChanged)
	penalty, err := CalculatePenalty(getSystemStateFromSchedulingMap(scheduler.nodes, scheduler.pods, scheduledMap), scheduler.networkPenalty, scheduler.memoryPenalty)
	s.Nil(err)
	s.Equal(0.0, penalty)
}
