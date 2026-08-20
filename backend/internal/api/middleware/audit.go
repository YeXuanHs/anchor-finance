package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/db"
	"anchorfinance/pkg/logger"

	"github.com/gin-gonic/gin"
)

// sensitiveFields 敏感字段列表（需要脱敏）
var sensitiveFields = []string{
	"password", "old_password", "new_password", "confirm_password",
	"token", "access_token", "refresh_token", "secret",
	"credit_card", "card_number", "cvv", "pin",
}

// AuditLogMiddleware 审计日志中间件
func AuditLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过不需要记录的路径
		path := c.Request.URL.Path
		if shouldSkipAudit(path) {
			c.Next()
			return
		}

		start := time.Now()

		// 读取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(start).Milliseconds()

		// 获取用户信息
		userID, _ := c.Get(ContextKeyUserID)
		isAdmin, _ := c.Get(ContextKeyIsAdmin)

		userType := "client"
		if isAdmin != nil && isAdmin.(bool) {
			userType = "admin"
		}

		// 脱敏请求数据
		sanitizedData := sanitizeRequestData(string(requestBody))

		// 确定操作类型
		action := determineAction(c.Request.Method, path)

		// 创建审计日志
		auditLog := &model.AuditLog{
			UserID:       getUserIDUint(userID),
			Username:     getUsername(c),
			UserType:     userType,
			Action:       action,
			Description:  generateDescription(c.Request.Method, path),
			Module:       getModule(path),
			Controller:   getController(path),
			Method:       c.Request.Method,
			IP:           c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			RequestData:  sanitizedData,
			ResponseCode: c.Writer.Status(),
			Duration:     duration,
			Status:       getStatus(c.Writer.Status()),
		}

		// 异步写入日志（不影响响应速度）
		go func() {
			if dbConn := db.GetDB(); dbConn != nil {
				log := logger.Default()
				if err := dbConn.Create(auditLog).Error; err != nil {
					log.Errorf("写入审计日志失败: %v", err)
				}
			}
		}()
	}
}

// shouldSkipAudit 判断是否跳过审计
func shouldSkipAudit(path string) bool {
	skipPaths := []string{
		"/api/v1/captcha",
		"/api/v1/geetest",
		"/api/public",
		"/health",
		"/metrics",
		"/favicon.ico",
	}

	for _, skipPath := range skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}

	// 跳过静态资源
	if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") ||
		strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".jpg") ||
		strings.HasSuffix(path, ".ico") || strings.HasSuffix(path, ".svg") {
		return true
	}

	return false
}

// sanitizeRequestData 脱敏请求数据
func sanitizeRequestData(data string) string {
	if data == "" {
		return ""
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(data), &jsonData); err != nil {
		// 非 JSON 格式，截断过长内容
		if len(data) > 1000 {
			return data[:1000] + "..."
		}
		return data
	}

	// 脱敏敏感字段
	for _, field := range sensitiveFields {
		if _, exists := jsonData[field]; exists {
			jsonData[field] = "***"
		}
	}

	sanitized, _ := json.Marshal(jsonData)
	if len(sanitized) > 1000 {
		return string(sanitized[:1000]) + "..."
	}
	return string(sanitized)
}

// determineAction 确定操作类型
func determineAction(method, path string) string {
	switch method {
	case "GET":
		if strings.Contains(path, "/list") || strings.Contains(path, "page=") {
			return "查询列表"
		}
		return "查询详情"
	case "POST":
		if strings.Contains(path, "/login") {
			return "登录"
		}
		if strings.Contains(path, "/register") {
			return "注册"
		}
		if strings.Contains(path, "/upload") {
			return "上传文件"
		}
		return "创建"
	case "PUT", "PATCH":
		return "更新"
	case "DELETE":
		return "删除"
	default:
		return method
	}
}

// generateDescription 生成操作描述
func generateDescription(method, path string) string {
	return method + " " + path
}

// getModule 获取模块名
func getModule(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 3 {
		return parts[2] // /api/v1/xxx -> v1
	}
	return ""
}

// getController 获取控制器名
func getController(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 4 {
		return parts[3]
	}
	return ""
}

// getUserIDUint 获取用户 ID
func getUserIDUint(userID interface{}) uint {
	if userID == nil {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	return 0
}

// getUsername 获取用户名
func getUsername(c *gin.Context) string {
	if username, exists := c.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name
		}
	}
	return ""
}

// getStatus 获取状态
func getStatus(code int) string {
	if code >= 200 && code < 300 {
		return "success"
	}
	return "failed"
}
