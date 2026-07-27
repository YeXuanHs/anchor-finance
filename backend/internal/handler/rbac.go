package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type RbacHandler struct {
	rbacSvc *service.RbacService
	log     *logger.Logger
}

func NewRbacHandler(rbacSvc *service.RbacService, log *logger.Logger) *RbacHandler {
	return &RbacHandler{rbacSvc: rbacSvc, log: log}
}

// GetRoles returns all roles.
// GET /admin/roles
func (h *RbacHandler) GetRoles(c *gin.Context) {
	roles, err := h.rbacSvc.GetRoles()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, roles)
}

// CreateRole creates a new role.
// POST /admin/roles
func (h *RbacHandler) CreateRole(c *gin.Context) {
	var req service.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	role, err := h.rbacSvc.CreateRole(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, role)
}

// UpdateRole updates an existing role.
// PUT /admin/roles/:id
func (h *RbacHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid role id")
		return
	}

	var req service.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	role, err := h.rbacSvc.UpdateRole(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, role)
}

// DeleteRole deletes a role.
// DELETE /admin/roles/:id
func (h *RbacHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid role id")
		return
	}

	if err := h.rbacSvc.DeleteRole(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "role deleted")
}

// GetPermissions returns all permissions grouped by module.
// GET /admin/permissions
func (h *RbacHandler) GetPermissions(c *gin.Context) {
	perms, err := h.rbacSvc.GetPermissions()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, perms)
}

// AssignRole assigns roles to a user.
// POST /admin/users/:id/roles
func (h *RbacHandler) AssignRole(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req service.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.rbacSvc.AssignRole(uint(userID), req.RoleIDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "roles assigned")
}

// GetUserRoles returns roles for a user.
// GET /admin/users/:id/roles
func (h *RbacHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	roles, err := h.rbacSvc.GetUserRoles(uint(userID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, roles)
}
