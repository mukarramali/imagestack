package src

import (
	"fmt"
	"imagestack/lib/rabbitmq_service"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"imagestack/lib/file_handler"
	"imagestack/lib/request"
)

// map to store image processing status
var (
	downloadQueueService *rabbitmq_service.RabbitMqService
)

func init() {
	downloadQueueService = rabbitmq_service.NewRabbitMqService("download_images", RABBITMQ_URL, 0)
}

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "image/webp")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Cache-Control", "public, max-age=4320")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	quality, _ := strconv.Atoi(r.FormValue("quality"))
	width, _ := strconv.Atoi(r.FormValue("width"))
	if quality == 0 {
		quality = 100
	}
	if quality > 100 {
		quality = 100
	}
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}
	setHeaders(&w)

	redisService := request.NewRequestService(REDIS_URL)
	existingRequest, _ := redisService.GetRequestByUrl(url)
	if existingRequest != nil && existingRequest.Quality == quality && existingRequest.Width == width {
		w.Header().Set("X-Cache", "HIT")
		http.ServeFile(w, r, existingRequest.LocalPathOptimized)
		return
	} else {
		w.Header().Set("X-Cache", "MISS")
	}

	imageId, _ := redisService.NewRequest(url, quality, width)

	downloadQueueService.Publish(imageId)

	futureUrl := filepath.Join(BASE_IMAGE_DIR, "compressed", fmt.Sprintf("%s.jpg", imageId))

	if file_handler.WaitForFile(futureUrl, 15*time.Second) {
		fmt.Println("Image compressed" + futureUrl)
		http.ServeFile(w, r, futureUrl)
	} else {
		fmt.Println("File could not be generated in time for " + url)
		http.Error(w, "File not ready", http.StatusGatewayTimeout)
	}
}
