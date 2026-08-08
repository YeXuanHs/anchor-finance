package repository

import (
	"github.com/anchorfinance/backend/internal/model"
	"github.com/anchorfinance/backend/pkg/database"
	"gorm.io/gorm"
)

type PluginRepository struct {
	db *gorm.DB
}

func NewPluginRepository() *PluginRepository {
	return &PluginRepository{
		db: database.GetDB(),
	}
}

// Create 创建插件
func (r *PluginRepository) Create(plugin *model.Plugin) error {
	return r.db.Create(plugin).Error
}

// FindByID 根据ID查找插件
func (r *PluginRepository) FindByID(id uint) (*model.Plugin, error) {
	var plugin model.Plugin
	err := r.db.First(&plugin, id).Error
	if err != nil {
		return nil, err
	}
	return &plugin, nil
}

// FindByCode 根据编码查找插件
func (r *PluginRepository) FindByCode(code string) (*model.Plugin, error) {
	var plugin model.Plugin
	err := r.db.Where("code = ?", code).First(&plugin).Error
	if err != nil {
		return nil, err
	}
	return &plugin, nil
}

// Update 更新插件
func (r *PluginRepository) Update(plugin *model.Plugin) error {
	return r.db.Save(plugin).Error
}

// Delete 删除插件
func (r *PluginRepository) Delete(id uint) error {
	return r.db.Delete(&model.Plugin{}, id).Error
}

// List 获取插件列表
func (r *PluginRepository) List(params *model.PluginQueryParams) ([]model.Plugin, int64, error) {
	var plugins []model.Plugin
	var total int64

	query := r.db.Model(&model.Plugin{})

	// 分类筛选
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}

	// 关键词搜索
	if params.Keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("category, name").Find(&plugins).Error; err != nil {
		return nil, 0, err
	}

	return plugins, total, nil
}

// FindEnabled 获取已启用的插件
func (r *PluginRepository) FindEnabled() ([]model.Plugin, error) {
	var plugins []model.Plugin
	err := r.db.Where("enabled = ?", true).Find(&plugins).Error
	return plugins, err
}

// FindByCategory 根据分类获取插件
func (r *PluginRepository) FindByCategory(category string) ([]model.Plugin, error) {
	var plugins []model.Plugin
	err := r.db.Where("category = ? AND enabled = ?", category, true).Find(&plugins).Error
	return plugins, err
}
