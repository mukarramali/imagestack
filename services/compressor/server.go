package main

import (
	"compressor/compress"
	"compressor/external/rabbitmq_service"
	"compressor/external/redis_service"
	"compressor/identifier"
	"compressor/shared"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rabbitmq/amqp091-go"
)

// map to store image processing status
var (
	compressQueueService *rabbitmq_service.RabbitMqService
)

func init() {
	err := os.MkdirAll(filepath.Join(shared.BASE_IMAGE_DIR, "raw"), os.ModePerm)
	shared.FailOnError(err, "Could not create images directory")
	err = os.MkdirAll(filepath.Join(shared.BASE_IMAGE_DIR, "compressed"), os.ModePerm)
	shared.FailOnError(err, "Could not create images directory")

	compressQueueService = rabbitmq_service.NewRabbitMqService("compress_images")

	go compressQueueService.Consume(func(msg amqp091.Delivery) {
		requestId := string(msg.Body)
		request, _ := identifier.GetRequest(requestId)

		outputPath := filepath.Join(shared.BASE_IMAGE_DIR, "compressed", fmt.Sprintf("%s.jpg", requestId))

		// compress
		err := compress.CompressImage(request.LocalPathUnOptimized, outputPath)

		if err != nil {
			shared.FailOnError(err, "Could not compress image")
		}

		identifier.SetLocalPathOptimized(requestId, outputPath)
	})
	go redis_service.GetRedisService()
}

func main() {
}
