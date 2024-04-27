package config

import (
	"github.com/redis/go-redis/v9"
	"os"
)

type Redis struct{}

func (r *Redis) ConnectRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
	return client
}
