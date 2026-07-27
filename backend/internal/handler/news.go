package handler

import (
	"strconv"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NewsHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewNewsHandler(db *gorm.DB, log *logger.Logger) *NewsHandler {
	return &NewsHandler{db: db, log: log}
}

// ---------- Public Endpoints ----------

// GetCategories returns all active news categories.
func (h *NewsHandler) GetCategories(c *gin.Context) {
	var categories []model.NewsCategory
	if err := h.db.Where("is_active = ?", true).Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, categories)
}

// GetList returns paginated published news with optional category and keyword filter.
func (h *NewsHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var categoryID uint
	if cid := c.Query("category_id"); cid != "" {
		v, _ := strconv.ParseUint(cid, 10, 64)
		categoryID = uint(v)
	}

	var news []model.News
	var total int64

	query := h.db.Model(&model.News{}).Where("is_published = ?", true)
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR keywords LIKE ? OR summary LIKE ?", q, q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("is_sticky DESC, published_at DESC, id DESC").
		Preload("Category").
		Find(&news).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, news, total, page, pageSize)
}

// GetDetail returns a single published news by slug and increments view count.
func (h *NewsHandler) GetDetail(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "slug is required")
		return
	}

	var news model.News
	if err := h.db.Preload("Category").
		Where("slug = ? AND is_published = ?", slug, true).
		First(&news).Error; err != nil {
		response.NotFound(c, "news not found")
		return
	}

	h.db.Model(&news).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	response.Success(c, news)
}

// ---------- Admin Endpoints ----------

// AdminGetCategories returns all news categories including inactive (admin).
func (h *NewsHandler) AdminGetCategories(c *gin.Context) {
	var categories []model.NewsCategory
	if err := h.db.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, categories)
}

type CreateNewsCategoryRequest struct {
	Name      string `json:"name" binding:"required,max=50"`
	Slug      string `json:"slug" binding:"required,max=50"`
	SortOrder int    `json:"sort_order"`
}

// AdminCreateCategory creates a news category (admin).
func (h *NewsHandler) AdminCreateCategory(c *gin.Context) {
	var req CreateNewsCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cat := model.NewsCategory{
		Name:      req.Name,
		Slug:      req.Slug,
		SortOrder: req.SortOrder,
		IsActive:  true,
	}
	if err := h.db.Create(&cat).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cat)
}

type UpdateNewsCategoryRequest struct {
	Name      *string `json:"name"`
	Slug      *string `json:"slug"`
	SortOrder *int    `json:"sort_order"`
	IsActive  *bool   `json:"is_active"`
}

// AdminUpdateCategory updates a news category (admin).
func (h *NewsHandler) AdminUpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	var req UpdateNewsCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := h.db.Model(&model.NewsCategory{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}
	response.SuccessMsg(c, "category updated")
}

// AdminDeleteCategory deletes a news category (admin).
func (h *NewsHandler) AdminDeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	var count int64
	h.db.Model(&model.News{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		response.BadRequest(c, "category has news articles, remove them first")
		return
	}

	if err := h.db.Delete(&model.NewsCategory{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "category deleted")
}

// AdminGetList returns all news including unpublished (admin).
func (h *NewsHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var categoryID uint
	if cid := c.Query("category_id"); cid != "" {
		v, _ := strconv.ParseUint(cid, 10, 64)
		categoryID = uint(v)
	}

	var news []model.News
	var total int64

	query := h.db.Model(&model.News{})
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR keywords LIKE ?", q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("id DESC").
		Preload("Category").
		Find(&news).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, news, total, page, pageSize)
}

// AdminGetDetail returns a single news by ID (admin, includes unpublished).
func (h *NewsHandler) AdminGetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid news id")
		return
	}

	var news model.News
	if err := h.db.Preload("Category").First(&news, id).Error; err != nil {
		response.NotFound(c, "news not found")
		return
	}
	response.Success(c, news)
}

type CreateNewsRequest struct {
	CategoryID  uint   `json:"category_id" binding:"required"`
	Title       string `json:"title" binding:"required,max=255"`
	Slug        string `json:"slug" binding:"required,max=255"`
	Summary     string `json:"summary" binding:"max=500"`
	Content     string `json:"content"`
	CoverImage  string `json:"cover_image" binding:"max=255"`
	Keywords    string `json:"keywords" binding:"max=255"`
	IsPublished bool   `json:"is_published"`
	IsSticky    bool   `json:"is_sticky"`
	AdminID     uint   `json:"admin_id"`
}

// AdminCreate creates a news article (admin).
func (h *NewsHandler) AdminCreate(c *gin.Context) {
	var req CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	news := model.News{
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Slug:        req.Slug,
		Summary:     req.Summary,
		Content:     req.Content,
		CoverImage:  req.CoverImage,
		Keywords:    req.Keywords,
		IsPublished: req.IsPublished,
		IsSticky:    req.IsSticky,
		AdminID:     req.AdminID,
	}

	if req.IsPublished {
		now := time.Now()
		news.PublishedAt = &now
	}

	if err := h.db.Create(&news).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, news)
}

type UpdateNewsRequest struct {
	CategoryID  *uint   `json:"category_id"`
	Title       *string `json:"title"`
	Slug        *string `json:"slug"`
	Summary     *string `json:"summary"`
	Content     *string `json:"content"`
	CoverImage  *string `json:"cover_image"`
	Keywords    *string `json:"keywords"`
	IsPublished *bool   `json:"is_published"`
	IsSticky    *bool   `json:"is_sticky"`
}

// AdminUpdate updates a news article (admin).
func (h *NewsHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid news id")
		return
	}

	var req UpdateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
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
	if req.Summary != nil {
		updates["summary"] = *req.Summary
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.CoverImage != nil {
		updates["cover_image"] = *req.CoverImage
	}
	if req.Keywords != nil {
		updates["keywords"] = *req.Keywords
	}
	if req.IsPublished != nil {
		updates["is_published"] = *req.IsPublished
	}
	if req.IsSticky != nil {
		updates["is_sticky"] = *req.IsSticky
	}

	if len(updates) > 0 {
		if err := h.db.Model(&model.News{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}
	response.SuccessMsg(c, "news updated")
}

// AdminDelete soft-deletes a news article (admin).
func (h *NewsHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid news id")
		return
	}

	if err := h.db.Delete(&model.News{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "news deleted")
}
