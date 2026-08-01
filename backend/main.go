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

	"gorm.io/gorm"
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

	// 自动迁移插件表
	if err := db.DB().AutoMigrate(
		&model.MarketplaceListing{},
		&model.MarketplaceOrder{},
		&model.MarketplaceChat{},
		&model.MarketplaceChatSession{},
		&model.MarketplaceConfig{},
		&model.Nav{},
		&model.MenuActive{},
	); err != nil {
		logger.Warnf("表迁移失败: %v", err)
	}

	// 插入默认菜单数据（如果表为空）
	initDefaultMenus(db.DB())

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

// initDefaultMenus 初始化默认菜单数据
func initDefaultMenus(db *gorm.DB) {
	var count int64
	db.Model(&model.Nav{}).Count(&count)
	if count > 0 {
		return // 已有数据
	}

	// 激活菜单
	db.Create(&model.MenuActive{MenuType: 1, MenuID: 1})

	// 用户中心默认菜单
	defaultNavs := []model.Nav{
		{Name: "控制台", URL: "/user/dashboard", ParentID: 0, Order: 0, FaIcon: "bx bx-home-circle", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "产品与服务", URL: "", ParentID: 0, Order: 1, FaIcon: "bx bxs-grid-alt", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "订购产品", URL: "/products", ParentID: 2, Order: 0, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "我的服务", URL: "/user/products", ParentID: 2, Order: 1, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "订单管理", URL: "/user/orders", ParentID: 2, Order: 2, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "产品升降级", URL: "/user/upgrade", ParentID: 2, Order: 3, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "账户管理", URL: "", ParentID: 0, Order: 2, FaIcon: "bx bx-user", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "个人信息", URL: "/user/profile", ParentID: 7, Order: 0, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "安全中心", URL: "/user/security", ParentID: 7, Order: 1, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "实名认证", URL: "/user/verification", ParentID: 7, Order: 2, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "消息中心", URL: "/user/system-message", ParentID: 7, Order: 3, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "联系人管理", URL: "/user/contacts", ParentID: 7, Order: 4, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "第三方登录", URL: "/user/oauth-bind", ParentID: 7, Order: 5, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "财务管理", URL: "", ParentID: 0, Order: 3, FaIcon: "bx bx-dollar-circle", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "账单列表", URL: "/user/invoices", ParentID: 14, Order: 0, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "账户充值", URL: "/user/wallet", ParentID: 14, Order: 1, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "优惠券", URL: "/user/coupons", ParentID: 14, Order: 2, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "技术支持", URL: "", ParentID: 0, Order: 4, FaIcon: "bx bx-detail", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "工单列表", URL: "/user/tickets", ParentID: 18, Order: 0, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "提交工单", URL: "/user/tickets/create", ParentID: 18, Order: 1, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "帮助中心", URL: "/knowledge-base", ParentID: 18, Order: 2, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "资源下载", URL: "/downloads", ParentID: 18, Order: 3, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "新闻中心", URL: "/news", ParentID: 18, Order: 4, FaIcon: "", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "推介计划", URL: "/user/referral", ParentID: 0, Order: 5, FaIcon: "bx bxs-paper-plane", MenuType: 1, NavType: 0, MenuID: 1},
		{Name: "交易市场", URL: "/user/marketplace", ParentID: 0, Order: 6, FaIcon: "bx bx-store", MenuType: 1, NavType: 0, MenuID: 1},
	}

	for _, nav := range defaultNavs {
		db.Create(&nav)
	}
}
