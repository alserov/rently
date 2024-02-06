package rabbit

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

//type Rabbit struct {
//	Conn *amqp.Connection
//	Ch   *amqp.Channel
//
//	Consumer
//}
//
//func NewBroker(brokerAddr string) *Rabbit {
//	conn, ch, err := Dial(brokerAddr)
//	if err != nil {
//		panic(fmt.Errorf("failed to init new rabbit: %w", err))
//	}
//
//	return &Rabbit{
//		Conn: conn,
//		Ch:   ch,
//	}
//}

func Dial(addr string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial rabbit server: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get rabbit channel: " + err.Error())
	}

	return conn, ch, nil
}
