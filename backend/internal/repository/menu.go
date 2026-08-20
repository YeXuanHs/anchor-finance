package repository

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/database"
	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{
		db: database.GetDB(),
	}
}

// Create 创建菜单
func (r *MenuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

// FindByID 根据ID查找菜单
func (r *MenuRepository) FindByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// Update 更新菜单
func (r *MenuRepository) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

// Delete 删除菜单
func (r *MenuRepository) Delete(id uint) error {
	return r.db.Delete(&model.Menu{}, id).Error
}

// FindAll 获取所有菜单
func (r *MenuRepository) FindAll() ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Order("sort_order, id").Find(&menus).Error
	return menus, err
}

// FindActive 获取所有激活的菜单（包括 is_visible=0 的隐藏菜单，用于路由注册）
func (r *MenuRepository) FindActive() ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Where("is_active = ?", true).
		Order("sort_order, id").
		Find(&menus).Error
	return menus, err
}

// FindByParentID 根据父级ID查找子菜单
func (r *MenuRepository) FindByParentID(parentID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Where("parent_id = ?", parentID).Order("sort_order, id").Find(&menus).Error
	return menus, err
}
