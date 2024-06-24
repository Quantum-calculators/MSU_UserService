package testbroker

import (
	broker "github.com/Quantum-calculators/MSU_UserService/internal/messageBroker"
)

type Broker struct {
	channel      map[string][]byte
	verification *Verification
}

func New() *Broker {
	return &Broker{
		channel: make(map[string][]byte),
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
