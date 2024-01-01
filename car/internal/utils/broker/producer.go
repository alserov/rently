package broker

import "github.com/IBM/sarama"

func NewProducer(brokerAddr string) sarama.SyncProducer {
	p, err := sarama.NewSyncProducer([]string{brokerAddr}, &sarama.Config{})
	if err != nil {
		panic("failed to init producer: " + err.Error())
	}

	return p
}
