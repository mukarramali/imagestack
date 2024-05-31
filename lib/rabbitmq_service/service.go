package rabbitmq_service

import (
	"imagestack/lib/error_handler"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqService struct {
	Client *amqp.Channel
	Queue  *amqp.Queue
}

// prefetchCountPerConsumer should be set to 0 when being used by a publisher
func NewRabbitMqService(queueName string, amqpUrl string, prefetchCountPerConsumer int) *RabbitMqService {
	client := getClient(amqpUrl)
	queue, err := client.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	client.Qos(prefetchCountPerConsumer, 0, false)

	error_handler.FailOnError(err, "Failed to create a queue")
	return &RabbitMqService{
		Client: client,
		Queue:  &queue,
	}
}
