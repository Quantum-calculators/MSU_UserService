package messagebroker

import (
	broker "github.com/Quantum-calculators/MSU_UserService/internal/messageBroker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Broker struct {
	channel      *amqp.Channel
	verification *Verification
}

func New(channel *amqp.Channel) *Broker {
	return &Broker{
		channel: channel,
	}
}

func (b *Broker) Message() broker.Verification {
	if b.verification != nil {
		return b.verification
	}

	b.verification = &Verification{
		broker: b,
	}

	return b.verification
}
