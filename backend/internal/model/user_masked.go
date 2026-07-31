package model

import (
	"time"

	"anchorfinance/internal/security"
)

// MaskedUser 脱敏后的用户信息（前台显示用）
// 后台管理员看到完整的 User 信息，前台用户看到脱敏后的信息
type MaskedUser struct {
	ID          uint      `json:"id"`
	Username    string    `json:"username"`
	Nickname    string    `json:"nickname"`
	Avatar      string    `json:"avatar"`
	Phone       string    `json:"phone,omitempty"`       // 脱敏: 138****5678
	Email       string    `json:"email,omitempty"`       // 脱敏: tes***@example.com
	GroupID     uint      `json:"group_id"`
	GroupName   string    `json:"group_name,omitempty"`
	Status      int16     `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// MaskedLoginLog 脱敏后的登录日志（前台显示用）
type MaskedLoginLog struct {
	ID        uint      `json:"id"`
	IP        string    `json:"ip"`          // 脱敏: 192.168.*.*
	UserAgent string    `json:"user_agent"`
	Location  string    `json:"location"`
	LoginAt   time.Time `json:"login_at"`
	Status    int16     `json:"status"`      // 1=成功 0=失败
}

// MaskedCertification 脱敏后的实名认证信息（前台显示用）
type MaskedCertification struct {
	ID        uint      `json:"id"`
	RealName  string    `json:"real_name"`   // 脱敏: 张*
	IDCard    string    `json:"id_card"`     // 脱敏: 1101**********1234
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ────────────────────────────────────────────────────────────
// 转换函数
// ────────────────────────────────────────────────────────────

// ToMaskedUser 将 User 转换为脱敏后的 MaskedUser（前台显示）
func (u *User) ToMaskedUser() *MaskedUser {
	return &MaskedUser{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Phone:     security.MaskPhone(u.Phone),
		Email:     security.MaskEmail(u.Email),
		GroupID:   u.GroupID,
		GroupName: u.Group.Name,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}

// ToMaskedUserWithRole 根据角色决定是否脱敏
// isAdmin=true 时返回完整信息，isAdmin=false 时返回脱敏信息
func (u *User) ToMaskedUserWithRole(isAdmin bool) interface{} {
	if isAdmin {
		// 管理员看到完整信息
		return u
	}
	// 普通用户看到脱敏信息
	return u.ToMaskedUser()
}

// ToMaskedLoginLog 将 LoginLog 转换为脱敏后的 MaskedLoginLog（前台显示）
func (l *LoginLog) ToMaskedLoginLog() *MaskedLoginLog {
	return &MaskedLoginLog{
		ID:        l.ID,
		IP:        security.MaskIP(l.IP),
		UserAgent: l.UserAgent,
		Location:  l.Location,
		LoginAt:   l.CreatedAt,
		Status:    l.Status,
	}
}

// ────────────────────────────────────────────────────────────
// 批量转换
// ────────────────────────────────────────────────────────────

// ToMaskedUsers 批量转换用户列表（前台显示）
func ToMaskedUsers(users []User) []MaskedUser {
	result := make([]MaskedUser, len(users))
	for i, u := range users {
		result[i] = *u.ToMaskedUser()
	}
	return result
}

// ToMaskedUsersWithRole 根据角色批量转换用户列表
func ToMaskedUsersWithRole(users []User, isAdmin bool) interface{} {
	if isAdmin {
		return users
	}
	return ToMaskedUsers(users)
}

// ToMaskedLoginLogs 批量转换登录日志（前台显示）
func ToMaskedLoginLogs(logs []LoginLog) []MaskedLoginLog {
	result := make([]MaskedLoginLog, len(logs))
	for i, l := range logs {
		result[i] = *l.ToMaskedLoginLog()
	}
	return result
}
