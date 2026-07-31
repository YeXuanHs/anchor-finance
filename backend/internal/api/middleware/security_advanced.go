package middleware

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"anchorfinance/internal/security"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// SecurityConfig 安全中间件配置
type SecurityConfig struct {
	EnableXSSFilter   bool     // 启用 XSS 过滤
	EnableSQLFilter   bool     // 启用 SQL 注入过滤
	BlockedIPs        []string // 黑名单 IP
	AllowedMethods    []string // 允许的 HTTP 方法
	MaxRequestSize    int64    // 最大请求体大小（字节）
}

// DefaultSecurityConfig 默认安全配置
func DefaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		EnableXSSFilter:   true,
		EnableSQLFilter:   true,
		MaxRequestSize:    10 << 20, // 10MB
		AllowedMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	}
}

// Security 安全中间件
func Security(config *SecurityConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultSecurityConfig()
	}

	return func(c *gin.Context) {
		// 1. 检查 IP 黑名单
		if isIPBlocked(c.ClientIP(), config.BlockedIPs) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "访问被拒绝",
			})
			return
		}

		// 2. 检查 HTTP 方法
		if !isMethodAllowed(c.Request.Method, config.AllowedMethods) {
			c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{
				"code":    405,
				"message": "不支持的请求方法",
			})
			return
		}

		// 3. 检查请求体大小
		if c.Request.ContentLength > config.MaxRequestSize {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"code":    413,
				"message": "请求体过大",
			})
			return
		}

		// 4. XSS 防护：检查查询参数
		if config.EnableXSSFilter {
			if detectXSS(c.Request.URL.RawQuery) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "检测到潜在的 XSS 攻击",
				})
				return
			}
		}

		// 5. SQL 注入防护：检查查询参数
		if config.EnableSQLFilter {
			if security.DetectSQLInjection(c.Request.URL.RawQuery) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"code":    400,
					"message": "检测到潜在的 SQL 注入攻击",
				})
				return
			}

			// 检查表单参数
			if c.Request.Method == "POST" || c.Request.Method == "PUT" {
				if err := c.Request.ParseForm(); err == nil {
					for _, values := range c.Request.Form {
						for _, v := range values {
							if security.DetectSQLInjection(v) {
								c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
									"code":    400,
									"message": "检测到潜在的 SQL 注入攻击",
								})
								return
							}
						}
					}
				}
			}
		}

		// 6. 添加安全头
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}

// SecurityHeaders 安全响应头中间件
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self' 'unsafe-inline' 'unsafe-eval' https: http:;")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}

// RequestBodySanitizer 请求体清理中间件
func RequestBodySanitizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只对 POST/PUT/PATCH 请求处理
		if c.Request.Method != "POST" && c.Request.Method != "PUT" && c.Request.Method != "PATCH" {
			c.Next()
			return
		}

		// 读取请求体
		if c.Request.Body == nil {
			c.Next()
			return
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}

		// 检查请求体是否包含 XSS 攻击
		bodyStr := string(bodyBytes)
		if detectXSS(bodyStr) {
			response.BadRequest(c, "检测到潜在的安全威胁")
			c.Abort()
			return
		}

		// 恢复请求体
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		c.Next()
	}
}

// detectXSS 检测 XSS 攻击
func detectXSS(input string) bool {
	input = strings.ToLower(input)
	xssPatterns := []string{
		"<script",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
		"onclick=",
		"onmouseover=",
		"onfocus=",
		"onblur=",
		"onsubmit=",
		"<iframe",
		"<object",
		"<embed",
		"<applet",
	}

	for _, pattern := range xssPatterns {
		if strings.Contains(input, pattern) {
			return true
		}
	}
	return false
}

// isIPBlocked 检查 IP 是否在黑名单中
func isIPBlocked(ip string, blockedIPs []string) bool {
	for _, blockedIP := range blockedIPs {
		if ip == blockedIP {
			return true
		}
	}
	return false
}

// isMethodAllowed 检查 HTTP 方法是否允许
func isMethodAllowed(method string, allowedMethods []string) bool {
	for _, allowed := range allowedMethods {
		if method == allowed {
			return true
		}
	}
	return false
}
