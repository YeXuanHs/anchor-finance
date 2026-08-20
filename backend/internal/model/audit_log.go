package model

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog 审计日志模型（移植自 zjmf activity_log）
type AuditLog struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	UserID       uint           `json:"user_id" gorm:"index"`                              // 操作用户 ID
	Username     string         `json:"username" gorm:"size:100"`                           // 用户名
	UserType     string         `json:"user_type" gorm:"size:20"`                           // 用户类型: admin, client, system
	Action       string         `json:"action" gorm:"size:100;index"`                       // 操作类型
	Description  string         `json:"description" gorm:"size:500"`                        // 操作描述
	Module       string         `json:"module" gorm:"size:50"`                              // 模块
	Controller   string         `json:"controller" gorm:"size:50"`                          // 控制器
	Method       string         `json:"method" gorm:"size:50"`                              // 方法
	IP           string         `json:"ip" gorm:"size:50"`                                  // IP 地址
	UserAgent    string         `json:"user_agent" gorm:"size:500"`                         // User Agent
	RequestData  string         `json:"request_data" gorm:"type:text"`                      // 请求数据（敏感信息已脱敏）
	ResponseCode int            `json:"response_code"`                                      // 响应码
	TargetType   string         `json:"target_type" gorm:"size:50"`                         // 目标类型
	TargetID     uint           `json:"target_id" gorm:"index"`                             // 目标 ID
	Duration     int64          `json:"duration"`                                           // 耗时（毫秒）
	Status       string         `json:"status" gorm:"size:20"`                              // 状态: success, failed
	Remark       string         `json:"remark" gorm:"type:text"`                            // 备注
	CreatedAt    time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditLogService 审计日志服务
type AuditLogService struct {
	db *gorm.DB
}

// NewAuditLogService 创建审计日志服务
func NewAuditLogService(db *gorm.DB) *AuditLogService {
	return &AuditLogService{db: db}
}

// Create 创建审计日志
func (s *AuditLogService) Create(log *AuditLog) error {
	return s.db.Create(log).Error
}

// List 查询审计日志列表
func (s *AuditLogService) List(userID uint, userType, action string, page, pageSize int) ([]AuditLog, int64, error) {
	var logs []AuditLog
	var total int64

	query := s.db.Model(&AuditLog{})

	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if userType != "" {
		query = query.Where("user_type = ?", userType)
	}
	if action != "" {
		query = query.Where("action = ?", action)
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

// GetByID 根据 ID 获取审计日志
func (s *AuditLogService) GetByID(id uint) (*AuditLog, error) {
	var log AuditLog
	if err := s.db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// GetStats 获取审计日志统计
func (s *AuditLogService) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总数
	var total int64
	if err := s.db.Model(&AuditLog{}).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = total

	// 今日数量
	var todayCount int64
	today := time.Now().Truncate(24 * time.Hour)
	if err := s.db.Model(&AuditLog{}).Where("created_at >= ?", today).Count(&todayCount).Error; err != nil {
		return nil, err
	}
	stats["today"] = todayCount

	// 按操作类型统计
	var actionStats []struct {
		Action string `json:"action"`
		Count  int64  `json:"count"`
	}
	if err := s.db.Model(&AuditLog{}).
		Select("action, count(*) as count").
		Group("action").
		Order("count DESC").
		Limit(10).
		Find(&actionStats).Error; err != nil {
		return nil, err
	}
	stats["by_action"] = actionStats

	// 按用户类型统计
	var userStats []struct {
		UserType string `json:"user_type"`
		Count    int64  `json:"count"`
	}
	if err := s.db.Model(&AuditLog{}).
		Select("user_type, count(*) as count").
		Group("user_type").
		Find(&userStats).Error; err != nil {
		return nil, err
	}
	stats["by_user_type"] = userStats

	return stats, nil
}

// CleanOldLogs 清理旧日志
func (s *AuditLogService) CleanOldLogs(days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := s.db.Where("created_at < ?", cutoff).Delete(&AuditLog{})
	return result.RowsAffected, result.Error
}
