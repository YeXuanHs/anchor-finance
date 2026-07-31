package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// LogCleanerHandler 日志清理处理器
type LogCleanerHandler struct {
	cleaner *service.LogCleaner
	log     *logger.Logger
}

// NewLogCleanerHandler 创建日志清理处理器
func NewLogCleanerHandler(cleaner *service.LogCleaner, log *logger.Logger) *LogCleanerHandler {
	return &LogCleanerHandler{cleaner: cleaner, log: log}
}

// GetStats 获取日志统计
func (h *LogCleanerHandler) GetStats(c *gin.Context) {
	stats := h.cleaner.GetLogStats()
	response.Success(c, stats)
}

// CleanByDays 按天数清理
func (h *LogCleanerHandler) CleanByDays(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days < 1 {
		days = 30
	}

	count, err := h.cleaner.CleanByDays(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned": count,
		"days":    days,
	})
}

// CleanByCount 按数量保留
func (h *LogCleanerHandler) CleanByCount(c *gin.Context) {
	keepCount, _ := strconv.Atoi(c.DefaultQuery("keep_count", "10000"))
	if keepCount < 100 {
		keepCount = 100
	}

	count, err := h.cleaner.CleanByCount(keepCount)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned":     count,
		"keep_count":  keepCount,
	})
}

// CleanByModule 按模块清理
func (h *LogCleanerHandler) CleanByModule(c *gin.Context) {
	module := c.Query("module")
	if module == "" {
		response.BadRequest(c, "模块名不能为空")
		return
	}

	count, err := h.cleaner.CleanByModule(module)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned": count,
		"module":  module,
	})
}

// CleanByStatus 按状态清理
func (h *LogCleanerHandler) CleanByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		response.BadRequest(c, "状态不能为空")
		return
	}

	count, err := h.cleaner.CleanByStatus(status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned": count,
		"status":  status,
	})
}

// CleanExpired 清理过期日志
func (h *LogCleanerHandler) CleanExpired(c *gin.Context) {
	count, err := h.cleaner.CleanExpired()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned": count,
	})
}

// CleanAll 清理所有日志（危险操作）
func (h *LogCleanerHandler) CleanAll(c *gin.Context) {
	// 需要二次确认
	confirm := c.Query("confirm")
	if confirm != "yes" {
		response.BadRequest(c, "请确认清理操作，添加参数 confirm=yes")
		return
	}

	count, err := h.cleaner.CleanByDays(0) // 0天表示清理所有
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned": count,
		"warning": "已清理所有日志",
	})
}
