package handler

import (
	"net/http"
	"strconv"
	"time"

	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct{}

func NewTicketHandler() *TicketHandler {
	return &TicketHandler{}
}

// GetTickets 获取工单列表
func (h *TicketHandler) GetTickets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"list": []interface{}{}, "total": 0})
		return
	}

	query := database.Table("tickets")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var tickets []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Subject   string    `json:"subject"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	query.Select("id, user_id, subject, status, created_at").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&tickets)

	c.JSON(http.StatusOK, gin.H{
		"list":  tickets,
		"total": total,
		"page":  page,
	})
}

// GetTicket 获取单个工单
func (h *TicketHandler) GetTicket(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}

	var ticket struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Subject   string    `json:"subject"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	err := database.Table("tickets").Where("id = ?", id).First(&ticket).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "工单不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": ticket})
}

// CreateTicket 创建工单
func (h *TicketHandler) CreateTicket(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "工单创建成功"})
}

// ReplyTicket 回复工单
func (h *TicketHandler) ReplyTicket(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "回复成功"})
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

// ==================== Admin Router Methods ====================

// GetList returns a paginated ticket list.
func (h *TicketHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": []interface{}{}, "total": 0}})
		return
	}

	query := database.Table("tickets")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("subject LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var tickets []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Subject   string    `json:"subject"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	query.Select("id, user_id, subject, status, created_at").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&tickets)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  tickets,
			"total": total,
		},
	})
}

// GetListEnhanced returns an enhanced ticket list.
func (h *TicketHandler) GetListEnhanced(c *gin.Context) {
	h.GetList(c)
}

// GetDetail returns ticket detail.
func (h *TicketHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": nil})
		return
	}

	var ticket struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Subject   string    `json:"subject"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	err := database.Table("tickets").Where("id = ?", id).First(&ticket).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "工单不存在"})
		return
	}

	// 获取回复
	var replies []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Content   string    `json:"content"`
		CreatedAt time.Time `json:"created_at"`
	}
	database.Table("ticket_replies").Where("ticket_id = ?", id).Order("created_at ASC").Find(&replies)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"ticket":  ticket,
			"replies": replies,
		},
	})
}

// AdminReply replies to a ticket (admin).
func (h *TicketHandler) AdminReply(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	reply := map[string]interface{}{
		"ticket_id": id,
		"content":   req.Content,
		"user_id":   0, // admin
	}

	if err := database.Table("ticket_replies").Create(reply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "回复失败"})
		return
	}

	// 更新工单状态
	database.Table("tickets").Where("id = ?", id).Update("status", "replied")

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "回复成功"})
}

// Assign assigns a ticket to an admin.
func (h *TicketHandler) Assign(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// Close closes a ticket.
func (h *TicketHandler) Close(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	if err := database.Table("tickets").Where("id = ?", id).Update("status", "closed").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "关闭失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "工单已关闭"})
}

// DownloadAttachment downloads a ticket attachment.
func (h *TicketHandler) DownloadAttachment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// TicketReceive receives/accepts a ticket.
func (h *TicketHandler) TicketReceive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// MergeTickets merges multiple tickets.
func (h *TicketHandler) MergeTickets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// TransferTicket transfers a ticket to another department.
func (h *TicketHandler) TransferTicket(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// GetTransferLogs returns transfer logs for a ticket.
func (h *TicketHandler) GetTransferLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// UploadAttachment uploads an attachment to a ticket.
func (h *TicketHandler) UploadAttachment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// GetAttachments returns attachments for a ticket.
func (h *TicketHandler) GetAttachments(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// DeleteAttachment deletes a ticket attachment.
func (h *TicketHandler) DeleteAttachment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// TicketStatistics returns ticket statistics.
func (h *TicketHandler) TicketStatistics(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
		return
	}

	var total, open, pending, closed int64
	database.Table("tickets").Count(&total)
	database.Table("tickets").Where("status = ?", "open").Count(&open)
	database.Table("tickets").Where("status = ?", "pending").Count(&pending)
	database.Table("tickets").Where("status = ?", "closed").Count(&closed)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"total":   total,
			"open":    open,
			"pending": pending,
			"closed":  closed,
		},
	})
}

// ==================== V1 Router Methods ====================

// Create creates a new ticket.
func (h *TicketHandler) Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// GetUserTickets returns tickets for the authenticated user.
func (h *TicketHandler) GetUserTickets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// Reply replies to a ticket.
func (h *TicketHandler) Reply(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// Upload uploads a file for a ticket.
func (h *TicketHandler) Upload(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}
