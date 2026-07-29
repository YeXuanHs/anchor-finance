package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type AgentService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAgentService(db *gorm.DB, log *logger.Logger) *AgentService {
	return &AgentService{db: db, log: log}
}

func (s *AgentService) GetByID(id uint) (*model.Agent, error) {
	var agent model.Agent
	if err := s.db.First(&agent, id).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

func (s *AgentService) GetByUserID(userID uint) (*model.Agent, error) {
	var agent model.Agent
	if err := s.db.Where("user_id = ?", userID).First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

func (s *AgentService) GetList(page, pageSize int) ([]model.Agent, int64, error) {
	var items []model.Agent
	var total int64

	s.db.Model(&model.Agent{}).Count(&total)
	offset := (page - 1) * pageSize
	s.db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	return items, total, nil
}

func (s *AgentService) Create(agent *model.Agent) error {
	return s.db.Create(agent).Error
}

func (s *AgentService) Update(agent *model.Agent) error {
	return s.db.Save(agent).Error
}

func (s *AgentService) Delete(id uint) error {
	return s.db.Delete(&model.Agent{}, id).Error
}
