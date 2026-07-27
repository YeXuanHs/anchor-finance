package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserRemarkHandler handles user remark HTTP requests.
type UserRemarkHandler struct {
	svc *service.UserRemarkService
	log *logger.Logger
}

// NewUserRemarkHandler creates a new UserRemarkHandler.
func NewUserRemarkHandler(svc *service.UserRemarkService, log *logger.Logger) *UserRemarkHandler {
	return &UserRemarkHandler{svc: svc, log: log}
}

// Add creates a new remark on a user.
// POST /user-remarks
func (h *UserRemarkHandler) Add(c *gin.Context) {
	adminID := c.GetUint("user_id")

	var req service.AddRemarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	remark, err := h.svc.Add(adminID, req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, remark)
}

// List returns remarks for a user with pagination.
// GET /user-remarks
func (h *UserRemarkHandler) List(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		response.BadRequest(c, "user_id is required")
		return
	}
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user_id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	remarkType := c.Query("type")

	items, total, err := h.svc.GetByUser(uint(userID), remarkType, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// Delete removes a remark created by the current admin.
// DELETE /user-remarks/:id
func (h *UserRemarkHandler) Delete(c *gin.Context) {
	adminID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid remark id")
		return
	}

	if err := h.svc.Delete(uint(id), adminID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "remark deleted")
}

// AdminDelete removes any remark (super-admin).
// DELETE /admin/user-remarks/:id
func (h *UserRemarkHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid remark id")
		return
	}

	if err := h.svc.AdminDelete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "remark deleted")
}
