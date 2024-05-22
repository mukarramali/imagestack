package main

import (
	"compressor/compress"
	"compressor/external/rabbitmq_service"
	"compressor/external/redis_service"
	"compressor/identifier"
	"compressor/load"
	"compressor/shared"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

// map to store image processing status
var (
	downloadQueueService *rabbitmq_service.RabbitMqService
	compressQueueService *rabbitmq_service.RabbitMqService
	cleanupQueueService  *rabbitmq_service.RabbitMqService
)

func init() {
	os.MkdirAll(shared.BASE_IMAGE_DIR+"raw", os.ModePerm)
	os.MkdirAll(shared.BASE_IMAGE_DIR+"compressed", os.ModePerm)

	downloadQueueService = rabbitmq_service.NewRabbitMqService("download_images")
	compressQueueService = rabbitmq_service.NewRabbitMqService("compress_images")
	cleanupQueueService = rabbitmq_service.NewRabbitMqService("cleanup_images")

	go downloadQueueService.Consume(func(msg amqp091.Delivery) {
		requestId := string(msg.Body)
		request, _ := identifier.GetRequest(requestId)

		// download image
		localPath := filepath.Join(shared.BASE_IMAGE_DIR, "raw", fmt.Sprintf("%s.jpg", requestId))
		err := load.DownloadImage(request.Url, localPath)
		shared.FailOnError(err, "Could not download image")

		// update request with local un optimized image path
		identifier.SetLocalPathUnOptimized(requestId, localPath)

		// send event for compressing
		compressQueueService.Publish(requestId)
	})
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

		// send event for cleanup
		cleanupQueueService.Publish(requestId)
	})
	go cleanupQueueService.Consume(func(msg amqp091.Delivery) {
		requestId := string(msg.Body)
		request, _ := identifier.GetRequest(requestId)
		os.Remove(request.LocalPathUnOptimized)
		identifier.SetStatus(requestId, "completed")
	})

	go redis_service.GetRedisService()
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	existingRequest, _ := identifier.GetRequestByUrl(url)
	if existingRequest != nil {
		fmt.Fprintf(w, "Go to:%s", existingRequest.LocalPathOptimized)
		return
	}

	imageId, _ := identifier.NewRequest(url)

	downloadQueueService.Publish(imageId)

	futureUrl := filepath.Join(shared.BASE_IMAGE_DIR, "compressed", fmt.Sprintf("%s.jpg", imageId))
	fmt.Fprintf(w, "Go to:%s", futureUrl)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	request, err := identifier.GetRequestByUrl(url)

	if err == redis.Nil {
		fmt.Fprintln(w, "Image never existed")
		return
	} else if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	if request.Status != "completed" {
		fmt.Fprintln(w, "Image processing status"+request.Status)
	} else {
		fmt.Fprintf(w, "Image processing complete. Download from %s", request.LocalPathOptimized)
	}
}

func main() {
	http.HandleFunc("/", submitHandler)
	http.ListenAndServe(":8080", nil)
}
