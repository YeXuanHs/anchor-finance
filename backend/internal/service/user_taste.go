package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type UserTasteService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewUserTasteService(db *gorm.DB, log *logger.Logger) *UserTasteService {
	return &UserTasteService{db: db, log: log}
}

// GetByUserID returns user taste settings by user ID.
func (s *UserTasteService) GetByUserID(userID uint) (*model.UserTaste, error) {
	var taste model.UserTaste
	if err := s.db.Where("user_id = ?", userID).First(&taste).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return default taste
			return &model.UserTaste{
				UserID:   userID,
				Theme:    "default",
				Language: "zh-CN",
				Layout:   "default",
			}, nil
		}
		return nil, err
	}
	return &taste, nil
}

// Update updates user taste settings.
func (s *UserTasteService) Update(userID uint, updates map[string]interface{}) error {
	var taste model.UserTaste
	err := s.db.Where("user_id = ?", userID).First(&taste).Error
	if err == gorm.ErrRecordNotFound {
		// Create new taste
		taste = model.UserTaste{
			UserID:   userID,
			Theme:    "default",
			Language: "zh-CN",
			Layout:   "default",
		}
		if theme, ok := updates["theme"]; ok {
			taste.Theme = theme.(string)
		}
		if lang, ok := updates["language"]; ok {
			taste.Language = lang.(string)
		}
		if layout, ok := updates["layout"]; ok {
			taste.Layout = layout.(string)
		}
		if settings, ok := updates["settings"]; ok {
			taste.Settings = toJSONMap(settings)
		}
		return s.db.Create(&taste).Error
	}
	if err != nil {
		return err
	}
	return s.db.Model(&taste).Updates(updates).Error
}
