package service

import (
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type APILogService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewApiLogService(db *gorm.DB, log *logger.Logger) *APILogService {
	return &APILogService{db: db, log: log}
}

type ApiLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	APIKey    string    `gorm:"type:varchar(64);index" json:"api_key"`
	Method    string    `gorm:"type:varchar(8)" json:"method"`
	Path      string    `gorm:"type:varchar(256)" json:"path"`
	IP        string    `gorm:"type:varchar(45)" json:"ip"`
	Status    int       `json:"status"`
	Latency   int64     `json:"latency"` // ms
	Request   string    `gorm:"type:text" json:"request"`
	Response  string    `gorm:"type:text" json:"response"`
	CreatedAt time.Time `json:"created_at"`
}

func (s *APILogService) List(page, pageSize int, userID uint, method string, startTime, endTime *time.Time) ([]ApiLog, int64, error) {
	var items []ApiLog
	var total int64
	query := s.db.Model(&ApiLog{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if method != "" {
		query = query.Where("method = ?", method)
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

func (s *APILogService) GetByID(id uint) (*ApiLog, error) {
	var item ApiLog
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *APILogService) Cleanup(days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := s.db.Where("created_at < ?", cutoff).Delete(&ApiLog{})
	return result.RowsAffected, result.Error
}
