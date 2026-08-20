package handler

import (
	"net/http"
	"strconv"

	"anchorfinance/internal/service"

	"github.com/gin-gonic/gin"
)

// UpgradeEnhancedHandler 升级增强处理器
type UpgradeEnhancedHandler struct {
	upgradeSvc *service.UpgradeEnhancedService
}

// NewUpgradeEnhancedHandler 创建升级增强处理器
func NewUpgradeEnhancedHandler(upgradeSvc *service.UpgradeEnhancedService) *UpgradeEnhancedHandler {
	return &UpgradeEnhancedHandler{upgradeSvc: upgradeSvc}
}

// GetUpgradeConfig 获取升级配置
func (h *UpgradeEnhancedHandler) GetUpgradeConfig(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	config, err := h.upgradeSvc.GetProductUpgradeConfig(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": config})
}

// CheckChange 检查升级变更
func (h *UpgradeEnhancedHandler) CheckChange(c *gin.Context) {
	hostIDStr := c.Param("host_id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	targetProductIDStr := c.Query("target_product_id")
	targetProductID, err := strconv.ParseUint(targetProductIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target product ID"})
		return
	}

	change, err := h.upgradeSvc.CheckChange(uint(hostID), uint(targetProductID), c.Query("cycle"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": change})
}

// FilterConfigOptions 过滤配置选项
func (h *UpgradeEnhancedHandler) FilterConfigOptions(c *gin.Context) {
	productIDStr := c.Param("product_id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var currentConfig map[string]interface{}
	c.ShouldBindJSON(&currentConfig)

	options, err := h.upgradeSvc.FilterConfigOptions(uint(productID), currentConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": options})
}

// CheckUpgradePromo 检查升级促销
func (h *UpgradeEnhancedHandler) CheckUpgradePromo(c *gin.Context) {
	hostIDStr := c.Param("host_id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	targetProductIDStr := c.Query("target_product_id")
	targetProductID, err := strconv.ParseUint(targetProductIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target product ID"})
		return
	}

	discount, err := h.upgradeSvc.CheckUpgradePromo(uint(hostID), uint(targetProductID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"discount": discount}})
}

// AllowUpgradeProducts 获取可升级产品列表
func (h *UpgradeEnhancedHandler) AllowUpgradeProducts(c *gin.Context) {
	hostIDStr := c.Param("host_id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	products, err := h.upgradeSvc.AllowUpgradeProducts(uint(hostID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

// DoUpgrade 执行升级
func (h *UpgradeEnhancedHandler) DoUpgrade(c *gin.Context) {
	hostIDStr := c.Param("host_id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	var req struct {
		TargetProductID uint   `json:"target_product_id" binding:"required"`
		Cycle           string `json:"cycle"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("user_id")

	validation, err := h.upgradeSvc.DoUpgrade(uint(hostID), req.TargetProductID, req.Cycle, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "validation": validation})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": validation})
}

// UpgradeConfigAdmin 管理员升级配置
func (h *UpgradeEnhancedHandler) UpgradeConfigAdmin(c *gin.Context) {
	hostIDStr := c.Param("host_id")
	hostID, err := strconv.ParseUint(hostIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid host ID"})
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	adminID := c.GetUint("admin_id")

	if err := h.upgradeSvc.UpgradeConfigAdmin(uint(hostID), adminID, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Upgrade config updated"})
}
