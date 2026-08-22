package service

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/redis/go-redis/v9"
	"context"
)

// RedisService Redis服务（可选启用，配置存数据库）
type RedisService struct {
	enabled  bool
	host     string
	port     int
	password string
	db       int
}

// NewRedisService 从数据库读取Redis配置并初始化
func NewRedisService() *RedisService {
	s := &RedisService{}
	s.loadConfig()
	return s
}

// loadConfig 从settings表加载Redis配置（按key查询，兼容所有group）
func (s *RedisService) loadConfig() {
	db := database.GetDB()

	keys := []string{"redis_enabled", "redis_host", "redis_port", "redis_password", "redis_db"}
	configMap := make(map[string]string)
	for _, key := range keys {
		var setting model.Setting
		if err := db.Where("`key` = ?", key).First(&setting).Error; err == nil {
			configMap[key] = setting.Value
		}
	}

	// 检查是否启用
	if enabled, ok := configMap["redis_enabled"]; ok && enabled == "1" {
		s.enabled = true
	} else {
		s.enabled = false
		return
	}

	s.host = configMap["redis_host"]
	if s.host == "" {
		s.host = "127.0.0.1"
	}

	if port, err := strconv.Atoi(configMap["redis_port"]); err == nil && port > 0 {
		s.port = port
	} else {
		s.port = 6379
	}

	s.password = configMap["redis_password"]

	if dbNum, err := strconv.Atoi(configMap["redis_db"]); err == nil {
		s.db = dbNum
	}
}

// IsEnabled 是否启用Redis
func (s *RedisService) IsEnabled() bool {
	return s.enabled
}

// GetConfig 获取Redis连接信息
func (s *RedisService) GetConfig() map[string]interface{} {
	return map[string]interface{}{
		"enabled":  s.enabled,
		"host":     s.host,
		"port":     s.port,
		"password": "***",
		"db":       s.db,
	}
}

// HealthCheck 检查Redis连接状态
func (s *RedisService) HealthCheck() map[string]interface{} {
	if !s.enabled {
		return map[string]interface{}{
			"status":  "disabled",
			"message": "Redis未启用",
		}
	}

	// 尝试连接
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	// 简单的TCP连接检查
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return map[string]interface{}{
			"status":  "error",
			"message": fmt.Sprintf("连接失败: %v", err),
			"addr":    addr,
		}
	}
	conn.Close()

	return map[string]interface{}{
		"status":  "ok",
		"message": "连接正常",
		"addr":    addr,
	}
}

// Set 设置缓存（如果Redis启用则存Redis，否则跳过）
func (s *RedisService) Set(key string, value interface{}, ttl time.Duration) error {
	if !s.enabled {
		return nil
	}
	ctx := context.Background()
	client := s.getClient()
	defer client.Close()
	return client.Set(ctx, key, fmt.Sprintf("%v", value), ttl).Err()
}

// Get 获取缓存
func (s *RedisService) Get(key string) (string, error) {
	if !s.enabled {
		return "", fmt.Errorf("redis未启用")
	}
	ctx := context.Background()
	client := s.getClient()
	defer client.Close()
	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key不存在")
	}
	return val, err
}

// Delete 删除缓存
func (s *RedisService) Delete(key string) error {
	if !s.enabled {
		return nil
	}
	ctx := context.Background()
	client := s.getClient()
	defer client.Close()
	return client.Del(ctx, key).Err()
}

// getClient 获取Redis客户端连接
func (s *RedisService) getClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", s.host, s.port),
		Password: s.password,
		DB:       s.db,
	})
}
