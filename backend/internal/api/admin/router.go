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

		// 仪表盘
		authenticated.GET("/dashboard/stats", getDashboardStats)

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

// 仪表盘统计（占位）
func getDashboardStats(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}
