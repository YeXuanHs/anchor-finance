package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// RateLimiter 简单内存限流器（按IP）
type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int           // 窗口内允许请求数
	window   time.Duration // 时间窗口
}

type visitor struct {
	count    int
	lastSeen time.Time
}

// NewRateLimiter 创建限流器
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		window:   window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow 检查请求是否允许
func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[key]
	if !exists {
		rl.visitors[key] = &visitor{count: 1, lastSeen: time.Now()}
		return true
	}

	if time.Since(v.lastSeen) > rl.window {
		v.count = 1
		v.lastSeen = time.Now()
		return true
	}

	if v.count >= rl.rate {
		return false
	}

	v.count++
	v.lastSeen = time.Now()
	return true
}

// getRateLimit 从settings表读取限流配置（防抖频率）
func getRateLimit() (int, int) {
	rate := 3      // 默认每分钟3次
	window := 60   // 默认60秒
	db := database.GetDB()
	if db == nil {
		return rate, window
	}
	var setting model.Setting
	if err := db.Where("`key` = ?", "api_rate_limit").First(&setting).Error; err == nil {
		var v int
		if _, err := fmt.Sscanf(setting.Value, "%d", &v); err == nil && v > 0 {
			rate = v
		}
	}
	if err := db.Where("`key` = ?", "api_rate_window").First(&setting).Error; err == nil {
		var v int
		if _, err := fmt.Sscanf(setting.Value, "%d", &v); err == nil && v > 0 {
			window = v
		}
	}
	return rate, window
}

// RateLimit 限流中间件工厂（从settings表读取配置）
func RateLimit(rate int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, window)
	return func(c *gin.Context) {
		// M6修复：用RemoteIP()取TCP连接真实IP，不读X-Forwarded-For等可伪造头
		key := c.RemoteIP()
		if !limiter.Allow(key) {
			c.JSON(http.StatusOK, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
				"data":    nil,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimitFromSettings 限流中间件工厂（从settings表读取配置，无需传参）
func RateLimitFromSettings() gin.HandlerFunc {
	rate, window := getRateLimit()
	return RateLimit(rate, time.Duration(window)*time.Second)
}
