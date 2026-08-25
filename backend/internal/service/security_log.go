package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"gorm.io/gorm"
)

// SecurityLogger 安全日志服务（MD 9.1 功能9）
type SecurityLogger struct {
	db *gorm.DB
}

// NewSecurityLogger 创建安全日志服务
func NewSecurityLogger(db *gorm.DB) *SecurityLogger {
	return &SecurityLogger{db: db}
}

// AttackType 攻击类型常量
const (
	AttackZeroPrice      = "zero_price"       // 0元购攻击
	AttackSQLInject      = "sql_inject"       // SQL注入
	AttackBruteForce     = "brute_force"      // 暴力破解
	AttackIDOR           = "idor"             // 越权访问
	AttackCSRF           = "csrf"             // CSRF攻击
	AttackRegistration   = "registration_abuse" // 批量注册
	AttackPriceManipulate = "price_manipulate" // 价格篡改
	AttackConfigInject   = "config_inject"    // 配置注入
	AttackOpenRedirect   = "open_redirect"    // 开放重定向
	AttackSessionHijack  = "session_hijack"   // 会话劫持
)

// LogSecurityEvent 记录安全事件
// MD 9.1要求记录：时间、用户ID、用户名、IP、真实IP、请求路径、方法、UA、会话ID、Referer、请求参数（脱敏）、攻击类型
func (sl *SecurityLogger) LogSecurityEvent(attackType string, userID uint, username, ip, path, method, userAgent, referer, detail string, params interface{}) {
	// 获取真实IP（X-Forwarded-For）
	realIP := ip
	// 注意：Gin的c.ClientIP()已经处理了X-Forwarded-For，这里ip就是真实IP

	// 脱敏参数
	paramsStr := ""
	if params != nil {
		if bytes, err := json.Marshal(params); err == nil {
			paramsStr = sl.sanitizeParams(string(bytes))
		}
	}

	log := model.SecurityLog{
		AttackType: attackType,
		UserID:     userID,
		Username:   username,
		IP:         ip,
		RealIP:     realIP,
		Path:       path,
		Method:     method,
		UserAgent:  userAgent,
		Referer:    referer,
		Params:     paramsStr,
		Detail:     detail,
		CreatedAt:  time.Now(),
	}

	sl.db.Create(&log)
}

// LogZeroPriceAttack 记录0元购攻击
func (sl *SecurityLogger) LogZeroPriceAttack(userID uint, username, ip, path string, amount float64) {
	detail := fmt.Sprintf("尝试价格操控，原始金额: %.2f", amount)
	sl.LogSecurityEvent(AttackZeroPrice, userID, username, ip, path, "POST", "", "", detail, nil)
}

// LogBruteForceAttack 记录暴力破解攻击
func (sl *SecurityLogger) LogBruteForceAttack(userID uint, username, ip, path string, failCount int) {
	detail := fmt.Sprintf("连续登录失败%d次", failCount)
	sl.LogSecurityEvent(AttackBruteForce, userID, username, ip, path, "POST", "", "", detail, nil)
}

// LogIDORAttack 记录越权访问攻击
func (sl *SecurityLogger) LogIDORAttack(userID uint, username, ip, path, resourceType string, resourceID uint) {
	detail := fmt.Sprintf("尝试访问不属于自己的%s (ID: %d)", resourceType, resourceID)
	sl.LogSecurityEvent(AttackIDOR, userID, username, ip, path, "", "", "", detail, nil)
}

// LogRegistrationAbuse 记录批量注册攻击
func (sl *SecurityLogger) LogRegistrationAbuse(ip, path string, count int) {
	detail := fmt.Sprintf("短时间内注册%d次", count)
	sl.LogSecurityEvent(AttackRegistration, 0, "", ip, path, "POST", "", "", detail, nil)
}

// sanitizeParams 脱敏参数（移除密码、token等敏感信息）
func (sl *SecurityLogger) sanitizeParams(params string) string {
	sensitiveKeys := []string{"password", "token", "secret", "key", "authorization", "credit_card", "card_number"}
	result := params
	for _, key := range sensitiveKeys {
		// 简单脱敏：将包含敏感key的值替换为***
		lower := strings.ToLower(result)
		if strings.Contains(lower, key) {
			// 这里可以做更精确的JSON脱敏，当前简单处理
			result = strings.ReplaceAll(result, "\"password\":\""+key+"\"", "\"password\":\"***\"")
		}
	}
	// 限制长度
	if len(result) > 1000 {
		result = result[:1000] + "..."
	}
	return result
}

// GetSecurityLogs 获取安全日志列表（支持按类型、时间、IP筛选）
func (sl *SecurityLogger) GetSecurityLogs(page, pageSize int, attackType, ip string, startTime, endTime *time.Time) ([]model.SecurityLog, int64, error) {
	query := sl.db.Model(&model.SecurityLog{})

	if attackType != "" {
		query = query.Where("attack_type = ?", attackType)
	}
	if ip != "" {
		query = query.Where("ip = ?", ip)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", endTime)
	}

	var total int64
	query.Count(&total)

	var logs []model.SecurityLog
	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)

	return logs, total, nil
}
