package car

import (
	"fmt"
	"github.com/alserov/rently/metrics/internal/log"
	"github.com/alserov/rently/metrics/internal/utils/broker"
	"github.com/prometheus/client_golang/prometheus"
)

type Gauge interface {
	prometheus.Collector
	GetMetrics() []prometheus.Collector
	IncreaseActiveRents(q string) error
	DecreaseActiveRents(q string) error
	NotifyBrandDemand(q string) error
}

func NewGaugeMetric() Gauge {
	return &gauge{
		log: log.GetLogger(),
		rentsAmountMetric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "current_rents_amount",
			Name:      "current_rents",
			Help:      "monitoring of cars that are being rented",
		}, []string{}),
	}
}

type gauge struct {
	prometheus.GaugeVec

	log log.Logger

	consumer broker.Consumer

	rentsAmountMetric *prometheus.GaugeVec
	brandDemandMetric *prometheus.CounterVec
}

func (g *gauge) GetMetrics() []prometheus.Collector {
	return []prometheus.Collector{g.rentsAmountMetric, g.brandDemandMetric}
}

func (g *gauge) IncreaseActiveRents(q string) error {
	msgs, err := g.consumer.Subscribe(q)
	if err != nil {
		return fmt.Errorf("failed to start IncreaseActiveRents metric: %v", err)
	}

	go func() {
		for range msgs {
			g.rentsAmountMetric.With(prometheus.Labels{}).Inc()
		}
	}()

	return nil
}

func (g *gauge) DecreaseActiveRents(q string) error {
	msgs, err := g.consumer.Subscribe(q)
	if err != nil {
		return fmt.Errorf("failed to start DecreaseActiveRents metric: %v", err)
	}

	go func() {
		for range msgs {
			g.rentsAmountMetric.With(prometheus.Labels{}).Dec()
		}
	}()

	return nil
}

func (g *gauge) NotifyBrandDemand(q string) error {
	msgs, err := g.consumer.Subscribe(q)
	if err != nil {
		return fmt.Errorf("failed to start DecreaseActiveRents metric: %v", err)
	}

	go func() {
		for range msgs {
			g.brandDemandMetric.With(prometheus.Labels{"": ""})
		}
	}()

	return nil
}
