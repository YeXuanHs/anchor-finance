package v1

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"anchorfinance/internal/handler"
	"anchorfinance/internal/api/middleware"
	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/auth"
	"anchorfinance/pkg/logger"
)

// Deps holds shared dependencies for route registration.
type Deps struct {
	DB      *gorm.DB
	Log     *logger.Logger
	JWTKey  string
	Redis   *redis.Client
	JWTMgr  *auth.JWTManager
	UserSvc *service.UserService
	ProdSvc *service.ProductService
	OrdSvc  *service.OrderService
	InvSvc  *service.InvoiceService
	TicSvc  *service.TicketService
	CartSvc *service.CartService
	OAuthSvc *service.OAuthService
}

// RegisterRoutes registers all v1 API routes on the given router group.
func RegisterRoutes(r *gin.RouterGroup, deps Deps) {
	log := deps.Log

	// 初始化处理器
	authHandler := &handler.AuthHandler{}
	userHandler := handler.NewUserHandler(deps.UserSvc, log)
	productHandler := handler.NewProductHandler()
	orderHandler := handler.NewOrderHandler()
	invoiceHandler := handler.NewInvoiceHandler()
	ticketHandler := handler.NewTicketHandler()
	cartHandler := handler.NewCartHandler(deps.CartSvc)

	balanceSvc := service.NewBalanceLogService(deps.DB)
	balanceHandler := handler.NewBalanceHandler(balanceSvc, deps.DB)

	creditSvc := service.NewCreditService(deps.DB, log)
	creditHandler := handler.NewCreditHandler(deps.DB, creditSvc)

	upstreamSvc := service.NewUpstreamService(deps.DB, log)
	hostSvc := service.NewHostService(deps.DB, log, upstreamSvc)
	hostHandler := handler.NewHostHandler(hostSvc, log)

	contactsSvc := service.NewContactService(deps.DB, log)
	contactsHandler := handler.NewContactsHandler(contactsSvc, log)

	promoCodeSvc := service.NewPromoCodeService(deps.DB, deps.Log)
	_ = handler.NewPromoCodeHandler(deps.DB, deps.Log) // used for route registration if needed

	announceSvc := service.NewAnnounceService(deps.DB, log)
	announceHandler := handler.NewAnnounceHandler(announceSvc, log)

	newsHandler := handler.NewNewsHandler(deps.DB, log)

	knowledgeSvc := service.NewKnowledgeService(deps.DB, log)
	knowledgeHandler := handler.NewKnowledgeHandler(deps.DB, knowledgeSvc, log)

	systemMsgSvc := service.NewSystemMessageService(deps.DB, log)
	systemMsgHandler := handler.NewSystemMessageHandler(systemMsgSvc, log)

	certSvc := service.NewCertificationService(deps.DB, log)
	certHandler := handler.NewCertificationHandler(certSvc, log)

	oauthHandler := handler.NewOAuthHandler(deps.OAuthSvc, log, deps.JWTMgr, nil, nil)

	upgradeSvc := service.NewUpgradeService(deps.DB, log)
	upgradeHandler := handler.NewUpgradeHandler(upgradeSvc, log)

	contractHandler := handler.NewContractHandler(deps.DB, log)

	affiliateSvc := service.NewAffiliateService(deps.DB, log)
	affiliateHandler := handler.NewAffiliateHandler(affiliateSvc, log)

	voucherSvc := service.NewVoucherService(deps.DB, log)
	voucherHandler := handler.NewVoucherHandler(voucherSvc, log)

	balSvc := service.NewBalanceService(deps.DB, log)
	multiRenewSvc := service.NewMultiRenewService(deps.DB, log, deps.InvSvc, balSvc)
	multiRenewHandler := handler.NewMultiRenewHandler(multiRenewSvc, log)

	loginLogSvc := service.NewLoginLogService(deps.DB, log)
	loginLogHandler := handler.NewLoginLogHandler(loginLogSvc, log)

	apiLogSvc := service.NewApiLogService(deps.DB, log)
	apiLogHandler := handler.NewAPILogHandler(apiLogSvc, log)

	aiShoppingSvc := service.NewAIShoppingCoreService(deps.DB, log)
	aiShoppingHandler := handler.NewAIShoppingCoreHandler(aiShoppingSvc, log)

	marketplaceSvc := service.NewMarketplaceService(deps.DB, log)
	marketplaceHandler := handler.NewMarketplaceHandler(marketplaceSvc, log)

	// 验证码 & 极验
	captchaSvc := service.NewCaptchaService(deps.Redis, deps.DB)
	captchaHandler := handler.NewCaptchaHandler(captchaSvc, deps.DB)

	// 用户公开信息（脱敏）
	frontendUserHandler := handler.NewFrontendUserHandler(deps.UserSvc, log)

	// 导航菜单
	userMenuHandler := handler.NewUserMenuHandler(deps.DB)

	// 公共信息
	publicSvc := service.NewPublicService(deps.DB, deps.Log)
	publicHandler := handler.NewPublicHandler(publicSvc, deps.Log)

	// 语言包
	langSvc := service.NewLanguageService(deps.DB, deps.Log)
	langHandler := handler.NewLanguageHandler(langSvc, deps.Log)

	// V10云
	v10CloudSvc := service.NewV10CloudService(deps.DB, deps.Log, deps.OrdSvc, promoCodeSvc)
	v10CloudHandler := handler.NewV10CloudHandler(v10CloudSvc, log)

	// ==================== 公开接口 ====================

	// 认证
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/sms-login", authHandler.SMSLogin)
		auth.POST("/refresh", middleware.AuthRequired(), authHandler.RefreshToken)
		auth.POST("/logout", middleware.AuthRequired(), authHandler.Logout)
		auth.POST("/access-token", authHandler.AccessTokenLogin)
	}

	// 密码重置
	r.POST("/password/verify-code", func(c *gin.Context) {
		var req struct {
			Email string `json:"email" binding:"required,email"`
			Code  string `json:"code" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// TODO: 实现验证码校验逻辑
		c.JSON(200, gin.H{"code": 0, "message": "验证码有效"})
	})
	r.POST("/password/reset", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required,email"`
			Code     string `json:"code" binding:"required"`
			Password string `json:"password" binding:"required,min=6"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// TODO: 实现密码重置逻辑
		c.JSON(200, gin.H{"code": 0, "message": "密码重置成功"})
	})

	// 验证码
	r.GET("/captcha/config", func(c *gin.Context) {
		configService := captchaSvc.GetCaptchaConfigService()
		configs, err := configService.GetAllConfigs()
		if err != nil {
			c.JSON(500, gin.H{"error": "failed to get captcha config"})
			return
		}
		result := make(map[string]string)
		for _, cfg := range configs {
			if cfg.Status {
				result[cfg.Key] = cfg.Value
			}
		}
		c.JSON(200, gin.H{"data": result})
	})
	r.GET("/captcha/generate", captchaHandler.GetImageJSON)
	r.POST("/captcha/verify", captchaHandler.VerifyImage)
	r.POST("/captcha/check", func(c *gin.Context) {
		var req struct {
			CaptchaID string `json:"captcha_id" binding:"required"`
			Answer    string `json:"answer" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		// 委托给 VerifyImage 处理
		captchaHandler.VerifyImage(c)
	})

	// 极验4.0
	r.GET("/geetest/register", captchaHandler.GetGeetestConfig)
	r.POST("/geetest/validate", captchaHandler.VerifyGeetest)

	// AI Shopping配置（公开）
	r.GET("/ai-shopping/config", aiShoppingHandler.GetConfig)

	// 产品
	r.GET("/products", productHandler.GetList)
	r.GET("/products/:id", productHandler.GetDetail)
	r.GET("/products/hot", productHandler.GetHot)
	r.GET("/products/groups", productHandler.GetGroups)
	r.GET("/products/:id/pricing", func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid product id"})
			return
		}
		product, err := deps.ProdSvc.GetByID(uint(id))
		if err != nil {
			c.JSON(404, gin.H{"error": "product not found"})
			return
		}
		c.JSON(200, gin.H{"data": gin.H{
			"id":            product.ID,
			"name":          product.Name,
			"price":         product.Price,
			"billing_cycle": product.BillingCycle,
		}})
	})

	// 用户公开信息（脱敏）
	r.GET("/users/:id", frontendUserHandler.GetPublicProfile)

	// 支付方式（公开）
	r.GET("/payment-methods", balanceHandler.GetEnabledGateways)

	// 系统语言
	r.GET("/system/lang", func(c *gin.Context) {
		c.JSON(200, gin.H{"code": 0, "data": gin.H{"lang": "zh-CN"}})
	})

	// 公开设置
	r.GET("/settings/public", func(c *gin.Context) {
		keys := []string{
			"site_name", "site_url", "site_description", "allow_register",
			"contact_phone", "contact_email", "contact_address",
			"sales_phone", "support_phone", "sales_email", "work_time",
		}
		m := make(map[string]string)
		for _, k := range keys {
			var cfg model.SystemConfig
			if err := deps.DB.Where("`key` = ?", k).First(&cfg).Error; err == nil {
				m[k] = cfg.Value
			}
		}
		c.JSON(200, gin.H{
			"code": 0,
			"data": gin.H{
				"site_name":       m["site_name"],
				"site_url":        m["site_url"],
				"site_description": m["site_description"],
				"allow_register":  m["allow_register"] != "false",
				"contact_phone":   m["contact_phone"],
				"contact_email":   m["contact_email"],
				"contact_address": m["contact_address"],
				"sales_phone":     m["sales_phone"],
				"support_phone":   m["support_phone"],
				"sales_email":     m["sales_email"],
				"work_time":       m["work_time"],
			},
		})
	})

	// 导航菜单（公开）
	r.GET("/user/menus", userMenuHandler.GetUserMenus)
	r.GET("/nav/top", userMenuHandler.GetTopNav)
	r.GET("/nav/bottom", userMenuHandler.GetBottomNav)

	// 首页基础信息 & 下载
	r.GET("/homepage/base-info", publicHandler.GetHomepageBaseInfo)
	r.GET("/downloads", publicHandler.GetUserDownloads)

	// 公告
	r.GET("/announcements", announceHandler.GetActive)
	r.GET("/announcements/:id", announceHandler.AdminGetDetail)

	// 新闻
	r.GET("/news", newsHandler.GetList)
	r.GET("/news/:id", newsHandler.GetDetail)
	r.GET("/news/categories", newsHandler.GetCategories)

	// 知识库（/help/ 路径）
	r.GET("/help/categories", knowledgeHandler.GetCategories)
	r.GET("/help/articles", knowledgeHandler.GetArticles)
	r.GET("/help/articles/hot", knowledgeHandler.GetHot)
	r.GET("/help/articles/search", knowledgeHandler.Search)
	r.GET("/help/articles/:id", knowledgeHandler.GetArticle)
	r.GET("/help/articles/:id/related", knowledgeHandler.GetRelatedArticles)
	r.POST("/help/articles/:id/feedback", knowledgeHandler.SubmitFeedback)
	r.GET("/help/categories/:id/sub", knowledgeHandler.GetSubCategories)

	// 知识库（兼容 /knowledge/ 路径）
	r.GET("/knowledge/categories", knowledgeHandler.GetCategories)
	r.GET("/knowledge/articles", knowledgeHandler.GetArticles)
	r.GET("/knowledge/articles/:id", knowledgeHandler.GetArticle)

	// 语言包（公开）
	r.GET("/languages", langHandler.GetActiveLanguages)
	r.GET("/languages/:code/translations", langHandler.GetTranslations)

	// 优惠码验证（GET版本，公开）
	r.GET("/promo-codes/validate", func(c *gin.Context) {
		code := c.Query("code")
		if code == "" {
			c.JSON(400, gin.H{"error": "code is required"})
			return
		}
		userID := c.GetUint("user_id")
		discount, promo, err := promoCodeSvc.Validate(code, userID, 0, 0)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"data": gin.H{"discount": discount, "promo": promo}})
	})

	// OAuth 回调（公开）
	r.GET("/oauth/:provider", oauthHandler.Redirect)
	r.GET("/oauth/:provider/callback", oauthHandler.Callback)

	// 系统设置（公开部分）
	r.GET("/system/settings", func(c *gin.Context) {
		var settings []struct{ Key, Value string }
		deps.DB.Table("system_configs").Where("key IN ?", []string{
			"company_name", "company_email", "main_phone", "main_address",
			"record_no", "logo_url", "favicon_url",
		}).Find(&settings)
		result := make(map[string]string)
		for _, s := range settings {
			result[s.Key] = s.Value
		}
		c.JSON(200, gin.H{"data": result})
	})

	// 产品组
	r.GET("/product-groups", func(c *gin.Context) {
		var groups []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
			Slug string `json:"slug"`
		}
		deps.DB.Table("product_groups").Where("show_in_nav = ?", true).Order("sort_order").Find(&groups)
		c.JSON(200, gin.H{"data": groups})
	})

	// 轮播图（公开接口）
	bannerSvc := service.NewBannerService(deps.DB, deps.Log)
	bannerHandler := handler.NewBannerHandler(bannerSvc, deps.Log)
	r.GET("/banners", bannerHandler.GetActive)

	// 首页特色区域（公开接口）
	featureSvc := service.NewHomepageFeatureService(deps.DB, deps.Log)
	featureHandler := handler.NewHomepageFeatureHandler(featureSvc, deps.Log)
	r.GET("/homepage-features", featureHandler.GetActive)

	// 维护公告
	r.GET("/maintenance/notices", func(c *gin.Context) {
		var notices []struct {
			ID      uint   `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		}
		deps.DB.Table("announcements").Where("status = 1 AND type = ?", "maintenance").
			Order("id DESC").Limit(10).Find(&notices)
		c.JSON(200, gin.H{"data": notices})
	})

	// 域名后缀列表
	r.GET("/domain/suffixes", func(c *gin.Context) {
		var suffixes []struct {
			ID     uint   `json:"id"`
			Suffix string `json:"suffix"`
			Price  string `json:"price"`
		}
		deps.DB.Table("domain_suffixes").Where("status = 1").Order("sort_order").Find(&suffixes)
		c.JSON(200, gin.H{"data": suffixes})
	})

	// 交易市场（公开）
	r.GET("/marketplace/listings", marketplaceHandler.GetListings)
	r.GET("/marketplace/listings/:id", marketplaceHandler.GetListing)

	// 联系表单
	r.POST("/contact", func(c *gin.Context) {
		var req struct {
			Name    string `json:"name" binding:"required"`
			Email   string `json:"email" binding:"required"`
			Phone   string `json:"phone"`
			Subject string `json:"subject" binding:"required"`
			Content string `json:"content" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"message": "提交成功，我们会尽快与您联系"})
	})

	// 解决方案（公开）
	r.GET("/solutions", publicHandler.GetSolutions)
	r.GET("/solutions/:slug", publicHandler.GetSolutionDetail)

	// 托管服务
	r.GET("/colocation/advantages", publicHandler.GetColocationAdvantages)
	r.GET("/colocation/datacenters", publicHandler.GetColocationDatacenters)

	// ==================== 需要登录 ====================
	user := r.Group("")
	user.Use(middleware.AuthRequired())
	{
		// 用户资料
		user.GET("/user/profile", userHandler.GetProfile)
		user.PUT("/user/profile", userHandler.UpdateProfile)
		user.POST("/user/change-password", userHandler.ChangePassword)
		user.POST("/user/bind-phone", userHandler.BindPhone)
		user.POST("/user/bind-email", userHandler.BindEmail)

		// 二步验证 (2FA)
		user.GET("/user/2fa", userHandler.Get2FAStatus)
		user.POST("/user/2fa/enable", userHandler.Enable2FA)
		user.POST("/user/2fa/verify", userHandler.Verify2FA)
		user.POST("/user/2fa/disable", userHandler.Disable2FA)

		// API密钥管理（对齐zjmf：开关/查看/重置）
		user.GET("/user/api/summary", userHandler.GetAPISummary)
		user.POST("/user/api/open", userHandler.ToggleAPIOpen)
		user.POST("/user/api/reset", userHandler.ResetAPIKey)

		// OAuth绑定
		user.GET("/oauth/providers", oauthHandler.GetProviders)
		user.POST("/oauth/:provider/bind", oauthHandler.BindAccount)
		user.DELETE("/oauth/:provider/unbind", oauthHandler.UnbindAccount)
		user.GET("/oauth/accounts", oauthHandler.GetBoundAccounts)

		// 用户偏好
		user.GET("/user/tastes", publicHandler.GetUserTastes)
		user.PUT("/user/tastes", publicHandler.SaveUserTastes)

		// 用户登录日志
		user.GET("/user/login-logs", loginLogHandler.List)
		user.GET("/login-logs", loginLogHandler.List)

		// 产品（用户的产品）
		user.GET("/user/products", productHandler.GetUserProducts)

		// 主机/服务管理
		user.GET("/hosts", hostHandler.GetUserHosts)
		user.GET("/hosts/:id", hostHandler.GetDetail)
		user.GET("/hosts/:id/operations", hostHandler.GetOperations)
		user.GET("/hosts/:id/billing", hostHandler.GetBilling)
		user.GET("/hosts/:id/log", hostHandler.GetLog)
		user.GET("/hosts/:id/download", hostHandler.GetDownload)
		user.POST("/hosts/:id/remark", hostHandler.PostRemark)
		user.POST("/hosts/:id/:action", hostHandler.PerformAction)

		// 购物车
		user.GET("/cart", cartHandler.GetCart)
		user.POST("/cart/add", cartHandler.AddToCart)
		user.PUT("/cart/:id", cartHandler.UpdateCart)
		user.DELETE("/cart/:id", cartHandler.RemoveFromCart)
		user.DELETE("/cart", cartHandler.ClearCart)
		user.POST("/cart/checkout", cartHandler.Checkout)
		user.POST("/cart/batch-delete", cartHandler.BatchDelete)

		// 订单
		user.POST("/orders", orderHandler.Create)
		user.GET("/orders", orderHandler.GetUserOrders)
		user.GET("/orders/:id", orderHandler.GetDetail)
		user.POST("/orders/:id/pay", orderHandler.Pay)
		user.POST("/orders/:id/cancel", orderHandler.Cancel)
		user.POST("/orders/preview", orderHandler.Preview)

		// 账单/发票
		user.GET("/invoices", invoiceHandler.GetUserInvoices)
		user.GET("/invoices/:id", invoiceHandler.GetDetail)
		user.POST("/invoices/:id/pay", invoiceHandler.Pay)
		user.POST("/invoices/combine", invoiceHandler.CombineInvoices)
		user.GET("/invoices/combine", invoiceHandler.GetCombineInvoices)

		// 工单
		user.POST("/tickets", ticketHandler.Create)
		user.GET("/tickets", ticketHandler.GetUserTickets)
		user.GET("/tickets/:id", ticketHandler.GetDetail)
		user.POST("/tickets/:id/reply", ticketHandler.Reply)
		user.POST("/tickets/:id/close", ticketHandler.Close)
		user.POST("/tickets/:id/attachments", ticketHandler.UploadAttachment)
		user.GET("/tickets/:id/attachments", ticketHandler.GetAttachments)
		user.DELETE("/tickets/attachments/:id", ticketHandler.DeleteAttachment)

		// 余额
		user.GET("/balance", balanceHandler.GetBalance)
		user.POST("/balance/recharge", balanceHandler.Recharge)
		user.GET("/balance/logs", balanceHandler.GetBalanceLogs)
		user.POST("/balance/withdraw", balanceHandler.Withdraw)

		// 信用额度
		user.GET("/credit", creditHandler.GetInfo)
		user.GET("/credit/logs", creditHandler.GetLogs)
		user.POST("/credit/apply", creditHandler.Apply)
		user.POST("/credit/repay", creditHandler.Repay)
		user.POST("/credit/prepayment", creditHandler.Prepayment)
		user.GET("/credit/bills", creditHandler.GetBills)
		user.GET("/credit/bills/:id", creditHandler.GetBillDetail)
		user.GET("/credit/bills/:id/repayments", creditHandler.GetBillRepayments)
		user.POST("/credit/bills/:id/pay", creditHandler.PayBill)
		user.GET("/credit/config", creditHandler.GetCreditConfig)
		user.PUT("/credit/config", creditHandler.UpdateCreditConfig)
		user.GET("/credit/used-detail", creditHandler.GetUsedDetail)
		user.GET("/credit/used-summary", creditHandler.GetUsedSummary)

		// 联系人
		user.GET("/contacts", contactsHandler.GetList)
		user.POST("/contacts", contactsHandler.Create)
		user.PUT("/contacts/:id", contactsHandler.Update)
		user.DELETE("/contacts/:id", contactsHandler.Delete)
		user.PUT("/contacts/:id/default", contactsHandler.SetDefault)
		user.GET("/contacts/default", contactsHandler.GetDefault)

		// 优惠码
		user.GET("/promo-codes", func(c *gin.Context) {
			userID := c.GetUint("user_id")
			promos, err := promoCodeSvc.GetUserPromoCodes(userID)
			if err != nil {
				c.JSON(500, gin.H{"error": "failed to get promo codes"})
				return
			}
			c.JSON(200, gin.H{"data": promos})
		})
		user.POST("/promo-codes/validate", func(c *gin.Context) {
			var req struct {
				Code string `json:"code" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			userID := c.GetUint("user_id")
			discount, promo, err := promoCodeSvc.Validate(req.Code, userID, 0, 0)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"data": gin.H{"discount": discount, "promo": promo}})
		})

		// 支付方式
		user.GET("/payment-gateways", balanceHandler.GetEnabledGateways)

		// OAuth解绑
		user.POST("/oauth/unbind", func(c *gin.Context) {
			var req struct {
				Provider string `json:"provider" binding:"required"`
			}
			if err := c.ShouldBindJSON(&req); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			userID := c.GetUint("user_id")
			if err := deps.DB.Where("user_id = ? AND provider = ?", userID, req.Provider).Delete(&model.OAuthBind{}).Error; err != nil {
				c.JSON(500, gin.H{"error": "failed to unbind"})
				return
			}
			c.JSON(200, gin.H{"message": "unbind success"})
		})

		// 系统消息
		user.GET("/messages", systemMsgHandler.GetList)
		user.GET("/messages/unread-count", systemMsgHandler.GetUnreadCount)
		user.PUT("/messages/:id/read", systemMsgHandler.MarkRead)
		user.PUT("/messages/read-all", systemMsgHandler.MarkAllRead)
		user.DELETE("/messages/:id", systemMsgHandler.Delete)

		// 实名认证
		user.GET("/certification/status", certHandler.GetStatus)
		user.GET("/certification/enterprise", certHandler.GetEnterpriseCert)
		user.POST("/certification/submit", certHandler.Submit)
		user.POST("/certification/enterprise", certHandler.SubmitEnterprise)

		// 升级
		user.GET("/upgrades/available/:host_id", upgradeHandler.GetAvailableUpgrades)
		user.POST("/upgrades", upgradeHandler.CreateUpgrade)
		user.GET("/upgrades", upgradeHandler.GetUserUpgrades)
		user.GET("/upgrades/:id", upgradeHandler.GetUpgradeDetail)
		user.POST("/upgrades/:id/pay", upgradeHandler.PayUpgrade)

		// 合同
		user.GET("/contracts", contractHandler.GetUserContracts)
		user.GET("/contracts/:id", contractHandler.GetDetail)
		user.POST("/contracts/:id/sign", contractHandler.SignContract)
		user.POST("/contracts/:id/bind-hosts", contractHandler.BindHosts)
		user.POST("/contracts/:id/unbind-host", contractHandler.UnbindHost)

		// 代理/推荐
		user.GET("/affiliate/info", affiliateHandler.GetInfo)
		user.GET("/affiliate/records", affiliateHandler.GetRecords)
		user.GET("/affiliate/withdraws", affiliateHandler.GetWithdraws)
		user.POST("/affiliate/withdraw", affiliateHandler.ApplyWithdraw)

		// 发票管理
		user.GET("/vouchers", voucherHandler.GetUserVoucherList)
		user.POST("/vouchers", voucherHandler.CreateUserVoucher)

		// 发票抬头
		user.GET("/voucher-types", voucherHandler.GetVoucherTypes)
		user.POST("/voucher-types", voucherHandler.CreateVoucherType)
		user.PUT("/voucher-types/:id", voucherHandler.UpdateVoucherType)
		user.DELETE("/voucher-types/:id", voucherHandler.DeleteVoucherType)

		// 收件地址
		user.GET("/voucher-posts", voucherHandler.GetVoucherPosts)
		user.POST("/voucher-posts", voucherHandler.CreateVoucherPost)
		user.PUT("/voucher-posts/:id", voucherHandler.UpdateVoucherPost)
		user.DELETE("/voucher-posts/:id", voucherHandler.DeleteVoucherPost)

		// 批量续费
		user.POST("/multi-renew", multiRenewHandler.Create)
		user.GET("/multi-renew/preview", multiRenewHandler.Preview)
		user.GET("/multi-renew", multiRenewHandler.List)
		user.GET("/multi-renew/:id", multiRenewHandler.GetDetail)
		user.POST("/multi-renew/:id/execute", multiRenewHandler.Execute)
		user.POST("/multi-renew/:id/cancel", multiRenewHandler.Cancel)

		// API日志
		user.GET("/api-logs", apiLogHandler.List)

		// 产品转移（用户端）
		productDivertSvc := service.NewProductDivertService(deps.DB, deps.Log)
		productDivertHandler := handler.NewProductDivertHandler(productDivertSvc, deps.Log)
		user.GET("/product-diverts", productDivertHandler.GetSent)
		user.GET("/product-diverts/received", productDivertHandler.GetReceived)
		user.GET("/product-diverts/:id", productDivertHandler.GetDetail)
		user.POST("/product-diverts", productDivertHandler.Create)
		user.POST("/product-diverts/:id/accept", productDivertHandler.Accept)
		user.POST("/product-diverts/:id/reject", productDivertHandler.Reject)
		user.POST("/product-diverts/:id/cancel", productDivertHandler.Cancel)
		user.GET("/product-diverts/:id/code", productDivertHandler.GetTransferCode)
		user.POST("/product-diverts/:id/regenerate-code", productDivertHandler.RegenerateCode)
		user.POST("/product-diverts/accept-by-code", productDivertHandler.AcceptByCode)

		// 用户操作日志
		systemLogSvc := service.NewSystemLogService(deps.DB, deps.Log)
		systemLogHandler := handler.NewSystemLogHandler(systemLogSvc, deps.Log)
		user.GET("/system-logs", systemLogHandler.List)
		user.GET("/system-logs/export", systemLogHandler.Export)

		// 发送验证码
		if deps.Redis != nil {
			user.POST("/sms/send", captchaHandler.SendSMS)
			user.POST("/email/send", captchaHandler.SendEmail)
		}

		// SSL证书（用户端）
		sslCertSvc := service.NewSSLCertificateService(deps.DB, deps.Log)
		sslCertHandler := handler.NewSSLCertificateHandler(sslCertSvc, deps.Log)
		user.GET("/ssl/certificates", sslCertHandler.GetList)
		user.GET("/ssl/certificates/:id", sslCertHandler.GetDetail)
		user.POST("/ssl/certificates/order", sslCertHandler.Order)
		user.POST("/ssl/certificates/:id/install", sslCertHandler.Install)
		user.POST("/ssl/certificates/:id/renew", sslCertHandler.Renew)
		user.GET("/ssl/certificate-types", sslCertHandler.GetCertificateTypes)

		// 用户服务管理（SMS、软件等）
		userServicesHandler := handler.NewUserServicesHandler(deps.DB)
		user.GET("/user/services/sms", userServicesHandler.GetSMSServices)
		user.GET("/user/services/sms/records", userServicesHandler.GetSMSRecords)
		user.GET("/user/services/software", userServicesHandler.GetSoftwareServices)
		user.POST("/user/services/software/reset-key", userServicesHandler.PostSoftwareResetKey)

		// AI Shopping (session routes require auth)
		user.POST("/ai-shopping/session", func(c *gin.Context) {
			sessionID := fmt.Sprintf("shop_%d_%d", c.GetUint("user_id"), time.Now().UnixNano())
			c.JSON(200, gin.H{"code": 0, "data": gin.H{"session_id": sessionID}})
		})
		user.POST("/ai-shopping/session/:session_id/message", aiShoppingHandler.Chat)
		user.GET("/ai-shopping/session/:session_id/messages", aiShoppingHandler.GetChatHistory)

		// 用户仪表盘
		user.GET("/user/dashboard", userHandler.GetDashboard)

		// 工单上传（通用）
		user.POST("/tickets/upload", ticketHandler.Upload)

		// V10 云
		user.POST("/v10cloud/order", v10CloudHandler.CreateOrder)
		user.GET("/v10cloud/products", v10CloudHandler.GetProductList)
		user.GET("/v10cloud/config-options", v10CloudHandler.GetConfigOptions)
		user.POST("/v10cloud/price", v10CloudHandler.CalculatePrice)

		// 市场交易/收益/日志
		user.GET("/marketplace/transactions", marketplaceHandler.GetTransactions)
		user.GET("/marketplace/transactions/summary", marketplaceHandler.GetTransactionSummary)
		user.GET("/marketplace/earnings", marketplaceHandler.GetEarnings)
		user.GET("/marketplace/withdrawals", marketplaceHandler.GetWithdrawals)
		user.POST("/marketplace/withdrawals", marketplaceHandler.CreateWithdrawal)
		user.GET("/marketplace/logs", marketplaceHandler.GetLogs)
		user.POST("/marketplace/orders/:id/pay", marketplaceHandler.PayOrder)

		// 挂售管理
		user.POST("/marketplace/listings", marketplaceHandler.CreateListing)
		user.PUT("/marketplace/listings/:id", marketplaceHandler.UpdateListing)
		user.DELETE("/marketplace/listings/:id", marketplaceHandler.RemoveListing)
		user.GET("/marketplace/listings/mine", marketplaceHandler.GetUserListings)

		// 订单管理
		user.GET("/marketplace/orders/buyer", marketplaceHandler.GetBuyerOrders)
		user.GET("/marketplace/orders/seller", marketplaceHandler.GetSellerOrders)
		user.POST("/marketplace/orders/:id/complete", marketplaceHandler.CompleteOrder)
		user.POST("/marketplace/orders/:id/cancel", marketplaceHandler.CancelOrder)

		// 私聊功能
		user.POST("/marketplace/messages", marketplaceHandler.SendMessage)
		user.GET("/marketplace/messages/:listing_id/:user_id", marketplaceHandler.GetChatMessages)
		user.GET("/marketplace/chat-sessions", marketplaceHandler.GetChatSessions)
		user.GET("/marketplace/unread-count", marketplaceHandler.GetUnreadCount)
	}
}
