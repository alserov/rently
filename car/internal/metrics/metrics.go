package metrics

import (
	"github.com/IBM/sarama"
	"github.com/alserov/rently/car/internal/utils/broker"
	"log/slog"
)

type Metrics interface {
	IncreaseActiveRentsAmount() // active rents increment
	DecreaseActiveRentsAmount() // active rents decrement

	IncreaseRentsCancel() // how many rents were canceled

	NotifyBrandDemand(brand string) // what brand is the most popular
}

func NewMetrics(producer sarama.SyncProducer, topics broker.MetricTopics, log *slog.Logger) Metrics {
	return &metrics{
		log:    log,
		p:      producer,
		topics: topics,
	}
}

type metrics struct {
	log *slog.Logger

	p sarama.SyncProducer

	topics broker.MetricTopics
}

func (m *metrics) IncreaseRentsCancel() {
	op := "metrics.IncreaseRentsCancel"
	_, _, err := m.p.SendMessage(&sarama.ProducerMessage{
		Topic: m.topics.IncreaseRentsCancel,
	})
	if err != nil {
		m.log.Error("failed to send message: ", err.Error(), slog.String("op", op))
	}
}

func (m *metrics) NotifyBrandDemand(brand string) {
	op := "metrics.NotifyBrandDemand"
	_, _, err := m.p.SendMessage(&sarama.ProducerMessage{
		Topic: m.topics.NotifyBrandDemand,
		Value: sarama.StringEncoder(brand),
	})
	m.log.Error("failed to send message: ", err.Error(), slog.String("op", op))
}

func (m *metrics) DecreaseActiveRentsAmount() {
	op := "metrics.DecreaseActiveRentsAmount"
	_, _, err := m.p.SendMessage(&sarama.ProducerMessage{
		Topic: m.topics.DecreaseActiveRentsAmount,
	})
	m.log.Error("failed to send message: ", err.Error(), slog.String("op", op))
}

func (m *metrics) IncreaseActiveRentsAmount() {
	op := "metrics.IncreaseActiveRentsAmount"
	_, _, err := m.p.SendMessage(&sarama.ProducerMessage{
		Topic: m.topics.IncreaseActiveRentsAmount,
	})
	m.log.Error("failed to send message: ", err.Error(), slog.String("op", op))
}
