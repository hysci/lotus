package sectorstorage

import (
	"context"

	"github.com/filecoin-project/lotus/extern/sector-storage/storiface"
	"golang.org/x/xerrors"
)

func (m *Manager) WorkerStats() map[uint64]storiface.WorkerStats {
	m.sched.workersLk.RLock()
	defer m.sched.workersLk.RUnlock()

	out := map[uint64]storiface.WorkerStats{}

	for id, handle := range m.sched.workers {
		out[uint64(id)] = storiface.WorkerStats{
			Info:       handle.info,
			MemUsedMin: handle.active.memUsedMin,
			MemUsedMax: handle.active.memUsedMax,
			GpuUsed:    handle.active.gpuUsed,
			CpuUse:     handle.active.cpuUse,
		}
	}

	return out
}

func (m *Manager) WorkerJobs() map[uint64][]storiface.WorkerJob {
	m.sched.workersLk.RLock()
	defer m.sched.workersLk.RUnlock()

	out := map[uint64][]storiface.WorkerJob{}

	for id, handle := range m.sched.workers {
		out[uint64(id)] = handle.wt.Running()

		handle.wndLk.Lock()
		for wi, window := range handle.activeWindows {
			for _, request := range window.todo {
				out[uint64(id)] = append(out[uint64(id)], storiface.WorkerJob{
					ID:      0,
					Sector:  request.sector,
					Task:    request.taskType,
					RunWait: wi + 1,
					Start:   request.start,
				})
			}
		}
		handle.wndLk.Unlock()
	}

	return out
}

func (m *Manager) GetWorker(ctx context.Context) map[uint64]WorkerInfo {
	m.sched.workersLk.Lock()
	defer m.sched.workersLk.Unlock()

	out := map[uint64]WorkerInfo{}

	for id, handle := range m.sched.workers {
		info := handle.w.GetWorkerInfo(ctx)
		out[uint64(id)] = info
	}
	return out
}

func (m *Manager) SetWorkerParam(ctx context.Context, worker uint64, key string, value string) error {
	m.sched.workersLk.Lock()
	defer m.sched.workersLk.Unlock()

	w, exist := m.sched.workers[WorkerID(worker)]
	if !exist {
		return xerrors.Errorf("worker not found: %s", key)
	}
	return w.w.SetWorkerParams(ctx, key, value)
}
