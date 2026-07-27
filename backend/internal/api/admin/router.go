package admin

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"anchorfinance/internal/handler"
	"anchorfinance/internal/api/middleware"
	"anchorfinance/internal/service"
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
}

func RegisterRoutes(r *gin.RouterGroup, deps Deps) {
	authHandler := handler.NewAuthHandler(deps.UserSvc, deps.Log, deps.JWTKey)
	userHandler := handler.NewUserHandler(deps.UserSvc, deps.Log)
	productHandler := handler.NewProductHandler(deps.ProdSvc, deps.Log)
	orderHandler := handler.NewOrderHandler(deps.OrdSvc, deps.Log)
	invoiceHandler := handler.NewInvoiceHandler(deps.InvSvc, deps.Log)
	ticketHandler := handler.NewTicketHandler(deps.TicSvc, deps.Log)
	adminHandler := handler.NewAdminHandler(deps.DB, deps.Log)

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

		// 公告管理
		admin.GET("/announcements", adminHandler.GetAnnouncements)
		admin.POST("/announcements", adminHandler.CreateAnnouncement)
		admin.PUT("/announcements/:id", adminHandler.UpdateAnnouncement)
		admin.DELETE("/announcements/:id", adminHandler.DeleteAnnouncement)

		// 系统设置
		admin.GET("/settings", adminHandler.GetSettings)
		admin.PUT("/settings", adminHandler.UpdateSettings)
		admin.GET("/settings/:group", adminHandler.GetSettingsByGroup)

		// 操作日志
		admin.GET("/logs", adminHandler.GetLogs)

		// 轮播图管理
		bannerHandler := handler.NewBannerHandler(deps.DB, deps.Log)
		admin.GET("/banners", bannerHandler.List)
		admin.GET("/banners/:id", bannerHandler.GetDetail)
		admin.POST("/banners", bannerHandler.Create)
		admin.PUT("/banners/:id", bannerHandler.Update)
		admin.DELETE("/banners/:id", bannerHandler.Delete)
		admin.POST("/banners/:id/status", bannerHandler.SetStatus)

		// 支付方式管理
		paymentGwHandler := handler.NewPaymentGatewayHandler(deps.DB, deps.Log)
		admin.GET("/payment-gateways", paymentGwHandler.List)
		admin.GET("/payment-gateways/:id", paymentGwHandler.GetDetail)
		admin.POST("/payment-gateways", paymentGwHandler.Create)
		admin.PUT("/payment-gateways/:id", paymentGwHandler.Update)
		admin.DELETE("/payment-gateways/:id", paymentGwHandler.Delete)
		admin.POST("/payment-gateways/:id/status", paymentGwHandler.SetStatus)
		admin.POST("/payment-gateways/:id/test", paymentGwHandler.TestConnection)

		// OAuth提供商管理
		oauthProvHandler := handler.NewOAuthProviderHandler(deps.DB, deps.Log)
		admin.GET("/oauth-providers", oauthProvHandler.List)
		admin.GET("/oauth-providers/:id", oauthProvHandler.GetDetail)
		admin.POST("/oauth-providers", oauthProvHandler.Create)
		admin.PUT("/oauth-providers/:id", oauthProvHandler.Update)
		admin.DELETE("/oauth-providers/:id", oauthProvHandler.Delete)
		admin.POST("/oauth-providers/:id/status", oauthProvHandler.SetStatus)

		// Cron定时任务
		cronHandler := handler.NewCronHandler(deps.DB, deps.Log)
		admin.GET("/cron-tasks", cronHandler.GetList)
		admin.GET("/cron-tasks/:id", cronHandler.GetDetail)
		admin.POST("/cron-tasks", cronHandler.Create)
		admin.PUT("/cron-tasks/:id", cronHandler.Update)
		admin.DELETE("/cron-tasks/:id", cronHandler.Delete)
		admin.POST("/cron-tasks/:id/status", cronHandler.SetStatus)
		admin.POST("/cron-tasks/:id/run", cronHandler.RunTask)
		admin.GET("/cron-tasks/:id/logs", cronHandler.GetLogs)

		// DCIM管理
		dcimHandler := handler.NewDcimHandler(deps.DB, deps.Log)
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
		appstoreHandler := handler.NewAppstoreHandler(deps.DB, deps.Log)
		admin.GET("/apps", appstoreHandler.AdminList)
		admin.POST("/apps", appstoreHandler.AdminCreate)
		admin.PUT("/apps/:id", appstoreHandler.AdminUpdate)
		admin.DELETE("/apps/:id", appstoreHandler.AdminDelete)

		// 插件管理
		pluginHandler := handler.NewPluginHandler(deps.DB, deps.Log)
		admin.GET("/plugins", pluginHandler.List)
		admin.POST("/plugins", pluginHandler.Install)
		admin.DELETE("/plugins/:id", pluginHandler.Uninstall)
		admin.POST("/plugins/:id/enable", pluginHandler.Enable)
		admin.POST("/plugins/:id/disable", pluginHandler.Disable)
		admin.PUT("/plugins/:id/config", pluginHandler.UpdateConfig)

		// 模块供给
		provisionHandler := handler.NewProvisionHandler(deps.DB, deps.Log)
		admin.GET("/provisions", provisionHandler.List)
		admin.GET("/provisions/:id", provisionHandler.GetDetail)
		admin.POST("/provisions", provisionHandler.Create)
		admin.PUT("/provisions/:id", provisionHandler.Update)
		admin.DELETE("/provisions/:id", provisionHandler.Delete)
		admin.POST("/provisions/:id/test", provisionHandler.TestConnection)

		// 客户分组
		clientGroupHandler := handler.NewClientGroupHandler(deps.DB, deps.Log)
		admin.GET("/client-groups", clientGroupHandler.List)
		admin.GET("/client-groups/:id", clientGroupHandler.Get)
		admin.POST("/client-groups", clientGroupHandler.Create)
		admin.PUT("/client-groups/:id", clientGroupHandler.Update)
		admin.DELETE("/client-groups/:id", clientGroupHandler.Delete)

		// 客户服务
		clientSvcHandler := handler.NewClientServiceHandler(deps.DB, deps.Log)
		admin.GET("/client-services", clientSvcHandler.List)
		admin.GET("/client-services/:id", clientSvcHandler.Get)
		admin.POST("/client-services/:id/suspend", clientSvcHandler.Suspend)
		admin.POST("/client-services/:id/terminate", clientSvcHandler.Terminate)
		admin.POST("/client-services/:id/renew", clientSvcHandler.Renew)

		// 用户管理
		userManageHandler := handler.NewUserManageHandler(deps.DB, deps.Log)
		admin.GET("/user-manage/search", userManageHandler.Search)
		admin.POST("/user-manage/:id/ban", userManageHandler.Ban)
		admin.POST("/user-manage/:id/unban", userManageHandler.Unban)
		admin.POST("/user-manage/:id/balance", userManageHandler.AdjustBalance)
		admin.POST("/user-manage/:id/reset-password", userManageHandler.ResetPassword)

		// 用户备注
		userRemarkHandler := handler.NewUserRemarkHandler(deps.DB, deps.Log)
		admin.GET("/user-remarks", userRemarkHandler.List)
		admin.POST("/user-remarks", userRemarkHandler.Add)
		admin.DELETE("/user-remarks/:id", userRemarkHandler.AdminDelete)

		// 邮件模板
		emailTplHandler := handler.NewEmailTemplateHandler(deps.DB, deps.Log)
		admin.GET("/email-templates", emailTplHandler.List)
		admin.GET("/email-templates/:id", emailTplHandler.GetDetail)
		admin.POST("/email-templates", emailTplHandler.Create)
		admin.PUT("/email-templates/:id", emailTplHandler.Update)
		admin.DELETE("/email-templates/:id", emailTplHandler.Delete)
		admin.POST("/email-templates/:id/test", emailTplHandler.SendTest)

		// 菜单管理
		menuHandler := handler.NewMenuHandler(deps.DB, deps.Log)
		admin.GET("/menus", menuHandler.List)
		admin.GET("/menus/tree", menuHandler.GetTree)
		admin.POST("/menus", menuHandler.Create)
		admin.PUT("/menus/:id", menuHandler.Update)
		admin.DELETE("/menus/:id", menuHandler.Delete)
		admin.POST("/menus/sort", menuHandler.Sort)

		// 日志记录
		logHandler := handler.NewLogRecordHandler(deps.DB, deps.Log)
		admin.GET("/log-records", logHandler.List)
		admin.GET("/log-records/:id", logHandler.GetDetail)
		admin.GET("/log-records/stats", logHandler.Stats)
		admin.POST("/log-records/export", logHandler.Export)

		// 规则管理
		ruleHandler := handler.NewRuleHandler(deps.DB, deps.Log)
		admin.GET("/rules", ruleHandler.List)
		admin.GET("/rules/:id", ruleHandler.GetDetail)
		admin.POST("/rules", ruleHandler.Create)
		admin.PUT("/rules/:id", ruleHandler.Update)
		admin.DELETE("/rules/:id", ruleHandler.Delete)
		admin.POST("/rules/:id/enable", ruleHandler.Enable)
		admin.POST("/rules/:id/disable", ruleHandler.Disable)
		admin.POST("/rules/:id/test", ruleHandler.Test)

		// 促销管理
		saleHandler := handler.NewSaleHandler(deps.DB, deps.Log)
		admin.GET("/sales", saleHandler.List)
		admin.GET("/sales/:id", saleHandler.GetDetail)
		admin.POST("/sales", saleHandler.Create)
		admin.PUT("/sales/:id", saleHandler.Update)
		admin.DELETE("/sales/:id", saleHandler.Delete)
		admin.POST("/sales/:id/status", saleHandler.SetStatus)

		// 消息发送
		sendMsgHandler := handler.NewSendMessageHandler(deps.DB, deps.Log)
		admin.POST("/messages/send", sendMsgHandler.Send)
		admin.POST("/messages/batch", sendMsgHandler.BatchSend)
		admin.GET("/messages/records", sendMsgHandler.GetRecords)

		// 工单部门
		ticketDeptHandler := handler.NewTicketDeptHandler(deps.DB, deps.Log)
		admin.GET("/ticket-depts", ticketDeptHandler.List)
		admin.POST("/ticket-depts", ticketDeptHandler.Create)
		admin.PUT("/ticket-depts/:id", ticketDeptHandler.Update)
		admin.DELETE("/ticket-depts/:id", ticketDeptHandler.Delete)

		// 主机管理
		hostHandler := handler.NewHostHandler(deps.DB, deps.Log)
		admin.GET("/hosts", hostHandler.List)
		admin.GET("/hosts/:id", hostHandler.GetDetail)
		admin.POST("/hosts/:id/boot", hostHandler.Boot)
		admin.POST("/hosts/:id/shutdown", hostHandler.Shutdown)
		admin.POST("/hosts/:id/reboot", hostHandler.Reboot)

		// 系统配置
		configServerHandler := handler.NewConfigServerHandler(deps.DB, deps.Log)
		admin.GET("/config/servers", configServerHandler.List)
		admin.POST("/config/servers", configServerHandler.Create)
		admin.PUT("/config/servers/:id", configServerHandler.Update)
		admin.DELETE("/config/servers/:id", configServerHandler.Delete)

		configOptionHandler := handler.NewConfigOptionHandler(deps.DB, deps.Log)
		admin.GET("/config/options", configOptionHandler.List)
		admin.POST("/config/options", configOptionHandler.Create)
		admin.PUT("/config/options/:id", configOptionHandler.Update)
		admin.DELETE("/config/options/:id", configOptionHandler.Delete)

		configGeneralHandler := handler.NewConfigGeneralHandler(deps.DB, deps.Log)
		admin.GET("/config/general", configGeneralHandler.Get)
		admin.PUT("/config/general", configGeneralHandler.Update)

		configCertifiHandler := handler.NewConfigCertifiHandler(deps.DB, deps.Log)
		admin.GET("/config/certifi", configCertifiHandler.Get)
		admin.PUT("/config/certifi", configCertifiHandler.Update)
	}
}

// RegisterPublicRoutes registers public-facing API routes for banners, payment methods, etc.
func RegisterPublicRoutes(r *gin.RouterGroup, db *gorm.DB, log *logger.Logger) {
	// 公开轮播图
	bannerHandler := handler.NewBannerHandler(db, log)
	r.GET("/banners", bannerHandler.GetActive)

	// 公开支付方式列表
	paymentGwHandler := handler.NewPaymentGatewayHandler(db, log)
	r.GET("/payment-gateways", paymentGwHandler.GetEnabled)

	// 公开OAuth提供商列表
	oauthProvHandler := handler.NewOAuthProviderHandler(db, log)
	r.GET("/oauth-providers", oauthProvHandler.GetEnabled)
}
