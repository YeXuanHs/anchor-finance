package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetTicketPrereplyList 获取工单预回复列表
// GET /api/admin/ticket-prereplies
func GetTicketPrereplyList(c *gin.Context) {
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
	query := db.Model(&model.TicketPrereply{}).Where("status = ?", "active")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	var total int64
	query.Count(&total)

	var prereplies []model.TicketPrereply
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("sort_order ASC").Find(&prereplies)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      prereplies,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetTicketPrereplyCategoryList 获取工单预回复分类列表
// GET /api/admin/ticket-prereply-categories
func GetTicketPrereplyCategoryList(c *gin.Context) {
	db := database.GetDB()
	var categories []model.TicketPrereplyCategory
	db.Where("status = ?", "active").Order("sort_order ASC").Find(&categories)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    categories,
	})
}

// CreateTicketPrereplyCategory 创建工单预回复分类
// POST /api/admin/ticket-prereply-categories
func CreateTicketPrereplyCategory(c *gin.Context) {
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
	category := model.TicketPrereplyCategory{
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

// UpdateTicketPrereplyCategory 更新工单预回复分类
// PUT /api/admin/ticket-prereply-categories/:id
func UpdateTicketPrereplyCategory(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name"`
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
	var category model.TicketPrereplyCategory
	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "分类不存在",
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

	db.Model(&category).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteTicketPrereplyCategory 删除工单预回复分类
// DELETE /api/admin/ticket-prereply-categories/:id
func DeleteTicketPrereplyCategory(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var category model.TicketPrereplyCategory
	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "分类不存在",
		})
		return
	}

	db.Delete(&category)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// CreateTicketPrereply 创建工单预回复
// POST /api/admin/ticket-prereplies
func CreateTicketPrereply(c *gin.Context) {
	var req struct {
		CategoryID uint   `json:"category_id"`
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		SortOrder  int    `json:"sort_order"`
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
	prereply := model.TicketPrereply{
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Content:    req.Content,
		SortOrder:  req.SortOrder,
		Status:     "active",
	}

	if err := db.Create(&prereply).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建预回复失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": prereply.ID,
		},
	})
}

// UpdateTicketPrereply 更新工单预回复
// PUT /api/admin/ticket-prereplies/:id
func UpdateTicketPrereply(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		CategoryID uint   `json:"category_id"`
		Title      string `json:"title"`
		Content    string `json:"content"`
		SortOrder  int    `json:"sort_order"`
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
	var prereply model.TicketPrereply
	if err := db.First(&prereply, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "预回复不存在",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.CategoryID > 0 {
		updates["category_id"] = req.CategoryID
	}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}

	db.Model(&prereply).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteTicketPrereply 删除工单预回复
// DELETE /api/admin/ticket-prereplies/:id
func DeleteTicketPrereply(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var prereply model.TicketPrereply
	if err := db.First(&prereply, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "预回复不存在",
		})
		return
	}

	db.Delete(&prereply)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// SearchTicketPrereply 搜索工单预回复
// POST /api/admin/ticket-prereplies/search
func SearchTicketPrereply(c *gin.Context) {
	var req struct {
		Keyword string `json:"keyword" binding:"required"`
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
	var prereplies []model.TicketPrereply
	db.Where("status = ? AND (title LIKE ? OR content LIKE ?)", "active",
		"%"+req.Keyword+"%", "%"+req.Keyword+"%").
		Order("sort_order ASC").
		Limit(20).
		Find(&prereplies)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    prereplies,
	})
}
