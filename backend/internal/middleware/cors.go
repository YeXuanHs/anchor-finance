package middleware

import (
	"net/http"
	"os"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// CORS 跨域中间件（优先从settings表读取，回退到环境变量）
func CORS() gin.HandlerFunc {
	// 优先从数据库settings表读取 cors_allowed_origin
	allowedOrigin := ""
	db := database.GetDB()
	if db != nil {
		var setting model.Setting
		if err := db.Where("`key` = ?", "cors_allowed_origin").First(&setting).Error; err == nil && setting.Value != "" {
			allowedOrigin = setting.Value
		}
	}
	// 回退到环境变量
	if allowedOrigin == "" {
		allowedOrigin = os.Getenv("CORS_ORIGIN")
	}
	// 最终默认值
	if allowedOrigin == "" {
		allowedOrigin = "*" // 开发环境默认允许所有
	}

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", allowedOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.JSON(http.StatusNoContent, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
