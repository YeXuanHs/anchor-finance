package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"anchorfinance/internal/backup"
	"anchorfinance/internal/handler"
	"anchorfinance/internal/api/middleware"
	"anchorfinance/internal/model"
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
	authHandler := handler.NewAuthHandler(deps.UserSvc, deps.Log, deps.JWTMgr)
	userHandler := handler.NewUserHandler(deps.UserSvc, deps.Log)
	productHandler := handler.NewProductHandler(deps.ProdSvc, deps.Log)
	orderHandler := handler.NewOrderHandlerWithDB(deps.OrdSvc, deps.DB, deps.Log)
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

		// 仪表盘新增接口
		dashboardHandler := handler.NewDashboardHandler(deps.DB)
		admin.GET("/dashboard/income-trend", dashboardHandler.GetIncomeTrend)
		admin.GET("/dashboard/product-distribution", dashboardHandler.GetProductDistribution)
		admin.GET("/dashboard/global-search", dashboardHandler.GlobalSearch)
		admin.GET("/dashboard/admin-index", dashboardHandler.GetAdminIndex)

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
		admin.POST("/products/duplicate", productHandler.DuplicateProduct)
		admin.PUT("/products/:id/stock", productHandler.EditStock)
		admin.POST("/products/batch-update", productHandler.BatchUpdate)
		admin.POST("/products/batch-delete", productHandler.BatchDelete)
		admin.POST("/products/sort", productHandler.UpdateProductSort)
		admin.GET("/products/:id/discounts", productHandler.GetDiscountList)
		admin.GET("/products/:id/downloads", productHandler.GetProductDownloads)
		admin.POST("/products/:id/upstream-price", productHandler.GetUpstreamPrice)
		// P1-6: 资源产品编辑
		admin.PUT("/products/:id/res", productHandler.EditResProduct)
		// P3-16: 下拉选项/分类管理
		admin.GET("/products/select-type", productHandler.SelectType)
		admin.GET("/products/selectcates", productHandler.Selectcates)
		admin.GET("/products/downloadcates", productHandler.Downloadcates)
		admin.POST("/products/add-downloadcats", productHandler.AddDownloadcats)

		// 一级分组
		admin.GET("/product-first-groups", productHandler.GetFirstGroups)
		admin.POST("/product-first-groups", productHandler.CreateFirstGroup)
		admin.PUT("/product-first-groups/:id", productHandler.UpdateFirstGroup)
		admin.DELETE("/product-first-groups/:id", productHandler.DeleteFirstGroup)
		admin.POST("/product-first-groups/sort", productHandler.UpdateFirstGroupSort)

		// 二级分组
		admin.GET("/product-groups", productHandler.GetGroups)
		admin.POST("/product-groups", productHandler.CreateGroup)
		admin.PUT("/product-groups/:id", productHandler.UpdateGroup)
		admin.DELETE("/product-groups/:id", productHandler.DeleteGroup)
		admin.POST("/product-groups/sort", productHandler.UpdateGroupSort)
		admin.GET("/product-groups/check-alias", productHandler.CheckAlias)

		// 产品下载管理
		admin.POST("/product-downloads/manage", productHandler.ManageDownloads)
		admin.POST("/product-downloads/add-file", productHandler.AddDownloadFile)

		// 产品自定义字段
		admin.DELETE("/product-custom-fields/:id", productHandler.DeleteCustomField)

		// 订单管理
		admin.GET("/orders", orderHandler.GetList)
		admin.GET("/orders/:id", orderHandler.GetDetail)
		admin.POST("/orders/:id/status", orderHandler.UpdateStatus)
		admin.POST("/orders", orderHandler.AdminCreateOrder)
		admin.GET("/orders/:id/check", orderHandler.CheckOrder)
		admin.POST("/orders/batch", orderHandler.BatchUpdate)
		admin.DELETE("/orders/:id", orderHandler.Delete)
		admin.POST("/orders/:id/notes", orderHandler.AddNote)
		admin.GET("/orders/:id/notes", orderHandler.GetNotes)
		admin.GET("/orders/sale", orderHandler.GetSaleOrders)
		admin.POST("/orders/:id/activate", orderHandler.ActivateOrder)
		admin.POST("/orders/:id/change-status", orderHandler.ChangeStatus)
		admin.POST("/orders/multi-total", orderHandler.GetMultiTotal)
		admin.POST("/orders/custom-promo", orderHandler.ApplyCustomPromo)
		admin.GET("/orders/search-page", orderHandler.SearchPage)
		// P0-5: 产品试用/资格校验
		admin.POST("/orders/check-product", orderHandler.CheckProduct)
		// P3-17: 下单配置和客户下拉
		admin.GET("/orders/set-config", orderHandler.SetConfig)
		admin.GET("/orders/clients", orderHandler.GetClients)
		// P2-13: 增强订单列表
		admin.GET("/orders-enhanced", orderHandler.GetListEnhanced)

		// 账单管理（静态路由必须在参数路由之前注册）
		admin.GET("/invoices", invoiceHandler.GetList)
		// P2-14: 增强发票列表
		admin.GET("/invoices-enhanced", invoiceHandler.GetListEnhanced)

		// 账单增强功能
		invoiceEnhancedSvc := service.NewInvoiceEnhancedService(deps.DB, deps.Log)
		emailEnhancedSvc := service.NewEmailEnhancedService(deps.DB, deps.Log)
		invoiceEnhancedHandler := handler.NewInvoiceEnhancedHandler(invoiceEnhancedSvc, emailEnhancedSvc, deps.Log)

		// 静态路由（必须在 /invoices/:id 之前）
		admin.GET("/invoices/paid", invoiceEnhancedHandler.GetPaidInvoices)
		admin.GET("/invoices/unpaid", invoiceEnhancedHandler.GetUnpaidInvoices)
		admin.GET("/invoices/cancelled", invoiceEnhancedHandler.GetCancelledInvoices)
		admin.GET("/invoices/overdue", invoiceEnhancedHandler.GetOverdueInvoices)
		admin.GET("/invoices/summary", invoiceEnhancedHandler.GetInvoiceSummary)
		admin.GET("/invoices/search", invoiceEnhancedHandler.SearchInvoices)
		admin.GET("/invoices/search-page", invoiceEnhancedHandler.SearchPage)
		admin.GET("/invoices/renew-list", invoiceEnhancedHandler.GetRenewInvoices)
		admin.POST("/invoices/combine", invoiceEnhancedHandler.CombineInvoices)
		admin.POST("/invoices/renew", invoiceEnhancedHandler.CreateRenewInvoice)

		// 参数路由
		admin.GET("/invoices/status/:status", invoiceEnhancedHandler.GetByStatus)
		admin.GET("/invoices/combine/user/:user_id", invoiceEnhancedHandler.GetCombineInvoices)
		admin.DELETE("/invoices/notes/:note_id", invoiceEnhancedHandler.DeleteNote)

		// /invoices/:id 相关路由（必须在最后）
		admin.GET("/invoices/:id", invoiceHandler.GetDetail)
		admin.POST("/invoices/:id/cancel", invoiceHandler.Cancel)
		admin.POST("/invoices/:id/pay", invoiceEnhancedHandler.AddPayInvoice)
		admin.POST("/invoices/:id/apply-credit-limit", invoiceEnhancedHandler.ApplyCreditLimit)
		admin.POST("/invoices/:id/execute-renew", invoiceEnhancedHandler.ExecuteRenew)
		admin.GET("/invoices/:id/notes-page", invoiceHandler.NotesPage)
		admin.PUT("/invoices/:id/notes-page", invoiceHandler.Notes)
		admin.DELETE("/invoices/:id/pay", invoiceEnhancedHandler.DeletePayInvoice)
		admin.POST("/invoices/:id/refund", invoiceEnhancedHandler.RefundInvoice)
		admin.GET("/invoices/:id/refund-page", invoiceEnhancedHandler.GetRefundPage)
		admin.POST("/invoices/:id/notes", invoiceEnhancedHandler.AddNote)
		admin.GET("/invoices/:id/notes", invoiceEnhancedHandler.GetNotes)
		admin.POST("/invoices/:id/email", invoiceEnhancedHandler.SendInvoiceEmail)
		admin.GET("/invoices/:id/email-template", invoiceEnhancedHandler.InvoiceEmail)
		admin.GET("/invoices/:id/logs", invoiceEnhancedHandler.GetInvoiceLog)
		admin.POST("/invoices/:id/duplicate", invoiceEnhancedHandler.DuplicateInvoice)
		admin.PUT("/invoices/items/:id", invoiceHandler.EditItem)
		admin.DELETE("/invoices/items", invoiceHandler.DeleteItems)
		admin.DELETE("/invoices/:id/account", invoiceHandler.DelAccount)

		// 工单管理
		admin.GET("/tickets", ticketHandler.GetList)
		// P2-15: 增强工单列表
		admin.GET("/tickets-enhanced", ticketHandler.GetListEnhanced)
		admin.GET("/tickets/:id", ticketHandler.GetDetail)
		admin.POST("/tickets/:id/reply", ticketHandler.AdminReply)
		admin.POST("/tickets/:id/assign", ticketHandler.Assign)
		admin.POST("/tickets/:id/close", ticketHandler.Close)
		// P1-8: 附件下载
		admin.GET("/tickets/:id/attachment/:aid", ticketHandler.DownloadAttachment)
		// P1-9: 工单接单
		admin.POST("/tickets/:id/receive", ticketHandler.TicketReceive)
		admin.POST("/tickets/merge", ticketHandler.MergeTickets)
		admin.POST("/tickets/:id/transfer", ticketHandler.TransferTicket)
		admin.GET("/tickets/:id/transfer-logs", ticketHandler.GetTransferLogs)
		admin.POST("/tickets/:id/attachments", ticketHandler.UploadAttachment)
		admin.GET("/tickets/:id/attachments", ticketHandler.GetAttachments)
		admin.DELETE("/tickets/attachments/:id", ticketHandler.DeleteAttachment)
		admin.GET("/tickets/statistics", ticketHandler.TicketStatistics)

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

		// ==================== 管理员管理 ====================
		adminManageSvc := service.NewAdminManageService(deps.DB, deps.Log)
		adminManageHandler := handler.NewAdminManageHandler(adminManageSvc, deps.Log)
		admin.GET("/admins", adminManageHandler.List)
		admin.GET("/admins/:id", adminManageHandler.Get)
		admin.POST("/admins", adminManageHandler.Create)
		admin.PUT("/admins/:id", adminManageHandler.Update)
		admin.DELETE("/admins/:id", adminManageHandler.Delete)
		admin.POST("/admins/:id/status", adminManageHandler.SetStatus)
		admin.POST("/admins/:id/reset-password", adminManageHandler.ResetPassword)
		admin.GET("/admins/operation-logs", adminManageHandler.GetOperationLogs)

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
		admin.GET("/payment-gateways/supported", paymentGwHandler.GetSupportedInfo)
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
		dcimAdvHandler := handler.NewDcimAdvancedHandler(dcimSvc, deps.Log)
		admin.GET("/dcim/servers", dcimHandler.GetServerList)
		admin.GET("/dcim/servers/:id", dcimHandler.GetServerDetail)
		admin.POST("/dcim/servers", dcimHandler.CreateServer)
		admin.PUT("/dcim/servers/:id", dcimHandler.UpdateServer)
		admin.DELETE("/dcim/servers/:id", dcimHandler.DeleteServer)
		admin.POST("/dcim/servers/:id/boot", dcimHandler.BootServer)
		admin.POST("/dcim/servers/:id/shutdown", dcimHandler.ShutdownServer)
		admin.POST("/dcim/servers/:id/reboot", dcimHandler.RebootServer)
		admin.GET("/dcim/datacenters", dcimHandler.GetDatacenterList)

		// DCIM扩展功能
		admin.GET("/dcim/servers/:id/flow-packets", dcimHandler.ListFlowPacket)
		admin.POST("/dcim/servers/:id/flow-packets", dcimHandler.AddFlowPacket)
		admin.PUT("/dcim/flow-packets/:id", dcimHandler.EditFlowPacket)
		admin.DELETE("/dcim/flow-packets/:id", dcimHandler.DelFlowPacket)
		admin.POST("/dcim/servers/:id/assign", dcimHandler.AssignServer)
		admin.GET("/dcim/sales-servers", dcimHandler.GetSalesServer)
		admin.POST("/dcim/servers/:id/refresh", dcimHandler.RefreshServerStatus)
		admin.GET("/dcim/buy-records", dcimHandler.ListBuyRecord)
		admin.DELETE("/dcim/buy-records/:id", dcimHandler.DelRecord)
		admin.GET("/dcim/flow-packets/create", dcimHandler.AddFlowPacketPage)
		admin.GET("/dcim/flow-packets/:id/edit", dcimHandler.EditFlowPacketPage)
		admin.POST("/dcim/servers/:id/unsuspend-reload", dcimHandler.UnsuspendReload)
		admin.GET("/dcim/servers/:id/detail", dcimHandler.Detail)

		// DCIM高级操作 - KVM/IPMI/BMC
		admin.GET("/dcim/servers/:id/kvm", dcimAdvHandler.GetKVMURL)
		admin.GET("/dcim/servers/:id/bmc", dcimAdvHandler.GetBMCInfo)
		admin.GET("/dcim/servers/:id/novnc", dcimAdvHandler.GetNoVNCURL)
		// DCIM高级操作 - 救援系统
		admin.POST("/dcim/servers/:id/rescue", dcimAdvHandler.BootRescue)
		admin.POST("/dcim/servers/:id/crack-password", dcimAdvHandler.CrackPassword)
		admin.GET("/dcim/servers/:id/rescue-status", dcimAdvHandler.GetRescueStatus)
		admin.GET("/dcim/rescue-logs", dcimAdvHandler.GetRescueLogs)
		// DCIM高级操作 - 流量监控
		admin.GET("/dcim/servers/:id/traffic", dcimAdvHandler.GetTrafficUsage)
		admin.GET("/dcim/servers/:id/traffic/chart", dcimAdvHandler.GetTrafficChart)
		admin.POST("/dcim/servers/:id/traffic/reset", dcimAdvHandler.ResetTrafficCounter)
		// DCIM高级操作 - 快照
		admin.POST("/dcim/servers/:id/snapshots", dcimAdvHandler.CreateSnapshot)
		admin.GET("/dcim/servers/:id/snapshots", dcimAdvHandler.GetSnapshots)
		admin.POST("/dcim/snapshots/:id/restore", dcimAdvHandler.RestoreSnapshot)
		admin.DELETE("/dcim/snapshots/:id", dcimAdvHandler.DeleteSnapshot)
		// DCIM高级操作 - 备份
		admin.POST("/dcim/servers/:id/backups", dcimAdvHandler.CreateBackup)
		admin.GET("/dcim/servers/:id/backups", dcimAdvHandler.GetBackups)
		admin.POST("/dcim/backups/:id/restore", dcimAdvHandler.RestoreBackup)
		admin.DELETE("/dcim/backups/:id", dcimAdvHandler.DeleteBackup)
		// DCIM高级操作 - 电源管理
		admin.GET("/dcim/servers/:id/power", dcimAdvHandler.GetPowerStatus)
		admin.POST("/dcim/servers/:id/power/refresh", dcimAdvHandler.RefreshPowerStatus)
		// DCIM高级操作 - 重装增强
		admin.GET("/dcim/servers/:id/reinstall-status", dcimAdvHandler.GetReinstallStatus)
		admin.POST("/dcim/servers/:id/cancel-reinstall", dcimAdvHandler.CancelReinstall)
		admin.GET("/dcim/os-list", dcimAdvHandler.GetOSList)

		// 插件管理
		pluginSvc := service.NewPluginService(deps.DB, deps.Log)
		pluginHandler := handler.NewPluginHandler(pluginSvc, deps.Log)
		admin.GET("/plugins", pluginHandler.List)
		admin.GET("/plugins/types", pluginHandler.GetTypes)
		admin.POST("/plugins", pluginHandler.Install)
		admin.DELETE("/plugins/:id", pluginHandler.Uninstall)
		admin.POST("/plugins/:id/enable", pluginHandler.Enable)
		admin.POST("/plugins/:id/disable", pluginHandler.Disable)
		admin.PUT("/plugins/:id/config", pluginHandler.UpdateConfig)
		admin.GET("/plugins/type/:type", pluginHandler.GetByType)

		// 服务器模块管理
		admin.GET("/server-modules", pluginHandler.ListServerModules)
		admin.POST("/server-modules", pluginHandler.CreateServerModule)
		admin.PUT("/server-modules/:id", pluginHandler.UpdateServerModule)
		admin.DELETE("/server-modules/:id", pluginHandler.DeleteServerModule)

		// 服务器分组管理
		admin.GET("/server-groups", pluginHandler.ListServerGroups)
		admin.POST("/server-groups", pluginHandler.CreateServerGroup)
		admin.PUT("/server-groups/:id", pluginHandler.UpdateServerGroup)
		admin.DELETE("/server-groups/:id", pluginHandler.DeleteServerGroup)

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

		// 模块供给扩展（Provision Module System）
		pmHandler := handler.NewProvisionModuleHandler(provisionSvc, deps.Log)
		// 模块管理
		admin.GET("/provision-modules", pmHandler.GetModules)
		admin.GET("/provision-modules/default-buttons", pmHandler.GetDefaultButtons)
		admin.GET("/provision-modules/:id", pmHandler.GetModule)
		admin.POST("/provision-modules", pmHandler.CreateModule)
		admin.PUT("/provision-modules/:id", pmHandler.UpdateModule)
		admin.DELETE("/provision-modules/:id", pmHandler.DeleteModule)
		admin.POST("/provision-modules/:id/test", pmHandler.TestModule)
		// 按钮管理
		admin.GET("/provision-modules/:id/buttons", pmHandler.GetButtons)
		admin.POST("/provision-modules/buttons", pmHandler.CreateButton)
		admin.PUT("/provision-modules/buttons/:id", pmHandler.UpdateButton)
		admin.DELETE("/provision-modules/buttons/:id", pmHandler.DeleteButton)
		// 自定义函数
		admin.GET("/provision-modules/:id/functions", pmHandler.GetCustomFunctions)
		admin.POST("/provision-modules/functions", pmHandler.CreateCustomFunction)
		admin.POST("/provision-modules/functions/:id/execute", pmHandler.ExecuteCustomFunction)
		admin.DELETE("/provision-modules/functions/:id", pmHandler.DeleteCustomFunction)
		// 图表
		admin.GET("/provision-modules/charts/:host_id", pmHandler.GetCharts)
		admin.GET("/provision-modules/charts/:host_id/:chart_id", pmHandler.GetChartData)
		admin.POST("/provision-modules/charts", pmHandler.CreateChart)
		admin.PUT("/provision-modules/charts/:id", pmHandler.UpdateChart)
		// 管理端区域
		admin.GET("/provision-modules/admin-area/:host_id", pmHandler.RenderAdminArea)
		admin.GET("/provision-modules/admin-area/:host_id/buttons", pmHandler.GetAdminButtons)
		admin.POST("/provision-modules/admin-area/:host_id/execute", pmHandler.ExecuteAdminButton)
		// 客户端区域
		admin.GET("/provision-modules/client-area/:host_id", pmHandler.RenderClientArea)
		admin.GET("/provision-modules/client-area/:host_id/detail", pmHandler.RenderClientAreaDetail)
		admin.GET("/provision-modules/client-area/:host_id/buttons", pmHandler.GetClientButtons)
		admin.POST("/provision-modules/client-area/:host_id/execute", pmHandler.ExecuteClientButton)
		// 用量追踪
		admin.GET("/provision-modules/usage/:host_id", pmHandler.GetUsage)
		admin.PUT("/provision-modules/usage/:host_id", pmHandler.UpdateUsage)
		admin.GET("/provision-modules/usage/:host_id/check", pmHandler.CheckDefineUsage)
		admin.GET("/provision-modules/usage/:host_id/traffic", pmHandler.TrafficUsage)
		// SSL/下载
		admin.GET("/provision-modules/ssl/:host_id", pmHandler.SSLButton)
		admin.GET("/provision-modules/download/:host_id/:resource_id", pmHandler.DownloadResource)
		// 流量包
		admin.POST("/provision-modules/flow-packet/paid", pmHandler.AfterFlowPacketPaid)

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
		// 新增缺失路由
		admin.POST("/client-services/transfer", clientSvcHandler.PostTransfer)
		admin.DELETE("/client-services/:id/delete-host", clientSvcHandler.DeleteHost)
		admin.GET("/client-services/batch-renew-page", clientSvcHandler.PostBatchRenewPage)
		admin.POST("/client-services/batch-renew", clientSvcHandler.PostBatchRenew)
		admin.GET("/client-services/apply-credit-page", clientSvcHandler.GetApplyCreditPage)
		admin.POST("/client-services/apply-credit", clientSvcHandler.ApplyCredit)
		admin.POST("/client-services/search-client", clientSvcHandler.PostSearchClient)
		admin.GET("/client-services/:id/refund-page", clientSvcHandler.GetRefundPage)
		admin.POST("/client-services/:id/refund", clientSvcHandler.Refund)
		// P1-10: 单台主机续费
		admin.POST("/client-services/:id/host-renew", clientSvcHandler.HostRenew)
		// P1-11: 升级配置
		admin.POST("/client-services/:id/upgrade-config", clientSvcHandler.UpgradeConfig)
		// P3-20: 联动筛选
		admin.GET("/client-services/linkage-list", clientSvcHandler.AdminGetLinkAgeList)
		admin.GET("/client-services/product-list", clientSvcHandler.GetProductList)

		// 用户管理
		userManageSvc := service.NewUserManageService(deps.DB, deps.Log)
		userManageHandler := handler.NewUserManageHandler(userManageSvc, deps.Log)
		// Search & Filter
		admin.GET("/user-manage/search", userManageHandler.Search)
		admin.GET("/user-manage/filter", userManageHandler.Filter)
		admin.GET("/user-manage/summary", userManageHandler.GetSummary)
		// Client Lifecycle
		admin.POST("/user-manage", userManageHandler.Create)
		admin.POST("/user-manage/:id/close", userManageHandler.Close)
		admin.DELETE("/user-manage/:id", userManageHandler.Delete)
		admin.POST("/user-manage/:id/ban", userManageHandler.Ban)
		admin.POST("/user-manage/:id/unban", userManageHandler.Unban)
		admin.POST("/user-manage/:id/cancel-ban", userManageHandler.CancelBan)
		// Client Profile
		admin.GET("/user-manage/:id/profile", userManageHandler.GetProfile)
		admin.PUT("/user-manage/:id/profile", userManageHandler.UpdateProfile)
		admin.GET("/user-manage/:id/hosts", userManageHandler.GetHosts)
		admin.GET("/user-manage/:id/invoices", userManageHandler.GetInvoices)
		admin.GET("/user-manage/:id/orders", userManageHandler.GetOrders)
		admin.GET("/user-manage/:id/tickets", userManageHandler.GetTickets)
		admin.GET("/user-manage/:id/client-logs", userManageHandler.GetLogs)
		// Client Notes
		admin.POST("/user-manage/:id/notes", userManageHandler.AddNote)
		admin.GET("/user-manage/:id/notes", userManageHandler.GetNotes)
		admin.DELETE("/user-manage/notes/:id", userManageHandler.DeleteNote)
		// Client Authorization
		admin.POST("/user-manage/:id/authorize", userManageHandler.Authorize)
		admin.GET("/user-manage/:id/auth", userManageHandler.GetAuth)
		// Client Group
		admin.POST("/user-manage/:id/group", userManageHandler.AssignGroup)
		admin.DELETE("/user-manage/:id/group/:group_id", userManageHandler.RemoveFromGroup)
		// Certification
		admin.GET("/user-manage/:id/certification", userManageHandler.GetCertificationStatus)
		admin.POST("/user-manage/:id/certification/review", userManageHandler.ReviewCertification)
		// Cancel Requests
		admin.GET("/user-manage/cancel-requests", userManageHandler.GetCancelRequests)
		admin.POST("/user-manage/cancel-requests/:id", userManageHandler.ProcessCancelRequest)
		// Balance & Password
		admin.POST("/user-manage/:id/balance", userManageHandler.AdjustBalance)
		admin.POST("/user-manage/:id/reset-password", userManageHandler.ResetPassword)
		admin.GET("/user-manage/:id/status", userManageHandler.GetStatus)
		admin.GET("/user-manage/:id/operation-logs", userManageHandler.GetOperationLogs)
		// 新增缺失路由
		admin.GET("/user-manage/:uid/hosts-by-uid", userManageHandler.HostByUid)
		admin.GET("/user-manage/certification/list", userManageHandler.CerifyList)
		admin.GET("/user-manage/certification/log-list", userManageHandler.CerifyLogList)
		admin.GET("/user-manage/certification/person/:id", userManageHandler.CertifiPersonDetail)
		admin.PUT("/user-manage/certification/person/:id", userManageHandler.CertifiPersonModify)
		admin.GET("/user-manage/certification/company/:id", userManageHandler.CertifiCompanyDetail)
		admin.PUT("/user-manage/certification/company/:id", userManageHandler.CertifiCompanyModify)
		admin.GET("/user-manage/:id/product-accounts", userManageHandler.UserProductaccounts)
		admin.POST("/user-manage/:id/login-as", userManageHandler.LoginByUser)
		admin.DELETE("/user-manage/cancel-requests/:id", userManageHandler.DeleteCancelRequest)
		admin.POST("/user-manage/:id/record-log", userManageHandler.AddRecordLog)
		admin.POST("/user-manage/:id/remark-log", userManageHandler.AddRemarkLog)
		// P1-7: 黑名单管理
		admin.GET("/user-manage/black-list", userManageHandler.GetBlackList)
		admin.DELETE("/user-manage/black-list/:id", userManageHandler.RemoveBlackList)
		// P3-21: 用户发票列表
		admin.GET("/user-manage/:id/user-invoice", userManageHandler.UserInvoice)

		// 用户备注
		userRemarkSvc := service.NewUserRemarkService(deps.DB, deps.Log)
		userRemarkHandler := handler.NewUserRemarkHandler(userRemarkSvc, deps.Log)
		admin.GET("/user-remarks", userRemarkHandler.List)
		admin.POST("/user-remarks", userRemarkHandler.Add)
		admin.PUT("/user-remarks/:id", userRemarkHandler.Update)
		admin.DELETE("/user-remarks/:id", userRemarkHandler.AdminDelete)

		// 邮件模板
		emailTplSvc := service.NewEmailTemplateService(deps.DB, deps.Log)
		emailTplHandler := handler.NewEmailTemplateHandler(emailTplSvc, deps.Log, deps.DB)
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

		// 语言管理
		langSvc := service.NewLanguageService(deps.DB, deps.Log)
		langHandler := handler.NewLanguageHandler(langSvc, deps.Log)
		admin.GET("/languages", langHandler.GetLanguages)
		admin.POST("/languages", langHandler.CreateLanguage)
		admin.PUT("/languages/:id", langHandler.UpdateLanguage)
		admin.DELETE("/languages/:id", langHandler.DeleteLanguage)
		admin.POST("/languages/:id/default", langHandler.SetDefaultLanguage)
		admin.GET("/languages/:code/translations", langHandler.GetTranslations)
		admin.POST("/languages/:code/translations", langHandler.SaveTranslations)
		admin.POST("/languages/:code/import", langHandler.ImportTranslations)
		admin.GET("/lang-keys", langHandler.GetLangKeys)

		// 优惠码管理
		promoHandler := handler.NewPromoCodeHandler(deps.DB, deps.Log)
		admin.GET("/promo-codes", promoHandler.GetList)
		admin.GET("/promo-codes/:id", promoHandler.GetDetail)
		admin.POST("/promo-codes", promoHandler.Create)
		admin.PUT("/promo-codes/:id", promoHandler.Update)
		admin.DELETE("/promo-codes/:id", promoHandler.Delete)
		admin.POST("/promo-codes/:id/status", promoHandler.SetStatus)
		admin.POST("/promo-codes/:id/expire", promoHandler.ExpireImmediately)
		admin.GET("/promo-codes/auto-generate", promoHandler.AutoGenerate)
		admin.GET("/promo-codes/:id/logs", promoHandler.GetUsageLogs)

		// 前台导航分组管理
		navGroupHandler := handler.NewNavGroupHandler(deps.DB)
		admin.GET("/nav-groups", navGroupHandler.List)
		admin.POST("/nav-groups", navGroupHandler.Create)
		admin.PUT("/nav-groups/:id", navGroupHandler.Update)
		admin.DELETE("/nav-groups/:id", navGroupHandler.Delete)
		admin.GET("/nav-groups/:id/products", navGroupHandler.GetProducts)
		admin.PUT("/nav-groups/:id/products", navGroupHandler.UpdateProducts)

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
		admin.GET("/sales/statistics", saleHandler.GetStatistics)
		admin.GET("/sales/records", saleHandler.GetSaleRecords)
		admin.GET("/sales/:id/users", saleHandler.GetSaleUsers)
		admin.GET("/sales/admin-list", saleHandler.GetAdminList)
		admin.GET("/sales/status", saleHandler.GetSaleStatus)
		admin.POST("/sales/status", saleHandler.SetSaleStatus)

		// 销售扩展功能
		admin.GET("/sale/timetypes", saleHandler.GetTimetype)
		admin.DELETE("/sale/ladder/:id", saleHandler.DelSaleLadder)

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

		// 短信详细管理
		smsSvc := service.NewSMSService(deps.DB, deps.Log)
		smsHandler := handler.NewSMSHandler(smsSvc, deps.Log)
		// 运营商检测
		admin.GET("/sms/detect-operator", smsHandler.DetectOperator)
		admin.GET("/sms/validate-phone", smsHandler.ValidatePhone)
		// 模板管理
		admin.GET("/sms/templates", smsHandler.GetTemplates)
		admin.GET("/sms/templates/:id", smsHandler.GetTemplate)
		admin.POST("/sms/templates", smsHandler.CreateTemplate)
		admin.PUT("/sms/templates/:id", smsHandler.UpdateTemplate)
		admin.DELETE("/sms/templates/:id", smsHandler.DeleteTemplate)
		// 发送
		admin.POST("/sms/send", smsHandler.SendSMS)
		admin.POST("/sms/send-batch", smsHandler.SendBatchSMS)
		admin.POST("/sms/send-marketing", smsHandler.SendMarketingSMS)
		// 日志
		admin.GET("/sms/logs", smsHandler.GetSMSLogs)
		admin.GET("/sms/logs/:id", smsHandler.GetSMSLog)
		admin.GET("/sms/logs/phone/:phone", smsHandler.GetSMSLogByPhone)
		admin.GET("/sms/logs/user/:user_id", smsHandler.GetSMSLogByUser)
		// 统计
		admin.GET("/sms/stats", smsHandler.GetSMSStats)
		admin.GET("/sms/stats/operator", smsHandler.GetOperatorStats)
		// 批次
		admin.POST("/sms/batches", smsHandler.CreateBatch)
		admin.GET("/sms/batches", smsHandler.GetBatches)
		admin.POST("/sms/batches/:id/execute", smsHandler.ExecuteBatch)

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
		admin.GET("/hosts", hostHandler.GetList)
		admin.GET("/hosts/:id", hostHandler.GetDetail)
		admin.POST("/hosts/:id/boot", hostHandler.Boot)
		admin.POST("/hosts/:id/shutdown", hostHandler.Shutdown)
		admin.POST("/hosts/:id/reboot", hostHandler.Reboot)
		// P3-18: 时间类型选项
		admin.GET("/hosts/timetype", hostHandler.GetTimetype)

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

		// ConfigOptions admin methods (from zjmf)
		admin.GET("/config-options/groups-list", configOptionHandler.GroupsList)
		admin.GET("/config-options/search-page", configOptionHandler.SearchPage)
		admin.GET("/config-options/create-groups", configOptionHandler.CreateGroups)
		admin.POST("/config-options/create-groups", configOptionHandler.CreateGroupsPost)
		admin.GET("/config-options/edit-groups", configOptionHandler.EditGroups)
		admin.POST("/config-options/edit-groups", configOptionHandler.EditGroupsPost)
		admin.GET("/config-options/add-options-page", configOptionHandler.AddOptionsPage)
		admin.POST("/config-options/add-options", configOptionHandler.AddOptions)
		admin.DELETE("/config-options/sub-options", configOptionHandler.DeleteSubOptions)
		admin.DELETE("/config-options/options", configOptionHandler.DeleteOptions)
		admin.DELETE("/config-options/groups", configOptionHandler.DeleteGroups)
		admin.GET("/config-options/duplicate-groups", configOptionHandler.DuplicateGroups)
		admin.POST("/config-options/duplicate-groups", configOptionHandler.DuplicateGroupPost)
		admin.GET("/config-options/check-os", configOptionHandler.ConfigOptionsCheckOs)
		admin.GET("/config-options/os", configOptionHandler.ConfigOptionsOs)
		admin.GET("/config-options/edit-config", configOptionHandler.EditConfig)
		admin.GET("/config-options/next-linkage-list", configOptionHandler.GetNextLinkAgeList)
		admin.POST("/config-options/edit-config", configOptionHandler.EditConfigPost)
		admin.POST("/config-options/save-linkage-level", configOptionHandler.SaveLinkAgeLevel)
		admin.POST("/config-options/save-config-option-info", configOptionHandler.SaveConfigOptionInfo)
		admin.POST("/config-options/save-linkage-order", configOptionHandler.SaveLinkAgeOrder)
		admin.DELETE("/config-options/linkage-sub", configOptionHandler.DelLinkAgeSub)

		// ConfigServers admin methods (from zjmf)
		admin.GET("/config-servers/server-list", configServerHandler.ServerList)
		admin.GET("/config-servers/add", configServerHandler.AddServers)
		admin.GET("/config-servers/modules-group", configServerHandler.GetModulesGroup)
		admin.POST("/config-servers/add", configServerHandler.AddServersPost)
		admin.GET("/config-servers/edit/:id", configServerHandler.EditServers)
		admin.POST("/config-servers/edit", configServerHandler.EditServersPost)
		admin.DELETE("/config-servers/delete/:id", configServerHandler.DeleteServers)
		admin.GET("/config-servers/groups-list", configServerHandler.GroupsList)
		admin.GET("/config-servers/create-groups", configServerHandler.CreateGroups)
		admin.POST("/config-servers/create-groups", configServerHandler.CreateGroupsPost)
		admin.GET("/config-servers/edit-groups/:id", configServerHandler.EditServerGroups)
		admin.POST("/config-servers/edit-groups", configServerHandler.EditServerGroupsPost)
		admin.DELETE("/config-servers/delete-groups/:id", configServerHandler.DeleteServerGroups)
		admin.GET("/config-servers/test-link/:id", configServerHandler.TestLink)

		configGeneralSvc := service.NewConfigGeneralService(deps.DB, deps.Log)
		configGeneralHandler := handler.NewConfigGeneralHandler(configGeneralSvc, deps.Log)
		admin.GET("/config/general", configGeneralHandler.Get)
		admin.PUT("/config/general", configGeneralHandler.Update)
		admin.GET("/config/email", configGeneralHandler.GetEmailConfig)
		admin.PUT("/config/email", configGeneralHandler.UpdateEmailConfig)
		admin.GET("/config/email-support", configGeneralHandler.GetEmailSupport)
		admin.PUT("/config/email-support", configGeneralHandler.UpdateEmailSupport)
		admin.GET("/config/affiliate-ladders", configGeneralHandler.GetAffiliateLadders)
		admin.PUT("/config/affiliate-ladders", configGeneralHandler.UpdateAffiliateLadders)
		admin.GET("/config/safe", configGeneralHandler.GetSafeConfig)
		admin.PUT("/config/safe", configGeneralHandler.UpdateSafeConfig)
		admin.GET("/config/recharge", configGeneralHandler.GetRechargeConfig)
		admin.PUT("/config/recharge", configGeneralHandler.UpdateRechargeConfig)
		admin.GET("/config/invoice", configGeneralHandler.GetInvoiceConfig)
		admin.PUT("/config/invoice", configGeneralHandler.UpdateInvoiceConfig)
		admin.GET("/config/register", configGeneralHandler.GetRegisterConfig)
		admin.PUT("/config/register", configGeneralHandler.UpdateRegisterConfig)
		admin.GET("/config/login", configGeneralHandler.GetLoginConfig)
		admin.PUT("/config/login", configGeneralHandler.UpdateLoginConfig)
		admin.GET("/config/api", configGeneralHandler.GetAPIConfig)
		admin.PUT("/config/api", configGeneralHandler.UpdateAPIConfig)
		admin.GET("/config/2fa", configGeneralHandler.GetTwoFactorConfig)
		admin.PUT("/config/2fa", configGeneralHandler.UpdateTwoFactorConfig)
		admin.GET("/config/debug-mode", configGeneralHandler.GetDebugMode)
		admin.PUT("/config/debug-mode", configGeneralHandler.SetDebugMode)
		admin.POST("/config/smtp-test", configGeneralHandler.TestSMTP)
		admin.POST("/config/sms-test", configGeneralHandler.TestSMS)
		admin.GET("/config/payment", configGeneralHandler.GetPaymentConfig)
		admin.PUT("/config/payment", configGeneralHandler.UpdatePaymentConfig)
		admin.GET("/config/sms", configGeneralHandler.GetSmsConfig)
		admin.PUT("/config/sms", configGeneralHandler.UpdateSmsConfig)
		admin.GET("/config/security", configGeneralHandler.GetSecurityConfig)
		admin.PUT("/config/security", configGeneralHandler.UpdateSecurityConfig)
		admin.GET("/config/local", configGeneralHandler.GetLocalConfig)
		admin.PUT("/config/local", configGeneralHandler.UpdateLocalConfig)
		admin.GET("/config/affiliate", configGeneralHandler.GetAffiliateConfig)
		admin.PUT("/config/affiliate", configGeneralHandler.UpdateAffiliateConfig)
		admin.GET("/config/captcha", configGeneralHandler.GetCaptchaConfig)
		admin.PUT("/config/captcha", configGeneralHandler.UpdateCaptchaConfig)
		admin.GET("/config/buy-product", configGeneralHandler.GetBuyProductConfig)
		admin.PUT("/config/buy-product", configGeneralHandler.UpdateBuyProductConfig)
		admin.GET("/config/second-verify", configGeneralHandler.GetSecondVerifyConfig)
		admin.PUT("/config/second-verify", configGeneralHandler.UpdateSecondVerifyConfig)
		admin.GET("/config/language", configGeneralHandler.GetLanguageConfig)
		admin.PUT("/config/language", configGeneralHandler.SetAdminLanguage)
		admin.GET("/config/header-footer", configGeneralHandler.GetHeaderFooter)
		admin.PUT("/config/header-footer", configGeneralHandler.UpdateHeaderFooter)
		admin.GET("/config/login-page", configGeneralHandler.GetNewLoginPageConfig)
		admin.PUT("/config/login-page", configGeneralHandler.UpdateNewLoginPageConfig)
		admin.GET("/config/nav-groups", configGeneralHandler.GetNavGroups)
		admin.POST("/config/nav-groups", configGeneralHandler.CreateNavGroup)
		admin.PUT("/config/nav-groups/:id", configGeneralHandler.UpdateNavGroup)
		admin.DELETE("/config/nav-groups/:id", configGeneralHandler.DeleteNavGroup)

		configCertifiSvc := service.NewConfigCertifiService(deps.DB, deps.Log)
		configCertifiHandler := handler.NewConfigCertifiHandler(configCertifiSvc, deps.Log, deps.DB)
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
		linkKnowledgeHandler := handler.NewLinkKnowledgeHandler(deps.DB, deps.Log)
		admin.GET("/link-knowledge", linkKnowledgeHandler.Index)
		admin.GET("/link-knowledge/:id", linkKnowledgeHandler.Edit)
		admin.POST("/link-knowledge", linkKnowledgeHandler.Create)
		admin.PUT("/link-knowledge/:id", linkKnowledgeHandler.Save)
		admin.DELETE("/link-knowledge/:id", linkKnowledgeHandler.Delete)
		admin.GET("/link-knowledge/index", linkKnowledgeHandler.Index)
		admin.GET("/link-knowledge/:id/edit", linkKnowledgeHandler.Edit)
		admin.POST("/link-knowledge/save", linkKnowledgeHandler.Save)
		admin.POST("/link-knowledge/add", linkKnowledgeHandler.Add)

		// 交易流水管理
		accountHandler := handler.NewAccountHandler(deps.DB, deps.Log)
		admin.GET("/accounts", accountHandler.Index)
		admin.GET("/accounts/create", accountHandler.Create)
		admin.POST("/accounts", accountHandler.Save)
		admin.GET("/accounts/:id", accountHandler.Read)
		admin.PUT("/accounts/:id", accountHandler.Update)
		admin.DELETE("/accounts/:id", accountHandler.Delete)
		admin.GET("/accounts/search-page", accountHandler.SearchPage)
		admin.POST("/accounts/create-invoice", accountHandler.CreateInvoice)

		// 公共接口
		publicSvc := service.NewPublicService(deps.DB, deps.Log)
		publicHandler := handler.NewPublicHandler(publicSvc, deps.Log)
		admin.GET("/public/system-info", publicHandler.GetSystemInfo)
		admin.GET("/public/config/:key", publicHandler.GetConfig)
		admin.GET("/public/configs", publicHandler.GetConfigs)
		admin.POST("/backup", publicHandler.BackupNow)
		admin.POST("/backup/:id/cancel", publicHandler.CancelBackup)

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
		sendMsgBatchHandler := handler.NewSendMessageBatchHandler(sendMsgBatchSvc, deps.Log, deps.DB)
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
		admin.GET("/system/info", systemHandler.GetSystemInfo)
		admin.GET("/system/database", systemHandler.GetDatabaseInfo)
		admin.POST("/system/database/optimize", systemHandler.OptimizeTables)
		admin.POST("/system/database/backup", systemHandler.BackupDatabase)
		admin.GET("/system/auto-update", systemHandler.GetAutoUpdateConfig)
		admin.PUT("/system/auto-update", systemHandler.UpdateAutoUpdateConfig)
		admin.GET("/system/authorize", systemHandler.GetAuthorizeInfo)
		admin.PUT("/system/license", systemHandler.SetLicense)
		admin.GET("/system/data-migration", systemHandler.GetDataMigration)
		admin.POST("/system/data-migration", systemHandler.StartDataMigration)
		admin.GET("/system/auth-rules", systemHandler.GetSystemAuthRules)
		admin.GET("/system/language", systemHandler.GetSystemLanguage)

		// 设置模块
		settingHandler := handler.NewSettingHandler(deps.Log, deps.DB)
		admin.GET("/setting/notification", settingHandler.GetNotificationSettings)
		admin.PUT("/setting/notification", settingHandler.SaveNotificationSettings)
		admin.GET("/setting/maintenance", settingHandler.GetMaintenanceMode)
		admin.PUT("/setting/maintenance", settingHandler.SetMaintenanceMode)
		admin.GET("/setting/cron", settingHandler.GetCronSettings)
		admin.PUT("/setting/cron", settingHandler.SaveCronSettings)
		admin.GET("/setting/site", settingHandler.GetSiteSettings)
		admin.PUT("/setting/site", settingHandler.SaveSiteSettings)
		admin.GET("/setting/payment", settingHandler.GetPaymentSettings)
		admin.PUT("/setting/payment", settingHandler.SavePaymentSettings)

		// 验证码配置管理
		captchaSvc := service.NewCaptchaService(deps.Redis, deps.DB)
		captchaConfigHandler := handler.NewCaptchaConfigHandler(captchaSvc)
		admin.GET("/captcha-config", captchaConfigHandler.GetConfigs)
		admin.PUT("/captcha-config/basic", captchaConfigHandler.UpdateBasicConfig)
		admin.PUT("/captcha-config/scenes", captchaConfigHandler.UpdateSceneConfig)
		admin.POST("/captcha-config/init", captchaConfigHandler.InitDefaultConfigs)

		// 系统配置管理（统一配置）
		configHandler := handler.NewConfigHandler(deps.DB)
		admin.GET("/config/groups", configHandler.GetGroups)
		admin.GET("/config/group/:group", configHandler.GetByGroup)
		admin.GET("/config/all", configHandler.GetAll)
		admin.PUT("/config/:key", configHandler.UpdateConfig)
		admin.PUT("/config/batch", configHandler.BatchUpdateConfig)
		admin.POST("/config/init", configHandler.InitDefaultConfigs)
		admin.GET("/config/admin-path", configHandler.GetAdminPath)
		admin.PUT("/config/admin-path", configHandler.UpdateAdminPath)

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

		// 上游透传管理
		admin.GET("/ticket-deliver/upstreams", ticketDeliverHandler.GetUpstreams)
		admin.GET("/ticket-deliver/upstreams/types", ticketDeliverHandler.GetSupportedTypes)
		admin.GET("/ticket-deliver/upstreams/:id", ticketDeliverHandler.GetUpstream)
		admin.POST("/ticket-deliver/upstreams", ticketDeliverHandler.CreateUpstream)
		admin.PUT("/ticket-deliver/upstreams/:id", ticketDeliverHandler.UpdateUpstream)
		admin.DELETE("/ticket-deliver/upstreams/:id", ticketDeliverHandler.DeleteUpstream)
		admin.POST("/ticket-deliver/upstreams/:id/test", ticketDeliverHandler.TestUpstreamConnection)

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
		admin.DELETE("/system-logs/clear-level", systemLogHandler.ClearByLevel)

		// LogRecord admin methods (from zjmf)
		admin.GET("/log-record/system-log", systemLogHandler.GetSystemLog)
		admin.GET("/log-record/cron-system-log", systemLogHandler.GetCronSystemLog)
		admin.GET("/log-record/admin-log", systemLogHandler.GetAdminLog)
		admin.GET("/log-record/notify-log", systemLogHandler.GetNotifyLog)
		admin.GET("/log-record/email-log", systemLogHandler.GetEmailLog)
		admin.GET("/log-record/email-detail", systemLogHandler.GetEmailDetail)
		admin.GET("/log-record/wechat-log", systemLogHandler.GetWechatLog)
		admin.GET("/log-record/sms-log", systemLogHandler.GetSmsLog)
		admin.GET("/log-record/sms-log-m", systemLogHandler.GetSmsLogM)
		admin.GET("/log-record/system-message-log", systemLogHandler.GetSystemMessageLog)
		admin.GET("/log-record/api-log", systemLogHandler.GetApiLog)
		admin.GET("/log-record/delete-page", systemLogHandler.GetDeleteLogPage)
		admin.GET("/log-record/affirm-delete-page", systemLogHandler.GetAffirmDeleteLogPage)
		admin.DELETE("/log-record/delete-log", systemLogHandler.DeleteLog)
		admin.POST("/system-logs/clear-by-level", systemLogHandler.ClearByLevel)

		// 日志清理（增强版）
		logCleanerSvc := service.NewLogCleaner(deps.DB, deps.Log)
		logCleanerHandler := handler.NewLogCleanerHandler(logCleanerSvc, deps.Log)
		admin.GET("/log-cleaner/stats", logCleanerHandler.GetStats)
		admin.POST("/log-cleaner/clean-by-days", logCleanerHandler.CleanByDays)
		admin.POST("/log-cleaner/clean-by-count", logCleanerHandler.CleanByCount)
		admin.POST("/log-cleaner/clean-by-module", logCleanerHandler.CleanByModule)
		admin.POST("/log-cleaner/clean-by-status", logCleanerHandler.CleanByStatus)
		admin.POST("/log-cleaner/clean-expired", logCleanerHandler.CleanExpired)
		admin.POST("/log-cleaner/clean-all", logCleanerHandler.CleanAll)

		// 通知管理（去重）
		notifySvc := service.NewNotificationService(deps.DB, deps.Log, nil)
		notifyManageHandler := handler.NewNotificationManageHandler(notifySvc, deps.Log)
		admin.GET("/notifications/stats", notifyManageHandler.GetStats)
		admin.POST("/notifications/reset-event", notifyManageHandler.ResetEvent)
		admin.POST("/notifications/clean-all", notifyManageHandler.CleanAll)

		// 邮箱后缀白名单
		emailSuffixSvc := service.NewEmailSuffixWhitelistService(deps.DB, deps.Log)
		emailSuffixHandler := handler.NewEmailSuffixWhitelistHandler(emailSuffixSvc, deps.Log)
		admin.GET("/email-suffixes", emailSuffixHandler.List)
		admin.POST("/email-suffixes", emailSuffixHandler.Add)
		admin.PUT("/email-suffixes/:id", emailSuffixHandler.Update)
		admin.DELETE("/email-suffixes/:id", emailSuffixHandler.Delete)
		admin.POST("/email-suffixes/batch-delete", emailSuffixHandler.BatchDelete)
		admin.POST("/email-suffixes/import-defaults", emailSuffixHandler.ImportDefaults)

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

		// ==================== 自定义字段 ====================
		customFieldSvc := service.NewCustomFieldService(deps.DB, deps.Log)
		customFieldHandler := handler.NewCustomFieldHandler(customFieldSvc, deps.Log)
		// 字段管理
		admin.GET("/custom-fields", customFieldHandler.GetFields)
		admin.GET("/custom-fields/:id", customFieldHandler.GetFieldDetail)
		admin.POST("/custom-fields", customFieldHandler.CreateField)
		admin.PUT("/custom-fields/:id", customFieldHandler.UpdateField)
		admin.DELETE("/custom-fields/:id", customFieldHandler.DeleteField)
		admin.POST("/custom-fields/reorder", customFieldHandler.ReorderFields)
		// 字段值
		admin.GET("/custom-field-values", customFieldHandler.GetValues)
		admin.POST("/custom-field-values", customFieldHandler.SaveValues)
		admin.DELETE("/custom-field-values", customFieldHandler.DeleteValues)
		admin.GET("/custom-field-values/single", customFieldHandler.GetSingleValue)
		// 分组管理
		admin.GET("/custom-field-groups", customFieldHandler.GetGroups)
		admin.POST("/custom-field-groups", customFieldHandler.CreateGroup)
		admin.PUT("/custom-field-groups/:id", customFieldHandler.UpdateGroup)
		admin.DELETE("/custom-field-groups/:id", customFieldHandler.DeleteGroup)
		// 验证
		admin.POST("/custom-fields/validate", customFieldHandler.ValidateFields)
		// 按类型获取
		admin.GET("/custom-fields/cart", customFieldHandler.GetCartCustomFields)
		admin.GET("/custom-fields/product/:product_id", customFieldHandler.GetProductCustomFields)
		admin.GET("/custom-fields/client", customFieldHandler.GetClientCustomFields)
		admin.GET("/custom-fields/host/:host_id", customFieldHandler.GetHostCustomFields)
		// 批量操作
		admin.POST("/custom-fields/copy", customFieldHandler.CopyFields)
		admin.POST("/custom-fields/import", customFieldHandler.ImportFields)
		admin.GET("/custom-fields/export", customFieldHandler.ExportFields)

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
		admin.GET("/affiliate/:id", affiliateHandler.AdminGetByID)
		admin.PUT("/affiliate/:id", affiliateHandler.AdminUpdate)
		admin.POST("/affiliate/records/:id/confirm", affiliateHandler.AdminConfirmRecord)
		admin.POST("/affiliate/withdraws/:id/process", affiliateHandler.AdminProcessWithdraw)

		// 推广联盟扩展功能
		admin.GET("/affiliate/gateways", affiliateHandler.GatewayList)
		admin.GET("/affiliate/product/:pid", affiliateHandler.ProductAffiPage)
		admin.POST("/affiliate/product", affiliateHandler.ProductAffiPost)
		admin.GET("/affiliate/buy-records", affiliateHandler.UserAffiBuyRecord)
		admin.GET("/affiliate/user-affi-page", affiliateHandler.UserAffiPage)
		admin.POST("/affiliate/user-affi-balance", affiliateHandler.UserAffiBalance)
		admin.POST("/affiliate/user-affi-post", affiliateHandler.UserAffiPost)
		admin.GET("/affiliate/user-affi-list", affiliateHandler.UserAffiList)
		admin.GET("/affiliate/user-affi-record", affiliateHandler.UserAffiRecord)
		admin.GET("/affiliate/time-type", affiliateHandler.GetTimeType)
		admin.GET("/affiliate/get-ids", affiliateHandler.GetIDs)
		admin.GET("/affiliate/withdraw-records", affiliateHandler.AffiWithdrawRecord)
		admin.POST("/affiliate/withdraw-sh", affiliateHandler.AffiWithdrawSH)

		// ==================== 代理商 ====================
		agentHandler := handler.NewAgentHandler(deps.DB)
		admin.GET("/agents", agentHandler.AdminGetList)
		admin.POST("/agents", agentHandler.AdminCreate)
		admin.PUT("/agents/:id", agentHandler.AdminUpdate)
		admin.POST("/agents/commissions/:id/confirm", agentHandler.AdminConfirmCommission)

		// ==================== 代理商增强 ====================
		agentEnhancedSvc := service.NewAgentEnhancedService(deps.DB, deps.Log)
		agentEnhancedHandler := handler.NewAgentEnhancedHandler(agentEnhancedSvc, deps.Log)
		// 资源管理
		admin.GET("/agent/resources", agentEnhancedHandler.GetResourceInfo)
		admin.POST("/agent/resources", agentEnhancedHandler.PostResourceInfo)
		admin.GET("/agent/products", agentEnhancedHandler.GetProducts)
		admin.GET("/agent/hosts", agentEnhancedHandler.GetHostLists)
		// 巡检
		admin.GET("/agent/inspections", agentEnhancedHandler.GetInspectionLists)
		admin.POST("/agent/inspections", agentEnhancedHandler.CreateInspection)
		admin.GET("/agent/inspections/ips", agentEnhancedHandler.GetInspectionIPs)
		admin.GET("/agent/inspections/:id", agentEnhancedHandler.GetInspectionDetail)
		admin.POST("/agent/inspections/:id/upload", agentEnhancedHandler.PostUpload)
		// 订单与财务
		admin.GET("/agent/orders", agentEnhancedHandler.GetOrders)
		admin.GET("/agent/orders/search", agentEnhancedHandler.GetOrderSearchPage)
		admin.GET("/agent/renews", agentEnhancedHandler.GetRenews)
		admin.GET("/agent/renews/search", agentEnhancedHandler.GetRenewSearchPage)
		admin.GET("/agent/income", agentEnhancedHandler.GetIncome)
		admin.GET("/agent/consumption", agentEnhancedHandler.GetConsumption)
		admin.GET("/agent/logs", agentEnhancedHandler.GetAgentLogs)
		// 售后
		admin.GET("/agent/aftersale/:id", agentEnhancedHandler.GetAfterSaleDetail)
		admin.POST("/agent/aftersale/:id", agentEnhancedHandler.PostAfterSale)
		admin.POST("/agent/aftersale/:id/cancel", agentEnhancedHandler.PostUnAfterSale)
		admin.GET("/agent/refunds/:id", agentEnhancedHandler.GetRefundDetail)
		admin.POST("/agent/refunds/:id", agentEnhancedHandler.PostRefund)
		// 工单
		admin.GET("/agent/tickets", agentEnhancedHandler.GetTickets)
		// 评价
		admin.POST("/agent/evaluation", agentEnhancedHandler.PostEvaluation)
		admin.GET("/agent/run-maps", agentEnhancedHandler.GetRunMapLists)
		// 令牌管理
		admin.GET("/agent/token", agentEnhancedHandler.GetToken)
		admin.POST("/agent/token", agentEnhancedHandler.SetToken)
		admin.POST("/agent/token/check", agentEnhancedHandler.CheckToken)
		// 基本信息
		admin.GET("/agent/base-info", agentEnhancedHandler.GetBaseInfo)

		// 资源池分配与管理
		admin.GET("/agent/resource-pools", agentEnhancedHandler.GetResourcePools)
		admin.POST("/agent/resource-pools", agentEnhancedHandler.CreateResourcePool)
		admin.PUT("/agent/resource-pools/:id", agentEnhancedHandler.UpdateResourcePool)
		admin.DELETE("/agent/resource-pools/:id", agentEnhancedHandler.DeleteResourcePool)

		// 代理商等级体系
		admin.GET("/agent/levels", agentEnhancedHandler.GetAgentLevels)
		admin.POST("/agent/levels", agentEnhancedHandler.CreateAgentLevel)
		admin.PUT("/agent/levels/:id", agentEnhancedHandler.UpdateAgentLevel)
		admin.DELETE("/agent/levels/:id", agentEnhancedHandler.DeleteAgentLevel)

		// 代理商业绩报表
		admin.GET("/agent/performance", agentEnhancedHandler.GetPerformanceReport)
		admin.GET("/agent/performance/chart", agentEnhancedHandler.GetPerformanceChart)

		// 代理商佣金结算增强
		admin.GET("/agent/commission/settlements", agentEnhancedHandler.GetCommissionSettlements)
		admin.POST("/agent/commission/settle", agentEnhancedHandler.SettleCommission)
		admin.GET("/agent/commission/rules", agentEnhancedHandler.GetCommissionRules)
		admin.PUT("/agent/commission/rules", agentEnhancedHandler.UpdateCommissionRules)

		// 开发者管理
		developerSvc := service.NewDeveloperService(deps.DB, deps.Log)
		developerHandler := handler.NewDeveloperHandler(developerSvc, deps.Log)
		admin.GET("/developers", developerHandler.List)
		admin.GET("/developers/:id", developerHandler.GetDetail)
		admin.POST("/developers", developerHandler.Create)
		admin.PUT("/developers/:id", developerHandler.Update)
		admin.DELETE("/developers/:id", developerHandler.Delete)
		admin.POST("/developers/:id/approve", developerHandler.Approve)
		admin.POST("/developers/:id/reject", developerHandler.Reject)
		admin.GET("/developers/:id/api-keys", developerHandler.GetAPIKeys)
		admin.POST("/developers/:id/api-keys", developerHandler.GenerateAPIKey)
		admin.DELETE("/developers/:id/api-keys/:key_id", developerHandler.RevokeAPIKey)
		admin.GET("/developers/docs", developerHandler.GetDocs)
		admin.PUT("/developers/docs", developerHandler.UpdateDocs)
		admin.GET("/developers/billing", developerHandler.GetBilling)
		admin.POST("/developers/billing/settle", developerHandler.SettleBilling)

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
			captchaSvc := service.NewCaptchaService(deps.Redis, deps.DB)
			captchaHandler := handler.NewCaptchaHandler(captchaSvc, deps.DB)
			admin.GET("/captcha/image", captchaHandler.GetImage)
			admin.POST("/captcha/sms", captchaHandler.SendSMS)
			admin.POST("/captcha/email", captchaHandler.SendEmail)
			admin.POST("/system/captcha/test", func(c *gin.Context) {
				var req struct {
					Type  string `json:"type" binding:"required"` // sms, email, image
					Phone string `json:"phone"`
					Email string `json:"email"`
				}
				if err := c.ShouldBindJSON(&req); err != nil {
					c.JSON(400, gin.H{"error": err.Error()})
					return
				}
				switch req.Type {
				case "sms":
					if req.Phone == "" {
						c.JSON(400, gin.H{"error": "phone is required"})
						return
					}
					code, err := captchaSvc.GenerateSMS(req.Phone)
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						return
					}
					c.JSON(200, gin.H{"message": "验证码已生成", "code": code})
				case "email":
					if req.Email == "" {
						c.JSON(400, gin.H{"error": "email is required"})
						return
					}
					code, err := captchaSvc.GenerateEmail(req.Email)
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						return
					}
					c.JSON(200, gin.H{"message": "验证码已生成", "code": code})
				case "image":
					captchaID, imgBytes, err := captchaSvc.GenerateImage("test")
					if err != nil {
						c.JSON(500, gin.H{"error": err.Error()})
						return
					}
					c.JSON(200, gin.H{"captcha_id": captchaID, "image_size": len(imgBytes)})
				default:
					c.JSON(400, gin.H{"error": "invalid type"})
				}
			})
		}

		// ==================== 购物车 ====================
		promoCodeSvcForCart := service.NewPromoCodeService(deps.DB, deps.Log)
		cartSvc := service.NewCartService(deps.DB, deps.Log, deps.OrdSvc, promoCodeSvcForCart)
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
		admin.GET("/contracts/setting", contractHandler.Setting)
		admin.POST("/contracts/setting", contractHandler.SettingPost)
		admin.GET("/contracts/tpl", contractHandler.Tpl)
		admin.DELETE("/contracts/tpl/:id", contractHandler.DeleteTpl)
		admin.GET("/contracts/page", contractHandler.ContractPage)
		admin.POST("/contracts/page", contractHandler.ContractPagePost)

		// ==================== 信用额度 ====================
		creditSvc := service.NewCreditService(deps.DB, deps.Log)
		creditHandler := handler.NewCreditHandler(deps.DB, creditSvc)
		admin.GET("/credit", creditHandler.GetInfo)
		admin.GET("/credit/logs", creditHandler.AdminGetLogs)
		admin.POST("/credit/users/:id/adjust", creditHandler.AdminAdjust)
		admin.POST("/credit/bills/generate", creditHandler.AdminGenerateBills)
		admin.GET("/credit/bills", creditHandler.AdminGetBills)
		admin.POST("/credit/bills/:id/waive-fee", creditHandler.AdminWaiveLateFee)

		// 信用额度 - 用户管理
		admin.GET("/credit/clients", creditHandler.AdminGetClientList)
		admin.POST("/credit/users/enable", creditHandler.AdminEnableCredit)
		admin.POST("/credit/users/:id/disable", creditHandler.AdminDisableCredit)
		admin.PUT("/credit/users/:id/settings", creditHandler.AdminUpdateCreditSettings)

		// 信用额度 - 用户账单
		admin.GET("/credit/users/:id/invoices", creditHandler.AdminGetUserCreditInvoices)
		admin.GET("/credit/invoices/:id/items", creditHandler.AdminGetCreditInvoiceSubItems)

		// 信用额度 - 全局配置
		admin.GET("/credit/config", creditHandler.AdminGetGlobalCreditConfig)
		admin.PUT("/credit/config", creditHandler.AdminUpdateGlobalCreditConfig)

		// 信用额度扩展功能
		admin.DELETE("/credit/users/:uid", creditHandler.Delete)
		admin.GET("/credit/search", creditHandler.GetSearch)
		admin.GET("/credit/index", creditHandler.AdminIndex)
		admin.GET("/credit/log", creditHandler.AdminCreditLog)
		admin.GET("/credit/invoices", creditHandler.AdminCreditInvoiceList)

		// ==================== 货币 ====================
		currencyHandler := handler.NewCurrencyHandler(deps.DB)
		admin.GET("/currencies", currencyHandler.AdminGetList)
		admin.GET("/currencies/all", currencyHandler.GetAll)
		admin.POST("/currencies", currencyHandler.AdminCreate)
		admin.PUT("/currencies/:id", currencyHandler.AdminUpdate)
		admin.DELETE("/currencies/:id", currencyHandler.AdminDelete)
		admin.PUT("/currencies/:id/rate", currencyHandler.AdminUpdateRate)
		admin.PUT("/currencies/:id/default", currencyHandler.AdminSetDefault)
		admin.POST("/currencies/update-prices", currencyHandler.AdminUpdateAllPrices)

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

		// 下载管理扩展功能
		admin.GET("/downloads/config", downloadHandler.DownloadsConfig)
		admin.POST("/downloads/sort", downloadHandler.UpdateFileSort)

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

		// News admin methods (from zjmf)
		admin.GET("/news/cats-page", newsHandler.GetCatsPage)
		admin.GET("/news/cate-list", newsHandler.GetCateList)
		admin.GET("/news/cat-data", newsHandler.GetCatData)
		admin.POST("/news/edit-cat", newsHandler.PostEditCat)
		admin.GET("/news/check-alias", newsHandler.GetCheckalias)
		admin.DELETE("/news/delete-cat", newsHandler.DeleteCat)
		admin.GET("/news/content", newsHandler.GetContent)
		admin.POST("/news/edit-content", newsHandler.PostEditContent)
		admin.DELETE("/news/delete-content", newsHandler.DeleteContent)
		admin.GET("/news/custom-param", newsHandler.GetGetCustomParam)
		admin.POST("/news/add-custom-param", newsHandler.GetAddCustomParam)
		admin.POST("/news/update-custom-param", newsHandler.GetUpdateCustomParam)
		admin.DELETE("/news/del-custom-param", newsHandler.GetDelCustomParam)
		admin.GET("/news/custom-update-val", newsHandler.GetGetCustomUpdateVal)

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

		// 客户通知（admin 查看指定客户的通知）
		admin.GET("/clients/:id/notifications", notificationHandler.AdminGetClientNotifications)

		// ==================== OAuth登录 ====================
		oauthSvc := service.NewOAuthService(deps.DB, deps.Log, deps.UserSvc, deps.FrontendURL)
		oauthHandler := handler.NewOAuthHandler(oauthSvc, deps.Log, deps.JWTMgr)
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
		admin.GET("/product-diverts/:id/code", productDivertHandler.GetTransferCode)
		admin.POST("/product-diverts/:id/regenerate-code", productDivertHandler.RegenerateCode)
		admin.POST("/product-diverts/accept-by-code", productDivertHandler.AcceptByCode)

		// 管理员产品转移配置
		admin.GET("/product-transfer/config", productDivertHandler.AdminGetConfig)
		admin.PUT("/product-transfer/config", productDivertHandler.AdminSaveConfig)
		admin.GET("/product-transfers", productDivertHandler.AdminGetAllTransfers)

		// ==================== RBAC权限 ====================
		rbacSvc := service.NewRbacService(deps.DB, deps.Log)
		rbacHandler := handler.NewRbacHandler(rbacSvc, deps.Log)
		admin.GET("/rbac/roles", rbacHandler.GetRoles)
		admin.POST("/rbac/roles", rbacHandler.CreateRole)
		admin.PUT("/rbac/roles/:id", rbacHandler.UpdateRole)
		admin.DELETE("/rbac/roles/:id", rbacHandler.DeleteRole)
		admin.GET("/rbac/permissions", rbacHandler.GetPermissions)
		admin.PUT("/rbac/permissions", rbacHandler.UpdatePermissions)
		admin.POST("/rbac/users/:id/roles", rbacHandler.AssignRole)
		admin.GET("/rbac/users/:id/roles", rbacHandler.GetUserRoles)

		// RBAC admin methods (from zjmf)
		admin.GET("/rbac/index", rbacHandler.Index)
		admin.GET("/rbac/add-role-page", rbacHandler.AddRolePage)
		admin.POST("/rbac/add-role", rbacHandler.AddRole)
		admin.GET("/rbac/edit-role-page", rbacHandler.EditRolePage)
		admin.POST("/rbac/edit-role", rbacHandler.EditRole)
		admin.DELETE("/rbac/delete/:id", rbacHandler.Delete)
		admin.POST("/rbac/copy-role", rbacHandler.CopyRole)

		// ==================== 报表 ====================
		reportSvc := service.NewReportService(deps.DB, deps.Log)
		reportHandler := handler.NewReportHandler(deps.DB, reportSvc, deps.Log)
		admin.GET("/reports/dashboard", reportHandler.GetDashboard)
		admin.GET("/reports/daily", reportHandler.GetDailyReport)
		admin.GET("/reports/monthly", reportHandler.GetMonthlyReport)
		admin.GET("/reports/revenue", reportHandler.GetRevenueChart)
		admin.GET("/reports/users", reportHandler.GetUserStats)
		admin.GET("/reports/orders", reportHandler.GetOrderStats)
		admin.GET("/reports/top-clients", reportHandler.GetTopClients)
		admin.GET("/reports/product-income", reportHandler.GetProductIncome)

		// 报表扩展功能
		admin.GET("/reports/modules", reportHandler.GetSystemInfoModulesList)
		admin.POST("/reports/modules/sort", reportHandler.UpdateSystemInfoModulesSort)
		admin.GET("/reports/year-income-statistics", reportHandler.GetYearIncomeStatistics)
		admin.GET("/reports/year-income-statistics-chart", reportHandler.GetYearIncomeStatisticsForChart)
		admin.GET("/reports/new-client-statistics", reportHandler.GetNewClientStatistics)
		admin.GET("/reports/revenue-ranking", reportHandler.GetRevenueRanking)
		admin.GET("/reports/product-income/trend", reportHandler.GetProductIncomeTrend)
		admin.GET("/reports/product-income/comparison", reportHandler.GetProductIncomeComparison)
		admin.GET("/reports/revenue-ranking/comparison", reportHandler.GetRevenueRankingComparison)

		// ==================== 黑名单管理 ====================
		blacklistSvc := service.NewBlacklistService(deps.DB, deps.Log)
		blacklistHandler := handler.NewBlacklistHandler(blacklistSvc, deps.Log)
		admin.GET("/blacklist", blacklistHandler.List)
		admin.POST("/blacklist", blacklistHandler.Create)
		admin.PUT("/blacklist/:id", blacklistHandler.Update)
		admin.DELETE("/blacklist/:id", blacklistHandler.Delete)

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

		// 上游产品对接（分组+多线程）
		admin.GET("/upstream/providers/:id/products", upstreamHandler.GetUpstreamProducts)
		admin.GET("/upstream/providers/:id/groups", upstreamHandler.GetUpstreamGroups)
		admin.GET("/upstream/local-groups", upstreamHandler.GetLocalGroups)
		admin.POST("/upstream/local-groups", upstreamHandler.CreateLocalGroup)
		admin.POST("/upstream/dock/products", upstreamHandler.DockProducts)
		admin.POST("/upstream/dock/group", upstreamHandler.DockGroup)
		admin.POST("/upstream/products/:product_id/sync", upstreamHandler.SyncSingleProduct)
		admin.POST("/upstream/providers/:id/sync-all", upstreamHandler.SyncAllProducts)
		// P3-19: 清空对接
		admin.POST("/upstream/providers/:id/empty", upstreamHandler.EmptyUpper)

		// ==================== 上游操作 ====================
		upstreamOpsSvc := service.NewUpstreamService(deps.DB, deps.Log)
		upstreamOpsHandler := handler.NewUpstreamOpsHandler(upstreamOpsSvc, deps.Log)

		// 电源操作
		admin.POST("/upstream/:provider_id/host/:host_id/boot", upstreamOpsHandler.Boot)
		admin.POST("/upstream/:provider_id/host/:host_id/shutdown", upstreamOpsHandler.Shutdown)
		admin.POST("/upstream/:provider_id/host/:host_id/reboot", upstreamOpsHandler.Reboot)
		admin.GET("/upstream/:provider_id/host/:host_id/status", upstreamOpsHandler.GetStatus)

		// 控制台
		admin.POST("/upstream/:provider_id/host/:host_id/vnc", upstreamOpsHandler.VNC)
		admin.POST("/upstream/:provider_id/host/:host_id/kvm", upstreamOpsHandler.KVM)
		admin.GET("/upstream/:provider_id/host/:host_id/ipmi/status", upstreamOpsHandler.IPMIStatus)
		admin.POST("/upstream/:provider_id/host/:host_id/ipmi/on", upstreamOpsHandler.IPMIOn)
		admin.POST("/upstream/:provider_id/host/:host_id/ipmi/off", upstreamOpsHandler.IPMIOff)
		admin.POST("/upstream/:provider_id/host/:host_id/ipmi/reboot", upstreamOpsHandler.IPMIReboot)
		admin.POST("/upstream/:provider_id/host/:host_id/ipmi/vnc", upstreamOpsHandler.IPMIVNC)

		// 重装
		admin.POST("/upstream/:provider_id/host/:host_id/reinstall", upstreamOpsHandler.Reinstall)
		admin.GET("/upstream/:provider_id/host/:host_id/reinstall/status", upstreamOpsHandler.GetReinstallStatus)
		admin.POST("/upstream/:provider_id/host/:host_id/reinstall/cancel", upstreamOpsHandler.CancelReinstall)
		admin.GET("/upstream/:provider_id/os-list", upstreamOpsHandler.GetOSList)
		admin.POST("/upstream/:provider_id/host/:host_id/crack-password", upstreamOpsHandler.CrackPassword)

		// DCIM客户端操作
		admin.GET("/upstream/:provider_id/host/:host_id/dcim/status", upstreamOpsHandler.DcimClientStatus)
		admin.POST("/upstream/:provider_id/host/:host_id/dcim/on", upstreamOpsHandler.DcimClientOn)
		admin.POST("/upstream/:provider_id/host/:host_id/dcim/off", upstreamOpsHandler.DcimClientOff)
		admin.POST("/upstream/:provider_id/host/:host_id/dcim/reboot", upstreamOpsHandler.DcimClientReboot)
		admin.POST("/upstream/:provider_id/host/:host_id/dcim/vnc", upstreamOpsHandler.DcimClientVNC)
		admin.POST("/upstream/:provider_id/host/:host_id/dcim/reinstall", upstreamOpsHandler.DcimClientReinstall)
		admin.POST("/upstream/:provider_id/host/:host_id/dcim/crack-pass", upstreamOpsHandler.DcimClientCrackPass)
		admin.POST("/upstream/:provider_id/host/:host_id/dcim/reinstall/cancel", upstreamOpsHandler.DcimClientCancelReinstall)
		admin.GET("/upstream/:provider_id/host/:host_id/dcim/reinstall/status", upstreamOpsHandler.DcimClientReinstallStatus)
		admin.GET("/upstream/:provider_id/dcim/os-list", upstreamOpsHandler.DcimClientGetOS)

		// 模块按钮
		admin.POST("/upstream/:provider_id/host/:host_id/module/client-button", upstreamOpsHandler.ModuleClientButton)
		admin.POST("/upstream/:provider_id/host/:host_id/module/admin-button", upstreamOpsHandler.ModuleAdminButton)
		admin.GET("/upstream/:provider_id/host/:host_id/module/power-status", upstreamOpsHandler.ModulePowerStatus)

		// ==================== 用户等级 ====================
		userLevelHandler := handler.NewUserLevelHandler(deps.DB)
		admin.GET("/user-levels", userLevelHandler.AdminGetList)
		admin.POST("/user-levels", userLevelHandler.AdminCreate)
		admin.PUT("/user-levels/:id", userLevelHandler.AdminUpdate)
		admin.DELETE("/user-levels/:id", userLevelHandler.AdminDelete)

		// ==================== v10购物车 ====================
		promoCodeSvcForV10 := service.NewPromoCodeService(deps.DB, deps.Log)
		v10CartSvc := service.NewV10CartService(deps.DB, deps.Log, deps.OrdSvc, promoCodeSvcForV10)
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

		// ==================== V10云管理 ====================
		promoCodeSvcForV10Cloud := service.NewPromoCodeService(deps.DB, deps.Log)
		v10CloudSvc := service.NewV10CloudService(deps.DB, deps.Log, deps.OrdSvc, promoCodeSvcForV10Cloud)
		v10CloudHandler := handler.NewV10CloudHandler(v10CloudSvc, deps.Log)

		// 产品浏览
		admin.GET("/v10/cloud/products", v10CloudHandler.GetProductList)
		admin.GET("/v10/cloud/products/:id", v10CloudHandler.GetProductDetail)
		admin.GET("/v10/cloud/regions", v10CloudHandler.GetRegions)
		admin.GET("/v10/cloud/os-types", v10CloudHandler.GetOSTypes)

		// 配置选项
		admin.GET("/v10/cloud/products/:id/config", v10CloudHandler.GetConfigOptions)
		admin.POST("/v10/cloud/calculate-price", v10CloudHandler.CalculatePrice)
		admin.GET("/v10/cloud/products/:id/linkage", v10CloudHandler.GetLinkAgeList)
		admin.GET("/v10/cloud/products/:id/config/filter", v10CloudHandler.FilterConfigOptions)

		// 购物车操作
		admin.GET("/v10/cloud/cart", v10CloudHandler.GetCartSummary)
		admin.GET("/v10/cloud/cart/items", v10CloudHandler.GetCartItems)
		admin.POST("/v10/cloud/cart", v10CloudHandler.AddToCart)
		admin.PUT("/v10/cloud/cart/:id", v10CloudHandler.UpdateCartItem)
		admin.POST("/v10/cloud/cart/settle", v10CloudHandler.SettleCart)

		// 订单流程
		admin.POST("/v10/cloud/orders", v10CloudHandler.CreateOrder)
		admin.GET("/v10/cloud/orders/:id", v10CloudHandler.GetOrderDetail)
		admin.POST("/v10/cloud/orders/:id/pay", v10CloudHandler.PayOrder)

		// 主机管理
		admin.GET("/v10/cloud/hosts/:id", v10CloudHandler.GetHostInfo)
		admin.GET("/v10/cloud/hosts/:id/config", v10CloudHandler.GetHostConfig)
		admin.GET("/v10/cloud/hosts/:id/traffic", v10CloudHandler.GetTrafficUsage)
		admin.GET("/v10/cloud/hosts/:id/os", v10CloudHandler.GetOSList)

		// ==================== 发票管理 ====================
		voucherSvc := service.NewVoucherService(deps.DB, deps.Log)
		voucherHandler := handler.NewVoucherHandler(voucherSvc, deps.Log)
		admin.GET("/voucher-rate", voucherHandler.GetRate)
		admin.POST("/voucher-rate", voucherHandler.PostRate)
		admin.GET("/voucher-list", voucherHandler.GetVoucherList)
		admin.GET("/voucher-detail/:id", voucherHandler.GetVoucherDetail)
		admin.POST("/voucher-status", voucherHandler.PostVoucherStatus)

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
		admin.GET("/zjmf-api/summary", zjmfHandler.GetSummary)
		admin.GET("/zjmf-api/downstream-summary", zjmfHandler.GetDownstreamSummary)
		admin.GET("/zjmf-api/:id", zjmfHandler.GetApi)
		admin.POST("/zjmf-api", zjmfHandler.CreateApi)
		admin.PUT("/zjmf-api/:id", zjmfHandler.UpdateApi)
		admin.DELETE("/zjmf-api/:id", zjmfHandler.DeleteApi)
		admin.POST("/zjmf-api/:id/test", zjmfHandler.TestConnection)
		admin.POST("/zjmf-api/:id/sync", zjmfHandler.SyncProducts)
		admin.POST("/zjmf-api/:id/refresh", zjmfHandler.RefreshStatus)
		admin.POST("/zjmf-api/:id/toggle", zjmfHandler.ToggleApi)
		admin.GET("/zjmf-api/:id/products", zjmfHandler.GetApiProducts)
		admin.GET("/zjmf-api/:id/orders", zjmfHandler.GetApiOrders)
		admin.GET("/zjmf-api/:id/hosts", zjmfHandler.GetApiHosts)
		admin.GET("/zjmf-api/:id/upstream-hosts", zjmfHandler.GetUpstreamHosts)
		admin.GET("/zjmf-api/:id/logs", zjmfHandler.GetApiLogs)
		admin.POST("/zjmf-api/:id/import", zjmfHandler.ImportProducts)
		admin.GET("/zjmf-api/:id/manual-hosts", zjmfHandler.GetManualHosts)
		admin.POST("/zjmf-api/:id/manual-hosts", zjmfHandler.PostManualHost)
		admin.GET("/zjmf-api/:id/credit", zjmfHandler.GetUpstreamCredit)
		admin.GET("/zjmf-api/:id/renew", zjmfHandler.GetRenewInfo)

		// ==================== SSL证书管理 ====================
		sslCertSvc := service.NewSSLCertificateService(deps.DB, deps.Log)
		sslCertHandler := handler.NewSSLCertificateHandler(sslCertSvc, deps.Log)
		admin.GET("/ssl-certificates", sslCertHandler.AdminGetList)
		admin.GET("/ssl-certificates/:id", sslCertHandler.AdminGetDetail)
		admin.PUT("/ssl-certificates/:id", sslCertHandler.AdminUpdate)
		admin.DELETE("/ssl-certificates/:id", sslCertHandler.AdminDelete)

		// ==================== PDF生成 ====================
		pdfSvc := service.NewPDFService(deps.DB)
		pdfHandler := handler.NewPDFHandler(pdfSvc)
		admin.GET("/contracts/:id/pdf", pdfHandler.GenerateContractPDF)
		admin.GET("/invoices/:id/pdf", pdfHandler.GenerateInvoicePDF)
		admin.GET("/contracts/:id/pdf/seal", pdfHandler.GenerateContractWithSeal)

		// ==================== 邮件增强 ====================
		emailEnhancedSvc := service.NewEmailEnhancedService(deps.DB, deps.Log)
		emailEnhancedHandler := handler.NewEmailEnhancedHandler(emailEnhancedSvc)
		admin.POST("/email/test", emailEnhancedHandler.SendTestEmail)
		admin.POST("/email/batch", emailEnhancedHandler.SendBatchEmail)
		admin.GET("/email-logs", emailEnhancedHandler.GetEmailLogs)
		admin.GET("/email-stats", emailEnhancedHandler.GetEmailStats)
		admin.GET("/email-support", emailEnhancedHandler.GetSupportConfig)
		admin.PUT("/email-support", emailEnhancedHandler.UpdateSupportConfig)

		// ==================== 升级增强 ====================
		upgradeEnhancedSvc := service.NewUpgradeEnhancedService(deps.DB)
		upgradeEnhancedHandler := handler.NewUpgradeEnhancedHandler(upgradeEnhancedSvc)
		admin.GET("/upgrade-config/:product_id", upgradeEnhancedHandler.GetUpgradeConfig)
		admin.POST("/upgrade-config/:host_id", upgradeEnhancedHandler.UpgradeConfigAdmin)
		admin.GET("/upgrade/check-change/:host_id", upgradeEnhancedHandler.CheckChange)
		admin.GET("/upgrade/filter-options/:product_id", upgradeEnhancedHandler.FilterConfigOptions)
		admin.GET("/upgrade/promo/:host_id", upgradeEnhancedHandler.CheckUpgradePromo)

		// ==================== 续费增强 ====================
		renewEnhancedSvc := service.NewRenewEnhancedService(deps.DB)
		renewEnhancedHandler := handler.NewRenewEnhancedHandler(renewEnhancedSvc)
		admin.GET("/renewal/page/:host_id", renewEnhancedHandler.GetRenewalPage)
		admin.GET("/renewal/price/:host_id", renewEnhancedHandler.GetRenewalPrice)
		admin.POST("/renewal/submit", renewEnhancedHandler.SubmitRenewal)
		admin.POST("/renewal/batch", renewEnhancedHandler.BatchRenew)
		admin.PUT("/renewal/auto-renew/:host_id", renewEnhancedHandler.SetAutoRenew)
		admin.PUT("/renewal/pay-type/:host_id", renewEnhancedHandler.SetPayType)
		admin.GET("/renewal/pay-type/:host_id", renewEnhancedHandler.GetPayType)
		admin.DELETE("/renewal/invoice/:id", renewEnhancedHandler.DeleteRenewInvoice)

		// ==================== 菜单管理增强 ====================
		menuEnhancedSvc := service.NewMenuEnhancedService(deps.DB)
		menuEnhancedHandler := handler.NewMenuEnhancedHandler(menuEnhancedSvc)
		admin.GET("/web-navs", menuEnhancedHandler.GetWebNavs)
		admin.POST("/web-navs", menuEnhancedHandler.CreateWebNav)
		admin.PUT("/web-navs/:id", menuEnhancedHandler.UpdateWebNav)
		admin.DELETE("/web-navs/:id", menuEnhancedHandler.DeleteWebNav)
		admin.GET("/web-navs/default", menuEnhancedHandler.GetDefaultSenior)
		admin.GET("/web-navs/list", menuEnhancedHandler.GetWebNavList)
		admin.POST("/web-navs/list", menuEnhancedHandler.SetWebNavList)
		admin.GET("/web-navs/:id", menuEnhancedHandler.GetOneNavs)
		admin.POST("/custom-pages", menuEnhancedHandler.AddCustomPage)
		admin.POST("/product-pages", menuEnhancedHandler.AddProductPage)
		admin.GET("/system-nav", menuEnhancedHandler.GetSystemNav)
		admin.GET("/product-menus", menuEnhancedHandler.GetProductList)
		admin.POST("/web-pages", menuEnhancedHandler.CreateWebPage)
		admin.GET("/menu-types", menuEnhancedHandler.GetMenuType)
		admin.GET("/other-menus", menuEnhancedHandler.GetOtherMenu)
		admin.DELETE("/two-menus/:id", menuEnhancedHandler.DelTwoMenu)
		admin.GET("/all-menus", menuEnhancedHandler.GetTypeAllMenu)
		admin.PUT("/menu-active/:id", menuEnhancedHandler.EditMenuActive)
		admin.GET("/nav-types", menuEnhancedHandler.GetNavType)
		admin.GET("/create-web-data", menuEnhancedHandler.GetCreateWebData)
		admin.DELETE("/direct/:id", menuEnhancedHandler.DirectDel)
		admin.POST("/hook-menus", menuEnhancedHandler.AddHookMenu)
		admin.DELETE("/hook-menus/:id", menuEnhancedHandler.DelHookMenu)
		admin.GET("/hook-menus", menuEnhancedHandler.GetHookMenus)
		admin.POST("/nav-links/:nav_id", menuEnhancedHandler.SaveLinks)
		admin.DELETE("/nav-links/:nav_id", menuEnhancedHandler.DeleteLinks)
		admin.GET("/nav-links", menuEnhancedHandler.AllLinks)

		// ==================== 客户跟踪管理 ====================
		clientTrackHandler := handler.NewClientTrackHandler(deps.DB)
		admin.GET("/client-tracks", clientTrackHandler.List)
		admin.GET("/client-tracks/:id", clientTrackHandler.GetDetail)
		admin.POST("/client-tracks", clientTrackHandler.Create)
		admin.PUT("/client-tracks/:id", clientTrackHandler.Update)
		admin.DELETE("/client-tracks/:id", clientTrackHandler.Delete)
		admin.POST("/client-tracks/:id/remarks", clientTrackHandler.AddRemark)
		admin.DELETE("/client-tracks/remarks/:id", clientTrackHandler.DeleteRemark)
		admin.PUT("/client-tracks/status", clientTrackHandler.UpdateTrackStatus)
		admin.GET("/client-tracks/notes/:uid", clientTrackHandler.GetClientNotes)
		admin.PUT("/client-tracks/notes/:uid", clientTrackHandler.UpdateClientNotes)

		// ==================== 快递管理 ====================
		expressHandler := handler.NewExpressHandler(deps.DB)
		admin.GET("/expresses", expressHandler.List)
		admin.POST("/expresses", expressHandler.Create)
		admin.PUT("/expresses/:id", expressHandler.Update)
		admin.DELETE("/expresses/:id", expressHandler.Delete)

		// ==================== 取消原因管理 ====================
		cancelReasonHandler := handler.NewCancelReasonHandler(deps.DB)
		admin.GET("/cancel-reasons", cancelReasonHandler.List)
		admin.POST("/cancel-reasons", cancelReasonHandler.Create)
		admin.PUT("/cancel-reasons/:id", cancelReasonHandler.Update)
		admin.DELETE("/cancel-reasons/:id", cancelReasonHandler.Delete)

		// ==================== 首页基本信息管理 ====================
		baseInfoHandler := handler.NewBaseInfoHandler(deps.DB)
		admin.GET("/base-infos", baseInfoHandler.List)
		admin.POST("/base-infos", baseInfoHandler.Create)
		admin.PUT("/base-infos/:id", baseInfoHandler.Update)
		admin.DELETE("/base-infos/:id", baseInfoHandler.Delete)

		// ==================== 用户专属下载管理 ====================
		userDownloadHandler := handler.NewUserDownloadHandler(deps.DB)
		admin.GET("/user-downloads", userDownloadHandler.AdminList)
		admin.POST("/user-downloads", userDownloadHandler.AdminCreate)
		admin.DELETE("/user-downloads/:id", userDownloadHandler.AdminDelete)

		// ==================== 审计日志 ====================
		auditLogSvc := model.NewAuditLogService(deps.DB)
		auditLogHandler := handler.NewAuditLogHandler(auditLogSvc, deps.Log)
		admin.GET("/audit-logs", auditLogHandler.List)
		admin.GET("/audit-logs/:id", auditLogHandler.Get)
		admin.GET("/audit-logs/stats", auditLogHandler.Stats)
		admin.POST("/audit-logs/clean", auditLogHandler.CleanOldLogs)

		// ==================== 数据库备份 ====================
		backupSvc := backup.NewService("./backups", deps.Log)
		backupHandler := handler.NewBackupHandler(backupSvc, deps.Log)
		admin.GET("/backups", backupHandler.ListBackups)
		admin.POST("/backups", backupHandler.CreateBackup)
		admin.DELETE("/backups/:filename", backupHandler.DeleteBackup)
		admin.POST("/backups/clean", backupHandler.CleanOldBackups)
		admin.POST("/backups/restore", backupHandler.RestoreBackup)

		// ==================== 客服聊天系统（anchor_cloud_finance_pro） ====================
		csSvc := service.NewCSChatService(deps.DB, deps.Log)
		csHandler := handler.NewCSChatHandler(csSvc, deps.Log)
		chat := admin.Group("/cs")
		{
			chat.GET("/ai-config", csHandler.GetAIConfig)
			chat.PUT("/ai-config", csHandler.SaveAIConfig)
			chat.GET("/appearance", csHandler.GetAppearanceConfig)
			chat.PUT("/appearance", csHandler.SaveAppearanceConfig)
			chat.GET("/working-hours", csHandler.GetWorkingHours)
			chat.PUT("/working-hours", csHandler.SaveWorkingHours)
			chat.GET("/sessions", csHandler.ListSessions)
			chat.GET("/sessions/:id", csHandler.GetSession)
			chat.POST("/sessions/:id/reply", csHandler.SendReply)
			chat.POST("/sessions/:id/transfer", csHandler.TransferToHuman)
			chat.POST("/sessions/:id/close", csHandler.CloseSession)
			chat.POST("/sessions/:id/rate", csHandler.RateSession)
			chat.GET("/stats", csHandler.GetStats)
		}

		// ==================== AI 工单核心（mianyu_ai_ticket） ====================
		aiTicketSvc := service.NewAITicketCoreService(deps.DB, deps.Log)
		aiTicketHandler := handler.NewAITicketCoreHandler(aiTicketSvc, deps.Log)
		ticket := admin.Group("/ai-ticket")
		{
			ticket.GET("/dashboard", aiTicketHandler.GetDashboard)
			ticket.PUT("/dashboard", aiTicketHandler.SaveDashboard)
			ticket.GET("/knowledge", aiTicketHandler.ListKnowledge)
			ticket.POST("/knowledge", aiTicketHandler.CreateKnowledge)
			ticket.PUT("/knowledge/:id", aiTicketHandler.UpdateKnowledge)
			ticket.DELETE("/knowledge/:id", aiTicketHandler.DeleteKnowledge)
			ticket.POST("/knowledge/import", aiTicketHandler.ImportDefaultKnowledge)
			ticket.GET("/rules", aiTicketHandler.ListRules)
			ticket.POST("/rules", aiTicketHandler.CreateRule)
			ticket.PUT("/rules/:id", aiTicketHandler.UpdateRule)
			ticket.DELETE("/rules/:id", aiTicketHandler.DeleteRule)
			ticket.GET("/queue", aiTicketHandler.ListQueue)
			ticket.GET("/queue/stats", aiTicketHandler.GetQueueStats)
			ticket.GET("/process-logs", aiTicketHandler.ListProcessLogs)
			ticket.GET("/notify-logs", aiTicketHandler.ListNotifyLogs)
			ticket.GET("/mode/:ticket_id", aiTicketHandler.GetTicketMode)
			ticket.PUT("/mode/:ticket_id", aiTicketHandler.SetTicketMode)
			ticket.POST("/test", aiTicketHandler.TestAutoReply)

			// Agent（Function Calling）工具管理
			ticket.GET("/tools", aiTicketHandler.ListTools)
			ticket.PUT("/tools/:name", aiTicketHandler.SetToolEnabled)
			ticket.GET("/tools/execution-logs", aiTicketHandler.ListToolExecutionLogs)
		}

		// ==================== AI 购物助手（mahiru_ai_shopping） ====================
		aiShopSvc := service.NewAIShoppingCoreService(deps.DB, deps.Log)
		aiShopHandler := handler.NewAIShoppingCoreHandler(aiShopSvc, deps.Log)
		shopping := admin.Group("/ai-shopping")
		{
			shopping.GET("/config", aiShopHandler.GetConfig)
			shopping.PUT("/config", aiShopHandler.SaveConfig)
		}

		// ==================== ACFP 模块（anchor_cloud_finance_pro） ====================
		acfpModulesSvc := service.NewACFPService(deps.DB, deps.Log)
		acfpModulesHandler := handler.NewACFPHandler(acfpModulesSvc, deps.Log)
		acfp := admin.Group("/acfp")
		{
			// 通用模块配置
			acfp.GET("/module/:key", acfpModulesHandler.GetModuleConfig)
			acfp.POST("/module/:key/toggle", acfpModulesHandler.ToggleModule)

			// 失败通知
			acfp.GET("/fail-notify/config", acfpModulesHandler.GetFailNotifyConfig)
			acfp.PUT("/fail-notify/config", acfpModulesHandler.SaveFailNotifyConfig)

			// 状态对账
			acfp.GET("/status-sync/config", acfpModulesHandler.GetStatusSyncConfig)
			acfp.PUT("/status-sync/config", acfpModulesHandler.SaveStatusSyncConfig)
			acfp.GET("/status-sync/cache/:host_id", acfpModulesHandler.GetUpstreamCache)
			acfp.GET("/status-sync/cron-statuses", acfpModulesHandler.GetCronStatuses)

			// IP 记录
			acfp.GET("/ip-history", acfpModulesHandler.GetIPHistory)

			// 限量发售
			acfp.GET("/limited-sale", acfpModulesHandler.ListLimitedSales)
			acfp.POST("/limited-sale", acfpModulesHandler.SetLimitedSale)
			acfp.PUT("/limited-sale/:id", acfpModulesHandler.UpdateLimitedSale)
			acfp.DELETE("/limited-sale/:id", acfpModulesHandler.DeleteLimitedSale)
			acfp.POST("/limited-sale/:id/reset-quota", acfpModulesHandler.ResetLimitedSaleQuota)

			// 价格锁定
			acfp.GET("/price-lock", acfpModulesHandler.ListPriceLocks)
			acfp.POST("/price-lock", acfpModulesHandler.SetPriceLock)
			acfp.DELETE("/price-lock/:id", acfpModulesHandler.DeletePriceLock)

			// 批量商品修改
			acfp.POST("/batch-product", acfpModulesHandler.BatchUpdateProducts)

			// 实名认证 Pro
			acfp.GET("/cert-pro/config", acfpModulesHandler.GetCertProConfig)
			acfp.PUT("/cert-pro/config", acfpModulesHandler.SaveCertProConfig)
			acfp.GET("/cert-pro/reviews", acfpModulesHandler.GetCertReviewList)
			acfp.POST("/cert-pro/reviews/:id/review", acfpModulesHandler.ReviewCert)
			acfp.GET("/cert-pro/scan-minors", acfpModulesHandler.ScanMinorCerts)
			acfp.POST("/cert-pro/reject-minors", acfpModulesHandler.RejectUnderageSubmissions)

			// 缓存预热
			acfp.GET("/cache-warm/status", acfpModulesHandler.GetCacheWarmStatus)
			acfp.POST("/cache-warm/trigger", acfpModulesHandler.WarmCache)

			// 系统日志
			acfp.GET("/logs", acfpModulesHandler.ListLogs)
			acfp.POST("/logs/clean", acfpModulesHandler.CleanLogs)

		// 业务列表 Pro
		acfp.GET("/business-list", acfpModulesHandler.GetBusinessList)
		acfp.GET("/business-list/:host_id", acfpModulesHandler.GetBusinessRow)
		acfp.GET("/business-list/:host_id/snapshot", acfpModulesHandler.GetBusinessSnapshot)
		acfp.POST("/business-list/:host_id/sync", acfpModulesHandler.SyncOneBusiness)
		acfp.POST("/business-list/:host_id/suspend", acfpModulesHandler.SuspendOneBusiness)
		acfp.POST("/business-list/:host_id/unsuspend", acfpModulesHandler.UnsuspendOneBusiness)
		acfp.POST("/business-list/:host_id/delete", acfpModulesHandler.DeleteOneBusiness)
		acfp.POST("/business-list/:host_id/provision", acfpModulesHandler.ProvisionOneBusiness)
		acfp.GET("/business-list/stats", acfpModulesHandler.GetBusinessStats)

		// ==================== 交易市场 ====================
		marketplaceSvc := service.NewMarketplaceService(deps.DB, deps.Log)
		marketplaceHandler := handler.NewMarketplaceHandler(marketplaceSvc, deps.Log)
		marketplace := admin.Group("/marketplace")
		{
			// 配置
			marketplace.GET("/config", marketplaceHandler.AdminGetConfig)
			marketplace.PUT("/config", marketplaceHandler.AdminSaveConfig)

			// 挂售管理
			marketplace.GET("/listings", marketplaceHandler.AdminGetAllListings)

			// 订单管理
			marketplace.GET("/orders", marketplaceHandler.AdminGetAllOrders)
		}

		// ==================== 公共接口 ====================
		commonHandler := handler.NewCommonHandler(deps.DB, deps.Log)
		{
			admin.GET("/common", commonHandler.Common)
			admin.GET("/info-notice", commonHandler.InfoNotice)
			admin.GET("/gateways", commonHandler.GetGetways)
			admin.GET("/email-templates", commonHandler.GetEmailTem)
			admin.GET("/sms-countries", commonHandler.GetSmsCountry)
		}
	}
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

	// 客服聊天（前台）
	csSvc := service.NewCSChatService(db, log)
	csHandler := handler.NewCSChatHandler(csSvc, log)
	r.POST("/cs/send", func(c *gin.Context) {
		var req struct {
			Content   string `json:"content" binding:"required"`
			VisitorID string `json:"visitor_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "消息不能为空"})
			return
		}
		session, err := csSvc.GetOrCreateSession(req.VisitorID)
		if err != nil {
			c.JSON(500, gin.H{"error": "会话创建失败"})
			return
		}
		if err := csSvc.SendMessage(session.ID, "user", req.Content); err != nil {
			c.JSON(500, gin.H{"error": "消息发送失败"})
			return
		}
		reply, err := csSvc.AIReply(session.ID, req.Content)
		if err != nil {
			c.JSON(500, gin.H{"error": "AI回复失败"})
			return
		}
		c.JSON(200, gin.H{"code": 0, "data": gin.H{"reply": reply, "session_id": session.ID}})
	})
	r.GET("/cs/history", func(c *gin.Context) {
		visitorID := c.Query("visitor_id")
		if visitorID == "" {
			c.JSON(400, gin.H{"error": "缺少visitor_id"})
			return
		}
		session, err := csSvc.GetOrCreateSession(visitorID)
		if err != nil {
			c.JSON(500, gin.H{"error": "会话创建失败"})
			return
		}
		messages, err := csSvc.GetSessionMessages(session.ID)
		if err != nil {
			c.JSON(500, gin.H{"error": "获取消息失败"})
			return
		}
		c.JSON(200, gin.H{"code": 0, "data": messages})
	})

	// 知识库搜索（前台）
	kbSvc := service.NewAITicketCoreService(db, log)
	r.GET("/kb/search", func(c *gin.Context) {
		keyword := c.Query("keyword")
		results := kbSvc.SearchKnowledge(keyword)
		c.JSON(200, gin.H{"code": 0, "data": results})
	})

	// AI 购物助手（前台）
	aiShopSvc := service.NewAIShoppingCoreService(db, log)
	aiShopHandler := handler.NewAIShoppingCoreHandler(aiShopSvc, log)
	r.POST("/ai-shopping/chat/:session_id", aiShopHandler.Chat)
	r.GET("/ai-shopping/history/:session_id", aiShopHandler.GetChatHistory)
}
