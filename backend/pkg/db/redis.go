package db

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// RedisConfig holds Redis connection parameters.
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

var redisClient *redis.Client

// InitRedis creates and tests a Redis client connection.
func InitRedis(cfg RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	redisClient = rdb
	return rdb, nil
}

// InitRedisFromDB reads Redis config from database and initializes.
func InitRedisFromDB() error {
	host := GetSystemSetting("redis_host")
	port := GetSystemSetting("redis_port")
	password := GetSystemSetting("redis_password")
	dbStr := GetSystemSetting("redis_db")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "6379"
	}
	dbNum := 0
	if dbStr != "" {
		fmt.Sscanf(dbStr, "%d", &dbNum)
	}

	_, err := InitRedis(RedisConfig{
		Host:     host,
		Port:     port,
		Password: password,
		DB:       dbNum,
	})
	return err
}

// GetRedis returns the Redis client (may be nil if not enabled).
func GetRedis() *redis.Client {
	return redisClient
}

