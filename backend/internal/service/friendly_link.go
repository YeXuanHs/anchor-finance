package service

import (
	"errors"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type FriendlyLinkService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewFriendlyLinkService(db *gorm.DB, log *logger.Logger) *FriendlyLinkService {
	return &FriendlyLinkService{db: db, log: log}
}

func (s *FriendlyLinkService) Create(link *model.FriendlyLink) error {
	return s.db.Create(link).Error
}

func (s *FriendlyLinkService) Update(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&model.FriendlyLink{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		return errors.New("link not found")
	}
	return result.Error
}

func (s *FriendlyLinkService) Delete(id uint) error {
	result := s.db.Delete(&model.FriendlyLink{}, id)
	if result.RowsAffected == 0 {
		return errors.New("link not found")
	}
	return result.Error
}

func (s *FriendlyLinkService) GetList(page, pageSize int, status int, group string) ([]model.FriendlyLink, int64, error) {
	var items []model.FriendlyLink
	var total int64

	query := s.db.Model(&model.FriendlyLink{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	if group != "" {
		query = query.Where("`group` = ?", group)
	}

	query.Count(&total)
	offset, limit := util.Paginate(page, pageSize)
	query.Offset(offset).Limit(limit).Order("sort_order ASC, id DESC").Find(&items)
	return items, total, nil
}

func (s *FriendlyLinkService) GetActive(group string) ([]model.FriendlyLink, error) {
	var items []model.FriendlyLink
	query := s.db.Where("status = 1")
	if group != "" {
		query = query.Where("`group` = ?", group)
	}
	query.Order("sort_order ASC, id DESC").Find(&items)
	return items, nil
}
