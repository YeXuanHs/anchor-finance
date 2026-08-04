package service

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// EmailEnhancedService 邮件增强服务
type EmailEnhancedService struct {
	db *gorm.DB
}

// NewEmailEnhancedService 创建邮件增强服务
func NewEmailEnhancedService(db *gorm.DB) *EmailEnhancedService {
	return &EmailEnhancedService{db: db}
}

// EmailTemplate 邮件模板
type EmailTemplate struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"size:128;not null"`
	Code      string         `gorm:"size:64;uniqueIndex"`
	Subject   string         `gorm:"size:256;not null"`
	Body      string         `gorm:"type:text;not null"`
	Params    datatypes.JSON `gorm:"type:json"`
	Type      string         `gorm:"size:32"` // system/marketing/notification
	Enabled   bool           `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// EmailLog 邮件日志
type EmailLog struct {
	ID         uint           `gorm:"primaryKey"`
	To         string         `gorm:"size:256;index"`
	Subject    string         `gorm:"size:256"`
	Body       string         `gorm:"type:text"`
	TemplateID *uint          `gorm:"index"`
	Params     datatypes.JSON `gorm:"type:json"`
	Status     string         `gorm:"size:32"` // pending/sent/failed/bounced
	Response   string         `gorm:"type:text"`
	UserID     *uint          `gorm:"index"`
	BatchID    *uint          `gorm:"index"`
	SentAt     *time.Time
	CreatedAt  time.Time
}

// EmailBatch 邮件批量发送
type EmailBatch struct {
	ID          uint       `gorm:"primaryKey"`
	Name        string     `gorm:"size:128"`
	TemplateID  uint       `gorm:"index"`
	TargetGroup string     `gorm:"size:64"` // all/new/active/vip/custom
	TotalCount  int
	SentCount   int
	FailedCount int
	Status      string     `gorm:"size:32"` // pending/sending/completed/failed
	CreatedAt   time.Time
	CompletedAt *time.Time
}

// EmailSupport 支持邮箱配置
type EmailSupport struct {
	SupportEmail     string `json:"support_email"`
	SalesEmail       string `json:"sales_email"`
	AbuseEmail       string `json:"abuse_email"`
	BillingEmail     string `json:"billing_email"`
	NoReplyEmail     string `json:"no_reply_email"`
	NoReplyName      string `json:"no_reply_name"`
}

// GetTemplates 获取邮件模板列表
func (s *EmailEnhancedService) GetTemplates() ([]EmailTemplate, error) {
	var templates []EmailTemplate
	err := s.db.Order("id").Find(&templates).Error
	return templates, err
}

// GetTemplateByCode 根据代码获取模板
func (s *EmailEnhancedService) GetTemplateByCode(code string) (*EmailTemplate, error) {
	var tpl EmailTemplate
	err := s.db.Where("code = ?", code).First(&tpl).Error
	if err != nil {
		return nil, err
	}
	return &tpl, nil
}

// CreateTemplate 创建邮件模板
func (s *EmailEnhancedService) CreateTemplate(tpl *EmailTemplate) error {
	return s.db.Create(tpl).Error
}

// UpdateTemplate 更新邮件模板
func (s *EmailEnhancedService) UpdateTemplate(id uint, updates map[string]interface{}) error {
	return s.db.Model(&EmailTemplate{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteTemplate 删除邮件模板
func (s *EmailEnhancedService) DeleteTemplate(id uint) error {
	return s.db.Delete(&EmailTemplate{}, id).Error
}

// ReplaceTemplateParams 替换模板参数
func (s *EmailEnhancedService) ReplaceTemplateParams(tplCode string, params map[string]string) (subject, body string, err error) {
	tpl, err := s.GetTemplateByCode(tplCode)
	if err != nil {
		return "", "", err
	}

	subject = tpl.Subject
	body = tpl.Body

	// 替换 {param} 格式的占位符
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		subject = strings.ReplaceAll(subject, placeholder, value)
		body = strings.ReplaceAll(body, placeholder, value)
	}

	// 替换系统变量
	systemVars := map[string]string{
		"{site_name}": "锚点财务",
		"{site_url}":  "https://example.com",
		"{year}":      fmt.Sprintf("%d", time.Now().Year()),
		"{date}":      time.Now().Format("2006-01-02"),
		"{time}":      time.Now().Format("15:04:05"),
	}
	for key, value := range systemVars {
		subject = strings.ReplaceAll(subject, key, value)
		body = strings.ReplaceAll(body, key, value)
	}

	return subject, body, nil
}

// SendEmailWithTemplate 使用模板发送邮件
func (s *EmailEnhancedService) SendEmailWithTemplate(to string, tplCode string, params map[string]string) error {
	subject, body, err := s.ReplaceTemplateParams(tplCode, params)
	if err != nil {
		return err
	}

	return s.SendEmail(to, subject, body, nil)
}

// SendEmail 发送邮件
func (s *EmailEnhancedService) SendEmail(to, subject, body string, userID *uint) error {
	// 获取SMTP配置
	smtpConfig, err := s.getSMTPConfig()
	if err != nil {
		return fmt.Errorf("failed to get SMTP config: %w", err)
	}

	if !smtpConfig.Enabled {
		return fmt.Errorf("SMTP is disabled")
	}

	// 记录日志
	log := &EmailLog{
		To:      to,
		Subject: subject,
		Body:    body,
		UserID:  userID,
		Status:  "pending",
	}
	s.db.Create(log)

	// 实际发送（调用已有的邮件发送功能）
	err = s.sendViaSMTP(smtpConfig, to, subject, body)
	if err != nil {
		s.db.Model(log).Updates(map[string]interface{}{
			"status":   "failed",
			"response": err.Error(),
		})
		return err
	}

	now := time.Now()
	s.db.Model(log).Updates(map[string]interface{}{
		"status":   "sent",
		"sent_at":  &now,
		"response": "OK",
	})

	return nil
}

// SendBatchEmail 批量发送邮件
func (s *EmailEnhancedService) SendBatchEmail(templateID uint, targetGroup string, customParams map[string]string) (*EmailBatch, error) {
	batch := &EmailBatch{
		TemplateID:  templateID,
		TargetGroup: targetGroup,
		Status:      "pending",
	}
	s.db.Create(batch)

	// 获取目标用户
	users, err := s.getTargetUsers(targetGroup)
	if err != nil {
		s.db.Model(batch).Update("status", "failed")
		return nil, err
	}

	batch.TotalCount = len(users)
	s.db.Model(batch).Update("total_count", len(users))

	// 异步发送
	go s.executeBatch(batch, users, customParams)

	return batch, nil
}

// executeBatch 执行批量发送
func (s *EmailEnhancedService) executeBatch(batch *EmailBatch, users []struct{ Email string; ID uint }, params map[string]string) {
	s.db.Model(batch).Update("status", "sending")

	tpl, err := s.GetTemplateByCode(fmt.Sprintf("%d", batch.TemplateID))
	if err != nil {
		s.db.Model(batch).Update("status", "failed")
		return
	}

	sentCount := 0
	failedCount := 0

	for _, user := range users {
		// 为每个用户替换参数
		userParams := make(map[string]string)
		for k, v := range params {
			userParams[k] = v
		}
		userParams["email"] = user.Email

		subject, body, _ := s.ReplaceTemplateParams(tpl.Code, userParams)

		err := s.SendEmail(user.Email, subject, body, &user.ID)
		if err != nil {
			failedCount++
		} else {
			sentCount++
		}

		// 更新进度
		s.db.Model(batch).Updates(map[string]interface{}{
			"sent_count":   sentCount,
			"failed_count": failedCount,
		})
	}

	now := time.Now()
	s.db.Model(batch).Updates(map[string]interface{}{
		"status":       "completed",
		"completed_at": &now,
	})
}

// getTargetUsers 获取目标用户
func (s *EmailEnhancedService) getTargetUsers(group string) ([]struct{ Email string; ID uint }, error) {
	var users []struct{ Email string; ID uint }

	query := s.db.Table("users").Select("id, email").Where("email != ''")

	switch group {
	case "new":
		// 最近30天注册
		query = query.Where("created_at > ?", time.Now().AddDate(0, 0, -30))
	case "active":
		// 最近30天活跃
		query = query.Where("last_login_at > ?", time.Now().AddDate(0, 0, -30))
	case "vip":
		// VIP用户组
		query = query.Joins("LEFT JOIN user_groups ON user_groups.user_id = users.id").
			Where("user_groups.group_id IN (SELECT id FROM client_groups WHERE name LIKE '%VIP%')")
	case "all":
		// 所有用户
	}

	err := query.Find(&users).Error
	return users, err
}

// getSMTPConfig 获取SMTP配置
func (s *EmailEnhancedService) getSMTPConfig() (*smtpConfig, error) {
	var setting struct {
		Value string
	}
	err := s.db.Table("system_configs").Where("key = ?", "email_config").Select("value").First(&setting).Error
	if err != nil {
		return nil, err
	}

	// 解析JSON配置
	config := &smtpConfig{
		Enabled:  true,
		Host:     "smtp.example.com",
		Port:     465,
		Username: "",
		Password: "",
		FromName: "锚点财务",
		FromAddr: "noreply@example.com",
	}

	// 实际应该从setting.Value解析JSON
	_ = setting.Value

	return config, nil
}

// sendViaSMTP 通过SMTP发送邮件
func (s *EmailEnhancedService) sendViaSMTP(config *smtpConfig, to, subject, body string) error {
	// 使用Go标准库的net/smtp发送邮件
	// 这里简化实现，实际应该使用完整的SMTP客户端
	logger.Info("Sending email via SMTP", "to", to, "subject", subject)

	// 构建邮件头
	headers := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", config.FromName, config.FromAddr),
		"To":           to,
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
		"Date":         time.Now().Format(time.RFC1123Z),
	}

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	// 实际发送
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	return smtp.SendMail(addr, auth, config.FromAddr, []string{to}, msg.Bytes())
}

type smtpConfig struct {
	Enabled  bool
	Host     string
	Port     int
	Username string
	Password string
	FromName string
	FromAddr string
}

// GetEmailLogs 获取邮件日志
func (s *EmailEnhancedService) GetEmailLogs(page, pageSize int) ([]EmailLog, int64, error) {
	var logs []EmailLog
	var total int64

	s.db.Model(&EmailLog{}).Count(&total)
	err := s.db.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs).Error

	return logs, total, err
}

// GetEmailStats 获取邮件统计
func (s *EmailEnhancedService) GetEmailStats(period string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var totalSent, totalFailed int64
	s.db.Model(&EmailLog{}).Where("status = ?", "sent").Count(&totalSent)
	s.db.Model(&EmailLog{}).Where("status = ?", "failed").Count(&totalFailed)

	stats["total_sent"] = totalSent
	stats["total_failed"] = totalFailed
	stats["success_rate"] = float64(0)
	if totalSent+totalFailed > 0 {
		stats["success_rate"] = float64(totalSent) / float64(totalSent+totalFailed) * 100
	}

	return stats, nil
}

// GetSupportConfig 获取支持邮箱配置
func (s *EmailEnhancedService) GetSupportConfig() (*EmailSupport, error) {
	var setting struct {
		Value string
	}
	err := s.db.Table("system_configs").Where("key = ?", "email_support").Select("value").First(&setting).Error
	if err != nil {
		return &EmailSupport{
			SupportEmail: "support@example.com",
			SalesEmail:   "sales@example.com",
			AbuseEmail:   "abuse@example.com",
			BillingEmail: "billing@example.com",
			NoReplyEmail: "noreply@example.com",
			NoReplyName:  "No Reply",
		}, nil
	}

	config := &EmailSupport{}
	// 实际应该从setting.Value解析JSON
	_ = setting.Value

	return config, nil
}

// UpdateSupportConfig 更新支持邮箱配置
func (s *EmailEnhancedService) UpdateSupportConfig(config *EmailSupport) error {
	return s.db.Table("system_configs").Where("key = ?", "email_support").
		Assign(map[string]interface{}{"value": fmt.Sprintf(`{"support_email":"%s","sales_email":"%s","abuse_email":"%s","billing_email":"%s","no_reply_email":"%s","no_reply_name":"%s"}`,
			config.SupportEmail, config.SalesEmail, config.AbuseEmail, config.BillingEmail, config.NoReplyEmail, config.NoReplyName)}).
		FirstOrCreate(&struct{ Key, Value string }{}).Error
}

// SendEmailDirect 直接发送邮件（不使用模板）
func (s *EmailEnhancedService) SendEmailDirect(to, subject, body string) error {
	return s.SendEmail(to, subject, body, nil)
}

// SendEmailDiy 自定义邮件发送
func (s *EmailEnhancedService) SendEmailDiy(to string, tplCode string, params map[string]interface{}) error {
	// 转换参数
	strParams := make(map[string]string)
	for k, v := range params {
		strParams[k] = fmt.Sprintf("%v", v)
	}
	return s.SendEmailWithTemplate(to, tplCode, strParams)
}

// ReplaceEmailContentParams 替换邮件内容参数
func (s *EmailEnhancedService) ReplaceEmailContentParams(content string, params map[string]string) string {
	for key, value := range params {
		content = strings.ReplaceAll(content, "{"+key+"}", value)
	}
	return content
}

// GetBaseArg 获取基础参数
func (s *EmailEnhancedService) GetBaseArg() map[string]string {
	return map[string]string{
		"site_name": "锚点财务",
		"site_url":  "https://example.com",
		"year":      fmt.Sprintf("%d", time.Now().Year()),
		"date":      time.Now().Format("2006-01-02"),
		"time":      time.Now().Format("15:04:05"),
	}
}

// SendEmailCode 发送邮箱验证码
func (s *EmailEnhancedService) SendEmailCode(email string) (string, error) {
	code := generateVerifyCode()
	params := map[string]string{
		"code": code,
	}
	return code, s.SendEmailWithTemplate(email, "verify_email", params)
}

// SendEmailBind 发送绑定邮箱验证码
func (s *EmailEnhancedService) SendEmailBind(userID uint, email string) (string, error) {
	code := generateVerifyCode()
	params := map[string]string{
		"code": code,
	}
	return code, s.SendEmailWithTemplate(email, "bind_email", params)
}

// generateVerifyCode 生成验证码
func generateVerifyCode() string {
	return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}

// Use Go template for complex email rendering
func renderTemplate(tplStr string, data interface{}) (string, error) {
	t, err := template.New("email").Parse(tplStr)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	return buf.String(), err
}
