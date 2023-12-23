package broker

import "github.com/IBM/sarama"

func NewConsumer(addr string) sarama.Consumer {
	c, err := sarama.NewConsumer([]string{addr}, &sarama.Config{})
	if err != nil {
		panic("failed to init consumer: " + err.Error())
	}

	return c
}
