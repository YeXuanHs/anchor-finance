package main

import (
	crand "crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"anchorfinance/internal/api"
	"anchorfinance/internal/config"
	"anchorfinance/internal/job"
	"anchorfinance/internal/model"
	"anchorfinance/pkg/auth"
	"anchorfinance/pkg/db"
	"anchorfinance/pkg/logger"
)

func main() {
	// 检查是否已安装（config.yaml 是否存在）
	configPath := filepath.Join("configs", "config.yaml")
	if !fileExists(configPath) {
		fmt.Println("========================================")
		fmt.Println("  锚点财务 - 未检测到配置文件")
		fmt.Println("  请先运行安装脚本: bash install.sh")
		fmt.Println("========================================")
		os.Exit(1)
	}

	// 已安装：正常启动
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库（MySQL）
	if _, err := db.InitDB(db.Config{
		Host:         cfg.Database.Host,
		Port:         cfg.Database.Port,
		User:         cfg.Database.User,
		Password:     cfg.Database.Password,
		DBName:       cfg.Database.DBName,
		Charset:      cfg.Database.Charset,
		MaxIdleConns: cfg.Database.MaxIdleConns,
		MaxOpenConns: cfg.Database.MaxOpenConns,
	}); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// Auto-migrate: 仅核心表，插件表在启用时动态创建
	// 核心表结构由 scripts/init.sql 管理

	// 从数据库读取日志配置并初始化
	logLevel := db.GetSystemSetting("log_level")
	if logLevel == "" {
		logLevel = "info"
	}
	logFormat := db.GetSystemSetting("log_format")
	if logFormat == "" {
		logFormat = "text"
	}
	logger.Init(logger.Config{
		Level:  logLevel,
		Format: logFormat,
		Output: "stdout",
	})

	// JWT 配置（优先从 config.yaml 读取，否则从数据库读取）
	jwtSecret := cfg.JWT.Secret
	if jwtSecret == "" {
		jwtSecret = db.GetSystemSetting("jwt_secret")
	}
	if jwtSecret == "" {
		// 首次启动：生成随机密钥并持久化到数据库
		jwtSecret = generateRandomSecret(64)
		if err := db.SetSystemSetting("jwt_secret", jwtSecret, "jwt", "JWT 签名密钥（自动生成）"); err != nil {
			log.Fatalf("无法保存 JWT 密钥: %v", err)
		}
		logger.Warn("已自动生成 JWT 密钥并保存到数据库")
	}
	jwtExpire := cfg.JWT.ExpireHours
	if jwtExpire == 0 {
		jwtExpireStr := db.GetSystemSetting("jwt_expire_hours")
		if v, err := strconv.Atoi(jwtExpireStr); err == nil && v > 0 {
			jwtExpire = v
		} else {
			jwtExpire = 72
		}
	}
	jwtMgr := auth.NewJWTManager(jwtSecret, jwtExpire)

	// Redis 可选：从数据库读取配置，如果启用了再初始化
	redisEnabled := db.GetSystemSetting("redis_enabled")
	if redisEnabled == "true" {
		if err := db.InitRedisFromDB(); err != nil {
			logger.Warnf("Redis 初始化失败（非致命）: %v", err)
		}
	}

	// 启动定时任务
	go job.Start()

	// 创建并启动 HTTP 服务（通过 api.NewServer 注册所有路由和中间件）
	srv := api.NewServer(cfg, jwtMgr)
	go func() {
		host := cfg.Server.Host
		if host == "" {
			host = "127.0.0.1"
		}
		addr := fmt.Sprintf("%s:%d", host, cfg.Server.Port)
		logger.Infof("锚点财务服务启动: http://%s", addr)
		if err := srv.Run(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("服务正在关闭...")
	job.StopAll()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func generateRandomSecret(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, n)
	crand.Read(b)
	for i := range b {
		b[i] = letters[int(b[i])%len(letters)]
	}
	return string(b)
}
