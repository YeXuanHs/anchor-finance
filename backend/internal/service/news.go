package service

import (
	"errors"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// NewsService 新闻管理服务
type NewsService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewNewsService(db *gorm.DB, log *logger.Logger) *NewsService {
	return &NewsService{db: db, log: log}
}

// CreateNewsRequest 创建新闻请求
type CreateNewsRequest struct {
	CategoryID uint    `json:"category_id" binding:"required"`
	Title      string  `json:"title" binding:"required,max=255"`
	Slug       string  `json:"slug" binding:"required,max=255"`
	Summary    string  `json:"summary" binding:"max=500"`
	Content    string  `json:"content"`
	CoverImage string  `json:"cover_image"`
	Keywords   string  `json:"keywords"`
	IsSticky   bool    `json:"is_sticky"`
	AdminID    uint    `json:"admin_id"`
}

// Create 创建新闻
func (s *NewsService) Create(req CreateNewsRequest) (*model.News, error) {
	var count int64
	s.db.Model(&model.News{}).Where("slug = ?", req.Slug).Count(&count)
	if count > 0 {
		return nil, errors.New("slug already exists")
	}

	news := model.News{
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Slug:        req.Slug,
		Summary:     req.Summary,
		Content:     req.Content,
		CoverImage:  req.CoverImage,
		Keywords:    req.Keywords,
		IsSticky:    req.IsSticky,
		IsPublished: false,
		AdminID:     req.AdminID,
	}

	if err := s.db.Create(&news).Error; err != nil {
		return nil, err
	}
	s.log.Infof("news created: id=%d title=%s", news.ID, news.Title)
	return &news, nil
}

// Get 获取新闻
func (s *NewsService) Get(id uint) (*model.News, error) {
	var news model.News
	if err := s.db.Preload("Category").First(&news, id).Error; err != nil {
		return nil, err
	}
	return &news, nil
}

// GetBySlug 根据slug获取新闻
func (s *NewsService) GetBySlug(slug string) (*model.News, error) {
	var news model.News
	if err := s.db.Preload("Category").Where("slug = ?", slug).First(&news).Error; err != nil {
		return nil, err
	}
	return &news, nil
}

// List 获取新闻列表
func (s *NewsService) List(page, pageSize int, publishedOnly bool) ([]model.News, int64, error) {
	var items []model.News
	var total int64

	query := s.db.Model(&model.News{})
	if publishedOnly {
		query = query.Where("is_published = ?", true)
	}

	query.Count(&total)
	offset, limit := util.Paginate(page, pageSize)
	query.Offset(offset).Limit(limit).Order("is_sticky DESC, id DESC").Preload("Category").Find(&items)
	return items, total, nil
}

// Update 更新新闻
func (s *NewsService) Update(id uint, updates map[string]interface{}) (*model.News, error) {
	var news model.News
	if err := s.db.First(&news, id).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&news).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.db.Preload("Category").First(&news, id)
	return &news, nil
}

// Delete 删除新闻
func (s *NewsService) Delete(id uint) error {
	result := s.db.Delete(&model.News{}, id)
	if result.RowsAffected == 0 {
		return errors.New("news not found")
	}
	return result.Error
}

// Publish 发布新闻
func (s *NewsService) Publish(id uint) (*model.News, error) {
	var news model.News
	if err := s.db.First(&news, id).Error; err != nil {
		return nil, err
	}
	now := time.Now()
	updates := map[string]interface{}{
		"is_published": true,
		"published_at": &now,
	}
	if err := s.db.Model(&news).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.db.Preload("Category").First(&news, id)
	s.log.Infof("news published: id=%d", id)
	return &news, nil
}

// GetByCategory 根据分类获取新闻
func (s *NewsService) GetByCategory(categoryID uint, page, pageSize int) ([]model.News, int64, error) {
	var items []model.News
	var total int64

	query := s.db.Model(&model.News{}).Where("category_id = ? AND is_published = ?", categoryID, true)
	query.Count(&total)
	offset, limit := util.Paginate(page, pageSize)
	query.Offset(offset).Limit(limit).Order("is_sticky DESC, id DESC").Preload("Category").Find(&items)
	return items, total, nil
}
