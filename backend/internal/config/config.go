package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 仅包含启动所需的最小配置（从 .env 读取），其余全部存数据库
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Host string
	Port int
	Mode string
}

type DatabaseConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}

// Load 从 .env 文件加载配置（仅数据库和服务器信息）
// .env 文件位置：当前目录或上一级目录
func Load() (*Config, error) {
	// 尝试加载 .env 文件（不强制存在，也支持环境变量）
	envPaths := []string{".env", "../.env"}
	for _, p := range envPaths {
		if abs, err := filepath.Abs(p); err == nil {
			if _, err := os.Stat(abs); err == nil {
				_ = godotenv.Load(abs)
				break
			}
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "127.0.0.1"),
			Port: getEnvInt("SERVER_PORT", 8080),
			Mode: getEnv("SERVER_MODE", "release"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "3306"),
			User:         getEnv("DB_USER", "root"),
			Password:     getEnv("DB_PASS", ""),
			DBName:       getEnv("DB_NAME", "anchorfinance"),
			Charset:      getEnv("DB_CHARSET", "utf8mb4"),
			MaxIdleConns: getEnvInt("DB_MAX_IDLE", 10),
			MaxOpenConns: getEnvInt("DB_MAX_OPEN", 100),
		},
	}

	// 校验必填项
	if cfg.Database.Host == "" || cfg.Database.User == "" || cfg.Database.DBName == "" {
		return nil, fmt.Errorf("数据库配置不完整，请检查 .env 文件（DB_HOST, DB_USER, DB_NAME）")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return defaultVal
}
