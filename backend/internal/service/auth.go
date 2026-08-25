package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/YeXuanHs/anchor-finance/config"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	db      *gorm.DB
	cfg     *config.JWTConfig
	risk    *LoginRiskControl
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB, cfg *config.JWTConfig) *AuthService {
	return &AuthService{
		db:   db,
		cfg:  cfg,
		risk: NewLoginRiskControl(db),
	}
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

// GenerateToken 生成JWT Token（Access + Refresh）
func (s *AuthService) GenerateToken(userID uint, username string, isAdmin bool) (string, string, error) {
	// Access Token（短期，2小时）
	accessClaims := Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "anchor-finance",
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", "", err
	}

	// Refresh Token（长期，7天）
	refreshClaims := Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "anchor-finance-refresh",
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshToken.SignedString([]byte(s.cfg.Secret))
	if err != nil {
		return "", "", err
	}

	// 存储Refresh Token到数据库
	expireAt := time.Now().Add(7 * 24 * time.Hour)
	if isAdmin {
		s.db.Model(&model.Admin{}).Where("id = ?", userID).Updates(map[string]interface{}{
			"refresh_token": refreshString,
			"refresh_expire": &expireAt,
		})
	} else {
		s.db.Model(&model.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
			"refresh_token": refreshString,
			"refresh_expire": &expireAt,
		})
	}

	return accessString, refreshString, nil
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

// AdminLogin 管理员登录（带防暴力破解+IP软锁）
func (s *AuthService) AdminLogin(username, password, ip string) (string, string, error) {
	// IP风控检查
	if locked, msg := s.risk.IsLocked(username, ip); locked {
		return "", "", errors.New(msg)
	}

	var admin model.Admin
	if err := s.db.Where("username = ?", username).First(&admin).Error; err != nil {
		s.risk.RecordFailure(username, ip)
		return "", "", errors.New("用户名或密码错误")
	}

	// 检查账号状态
	if admin.Status != "active" {
		return "", "", errors.New("账号已被禁用")
	}

	// 检查是否被锁定
	if admin.LockedUntil != nil && admin.LockedUntil.After(time.Now()) {
		remaining := time.Until(*admin.LockedUntil)
		hours := int(remaining.Hours())
		minutes := int(remaining.Minutes()) % 60
		return "", "", fmt.Errorf("账号已冻结，请%d小时%d分钟后重试", hours, minutes)
	}

	// 验证密码
	if !CheckPassword(password, admin.PasswordHash) {
		// 记录IP风控失败
		s.risk.RecordFailure(username, ip)

		// 增加失败次数
		admin.LoginFailCount++

		// 从settings读取配置（MD 9.2：后台可配置）
		maxFail := GetSettingInt("admin_login_max_fail", 5)
		lockHours := GetSettingInt("admin_login_lock_hours", 6)

		// 连续N次失败，冻结M小时
		if admin.LoginFailCount >= maxFail {
			lockedUntil := time.Now().Add(time.Duration(lockHours) * time.Hour)
			admin.LockedUntil = &lockedUntil
			admin.LoginFailCount = 0
			s.db.Save(&admin)

			s.logSecurityEvent(admin.ID, username, ip, "admin_login_locked")
			return "", "", fmt.Errorf("连续%d次密码错误，账号已冻结%d小时", maxFail, lockHours)
		}

		s.db.Save(&admin)
		s.logSecurityEvent(admin.ID, username, ip, "admin_login_fail")
		return "", "", fmt.Errorf("用户名或密码错误（已失败%d次，连续%d次将冻结）", admin.LoginFailCount, maxFail)
	}

	// 登录成功，重置失败次数+清除风控
	s.risk.ClearSuccess(username, ip)
	admin.LoginFailCount = 0
	admin.LockedUntil = nil
	now := time.Now()
	admin.LastLoginAt = &now
	admin.LastLoginIP = ip
	s.db.Save(&admin)

	return s.GenerateToken(admin.ID, admin.Username, true)
}

// UserLogin 用户登录（M2修复：带防暴力破解锁定）
func (s *AuthService) UserLogin(username, password, ip string) (string, string, error) {
	// IP风控检查
	if locked, msg := s.risk.IsLocked(username, ip); locked {
		return "", "", errors.New(msg)
	}

	var user model.User
	if err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		s.risk.RecordFailure(username, ip)
		return "", "", errors.New("用户名或密码错误")
	}

	if user.Status != "active" {
		return "", "", errors.New("账号已被禁用")
	}

	// M2修复：检查是否被锁定
	if user.LockedUntil != nil && user.LockedUntil.After(time.Now()) {
		remaining := time.Until(*user.LockedUntil)
		hours := int(remaining.Hours())
		minutes := int(remaining.Minutes()) % 60
		return "", "", fmt.Errorf("账号已冻结，请%d小时%d分钟后重试", hours, minutes)
	}

	if !CheckPassword(password, user.PasswordHash) {
		// 记录IP风控失败
		s.risk.RecordFailure(username, ip)

		// M2修复：增加失败次数
		user.LoginFailCount++

		// 从settings读取配置（MD 9.2：后台可配置）
		maxFail := GetSettingInt("client_login_max_fail", 5)
		lockHours := GetSettingInt("client_login_lock_hours", 6)

		if user.LoginFailCount >= maxFail {
			lockedUntil := time.Now().Add(time.Duration(lockHours) * time.Hour)
			user.LockedUntil = &lockedUntil
			user.LoginFailCount = 0
			s.db.Save(&user)
			return "", "", fmt.Errorf("连续%d次密码错误，账号已冻结%d小时", maxFail, lockHours)
		}

		s.db.Save(&user)
		return "", "", fmt.Errorf("用户名或密码错误（已失败%d次，连续%d次将冻结）", user.LoginFailCount, maxFail)
	}

	// 登录成功，重置失败次数+清除风控
	s.risk.ClearSuccess(username, ip)
	user.LoginFailCount = 0
	user.LockedUntil = nil
	now := time.Now()
	user.LastLoginAt = &now
	s.db.Save(&user)

	return s.GenerateToken(user.ID, user.Username, false)
}

// UnlockAdmin 解冻管理员账号
func (s *AuthService) UnlockAdmin(adminID uint) error {
	return s.db.Model(&model.Admin{}).Where("id = ?", adminID).Updates(map[string]interface{}{
		"login_fail_count": 0,
		"locked_until":     nil,
	}).Error
}

// logSecurityEvent 记录安全事件（MD 9.1：记录到security_logs表）
func (s *AuthService) logSecurityEvent(adminID uint, username, ip, eventType string) {
	// 同时记录到operation_logs（兼容）
	s.db.Create(&model.OperationLog{
		UserID:   adminID,
		Username: username,
		Action:   eventType,
		Resource: "security",
		Detail:   "安全事件: " + eventType,
		IP:       ip,
	})
}

// GenerateTokenStatic 静态版本GenerateToken（用于zjmf兼容登录等不需要AuthService实例的场景）
func GenerateTokenStatic(userID uint, username string, isAdmin bool) (string, error) {
	// 安全要求：强制从环境变量读取密钥，不再使用硬编码备用密钥
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable not set")
	}

	claims := Claims{
		UserID:   userID,
		Username: username,
		IsAdmin:  isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "anchor-finance",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
