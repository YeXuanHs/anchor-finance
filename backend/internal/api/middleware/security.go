package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
)

// ────────────────────────────────────────────────────────────
// 登录多层限流（移植自 zjmf ClientsModel::login_get/login_inc）
// ────────────────────────────────────────────────────────────

type loginWindow struct {
	count     int
	firstTry  time.Time
}

type loginLimiter struct {
	mu       sync.Mutex
	attempts map[string]*loginWindow // key = IP
}

var loginLimit = &loginLimiter{
	attempts: make(map[string]*loginWindow),
}

// checkLoginLimit 检查登录限流（两层）
// 第一层：短时间内频繁登录 → 返回 "登录过于频繁，请稍后再试"
// 第二层：累计失败次数过多 → 返回 "登录已被锁定，请稍后再试"
func CheckLoginLimit(ip string) (bool, string) {
	loginLimit.mu.Lock()
	defer loginLimit.mu.Unlock()

	// 从数据库读取配置（秒为单位）
	blockMinutes1 := getConfigInt("login_block_minutes_1", 5)   // 第一层锁定 5 分钟
	blockMinutes2 := getConfigInt("login_block_minutes_2", 30)  // 第二层锁定 30 分钟
	maxAttempts1   := getConfigInt("login_max_attempts_1", 5)   // 第一层：5分钟内5次
	maxAttempts2   := getConfigInt("login_max_attempts_2", 15)  // 第二层：累计15次

	attempt, exists := loginLimit.attempts[ip]
	if !exists {
		return true, ""
	}

	// 第二层锁定检查
	if attempt.count >= maxAttempts2 {
		elapsed := time.Since(attempt.firstTry)
		if elapsed < time.Duration(blockMinutes2)*time.Minute {
			return false, "登录已被锁定，请稍后再试"
		}
		// 锁定过期，重置
		delete(loginLimit.attempts, ip)
		return true, ""
	}

	// 第一层锁定检查
	if attempt.count >= maxAttempts1 {
		elapsed := time.Since(attempt.firstTry)
		if elapsed < time.Duration(blockMinutes1)*time.Minute {
			return false, "登录过于频繁，请稍后再试"
		}
		// 窗口过期但未达第二层阈值，重置计数但保留累计
		attempt.count = 0
		attempt.firstTry = time.Now()
	}

	return true, ""
}

// recordLoginAttempt 记录登录失败
func RecordLoginAttempt(ip string) {
	loginLimit.mu.Lock()
	defer loginLimit.mu.Unlock()

	attempt, exists := loginLimit.attempts[ip]
	if !exists {
		attempt = &loginWindow{firstTry: time.Now()}
		loginLimit.attempts[ip] = attempt
	}
	attempt.count++
}

// clearLoginAttempts 登录成功后清除
func ClearLoginAttempts(ip string) {
	loginLimit.mu.Lock()
	defer loginLimit.mu.Unlock()
	delete(loginLimit.attempts, ip)
}

// ────────────────────────────────────────────────────────────
// 请求频率限制（移植自 zjmf ClientsModel::checkIPRequestTimes）
// ────────────────────────────────────────────────────────────

type requestRecord struct {
	timestamps []time.Time
}

type requestLimiter struct {
	mu       sync.Mutex
	requests map[string]*requestRecord // key = IP
}

var reqLimit = &requestLimiter{
	requests: make(map[string]*requestRecord),
}

// checkRequestRateLimit 检查请求频率
// 规则（参考 zjmf）：
//   - 两次请求间隔 < 1秒 → 拒绝
//   - 10秒内 > 3次 → 拒绝
//   - 1分钟内 > 30次 → 拒绝
func CheckRequestRateLimit(ip string) (bool, string) {
	reqLimit.mu.Lock()
	defer reqLimit.mu.Unlock()

	rec, exists := reqLimit.requests[ip]
	if !exists {
		rec = &requestRecord{}
		reqLimit.requests[ip] = rec
	}

	now := time.Now()

	// 清理超过 1 分钟的记录
	var valid []time.Time
	for _, t := range rec.timestamps {
		if now.Sub(t) <= time.Minute {
			valid = append(valid, t)
		}
	}
	rec.timestamps = valid

	// 检查两次请求间隔 < 1秒
	if len(rec.timestamps) > 0 {
		last := rec.timestamps[len(rec.timestamps)-1]
		if now.Sub(last) < time.Second {
			return false, "请求过于频繁，请稍后再试"
		}
	}

	// 检查 10 秒内 > 3 次
	recent10s := 0
	for _, t := range rec.timestamps {
		if now.Sub(t) <= 10*time.Second {
			recent10s++
		}
	}
	if recent10s >= 3 {
		return false, "您在10秒内已提交超过3次请求，请稍后再试"
	}

	// 检查 1 分钟内 > 30 次
	if len(rec.timestamps) >= 30 {
		return false, "您在1分钟内已提交超过30次请求，请稍后再试"
	}

	// 记录本次请求
	rec.timestamps = append(rec.timestamps, now)
	return true, ""
}

// ────────────────────────────────────────────────────────────
// 短信/邮件发送限流（移植自 zjmf sendmsglimit）
// ────────────────────────────────────────────────────────────

type sendRecord struct {
	phone   string
	sentAt  time.Time
}

type sendLimiter struct {
	mu      sync.Mutex
	sends   []sendRecord // 同 IP 的发送记录
}

var sendLimit = &sendLimiter{}

// checkSendLimit 检查发送限制
// 规则（参考 zjmf）：
//   - 同一 IP 每天最多发送 30 次
//   - 同一 IP 每天最多发送给 5 个不同的手机号
func CheckSendLimit(ip, phone string) (bool, string) {
	sendLimit.mu.Lock()
	defer sendLimit.mu.Unlock()

	today := time.Now().Truncate(24 * time.Hour)

	// 清理昨天的记录
	var valid []sendRecord
	todayPhones := make(map[string]bool)
	todayCount := 0

	for _, r := range sendLimit.sends {
		if r.sentAt.After(today) {
			valid = append(valid, r)
			todayPhones[r.phone] = true
			todayCount++
		}
	}
	sendLimit.sends = valid

	// 检查同 IP 每天发送次数
	maxDaily := getConfigInt("send_max_daily", 30)
	if todayCount >= maxDaily {
		return false, "今日发送次数已达上限"
	}

	// 检查同 IP 每天发送给不同手机号数量
	maxPhones := getConfigInt("send_max_phones", 5)
	if _, exists := todayPhones[phone]; !exists && len(todayPhones) >= maxPhones {
		return false, "今日发送手机号已达上限"
	}

	return true, ""
}

// recordSend 记录发送
func RecordSend(ip, phone string) {
	sendLimit.mu.Lock()
	defer sendLimit.mu.Unlock()
	sendLimit.sends = append(sendLimit.sends, sendRecord{
		phone:  phone,
		sentAt: time.Now(),
	})
}

// ────────────────────────────────────────────────────────────
// 密码修改后 Token 失效（移植自 zjmf client_user_update_pass_）
// ────────────────────────────────────────────────────────────

// InvalidateUserTokens 密码修改后使用户所有 Token 失效
func InvalidateUserTokens(userID uint) {
	// 在 Redis 中记录密码修改时间
	if rdb := db.GetRedis(); rdb != nil {
		ctx := context.Background()
		key := "user:pass_changed:" + uintToStr(userID)
		rdb.Set(ctx, key, time.Now().Unix(), 0)
	}
}

// IsTokenValid 检查 Token 是否在密码修改之后签发
func IsTokenValid(userID uint, tokenIssuedAt int64) bool {
	if rdb := db.GetRedis(); rdb != nil {
		ctx := context.Background()
		key := "user:pass_changed:" + uintToStr(userID)
		val, err := rdb.Get(ctx, key).Int64()
		if err == nil && val > tokenIssuedAt {
			return false // 密码在 Token 签发之后修改
		}
	}
	return true
}

// ────────────────────────────────────────────────────────────
// IP 绑定检查（移植自 zjmf home_ip_check）
// ────────────────────────────────────────────────────────────

// CheckIPBinding 检查 Token 是否与签发 IP 一致
func CheckIPBinding(tokenIP, currentIP string) bool {
	ipCheck := db.GetSystemSetting("home_ip_check")
	if ipCheck != "1" {
		return true // 未启用 IP 绑定
	}
	return tokenIP == currentIP
}

// ────────────────────────────────────────────────────────────
// 全局请求限流中间件
// ────────────────────────────────────────────────────────────

// RequestRateLimit 全局请求频率限制中间件
func RequestRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		// 跳过静态资源
		path := c.Request.URL.Path
		if len(path) > 4 && (path[:4] == "/api" || path[:5] == "/api/") {
			ok, msg := CheckRequestRateLimit(ip)
			if !ok {
				c.JSON(http.StatusTooManyRequests, gin.H{
					"code":    429,
					"message": msg,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// ────────────────────────────────────────────────────────────
// 辅助函数
// ────────────────────────────────────────────────────────────

// getConfigInt 从数据库读取整数配置，带默认值
func getConfigInt(key string, defaultVal int) int {
	val := db.GetSystemSetting(key)
	if val == "" {
		return defaultVal
	}
	// 简单的字符串转整数
	n := 0
	for _, c := range val {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		} else {
			return defaultVal
		}
	}
	if n == 0 {
		return defaultVal
	}
	return n
}

// uintToStr 简单的 uint 转字符串
func uintToStr(n uint) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}

// CleanupExpiredEntries 清理过期的限流记录
func CleanupExpiredEntries() {
	// 清理登录限流记录
	loginLimit.mu.Lock()
	for ip, attempt := range loginLimit.attempts {
		if time.Since(attempt.firstTry) > 30*time.Minute {
			delete(loginLimit.attempts, ip)
		}
	}
	loginLimit.mu.Unlock()

	// 清理请求限流记录
	reqLimit.mu.Lock()
	for ip, rec := range reqLimit.requests {
		var valid []time.Time
		for _, t := range rec.timestamps {
			if time.Since(t) <= time.Minute {
				valid = append(valid, t)
			}
		}
		if len(valid) == 0 {
			delete(reqLimit.requests, ip)
		} else {
			rec.timestamps = valid
		}
	}
	reqLimit.mu.Unlock()

	// 清理发送限流记录
	sendLimit.mu.Lock()
	today := time.Now().Truncate(24 * time.Hour)
	var valid []sendRecord
	for _, r := range sendLimit.sends {
		if r.sentAt.After(today) {
			valid = append(valid, r)
		}
	}
	sendLimit.sends = valid
	sendLimit.mu.Unlock()
}
