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
	PodReconcileLoops                 int
	PodReconcileLoopDuration          onlinestats.Running
	NodeResourceReconcileLoops        int
	NodeResourceReconcileLoopDuration onlinestats.Running
	RestartedPods                     int
	Errors                            int
}

func (stat *ReconcileStatistics) GetData() ReconcileStatisticData {
	if stat == nil {
		return ReconcileStatisticData{}
	}
	stat.lock.RLock()
	defer stat.lock.RUnlock()
	return stat.data
}

func (stat *ReconcileStatistics) PodsReconciled(duration time.Duration) {
	if stat == nil {
		return
	}
	stat.lock.Lock()
	defer stat.lock.Unlock()
	stat.data.PodReconcileLoops++
	stat.data.PodReconcileLoopDuration.Push(float64(duration.Nanoseconds()))
}

func (stat *ReconcileStatistics) NodeResourcesReconciled(duration time.Duration) {
	if stat == nil {
		return
	}
	stat.lock.Lock()
	defer stat.lock.Unlock()
	stat.data.NodeResourceReconcileLoops++
	stat.data.NodeResourceReconcileLoopDuration.Push(float64(duration.Nanoseconds()))
}

func (stat *ReconcileStatistics) PodRespawned() {
	if stat == nil {
		return
	}
	stat.lock.Lock()
	defer stat.lock.Unlock()
	stat.data.RestartedPods++
}

func (stat *ReconcileStatistics) ErrorOccurred() {
	if stat == nil {
		return
	}
	stat.lock.Lock()
	defer stat.lock.Unlock()
	stat.data.Errors++
}
