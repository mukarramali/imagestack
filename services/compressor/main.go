package main

import (
	"compressor/src"
	"fmt"
	"math/rand"
	"net/http"
	"os"
)

func init() {
	if os.Getenv("NODE_ID") == "" {
		os.Setenv("NODE_ID", fmt.Sprintf("compressor:%d", rand.Intn(100)+1))
	}

	src.ConsumeQueues()
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Compressor Service is Healthy"))
}

func main() {
	http.HandleFunc("/health", healthCheckHandler)
	http.ListenAndServe(":8080", nil)
}
