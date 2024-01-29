package broker

import amqp "github.com/rabbitmq/amqp091-go"

type ConsumerParams struct {
	Queue string
}

func NewConsumer(brokerAddr string, cp ConsumerParams) <-chan amqp.Delivery {
	conn, ch := dial(brokerAddr)
	defer conn.Close()
	defer ch.Close()

	msgs, err := ch.Consume(
		cp.Queue,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil)
	if err != nil {
		panicOnError("failed to start consumer", err)
	}

	return msgs
}
