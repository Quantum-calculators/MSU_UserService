package messagebroker

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Verification struct {
	broker *Broker
}

func (v *Verification) SendMessage(message []byte, uri string) error {
	q, err := v.broker.channel.QueueDeclare(
		uri,   // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	header := amqp.Table{}
	header[uri] = uri
	err = v.broker.channel.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			Headers:      header,
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         message,
		})
	if err != nil {
		return err
	}
	return nil
}
