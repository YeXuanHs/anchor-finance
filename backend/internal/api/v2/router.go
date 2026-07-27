package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/anchor-finance/backend/internal/handler"
	"github.com/anchor-finance/backend/internal/api/middleware"
)

func RegisterRoutes(r *gin.RouterGroup) {
	authHandler := handler.NewAuthHandler()
	userHandler := handler.NewUserHandler()
	productHandler := handler.NewProductHandler()
	orderHandler := handler.NewOrderHandler()
	invoiceHandler := handler.NewInvoiceHandler()
	ticketHandler := handler.NewTicketHandler()

	// 认证
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.POST("/refresh", middleware.AuthRequired(), authHandler.RefreshToken)
		auth.POST("/logout", middleware.AuthRequired(), authHandler.Logout)
	}

	// 产品（公开）
	r.GET("/products", productHandler.GetList)
	r.GET("/products/:id", productHandler.GetDetail)

	// 需要登录
	user := r.Group("")
	user.Use(middleware.AuthRequired())
	{
		// 用户
		user.GET("/user/profile", userHandler.GetProfile)
		user.PUT("/user/profile", userHandler.UpdateProfile)
		user.POST("/user/change-password", userHandler.ChangePassword)

		// 订单
		user.POST("/orders", orderHandler.Create)
		user.GET("/orders", orderHandler.GetUserOrders)
		user.GET("/orders/:id", orderHandler.GetDetail)
		user.POST("/orders/:id/pay", orderHandler.Pay)
		user.POST("/orders/:id/cancel", orderHandler.Cancel)

		// 账单
		user.GET("/invoices", invoiceHandler.GetUserInvoices)
		user.GET("/invoices/:id", invoiceHandler.GetDetail)
		user.POST("/invoices/:id/pay", invoiceHandler.Pay)

		// 工单
		user.POST("/tickets", ticketHandler.Create)
		user.GET("/tickets", ticketHandler.GetUserTickets)
		user.GET("/tickets/:id", ticketHandler.GetDetail)
		user.POST("/tickets/:id/reply", ticketHandler.Reply)
		user.POST("/tickets/:id/close", ticketHandler.Close)

		// 用户产品
		user.GET("/user/products", productHandler.GetUserProducts)
	}
}
