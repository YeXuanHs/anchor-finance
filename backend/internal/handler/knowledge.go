package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KnowledgeHandler struct {
	db           *gorm.DB
	knowledgeSvc *service.KnowledgeService
	log          *logger.Logger
	uploadDir    string
}

func NewKnowledgeHandler(db *gorm.DB, knowledgeSvc *service.KnowledgeService, log *logger.Logger, uploadDir ...string) *KnowledgeHandler {
	h := &KnowledgeHandler{db: db, knowledgeSvc: knowledgeSvc, log: log}
	if len(uploadDir) > 0 && uploadDir[0] != "" {
		h.uploadDir = uploadDir[0]
	} else {
		h.uploadDir = "uploads/knowledge"
	}
	return h
}

// ---------- Public ----------

// GetCategories returns all active knowledge categories.
func (h *KnowledgeHandler) GetCategories(c *gin.Context) {
	cats, err := h.knowledgeSvc.GetCategories()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cats)
}

// GetArticles returns paginated published articles.
func (h *KnowledgeHandler) GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var categoryID uint
	if cid := c.Query("category_id"); cid != "" {
		v, _ := strconv.ParseUint(cid, 10, 64)
		categoryID = uint(v)
	}

	articles, total, err := h.knowledgeSvc.GetArticles(page, pageSize, categoryID, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, articles, total, page, pageSize)
}

// GetArticle returns a single article by slug.
func (h *KnowledgeHandler) GetArticle(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "slug is required")
		return
	}

	article, err := h.knowledgeSvc.GetArticleBySlug(slug)
	if err != nil {
		response.NotFound(c, "article not found")
		return
	}
	response.Success(c, article)
}

// GetHot returns the most viewed knowledge articles.
func (h *KnowledgeHandler) GetHot(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 50 {
		limit = 10
	}
	articles, err := h.knowledgeSvc.GetHotArticles(limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, articles)
}

// Search searches knowledge articles by keyword.
func (h *KnowledgeHandler) Search(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		keyword = c.Query("keyword")
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	articles, total, err := h.knowledgeSvc.SearchArticles(keyword, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, articles, total, page, pageSize)
}

// GetRelatedArticles returns related articles (same category, excluding self, limit 5).
func (h *KnowledgeHandler) GetRelatedArticles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid article id")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if limit < 1 || limit > 20 {
		limit = 5
	}

	articles, err := h.knowledgeSvc.GetRelatedArticles(uint(id), limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, articles)
}

// SubmitFeedback records user feedback for an article (helpful or not helpful).
func (h *KnowledgeHandler) SubmitFeedback(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid article id")
		return
	}

	var req struct {
		Helpful bool `json:"helpful"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.knowledgeSvc.SubmitFeedback(uint(id), req.Helpful); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "feedback submitted")
}

// GetSubCategories returns sub-categories of the given category.
func (h *KnowledgeHandler) GetSubCategories(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	cats, err := h.knowledgeSvc.GetSubCategories(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cats)
}

// MarkHelpful marks an article as helpful or not helpful.
func (h *KnowledgeHandler) MarkHelpful(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid article id")
		return
	}

	var req struct {
		Helpful bool `json:"helpful"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.knowledgeSvc.MarkHelpful(uint(id), req.Helpful); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "feedback recorded")
}

// ---------- Admin ----------

// AdminGetCategories returns all categories including inactive (admin).
func (h *KnowledgeHandler) AdminGetCategories(c *gin.Context) {
	cats, err := h.knowledgeSvc.GetCategories()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cats)
}

// AdminCreateCategory creates a category (admin).
func (h *KnowledgeHandler) AdminCreateCategory(c *gin.Context) {
	var req service.CreateKnowledgeCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cat, err := h.knowledgeSvc.CreateCategory(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cat)
}

// AdminUpdateCategory updates a category (admin).
func (h *KnowledgeHandler) AdminUpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	var req service.UpdateKnowledgeCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cat, err := h.knowledgeSvc.UpdateCategory(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cat)
}

// AdminDeleteCategory deletes a category (admin).
func (h *KnowledgeHandler) AdminDeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	if err := h.knowledgeSvc.DeleteCategory(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "category deleted")
}

// AdminGetArticles returns all articles including unpublished (admin).
func (h *KnowledgeHandler) AdminGetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var categoryID uint
	if cid := c.Query("category_id"); cid != "" {
		v, _ := strconv.ParseUint(cid, 10, 64)
		categoryID = uint(v)
	}

	articles, total, err := h.knowledgeSvc.AdminGetArticles(page, pageSize, categoryID, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, articles, total, page, pageSize)
}

// AdminGetArticle returns a single article by ID (admin).
func (h *KnowledgeHandler) AdminGetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid article id")
		return
	}

	article, err := h.knowledgeSvc.AdminGetArticle(uint(id))
	if err != nil {
		response.NotFound(c, "article not found")
		return
	}
	response.Success(c, article)
}

// AdminCreateArticle creates an article (admin).
func (h *KnowledgeHandler) AdminCreateArticle(c *gin.Context) {
	var req service.CreateKnowledgeArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	article, err := h.knowledgeSvc.CreateArticle(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, article)
}

// AdminUpdateArticle updates an article (admin).
func (h *KnowledgeHandler) AdminUpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid article id")
		return
	}

	var req service.UpdateKnowledgeArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	article, err := h.knowledgeSvc.UpdateArticle(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, article)
}

// AdminDeleteArticle deletes an article (admin).
func (h *KnowledgeHandler) AdminDeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid article id")
		return
	}

	if err := h.knowledgeSvc.DeleteArticle(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "article deleted")
}

// ==================== 新增缺失方法 ====================

// TagsList returns articles matching a tag/keyword with pagination (admin).
// GET /admin/knowledge/tags
func (h *KnowledgeHandler) TagsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	tag := c.Query("tag")
	if tag == "" {
		tag = c.Query("keyword")
	}

	var articles []model.KnowledgeArticle
	var total int64

	query := h.db.Model(&model.KnowledgeArticle{})
	if tag != "" {
		q := "%" + tag + "%"
		query = query.Where("keywords LIKE ? OR title LIKE ?", q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("id DESC").
		Preload("Category").
		Find(&articles).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// Collect all unique tags from articles
	tagSet := make(map[string]bool)
	for _, article := range articles {
		for _, kw := range strings.Split(article.Keywords, ",") {
			kw = strings.TrimSpace(kw)
			if kw != "" {
				tagSet[kw] = true
			}
		}
	}
	tags := make([]string, 0, len(tagSet))
	for t := range tagSet {
		tags = append(tags, t)
	}

	response.Success(c, gin.H{
		"articles": articles,
		"total":    total,
		"tags":     tags,
	})
}

// UploadHandle handles file upload for knowledge articles (admin).
// POST /admin/knowledge/upload
func (h *KnowledgeHandler) UploadHandle(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		response.BadRequest(c, "file size exceeds 10MB limit")
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".zip": true, ".rar": true, ".txt": true, ".md": true,
	}
	if !allowedExts[ext] {
		response.BadRequest(c, "unsupported file type: "+ext)
		return
	}

	// Create upload directory
	uploadDir := filepath.Join(h.uploadDir, time.Now().Format("200601"))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.ServerError(c, "failed to create upload directory")
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		response.ServerError(c, "failed to save file")
		return
	}

	// Return the relative URL
	urlPath := strings.ReplaceAll(savePath, "\\", "/")
	response.Success(c, gin.H{
		"url":      "/" + urlPath,
		"filename": file.Filename,
		"size":     file.Size,
	})
}
