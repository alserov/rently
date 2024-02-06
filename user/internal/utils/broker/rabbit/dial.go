package rabbit

import amqp "github.com/rabbitmq/amqp091-go"

func Dial(brokerAddr string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(brokerAddr)
	if err != nil {
		panic("failed to dial rabbit server: " + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		panic("failed to get rabbit channel: " + err.Error())
	}

	return conn, ch, nil
}
