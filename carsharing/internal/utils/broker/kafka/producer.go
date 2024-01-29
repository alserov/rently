package kafka

import (
	"github.com/IBM/sarama"
)

func NewSyncProducer(brokerAddr string) sarama.SyncProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer([]string{brokerAddr}, cfg)
	if err != nil {
		panic("failed to init producer: " + err.Error())
	}

	return p
}

func NewAsyncProducer(brokerAddr string) sarama.AsyncProducer {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.Return.Errors = true
	producerConfig.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer([]string{brokerAddr}, producerConfig)
	if err != nil {
		panic("failed to init async producer: " + err.Error())
	}

	return producer
}
