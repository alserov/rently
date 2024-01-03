package car

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/alserov/rently/metrics/internal/metrics"
	"github.com/alserov/rently/metrics/internal/utils/broker"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
)

type GaugeTopics struct {
	DecreaseActiveRentsAmount string
	IncreaseActiveRentsAmount string
}

func NewGaugeMetric(brokerAddr string, topics *GaugeTopics, log *slog.Logger) metrics.Metric {
	return &gaugeMetric{
		log:                            log,
		brokerAddr:                     brokerAddr,
		increaseActiveRentsAmountTopic: topics.IncreaseActiveRentsAmount,
		decreaseActiveRentsAmountTopic: topics.DecreaseActiveRentsAmount,
		metric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "active_car_rents",
			Name:      "active_rents",
			Help:      "monitoring of cars that are being rented",
		}, []string{}),
	}
}

type gaugeMetric struct {
	log *slog.Logger

	brokerAddr string

	metric *prometheus.GaugeVec

	decreaseActiveRentsAmountTopic string
	increaseActiveRentsAmountTopic string
	brandDemandTopic               string
}

func (g *gaugeMetric) Get() prometheus.Collector {
	return g.metric
}

func (g *gaugeMetric) Run(ctx context.Context) {
	cfg := sarama.NewConfig()
	cfg.Consumer.Return.Errors = true

	master, err := sarama.NewConsumer([]string{g.brokerAddr}, cfg)
	if err != nil {
		panic("failed to init consumer: " + err.Error())
	}
	defer master.Close()

	go g.increaseActiveRents(master, g.increaseActiveRentsAmountTopic)
	go g.decreaseActiveRents(master, g.decreaseActiveRentsAmountTopic)
	go g.notifyBrandDemand(master, g.brandDemandTopic)

	select {
	case <-ctx.Done():

	}
}

func (g *gaugeMetric) increaseActiveRents(master sarama.Consumer, topic string) {
	msgs := broker.Subscribe(master, topic)
	for range msgs {
		g.metric.With(prometheus.Labels{}).Inc()
	}
}

func (g *gaugeMetric) decreaseActiveRents(master sarama.Consumer, topic string) {
	msgs := broker.Subscribe(master, topic)
	for range msgs {
		g.metric.With(prometheus.Labels{}).Dec()
	}
}

func (g *gaugeMetric) notifyBrandDemand(master sarama.Consumer, topic string) {
	msgs := broker.SubscribeWithValue[string](master, topic, g.log)
	for m := range msgs {

	}
}
