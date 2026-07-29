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

func (s *ClientCareService) GetList(page, pageSize int, careType string) ([]model.ClientCare, int64, error) {
	var items []model.ClientCare
	var total int64

	query := s.db.Model(&model.ClientCare{})
	if careType != "" {
		query = query.Where("type = ?", careType)
	}

	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	return items, total, nil
}

func (s *ClientCareService) Create(care *model.ClientCare) error {
	return s.db.Create(care).Error
}

func (s *ClientCareService) Update(care *model.ClientCare) error {
	return s.db.Save(care).Error
}

func (s *ClientCareService) Delete(id uint) error {
	return s.db.Delete(&model.ClientCare{}, id).Error
}
