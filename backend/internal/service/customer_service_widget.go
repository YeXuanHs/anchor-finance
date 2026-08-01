package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type CustomerServiceWidgetService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCustomerServiceWidgetService(db *gorm.DB, log *logger.Logger) *CustomerServiceWidgetService {
	return &CustomerServiceWidgetService{db: db, log: log}
}

// List 获取客服列表
func (s *CustomerServiceWidgetService) List(showInactive bool) ([]model.CustomerServiceWidget, error) {
	var items []model.CustomerServiceWidget
	q := s.db.Order("sort_order ASC, id ASC")
	if !showInactive {
		q = q.Where("is_active = ?", true)
	}
	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetByID 获取单个客服
func (s *CustomerServiceWidgetService) GetByID(id uint) (*model.CustomerServiceWidget, error) {
	var item model.CustomerServiceWidget
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Create 创建客服入口
func (s *CustomerServiceWidgetService) Create(item *model.CustomerServiceWidget) error {
	return s.db.Create(item).Error
}

// Update 更新客服入口
func (s *CustomerServiceWidgetService) Update(id uint, updates map[string]interface{}) error {
	updates["updated_at"] = gorm.Expr("NOW()")
	return s.db.Model(&model.CustomerServiceWidget{}).Where("id = ?", id).Updates(updates).Error
}

// Delete 删除客服入口
func (s *CustomerServiceWidgetService) Delete(id uint) error {
	return s.db.Delete(&model.CustomerServiceWidget{}, id).Error
}

// GetSettings 获取全局设置
func (s *CustomerServiceWidgetService) GetSettings() *model.CustomerServiceWidgetSetting {
	var settings model.CustomerServiceWidgetSetting
	if err := s.db.First(&settings).Error; err != nil {
		// 返回默认值
		return &model.CustomerServiceWidgetSetting{
			Enabled:      true,
			Position:     "right",
			OffsetBottom: 100,
			ThemeColor:   "#1890ff",
			ShowOnMobile: true,
			Title:        "联系客服",
			WelcomeText:  "您好，请问有什么可以帮助您？",
		}
	}
	return &settings
}

// UpdateSettings 更新全局设置
func (s *CustomerServiceWidgetService) UpdateSettings(updates map[string]interface{}) error {
	var settings model.CustomerServiceWidgetSetting
	if err := s.db.First(&settings).Error; err != nil {
		// 不存在则创建
		settings = model.CustomerServiceWidgetSetting{
			Enabled:      true,
			Position:     "right",
			OffsetBottom: 100,
			ThemeColor:   "#1890ff",
			ShowOnMobile: true,
			Title:        "联系客服",
			WelcomeText:  "您好，请问有什么可以帮助您？",
		}
		return s.db.Create(&settings).Error
	}
	updates["updated_at"] = gorm.Expr("NOW()")
	return s.db.Model(&settings).Updates(updates).Error
}

// GetActiveForDisplay 前台获取活跃的客服列表
func (s *CustomerServiceWidgetService) GetActiveForDisplay() ([]model.CustomerServiceWidget, *model.CustomerServiceWidgetSetting) {
	var items []model.CustomerServiceWidget
	s.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&items)
	settings := s.GetSettings()
	return items, settings
}
