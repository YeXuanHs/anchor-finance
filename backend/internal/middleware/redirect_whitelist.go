package middleware

import (
	"net/http"
	"net/url"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// RedirectWhitelist 开放重定向白名单中间件（MD 9.1 漏洞6）
// 检查请求中的redirect/url参数是否在白名单内
func RedirectWhitelist() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查常见的重定向参数
		redirectParams := []string{"redirect", "redirect_url", "return_url", "url", "callback"}
		for _, param := range redirectParams {
			redirectURL := c.Query(param)
			if redirectURL == "" {
				continue
			}

			// 解析URL
			parsed, err := url.Parse(redirectURL)
			if err != nil {
				continue
			}

			// 空host或相对路径（站内）放行
			if parsed.Host == "" {
				continue
			}

			// 检查是否在白名单内
			if !isAllowedRedirectHost(parsed.Host) {
				c.JSON(http.StatusOK, gin.H{"code": 400, "message": "不允许跳转到外部域名", "data": nil})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// isAllowedRedirectHost 检查host是否在白名单内
func isAllowedRedirectHost(host string) bool {
	db := database.GetDB()
	var setting model.Setting
	if err := db.Where("`key` = ?", "redirect_whitelist").First(&setting).Error; err != nil {
		// 没有配置白名单，默认只允许站内
		return false
	}

	// 白名单格式：域名用逗号分隔，如 "example.com,cdn.example.com"
	allowedHosts := []string{}
	for _, h := range splitAndTrim(setting.Value) {
		allowedHosts = append(allowedHosts, h)
	}

	for _, allowed := range allowedHosts {
		if host == allowed {
			return true
		}
		// 支持通配符：*.example.com
		if len(allowed) > 2 && allowed[:2] == "*." {
			suffix := allowed[1:] // .example.com
			if len(host) > len(suffix) && host[len(host)-len(suffix):] == suffix {
				return true
			}
		}
	}
	return false
}

// splitAndTrim 按逗号分割并去除空白
func splitAndTrim(s string) []string {
	var result []string
	current := ""
	for _, c := range s {
		if c == ',' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
