package rabbitmq_service

import (
	"log"
)

var RABBITMQ_HOST = "amqp://user:password@rabbitmq:5672/"

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
