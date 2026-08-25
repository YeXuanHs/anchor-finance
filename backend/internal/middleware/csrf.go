package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CSRF CSRF防护中间件（双Token方案）
// 前端在登录时从响应头获取csrf_token，后续POST/PUT/DELETE请求必须携带
func CSRF() gin.HandlerFunc {
	return func(c *gin.Context) {
		// GET/HEAD/OPTIONS不校验
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// 从cookie获取csrf_token
		cookieToken, err := c.Cookie("csrf_token")
		if err != nil || cookieToken == "" {
			// 首次请求，设置csrf cookie
			token := generateCSRFToken()
			c.SetCookie("csrf_token", token, 3600*24, "/", "", false, false)
			c.Header("X-CSRF-Token", token)
			// 仅login路由首次POST不拦截
			path := c.Request.URL.Path
			if path == "/api/admin/login" || path == "/api/client/login" || path == "/api/client/register" || path == "/zjmf_api_login" {
				c.Next()
				return
			}
			// 其他路由首次POST也拦截
			c.JSON(http.StatusOK, gin.H{
				"code":    403,
				"message": "CSRF token验证失败",
				"data":    nil,
			})
			c.Abort()
			return
		}

		// 校验header中的csrf_token和cookie中的一致
		headerToken := c.GetHeader("X-CSRF-Token")
		if headerToken == "" || headerToken != cookieToken {
			c.JSON(http.StatusOK, gin.H{
				"code":    403,
				"message": "CSRF token验证失败",
				"data":    nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func generateCSRFToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
