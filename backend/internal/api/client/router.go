package client

import (
	"github.com/YeXuanHs/anchor-finance/internal/middleware"
	"github.com/YeXuanHs/anchor-finance/internal/service"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置用户前台路由
func SetupRoutes(r *gin.RouterGroup, authService *service.AuthService) {
	// 创建处理器
	authHandler := NewAuthHandler(authService)

	// 公开路由（不需要认证）
	public := r.Group("")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
		public.POST("/auth/reset-password", authHandler.ResetPassword)
		public.GET("/notices", GetNotices)
		public.GET("/help-articles", GetHelpArticles)
	}

	// 需要认证的路由
	authenticated := r.Group("")
	authenticated.Use(middleware.JWTAuth(authService))
	{
		// 认证相关
		authenticated.GET("/auth/info", authHandler.GetInfo)
		authenticated.POST("/auth/logout", authHandler.Logout)
		authenticated.PUT("/password", authHandler.UpdatePassword)
		authenticated.PUT("/auth/profile", authHandler.UpdateProfile)

		// 实名认证
		authenticated.GET("/verification/status", GetVerificationStatus)
		authenticated.POST("/verification/submit", SubmitVerification)

		// 服务管理
		authenticated.GET("/services", GetUserServices)
		authenticated.GET("/services/grouped-overview", GetUserServicesGroupedOverview)
		authenticated.GET("/services/:id", GetUserService)
		authenticated.PUT("/services/:id/name", UpdateServiceName)
		authenticated.PUT("/services/:id/remark", UpdateServiceRemark)
		authenticated.GET("/services/:id/renewals", GetServiceRenewPreview)
		authenticated.POST("/services/:id/renewals", CreateRenewOrder)
		authenticated.POST("/services/:id/power-actions", PowerService)
		authenticated.POST("/services/:id/password-resets", ResetServicePassword)
		authenticated.POST("/services/:id/reinstallations", ReinstallService)
		authenticated.GET("/services/:id/module-status", GetServiceStatus)

		// 订单管理
		authenticated.GET("/orders", GetUserOrders)
		authenticated.GET("/orders/:id", GetUserOrder)
		authenticated.POST("/orders/:id/cancel", CancelUserOrder)

		// 工单管理
		authenticated.GET("/tickets", GetUserTickets)
		authenticated.GET("/tickets/:id", GetUserTicket)
		authenticated.GET("/tickets/:id/replies", GetUserTicketReplies)
		authenticated.POST("/tickets", CreateUserTicket)
		authenticated.POST("/tickets/:id/reply", ReplyUserTicket)
		authenticated.POST("/tickets/:id/close", CloseUserTicket)

		// 账单管理
		authenticated.GET("/invoices", GetUserInvoices)
		authenticated.GET("/invoices/summary", GetInvoiceSummary)
		authenticated.GET("/invoices/:id", GetUserInvoice)
		authenticated.POST("/invoices/:id/cancellations", CancelUserInvoice)
		authenticated.POST("/invoices/:id/pay/balance", PayInvoiceByBalance)

		// 财务
		authenticated.GET("/balance-logs", GetBalanceLogs)
		authenticated.GET("/recharge/gateways", GetRechargeGateways)
		authenticated.POST("/recharge", CreateRecharge)

		// 优惠券
		authenticated.GET("/coupons", GetUserCoupons)
		authenticated.POST("/coupons/:id/claim", ClaimCoupon)

		// 通知
		authenticated.GET("/notifications", GetUserNotifications)
		authenticated.GET("/notifications/unread-count", GetNotificationUnreadCount)
		authenticated.PUT("/notifications/:id/read-state", MarkNotificationRead)
		authenticated.POST("/notifications/mark-all-read", MarkAllNotificationsRead)

		// 支付记录
		authenticated.GET("/payments", GetPaymentList)
		authenticated.GET("/payments/summary", GetPaymentSummary)
		authenticated.GET("/payments/:id", GetPaymentDetail)

		// 推介系统
		authenticated.GET("/referral/overview", GetUserReferralOverview)
		authenticated.GET("/referral/rewards", GetUserReferralRewards)
		authenticated.POST("/referral/withdrawals", ApplyReferralWithdrawal)
		authenticated.GET("/referral/withdrawals", GetUserReferralWithdrawals)
	}
}
