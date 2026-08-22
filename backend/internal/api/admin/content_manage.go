package admin

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// ========== 新闻管理CRUD ==========

// CreateNews 创建新闻
// POST /api/admin/news
func CreateNews(c *gin.Context) {
	var req struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content"`
		CategoryID uint   `json:"category_id"`
		Author     string `json:"author"`
		Status     string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	if req.Status == "" {
		req.Status = "draft"
	}

	db := database.GetDB()
	news := model.News{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		Author:     req.Author,
		Status:     req.Status,
	}

	if err := db.Create(&news).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建新闻失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": news.ID,
		},
	})
}

// UpdateNews 更新新闻
// PUT /api/admin/news/:id
func UpdateNews(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		CategoryID uint   `json:"category_id"`
		Author     string `json:"author"`
		Status     string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var news model.News
	if err := db.First(&news, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "新闻不存在",
			"data": nil,
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Author != "" {
		updates["author"] = req.Author
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&news).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteNews 删除新闻
// DELETE /api/admin/news/:id
func DeleteNews(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var news model.News
	if err := db.First(&news, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "新闻不存在",
			"data": nil,
		})
		return
	}

	db.Delete(&news)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// ========== 新闻分类CRUD ==========

// CreateNewsCategory 创建新闻分类
// POST /api/admin/news-categories
func CreateNewsCategory(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	category := model.NewsCategory{
		Name:      req.Name,
		SortOrder: req.SortOrder,
		Status:    "active",
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建分类失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": category.ID,
		},
	})
}

// UpdateNewsCategory 更新新闻分类
// PUT /api/admin/news-categories/:id
func UpdateNewsCategory(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name"`
		SortOrder int    `json:"sort_order"`
		Status    string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var category model.NewsCategory
	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "分类不存在",
			"data": nil,
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&category).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteNewsCategory 删除新闻分类
// DELETE /api/admin/news-categories/:id
func DeleteNewsCategory(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var category model.NewsCategory
	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "分类不存在",
			"data": nil,
		})
		return
	}

	db.Delete(&category)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// ========== 知识库分类CRUD ==========

// CreateKnowledgeCategory 创建知识库分类
// POST /api/admin/knowledge/categories
func CreateKnowledgeCategory(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	category := model.KnowledgeCategory{
		Name:      req.Name,
		SortOrder: req.SortOrder,
		Status:    "active",
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建分类失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": category.ID,
		},
	})
}

// ========== 知识库文章CRUD ==========

// CreateKnowledgeArticle 创建知识库文章
// POST /api/admin/knowledge/articles
func CreateKnowledgeArticle(c *gin.Context) {
	var req struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content"`
		CategoryID uint   `json:"category_id"`
		Status     string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	if req.Status == "" {
		req.Status = "draft"
	}

	db := database.GetDB()
	article := model.KnowledgeArticle{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		Status:     req.Status,
	}

	if err := db.Create(&article).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建文章失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": article.ID,
		},
	})
}

// UpdateKnowledgeArticle 更新知识库文章
// PUT /api/admin/knowledge/articles/:id
func UpdateKnowledgeArticle(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		CategoryID uint   `json:"category_id"`
		Status     string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var article model.KnowledgeArticle
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "文章不存在",
			"data": nil,
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&article).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteKnowledgeArticle 删除知识库文章
// DELETE /api/admin/knowledge/articles/:id
func DeleteKnowledgeArticle(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var article model.KnowledgeArticle
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "文章不存在",
			"data": nil,
		})
		return
	}

	db.Delete(&article)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// ========== 下载CRUD ==========

// CreateDownload 创建下载
// POST /api/admin/downloads
func CreateDownload(c *gin.Context) {
	var req struct {
		Title      string `json:"title" binding:"required"`
		FileName   string `json:"file_name"`
		FileURL    string `json:"file_url"`
		CategoryID uint   `json:"category_id"`
		Status     string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	if req.Status == "" {
		req.Status = "active"
	}

	db := database.GetDB()
	download := model.Download{
		Title:      req.Title,
		FileName:   req.FileName,
		FileURL:    req.FileURL,
		CategoryID: req.CategoryID,
		Status:     req.Status,
	}

	if err := db.Create(&download).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建下载失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": download.ID,
		},
	})
}

// UpdateDownload 更新下载
// PUT /api/admin/downloads/:id
func UpdateDownload(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Title      string `json:"title"`
		FileName   string `json:"file_name"`
		FileURL    string `json:"file_url"`
		CategoryID uint   `json:"category_id"`
		Status     string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var download model.Download
	if err := db.First(&download, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "下载不存在",
			"data": nil,
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.FileName != "" {
		updates["file_name"] = req.FileName
	}
	if req.FileURL != "" {
		updates["file_url"] = req.FileURL
	}
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&download).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteDownload 删除下载
// DELETE /api/admin/downloads/:id
func DeleteDownload(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var download model.Download
	if err := db.First(&download, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "下载不存在",
			"data": nil,
		})
		return
	}

	db.Delete(&download)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// ========== 下载分类CRUD ==========

// CreateDownloadCategory 创建下载分类
// POST /api/admin/downloads/categories
func CreateDownloadCategory(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	category := model.DownloadCategory{
		Name:      req.Name,
		SortOrder: req.SortOrder,
		Status:    "active",
	}

	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建分类失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": category.ID,
		},
	})
}
