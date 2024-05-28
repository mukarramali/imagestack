package rabbitmq_service

import (
	"imagestack/lib"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqService struct {
	client *amqp.Channel
	queue  *amqp.Queue
}

func NewRabbitMqService(queueName string) *RabbitMqService {
	client := getClient()
	queue, err := client.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	lib.FailOnError(err, "Failed to create a queue")
	return &RabbitMqService{
		client: client,
		queue:  &queue,
	}
}
