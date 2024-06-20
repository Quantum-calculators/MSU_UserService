package message_broker

type Broker interface {
	Message() Verification
}
