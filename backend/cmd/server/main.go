package main

import (
	"fmt"
	"log"
	"time"

	"github.com/YeXuanHs/anchor-finance/config"
	"github.com/YeXuanHs/anchor-finance/internal/api/admin"
	"github.com/YeXuanHs/anchor-finance/internal/api/client"
	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/middleware"
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
		&model.OrderItem{},
		&model.Invoice{},
		&model.Service{},
		&model.Ticket{},
		&model.TicketReply{},
		&model.TicketDepartment{},
		&model.TicketStatus{},
		&model.Product{},
		&model.ProductGroup{},
		&model.ProductType{},
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
		&model.Coupon{},
		&model.CouponCampaign{},
		&model.Payment{},
		&model.TicketPrereply{},
		&model.TicketPrereplyCategory{},
		&model.CreditLimit{},
		&model.CreditLimitLog{},
		&model.Recharge{},
		&model.MediaFile{},
		&model.HomeHero{},
		&model.Blacklist{},
		&model.EmailTemplate{},
		&model.SMSTemplate{},
		&model.FriendlyLink{},
		&model.Contract{},
		&model.ContractTemplate{},
		&model.MarketingPush{},
		&model.CancelRequest{},
		&model.ConfigurableOption{},
		&model.OAuthProvider{},
		&model.CustomTemplateField{},
		&model.TrafficPackage{},
		&model.TrafficLog{},
		&model.TaskQueue{},
		&model.TwoFactorConfig{},
		&model.UserTwoFactor{},
		&model.SalesConfig{},
		&model.SalesGroup{},
		&model.Theme{},
		&model.TicketRule{},
		&model.OrderConfig{},
		&model.CPUModel{},
		&model.InstanceSpec{},
		&model.ScheduleTask{},
		&model.ScheduleRun{},
		// 新增模型
		&model.CartItem{},
		&model.Captcha{},
		&model.TokenBlacklist{},
		&model.UserRemark{},
		&model.SMSConfig{},
		&model.Referral{},
		&model.UserCoupon{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 创建服务
	authService := service.NewAuthService(db, &cfg.JWT)

	// M5修复：定时清理过期的token黑名单和验证码（每小时执行一次）
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			// 清理过期的token黑名单
			result := db.Where("expires_at < ?", time.Now()).Delete(&model.TokenBlacklist{})
			if result.RowsAffected > 0 {
				log.Printf("[Cleanup] Deleted %d expired token blacklist entries", result.RowsAffected)
			}
			// 清理过期的验证码
			result = db.Where("expires_at < ?", time.Now()).Delete(&model.Captcha{})
			if result.RowsAffected > 0 {
				log.Printf("[Cleanup] Deleted %d expired captcha entries", result.RowsAffected)
			}
		}
	}()

	// 初始化Gin
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS中间件（使用独立的middleware包）
	r.Use(middleware.CORS())

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
			"status":  "ok",
			"name":    "锚点财务",
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
