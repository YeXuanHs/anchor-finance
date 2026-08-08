package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct{}

func NewTicketHandler() *TicketHandler {
	return &TicketHandler{}
}

// GetTickets 获取工单列表
func (h *TicketHandler) GetTickets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"list": []interface{}{}, "total": 0})
}

// GetTicket 获取单个工单
func (h *TicketHandler) GetTicket(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// CreateTicket 创建工单
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "工单创建成功"})
}

// ReplyTicket 回复工单
func (h *TicketHandler) ReplyTicket(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "回复成功"})
}

// RegisterRoutes 注册路由
func (h *TicketHandler) RegisterRoutes(r *gin.RouterGroup) {
	ticket := r.Group("/tickets")
	{
		ticket.GET("", h.GetTickets)
		ticket.GET("/:id", h.GetTicket)
		ticket.POST("", h.CreateTicket)
		ticket.POST("/:id/reply", h.ReplyTicket)
	}
}
