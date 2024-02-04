package metrics

import (
	"context"
	"encoding/json"
	"github.com/alserov/rently/carsharing/internal/config"
	"github.com/alserov/rently/carsharing/internal/log"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"
)

type Metrics interface {
	ResponseTime(duration time.Duration, op string, t string)
}

func NewMetrics(ch *amqp.Channel, topics config.Metrics) Metrics {
	return &metrics{
		log:    log.GetLogger(),
		ch:     ch,
		topics: topics,
	}
}

type metrics struct {
	log log.Logger

	ch *amqp.Channel

	topics config.Metrics
}

type responseTimeData struct {
	RequestType string        `json:"requestType"`
	Operation   string        `json:"operation"`
	Duration    time.Duration `json:"duration"`
}

func (m *metrics) ResponseTime(duration time.Duration, op string, t string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	b, err := json.Marshal(responseTimeData{
		Duration:    duration,
		Operation:   op,
		RequestType: t,
	})
	if err != nil {
		m.log.Error("failed to marshal message", slog.String("error", err.Error()), slog.String("op", op))
	}

	if err = m.ch.PublishWithContext(ctx, "", m.topics.ResponseTime, false, false, amqp.Publishing{
		MessageId: time.Now().String(),
		Body:      b,
	}); err != nil {
		m.log.Error("failed to send message", slog.String("error", err.Error()), slog.String("op", op))
	}
}

func (m *metrics) NotifyBrandDemand(brand string) {
	if brand == "" {
		return
	}

	op := "metrics.NotifyBrandDemand"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := m.ch.PublishWithContext(ctx, "", m.topics.NotifyBrandDemand, false, true, amqp.Publishing{
		MessageId: brand,
	}); err != nil {
		m.log.Error("failed to send message: ", err.Error(), slog.String("op", op))
	}
}
