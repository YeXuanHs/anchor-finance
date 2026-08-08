package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// GetAdmins 获取管理员列表
func (h *AdminHandler) GetAdmins(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"list": []interface{}{}, "total": 0})
}

// GetAdmin 获取单个管理员
func (h *AdminHandler) GetAdmin(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// CreateAdmin 创建管理员
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "管理员创建成功"})
}

// UpdateAdmin 更新管理员
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "管理员更新成功"})
}

// DeleteAdmin 删除管理员
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "管理员删除成功"})
}

// RegisterRoutes 注册路由
func (h *AdminHandler) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admins")
	{
		admin.GET("", h.GetAdmins)
		admin.GET("/:id", h.GetAdmin)
		admin.POST("", h.CreateAdmin)
		admin.PUT("/:id", h.UpdateAdmin)
		admin.DELETE("/:id", h.DeleteAdmin)
	}
}
