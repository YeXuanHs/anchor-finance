package admin

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
	adminHandler := handler.NewAdminHandler()

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
	}
}
