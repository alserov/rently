package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/alserov/rently/user/internal/utils/broker"
)

func NewProducer(topic string) broker.Producer {
	return &producer{
		topic: topic,
	}
}

type producer struct {
	p sarama.SyncProducer

	topic string
}

func (p producer) Produce(values ...any) error {
	for _, v := range values {
		b, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}

		_, _, err = p.p.SendMessage(&sarama.ProducerMessage{
			Topic: p.topic,
			Value: sarama.StringEncoder(b),
		})
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	return nil
}
