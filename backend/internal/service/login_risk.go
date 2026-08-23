package service

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// LoginRiskControl 登录风控（数据库存储）
type LoginRiskControl struct {
	db *gorm.DB
}

const (
	accountIPMaxAttempts = 5              // 同账号+同IP 5次失败
	accountIPDecay       = 15 * time.Minute // 锁15分钟
	ipMaxAttempts        = 15             // 同IP 15次失败
	ipDecay              = 30 * time.Minute // 锁30分钟
)

// LoginAttempt 登录尝试记录
type LoginAttempt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"size:100;index" json:"username"`
	IP        string    `gorm:"size:45;index" json:"ip"`
	Failures  int       `gorm:"default:0" json:"failures"`
	LastFail  time.Time `json:"last_fail"`
	LockedUntil *time.Time `json:"locked_until"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 指定表名
func (LoginAttempt) TableName() string {
	return "login_attempts"
}

// NewLoginRiskControl 创建登录风控
func NewLoginRiskControl(db *gorm.DB) *LoginRiskControl {
	return &LoginRiskControl{db: db}
}

// IsLocked 检查是否被锁定
func (lrc *LoginRiskControl) IsLocked(username, ip string) (bool, string) {
	now := time.Now()

	// 检查账号+IP锁定
	var accountIPAttempt LoginAttempt
	if err := lrc.db.Where("username = ? AND ip = ?", username, ip).First(&accountIPAttempt).Error; err == nil {
		if accountIPAttempt.LockedUntil != nil && now.Before(*accountIPAttempt.LockedUntil) {
			remaining := time.Until(*accountIPAttempt.LockedUntil)
			return true, fmt.Sprintf("该账号在此IP登录失败次数过多，请%d分钟后重试", int(remaining.Minutes())+1)
		}
	}

	// 检查IP锁定（同一IP对多个账号暴力破解）
	var ipAttempt LoginAttempt
	if err := lrc.db.Where("ip = ? AND locked_until IS NOT NULL AND locked_until > ?", ip, now).First(&ipAttempt).Error; err == nil {
		remaining := time.Until(*ipAttempt.LockedUntil)
		return true, fmt.Sprintf("该IP登录失败次数过多，请%d分钟后重试", int(remaining.Minutes())+1)
	}

	return false, ""
}

// RecordFailure 记录登录失败
func (lrc *LoginRiskControl) RecordFailure(username, ip string) {
	now := time.Now()

	// 更新账号+IP记录
	var accountIPAttempt LoginAttempt
	result := lrc.db.Where("username = ? AND ip = ?", username, ip).First(&accountIPAttempt)
	if result.Error != nil {
		// 不存在则创建
		accountIPAttempt = LoginAttempt{
			Username: username,
			IP:       ip,
			Failures: 1,
			LastFail: now,
		}
		lrc.db.Create(&accountIPAttempt)
	} else {
		// 存在则更新
		accountIPAttempt.Failures++
		accountIPAttempt.LastFail = now
		if accountIPAttempt.Failures >= accountIPMaxAttempts {
			lockUntil := now.Add(accountIPDecay)
			accountIPAttempt.LockedUntil = &lockUntil
		}
		lrc.db.Save(&accountIPAttempt)
	}

	// 统计该IP对不同账号的失败总数（防跨账号暴力破解）
	var ipFailures int64
	lrc.db.Model(&LoginAttempt{}).Where("ip = ? AND last_fail > ?", ip, now.Add(-ipDecay)).Count(&ipFailures)
	if ipFailures >= int64(ipMaxAttempts) {
		// IP级别锁定：更新该IP所有记录的锁定时间
		lockUntil := now.Add(ipDecay)
		lrc.db.Model(&LoginAttempt{}).Where("ip = ?", ip).Update("locked_until", &lockUntil)
	}
}

// ClearSuccess 登录成功时清除该账号+IP的失败记录
func (lrc *LoginRiskControl) ClearSuccess(username, ip string) {
	lrc.db.Where("username = ? AND ip = ?", username, ip).Delete(&LoginAttempt{})
}

// Cleanup 过期记录清理（定时任务调用）
func (lrc *LoginRiskControl) Cleanup() {
	lrc.db.Where("locked_until IS NOT NULL AND locked_until < ?", time.Now().Add(-1*time.Hour)).Delete(&LoginAttempt{})
	lrc.db.Where("last_fail < ?", time.Now().Add(-24*time.Hour)).Delete(&LoginAttempt{})
}
