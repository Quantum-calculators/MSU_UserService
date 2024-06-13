package messagebroker

import (
	"context"
	"time"

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

func (b *Broker) SendVerificationMessage(message []byte, uri string) error {
	q, err := b.ch.QueueDeclare(
		"EmailService", // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)

	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	header := amqp.Table{}
	header[uri] = uri

	err = b.ch.PublishWithContext(ctx,
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
