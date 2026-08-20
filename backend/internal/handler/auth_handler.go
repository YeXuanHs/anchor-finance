package handler

import (
	"net/http"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/auth"
	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler handles authentication requests.
type AuthHandler struct{}

// LoginRequest represents the login request body.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles admin login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"ok": false, "message": "参数错误: " + err.Error()})
		return
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "message": "数据库未初始化"})
		return
	}

	// Look up admin user in users table (admins are users with is_admin=true)
	var user model.User
	err := database.Where("username = ? AND is_admin = ?", req.Username, true).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "message": "用户名或密码错误"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "message": "查询失败"})
		return
	}

	// Check if user is active
	if user.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{"ok": false, "message": "账户已被禁用"})
		return
	}

	// Verify password (using PasswordHash field)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "message": "用户名或密码错误"})
		return
	}

	// Get JWT secret from system settings
	jwtSecret := db.GetSystemSetting("jwt_secret")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "message": "JWT配置错误"})
		return
	}

	// Generate JWT token
	jwtMgr := auth.NewJWTManager(jwtSecret, 72)
	token, err := jwtMgr.GenerateToken(user.ID, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "message": "生成token失败"})
		return
	}

	// Generate refresh token (longer expiry)
	refreshToken, err := jwtMgr.GenerateToken(user.ID, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "message": "生成refresh token失败"})
		return
	}

	// Update last login time
	now := time.Now()
	database.Model(&user).Update("last_login_at", &now)

	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"code": 0,
		"data": gin.H{
			"access_token":  token,
			"refresh_token": refreshToken,
			"id":            user.ID,
			"username":      user.Username,
			"email":         user.Email,
			"nickname":      user.Nickname,
			"avatar":        user.Avatar,
			"is_admin":      user.IsAdmin,
		},
	})
}

// SMSLogin handles SMS login.
func (h *AuthHandler) SMSLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true, "code": 0, "message": "ok"})
}

// RefreshToken refreshes the access token.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get current user from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"ok": false, "message": "未登录"})
		return
	}

	jwtSecret := db.GetSystemSetting("jwt_secret")
	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "message": "JWT配置错误"})
		return
	}

	jwtMgr := auth.NewJWTManager(jwtSecret, 72)
	token, err := jwtMgr.GenerateToken(userID.(uint), true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"ok": false, "message": "刷新token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok":   true,
		"code": 0,
		"data": gin.H{
			"access_token": token,
		},
	})
}

// Logout handles user logout.
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true, "code": 0, "message": "已登出"})
}

// AccessTokenLogin handles access token login.
func (h *AuthHandler) AccessTokenLogin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ok": true, "code": 0, "message": "ok"})
}
