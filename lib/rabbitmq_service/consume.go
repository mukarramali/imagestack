package rabbitmq_service

import (
	"imagestack/lib/error_handler"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

// uses NODE_ID from env as consumer name. If not provided, generates unique tag
// Ack true when msg is consumed successfully
func (rs *RabbitMqService) Consume(handler func(msg amqp091.Delivery)) {
	nodeId := os.Getenv("NODE_ID")
	msgs, err := rs.Client.Consume(
		rs.Queue.Name, // queue
		nodeId,        // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	error_handler.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Println(d)
			handler(d)
			d.Ack(true)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
