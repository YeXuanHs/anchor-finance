package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type ClientCareService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewClientCareService(db *gorm.DB, log *logger.Logger) *ClientCareService {
	return &ClientCareService{db: db, log: log}
}

// GetRules returns all client care rules with pagination.
func (s *ClientCareService) GetRules(page, pageSize int, careType string, keyword string) ([]model.ClientCareRule, int64, error) {
	var items []model.ClientCareRule
	var total int64

	query := s.db.Model(&model.ClientCareRule{})
	if careType != "" {
		query = query.Where("type = ?", careType)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// GetRuleByID returns a single care rule by ID.
func (s *ClientCareService) GetRuleByID(id uint) (*model.ClientCareRule, error) {
	var rule model.ClientCareRule
	if err := s.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// CreateRule creates a new client care rule.
func (s *ClientCareService) CreateRule(rule *model.ClientCareRule) error {
	return s.db.Create(rule).Error
}

// UpdateRule updates a client care rule.
func (s *ClientCareService) UpdateRule(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.ClientCareRule{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteRule deletes a client care rule and its related logs.
func (s *ClientCareService) DeleteRule(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("rule_id = ?", id).Delete(&model.ClientCareLog{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.ClientCareRule{}, id).Error
	})
}

// GetLogs returns client care logs with pagination.
func (s *ClientCareService) GetLogs(page, pageSize int, ruleID *uint, userID *uint) ([]model.ClientCareLog, int64, error) {
	var logs []model.ClientCareLog
	var total int64

	query := s.db.Model(&model.ClientCareLog{})
	if ruleID != nil {
		query = query.Where("rule_id = ?", *ruleID)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// CreateLog creates a new care log entry.
func (s *ClientCareService) CreateLog(log *model.ClientCareLog) error {
	return s.db.Create(log).Error
}
