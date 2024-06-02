package src

import (
	"fmt"
	"imagestack/lib/rabbitmq_service"
	"net/http"
	"path/filepath"
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

func Handler(w http.ResponseWriter, r *http.Request) {
	params, err := GetSafeParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	setHeaders(&w)

	redisService := request.NewRequestService(REDIS_URL)
	existingRequest, _ := redisService.GetRequestByUrl(params.url)
	if existingRequest != nil {
		if existingRequest.Quality == params.quality && existingRequest.Width == params.width {
			w.Header().Set("X-Cache", "HIT")
			http.ServeFile(w, r, existingRequest.LocalPathOptimized)
			return
		} else if existingRequest.Status == "error" {
			http.Error(w, "Check the source", http.StatusNotFound)
			return
		}
	}

	w.Header().Set("X-Cache", "MISS")

	imageId, _ := redisService.NewRequest(params.url, params.quality, params.width)

	downloadQueueService.Publish(imageId)

	futureUrl := filepath.Join(BASE_IMAGE_DIR, "compressed", fmt.Sprintf("%s.jpg", imageId))

	if file_handler.WaitForFile(futureUrl, 20*time.Second) {
		fmt.Println("Image compressed" + futureUrl)
		http.ServeFile(w, r, futureUrl)
	} else {
		fmt.Println("File could not be generated in time for traceId:" + imageId)
		http.Error(w, "File not ready", http.StatusGatewayTimeout)
	}
}
