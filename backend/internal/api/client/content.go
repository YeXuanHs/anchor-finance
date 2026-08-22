package client

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetClientHomeHero 用户前台获取首页Hero配置
// GET /api/client/home-hero
func GetClientHomeHero(c *gin.Context) {
	db := database.GetDB()
	var hero model.HomeHero
	if err := db.Where("is_active = ?", true).Order("id DESC").First(&hero).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 0, "message": "success",
			"data": gin.H{
				"slides":   []interface{}{},
				"features": []interface{}{},
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": hero.Config})
}

// GetContentOverview 获取首页内容概览
// GET /api/client/content/overview
func GetContentOverview(c *gin.Context) {
	db := database.GetDB()

	var notices []model.News
	db.Where("status = ?", "published").Order("created_at DESC").Limit(5).Find(&notices)

	var articles []model.KnowledgeArticle
	db.Where("status = ?", "published").Order("view_count DESC").Limit(5).Find(&articles)

	var newsList []model.News
	db.Where("status = ?", "published").Order("created_at DESC").Limit(5).Find(&newsList)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{
			"notices":       notices,
			"help_articles": articles,
			"news":          newsList,
		},
	})
}

// GetNoticeDetail 获取公告详情
// GET /api/client/notices/:id
func GetNoticeDetail(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var notice model.News
	if err := db.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "公告不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": notice})
}

// GetHelpArticleDetail 获取帮助文章详情
// GET /api/client/help-articles/:id
func GetHelpArticleDetail(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var article model.KnowledgeArticle
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "文章不存在"})
		return
	}

	// 增加浏览次数
	db.Model(&article).Update("view_count", gorm.Expr("view_count + 1"))

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": article})
}

// GetNoticesUnreadCount 获取公告未读数（需登录）
// GET /api/client/notices/unread-count
func GetNoticesUnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var totalNotices int64
	db.Model(&model.News{}).Where("status = ?", "published").Count(&totalNotices)

	var readCount int64
	db.Model(&model.OperationLog{}).Where("user_id = ? AND action = ?", userID, "read_notice").Count(&readCount)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{"unread_count": totalNotices - readCount},
	})
}

// MarkAllNoticesRead 标记全部公告已读
// POST /api/client/notices/mark-all-read
func MarkAllNoticesRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	db.Create(&model.OperationLog{
		UserID: userID.(uint),
		Action: "read_notice",
		Detail: "mark_all_read",
	})

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
