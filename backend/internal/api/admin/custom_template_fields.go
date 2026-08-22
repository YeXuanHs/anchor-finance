package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetCustomTemplateFieldList 获取官网自定义字段列表
// GET /api/admin/custom-template-fields
func GetCustomTemplateFieldList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	pageName := c.Query("page")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.CustomTemplateField{})

	if pageName != "" {
		query = query.Where("page = ?", pageName)
	}

	var total int64
	query.Count(&total)

	var fields []model.CustomTemplateField
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("sort_order ASC").Find(&fields)

	if fields == nil {
		fields = []model.CustomTemplateField{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      fields,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateCustomTemplateField 创建官网自定义字段
// POST /api/admin/custom-template-fields
func CreateCustomTemplateField(c *gin.Context) {
	var req struct {
		Page      string `json:"page" binding:"required"`
		Name      string `json:"name" binding:"required"`
		Key       string `json:"key" binding:"required"`
		Type      string `json:"type"`
		Value     string `json:"value"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if req.Type == "" {
		req.Type = "text"
	}

	db := database.GetDB()

	// 检查Page+Key是否重复
	var count int64
	db.Model(&model.CustomTemplateField{}).Where("page = ? AND `key` = ?", req.Page, req.Key).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "该页面已存在此字段", "data": nil})
		return
	}

	field := model.CustomTemplateField{
		Page:      req.Page,
		Name:      req.Name,
		Key:       req.Key,
		Type:      req.Type,
		Value:     req.Value,
		SortOrder: req.SortOrder,
		Status:    "active",
	}

	if err := db.Create(&field).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败", "data": nil})
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

// UpdateCustomTemplateField 更新官网自定义字段
// PUT /api/admin/custom-template-fields/:id
func UpdateCustomTemplateField(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name"`
		Value     string `json:"value"`
		SortOrder int    `json:"sort_order"`
		Status    string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var field model.CustomTemplateField
	if err := db.First(&field, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "字段不存在", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Value != "" {
		updates["value"] = req.Value
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&field).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteCustomTemplateField 删除官网自定义字段
// DELETE /api/admin/custom-template-fields/:id
func DeleteCustomTemplateField(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var field model.CustomTemplateField
	if err := db.First(&field, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "字段不存在", "data": nil})
		return
	}

	db.Delete(&field)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
