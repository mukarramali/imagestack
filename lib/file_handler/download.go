package file_handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadImage(url string, localPath string) error {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Could not get image from url:" + url)
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(localPath)
	if err != nil {
		fmt.Println("Could not create file for localPath:" + localPath)
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println("Could not copy image into localPath:" + localPath)
	}
	return err
}
