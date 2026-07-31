package middleware

import (
	"net/http"
	"strings"

	"anchorfinance/pkg/auth"
	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
)

const (
	// ContextKeyUserID is the gin context key for the authenticated user ID.
	ContextKeyUserID = "user_id"
	// ContextKeyIsAdmin is the gin context key for the admin flag.
	ContextKeyIsAdmin = "is_admin"
	// ContextKeyTokenIssuedAt is the gin context key for the token issued at time.
	ContextKeyTokenIssuedAt = "token_issued_at"
	// ContextKeyTokenIP is the gin context key for the token IP.
	ContextKeyTokenIP = "token_ip"
)

var defaultJWTManager *auth.JWTManager

// Init sets the default JWT manager used by AuthRequired.
func Init(jwtMgr *auth.JWTManager) {
	defaultJWTManager = jwtMgr
}

// AuthRequired returns a middleware that validates Bearer tokens using the default JWT manager.
func AuthRequired() gin.HandlerFunc {
	return JWTAuth(defaultJWTManager)
}

// JWTAuth returns a middleware that validates Bearer tokens and sets user info in context.
// 移植自 zjmf 的安全逻辑：
//   - Token 有效性验证
//   - 密码修改后 Token 失效（client_user_update_pass_）
//   - IP 绑定检查（home_ip_check）
func JWTAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok":      false,
				"message": "请登录后再试",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok":      false,
				"message": "授权格式错误，应为: Bearer <token>",
			})
			return
		}

		claims, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok":      false,
				"message": "Token 无效或已过期，请重新登录",
			})
			return
		}

		// 检查 Token 是否在密码修改之后签发（移植自 zjmf client_user_update_pass_）
		if !IsTokenValid(claims.UserID, claims.IssuedAt) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok":      false,
				"message": "密码已修改，请重新登录",
			})
			return
		}

		// 检查 IP 绑定（移植自 zjmf home_ip_check）
		if claims.IP != "" && !CheckIPBinding(claims.IP, c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok":      false,
				"message": "登录已失效，请重新登录",
			})
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyIsAdmin, claims.IsAdmin)
		c.Set(ContextKeyTokenIssuedAt, claims.IssuedAt)
		c.Set(ContextKeyTokenIP, claims.IP)
		c.Next()
	}
}

// AdminRequired returns a middleware that requires the authenticated user to be an admin.
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get(ContextKeyIsAdmin)
		if !exists || !isAdmin.(bool) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"ok":      false,
				"message": "需要管理员权限",
			})
			return
		}
		c.Next()
	}
}
