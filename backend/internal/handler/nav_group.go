package handler

import (
	"net/http"
	"strconv"

	"anchorfinance/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// NavGroupHandler 前台导航分组管理
type NavGroupHandler struct {
	db *gorm.DB
}

func NewNavGroupHandler(db *gorm.DB) *NavGroupHandler {
	return &NavGroupHandler{db: db}
}

// List 获取导航分组列表
// GET /admin/nav-groups
func (h *NavGroupHandler) List(c *gin.Context) {
	var groups []model.NavGroup
	if err := h.db.Order("`order` ASC, id ASC").Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	// 统计每个分组的产品数量
	for i := range groups {
		var count int64
		h.db.Model(&model.NavGroupProduct{}).Where("nav_group_id = ?", groups[i].ID).Count(&count)
		groups[i].ProductCount = int(count)
	}

	c.JSON(http.StatusOK, gin.H{"data": groups})
}

// Create 创建导航分组
// POST /admin/nav-groups
func (h *NavGroupHandler) Create(c *gin.Context) {
	var req struct {
		Groupname string `json:"groupname" binding:"required"`
		FaIcon    string `json:"fa_icon"`
		Order     int    `json:"order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	group := model.NavGroup{
		Groupname: req.Groupname,
		FaIcon:    req.FaIcon,
		Order:     req.Order,
	}

	if err := h.db.Create(&group).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": group, "message": "创建成功"})
}

// Update 更新导航分组
// PUT /admin/nav-groups/:id
func (h *NavGroupHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req struct {
		Groupname string `json:"groupname"`
		FaIcon    string `json:"fa_icon"`
		Order     int    `json:"order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.db.Model(&model.NavGroup{}).Where("id = ?", id).Updates(map[string]interface{}{
		"groupname": req.Groupname,
		"fa_icon":   req.FaIcon,
		"order":     req.Order,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除导航分组
// DELETE /admin/nav-groups/:id
func (h *NavGroupHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	// 删除关联关系
	h.db.Where("nav_group_id = ?", id).Delete(&model.NavGroupProduct{})

	// 删除分组
	if err := h.db.Delete(&model.NavGroup{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetProducts 获取分组关联的产品
// GET /admin/nav-groups/:id/products
func (h *NavGroupHandler) GetProducts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var products []model.Product
	if err := h.db.Joins("JOIN nav_group_products ON nav_group_products.product_id = products.id").
		Where("nav_group_products.nav_group_id = ?", id).
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取产品列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}

// UpdateProducts 更新分组关联的产品
// PUT /admin/nav-groups/:id/products
func (h *NavGroupHandler) UpdateProducts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req struct {
		ProductIDs []uint `json:"product_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	tx := h.db.Begin()

	// 删除旧关联
	if err := tx.Where("nav_group_id = ?", id).Delete(&model.NavGroupProduct{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	// 添加新关联
	for _, productID := range req.ProductIDs {
		if err := tx.Create(&model.NavGroupProduct{
			NavGroupID: uint(id),
			ProductID:  productID,
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// GetPublicNavGroups 获取公开的导航分组（前端用）
// GET /nav-groups
func (h *NavGroupHandler) GetPublicNavGroups(c *gin.Context) {
	var groups []model.NavGroup
	if err := h.db.Preload("Products", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", "active")
	}).Order("`order` ASC, id ASC").Find(&groups).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": groups})
}

// GetUserNavConfig 获取用户的导航配置
// GET /user/nav-config
func (h *NavGroupHandler) GetUserNavConfig(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var configs []model.NavGroupUser
	if err := h.db.Where("uid = ?", userID).Find(&configs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": configs})
}

// UpdateUserNavConfig 更新用户的导航配置
// PUT /user/nav-config
func (h *NavGroupHandler) UpdateUserNavConfig(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	var req []struct {
		GroupID uint `json:"group_id"`
		IsShow  bool `json:"is_show"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	tx := h.db.Begin()

	// 删除旧配置
	if err := tx.Where("uid = ?", userID).Delete(&model.NavGroupUser{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	// 添加新配置
	for _, config := range req {
		if err := tx.Create(&model.NavGroupUser{
			Uid:     userID.(uint),
			GroupID: config.GroupID,
			IsShow:  config.IsShow,
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
