package redis_service

import (
	"compressor/shared"
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var once sync.Once
var redisClient *redis.Client

func getClient() *redis.Client {
	redisOptions, err := redis.ParseURL(shared.REDIS_URL)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(redisOptions)
	return client
}

func GetRedisService() *redis.Client {
	once.Do(func() {
		redisClient = getClient()

		_, err := redisClient.Ping(context.Background()).Result()
		if err != nil {
			fmt.Println("Error connecting to Redis:", err)
		} else {
			fmt.Println("Redis connected")
		}
	})

	return redisClient
}
