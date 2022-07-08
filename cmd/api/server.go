package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/hellodhlyn/tenri-api/internal/database"
	"os"
)

type serverContext struct {
	rdb *redis.Client
}

func main() {
	redisURL, ok := os.LookupEnv("REDIS_URL")
	if !ok {
		redisURL = "127.0.0.1:6379"
	}

	ctx := serverContext{
		rdb: database.NewRedisClient(redisURL),
	}

	server := gin.Default()

	q := server.Group("/q/v1")
	q.GET("/tasks", getTasks(ctx))
	q.POST("/tasks", postTask(ctx))
	q.PATCH("/tasks/:uuid", patchTask(ctx))
	q.DELETE("/tasks/:uuid", deleteTask(ctx))

	fmt.Println(server.Run("0.0.0.0:8080"))
}
