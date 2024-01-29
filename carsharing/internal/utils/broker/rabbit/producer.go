package rabbit

import amqp "github.com/rabbitmq/amqp091-go"

func NewProducer(addr string) (*amqp.Channel, *amqp.Connection) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	return ch, conn
}
