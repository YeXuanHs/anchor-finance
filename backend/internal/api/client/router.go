package client

import (
	"github.com/YeXuanHs/anchor-finance/internal/middleware"
	"github.com/YeXuanHs/anchor-finance/internal/service"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置用户前台路由
func SetupRoutes(r *gin.RouterGroup, authService *service.AuthService) {
	// 公开路由（不需要认证）
	public := r.Group("")
	{
		public.POST("/login", userLogin(authService))
		public.POST("/register", userRegister)
		public.POST("/password/reset", resetPassword)
	}

	// 需要认证的路由
	authenticated := r.Group("")
	authenticated.Use(middleware.JWTAuth(authService))
	{
		// 用户信息
		authenticated.GET("/user/info", getUserInfo)
		authenticated.PUT("/user/info", updateUserInfo)
		authenticated.PUT("/user/password", changePassword)

		// 产品/服务
		authenticated.GET("/services", getUserServices)
		authenticated.GET("/services/:id", getUserService)
		authenticated.POST("/services/:id/renew", renewService)

		// 订单
		authenticated.GET("/orders", getUserOrders)
		authenticated.GET("/orders/:id", getUserOrder)
		authenticated.POST("/orders", createOrder)

		// 账单
		authenticated.GET("/invoices", getUserInvoices)
		authenticated.GET("/invoices/:id", getUserInvoice)
		authenticated.POST("/invoices/:id/pay", payInvoice)

		// 工单
		authenticated.GET("/tickets", getUserTickets)
		authenticated.GET("/tickets/:id", getUserTicket)
		authenticated.POST("/tickets", createTicket)
		authenticated.POST("/tickets/:id/reply", replyTicket)

		// 财务
		authenticated.GET("/finance/balance", getBalance)
		authenticated.POST("/finance/recharge", recharge)

		// 产品订购
		authenticated.GET("/products", getProducts)
		authenticated.GET("/products/:id", getProduct)
		authenticated.GET("/products/:id/config-options", getConfigOptions)
	}
}

// 以下为占位函数，后续实现具体业务逻辑

func userLogin(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: 实现用户登录
		c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
	}
}

func userRegister(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func resetPassword(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func getUserInfo(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func updateUserInfo(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func changePassword(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func getUserServices(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{"list": []interface{}{}, "total": 0}})
}

func getUserService(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func renewService(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func getUserOrders(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{"list": []interface{}{}, "total": 0}})
}

func getUserOrder(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func createOrder(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func getUserInvoices(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{"list": []interface{}{}, "total": 0}})
}

func getUserInvoice(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func payInvoice(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func getUserTickets(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{"list": []interface{}{}, "total": 0}})
}

func getUserTicket(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func createTicket(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func replyTicket(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func getBalance(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func recharge(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success"})
}

func getProducts(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{"list": []interface{}{}, "total": 0}})
}

func getProduct(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}

func getConfigOptions(c *gin.Context) {
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": gin.H{}})
}
