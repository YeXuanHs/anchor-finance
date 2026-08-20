package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"math"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/email"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/sms"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type EmailTemplateService struct {
	db       *gorm.DB
	log      *logger.Logger
	emailSnd *email.Sender
	smsSnd   *sms.Sender
}

func NewEmailTemplateService(db *gorm.DB, log *logger.Logger) *EmailTemplateService {
	return &EmailTemplateService{
		db:       db,
		log:      log,
		emailSnd: email.NewSender(db),
		smsSnd:   sms.NewSender(db),
	}
}

type CreateEmailTemplateRequest struct {
	Code      string                 `json:"code" binding:"required,max=64"`
	Name      string                 `json:"name" binding:"required,max=128"`
	Subject   string                 `json:"subject" binding:"required,max=256"`
	Body      string                 `json:"body" binding:"required"`
	Type      string                 `json:"type" binding:"required,oneof=email sms notice"`
	Variables []TemplateVariable     `json:"variables"`
	Format    string                 `json:"format" binding:"omitempty,oneof=html plain"`
	Language  string                 `json:"language"`
}

type UpdateEmailTemplateRequest struct {
	Name      string             `json:"name" binding:"omitempty,max=128"`
	Subject   string             `json:"subject" binding:"omitempty,max=256"`
	Body      string             `json:"body"`
	Type      string             `json:"type" binding:"omitempty,oneof=email sms notice"`
	Variables []TemplateVariable `json:"variables"`
	Format    string             `json:"format" binding:"omitempty,oneof=html plain"`
	Language  string             `json:"language"`
	Status    *int16             `json:"status"`
}

type TemplateVariable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Example     string `json:"example"`
	Required    bool   `json:"required"`
}

type SendTestRequest struct {
	Recipient string                 `json:"recipient" binding:"required"`
	Data      map[string]interface{} `json:"data"`
}

// List returns paginated email templates.
func (s *EmailTemplateService) List(page, pageSize int, templateType string, keyword string) ([]model.EmailTemplate, int64, error) {
	var templates []model.EmailTemplate
	var total int64

	query := s.db.Model(&model.EmailTemplate{})
	if templateType != "" {
		query = query.Where("type = ?", templateType)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := emailTemplatePaginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&templates).Error; err != nil {
		return nil, 0, err
	}
	return templates, total, nil
}

// GetByID returns a single template by ID.
func (s *EmailTemplateService) GetByID(id uint) (*model.EmailTemplate, error) {
	var tmpl model.EmailTemplate
	if err := s.db.First(&tmpl, id).Error; err != nil {
		return nil, err
	}
	return &tmpl, nil
}

// GetByCode returns a single template by code.
func (s *EmailTemplateService) GetByCode(code string) (*model.EmailTemplate, error) {
	var tmpl model.EmailTemplate
	if err := s.db.Where("code = ?", code).First(&tmpl).Error; err != nil {
		return nil, err
	}
	return &tmpl, nil
}

// Create creates a new email template.
func (s *EmailTemplateService) Create(req CreateEmailTemplateRequest) (*model.EmailTemplate, error) {
	var existing model.EmailTemplate
	if err := s.db.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		return nil, errors.New("template code already exists")
	}

	format := req.Format
	if format == "" {
		format = "html"
	}
	language := req.Language
	if language == "" {
		language = "zh-CN"
	}

	tmpl := &model.EmailTemplate{
		Code:      req.Code,
		Name:      req.Name,
		Subject:   req.Subject,
		Body:      req.Body,
		Type:      req.Type,
		Variables: toJSON(req.Variables),
		Format:    format,
		Language:  language,
		Status:    1,
	}

	if err := s.db.Create(tmpl).Error; err != nil {
		return nil, err
	}

	s.log.Infof("email template created: %s (code=%s)", tmpl.Name, tmpl.Code)
	return tmpl, nil
}

// Update updates an email template.
func (s *EmailTemplateService) Update(id uint, req UpdateEmailTemplateRequest) (*model.EmailTemplate, error) {
	var tmpl model.EmailTemplate
	if err := s.db.First(&tmpl, id).Error; err != nil {
		return nil, err
	}
	if tmpl.IsSystem {
		return nil, errors.New("system template cannot be modified")
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Subject != "" {
		updates["subject"] = req.Subject
	}
	if req.Body != "" {
		updates["body"] = req.Body
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Variables != nil {
		updates["variables"] = toJSON(req.Variables)
	}
	if req.Format != "" {
		updates["format"] = req.Format
	}
	if req.Language != "" {
		updates["language"] = req.Language
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := s.db.Model(&tmpl).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	s.log.Infof("email template updated: %s (id=%d)", tmpl.Name, tmpl.ID)
	return s.GetByID(id)
}

// Delete deletes an email template (system templates excluded).
func (s *EmailTemplateService) Delete(id uint) error {
	var tmpl model.EmailTemplate
	if err := s.db.First(&tmpl, id).Error; err != nil {
		return err
	}
	if tmpl.IsSystem {
		return errors.New("system template cannot be deleted")
	}
	return s.db.Delete(&tmpl).Error
}

// Preview renders a template with provided data and returns the result.
func (s *EmailTemplateService) Preview(id uint, data map[string]interface{}) (string, string, error) {
	tmpl, err := s.GetByID(id)
	if err != nil {
		return "", "", err
	}

	subject, err := s.renderString(tmpl.Subject, data)
	if err != nil {
		return "", "", fmt.Errorf("subject render error: %w", err)
	}
	body, err := s.renderString(tmpl.Body, data)
	if err != nil {
		return "", "", fmt.Errorf("body render error: %w", err)
	}
	return subject, body, nil
}

// SendTest sends a test email/sms/notice using the template.
func (s *EmailTemplateService) SendTest(id uint, req SendTestRequest) error {
	tmpl, err := s.GetByID(id)
	if err != nil {
		return err
	}
	if tmpl.Status != 1 {
		return errors.New("template is disabled")
	}

	subject, err := s.renderString(tmpl.Subject, req.Data)
	if err != nil {
		return fmt.Errorf("subject render error: %w", err)
	}
	content, err := s.renderString(tmpl.Body, req.Data)
	if err != nil {
		return fmt.Errorf("body render error: %w", err)
	}

	logEntry := &model.EmailTemplateLog{
		TemplateID: tmpl.ID,
		Recipient:  req.Recipient,
		Subject:    subject,
		Content:    content,
		Type:       tmpl.Type,
		Status:     1,
	}
	if err := s.db.Create(logEntry).Error; err != nil {
		return err
	}

	// 按类型分发（此处为占位，实际实现需接入邮件/短信/站内信SDK）
	if err := s.dispatch(tmpl.Type, req.Recipient, subject, content); err != nil {
		s.db.Model(logEntry).Updates(map[string]interface{}{
			"status": 3,
			"error":  err.Error(),
		})
		s.log.Errorf("template test send failed (template=%d, type=%s): %v", id, tmpl.Type, err)
		return err
	}

	now := time.Now()
	s.db.Model(logEntry).Updates(map[string]interface{}{
		"status":  2,
		"sent_at": &now,
	})

	// 更新模板发送计数
	s.db.Model(tmpl).Updates(map[string]interface{}{
		"send_count":   tmpl.SendCount + 1,
		"last_sent_at": &now,
	})

	s.log.Infof("template test sent (template=%d, type=%s, to=%s)", id, tmpl.Type, req.Recipient)
	return nil
}

// GetSendLogs returns send logs for a template.
func (s *EmailTemplateService) GetSendLogs(page, pageSize int, templateID *uint, logType string) ([]model.EmailTemplateLog, int64, error) {
	var logs []model.EmailTemplateLog
	var total int64

	query := s.db.Model(&model.EmailTemplateLog{})
	if templateID != nil {
		query = query.Where("template_id = ?", *templateID)
	}
	if logType != "" {
		query = query.Where("type = ?", logType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := emailTemplatePaginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// renderString parses and executes a Go template string.
func (s *EmailTemplateService) renderString(tmplStr string, data map[string]interface{}) (string, error) {
	t, err := template.New("").Parse(tmplStr)
	if err != nil {
		return tmplStr, nil
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return tmplStr, nil
	}
	return buf.String(), nil
}

// dispatch dispatches a message via the specified type.
func (s *EmailTemplateService) dispatch(msgType, to, subject, content string) error {
	switch msgType {
	case "email":
		if s.emailSnd == nil {
			return errors.New("email sender not configured")
		}
		return s.emailSnd.Send(to, subject, content)
	case "sms":
		if s.smsSnd == nil {
			return errors.New("sms sender not configured")
		}
		return s.smsSnd.Send(to, content)
	case "notice":
		// In-app notice: already logged by caller
		return nil
	default:
		return fmt.Errorf("unsupported message type: %s", msgType)
	}
}

func toJSON(v interface{}) datatypes.JSON {
	if v == nil {
		return nil
	}
	data, _ := json.Marshal(v)
	return datatypes.JSON(data)
}

func emailTemplatePaginate(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	pageSize = int(math.Min(float64(pageSize), 100))
	return (page - 1) * pageSize, pageSize
}
