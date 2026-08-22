package admin

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetHomeHero 获取首页英雄区
// GET /api/admin/site/home-hero
func GetHomeHero(c *gin.Context) {
	db := database.GetDB()
	var heroes []model.HomeHero
	db.Where("status = ?", "active").Order("sort_order ASC").Find(&heroes)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    heroes,
	})
}

// UpdateHomeHero 更新首页英雄区
// POST /api/admin/site/home-hero
func UpdateHomeHero(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Subtitle  string `json:"subtitle"`
		ImageURL  string `json:"image_url"`
		LinkURL   string `json:"link_url"`
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
	hero := model.HomeHero{
		Title:     req.Title,
		Subtitle:  req.Subtitle,
		ImageURL:  req.ImageURL,
		LinkURL:   req.LinkURL,
		SortOrder: req.SortOrder,
		Status:    "active",
	}

	if err := db.Create(&hero).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建失败: " + err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": hero.ID,
		},
	})
}
