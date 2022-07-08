package database

import (
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(url string) *redis.Client {
	client := redis.NewClient(&redis.Options{Addr: url})
	return client
}
