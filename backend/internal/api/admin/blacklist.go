package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetBlacklist 获取黑名单列表
// GET /api/admin/blacklist
func GetBlacklist(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	blacklistType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Blacklist{})

	if blacklistType != "" {
		query = query.Where("type = ?", blacklistType)
	}

	var total int64
	query.Count(&total)

	var items []model.Blacklist
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)

	if items == nil {
		items = []model.Blacklist{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      items,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateBlacklist 添加黑名单
// POST /api/admin/blacklist
func CreateBlacklist(c *gin.Context) {
	var req struct {
		Type      string `json:"type" binding:"required"`
		Value     string `json:"value" binding:"required"`
		Reason    string `json:"reason"`
		ExpiresAt string `json:"expires_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 验证类型
	validTypes := map[string]bool{
		"ip": true, "email": true, "phone": true, "username": true,
	}
	if !validTypes[req.Type] {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的类型，只支持: ip, email, phone, username"})
		return
	}

	adminID, _ := c.Get("user_id")

	db := database.GetDB()
	item := model.Blacklist{
		Type:    req.Type,
		Value:   req.Value,
		Reason:  req.Reason,
		AdminID: adminID.(uint),
	}

	if req.ExpiresAt != "" {
		t, err := time.Parse("2006-01-02 15:04:05", req.ExpiresAt)
		if err == nil {
			item.ExpiresAt = &t
		}
	}

	// 检查是否已存在
	var count int64
	db.Model(&model.Blacklist{}).Where("type = ? AND value = ?", req.Type, req.Value).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "该记录已存在"})
		return
	}

	if err := db.Create(&item).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "添加成功",
		"data": gin.H{
			"id": item.ID,
		},
	})
}

// DeleteBlacklist 删除黑名单
// DELETE /api/admin/blacklist/:id
func DeleteBlacklist(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var item model.Blacklist
	if err := db.First(&item, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "记录不存在", "data": nil})
		return
	}

	db.Delete(&item)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
