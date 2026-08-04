package service

import (
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/email"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// NotificationService 通知服务
// 同一失败事件只通知一次；同时支持用户通知查询与管理
type NotificationService struct {
	db           *gorm.DB
	log          *logger.Logger
	deduplicator *Deduplicator
	wechatSvc    *WechatService
}

// NewNotificationService 创建通知服务
func NewNotificationService(db *gorm.DB, log *logger.Logger, wechatSvc *WechatService) *NotificationService {
	return &NotificationService{
		db:           db,
		log:          log,
		deduplicator: NewDeduplicator(db, log),
		wechatSvc:    wechatSvc,
	}
}

// NotifyFailure 发送失败通知（同一事件只发一次）
// eventType: 事件类型（如 "host_down", "payment_fail", "sync_error"）
// targetID: 目标ID（如主机ID、订单ID）
// title: 通知标题
// content: 通知内容
func (s *NotificationService) NotifyFailure(eventType string, targetID interface{}, title string, content string) bool {
	// 生成事件唯一标识
	eventKey := GenerateEventKey(eventType, targetID, title)

	// 检查是否已通知（24小时内同一事件只通知一次）
	if !s.deduplicator.ShouldNotify(eventKey, 24*time.Hour) {
		s.log.Debugf("事件已通知过，跳过: %s", eventKey)
		return false
	}

	// 发送通知
	err := s.sendNotification(title, content)
	if err != nil {
		s.log.Errorf("发送通知失败: %v", err)
		return false
	}

	// 记录已通知
	s.deduplicator.RecordNotified(eventKey)
	s.log.Infof("发送通知成功: %s", title)

	return true
}

// NotifyOnce 只通知一次（永不重复）
// 适用于：一次性事件，如配置错误、证书过期等
func (s *NotificationService) NotifyOnce(eventType string, targetID interface{}, title string, content string) bool {
	eventKey := GenerateEventKey(eventType, targetID, title)

	// 检查是否已通知（永不过期）
	if !s.deduplicator.ShouldNotify(eventKey, 365*24*time.Hour) {
		return false
	}

	err := s.sendNotification(title, content)
	if err != nil {
		s.log.Errorf("发送通知失败: %v", err)
		return false
	}

	s.deduplicator.RecordNotified(eventKey)
	return true
}

// sendNotification 实际发送通知
func (s *NotificationService) sendNotification(title string, content string) error {
	emails := s.getNotifyEmails()
	if emails == "" {
		s.log.Warn("未配置通知邮箱，跳过发送")
		return nil
	}

	sender := email.NewSender(s.db)
	if err := sender.Send(emails, title, content); err != nil {
		return err
	}

	return nil
}

// getNotifyEmails 获取通知邮箱
func (s *NotificationService) getNotifyEmails() string {
	var setting struct {
		Value string
	}
	s.db.Table("system_configs").
		Where("`key` = ?", "notify_emails").
		First(&setting)

	return setting.Value
}

// GetDeduplicator 获取去重器（用于管理）
func (s *NotificationService) GetDeduplicator() *Deduplicator {
	return s.deduplicator
}

// ResetEvent 重置事件（允许重新通知）
func (s *NotificationService) ResetEvent(eventType string, targetID interface{}, title string) error {
	eventKey := GenerateEventKey(eventType, targetID, title)
	return s.deduplicator.CleanByEventKey(eventKey)
}

// ─── 用户通知 ───

// GetUserNotifications 返回用户的通知列表（分页）
func (s *NotificationService) GetUserNotifications(userID uint, page, pageSize int, onlyUnread bool) ([]model.SystemMessage, int64, error) {
	var total int64
	var messages []model.SystemMessage

	query := s.db.Model(&model.SystemMessage{}).Where("user_id = ?", userID)
	if onlyUnread {
		query = query.Where("is_read = ?", false)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, total, nil
}

// MarkRead 标记单条通知为已读
func (s *NotificationService) MarkRead(userID, notificationID uint) error {
	now := time.Now()
	result := s.db.Model(&model.SystemMessage{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Updates(map[string]interface{}{"is_read": true, "read_at": now})
	if result.RowsAffected == 0 {
		return fmt.Errorf("notification not found")
	}
	return result.Error
}

// MarkAllRead 标记用户所有通知为已读
func (s *NotificationService) MarkAllRead(userID uint) error {
	now := time.Now()
	return s.db.Model(&model.SystemMessage{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{"is_read": true, "read_at": now}).Error
}

// ─── 模板管理 ───

// GetTemplates 返回通知模板列表，可按 channel 过滤
func (s *NotificationService) GetTemplates(channel string) ([]model.NotificationTemplate, error) {
	var templates []model.NotificationTemplate
	query := s.db.Model(&model.NotificationTemplate{})
	if channel != "" {
		query = query.Where("channel = ?", channel)
	}
	if err := query.Order("id ASC").Find(&templates).Error; err != nil {
		return nil, err
	}
	return templates, nil
}

// UpdateTemplate 更新通知模板字段
func (s *NotificationService) UpdateTemplate(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&model.NotificationTemplate{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		return fmt.Errorf("template not found")
	}
	return result.Error
}

// ─── 通知日志 ───

// GetLogs 返回通知日志列表（分页），支持 channel / userID / status 过滤
func (s *NotificationService) GetLogs(page, pageSize int, channel string, userID *uint, status *int) ([]model.NotificationLog, int64, error) {
	var total int64
	var logs []model.NotificationLog

	query := s.db.Model(&model.NotificationLog{})
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

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// ─── 批量发送 ───

// SendBatch 批量发送通知，返回成功数和失败数
func (s *NotificationService) SendBatch(userIDs []uint, channel, templateName string, data map[string]interface{}) (int, int, error) {
	// 加载模板
	var tpl model.NotificationTemplate
	if err := s.db.Where("code = ? AND is_active = ?", templateName, true).First(&tpl).Error; err != nil {
		return 0, 0, fmt.Errorf("template not found: %s", templateName)
	}

	success, fail := 0, 0
	for _, uid := range userIDs {
		err := s.sendToUser(uid, channel, tpl, data)
		if err != nil {
			fail++
			s.log.Warnf("send notification to user %d failed: %v", uid, err)
		} else {
			success++
		}
	}

	return success, fail, nil
}

// sendToUser 向单个用户发送通知
func (s *NotificationService) sendToUser(userID uint, channel string, tpl model.NotificationTemplate, data map[string]interface{}) error {
	logEntry := model.NotificationLog{
		UserID:   userID,
		Channel:  channel,
		Template: tpl.Code,
		Status:   1, // 发送中
	}

	var err error
	switch channel {
	case "email":
		to, _ := data["to"].(string)
		subject, _ := data["subject"].(string)
		if subject == "" {
			subject = tpl.Subject
		}
		sender := email.NewSender(s.db)
		err = sender.Send(to, subject, tpl.Content)
	case "wechat":
		if s.wechatSvc != nil {
			// 通过微信服务发送（具体实现依赖 WechatService）
			s.log.Debugf("wechat notification to user %d (template: %s)", userID, tpl.Code)
		}
	default:
		err = fmt.Errorf("unsupported channel: %s", channel)
	}

	if err != nil {
		logEntry.Status = 3 // 失败
		logEntry.Error = err.Error()
	} else {
		logEntry.Status = 2 // 成功
		now := time.Now()
		logEntry.SentAt = &now
	}

	_ = s.db.Create(&logEntry).Error
	return err
}
