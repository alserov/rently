package rabbit

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func NewProducer(addr string) (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch, conn
}
