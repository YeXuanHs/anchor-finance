package v2

import (
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

// RegisterRoutes registers all v2 API routes on the given router group.
// v2 is AnchorFinance's native API format for same-system integration.
func RegisterRoutes(r *gin.RouterGroup, deps Deps) {
	log := deps.Log

	// 初始化处理器
	authHandler := handler.NewAuthHandler(deps.UserSvc, log, deps.JWTMgr)
	userHandler := handler.NewUserHandler(deps.UserSvc, log)
	productHandler := handler.NewProductHandler(deps.ProdSvc, log)
	orderHandler := handler.NewOrderHandlerWithDB(deps.OrdSvc, deps.DB, log)
	invoiceHandler := handler.NewInvoiceHandler(deps.InvSvc, log)
	ticketHandler := handler.NewTicketHandler(deps.TicSvc, log)
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
	promoCodeHandler := handler.NewPromoCodeHandler(deps.DB, deps.Log)

	announceSvc := service.NewAnnounceService(deps.DB, log)
	announceHandler := handler.NewAnnounceHandler(announceSvc, log)

	newsHandler := handler.NewNewsHandler(deps.DB, log)

	knowledgeSvc := service.NewKnowledgeService(deps.DB, log)
	knowledgeHandler := handler.NewKnowledgeHandler(knowledgeSvc, log)

	systemMsgSvc := service.NewSystemMessageService(deps.DB, log)
	systemMsgHandler := handler.NewSystemMessageHandler(systemMsgSvc, log)

	certSvc := service.NewCertificationService(deps.DB, log)
	certHandler := handler.NewCertificationHandler(certSvc, log)

	oauthHandler := handler.NewOAuthHandler(deps.OAuthSvc, log, deps.JWTMgr)

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

	// ==================== 公开接口 ====================

	// 认证
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", middleware.AuthRequired(), authHandler.RefreshToken)
		auth.POST("/logout", middleware.AuthRequired(), authHandler.Logout)
		auth.POST("/access-token", authHandler.AccessTokenLogin)
	}

	// 产品
	r.GET("/products", productHandler.GetList)
	r.GET("/products/:id", productHandler.GetDetail)
	r.GET("/products/hot", productHandler.GetHot)
	r.GET("/products/groups", productHandler.GetGroups)

	// 公告
	r.GET("/announcements", announceHandler.GetActive)
	r.GET("/announcements/:id", announceHandler.AdminGetDetail)

	// 新闻
	r.GET("/news", newsHandler.GetList)
	r.GET("/news/:id", newsHandler.GetDetail)
	r.GET("/news/categories", newsHandler.GetCategories)

	// 知识库
	r.GET("/help/categories", knowledgeHandler.GetCategories)
	r.GET("/help/articles", knowledgeHandler.GetArticles)
	r.GET("/help/articles/hot", knowledgeHandler.GetHot)
	r.GET("/help/articles/search", knowledgeHandler.Search)
	r.GET("/help/articles/:id", knowledgeHandler.GetArticle)
	r.GET("/help/articles/:id/related", knowledgeHandler.GetRelatedArticles)
	r.POST("/help/articles/:id/feedback", knowledgeHandler.SubmitFeedback)
	r.GET("/help/categories/:id/sub", knowledgeHandler.GetSubCategories)

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

		// API密钥管理
		user.GET("/user/api-keys", userHandler.GetAPIKeys)
		user.POST("/user/api-keys", userHandler.CreateAPIKey)
		user.PUT("/user/api-keys/:id/toggle", userHandler.ToggleAPIKey)
		user.DELETE("/user/api-keys/:id", userHandler.DeleteAPIKey)

		// OAuth绑定
		user.GET("/oauth/providers", oauthHandler.GetProviders)
		user.POST("/oauth/:provider/bind", oauthHandler.BindAccount)
		user.DELETE("/oauth/:provider/unbind", oauthHandler.UnbindAccount)

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
			promo, err := promoCodeSvc.Validate(req.Code, userID)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"data": promo})
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

		// 登录日志
		user.GET("/login-logs", loginLogHandler.List)

		// API日志
		user.GET("/api-logs", apiLogHandler.List)

		// 产品组（用户可见）
		user.GET("/product-groups", func(c *gin.Context) {
			var groups []struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
				Slug string `json:"slug"`
			}
			deps.DB.Table("product_groups").Order("sort_order").Find(&groups)
			c.JSON(200, gin.H{"data": groups})
		})

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
			captchaSvc := service.NewCaptchaService(deps.Redis, deps.DB)
			captchaHandler := handler.NewCaptchaHandler(captchaSvc, deps.DB)
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

		// AI Shopping 配置
		user.GET("/ai-shopping/config", aiShoppingHandler.GetConfig)

		// 用户仪表盘
		user.GET("/user/dashboard", userHandler.GetDashboard)

		// 工单上传（通用）
		user.POST("/tickets/upload", ticketHandler.Upload)

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

		// DDoS管理
		ddosHandler := handler.NewDDoSHandler(deps.DB)
		user.GET("/user/ddos/ips", ddosHandler.GetIPs)
		user.POST("/user/ddos/ips", ddosHandler.PostIP)
		user.DELETE("/user/ddos/ips/:id", ddosHandler.DeleteIP)
		user.GET("/user/ddos/ips/:id/rules", ddosHandler.GetIPRules)
		user.POST("/user/ddos/ips/:id/toggle", ddosHandler.PostIPToggle)
		user.DELETE("/user/ddos/rules/:id", ddosHandler.DeleteRule)
		user.PUT("/user/ddos/rules/:id", ddosHandler.PutRule)
		user.GET("/user/ddos/traffic", ddosHandler.GetTraffic)
		user.GET("/user/ddos/overview", ddosHandler.GetOverview)
	}

	// 解决方案（公开）
	publicSvc := service.NewPublicService(deps.DB, deps.Log)
	publicHandler := handler.NewPublicHandler(publicSvc, deps.Log)
	r.GET("/solutions", publicHandler.GetSolutions)
	r.GET("/solutions/:slug", publicHandler.GetSolutionDetail)

	// 托管服务
	r.GET("/colocation/advantages", publicHandler.GetColocationAdvantages)
	r.GET("/colocation/datacenters", publicHandler.GetColocationDatacenters)
}
