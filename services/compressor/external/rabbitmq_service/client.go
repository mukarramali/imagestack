package rabbitmq_service

import (
	"compressor/shared"

	amqp "github.com/rabbitmq/amqp091-go"
)

func getClient() *amqp.Channel {
	conn, err := amqp.Dial(shared.RABBITMQ_URL)
	shared.FailOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")
	return ch
}
