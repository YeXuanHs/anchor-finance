package service

import (
	"fmt"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type SystemLogService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewSystemLogService(db *gorm.DB, log *logger.Logger) *SystemLogService {
	return &SystemLogService{db: db, log: log}
}

func (s *SystemLogService) Record(level, module, message string, userID uint, ip string, details string) {
	entry := model.SystemLog{
		Level: level, Module: module, Message: message,
		UserID: userID, IP: ip, Details: details,
	}
	s.db.Create(&entry)
}

func (s *SystemLogService) List(page, pageSize int, level, module, keyword string, startTime, endTime *time.Time) ([]model.SystemLog, int64, error) {
	var items []model.SystemLog
	var total int64
	query := s.db.Model(&model.SystemLog{})
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if keyword != "" {
		query = query.Where("message LIKE ?", "%"+keyword+"%")
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

func (s *SystemLogService) GetByID(id uint) (*model.SystemLog, error) {
	var item model.SystemLog
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *SystemLogService) Delete(id uint) error {
	return s.db.Delete(&model.SystemLog{}, id).Error
}

func (s *SystemLogService) Cleanup(days int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	result := s.db.Where("created_at < ?", cutoff).Delete(&model.SystemLog{})
	return result.RowsAffected, result.Error
}

func (s *SystemLogService) GetStats(days int) (map[string]interface{}, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	var total int64
	s.db.Model(&model.SystemLog{}).Where("created_at >= ?", cutoff).Count(&total)
	return map[string]interface{}{"total": total, "days": days}, nil
}

func (s *SystemLogService) GetLevelStats(days int) ([]map[string]interface{}, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	var results []map[string]interface{}
	s.db.Model(&model.SystemLog{}).Where("created_at >= ?", cutoff).
		Select("level, COUNT(*) as count").Group("level").Scan(&results)
	return results, nil
}

func (s *SystemLogService) GetModuleStats(days int) ([]map[string]interface{}, error) {
	cutoff := time.Now().AddDate(0, 0, -days)
	var results []map[string]interface{}
	s.db.Model(&model.SystemLog{}).Where("created_at >= ?", cutoff).
		Select("module, COUNT(*) as count").Group("module").Scan(&results)
	return results, nil
}

func (s *SystemLogService) Export(level, module string, startTime, endTime *time.Time) (string, error) {
	var logs []model.SystemLog
	query := s.db.Model(&model.SystemLog{})
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", *startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", *endTime)
	}
	query.Order("created_at DESC").Limit(10000).Find(&logs)

	var buf strings.Builder
	buf.WriteString("ID,Level,Module,Message,IP,User,Data,CreatedAt\n")
	for _, l := range logs {
		msg := strings.ReplaceAll(l.Message, "\"", "\"\"")
		data := strings.ReplaceAll(l.Data, "\"", "\"\"")
		buf.WriteString(fmt.Sprintf("%d,%s,%s,\"%s\",%s,%s,\"%s\",%s\n",
			l.ID, l.Level, l.Module, msg, l.IP, l.User, data, l.CreatedAt.Format("2006-01-02 15:04:05")))
	}
	return buf.String(), nil
}

func (s *SystemLogService) ClearByLevel(level string) (int64, error) {
	result := s.db.Where("level = ?", level).Delete(&model.SystemLog{})
	return result.RowsAffected, result.Error
}
