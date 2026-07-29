package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type RuleMiddleService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewRuleMiddleService(db *gorm.DB, log *logger.Logger) *RuleMiddleService {
	return &RuleMiddleService{db: db, log: log}
}

// GetMenuList returns all rule middle menus.
func (s *RuleMiddleService) GetMenuList(page, pageSize int, keyword string) ([]model.RuleMiddle, int64, error) {
	var menus []model.RuleMiddle
	var total int64

	query := s.db.Model(&model.RuleMiddle{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("sort_order ASC, id ASC").Find(&menus).Error; err != nil {
		return nil, 0, err
	}
	return menus, total, nil
}

// GetByID returns a single menu by ID.
func (s *RuleMiddleService) GetByID(id uint) (*model.RuleMiddle, error) {
	var menu model.RuleMiddle
	if err := s.db.First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}

// Create creates a new rule middle menu.
func (s *RuleMiddleService) Create(menu *model.RuleMiddle) error {
	return s.db.Create(menu).Error
}

// Update updates a rule middle menu.
func (s *RuleMiddleService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.RuleMiddle{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes a rule middle menu.
func (s *RuleMiddleService) Delete(id uint) error {
	return s.db.Delete(&model.RuleMiddle{}, id).Error
}
