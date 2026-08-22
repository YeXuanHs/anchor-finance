package client

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
	"github.com/YeXuanHs/anchor-finance/internal/service"
	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login 用户登录
// POST /api/client/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误: " + err.Error(), "data": nil})
		return
	}

	token, err := h.authService.UserLogin(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": err.Error(), "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{"token": token},
	})
}

// LoginByCode 验证码登录
// POST /api/client/auth/login-by-code
func (h *AuthHandler) LoginByCode(c *gin.Context) {
	var req struct {
		Target string `json:"target" binding:"required"` // 手机号或邮箱
		Code   string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误: " + err.Error(), "data": nil})
		return
	}

	db := database.GetDB()

	// 验证验证码
	var captcha model.Captcha
	if err := db.Where("target = ? AND code = ? AND type = ? AND used = ? AND expires_at > ?",
		req.Target, req.Code, "login", false, time.Now()).First(&captcha).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "验证码无效或已过期", "data": nil})
		return
	}

	// 标记验证码已使用
	db.Model(&captcha).Update("used", true)

	// 查找用户
	var user model.User
	if err := db.Where("phone = ? OR email = ?", req.Target, req.Target).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}

	if user.Status != "active" {
		c.JSON(http.StatusOK, gin.H{"code": 403, "message": "账号已被禁用", "data": nil})
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Username, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "生成token失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{"token": token},
	})
}

// SendCaptcha 发送验证码
// POST /api/client/auth/captcha
func (h *AuthHandler) SendCaptcha(c *gin.Context) {
	var req struct {
		Target string `json:"target" binding:"required"` // 手机号或邮箱
		Type   string `json:"type" binding:"required"`   // register, login, reset_password
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误: " + err.Error(), "data": nil})
		return
	}

	// 频率限制：从settings表读取配置（captcha_rate=每秒几次，默认1次/秒）
	db := database.GetDB()
	var setting model.Setting
	ratePerSecond := 3.0 // 默认3次/秒
	if err := db.Where("`key` = ?", "captcha_rate").First(&setting).Error; err == nil {
		if v, err := strconv.ParseFloat(setting.Value, 64); err == nil && v > 0 {
			ratePerSecond = v
		}
	}
	intervalSeconds := 1.0 / ratePerSecond
	var recentCount int64
	db.Model(&model.Captcha{}).Where("target = ? AND created_at > ?", req.Target, time.Now().Add(-time.Duration(intervalSeconds*float64(time.Second)))).Count(&recentCount)
	if recentCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 429, "message": "发送过于频繁，请稍后再试", "data": nil})
		return
	}

	// L2修复：使用crypto/rand生成安全验证码
	codeNum, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	code := fmt.Sprintf("%06d", codeNum.Int64())

	// 保存验证码
	captcha := model.Captcha{
		Target:    req.Target,
		Code:      code,
		Type:      req.Type,
		Used:      false,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	db.Create(&captcha)

	// 通过插件引擎发送（短信或邮件）
	if req.Type == "register" || req.Type == "login" {
		// 尝试通过短信发送
		pluginengine.SendSMS(req.Target, fmt.Sprintf("您的验证码是：%s，10分钟内有效。", code))
	} else {
		// 通过邮件发送
		pluginengine.SendEmail(req.Target, "验证码", fmt.Sprintf("您的验证码是：%s，10分钟内有效。", code))
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "验证码已发送", "data": nil})
}

// Register 用户注册
// POST /api/client/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Phone    string `json:"phone"`
		Code     string `json:"code"` // 验证码（可选，后台配置是否强制）
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误: " + err.Error(), "data": nil})
		return
	}

	db := database.GetDB()

	// 检查注册是否需要验证码（默认强制，防批量注册）
	var setting model.Setting
	requireCode := true // 默认需要验证码
	if err := db.Where("`key` = ?", "register_require_captcha").First(&setting).Error; err == nil && setting.Value == "0" {
		requireCode = false
	}

	if requireCode && req.Code == "" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请输入验证码", "data": nil})
		return
	}

	if requireCode {
		// 验证验证码
		var captcha model.Captcha
		if err := db.Where("target = ? AND code = ? AND type = ? AND used = ? AND expires_at > ?",
			req.Email, req.Code, "register", false, time.Now()).First(&captcha).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 400, "message": "验证码无效或已过期", "data": nil})
			return
		}
		db.Model(&captcha).Update("used", true)
	}

	// 检查用户名是否已存在
	var count int64
	db.Model(&model.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "用户名或邮箱已存在", "data": nil})
		return
	}

	// 创建用户
	hashedPassword, err := service.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "密码加密失败", "data": nil})
		return
	}

	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		Status:       "active",
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "注册失败: " + err.Error(), "data": nil})
		return
	}

	token, err := h.authService.GenerateToken(user.ID, user.Username, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "生成token失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "注册成功",
		"data": gin.H{"token": token},
	})
}

// ResetPassword 重置密码
// POST /api/client/auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Code     string `json:"code" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误: " + err.Error(), "data": nil})
		return
	}

	db := database.GetDB()

	// 验证验证码
	var captcha model.Captcha
	if err := db.Where("target = ? AND code = ? AND type = ? AND used = ? AND expires_at > ?",
		req.Email, req.Code, "reset_password", false, time.Now()).First(&captcha).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "如果邮箱存在，重置链接已发送", "data": nil})
		return
	}
	db.Model(&captcha).Update("used", true)

	// 查找用户
	var user model.User
	if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "如果邮箱存在，重置链接已发送", "data": nil})
		return
	}

	// 更新密码
	hashedPassword, _ := service.HashPassword(req.Password)
	db.Model(&user).Update("password_hash", hashedPassword)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "密码重置成功", "data": nil})
}

// GetInfo 获取用户信息
// GET /api/client/auth/info
func (h *AuthHandler) GetInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")

	db := database.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{
			"id":          user.ID,
			"username":    user.Username,
			"email":       user.Email,
			"phone":       user.Phone,
			"company":     user.Company,
			"balance":     user.Balance,
			"is_verified": user.IsVerified,
			"created_at":  user.CreatedAt,
		},
	})
}

// Logout 用户登出（token加入黑名单）
// POST /api/client/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// 获取当前token并加入黑名单
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 {
		tokenStr := authHeader[7:]
		hash := sha256.Sum256([]byte(tokenStr))
		tokenHash := hex.EncodeToString(hash[:])

		// 解析token获取过期时间
		claims, _ := h.authService.ParseToken(tokenStr)
		expiresAt := time.Now().Add(24 * time.Hour) // 默认24小时
		if claims != nil {
			expiresAt = claims.ExpiresAt.Time
		}

		db := database.GetDB()
		db.Create(&model.TokenBlacklist{
			TokenHash: tokenHash,
			ExpiresAt: expiresAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": nil})
}

// UpdatePassword 修改密码
// PUT /api/client/password
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误: " + err.Error(), "data": nil})
		return
	}

	db := database.GetDB()
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "用户不存在", "data": nil})
		return
	}

	if !service.CheckPassword(req.OldPassword, user.PasswordHash) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "旧密码错误", "data": nil})
		return
	}

	hashedPassword, err := service.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "密码加密失败", "data": nil})
		return
	}

	db.Model(&user).Update("password_hash", hashedPassword)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "密码修改成功", "data": nil})
}

// UpdateProfile 更新个人资料
// PUT /api/client/auth/profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Phone   string `json:"phone"`
		Company string `json:"company"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误: " + err.Error(), "data": nil})
		return
	}

	db := database.GetDB()
	updates := map[string]interface{}{}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Company != "" {
		updates["company"] = req.Company
	}

	db.Model(&model.User{}).Where("id = ?", userID).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// UpdatePhone 更新手机号
// PUT /api/client/auth/phone
func (h *AuthHandler) UpdatePhone(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()

	// 验证验证码
	var captcha model.Captcha
	if err := db.Where("target = ? AND code = ? AND type = ? AND used = ? AND expires_at > ?",
		req.Phone, req.Code, "bindphone", false, time.Now()).First(&captcha).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "验证码无效", "data": nil})
		return
	}
	db.Model(&captcha).Update("used", true)

	db.Model(&model.User{}).Where("id = ?", userID).Update("phone", req.Phone)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "手机号更新成功", "data": nil})
}

// UpdateEmail 更新邮箱
// PUT /api/client/auth/email
func (h *AuthHandler) UpdateEmail(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()

	// 验证验证码
	var captcha model.Captcha
	if err := db.Where("target = ? AND code = ? AND type = ? AND used = ? AND expires_at > ?",
		req.Email, req.Code, "bindemail", false, time.Now()).First(&captcha).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "验证码无效", "data": nil})
		return
	}
	db.Model(&captcha).Update("used", true)

	db.Model(&model.User{}).Where("id = ?", userID).Update("email", req.Email)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "邮箱更新成功", "data": nil})
}
