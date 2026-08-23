package admin

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetCustomFieldList 获取自定义字段列表
// GET /api/admin/custom-fields
func GetCustomFieldList(c *gin.Context) {
	db := database.GetDB()
	var fields []model.CustomField
	db.Where("status = ?", "active").Order("sort_order ASC").Find(&fields)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    fields,
	})
}

// CreateCustomField 创建自定义字段
// POST /api/admin/custom-fields
func CreateCustomField(c *gin.Context) {
	var req struct {
		Name         string `json:"name" binding:"required"`
		Label        string `json:"label" binding:"required"`
		Type         string `json:"type" binding:"required"`
		Options      string `json:"options"`
		Required     bool   `json:"required"`
		DefaultValue string `json:"default_value"`
		SortOrder    int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	// 验证字段类型
	validTypes := map[string]bool{
		"text": true, "textarea": true, "select": true,
		"checkbox": true, "radio": true, "date": true, "number": true,
	}
	if !validTypes[req.Type] {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的字段类型",
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	field := model.CustomField{
		Name:         req.Name,
		Label:        req.Label,
		Type:         req.Type,
		Options:      req.Options,
		Required:     req.Required,
		DefaultValue: req.DefaultValue,
		SortOrder:    req.SortOrder,
		Status:       "active",
	}

	if err := db.Create(&field).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建字段失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": field.ID,
		},
	})
}

// UpdateCustomField 更新自定义字段
// PUT /api/admin/custom-fields/:id
func UpdateCustomField(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Label        string `json:"label"`
		Options      string `json:"options"`
		Required     *bool  `json:"required"`
		DefaultValue string `json:"default_value"`
		SortOrder    int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var field model.CustomField
	if err := db.First(&field, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "字段不存在",
			"data": nil,
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Label != "" {
		updates["label"] = req.Label
	}
	if req.Options != "" {
		updates["options"] = req.Options
	}
	if req.Required != nil {
		updates["required"] = *req.Required
	}
	if req.DefaultValue != "" {
		updates["default_value"] = req.DefaultValue
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}

	db.Model(&field).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteCustomField 删除自定义字段
// DELETE /api/admin/custom-fields/:id
func DeleteCustomField(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var field model.CustomField
	if err := db.First(&field, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "字段不存在",
			"data": nil,
		})
		return
	}

	db.Delete(&field)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}
