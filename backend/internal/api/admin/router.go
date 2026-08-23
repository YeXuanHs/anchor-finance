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

	// 需要认证的路由（必须是管理员）
	authenticated := r.Group("")
	authenticated.Use(middleware.JWTAuth(authService))
	authenticated.Use(middleware.AdminRequired())
	{
		// 认证相关
		authenticated.GET("/auth/info", authHandler.GetInfo)
		authenticated.PUT("/auth/profile", authHandler.UpdateProfile)
		authenticated.PUT("/auth/password", authHandler.UpdatePassword)
		authenticated.POST("/logout", authHandler.Logout)
		authenticated.POST("/auth/reset-password", authHandler.ResetAdminPassword)

		// 仪表盘
		authenticated.GET("/dashboard/stats", GetDashboardStats)
		authenticated.GET("/dashboard/income-trend", GetIncomeTrend)
		authenticated.GET("/dashboard/online-admins", GetOnlineAdmins)
		authenticated.GET("/dashboard/recent-invoices", GetRecentInvoices)
		authenticated.GET("/dashboard/monthly-revenue", GetMonthlyRevenue)

		// 客户管理
		authenticated.GET("/os-options", GetOSOptions)
		authenticated.GET("/cpu-model-catalog", GetCPUModelCatalog)
		authenticated.GET("/instance-spec-catalog", GetInstanceSpecCatalog)
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
		authenticated.GET("/users/:id/email-logs", GetUserEmailLogs)
		authenticated.GET("/users/:id/sms-logs", GetUserSmsLogs)
		authenticated.GET("/users/:id/invoices/:invoice_id", GetUserInvoiceDetail)
		authenticated.POST("/users/:id/recharges", RechargeUser)
		authenticated.POST("/users/:id/services/:service_id/refunds", RefundUserService)
		authenticated.GET("/users/:id/services/:service_id/connection", AdminGetServiceConnection)
		authenticated.GET("/users/:id/services/:service_id/remote-status", AdminGetServiceRemoteStatus)
		authenticated.PUT("/users/:id/services/:service_id/meta", AdminUpdateServiceMeta)
		authenticated.POST("/users/:id/services/:service_id/manual-provision", AdminManualProvision)
		authenticated.POST("/users/:id/services/:service_id/power-actions", AdminServicePowerAction)
		authenticated.POST("/users/:id/services/:service_id/password-resets", AdminResetServicePassword)
		// 新增：用户备注、登录为用户、刷新服务状态
		authenticated.GET("/users/:id/remarks", GetUserRemarks)
		authenticated.POST("/users/:id/remarks", AddUserRemark)
		authenticated.POST("/users/:id/login-as", LoginAsUser)
		authenticated.POST("/users/:id/services/refresh-statuses", RefreshUserServicesStatus)

		// 订单管理
		authenticated.GET("/orders", GetOrderList)
		authenticated.POST("/orders/search", SearchOrders)
		authenticated.GET("/orders/:id", GetOrder)
		authenticated.POST("/orders", CreateOrder)
		authenticated.PUT("/orders/:id", UpdateOrder)
		authenticated.POST("/orders/:id/activate", ActivateOrder)
		authenticated.POST("/orders/:id/cancel", CancelOrder)
		authenticated.POST("/orders/:id/notes", AddOrderNote)

		// 服务管理
		authenticated.GET("/services", GetServiceList)
		authenticated.GET("/services/:id", GetService)
		authenticated.PUT("/services/:id", UpdateService)
		authenticated.POST("/services/:id/suspend", SuspendService)
		authenticated.POST("/services/:id/unsuspend", UnsuspendService)
		authenticated.POST("/services/:id/terminate", TerminateService)

		// 账单管理
		authenticated.GET("/invoices", GetInvoiceList)
		authenticated.GET("/invoices/:id", GetInvoice)
		authenticated.POST("/invoices/:id/cancel", CancelInvoice)
		authenticated.POST("/invoices/:id/notes", AddInvoiceNote)
		authenticated.GET("/transactions", GetTransactionList)

		// 信用额管理
		authenticated.GET("/credit-limits", GetCreditLimitList)
		authenticated.GET("/credit-limits/config", GetCreditLimitConfig)
		authenticated.POST("/credit-limits/config", SaveCreditLimitConfig)
		authenticated.POST("/credit-limits", SaveCreditLimit)
		authenticated.PUT("/credit-limits/:id", UpdateCreditLimit)
		authenticated.DELETE("/credit-limits/:id", DeleteCreditLimit)
		authenticated.GET("/credit-limits/logs", GetCreditLimitLogs)

		// 财务报表
		authenticated.GET("/finance/new-customer-daily-summary", GetNewCustomerDailySummary)
		authenticated.GET("/finance/product-income-summary", GetProductIncomeSummary)
		authenticated.GET("/finance/ledger", GetFinanceLedger)
		authenticated.GET("/finance/ledger/:id", GetFinanceLedgerDetail)
		authenticated.GET("/finance/ledger/summary", GetFinanceLedgerSummary)
		authenticated.GET("/finance/recharges", GetRechargeList)
		authenticated.GET("/finance/recharges/summary", GetRechargeSummary)
		authenticated.GET("/finance/renewal-orders", GetRenewalOrders)
		authenticated.GET("/finance/upgrade-orders", GetUpgradeOrders)

		// 报表
		authenticated.GET("/reports/new-customers", GetNewCustomerStatistics)
		authenticated.GET("/reports/revenue-ranking", GetRevenueRanking)

		// 工单管理
		authenticated.GET("/tickets", GetTicketList)
		authenticated.GET("/tickets/summary", GetTicketSummary)
		authenticated.GET("/tickets/admin-users", GetTicketAdminUsers)
		authenticated.GET("/tickets/:id", GetTicket)
		authenticated.GET("/tickets/:id/replies", GetTicketReplies)
		authenticated.POST("/tickets/:id/reply", ReplyTicket)
		authenticated.POST("/tickets/:id/close", CloseTicket)
		authenticated.POST("/tickets/:id/reopen", ReopenTicket)
		authenticated.POST("/tickets/:id/receive", ReceiveTicket)
		authenticated.PUT("/tickets/:id/assignment", AssignTicket)
		authenticated.POST("/tickets/:id/replies/:reply_id/recalls", RecallTicketReply)
		authenticated.GET("/ticket-departments", GetTicketDepartments)
		authenticated.GET("/ticket-departments/:id", GetTicketDepartmentDetail)
		authenticated.POST("/ticket-departments/:id/move-up", MoveTicketDepartmentUp)
		authenticated.POST("/ticket-departments/:id/move-down", MoveTicketDepartmentDown)
		authenticated.GET("/ticket-statuses", GetTicketStatuses)
		authenticated.GET("/ticket-statuses/:id", GetTicketStatusDetail)

		// 工单上游投递（参考图拉财务TicketDeliveryService）
		authenticated.GET("/ticket-delivery-rules", GetTicketDeliveryRules)
		authenticated.POST("/ticket-delivery-rules", CreateTicketDeliveryRule)
		authenticated.PUT("/ticket-delivery-rules/:id", UpdateTicketDeliveryRule)
		authenticated.DELETE("/ticket-delivery-rules/:id", DeleteTicketDeliveryRule)
		authenticated.GET("/tickets/:id/upstream-delivery", GetTicketUpstreamDelivery)
		authenticated.GET("/tickets/:id/upstream-delivery/logs", GetTicketUpstreamDeliveryLogs)

		// 工单预回复
		authenticated.GET("/ticket-prereplies", GetTicketPrereplyList)
		authenticated.POST("/ticket-prereplies", CreateTicketPrereply)
		authenticated.PUT("/ticket-prereplies/:id", UpdateTicketPrereply)
		authenticated.DELETE("/ticket-prereplies/:id", DeleteTicketPrereply)
		authenticated.POST("/ticket-prereplies/search", SearchTicketPrereply)
		authenticated.GET("/ticket-prereply-categories", GetTicketPrereplyCategoryList)
		authenticated.POST("/ticket-prereply-categories", CreateTicketPrereplyCategory)
		authenticated.PUT("/ticket-prereply-categories/:id", UpdateTicketPrereplyCategory)
		authenticated.DELETE("/ticket-prereply-categories/:id", DeleteTicketPrereplyCategory)

		// 产品管理
		authenticated.GET("/products", GetProductList)
		authenticated.GET("/products/summary", GetProductSummary)
		authenticated.GET("/products/:id", GetProduct)
		authenticated.GET("/products/:id/owners", GetProductOwners)
		authenticated.POST("/products", CreateProduct)
		authenticated.PUT("/products/:id", UpdateProduct)
		authenticated.DELETE("/products/:id", DeleteProduct)
		authenticated.POST("/products/:id/restorations", RestoreProduct)
		authenticated.PATCH("/products/:id/status", UpdateProductStatus)
		authenticated.POST("/products/reorders", ReorderProducts)
		authenticated.POST("/products/category-batches", BatchUpdateProductCategory)
		authenticated.GET("/product-groups", GetProductGroups)
		authenticated.GET("/product-groups/tree", GetProductGroupTree)
		authenticated.GET("/product-groups/:id", GetProductGroupDetail)
		authenticated.GET("/product-groups/:id/children", GetProductGroupChildren)
		authenticated.POST("/product-groups", CreateProductGroup)
		authenticated.PUT("/product-groups/:id", UpdateProductGroup)
		authenticated.DELETE("/product-groups/:id", DeleteProductGroup)
		authenticated.POST("/product-groups/reorders", ReorderProductGroups)
		authenticated.DELETE("/products/:id/force", ForceDeleteProduct)
		authenticated.GET("/product-types", GetProductTypeList)
		authenticated.POST("/product-types", CreateProductType)
		authenticated.PUT("/product-types/:id", UpdateProductType)
		authenticated.DELETE("/product-types/:id", DeleteProductType)
		authenticated.POST("/product-types/reorders", ReorderProductTypes)

		// 设置管理
		authenticated.GET("/settings", GetSettings)
		authenticated.GET("/settings/:group", GetSettingsByGroup)
		authenticated.PUT("/settings", UpdateSettings)
		authenticated.GET("/settings/email", GetEmailConfig)
		authenticated.PUT("/settings/email", UpdateEmailConfig)
		authenticated.GET("/settings/sms", GetSMSConfig)
		authenticated.PUT("/settings/sms", UpdateSMSConfig)
		authenticated.GET("/settings/register-login", GetRegisterLoginConfig)
		authenticated.PUT("/settings/register-login", UpdateRegisterLoginConfig)
		authenticated.GET("/settings/captcha", GetCaptchaConfig)
		authenticated.PUT("/settings/captcha", UpdateCaptchaConfig)
		authenticated.GET("/settings/security", GetSecurityConfig)
		authenticated.PUT("/settings/security", UpdateSecurityConfig)
		authenticated.GET("/settings/general", GetGeneralConfig)
		authenticated.PUT("/settings/general", UpdateGeneralConfig)
		authenticated.GET("/settings/display", GetDisplayConfig)
		authenticated.PUT("/settings/display", UpdateDisplayConfig)
		authenticated.GET("/settings/invoice", GetInvoiceConfig)
		authenticated.PUT("/settings/invoice", UpdateInvoiceConfig)
		authenticated.GET("/settings/contract", GetContractConfig)
		authenticated.PUT("/settings/contract", UpdateContractConfig)
		authenticated.GET("/settings/credit-setting", GetCreditSettingConfig)
		authenticated.PUT("/settings/credit-setting", UpdateCreditSettingConfig)
		authenticated.GET("/settings/payment-gateway", GetPaymentGatewayConfig)
		authenticated.PUT("/settings/payment-gateway", UpdatePaymentGatewayConfig)

		// 通知模板
		authenticated.GET("/notification-templates", GetNotificationTemplates)
		authenticated.POST("/notification-templates", CreateNotificationTemplate)
		authenticated.POST("/notification-templates/test-send", TestNotificationTemplate)
		authenticated.PUT("/notification-templates/:id", UpdateNotificationTemplate)
		authenticated.DELETE("/notification-templates/:id", DeleteNotificationTemplate)

		// 菜单管理
		authenticated.GET("/menus", GetMenus)
		authenticated.GET("/menu-types", GetMenuTypeList)
		authenticated.POST("/menus", CreateMenu)
		authenticated.PUT("/menus/:id", UpdateMenu)
		authenticated.DELETE("/menus/:id", DeleteMenu)

		// 管理员管理
		authenticated.GET("/admins", GetAdminList)
		authenticated.POST("/admins", CreateAdmin)
		authenticated.PUT("/admins/:id", UpdateAdmin)

		// 员工管理
		authenticated.GET("/staff", GetStaffList)
		authenticated.GET("/staff/roles", GetStaffRoles)
		authenticated.GET("/staff/:id", GetStaffDetail)
		authenticated.POST("/staff", CreateStaff)
		authenticated.PUT("/staff/:id", UpdateStaff)
		authenticated.DELETE("/staff/:id", DeleteStaff)
		authenticated.PATCH("/staff/:id/status", UpdateStaffStatus)
		authenticated.POST("/staff/:id/password-resets", ResetStaffPassword)

		// 角色管理
		authenticated.GET("/roles", GetRoleList)
		authenticated.GET("/roles/:id", GetRoleDetail)
		authenticated.POST("/roles", CreateRole)
		authenticated.PUT("/roles/:id", UpdateRole)
		authenticated.DELETE("/roles/:id", DeleteRole)
		authenticated.POST("/roles/:id/copies", CopyRole)
		authenticated.GET("/permissions", GetPermissions)

		// 定时任务
		authenticated.GET("/cron-tasks", GetCronTasks)
		authenticated.GET("/schedules/overview", GetScheduleOverview)
		authenticated.GET("/schedule-runs", GetScheduleRunList)
		authenticated.GET("/schedule-runs/:id", GetScheduleRunDetail)
		authenticated.POST("/schedule-runs/:id/retry", RetryScheduleRun)
		authenticated.POST("/schedule-triggers", TriggerSchedule)

		// 数据库管理
		authenticated.GET("/database/status", GetDatabaseStatus)
		authenticated.POST("/database/optimizations", OptimizeDatabase)
		authenticated.POST("/database/backups", BackupDatabase)

		// 系统信息
		authenticated.GET("/system/info", GetSystemInfo)
		authenticated.GET("/system/modules", GetSystemModules)

		// 日志管理
		authenticated.GET("/system-logs", GetSystemLogs)
		authenticated.GET("/operation-logs", GetOperationLogs)
		authenticated.GET("/login-logs", GetLoginLogs)
		authenticated.GET("/log-cleanups/overview", GetLogCleanupOverview)
		authenticated.POST("/log-cleanups", CleanupLogs)
		authenticated.GET("/log-summaries/:channel", GetLogSummaryByChannel)

		// 分类日志
		authenticated.GET("/logs/sms", GetSMSLogs)
		authenticated.GET("/logs/email", GetEmailLogs)
		authenticated.GET("/logs/api", GetAPILogs)
		authenticated.GET("/logs/cron", GetCronLogs)
		authenticated.GET("/logs/admin-login", GetAdminLoginLogs)
		authenticated.GET("/logs/notification", GetNotificationLogs)

		// 内容管理
		authenticated.GET("/content/summary", GetContentSummary)

		// 主题模板 - 首页Hero（修正路径：MD文档定义 /api/admin/home-hero）
		authenticated.GET("/home-hero", GetHomeHero)
		authenticated.PUT("/home-hero", UpdateHomeHero)
		authenticated.GET("/home-hero/assets", GetHomeHeroAssets)

		authenticated.POST("/upload", UploadFile)
		authenticated.GET("/media-files", GetMediaFileList)
		authenticated.DELETE("/media-files/:id", DeleteMediaFile)
		authenticated.GET("/media-files/:id/references", GetMediaFileReferences)
		authenticated.GET("/news", GetNewsList)
		authenticated.GET("/news/:id", GetNewsDetail)
		authenticated.POST("/news", CreateNews)
		authenticated.PUT("/news/:id", UpdateNews)
		authenticated.DELETE("/news/:id", DeleteNews)
		authenticated.GET("/news-categories", GetNewsCategories)
		authenticated.POST("/news-categories", CreateNewsCategory)
		authenticated.PUT("/news-categories/:id", UpdateNewsCategory)
		authenticated.DELETE("/news-categories/:id", DeleteNewsCategory)
		authenticated.GET("/knowledge/categories", GetKnowledgeCategories)
		authenticated.POST("/knowledge/categories", CreateKnowledgeCategory)
		authenticated.PUT("/knowledge/categories/:id", UpdateKnowledgeCategory)
		authenticated.DELETE("/knowledge/categories/:id", DeleteKnowledgeCategory)
		authenticated.GET("/knowledge/articles", GetKnowledgeArticles)
		authenticated.GET("/knowledge/articles/:id", GetKnowledgeArticleDetail)
		authenticated.POST("/knowledge/articles", CreateKnowledgeArticle)
		authenticated.PUT("/knowledge/articles/:id", UpdateKnowledgeArticle)
		authenticated.DELETE("/knowledge/articles/:id", DeleteKnowledgeArticle)
		authenticated.GET("/downloads", GetDownloads)
		authenticated.POST("/downloads", CreateDownload)
		authenticated.PUT("/downloads/:id", UpdateDownload)
		authenticated.DELETE("/downloads/:id", DeleteDownload)
		authenticated.GET("/downloads/categories", GetDownloadCategories)
		authenticated.POST("/downloads/categories", CreateDownloadCategory)
		authenticated.PUT("/downloads/categories/:id", UpdateDownloadCategory)
		authenticated.DELETE("/downloads/categories/:id", DeleteDownloadCategory)

		// 货币管理
		authenticated.GET("/currencies", GetCurrencyList)
		authenticated.POST("/currencies", CreateCurrency)
		authenticated.PUT("/currencies/:id", UpdateCurrency)

		// 实名认证
		authenticated.GET("/verifications", GetVerificationList)
		authenticated.GET("/verifications/summary", GetVerificationSummary)
		authenticated.GET("/verifications/:id", GetVerificationDetail)
		authenticated.GET("/verifications/:id/history", GetVerificationHistory)
		authenticated.POST("/verifications/:id/approve", ApproveVerification)
		authenticated.POST("/verifications/:id/reject", RejectVerification)
		authenticated.POST("/verifications/:id/unbindings", UnbindVerificationByUser)
		authenticated.POST("/users/:id/unbind-verification", UnbindVerification)

		// 供应商管理
		authenticated.GET("/suppliers", GetSupplierList)
		authenticated.GET("/suppliers/summary", GetSupplierSummary)
		authenticated.GET("/suppliers/provider-types", GetSupplierProviderTypes)
		authenticated.GET("/suppliers/:id", GetSupplierDetail)
		authenticated.GET("/suppliers/:id/balance", GetSupplierBalance)
		authenticated.PATCH("/suppliers/:id/status", UpdateSupplierStatus)
		authenticated.POST("/suppliers", CreateSupplier)
		authenticated.PUT("/suppliers/:id", UpdateSupplier)
		authenticated.DELETE("/suppliers/:id", DeleteSupplier)
		authenticated.GET("/suppliers/:id/products", GetSupplierProducts)
		authenticated.POST("/suppliers/:id/tasks", RunSupplierTask)
		authenticated.POST("/suppliers/:id/sync-products", SyncSupplierProducts)
		authenticated.POST("/suppliers/:id/sync-prices", SyncSupplierPrices)
		authenticated.POST("/suppliers/:id/sync-stock", SyncSupplierStock)
		authenticated.GET("/suppliers/:id/secrets/:key", RevealSupplierSecret)

		// 供应商分组映射（MD 7.2.5）
		authenticated.GET("/suppliers/:id/group-mappings", GetSupplierGroupMappings)
		authenticated.POST("/suppliers/:id/group-mappings", CreateSupplierGroupMapping)
		authenticated.PUT("/suppliers/:id/group-mappings/:mapping_id", UpdateSupplierGroupMapping)
		authenticated.DELETE("/suppliers/:id/group-mappings/:mapping_id", DeleteSupplierGroupMapping)

		// 插件管理
		authenticated.GET("/plugins", GetPluginList)
		authenticated.POST("/plugins/install", InstallPlugin)
		authenticated.POST("/plugins/scan", ScanPlugins)
		authenticated.GET("/plugins/:id", GetPluginDetail)
		authenticated.POST("/plugins/:id/enable", EnablePlugin)
		authenticated.POST("/plugins/:id/disable", DisablePlugin)
		authenticated.DELETE("/plugins/:id", UninstallPlugin)
		authenticated.GET("/plugins/:id/config", GetPluginConfig)
		authenticated.PUT("/plugins/:id/config", UpdatePluginConfig)
		authenticated.POST("/plugins/:id/health", PluginHealthCheck)

		// 插件域API
		authenticated.GET("/payment-gateways", GetPaymentGateways)
		authenticated.GET("/sms-providers", GetSMSProviders)
		authenticated.GET("/mail-providers", GetMailProviders)
		authenticated.GET("/certification-providers", GetCertificationProviders)
		authenticated.GET("/server-modules", GetServerModules)

		// 优惠码管理
		authenticated.GET("/promo-codes", GetPromoCodeList)
		authenticated.POST("/promo-codes", CreatePromoCode)
		authenticated.PUT("/promo-codes/:id", UpdatePromoCode)
		authenticated.DELETE("/promo-codes/:id", DeletePromoCode)

		// 推介系统
		authenticated.GET("/referral/overview", GetReferralOverview)
		authenticated.GET("/referral/rewards", GetReferralRewards)
		authenticated.GET("/referral-withdrawals", GetReferralWithdrawals)
		authenticated.POST("/referral-withdrawals/:id/approve", ApproveReferralWithdrawal)
		authenticated.POST("/referral-withdrawals/:id/reject", RejectReferralWithdrawal)

		// 会员等级
		authenticated.GET("/member-levels", GetMemberLevelList)
		authenticated.POST("/member-levels", CreateMemberLevel)
		authenticated.PUT("/member-levels/:id", UpdateMemberLevel)
		authenticated.DELETE("/member-levels/:id", DeleteMemberLevel)

		// 自定义字段
		authenticated.GET("/custom-fields", GetCustomFieldList)
		authenticated.POST("/custom-fields", CreateCustomField)
		authenticated.PUT("/custom-fields/:id", UpdateCustomField)
		authenticated.DELETE("/custom-fields/:id", DeleteCustomField)

		// 优惠券
		authenticated.GET("/coupons", GetCouponList)
		authenticated.GET("/coupons/summary", GetCouponSummary)
		authenticated.GET("/coupon-product-groups", GetCouponProductGroups)
		authenticated.POST("/coupons", CreateCoupon)
		authenticated.PUT("/coupons/:id", UpdateCoupon)
		authenticated.DELETE("/coupons/:id", DeleteCoupon)
		authenticated.PATCH("/coupons/:id/status", UpdateCouponStatus)

		// 优惠券活动
		authenticated.GET("/coupon-campaigns", GetCouponCampaignList)
		authenticated.GET("/coupon-campaigns/summary", GetCouponCampaignSummary)
		authenticated.POST("/coupon-campaigns", CreateCouponCampaign)
		authenticated.PUT("/coupon-campaigns/:id", UpdateCouponCampaign)
		authenticated.DELETE("/coupon-campaigns/:id", DeleteCouponCampaign)
		authenticated.PATCH("/coupon-campaigns/:id/status", UpdateCouponCampaignStatus)

		// 发送消息
		authenticated.GET("/send-message/search-params", GetSendMessageSearchParams)
		authenticated.GET("/send-message/send-methods", GetSendMethodList)
		authenticated.GET("/send-message/search", SearchSendMessageList)

		// 黑名单管理
		authenticated.GET("/blacklist", GetBlacklist)
		authenticated.POST("/blacklist", CreateBlacklist)
		authenticated.DELETE("/blacklist/:id", DeleteBlacklist)

		// 邮件模板
		authenticated.GET("/email-templates", GetEmailTemplateList)
		authenticated.GET("/email-templates/:id", GetEmailTemplateDetail)
		authenticated.POST("/email-templates", CreateEmailTemplate)
		authenticated.PUT("/email-templates/:id", UpdateEmailTemplate)
		authenticated.DELETE("/email-templates/:id", DeleteEmailTemplate)

		// 短信模板
		authenticated.GET("/sms-templates", GetSMSTemplateList)
		authenticated.POST("/sms-templates", CreateSMSTemplate)
		authenticated.PUT("/sms-templates/:id", UpdateSMSTemplate)
		authenticated.DELETE("/sms-templates/:id", DeleteSMSTemplate)

		// 友情链接
		authenticated.GET("/friendly-links", GetFriendlyLinkList)
		authenticated.POST("/friendly-links", CreateFriendlyLink)
		authenticated.PUT("/friendly-links/:id", UpdateFriendlyLink)
		authenticated.DELETE("/friendly-links/:id", DeleteFriendlyLink)

		// 合同管理
		authenticated.GET("/contracts", GetContractList)
		authenticated.GET("/contracts/:id", GetContractDetail)
		authenticated.POST("/contracts", CreateContract)
		authenticated.PUT("/contracts/:id", UpdateContract)
		authenticated.DELETE("/contracts/:id", DeleteContract)
		authenticated.POST("/contracts/:id/sign", SignContract)
		authenticated.POST("/contracts/:id/cancel", CancelContract)
		authenticated.GET("/contract-templates", GetContractTemplateList)
		authenticated.POST("/contract-templates", CreateContractTemplate)
		authenticated.PUT("/contract-templates/:id", UpdateContractTemplate)
		authenticated.DELETE("/contract-templates/:id", DeleteContractTemplate)

		// 营销推送
		authenticated.GET("/marketing/pushes", GetMarketingPushList)
		authenticated.POST("/marketing/pushes", CreateMarketingPush)
		authenticated.POST("/marketing/pushes/:id/send", SendMarketingPush)
		authenticated.DELETE("/marketing/pushes/:id", DeleteMarketingPush)

		// 取消请求
		authenticated.GET("/cancel-requests", GetCancelRequestList)
		authenticated.POST("/cancel-requests/:id/approve", ApproveCancelRequest)
		authenticated.POST("/cancel-requests/:id/reject", RejectCancelRequest)

		// 销售统计
		authenticated.GET("/sales/statistics", GetSalesStatistics)
		authenticated.GET("/sales/records", GetSalesRecords)

		// 全局可配置项
		authenticated.GET("/configurable-options", GetConfigurableOptionList)
		authenticated.POST("/configurable-options", CreateConfigurableOption)
		authenticated.PUT("/configurable-options/:id", UpdateConfigurableOption)
		authenticated.DELETE("/configurable-options/:id", DeleteConfigurableOption)

		// 第三方登录
		authenticated.GET("/oauth-providers", GetOAuthProviderList)
		authenticated.POST("/oauth-providers", CreateOAuthProvider)
		authenticated.PUT("/oauth-providers/:id", UpdateOAuthProvider)
		authenticated.DELETE("/oauth-providers/:id", DeleteOAuthProvider)

		// 官网自定义字段
		authenticated.GET("/custom-template-fields", GetCustomTemplateFieldList)
		authenticated.POST("/custom-template-fields", CreateCustomTemplateField)
		authenticated.PUT("/custom-template-fields/:id", UpdateCustomTemplateField)
		authenticated.DELETE("/custom-template-fields/:id", DeleteCustomTemplateField)

		// 流量包管理
		authenticated.GET("/traffic-packages", GetTrafficPackageList)
		authenticated.POST("/traffic-packages", CreateTrafficPackage)
		authenticated.PUT("/traffic-packages/:id", UpdateTrafficPackage)
		authenticated.DELETE("/traffic-packages/:id", DeleteTrafficPackage)
		authenticated.GET("/traffic-logs", GetTrafficLogList)

		// 任务队列
		authenticated.GET("/task-queue/overview", GetTaskQueueOverview)
		authenticated.GET("/task-queue", GetTaskQueueList)
		authenticated.POST("/task-queue/:id/retry", RetryTask)
		authenticated.DELETE("/task-queue/:id", DeleteTask)

		// 二次验证配置
		authenticated.GET("/two-factor-config", GetTwoFactorConfig)
		authenticated.PUT("/two-factor-config", UpdateTwoFactorConfig)

		// 销售设置
		authenticated.GET("/sales-config", GetSalesConfig)
		authenticated.PUT("/sales-config", UpdateSalesConfig)
		authenticated.GET("/sales-groups", GetSalesGroupList)
		authenticated.POST("/sales-groups", CreateSalesGroup)
		authenticated.PUT("/sales-groups/:id", UpdateSalesGroup)
		authenticated.DELETE("/sales-groups/:id", DeleteSalesGroup)

		// 主题模板
		authenticated.GET("/themes", GetThemeList)
		authenticated.GET("/themes/active", GetActiveTheme)
		authenticated.POST("/themes", CreateTheme)
		authenticated.PUT("/themes/:id", UpdateTheme)
		authenticated.DELETE("/themes/:id", DeleteTheme)
		authenticated.POST("/themes/:id/set-default", SetDefaultTheme)

		// 工单传递规则
		authenticated.GET("/ticket-rules", GetTicketRuleList)
		authenticated.POST("/ticket-rules", CreateTicketRule)
		authenticated.PUT("/ticket-rules/:id", UpdateTicketRule)
		authenticated.DELETE("/ticket-rules/:id", DeleteTicketRule)

		// 商品订购/财务配置
		authenticated.GET("/order-config", GetOrderConfig)
		authenticated.PUT("/order-config", UpdateOrderConfig)
		authenticated.GET("/finance-config", GetFinanceConfig)
		authenticated.PUT("/finance-config", UpdateFinanceConfig)

		// 客户分组
		authenticated.GET("/customer-groups", GetCustomerGroupList)
		authenticated.POST("/customer-groups", CreateCustomerGroup)
		authenticated.PUT("/customer-groups/:id", UpdateCustomerGroup)
		authenticated.DELETE("/customer-groups/:id", DeleteCustomerGroup)

		// Redis配置
		authenticated.GET("/redis/config", GetRedisConfig)
		authenticated.GET("/redis/health", RedisHealthCheck)

		// AI配置
		authenticated.GET("/ai/config", GetAIConfig)
		authenticated.POST("/ai/test", TestAIConnection)
		authenticated.POST("/ai/generate-description", GenerateProductDescription)

		// AI工单系统（参考mianyu_ai_ticket插件）
		authenticated.GET("/ai-ticket/config", GetAITicketConfig)
		authenticated.GET("/ai-ticket/queue/stats", GetAITicketQueueStats)
		authenticated.GET("/ai-ticket/queue", GetAITicketQueueList)
		authenticated.POST("/ai-ticket/queue/process", ProcessAITicketQueue)
		authenticated.GET("/ai-ticket/knowledge", GetAITicketKnowledgeList)
		authenticated.POST("/ai-ticket/knowledge", CreateAITicketKnowledge)
		authenticated.PUT("/ai-ticket/knowledge/:id", UpdateAITicketKnowledge)
		authenticated.DELETE("/ai-ticket/knowledge/:id", DeleteAITicketKnowledge)
		authenticated.GET("/ai-ticket/rules", GetAITicketRuleList)
		authenticated.POST("/ai-ticket/rules", CreateAITicketRule)
		authenticated.PUT("/ai-ticket/rules/:id", UpdateAITicketRule)
		authenticated.DELETE("/ai-ticket/rules/:id", DeleteAITicketRule)
		authenticated.GET("/ai-ticket/logs", GetAITicketProcessLogs)
		authenticated.POST("/ai-ticket/tickets/:id/mode", SetAITicketMode)
		authenticated.POST("/ai/ticket-reply", AITicketReply)
	}
}
