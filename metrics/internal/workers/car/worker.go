package car

import (
	"context"
	"github.com/alserov/rently/metrics/internal/workers"
)

func NewRentWorker(brokerAddr string, gaugeTopics *GaugeTopics) workers.Worker {
	return &worker{
		gaugeTopics: gaugeTopics,
		gauge:       NewGaugeMetric(),
		brokerAddr:  brokerAddr,
	}
}

type worker struct {
	gauge       *GaugeMetric
	gaugeTopics *GaugeTopics

	brokerAddr string
}

func (w *worker) MustStart(ctx context.Context) {
	go func() {
		w.gauge.Run(w.gaugeTopics)
	}()

	select {
	case <-ctx.Done():

	}
}
