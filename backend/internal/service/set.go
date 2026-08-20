package service

import (
	"errors"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// SetService 系统设置服务
type SetService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewSetService(db *gorm.DB, log *logger.Logger) *SetService {
	return &SetService{db: db, log: log}
}

// Get 获取单个设置
func (s *SetService) Get(code string) (string, error) {
	var config model.ConfigOption
	if err := s.db.Where("code = ?", code).First(&config).Error; err != nil {
		return "", err
	}
	return config.Value, nil
}

// Set 设置单个设置
func (s *SetService) Set(code, value string) error {
	var config model.ConfigOption
	err := s.db.Where("code = ?", code).First(&config).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		config = model.ConfigOption{
			Code:  code,
			Value: value,
			Group: "general",
			Name:  code,
			Type:  "text",
		}
		return s.db.Create(&config).Error
	}
	if err != nil {
		return err
	}
	return s.db.Model(&config).Update("value", value).Error
}

// GetGroup 获取分组设置
func (s *SetService) GetGroup(group string) (map[string]string, error) {
	var configs []model.ConfigOption
	if err := s.db.Where("`group` = ?", group).Find(&configs).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string, len(configs))
	for _, c := range configs {
		result[c.Code] = c.Value
	}
	return result, nil
}

// SetGroup 批量设置分组
func (s *SetService) SetGroup(group string, settings map[string]string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for code, value := range settings {
			var config model.ConfigOption
			err := tx.Where("code = ?", code).First(&config).Error
			if errors.Is(err, gorm.ErrRecordNotFound) {
				config = model.ConfigOption{
					Code:  code,
					Value: value,
					Group: group,
					Name:  code,
					Type:  "text",
				}
				if err := tx.Create(&config).Error; err != nil {
					return err
				}
				continue
			}
			if err != nil {
				return err
			}
			if err := tx.Model(&config).Updates(map[string]interface{}{"value": value, "group": group}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetAll 获取所有设置
func (s *SetService) GetAll() (map[string]string, error) {
	var configs []model.ConfigOption
	if err := s.db.Find(&configs).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string, len(configs))
	for _, c := range configs {
		result[c.Code] = c.Value
	}
	return result, nil
}

// Reset 重置设置为默认值
func (s *SetService) Reset(code string) error {
	var config model.ConfigOption
	if err := s.db.Where("code = ?", code).First(&config).Error; err != nil {
		return err
	}
	if config.DefaultValue == "" {
		return s.db.Delete(&config).Error
	}
	return s.db.Model(&config).Update("value", config.DefaultValue).Error
}
