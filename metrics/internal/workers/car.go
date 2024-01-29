package workers

import (
	"github.com/alserov/rently/metrics/internal/config"
	"github.com/alserov/rently/metrics/internal/metrics/domains/car"
	"github.com/prometheus/client_golang/prometheus"
)

func NewCarWorker(brokerAddr string, topics config.CarsharingTopics) Worker {
	return &worker{
		gauge: car.NewGaugeMetric(),
	}
}

type worker struct {
	gauge car.Gauge

	metrics []prometheus.Collector

	topics config.CarsharingTopics
}

func (w *worker) Metrics() []prometheus.Collector {
	return []prometheus.Collector{w.gauge}
}

func (w *worker) MustStart() {
	// gauge
	panicOnErr(w.gauge.IncreaseActiveRents(w.topics.IncreaseActiveRentsAmount))
	panicOnErr(w.gauge.DecreaseActiveRents(w.topics.DecreaseActiveRentsAmount))
	panicOnErr(w.gauge.NotifyBrandDemand(w.topics.NotifyBrandDemand))

	select {}
}

func panicOnErr(err error) {
	if err != nil {
		panic("failed to start carsharing workers: " + err.Error())
	}
}
