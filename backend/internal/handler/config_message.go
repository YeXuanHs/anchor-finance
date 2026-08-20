package handler

import (
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigMessageHandler struct {
	svc *service.ConfigMessageService
	log *logger.Logger
}

func NewConfigMessageHandler(svc *service.ConfigMessageService, log *logger.Logger) *ConfigMessageHandler {
	return &ConfigMessageHandler{svc: svc, log: log}
}

// GetAll returns all message channel configs.
func (h *ConfigMessageHandler) GetAll(c *gin.Context) {
	items, err := h.svc.GetAll()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// GetByChannel returns a single message channel config.
func (h *ConfigMessageHandler) GetByChannel(c *gin.Context) {
	channel := c.Param("channel")
	if channel == "" {
		response.BadRequest(c, "channel is required")
		return
	}

	item, err := h.svc.GetByChannel(channel)
	if err != nil {
		response.NotFound(c, "message channel config not found")
		return
	}
	response.Success(c, item)
}

// Update updates a message channel config.
func (h *ConfigMessageHandler) Update(c *gin.Context) {
	channel := c.Param("channel")
	if channel == "" {
		response.BadRequest(c, "channel is required")
		return
	}

	var req service.UpdateMessageConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Update(channel, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// TestSend tests sending via a message channel.
func (h *ConfigMessageHandler) TestSend(c *gin.Context) {
	channel := c.Param("channel")
	if channel == "" {
		response.BadRequest(c, "channel is required")
		return
	}

	if err := h.svc.TestSend(channel); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "test message sent successfully")
}

// GetEnabled returns all enabled message channels.
func (h *ConfigMessageHandler) GetEnabled(c *gin.Context) {
	items, err := h.svc.GetEnabledChannels()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// CreateTemplate creates a new message template.
// POST /config/messages/templates
func (h *ConfigMessageHandler) CreateTemplate(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Channel     string `json:"channel" binding:"required"`
		Subject     string `json:"subject"`
		Content     string `json:"content" binding:"required"`
		Description string `json:"description"`
		Variables   string `json:"variables"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	template, err := h.svc.CreateTemplate(req.Name, req.Channel, req.Subject, req.Content, req.Description, req.Variables)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, template)
}

// UpdateTemplate updates an existing message template.
// PUT /config/messages/templates/:id
func (h *ConfigMessageHandler) UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "template id is required")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Subject     string `json:"subject"`
		Content     string `json:"content"`
		Description string `json:"description"`
		Variables   string `json:"variables"`
		Status      *int   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateTemplate(id, req.Name, req.Subject, req.Content, req.Description, req.Variables, req.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "template updated")
}

// DeleteTemplate deletes a message template.
// DELETE /config/messages/templates/:id
func (h *ConfigMessageHandler) DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "template id is required")
		return
	}

	if err := h.svc.DeleteTemplate(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "template deleted")
}

// GetTemplateDesc returns template description and variables.
// GET /config/messages/templates/:id/desc
func (h *ConfigMessageHandler) GetTemplateDesc(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "template id is required")
		return
	}

	desc, err := h.svc.GetTemplateDesc(id)
	if err != nil {
		response.NotFound(c, "template not found")
		return
	}
	response.Success(c, desc)
}

// SetSmsTemplate sets an SMS template.
// POST /config/messages/sms-template
func (h *ConfigMessageHandler) SetSmsTemplate(c *gin.Context) {
	var req struct {
		TemplateID string `json:"template_id" binding:"required"`
		SmsContent string `json:"sms_content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.SetSmsTemplate(req.TemplateID, req.SmsContent); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "SMS template set")
}

// BeforeSendMessageCheck checks if a message can be sent.
// POST /config/messages/check
func (h *ConfigMessageHandler) BeforeSendMessageCheck(c *gin.Context) {
	var req struct {
		Channel    string `json:"channel" binding:"required"`
		TemplateID string `json:"template_id"`
		UserID     uint   `json:"user_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.BeforeSendMessageCheck(req.Channel, req.TemplateID, req.UserID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

// ==================== 新增缺失方法 ====================

// GetMobileConfig returns mobile/SMS specific configuration (admin).
// GET /admin/config/message/mobile
func (h *ConfigMessageHandler) GetMobileConfig(c *gin.Context) {
	item, err := h.svc.GetByChannel("sms")
	if err != nil {
		// Return empty config if not found
		response.Success(c, map[string]interface{}{
			"channel":   "sms",
			"enabled":   false,
			"provider":  "",
			"config":    map[string]interface{}{},
			"signature": "",
		})
		return
	}

	response.Success(c, map[string]interface{}{
		"channel":      item.Channel,
		"enabled":      item.IsEnabled,
		"provider":     item.Provider,
		"config":       item.Config,
		"signature":    item.Signature,
		"rate_limit":   item.RateLimit,
		"daily_limit":  item.DailyLimit,
		"test_address": item.TestAddress,
		"status":       item.Status,
	})
}

// UpdateTemplateStatus updates the status of a message template (admin).
// PUT /admin/config/message/templates/:id/status
func (h *ConfigMessageHandler) UpdateTemplateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "template id is required")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Status != 0 && req.Status != 1 {
		response.BadRequest(c, "status must be 0 (disabled) or 1 (enabled)")
		return
	}

	if err := h.svc.UpdateTemplate(id, "", "", "", "", "", &req.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "template status updated")
}

// TestTemplate tests a message template by sending a test message (admin).
// POST /admin/config/message/test
func (h *ConfigMessageHandler) TestTemplate(c *gin.Context) {
	var req struct {
		TemplateID string `json:"template_id" binding:"required"`
		Channel    string `json:"channel"`
		Recipient  string `json:"recipient"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get template info
	desc, err := h.svc.GetTemplateDesc(req.TemplateID)
	if err != nil {
		response.BadRequest(c, "template not found")
		return
	}

	channel := req.Channel
	if channel == "" {
		if ch, ok := desc["channel"].(string); ok {
			channel = ch
		} else {
			channel = "site"
		}
	}

	// Use the test send functionality
	if err := h.svc.TestSend(channel); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"template_id": req.TemplateID,
		"channel":     channel,
		"recipient":   req.Recipient,
		"status":      "sent",
		"content":     desc["content"],
	})
}
