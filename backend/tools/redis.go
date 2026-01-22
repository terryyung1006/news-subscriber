package tools

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Printf("Warning: Failed to connect to Redis at %s: %v\n", redisAddr, err)
	} else {
		fmt.Printf("Connected to Redis at %s\n", redisAddr)
	}
}
