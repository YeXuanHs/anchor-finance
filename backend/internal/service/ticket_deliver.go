package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type TicketDeliverService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewTicketDeliverService(db *gorm.DB, log *logger.Logger) *TicketDeliverService {
	return &TicketDeliverService{db: db, log: log}
}

// GetAddPage returns data needed for adding a deliver rule.
func (s *TicketDeliverService) GetAddPage() (map[string]interface{}, error) {
	return map[string]interface{}{
		"departments": []interface{}{},
		"products":    []interface{}{},
	}, nil
}

// GetRules returns all deliver rules.
func (s *TicketDeliverService) GetRules(page, pageSize int, keyword string) ([]model.TicketDeliverRule, int64, error) {
	var rules []model.TicketDeliverRule
	var total int64

	query := s.db.Model(&model.TicketDeliverRule{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("priority DESC, id ASC").Find(&rules).Error; err != nil {
		return nil, 0, err
	}
	return rules, total, nil
}

// GetRuleByID returns a single rule by ID.
func (s *TicketDeliverService) GetRuleByID(id uint) (*model.TicketDeliverRule, error) {
	var rule model.TicketDeliverRule
	if err := s.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// CreateRule creates a new deliver rule.
func (s *TicketDeliverService) CreateRule(rule *model.TicketDeliverRule) error {
	return s.db.Create(rule).Error
}

// UpdateRule updates a deliver rule.
func (s *TicketDeliverService) UpdateRule(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.TicketDeliverRule{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteRule deletes a deliver rule.
func (s *TicketDeliverService) DeleteRule(id uint) error {
	return s.db.Delete(&model.TicketDeliverRule{}, id).Error
}

// GetLogs returns deliver logs for a ticket.
func (s *TicketDeliverService) GetLogs(ticketID uint, page, pageSize int) ([]model.TicketDeliverLog, int64, error) {
	var logs []model.TicketDeliverLog
	var total int64

	query := s.db.Model(&model.TicketDeliverLog{})
	if ticketID > 0 {
		query = query.Where("ticket_id = ?", ticketID)
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
