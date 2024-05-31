package rabbitmq_service

import (
	"fmt"
	"imagestack/lib/error_handler"
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

// uses NODE_ID from env as consumer name. If not provided, generates unique tag
// Ack true when msg is consumed successfully
func (rs *RabbitMqService) Consume(handler func(msg amqp091.Delivery) error) {
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
			log.Println("TraceId:" + string(d.Body))
			go func(msg amqp091.Delivery) {
				err := handler(msg)
				if err == nil {
					msg.Ack(false)
				} else {
					fmt.Printf("Msg failed by %s, error: %s\n", nodeId, err)
					msg.Nack(false, false)
				}
			}(d)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
