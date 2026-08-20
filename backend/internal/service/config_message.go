package service

import (
	"errors"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type ConfigMessageService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewConfigMessageService(db *gorm.DB, log *logger.Logger) *ConfigMessageService {
	return &ConfigMessageService{db: db, log: log}
}

type UpdateMessageConfigRequest struct {
	Name        *string `json:"name"`
	Provider    *string `json:"provider"`
	Config      *string `json:"config"`
	SenderName  *string `json:"sender_name"`
	SenderAddr  *string `json:"sender_addr"`
	Signature   *string `json:"signature"`
	RateLimit   *int    `json:"rate_limit"`
	DailyLimit  *int    `json:"daily_limit"`
	IsEnabled   *bool   `json:"is_enabled"`
	TestAddress *string `json:"test_address"`
	Status      *int16  `json:"status"`
	Remark      *string `json:"remark"`
}

func (s *ConfigMessageService) GetAll() ([]model.MessageConfig, error) {
	var items []model.MessageConfig
	if err := s.db.Order("channel ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ConfigMessageService) GetByChannel(channel string) (*model.MessageConfig, error) {
	var item model.MessageConfig
	if err := s.db.Where("channel = ?", channel).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ConfigMessageService) Update(channel string, req UpdateMessageConfigRequest) (*model.MessageConfig, error) {
	var item model.MessageConfig
	if err := s.db.Where("channel = ?", channel).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			item = model.MessageConfig{
				Channel: channel,
				Status:  1,
			}
			if req.Name != nil {
				item.Name = *req.Name
			}
			if req.Provider != nil {
				item.Provider = *req.Provider
			}
			if req.Config != nil {
				item.Config = []byte(*req.Config)
			}
			if req.SenderName != nil {
				item.SenderName = *req.SenderName
			}
			if req.SenderAddr != nil {
				item.SenderAddr = *req.SenderAddr
			}
			if req.Signature != nil {
				item.Signature = *req.Signature
			}
			if req.RateLimit != nil {
				item.RateLimit = *req.RateLimit
			}
			if req.DailyLimit != nil {
				item.DailyLimit = *req.DailyLimit
			}
			if req.IsEnabled != nil {
				item.IsEnabled = *req.IsEnabled
			}
			if req.TestAddress != nil {
				item.TestAddress = *req.TestAddress
			}
			if req.Status != nil {
				item.Status = *req.Status
			}
			if req.Remark != nil {
				item.Remark = *req.Remark
			}
			if err := s.db.Create(&item).Error; err != nil {
				return nil, err
			}
			s.log.Infof("message config created: channel=%s", channel)
			return &item, nil
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Provider != nil {
		updates["provider"] = *req.Provider
	}
	if req.Config != nil {
		updates["config"] = *req.Config
	}
	if req.SenderName != nil {
		updates["sender_name"] = *req.SenderName
	}
	if req.SenderAddr != nil {
		updates["sender_addr"] = *req.SenderAddr
	}
	if req.Signature != nil {
		updates["signature"] = *req.Signature
	}
	if req.RateLimit != nil {
		updates["rate_limit"] = *req.RateLimit
	}
	if req.DailyLimit != nil {
		updates["daily_limit"] = *req.DailyLimit
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	if req.TestAddress != nil {
		updates["test_address"] = *req.TestAddress
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != nil {
		updates["remark"] = *req.Remark
	}

	if len(updates) > 0 {
		if err := s.db.Model(&item).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.Where("channel = ?", channel).First(&item).Error; err != nil {
		return nil, err
	}
	s.log.Infof("message config updated: channel=%s", channel)
	return &item, nil
}

func (s *ConfigMessageService) TestSend(channel string) error {
	var config model.MessageConfig
	if err := s.db.Where("channel = ?", channel).First(&config).Error; err != nil {
		return errors.New("message channel not found")
	}

	if !config.IsEnabled {
		return errors.New("message channel is disabled")
	}

	if config.TestAddress == "" {
		return errors.New("test address not configured")
	}

	now := time.Now()
	ok := true

	switch channel {
	case "email":
		if config.SenderAddr == "" {
			ok = false
		}
	case "sms":
		if config.Signature == "" {
			ok = false
		}
	case "site":
		// 站内信不需要额外配置校验
	default:
		return errors.New("unsupported message channel: " + channel)
	}

	updates := map[string]interface{}{
		"last_test_at": &now,
		"last_test_ok": ok,
	}
	if err := s.db.Model(&config).Updates(updates).Error; err != nil {
		return err
	}

	if !ok {
		return errors.New("test send failed: incomplete configuration")
	}

	s.log.Infof("message config test send: channel=%s ok=%v", channel, ok)
	return nil
}

func (s *ConfigMessageService) GetEnabledChannels() ([]model.MessageConfig, error) {
	var items []model.MessageConfig
	if err := s.db.Where("is_enabled = ? AND status = 1", true).
		Order("channel ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// MessageTemplate represents a message template.
type MessageTemplate struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Channel     string `json:"channel"`
	Subject     string `json:"subject"`
	Content     string `json:"content"`
	Description string `json:"description"`
	Variables   string `json:"variables"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateTemplate creates a new message template.
func (s *ConfigMessageService) CreateTemplate(name, channel, subject, content, description, variables string) (*MessageTemplate, error) {
	template := map[string]interface{}{
		"name":        name,
		"channel":     channel,
		"subject":     subject,
		"content":     content,
		"description": description,
		"variables":   variables,
		"status":      1,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}
	if err := s.db.Table("message_templates").Create(&template).Error; err != nil {
		return nil, err
	}

	return &MessageTemplate{
		Name:        name,
		Channel:     channel,
		Subject:     subject,
		Content:     content,
		Description: description,
		Variables:   variables,
		Status:      1,
	}, nil
}

// UpdateTemplate updates a message template.
func (s *ConfigMessageService) UpdateTemplate(id, name, subject, content, description string, variables string, status *int) error {
	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if name != "" {
		updates["name"] = name
	}
	if subject != "" {
		updates["subject"] = subject
	}
	if content != "" {
		updates["content"] = content
	}
	if description != "" {
		updates["description"] = description
	}
	if variables != "" {
		updates["variables"] = variables
	}
	if status != nil {
		updates["status"] = *status
	}

	result := s.db.Table("message_templates").Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		return errors.New("template not found")
	}
	return result.Error
}

// DeleteTemplate deletes a message template.
func (s *ConfigMessageService) DeleteTemplate(id string) error {
	result := s.db.Table("message_templates").Where("id = ?", id).Delete(nil)
	if result.RowsAffected == 0 {
		return errors.New("template not found")
	}
	return result.Error
}

// GetTemplateDesc returns template description and variables.
func (s *ConfigMessageService) GetTemplateDesc(id string) (map[string]interface{}, error) {
	var template MessageTemplate
	if err := s.db.Table("message_templates").Where("id = ?", id).First(&template).Error; err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"description": template.Description,
		"variables":   template.Variables,
		"content":     template.Content,
	}, nil
}

// SetSmsTemplate sets an SMS template.
func (s *ConfigMessageService) SetSmsTemplate(templateID, smsContent string) error {
	updates := map[string]interface{}{
		"sms_content": smsContent,
		"updated_at":  time.Now(),
	}
	result := s.db.Table("message_templates").Where("id = ? AND channel = ?", templateID, "sms").Updates(updates)
	if result.RowsAffected == 0 {
		return errors.New("SMS template not found")
	}
	return result.Error
}

// BeforeSendMessageCheck checks if a message can be sent.
func (s *ConfigMessageService) BeforeSendMessageCheck(channel, templateID string, userID uint) (map[string]interface{}, error) {
	// Check channel config
	var config model.MessageConfig
	if err := s.db.Where("channel = ?", channel).First(&config).Error; err != nil {
		return nil, errors.New("message channel not found")
	}

	if !config.IsEnabled {
		return map[string]interface{}{"can_send": false, "reason": "channel disabled"}, nil
	}

	// Check rate limits
	if config.RateLimit > 0 {
		var recentCount int64
		s.db.Table("message_logs").Where("channel = ? AND created_at > ?", channel, time.Now().Add(-time.Minute)).Count(&recentCount)
		if int(recentCount) >= config.RateLimit {
			return map[string]interface{}{"can_send": false, "reason": "rate limit exceeded"}, nil
		}
	}

	// Check daily limits
	if config.DailyLimit > 0 {
		today := time.Now().Format("2006-01-02")
		var dailyCount int64
		s.db.Table("message_logs").Where("channel = ? AND DATE(created_at) = ?", channel, today).Count(&dailyCount)
		if int(dailyCount) >= config.DailyLimit {
			return map[string]interface{}{"can_send": false, "reason": "daily limit exceeded"}, nil
		}
	}

	return map[string]interface{}{"can_send": true, "reason": ""}, nil
}
