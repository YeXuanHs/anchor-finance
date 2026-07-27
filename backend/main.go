package main

import (
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
	"anchorfinance/pkg/auth"
	"anchorfinance/pkg/db"
	"anchorfinance/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// 检查是否已安装（config.yaml 是否存在）
	configPath := filepath.Join("configs", "config.yaml")
	installed := fileExists(configPath)

	if !installed {
		// 未安装：只启动安装页面
		startInstaller()
		return
	}

	// 已安装：正常启动
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	if err := db.InitDB(cfg.Database); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

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

	// 从数据库读取 JWT 配置
	jwtSecret := db.GetSystemSetting("jwt_secret")
	if jwtSecret == "" {
		jwtSecret = "default_secret_please_change"
	}
	jwtExpireStr := db.GetSystemSetting("jwt_expire_hours")
	jwtExpire := 72
	if v, err := strconv.Atoi(jwtExpireStr); err == nil && v > 0 {
		jwtExpire = v
	}
	jwtMgr := auth.NewJWTManager(jwtSecret, jwtExpire)

	// Redis 可选：从数据库读取配置，如果启用了再初始化
	redisEnabled := db.GetSystemSetting("redis_enabled")
	if redisEnabled == "true" {
		if err := db.InitRedisFromDB(); err != nil {
			logger.Warnf("Redis 初始化失败，将使用内存缓存: %v", err)
		}
	}

	// 启动定时任务
	go job.Start()

	// 启动 HTTP 服务
	srv := api.NewServer(cfg, jwtMgr)
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		logger.Infof("锚点财务服务启动: http://localhost%s", addr)
		if err := srv.Run(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务启动失败: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("服务正在关闭...")
}

// startInstaller 启动仅包含安装页面的最小服务
func startInstaller() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	installH := api.NewInstallHandler()
	installH.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("========================================\n")
	fmt.Printf("  锚点财务 - 安装向导\n")
	fmt.Printf("  请访问: http://localhost:%s/install\n", port)
	fmt.Printf("========================================\n")

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("安装服务启动失败: %v", err)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
