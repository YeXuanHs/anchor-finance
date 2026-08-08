package handler

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// GetSystemInfo 获取系统信息
func (h *SystemHandler) GetSystemInfo(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":     "1.0.0",
		"go_version":  runtime.Version(),
		"os":          runtime.GOOS + "/" + runtime.GOARCH,
		"server_time": "2026-08-08 17:00:00",
	})
}

// CheckUpdate 检查更新
func (h *SystemHandler) CheckUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"has_update":      false,
		"latest_version":  "1.0.0",
		"release_date":    "2026-08-08",
		"changelog":       "",
	})
}

// RegisterRoutes 注册路由
func (h *SystemHandler) RegisterRoutes(r *gin.RouterGroup) {
	system := r.Group("/system")
	{
		system.GET("/info", h.GetSystemInfo)
		system.GET("/check-update", h.CheckUpdate)
	}
}
