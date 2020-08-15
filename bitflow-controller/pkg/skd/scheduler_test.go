package skd

import (
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

	scheduledMap, err := scheduler.Schedule()

	s.Nil(err)
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

func (s *SkdTestSuite) Test_AdvancedScheduler_shouldMapPodsCorrectly() {
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
				name:             "p3",
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
				name:             "p4",
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
				name:             "p8",
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
				name:             "p9",
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
	}

	scheduledMap, err := scheduler.Schedule()

	s.Nil(err)
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
