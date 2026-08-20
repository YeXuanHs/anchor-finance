package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type TicketStatusService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewTicketStatusService(db *gorm.DB, log *logger.Logger) *TicketStatusService {
	return &TicketStatusService{db: db, log: log}
}

// GetStatuses returns all ticket statuses.
func (s *TicketStatusService) GetStatuses() ([]model.TicketStatus, error) {
	var statuses []model.TicketStatus
	if err := s.db.Where("status = 1").Order("`order` ASC, id ASC").Find(&statuses).Error; err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetAllStatuses returns all ticket statuses including disabled.
func (s *TicketStatusService) GetAllStatuses() ([]model.TicketStatus, error) {
	var statuses []model.TicketStatus
	if err := s.db.Order("`order` ASC, id ASC").Find(&statuses).Error; err != nil {
		return nil, err
	}
	return statuses, nil
}

// GetByID returns a single status by ID.
func (s *TicketStatusService) GetByID(id uint) (*model.TicketStatus, error) {
	var status model.TicketStatus
	if err := s.db.First(&status, id).Error; err != nil {
		return nil, err
	}
	return &status, nil
}

// AddStatus adds a new ticket status.
func (s *TicketStatusService) AddStatus(status *model.TicketStatus) error {
	return s.db.Create(status).Error
}

// UpdateStatus updates a ticket status.
func (s *TicketStatusService) UpdateStatus(id uint, updates map[string]interface{}) error {
	var ts model.TicketStatus
	if err := s.db.First(&ts, id).Error; err != nil {
		return err
	}
	if ts.IsSystem {
		return gorm.ErrInvalidData
	}
	return s.db.Model(&ts).Updates(updates).Error
}

// DeleteStatus deletes a ticket status.
func (s *TicketStatusService) DeleteStatus(id uint) error {
	var ts model.TicketStatus
	if err := s.db.First(&ts, id).Error; err != nil {
		return err
	}
	if ts.IsSystem {
		return gorm.ErrInvalidData
	}
	return s.db.Delete(&ts).Error
}

// GetDefaultStatuses returns default ticket statuses.
func (s *TicketStatusService) GetDefaultStatuses() []string {
	return []string{"Open", "Answered", "CustomerReply", "Closed"}
}

// GetStatusByCode returns a status by its code.
func (s *TicketStatusService) GetStatusByCode(code string) (*model.TicketStatus, error) {
	var status model.TicketStatus
	if err := s.db.Where("code = ?", code).First(&status).Error; err != nil {
		return nil, err
	}
	return &status, nil
}

// GetStatusByID returns a status by its ID (alias for GetByID).
func (s *TicketStatusService) GetStatusByID(id uint) (*model.TicketStatus, error) {
	return s.GetByID(id)
}
