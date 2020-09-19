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

//func (s *SkdTestSuite) Test_AdvancedScheduler_shouldScheduleRealisticScenarioWithNetworkPenalty() {
//	var scheduler Scheduler
//	nodes := []*NodeData{
//		{
//			name:                    "n1",
//			allocatableCpu:          4000,
//			memory:                  64,
//			initialNumberOfPodSlots: 2,
//			podSlotScalingFactor:    2,
//			resourceLimit:           0.1,
//		},
//		{
//			name:                    "n2",
//			allocatableCpu:          4000,
//			memory:                  64,
//			initialNumberOfPodSlots: 2,
//			podSlotScalingFactor:    2,
//			resourceLimit:           0.1,
//		},
//		{
//			name:                    "n3",
//			allocatableCpu:          4000,
//			memory:                  64,
//			initialNumberOfPodSlots: 2,
//			podSlotScalingFactor:    2,
//			resourceLimit:           0.1,
//		},
//		{
//			name:                    "n4",
//			allocatableCpu:          4000,
//			memory:                  64,
//			initialNumberOfPodSlots: 2,
//			podSlotScalingFactor:    2,
//			resourceLimit:           0.1,
//		},
//		{
//			name:                    "n5",
//			allocatableCpu:          4000,
//			memory:                  64,
//			initialNumberOfPodSlots: 2,
//			podSlotScalingFactor:    2,
//			resourceLimit:           0.1,
//		},
//	}
//	curve := Curve{
//		a: 6.71881241016441,
//		b: 0.0486498280492762,
//		c: 2.0417306475862214,
//		d: 15.899403720950454,
//	}
//	pods := []*PodData{
//		{
//			name:             "p1_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p2_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p3_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p4_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p5_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p6_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p7_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p8_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p9_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p10_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p11_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p12_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p13_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p14_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p15_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p16_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p17_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p18_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p19_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p20_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p21_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p22_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p23_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p24_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p25_1",
//			receivesDataFrom: []string{},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p1_2",
//			receivesDataFrom: []string{"p1_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p2_2",
//			receivesDataFrom: []string{"p2_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p3_2",
//			receivesDataFrom: []string{"p3_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p4_2",
//			receivesDataFrom: []string{"p4_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p5_2",
//			receivesDataFrom: []string{"p5_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p6_2",
//			receivesDataFrom: []string{"p6_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p7_2",
//			receivesDataFrom: []string{"p7_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p8_2",
//			receivesDataFrom: []string{"p8_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p9_2",
//			receivesDataFrom: []string{"p9_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p10_2",
//			receivesDataFrom: []string{"p10_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p11_2",
//			receivesDataFrom: []string{"p11_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p12_2",
//			receivesDataFrom: []string{"p12_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p13_2",
//			receivesDataFrom: []string{"p13_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p14_2",
//			receivesDataFrom: []string{"p14_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p15_2",
//			receivesDataFrom: []string{"p15_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p16_2",
//			receivesDataFrom: []string{"p16_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p17_2",
//			receivesDataFrom: []string{"p17_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p18_2",
//			receivesDataFrom: []string{"p18_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p19_2",
//			receivesDataFrom: []string{"p19_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p20_2",
//			receivesDataFrom: []string{"p20_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p21_2",
//			receivesDataFrom: []string{"p21_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p22_2",
//			receivesDataFrom: []string{"p22_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p23_2",
//			receivesDataFrom: []string{"p23_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p24_2",
//			receivesDataFrom: []string{"p24_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p25_2",
//			receivesDataFrom: []string{"p25_1"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p1_3",
//			receivesDataFrom: []string{"p1_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p2_3",
//			receivesDataFrom: []string{"p2_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p3_3",
//			receivesDataFrom: []string{"p3_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p4_3",
//			receivesDataFrom: []string{"p4_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p5_3",
//			receivesDataFrom: []string{"p5_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p6_3",
//			receivesDataFrom: []string{"p6_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p7_3",
//			receivesDataFrom: []string{"p7_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p8_3",
//			receivesDataFrom: []string{"p8_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p9_3",
//			receivesDataFrom: []string{"p9_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p10_3",
//			receivesDataFrom: []string{"p10_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p11_3",
//			receivesDataFrom: []string{"p11_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p12_3",
//			receivesDataFrom: []string{"p12_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p13_3",
//			receivesDataFrom: []string{"p13_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p14_3",
//			receivesDataFrom: []string{"p14_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p15_3",
//			receivesDataFrom: []string{"p15_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p16_3",
//			receivesDataFrom: []string{"p16_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p17_3",
//			receivesDataFrom: []string{"p17_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p18_3",
//			receivesDataFrom: []string{"p18_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p19_3",
//			receivesDataFrom: []string{"p19_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p20_3",
//			receivesDataFrom: []string{"p20_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p21_3",
//			receivesDataFrom: []string{"p21_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p22_3",
//			receivesDataFrom: []string{"p22_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p23_3",
//			receivesDataFrom: []string{"p23_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p24_3",
//			receivesDataFrom: []string{"p24_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//		{
//			name:             "p25_3",
//			receivesDataFrom: []string{"p25_2"},
//			curve:            curve,
//			minimumMemory:    16,
//		},
//	}
//	scheduler = AdvancedScheduler{
//		nodes:            nodes,
//		pods:             pods,
//		networkPenalty:   1_000,
//		thresholdPercent: 10,
//	}
//
//	schedulingChanged, scheduledMap, err := scheduler.Schedule()
//
//	s.Nil(err)
//	s.True(schedulingChanged)
//	s.Equal(scheduledMap["p2"], scheduledMap["p10"])
//	s.Equal(scheduledMap["p7"], scheduledMap["p1"])
//	s.Equal(scheduledMap["p8"], scheduledMap["p1"])
//	s.Equal(scheduledMap["p9"], scheduledMap["p3"])
//	s.Equal(scheduledMap["p9"], scheduledMap["p4"])
//	s.Equal(scheduledMap["p9"], scheduledMap["p6"])
//	s.Equal(scheduledMap["p9"], scheduledMap["p7"])
//	s.Equal(scheduledMap["p9"], scheduledMap["p8"])
//}

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

	_, _, err := scheduler.Schedule()

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

	_, _, err := scheduler.Schedule()

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

	_, _, err := scheduler.Schedule()

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

	_, _, err := scheduler.Schedule()

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

	_, _, err := scheduler.Schedule()

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

	_, _, err := scheduler.Schedule()

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

	_, _, err := scheduler.Schedule()

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

	_, _, err := scheduler.Schedule()

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

	_, _, err := scheduler.Schedule()

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
