package metrics

import (
	"github.com/IBM/sarama"
	"log/slog"
)

type Metrics interface {
	IncreaseActiveRentsAmount()
	DecreaseActiveRentsAmount()
}

type MetricTopics struct {
	DecreaseActiveRentsAmount string
	IncreaseActiveRentsAmount string
}

func NewMetrics(producer sarama.SyncProducer, topics MetricTopics, log *slog.Logger) Metrics {
	return &metrics{
		p:      producer,
		topics: topics,
	}
}

type metrics struct {
	log *slog.Logger

	p sarama.SyncProducer

	topics MetricTopics
}

func (m *metrics) DecreaseActiveRentsAmount() {
	_, _, err := m.p.SendMessage(&sarama.ProducerMessage{
		Topic: m.topics.DecreaseActiveRentsAmount,
	})
	if err != nil {
		m.log.Error("failed to send message: ", err.Error())
	}
}

func (m *metrics) IncreaseActiveRentsAmount() {
	_, _, err := m.p.SendMessage(&sarama.ProducerMessage{
		Topic: m.topics.IncreaseActiveRentsAmount,
	})
	if err != nil {
		m.log.Error("failed to send message: ", err.Error())
	}
}
