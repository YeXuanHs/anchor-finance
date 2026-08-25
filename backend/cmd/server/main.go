package main

import (
	"fmt"
	"log"
	"net/http"
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

	// 检查JWT密钥是否为默认弱密钥
	if cfg.JWT.Secret == "your-secret-key-change-in-production" {
		log.Println("[WARNING] JWT_SECRET is using default value! Set a strong secret in .env for production.")
	}

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
		&model.ReferralWithdrawal{},
		&service.LoginAttempt{},
		&model.TicketDeliveryRule{},
		// AI工单系统
		&model.AITicketQueue{},
		&model.AITicketMode{},
		&model.AITicketKnowledge{},
		&model.AITicketRule{},
		&model.AITicketProcessLog{},
		&model.SupplierGroupMapping{},
		// DCIM服务器管理
		&model.Server{},
		&model.DcimServer{},
		// Hook系统
		&model.HookDefinition{},
		&model.HookBinding{},
		// 产品配置选项（5张表，从zjmf搬的表结构）
		&model.ProductConfigGroup{},
		&model.ProductConfigOption{},
		&model.ProductConfigOptionSub{},
		&model.ProductConfigLink{},
		&model.ProductConfigPricing{},
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
			// 清理过期的登录风控记录
			riskSvc := service.NewLoginRiskControl(db)
			riskSvc.Cleanup()
		}
	}()

	// 启动统一定时任务引擎（合并供应商同步+AI工单+通用定时任务）
	cronEngine := service.NewCronEngine()
	cronEngine.Start()

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

	// zjmf兼容路由（根路径/v1/，让zjmf系统能对接我们）
	client.SetupZjmfCompatRoutes(r, authService)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"name":    "锚点财务",
			"en_name": "AnchorFinance",
		})
	})

	// 启动服务器（包装Gin引擎，清理双斜杠路径）
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// 清理路径中的连续斜杠（zjmf拼接URL时会产生//cart/all等路径）
		cleaned := cleanPath(req.URL.Path)
		if cleaned != req.URL.Path {
			req.URL.Path = cleaned
		}
		r.ServeHTTP(w, req)
	})
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// cleanPath 清理路径中的连续斜杠
func cleanPath(p string) string {
	if len(p) <= 1 {
		return p
	}
	result := make([]byte, 0, len(p))
	result = append(result, p[0])
	for i := 1; i < len(p); i++ {
		if p[i] == '/' && p[i-1] == '/' {
			continue
		}
		result = append(result, p[i])
	}
	return string(result)
}
