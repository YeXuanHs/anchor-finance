package handler

import (
	"net/http"

	"anchorfinance/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ExpressHandler 快递管理
type ExpressHandler struct {
	db *gorm.DB
}

func NewExpressHandler(db *gorm.DB) *ExpressHandler {
	return &ExpressHandler{db: db}
}

// List 获取快递列表
// GET /admin/expresses
func (h *ExpressHandler) List(c *gin.Context) {
	var expresses []model.Express
	if err := h.db.Order("sort_order ASC, id ASC").Find(&expresses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": expresses})
}

// Create 创建快递
// POST /admin/expresses
func (h *ExpressHandler) Create(c *gin.Context) {
	var req struct {
		Name      string  `json:"name" binding:"required"`
		Code      string  `json:"code"`
		Price     float64 `json:"price"`
		IsActive  *bool   `json:"is_active"`
		SortOrder int     `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	express := model.Express{
		Name:      req.Name,
		Code:      req.Code,
		Price:     req.Price,
		IsActive:  isActive,
		SortOrder: req.SortOrder,
	}
	if err := h.db.Create(&express).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": express, "message": "创建成功"})
}

// Update 更新快递
// PUT /admin/expresses/:id
func (h *ExpressHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name      string  `json:"name"`
		Code      string  `json:"code"`
		Price     float64 `json:"price"`
		IsActive  *bool   `json:"is_active"`
		SortOrder int     `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	updates := map[string]interface{}{
		"name":       req.Name,
		"code":       req.Code,
		"price":      req.Price,
		"sort_order": req.SortOrder,
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := h.db.Model(&model.Express{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除快递
// DELETE /admin/expresses/:id
func (h *ExpressHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&model.Express{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ========== 取消原因 ==========

// CancelReasonHandler 取消原因管理
type CancelReasonHandler struct {
	db *gorm.DB
}

func NewCancelReasonHandler(db *gorm.DB) *CancelReasonHandler {
	return &CancelReasonHandler{db: db}
}

// List 获取取消原因列表
// GET /admin/cancel-reasons
func (h *CancelReasonHandler) List(c *gin.Context) {
	var reasons []model.CancelReason
	if err := h.db.Order("sort_order ASC, id ASC").Find(&reasons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reasons})
}

// Create 创建取消原因
// POST /admin/cancel-reasons
func (h *CancelReasonHandler) Create(c *gin.Context) {
	var req struct {
		Reason    string `json:"reason" binding:"required"`
		IsActive  *bool  `json:"is_active"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	reason := model.CancelReason{
		Reason:    req.Reason,
		IsActive:  isActive,
		SortOrder: req.SortOrder,
	}
	if err := h.db.Create(&reason).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": reason, "message": "创建成功"})
}

// Update 更新取消原因
// PUT /admin/cancel-reasons/:id
func (h *CancelReasonHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Reason    string `json:"reason"`
		IsActive  *bool  `json:"is_active"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	updates := map[string]interface{}{
		"reason":     req.Reason,
		"sort_order": req.SortOrder,
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := h.db.Model(&model.CancelReason{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除取消原因
// DELETE /admin/cancel-reasons/:id
func (h *CancelReasonHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&model.CancelReason{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ========== BaseInfo Handler ==========

// BaseInfoHandler 首页基本信息管理
type BaseInfoHandler struct {
	db *gorm.DB
}

func NewBaseInfoHandler(db *gorm.DB) *BaseInfoHandler {
	return &BaseInfoHandler{db: db}
}

// List 获取首页信息列表
// GET /admin/base-infos
func (h *BaseInfoHandler) List(c *gin.Context) {
	var infos []model.BaseInfo
	if err := h.db.Order("sort_order ASC, id ASC").Find(&infos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": infos})
}

// Create 创建首页信息
// POST /admin/base-infos
func (h *BaseInfoHandler) Create(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		Desc      string `json:"desc"`
		Icon      string `json:"icon"`
		Content   string `json:"content"`
		SortOrder int    `json:"sort_order"`
		IsActive  *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	info := model.BaseInfo{
		Name:      req.Name,
		Desc:      req.Desc,
		Icon:      req.Icon,
		Content:   req.Content,
		SortOrder: req.SortOrder,
		IsActive:  isActive,
	}
	if err := h.db.Create(&info).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": info, "message": "创建成功"})
}

// Update 更新首页信息
// PUT /admin/base-infos/:id
func (h *BaseInfoHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name      string `json:"name"`
		Desc      string `json:"desc"`
		Icon      string `json:"icon"`
		Content   string `json:"content"`
		SortOrder int    `json:"sort_order"`
		IsActive  *bool  `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	updates := map[string]interface{}{
		"name":       req.Name,
		"desc":       req.Desc,
		"icon":       req.Icon,
		"content":    req.Content,
		"sort_order": req.SortOrder,
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := h.db.Model(&model.BaseInfo{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除首页信息
// DELETE /admin/base-infos/:id
func (h *BaseInfoHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&model.BaseInfo{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetActive 获取启用的首页信息（前台）
// GET /base-infos
func (h *BaseInfoHandler) GetActive(c *gin.Context) {
	var infos []model.BaseInfo
	if err := h.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&infos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": infos})
}
