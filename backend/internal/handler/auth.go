package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login 管理员登录
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: 实现登录逻辑
	c.JSON(http.StatusOK, gin.H{
		"token": "test-token",
		"user": gin.H{
			"id":       1,
			"username": req.Username,
			"email":    "admin@example.com",
		},
	})
}

// Logout 管理员登出
func Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}
