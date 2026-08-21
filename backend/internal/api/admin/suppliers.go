package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetSupplierList 获取供应商列表
// GET /api/admin/suppliers
func GetSupplierList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Supplier{})

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var suppliers []model.Supplier
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&suppliers)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      suppliers,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetSupplierSummary 获取供应商统计
// GET /api/admin/suppliers/summary
func GetSupplierSummary(c *gin.Context) {
	db := database.GetDB()

	var activeCount int64
	db.Model(&model.Supplier{}).Where("status = ?", "active").Count(&activeCount)

	var disabledCount int64
	db.Model(&model.Supplier{}).Where("status = ?", "disabled").Count(&disabledCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"active":   activeCount,
			"disabled": disabledCount,
			"total":    activeCount + disabledCount,
		},
	})
}

// GetSupplierDetail 获取供应商详情
// GET /api/admin/suppliers/:id
func GetSupplierDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的供应商ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "供应商不存在",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    supplier,
	})
}

// CreateSupplier 创建供应商
// POST /api/admin/suppliers
func CreateSupplier(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		APIURL      string `json:"api_url"`
		APIKey      string `json:"api_key"`
		APISecret   string `json:"api_secret"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	supplier := model.Supplier{
		Name:        req.Name,
		Type:        req.Type,
		APIURL:      req.APIURL,
		APIKey:      req.APIKey,
		APISecret:   req.APISecret,
		Description: req.Description,
		Status:      "active",
	}

	if err := db.Create(&supplier).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建供应商失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": supplier.ID,
		},
	})
}

// UpdateSupplier 更新供应商
// PUT /api/admin/suppliers/:id
func UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的供应商ID",
			"data":    nil,
		})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		APIURL      string `json:"api_url"`
		APIKey      string `json:"api_key"`
		APISecret   string `json:"api_secret"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "供应商不存在",
			"data":    nil,
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.APIURL != "" {
		updates["api_url"] = req.APIURL
	}
	if req.APIKey != "" {
		updates["api_key"] = req.APIKey
	}
	if req.APISecret != "" {
		updates["api_secret"] = req.APISecret
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	db.Model(&supplier).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    nil,
	})
}

// DeleteSupplier 删除供应商
// DELETE /api/admin/suppliers/:id
func DeleteSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的供应商ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "供应商不存在",
			"data":    nil,
		})
		return
	}

	db.Delete(&supplier)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"data":    nil,
	})
}

// GetSupplierProviderTypes 获取供应商类型列表
// GET /api/admin/suppliers/provider-types
func GetSupplierProviderTypes(c *gin.Context) {
	// 返回支持的供应商类型
	types := []gin.H{
		{"id": "manual", "name": "手动管理", "description": "不对接API，手动管理"},
		{"id": "zjmf", "name": "zjmf接口", "description": "对接zjmf API"},
		{"id": "v10", "name": "v10接口", "description": "对接v10 API"},
		{"id": "anchor", "name": "锚点接口", "description": "对接锚点财务自有API"},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    types,
	})
}

// GetSupplierProducts 获取供应商产品
// GET /api/admin/suppliers/:id/products
func GetSupplierProducts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的供应商ID",
			"data":    nil,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	db := database.GetDB()
	var total int64
	db.Model(&model.SupplierProduct{}).Where("supplier_id = ?", id).Count(&total)

	var products []model.SupplierProduct
	offset := (page - 1) * pageSize
	db.Where("supplier_id = ?", id).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      products,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
