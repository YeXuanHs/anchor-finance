package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct{}

func NewInvoiceHandler() *InvoiceHandler {
	return &InvoiceHandler{}
}

// GetInvoices 获取发票列表
func (h *InvoiceHandler) GetInvoices(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"list": []interface{}{}, "total": 0})
}

// RegisterRoutes 注册路由
func (h *InvoiceHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/invoices", h.GetInvoices)
}
