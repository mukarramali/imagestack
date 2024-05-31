package rabbitmq_service

import (
	"context"
	"imagestack/lib/error_handler"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (rs *RabbitMqService) Publish(body string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rs.Client.PublishWithContext(ctx,
		"",            // exchange
		rs.Queue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	error_handler.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
