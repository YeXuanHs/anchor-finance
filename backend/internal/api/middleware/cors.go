package middleware

import (
	"net/http"
	"strings"

	"anchorfinance/pkg/db"
	"github.com/gin-gonic/gin"
)

// CORS returns a middleware that restricts origins based on configuration.
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := isOriginAllowed(origin)

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Requested-With")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate")
		c.Header("Pragma", "no-cache")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isOriginAllowed checks if the origin is in the allowed list.
// Allowed origins come from DB setting "cors_allowed_origins" (comma-separated),
// plus the site_url. Empty means same-origin only.
func isOriginAllowed(origin string) bool {
	if origin == "" {
		return true // same-origin requests have no Origin header
	}

	siteURL := db.GetSystemSetting("site_url")
	extraOrigins := db.GetSystemSetting("cors_allowed_origins")

	allowed := []string{}
	if siteURL != "" {
		allowed = append(allowed, strings.TrimRight(siteURL, "/"))
	}
	if extraOrigins != "" {
		for _, o := range strings.Split(extraOrigins, ",") {
			o = strings.TrimSpace(o)
			if o != "" {
				allowed = append(allowed, strings.TrimRight(o, "/"))
			}
		}
	}

	// 如果没有配置任何来源，拒绝跨域
	if len(allowed) == 0 {
		return false
	}

	for _, a := range allowed {
		if strings.EqualFold(a, origin) {
			return true
		}
	}
	return false
}
