package admin

import (
	"encoding/json"
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetHomeHero 获取首页Hero配置（MD 11.3：GET /api/admin/home-hero）
// 返回单条记录的JSON config字段，包含slides和features
func GetHomeHero(c *gin.Context) {
	db := database.GetDB()
	var hero model.HomeHero

	// 获取默认配置（或第一条active配置）
	if err := db.Where("is_default = ? AND status = ?", true, "active").First(&hero).Error; err != nil {
		// 如果没有默认配置，尝试获取任意一条
		if err := db.Where("status = ?", "active").Order("id ASC").First(&hero).Error; err != nil {
			// 返回空配置
			c.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "success",
				"data": gin.H{
					"id":         0,
					"name":       "default",
					"config":     `{"slides":[],"features":[]}`,
					"is_default": true,
					"status":     "active",
				},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    hero,
	})
}

// UpdateHomeHero 更新首页Hero配置（MD 11.3：PUT /api/admin/home-hero）
// 接收完整的config JSON，包含slides和features
func UpdateHomeHero(c *gin.Context) {
	var req struct {
		ID     uint   `json:"id"`     // 可选，指定更新哪条配置
		Name   string `json:"name"`   // 配置名称
		Config string `json:"config"` // JSON格式的Hero配置
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	// 验证config是合法JSON
	if req.Config == "" {
		req.Config = `{"slides":[],"features":[]}`
	}
	var configCheck model.HomeHeroConfig
	if err := json.Unmarshal([]byte(req.Config), &configCheck); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "config字段必须是合法的JSON格式",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var hero model.HomeHero

	if req.ID > 0 {
		// 更新指定配置
		if err := db.First(&hero, req.ID).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 404, "message": "配置不存在", "data": nil})
			return
		}
		updates := map[string]interface{}{
			"config": req.Config,
		}
		if req.Name != "" {
			updates["name"] = req.Name
		}
		db.Model(&hero).Updates(updates)
	} else {
		// 更新默认配置（不存在则创建）
		if err := db.Where("is_default = ?", true).First(&hero).Error; err != nil {
			// 创建新配置
			if req.Name == "" {
				req.Name = "default"
			}
			hero = model.HomeHero{
				Name:      req.Name,
				Config:    req.Config,
				IsDefault: true,
				Status:    "active",
			}
			db.Create(&hero)
		} else {
			updates := map[string]interface{}{
				"config": req.Config,
			}
			if req.Name != "" {
				updates["name"] = req.Name
			}
			db.Model(&hero).Updates(updates)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data": gin.H{
			"id": hero.ID,
		},
	})
}

// GetHomeHeroAssets 获取可用的轮播资源文件列表（MD 11.3：GET /api/admin/home-hero/assets）
// 返回 hero/images/ 和 hero/videos/ 下所有文件
func GetHomeHeroAssets(c *gin.Context) {
	db := database.GetDB()
	var files []model.MediaFile
	db.Where("mime_type LIKE ? OR mime_type LIKE ?", "image/%", "video/%").Order("id DESC").Limit(100).Find(&files)

	images := []string{}
	videos := []string{}
	for _, f := range files {
		if f.MimeType != "" && len(f.MimeType) >= 5 && f.MimeType[:5] == "video" {
			videos = append(videos, f.Path)
		} else {
			images = append(images, f.Path)
		}
	}
	if images == nil {
		images = []string{}
	}
	if videos == nil {
		videos = []string{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"images": images,
			"videos": videos,
		},
	})
}
