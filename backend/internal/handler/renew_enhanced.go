package handler

import (
	"net/http"
	"strconv"

	"anchorfinance/internal/service"

	"github.com/gin-gonic/gin"
)

// RenewEnhancedHandler 续费增强处理器
type RenewEnhancedHandler struct {
	renewSvc *service.RenewEnhancedService
}

// NewRenewEnhancedHandler 创建续费增强处理器
func NewRenewEnhancedHandler(renewSvc *service.RenewEnhancedService) *RenewEnhancedHandler {
	return &RenewEnhancedHandler{renewSvc: renewSvc}
}

// GetRenewalPage 获取续费页面
func (h *RenewEnhancedHandler) GetRenewalPage(c *gin.Context) {
	hostIDStr := c.Param("host_id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	page, err := h.renewSvc.GetRenewalPage(uint(hostID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": page})
}

// GetRenewalPrice 获取续费价格
func (h *RenewEnhancedHandler) GetRenewalPrice(c *gin.Context) {
	hostIDStr := c.Param("host_id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	cycle := c.Query("cycle")
	if cycle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cycle is required"})
		return
	}

	price, err := h.renewSvc.CalculatedPrice(uint(hostID), cycle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"price": price, "cycle": cycle}})
}

// SubmitRenewal 提交续费
func (h *RenewEnhancedHandler) SubmitRenewal(c *gin.Context) {
	var params service.RenewParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	result, err := h.renewSvc.Renew(params, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

// BatchRenew 批量续费
func (h *RenewEnhancedHandler) BatchRenew(c *gin.Context) {
	var req struct {
		Items []service.RenewParams `json:"items" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	results, err := h.renewSvc.BatchRenew(req.Items, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}

// SetAutoRenew 设置自动续费
func (h *RenewEnhancedHandler) SetAutoRenew(c *gin.Context) {
	hostIDStr := c.Param("host_id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.renewSvc.SetAutoRenew(uint(hostID), req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Auto-renew setting updated"})
}

// SetPayType 设置支付方式
func (h *RenewEnhancedHandler) SetPayType(c *gin.Context) {
	hostIDStr := c.Param("host_id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	var req struct {
		PayType string `json:"pay_type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.renewSvc.SetPayType(uint(hostID), req.PayType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pay type updated"})
}

// GetPayType 获取支付方式
func (h *RenewEnhancedHandler) GetPayType(c *gin.Context) {
	hostIDStr := c.Param("host_id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	payType, err := h.renewSvc.GetPayType(uint(hostID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"pay_type": payType}})
}

// DeleteRenewInvoice 删除续费发票
func (h *RenewEnhancedHandler) DeleteRenewInvoice(c *gin.Context) {
	invoiceIDStr := c.Param("id")
	invoiceID, err := strconv.ParseUint(invoiceIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	if err := h.renewSvc.DeleteRenewInvoice(uint(invoiceID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Renew invoice deleted"})
}
