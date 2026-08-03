package service

import (
	"time"

	"anchorfinance/pkg/email"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// NotificationService 通知服务
// 同一失败事件只通知一次
type NotificationService struct {
	db         *gorm.DB
	log        *logger.Logger
	deduplicator *Deduplicator
}

// NewNotificationService 创建通知服务
func NewNotificationService(db *gorm.DB, log *logger.Logger) *NotificationService {
	return &NotificationService{
		db:         db,
		log:        log,
		deduplicator: NewDeduplicator(db, log),
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
	s.db.Table("system_settings").
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
