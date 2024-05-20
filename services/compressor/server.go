package main

import (
	"compressor/compress"
	"compressor/external/rabbitmq_service"
	"compressor/load"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// map to store image processing status
var (
	mu                   sync.Mutex
	images               map[string]string = make(map[string]string) // Maps image URL to local compressed file path
	baseDir              string            = "/data/images"
	downloadQueueService *rabbitmq_service.RabbitMqService
	compressQueueService *rabbitmq_service.RabbitMqService
	cleanupQueueService  *rabbitmq_service.RabbitMqService
)

func init() {
	os.MkdirAll(baseDir, os.ModePerm)
	downloadQueueService = rabbitmq_service.NewRabbitMqService("download_images")
	compressQueueService = rabbitmq_service.NewRabbitMqService("compress_images")
	cleanupQueueService = rabbitmq_service.NewRabbitMqService("cleanup_images")

	downloadQueueService.Consume()
	compressQueueService.Consume()
	cleanupQueueService.Consume()
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	go func() {

		downloadQueueService.Publish(url)

		fmt.Println("Downloading the image from ", url)
		localPath, err := load.DownloadImage(url)
		if err != nil {
			fmt.Println("Failed to download image:", err)
			return
		}

		outputPath := filepath.Join(baseDir, fmt.Sprintf("compressed_%d.jpg", time.Now().UnixNano()))

		fmt.Println("Compressing the image from ", url)
		err = compress.CompressImage(localPath, outputPath)
		if err != nil {
			fmt.Println("Failed to compress image:", err)
			return
		}

		mu.Lock()
		images[url] = outputPath
		mu.Unlock()
		fmt.Println("Compressed the image for ", url)
	}()
	fmt.Fprintf(w, "Image processing started. Check status at /status?url=%s", url)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	rabbitmq_service.Publish("images", url)

	mu.Lock()
	path, exists := images[url]
	mu.Unlock()

	if !exists {
		fmt.Fprintln(w, "No such image processing found or it might have been completed.")
		return
	}

	if path == "" {
		fmt.Fprintln(w, "Image processing in progress")
	} else {
		fmt.Fprintf(w, "Image processing complete. Download from %s", path)
	}
}

func main() {
	http.HandleFunc("/submit", imageHandler)
	http.HandleFunc("/status", statusHandler)
	http.ListenAndServe(":8080", nil)
}
