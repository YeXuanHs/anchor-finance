package client

import (
	"github.com/YeXuanHs/anchor-finance/internal/middleware"
	"github.com/YeXuanHs/anchor-finance/internal/service"
	"github.com/gin-gonic/gin"
)

// SetupZjmfCompatRoutes 注册zjmf兼容路由到根路径（/v1/）
func SetupZjmfCompatRoutes(r *gin.Engine, authService *service.AuthService) {
	// 公开端点
	v1Compat := r.Group("/v1")
	{
		v1Compat.POST("/login_api", ZjmfCompatLogin)
		v1Compat.GET("/products", ZjmfCompatProducts)
		v1Compat.GET("/hosts/cates", ZjmfCompatCategories)
	}

	// zjmf对接上游时调的端点
	r.POST("/zjmf_api_login", ZjmfCompatLogin) // zjmf.php line 193
	r.GET("/cart/all", ZjmfCompatCartAll)       // zjmf.php line 38
	r.GET("/api/product/proinfo", ZjmfCompatProductInfo) // zjmf.php line 48
	r.GET("/api/product/prodetail", ZjmfCompatProductDetail) // zjmf.php line 54
	r.GET("/cart/get_product_config", ZjmfCompatProductConfig) // zjmf.php line 70
	r.GET("/cart/ontrialmax", ZjmfCompatOnTrialMax) // CartController:1779
	r.GET("/cart/stock_control", ZjmfCompatStockControl) // CartController:2804
	r.POST("/apply_credit", middleware.JWTAuth(authService), ZjmfCompatApplyCredit) // common.php:2315
	r.GET("/host/trafficusage", middleware.JWTAuth(authService), ZjmfCompatTrafficUsage) // HostController:2192
	r.POST("/host/setdownstream", middleware.JWTAuth(authService), ZjmfCompatSetDownstream) // ClientsServicesController:415

	// v10类型端点（zjmf.php:190, ZjmfFinanceApiController:46）
	r.POST("/api/v1/auth", ZjmfCompatLogin)          // v10登录（和zjmf_api_login共用handler）
	r.GET("/api/v1/product", ZjmfCompatCartAll)       // v10商品列表
	r.GET("/api/v1/group/product", ZjmfCompatCartAll) // v10分组商品

	// 需要JWT认证的端点
	v1Auth := r.Group("/v1")
	v1Auth.Use(middleware.JWTAuth(authService))
	{
		v1Auth.GET("/user", ZjmfCompatUser)
		v1Auth.GET("/hosts/:id/module/status", ZjmfCompatModuleStatus)
		v1Auth.POST("/hosts/:id/module/suspend", ZjmfCompatModuleSuspend)
		v1Auth.POST("/hosts/:id/module/unsuspend", ZjmfCompatModuleUnsuspend)
		v1Auth.POST("/hosts/:id/module/terminate", ZjmfCompatModuleTerminate)
		v1Auth.GET("/hosts/:id/renew", ZjmfCompatRenew)
	}

	// zjmf兼容余额查询（/cart/credit，需JWT）
	r.GET("/cart/credit", middleware.JWTAuth(authService), ZjmfCompatBalance)
	// zjmf调host/header（需JWT）
	r.GET("/host/header", middleware.JWTAuth(authService), ZjmfCompatHostDetail)

	// DCIM硬件操作
	r.POST("/dcim/on", middleware.JWTAuth(authService), ZjmfCompatDcimOn)
	r.POST("/dcim/off", middleware.JWTAuth(authService), ZjmfCompatDcimOff)
	r.POST("/dcim/reboot", middleware.JWTAuth(authService), ZjmfCompatDcimReboot)
	r.POST("/dcim/rescue", middleware.JWTAuth(authService), ZjmfCompatDcimRescue)
	r.POST("/dcim/reinstall", middleware.JWTAuth(authService), ZjmfCompatDcimReinstall)
	r.POST("/dcim/crack_pass", middleware.JWTAuth(authService), ZjmfCompatDcimCrackPass)
	r.POST("/dcim/check_reinstall", middleware.JWTAuth(authService), ZjmfCompatDcimCheckReinstall)
	r.GET("/dcim/detail", middleware.JWTAuth(authService), ZjmfCompatDcimDetail)
	r.POST("/dcim/buy_reinstall_times", middleware.JWTAuth(authService), ZjmfCompatDcimBuyReinstallTimes)
	r.POST("/dcim/buy_flow_packet", middleware.JWTAuth(authService), ZjmfCompatDcimBuyFlowPacket)
	r.POST("/dcim/refresh_power_status", middleware.JWTAuth(authService), ZjmfCompatRefreshPowerStatus)
	r.POST("/dcim/refresh_all_power_status", middleware.JWTAuth(authService), ZjmfCompatRefreshAllPowerStatus)
	r.POST("/dcim/hide_result", middleware.JWTAuth(authService), ZjmfCompatDcimHideResult)

	// 升级操作
	r.POST("/upgrade/checkout_config_upgrade", middleware.JWTAuth(authService), ZjmfCompatUpgradeCheckoutConfig)
	r.POST("/upgrade/upgrade_product_post", middleware.JWTAuth(authService), ZjmfCompatUpgradeProductPost)
	r.POST("/upgrade/checkout_upgrade_product", middleware.JWTAuth(authService), ZjmfCompatUpgradeCheckoutProduct)

	// SSL证书
	r.POST("/provision/sslCertFunc", middleware.JWTAuth(authService), ZjmfCompatSslCertFunc)

	// 购物车/订单/信用额度操作
	r.POST("/cart/clear", middleware.JWTAuth(authService), ZjmfCompatClearCart)
	r.POST("/cart/add_to_shop", middleware.JWTAuth(authService), ZjmfCompatAddToShop)
	r.POST("/cart/settle", middleware.JWTAuth(authService), ZjmfCompatSettle)
	r.POST("/apply_credit_limit", middleware.JWTAuth(authService), ZjmfCompatApplyCreditLimit)
	r.GET("/user_info", middleware.JWTAuth(authService), ZjmfCompatUserInfo)

	// host/provision/upgrade 深度兼容端点
	r.POST("/host/cancel", middleware.JWTAuth(authService), ZjmfCompatHostCancel)
	r.POST("/host/renew", middleware.JWTAuth(authService), ZjmfCompatHostRenew)
	r.POST("/provision/default", middleware.JWTAuth(authService), ZjmfCompatProvisionDefault)
	r.POST("/provision/button", middleware.JWTAuth(authService), ZjmfCompatProvisionButton)
	r.POST("/provision/custom/:id", middleware.JWTAuth(authService), ZjmfCompatProvisionCustom)
	r.GET("/provision/chart/:id", middleware.JWTAuth(authService), ZjmfCompatProvisionChart)
	r.POST("/upgrade/upgrade_config_post", middleware.JWTAuth(authService), ZjmfCompatUpgradeConfigPost)

	// 补充缺失端点
	r.POST("/dcim/bmc", middleware.JWTAuth(authService), ZjmfCompatDcimBmc)
	r.POST("/dcim/cancel_task", middleware.JWTAuth(authService), ZjmfCompatDcimCancelTask)
	r.POST("/dcim/ikvm", middleware.JWTAuth(authService), ZjmfCompatDcimIkvm)
	r.POST("/dcim/kvm", middleware.JWTAuth(authService), ZjmfCompatDcimKvm)
	r.POST("/dcim/novnc", middleware.JWTAuth(authService), ZjmfCompatDcimNovnc)
	r.GET("/dcim/resintall_status", middleware.JWTAuth(authService), ZjmfCompatDcimReinstallStatus)
	r.POST("/dcim/traffic", middleware.JWTAuth(authService), ZjmfCompatDcimTraffic)
	r.GET("/dcim/traffic_usage", middleware.JWTAuth(authService), ZjmfCompatDcimTrafficUsage)
	r.GET("/cart/hostinfo", middleware.JWTAuth(authService), ZjmfCompatCartHostinfo)
	r.GET("/cart/summary", middleware.JWTAuth(authService), ZjmfCompatCartSummary)
}

// SetupRoutes 设置用户前台路由
func SetupRoutes(r *gin.RouterGroup, authService *service.AuthService) {
	// 创建处理器
	authHandler := NewAuthHandler(authService)

	// 公开路由（不需要认证）
	public := r.Group("")
	public.Use(middleware.RedirectWhitelist()) // MD 9.1漏洞6：跳转URL白名单
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
		public.POST("/auth/reset-password", authHandler.ResetPassword)
		public.POST("/auth/refresh", authHandler.RefreshToken)
		public.POST("/auth/captcha", authHandler.SendCaptcha)
		public.POST("/auth/login-by-code", authHandler.LoginByCode)
		public.GET("/notices", GetNotices)
		public.GET("/notices/:id", GetNoticeDetail)
		public.GET("/help-articles", GetHelpArticles)
		public.GET("/help-articles/:id", GetHelpArticleDetail)
		public.GET("/content/overview", GetContentOverview)
		public.GET("/home-hero", GetClientHomeHero)
		public.POST("/payment/notify/:gateway", PaymentNotify)
		public.POST("/tickets/upstream/replies", TicketUpstreamReply)
		public.POST("/auth/login-api", LoginByAPIKeyHandler) // 锚点自有API密钥登录
	}

	// 需要认证的路由（禁止admin token访问，防越权）
	authenticated := r.Group("")
	authenticated.Use(middleware.JWTAuth(authService))
	authenticated.Use(middleware.ClientRequired())
	{
		// 认证相关
		authenticated.GET("/auth/info", authHandler.GetInfo)
		authenticated.POST("/auth/logout", authHandler.Logout)
		authenticated.PUT("/password", authHandler.UpdatePassword)
		authenticated.PUT("/auth/profile", authHandler.UpdateProfile)
		authenticated.PUT("/auth/phone", authHandler.UpdatePhone)
		authenticated.PUT("/auth/email", authHandler.UpdateEmail)
		authenticated.GET("/auth/notification-preferences", GetNotificationPreferences)
		authenticated.PUT("/auth/notification-preferences", UpdateNotificationPreferences)

		// API密钥管理
		authenticated.GET("/api-key/status", GetAPIKeyStatus)
		authenticated.POST("/api-key/enable", EnableAPIKey)
		authenticated.POST("/api-key/reset", ResetAPIKey)
		authenticated.POST("/api-key/disable", DisableAPIKey)

		// 实名认证
		authenticated.GET("/verification/status", GetVerificationStatus)
		authenticated.POST("/verification/submit", SubmitVerification)

		// 服务管理
		authenticated.GET("/services", GetUserServices)
		authenticated.GET("/services/grouped-overview", GetUserServicesGroupedOverview)
		authenticated.GET("/services/:id", GetUserService)
		authenticated.GET("/services/:id/connection", GetServiceConnection)
		authenticated.GET("/services/:id/runtime", GetServiceRuntime)
		authenticated.PUT("/services/:id/name", UpdateServiceName)
		authenticated.PUT("/services/:id/remark", UpdateServiceRemark)
		authenticated.GET("/services/:id/renewals", GetServiceRenewPreview)
		authenticated.POST("/services/:id/renewals", CreateRenewOrder)
		authenticated.POST("/services/:id/power-actions", PowerService)
		authenticated.POST("/services/:id/password-resets", ResetServicePassword)
		authenticated.POST("/services/:id/reinstallations", ReinstallService)
		authenticated.GET("/services/:id/module-status", GetServiceStatus)
		authenticated.GET("/services/:id/operation-logs", GetServiceOperationLogs)
		authenticated.GET("/services/:id/config", GetServiceConfig)
		authenticated.GET("/services/:id/reinstallations/options", GetServiceReinstallOptions)
		authenticated.GET("/services/:id/upgrades", GetServiceUpgradePreview)
		authenticated.POST("/services/:id/upgrades/quotes", QuoteServiceUpgrade)
		authenticated.POST("/services/:id/upgrades/orders", CreateServiceUpgradeOrder)
		authenticated.PUT("/services/:id/renewals/auto", UpdateAutoRenew)

		// 产品浏览
		authenticated.GET("/products", GetClientProducts)
		authenticated.GET("/products/categories", GetClientProductCategories)
		authenticated.GET("/products/:id", GetClientProductDetail)

		// 订单管理
		authenticated.GET("/orders", GetUserOrders)
		authenticated.GET("/orders/summary", GetOrderSummary)
		authenticated.GET("/orders/:id", GetUserOrder)
		authenticated.POST("/orders/:id/cancel", CancelUserOrder)

		// 工单管理
		authenticated.GET("/tickets", GetUserTickets)
		authenticated.GET("/tickets/service-options", GetTicketServiceOptions)
		authenticated.GET("/tickets/:id", GetUserTicket)
		authenticated.GET("/tickets/:id/replies", GetUserTicketReplies)
		authenticated.POST("/tickets", CreateUserTicket)
		authenticated.POST("/tickets/upload-images", UploadTicketImages)
		authenticated.POST("/tickets/:id/reply", ReplyUserTicket)
		authenticated.POST("/tickets/:id/close", CloseUserTicket)

		// 账单管理
		authenticated.GET("/invoices", GetUserInvoices)
		authenticated.GET("/invoices/summary", GetInvoiceSummary)
		authenticated.GET("/invoices/:id", GetUserInvoice)
		authenticated.POST("/invoices/:id/cancellations", CancelUserInvoice)
		authenticated.POST("/invoices/:id/pay/balance", PayInvoiceByBalance)
		authenticated.POST("/invoices/combines", CombineInvoices)
		authenticated.POST("/invoices/:id/fund", FundInvoice)

		// 财务
		authenticated.GET("/balance-logs", GetBalanceLogs)
		authenticated.GET("/balance-logs/summary", GetBalanceLogsSummary)
		authenticated.GET("/recharge/gateways", GetRechargeGateways)
		authenticated.POST("/recharge", CreateRecharge)
		authenticated.GET("/recharge/:paymentNo/status", GetRechargeStatus)

		// 公告系统（需登录）
		authenticated.GET("/notices/unread-count", GetNoticesUnreadCount)
		authenticated.POST("/notices/mark-all-read", MarkAllNoticesRead)

		// 通知
		authenticated.GET("/notifications", GetUserNotifications)
		authenticated.GET("/notifications/unread-count", GetNotificationUnreadCount)
		authenticated.PUT("/notifications/:id/read-state", MarkNotificationRead)
		authenticated.POST("/notifications/mark-all-read", MarkAllNotificationsRead)

		// 支付记录
		authenticated.GET("/payments", GetPaymentList)
		authenticated.GET("/payments/summary", GetPaymentSummary)
		authenticated.GET("/payments/:id", GetPaymentDetail)

		// 购物车
		authenticated.GET("/cart", GetCart)
		authenticated.POST("/cart", AddToCart)
		authenticated.PUT("/cart/:id", UpdateCartItem)
		authenticated.DELETE("/cart/:id", RemoveCartItem)
		authenticated.DELETE("/cart", ClearCart)
		authenticated.POST("/cart/checkout", Checkout)

		// 推介系统
		authenticated.GET("/referral/overview", GetUserReferralOverview)
		authenticated.GET("/referral/rewards", GetUserReferralRewards)
		authenticated.GET("/referral/direct-referrals", GetUserDirectReferrals)
		authenticated.GET("/referral/account-logs", GetUserReferralAccountLogs)
		authenticated.POST("/referral/withdrawals", ApplyReferralWithdrawal)
		authenticated.GET("/referral/withdrawals", GetUserReferralWithdrawals)
	}
}
