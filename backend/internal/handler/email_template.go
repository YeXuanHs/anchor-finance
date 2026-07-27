package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type EmailTemplateHandler struct {
	emailTemplateSvc *service.EmailTemplateService
	log              *logger.Logger
}

func NewEmailTemplateHandler(emailTemplateSvc *service.EmailTemplateService, log *logger.Logger) *EmailTemplateHandler {
	return &EmailTemplateHandler{emailTemplateSvc: emailTemplateSvc, log: log}
}

// List returns paginated email templates.
func (h *EmailTemplateHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	templateType := c.Query("type")
	keyword := c.Query("keyword")

	templates, total, err := h.emailTemplateSvc.List(page, pageSize, templateType, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, templates, total, page, pageSize)
}

// GetDetail returns a single email template by ID.
func (h *EmailTemplateHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	tmpl, err := h.emailTemplateSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "template not found")
		return
	}
	response.Success(c, tmpl)
}

// Create creates a new email template.
func (h *EmailTemplateHandler) Create(c *gin.Context) {
	var req service.CreateEmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tmpl, err := h.emailTemplateSvc.Create(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, tmpl)
}

// Update updates an email template.
func (h *EmailTemplateHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	var req service.UpdateEmailTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	tmpl, err := h.emailTemplateSvc.Update(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, tmpl)
}

// Delete deletes an email template.
func (h *EmailTemplateHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	if err := h.emailTemplateSvc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "template deleted")
}

// Preview renders a template with provided data.
func (h *EmailTemplateHandler) Preview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	var req struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	subject, body, err := h.emailTemplateSvc.Preview(uint(id), req.Data)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"subject": subject,
		"body":    body,
	})
}

// SendTest sends a test email/sms/notice.
func (h *EmailTemplateHandler) SendTest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	var req service.SendTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.emailTemplateSvc.SendTest(uint(id), req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "test message sent")
}

// GetSendLogs returns send logs with pagination.
func (h *EmailTemplateHandler) GetSendLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	logType := c.Query("type")

	var templateID *uint
	if tid := c.Query("template_id"); tid != "" {
		v, _ := strconv.ParseUint(tid, 10, 64)
		id := uint(v)
		templateID = &id
	}

	logs, total, err := h.emailTemplateSvc.GetSendLogs(page, pageSize, templateID, logType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}
