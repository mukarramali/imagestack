package main

import (
	"compressor/compress"
	"compressor/shared"
	"fmt"
	"imagestack/lib/rabbitmq_service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"imagestack/lib/error_handler"
	"imagestack/lib/file_handler"
	"imagestack/lib/request"

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
	error_handler.FailOnError(err, "Could not create images directory")
	err = os.MkdirAll(filepath.Join(shared.BASE_IMAGE_DIR, "compressed"), os.ModePerm)
	error_handler.FailOnError(err, "Could not create images directory")

	downloadQueueService = rabbitmq_service.NewRabbitMqService("download_images", shared.RABBITMQ_URL)
	compressQueueService = rabbitmq_service.NewRabbitMqService("compress_images", shared.RABBITMQ_URL)
	cleanupQueueService = rabbitmq_service.NewRabbitMqService("cleanup_images", shared.RABBITMQ_URL)

	redisService := request.NewRequestService(shared.REDIS_URL)

	go downloadQueueService.Consume(func(msg amqp091.Delivery) {
		requestId := string(msg.Body)
		request, _ := redisService.GetRequest(requestId)

		// download image
		localPath := filepath.Join(shared.BASE_IMAGE_DIR, "raw", fmt.Sprintf("%s.jpg", requestId))
		err := file_handler.DownloadImage(request.Url, localPath)
		error_handler.FailOnError(err, "Could not download image")

		// update request with local un optimized image path
		redisService.SetLocalPathUnOptimized(requestId, localPath)

		// send event for compressing
		compressQueueService.Publish(requestId)
	})
	go compressQueueService.Consume(func(msg amqp091.Delivery) {
		requestId := string(msg.Body)
		request, _ := redisService.GetRequest(requestId)

		outputPath := filepath.Join(shared.BASE_IMAGE_DIR, "compressed", fmt.Sprintf("%s.jpg", requestId))

		// compress
		err := compress.CompressImage(request.LocalPathUnOptimized, outputPath, request.Quality, request.Width)

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

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "image/webp")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Cache-Control", "public, max-age=4320")
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	quality, _ := strconv.Atoi(r.FormValue("quality"))
	width, _ := strconv.Atoi(r.FormValue("width"))
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

	redisService := request.NewRequestService(shared.REDIS_URL)
	existingRequest, _ := redisService.GetRequestByUrl(url)
	if existingRequest != nil && existingRequest.Quality == quality && existingRequest.Width == width {
		http.ServeFile(w, r, existingRequest.LocalPathOptimized)
		return
	}

	imageId, _ := redisService.NewRequest(url, quality, width)

	downloadQueueService.Publish(imageId)

	futureUrl := filepath.Join(shared.BASE_IMAGE_DIR, "compressed", fmt.Sprintf("%s.jpg", imageId))

	if file_handler.WaitForFile(futureUrl, 15*time.Second) {
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
