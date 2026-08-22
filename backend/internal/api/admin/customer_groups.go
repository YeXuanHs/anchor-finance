package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetCustomerGroupList 获取客户分组列表（复用MemberLevel）
// GET /api/admin/customer-groups
func GetCustomerGroupList(c *gin.Context) {
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
	db.Model(&model.MemberLevel{}).Count(&total)

	var groups []model.MemberLevel
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("sort_order ASC").Find(&groups)

	if groups == nil {
		groups = []model.MemberLevel{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      groups,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateCustomerGroup 创建客户分组
// POST /api/admin/customer-groups
func CreateCustomerGroup(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Discount    float64 `json:"discount"`
		SortOrder   int     `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	// 折扣比例校验（0-100）
	if req.Discount < 0 || req.Discount > 100 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "折扣需在0-100之间", "data": nil})
		return
	}

	db := database.GetDB()
	group := model.MemberLevel{
		Name:        req.Name,
		Description: req.Description,
		Discount:    req.Discount,
		SortOrder:   req.SortOrder,
		Status:      "active",
	}

	if err := db.Create(&group).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": group.ID,
		},
	})
}

// UpdateCustomerGroup 更新客户分组
// PUT /api/admin/customer-groups/:id
func UpdateCustomerGroup(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Discount    float64 `json:"discount"`
		SortOrder   int     `json:"sort_order"`
		Status      string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var group model.MemberLevel
	if err := db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "分组不存在", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Discount >= 0 && req.Discount <= 100 {
		updates["discount"] = req.Discount
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&group).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteCustomerGroup 删除客户分组
// DELETE /api/admin/customer-groups/:id
func DeleteCustomerGroup(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var group model.MemberLevel
	if err := db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "分组不存在", "data": nil})
		return
	}

	db.Delete(&group)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
