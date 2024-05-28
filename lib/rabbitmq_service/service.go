package rabbitmq_service

import (
	"imagestack/lib/error_handler"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqService struct {
	client *amqp.Channel
	queue  *amqp.Queue
}

func NewRabbitMqService(queueName string, amqpUrl string) *RabbitMqService {
	client := getClient(amqpUrl)
	queue, err := client.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	error_handler.FailOnError(err, "Failed to create a queue")
	return &RabbitMqService{
		client: client,
		queue:  &queue,
	}
}
