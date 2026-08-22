package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetFriendlyLinkList 获取友情链接列表
// GET /api/admin/friendly-links
func GetFriendlyLinkList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	var total int64
	db.Model(&model.FriendlyLink{}).Count(&total)

	var links []model.FriendlyLink
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("sort_order ASC").Find(&links)

	if links == nil {
		links = []model.FriendlyLink{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      links,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateFriendlyLink 创建友情链接
// POST /api/admin/friendly-links
func CreateFriendlyLink(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		URL       string `json:"url" binding:"required"`
		Logo      string `json:"logo"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	db := database.GetDB()
	link := model.FriendlyLink{
		Name:      req.Name,
		URL:       req.URL,
		Logo:      req.Logo,
		SortOrder: req.SortOrder,
		Status:    "active",
	}

	if err := db.Create(&link).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": link.ID,
		},
	})
}

// UpdateFriendlyLink 更新友情链接
// PUT /api/admin/friendly-links/:id
func UpdateFriendlyLink(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name"`
		URL       string `json:"url"`
		Logo      string `json:"logo"`
		SortOrder int    `json:"sort_order"`
		Status    string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	db := database.GetDB()
	var link model.FriendlyLink
	if err := db.First(&link, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "链接不存在", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.URL != "" {
		updates["url"] = req.URL
	}
	if req.Logo != "" {
		updates["logo"] = req.Logo
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&link).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteFriendlyLink 删除友情链接
// DELETE /api/admin/friendly-links/:id
func DeleteFriendlyLink(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var link model.FriendlyLink
	if err := db.First(&link, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "链接不存在", "data": nil})
		return
	}

	db.Delete(&link)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
