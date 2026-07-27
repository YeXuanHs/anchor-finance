package v1

import (
	"github.com/gin-gonic/gin"
	"anchorfinance/internal/handler"
	"anchorfinance/internal/api/middleware"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"gorm.io/gorm"
)

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
