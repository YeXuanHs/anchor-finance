package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type RbacPageService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewRbacPageService(db *gorm.DB, log *logger.Logger) *RbacPageService {
	return &RbacPageService{db: db, log: log}
}

// GetList returns paginated RBAC pages.
func (s *RbacPageService) GetList(page, pageSize int, module string, keyword string) ([]model.RbacPage, int64, error) {
	var pages []model.RbacPage
	var total int64

	query := s.db.Model(&model.RbacPage{})
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR path LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("sort_order ASC, id ASC").Find(&pages).Error; err != nil {
		return nil, 0, err
	}
	return pages, total, nil
}

// GetByID returns a single page by ID.
func (s *RbacPageService) GetByID(id uint) (*model.RbacPage, error) {
	var page model.RbacPage
	if err := s.db.First(&page, id).Error; err != nil {
		return nil, err
	}
	return &page, nil
}

// GetTree returns RBAC pages as a tree structure.
func (s *RbacPageService) GetTree(module string) ([]model.RbacPage, error) {
	var pages []model.RbacPage
	query := s.db.Where("parent_id IS NULL AND status = 1")
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if err := query.Preload("Children", "status = 1").Order("sort_order ASC").Find(&pages).Error; err != nil {
		return nil, err
	}
	return pages, nil
}

// Create creates a new RBAC page.
func (s *RbacPageService) Create(page *model.RbacPage) error {
	return s.db.Create(page).Error
}

// Update updates an RBAC page.
func (s *RbacPageService) Update(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.RbacPage{}).Where("id = ?", id).Updates(updates).Error
}

// Delete deletes an RBAC page.
func (s *RbacPageService) Delete(id uint) error {
	return s.db.Delete(&model.RbacPage{}, id).Error
}

// GetAuthTree returns the auth rule tree for pages.
func (s *RbacPageService) GetAuthTree() ([]model.RbacPage, error) {
	var pages []model.RbacPage
	if err := s.db.Where("parent_id IS NULL AND status = 1").
		Preload("Children", "status = 1").
		Order("sort_order ASC").Find(&pages).Error; err != nil {
		return nil, err
	}
	return pages, nil
}

// SetPageAuth sets auth rules for a page and role.
func (s *RbacPageService) SetPageAuth(pageID, roleID uint, canView, canEdit, canDelete bool) error {
	auth := model.RbacPageAuth{
		PageID:    pageID,
		RoleID:    roleID,
		CanView:   canView,
		CanEdit:   canEdit,
		CanDelete: canDelete,
	}
	return s.db.Where("page_id = ? AND role_id = ?", pageID, roleID).
		Assign(auth).FirstOrCreate(&auth).Error
}
