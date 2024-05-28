package rabbitmq_service

import (
	"imagestack/lib/error_handler"

	amqp "github.com/rabbitmq/amqp091-go"
)

func getClient(url string) *amqp.Channel {
	conn, err := amqp.Dial(url)
	error_handler.FailOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	error_handler.FailOnError(err, "Failed to open a channel")
	return ch
}
