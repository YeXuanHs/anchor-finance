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

// GetContentSummary 获取内容统计
// GET /api/admin/content/summary
func GetContentSummary(c *gin.Context) {
	db := database.GetDB()

	// 统计新闻数量
	var newsCount int64
	db.Model(&model.News{}).Count(&newsCount)

	// 统计知识库文章数量
	var articleCount int64
	db.Model(&model.KnowledgeArticle{}).Count(&articleCount)

	// 统计下载数量
	var downloadCount int64
	db.Model(&model.Download{}).Count(&downloadCount)

	// 统计新闻分类数量
	var newsCategoryCount int64
	db.Model(&model.NewsCategory{}).Count(&newsCategoryCount)

	// 统计知识库分类数量
	var knowledgeCategoryCount int64
	db.Model(&model.KnowledgeCategory{}).Count(&knowledgeCategoryCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"news_count":               newsCount,
			"article_count":            articleCount,
			"download_count":           downloadCount,
			"news_category_count":      newsCategoryCount,
			"knowledge_category_count": knowledgeCategoryCount,
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
