package redis_service

import (
	"context"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var once sync.Once
var redisClient *redis.Client

func getClient(redisUrl string) *redis.Client {
	redisOptions, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(redisOptions)
	return client
}

func GetRedisService(redisUrl string) *redis.Client {
	once.Do(func() {
		redisClient = getClient(redisUrl)

		_, err := redisClient.Ping(context.Background()).Result()
		if err != nil {
			fmt.Println("Error connecting to Redis:", err)
		} else {
			fmt.Println("Redis connected")
		}
	})

	return redisClient
}
