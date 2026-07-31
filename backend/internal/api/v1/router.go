package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"anchorfinance/internal/handler"
	"anchorfinance/internal/api/middleware"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/auth"
	"anchorfinance/pkg/logger"
	"gorm.io/gorm"
)

// Deps holds shared dependencies for user-facing API route registration.
type Deps struct {
	DB         *gorm.DB
	Redis      *redis.Client
	Log        *logger.Logger
	JWTKey     string
	JWTManager *auth.JWTManager
	UserSvc    *service.UserService
	ProdSvc    *service.ProductService
	OrdSvc     *service.OrderService
	InvSvc     *service.InvoiceService
	TicSvc     *service.TicketService
	CartSvc    *service.CartService
	OAuthSvc   *service.OAuthService
}

// RegisterRoutes registers all user-facing v1 API routes.
func RegisterRoutes(r *gin.RouterGroup, deps Deps) {
	// ─── 核心 handlers ───
	captchaSvcForAuth := service.NewCaptchaService(deps.Redis, deps.DB)
	authHandler := handler.NewAuthHandlerWithCaptcha(deps.UserSvc, captchaSvcForAuth, deps.Log, deps.JWTManager)
	userHandler := handler.NewUserHandlerWithCaptcha(deps.UserSvc, captchaSvcForAuth, deps.Log)
	productHandler := handler.NewProductHandler(deps.ProdSvc, deps.Log)
	orderHandler := handler.NewOrderHandlerWithDB(deps.OrdSvc, deps.DB, deps.Log)
	invoiceHandler := handler.NewInvoiceHandler(deps.InvSvc, deps.Log)
	ticketHandler := handler.NewTicketHandler(deps.TicSvc, deps.Log)
	cartHandler := handler.NewCartHandler(deps.CartSvc)

	// ─── V10云 ───
	couponSvcForV10Cloud := service.NewCouponService(deps.DB, deps.Log)
	v10CloudSvc := service.NewV10CloudService(deps.DB, deps.Log, deps.OrdSvc, couponSvcForV10Cloud)
	v10CloudHandler := handler.NewV10CloudHandler(v10CloudSvc, deps.Log)

	// ─── 验证码 ───
	captchaSvc := service.NewCaptchaService(deps.Redis, deps.DB)
	captchaHandler := handler.NewCaptchaHandler(captchaSvc, deps.DB)
	captchaConfigHandler := handler.NewCaptchaConfigHandler(captchaSvc)

	// 初始化极验服务（如果配置了）
	geetestHandler := handler.NewGeetestHandler(captchaSvc)

	// ─── 公共路由 ───
	r.POST("/login", authHandler.Login)
	r.POST("/login/sms", authHandler.SMSLogin)
	r.POST("/register", authHandler.Register)
	r.POST("/password/verify-code", authHandler.VerifyResetCode)
	r.POST("/password/reset", authHandler.ResetPassword)

	// 验证码
	r.GET("/captcha/config", captchaConfigHandler.GetConfig)
	r.GET("/captcha/generate", captchaHandler.Generate)
	r.POST("/captcha/verify", captchaHandler.Verify)
	r.POST("/captcha/check", captchaHandler.Check)

	// 极验4.0
	r.GET("/geetest/register", geetestHandler.Register)
	r.POST("/geetest/validate", geetestHandler.Validate)

	// 产品
	r.GET("/products", productHandler.List)
	r.GET("/products/:id", productHandler.Get)
	r.GET("/products/:id/pricing", productHandler.GetPricing)

	// 用户公开信息（脱敏）
	frontendUserHandler := handler.NewFrontendUserHandler(deps.UserSvc, deps.Log)
	r.GET("/users/:id", frontendUserHandler.GetPublicProfile)

	// 通用
	r.GET("/payment-methods", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "data": []interface{}{}})
	})
	r.GET("/system/lang", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "data": gin.H{"lang": "zh-CN"}})
	})

	// 公开系统设置
	r.GET("/settings/public", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"data": gin.H{
				"site_name":     "锚点财务",
				"site_url":      "",
				"allow_register": true,
			},
		})
	})

	// 需要认证的路由
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		// 用户信息
		auth.GET("/user/profile", userHandler.GetProfile)
		auth.PUT("/user/profile", userHandler.UpdateProfile)
		auth.POST("/user/change-password", authHandler.ChangePassword)

		// 用户登录日志（脱敏IP）
		auth.GET("/user/login-logs", frontendUserHandler.GetMyLoginLogs)

	// OAuth
	oauthHandler := handler.NewOAuthHandler(deps.OAuthSvc, deps.Log, deps.JWTManager)
	auth.GET("/oauth/:provider", oauthHandler.Redirect)
	auth.GET("/oauth/:provider/callback", oauthHandler.Callback)
	auth.POST("/oauth/unbind", oauthHandler.Unbind)
	auth.GET("/oauth/accounts", oauthHandler.GetBoundAccounts)

		// 产品
		auth.GET("/products", productHandler.List)
		auth.GET("/products/:id", productHandler.Get)
		auth.GET("/products/:id/pricing", productHandler.GetPricing)

		// 订单
		auth.POST("/orders", orderHandler.Create)
		auth.GET("/orders", orderHandler.List)
		auth.GET("/orders/:id", orderHandler.Get)

		// 发票
		auth.GET("/invoices", invoiceHandler.List)
		auth.GET("/invoices/:id", invoiceHandler.Get)
		auth.POST("/invoices/:id/pay", invoiceHandler.Pay)

		// 工单
		auth.POST("/tickets", ticketHandler.Create)
		auth.GET("/tickets", ticketHandler.List)
		auth.GET("/tickets/:id", ticketHandler.Get)
		auth.POST("/tickets/:id/reply", ticketHandler.Reply)
		auth.POST("/tickets/:id/close", ticketHandler.Close)

		// 购物车
		auth.POST("/cart/add", cartHandler.Add)
		auth.GET("/cart", cartHandler.Get)
		auth.PUT("/cart/:id", cartHandler.Update)
		auth.DELETE("/cart/:id", cartHandler.Remove)
		auth.POST("/cart/checkout", cartHandler.Checkout)
		auth.POST("/cart/coupon", cartHandler.ApplyCoupon)

		// V10云
		auth.POST("/v10cloud/order", v10CloudHandler.CreateOrder)
		auth.GET("/v10cloud/products", v10CloudHandler.GetProducts)
		auth.GET("/v10cloud/config-options", v10CloudHandler.GetConfigOptions)
		auth.POST("/v10cloud/price", v10CloudHandler.CalculatePrice)
	}
}
