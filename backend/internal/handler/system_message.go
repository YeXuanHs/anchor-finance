package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// SystemMessageHandler handles system message HTTP requests.
type SystemMessageHandler struct {
	messageSvc *service.SystemMessageService
	log        *logger.Logger
}

// NewSystemMessageHandler creates a new SystemMessageHandler.
func NewSystemMessageHandler(messageSvc *service.SystemMessageService, log *logger.Logger) *SystemMessageHandler {
	return &SystemMessageHandler{messageSvc: messageSvc, log: log}
}

// GetList returns paginated system messages for the authenticated user.
// GET /messages
func (h *SystemMessageHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	msgType := c.Query("type")
	onlyUnread := c.Query("unread") == "1"

	messages, total, err := h.messageSvc.GetList(userID, page, pageSize, msgType, onlyUnread)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, messages, total, page, pageSize)
}

// GetDetail returns a single message.
// GET /messages/:id
func (h *SystemMessageHandler) GetDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid message id")
		return
	}

	msg, err := h.messageSvc.GetByID(userID, uint(id))
	if err != nil {
		response.NotFound(c, "message not found")
		return
	}
	response.Success(c, msg)
}

// MarkRead marks a single message as read.
// POST /messages/:id/read
func (h *SystemMessageHandler) MarkRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid message id")
		return
	}

	if err := h.messageSvc.MarkRead(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "message marked as read")
}

// MarkAllRead marks all messages as read.
// POST /messages/read-all
func (h *SystemMessageHandler) MarkAllRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, err := h.messageSvc.MarkAllRead(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"marked_count": count})
}

// Delete deletes a single message.
// DELETE /messages/:id
func (h *SystemMessageHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid message id")
		return
	}

	if err := h.messageSvc.Delete(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "message deleted")
}

// DeleteAll deletes all messages for the user.
// DELETE /messages
func (h *SystemMessageHandler) DeleteAll(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, err := h.messageSvc.DeleteAll(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"deleted_count": count})
}

// GetUnreadCount returns the number of unread messages.
// GET /messages/unread-count
func (h *SystemMessageHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, err := h.messageSvc.GetUnreadCount(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"unread_count": count})
}

// GetTypes returns distinct message types for the user.
// GET /messages/types
func (h *SystemMessageHandler) GetTypes(c *gin.Context) {
	userID := c.GetUint("user_id")

	types, err := h.messageSvc.GetTypes(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"types": types})
}
