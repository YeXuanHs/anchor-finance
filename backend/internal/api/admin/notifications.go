package admin

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetNotificationTemplates 获取通知模板列表
// GET /api/admin/notification-templates
func GetNotificationTemplates(c *gin.Context) {
	db := database.GetDB()
	var templates []model.NotificationTemplate
	db.Order("type ASC, name ASC").Find(&templates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    templates,
	})
}

// CreateNotificationTemplate 创建通知模板
// POST /api/admin/notification-templates
func CreateNotificationTemplate(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Type      string `json:"type" binding:"required"`
		Subject   string `json:"subject"`
		Content   string `json:"content" binding:"required"`
		Variables string `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),,
			"data": nil})
		return
	}

	db := database.GetDB()
	template := model.NotificationTemplate{
		Name:      req.Name,
		Type:      req.Type,
		Subject:   req.Subject,
		Content:   req.Content,
		Variables: req.Variables,
		Status:    "active",
	}

	if err := db.Create(&template).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建模板失败: " + err.Error(),,
			"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": template.ID,
		},
	})
}

// UpdateNotificationTemplate 更新通知模板
// PUT /api/admin/notification-templates/:id
func UpdateNotificationTemplate(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name"`
		Subject   string `json:"subject"`
		Content   string `json:"content"`
		Variables string `json:"variables"`
		Status    string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),,
			"data": nil})
		return
	}

	db := database.GetDB()
	var template model.NotificationTemplate
	if err := db.First(&template, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "模板不存在",,
			"data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Subject != "" {
		updates["subject"] = req.Subject
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Variables != "" {
		updates["variables"] = req.Variables
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&template).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteNotificationTemplate 删除通知模板
// DELETE /api/admin/notification-templates/:id
func DeleteNotificationTemplate(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var template model.NotificationTemplate
	if err := db.First(&template, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "模板不存在",,
			"data": nil})
		return
	}

	db.Delete(&template)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}
