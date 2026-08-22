package admin

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetMemberLevelList 获取会员等级列表
// GET /api/admin/member-levels
func GetMemberLevelList(c *gin.Context) {
	db := database.GetDB()
	var levels []model.MemberLevel
	db.Where("status = ?", "active").Order("sort_order ASC").Find(&levels)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    levels,
	})
}

// CreateMemberLevel 创建会员等级
// POST /api/admin/member-levels
func CreateMemberLevel(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		Discount    float64 `json:"discount"`
		MinPoints   int     `json:"min_points"`
		SortOrder   int     `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data": nil,
		})
		return
	}

	if req.Discount == 0 {
		req.Discount = 100
	}

	db := database.GetDB()
	level := model.MemberLevel{
		Name:        req.Name,
		Description: req.Description,
		Discount:    req.Discount,
		MinPoints:   req.MinPoints,
		SortOrder:   req.SortOrder,
		Status:      "active",
	}

	if err := db.Create(&level).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建等级失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": level.ID,
		},
	})
}

// UpdateMemberLevel 更新会员等级
// PUT /api/admin/member-levels/:id
func UpdateMemberLevel(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Discount    float64 `json:"discount"`
		MinPoints   int     `json:"min_points"`
		SortOrder   int     `json:"sort_order"`
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
	var level model.MemberLevel
	if err := db.First(&level, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "等级不存在",
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Discount > 0 {
		updates["discount"] = req.Discount
	}
	if req.MinPoints > 0 {
		updates["min_points"] = req.MinPoints
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}

	db.Model(&level).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteMemberLevel 删除会员等级
// DELETE /api/admin/member-levels/:id
func DeleteMemberLevel(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var level model.MemberLevel
	if err := db.First(&level, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "等级不存在",
		})
		return
	}

	db.Delete(&level)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}
