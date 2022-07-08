package database

import (
	"github.com/go-redis/redis/v8"
	"os"
)

func NewRedisClient() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "127.0.0.1:6379"
	}

	client := redis.NewClient(&redis.Options{Addr: redisURL})
	return client
}
