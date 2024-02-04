package workers

import (
	"github.com/prometheus/client_golang/prometheus"
)

type WorkerManager interface {
	Add(worker Worker)
	Metrics() []prometheus.Collector
}

type Worker interface {
	MustStart()
	Metrics() []prometheus.Collector
}

func NewWorkerManager() WorkerManager {
	return &workerManager{}
}

type workerManager struct {
	workers []Worker
}

func (w *workerManager) Metrics() []prometheus.Collector {
	var m []prometheus.Collector
	for _, wkr := range w.workers {
		m = append(m, wkr.Metrics()...)
	}
	return m
}

func (w *workerManager) Add(worker Worker) {
	w.workers = append(w.workers, worker)
	go worker.MustStart()
}
