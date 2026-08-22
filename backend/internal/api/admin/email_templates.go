package admin

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetEmailTemplateList 获取邮件模板列表
// GET /api/admin/email-templates
func GetEmailTemplateList(c *gin.Context) {
	db := database.GetDB()
	var templates []model.EmailTemplate
	db.Where("status = ?", "active").Order("id ASC").Find(&templates)

	if templates == nil {
		templates = []model.EmailTemplate{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    templates,
	})
}

// GetEmailTemplateDetail 获取邮件模板详情
// GET /api/admin/email-templates/:id
func GetEmailTemplateDetail(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var template model.EmailTemplate
	if err := db.First(&template, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "模板不存在", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    template,
	})
}

// CreateEmailTemplate 创建邮件模板
// POST /api/admin/email-templates
func CreateEmailTemplate(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Subject   string `json:"subject"`
		Content   string `json:"content" binding:"required"`
		Variables string `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	template := model.EmailTemplate{
		Name:      req.Name,
		Subject:   req.Subject,
		Content:   req.Content,
		Variables: req.Variables,
		Status:    "active",
	}

	if err := db.Create(&template).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败", "data": nil})
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

// UpdateEmailTemplate 更新邮件模板
// PUT /api/admin/email-templates/:id
func UpdateEmailTemplate(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name"`
		Subject   string `json:"subject"`
		Content   string `json:"content"`
		Variables string `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var template model.EmailTemplate
	if err := db.First(&template, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "模板不存在", "data": nil})
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

	db.Model(&template).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteEmailTemplate 删除邮件模板
// DELETE /api/admin/email-templates/:id
func DeleteEmailTemplate(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var template model.EmailTemplate
	if err := db.First(&template, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "模板不存在", "data": nil})
		return
	}

	db.Delete(&template)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
