package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/YeXuanHs/anchor-finance/config"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	db  *gorm.DB
	cfg *config.JWTConfig
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB, cfg *config.JWTConfig) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// Claims JWT声明
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// HashPassword 密码哈希
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken 生成JWT Token
func (s *AuthService) GenerateToken(userID uint, username string, isAdmin bool) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.cfg.ExpireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "anchor-finance",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.Secret))
}

// ParseToken 解析JWT Token
func (s *AuthService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// AdminLogin 管理员登录（带防暴力破解）
func (s *AuthService) AdminLogin(username, password, ip string) (string, error) {
	var admin model.Admin
	if err := s.db.Where("username = ?", username).First(&admin).Error; err != nil {
		return "", errors.New("用户名或密码错误")
	}

	// 检查账号状态
	if admin.Status != "active" {
		return "", errors.New("账号已被禁用")
	}

	// 检查是否被锁定
	if admin.LockedUntil != nil && admin.LockedUntil.After(time.Now()) {
		remaining := time.Until(*admin.LockedUntil)
		hours := int(remaining.Hours())
		minutes := int(remaining.Minutes()) % 60
		return "", fmt.Errorf("账号已冻结，请%d小时%d分钟后重试", hours, minutes)
	}

	// 验证密码
	if !CheckPassword(password, admin.PasswordHash) {
		// 增加失败次数
		admin.LoginFailCount++

		// 连续5次失败，冻结6小时
		if admin.LoginFailCount >= 5 {
			lockedUntil := time.Now().Add(6 * time.Hour)
			admin.LockedUntil = &lockedUntil
			admin.LoginFailCount = 0 // 重置计数
			s.db.Save(&admin)

			// 记录安全日志
			s.logSecurityEvent(admin.ID, username, ip, "admin_login_locked")

			return "", errors.New("连续5次密码错误，账号已冻结6小时")
		}

		s.db.Save(&admin)

		// 记录失败日志
		s.logSecurityEvent(admin.ID, username, ip, "admin_login_fail")

		return "", fmt.Errorf("用户名或密码错误（已失败%d次，连续5次将冻结）", admin.LoginFailCount)
	}

	// 登录成功，重置失败次数
	admin.LoginFailCount = 0
	admin.LockedUntil = nil
	now := time.Now()
	admin.LastLoginAt = &now
	admin.LastLoginIP = ip
	s.db.Save(&admin)

	return s.GenerateToken(admin.ID, admin.Username, true)
}

// UserLogin 用户登录
func (s *AuthService) UserLogin(username, password string) (string, error) {
	var user model.User
	if err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		return "", errors.New("用户名或密码错误")
	}

	if user.Status != "active" {
		return "", errors.New("账号已被禁用")
	}

	if !CheckPassword(password, user.PasswordHash) {
		return "", errors.New("用户名或密码错误")
	}

	// 更新最后登录时间
	now := time.Now()
	s.db.Model(&user).Updates(map[string]interface{}{
		"last_login_at": &now,
	})

	return s.GenerateToken(user.ID, user.Username, false)
}

// UnlockAdmin 解冻管理员账号
func (s *AuthService) UnlockAdmin(adminID uint) error {
	return s.db.Model(&model.Admin{}).Where("id = ?", adminID).Updates(map[string]interface{}{
		"login_fail_count": 0,
		"locked_until":     nil,
	}).Error
}

// logSecurityEvent 记录安全事件
func (s *AuthService) logSecurityEvent(adminID uint, username, ip, eventType string) {
	// 这里可以写入security_logs表
	// 暂时只打印日志
	fmt.Printf("[SECURITY] %s | admin_id=%d | username=%s | ip=%s\n",
		eventType, adminID, username, ip)
}
