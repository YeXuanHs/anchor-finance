package admin

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

// Deps holds dependencies for admin route handlers.
type Deps struct {
	DB      *gorm.DB
	Log     *logger.Logger
	JWTKey  string
	UserSvc *service.UserService
	ProdSvc *service.ProductService
	OrdSvc  *service.OrderService
	InvSvc  *service.InvoiceService
	TicSvc  *service.TicketService

	// 新增依赖
	Redis              *redis.Client          // 验证码
	JWTMgr             *auth.JWTManager       // 聚合登录
	FrontendURL        string                 // OAuth回调
	BaseURL            string                 // 聚合登录
	UploadDir          string                 // 文件上传
	WechatAppID        string                 // 微信
	WechatAppSecret    string                 // 微信
	WechatMchID        string                 // 微信
	WechatMchKey       string                 // 微信
	WechatNotifyURL    string                 // 微信
	WechatTemplateID   string                 // 微信
	PaymentTxStore     service.TransactionStore   // 支付
	PaymentInvStore    service.InvoiceStore        // 支付
}

func RegisterRoutes(r *gin.RouterGroup, deps Deps) {
	authHandler := handler.NewAuthHandler(deps.UserSvc, deps.Log, deps.JWTKey)
	userHandler := handler.NewUserHandler(deps.UserSvc, deps.Log)
	productHandler := handler.NewProductHandler(deps.ProdSvc, deps.Log)
	orderHandler := handler.NewOrderHandler(deps.OrdSvc, deps.Log)
	invoiceHandler := handler.NewInvoiceHandler(deps.InvSvc, deps.Log)
	ticketHandler := handler.NewTicketHandler(deps.TicSvc, deps.Log)
	adminHandler := handler.NewAdminHandler(deps.DB, deps.UserSvc, deps.OrdSvc, deps.InvSvc, deps.Log)

	// 管理员登录
	r.POST("/login", authHandler.Login)

	// 需要管理员权限
	admin := r.Group("")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		// 仪表盘
		admin.GET("/dashboard", adminHandler.Dashboard)
		admin.GET("/dashboard/stats", adminHandler.Stats)

		// 用户管理
		admin.GET("/users", userHandler.GetUserList)
		admin.GET("/users/:id", userHandler.GetUserDetail)
		admin.PUT("/users/:id", userHandler.UpdateUser)
		admin.POST("/users/:id/status", userHandler.UpdateUserStatus)
		admin.DELETE("/users/:id", userHandler.DeleteUser)

		// 产品管理
		admin.GET("/products", productHandler.GetList)
		admin.POST("/products", productHandler.Create)
		admin.PUT("/products/:id", productHandler.Update)
		admin.DELETE("/products/:id", productHandler.Delete)
		admin.GET("/product-groups", productHandler.GetGroups)
		admin.POST("/product-groups", productHandler.CreateGroup)
		admin.PUT("/product-groups/:id", productHandler.UpdateGroup)
		admin.DELETE("/product-groups/:id", productHandler.DeleteGroup)

		// 订单管理
		admin.GET("/orders", orderHandler.GetList)
		admin.GET("/orders/:id", orderHandler.GetDetail)
		admin.POST("/orders/:id/status", orderHandler.UpdateStatus)

		// 账单管理
		admin.GET("/invoices", invoiceHandler.GetList)
		admin.GET("/invoices/:id", invoiceHandler.GetDetail)
		admin.POST("/invoices/:id/cancel", invoiceHandler.Cancel)

		// 工单管理
		admin.GET("/tickets", ticketHandler.GetList)
		admin.GET("/tickets/:id", ticketHandler.GetDetail)
		admin.POST("/tickets/:id/reply", ticketHandler.Reply)
		admin.POST("/tickets/:id/assign", ticketHandler.Assign)
		admin.POST("/tickets/:id/close", ticketHandler.Close)
		admin.POST("/tickets/merge", ticketHandler.MergeTickets)
		admin.POST("/tickets/:id/transfer", ticketHandler.TransferTicket)
		admin.GET("/tickets/:id/transfer-logs", ticketHandler.GetTransferLogs)
		admin.POST("/tickets/:id/attachments", ticketHandler.UploadAttachment)
		admin.GET("/tickets/:id/attachments", ticketHandler.GetAttachments)
		admin.DELETE("/tickets/attachments/:id", ticketHandler.DeleteAttachment)

		// 公告管理
		announceSvc := service.NewAnnounceService(deps.DB, deps.Log)
		announceHandler := handler.NewAnnounceHandler(announceSvc, deps.Log)
		admin.GET("/announcements", announceHandler.AdminGetList)
		admin.GET("/announcements/:id", announceHandler.AdminGetDetail)
		admin.POST("/announcements", announceHandler.AdminCreate)
		admin.PUT("/announcements/:id", announceHandler.AdminUpdate)
		admin.DELETE("/announcements/:id", announceHandler.AdminDelete)

		// 友情链接
		friendlyLinkSvc := service.NewFriendlyLinkService(deps.DB, deps.Log)
		friendlyLinkHandler := handler.NewFriendlyLinkHandler(friendlyLinkSvc, deps.Log)
		admin.GET("/friendly-links", friendlyLinkHandler.AdminGetList)
		admin.POST("/friendly-links", friendlyLinkHandler.AdminCreate)
		admin.PUT("/friendly-links/:id", friendlyLinkHandler.AdminUpdate)
		admin.DELETE("/friendly-links/:id", friendlyLinkHandler.AdminDelete)

		// 系统设置
		admin.GET("/settings", adminHandler.GetSettings)
		admin.PUT("/settings", adminHandler.UpdateSettings)
		admin.GET("/settings/:group", adminHandler.GetSettingsByGroup)

		// 操作日志
		admin.GET("/logs", adminHandler.GetLogs)

		// 轮播图管理
		bannerSvc := service.NewBannerService(deps.DB, deps.Log)
		bannerHandler := handler.NewBannerHandler(bannerSvc, deps.Log)
		admin.GET("/banners", bannerHandler.List)
		admin.GET("/banners/:id", bannerHandler.GetDetail)
		admin.POST("/banners", bannerHandler.Create)
		admin.PUT("/banners/:id", bannerHandler.Update)
		admin.DELETE("/banners/:id", bannerHandler.Delete)
		admin.POST("/banners/:id/status", bannerHandler.SetStatus)

		// 支付方式管理
		paymentGwSvc := service.NewPaymentGatewayService(deps.DB, deps.Log)
		paymentGwHandler := handler.NewPaymentGatewayHandler(paymentGwSvc, deps.Log)
		admin.GET("/payment-gateways", paymentGwHandler.List)
		admin.GET("/payment-gateways/:id", paymentGwHandler.GetDetail)
		admin.POST("/payment-gateways", paymentGwHandler.Create)
		admin.PUT("/payment-gateways/:id", paymentGwHandler.Update)
		admin.DELETE("/payment-gateways/:id", paymentGwHandler.Delete)
		admin.POST("/payment-gateways/:id/status", paymentGwHandler.SetStatus)
		admin.POST("/payment-gateways/:id/test", paymentGwHandler.TestConnection)

		// OAuth提供商管理
		oauthProvSvc := service.NewOAuthProviderService(deps.DB, deps.Log)
		oauthProvHandler := handler.NewOAuthProviderHandler(oauthProvSvc, deps.Log)
		admin.GET("/oauth-providers", oauthProvHandler.List)
		admin.GET("/oauth-providers/:id", oauthProvHandler.GetDetail)
		admin.POST("/oauth-providers", oauthProvHandler.Create)
		admin.PUT("/oauth-providers/:id", oauthProvHandler.Update)
		admin.DELETE("/oauth-providers/:id", oauthProvHandler.Delete)
		admin.POST("/oauth-providers/:id/status", oauthProvHandler.SetStatus)

		// Cron定时任务
		provisionSvcForCron := service.NewProvisionService(deps.DB, deps.Log)
		cronSvc := service.NewCronService(deps.DB, deps.Log, provisionSvcForCron)
		cronHandler := handler.NewCronHandler(cronSvc, deps.Log)
		admin.GET("/cron-tasks", cronHandler.GetList)
		admin.GET("/cron-tasks/:id", cronHandler.GetDetail)
		admin.POST("/cron-tasks", cronHandler.Create)
		admin.PUT("/cron-tasks/:id", cronHandler.Update)
		admin.DELETE("/cron-tasks/:id", cronHandler.Delete)
		admin.POST("/cron-tasks/:id/status", cronHandler.SetStatus)
		admin.POST("/cron-tasks/:id/run", cronHandler.RunTask)
		admin.GET("/cron-tasks/:id/logs", cronHandler.GetLogs)

		// DCIM管理
		dcimSvc := service.NewDcimService(deps.DB, deps.Log)
		dcimHandler := handler.NewDcimHandler(dcimSvc, deps.Log)
		admin.GET("/dcim/servers", dcimHandler.GetServerList)
		admin.GET("/dcim/servers/:id", dcimHandler.GetServerDetail)
		admin.POST("/dcim/servers", dcimHandler.CreateServer)
		admin.PUT("/dcim/servers/:id", dcimHandler.UpdateServer)
		admin.DELETE("/dcim/servers/:id", dcimHandler.DeleteServer)
		admin.POST("/dcim/servers/:id/boot", dcimHandler.BootServer)
		admin.POST("/dcim/servers/:id/shutdown", dcimHandler.ShutdownServer)
		admin.POST("/dcim/servers/:id/reboot", dcimHandler.RebootServer)
		admin.GET("/dcim/datacenters", dcimHandler.GetDatacenterList)

		// 应用商店
		appstoreSvc := service.NewAppStoreService(deps.DB, deps.Log)
		appstoreHandler := handler.NewAppStoreHandler(appstoreSvc, deps.Log)
		admin.GET("/apps", appstoreHandler.AdminList)
		admin.POST("/apps", appstoreHandler.AdminCreate)
		admin.PUT("/apps/:id", appstoreHandler.AdminUpdate)
		admin.DELETE("/apps/:id", appstoreHandler.AdminDelete)

		// 插件管理
		pluginSvc := service.NewPluginService(deps.DB, deps.Log)
		pluginHandler := handler.NewPluginHandler(pluginSvc, deps.Log)
		admin.GET("/plugins", pluginHandler.List)
		admin.POST("/plugins", pluginHandler.Install)
		admin.DELETE("/plugins/:id", pluginHandler.Uninstall)
		admin.POST("/plugins/:id/enable", pluginHandler.Enable)
		admin.POST("/plugins/:id/disable", pluginHandler.Disable)
		admin.PUT("/plugins/:id/config", pluginHandler.UpdateConfig)

		// 模块供给
		provisionSvc := service.NewProvisionService(deps.DB, deps.Log)
		provisionHandler := handler.NewProvisionHandler(provisionSvc, deps.Log)
		admin.GET("/provisions", provisionHandler.List)
		admin.GET("/provisions/:id", provisionHandler.GetDetail)
		admin.POST("/provisions", provisionHandler.Create)
		admin.PUT("/provisions/:id", provisionHandler.Update)
		admin.DELETE("/provisions/:id", provisionHandler.Delete)
		admin.POST("/provisions/:id/test", provisionHandler.TestConnection)

		// 服务开通操作
		admin.POST("/provision/execute", provisionHandler.Execute)
		admin.POST("/provision/:id/suspend", provisionHandler.Suspend)
		admin.POST("/provision/:id/terminate", provisionHandler.Terminate)
		admin.POST("/provision/:id/unsuspend", provisionHandler.Unsuspend)

		// 客户分组
		clientGroupSvc := service.NewClientGroupService(deps.DB)
		clientGroupHandler := handler.NewClientGroupHandler(clientGroupSvc, deps.Log)
		admin.GET("/client-groups", clientGroupHandler.List)
		admin.GET("/client-groups/:id", clientGroupHandler.Get)
		admin.POST("/client-groups", clientGroupHandler.Create)
		admin.PUT("/client-groups/:id", clientGroupHandler.Update)
		admin.DELETE("/client-groups/:id", clientGroupHandler.Delete)

		// 客户服务
		clientSvcSvc := service.NewClientServiceService(deps.DB)
		clientSvcHandler := handler.NewClientServiceHandler(clientSvcSvc, deps.Log)
		admin.GET("/client-services", clientSvcHandler.List)
		admin.GET("/client-services/:id", clientSvcHandler.Get)
		admin.POST("/client-services/:id/suspend", clientSvcHandler.Suspend)
		admin.POST("/client-services/:id/terminate", clientSvcHandler.Terminate)
		admin.POST("/client-services/:id/renew", clientSvcHandler.Renew)

		// 用户管理
		userManageSvc := service.NewUserManageService(deps.DB, deps.Log)
		userManageHandler := handler.NewUserManageHandler(userManageSvc, deps.Log)
		admin.GET("/user-manage/search", userManageHandler.Search)
		admin.POST("/user-manage/:id/ban", userManageHandler.Ban)
		admin.POST("/user-manage/:id/unban", userManageHandler.Unban)
		admin.POST("/user-manage/:id/balance", userManageHandler.AdjustBalance)
		admin.POST("/user-manage/:id/reset-password", userManageHandler.ResetPassword)

		// 用户备注
		userRemarkSvc := service.NewUserRemarkService(deps.DB, deps.Log)
		userRemarkHandler := handler.NewUserRemarkHandler(userRemarkSvc, deps.Log)
		admin.GET("/user-remarks", userRemarkHandler.List)
		admin.POST("/user-remarks", userRemarkHandler.Add)
		admin.DELETE("/user-remarks/:id", userRemarkHandler.AdminDelete)

		// 邮件模板
		emailTplSvc := service.NewEmailTemplateService(deps.DB, deps.Log)
		emailTplHandler := handler.NewEmailTemplateHandler(emailTplSvc, deps.Log)
		admin.GET("/email-templates", emailTplHandler.List)
		admin.GET("/email-templates/:id", emailTplHandler.GetDetail)
		admin.POST("/email-templates", emailTplHandler.Create)
		admin.PUT("/email-templates/:id", emailTplHandler.Update)
		admin.DELETE("/email-templates/:id", emailTplHandler.Delete)
		admin.POST("/email-templates/:id/test", emailTplHandler.SendTest)

		// 菜单管理
		menuSvc := service.NewMenuService(deps.DB, deps.Log)
		menuHandler := handler.NewMenuHandler(menuSvc, deps.Log)
		admin.GET("/menus", menuHandler.List)
		admin.GET("/menus/tree", menuHandler.GetTree)
		admin.POST("/menus", menuHandler.Create)
		admin.PUT("/menus/:id", menuHandler.Update)
		admin.DELETE("/menus/:id", menuHandler.Delete)
		admin.POST("/menus/sort", menuHandler.Sort)

		// 日志记录
		logRecordSvc := service.NewLogRecordService(deps.DB, deps.Log)
		logHandler := handler.NewLogRecordHandler(logRecordSvc, deps.Log)
		admin.GET("/log-records", logHandler.List)
		admin.GET("/log-records/:id", logHandler.GetDetail)
		admin.GET("/log-records/stats", logHandler.Stats)
		admin.POST("/log-records/export", logHandler.Export)

		// 规则管理
		ruleSvc := service.NewRuleService(deps.DB, deps.Log)
		ruleHandler := handler.NewRuleHandler(ruleSvc, deps.Log)
		admin.GET("/rules", ruleHandler.List)
		admin.GET("/rules/:id", ruleHandler.GetDetail)
		admin.POST("/rules", ruleHandler.Create)
		admin.PUT("/rules/:id", ruleHandler.Update)
		admin.DELETE("/rules/:id", ruleHandler.Delete)
		admin.POST("/rules/:id/enable", ruleHandler.Enable)
		admin.POST("/rules/:id/disable", ruleHandler.Disable)
		admin.POST("/rules/:id/test", ruleHandler.Test)

		// 促销管理
		saleSvc := service.NewSaleService(deps.DB, deps.Log)
		saleHandler := handler.NewSaleHandler(saleSvc, deps.Log)
		admin.GET("/sales", saleHandler.List)
		admin.GET("/sales/:id", saleHandler.GetDetail)
		admin.POST("/sales", saleHandler.Create)
		admin.PUT("/sales/:id", saleHandler.Update)
		admin.DELETE("/sales/:id", saleHandler.Delete)
		admin.POST("/sales/:id/status", saleHandler.SetStatus)
		admin.GET("/sales/:id/commission-ladder", saleHandler.GetCommissionLadder)
		admin.PUT("/sales/:id/commission-ladder", saleHandler.SetCommissionLadder)
		admin.POST("/sales/:id/calculate-commission", saleHandler.CalculateCommission)
		admin.POST("/sales/:id/validate-usage", saleHandler.ValidateUsage)

		// 消息发送
		sendMsgSvc := service.NewSendMessageService(deps.DB, deps.Log)
		sendMsgHandler := handler.NewSendMessageHandler(sendMsgSvc, deps.Log)
		admin.POST("/messages/send", sendMsgHandler.SendEmail)
		admin.POST("/messages/send-sms", sendMsgHandler.SendSMS)
		admin.POST("/messages/send-site", sendMsgHandler.SendSiteMessage)
		admin.POST("/messages/batch", sendMsgHandler.BatchSend)
		admin.GET("/messages", sendMsgHandler.GetList)
		admin.GET("/messages/:id", sendMsgHandler.GetDetail)
		admin.GET("/messages/batch/:batch_id", sendMsgHandler.GetByBatchID)
		admin.POST("/messages/retry-failed", sendMsgHandler.RetryFailed)

		// 工单部门
		ticketDeptSvc := service.NewTicketDepartmentService(deps.DB, deps.Log)
		ticketDeptHandler := handler.NewTicketDeptHandler(ticketDeptSvc, deps.Log)
		admin.GET("/ticket-depts", ticketDeptHandler.List)
		admin.GET("/ticket-depts/tree", ticketDeptHandler.GetTree)
		admin.POST("/ticket-depts", ticketDeptHandler.Create)
		admin.GET("/ticket-depts/:id", ticketDeptHandler.GetDetail)
		admin.PUT("/ticket-depts/:id", ticketDeptHandler.Update)
		admin.DELETE("/ticket-depts/:id", ticketDeptHandler.Delete)
		admin.POST("/ticket-depts/:id/enable", ticketDeptHandler.Enable)
		admin.POST("/ticket-depts/:id/disable", ticketDeptHandler.Disable)
		admin.GET("/ticket-depts/:id/members", ticketDeptHandler.GetMembers)
		admin.POST("/ticket-depts/:id/members", ticketDeptHandler.AddMember)
		admin.DELETE("/ticket-depts/:id/members/:user_id", ticketDeptHandler.RemoveMember)
		admin.POST("/ticket-depts/:id/managers", ticketDeptHandler.SetManagers)

		// 主机管理
		upstreamSvc := service.NewUpstreamService(deps.DB, deps.Log)
		hostSvc := service.NewHostService(deps.DB, deps.Log, upstreamSvc)
		hostHandler := handler.NewHostHandler(hostSvc, deps.Log)
		admin.GET("/hosts", hostHandler.List)
		admin.GET("/hosts/:id", hostHandler.GetDetail)
		admin.POST("/hosts/:id/boot", hostHandler.Boot)
		admin.POST("/hosts/:id/shutdown", hostHandler.Shutdown)
		admin.POST("/hosts/:id/reboot", hostHandler.Reboot)

		// 系统配置
		configServerSvc := service.NewConfigServerService(deps.DB, deps.Log)
		configServerHandler := handler.NewConfigServerHandler(configServerSvc, deps.Log)
		admin.GET("/config/servers", configServerHandler.List)
		admin.POST("/config/servers", configServerHandler.Create)
		admin.PUT("/config/servers/:id", configServerHandler.Update)
		admin.DELETE("/config/servers/:id", configServerHandler.Delete)

		configOptionSvc := service.NewConfigOptionService(deps.DB, deps.Log)
		configOptionHandler := handler.NewConfigOptionHandler(configOptionSvc, deps.Log)
		admin.GET("/config/options", configOptionHandler.List)
		admin.POST("/config/options", configOptionHandler.Create)
		admin.PUT("/config/options/:id", configOptionHandler.Update)
		admin.DELETE("/config/options/:id", configOptionHandler.Delete)
		admin.GET("/product-config-groups", configOptionHandler.GetProductConfigGroups)
		admin.POST("/product-config-groups", configOptionHandler.CreateProductConfigGroup)
		admin.PUT("/product-config-groups/:id", configOptionHandler.UpdateProductConfigGroup)
		admin.DELETE("/product-config-groups/:id", configOptionHandler.DeleteProductConfigGroup)
		admin.POST("/product-config-groups/:id/link", configOptionHandler.LinkGroupToProduct)
		admin.DELETE("/product-config-groups/:id/link/:product_id", configOptionHandler.UnlinkGroupFromProduct)

		configGeneralSvc := service.NewConfigGeneralService(deps.DB, deps.Log)
		configGeneralHandler := handler.NewConfigGeneralHandler(configGeneralSvc, deps.Log)
		admin.GET("/config/general", configGeneralHandler.Get)
		admin.PUT("/config/general", configGeneralHandler.Update)

		configCertifiSvc := service.NewConfigCertifiService(deps.DB, deps.Log)
		configCertifiHandler := handler.NewConfigCertifiHandler(configCertifiSvc, deps.Log)
		admin.GET("/config/certifi", configCertifiHandler.Get)
		admin.PUT("/config/certifi", configCertifiHandler.Update)

		// ==================== 新增模块 ====================

		// 高级配置选项
		advancedOptSvc := service.NewAdvancedOptionsService(deps.DB, deps.Log)
		advancedOptHandler := handler.NewAdvancedOptionsHandler(advancedOptSvc, deps.Log)
		admin.GET("/advanced-options", advancedOptHandler.GetOptions)
		admin.POST("/advanced-options", advancedOptHandler.CreateOption)
		admin.PUT("/advanced-options/:id", advancedOptHandler.UpdateOption)
		admin.DELETE("/advanced-options/:id", advancedOptHandler.DeleteOption)
		admin.GET("/advanced-options/links", advancedOptHandler.GetLinks)
		admin.POST("/advanced-options/links", advancedOptHandler.CreateLink)
		admin.PUT("/advanced-options/links/:id", advancedOptHandler.UpdateLink)
		admin.DELETE("/advanced-options/links/:id", advancedOptHandler.DeleteLink)

		// 批量同步
		batchSyncSvc := service.NewBatchSyncService(deps.DB, deps.Log)
		batchSyncHandler := handler.NewBatchSyncHandler(batchSyncSvc, deps.Log)
		admin.GET("/batch-sync", batchSyncHandler.GetList)
		admin.GET("/batch-sync/:id", batchSyncHandler.GetDetail)
		admin.POST("/batch-sync/:id/execute", batchSyncHandler.Execute)
		admin.GET("/batch-sync/:id/logs", batchSyncHandler.GetLogs)

		// 社区管理
		communitySvc := service.NewCommunityService(deps.DB, deps.Log)
		communityHandler := handler.NewCommunityHandler(communitySvc, deps.Log)
		admin.GET("/community/posts", communityHandler.GetPostList)
		admin.GET("/community/posts/:id", communityHandler.GetPost)
		admin.POST("/community/posts", communityHandler.CreatePost)
		admin.PUT("/community/posts/:id", communityHandler.UpdatePost)
		admin.DELETE("/community/posts/:id", communityHandler.DeletePost)
		admin.GET("/community/posts/:id/comments", communityHandler.GetComments)
		admin.POST("/community/posts/:id/comments", communityHandler.CreateComment)
		admin.DELETE("/community/comments/:id", communityHandler.DeleteComment)

		// URL定时任务
		cronURLSvc := service.NewCronURLService(deps.DB, deps.Log)
		cronURLHandler := handler.NewCronURLHandler(cronURLSvc, deps.Log)
		admin.GET("/cron-url", cronURLHandler.GetList)
		admin.GET("/cron-url/:id", cronURLHandler.GetDetail)
		admin.POST("/cron-url", cronURLHandler.Create)
		admin.PUT("/cron-url/:id", cronURLHandler.Update)
		admin.DELETE("/cron-url/:id", cronURLHandler.Delete)
		admin.POST("/cron-url/:id/status", cronURLHandler.SetStatus)
		admin.POST("/cron-url/:id/run", cronURLHandler.RunCron)
		admin.GET("/cron-url/logs", cronURLHandler.GetLogs)

		// 魔方云对接
		dcimCloudSvc := service.NewDcimCloudService(deps.DB, deps.Log)
		dcimCloudHandler := handler.NewDcimCloudHandler(dcimCloudSvc, deps.Log)
		admin.GET("/dcim-cloud/servers", dcimCloudHandler.GetServers)
		admin.GET("/dcim-cloud/servers/:id", dcimCloudHandler.GetServer)
		admin.POST("/dcim-cloud/servers", dcimCloudHandler.CreateServer)
		admin.PUT("/dcim-cloud/servers/:id", dcimCloudHandler.UpdateServer)
		admin.DELETE("/dcim-cloud/servers/:id", dcimCloudHandler.DeleteServer)
		admin.POST("/dcim-cloud/servers/:id/test", dcimCloudHandler.TestConnection)
		admin.POST("/dcim-cloud/servers/:id/sync", dcimCloudHandler.SyncServer)
		admin.GET("/dcim-cloud/logs", dcimCloudHandler.GetOperationLogs)

		// 钩子系统
		hookSvc := service.NewHookService(deps.DB, deps.Log)
		hookHandler := handler.NewHookHandler(hookSvc, deps.Log)
		admin.GET("/hooks", hookHandler.GetList)
		admin.GET("/hooks/:id", hookHandler.GetDetail)
		admin.POST("/hooks", hookHandler.Create)
		admin.PUT("/hooks/:id", hookHandler.Update)
		admin.DELETE("/hooks/:id", hookHandler.Delete)
		admin.POST("/hooks/:id/status", hookHandler.SetStatus)
		admin.POST("/hooks/:id/trigger", hookHandler.Trigger)
		admin.GET("/hooks/logs", hookHandler.GetLogs)

		// 关联原因
		linkCauseSvc := service.NewLinkCauseService(deps.DB, deps.Log)
		linkCauseHandler := handler.NewLinkCauseHandler(linkCauseSvc, deps.Log)
		admin.GET("/link-causes", linkCauseHandler.GetCauses)
		admin.GET("/link-causes/tree", linkCauseHandler.GetTree)
		admin.GET("/link-causes/:id", linkCauseHandler.GetCause)
		admin.POST("/link-causes", linkCauseHandler.CreateCause)
		admin.PUT("/link-causes/:id", linkCauseHandler.UpdateCause)
		admin.DELETE("/link-causes/:id", linkCauseHandler.DeleteCause)

		// 关联知识
		linkKnowledgeSvc := service.NewLinkKnowledgeService(deps.DB, deps.Log)
		linkKnowledgeHandler := handler.NewLinkKnowledgeHandler(linkKnowledgeSvc, deps.Log)
		admin.GET("/link-knowledge", linkKnowledgeHandler.GetKnowledges)
		admin.GET("/link-knowledge/:id", linkKnowledgeHandler.GetKnowledge)
		admin.POST("/link-knowledge", linkKnowledgeHandler.CreateKnowledge)
		admin.PUT("/link-knowledge/:id", linkKnowledgeHandler.UpdateKnowledge)
		admin.DELETE("/link-knowledge/:id", linkKnowledgeHandler.DeleteKnowledge)

		// 公共接口
		publicSvc := service.NewPublicService(deps.DB, deps.Log)
		publicHandler := handler.NewPublicHandler(publicSvc, deps.Log)
		admin.GET("/public/system-info", publicHandler.GetSystemInfo)
		admin.GET("/public/config/:key", publicHandler.GetConfig)
		admin.GET("/public/configs", publicHandler.GetConfigs)

		// RBAC页面权限
		rbacPageSvc := service.NewRbacPageService(deps.DB, deps.Log)
		rbacPageHandler := handler.NewRbacPageHandler(rbacPageSvc, deps.Log)
		admin.GET("/rbac-pages", rbacPageHandler.GetPages)
		admin.GET("/rbac-pages/tree", rbacPageHandler.GetTree)
		admin.GET("/rbac-pages/:id", rbacPageHandler.GetPage)
		admin.POST("/rbac-pages", rbacPageHandler.CreatePage)
		admin.PUT("/rbac-pages/:id", rbacPageHandler.UpdatePage)
		admin.DELETE("/rbac-pages/:id", rbacPageHandler.DeletePage)
		admin.GET("/rbac-pages/auth-tree", rbacPageHandler.GetAuthTree)

		// 规则中间件
		ruleMiddleSvc := service.NewRuleMiddleService(deps.DB, deps.Log)
		ruleMiddleHandler := handler.NewRuleMiddleHandler(ruleMiddleSvc, deps.Log)
		admin.GET("/rule-middle/menus", ruleMiddleHandler.GetMenuList)
		admin.POST("/rule-middle/menus", ruleMiddleHandler.AddMenu)
		admin.PUT("/rule-middle/menus/:id", ruleMiddleHandler.UpdateMenu)
		admin.DELETE("/rule-middle/menus/:id", ruleMiddleHandler.DeleteMenu)

		// 批量发送消息
		sendMsgBatchSvc := service.NewSendMessageBatchService(deps.DB, deps.Log)
		sendMsgBatchHandler := handler.NewSendMessageBatchHandler(sendMsgBatchSvc, deps.Log)
		admin.GET("/messages/batch/search-params", sendMsgBatchHandler.GetSearchParams)
		admin.GET("/messages/batch/list", sendMsgBatchHandler.GetBatches)
		admin.POST("/messages/batch/send", sendMsgBatchHandler.SendBatch)
		admin.GET("/messages/batch/progress/:id", sendMsgBatchHandler.GetProgress)
		admin.GET("/messages/batch/records", sendMsgBatchHandler.GetRecords)

		// 系统管理
		systemSvc := service.NewSystemService(deps.DB, deps.Log, deps.Redis)
		systemHandler := handler.NewSystemHandler(systemSvc, deps.Log)
		admin.GET("/system/common-info", systemHandler.GetCommonInfo)
		admin.GET("/system/update-content", systemHandler.GetUpdateContent)
		admin.GET("/system/check-update", systemHandler.CheckUpdate)
		admin.GET("/system/updates", systemHandler.GetUpdateList)
		admin.POST("/system/updates/:id/install", systemHandler.InstallUpdate)
		admin.GET("/system/logs", systemHandler.GetSystemLog)
		admin.POST("/system/clear-cache", systemHandler.ClearCache)

		// 工单传递
		ticketDeliverSvc := service.NewTicketDeliverService(deps.DB, deps.Log)
		ticketDeliverHandler := handler.NewTicketDeliverHandler(ticketDeliverSvc, deps.Log)
		admin.GET("/ticket-deliver/add-page", ticketDeliverHandler.GetAddPage)
		admin.GET("/ticket-deliver/rules", ticketDeliverHandler.GetRules)
		admin.GET("/ticket-deliver/rules/:id", ticketDeliverHandler.GetRule)
		admin.POST("/ticket-deliver/rules", ticketDeliverHandler.CreateRule)
		admin.PUT("/ticket-deliver/rules/:id", ticketDeliverHandler.UpdateRule)
		admin.DELETE("/ticket-deliver/rules/:id", ticketDeliverHandler.DeleteRule)
		admin.GET("/ticket-deliver/logs", ticketDeliverHandler.GetLogs)

		// 工单预设回复
		ticketPrereplySvc := service.NewTicketPrereplyService(deps.DB, deps.Log)
		ticketPrereplyHandler := handler.NewTicketPrereplyHandler(ticketPrereplySvc, deps.Log)
		admin.GET("/ticket-prereply", ticketPrereplyHandler.GetReplyList)
		admin.POST("/ticket-prereply/categories", ticketPrereplyHandler.AddCategory)
		admin.PUT("/ticket-prereply/categories/:id", ticketPrereplyHandler.UpdateCategory)
		admin.DELETE("/ticket-prereply/categories/:id", ticketPrereplyHandler.DeleteCategory)
		admin.POST("/ticket-prereply/replies", ticketPrereplyHandler.AddReply)
		admin.PUT("/ticket-prereply/replies/:id", ticketPrereplyHandler.UpdateReply)
		admin.DELETE("/ticket-prereply/replies/:id", ticketPrereplyHandler.DeleteReply)

		// 工单状态
		ticketStatusSvc := service.NewTicketStatusService(deps.DB, deps.Log)
		ticketStatusHandler := handler.NewTicketStatusHandler(ticketStatusSvc, deps.Log)
		admin.GET("/ticket-statuses", ticketStatusHandler.GetStatuses)
		admin.GET("/ticket-statuses/default", ticketStatusHandler.GetDefaultStatuses)
		admin.POST("/ticket-statuses", ticketStatusHandler.AddStatus)
		admin.PUT("/ticket-statuses/:id", ticketStatusHandler.UpdateStatus)
		admin.DELETE("/ticket-statuses/:id", ticketStatusHandler.DeleteStatus)

		// 用户偏好
		userTasteSvc := service.NewUserTasteService(deps.DB, deps.Log)
		userTasteHandler := handler.NewUserTastesHandler(userTasteSvc, deps.Log)
		admin.GET("/user-tastes", userTasteHandler.GetUserTastes)
		admin.PUT("/user-tastes", userTasteHandler.UpdateUserTastes)

		// ==================== 补充模块 ====================

		// 登录日志
		loginLogSvc := service.NewLoginLogService(deps.DB, deps.Log)
		loginLogHandler := handler.NewLoginLogHandler(loginLogSvc, deps.Log)
		admin.GET("/login-logs", loginLogHandler.List)
		admin.GET("/login-logs/:id", loginLogHandler.GetDetail)
		admin.DELETE("/login-logs/:id", loginLogHandler.Delete)
		admin.POST("/login-logs/cleanup", loginLogHandler.Cleanup)
		admin.GET("/login-logs/stats", loginLogHandler.GetStats)
		admin.GET("/login-logs/failed-attempts", loginLogHandler.GetFailedAttempts)
		admin.POST("/login-logs/export", loginLogHandler.Export)

		// 系统日志
		systemLogSvc := service.NewSystemLogService(deps.DB, deps.Log)
		systemLogHandler := handler.NewSystemLogHandler(systemLogSvc, deps.Log)
		admin.GET("/system-logs", systemLogHandler.List)
		admin.GET("/system-logs/:id", systemLogHandler.GetDetail)
		admin.DELETE("/system-logs/:id", systemLogHandler.Delete)
		admin.POST("/system-logs/cleanup", systemLogHandler.Cleanup)
		admin.GET("/system-logs/stats", systemLogHandler.GetStats)
		admin.GET("/system-logs/level-stats", systemLogHandler.GetLevelStats)
		admin.GET("/system-logs/module-stats", systemLogHandler.GetModuleStats)
		admin.POST("/system-logs/export", systemLogHandler.Export)
		admin.POST("/system-logs/clear-by-level", systemLogHandler.ClearByLevel)

		// API管理
		apiManageSvc := service.NewAPIManageService(deps.DB, deps.Log)
		apiManageHandler := handler.NewAPIManageHandler(apiManageSvc, deps.Log)
		admin.GET("/api-keys", apiManageHandler.List)
		admin.GET("/api-keys/:id", apiManageHandler.GetDetail)
		admin.POST("/api-keys", apiManageHandler.Create)
		admin.PUT("/api-keys/:id", apiManageHandler.Update)
		admin.DELETE("/api-keys/:id", apiManageHandler.Delete)
		admin.POST("/api-keys/:id/status", apiManageHandler.SetStatus)
		admin.POST("/api-keys/:id/regenerate", apiManageHandler.Regenerate)
		admin.GET("/api-keys/permissions", apiManageHandler.GetPermissions)

		// API调用日志
		apiLogSvc := service.NewApiLogService(deps.DB, deps.Log)
		apiLogHandler := handler.NewAPILogHandler(apiLogSvc, deps.Log)
		admin.GET("/api-logs", apiLogHandler.List)
		admin.GET("/api-logs/:id", apiLogHandler.GetDetail)
		admin.DELETE("/api-logs/:id", apiLogHandler.Delete)
		admin.POST("/api-logs/cleanup", apiLogHandler.Cleanup)
		admin.GET("/api-logs/stats", apiLogHandler.GetStats)
		admin.GET("/api-logs/top-endpoints", apiLogHandler.GetTopEndpoints)
		admin.GET("/api-logs/slow-requests", apiLogHandler.GetSlowRequests)
		admin.GET("/api-logs/error-rate", apiLogHandler.GetErrorRate)
		admin.POST("/api-logs/export", apiLogHandler.Export)

		// 批量续费
		balSvc := service.NewBalanceService(deps.DB, deps.Log)
		multiRenewSvc := service.NewMultiRenewService(deps.DB, deps.Log, deps.InvSvc, balSvc)
		multiRenewHandler := handler.NewMultiRenewHandler(multiRenewSvc, deps.Log)
		admin.GET("/multi-renew", multiRenewHandler.List)
		admin.GET("/multi-renew/:id", multiRenewHandler.GetDetail)
		admin.POST("/multi-renew", multiRenewHandler.Create)
		admin.POST("/multi-renew/:id/execute", multiRenewHandler.Execute)
		admin.POST("/multi-renew/:id/cancel", multiRenewHandler.Cancel)
		admin.DELETE("/multi-renew/:id", multiRenewHandler.Delete)
		admin.GET("/multi-renew/:id/logs", multiRenewHandler.GetLogs)
		admin.GET("/multi-renew/stats", multiRenewHandler.GetStats)
		admin.POST("/multi-renew/preview", multiRenewHandler.Preview)

		// 账号绑定
		bindSvc := service.NewBindService(deps.DB, deps.Log)
		bindHandler := handler.NewBindHandler(bindSvc, deps.Log)
		admin.GET("/binds", bindHandler.List)
		admin.GET("/binds/:id", bindHandler.GetDetail)
		admin.GET("/binds/user/:user_id", bindHandler.GetUserBindings)
		admin.DELETE("/binds/:id", bindHandler.Delete)
		admin.POST("/binds/unbind", bindHandler.Unbind)
		admin.GET("/binds/providers", bindHandler.GetProviders)
		admin.GET("/binds/stats", bindHandler.GetStats)
		admin.POST("/binds/batch-unbind", bindHandler.BatchUnbind)

		// 维护模式
		maintenanceSvc := service.NewMaintenanceService(deps.DB, deps.Log)
		maintenanceHandler := handler.NewMaintenanceHandler(maintenanceSvc, deps.Log)
		admin.GET("/maintenance/status", maintenanceHandler.GetStatus)
		admin.POST("/maintenance/enable", maintenanceHandler.Enable)
		admin.POST("/maintenance/disable", maintenanceHandler.Disable)
		admin.PUT("/maintenance/settings", maintenanceHandler.Update)
		admin.GET("/maintenance/history", maintenanceHandler.GetHistory)
		admin.POST("/maintenance/allowed-ips", maintenanceHandler.AddAllowedIP)
		admin.DELETE("/maintenance/allowed-ips", maintenanceHandler.RemoveAllowedIP)
		admin.GET("/maintenance/allowed-ips", maintenanceHandler.GetAllowedIPs)
		admin.POST("/maintenance/test", maintenanceHandler.TestMode)

		// 服务详情
		serviceDetailSvc := service.NewServiceDetailService(deps.DB, deps.Log)
		serviceDetailHandler := handler.NewServiceDetailHandler(serviceDetailSvc, deps.Log)
		admin.GET("/service-details", serviceDetailHandler.List)
		admin.GET("/service-details/:id", serviceDetailHandler.GetDetail)
		admin.GET("/service-details/user/:user_id", serviceDetailHandler.GetByUser)
		admin.POST("/service-details", serviceDetailHandler.Create)
		admin.PUT("/service-details/:id", serviceDetailHandler.Update)
		admin.DELETE("/service-details/:id", serviceDetailHandler.Delete)
		admin.POST("/service-details/:id/suspend", serviceDetailHandler.Suspend)
		admin.POST("/service-details/:id/unsuspend", serviceDetailHandler.Unsuspend)
		admin.POST("/service-details/:id/terminate", serviceDetailHandler.Terminate)
		admin.POST("/service-details/:id/renew", serviceDetailHandler.Renew)
		admin.GET("/service-details/stats", serviceDetailHandler.GetStats)
		admin.GET("/service-details/:id/logs", serviceDetailHandler.GetServiceLogs)

		// ==================== 推广联盟 ====================
		affiliateSvc := service.NewAffiliateService(deps.DB, deps.Log)
		affiliateHandler := handler.NewAffiliateHandler(affiliateSvc, deps.Log)
		admin.GET("/affiliate", affiliateHandler.AdminGetList)
		admin.POST("/affiliate/records/:id/confirm", affiliateHandler.AdminConfirmRecord)
		admin.POST("/affiliate/withdraws/:id/process", affiliateHandler.AdminProcessWithdraw)

		// ==================== 代理商 ====================
		agentHandler := handler.NewAgentHandler(deps.DB)
		admin.GET("/agents", agentHandler.AdminGetList)
		admin.POST("/agents", agentHandler.AdminCreate)
		admin.PUT("/agents/:id", agentHandler.AdminUpdate)
		admin.POST("/agents/commissions/:id/confirm", agentHandler.AdminConfirmCommission)

		// ==================== 聚合登录 ====================
		if deps.JWTMgr != nil {
			aggregateLoginSvc := service.NewAggregateLoginService(deps.DB, deps.Log, deps.UserSvc, deps.JWTMgr, deps.BaseURL)
			oauthSvcForAgg := service.NewOAuthService(deps.DB, deps.Log, deps.UserSvc, deps.FrontendURL)
			aggregateLoginHandler := handler.NewAggregateLoginHandler(aggregateLoginSvc, oauthSvcForAgg, deps.UserSvc, deps.Log, deps.JWTKey)
			admin.GET("/aggregate-login/providers", aggregateLoginHandler.GetProviders)
			admin.GET("/aggregate-login/:code", aggregateLoginHandler.Login)
			admin.GET("/aggregate-login/:code/callback", aggregateLoginHandler.Callback)
			admin.POST("/aggregate-login/bind", aggregateLoginHandler.BindAccount)
			admin.GET("/aggregate-login/accounts", aggregateLoginHandler.GetBoundAccounts)
			admin.POST("/aggregate-login/unbind", aggregateLoginHandler.UnbindAccount)
		}

		// ==================== 余额 ====================
		balanceSvc := service.NewBalanceLogService(deps.DB)
		balanceHandler := handler.NewBalanceHandler(balanceSvc, deps.DB)
		admin.GET("/balances", balanceHandler.GetBalance)
		admin.GET("/balances/logs", balanceHandler.GetBalanceLogs)
		admin.POST("/balances/recharge", balanceHandler.Recharge)

		// ==================== 验证码 ====================
		if deps.Redis != nil {
			captchaSvc := service.NewCaptchaService(deps.Redis)
			captchaHandler := handler.NewCaptchaHandler(captchaSvc, deps.DB)
			admin.GET("/captcha/image", captchaHandler.GetImage)
			admin.POST("/captcha/sms", captchaHandler.SendSMS)
			admin.POST("/captcha/email", captchaHandler.SendEmail)
		}

		// ==================== 购物车 ====================
		couponSvcForCart := service.NewCouponService(deps.DB, deps.Log)
		cartSvc := service.NewCartService(deps.DB, deps.Log, deps.OrdSvc, couponSvcForCart)
		cartHandler := handler.NewCartHandler(cartSvc)
		admin.GET("/cart", cartHandler.GetCart)
		admin.POST("/cart", cartHandler.AddToCart)
		admin.PUT("/cart/:id", cartHandler.UpdateCart)
		admin.DELETE("/cart/:id", cartHandler.RemoveFromCart)
		admin.DELETE("/cart/clear", cartHandler.ClearCart)
		admin.POST("/cart/checkout", cartHandler.Checkout)

		// ==================== 实名认证 ====================
		certificationSvc := service.NewCertificationService(deps.DB, deps.Log)
		certificationHandler := handler.NewCertificationHandler(certificationSvc, deps.Log)
		admin.POST("/certifications", certificationHandler.Submit)
		admin.GET("/certifications/status", certificationHandler.GetStatus)
		admin.GET("/certifications", certificationHandler.AdminGetList)
		admin.POST("/certifications/:id/review", certificationHandler.AdminReview)

		// ==================== 客户关怀 ====================
		clientCareHandler := handler.NewClientCareHandler(deps.DB, deps.Log)
		admin.GET("/client-care/rules", clientCareHandler.GetRules)
		admin.POST("/client-care/rules", clientCareHandler.CreateRule)
		admin.PUT("/client-care/rules/:id", clientCareHandler.UpdateRule)
		admin.DELETE("/client-care/rules/:id", clientCareHandler.DeleteRule)
		admin.GET("/client-care/logs", clientCareHandler.GetLogs)

		// ==================== 客户联系人 ====================
		clientContactSvc := service.NewClientContactService(deps.DB)
		clientContactHandler := handler.NewClientContactHandler(clientContactSvc, deps.Log)
		admin.GET("/client-contacts", clientContactHandler.List)
		admin.GET("/client-contacts/:id", clientContactHandler.Get)
		admin.POST("/client-contacts", clientContactHandler.Create)
		admin.PUT("/client-contacts/:id", clientContactHandler.Update)
		admin.DELETE("/client-contacts/:id", clientContactHandler.Delete)

		// ==================== 消息配置 ====================
		configMessageSvc := service.NewConfigMessageService(deps.DB, deps.Log)
		configMessageHandler := handler.NewConfigMessageHandler(configMessageSvc, deps.Log)
		admin.GET("/config/messages", configMessageHandler.GetAll)
		admin.GET("/config/messages/:channel", configMessageHandler.GetByChannel)
		admin.PUT("/config/messages/:channel", configMessageHandler.Update)
		admin.POST("/config/messages/:channel/test", configMessageHandler.TestSend)
		admin.GET("/config/messages/enabled", configMessageHandler.GetEnabled)

		// ==================== 联系人 ====================
		contactsSvc := service.NewContactService(deps.DB, deps.Log)
		contactsHandler := handler.NewContactsHandler(contactsSvc, deps.Log)
		admin.GET("/contacts", contactsHandler.GetList)
		admin.GET("/contacts/default", contactsHandler.GetDefault)
		admin.GET("/contacts/:id", contactsHandler.GetDetail)
		admin.POST("/contacts", contactsHandler.Create)
		admin.PUT("/contacts/:id", contactsHandler.Update)
		admin.DELETE("/contacts/:id", contactsHandler.Delete)
		admin.POST("/contacts/:id/default", contactsHandler.SetDefault)

		// ==================== 合同 ====================
		contractHandler := handler.NewContractHandler(deps.DB, deps.Log)
		admin.GET("/contracts", contractHandler.AdminGetList)
		admin.GET("/contracts/:id", contractHandler.AdminGetDetail)
		admin.POST("/contracts", contractHandler.AdminCreate)
		admin.PUT("/contracts/:id", contractHandler.AdminUpdate)
		admin.DELETE("/contracts/:id", contractHandler.AdminDelete)
		admin.POST("/contracts/:id/sign", contractHandler.AdminSign)

		// ==================== 优惠券 ====================
		couponSvc := service.NewCouponService(deps.DB, deps.Log)
		couponHandler := handler.NewCouponHandler(couponSvc)
		admin.POST("/coupons/validate", couponHandler.ValidateCoupon)
		admin.GET("/coupons", couponHandler.ListCoupons)
		admin.POST("/coupons", couponHandler.CreateCoupon)
		admin.PUT("/coupons/:id", couponHandler.UpdateCoupon)
		admin.DELETE("/coupons/:id", couponHandler.DeleteCoupon)

		// ==================== 信用额度 ====================
		creditSvc := service.NewCreditService(deps.DB, deps.Log)
		creditHandler := handler.NewCreditHandler(deps.DB, creditSvc)
		admin.GET("/credit", creditHandler.GetInfo)
		admin.GET("/credit/logs", creditHandler.AdminGetLogs)
		admin.POST("/credit/users/:id/adjust", creditHandler.AdminAdjust)
		admin.POST("/credit/bills/generate", creditHandler.AdminGenerateBills)
		admin.GET("/credit/bills", creditHandler.AdminGetBills)
		admin.POST("/credit/bills/:id/waive-fee", creditHandler.AdminWaiveLateFee)

		// ==================== 货币 ====================
		currencyHandler := handler.NewCurrencyHandler(deps.DB)
		admin.GET("/currencies", currencyHandler.GetAll)
		admin.GET("/currencies", currencyHandler.AdminGetList)
		admin.POST("/currencies", currencyHandler.AdminCreate)
		admin.PUT("/currencies/:id", currencyHandler.AdminUpdate)
		admin.DELETE("/currencies/:id", currencyHandler.AdminDelete)
		admin.PUT("/currencies/:id/rate", currencyHandler.AdminUpdateRate)

		// ==================== 下载 ====================
		downloadSvc := service.NewDownloadService(deps.DB, deps.Log)
		downloadHandler := handler.NewDownloadHandler(downloadSvc, deps.Log)
		admin.GET("/downloads/categories", downloadHandler.AdminGetCategories)
		admin.POST("/downloads/categories", downloadHandler.AdminCreateCategory)
		admin.PUT("/downloads/categories/:id", downloadHandler.AdminUpdateCategory)
		admin.DELETE("/downloads/categories/:id", downloadHandler.AdminDeleteCategory)
		admin.GET("/downloads", downloadHandler.AdminGetFiles)
		admin.GET("/downloads/:id", downloadHandler.AdminGetFile)
		admin.POST("/downloads", downloadHandler.AdminCreateFile)
		admin.PUT("/downloads/:id", downloadHandler.AdminUpdateFile)
		admin.DELETE("/downloads/:id", downloadHandler.AdminDeleteFile)
		admin.POST("/downloads/upload", downloadHandler.AdminUploadFile)
		admin.GET("/downloads/:id/download", downloadHandler.Download)

		// ==================== 获取用户 ====================
		getUserHandler := handler.NewGetUserHandler(deps.DB, deps.Log)
		admin.GET("/get-user/:id", getUserHandler.GetUser)
		admin.GET("/get-users", getUserHandler.GetUsers)
		admin.GET("/get-user/:id/access", getUserHandler.CheckAccess)

		// ==================== 交流/工单流转 ====================
		interflowSvc := service.NewInterflowService(deps.DB, deps.Log)
		interflowHandler := handler.NewInterflowHandler(interflowSvc, deps.Log)
		admin.GET("/interflows", interflowHandler.GetList)
		admin.GET("/interflows/:id", interflowHandler.GetDetail)
		admin.GET("/interflows/product/:product_id", interflowHandler.GetByProduct)
		admin.POST("/interflows", interflowHandler.Create)
		admin.PUT("/interflows/:id", interflowHandler.Update)
		admin.DELETE("/interflows/:id", interflowHandler.Delete)
		admin.POST("/interflows/:id/toggle", interflowHandler.ToggleStatus)

		// ==================== 账单明细 ====================
		invoiceItemSvc := service.NewInvoiceItemService(deps.DB, deps.Log)
		invoiceItemHandler := handler.NewInvoiceItemHandler(invoiceItemSvc, deps.Log)
		admin.GET("/invoice-items/:id", invoiceItemHandler.GetDetail)
		admin.GET("/invoice-items/invoice/:invoice_id", invoiceItemHandler.GetByInvoiceID)
		admin.POST("/invoice-items", invoiceItemHandler.Create)
		admin.PUT("/invoice-items/:id", invoiceItemHandler.Update)
		admin.DELETE("/invoice-items/:id", invoiceItemHandler.Delete)
		admin.POST("/invoice-items/batch-create", invoiceItemHandler.BatchCreate)
		admin.POST("/invoice-items/batch-delete", invoiceItemHandler.BatchDelete)
		admin.POST("/invoice-items/batch-update", invoiceItemHandler.BatchUpdate)
		admin.GET("/invoice-items/invoice/:invoice_id/total", invoiceItemHandler.CalculateInvoiceTotal)
		admin.POST("/invoice-items/:id/discount", invoiceItemHandler.AddDiscount)
		admin.POST("/invoice-items/:id/tax", invoiceItemHandler.AddTax)

		// ==================== 知识库 ====================
		knowledgeSvc := service.NewKnowledgeService(deps.DB, deps.Log)
		knowledgeHandler := handler.NewKnowledgeHandler(knowledgeSvc, deps.Log)
		admin.GET("/knowledge/categories", knowledgeHandler.AdminGetCategories)
		admin.POST("/knowledge/categories", knowledgeHandler.AdminCreateCategory)
		admin.PUT("/knowledge/categories/:id", knowledgeHandler.AdminUpdateCategory)
		admin.DELETE("/knowledge/categories/:id", knowledgeHandler.AdminDeleteCategory)
		admin.GET("/knowledge/articles", knowledgeHandler.AdminGetArticles)
		admin.GET("/knowledge/articles/:id", knowledgeHandler.AdminGetArticle)
		admin.POST("/knowledge/articles", knowledgeHandler.AdminCreateArticle)
		admin.PUT("/knowledge/articles/:id", knowledgeHandler.AdminUpdateArticle)
		admin.DELETE("/knowledge/articles/:id", knowledgeHandler.AdminDeleteArticle)

		// ==================== 菜单组 ====================
		menusHandler := handler.NewMenusHandler(deps.DB, deps.Log)
		admin.GET("/menu-groups", menusHandler.GetMenuList)
		admin.GET("/menu-groups/:id", menusHandler.GetMenu)
		admin.POST("/menu-groups", menusHandler.CreateMenu)
		admin.PUT("/menu-groups/:id", menusHandler.UpdateMenu)
		admin.DELETE("/menu-groups/:id", menusHandler.DeleteMenu)
		admin.POST("/menu-groups/nav-list", menusHandler.SetNavList)
		admin.POST("/menu-groups/:id/active", menusHandler.SetActive)

		// ==================== 新闻 ====================
		newsHandler := handler.NewNewsHandler(deps.DB, deps.Log)
		admin.GET("/news/categories", newsHandler.AdminGetCategories)
		admin.POST("/news/categories", newsHandler.AdminCreateCategory)
		admin.PUT("/news/categories/:id", newsHandler.AdminUpdateCategory)
		admin.DELETE("/news/categories/:id", newsHandler.AdminDeleteCategory)
		admin.GET("/news", newsHandler.AdminGetList)
		admin.GET("/news/:id", newsHandler.AdminGetDetail)
		admin.POST("/news", newsHandler.AdminCreate)
		admin.PUT("/news/:id", newsHandler.AdminUpdate)
		admin.DELETE("/news/:id", newsHandler.AdminDelete)

		// ==================== 通知 ====================
		var wechatSvcForNotif *service.WechatService
		if deps.WechatAppID != "" {
			wechatSvcForNotif = service.NewWechatService(deps.DB, deps.Log, deps.UserSvc,
				deps.WechatAppID, deps.WechatAppSecret, deps.WechatMchID,
				deps.WechatMchKey, deps.WechatNotifyURL, deps.WechatTemplateID)
		}
		notificationSvc := service.NewNotificationService(deps.DB, deps.Log, wechatSvcForNotif)
		notificationHandler := handler.NewNotificationHandler(notificationSvc, deps.Log)
		admin.GET("/notifications", notificationHandler.GetUserNotifications)
		admin.POST("/notifications/:id/read", notificationHandler.MarkRead)
		admin.POST("/notifications/read-all", notificationHandler.MarkAllRead)
		admin.GET("/notifications/templates", notificationHandler.AdminGetTemplates)
		admin.PUT("/notifications/templates/:id", notificationHandler.AdminUpdateTemplate)
		admin.GET("/notifications/logs", notificationHandler.AdminGetLogs)
		admin.POST("/notifications/batch", notificationHandler.AdminSendBatch)

		// ==================== OAuth登录 ====================
		oauthSvc := service.NewOAuthService(deps.DB, deps.Log, deps.UserSvc, deps.FrontendURL)
		oauthHandler := handler.NewOAuthHandler(oauthSvc, deps.Log, deps.JWTKey)
		admin.GET("/oauth/providers", oauthHandler.GetProviders)
		admin.GET("/oauth/:provider", oauthHandler.Login)
		admin.GET("/oauth/:provider/callback", oauthHandler.Callback)
		admin.POST("/oauth/bind", oauthHandler.BindAccount)
		admin.GET("/oauth/accounts", oauthHandler.GetBoundAccounts)
		admin.POST("/oauth/unbind", oauthHandler.UnbindAccount)

		// ==================== OAuth绑定 ====================
		oauthBindSvc := service.NewOAuthBindService(deps.DB, deps.Log, deps.UserSvc)
		oauthBindHandler := handler.NewOAuthBindHandler(oauthBindSvc, deps.Log)
		admin.GET("/oauth-bind/providers", oauthBindHandler.GetProviders)
		admin.POST("/oauth-bind/bind", oauthBindHandler.Bind)
		admin.POST("/oauth-bind/unbind", oauthBindHandler.Unbind)
		admin.GET("/oauth-bind/accounts", oauthBindHandler.GetBoundAccounts)
		admin.GET("/oauth-bind/check/:provider", oauthBindHandler.CheckBinding)

		// ==================== 支付 ====================
		if deps.PaymentTxStore != nil && deps.PaymentInvStore != nil {
			paymentSvc := service.NewPaymentService(deps.PaymentTxStore, deps.PaymentInvStore)
			paymentHandler := handler.NewPaymentHandler(paymentSvc)
			admin.GET("/payment/gateways", gin.WrapF(paymentHandler.GetGateways))
			admin.POST("/payment/create", gin.WrapF(paymentHandler.CreatePayment))
			admin.POST("/payment/alipay/notify", gin.WrapF(paymentHandler.AlipayNotify))
			admin.POST("/payment/wechat/notify", gin.WrapF(paymentHandler.WechatNotify))
			admin.GET("/payment/return", gin.WrapF(paymentHandler.ReturnURL))
		}

		// ==================== 产品转移 ====================
		productDivertSvc := service.NewProductDivertService(deps.DB, deps.Log)
		productDivertHandler := handler.NewProductDivertHandler(productDivertSvc, deps.Log)
		admin.GET("/product-diverts", productDivertHandler.GetSent)
		admin.GET("/product-diverts/received", productDivertHandler.GetReceived)
		admin.GET("/product-diverts/:id", productDivertHandler.GetDetail)
		admin.POST("/product-diverts", productDivertHandler.Create)
		admin.POST("/product-diverts/:id/accept", productDivertHandler.Accept)
		admin.POST("/product-diverts/:id/reject", productDivertHandler.Reject)
		admin.POST("/product-diverts/:id/cancel", productDivertHandler.Cancel)

		// ==================== RBAC权限 ====================
		rbacSvc := service.NewRbacService(deps.DB, deps.Log)
		rbacHandler := handler.NewRbacHandler(rbacSvc, deps.Log)
		admin.GET("/rbac/roles", rbacHandler.GetRoles)
		admin.POST("/rbac/roles", rbacHandler.CreateRole)
		admin.PUT("/rbac/roles/:id", rbacHandler.UpdateRole)
		admin.DELETE("/rbac/roles/:id", rbacHandler.DeleteRole)
		admin.GET("/rbac/permissions", rbacHandler.GetPermissions)
		admin.POST("/rbac/users/:id/roles", rbacHandler.AssignRole)
		admin.GET("/rbac/users/:id/roles", rbacHandler.GetUserRoles)

		// ==================== 报表 ====================
		reportSvc := service.NewReportService(deps.DB, deps.Log)
		reportHandler := handler.NewReportHandler(reportSvc, deps.Log)
		admin.GET("/reports/dashboard", reportHandler.GetDashboard)
		admin.GET("/reports/daily", reportHandler.GetDailyReport)
		admin.GET("/reports/monthly", reportHandler.GetMonthlyReport)
		admin.GET("/reports/revenue", reportHandler.GetRevenueChart)
		admin.GET("/reports/users", reportHandler.GetUserStats)
		admin.GET("/reports/orders", reportHandler.GetOrderStats)
		admin.GET("/reports/top-clients", reportHandler.GetTopClients)
		admin.GET("/reports/product-income", reportHandler.GetProductIncome)

		// ==================== 运行映射 ====================
		runMapHandler := handler.NewRunMapHandler(deps.DB, deps.Log)
		admin.GET("/run-map", runMapHandler.GetRunMapList)
		admin.GET("/run-map/:id", runMapHandler.GetRunMap)
		admin.POST("/run-map/:id/repeat", runMapHandler.RepeatTask)
		admin.GET("/run-map/types", runMapHandler.GetTaskTypes)

		// ==================== 设置 ====================
		setHandler := handler.NewSetHandler(deps.DB, deps.Log)
		admin.GET("/site-settings", setHandler.GetSiteSettings)
		admin.PUT("/site-settings", setHandler.UpdateSiteSettings)
		admin.GET("/site-settings/themes", setHandler.GetAdminThemes)
		admin.POST("/site-settings/themes", setHandler.SetAdminTheme)

		// ==================== 系统消息 ====================
		systemMessageSvc := service.NewSystemMessageService(deps.DB, deps.Log)
		systemMessageHandler := handler.NewSystemMessageHandler(systemMessageSvc, deps.Log)
		admin.GET("/system-messages", systemMessageHandler.GetList)
		admin.GET("/system-messages/unread-count", systemMessageHandler.GetUnreadCount)
		admin.GET("/system-messages/types", systemMessageHandler.GetTypes)
		admin.GET("/system-messages/:id", systemMessageHandler.GetDetail)
		admin.POST("/system-messages/:id/read", systemMessageHandler.MarkRead)
		admin.POST("/system-messages/read-all", systemMessageHandler.MarkAllRead)
		admin.DELETE("/system-messages/:id", systemMessageHandler.Delete)
		admin.DELETE("/system-messages", systemMessageHandler.DeleteAll)

		// ==================== 升级 ====================
		upgradeSvc := service.NewUpgradeService(deps.DB, deps.Log)
		upgradeHandler := handler.NewUpgradeHandler(upgradeSvc, deps.Log)
		admin.GET("/upgrades", upgradeHandler.GetUserUpgrades)
		admin.GET("/upgrades/:id", upgradeHandler.GetUpgradeDetail)
		admin.GET("/upgrades/products/:id", upgradeHandler.GetAvailableUpgrades)
		admin.POST("/upgrades", upgradeHandler.CreateUpgrade)
		admin.POST("/upgrades/:id/pay", upgradeHandler.PayUpgrade)

		// ==================== 上传 ====================
		uploadSvc := service.NewUploadService(deps.DB, deps.Log, deps.UploadDir, deps.BaseURL)
		uploadHandler := handler.NewUploadHandler(uploadSvc, deps.Log)
		admin.POST("/upload", uploadHandler.Upload)
		admin.POST("/upload/avatar", uploadHandler.UploadAvatar)
		admin.GET("/upload", uploadHandler.GetList)
		admin.GET("/upload/:id", uploadHandler.GetDetail)
		admin.DELETE("/upload/:id", uploadHandler.Delete)

		// ==================== 上游供应商 ====================
		upstreamHandler := handler.NewUpstreamHandler(deps.DB, deps.Log)
		admin.GET("/upstream/providers", upstreamHandler.GetProviders)
		admin.POST("/upstream/providers", upstreamHandler.CreateProvider)
		admin.PUT("/upstream/providers/:id", upstreamHandler.UpdateProvider)
		admin.DELETE("/upstream/providers/:id", upstreamHandler.DeleteProvider)
		admin.POST("/upstream/providers/:id/test", upstreamHandler.TestConnection)
		admin.POST("/upstream/providers/:id/sync", upstreamHandler.SyncProducts)
		admin.GET("/upstream/providers/:id/logs", upstreamHandler.GetSyncLogs)

		// ==================== 用户等级 ====================
		userLevelHandler := handler.NewUserLevelHandler(deps.DB)
		admin.GET("/user-levels", userLevelHandler.AdminGetList)
		admin.POST("/user-levels", userLevelHandler.AdminCreate)
		admin.PUT("/user-levels/:id", userLevelHandler.AdminUpdate)
		admin.DELETE("/user-levels/:id", userLevelHandler.AdminDelete)

		// ==================== v10购物车 ====================
		couponSvcForV10 := service.NewCouponService(deps.DB, deps.Log)
		v10CartSvc := service.NewV10CartService(deps.DB, deps.Log, deps.OrdSvc, couponSvcForV10)
		v10CartHandler := handler.NewV10CartHandler(v10CartSvc, deps.Log)
		admin.GET("/v10/cart", v10CartHandler.GetCart)
		admin.GET("/v10/cart/count", v10CartHandler.GetItemCount)
		admin.POST("/v10/cart", v10CartHandler.AddItem)
		admin.PUT("/v10/cart/:id", v10CartHandler.UpdateItem)
		admin.DELETE("/v10/cart/:id", v10CartHandler.RemoveItem)
		admin.DELETE("/v10/cart/clear", v10CartHandler.ClearCart)
		admin.POST("/v10/cart/coupon", v10CartHandler.ApplyCoupon)
		admin.DELETE("/v10/cart/coupon", v10CartHandler.RemoveCoupon)
		admin.POST("/v10/cart/checkout", v10CartHandler.Checkout)

		// ==================== 代金券 ====================
		voucherSvc := service.NewVoucherService(deps.DB, deps.Log)
		voucherHandler := handler.NewVoucherHandler(voucherSvc, deps.Log)
		admin.GET("/vouchers", voucherHandler.AdminGetList)
		admin.POST("/vouchers", voucherHandler.AdminCreate)
		admin.PUT("/vouchers/:id", voucherHandler.AdminUpdate)
		admin.DELETE("/vouchers/:id", voucherHandler.AdminDelete)

		// ==================== 微信 ====================
		if deps.WechatAppID != "" {
			wechatSvc := service.NewWechatService(deps.DB, deps.Log, deps.UserSvc,
				deps.WechatAppID, deps.WechatAppSecret, deps.WechatMchID,
				deps.WechatMchKey, deps.WechatNotifyURL, deps.WechatTemplateID)
			wechatHandler := handler.NewWechatHandler(wechatSvc, deps.Log)
			admin.GET("/wechat/auth", wechatHandler.GetAuthURL)
			admin.GET("/wechat/callback", wechatHandler.Login)
			admin.POST("/wechat/pay/notify", wechatHandler.PayNotify)
			admin.POST("/wechat/message/send", wechatHandler.SendTemplateMessage)
			admin.POST("/wechat/pay/create", wechatHandler.CreatePayOrder)
		}

		// ==================== 智简魔方兼容API ====================
		zjmfHandler := handler.NewZjmfFinanceApiHandler(deps.DB, deps.Log)
		admin.GET("/zjmf-api", zjmfHandler.GetApis)
		admin.GET("/zjmf-api/:id", zjmfHandler.GetApi)
		admin.POST("/zjmf-api", zjmfHandler.CreateApi)
		admin.PUT("/zjmf-api/:id", zjmfHandler.UpdateApi)
		admin.DELETE("/zjmf-api/:id", zjmfHandler.DeleteApi)
		admin.POST("/zjmf-api/:id/test", zjmfHandler.TestConnection)
		admin.POST("/zjmf-api/:id/sync", zjmfHandler.SyncProducts)

		// ==================== 域名管理 ====================
		domainSvc := service.NewDomainService(deps.DB, deps.Log)
		domainHandler := handler.NewDomainHandler(domainSvc, deps.DB, deps.Log)
		admin.GET("/domains", domainHandler.AdminGetList)
		admin.GET("/domains/:id", domainHandler.AdminGetDetail)
		admin.POST("/domains", domainHandler.AdminCreate)
		admin.PUT("/domains/:id", domainHandler.AdminUpdate)
		admin.DELETE("/domains/:id", domainHandler.AdminDelete)
		admin.GET("/domains/transfers", domainHandler.AdminGetTransfers)
		admin.PUT("/domains/transfers/:id", domainHandler.AdminUpdateTransfer)

		// ==================== SSL证书管理 ====================
		sslCertSvc := service.NewSSLCertificateService(deps.DB, deps.Log)
		sslCertHandler := handler.NewSSLCertificateHandler(sslCertSvc, deps.Log)
		admin.GET("/ssl-certificates", sslCertHandler.AdminGetList)
		admin.GET("/ssl-certificates/:id", sslCertHandler.AdminGetDetail)
		admin.PUT("/ssl-certificates/:id", sslCertHandler.AdminUpdate)
		admin.DELETE("/ssl-certificates/:id", sslCertHandler.AdminDelete)
	}
}

// RegisterPublicRoutes registers public-facing API routes.
func RegisterPublicRoutes(r *gin.RouterGroup, db *gorm.DB, log *logger.Logger) {
	publicSvc := service.NewPublicService(db, log)
	publicHandler := handler.NewPublicHandler(publicSvc, log)
	r.GET("/system-info", publicHandler.GetSystemInfo)
	r.GET("/config/:key", publicHandler.GetConfig)
	r.GET("/configs", publicHandler.GetConfigs)

	// 公告（前台）
	announceSvc := service.NewAnnounceService(db, log)
	announceHandler := handler.NewAnnounceHandler(announceSvc, log)
	r.GET("/announcements/active", announceHandler.GetActive)

	// 友情链接（前台）
	friendlyLinkSvc := service.NewFriendlyLinkService(db, log)
	friendlyLinkHandler := handler.NewFriendlyLinkHandler(friendlyLinkSvc, log)
	r.GET("/friendly-links/active", friendlyLinkHandler.GetActive)
}
