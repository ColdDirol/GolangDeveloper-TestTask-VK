package storage

import (
	"GolangDeveloper-TestTask-VK/internal/config"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func InitRedis(redisConfig config.Redis) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
}
