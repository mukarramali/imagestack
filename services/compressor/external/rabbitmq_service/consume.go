package rabbitmq_service

import (
	"log"
)

func (rs *RabbitMqService) Consume() {
	msgs, err := rs.client.Consume(
		rs.queue.Name, // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func handler(msg []byte) {
	log.Printf("Received a message: %s", msg)
}
