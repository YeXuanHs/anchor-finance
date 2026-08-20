package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type LogHandler struct{}

func NewLogHandler() *LogHandler {
	return &LogHandler{}
}

// GetLogs 获取日志列表
func (h *LogHandler) GetLogs(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"list": []interface{}{}, "total": 0})
}

// RegisterRoutes 注册路由
func (h *LogHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/logs", h.GetLogs)
}
