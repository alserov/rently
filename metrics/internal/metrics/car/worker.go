package car

import (
	"context"
	"github.com/alserov/rently/metrics/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
)

func NewWorker(brokerAddr string, gaugeTopics *GaugeTopics, log *slog.Logger) (metrics.Worker, []prometheus.Collector) {
	gauge := NewGaugeMetric(brokerAddr, gaugeTopics, log)

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
