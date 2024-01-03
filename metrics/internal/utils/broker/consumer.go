package broker

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log/slog"
)

func NewConsumer(addr string) sarama.Consumer {
	c, err := sarama.NewConsumer([]string{addr}, &sarama.Config{})
	if err != nil {
		panic("failed to init consumer: " + err.Error())
	}

	return c
}

func Subscribe(master sarama.Consumer, topic string) chan struct{} {
	partitions, err := master.Partitions(topic)
	if err != nil {
		panic("failed to get partitions: " + err.Error())
	}

	messages := make(chan struct{}, 3)
	for _, p := range partitions {
		consumer, err := master.ConsumePartition(topic, p, sarama.OffsetNewest)
		if err != nil {
			panic("failed to init partition consumer: " + err.Error())
		}

		for range consumer.Messages() {
			messages <- struct{}{}
		}
	}
	return messages
}

func SubscribeWithValue[T comparable](master sarama.Consumer, topic string, log *slog.Logger) chan T {
	partitions, err := master.Partitions(topic)
	if err != nil {
		panic("failed to get partitions: " + err.Error())
	}

	messages := make(chan T, 3)
	for _, p := range partitions {
		consumer, err := master.ConsumePartition(topic, p, sarama.OffsetNewest)
		if err != nil {
			panic("failed to init partition consumer: " + err.Error())
		}

		for m := range consumer.Messages() {
			var msg T
			if err = json.Unmarshal(m.Value, &msg); err != nil {
				log.Error("failed to unmarshal message: ", err)
			}
			messages <- msg
		}
	}
	return messages
}
