package service

import (
	"errors"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type DownloadService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewDownloadService(db *gorm.DB, log *logger.Logger) *DownloadService {
	return &DownloadService{db: db, log: log}
}

// ---------- Category ----------

// GetCategories returns all active download categories.
func (s *DownloadService) GetCategories() ([]model.DownloadCategory, error) {
	var cats []model.DownloadCategory
	err := s.db.Where("is_active = ?", true).
		Order("sort_order ASC, id ASC").
		Find(&cats).Error
	return cats, err
}

// CreateCategory creates a download category.
func (s *DownloadService) CreateCategory(req CreateDownloadCategoryRequest) (*model.DownloadCategory, error) {
	var parentID uint
	if req.ParentID != nil {
		parentID = *req.ParentID
	}
	cat := &model.DownloadCategory{
		Name:      req.Name,
		ParentID:  parentID,
		SortOrder: req.SortOrder,
		IsActive:  true,
	}
	if err := s.db.Create(cat).Error; err != nil {
		return nil, err
	}
	return cat, nil
}

// UpdateCategory updates a download category.
func (s *DownloadService) UpdateCategory(id uint, req UpdateDownloadCategoryRequest) (*model.DownloadCategory, error) {
	var cat model.DownloadCategory
	if err := s.db.First(&cat, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.ParentID != nil {
		updates["parent_id"] = req.ParentID
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := s.db.Model(&cat).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&cat, id).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

// DeleteCategory soft-deletes a download category.
func (s *DownloadService) DeleteCategory(id uint) error {
	var count int64
	s.db.Model(&model.DownloadFile{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("category has files, remove them first")
	}
	return s.db.Delete(&model.DownloadCategory{}, id).Error
}

// ---------- File ----------

// GetFiles returns paginated published download files with optional category filter.
func (s *DownloadService) GetFiles(page, pageSize int, categoryID uint) ([]model.DownloadFile, int64, error) {
	var files []model.DownloadFile
	var total int64

	query := s.db.Model(&model.DownloadFile{}).Where("is_published = ?", true)
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("sort_order ASC, id DESC").
		Preload("Category").
		Find(&files).Error; err != nil {
		return nil, 0, err
	}
	return files, total, nil
}

// GetFileByID returns a single published download file.
func (s *DownloadService) GetFileByID(id uint) (*model.DownloadFile, error) {
	var file model.DownloadFile
	if err := s.db.Preload("Category").
		Where("id = ? AND is_published = ?", id, true).
		First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// IncrementDownload increments the download count for a file.
func (s *DownloadService) IncrementDownload(id uint) error {
	return s.db.Model(&model.DownloadFile{}).
		Where("id = ?", id).
		UpdateColumn("download_count", gorm.Expr("download_count + 1")).Error
}

// AdminGetFiles returns all files (including unpublished) with pagination.
func (s *DownloadService) AdminGetFiles(page, pageSize int, categoryID uint) ([]model.DownloadFile, int64, error) {
	var files []model.DownloadFile
	var total int64

	query := s.db.Model(&model.DownloadFile{})
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("id DESC").
		Preload("Category").
		Find(&files).Error; err != nil {
		return nil, 0, err
	}
	return files, total, nil
}

// AdminGetFile returns a single download file by ID (admin, includes unpublished).
func (s *DownloadService) AdminGetFile(id uint) (*model.DownloadFile, error) {
	var file model.DownloadFile
	if err := s.db.Preload("Category").First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// CreateFile creates a download file record.
func (s *DownloadService) CreateFile(req CreateDownloadFileRequest) (*model.DownloadFile, error) {
	file := &model.DownloadFile{
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Description: req.Description,
		FilePath:    req.FilePath,
		FileSize:    req.FileSize,
		FileType:    req.FileType,
		IsPublished: req.IsPublished,
		SortOrder:   req.SortOrder,
		AdminID:     req.AdminID,
	}
	if err := s.db.Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

// UpdateFile updates a download file record.
func (s *DownloadService) UpdateFile(id uint, req UpdateDownloadFileRequest) (*model.DownloadFile, error) {
	var file model.DownloadFile
	if err := s.db.First(&file, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.FilePath != nil {
		updates["file_path"] = *req.FilePath
	}
	if req.FileSize != nil {
		updates["file_size"] = *req.FileSize
	}
	if req.FileType != nil {
		updates["file_type"] = *req.FileType
	}
	if req.IsPublished != nil {
		updates["is_published"] = *req.IsPublished
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) > 0 {
		if err := s.db.Model(&file).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.Preload("Category").First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// DeleteFile soft-deletes a download file.
func (s *DownloadService) DeleteFile(id uint) error {
	return s.db.Delete(&model.DownloadFile{}, id).Error
}

// UpdateSort updates the sort order of a download file.
func (s *DownloadService) UpdateSort(id uint, sortOrder int) error {
	return s.db.Model(&model.DownloadFile{}).Where("id = ?", id).Update("sort_order", sortOrder).Error
}

// ---------- Request DTOs ----------

type CreateDownloadCategoryRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	ParentID  *uint  `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
}

type UpdateDownloadCategoryRequest struct {
	Name      *string `json:"name"`
	ParentID  *uint   `json:"parent_id"`
	SortOrder *int    `json:"sort_order"`
	IsActive  *bool   `json:"is_active"`
}

type CreateDownloadFileRequest struct {
	CategoryID  uint   `json:"category_id" binding:"required"`
	Title       string `json:"title" binding:"required,max=255"`
	Description string `json:"description"`
	FilePath    string `json:"file_path" binding:"required,max=500"`
	FileSize    int64  `json:"file_size" binding:"required"`
	FileType    string `json:"file_type" binding:"max=50"`
	IsPublished bool   `json:"is_published"`
	SortOrder   int    `json:"sort_order"`
	AdminID     uint   `json:"admin_id"`
}

type UpdateDownloadFileRequest struct {
	CategoryID  *uint   `json:"category_id"`
	Title       *string `json:"title"`
	Description *string `json:"description"`
	FilePath    *string `json:"file_path"`
	FileSize    *int64  `json:"file_size"`
	FileType    *string `json:"file_type"`
	IsPublished *bool   `json:"is_published"`
	SortOrder   *int    `json:"sort_order"`
}
