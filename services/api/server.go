package main

import (
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
)

// map to store image processing status
var (
	downloadQueueService *rabbitmq_service.RabbitMqService
)

func init() {
	err := os.MkdirAll(filepath.Join(shared.BASE_IMAGE_DIR, "raw"), os.ModePerm)
	error_handler.FailOnError(err, "Could not create images directory")
	err = os.MkdirAll(filepath.Join(shared.BASE_IMAGE_DIR, "compressed"), os.ModePerm)
	error_handler.FailOnError(err, "Could not create images directory")

	downloadQueueService = rabbitmq_service.NewRabbitMqService("download_images", shared.RABBITMQ_URL)
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
	w.Write([]byte("API is Healthy"))
}

func main() {
	http.HandleFunc("/", submitHandler)
	http.HandleFunc("/health", healthCheckHandler)
	http.ListenAndServe(":8080", nil)
}
