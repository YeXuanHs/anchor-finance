package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetThemeList 获取主题列表
// GET /api/admin/themes
func GetThemeList(c *gin.Context) {
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
	db.Model(&model.Theme{}).Count(&total)

	var themes []model.Theme
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("id ASC").Find(&themes)

	if themes == nil {
		themes = []model.Theme{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      themes,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetActiveTheme 获取当前激活主题
// GET /api/admin/themes/active
func GetActiveTheme(c *gin.Context) {
	db := database.GetDB()
	var theme model.Theme
	if err := db.Where("is_default = ?", true).First(&theme).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "未设置默认主题", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": theme})
}

// SetDefaultTheme 设置默认主题
// POST /api/admin/themes/:id/set-default
func SetDefaultTheme(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的主题ID", "data": nil})
		return
	}

	db := database.GetDB()
	var theme model.Theme
	if err := db.First(&theme, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "主题不存在", "data": nil})
		return
	}

	// 事务：取消所有默认，再设置新默认
	db.Transaction(func(tx *gorm.DB) error {
		db.Model(&model.Theme{}).Where("is_default = ?", true).Update("is_default", false)
		db.Model(&theme).Update("is_default", true)
		return nil
	})

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "设置成功", "data": nil})
}

// CreateTheme 创建主题
// POST /api/admin/themes
func CreateTheme(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Description string `json:"description"`
		Thumbnail   string `json:"thumbnail"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()

	// 检查Code是否重复
	var count int64
	db.Model(&model.Theme{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "主题代码已存在", "data": nil})
		return
	}

	theme := model.Theme{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Thumbnail:   req.Thumbnail,
		Status:      "active",
	}

	if err := db.Create(&theme).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": theme.ID,
		},
	})
}

// UpdateTheme 更新主题
// PUT /api/admin/themes/:id
func UpdateTheme(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Thumbnail   string `json:"thumbnail"`
		Status      string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var theme model.Theme
	if err := db.First(&theme, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "主题不存在", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Thumbnail != "" {
		updates["thumbnail"] = req.Thumbnail
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&theme).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteTheme 删除主题
// DELETE /api/admin/themes/:id
func DeleteTheme(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var theme model.Theme
	if err := db.First(&theme, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "主题不存在", "data": nil})
		return
	}

	if theme.IsDefault {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "不能删除默认主题", "data": nil})
		return
	}

	db.Delete(&theme)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
