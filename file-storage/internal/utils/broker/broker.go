package broker

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"log/slog"
)

type Broker struct {
	Addr   string
	Topics Topics
}

type Topics struct {
	SaveImages   string
	DeleteImages string
}

func Consume[MessageType any](workersAmount int, addr string, topic string, consumerConfig *sarama.Config, log *slog.Logger) chan MessageType {
	master, err := sarama.NewConsumer([]string{addr}, consumerConfig)
	if err != nil {
		panic("failed to init consumer: " + err.Error())
	}

	partitions, err := master.Partitions(topic)
	if err != nil {
		panic("failed to getLinks partitions: " + err.Error())
	}

	messages := make(chan MessageType, workersAmount)

	for _, p := range partitions {
		go func(p int32) {
			consumer, err := master.ConsumePartition(topic, p, sarama.OffsetNewest)
			if err != nil {
				panic("failed to consume partition: " + err.Error())
			}
			for m := range consumer.Messages() {
				var msg MessageType
				if err = json.Unmarshal(m.Value, &msg); err != nil {
					log.Error("failed to unmarshal fileMessage", slog.String("error", err.Error()))
				}
				messages <- msg
			}
		}(p)
	}

	return messages
}
