package bitflow

import (
	"sync"
	"time"

	"github.com/antongulenko/go-onlinestats"
)

type ReconcileStatistics struct {
	data ReconcileStatisticData
	lock sync.RWMutex
}

type ReconcileStatisticData struct {
	PodUpdateRoutines int
	PodUpdateDuration onlinestats.Running
	PodSpawnRoutines  int
	PodSpawnDuration  onlinestats.Running
	RestartedPods     int
	Errors            int
}

func (stat *ReconcileStatistics) GetData() ReconcileStatisticData {
	if stat == nil {
		return ReconcileStatisticData{}
	}
	stat.lock.RLock()
	defer stat.lock.RUnlock()
	return stat.data
}

func (stat *ReconcileStatistics) update(f func()) {
	if stat == nil {
		return
	}
	stat.lock.Lock()
	defer stat.lock.Unlock()
	f()
}

func (stat *ReconcileStatistics) PodsUpdated(duration time.Duration) {
	stat.update(func() {
		stat.data.PodUpdateRoutines++
		stat.data.PodUpdateDuration.Push(float64(duration.Nanoseconds()))
	})
}

func (stat *ReconcileStatistics) PodsSpawned(duration time.Duration) {
	stat.update(func() {
		stat.data.PodSpawnRoutines++
		stat.data.PodSpawnDuration.Push(float64(duration.Nanoseconds()))
	})
}

func (stat *ReconcileStatistics) PodRespawned() {
	stat.update(func() {
		stat.data.RestartedPods++
	})
}

func (stat *ReconcileStatistics) ErrorOccurred() {
	stat.update(func() {
		stat.data.Errors++
	})
}
