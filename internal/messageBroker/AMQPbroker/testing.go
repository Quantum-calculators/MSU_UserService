package messagebroker

import (
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	Host     = "localhost"
	Port     = "5672"
	User     = "guest"
	Password = "guest"
)

func TestBroker() (*amqp.Channel, error) {

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", User, Password, Host, Port))
	if err != nil {
		return nil, errors.New("failed to connect to AMPQ")
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.New("failed to connect to AMPQ")
	}
	return ch, nil
}
