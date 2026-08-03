package handler

import (
	"strconv"
	"time"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type EmailTemplateHandler struct {
	emailTemplateSvc *service.EmailTemplateService
	log              *logger.Logger
	db               *gorm.DB
}

func NewEmailTemplateHandler(emailTemplateSvc *service.EmailTemplateService, log *logger.Logger, db *gorm.DB) *EmailTemplateHandler {
	return &EmailTemplateHandler{emailTemplateSvc: emailTemplateSvc, log: log, db: db}
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
	response.SuccessMsg(c, "删除成功")
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
	response.SuccessMsg(c, "发送成功")
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

// SwitchOperator 邮件运营商切换
// POST /admin/email-template/operator-switch
func (h *EmailTemplateHandler) SwitchOperator(c *gin.Context) {
	var req struct {
		EmailOperator string `json:"email_operator"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	operator := req.EmailOperator
	if operator != "" {
		// 转小写
		for i := range operator {
			if operator[i] >= 'A' && operator[i] <= 'Z' {
				operator[i] = operator[i] + 32
			}
		}
	}

	var count int64
	h.db.Table("system_configs").Where("setting = ?", "email_operator").Count(&count)
	if count > 0 {
		h.db.Table("system_configs").Where("setting = ?", "email_operator").Update("value", operator)
	} else {
		h.db.Table("system_configs").Create(&map[string]interface{}{
			"setting": "email_operator",
			"value":   operator,
		})
	}

	response.SuccessMsg(c, "修改成功")
}

// GetLanguages 获取模板多语言列表
// GET /admin/email-template/languages
func (h *EmailTemplateHandler) GetLanguages(c *gin.Context) {
	// 所有支持的语言
	allLangs := []string{
		"zh-CN", "en-US", "ja-JP", "ko-KR", "zh-TW",
		"de-DE", "fr-FR", "es-ES", "pt-BR", "ru-RU",
	}

	// 已使用的语言
	type LangUsed struct {
		Language string `json:"language"`
	}
	var usedLangs []LangUsed
	h.db.Table("email_templates").
		Select("DISTINCT language").
		Where("language != '' AND language IS NOT NULL").
		Scan(&usedLangs)

	used := make([]string, 0)
	for _, l := range usedLangs {
		used = append(used, l.Language)
	}

	response.Success(c, gin.H{
		"langs":     allLangs,
		"lang_used": used,
	})
}

// SaveLanguages 保存模板多语言
// POST /admin/email-template/languages
func (h *EmailTemplateHandler) SaveLanguages(c *gin.Context) {
	var req struct {
		Language string `json:"language" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 查找所有默认语言（language为空）的模板
	type DefaultTemplate struct {
		ID      uint   `json:"id"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}
	var defaults []DefaultTemplate
	h.db.Table("email_templates").Where("language = '' OR language IS NULL").Find(&defaults)

	// 为每个默认模板创建指定语言版本
	for _, tmpl := range defaults {
		// 检查是否已存在该语言版本
		var existCount int64
		h.db.Table("email_templates").Where("name = ? AND language = ?", tmpl.Name, req.Language).Count(&existCount)
		if existCount == 0 {
			h.db.Table("email_templates").Create(&map[string]interface{}{
				"name":     tmpl.Name,
				"type":     tmpl.Type,
				"subject":  tmpl.Subject,
				"body":     tmpl.Body,
				"language": req.Language,
				"is_system": false,
				"status":   1,
			})
		}
	}

	response.SuccessMsg(c, "添加成功")
}

// ToggleEnabled 启用/禁用邮件模板
// POST /admin/email-template/toggle
func (h *EmailTemplateHandler) ToggleEnabled(c *gin.Context) {
	var req struct {
		ID       uint `json:"id" binding:"required"`
		Disabled int  `json:"disabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Disabled != 0 && req.Disabled != 1 {
		req.Disabled = 0
	}

	h.db.Table("email_templates").Where("id = ?", req.ID).Updates(map[string]interface{}{
		"status":      req.Disabled,
		"updated_at":  time.Now(),
	})

	response.SuccessMsg(c, "操作成功")
}

// EditTemplate 编辑邮件模板（含模板变量、多语言编辑）
// GET /admin/email-template/edit/:id
func (h *EmailTemplateHandler) EditTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	type TemplateDetail struct {
		ID       uint   `json:"id"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		Subject  string `json:"subject"`
		Body     string `json:"body"`
		Status   int16  `json:"status"`
		Language string `json:"language"`
	}
	var tmpl TemplateDetail
	h.db.Table("email_templates").Where("id = ?", id).Scan(&tmpl)
	if tmpl.ID == 0 {
		response.NotFound(c, "template not found")
		return
	}

	// 获取同名的所有语言版本
	var childVersions []TemplateDetail
	h.db.Table("email_templates").Where("name = ? AND id != ?", tmpl.Name, id).Scan(&childVersions)

	// 模板变量定义
	baseArgs := []gin.H{
		{"label": "{SYSTEM_COMPANYNAME}", "name": "公司名称"},
		{"label": "{COMPANY_DOMAIN}", "name": "公司域名"},
		{"label": "{TEMPLATE_DATE}", "name": "模板日期"},
		{"label": "{TEMPLATE_TIME}", "name": "模板时间"},
		{"label": "{CODE}", "name": "验证码"},
		{"label": "{SEND_TIME}", "name": "发送时间"},
		{"label": "{SYSTEM_URL}", "name": "系统URL"},
		{"label": "{SYSTEM_WEB_URL}", "name": "网站URL"},
	}

	clientArgs := []gin.H{
		{"label": "{CLIENT_ID}", "name": "客户ID"},
		{"label": "{USERNAME}", "name": "用户名"},
		{"label": "{ACCOUNT_EMAIL}", "name": "邮箱"},
		{"label": "{CLIENT_SIGNUP_DATE}", "name": "注册时间"},
		{"label": "{CLIENT_STATUS}", "name": "客户状态"},
		{"label": "{CLIENT_GROUP_NAME}", "name": "客户组"},
	}

	productArgs := []gin.H{
		{"label": "{PRODUCT_NAME}", "name": "产品名称"},
		{"label": "{HOSTNAME}", "name": "主机名"},
		{"label": "{PRODUCT_MAINIP}", "name": "主IP"},
		{"label": "{PRODUCT_FIRST_TIME}", "name": "开通时间"},
		{"label": "{PRODUCT_END_TIME}", "name": "到期时间"},
		{"label": "{ORDER_ID}", "name": "订单ID"},
		{"label": "{ORDER_TOTAL_FEE}", "name": "订单金额"},
	}

	// 根据模板类型返回对应的变量
	var combine []gin.H
	if tmpl.Type == "product" || tmpl.Type == "invoice" {
		combine = append(combine, gin.H{"label": "args_product", "name": "产品/服务相关", "list": productArgs})
	}
	combine = append(combine, gin.H{"label": "args_clients", "name": "客户相关", "list": clientArgs})
	combine = append(combine, gin.H{"label": "args_base", "name": "其他", "list": baseArgs})

	response.Success(c, gin.H{
		"emailtemplate": tmpl,
		"child":         childVersions,
		"combine":       combine,
	})
}

// SaveTemplate 保存编辑的模板
// POST /admin/email-template/save
func (h *EmailTemplateHandler) SaveTemplate(c *gin.Context) {
	var req struct {
		IDs      map[string]struct {
			Subject string `json:"subject"`
			Message string `json:"message"`
		} `json:"ids"`
		FromName        string `json:"fromname"`
		FromEmail       string `json:"fromemail"`
		CopyTo          string `json:"copyto"`
		BlindCopyTo     string `json:"blind_copy_to"`
		Plaintext       bool   `json:"plaintext"`
		Disabled        int    `json:"disabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	for idStr, content := range req.IDs {
		id, _ := strconv.ParseUint(idStr, 10, 64)
		if id == 0 {
			continue
		}

		updates := map[string]interface{}{
			"subject":   content.Subject,
			"body":      content.Message,
			"from_name": req.FromName,
			"from_email": req.FromEmail,
			"copy_to":   req.CopyTo,
			"blind_copy_to": req.BlindCopyTo,
			"format":    map[bool]string{true: "plain", false: "html"}[req.Plaintext],
			"status":    req.Disabled,
			"updated_at": time.Now(),
		}
		h.db.Table("email_templates").Where("id = ?", id).Updates(updates)
	}

	response.SuccessMsg(c, "编辑成功")
}
