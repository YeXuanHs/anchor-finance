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

// Delete deletes a specific API log entry.
func (s *APILogService) Delete(id uint) error {
	return s.db.Delete(&ApiLog{}, id).Error
}

// GetStats returns API log statistics for the given number of days.
func (s *APILogService) GetStats(days int) (map[string]interface{}, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	var total int64
	s.db.Model(&ApiLog{}).Where("created_at >= ?", cutoff).Count(&total)

	var successCount int64
	s.db.Model(&ApiLog{}).Where("created_at >= ? AND status_code >= 200 AND status_code < 300", cutoff).Count(&successCount)

	var errorCount int64
	s.db.Model(&ApiLog{}).Where("created_at >= ? AND status_code >= 400", cutoff).Count(&errorCount)

	var avgResponseTime float64
	s.db.Model(&ApiLog{}).Where("created_at >= ?", cutoff).Select("COALESCE(AVG(response_time), 0)").Scan(&avgResponseTime)

	return map[string]interface{}{
		"total":             total,
		"success_count":     successCount,
		"error_count":       errorCount,
		"avg_response_time": avgResponseTime,
	}, nil
}

// GetTopEndpoints returns top API endpoints by call count.
func (s *APILogService) GetTopEndpoints(limit, days int) ([]map[string]interface{}, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	var results []map[string]interface{}
	s.db.Model(&ApiLog{}).
		Select("endpoint, COUNT(*) as call_count, AVG(response_time) as avg_time").
		Where("created_at >= ?", cutoff).
		Group("endpoint").
		Order("call_count DESC").
		Limit(limit).
		Find(&results)
	return results, nil
}

// GetSlowRequests returns slowest API requests.
func (s *APILogService) GetSlowRequests(limit, days int) ([]ApiLog, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	var logs []ApiLog
	s.db.Where("created_at >= ?", cutoff).
		Order("response_time DESC").
		Limit(limit).
		Find(&logs)
	return logs, nil
}

// GetErrorRate returns API error rate statistics.
func (s *APILogService) GetErrorRate(days int) (map[string]interface{}, error) {
	cutoff := time.Now().AddDate(0, 0, -days)

	var total int64
	s.db.Model(&ApiLog{}).Where("created_at >= ?", cutoff).Count(&total)

	var errorCount int64
	s.db.Model(&ApiLog{}).Where("created_at >= ? AND status_code >= 400", cutoff).Count(&errorCount)

	rate := float64(0)
	if total > 0 {
		rate = float64(errorCount) / float64(total) * 100
	}

	return map[string]interface{}{
		"total":       total,
		"error_count": errorCount,
		"error_rate":  rate,
	}, nil
}
