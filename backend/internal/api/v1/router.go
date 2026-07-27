package v1

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

	// 兼容智简魔方 v1 API 格式
	r.POST("/user/login", authHandler.Login)
	r.POST("/user/register", authHandler.Register)
	r.GET("/user/info", middleware.AuthRequired(), userHandler.GetProfile)

	r.GET("/products", productHandler.GetList)
	r.GET("/products/:id", productHandler.GetDetail)

	// 需要登录
	auth := r.Group("")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/user/orders", orderHandler.GetUserOrders)
		auth.GET("/user/invoices", invoiceHandler.GetUserInvoices)
		auth.GET("/user/tickets", ticketHandler.GetUserTickets)
		auth.POST("/orders", orderHandler.Create)
		auth.GET("/orders/:id", orderHandler.GetDetail)
		auth.POST("/orders/:id/pay", orderHandler.Pay)
		auth.POST("/tickets", ticketHandler.Create)
		auth.GET("/tickets/:id", ticketHandler.GetDetail)
		auth.POST("/tickets/:id/reply", ticketHandler.Reply)
	}
}
