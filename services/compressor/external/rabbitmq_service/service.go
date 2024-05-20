package rabbitmq_service

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqService struct {
	client *amqp.Channel
	queue  *amqp.Queue
}
