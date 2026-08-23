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
		v1Auth.GET("/host/header", ZjmfCompatHostDetail)
	}
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

		// 优惠券
		authenticated.GET("/coupons", GetUserCoupons)
		authenticated.POST("/coupons/:id/claim", ClaimCoupon)

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
