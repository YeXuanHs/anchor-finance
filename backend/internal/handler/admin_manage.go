package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// AdminManageHandler handles admin management HTTP requests.
type AdminManageHandler struct {
	svc *service.AdminManageService
	log *logger.Logger
}

// NewAdminManageHandler creates a new AdminManageHandler.
func NewAdminManageHandler(svc *service.AdminManageService, log *logger.Logger) *AdminManageHandler {
	return &AdminManageHandler{svc: svc, log: log}
}

// ==================== Admin CRUD ====================

// List returns a paginated list of admins.
// GET /admins
func (h *AdminManageHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	admins, total, err := h.svc.GetList(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, admins, total, page, pageSize)
}

// Get returns a single admin by ID.
// GET /admins/:id
func (h *AdminManageHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid admin id")
		return
	}

	admin, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, admin)
}

// Create creates a new admin.
// POST /admins
func (h *AdminManageHandler) Create(c *gin.Context) {
	var req service.CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	admin, err := h.svc.Create(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Log the operation
	adminID := c.GetUint("user_id")
	h.svc.CreateOperationLog(&model.AdminLog{
		AdminID:    adminID,
		Action:     "create",
		Module:     "admin",
		TargetID:   admin.ID,
		TargetType: "admin",
		Detail:     "Created admin: " + req.Username,
		IP:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
	})

	response.Success(c, admin)
}

// Update modifies an existing admin.
// PUT /admins/:id
func (h *AdminManageHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid admin id")
		return
	}

	var req service.UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Update(uint(id), req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Log the operation
	adminID := c.GetUint("user_id")
	h.svc.CreateOperationLog(&model.AdminLog{
		AdminID:    adminID,
		Action:     "update",
		Module:     "admin",
		TargetID:   uint(id),
		TargetType: "admin",
		Detail:     "Updated admin profile",
		IP:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
	})

	response.SuccessMsg(c, "admin updated")
}

// Delete removes an admin.
// DELETE /admins/:id
func (h *AdminManageHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid admin id")
		return
	}

	// Prevent self-deletion
	currentAdminID := c.GetUint("user_id")
	if uint(id) == currentAdminID {
		response.BadRequest(c, "cannot delete yourself")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Log the operation
	h.svc.CreateOperationLog(&model.AdminLog{
		AdminID:    currentAdminID,
		Action:     "delete",
		Module:     "admin",
		TargetID:   uint(id),
		TargetType: "admin",
		Detail:     "Deleted admin",
		IP:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
	})

	response.SuccessMsg(c, "admin deleted")
}

// ==================== Status Management ====================

// SetStatus enables or disables an admin.
// POST /admins/:id/status
func (h *AdminManageHandler) SetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid admin id")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Prevent self-disable
	currentAdminID := c.GetUint("user_id")
	if uint(id) == currentAdminID && req.Status == 0 {
		response.BadRequest(c, "cannot disable yourself")
		return
	}

	if err := h.svc.SetStatus(uint(id), req.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Log the operation
	action := "enable"
	if req.Status == 0 {
		action = "disable"
	}
	h.svc.CreateOperationLog(&model.AdminLog{
		AdminID:    currentAdminID,
		Action:     action,
		Module:     "admin",
		TargetID:   uint(id),
		TargetType: "admin",
		Detail:     action + " admin",
		IP:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
	})

	response.SuccessMsg(c, "admin status updated")
}

// ==================== Password Management ====================

// ResetPassword sets a new password for an admin.
// POST /admins/:id/reset-password
func (h *AdminManageHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid admin id")
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

	// Log the operation
	adminID := c.GetUint("user_id")
	h.svc.CreateOperationLog(&model.AdminLog{
		AdminID:    adminID,
		Action:     "reset_password",
		Module:     "admin",
		TargetID:   uint(id),
		TargetType: "admin",
		Detail:     "Reset admin password",
		IP:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
	})

	response.SuccessMsg(c, "password reset successfully")
}

// ==================== Operation Logs ====================

// GetOperationLogs returns admin operation logs.
// GET /admins/operation-logs
func (h *AdminManageHandler) GetOperationLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var adminID uint
	if uid := c.Query("admin_id"); uid != "" {
		v, _ := strconv.ParseUint(uid, 10, 64)
		adminID = uint(v)
	}
	module := c.Query("module")

	logs, total, err := h.svc.GetOperationLogs(page, pageSize, adminID, module)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}
