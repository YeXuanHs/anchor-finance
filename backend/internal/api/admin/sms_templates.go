package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetSMSTemplateList 获取短信模板列表
// GET /api/admin/sms-templates
func GetSMSTemplateList(c *gin.Context) {
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
	db.Model(&model.SMSTemplate{}).Count(&total)

	var templates []model.SMSTemplate
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("id ASC").Find(&templates)

	if templates == nil {
		templates = []model.SMSTemplate{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      templates,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateSMSTemplate 创建短信模板
// POST /api/admin/sms-templates
func CreateSMSTemplate(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Code      string `json:"code" binding:"required"`
		Content   string `json:"content" binding:"required"`
		Variables string `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()

	// 检查Code是否重复
	var count int64
	db.Model(&model.SMSTemplate{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "模板编码已存在", "data": nil})
		return
	}

	template := model.SMSTemplate{
		Name:      req.Name,
		Code:      req.Code,
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

// UpdateSMSTemplate 更新短信模板
// PUT /api/admin/sms-templates/:id
func UpdateSMSTemplate(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name"`
		Content   string `json:"content"`
		Variables string `json:"variables"`
		Status    string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var template model.SMSTemplate
	if err := db.First(&template, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "模板不存在", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
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

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteSMSTemplate 删除短信模板
// DELETE /api/admin/sms-templates/:id
func DeleteSMSTemplate(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var template model.SMSTemplate
	if err := db.First(&template, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "模板不存在", "data": nil})
		return
	}

	db.Delete(&template)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
