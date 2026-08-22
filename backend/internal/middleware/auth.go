package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/service"
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件（含黑名单检查）
// MD规范：统一返回HTTP 200，用code字段区分错误类型
func JWTAuth(authService *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "未提供认证令牌", "data": nil})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "认证格式错误", "data": nil})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := authService.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "无效的认证令牌", "data": nil})
			c.Abort()
			return
		}

		// 检查token是否在黑名单中（登出后失效）
		hash := sha256.Sum256([]byte(tokenStr))
		tokenHash := hex.EncodeToString(hash[:])
		var count int64
		database.GetDB().Model(&model.TokenBlacklist{}).Where("token_hash = ?", tokenHash).Count(&count)
		if count > 0 {
			c.JSON(http.StatusOK, gin.H{"code": 401, "message": "令牌已失效，请重新登录", "data": nil})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)
		c.Next()
	}
}

// AdminRequired 管理员权限中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.JSON(http.StatusOK, gin.H{"code": 403, "message": "需要管理员权限", "data": nil})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ClientRequired 客户权限中间件（禁止admin token访问客户端路由，防越权）
func ClientRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if exists && isAdmin.(bool) {
			c.JSON(http.StatusOK, gin.H{"code": 403, "message": "请使用客户账号登录", "data": nil})
			c.Abort()
			return
		}
		c.Next()
	}
}
