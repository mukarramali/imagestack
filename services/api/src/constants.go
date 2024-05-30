package src

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	REDIS_URL      string
	RABBITMQ_URL   string
	BASE_IMAGE_DIR string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Didn't load .env file, make sure you have passed env from somewhere else.")
	}
	REDIS_URL = os.Getenv("REDIS_URL")
	RABBITMQ_URL = os.Getenv("RABBITMQ_URL")
	BASE_IMAGE_DIR = os.Getenv("BASE_IMAGE_DIR")
}
