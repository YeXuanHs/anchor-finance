package handler

import (
	"strconv"

	"github.com/anchor-finance/backend/internal/service"
	"github.com/anchor-finance/backend/pkg/logger"
	"github.com/anchor-finance/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notifSvc *service.NotificationService
	log      *logger.Logger
}

func NewNotificationHandler(notifSvc *service.NotificationService, log *logger.Logger) *NotificationHandler {
	return &NotificationHandler{notifSvc: notifSvc, log: log}
}

// GetUserNotifications returns paginated notifications for the current user.
// GET /user/notifications
func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	onlyUnread := c.Query("unread") == "1"

	logs, total, err := h.notifSvc.GetUserNotifications(userID, page, pageSize, onlyUnread)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// MarkRead marks a single notification as read.
// POST /user/notifications/:id/read
func (h *NotificationHandler) MarkRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid notification id")
		return
	}

	if err := h.notifSvc.MarkRead(userID, uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "marked as read")
}

// MarkAllRead marks all notifications as read for the current user.
// POST /user/notifications/read-all
func (h *NotificationHandler) MarkAllRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := h.notifSvc.MarkAllRead(userID); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "all notifications marked as read")
}

// AdminGetTemplates returns all notification templates.
// GET /admin/notifications/templates
func (h *NotificationHandler) AdminGetTemplates(c *gin.Context) {
	channel := c.Query("channel")

	templates, err := h.notifSvc.GetTemplates(channel)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, templates)
}

// AdminUpdateTemplate updates a notification template.
// PUT /admin/notifications/templates/:id
func (h *NotificationHandler) AdminUpdateTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	var req struct {
		Name      *string `json:"name"`
		Subject   *string `json:"subject"`
		Content   *string `json:"content"`
		Variables *string `json:"variables"`
		IsActive  *bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Subject != nil {
		updates["subject"] = *req.Subject
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Variables != nil {
		updates["variables"] = *req.Variables
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}

	if err := h.notifSvc.UpdateTemplate(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "template updated")
}

// AdminGetLogs returns notification logs with filters.
// GET /admin/notifications/logs
func (h *NotificationHandler) AdminGetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	channel := c.Query("channel")

	var userID *uint
	if u := c.Query("user_id"); u != "" {
		v, _ := strconv.ParseUint(u, 10, 64)
		uid := uint(v)
		userID = &uid
	}

	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}

	logs, total, err := h.notifSvc.GetLogs(page, pageSize, channel, userID, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// AdminSendBatch sends a batch notification.
// POST /admin/notifications/batch
func (h *NotificationHandler) AdminSendBatch(c *gin.Context) {
	var req struct {
		UserIDs  []uint                 `json:"user_ids" binding:"required"`
		Channel  string                 `json:"channel" binding:"required"`
		Template string                 `json:"template" binding:"required"`
		To       string                 `json:"to"`
		Subject  string                 `json:"subject"`
		Data     map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Data == nil {
		req.Data = make(map[string]interface{})
	}
	if req.To != "" {
		req.Data["to"] = req.To
	}
	if req.Subject != "" {
		req.Data["subject"] = req.Subject
	}

	success, fail, err := h.notifSvc.SendBatch(req.UserIDs, req.Channel, req.Template, req.Data)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"success": success,
		"fail":    fail,
		"total":   len(req.UserIDs),
	})
}
