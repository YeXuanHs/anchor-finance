package service

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// Deduplicator 通知去重器
// 同一个失败事件只通知一次
type Deduplicator struct {
	db    *gorm.DB
	log   *logger.Logger
	mu    sync.Mutex
	cache map[string]time.Time // 内存缓存，加速查询
}

// NewDeduplicator 创建去重器
func NewDeduplicator(db *gorm.DB, log *logger.Logger) *Deduplicator {
	d := &Deduplicator{
		db:    db,
		log:   log,
		cache: make(map[string]time.Time),
	}

	// 确保表存在
	db.AutoMigrate(&NotifiedEvent{})

	// 启动时加载最近的记录到缓存
	d.loadCache()

	// 定期清理过期记录
	go d.cleanup()

	return d
}

// NotifiedEvent 已通知事件记录
type NotifiedEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EventKey  string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"event_key"`
	CreatedAt time.Time `json:"created_at"`
}

func (NotifiedEvent) TableName() string {
	return "notified_events"
}

// ShouldNotify 检查是否应该通知
// eventKey: 事件唯一标识（如 "host_down:123", "payment_fail:order456"）
// dedupWindow: 去重时间窗口（如 24小时内同一事件只通知一次）
func (d *Deduplicator) ShouldNotify(eventKey string, dedupWindow time.Duration) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 先检查内存缓存
	if lastNotify, ok := d.cache[eventKey]; ok {
		if time.Since(lastNotify) < dedupWindow {
			d.log.Debugf("事件 %s 在去重窗口内，跳过通知", eventKey)
			return false
		}
	}

	// 检查数据库
	var event NotifiedEvent
	err := d.db.Where("event_key = ?", eventKey).First(&event).Error
	if err == nil && time.Since(event.CreatedAt) < dedupWindow {
		// 更新缓存
		d.cache[eventKey] = event.CreatedAt
		return false
	}

	return true
}

// RecordNotified 记录已通知
func (d *Deduplicator) RecordNotified(eventKey string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()

	// 更新缓存
	d.cache[eventKey] = now

	// 写入数据库（存在则更新，不存在则插入）
	d.db.Where("event_key = ?", eventKey).
		Assign(NotifiedEvent{CreatedAt: now}).
		FirstOrCreate(&NotifiedEvent{EventKey: eventKey, CreatedAt: now})
}

// GenerateEventKey 生成事件唯一标识
func GenerateEventKey(eventType string, targetID interface{}, detail string) string {
	raw := fmt.Sprintf("%s:%v:%s", eventType, targetID, detail)
	hash := md5.Sum([]byte(raw))
	return fmt.Sprintf("%s:%x", eventType, hash[:8])
}

// loadCache 加载最近的记录到缓存
func (d *Deduplicator) loadCache() {
	var events []NotifiedEvent
	// 只加载最近7天的记录
	d.db.Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).Find(&events)

	for _, e := range events {
		d.cache[e.EventKey] = e.CreatedAt
	}

	d.log.Infof("加载 %d 条通知记录到缓存", len(events))
}

// cleanup 定期清理过期记录
func (d *Deduplicator) cleanup() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		// 清理30天前的记录
		result := d.db.Where("created_at < ?", time.Now().AddDate(0, 0, -30)).Delete(&NotifiedEvent{})
		if result.Error == nil && result.RowsAffected > 0 {
			d.log.Infof("清理 %d 条过期通知记录", result.RowsAffected)
		}

		// 清理内存缓存
		d.mu.Lock()
		for key, t := range d.cache {
			if time.Since(t) > 7*24*time.Hour {
				delete(d.cache, key)
			}
		}
		d.mu.Unlock()
	}
}

// GetStats 获取统计信息
func (d *Deduplicator) GetStats() map[string]interface{} {
	var total int64
	d.db.Model(&NotifiedEvent{}).Count(&total)

	var todayCount int64
	d.db.Model(&NotifiedEvent{}).
		Where("created_at >= ?", time.Now().Format("2006-01-02")).
		Count(&todayCount)

	return map[string]interface{}{
		"total_records":  total,
		"today_count":    todayCount,
		"cache_size":     len(d.cache),
	}
}

// CleanByEventKey 删除指定事件的记录（允许重新通知）
func (d *Deduplicator) CleanByEventKey(eventKey string) error {
	d.mu.Lock()
	delete(d.cache, eventKey)
	d.mu.Unlock()

	return d.db.Where("event_key = ?", eventKey).Delete(&NotifiedEvent{}).Error
}

// CleanAll 清空所有记录
func (d *Deduplicator) CleanAll() error {
	d.mu.Lock()
	d.cache = make(map[string]time.Time)
	d.mu.Unlock()

	return d.db.Where("1 = 1").Delete(&NotifiedEvent{}).Error
}
