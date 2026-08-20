package service

import (
	"errors"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type UserLevelService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewUserLevelService(db *gorm.DB, log *logger.Logger) *UserLevelService {
	return &UserLevelService{db: db, log: log}
}

func (s *UserLevelService) GetByID(id uint) (*model.UserLevel, error) {
	var level model.UserLevel
	if err := s.db.First(&level, id).Error; err != nil {
		return nil, err
	}
	return &level, nil
}

func (s *UserLevelService) GetList() ([]model.UserLevel, error) {
	var levels []model.UserLevel
	err := s.db.Order("priority DESC, id ASC").Find(&levels).Error
	return levels, err
}

func (s *UserLevelService) Create(level *model.UserLevel) error {
	return s.db.Create(level).Error
}

func (s *UserLevelService) Update(level *model.UserLevel) error {
	return s.db.Save(level).Error
}

func (s *UserLevelService) Delete(id uint) error {
	result := s.db.Delete(&model.UserLevel{}, id)
	if result.RowsAffected == 0 {
		return errors.New("level not found")
	}
	return result.Error
}

func (s *UserLevelService) CalculateLevel(amount float64) (*model.UserLevel, error) {
	var level model.UserLevel
	err := s.db.Where("min_amount <= ?", amount).Order("min_amount DESC").First(&level).Error
	if err != nil {
		return nil, err
	}
	return &level, nil
}
