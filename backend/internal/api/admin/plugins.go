package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetPluginList 获取插件列表
// GET /api/admin/plugins
func GetPluginList(c *gin.Context) {
	domain := c.Query("domain")
	status := c.Query("status")

	db := database.GetDB()
	query := db.Model(&model.Plugin{})

	if domain != "" {
		query = query.Where("domain = ?", domain)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var plugins []model.Plugin
	query.Order("domain ASC, name ASC").Find(&plugins)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    plugins,
	})
}

// GetPluginDetail 获取插件详情
// GET /api/admin/plugins/:id
func GetPluginDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的插件ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var plugin model.Plugin
	if err := db.First(&plugin, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "插件不存在",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    plugin,
	})
}

// EnablePlugin 启用插件
// POST /api/admin/plugins/:id/enable
func EnablePlugin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的插件ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var plugin model.Plugin
	if err := db.First(&plugin, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "插件不存在",
			"data":    nil,
		})
		return
	}

	if plugin.Status == "active" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "插件已启用",
			"data":    nil,
		})
		return
	}

	db.Model(&plugin).Update("status", "active")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "启用成功",
		"data":    nil,
	})
}

// DisablePlugin 禁用插件
// POST /api/admin/plugins/:id/disable
func DisablePlugin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的插件ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var plugin model.Plugin
	if err := db.First(&plugin, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "插件不存在",
			"data":    nil,
		})
		return
	}

	if plugin.Status == "inactive" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "插件已禁用",
			"data":    nil,
		})
		return
	}

	db.Model(&plugin).Update("status", "inactive")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "禁用成功",
		"data":    nil,
	})
}

// UninstallPlugin 卸载插件
// DELETE /api/admin/plugins/:id
func UninstallPlugin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的插件ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var plugin model.Plugin
	if err := db.First(&plugin, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "插件不存在",
			"data":    nil,
		})
		return
	}

	db.Delete(&plugin)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "卸载成功",
		"data":    nil,
	})
}

// GetPluginConfig 获取插件配置
// GET /api/admin/plugins/:id/config
func GetPluginConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的插件ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var plugin model.Plugin
	if err := db.First(&plugin, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "插件不存在",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"config_json": plugin.ConfigJSON,
		},
	})
}

// UpdatePluginConfig 更新插件配置
// PUT /api/admin/plugins/:id/config
func UpdatePluginConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的插件ID",
			"data":    nil,
		})
		return
	}

	var req struct {
		ConfigJSON string `json:"config_json"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var plugin model.Plugin
	if err := db.First(&plugin, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "插件不存在",
			"data":    nil,
		})
		return
	}

	db.Model(&plugin).Update("config_json", req.ConfigJSON)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "配置更新成功",
		"data":    nil,
	})
}
