package repository

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/database"
	"gorm.io/gorm"
)

type SupplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository() *SupplierRepository {
	return &SupplierRepository{
		db: database.GetDB(),
	}
}

// Create 创建供应商
func (r *SupplierRepository) Create(supplier *model.Supplier) error {
	return r.db.Create(supplier).Error
}

// FindByID 根据ID查找供应商
func (r *SupplierRepository) FindByID(id uint) (*model.Supplier, error) {
	var supplier model.Supplier
	err := r.db.First(&supplier, id).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

// Update 更新供应商
func (r *SupplierRepository) Update(supplier *model.Supplier) error {
	return r.db.Save(supplier).Error
}

// Delete 删除供应商（软删除）
func (r *SupplierRepository) Delete(id uint) error {
	return r.db.Delete(&model.Supplier{}, id).Error
}

// List 获取供应商列表
func (r *SupplierRepository) List(params *model.SupplierQueryParams) ([]model.Supplier, int64, error) {
	var suppliers []model.Supplier
	var total int64

	query := r.db.Model(&model.Supplier{})

	// 关键词搜索
	if params.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+params.Keyword+"%")
	}

	// API类型筛选
	if params.APIType != "" {
		query = query.Where("api_type = ?", params.APIType)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("id DESC").Find(&suppliers).Error; err != nil {
		return nil, 0, err
	}

	return suppliers, total, nil
}

// FindByAPIType 根据API类型查找供应商
func (r *SupplierRepository) FindByAPIType(apiType string) ([]model.Supplier, error) {
	var suppliers []model.Supplier
	err := r.db.Where("api_type = ? AND status = ?", apiType, "active").Find(&suppliers).Error
	return suppliers, err
}
