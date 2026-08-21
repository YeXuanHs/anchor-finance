package main

import (
	"fmt"
	"log"

	"github.com/YeXuanHs/anchor-finance/config"
	"github.com/YeXuanHs/anchor-finance/internal/api/admin"
	"github.com/YeXuanHs/anchor-finance/internal/api/client"
	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	database.Init(&cfg.Database)

	// 自动迁移数据库表
	db := database.GetDB()
	err := db.AutoMigrate(
		&model.User{},
		&model.Admin{},
		&model.Role{},
		&model.Order{},
		&model.Invoice{},
		&model.Service{},
		&model.Ticket{},
		&model.TicketReply{},
		&model.TicketDepartment{},
		&model.TicketStatus{},
		&model.Product{},
		&model.ProductGroup{},
		&model.Plugin{},
		&model.Setting{},
		&model.Menu{},
		&model.SystemLog{},
		&model.OperationLog{},
		&model.LoginLog{},
		&model.News{},
		&model.NewsCategory{},
		&model.KnowledgeCategory{},
		&model.KnowledgeArticle{},
		&model.Download{},
		&model.DownloadCategory{},
		&model.Currency{},
		&model.PromoCode{},
		&model.Verification{},
		&model.Supplier{},
		&model.SupplierProduct{},
		&model.NotificationTemplate{},
		&model.UserNotification{},
		&model.Staff{},
		&model.MemberLevel{},
		&model.CustomField{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建服务
	authService := service.NewAuthService(db, &cfg.JWT)

	// 初始化Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS中间件
	r.Use(corsMiddleware())

	// API路由
	api := r.Group("/api")

	// 管理后台路由
	adminGroup := api.Group("/admin")
	admin.SetupRoutes(adminGroup, authService)

	// 用户前台路由
	clientGroup := api.Group("/client")
	client.SetupRoutes(clientGroup, authService)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"name":   "锚点财务",
			"en_name": "AnchorFinance",
		})
	})

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
