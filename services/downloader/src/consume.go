package src

import (
	"fmt"
	"imagestack/lib/rabbitmq_service"
	"path/filepath"

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
	downloadQueueService = rabbitmq_service.NewRabbitMqService("download_images", RABBITMQ_URL, 10)
	compressQueueService = rabbitmq_service.NewRabbitMqService("compress_images", RABBITMQ_URL, 0)

	redisService = request.NewRequestService(REDIS_URL)
}

func ConsumeQueues() {
	go downloadQueueService.Consume(func(msg amqp091.Delivery) error {
		requestId := string(msg.Body)
		request, _ := redisService.GetRequest(requestId)

		// download image
		localPath := filepath.Join(BASE_IMAGE_DIR, "raw", fmt.Sprintf("%s.jpg", requestId))
		err := file_handler.DownloadImage(request.Url, localPath)
		if err != nil {
			return err
		}

		// update request with local un optimized image path
		redisService.SetLocalPathUnOptimized(requestId, localPath)

		// send event for compressing
		compressQueueService.Publish(requestId)
		return nil
	})
}
