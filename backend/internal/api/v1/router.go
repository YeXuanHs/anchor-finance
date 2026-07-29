package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"anchorfinance/internal/handler"
	"anchorfinance/internal/api/middleware"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"gorm.io/gorm"
)

// Deps holds shared dependencies for user-facing API route registration.
type Deps struct {
	DB       *gorm.DB
	Redis    *redis.Client
	Log      *logger.Logger
	JWTKey   string
	UserSvc  *service.UserService
	ProdSvc  *service.ProductService
	OrdSvc   *service.OrderService
	InvSvc   *service.InvoiceService
	TicSvc   *service.TicketService
	CartSvc  *service.CartService
	OAuthSvc *service.OAuthService
}

// RegisterRoutes registers all user-facing v1 API routes.
func RegisterRoutes(r *gin.RouterGroup, deps Deps) {
	// ─── 核心 handlers ───
	captchaSvcForAuth := service.NewCaptchaService(deps.Redis)
	authHandler := handler.NewAuthHandlerWithCaptcha(deps.UserSvc, captchaSvcForAuth, deps.Log, deps.JWTKey)
	userHandler := handler.NewUserHandlerWithCaptcha(deps.UserSvc, captchaSvcForAuth, deps.Log)
	productHandler := handler.NewProductHandler(deps.ProdSvc, deps.Log)
	orderHandler := handler.NewOrderHandlerWithDB(deps.OrdSvc, deps.DB, deps.Log)
	invoiceHandler := handler.NewInvoiceHandler(deps.InvSvc, deps.Log)
	ticketHandler := handler.NewTicketHandler(deps.TicSvc, deps.Log)
	cartHandler := handler.NewCartHandler(deps.CartSvc)

	// ─── 验证码 ───
	captchaHandler := handler.NewCaptchaHandler(service.NewCaptchaService(deps.Redis), deps.DB)

	// ─── 余额 ───
	balanceSvc := service.NewBalanceLogService(deps.DB)
	balanceHandler := handler.NewBalanceHandler(balanceSvc, deps.DB)

	// ─── 信用额度 ───
	creditSvc := service.NewCreditService(deps.DB, deps.Log)
	creditHandler := handler.NewCreditHandler(deps.DB, creditSvc)

	// ─── 联系人 ───
	contactSvc := service.NewContactService(deps.DB, deps.Log)
	contactsHandler := handler.NewContactsHandler(contactSvc, deps.Log)

	// ─── 合同 ───
	contractHandler := handler.NewContractHandler(deps.DB, deps.Log)

	// ─── 主机 ───
	upstreamSvc := service.NewUpstreamService(deps.DB, deps.Log)
	hostSvc := service.NewHostService(deps.DB, deps.Log, upstreamSvc)
	hostHandler := handler.NewHostHandler(hostSvc, deps.Log)

	// ─── 多续费 ───
	balSvc := service.NewBalanceService(deps.DB, deps.Log)
	multiRenewSvc := service.NewMultiRenewService(deps.DB, deps.Log, deps.InvSvc, balSvc)
	multiRenewHandler := handler.NewMultiRenewHandler(multiRenewSvc, deps.Log)

	// ─── 维护公告 ───
	maintenanceSvc := service.NewMaintenanceService(deps.DB, deps.Log)
	maintenanceHandler := handler.NewMaintenanceHandler(maintenanceSvc, deps.Log)

	// ─── 日志 ───
	loginLogSvc := service.NewLoginLogService(deps.DB, deps.Log)
	loginLogHandler := handler.NewLoginLogHandler(loginLogSvc, deps.Log)
	systemLogSvc := service.NewSystemLogService(deps.DB, deps.Log)
	systemLogHandler := handler.NewSystemLogHandler(systemLogSvc, deps.Log)
	apiLogSvc := service.NewApiLogService(deps.DB, deps.Log)
	apiLogHandler := handler.NewAPILogHandler(apiLogSvc, deps.Log)

	// ─── OAuth ───
	oauthHandler := handler.NewOAuthHandler(deps.OAuthSvc, deps.Log, deps.JWTKey)

	// ─── 公告 ───
	announceSvc := service.NewAnnounceService(deps.DB, deps.Log)
	announceHandler := handler.NewAnnounceHandler(announceSvc, deps.Log)

	// ─── 新闻 ───
	newsHandler := handler.NewNewsHandler(deps.DB, deps.Log)

	// ─── 知识库 ───
	knowledgeSvc := service.NewKnowledgeService(deps.DB, deps.Log)
	knowledgeHandler := handler.NewKnowledgeHandler(knowledgeSvc, deps.Log)

	// ─── Banner ───
	bannerSvc := service.NewBannerService(deps.DB, deps.Log)
	bannerHandler := handler.NewBannerHandler(bannerSvc, deps.Log)

	// ─── 优惠券 ───
	couponSvc := service.NewCouponService(deps.DB, deps.Log)
	couponHandler := handler.NewCouponHandler(couponSvc)

	// ─── 实名认证 ───
	certSvc := service.NewCertificationService(deps.DB, deps.Log)
	certHandler := handler.NewCertificationHandler(certSvc, deps.Log)

	// ─── 系统消息 ───
	sysMsgSvc := service.NewSystemMessageService(deps.DB, deps.Log)
	sysMsgHandler := handler.NewSystemMessageHandler(sysMsgSvc, deps.Log)

	// ─── 升级 ───
	upgradeSvc := service.NewUpgradeService(deps.DB, deps.Log)
	upgradeHandler := handler.NewUpgradeHandler(upgradeSvc, deps.Log)

	// ─── 下载 ───
	downloadSvc := service.NewDownloadService(deps.DB, deps.Log)
	downloadHandler := handler.NewDownloadHandler(downloadSvc, deps.Log)

	// ─── 用户等级 ───
	userLevelHandler := handler.NewUserLevelHandler(deps.DB)

	// ─── 友情链接 ───
	friendlyLinkSvc := service.NewFriendlyLinkService(deps.DB, deps.Log)
	friendlyLinkHandler := handler.NewFriendlyLinkHandler(friendlyLinkSvc, deps.Log)

	// ─── 公共 ───
	publicSvc := service.NewPublicService(deps.DB, deps.Log)
	publicHandler := handler.NewPublicHandlerWithDB(publicSvc, deps.DB, deps.Log)

	// ─── 服务详情 ───
	serviceDetailSvc := service.NewServiceDetailService(deps.DB, deps.Log)
	serviceDetailHandler := handler.NewServiceDetailHandler(serviceDetailSvc, deps.Log)

	// ─── 推介 ───
	affiliateSvc := service.NewAffiliateService(deps.DB, deps.Log)
	affiliateHandler := handler.NewAffiliateHandler(affiliateSvc, deps.Log)

	// ─── 代金券 ───
	voucherSvc := service.NewVoucherService(deps.DB, deps.Log)
	voucherHandler := handler.NewVoucherHandler(voucherSvc, deps.Log)

	// ═══════════════════ 公开接口（无需登录） ═══════════════════

	// 认证
	r.POST("/login", authHandler.Login)
	r.POST("/login/sms", authHandler.SMSLogin)
	r.POST("/register", authHandler.Register)
	r.POST("/password/verify-code", authHandler.VerifyResetCode)
	r.POST("/password/reset", authHandler.ResetPassword)

	// 验证码
	r.POST("/captcha/sms", captchaHandler.SendSMS)
	r.POST("/captcha/email", captchaHandler.SendEmail)
	r.GET("/captcha/image", captchaHandler.GetImage)

	// 产品
	r.GET("/products", productHandler.GetList)
	r.GET("/products/:id", productHandler.GetDetail)
	r.GET("/products/hot", productHandler.GetHot)
	r.GET("/product-groups", productHandler.GetGroups)

	// 产品配置选项
	configOptionSvcV1 := service.NewConfigOptionService(deps.DB, deps.Log)
	configOptionHandlerV1 := handler.NewConfigOptionHandler(configOptionSvcV1, deps.Log)
	r.GET("/products/:id/config", configOptionHandlerV1.GetProductConfig)

	// 合作伙伴
	r.GET("/partners", publicHandler.GetPartners)

	// 新闻
	r.GET("/news", newsHandler.GetList)
	r.GET("/news/:id", newsHandler.GetDetail)
	r.GET("/news/categories", newsHandler.GetCategories)
	r.GET("/news/search", newsHandler.Search)

	// 公告
	r.GET("/announcements", announceHandler.GetActive)

	// 帮助中心/知识库
	r.GET("/help/categories", knowledgeHandler.GetCategories)
	r.GET("/help/articles", knowledgeHandler.GetArticles)
	r.GET("/help/hot", knowledgeHandler.GetHot)
	r.GET("/help/search", knowledgeHandler.Search)
	r.GET("/help/:slug", knowledgeHandler.GetArticle)

	// Banner
	r.GET("/banners", bannerHandler.GetActive)

	// 友情链接
	r.GET("/friendly-links", friendlyLinkHandler.GetActive)

	// 系统设置
	r.GET("/settings", publicHandler.GetConfigs)
	r.GET("/settings/:key", publicHandler.GetConfig)

	// 域名/SSL
	r.GET("/domain/suffixes", publicHandler.GetDomainSuffixes)
	r.GET("/ssl/certificates", publicHandler.GetSSLCertificates)

	// 解决方案
	r.GET("/solutions", publicHandler.GetSolutions)
	r.GET("/solutions/:id", publicHandler.GetSolutionDetail)

	// Anti-DDoS
	r.GET("/anti-ddos/capabilities", publicHandler.GetAntiDDoSCapabilities)
	r.GET("/anti-ddos/plans", publicHandler.GetAntiDDoSPlans)

	// 托管
	r.GET("/colocation/advantages", publicHandler.GetColocationAdvantages)
	r.GET("/colocation/datacenters", publicHandler.GetColocationDatacenters)

	// CDN
	r.GET("/cdn/advantages", publicHandler.GetCDNAdvantages)
	r.GET("/cdn/plans", publicHandler.GetCDNPlans)

	// 联系我们
	r.POST("/contact", publicHandler.SubmitContact)

	// ═══════════════════ 需要登录 ═══════════════════
	auth := r.Group("")
	auth.Use(middleware.AuthRequired())
	{
		// 用户
		auth.GET("/user/profile", userHandler.GetProfile)
		auth.PUT("/user/profile", userHandler.UpdateProfile)
		auth.POST("/user/change-password", userHandler.ChangePassword)
		auth.POST("/user/bind-phone", userHandler.BindPhone)
		auth.POST("/user/bind-email", userHandler.BindEmail)
		auth.GET("/user/dashboard", serviceDetailHandler.GetUserDashboard)
		auth.GET("/user/products", serviceDetailHandler.GetByUser)

		// OAuth
		auth.GET("/oauth/:provider", oauthHandler.Redirect)
		auth.GET("/oauth/:provider/callback", oauthHandler.Callback)
		auth.POST("/oauth/unbind", oauthHandler.Unbind)
		auth.GET("/oauth/accounts", oauthHandler.GetBoundAccounts)

		// 订单
		auth.POST("/order/preview", orderHandler.Preview)
		auth.POST("/order/create", orderHandler.Create)
		auth.GET("/orders", orderHandler.GetUserOrders)
		auth.GET("/orders/:id", orderHandler.GetDetail)
		auth.POST("/orders/:id/pay", orderHandler.Pay)
		auth.POST("/orders/:id/cancel", orderHandler.Cancel)

		// 账单
		auth.GET("/invoices", invoiceHandler.GetUserInvoices)
		auth.GET("/invoices/:id", invoiceHandler.GetDetail)
		auth.POST("/invoices/:id/pay", invoiceHandler.Pay)

		// 工单
		auth.POST("/tickets", ticketHandler.Create)
		auth.GET("/tickets", ticketHandler.GetUserTickets)
		auth.GET("/tickets/:id", ticketHandler.GetDetail)
		auth.POST("/tickets/:id/reply", ticketHandler.Reply)
		auth.POST("/tickets/:id/close", ticketHandler.Close)
		auth.POST("/tickets/:id/attachments", ticketHandler.UploadAttachment)
		auth.GET("/tickets/:id/attachments", ticketHandler.GetAttachments)
		auth.DELETE("/tickets/attachments/:id", ticketHandler.DeleteAttachment)

		// 购物车
		auth.GET("/cart", cartHandler.GetCart)
		auth.POST("/cart/add", cartHandler.AddToCart)
		auth.PUT("/cart/:id", cartHandler.UpdateCart)
		auth.DELETE("/cart/:id", cartHandler.RemoveFromCart)
		auth.DELETE("/cart/clear", cartHandler.ClearCart)
		auth.POST("/cart/checkout", cartHandler.Checkout)

		// 余额
		auth.GET("/balances", balanceHandler.GetBalance)
		auth.GET("/balances/logs", balanceHandler.GetBalanceLogs)
		auth.POST("/balances/recharge", balanceHandler.Recharge)
		auth.POST("/balances/withdraw", balanceHandler.Withdraw)

		// 信用额度
		auth.GET("/credit", creditHandler.GetInfo)
		auth.GET("/credit/logs", creditHandler.GetLogs)
		auth.POST("/credit/apply", creditHandler.Apply)
		auth.POST("/credit/repay", creditHandler.Repay)

		// 信用账单
		auth.GET("/credit/bills", creditHandler.GetBills)
		auth.GET("/credit/bills/:id", creditHandler.GetBillDetail)
		auth.POST("/credit/bills/:id/pay", creditHandler.PayBill)
		auth.GET("/credit/config", creditHandler.GetCreditConfig)
		auth.PUT("/credit/config", creditHandler.UpdateCreditConfig)

		// 联系人
		auth.GET("/contacts", contactsHandler.GetList)
		auth.POST("/contacts", contactsHandler.Create)
		auth.PUT("/contacts/:id", contactsHandler.Update)
		auth.DELETE("/contacts/:id", contactsHandler.Delete)

		// 合同
		auth.GET("/contracts", contractHandler.GetUserContracts)
		auth.GET("/contracts/:id", contractHandler.GetDetail)
		auth.POST("/contracts/:id/sign", contractHandler.SignContract)

		// 下载
		auth.GET("/downloads/categories", downloadHandler.GetCategories)
		auth.GET("/downloads", downloadHandler.GetFiles)
		auth.GET("/downloads/:id", downloadHandler.GetFile)

		// 用户等级
		auth.GET("/user-levels", userLevelHandler.GetAll)

		// 主机/服务详情
		auth.GET("/host/:id", serviceDetailHandler.GetDetail)
		auth.GET("/host/:id/log", serviceDetailHandler.GetServiceLogs)
		auth.GET("/host/:id/upgrade", upgradeHandler.GetAvailable)
		auth.POST("/host/:id/upgrade", upgradeHandler.Submit)
		auth.POST("/host/:id/:action", hostHandler.PerformAction)

		// 批量续费
		auth.POST("/multi-renew", multiRenewHandler.Create)

		// 维护公告
		auth.GET("/maintenance/notices", maintenanceHandler.GetStatus)

		// 系统消息
		auth.GET("/messages", sysMsgHandler.GetList)
		auth.PUT("/messages/:id/read", sysMsgHandler.MarkRead)
		auth.PUT("/messages/read-all", sysMsgHandler.MarkAllRead)
		auth.GET("/messages/unread-count", sysMsgHandler.GetUnreadCount)

		// 日志
		auth.GET("/api-logs", apiLogHandler.List)
		auth.GET("/login-logs", loginLogHandler.List)
		auth.GET("/system-logs", systemLogHandler.List)
		auth.GET("/system-logs/export", systemLogHandler.Export)

		// 优惠券
		auth.GET("/coupons", couponHandler.GetUserCoupons)
		auth.POST("/coupon/verify", couponHandler.ValidateCoupon)
		auth.POST("/coupons/verify", couponHandler.ValidateCoupon)

		// 代金券
		auth.GET("/vouchers", voucherHandler.GetUserVouchers)

		// 推介
		auth.GET("/affiliate/info", affiliateHandler.GetInfo)
		auth.GET("/affiliate/records", affiliateHandler.GetRecords)
		auth.GET("/affiliate/withdraws", affiliateHandler.GetWithdraws)
		auth.POST("/affiliate/withdraws", affiliateHandler.ApplyWithdraw)

	// 实名认证
	auth.POST("/certification/submit", certHandler.Submit)
	auth.GET("/certification/status", certHandler.GetStatus)

	// 域名管理
	domainSvc := service.NewDomainService(deps.DB, deps.Log)
	domainHandler := handler.NewDomainHandler(domainSvc, deps.DB, deps.Log)
	r.GET("/domain/check", domainHandler.CheckAvailability)
	auth.GET("/domains", domainHandler.GetList)
	auth.GET("/domains/:id", domainHandler.GetDetail)
	auth.POST("/domains/:id/renew", domainHandler.Renew)
	auth.GET("/domains/:id/dns", domainHandler.GetDNSRecords)
	auth.POST("/domains/:id/dns", domainHandler.AddDNSRecord)
	auth.PUT("/domains/:id/dns/:record_id", domainHandler.UpdateDNSRecord)
	auth.DELETE("/domains/:id/dns/:record_id", domainHandler.DeleteDNSRecord)
	auth.POST("/domains/transfer", domainHandler.InitiateTransfer)
	auth.GET("/domains/transfer", domainHandler.GetTransfers)

	// SSL证书
	sslCertSvc := service.NewSSLCertificateService(deps.DB, deps.Log)
	sslCertHandler := handler.NewSSLCertificateHandler(sslCertSvc, deps.Log)
	auth.GET("/ssl-certificates", sslCertHandler.GetList)
	auth.GET("/ssl-certificates/:id", sslCertHandler.GetDetail)
	auth.POST("/ssl-certificates/order", sslCertHandler.Order)
	auth.POST("/ssl-certificates/generate-csr", sslCertHandler.GenerateCSR)
	auth.POST("/ssl-certificates/:id/validate", sslCertHandler.Validate)
	auth.POST("/ssl-certificates/:id/install", sslCertHandler.Install)
	auth.POST("/ssl-certificates/:id/renew", sslCertHandler.Renew)
	auth.POST("/ssl-certificates/:id/revoke", sslCertHandler.Revoke)
	}

	// SSL证书类型（公开）
	sslCertSvcPub := service.NewSSLCertificateService(deps.DB, deps.Log)
	sslCertHandlerPub := handler.NewSSLCertificateHandler(sslCertSvcPub, deps.Log)
	r.GET("/ssl-certificates/types", sslCertHandlerPub.GetCertificateTypes)
}
