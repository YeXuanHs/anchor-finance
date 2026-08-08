package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct{}

func NewSettingHandler() *SettingHandler {
	return &SettingHandler{}
}

// GetSettings 获取设置
func (h *SettingHandler) GetSettings(c *gin.Context) {
	group := c.Param("group")
	c.JSON(http.StatusOK, gin.H{"group": group})
}

// UpdateSettings 更新设置
func (h *SettingHandler) UpdateSettings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "设置更新成功"})
}

// RegisterRoutes 注册路由
func (h *SettingHandler) RegisterRoutes(r *gin.RouterGroup) {
	setting := r.Group("/settings")
	{
		setting.GET("/:group", h.GetSettings)
		setting.POST("/:group", h.UpdateSettings)
	}
}
