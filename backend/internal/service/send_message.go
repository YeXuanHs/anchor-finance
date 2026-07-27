package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"anchorfinance/pkg/logger"
)

// SendMessage 消息发送
type SendMessage struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Type        string     `gorm:"size:32;not null;index" json:"type"` // email/sms/site_message
	Channel     string     `gorm:"size:32" json:"channel"`            // 通道
	To          string     `gorm:"size:512;not null" json:"to"`        // 收件人（逗号分隔）
	CC          string     `gorm:"size:512" json:"cc"`
	BCC         string     `gorm:"size:512" json:"bcc"`
	Subject     string     `gorm:"size:256" json:"subject"`
	Content     string     `gorm:"type:text;not null" json:"content"`
	ContentType string     `gorm:"size:32;default:text" json:"content_type"` // text/html
	TemplateID  *uint      `gorm:"index" json:"template_id"`
	TemplateData string    `gorm:"type:text" json:"template_data"` // 模板变量JSON
	Status      int        `gorm:"default:0;index" json:"status"` // 0=待发送 1=发送中 2=成功 3=失败 4=部分失败
	RetryCount  int        `gorm:"default:0" json:"retry_count"`
	MaxRetries  int        `gorm:"default:3" json:"max_retries"`
	SentAt      *time.Time `json:"sent_at"`
	ErrorMsg    string     `gorm:"type:text" json:"error_msg"`
	UserID      *uint      `gorm:"index" json:"user_id"` // 关联用户
	BatchID     string     `gorm:"size:64;index" json:"batch_id"` // 批次ID
	Source      string     `gorm:"size:32" json:"source"` // system/admin/api/cron
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// SendMessageQueue 消息发送队列
type SendMessageQueue struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	MessageID   uint       `gorm:"index;not null" json:"message_id"`
	Priority    int        `gorm:"default:5" json:"priority"` // 1=最高 10=最低
	ScheduledAt time.Time  `gorm:"index;not null" json:"scheduled_at"`
	Attempts    int        `gorm:"default:0" json:"attempts"`
	Status      int        `gorm:"default:0;index" json:"status"` // 0=待处理 1=处理中 2=完成 3=失败
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// SendMessageTemplate 消息模板
type SendMessageTemplate struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:128;not null;uniqueIndex" json:"name"`
	Type        string `gorm:"size:32;not null" json:"type"` // email/sms/site_message
	Subject     string `gorm:"size:256" json:"subject"`
	Content     string `gorm:"type:text;not null" json:"content"`
	ContentType string `gorm:"size:32;default:text" json:"content_type"`
	Variables   string `gorm:"type:text" json:"variables"` // 可用变量说明JSON
	Status      int    `gorm:"default:1" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SendMessageService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewSendMessageService(db *gorm.DB, log *logger.Logger) *SendMessageService {
	return &SendMessageService{db: db, log: log}
}

type SendEmailRequest struct {
	To          string `json:"to" binding:"required"`
	CC          string `json:"cc"`
	BCC         string `json:"bcc"`
	Subject     string `json:"subject" binding:"required"`
	Content     string `json:"content" binding:"required"`
	ContentType string `json:"content_type"`
	TemplateID  *uint  `json:"template_id"`
	UserID      *uint  `json:"user_id"`
}

type SendSMSRequest struct {
	To      string `json:"to" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserID  *uint  `json:"user_id"`
}

type SendSiteMessageRequest struct {
	To      []uint `json:"to" binding:"required"`
	Subject string `json:"subject" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserID  *uint  `json:"user_id"`
}

type BatchSendRequest struct {
	Type    string   `json:"type" binding:"required,oneof=email sms site_message"`
	To      []string `json:"to" binding:"required"`
	Subject string   `json:"subject"`
	Content string   `json:"content" binding:"required"`
}

// SendEmail sends an email message.
func (s *SendMessageService) SendEmail(req SendEmailRequest) (*SendMessage, error) {
	contentType := req.ContentType
	if contentType == "" {
		contentType = "html"
	}

	msg := &SendMessage{
		Type:        "email",
		To:          req.To,
		CC:          req.CC,
		BCC:         req.BCC,
		Subject:     req.Subject,
		Content:     req.Content,
		ContentType: contentType,
		TemplateID:  req.TemplateID,
		UserID:      req.UserID,
		Status:      0,
		Source:      "admin",
	}

	if err := s.db.Create(msg).Error; err != nil {
		return nil, err
	}

	// Add to queue
	queue := &SendMessageQueue{
		MessageID:   msg.ID,
		Priority:    5,
		ScheduledAt: time.Now(),
		Status:      0,
	}
	s.db.Create(queue)

	s.log.Infof("email queued: to=%s subject=%s", req.To, req.Subject)
	return msg, nil
}

// SendSMS sends an SMS message.
func (s *SendMessageService) SendSMS(req SendSMSRequest) (*SendMessage, error) {
	msg := &SendMessage{
		Type:    "sms",
		To:      req.To,
		Content: req.Content,
		UserID:  req.UserID,
		Status:  0,
		Source:  "admin",
	}

	if err := s.db.Create(msg).Error; err != nil {
		return nil, err
	}

	queue := &SendMessageQueue{
		MessageID:   msg.ID,
		Priority:    3,
		ScheduledAt: time.Now(),
		Status:      0,
	}
	s.db.Create(queue)

	s.log.Infof("sms queued: to=%s", req.To)
	return msg, nil
}

// SendSiteMessage sends a site message.
func (s *SendMessageService) SendSiteMessage(req SendSiteMessageRequest) ([]SendMessage, error) {
	var messages []SendMessage

	for _, userID := range req.To {
		msg := SendMessage{
			Type:    "site_message",
			To:      "",
			Subject: req.Subject,
			Content: req.Content,
			UserID:  &userID,
			Status:  2, // 立即成功
			SentAt:  timePtr(time.Now()),
			Source:  "admin",
		}
		messages = append(messages, msg)
	}

	if err := s.db.CreateInBatches(messages, 100).Error; err != nil {
		return nil, err
	}

	s.log.Infof("site messages sent: count=%d subject=%s", len(messages), req.Subject)
	return messages, nil
}

// BatchSend sends messages to multiple recipients.
func (s *SendMessageService) BatchSend(req BatchSendRequest) (string, []SendMessage, error) {
	if len(req.To) == 0 {
		return "", nil, errors.New("recipients list is empty")
	}

	batchID := generateBatchID()
	var messages []SendMessage

	for _, to := range req.To {
		msg := SendMessage{
			Type:    req.Type,
			To:      to,
			Subject: req.Subject,
			Content: req.Content,
			Status:  0,
			BatchID: batchID,
			Source:  "admin",
		}
		messages = append(messages, msg)
	}

	if err := s.db.CreateInBatches(messages, 100).Error; err != nil {
		return "", nil, err
	}

	// Add to queue
	for _, msg := range messages {
		queue := &SendMessageQueue{
			MessageID:   msg.ID,
			Priority:    5,
			ScheduledAt: time.Now(),
			Status:      0,
		}
		s.db.Create(queue)
	}

	s.log.Infof("batch send created: batch=%s type=%s count=%d", batchID, req.Type, len(messages))
	return batchID, messages, nil
}

// GetByID returns a single message by ID.
func (s *SendMessageService) GetByID(id uint) (*SendMessage, error) {
	var msg SendMessage
	if err := s.db.First(&msg, id).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

// GetList returns all messages with pagination.
func (s *SendMessageService) GetList(page, pageSize int, msgType *string, status *int) ([]SendMessage, int64, error) {
	var messages []SendMessage
	var total int64

	query := s.db.Model(&SendMessage{})
	if msgType != nil {
		query = query.Where("type = ?", *msgType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&messages).Error; err != nil {
		return nil, 0, err
	}
	return messages, total, nil
}

// GetByBatchID returns all messages in a batch.
func (s *SendMessageService) GetByBatchID(batchID string) ([]SendMessage, error) {
	var messages []SendMessage
	if err := s.db.Where("batch_id = ?", batchID).Order("id ASC").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

// GetUserMessages returns messages for a specific user.
func (s *SendMessageService) GetUserMessages(userID uint, page, pageSize int) ([]SendMessage, int64, error) {
	var messages []SendMessage
	var total int64

	query := s.db.Model(&SendMessage{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&messages).Error; err != nil {
		return nil, 0, err
	}
	return messages, total, nil
}

// GetQueuePending returns pending queue items.
func (s *SendMessageService) GetQueuePending(limit int) ([]SendMessageQueue, error) {
	var items []SendMessageQueue
	if err := s.db.Where("status = 0 AND scheduled_at <= ?", time.Now()).
		Order("priority ASC, id ASC").Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// MarkSending marks a message as sending.
func (s *SendMessageService) MarkSending(id uint) error {
	return s.db.Model(&SendMessage{}).Where("id = ?", id).Update("status", 1).Error
}

// MarkSent marks a message as sent successfully.
func (s *SendMessageService) MarkSent(id uint) error {
	now := time.Now()
	return s.db.Model(&SendMessage{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":  2,
		"sent_at": &now,
	}).Error
}

// MarkFailed marks a message as failed.
func (s *SendMessageService) MarkFailed(id uint, errMsg string) error {
	return s.db.Model(&SendMessage{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     3,
		"error_msg":  errMsg,
		"retry_count": gorm.Expr("retry_count + 1"),
	}).Error
}

// RetryFailed requeues failed messages for retry.
func (s *SendMessageService) RetryFailed(limit int) error {
	var messages []SendMessage
	s.db.Where("status = 3 AND retry_count < max_retries").Limit(limit).Find(&messages)

	for _, msg := range messages {
		s.db.Model(&msg).Update("status", 0)
		queue := &SendMessageQueue{
			MessageID:   msg.ID,
			Priority:    5,
			ScheduledAt: time.Now(),
			Status:      0,
		}
		s.db.Create(queue)
	}

	return nil
}

// Helper functions
func timePtr(t time.Time) *time.Time {
	return &t
}

func generateBatchID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
