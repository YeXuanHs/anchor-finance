package service

import (
	"context"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SystemService struct {
	db    *gorm.DB
	log   *logger.Logger
	redis *redis.Client
}

func NewSystemService(db *gorm.DB, log *logger.Logger, redis *redis.Client) *SystemService {
	return &SystemService{db: db, log: log, redis: redis}
}

// GetCommonInfo returns common system information.
func (s *SystemService) GetCommonInfo() (*model.SystemInfo, error) {
	var info model.SystemInfo
	if err := s.db.First(&info).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			info = model.SystemInfo{
				Version:     "1.0.0",
				LicenseType: 0,
			}
			return &info, nil
		}
		return nil, err
	}
	return &info, nil
}

// CheckUpdate checks for available system updates.
func (s *SystemService) CheckUpdate() (*model.SystemUpdate, error) {
	var update model.SystemUpdate
	if err := s.db.Where("status = 0").Order("id DESC").First(&update).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &update, nil
}

// GetUpdateContent returns update content.
func (s *SystemService) GetUpdateContent() (*model.SystemUpdate, error) {
	return s.CheckUpdate()
}

// GetSystemLog returns paginated system logs.
func (s *SystemService) GetSystemLog(page, pageSize int, level string, module string) ([]model.SystemLog, int64, error) {
	var logs []model.SystemLog
	var total int64

	query := s.db.Model(&model.SystemLog{})
	if level != "" {
		query = query.Where("level = ?", level)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// ClearCache clears system cache by flushing Redis.
func (s *SystemService) ClearCache() error {
	if s.redis != nil {
		if err := s.redis.FlushDB(context.Background()).Err(); err != nil {
			s.log.Errorf("redis flush failed: %v", err)
			return err
		}
		s.log.Info("redis cache cleared")
	}
	return nil
}

// GetUpdateList returns paginated system updates.
func (s *SystemService) GetUpdateList(page, pageSize int) ([]model.SystemUpdate, int64, error) {
	var updates []model.SystemUpdate
	var total int64

	query := s.db.Model(&model.SystemUpdate{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&updates).Error; err != nil {
		return nil, 0, err
	}
	return updates, total, nil
}

// InstallUpdate marks an update as installed.
func (s *SystemService) InstallUpdate(id uint) error {
	return s.db.Model(&model.SystemUpdate{}).Where("id = ?", id).
		Update("status", 1).Error
}
