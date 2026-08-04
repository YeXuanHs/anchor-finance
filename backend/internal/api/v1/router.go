package v1

import (
	"fmt"
	"time"

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
	// 邮箱后缀白名单（注册校验）
	emailSuffixSvc := service.NewEmailSuffixWhitelistService(deps.DB, deps.Log)
	authHandler.SetEmailSuffixService(emailSuffixSvc)
	userHandler := handler.NewUserHandlerWithCaptcha(deps.UserSvc, captchaSvcForAuth, deps.Log)
	productHandler := handler.NewProductHandler(deps.ProdSvc, deps.Log)
	orderHandler := handler.NewOrderHandlerWithDB(deps.OrdSvc, deps.DB, deps.Log)
	invoiceHandler := handler.NewInvoiceHandler(deps.InvSvc, deps.Log)
	ticketHandler := handler.NewTicketHandler(deps.TicSvc, deps.Log)
	cartHandler := handler.NewCartHandler(deps.CartSvc)

	// ─── V10云 ───
	promoCodeSvcForV10Cloud := service.NewPromoCodeService(deps.DB, deps.Log)
	v10CloudSvc := service.NewV10CloudService(deps.DB, deps.Log, deps.OrdSvc, promoCodeSvcForV10Cloud)
	v10CloudHandler := handler.NewV10CloudHandler(v10CloudSvc, deps.Log)

	// ─── 验证码 ───
	captchaSvc := service.NewCaptchaService(deps.Redis, deps.DB)
	captchaHandler := handler.NewCaptchaHandler(captchaSvc, deps.DB)
	captchaConfigHandler := handler.NewCaptchaConfigHandler(captchaSvc)

	// 极验服务（通过 CaptchaHandler 集成）
	geetestSvc := service.NewGeetestService(deps.DB, deps.Log)
	captchaHandler.SetGeetestService(geetestSvc)

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
	r.GET("/geetest/register", captchaHandler.GetGeetestConfig)
	r.POST("/geetest/validate", captchaHandler.VerifyGeetest)

	// 产品
	r.GET("/products", productHandler.List)
	r.GET("/products/:id", productHandler.Get)
	r.GET("/products/:id/pricing", productHandler.GetPricing)

	// 用户公开信息（脱敏）
	frontendUserHandler := handler.NewFrontendUserHandler(deps.UserSvc, deps.Log)
	r.GET("/users/:id", frontendUserHandler.GetPublicProfile)

	// 支付方式
	balanceSvc := service.NewBalanceLogService(deps.DB)
	balanceHandler := handler.NewBalanceHandler(balanceSvc, deps.DB)
	r.GET("/payment-methods", balanceHandler.GetEnabledGateways)
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

	// 用户中心菜单（公开接口）
	userMenuHandler := handler.NewUserMenuHandler(deps.DB)
	r.GET("/user/menus", userMenuHandler.GetUserMenus)
	r.GET("/nav/top", userMenuHandler.GetTopNav)
	r.GET("/nav/bottom", userMenuHandler.GetBottomNav)

	// 公共信息
	publicSvc := service.NewPublicService(deps.DB, deps.Log)
	publicHandler := handler.NewPublicHandler(publicSvc, deps.Log)
	r.GET("/homepage/base-info", publicHandler.GetHomepageBaseInfo)
	r.GET("/downloads", publicHandler.GetUserDownloads)

	// 知识库（公开接口）
	knowledgeBaseSvc := service.NewKnowledgeBaseService(deps.DB, deps.Log)
	knowledgeBaseHandler := handler.NewKnowledgeBaseHandler(knowledgeBaseSvc, deps.Log)
	r.GET("/knowledge/categories", knowledgeBaseHandler.ListCategories)
	r.GET("/knowledge/articles", knowledgeBaseHandler.ListArticles)
	r.GET("/knowledge/articles/:id", knowledgeBaseHandler.GetArticle)

	// 公告/新闻（公开接口）
	newsHandler := handler.NewNewsHandler(deps.DB, deps.Log)
	r.GET("/news", newsHandler.GetPublishedList)
	r.GET("/news/:id", newsHandler.GetPublishedDetail)

	// 语言包（公开接口）
	langSvc := service.NewLanguageService(deps.DB, deps.Log)
	langHandler := handler.NewLanguageHandler(langSvc, deps.Log)
	r.GET("/languages", langHandler.GetActiveLanguages)
	r.GET("/languages/:code/translations", langHandler.GetTranslations)

	// 优惠码验证
	promoHandler := handler.NewPromoCodeHandler(deps.DB, deps.Log)
	r.GET("/promo-codes/validate", promoHandler.Validate)

	// 需要认证的路由
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		// 用户信息
		auth.GET("/user/profile", userHandler.GetProfile)
		auth.PUT("/user/profile", userHandler.UpdateProfile)
		auth.POST("/user/change-password", authHandler.ChangePassword)

		// 二步验证 (2FA)
		auth.GET("/user/2fa", userHandler.Get2FAStatus)
		auth.POST("/user/2fa/enable", userHandler.Enable2FA)
		auth.POST("/user/2fa/verify", userHandler.Verify2FA)
		auth.POST("/user/2fa/disable", userHandler.Disable2FA)

		// API密钥管理
		auth.GET("/user/api-keys", userHandler.GetAPIKeys)
		auth.POST("/user/api-keys", userHandler.CreateAPIKey)
		auth.PUT("/user/api-keys/:id/toggle", userHandler.ToggleAPIKey)
		auth.DELETE("/user/api-keys/:id", userHandler.DeleteAPIKey)

		// 用户登录日志（脱敏IP）
		auth.GET("/user/login-logs", frontendUserHandler.GetMyLoginLogs)

		// 用户偏好设置
		auth.GET("/user/tastes", publicHandler.GetUserTastes)
		auth.PUT("/user/tastes", publicHandler.SaveUserTastes)

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

		// 发票申请（Voucher）
		voucherSvc := service.NewVoucherService(deps.DB, deps.Log)
		voucherHandler := handler.NewVoucherHandler(voucherSvc, deps.Log)
		auth.GET("/user/vouchers", voucherHandler.GetUserVoucherList)
		auth.POST("/user/vouchers", voucherHandler.CreateUserVoucher)
		// 发票抬头
		auth.GET("/user/voucher-types", voucherHandler.GetVoucherTypes)
		auth.POST("/user/voucher-types", voucherHandler.CreateVoucherType)
		auth.PUT("/user/voucher-types/:id", voucherHandler.UpdateVoucherType)
		auth.DELETE("/user/voucher-types/:id", voucherHandler.DeleteVoucherType)
		// 收件地址
		auth.GET("/user/voucher-posts", voucherHandler.GetVoucherPosts)
		auth.POST("/user/voucher-posts", voucherHandler.CreateVoucherPost)
		auth.PUT("/user/voucher-posts/:id", voucherHandler.UpdateVoucherPost)
		auth.DELETE("/user/voucher-posts/:id", voucherHandler.DeleteVoucherPost)

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

		// V10云
		auth.POST("/v10cloud/order", v10CloudHandler.CreateOrder)
		auth.GET("/v10cloud/products", v10CloudHandler.GetProducts)
		auth.GET("/v10cloud/config-options", v10CloudHandler.GetConfigOptions)
		auth.POST("/v10cloud/price", v10CloudHandler.CalculatePrice)

		// AI 购物助手会话
		aiShopSvc := service.NewAIShoppingCoreService(deps.DB, deps.Log)
		aiShopHandler := handler.NewAIShoppingCoreHandler(aiShopSvc, deps.Log)
		auth.POST("/ai-shopping/session", func(c *gin.Context) {
			sessionID := fmt.Sprintf("shop_%d_%d", c.GetUint("user_id"), time.Now().UnixNano())
			c.JSON(200, gin.H{"code": 0, "data": gin.H{"session_id": sessionID}})
		})
		auth.POST("/ai-shopping/session/:session_id/message", aiShopHandler.Chat)
		auth.GET("/ai-shopping/session/:session_id/messages", aiShopHandler.GetChatHistory)

		// ─── 交易市场 ───
		marketplaceSvc := service.NewMarketplaceService(deps.DB, deps.Log)
		marketplaceHandler := handler.NewMarketplaceHandler(marketplaceSvc, deps.Log)

		// 挂售管理
		auth.POST("/marketplace/listings", marketplaceHandler.CreateListing)
		auth.PUT("/marketplace/listings/:id", marketplaceHandler.UpdateListing)
		auth.DELETE("/marketplace/listings/:id", marketplaceHandler.RemoveListing)
		auth.GET("/marketplace/listings/mine", marketplaceHandler.GetUserListings)

		// 订单管理
		auth.POST("/marketplace/orders", marketplaceHandler.CreateOrder)
		auth.GET("/marketplace/orders/buyer", marketplaceHandler.GetBuyerOrders)
		auth.GET("/marketplace/orders/seller", marketplaceHandler.GetSellerOrders)
		auth.POST("/marketplace/orders/:id/complete", marketplaceHandler.CompleteOrder)
		auth.POST("/marketplace/orders/:id/cancel", marketplaceHandler.CancelOrder)

		// 私聊功能
		auth.POST("/marketplace/messages", marketplaceHandler.SendMessage)
		auth.GET("/marketplace/messages/:listing_id/:user_id", marketplaceHandler.GetChatMessages)
		auth.GET("/marketplace/chat-sessions", marketplaceHandler.GetChatSessions)
		auth.GET("/marketplace/unread-count", marketplaceHandler.GetUnreadCount)
	}

	// 公开接口（不需要登录）
	r.GET("/marketplace/listings", marketplaceHandler.GetListings)
	r.GET("/marketplace/listings/:id", marketplaceHandler.GetListing)
}
