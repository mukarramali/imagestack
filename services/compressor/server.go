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
	"strconv"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

// map to store image processing status
var (
	downloadQueueService *rabbitmq_service.RabbitMqService
	compressQueueService *rabbitmq_service.RabbitMqService
	cleanupQueueService  *rabbitmq_service.RabbitMqService
)

func init() {
	err := os.MkdirAll(filepath.Join(shared.BASE_IMAGE_DIR, "raw"), os.ModePerm)
	shared.FailOnError(err, "Could not create images directory")
	err = os.MkdirAll(filepath.Join(shared.BASE_IMAGE_DIR, "compressed"), os.ModePerm)
	shared.FailOnError(err, "Could not create images directory")

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
		err := compress.CompressImage(request.LocalPathUnOptimized, outputPath, request.Quality)

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
		err := os.Remove(request.LocalPathUnOptimized)
		shared.FailOnError(err, "Couldn't delete uncompressed image")
		identifier.SetStatus(requestId, "completed")
	})

	go redis_service.GetRedisService()
}

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "image/webp")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Cache-Control", "public, max-age=4320")
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	quality, _ := strconv.Atoi(r.FormValue("quality"))
	if quality == 0 {
		quality = 10
	}
	if quality > 100 {
		quality = 100
	}
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	setHeaders(&w)

	existingRequest, _ := identifier.GetRequestByUrl(url)
	if existingRequest != nil {
		http.ServeFile(w, r, existingRequest.LocalPathOptimized)
		return
	}

	imageId, _ := identifier.NewRequest(url, quality)

	downloadQueueService.Publish(imageId)

	futureUrl := filepath.Join(shared.BASE_IMAGE_DIR, "compressed", fmt.Sprintf("%s.jpg", imageId))

	if shared.WaitForFile(futureUrl, 5*time.Second) {
		fmt.Println("Image compressed" + futureUrl)
		http.ServeFile(w, r, futureUrl)
	} else {
		fmt.Println("File could not be generated in time for " + url)
		http.Error(w, "File not ready", http.StatusGatewayTimeout)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Healthy"))
}

func main() {
	http.HandleFunc("/", submitHandler)
	http.HandleFunc("/health", healthCheckHandler)
	http.ListenAndServe(":8080", nil)
}
