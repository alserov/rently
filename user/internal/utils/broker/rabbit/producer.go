package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alserov/rently/user/internal/utils/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewProducer(ch *amqp.Channel) broker.Producer {
	return &producer{
		ch: ch,
	}
}

type producer struct {
	ch *amqp.Channel
}

func (p producer) Produce(ctx context.Context, value any, id string, q string) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = p.ch.PublishWithContext(ctx, "", q, false, false, amqp.Publishing{
		CorrelationId: id,
		Body:          b,
	})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}
