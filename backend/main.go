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
	"anchorfinance/internal/model"
	"anchorfinance/pkg/auth"
	"anchorfinance/pkg/db"
	"anchorfinance/pkg/logger"
	"github.com/gin-gonic/gin"
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
	if _, err := db.InitDB(cfg.Database.ToDBConfig()); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// Auto-migrate new models/columns
	if conn := db.GetDB(); conn != nil {
		conn.AutoMigrate(&model.TicketTransferLog{})
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

	// JWT 配置（优先从 config.yaml 读取，否则从数据库读取）
	jwtSecret := cfg.JWT.Secret
	if jwtSecret == "" {
		jwtSecret = db.GetSystemSetting("jwt_secret")
	}
	if jwtSecret == "" {
		jwtSecret = "default_secret_please_change"
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
			log.Printf("Redis 初始化失败（非致命）: %v", err)
		}
	}

	// 创建 Gin 引擎
	r := gin.Default()

	// 注册静态资源（前端编译产物）
	frontendDist := filepath.Join("frontend", "dist")
	if _, err := os.Stat(frontendDist); err == nil {
		r.Static("/assets", filepath.Join(frontendDist, "assets"))
		r.StaticFile("/favicon.ico", filepath.Join(frontendDist, "favicon.ico"))
		r.StaticFile("/logo.png", filepath.Join(frontendDist, "logo.png"))
	}

	// 注册 API 路由
	api.RegisterRoutes(r, jwtMgr)

	// 注册后台管理 API
	adminPath := db.GetSystemSetting("admin_path")
	if adminPath == "" {
		adminPath = "/admin"
	}
	api.RegisterAdminRoutes(r, jwtMgr, adminPath)

	// SPA 回退：所有非 API 请求返回前端 index.html
	r.NoRoute(func(c *gin.Context) {
		// API 请求返回 404
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"error": "接口不存在"})
			return
		}
		// 后台请求返回后台 index.html
		if len(c.Request.URL.Path) >= len(adminPath) && c.Request.URL.Path[:len(adminPath)] == adminPath {
			c.File(filepath.Join(frontendDist, "admin", "index.html"))
			return
		}
		// 其他请求返回前台 index.html
		c.File(filepath.Join(frontendDist, "index.html"))
	})

	// 启动定时任务
	job.StartAll()

	// 优雅关闭
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("收到关闭信号，正在停止...")
		job.StopAll()
		os.Exit(0)
	}()

	// 启动服务
	host := cfg.Server.Host
	if host == "" {
		host = "127.0.0.1"
	}
	port := cfg.Server.Port
	if port == 0 {
		port = 8080
	}

	fmt.Printf("========================================\n")
	fmt.Printf("  锚点财务 启动成功\n")
	fmt.Printf("  监听地址: %s:%d\n", host, port)
	fmt.Printf("  站点地址: http://localhost:%d\n", port)
	fmt.Printf("  后台地址: http://localhost:%d%s\n", port, adminPath)
	fmt.Printf("========================================\n")

	if err := r.Run(fmt.Sprintf("%s:%d", host, port)); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
