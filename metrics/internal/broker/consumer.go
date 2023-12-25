package broker

import (
	"github.com/IBM/sarama"
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
