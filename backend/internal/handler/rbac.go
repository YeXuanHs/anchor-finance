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
	response.SuccessMsg(c, "删除成功")
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
	response.SuccessMsg(c, "请求成功")
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

// ─── RBAC Admin Methods (from zjmf RbacController) ───

// Index returns all roles with admin users.
func (h *RbacHandler) Index(c *gin.Context) {
	order := c.DefaultQuery("order", "a.id")
	sort := c.DefaultQuery("sort", "DESC")

	roles, err := h.rbacSvc.AdminGetRoles(order, sort)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"roles": roles})
}

// AddRolePage returns auth rules tree for the add role page.
func (h *RbacHandler) AddRolePage(c *gin.Context) {
	auths, err := h.rbacSvc.AdminGetAuthTree()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"auths": auths})
}

// AddRole creates a new role with auth rules.
func (h *RbacHandler) AddRole(c *gin.Context) {
	var req struct {
		Name   string `json:"name" binding:"required,max=50"`
		Remark string `json:"remark"`
		Status int    `json:"status"`
		Auth   []uint `json:"auth"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.rbacSvc.AdminAddRole(req.Name, req.Remark, req.Status, req.Auth); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "请求成功")
}

// EditRolePage returns data for editing a role.
func (h *RbacHandler) EditRolePage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	data, err := h.rbacSvc.AdminEditRolePage(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

// EditRole updates a role with auth rules.
func (h *RbacHandler) EditRole(c *gin.Context) {
	var req struct {
		ID     uint   `json:"id" binding:"required"`
		Name   string `json:"name" binding:"required,max=50"`
		Remark string `json:"remark"`
		Status int    `json:"status"`
		Auth   []uint `json:"auth"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.rbacSvc.AdminEditRole(req.ID, req.Name, req.Remark, req.Status, req.Auth); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "请求成功")
}

// Delete deletes a role (system roles excluded).
func (h *RbacHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		id, _ = strconv.ParseUint(c.Query("id"), 10, 64)
	}
	if id == 0 {
		response.BadRequest(c, "invalid role id")
		return
	}

	if err := h.rbacSvc.AdminDeleteRole(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// CopyRole duplicates a role.
func (h *RbacHandler) CopyRole(c *gin.Context) {
	var req struct {
		RoleID    uint   `json:"role_id" binding:"required"`
		RoleName  string `json:"role_name" binding:"required"`
		RoleRemark string `json:"role_remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.rbacSvc.AdminCopyRole(req.RoleID, req.RoleName, req.RoleRemark); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "请求成功")
}
