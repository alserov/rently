package metrics

import (
	"context"
	"encoding/json"
	"github.com/alserov/rently/carsharing/internal/config"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"time"
)

type Metrics interface {
	IncreaseActiveRentsAmount() // active rents increment
	DecreaseActiveRentsAmount() // active rents decrement

	IncreaseRentsCancel() // how many rents were canceled

	ResponseTime(duration time.Duration, path string)

	NotifyBrandDemand(brand string) // what brand is the most popular
}

func NewMetrics(ch *amqp.Channel, topics config.Metrics, log *slog.Logger) Metrics {
	return &metrics{
		log:    log,
		ch:     ch,
		topics: topics,
	}
}

type metrics struct {
	log *slog.Logger

	ch *amqp.Channel

	topics config.Metrics
}

type responseTime struct {
	Duration time.Duration `json:"duration"`
	Path     string        `json:"path"`
}

func (m *metrics) ResponseTime(duration time.Duration, path string) {
	op := "metrics.ResponseTime"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	b, err := json.Marshal(responseTime{
		Duration: duration,
		Path:     path,
	})
	if err != nil {
		m.log.Error("failed to marshal message", slog.String("error", err.Error()), slog.String("op", op))
	}

	if err = m.ch.PublishWithContext(ctx, "", m.topics.ResponseTime, false, true, amqp.Publishing{
		MessageId: path,
		Body:      b,
	}); err != nil {
		m.log.Error("failed to send message", slog.String("error", err.Error()), slog.String("op", op))
	}
}

func (m *metrics) IncreaseRentsCancel() {
	op := "metrics.IncreaseRentsCancel"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := m.ch.PublishWithContext(ctx, "", m.topics.IncreaseRentsCancel, false, true, amqp.Publishing{}); err != nil {
		m.log.Error("failed to send message: ", err.Error(), slog.String("op", op))
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

func (m *metrics) DecreaseActiveRentsAmount() {
	op := "metrics.DecreaseActiveRentsAmount"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := m.ch.PublishWithContext(ctx, "", m.topics.DecreaseActiveRentsAmount, false, true, amqp.Publishing{}); err != nil {
		m.log.Error("failed to send message: ", err.Error(), slog.String("op", op))
	}
}

func (m *metrics) IncreaseActiveRentsAmount() {
	op := "metrics.IncreaseActiveRentsAmount"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := m.ch.PublishWithContext(ctx, "", m.topics.IncreaseActiveRentsAmount, false, true, amqp.Publishing{}); err != nil {
		m.log.Error("failed to send message: ", err.Error(), slog.String("op", op))
	}
}
