package testbroker

import "errors"

type Verification struct {
	broker *Broker
}

func (v *Verification) SendMessage(message []byte, uri string) error {
	v.broker.channel[uri] = message

	_, ok := v.broker.channel[uri]
	if !ok {
		return errors.New("error sending the message")
	}
	return nil
}
