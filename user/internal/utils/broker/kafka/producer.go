package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/alserov/rently/user/internal/utils/broker"
)

func NewProducer() broker.Producer {
	return &producer{}
}

type producer struct {
	p sarama.SyncProducer
}

func (p producer) Produce(ctx context.Context, value any, id string, q string) error {
	b, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	_, _, err = p.p.SendMessage(&sarama.ProducerMessage{
		Topic: q,
		Key:   sarama.StringEncoder(id),
		Value: sarama.StringEncoder(b),
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
