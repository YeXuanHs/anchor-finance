package handler

import (
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// NotificationHandler 通知管理处理器
type NotificationHandler struct {
	svc *service.NotificationService
	log *logger.Logger
}

// NewNotificationHandler 创建通知处理器
func NewNotificationHandler(svc *service.NotificationService, log *logger.Logger) *NotificationHandler {
	return &NotificationHandler{svc: svc, log: log}
}

// GetStats 获取通知统计
func (h *NotificationHandler) GetStats(c *gin.Context) {
	stats := h.svc.GetDeduplicator().GetStats()
	response.Success(c, stats)
}

// ResetEvent 重置事件（允许重新通知）
func (h *NotificationHandler) ResetEvent(c *gin.Context) {
	var req struct {
		EventType string `json:"event_type" binding:"required"`
		TargetID  string `json:"target_id" binding:"required"`
		Title     string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.ResetEvent(req.EventType, req.TargetID, req.Title); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "事件已重置，下次失败会重新通知")
}

// CleanAll 清空所有通知记录
func (h *NotificationHandler) CleanAll(c *gin.Context) {
	if err := h.svc.GetDeduplicator().CleanAll(); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "已清空所有通知记录")
}
