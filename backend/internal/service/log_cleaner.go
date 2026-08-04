package service

import (
	"fmt"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// LogCleaner 日志清理服务
// 支持多种清理策略，解决日志积压问题
type LogCleaner struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewLogCleaner 创建日志清理服务
func NewLogCleaner(db *gorm.DB, log *logger.Logger) *LogCleaner {
	return &LogCleaner{db: db, log: log}
}

// CleanByDays 按天数清理日志
func (lc *LogCleaner) CleanByDays(days int) (int64, error) {
	if days < 1 {
		return 0, fmt.Errorf("天数必须大于0")
	}

	expire := time.Now().AddDate(0, 0, -days)
	result := lc.db.Where("created_at < ?", expire).Delete(&AuditLog{})
	if result.Error != nil {
		return 0, result.Error
	}

	lc.log.Infof("清理 %d 天前的日志，共 %d 条", days, result.RowsAffected)
	return result.RowsAffected, nil
}

// CleanByCount 按数量保留日志（保留最新的N条）
func (lc *LogCleaner) CleanByCount(keepCount int) (int64, error) {
	if keepCount < 1 {
		return 0, fmt.Errorf("保留数量必须大于0")
	}

	var total int64
	lc.db.Model(&AuditLog{}).Count(&total)

	if total <= int64(keepCount) {
		return 0, nil
	}

	// 找到第 keepCount 条的 ID
	var lastID uint
	lc.db.Model(&AuditLog{}).
		Order("id DESC").
		Offset(keepCount - 1).
		Limit(1).
		Pluck("id", &lastID)

	// 删除比这个 ID 更早的日志
	result := lc.db.Where("id < ?", lastID).Delete(&AuditLog{})
	if result.Error != nil {
		return 0, result.Error
	}

	lc.log.Infof("保留最新 %d 条日志，清理 %d 条", keepCount, result.RowsAffected)
	return result.RowsAffected, nil
}

// CleanByModule 按模块清理日志
func (lc *LogCleaner) CleanByModule(module string) (int64, error) {
	result := lc.db.Where("module = ?", module).Delete(&AuditLog{})
	if result.Error != nil {
		return 0, result.Error
	}

	lc.log.Infof("清理模块 %s 的日志，共 %d 条", module, result.RowsAffected)
	return result.RowsAffected, nil
}

// CleanByStatus 按状态清理日志
func (lc *LogCleaner) CleanByStatus(status string) (int64, error) {
	result := lc.db.Where("status = ?", status).Delete(&AuditLog{})
	if result.Error != nil {
		return 0, result.Error
	}

	lc.log.Infof("清理状态 %s 的日志，共 %d 条", status, result.RowsAffected)
	return result.RowsAffected, nil
}

// CleanBySize 按大小清理日志（保留指定大小）
func (lc *LogCleaner) CleanBySize(keepSizeMB int) (int64, error) {
	if keepSizeMB < 1 {
		return 0, fmt.Errorf("保留大小必须大于0")
	}

	// 估算当前日志大小
	var count int64
	lc.db.Model(&AuditLog{}).Count(&count)

	// 假设每条日志约 1KB
	estimatedSizeMB := count / 1024
	if estimatedSizeMB <= int64(keepSizeMB) {
		return 0, nil
	}

	// 计算需要删除的数量
	deleteCount := count - int64(keepSizeMB)*1024

	// 找到要删除的 ID 范围
	var lastID uint
	lc.db.Model(&AuditLog{}).
		Order("id DESC").
		Offset(int(deleteCount) - 1).
		Limit(1).
		Pluck("id", &lastID)

	result := lc.db.Where("id < ?", lastID).Delete(&AuditLog{})
	if result.Error != nil {
		return 0, result.Error
	}

	lc.log.Infof("保留 %dMB 日志，清理 %d 条", keepSizeMB, result.RowsAffected)
	return result.RowsAffected, nil
}

// CleanExpired 清理过期日志（根据配置的保留天数）
func (lc *LogCleaner) CleanExpired() (int64, error) {
	// 从配置中获取保留天数
	keepDays := 30 // 默认30天

	var setting struct {
		Value string
	}
	lc.db.Table("system_configs").
		Where("`key` = ?", "audit_log_retention_days").
		First(&setting)

	if setting.Value != "" {
		fmt.Sscanf(setting.Value, "%d", &keepDays)
	}

	return lc.CleanByDays(keepDays)
}

// GetLogStats 获取日志统计
func (lc *LogCleaner) GetLogStats() map[string]interface{} {
	var total int64
	lc.db.Model(&AuditLog{}).Count(&total)

	// 按模块统计
	var moduleStats []struct {
		Module string
		Count  int64
	}
	lc.db.Model(&AuditLog{}).
		Select("module, COUNT(*) as count").
		Group("module").
		Find(&moduleStats)

	// 按状态统计
	var statusStats []struct {
		Status string
		Count  int64
	}
	lc.db.Model(&AuditLog{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Find(&statusStats)

	// 最近7天的日志数量
	var recentCount int64
	lc.db.Model(&AuditLog{}).
		Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).
		Count(&recentCount)

	return map[string]interface{}{
		"total":         total,
		"module_stats":  moduleStats,
		"status_stats":  statusStats,
		"recent_7_days": recentCount,
	}
}

// CleanupScheduler 定时清理调度器
type CleanupScheduler struct {
	cleaner *LogCleaner
	log     *logger.Logger
}

// NewCleanupScheduler 创建清理调度器
func NewCleanupScheduler(cleaner *LogCleaner, log *logger.Logger) *CleanupScheduler {
	return &CleanupScheduler{cleaner: cleaner, log: log}
}

// Start 启动定时清理
func (cs *CleanupScheduler) Start() {
	// 每天凌晨3点执行清理
	go func() {
		for {
			now := time.Now()
			next := time.Date(now.Year(), now.Month(), now.Day()+1, 3, 0, 0, 0, now.Location())
			duration := next.Sub(now)

			time.Sleep(duration)

			cs.log.Info("开始执行定时日志清理...")
			count, err := cs.cleaner.CleanExpired()
			if err != nil {
				cs.log.Errorf("定时清理失败: %v", err)
			} else {
				cs.log.Infof("定时清理完成，清理 %d 条日志", count)
			}
		}
	}()
}
