package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Consume(q string) <-chan amqp.Delivery
}

func NewConsumer(ch *amqp.Channel) Consumer {
	return &consumer{
		ch: ch,
	}
}

type consumer struct {
	ch *amqp.Channel
}

func (c *consumer) Consume(q string) <-chan amqp.Delivery {
	queue, err := c.ch.QueueDeclare(q, false, false, false, false, nil)
	if err != nil {
		panic("failed to declare a queue: " + err.Error())
	}

	msgs, err := c.ch.Consume(
		queue.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil)
	if err != nil {
		if err != nil {
			panic("failed to start consumer: " + err.Error())
		}
	}

	return msgs
}
