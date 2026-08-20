package service

import (
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type PublicService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewPublicService(db *gorm.DB, log *logger.Logger) *PublicService {
	return &PublicService{db: db, log: log}
}

// GetPublicConfig returns a public config value by key.
func (s *PublicService) GetPublicConfig(key string) (string, error) {
	var config struct {
		Value string
	}
	if err := s.db.Table("public_configs").Where("`key` = ?", key).Select("value").Scan(&config).Error; err != nil {
		return "", err
	}
	return config.Value, nil
}

// GetPublicConfigs returns all public configs for a group.
func (s *PublicService) GetPublicConfigs(group string) (map[string]string, error) {
	var configs []struct {
		Key   string
		Value string
	}
	query := s.db.Table("public_configs")
	if group != "" {
		query = query.Where("`group` = ?", group)
	}
	if err := query.Find(&configs).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, c := range configs {
		result[c.Key] = c.Value
	}
	return result, nil
}

// GetSystemInfo returns basic system info for public display.
func (s *PublicService) GetSystemInfo() (map[string]interface{}, error) {
	return map[string]interface{}{
		"version":  "1.0.0",
		"name":     "智简魔方",
		"powered":  true,
	}, nil
}
