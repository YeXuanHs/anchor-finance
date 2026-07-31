package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"anchorfinance/internal/handler"
	"anchorfinance/internal/api/middleware"
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
	orderHandler := handler.NewOrderHandler(deps.OrdSvc, log)
	invoiceHandler := handler.NewInvoiceHandler(deps.InvSvc, log)
	ticketHandler := handler.NewTicketHandler(deps.TicSvc, log)
	cartHandler := handler.NewCartHandler(deps.CartSvc, log)

	balanceSvc := service.NewBalanceLogService(deps.DB)
	balanceHandler := handler.NewBalanceHandler(balanceSvc, deps.DB)

	creditSvc := service.NewCreditService(deps.DB)
	creditHandler := handler.NewCreditHandler(deps.DB, creditSvc)

	hostSvc := service.NewHostService(deps.DB, log)
	hostHandler := handler.NewHostHandler(hostSvc, log)

	contactsSvc := service.NewClientContactService(deps.DB)
	contactsHandler := handler.NewContactsHandler(contactsSvc, log)

	couponSvc := service.NewCouponService(deps.DB)
	couponHandler := handler.NewCouponHandler(couponSvc)

	announceSvc := service.NewAnnounceService(deps.DB, log)
	announceHandler := handler.NewAnnounceHandler(announceSvc, log)

	newsSvc := service.NewNewsService(deps.DB, log)
	newsHandler := handler.NewNewsHandler(newsSvc, log)

	knowledgeSvc := service.NewKnowledgeService(deps.DB, log)
	knowledgeHandler := handler.NewKnowledgeHandler(knowledgeSvc, log)

	systemMsgSvc := service.NewSystemMessageService(deps.DB)
	systemMsgHandler := handler.NewSystemMessageHandler(systemMsgSvc, log)

	certSvc := service.NewCertificationService(deps.DB)
	certHandler := handler.NewCertificationHandler(certSvc, log)

	oauthHandler := handler.NewOAuthHandler(deps.OAuthSvc, log, deps.JWTMgr)

	upgradeSvc := service.NewUpgradeService(deps.DB)
	upgradeHandler := handler.NewUpgradeHandler(upgradeSvc, log)

	contractHandler := handler.NewContractHandler(deps.DB, log)

	affiliateSvc := service.NewAffiliateService(deps.DB)
	affiliateHandler := handler.NewAffiliateHandler(affiliateSvc, log)

	voucherSvc := service.NewVoucherService(deps.DB)
	voucherHandler := handler.NewVoucherHandler(voucherSvc, log)

	multiRenewSvc := service.NewMultiRenewService(deps.DB)
	multiRenewHandler := handler.NewMultiRenewHandler(multiRenewSvc, log)

	loginLogSvc := service.NewLoginLogService(deps.DB)
	loginLogHandler := handler.NewLoginLogHandler(loginLogSvc, log)

	apiLogSvc := service.NewAPILogService(deps.DB)
	apiLogHandler := handler.NewApiLogHandler(apiLogSvc, log)

	// ==================== 公开接口 ====================

	// 认证
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", middleware.AuthRequired(), authHandler.RefreshToken)
		auth.POST("/logout", middleware.AuthRequired(), authHandler.Logout)
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

	// 系统设置（公开部分）
	r.GET("/system/settings", func(c *gin.Context) {
		var settings []struct{ Key, Value string }
		deps.DB.Table("system_settings").Where("key IN ?", []string{
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
		user.POST("/hosts/:id/:action", hostHandler.PerformAction)

		// 购物车
		user.GET("/cart", cartHandler.GetCart)
		user.POST("/cart/add", cartHandler.AddToCart)
		user.PUT("/cart/:id", cartHandler.UpdateCart)
		user.DELETE("/cart/:id", cartHandler.RemoveFromCart)
		user.DELETE("/cart", cartHandler.ClearCart)
		user.POST("/cart/checkout", cartHandler.Checkout)

		// 订单
		user.POST("/orders", orderHandler.Create)
		user.GET("/orders", orderHandler.GetUserOrders)
		user.GET("/orders/:id", orderHandler.GetDetail)
		user.POST("/orders/:id/pay", orderHandler.Pay)
		user.POST("/orders/:id/cancel", orderHandler.Cancel)

		// 账单/发票
		user.GET("/invoices", invoiceHandler.GetUserInvoices)
		user.GET("/invoices/:id", invoiceHandler.GetDetail)
		user.POST("/invoices/:id/pay", invoiceHandler.Pay)

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

		// 信用额度
		user.GET("/credit", creditHandler.GetInfo)
		user.GET("/credit/logs", creditHandler.GetLogs)
		user.POST("/credit/apply", creditHandler.Apply)
		user.POST("/credit/repay", creditHandler.Repay)

		// 联系人
		user.GET("/contacts", contactsHandler.GetList)
		user.POST("/contacts", contactsHandler.Create)
		user.PUT("/contacts/:id", contactsHandler.Update)
		user.DELETE("/contacts/:id", contactsHandler.Delete)
		user.PUT("/contacts/:id/default", contactsHandler.SetDefault)
		user.GET("/contacts/default", contactsHandler.GetDefault)

		// 优惠券
		user.GET("/coupons", couponHandler.GetUserCoupons)

		// 系统消息
		user.GET("/messages", systemMsgHandler.GetList)
		user.GET("/messages/unread-count", systemMsgHandler.GetUnreadCount)
		user.PUT("/messages/:id/read", systemMsgHandler.MarkRead)
		user.PUT("/messages/read-all", systemMsgHandler.MarkAllRead)
		user.DELETE("/messages/:id", systemMsgHandler.Delete)

		// 实名认证
		user.GET("/certification/status", certHandler.GetStatus)
		user.POST("/certification/submit", certHandler.Submit)

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

		// 代理/推荐
		user.GET("/affiliate/info", affiliateHandler.GetInfo)
		user.GET("/affiliate/records", affiliateHandler.GetRecords)
		user.GET("/affiliate/withdraws", affiliateHandler.GetWithdraws)
		user.POST("/affiliate/withdraw", affiliateHandler.ApplyWithdraw)

		// 代金券
		user.GET("/vouchers", voucherHandler.GetUserVouchers)
		user.POST("/vouchers/claim", voucherHandler.ClaimVoucher)

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
	}
}
