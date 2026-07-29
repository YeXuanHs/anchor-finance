package middleware

import (
	"net/http"
	"strings"

	"anchorfinance/pkg/auth"
	"github.com/gin-gonic/gin"
)

const (
	// ContextKeyUserID is the gin context key for the authenticated user ID.
	ContextKeyUserID = "user_id"
	// ContextKeyIsAdmin is the gin context key for the admin flag.
	ContextKeyIsAdmin = "is_admin"
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
func JWTAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok":      false,
				"message": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok":      false,
				"message": "invalid authorization format, expected: Bearer <token>",
			})
			return
		}

		claims, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"ok":      false,
				"message": "invalid or expired token",
			})
			return
		}

		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyIsAdmin, claims.IsAdmin)
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
				"message": "admin access required",
			})
			return
		}
		c.Next()
	}
}
