package src

import (
	"compressor/shared"
	"fmt"
	"imagestack/lib/rabbitmq_service"
	"os"
	"path/filepath"

	"imagestack/lib/error_handler"
	"imagestack/lib/request"

	"github.com/rabbitmq/amqp091-go"
)

var (
	compressQueueService *rabbitmq_service.RabbitMqService
	cleanupQueueService  *rabbitmq_service.RabbitMqService
	redisService         *request.RequestService
)

func init() {
	compressQueueService = rabbitmq_service.NewRabbitMqService("compress_images", shared.RABBITMQ_URL)
	cleanupQueueService = rabbitmq_service.NewRabbitMqService("cleanup_images", shared.RABBITMQ_URL)

	redisService = request.NewRequestService(shared.REDIS_URL)
}

func RegisterQueues() {
	go compressQueueService.Consume(func(msg amqp091.Delivery) {
		requestId := string(msg.Body)
		request, _ := redisService.GetRequest(requestId)

		outputPath := filepath.Join(shared.BASE_IMAGE_DIR, "compressed", fmt.Sprintf("%s.jpg", requestId))

		// compress
		err := CompressImage(request.LocalPathUnOptimized, outputPath, request.Quality, request.Width)

		if err != nil {
			error_handler.FailOnError(err, "Could not compress image")
		}

		redisService.SetLocalPathOptimized(requestId, outputPath)

		// send event for cleanup
		cleanupQueueService.Publish(requestId)
	})
	go cleanupQueueService.Consume(func(msg amqp091.Delivery) {
		requestId := string(msg.Body)
		request, _ := redisService.GetRequest(requestId)
		err := os.Remove(request.LocalPathUnOptimized)
		error_handler.FailOnError(err, "Couldn't delete uncompressed image")
		redisService.SetStatus(requestId, "completed")
	})
}
