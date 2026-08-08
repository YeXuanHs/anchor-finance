package router

import (
	"github.com/gin-gonic/gin"
	"github.com/anchorfinance/backend/internal/handler"
	"github.com/anchorfinance/backend/internal/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS中间件
	r.Use(middleware.CORS())

	// API路由组
	api := r.Group("/api")
	{
		// 公开路由
		public := api.Group("")
		{
			public.POST("/login", handler.Login)
			public.POST("/logout", handler.Logout)
		}

		// 需要认证的路由
		admin := api.Group("/admin")
		admin.Use(middleware.Auth())
		{
			// 仪表盘
			admin.GET("/dashboard/stats", handler.GetDashboardStats)
			admin.GET("/dashboard/income", handler.GetIncomeData)
			admin.GET("/dashboard/recent-orders", handler.GetRecentOrders)
			admin.GET("/dashboard/recent-tickets", handler.GetRecentTickets)

			// 客户管理
			clientHandler := handler.NewClientHandler()
			clientHandler.RegisterRoutes(admin)

			// 订单管理
			orderHandler := handler.NewOrderHandler()
			orderHandler.RegisterRoutes(admin)

			// 工单管理
			ticketHandler := handler.NewTicketHandler()
			ticketHandler.RegisterRoutes(admin)

			// 交易管理
			transactionHandler := handler.NewTransactionHandler()
			transactionHandler.RegisterRoutes(admin)

			// 发票管理
			invoiceHandler := handler.NewInvoiceHandler()
			invoiceHandler.RegisterRoutes(admin)

			// 产品管理
			productHandler := handler.NewProductHandler()
			productHandler.RegisterRoutes(admin)

			// 供应商管理
			supplierHandler := handler.NewSupplierHandler()
			supplierHandler.RegisterRoutes(admin)

			// 插件管理
			pluginHandler := handler.NewPluginHandler()
			pluginHandler.RegisterRoutes(admin)

			// 菜单管理
			menuHandler := handler.NewMenuHandler()
			menuHandler.RegisterRoutes(admin)

			// 管理员管理
			adminHandler := handler.NewAdminHandler()
			adminHandler.RegisterRoutes(admin)

			// 角色管理
			roleHandler := handler.NewRoleHandler()
			roleHandler.RegisterRoutes(admin)

			// 系统设置
			settingHandler := handler.NewSettingHandler()
			settingHandler.RegisterRoutes(admin)

			// 日志管理
			logHandler := handler.NewLogHandler()
			logHandler.RegisterRoutes(admin)

			// 系统信息
			systemHandler := handler.NewSystemHandler()
			systemHandler.RegisterRoutes(admin)
		}
	}

	return r
}
