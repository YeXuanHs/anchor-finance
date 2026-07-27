package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserManageHandler handles admin user management HTTP requests.
type UserManageHandler struct {
	svc *service.UserManageService
	log *logger.Logger
}

// NewUserManageHandler creates a new UserManageHandler.
func NewUserManageHandler(svc *service.UserManageService, log *logger.Logger) *UserManageHandler {
	return &UserManageHandler{svc: svc, log: log}
}

// Search searches users by keyword.
// GET /manage/users
func (h *UserManageHandler) Search(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	users, total, err := h.svc.SearchUsers(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, users, total, page, pageSize)
}

// Ban disables a user account.
// POST /manage/users/:id/ban
func (h *UserManageHandler) Ban(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	if err := h.svc.Ban(uint(id), req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "user banned")
}

// Unban re-enables a user account.
// POST /manage/users/:id/unban
func (h *UserManageHandler) Unban(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.svc.Unban(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "user unbanned")
}

// AdjustBalance adds or deducts balance for a user.
// POST /manage/users/:id/balance
func (h *UserManageHandler) AdjustBalance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		Amount      float64 `json:"amount" binding:"required"`
		Description string  `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AdjustBalance(uint(id), req.Amount, req.Description); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "balance adjusted")
}

// ResetPassword sets a new password for a user.
// POST /manage/users/:id/reset-password
func (h *UserManageHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.ResetPassword(uint(id), req.NewPassword); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "password reset")
}

// GetOperationLogs returns operation logs for a user.
// GET /manage/users/:id/logs
func (h *UserManageHandler) GetOperationLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.svc.GetOperationLogs(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// GetStatus returns user status information.
// GET /manage/users/:id/status
func (h *UserManageHandler) GetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	info, err := h.svc.GetUserStatus(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, info)
}
