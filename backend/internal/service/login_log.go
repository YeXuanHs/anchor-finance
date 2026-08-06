package service

import (
	"errors"
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

func (s *LoginLogService) List(page, pageSize int, userID *uint, username, ip, status string, startTime, endTime *time.Time) ([]LoginLog, int64, error) {
	var items []LoginLog
	var total int64
	query := s.db.Model(&LoginLog{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if ip != "" {
		query = query.Where("ip LIKE ?", "%"+ip+"%")
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

// GetByID returns a login log by ID.
func (s *LoginLogService) GetByID(id uint) (*LoginLog, error) {
	var item LoginLog
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Delete deletes a login log by ID.
func (s *LoginLogService) Delete(id uint) error {
	result := s.db.Delete(&LoginLog{}, id)
	if result.RowsAffected == 0 {
		return errors.New("login log not found")
	}
	return result.Error
}

// GetStats returns login statistics for the given number of days.
func (s *LoginLogService) GetStats(days int) (map[string]interface{}, error) {
	cutoff := time.Now().AddDate(0, 0, -days)

	var totalLogins int64
	s.db.Model(&LoginLog{}).Where("created_at >= ?", cutoff).Count(&totalLogins)

	var successLogins int64
	s.db.Model(&LoginLog{}).Where("created_at >= ? AND status = ?", cutoff, "success").Count(&successLogins)

	var failedLogins int64
	s.db.Model(&LoginLog{}).Where("created_at >= ? AND status = ?", cutoff, "failed").Count(&failedLogins)

	var uniqueUsers int64
	s.db.Model(&LoginLog{}).Where("created_at >= ?", cutoff).Distinct("user_id").Count(&uniqueUsers)

	return map[string]interface{}{
		"total_logins":   totalLogins,
		"success_logins": successLogins,
		"failed_logins":  failedLogins,
		"unique_users":   uniqueUsers,
		"days":           days,
	}, nil
}

// GetFailedAttempts returns failed login attempts count.
func (s *LoginLogService) GetFailedAttempts(userID *uint, ip string, minutes int) (int64, error) {
	cutoff := time.Now().Add(-time.Duration(minutes) * time.Minute)
	query := s.db.Model(&LoginLog{}).Where("created_at >= ? AND status = ?", cutoff, "failed")
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if ip != "" {
		query = query.Where("ip = ?", ip)
	}
	var count int64
	err := query.Count(&count).Error
	return count, err
}

// Export exports login logs to a CSV file.
func (s *LoginLogService) Export(username, ip string, startTime, endTime *time.Time) (string, error) {
	// Placeholder - implement CSV export
	return "", errors.New("export not implemented")
}
