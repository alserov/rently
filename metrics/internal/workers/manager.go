package workers

import (
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
)

type WorkerManager interface {
	Add(worker Worker)
	Metrics() []prometheus.Collector
}

func NewWorkerManager(log *slog.Logger) WorkerManager {
	return &workerManager{}
}

type workerManager struct {
	log *slog.Logger

	workers []Worker
}

func (w workerManager) Metrics() []prometheus.Collector {
	var m []prometheus.Collector
	for _, wkr := range w.workers {
		m = append(m, wkr.Metrics()...)
	}
	return m
}

func (w workerManager) Add(worker Worker) {
	worker.MustStart()
	w.workers = append(w.workers, worker)
}
