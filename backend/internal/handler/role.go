package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct{}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{}
}

// GetRoles 获取角色列表
func (h *RoleHandler) GetRoles(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"list": []interface{}{}, "total": 0})
}

// CreateRole 创建角色
func (h *RoleHandler) CreateRole(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "角色创建成功"})
}

// UpdateRole 更新角色
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "角色更新成功"})
}

// DeleteRole 删除角色
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "角色删除成功"})
}

// RegisterRoutes 注册路由
func (h *RoleHandler) RegisterRoutes(r *gin.RouterGroup) {
	role := r.Group("/roles")
	{
		role.GET("", h.GetRoles)
		role.POST("", h.CreateRole)
		role.PUT("/:id", h.UpdateRole)
		role.DELETE("/:id", h.DeleteRole)
	}
}
