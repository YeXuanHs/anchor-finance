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
		authenticated.POST("/logout", authHandler.Logout)

		// 仪表盘 - 已实现
		authenticated.GET("/dashboard/stats", GetDashboardStats)
		authenticated.GET("/dashboard/income-trend", GetIncomeTrend)
		authenticated.GET("/dashboard/online-admins", GetOnlineAdmins)

		// 客户管理 - 已实现
		authenticated.GET("/users", GetUserList)
		authenticated.GET("/users/:id", GetUser)
		authenticated.POST("/users", CreateUser)
		authenticated.PUT("/users/:id", UpdateUser)
		authenticated.DELETE("/users/:id", DeleteUser)
		authenticated.GET("/users/:id/orders", GetUserOrders)
		authenticated.GET("/users/:id/invoices", GetUserInvoices)
		authenticated.GET("/users/:id/tickets", GetUserTickets)
		authenticated.GET("/users/:id/services", GetUserServices)

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

		// 工单管理 - 已实现
		authenticated.GET("/tickets", GetTicketList)
		authenticated.GET("/tickets/:id", GetTicket)
		authenticated.POST("/tickets/:id/reply", ReplyTicket)
		authenticated.POST("/tickets/:id/close", CloseTicket)
		authenticated.GET("/ticket-departments", GetTicketDepartments)
		authenticated.GET("/ticket-statuses", GetTicketStatuses)

		// 产品管理 - 已实现
		authenticated.GET("/products", GetProductList)
		authenticated.GET("/products/:id", GetProduct)
		authenticated.POST("/products", CreateProduct)
		authenticated.PUT("/products/:id", UpdateProduct)
		authenticated.DELETE("/products/:id", DeleteProduct)
		authenticated.GET("/product-groups", GetProductGroups)

		// 设置管理 - 已实现
		authenticated.GET("/settings", GetSettings)
		authenticated.GET("/settings/:group", GetSettingsByGroup)
		authenticated.PUT("/settings", UpdateSettings)

		// 菜单管理 - 已实现
		authenticated.GET("/menus", GetMenus)
		authenticated.POST("/menus", CreateMenu)
		authenticated.PUT("/menus/:id", UpdateMenu)
		authenticated.DELETE("/menus/:id", DeleteMenu)

		// 管理员管理 - 已实现
		authenticated.GET("/admins", GetAdminList)
		authenticated.GET("/roles", GetRoleList)
		authenticated.POST("/admins", CreateAdmin)
		authenticated.PUT("/admins/:id", UpdateAdmin)
		authenticated.GET("/cron-tasks", GetCronTasks)

		// 日志管理 - 已实现
		authenticated.GET("/system-logs", GetSystemLogs)
		authenticated.GET("/operation-logs", GetOperationLogs)
		authenticated.GET("/login-logs", GetLoginLogs)

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


