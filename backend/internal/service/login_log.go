package service

import (
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type LoginLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"type:varchar(64)" json:"username"`
	IP        string    `gorm:"type:varchar(45)" json:"ip"`
	UserAgent string    `gorm:"type:varchar(256)" json:"user_agent"`
	Status    string    `gorm:"type:varchar(16)" json:"status"` // success/failed
	Reason    string    `gorm:"type:varchar(128)" json:"reason"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

type LoginLogService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewLoginLogService(db *gorm.DB, log *logger.Logger) *LoginLogService {
	return &LoginLogService{db: db, log: log}
}

func (s *LoginLogService) Record(userID uint, username, ip, userAgent, status, reason string) error {
	entry := LoginLog{
		UserID: userID, Username: username, IP: ip,
		UserAgent: userAgent, Status: status, Reason: reason,
	}
	return s.db.Create(&entry).Error
}

func (s *LoginLogService) List(page, pageSize int, userID uint, status string, startTime, endTime *time.Time) ([]LoginLog, int64, error) {
	var items []LoginLog
	var total int64
	query := s.db.Model(&LoginLog{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", endTime)
	}
	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	return items, total, nil
}

func (s *LoginLogService) Cleanup(days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := s.db.Where("created_at < ?", cutoff).Delete(&LoginLog{})
	return result.RowsAffected, result.Error
}
