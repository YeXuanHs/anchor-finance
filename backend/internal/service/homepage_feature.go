package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type HomepageFeatureService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewHomepageFeatureService(db *gorm.DB, log *logger.Logger) *HomepageFeatureService {
	return &HomepageFeatureService{db: db, log: log}
}

type CreateFeatureRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	LinkURL     string `json:"link_url"`
	SortOrder   int    `json:"sort_order"`
	Status      *int16 `json:"status"`
	Position    string `json:"position"`
}

type UpdateFeatureRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	LinkURL     *string `json:"link_url"`
	SortOrder   *int    `json:"sort_order"`
	Status      *int16  `json:"status"`
	Position    *string `json:"position"`
}

func (s *HomepageFeatureService) GetList(page, pageSize int, position string, status int) ([]model.HomepageFeature, int64, error) {
	var items []model.HomepageFeature
	var total int64

	query := s.db.Model(&model.HomepageFeature{})
	if position != "" {
		query = query.Where("position = ?", position)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Order("sort_order ASC, id DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (s *HomepageFeatureService) GetByID(id uint) (*model.HomepageFeature, error) {
	var item model.HomepageFeature
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *HomepageFeatureService) GetActive(position string) ([]model.HomepageFeature, error) {
	var items []model.HomepageFeature
	query := s.db.Where("status = 1")
	if position != "" {
		query = query.Where("position = ?", position)
	}
	if err := query.Order("sort_order ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *HomepageFeatureService) Create(req CreateFeatureRequest) (*model.HomepageFeature, error) {
	status := int16(1)
	if req.Status != nil {
		status = *req.Status
	}
	position := req.Position
	if position == "" {
		position = "home"
	}

	item := &model.HomepageFeature{
		Title:       req.Title,
		Description: req.Description,
		Icon:        req.Icon,
		LinkURL:     req.LinkURL,
		SortOrder:   req.SortOrder,
		Status:      status,
		Position:    position,
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (s *HomepageFeatureService) Update(id uint, req UpdateFeatureRequest) (*model.HomepageFeature, error) {
	var item model.HomepageFeature
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.LinkURL != nil {
		updates["link_url"] = *req.LinkURL
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Position != nil {
		updates["position"] = *req.Position
	}

	if err := s.db.Model(&item).Updates(updates).Error; err != nil {
		return nil, err
	}

	return &item, nil
}

func (s *HomepageFeatureService) Delete(id uint) error {
	return s.db.Delete(&model.HomepageFeature{}, id).Error
}

func (s *HomepageFeatureService) ToggleStatus(id uint) error {
	var item model.HomepageFeature
	if err := s.db.First(&item, id).Error; err != nil {
		return err
	}

	newStatus := int16(1)
	if item.Status == 1 {
		newStatus = 0
	}

	return s.db.Model(&item).Update("status", newStatus).Error
}

func (s *HomepageFeatureService) UpdateSort(id uint, sortOrder int) error {
	return s.db.Model(&model.HomepageFeature{}).Where("id = ?", id).Update("sort_order", sortOrder).Error
}
