package rabbitmq_service

import (
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

var once sync.Once
var rabbitmqClient *amqp.Channel

func getClient() *amqp.Channel {
	conn, err := amqp.Dial(RABBITMQ_HOST)
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch
}

func GetRabbitMqChannel() *amqp.Channel {
	once.Do(func() {
		rabbitmqClient = getClient()
	})

	return rabbitmqClient
}
