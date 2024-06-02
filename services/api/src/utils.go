package src

import (
	"errors"
	"net/http"
	"strconv"
)

func setHeaders(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "image/webp")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Cache-Control", "public, max-age=4320")
}

type params struct {
	url     string
	quality int
	width   int
}

func GetSafeParams(r *http.Request) (*params, error) {
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
		return nil, errors.New("URL is required")
	}
	return &params{url, quality, width}, nil
}
