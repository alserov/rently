package broker

import (
	"github.com/IBM/sarama"
)

func NewProducer(brokerAddr string) sarama.SyncProducer {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer([]string{brokerAddr}, cfg)
	if err != nil {
		panic("failed to init producer: " + err.Error())
	}

	return p
}
