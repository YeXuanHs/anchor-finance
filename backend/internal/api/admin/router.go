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

		// 客户管理
		authenticated.GET("/users", getUsers)
		authenticated.GET("/users/:id", getUser)
		authenticated.POST("/users", createUser)
		authenticated.PUT("/users/:id", updateUser)
		authenticated.DELETE("/users/:id", deleteUser)

		// 订单管理
		authenticated.GET("/orders", getOrders)
		authenticated.GET("/orders/:id", getOrder)

		// 服务管理
		authenticated.GET("/services", getServices)
		authenticated.GET("/services/:id", getService)

		// 账单管理
		authenticated.GET("/invoices", getInvoices)
		authenticated.GET("/invoices/:id", getInvoice)

		// 工单管理
		authenticated.GET("/tickets", getTickets)
		authenticated.GET("/tickets/:id", getTicket)

		// 产品管理
		authenticated.GET("/products", getProducts)
		authenticated.GET("/products/:id", getProduct)

		// 插件管理
		authenticated.GET("/plugins", getPlugins)

		// 设置
		authenticated.GET("/settings", getSettings)
		authenticated.PUT("/settings", updateSettings)

		// 菜单
		authenticated.GET("/menus", getMenus)
	}
}

// 以下为占位函数，后续实现具体业务逻辑

func getDashboardStats(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func getUsers(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{"list": []interface{}{}, "total": 0}})
}

func getUser(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func createUser(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func updateUser(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func deleteUser(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func getOrders(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{"list": []interface{}{}, "total": 0}})
}

func getOrder(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func getServices(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{"list": []interface{}{}, "total": 0}})
}

func getService(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func getInvoices(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{"list": []interface{}{}, "total": 0}})
}

func getInvoice(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func getTickets(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{"list": []interface{}{}, "total": 0}})
}

func getTicket(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func getProducts(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{"list": []interface{}{}, "total": 0}})
}

func getProduct(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func getPlugins(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{"list": []interface{}{}, "total": 0}})
}

func getSettings(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func updateSettings(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func getMenus(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}
