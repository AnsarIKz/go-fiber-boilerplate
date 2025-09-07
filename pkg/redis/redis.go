package redis

import (
	"context"
	"fmt"
	"nodabackend/pkg/env"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	redisOnce   sync.Once
	redisClient *redis.Client
	errRedis    error
)

// RedisConfig содержит конфигурацию для подключения к Redis
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// NewClient создает и возвращает новый синглтон-клиент Redis
func NewClient() (*redis.Client, error) {
	redisOnce.Do(func() {
		redisDBStr := env.GetEnvOrDefault("REDIS_DB", "0")
		redisDB, err := strconv.Atoi(redisDBStr)
		if err != nil {
			errRedis = fmt.Errorf("invalid REDIS_DB value: '%s'", redisDBStr)
			return
		}

		config := RedisConfig{
			Addr:     env.GetEnvOrDefault("REDIS_ADDR", "localhost:6379"),
			Password: env.GetEnvOrDefault("REDIS_PASSWORD", ""),
			DB:       redisDB,
		}

		rdb := redis.NewClient(&redis.Options{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := rdb.Ping(ctx).Err(); err != nil {
			errRedis = fmt.Errorf("failed to connect to Redis at %s: %w", config.Addr, err)
			return
		}

		fmt.Println("[DB] Successfully connected to Redis")
		redisClient = rdb
	})

	return redisClient, errRedis
}
