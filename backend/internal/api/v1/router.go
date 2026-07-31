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
	captchaSvcForAuth := service.NewCaptchaService(deps.Redis, deps.DB)
	authHandler := handler.NewAuthHandlerWithCaptcha(deps.UserSvc, captchaSvcForAuth, deps.Log, deps.JWTKey)
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
	captchaConfigSvc := service.NewCaptchaConfigService(deps.DB)
	if captchaConfigSvc.IsGeetestEnabled() {
		captchaID, captchaKey := captchaConfigSvc.GetGeetestConfig()
		geetestSvc := service.NewGeetestService(captchaID, captchaKey, deps.Log.Desugar())
		captchaHandler.SetGeetestService(geetestSvc)
	}

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

	// ─── 主机增强管理 ───
	hostEnhancedSvc := service.NewHostEnhancedService(deps.DB, deps.Log, hostSvc, deps.InvSvc, balSvc, upstreamSvc)
	hostEnhancedHandler := handler.NewHostEnhancedHandler(hostEnhancedSvc, deps.Log)

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
	r.POST("/captcha/sms/verify", captchaHandler.VerifySMS)
	r.POST("/captcha/email", captchaHandler.SendEmail)
	r.POST("/captcha/email/verify", captchaHandler.VerifyEmail)
	r.GET("/captcha/image", captchaHandler.GetImage)
	r.GET("/captcha/image/json", captchaHandler.GetImageJSON)
	r.POST("/captcha/image/verify", captchaHandler.VerifyImage)
	r.GET("/captcha/status", captchaConfigHandler.GetSceneStatus) // 获取验证码场景状态

	// 极验验证码
	r.GET("/captcha/geetest/config", captchaHandler.GetGeetestConfig)
	r.POST("/captcha/geetest/verify", captchaHandler.VerifyGeetest)

	// 公开配置
	configHandler := handler.NewConfigHandler(deps.DB)
	r.GET("/config/public", configHandler.GetPublicConfig)
	r.GET("/config/login", configHandler.GetLoginConfig)
	r.GET("/config/maintenance", configHandler.GetMaintenanceStatus)

	// 短信验证码（用户端）
	smsSvcForVerify := service.NewSMSService(deps.DB, deps.Log)
	smsHandlerForVerify := handler.NewSMSHandler(smsSvcForVerify, deps.Log)
	r.POST("/sms/send-verify", smsHandlerForVerify.SendVerifyCode)

	// 产品
	r.GET("/products", productHandler.GetList)
	r.GET("/products/:id", productHandler.GetDetail)
	r.GET("/products/hot", productHandler.GetHot)
	r.GET("/product-groups", productHandler.GetGroups)

	// 前台导航分组
	navGroupHandler := handler.NewNavGroupHandler(deps.DB)
	r.GET("/nav-groups", navGroupHandler.GetPublicNavGroups)

	// 产品配置选项
	configOptionSvcV1 := service.NewConfigOptionService(deps.DB, deps.Log)
	configOptionHandlerV1 := handler.NewConfigOptionHandler(configOptionSvcV1, deps.Log)
	r.GET("/products/:id/config", configOptionHandlerV1.GetProductConfig)

	// 自定义字段（公开：产品字段）
	customFieldSvc := service.NewCustomFieldService(deps.DB, deps.Log)
	customFieldHandler := handler.NewCustomFieldHandler(customFieldSvc, deps.Log)
	r.GET("/custom-fields/product/:product_id", customFieldHandler.GetProductCustomFields)
	r.GET("/custom-fields/cart", customFieldHandler.GetCartCustomFields)

	// V10云产品（公开）
	r.GET("/v10/cloud/products", v10CloudHandler.GetProductList)
	r.GET("/v10/cloud/products/:id", v10CloudHandler.GetProductDetail)
	r.GET("/v10/cloud/products/:id/config", v10CloudHandler.GetConfigOptions)
	r.GET("/v10/cloud/products/:id/linkage", v10CloudHandler.GetLinkAgeList)
	r.GET("/v10/cloud/products/:id/config/filter", v10CloudHandler.FilterConfigOptions)
	r.GET("/v10/cloud/regions", v10CloudHandler.GetRegions)
	r.GET("/v10/cloud/os-types", v10CloudHandler.GetOSTypes)
	r.POST("/v10/cloud/calculate-price", v10CloudHandler.CalculatePrice)

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

	// 支付方式（公开，无需登录）
	r.GET("/payments/gateways", balanceHandler.GetEnabledGateways)

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

		// V10云购物车
		auth.GET("/v10/cloud/cart", v10CloudHandler.GetCartSummary)
		auth.GET("/v10/cloud/cart/items", v10CloudHandler.GetCartItems)
		auth.POST("/v10/cloud/cart", v10CloudHandler.AddToCart)
		auth.PUT("/v10/cloud/cart/:id", v10CloudHandler.UpdateCartItem)
		auth.POST("/v10/cloud/cart/settle", v10CloudHandler.SettleCart)

		// V10云订单
		auth.POST("/v10/cloud/orders", v10CloudHandler.CreateOrder)
		auth.GET("/v10/cloud/orders/:id", v10CloudHandler.GetOrderDetail)
		auth.POST("/v10/cloud/orders/:id/pay", v10CloudHandler.PayOrder)

		// V10云主机管理
		auth.GET("/v10/cloud/hosts/:id", v10CloudHandler.GetHostInfo)
		auth.GET("/v10/cloud/hosts/:id/config", v10CloudHandler.GetHostConfig)
		auth.GET("/v10/cloud/hosts/:id/traffic", v10CloudHandler.GetTrafficUsage)
		auth.GET("/v10/cloud/hosts/:id/os", v10CloudHandler.GetOSList)

		// 余额
		auth.GET("/balances", balanceHandler.GetBalance)
		auth.GET("/balances/logs", balanceHandler.GetBalanceLogs)
		auth.POST("/balances/recharge", balanceHandler.Recharge)
		auth.POST("/balances/withdraw", balanceHandler.Withdraw)

		// 支付方式（用户可见）
		auth.GET("/payments/gateways", balanceHandler.GetEnabledGateways)

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

		// 主机增强管理 - 续费
		auth.GET("/host/:id/renewal", hostEnhancedHandler.GetRenewalPage)
		auth.GET("/host/:id/renewal/price", hostEnhancedHandler.GetRenewalPrice)
		auth.POST("/host/:id/renewal", hostEnhancedHandler.SubmitRenewal)
		auth.PUT("/host/:id/auto-renew", hostEnhancedHandler.SetAutoRenew)
		auth.POST("/host/batch-renew", hostEnhancedHandler.BatchRenew)

		// 主机增强管理 - 转让
		auth.POST("/host/:id/transfer", hostEnhancedHandler.TransferHost)
		auth.GET("/host/:id/transfer/history", hostEnhancedHandler.GetTransferHistory)

		// 主机增强管理 - 分类
		auth.GET("/host/categories", hostEnhancedHandler.GetHostCategories)
		auth.POST("/host/categories", hostEnhancedHandler.CreateHostCategory)
		auth.POST("/host/:id/category", hostEnhancedHandler.AssignCategory)
		auth.DELETE("/host/categories/:category_id", hostEnhancedHandler.DeleteCategory)

		// 主机增强管理 - SSL
		auth.GET("/host/:id/ssl", hostEnhancedHandler.GetSSLConfig)
		auth.PUT("/host/:id/ssl", hostEnhancedHandler.SetSSLConfig)
		auth.POST("/host/:id/ssl/install", hostEnhancedHandler.InstallSSL)

		// 主机增强管理 - 下游/二次验证
		auth.PUT("/host/:id/downstream", hostEnhancedHandler.SetDownstream)
		auth.PUT("/host/:id/second-verify", hostEnhancedHandler.SetSecondVerify)
		auth.POST("/host/:id/second-verify", hostEnhancedHandler.VerifySecond)

		// 主机增强管理 - 流量与电源
		auth.GET("/host/:id/traffic", hostEnhancedHandler.GetTrafficUsage)
		auth.GET("/host/:id/traffic/chart", hostEnhancedHandler.GetTrafficChart)
		auth.POST("/host/:id/power/refresh", hostEnhancedHandler.RefreshPowerStatus)

		// 主机增强管理 - 详情与状态
		auth.GET("/host/:id/detail", hostEnhancedHandler.GetHostDetail)
		auth.GET("/host/:id/status", hostEnhancedHandler.GetHostStatus)
		auth.POST("/host/:id/remark", hostEnhancedHandler.PostRemark)
		auth.POST("/host/:id/hide", hostEnhancedHandler.HideHost)

		// 主机增强管理 - 取消/终止
		auth.GET("/host/:id/cancel", hostEnhancedHandler.GetCancelPage)
		auth.POST("/host/:id/cancel", hostEnhancedHandler.SubmitCancel)
		auth.DELETE("/host/:id/cancel", hostEnhancedHandler.DeleteCancel)

		// 主机增强管理 - 独立服务器
		auth.GET("/host/:id/dedicated", hostEnhancedHandler.GetDedicatedServer)

		// 主机增强管理 - 流量包
		auth.GET("/host/:id/flow-packets", hostEnhancedHandler.GetFlowPackets)
		auth.POST("/host/:id/flow-packets/:packet_id", hostEnhancedHandler.BuyFlowPacket)

		// 主机增强管理 - 重装增强
		auth.GET("/host/:id/reinstall/check", hostEnhancedHandler.CheckReinstall)
		auth.GET("/host/:id/reinstall/status", hostEnhancedHandler.GetReinstallStatus)
		auth.POST("/host/:id/reinstall/cancel", hostEnhancedHandler.CancelReinstall)

		// 主机增强管理 - 主机充值
		auth.GET("/host/:id/recharge", hostEnhancedHandler.GetHostRecharge)

		// 模块供给（客户端）
		provisionSvcForModule := service.NewProvisionService(deps.DB, deps.Log)
		pmClientHandler := handler.NewProvisionModuleHandler(provisionSvcForModule, deps.Log)
		auth.GET("/provision/client-area/:host_id", pmClientHandler.RenderClientArea)
		auth.GET("/provision/client-area/:host_id/detail", pmClientHandler.RenderClientAreaDetail)
		auth.GET("/provision/client-area/:host_id/buttons", pmClientHandler.GetClientButtons)
		auth.POST("/provision/client-area/:host_id/execute", pmClientHandler.ExecuteClientButton)
		auth.GET("/provision/charts/:host_id", pmClientHandler.GetCharts)
		auth.GET("/provision/charts/:host_id/:chart_id", pmClientHandler.GetChartData)
		auth.GET("/provision/usage/:host_id", pmClientHandler.GetUsage)
		auth.GET("/provision/usage/:host_id/traffic", pmClientHandler.TrafficUsage)
		auth.GET("/provision/ssl/:host_id", pmClientHandler.SSLButton)

		// 主机操作（catch-all，放最后）
		auth.POST("/host/:id/:action", hostHandler.PerformAction)

		// 魔方云高级操作（用户）
		dcimCloudSvcV1 := service.NewDcimCloudService(deps.DB, deps.Log)
		dcimCloudHandlerV1 := handler.NewDcimCloudHandler(dcimCloudSvcV1, deps.Log)
		auth.GET("/cloud/:id/nat", dcimCloudHandlerV1.GetNATRules)
		auth.POST("/cloud/:id/nat", dcimCloudHandlerV1.CreateNATRule)
		auth.PUT("/cloud/nat/:rule_id", dcimCloudHandlerV1.UpdateNATRule)
		auth.DELETE("/cloud/nat/:rule_id", dcimCloudHandlerV1.DeleteNATRule)
		auth.GET("/cloud/:id/security-groups", dcimCloudHandlerV1.GetSecurityGroups)
		auth.POST("/cloud/:id/security-groups", dcimCloudHandlerV1.CreateSecurityGroup)
		auth.PUT("/cloud/security-groups/:group_id", dcimCloudHandlerV1.UpdateSecurityGroup)
		auth.DELETE("/cloud/security-groups/:group_id", dcimCloudHandlerV1.DeleteSecurityGroup)
		auth.GET("/cloud/security-groups/:group_id/rules", dcimCloudHandlerV1.GetSecurityGroupRules)
		auth.POST("/cloud/security-groups/:group_id/rules", dcimCloudHandlerV1.AddSecurityGroupRule)
		auth.PUT("/cloud/security-rules/:rule_id", dcimCloudHandlerV1.UpdateSecurityGroupRule)
		auth.DELETE("/cloud/security-rules/:rule_id", dcimCloudHandlerV1.DeleteSecurityGroupRule)
		auth.GET("/cloud/isos", dcimCloudHandlerV1.GetISOList)
		auth.POST("/cloud/:id/mount-iso", dcimCloudHandlerV1.MountISO)
		auth.POST("/cloud/:id/unmount-iso", dcimCloudHandlerV1.UnmountISO)
		auth.GET("/cloud/:id/vnc", dcimCloudHandlerV1.GetVNCURL)
		auth.GET("/cloud/:id/vnc-page", dcimCloudHandlerV1.GetVNCPage)
		auth.GET("/cloud/:id/charts", dcimCloudHandlerV1.GetResourceChart)
		auth.GET("/cloud/:id/resources", dcimCloudHandlerV1.GetResourceInfo)
		auth.GET("/cloud/:id/flow-packets", dcimCloudHandlerV1.GetFlowPackets)
		auth.POST("/cloud/:id/flow-packets/buy", dcimCloudHandlerV1.BuyFlowPacket)
		auth.GET("/cloud/:id/flow-usage", dcimCloudHandlerV1.GetFlowPacketUsage)
		auth.GET("/cloud/:id/status", dcimCloudHandlerV1.GetCloudStatus)
		auth.POST("/cloud/:id/reset-password", dcimCloudHandlerV1.ResetCloudPassword)

		// DCIM服务器操作（用户端）
		dcimSvc := service.NewDcimService(deps.DB, deps.Log)
		dcimAdvHandler := handler.NewDcimAdvancedHandler(dcimSvc, deps.Log)
		// KVM/远程控制
		auth.GET("/dcim/server/:id/kvm", dcimAdvHandler.GetKVMURL)
		auth.GET("/dcim/server/:id/bmc", dcimAdvHandler.GetBMCInfo)
		auth.GET("/dcim/server/:id/novnc", dcimAdvHandler.GetNoVNCURL)
		// 救援系统
		auth.POST("/dcim/server/:id/rescue", dcimAdvHandler.BootRescue)
		auth.POST("/dcim/server/:id/crack-password", dcimAdvHandler.CrackPassword)
		auth.GET("/dcim/server/:id/rescue-status", dcimAdvHandler.GetRescueStatus)
		// 流量监控
		auth.GET("/dcim/server/:id/traffic", dcimAdvHandler.GetTrafficUsage)
		auth.GET("/dcim/server/:id/traffic/chart", dcimAdvHandler.GetTrafficChart)
		// 快照
		auth.POST("/dcim/server/:id/snapshots", dcimAdvHandler.CreateSnapshot)
		auth.GET("/dcim/server/:id/snapshots", dcimAdvHandler.GetSnapshots)
		auth.POST("/dcim/snapshot/:id/restore", dcimAdvHandler.RestoreSnapshot)
		auth.DELETE("/dcim/snapshot/:id", dcimAdvHandler.DeleteSnapshot)
		// 备份
		auth.POST("/dcim/server/:id/backups", dcimAdvHandler.CreateBackup)
		auth.GET("/dcim/server/:id/backups", dcimAdvHandler.GetBackups)
		auth.POST("/dcim/backup/:id/restore", dcimAdvHandler.RestoreBackup)
		auth.DELETE("/dcim/backup/:id", dcimAdvHandler.DeleteBackup)
		// 电源状态
		auth.GET("/dcim/server/:id/power", dcimAdvHandler.GetPowerStatus)
		// 重装状态
		auth.GET("/dcim/server/:id/reinstall-status", dcimAdvHandler.GetReinstallStatus)
		auth.GET("/dcim/os-list", dcimAdvHandler.GetOSList)

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

		// 自定义字段（用户端）
		auth.GET("/custom-fields/client", customFieldHandler.GetClientCustomFields)
		auth.GET("/custom-fields/host/:host_id", customFieldHandler.GetHostCustomFields)
		auth.GET("/custom-field-values", customFieldHandler.GetValues)
		auth.POST("/custom-field-values", customFieldHandler.SaveValues)
		auth.POST("/custom-fields/validate", customFieldHandler.ValidateFields)

		// 推介
		auth.GET("/affiliate/info", affiliateHandler.GetInfo)
		auth.GET("/affiliate/records", affiliateHandler.GetRecords)
		auth.GET("/affiliate/withdraws", affiliateHandler.GetWithdraws)
		auth.POST("/affiliate/withdraws", affiliateHandler.ApplyWithdraw)

	// 实名认证
	auth.POST("/certification/submit", certHandler.Submit)
	auth.GET("/certification/status", certHandler.GetStatus)

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

		// 客户跟踪
		auth.GET("/user/tracks", handler.NewClientTrackHandler(deps.DB).GetUserTracks)
	}

	// SSL证书类型（公开）
	sslCertSvcPub := service.NewSSLCertificateService(deps.DB, deps.Log)
	sslCertHandlerPub := handler.NewSSLCertificateHandler(sslCertSvcPub, deps.Log)
	r.GET("/ssl-certificates/types", sslCertHandlerPub.GetCertificateTypes)

	// 首页基本信息（公开）
	baseInfoHandlerPub := handler.NewBaseInfoHandler(deps.DB)
	r.GET("/base-infos", baseInfoHandlerPub.GetActive)

	// 取消原因（公开）
	cancelReasonHandlerPub := handler.NewCancelReasonHandler(deps.DB)
	r.GET("/cancel-reasons", cancelReasonHandlerPub.List)

	// 需要登录的接口
	auth := r.Group("")
	auth.Use(middleware.AuthMiddleware(deps.JWTMgr))
	{
		// 用户专属下载
		userDownloadHandler := handler.NewUserDownloadHandler(deps.DB)
		auth.GET("/user/downloads", userDownloadHandler.GetUserDownloads)

		// 用户偏好设置
		userTasteHandler := handler.NewUserTasteHandler(deps.DB)
		auth.GET("/user/tastes", userTasteHandler.Get)
		auth.PUT("/user/tastes", userTasteHandler.Update)
	}
}
