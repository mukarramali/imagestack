package rabbitmq_service

import (
	"context"
	"imagestack/lib"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (rs *RabbitMqService) Publish(body string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := rs.client.PublishWithContext(ctx,
		"",            // exchange
		rs.queue.Name, // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	lib.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s\n", body)
}
