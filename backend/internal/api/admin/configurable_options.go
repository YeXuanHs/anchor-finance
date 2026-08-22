package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetConfigurableOptionList 获取全局可配置项列表
// GET /api/admin/configurable-options
func GetConfigurableOptionList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	var total int64
	db.Model(&model.ConfigurableOption{}).Count(&total)

	var options []model.ConfigurableOption
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("sort_order ASC").Find(&options)

	if options == nil {
		options = []model.ConfigurableOption{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      options,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateConfigurableOption 创建全局可配置项
// POST /api/admin/configurable-options
func CreateConfigurableOption(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Type      string `json:"type" binding:"required"`
		Options   string `json:"options"`
		Default   string `json:"default_value"`
		Required  bool   `json:"required"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	option := model.ConfigurableOption{
		Name:      req.Name,
		Type:      req.Type,
		Options:   req.Options,
		Default:   req.Default,
		Required:  req.Required,
		SortOrder: req.SortOrder,
		Status:    "active",
	}

	if err := db.Create(&option).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": option.ID,
		},
	})
}

// UpdateConfigurableOption 更新全局可配置项
// PUT /api/admin/configurable-options/:id
func UpdateConfigurableOption(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name"`
		Type      string `json:"type"`
		Options   string `json:"options"`
		Default   string `json:"default_value"`
		Required  bool   `json:"required"`
		SortOrder int    `json:"sort_order"`
		Status    string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var option model.ConfigurableOption
	if err := db.First(&option, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "配置项不存在", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Options != "" {
		updates["options"] = req.Options
	}
	if req.Default != "" {
		updates["default"] = req.Default
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&option).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteConfigurableOption 删除全局可配置项
// DELETE /api/admin/configurable-options/:id
func DeleteConfigurableOption(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var option model.ConfigurableOption
	if err := db.First(&option, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "配置项不存在", "data": nil})
		return
	}

	db.Delete(&option)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
