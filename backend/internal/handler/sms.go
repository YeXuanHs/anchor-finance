package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type SMSHandler struct {
	smsSvc *service.SMSService
	log    *logger.Logger
}

func NewSMSHandler(smsSvc *service.SMSService, log *logger.Logger) *SMSHandler {
	return &SMSHandler{smsSvc: smsSvc, log: log}
}

// ─── Operator Detection ───

// DetectOperator detects the mobile operator from a phone number.
func (h *SMSHandler) DetectOperator(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		response.BadRequest(c, "phone is required")
		return
	}

	operator := h.smsSvc.DetectOperator(phone)
	response.Success(c, gin.H{"phone": phone, "operator": operator})
}

// ValidatePhone validates a phone number format.
func (h *SMSHandler) ValidatePhone(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		response.BadRequest(c, "phone is required")
		return
	}

	valid := h.smsSvc.ValidatePhone(phone)
	response.Success(c, gin.H{"phone": phone, "valid": valid})
}

// ─── Template Management ───

// GetTemplates returns all SMS templates (admin).
func (h *SMSHandler) GetTemplates(c *gin.Context) {
	templates, err := h.smsSvc.GetTemplates()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, templates)
}

// GetTemplate returns a single template by ID (admin).
func (h *SMSHandler) GetTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	tpl, err := h.smsSvc.GetTemplateByID(uint(id))
	if err != nil {
		response.NotFound(c, "template not found")
		return
	}
	response.Success(c, tpl)
}

// CreateTemplate creates a new SMS template (admin).
func (h *SMSHandler) CreateTemplate(c *gin.Context) {
	var req struct {
		Name    string   `json:"name" binding:"required"`
		Code    string   `json:"code" binding:"required"`
		Content string   `json:"content" binding:"required"`
		Params  []string `json:"params"`
		Type    string   `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tpl := &model.SMSTemplate{
		Name:    req.Name,
		Code:    req.Code,
		Content: req.Content,
		Type:    req.Type,
		Enabled: true,
	}

	if err := h.smsSvc.CreateTemplate(tpl); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, tpl)
}

// UpdateTemplate updates an existing SMS template (admin).
func (h *SMSHandler) UpdateTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	var req struct {
		Name    string `json:"name"`
		Content string `json:"content"`
		Type    string `json:"type"`
		Enabled *bool  `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Enabled != nil {
		updates["enabled"] = *req.Enabled
	}

	if err := h.smsSvc.UpdateTemplate(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "template updated")
}

// DeleteTemplate deletes an SMS template (admin).
func (h *SMSHandler) DeleteTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	if err := h.smsSvc.DeleteTemplate(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "template deleted")
}

// ─── Enhanced Sending ───

// SendSMS sends an SMS using a template (admin).
func (h *SMSHandler) SendSMS(c *gin.Context) {
	var req struct {
		Phone   string            `json:"phone" binding:"required"`
		TplCode string            `json:"tpl_code" binding:"required"`
		Params  map[string]string `json:"params"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	logEntry, err := h.smsSvc.SendSMSSend(req.Phone, req.TplCode, req.Params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, logEntry)
}

// SendBatchSMS sends SMS to multiple phones (admin).
func (h *SMSHandler) SendBatchSMS(c *gin.Context) {
	var req struct {
		Phones  []string          `json:"phones" binding:"required"`
		TplCode string            `json:"tpl_code" binding:"required"`
		Params  map[string]string `json:"params"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	h.smsSvc.SendBatchSMS(req.Phones, req.TplCode, req.Params, nil)
	response.SuccessMsg(c, "batch SMS sending started")
}

// SendMarketingSMS sends marketing SMS to a target group (admin).
func (h *SMSHandler) SendMarketingSMS(c *gin.Context) {
	var req struct {
		TargetGroup string            `json:"target_group" binding:"required"`
		TplCode     string            `json:"tpl_code" binding:"required"`
		Params      map[string]string `json:"params"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	batch, err := h.smsSvc.SendMarketingSMS(req.TargetGroup, req.TplCode, req.Params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, batch)
}

// SendVerifyCode sends a verification code (user-facing).
func (h *SMSHandler) SendVerifyCode(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	logEntry, err := h.smsSvc.SendVerifyCode(req.Phone)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "verification code sent", "log_id": logEntry.ID})
}

// ─── Logging ───

// GetSMSLogs returns paginated SMS logs (admin).
func (h *SMSHandler) GetSMSLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.smsSvc.GetSMSLogs(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// GetSMSLog returns a single SMS log by ID (admin).
func (h *SMSHandler) GetSMSLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid log id")
		return
	}

	logEntry, err := h.smsSvc.GetSMSLogByID(uint(id))
	if err != nil {
		response.NotFound(c, "log not found")
		return
	}
	response.Success(c, logEntry)
}

// GetSMSLogByPhone returns SMS logs for a phone number (admin).
func (h *SMSHandler) GetSMSLogByPhone(c *gin.Context) {
	phone := c.Param("phone")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.smsSvc.GetSMSLogByPhone(phone, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// GetSMSLogByUser returns SMS logs for a user (admin).
func (h *SMSHandler) GetSMSLogByUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.smsSvc.GetSMSLogByUser(uint(userID), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// ─── Statistics ───

// GetSMSStats returns SMS sending statistics (admin).
func (h *SMSHandler) GetSMSStats(c *gin.Context) {
	period := c.DefaultQuery("period", "all")

	stats, err := h.smsSvc.GetSMSStats(period)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// GetOperatorStats returns SMS statistics by operator (admin).
func (h *SMSHandler) GetOperatorStats(c *gin.Context) {
	stats, err := h.smsSvc.GetOperatorStats()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// ─── Batch Operations ───

// CreateBatch creates a new SMS batch job (admin).
func (h *SMSHandler) CreateBatch(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		TemplateID  uint   `json:"template_id" binding:"required"`
		TargetGroup string `json:"target_group" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	batch, err := h.smsSvc.CreateBatch(req.Name, req.TemplateID, req.TargetGroup)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, batch)
}

// GetBatches returns all SMS batch jobs (admin).
func (h *SMSHandler) GetBatches(c *gin.Context) {
	batches, err := h.smsSvc.GetBatches()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, batches)
}

// ExecuteBatch executes an SMS batch job (admin).
func (h *SMSHandler) ExecuteBatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid batch id")
		return
	}

	if err := h.smsSvc.ExecuteBatch(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "batch execution started")
}
