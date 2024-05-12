package main

import (
	"compressor/compress"
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
	mu      sync.Mutex
	images  map[string]string = make(map[string]string) // Maps image URL to local compressed file path
	baseDir string            = "compressed_images"
)

func init() {
	os.MkdirAll(baseDir, os.ModePerm)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	_, exists := images[url]
	if exists {
		mu.Unlock()
		fmt.Fprintf(w, "Image is already being processed")
		return
	}
	images[url] = ""
	mu.Unlock()

	go func() {
		defer func() {
			mu.Lock()
			delete(images, url)
			mu.Unlock()
		}()

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
