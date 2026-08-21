package admin

import (
	"github.com/YeXuanHs/anchor-finance/internal/middleware"
	"github.com/YeXuanHs/anchor-finance/internal/service"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置管理后台路由
func SetupRoutes(r *gin.RouterGroup, authService *service.AuthService) {
	// 创建处理器
	authHandler := NewAuthHandler(authService)

	// 公开路由（不需要认证）
	public := r.Group("")
	{
		public.POST("/login", authHandler.Login)
	}

	// 需要认证的路由
	authenticated := r.Group("")
	authenticated.Use(middleware.JWTAuth(authService))
	{
		// 认证相关
		authenticated.GET("/auth/info", authHandler.GetInfo)
		authenticated.PUT("/auth/profile", authHandler.UpdateProfile)
		authenticated.PUT("/auth/password", authHandler.UpdatePassword)
		authenticated.POST("/logout", authHandler.Logout)

		// 仪表盘 - 已实现
		authenticated.GET("/dashboard/stats", GetDashboardStats)
		authenticated.GET("/dashboard/income-trend", GetIncomeTrend)
		authenticated.GET("/dashboard/online-admins", GetOnlineAdmins)
		authenticated.GET("/dashboard/recent-invoices", GetRecentInvoices)
		authenticated.GET("/dashboard/monthly-revenue", GetMonthlyRevenue)

		// 客户管理 - 已实现
		authenticated.GET("/os-options", GetOSOptions)
		authenticated.GET("/users", GetUserList)
		authenticated.GET("/users/:id", GetUser)
		authenticated.POST("/users", CreateUser)
		authenticated.PUT("/users/:id", UpdateUser)
		authenticated.DELETE("/users/:id", DeleteUser)
		authenticated.PATCH("/users/:id/status", UpdateUserStatus)
		authenticated.GET("/users/:id/orders", GetUserOrders)
		authenticated.GET("/users/:id/invoices", GetUserInvoices)
		authenticated.GET("/users/:id/tickets", GetUserTickets)
		authenticated.GET("/users/:id/services", GetUserServices)
		authenticated.GET("/users/:id/balance-logs", GetUserBalanceLogs)
		authenticated.GET("/users/:id/operation-logs", GetUserOperationLogs)
		authenticated.POST("/users/:id/recharges", RechargeUser)
		authenticated.POST("/users/:id/services/:service_id/refunds", RefundUserService)

		// 订单管理 - 已实现
		authenticated.GET("/orders", GetOrderList)
		authenticated.GET("/orders/:id", GetOrder)
		authenticated.POST("/orders", CreateOrder)
		authenticated.PUT("/orders/:id", UpdateOrder)
		authenticated.POST("/orders/:id/activate", ActivateOrder)
		authenticated.POST("/orders/:id/cancel", CancelOrder)

		// 服务管理 - 已实现
		authenticated.GET("/services", GetServiceList)
		authenticated.GET("/services/:id", GetService)
		authenticated.PUT("/services/:id", UpdateService)
		authenticated.POST("/services/:id/suspend", SuspendService)
		authenticated.POST("/services/:id/unsuspend", UnsuspendService)
		authenticated.POST("/services/:id/terminate", TerminateService)

		// 账单管理 - 已实现
		authenticated.GET("/invoices", GetInvoiceList)
		authenticated.GET("/invoices/:id", GetInvoice)
		authenticated.POST("/invoices/:id/cancel", CancelInvoice)
		authenticated.GET("/transactions", GetTransactionList)

		// 财务报表 - 已实现
		authenticated.GET("/finance/new-customer-daily-summary", GetNewCustomerDailySummary)
		authenticated.GET("/finance/product-income-summary", GetProductIncomeSummary)
		authenticated.GET("/finance/ledger", GetFinanceLedger)

		// 工单管理 - 已实现
		authenticated.GET("/tickets", GetTicketList)
		authenticated.GET("/tickets/summary", GetTicketSummary)
		authenticated.GET("/tickets/:id", GetTicket)
		authenticated.GET("/tickets/:id/replies", GetTicketReplies)
		authenticated.POST("/tickets/:id/reply", ReplyTicket)
		authenticated.POST("/tickets/:id/close", CloseTicket)
		authenticated.POST("/tickets/:id/reopen", ReopenTicket)
		authenticated.PUT("/tickets/:id/assignment", AssignTicket)
		authenticated.GET("/ticket-departments", GetTicketDepartments)
		authenticated.GET("/ticket-statuses", GetTicketStatuses)

		// 产品管理 - 已实现
		authenticated.GET("/products", GetProductList)
		authenticated.GET("/products/summary", GetProductSummary)
		authenticated.GET("/products/:id", GetProduct)
		authenticated.GET("/products/:id/owners", GetProductOwners)
		authenticated.POST("/products", CreateProduct)
		authenticated.PUT("/products/:id", UpdateProduct)
		authenticated.DELETE("/products/:id", DeleteProduct)
		authenticated.PATCH("/products/:id/status", UpdateProductStatus)
		authenticated.GET("/product-groups", GetProductGroups)
		authenticated.GET("/product-groups/tree", GetProductGroupTree)
		authenticated.POST("/product-groups", CreateProductGroup)
		authenticated.PUT("/product-groups/:id", UpdateProductGroup)
		authenticated.DELETE("/product-groups/:id", DeleteProductGroup)
		authenticated.GET("/product-types", GetProductTypeList)
		authenticated.POST("/product-types", CreateProductType)
		authenticated.PUT("/product-types/:id", UpdateProductType)
		authenticated.DELETE("/product-types/:id", DeleteProductType)

		// 设置管理 - 已实现
		authenticated.GET("/settings", GetSettings)
		authenticated.GET("/settings/:group", GetSettingsByGroup)
		authenticated.PUT("/settings", UpdateSettings)

		// 通知模板 - 已实现
		authenticated.GET("/notification-templates", GetNotificationTemplates)
		authenticated.POST("/notification-templates", CreateNotificationTemplate)
		authenticated.PUT("/notification-templates/:id", UpdateNotificationTemplate)
		authenticated.DELETE("/notification-templates/:id", DeleteNotificationTemplate)

		// 菜单管理 - 已实现
		authenticated.GET("/menus", GetMenus)
		authenticated.POST("/menus", CreateMenu)
		authenticated.PUT("/menus/:id", UpdateMenu)
		authenticated.DELETE("/menus/:id", DeleteMenu)

		// 管理员管理 - 已实现
		authenticated.GET("/admins", GetAdminList)
		authenticated.POST("/admins", CreateAdmin)
		authenticated.PUT("/admins/:id", UpdateAdmin)

		// 员工管理 - 已实现
		authenticated.GET("/staff", GetStaffList)
		authenticated.GET("/staff/:id", GetStaffDetail)
		authenticated.POST("/staff", CreateStaff)
		authenticated.PUT("/staff/:id", UpdateStaff)
		authenticated.DELETE("/staff/:id", DeleteStaff)
		authenticated.PATCH("/staff/:id/status", UpdateStaffStatus)
		authenticated.POST("/staff/:id/password-resets", ResetStaffPassword)

		// 角色管理 - 已实现
		authenticated.GET("/roles", GetRoleList)
		authenticated.GET("/roles/:id", GetRoleDetail)
		authenticated.POST("/roles", CreateRole)
		authenticated.PUT("/roles/:id", UpdateRole)
		authenticated.DELETE("/roles/:id", DeleteRole)
		authenticated.GET("/permissions", GetPermissions)

		// 定时任务
		authenticated.GET("/cron-tasks", GetCronTasks)

		// 数据库管理 - 已实现
		authenticated.GET("/database/status", GetDatabaseStatus)
		authenticated.POST("/database/optimizations", OptimizeDatabase)
		authenticated.POST("/database/backups", BackupDatabase)

		// 系统信息 - 已实现
		authenticated.GET("/system/info", GetSystemInfo)
		authenticated.GET("/system/modules", GetSystemModules)

		// 日志管理 - 已实现
		authenticated.GET("/system-logs", GetSystemLogs)
		authenticated.GET("/operation-logs", GetOperationLogs)
		authenticated.GET("/login-logs", GetLoginLogs)

		// 内容管理 - 已实现
		authenticated.GET("/content/summary", GetContentSummary)
		authenticated.GET("/news", GetNewsList)
		authenticated.POST("/news", CreateNews)
		authenticated.PUT("/news/:id", UpdateNews)
		authenticated.DELETE("/news/:id", DeleteNews)
		authenticated.GET("/news-categories", GetNewsCategories)
		authenticated.POST("/news-categories", CreateNewsCategory)
		authenticated.PUT("/news-categories/:id", UpdateNewsCategory)
		authenticated.DELETE("/news-categories/:id", DeleteNewsCategory)
		authenticated.GET("/knowledge/categories", GetKnowledgeCategories)
		authenticated.POST("/knowledge/categories", CreateKnowledgeCategory)
		authenticated.GET("/knowledge/articles", GetKnowledgeArticles)
		authenticated.POST("/knowledge/articles", CreateKnowledgeArticle)
		authenticated.PUT("/knowledge/articles/:id", UpdateKnowledgeArticle)
		authenticated.DELETE("/knowledge/articles/:id", DeleteKnowledgeArticle)
		authenticated.GET("/downloads", GetDownloads)
		authenticated.POST("/downloads", CreateDownload)
		authenticated.PUT("/downloads/:id", UpdateDownload)
		authenticated.DELETE("/downloads/:id", DeleteDownload)
		authenticated.GET("/downloads/categories", GetDownloadCategories)
		authenticated.POST("/downloads/categories", CreateDownloadCategory)

		// 货币管理 - 已实现
		authenticated.GET("/currencies", GetCurrencyList)
		authenticated.POST("/currencies", CreateCurrency)
		authenticated.PUT("/currencies/:id", UpdateCurrency)

		// 实名认证 - 已实现
		authenticated.GET("/verifications", GetVerificationList)
		authenticated.GET("/verifications/summary", GetVerificationSummary)
		authenticated.GET("/verifications/:id", GetVerificationDetail)
		authenticated.POST("/verifications/:id/approve", ApproveVerification)
		authenticated.POST("/verifications/:id/reject", RejectVerification)

		// 供应商管理 - 已实现
		authenticated.GET("/suppliers", GetSupplierList)
		authenticated.GET("/suppliers/summary", GetSupplierSummary)
		authenticated.GET("/suppliers/:id", GetSupplierDetail)
		authenticated.POST("/suppliers", CreateSupplier)
		authenticated.PUT("/suppliers/:id", UpdateSupplier)
		authenticated.DELETE("/suppliers/:id", DeleteSupplier)
		authenticated.GET("/suppliers/:id/products", GetSupplierProducts)

		// 插件管理 - 已实现
		authenticated.GET("/plugins", GetPluginList)
		authenticated.GET("/plugins/:id", GetPluginDetail)
		authenticated.POST("/plugins/:id/enable", EnablePlugin)
		authenticated.POST("/plugins/:id/disable", DisablePlugin)
		authenticated.DELETE("/plugins/:id", UninstallPlugin)
		authenticated.GET("/plugins/:id/config", GetPluginConfig)
		authenticated.PUT("/plugins/:id/config", UpdatePluginConfig)

		// 优惠码管理 - 已实现
		authenticated.GET("/promo-codes", GetPromoCodeList)
		authenticated.POST("/promo-codes", CreatePromoCode)
		authenticated.PUT("/promo-codes/:id", UpdatePromoCode)
		authenticated.DELETE("/promo-codes/:id", DeletePromoCode)

		// 推介系统 - 已实现
		authenticated.GET("/referral/overview", GetReferralOverview)
		authenticated.GET("/referral/rewards", GetReferralRewards)
		authenticated.GET("/referral-withdrawals", GetReferralWithdrawals)
		authenticated.POST("/referral-withdrawals/:id/approve", ApproveReferralWithdrawal)
		authenticated.POST("/referral-withdrawals/:id/reject", RejectReferralWithdrawal)

		// 会员等级 - 已实现
		authenticated.GET("/member-levels", GetMemberLevelList)
		authenticated.POST("/member-levels", CreateMemberLevel)
		authenticated.PUT("/member-levels/:id", UpdateMemberLevel)
		authenticated.DELETE("/member-levels/:id", DeleteMemberLevel)

		// 自定义字段 - 已实现
		authenticated.GET("/custom-fields", GetCustomFieldList)
		authenticated.POST("/custom-fields", CreateCustomField)
		authenticated.PUT("/custom-fields/:id", UpdateCustomField)
		authenticated.DELETE("/custom-fields/:id", DeleteCustomField)

		// 优惠券 - 已实现
		authenticated.GET("/coupons", GetCouponList)
		authenticated.GET("/coupons/summary", GetCouponSummary)
		authenticated.POST("/coupons", CreateCoupon)
		authenticated.PUT("/coupons/:id", UpdateCoupon)
		authenticated.DELETE("/coupons/:id", DeleteCoupon)
		authenticated.PATCH("/coupons/:id/status", UpdateCouponStatus)

		// 优惠券活动 - 已实现
		authenticated.GET("/coupon-campaigns", GetCouponCampaignList)
		authenticated.GET("/coupon-campaigns/summary", GetCouponCampaignSummary)
		authenticated.POST("/coupon-campaigns", CreateCouponCampaign)
		authenticated.PUT("/coupon-campaigns/:id", UpdateCouponCampaign)
		authenticated.DELETE("/coupon-campaigns/:id", DeleteCouponCampaign)
		authenticated.PATCH("/coupon-campaigns/:id/status", UpdateCouponCampaignStatus)

		// TODO: 以下功能待实现

		// 服务管理
		// authenticated.GET("/services", getServices)
		// authenticated.GET("/services/:id", getService)

		// 账单管理
		// authenticated.GET("/invoices", getInvoices)
		// authenticated.GET("/invoices/:id", getInvoice)

		// 工单管理
		// authenticated.GET("/tickets", getTickets)
		// authenticated.GET("/tickets/:id", getTicket)

		// 产品管理
		// authenticated.GET("/products", getProducts)
		// authenticated.GET("/products/:id", getProduct)

		// 插件管理
		// authenticated.GET("/plugins", getPlugins)

		// 设置
		// authenticated.GET("/settings", getSettings)
		// authenticated.PUT("/settings", updateSettings)

		// 菜单
		// authenticated.GET("/menus", getMenus)
	}
}


