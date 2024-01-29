package broker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func dial(brokerAddr string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(brokerAddr)
	panicOnError("failed to dial broker", err)

	ch, err := conn.Channel()
	panicOnError("failed to open channel", err)

	return conn, ch
}

func panicOnError(msg string, err error) {
	if err != nil {
		panic(msg + ": " + err.Error())
	}
}
