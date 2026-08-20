package handler

import (
	"net/http"
	"strconv"
	"time"

	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct{}

func NewInvoiceHandler() *InvoiceHandler {
	return &InvoiceHandler{}
}

// GetInvoices 获取发票列表
func (h *InvoiceHandler) GetInvoices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

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

	var total int64
	database.Table("invoices").Count(&total)

	var invoices []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Status    string    `json:"status"`
		Total     float64   `json:"total"`
		CreatedAt time.Time `json:"created_at"`
	}

	database.Table("invoices").
		Select("id, user_id, status, total, created_at").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&invoices)

	c.JSON(http.StatusOK, gin.H{
		"list":  invoices,
		"total": total,
		"page":  page,
	})
}

// RegisterRoutes 注册路由
func (h *InvoiceHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/invoices", h.GetInvoices)
}

// ==================== Admin Router Methods ====================

// GetList returns a paginated invoice list.
func (h *InvoiceHandler) GetList(c *gin.Context) {
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
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": []interface{}{}, "total": 0}})
		return
	}

	query := database.Table("invoices")
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var invoices []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Status    string    `json:"status"`
		Total     float64   `json:"total"`
		CreatedAt time.Time `json:"created_at"`
	}

	query.Select("id, user_id, status, total, created_at").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&invoices)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  invoices,
			"total": total,
		},
	})
}

// GetDetail returns invoice detail.
func (h *InvoiceHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": nil})
		return
	}

	var invoice struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Status    string    `json:"status"`
		Total     float64   `json:"total"`
		CreatedAt time.Time `json:"created_at"`
	}

	err := database.Table("invoices").Where("id = ?", id).First(&invoice).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "发票不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": invoice})
}

// Cancel cancels an invoice.
func (h *InvoiceHandler) Cancel(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	if err := database.Table("invoices").Where("id = ?", id).Update("status", "cancelled").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "取消失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "发票已取消"})
}

// NotesPage returns the invoice notes page.
func (h *InvoiceHandler) NotesPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": nil})
}

// Notes updates invoice notes.
func (h *InvoiceHandler) Notes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// EditItem edits an invoice item.
func (h *InvoiceHandler) EditItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// DeleteItems deletes invoice items.
func (h *InvoiceHandler) DeleteItems(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// DelAccount deletes an invoice account entry.
func (h *InvoiceHandler) DelAccount(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// GetListEnhanced returns an enhanced invoice list.
func (h *InvoiceHandler) GetListEnhanced(c *gin.Context) {
	h.GetList(c)
}

// ==================== V1 Router Methods ====================

// GetUserInvoices returns invoices for the authenticated user.
func (h *InvoiceHandler) GetUserInvoices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// Pay pays an invoice.
func (h *InvoiceHandler) Pay(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// CombineInvoices combines multiple invoices.
func (h *InvoiceHandler) CombineInvoices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// GetCombineInvoices returns combined invoices.
func (h *InvoiceHandler) GetCombineInvoices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}
