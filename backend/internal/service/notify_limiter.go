package service

import (
	"sync"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// NotificationLimiter 通知频率限制器
// 防止失败通知轰炸邮箱
type NotificationLimiter struct {
	db          *gorm.DB
	log         *logger.Logger
	mu          sync.Mutex
	lastSent    map[string]time.Time // key -> 上次发送时间
	counts      map[string]int       // key -> 当前窗口内发送次数
	windowStart map[string]time.Time // key -> 窗口开始时间
}

// NewNotificationLimiter 创建通知限制器
func NewNotificationLimiter(db *gorm.DB, log *logger.Logger) *NotificationLimiter {
	limiter := &NotificationLimiter{
		db:          db,
		log:         log,
		lastSent:    make(map[string]time.Time),
		counts:      make(map[string]int),
		windowStart: make(map[string]time.Time),
	}

	// 启动清理协程
	go limiter.cleanup()

	return limiter
}

// CanSend 检查是否可以发送通知
// key: 通知类型+目标（如 "fail_notify:admin@example.com"）
// interval: 最小发送间隔
// maxPerWindow: 窗口内最大发送次数
// windowSize: 窗口大小
func (l *NotificationLimiter) CanSend(key string, interval time.Duration, maxPerWindow int, windowSize time.Duration) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	// 检查最小间隔
	if lastSend, ok := l.lastSent[key]; ok {
		if now.Sub(lastSend) < interval {
			l.log.Debugf("通知被限制: %s, 距上次发送 %v, 需要等待 %v", key, now.Sub(lastSend), interval)
			return false
		}
	}

	// 检查窗口内发送次数
	windowStart, ok := l.windowStart[key]
	if !ok || now.Sub(windowStart) > windowSize {
		// 新窗口
		l.windowStart[key] = now
		l.counts[key] = 0
	}

	if l.counts[key] >= maxPerWindow {
		l.log.Debugf("通知被限制: %s, 当前窗口已发送 %d 次, 最大 %d 次", key, l.counts[key], maxPerWindow)
		return false
	}

	return true
}

// RecordSend 记录发送
func (l *NotificationLimiter) RecordSend(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.lastSent[key] = time.Now()
	l.counts[key]++
}

// cleanup 定期清理过期记录
func (l *NotificationLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		l.mu.Lock()
		now := time.Now()
		for key, lastSend := range l.lastSent {
			if now.Sub(lastSend) > 1*time.Hour {
				delete(l.lastSent, key)
				delete(l.counts, key)
				delete(l.windowStart, key)
			}
		}
		l.mu.Unlock()
	}
}

// GetNotificationConfig 获取通知配置
func (l *NotificationLimiter) GetNotificationConfig() map[string]interface{} {
	// 从数据库读取配置
	var config struct {
		Enabled         bool   `json:"enabled"`
		MinInterval     int    `json:"min_interval"`      // 最小间隔（秒）
		MaxPerHour      int    `json:"max_per_hour"`      // 每小时最大发送次数
		MaxPerDay       int    `json:"max_per_day"`       // 每天最大发送次数
		NotifyEmails    string `json:"notify_emails"`     // 通知邮箱（逗号分隔）
		NotifyTypes     string `json:"notify_types"`      // 通知类型（逗号分隔）
		QuietHoursStart int    `json:"quiet_hours_start"` // 免打扰开始时间（小时）
		QuietHoursEnd   int    `json:"quiet_hours_end"`   // 免打扰结束时间（小时）
	}

	// 默认配置
	config.Enabled = true
	config.MinInterval = 300    // 5分钟
	config.MaxPerHour = 10      // 每小时最多10次
	config.MaxPerDay = 50       // 每天最多50次
	config.QuietHoursStart = 23 // 23:00开始免打扰
	config.QuietHoursEnd = 8    // 08:00结束免打扰

	// 从数据库读取
	l.db.Table("system_settings").
		Where("`key` IN ?", []string{
			"notify_enabled", "notify_min_interval", "notify_max_per_hour",
			"notify_max_per_day", "notify_emails", "notify_types",
			"notify_quiet_start", "notify_quiet_end",
		}).
		Select("`key`, `value`").
		Find(nil)

	return map[string]interface{}{
		"enabled":           config.Enabled,
		"min_interval":      config.MinInterval,
		"max_per_hour":      config.MaxPerHour,
		"max_per_day":       config.MaxPerDay,
		"notify_emails":     config.NotifyEmails,
		"notify_types":      config.NotifyTypes,
		"quiet_hours_start": config.QuietHoursStart,
		"quiet_hours_end":   config.QuietHoursEnd,
	}
}

// IsInQuietHours 检查是否在免打扰时间段
func (l *NotificationLimiter) IsInQuietHours() bool {
	config := l.GetNotificationConfig()
	hour := time.Now().Hour()

	start := config["quiet_hours_start"].(int)
	end := config["quiet_hours_end"].(int)

	if start < end {
		return hour >= start && hour < end
	}
	// 跨午夜的情况（如 23:00 - 08:00）
	return hour >= start || hour < end
}

// ShouldNotify 检查是否应该发送通知
func (l *NotificationLimiter) ShouldNotify(notifyType string) bool {
	config := l.GetNotificationConfig()

	// 检查是否启用
	if !config["enabled"].(bool) {
		return false
	}

	// 检查是否在免打扰时间段
	if l.IsInQuietHours() {
		l.log.Debugf("当前在免打扰时间段，跳过通知: %s", notifyType)
		return false
	}

	// 检查通知类型是否启用
	// TODO: 实现类型过滤

	// 检查频率限制
	key := "fail_notify:" + notifyType
	minInterval := time.Duration(config["min_interval"].(int)) * time.Second
	maxPerHour := config["max_per_hour"].(int)

	return l.CanSend(key, minInterval, maxPerHour, 1*time.Hour)
}

// SendNotification 发送通知（带频率限制）
func (l *NotificationLimiter) SendNotification(notifyType string, subject string, content string) bool {
	if !l.ShouldNotify(notifyType) {
		return false
	}

	// 记录发送
	key := "fail_notify:" + notifyType
	l.RecordSend(key)

	// TODO: 实际发送邮件
	l.log.Infof("发送通知: type=%s, subject=%s", notifyType, subject)

	return true
}
