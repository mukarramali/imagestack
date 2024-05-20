package main

import (
	"compressor/compress"
	"compressor/external/rabbitmq_service"
	"compressor/identifier"
	"compressor/load"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

// map to store image processing status
var (
	baseDir              string = "/data/images"
	downloadQueueService *rabbitmq_service.RabbitMqService
	compressQueueService *rabbitmq_service.RabbitMqService
	cleanupQueueService  *rabbitmq_service.RabbitMqService
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func init() {
	os.MkdirAll(baseDir, os.ModePerm)

	downloadQueueService = rabbitmq_service.NewRabbitMqService("download_images")
	compressQueueService = rabbitmq_service.NewRabbitMqService("compress_images")
	cleanupQueueService = rabbitmq_service.NewRabbitMqService("cleanup_images")

	downloadQueueService.Consume(func(msg amqp091.Delivery) {
		requestId := string(msg.Body)
		request, _ := identifier.GetRequest(requestId)

		// download image
		localPath, err := load.DownloadImage(request.Url)
		failOnError(err, "Could not download image")

		// update request with local un optimized image path
		identifier.SetLocalPathUnOptimized(requestId, localPath)

		// send event for compressing
		compressQueueService.Publish(requestId)
	})
	compressQueueService.Consume(func(msg amqp091.Delivery) {
		requestId := string(msg.Body)
		request, _ := identifier.GetRequest(requestId)

		outputPath := filepath.Join(baseDir, fmt.Sprintf("compressed_%s.jpg", requestId))

		// compress
		compress.CompressImage(request.LocalPathUnOptimized, outputPath)

		// send event for cleanup
		cleanupQueueService.Publish(requestId)
	})
	cleanupQueueService.Consume(func(msg amqp091.Delivery) {
		requestId := string(msg.Body)
		request, _ := identifier.GetRequest(requestId)
		os.Remove(request.LocalPathUnOptimized)
		identifier.SetStatus(requestId, "completed")
	})
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	imageId, _ := identifier.NewRequest(url)

	downloadQueueService.Publish(imageId)
	fmt.Fprintf(w, "Image processing started. Check status at /status?url=%s", url)
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
	}

	failOnError(err, "Could not check")

	if request.Status != "completed" {
		fmt.Fprintln(w, "Image processing status"+request.Status)
	} else {
		fmt.Fprintf(w, "Image processing complete. Download from %s", request.LocalPathOptimized)
	}
}

func main() {
	http.HandleFunc("/submit", imageHandler)
	http.HandleFunc("/status", statusHandler)
	http.ListenAndServe(":8080", nil)
}
