package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ClientHandler struct{}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{}
}

// GetClients 获取客户列表
func (h *ClientHandler) GetClients(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	// TODO: 实现客户列表查询
	c.JSON(http.StatusOK, gin.H{
		"list":      []interface{}{},
		"total":     0,
		"page":      page,
		"page_size": pageSize,
		"keyword":   keyword,
		"status":    status,
	})
}

// GetClient 获取单个客户
func (h *ClientHandler) GetClient(c *gin.Context) {
	id := c.Param("id")
	// TODO: 实现客户详情查询
	c.JSON(http.StatusOK, gin.H{"id": id, "username": "test", "email": "test@example.com"})
}

// CreateClient 创建客户
func (h *ClientHandler) CreateClient(c *gin.Context) {
	// TODO: 实现客户创建
	c.JSON(http.StatusCreated, gin.H{"message": "客户创建成功"})
}

// UpdateClient 更新客户
func (h *ClientHandler) UpdateClient(c *gin.Context) {
	// TODO: 实现客户更新
	c.JSON(http.StatusOK, gin.H{"message": "客户更新成功"})
}

// DeleteClient 删除客户
func (h *ClientHandler) DeleteClient(c *gin.Context) {
	// TODO: 实现客户删除
	c.JSON(http.StatusOK, gin.H{"message": "客户删除成功"})
}

// RegisterRoutes 注册路由
func (h *ClientHandler) RegisterRoutes(r *gin.RouterGroup) {
	client := r.Group("/clients")
	{
		client.GET("", h.GetClients)
		client.GET("/:id", h.GetClient)
		client.POST("", h.CreateClient)
		client.PUT("/:id", h.UpdateClient)
		client.DELETE("/:id", h.DeleteClient)
	}
}
