package infra

import (
	"context"
	"log"
	"shortvid-backend/app/shortvid-service/internal/conf"

	"github.com/redis/go-redis/v9"
)

// 初始化Redis客户端
func NewRedis(c *conf.Data) *redis.Client {
	opts := &redis.Options{
		Addr:         c.Redis.Addr,
		Password:     c.Redis.Password,
		ReadTimeout:  c.Redis.ReadTimeout.AsDuration(),
		WriteTimeout: c.Redis.WriteTimeout.AsDuration(),
	}
	client := redis.NewClient(opts)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Connect Redis failed: %v", err)
		panic(err)
	}
	log.Printf("Redis connect success...")
	return client
}
