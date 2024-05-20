package rabbitmq_service

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func getClient() *amqp.Channel {
	conn, err := amqp.Dial(RABBITMQ_HOST)
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch
}
