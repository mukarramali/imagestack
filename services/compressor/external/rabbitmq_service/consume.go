package rabbitmq_service

import (
	"imagestack/lib"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func (rs *RabbitMqService) Consume(handler func(msg amqp091.Delivery)) {
	msgs, err := rs.client.Consume(
		rs.queue.Name, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	lib.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Println(d)
			handler(d)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
