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
	var scheduler Scheduler
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
				receivesDataFrom: []string{"p2, p3"},
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
		},
		networkPenalty:   0,
		thresholdPercent: 10,
	}

	schedulingChanged, scheduledMap, err := scheduler.Schedule()

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
	var scheduler Scheduler
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
		thresholdPercent: 10,
	}

	schedulingChanged, scheduledMap, err := scheduler.Schedule()

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
	var scheduler Scheduler
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

	schedulingChanged1, scheduledMap, err1 := scheduler.Schedule()

	s.Nil(err1)
	s.True(schedulingChanged1)

	scheduler = AdvancedScheduler{
		nodes:              nodes,
		pods:               pods,
		thresholdPercent:   10,
		previousScheduling: scheduledMap,
	}

	schedulingChanged2, _, err2 := scheduler.Schedule()

	s.Nil(err2)
	s.False(schedulingChanged2)
}

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldRecognizeSchedulingHasChangedWithoutThreshold() {
	var scheduler Scheduler
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

	schedulingChanged1, scheduledMap, err1 := scheduler.Schedule()

	s.Nil(err1)
	s.True(schedulingChanged1)

	scheduler = AdvancedScheduler{
		nodes:              nodes,
		pods:               pods,
		previousScheduling: scheduledMap,
	}

	schedulingChanged2, _, err2 := scheduler.Schedule()

	s.Nil(err2)
	s.True(schedulingChanged2)
}

func (s *SkdTestSuite) testNewDistributionPenaltyLowerConsideringThreshold(previousPenalty float64, newPenalty float64, thresholdPercent float64, expectedOutcome bool) {
	s.SubTest(fmt.Sprintf("previous%f:new%f:threshold%f->%v", previousPenalty, newPenalty, thresholdPercent, expectedOutcome), func() {
		actualOutcome := NewDistributionPenaltyLowerConsideringThreshold(previousPenalty, newPenalty, thresholdPercent)
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
