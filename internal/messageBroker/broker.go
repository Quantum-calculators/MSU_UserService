package messagebroker

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	ch *amqp.Channel
}

func NewAMQL(AMQPaddr string) *Broker {
	conn, err := amqp.Dial(AMQPaddr)
	if err != nil {
		panic("Failed to connect to RabbitMQ")
	}
	ch, err := conn.Channel()
	if err != nil {
		panic("Failed to connect to RabbitMQ")
	}
	return &Broker{
		ch: ch,
	}
}

func (b *Broker) CloseAMQL() {
	if err := b.ch.Close(); err != nil {
		panic("Failed to close connection to RabbitMQ")
	}
}
