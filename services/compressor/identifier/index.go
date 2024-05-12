package identifier

import (
	"compressor/data/redis_service"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func GetUuid(url string) (uuid.UUID, error) {
	redisService := redis_service.GetRedisService()

	existingId, err := redisService.Get(context.Background(), url).Result()
	if err != redis.Nil {
		id, _ := uuid.Parse(existingId)
		return id, nil
	}

	newId := uuid.New()

	err = redisService.Set(context.Background(), url, newId, 0).Err()

	if err != nil {
		fmt.Println("Error setting key in Redis:", err)
		return newId, err
	}

	return newId, nil
}
