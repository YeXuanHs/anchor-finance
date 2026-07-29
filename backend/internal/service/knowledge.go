package service

import (
	"errors"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type KnowledgeService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewKnowledgeService(db *gorm.DB, log *logger.Logger) *KnowledgeService {
	return &KnowledgeService{db: db, log: log}
}

// ---------- Category ----------

// GetCategories returns all active categories as a flat list (caller can build tree).
func (s *KnowledgeService) GetCategories() ([]model.KnowledgeCategory, error) {
	var cats []model.KnowledgeCategory
	err := s.db.Where("is_active = ?", true).
		Order("sort_order ASC, id ASC").
		Preload("Children", "is_active = ?", true).
		Find(&cats).Error
	return cats, err
}

// CreateCategory creates a knowledge category.
func (s *KnowledgeService) CreateCategory(req CreateKnowledgeCategoryRequest) (*model.KnowledgeCategory, error) {
	cat := &model.KnowledgeCategory{
		Name:      req.Name,
		Slug:      req.Slug,
		ParentID:  req.ParentID,
		SortOrder: req.SortOrder,
		IsActive:  true,
	}
	if err := s.db.Create(cat).Error; err != nil {
		return nil, err
	}
	return cat, nil
}

// UpdateCategory updates a knowledge category.
func (s *KnowledgeService) UpdateCategory(id uint, req UpdateKnowledgeCategoryRequest) (*model.KnowledgeCategory, error) {
	var cat model.KnowledgeCategory
	if err := s.db.First(&cat, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
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

// DeleteCategory soft-deletes a knowledge category.
func (s *KnowledgeService) DeleteCategory(id uint) error {
	var count int64
	s.db.Model(&model.KnowledgeArticle{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return errors.New("category has articles, remove them first")
	}
	return s.db.Delete(&model.KnowledgeCategory{}, id).Error
}

// ---------- Article ----------

// GetArticles returns paginated published articles with optional search and category filter.
func (s *KnowledgeService) GetArticles(page, pageSize int, categoryID uint, keyword string) ([]model.KnowledgeArticle, int64, error) {
	var articles []model.KnowledgeArticle
	var total int64

	query := s.db.Model(&model.KnowledgeArticle{}).Where("is_published = ?", true)
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR keywords LIKE ? OR summary LIKE ?", q, q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("sort_order ASC, id DESC").
		Preload("Category").
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

// GetArticleBySlug returns a single article by slug and increments view count.
func (s *KnowledgeService) GetArticleBySlug(slug string) (*model.KnowledgeArticle, error) {
	var article model.KnowledgeArticle
	if err := s.db.Preload("Category").
		Where("slug = ? AND is_published = ?", slug, true).
		First(&article).Error; err != nil {
		return nil, err
	}
	s.db.Model(&article).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	return &article, nil
}

// MarkHelpful records whether an article was helpful.
func (s *KnowledgeService) MarkHelpful(articleID uint, helpful bool) error {
	var article model.KnowledgeArticle
	if err := s.db.First(&article, articleID).Error; err != nil {
		return err
	}
	if helpful {
		return s.db.Model(&article).UpdateColumn("help_count", gorm.Expr("help_count + 1")).Error
	}
	return s.db.Model(&article).UpdateColumn("not_help_count", gorm.Expr("not_help_count + 1")).Error
}

// GetHotArticles returns the most viewed published articles.
func (s *KnowledgeService) GetHotArticles(limit int) ([]model.KnowledgeArticle, error) {
	var articles []model.KnowledgeArticle
	if err := s.db.Where("is_published = ?", true).
		Order("view_count DESC").
		Limit(limit).
		Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// SearchArticles searches published articles by keyword with pagination.
func (s *KnowledgeService) SearchArticles(keyword string, page, pageSize int) ([]model.KnowledgeArticle, int64, error) {
	var articles []model.KnowledgeArticle
	var total int64

	query := s.db.Model(&model.KnowledgeArticle{}).Where("is_published = ?", true)
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR keywords LIKE ? OR summary LIKE ? OR content LIKE ?", q, q, q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("view_count DESC, id DESC").
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

// AdminGetArticles returns all articles (including unpublished) with pagination.
func (s *KnowledgeService) AdminGetArticles(page, pageSize int, categoryID uint, keyword string) ([]model.KnowledgeArticle, int64, error) {
	var articles []model.KnowledgeArticle
	var total int64

	query := s.db.Model(&model.KnowledgeArticle{})
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR keywords LIKE ?", q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("id DESC").
		Preload("Category").
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

// AdminGetArticle returns a single article by ID (admin, includes unpublished).
func (s *KnowledgeService) AdminGetArticle(id uint) (*model.KnowledgeArticle, error) {
	var article model.KnowledgeArticle
	if err := s.db.Preload("Category").First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// CreateArticle creates a knowledge article.
func (s *KnowledgeService) CreateArticle(req CreateKnowledgeArticleRequest) (*model.KnowledgeArticle, error) {
	article := &model.KnowledgeArticle{
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Slug:        req.Slug,
		Content:     req.Content,
		Summary:     req.Summary,
		Keywords:    req.Keywords,
		IsPublished: req.IsPublished,
		SortOrder:   req.SortOrder,
		AdminID:     req.AdminID,
	}
	if err := s.db.Create(article).Error; err != nil {
		return nil, err
	}
	return article, nil
}

// UpdateArticle updates a knowledge article.
func (s *KnowledgeService) UpdateArticle(id uint, req UpdateKnowledgeArticleRequest) (*model.KnowledgeArticle, error) {
	var article model.KnowledgeArticle
	if err := s.db.First(&article, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Summary != nil {
		updates["summary"] = *req.Summary
	}
	if req.Keywords != nil {
		updates["keywords"] = *req.Keywords
	}
	if req.IsPublished != nil {
		updates["is_published"] = *req.IsPublished
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) > 0 {
		if err := s.db.Model(&article).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.Preload("Category").First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

// DeleteArticle soft-deletes a knowledge article.
func (s *KnowledgeService) DeleteArticle(id uint) error {
	return s.db.Delete(&model.KnowledgeArticle{}, id).Error
}

// ---------- Request DTOs ----------

type CreateKnowledgeCategoryRequest struct {
	Name      string `json:"name" binding:"required,max=100"`
	Slug      string `json:"slug" binding:"required,max=100"`
	ParentID  *uint  `json:"parent_id"`
	SortOrder int    `json:"sort_order"`
}

type UpdateKnowledgeCategoryRequest struct {
	Name      *string `json:"name"`
	Slug      *string `json:"slug"`
	ParentID  *uint   `json:"parent_id"`
	SortOrder *int    `json:"sort_order"`
	IsActive  *bool   `json:"is_active"`
}

type CreateKnowledgeArticleRequest struct {
	CategoryID  uint   `json:"category_id" binding:"required"`
	Title       string `json:"title" binding:"required,max=255"`
	Slug        string `json:"slug" binding:"required,max=255"`
	Content     string `json:"content"`
	Summary     string `json:"summary" binding:"max=500"`
	Keywords    string `json:"keywords" binding:"max=255"`
	IsPublished bool   `json:"is_published"`
	SortOrder   int    `json:"sort_order"`
	AdminID     uint   `json:"admin_id"`
}

type UpdateKnowledgeArticleRequest struct {
	CategoryID  *uint   `json:"category_id"`
	Title       *string `json:"title"`
	Slug        *string `json:"slug"`
	Content     *string `json:"content"`
	Summary     *string `json:"summary"`
	Keywords    *string `json:"keywords"`
	IsPublished *bool   `json:"is_published"`
	SortOrder   *int    `json:"sort_order"`
}
