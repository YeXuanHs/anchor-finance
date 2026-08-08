package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/anchorfinance/backend/internal/model"
	"github.com/anchorfinance/backend/internal/service"
)

type SupplierHandler struct {
	supplierService *service.SupplierService
}

func NewSupplierHandler() *SupplierHandler {
	return &SupplierHandler{
		supplierService: service.NewSupplierService(),
	}
}

// GetSuppliers 获取供应商列表
func (h *SupplierHandler) GetSuppliers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	apiType := c.Query("api_type")

	params := &model.SupplierQueryParams{
		Page:     page,
		PageSize: pageSize,
		Keyword:  keyword,
		APIType:  apiType,
	}

	result, err := h.supplierService.GetSuppliers(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetSupplier 获取单个供应商
func (h *SupplierHandler) GetSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	supplier, err := h.supplierService.GetSupplierByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "供应商不存在"})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// CreateSupplier 创建供应商
func (h *SupplierHandler) CreateSupplier(c *gin.Context) {
	var req model.CreateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证API类型
	validTypes := map[string]bool{
		"manual": true,
		"zjmf":   true,
		"v10":    true,
		"anchor": true,
	}
	if !validTypes[req.APIType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的API类型"})
		return
	}

	supplier, err := h.supplierService.CreateSupplier(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

// UpdateSupplier 更新供应商
func (h *SupplierHandler) UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req model.UpdateSupplierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	supplier, err := h.supplierService.UpdateSupplier(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, supplier)
}

// DeleteSupplier 删除供应商
func (h *SupplierHandler) DeleteSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	if err := h.supplierService.DeleteSupplier(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// TestConnection 测试连接
func (h *SupplierHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	result, err := h.supplierService.TestConnection(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "success": false})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": result})
}

// SyncProducts 同步产品
func (h *SupplierHandler) SyncProducts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	result, err := h.supplierService.SyncProducts(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "同步完成", "synced_count": result})
}

// ToggleStatus 切换状态
func (h *SupplierHandler) ToggleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.supplierService.ToggleStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

// RegisterRoutes 注册路由
func (h *SupplierHandler) RegisterRoutes(r *gin.RouterGroup) {
	supplier := r.Group("/suppliers")
	{
		supplier.GET("", h.GetSuppliers)
		supplier.GET("/:id", h.GetSupplier)
		supplier.POST("", h.CreateSupplier)
		supplier.PUT("/:id", h.UpdateSupplier)
		supplier.DELETE("/:id", h.DeleteSupplier)
		supplier.POST("/:id/test-connection", h.TestConnection)
		supplier.POST("/:id/sync-products", h.SyncProducts)
		supplier.POST("/:id/status", h.ToggleStatus)
	}
}
