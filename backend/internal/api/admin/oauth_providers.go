package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetOAuthProviderList 获取第三方登录提供商列表
// GET /api/admin/oauth-providers
func GetOAuthProviderList(c *gin.Context) {
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
	db.Model(&model.OAuthProvider{}).Count(&total)

	var providers []model.OAuthProvider
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("sort_order ASC").Find(&providers)

	if providers == nil {
		providers = []model.OAuthProvider{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      providers,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateOAuthProvider 创建第三方登录提供商
// POST /api/admin/oauth-providers
func CreateOAuthProvider(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Code      string `json:"code" binding:"required"`
		Icon      string `json:"icon"`
		Config    string `json:"config"`
		SortOrder int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	// 检查Code是否重复
	db := database.GetDB()
	var count int64
	db.Model(&model.OAuthProvider{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "提供商代码已存在", "data": nil})
		return
	}

	provider := model.OAuthProvider{
		Name:      req.Name,
		Code:      req.Code,
		Icon:      req.Icon,
		Config:    req.Config,
		SortOrder: req.SortOrder,
		Status:    "active",
	}

	if err := db.Create(&provider).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": provider.ID,
		},
	})
}

// UpdateOAuthProvider 更新第三方登录提供商
// PUT /api/admin/oauth-providers/:id
func UpdateOAuthProvider(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string `json:"name"`
		Icon      string `json:"icon"`
		Config    string `json:"config"`
		SortOrder int    `json:"sort_order"`
		Status    string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var provider model.OAuthProvider
	if err := db.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "提供商不存在", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.Config != "" {
		updates["config"] = req.Config
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&provider).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteOAuthProvider 删除第三方登录提供商
// DELETE /api/admin/oauth-providers/:id
func DeleteOAuthProvider(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var provider model.OAuthProvider
	if err := db.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "提供商不存在", "data": nil})
		return
	}

	db.Delete(&provider)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
