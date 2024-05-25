package shared

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	REDIS_URL      string
	RABBITMQ_URL   string
	BASE_IMAGE_DIR string = "/data/images"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	REDIS_URL = os.Getenv("REDIS_URL")
	RABBITMQ_URL = os.Getenv("RABBITMQ_URL")
}
