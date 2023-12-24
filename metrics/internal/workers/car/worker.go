package car

import (
	"context"
	"github.com/alserov/rently/metrics/internal/workers"
)

func NewRentWorker(brokerAddr string, gaugeTopics *GaugeTopics) workers.Worker {
	return &worker{
		gauge:      NewGaugeMetric(gaugeTopics),
		brokerAddr: brokerAddr,
	}
}

type worker struct {
	gauge *GaugeMetric

	brokerAddr string
}

func (w *worker) MustStart(ctx context.Context) {
	go func() {
		w.gauge.Run()
	}()

	select {
	case <-ctx.Done():

	}
}
