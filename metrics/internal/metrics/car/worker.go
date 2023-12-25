package car

import (
	"context"
	"github.com/alserov/rently/metrics/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

func NewCarWorker(brokerAddr string, gaugeTopics *GaugeTopics) (metrics.Worker, []prometheus.Collector) {
	gauge := NewGaugeMetric(brokerAddr, gaugeTopics)

	return &worker{
		gauge:      gauge,
		brokerAddr: brokerAddr,
	}, []prometheus.Collector{gauge.Get()}
}

type worker struct {
	gauge metrics.Metric

	brokerAddr string
}

func (w *worker) MustStart(ctx context.Context) {
	go func() {
		w.gauge.Run(ctx)
	}()

	select {
	case <-ctx.Done():
	}
}
