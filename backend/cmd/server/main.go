package main

import (
	"log"
	"os"

	"github.com/anchorfinance/backend/internal/router"
	"github.com/anchorfinance/backend/pkg/database"
)

func main() {
	// 初始化数据库
	dbCfg := &database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "anchorfinance"),
	}

	if err := database.Init(dbCfg); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	defer database.Close()

	// 设置路由
	r := router.SetupRouter()

	// 启动服务器
	port := getEnv("PORT", "8080")
	log.Printf("服务器启动在端口 %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
