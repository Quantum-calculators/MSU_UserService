package message_broker

type Verification interface {
	SendMessage([]byte, string) error
}
