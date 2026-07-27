package service

import (
	"errors"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type BannerService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewBannerService(db *gorm.DB, log *logger.Logger) *BannerService {
	return &BannerService{db: db, log: log}
}

type CreateBannerRequest struct {
	Title       string  `json:"title" binding:"required,max=256"`
	Description string  `json:"description" binding:"max=512"`
	Type        string  `json:"type" binding:"required,oneof=image video"`
	MediaURL    string  `json:"media_url" binding:"required,max=512"`
	LinkURL     string  `json:"link_url" binding:"max=512"`
	ButtonText  string  `json:"button_text" binding:"max=64"`
	OpenNew     bool    `json:"open_new"`
	SortOrder   int     `json:"sort_order"`
	Status      int16   `json:"status"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
	Position    string  `json:"position" binding:"max=32"`
}

type UpdateBannerRequest struct {
	Title       *string `json:"title" binding:"omitempty,max=256"`
	Description *string `json:"description" binding:"omitempty,max=512"`
	Type        *string `json:"type" binding:"omitempty,oneof=image video"`
	MediaURL    *string `json:"media_url" binding:"omitempty,max=512"`
	LinkURL     *string `json:"link_url" binding:"omitempty,max=512"`
	ButtonText  *string `json:"button_text" binding:"omitempty,max=64"`
	OpenNew     *bool   `json:"open_new"`
	SortOrder   *int    `json:"sort_order"`
	Status      *int16  `json:"status"`
	StartTime   *string `json:"start_time"`
	EndTime     *string `json:"end_time"`
	Position    *string `json:"position" binding:"omitempty,max=32"`
}

func (s *BannerService) Create(req CreateBannerRequest) (*model.Banner, error) {
	banner := model.Banner{
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		MediaURL:    req.MediaURL,
		LinkURL:     req.LinkURL,
		ButtonText:  req.ButtonText,
		OpenNew:     req.OpenNew,
		SortOrder:   req.SortOrder,
		Status:      req.Status,
		Position:    req.Position,
	}
	if banner.Status == 0 {
		banner.Status = 1
	}
	if banner.Position == "" {
		banner.Position = "home"
	}

	if req.StartTime != nil && *req.StartTime != "" {
		t, err := util.ParseTime(*req.StartTime)
		if err != nil {
			return nil, errors.New("invalid start_time format")
		}
		banner.StartTime = &t
	}
	if req.EndTime != nil && *req.EndTime != "" {
		t, err := util.ParseTime(*req.EndTime)
		if err != nil {
			return nil, errors.New("invalid end_time format")
		}
		banner.EndTime = &t
	}

	if err := s.db.Create(&banner).Error; err != nil {
		return nil, err
	}
	s.log.Infof("banner created: id=%d title=%s", banner.ID, banner.Title)
	return &banner, nil
}

func (s *BannerService) Update(id uint, req UpdateBannerRequest) (*model.Banner, error) {
	var banner model.Banner
	if err := s.db.First(&banner, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.MediaURL != nil {
		updates["media_url"] = *req.MediaURL
	}
	if req.LinkURL != nil {
		updates["link_url"] = *req.LinkURL
	}
	if req.ButtonText != nil {
		updates["button_text"] = *req.ButtonText
	}
	if req.OpenNew != nil {
		updates["open_new"] = *req.OpenNew
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
	if req.StartTime != nil {
		if *req.StartTime == "" {
			updates["start_time"] = nil
		} else {
			t, err := util.ParseTime(*req.StartTime)
			if err != nil {
				return nil, errors.New("invalid start_time format")
			}
			updates["start_time"] = t
		}
	}
	if req.EndTime != nil {
		if *req.EndTime == "" {
			updates["end_time"] = nil
		} else {
			t, err := util.ParseTime(*req.EndTime)
			if err != nil {
				return nil, errors.New("invalid end_time format")
			}
			updates["end_time"] = t
		}
	}

	if len(updates) > 0 {
		if err := s.db.Model(&banner).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&banner, id).Error; err != nil {
		return nil, err
	}
	s.log.Infof("banner updated: id=%d", id)
	return &banner, nil
}

func (s *BannerService) Delete(id uint) error {
	result := s.db.Delete(&model.Banner{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("banner not found")
	}
	s.log.Infof("banner deleted: id=%d", id)
	return nil
}

func (s *BannerService) GetByID(id uint) (*model.Banner, error) {
	var banner model.Banner
	if err := s.db.First(&banner, id).Error; err != nil {
		return nil, err
	}
	return &banner, nil
}

func (s *BannerService) GetList(page, pageSize int, position string, status int) ([]model.Banner, int64, error) {
	var banners []model.Banner
	var total int64

	query := s.db.Model(&model.Banner{})
	if position != "" {
		query = query.Where("position = ?", position)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := util.Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("sort_order ASC, id DESC").
		Find(&banners).Error; err != nil {
		return nil, 0, err
	}

	return banners, total, nil
}

// GetActive returns currently active banners for frontend display.
func (s *BannerService) GetActive(position string) ([]model.Banner, error) {
	var banners []model.Banner
	now := time.Now()

	query := s.db.Where("status = 1")
	if position != "" {
		query = query.Where("position = ?", position)
	}
	query = query.Where("(start_time IS NULL OR start_time <= ?)", now)
	query = query.Where("(end_time IS NULL OR end_time >= ?)", now)

	if err := query.Order("sort_order ASC, id DESC").Find(&banners).Error; err != nil {
		return nil, err
	}
	return banners, nil
}

func (s *BannerService) UpdateSort(id uint, sortOrder int) error {
	result := s.db.Model(&model.Banner{}).Where("id = ?", id).Update("sort_order", sortOrder)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("banner not found")
	}
	return nil
}

func (s *BannerService) ToggleStatus(id uint) error {
	var banner model.Banner
	if err := s.db.First(&banner, id).Error; err != nil {
		return err
	}
	newStatus := int16(1)
	if banner.Status == 1 {
		newStatus = 0
	}
	return s.db.Model(&banner).Update("status", newStatus).Error
}

func (s *BannerService) IncrementClick(id uint) error {
	return s.db.Model(&model.Banner{}).Where("id = ?", id).
		Update("click_count", gorm.Expr("click_count + 1")).Error
}
