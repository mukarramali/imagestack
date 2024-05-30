package main

import (
	"api/shared"
	"api/src"
	"net/http"
	"os"
	"path/filepath"

	"imagestack/lib/error_handler"
)

func init() {
	err := os.MkdirAll(filepath.Join(shared.BASE_IMAGE_DIR, "raw"), os.ModePerm)
	error_handler.FailOnError(err, "Could not create images directory")
	err = os.MkdirAll(filepath.Join(shared.BASE_IMAGE_DIR, "compressed"), os.ModePerm)
	error_handler.FailOnError(err, "Could not create images directory")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is Healthy"))
}

func main() {
	http.HandleFunc("/", src.Handler)
	http.HandleFunc("/health", healthCheckHandler)
	http.ListenAndServe(":8080", nil)
}
