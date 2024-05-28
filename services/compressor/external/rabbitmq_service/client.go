package rabbitmq_service

import (
	"compressor/shared"
	"imagestack/lib"

	amqp "github.com/rabbitmq/amqp091-go"
)

func getClient() *amqp.Channel {
	conn, err := amqp.Dial(shared.RABBITMQ_URL)
	lib.FailOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	lib.FailOnError(err, "Failed to open a channel")
	return ch
}
