package service

import (
	"bytes"
	"errors"
	"html/template"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// NotificationTemplate 消息模板
type NotificationTemplate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Channel   string    `gorm:"type:varchar(20);not null" json:"channel"`
	Subject   string    `gorm:"type:varchar(255)" json:"subject"`
	Content   string    `gorm:"type:text" json:"content"`
	Variables string    `gorm:"type:text" json:"variables"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NotificationLog 通知日志
type NotificationLog struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index" json:"user_id"`
	Channel   string     `gorm:"type:varchar(20);not null" json:"channel"`
	Template  string     `gorm:"type:varchar(50)" json:"template"`
	To        string     `gorm:"type:varchar(100)" json:"to"`
	Subject   string     `gorm:"type:varchar(255)" json:"subject"`
	Content   string     `gorm:"type:text" json:"content"`
	Status    int8       `gorm:"type:smallint;default:1" json:"status"`
	Error     string     `gorm:"type:text" json:"error"`
	SentAt    *time.Time `json:"sent_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type NotificationService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewNotificationService(db *gorm.DB, log *logger.Logger) *NotificationService {
	return &NotificationService{db: db, log: log}
}

type SendRequest struct {
	To      string                 `json:"to"`
	Subject string                 `json:"subject"`
	Data    map[string]interface{} `json:"data"`
}

type BatchSendRequest struct {
	UserIDs []uint `json:"user_ids" binding:"required"`
	Channel string `json:"channel" binding:"required"`
	To      string `json:"to"` // 预留目标地址，批量时按用户配置
	Subject string `json:"subject"`
	Data    map[string]interface{} `json:"data"`
}

// Send sends a notification to a user via the specified channel and template.
func (s *NotificationService) Send(userID uint, channel, templateCode string, data map[string]interface{}) error {
	tmpl, err := s.getTemplate(channel, templateCode)
	if err != nil {
		return err
	}

	subject, err := s.renderString(tmpl.Subject, data)
	if err != nil {
		return err
	}
	content, err := s.renderString(tmpl.Content, data)
	if err != nil {
		return err
	}

	to := ""
	if v, ok := data["to"].(string); ok {
		to = v
	}

	logEntry := &NotificationLog{
		UserID:   userID,
		Channel:  channel,
		Template: templateCode,
		To:       to,
		Subject:  subject,
		Content:  content,
		Status:   1,
	}
	if err := s.db.Create(logEntry).Error; err != nil {
		return err
	}

	// 实际发送逻辑按 channel 分发（此处为占位）
	if err := s.dispatch(channel, to, subject, content); err != nil {
		now := time.Now()
		s.db.Model(logEntry).Updates(map[string]interface{}{
			"status": 3,
			"error":  err.Error(),
		})
		s.log.Errorf("notification send failed (user=%d, channel=%s): %v", userID, channel, err)
		return err
	}

	now := time.Now()
	s.db.Model(logEntry).Updates(map[string]interface{}{
		"status": 2,
		"sent_at": &now,
	})
	s.log.Infof("notification sent (user=%d, channel=%s, template=%s)", userID, channel, templateCode)
	return nil
}

// SendBatch sends a notification to multiple users.
func (s *NotificationService) SendBatch(userIDs []uint, channel, templateCode string, data map[string]interface{}) (int, int, error) {
	var success, fail int
	for _, uid := range userIDs {
		if err := s.Send(uid, channel, templateCode, data); err != nil {
			fail++
		} else {
			success++
		}
	}
	return success, fail, nil
}

// GetUserNotifications returns paginated notifications for a user.
func (s *NotificationService) GetUserNotifications(userID uint, page, pageSize int, onlyUnread bool) ([]NotificationLog, int64, error) {
	var logs []NotificationLog
	var total int64

	query := s.db.Model(&NotificationLog{}).Where("user_id = ?", userID)
	if onlyUnread {
		query = query.Where("status != 2")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// MarkRead marks a single notification as read (status=2).
func (s *NotificationService) MarkRead(userID, notifID uint) error {
	return s.db.Model(&NotificationLog{}).
		Where("id = ? AND user_id = ?", notifID, userID).
		Update("status", 2).Error
}

// MarkAllRead marks all notifications for a user as read.
func (s *NotificationService) MarkAllRead(userID uint) error {
	return s.db.Model(&NotificationLog{}).
		Where("user_id = ? AND status = 1", userID).
		Update("status", 2).Error
}

// GetTemplates returns all notification templates, optionally filtered by channel.
func (s *NotificationService) GetTemplates(channel string) ([]NotificationTemplate, error) {
	var templates []NotificationTemplate
	query := s.db.Order("id ASC")
	if channel != "" {
		query = query.Where("channel = ?", channel)
	}
	if err := query.Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// UpdateTemplate updates a notification template.
func (s *NotificationService) UpdateTemplate(tmplID uint, updates map[string]interface{}) error {
	return s.db.Model(&NotificationTemplate{}).Where("id = ?", tmplID).Updates(updates).Error
}

// GetLogs returns notification logs with pagination and optional filters.
func (s *NotificationService) GetLogs(page, pageSize int, channel string, userID *uint, status *int) ([]NotificationLog, int64, error) {
	var logs []NotificationLog
	var total int64

	query := s.db.Model(&NotificationLog{})
	if channel != "" {
		query = query.Where("channel = ?", channel)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// RenderTemplate renders a template string with provided data.
func (s *NotificationService) RenderTemplate(tmplStr string, data map[string]interface{}) (string, error) {
	return s.renderString(tmplStr, data)
}

// getTemplate fetches an active template by channel and code.
func (s *NotificationService) getTemplate(channel, code string) (*NotificationTemplate, error) {
	var tmpl NotificationTemplate
	if err := s.db.Where("channel = ? AND code = ? AND is_active = true", channel, code).First(&tmpl).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("template not found or inactive")
		}
		return nil, err
	}
	return &tmpl, nil
}

// renderString parses and executes a Go template string.
func (s *NotificationService) renderString(tmplStr string, data map[string]interface{}) (string, error) {
	t, err := template.New("").Parse(tmplStr)
	if err != nil {
		return tmplStr, nil // 解析失败时原样返回
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return tmplStr, nil
	}
	return buf.String(), nil
}

// dispatch dispatches a message via the specified channel (placeholder).
func (s *NotificationService) dispatch(channel, to, subject, content string) error {
	// 实际实现中按 channel 调用邮件/短信/微信/Webhook SDK
	s.log.Infof("dispatch[%s] to=%s subject=%s", channel, to, subject)
	return nil
}
