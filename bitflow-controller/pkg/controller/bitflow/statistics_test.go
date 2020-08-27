package bitflow

import (
	"sync"
	"testing"
	"time"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/stretchr/testify/suite"
)

const (
	loopTimeNanos = 102334
	loopTime      = time.Duration(loopTimeNanos)
)

type StatisticsTestSuite struct {
	common.AbstractTestSuite
}

func TestStatistics(t *testing.T) {
	suite.Run(t, new(StatisticsTestSuite))
}

func (s *StatisticsTestSuite) writeStats(stat *ReconcileStatistics, numberOfLogs ...int) {
	var wg sync.WaitGroup
	wg.Add(len(numberOfLogs))
	for _, numLogs := range numberOfLogs {
		go func(numLogs int) {
			defer wg.Done()
			for i := 0; i < numLogs; i++ {
				stat.PodsUpdated(loopTime)
				stat.PodsSpawned(loopTime)
				stat.PodRespawned()
				stat.ErrorOccurred()
			}
		}(numLogs)
	}
	wg.Wait()
}

func (s *StatisticsTestSuite) TestConcurrentStatistics() {
	stat := new(ReconcileStatistics)
	s.writeStats(stat, 5, 3, 9, 3, 2)

	data := stat.GetData()
	s.Equal(22, data.PodUpdateRoutines)
	s.Equal(22, data.PodSpawnRoutines)
	s.Equal(22, data.RestartedPods)
	s.Equal(22, data.Errors)
	s.InEpsilon(loopTimeNanos, data.PodUpdateDuration.Mean(), 0.00001)
	s.InEpsilon(loopTimeNanos, data.PodSpawnDuration.Mean(), 0.00001)
}

func (s *StatisticsTestSuite) TestNilStatistics() {
	var stat *ReconcileStatistics // nil
	s.writeStats(stat, 2, 4, 6)
	s.Zero(stat.GetData().PodUpdateRoutines, "nil-statistics should remain empty")
}
