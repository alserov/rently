package broker

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Subscribe(q string) (<-chan amqp.Delivery, error)
}

func NewConsumer(addr string) Consumer {
	conn, err := amqp.Dial(addr)
	if err != nil {
		panic("failed to init consumer: " + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		panic("failed to open a channel: " + err.Error())
	}

	return &consumer{
		ch: ch,
	}
}

type consumer struct {
	ch *amqp.Channel
}

func (c consumer) Subscribe(q string) (<-chan amqp.Delivery, error) {
	msgs, err := c.ch.Consume(q, "", true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to queue: %v", err)
	}

	return msgs, nil
}
