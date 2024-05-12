package load

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	baseDir string = "compressed_images"
)

func DownloadImage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	localPath := filepath.Join(baseDir, fmt.Sprintf("%d.jpg", time.Now().UnixNano()))
	file, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return localPath, err
}
