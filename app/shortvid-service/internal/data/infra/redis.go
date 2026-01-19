package infra

import (
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
	return redis.NewClient(opts)
}
