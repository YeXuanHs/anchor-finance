package handler

import (
	"net/http"
	"strconv"

	"anchorfinance/internal/service"

	"github.com/gin-gonic/gin"
)

// EmailEnhancedHandler 邮件增强处理器
type EmailEnhancedHandler struct {
	emailSvc *service.EmailEnhancedService
}

// NewEmailEnhancedHandler 创建邮件增强处理器
func NewEmailEnhancedHandler(emailSvc *service.EmailEnhancedService) *EmailEnhancedHandler {
	return &EmailEnhancedHandler{emailSvc: emailSvc}
}

// GetTemplates 获取邮件模板列表
func (h *EmailEnhancedHandler) GetTemplates(c *gin.Context) {
	templates, err := h.emailSvc.GetTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": templates})
}

// CreateTemplate 创建邮件模板
func (h *EmailEnhancedHandler) CreateTemplate(c *gin.Context) {
	var tpl service.EmailTemplate
	if err := c.ShouldBindJSON(&tpl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.emailSvc.CreateTemplate(&tpl); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": tpl})
}

// UpdateTemplate 更新邮件模板
func (h *EmailEnhancedHandler) UpdateTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.emailSvc.UpdateTemplate(uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template updated"})
}

// DeleteTemplate 删除邮件模板
func (h *EmailEnhancedHandler) DeleteTemplate(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid template ID"})
		return
	}

	if err := h.emailSvc.DeleteTemplate(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted"})
}

// SendTestEmail 发送测试邮件
func (h *EmailEnhancedHandler) SendTestEmail(c *gin.Context) {
	var req struct {
		To      string `json:"to" binding:"required"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.emailSvc.SendEmailDirect(req.To, req.Subject, req.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Test email sent"})
}

// SendBatchEmail 批量发送邮件
func (h *EmailEnhancedHandler) SendBatchEmail(c *gin.Context) {
	var req struct {
		TemplateID  uint   `json:"template_id" binding:"required"`
		TargetGroup string `json:"target_group" binding:"required"`
		Params      map[string]string `json:"params"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	batch, err := h.emailSvc.SendBatchEmail(req.TemplateID, req.TargetGroup, req.Params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": batch})
}

// GetEmailLogs 获取邮件日志
func (h *EmailEnhancedHandler) GetEmailLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.emailSvc.GetEmailLogs(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  page,
	})
}

// GetEmailStats 获取邮件统计
func (h *EmailEnhancedHandler) GetEmailStats(c *gin.Context) {
	period := c.DefaultQuery("period", "month")

	stats, err := h.emailSvc.GetEmailStats(period)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}

// GetSupportConfig 获取支持邮箱配置
func (h *EmailEnhancedHandler) GetSupportConfig(c *gin.Context) {
	config, err := h.emailSvc.GetSupportConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": config})
}

// UpdateSupportConfig 更新支持邮箱配置
func (h *EmailEnhancedHandler) UpdateSupportConfig(c *gin.Context) {
	var config service.EmailSupport
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.emailSvc.UpdateSupportConfig(&config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Support config updated"})
}
