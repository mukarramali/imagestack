package src

import (
	"fmt"
	"imagestack/lib/rabbitmq_service"
	"path/filepath"

	"imagestack/lib/error_handler"
	"imagestack/lib/file_handler"
	"imagestack/lib/request"

	"github.com/rabbitmq/amqp091-go"
)

var (
	downloadQueueService *rabbitmq_service.RabbitMqService
	compressQueueService *rabbitmq_service.RabbitMqService
	redisService         *request.RequestService
)

func init() {
	downloadQueueService = rabbitmq_service.NewRabbitMqService("download_images", RABBITMQ_URL)
	compressQueueService = rabbitmq_service.NewRabbitMqService("compress_images", RABBITMQ_URL)

	redisService = request.NewRequestService(REDIS_URL)
}

func ConsumeQueues() {
	go downloadQueueService.Consume(func(msg amqp091.Delivery) {
		requestId := string(msg.Body)
		request, _ := redisService.GetRequest(requestId)

		// download image
		localPath := filepath.Join(BASE_IMAGE_DIR, "raw", fmt.Sprintf("%s.jpg", requestId))
		err := file_handler.DownloadImage(request.Url, localPath)
		error_handler.FailOnError(err, "Could not download image")

		// update request with local un optimized image path
		redisService.SetLocalPathUnOptimized(requestId, localPath)

		// send event for compressing
		compressQueueService.Publish(requestId)
	})
}
