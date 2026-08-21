package admin

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetProductTypeList 获取产品类型列表
// GET /api/admin/product-types
func GetProductTypeList(c *gin.Context) {
	db := database.GetDB()
	var types []model.ProductType
	db.Where("status = ?", "active").Order("sort_order ASC").Find(&types)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    types,
	})
}

// CreateProductType 创建产品类型
// POST /api/admin/product-types
func CreateProductType(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	db := database.GetDB()
	productType := model.ProductType{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      "active",
	}

	if err := db.Create(&productType).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建类型失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": productType.ID,
		},
	})
}

// UpdateProductType 更新产品类型
// PUT /api/admin/product-types/:id
func UpdateProductType(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	db := database.GetDB()
	var productType model.ProductType
	if err := db.First(&productType, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "类型不存在",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}

	db.Model(&productType).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteProductType 删除产品类型
// DELETE /api/admin/product-types/:id
func DeleteProductType(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var productType model.ProductType
	if err := db.First(&productType, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "类型不存在",
		})
		return
	}

	db.Delete(&productType)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}
