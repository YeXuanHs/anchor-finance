package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// AuditLogHandler 审计日志处理器
type AuditLogHandler struct {
	auditSvc *model.AuditLogService
	log      *logger.Logger
}

// NewAuditLogHandler 创建审计日志处理器
func NewAuditLogHandler(auditSvc *model.AuditLogService, log *logger.Logger) *AuditLogHandler {
	return &AuditLogHandler{
		auditSvc: auditSvc,
		log:      log,
	}
}

// List 查询审计日志列表
func (h *AuditLogHandler) List(c *gin.Context) {
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 32)
	userType := c.Query("user_type")
	action := c.Query("action")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := h.auditSvc.List(uint(userID), userType, action, page, pageSize)
	if err != nil {
		h.log.Errorf("查询审计日志失败: %v", err)
		response.ServerError(c, "查询失败")
		return
	}

	response.Success(c, gin.H{
		"list":      logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Get 获取审计日志详情
func (h *AuditLogHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "无效的 ID")
		return
	}

	log, err := h.auditSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "日志不存在")
		return
	}

	response.Success(c, log)
}

// Stats 获取审计日志统计
func (h *AuditLogHandler) Stats(c *gin.Context) {
	stats, err := h.auditSvc.GetStats()
	if err != nil {
		h.log.Errorf("获取审计日志统计失败: %v", err)
		response.ServerError(c, "获取统计失败")
		return
	}

	response.Success(c, stats)
}

// CleanOldLogs 清理旧日志
func (h *AuditLogHandler) CleanOldLogs(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))
	if days < 30 {
		days = 30
	}

	count, err := h.auditSvc.CleanOldLogs(days)
	if err != nil {
		h.log.Errorf("清理审计日志失败: %v", err)
		response.ServerError(c, "清理失败")
		return
	}

	response.Success(c, gin.H{
		"deleted": count,
		"days":    days,
	})
}
