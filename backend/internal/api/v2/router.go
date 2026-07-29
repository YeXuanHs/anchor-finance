package v2

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"anchorfinance/internal/handler"
	"anchorfinance/internal/api/middleware"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
)

// Deps holds shared dependencies for route registration.
type Deps struct {
	DB      *gorm.DB
	Log     *logger.Logger
	JWTKey  string
	UserSvc *service.UserService
	ProdSvc *service.ProductService
	OrdSvc  *service.OrderService
	InvSvc  *service.InvoiceService
	TicSvc  *service.TicketService
	CartSvc *service.CartService
}

// RegisterRoutes registers all v2 API routes on the given router group.
func RegisterRoutes(r *gin.RouterGroup, deps Deps) {
	authHandler := handler.NewAuthHandler(deps.UserSvc, deps.Log, deps.JWTKey)
	userHandler := handler.NewUserHandler(deps.UserSvc, deps.Log)
	productHandler := handler.NewProductHandler(deps.ProdSvc, deps.Log)
	orderHandler := handler.NewOrderHandler(deps.OrdSvc, deps.Log)
	invoiceHandler := handler.NewInvoiceHandler(deps.InvSvc, deps.Log)
	ticketHandler := handler.NewTicketHandler(deps.TicSvc, deps.Log)

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
		user.POST("/tickets/:id/attachments", ticketHandler.UploadAttachment)
		user.GET("/tickets/:id/attachments", ticketHandler.GetAttachments)
		user.DELETE("/tickets/attachments/:id", ticketHandler.DeleteAttachment)

		// 用户产品
		user.GET("/user/products", productHandler.GetUserProducts)
	}
}
