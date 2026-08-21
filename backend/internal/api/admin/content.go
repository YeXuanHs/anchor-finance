package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// ========== 新闻管理 ==========

// GetNewsList 获取新闻列表
// GET /api/admin/news
func GetNewsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	categoryID := c.Query("category_id")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.News{})

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var news []model.News
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&news)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      news,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetNewsCategories 获取新闻分类
// GET /api/admin/news-categories
func GetNewsCategories(c *gin.Context) {
	db := database.GetDB()
	var categories []model.NewsCategory
	db.Where("status = ?", "active").Order("sort_order ASC").Find(&categories)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    categories,
	})
}

// ========== 知识库管理 ==========

// GetKnowledgeCategories 获取知识库分类
// GET /api/admin/knowledge/categories
func GetKnowledgeCategories(c *gin.Context) {
	db := database.GetDB()
	var categories []model.KnowledgeCategory
	db.Where("status = ?", "active").Order("sort_order ASC").Find(&categories)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    categories,
	})
}

// GetKnowledgeArticles 获取知识库文章
// GET /api/admin/knowledge/articles
func GetKnowledgeArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	categoryID := c.Query("category_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.KnowledgeArticle{})

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	var total int64
	query.Count(&total)

	var articles []model.KnowledgeArticle
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&articles)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      articles,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ========== 下载管理 ==========

// GetDownloads 获取下载列表
// GET /api/admin/downloads
func GetDownloads(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	categoryID := c.Query("category_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Download{})

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	var total int64
	query.Count(&total)

	var downloads []model.Download
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&downloads)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      downloads,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetDownloadCategories 获取下载分类
// GET /api/admin/downloads/categories
func GetDownloadCategories(c *gin.Context) {
	db := database.GetDB()
	var categories []model.DownloadCategory
	db.Where("status = ?", "active").Order("sort_order ASC").Find(&categories)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    categories,
	})
}
