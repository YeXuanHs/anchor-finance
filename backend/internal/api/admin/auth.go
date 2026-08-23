package admin

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
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

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 管理员登录
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"data": nil,
		})
		return
	}

	// 获取客户端IP
	ip := c.ClientIP()

	token, err := h.authService.AdminLogin(req.Username, req.Password, ip)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    401,
			"message": "操作失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"token": token,
		},
	})
}

// GetInfo 获取管理员信息
func (h *AuthHandler) GetInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"user_id":  userID,
			"username": username,
			"avatar":   "",
			"role":     "admin",
		},
	})
}

// Logout 管理员登出
func (h *AuthHandler) Logout(c *gin.Context) {
	// 将token加入黑名单
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 {
		tokenStr := authHeader[7:]
		hash := sha256.Sum256([]byte(tokenStr))
		tokenHash := hex.EncodeToString(hash[:])

		claims, _ := h.authService.ParseToken(tokenStr)
		expiresAt := time.Now().Add(24 * time.Hour)
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

// UpdateProfile 更新个人资料
// PUT /api/admin/auth/profile
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	// 获取当前用户ID
	userID, _ := c.Get("user_id")

	var req struct {
		RealName string `json:"real_name"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	updates := map[string]interface{}{}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}

	db.Model(&model.Admin{}).Where("id = ?", userID).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    nil,
	})
}

// UpdatePassword 修改密码
// PUT /api/admin/auth/password
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	// 获取当前用户ID
	userID, _ := c.Get("user_id")

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	// 查询管理员
	db := database.GetDB()
	var admin model.Admin
	if err := db.First(&admin, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "管理员不存在",
			"data": nil,
		})
		return
	}

	// 验证旧密码
	if !service.CheckPassword(req.OldPassword, admin.PasswordHash) {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "旧密码错误",
			"data": nil,
		})
		return
	}

	// 生成新密码hash
	hashedPassword, err := service.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "密码加密失败",
			"data": nil,
		})
		return
	}

	// 更新密码
	db.Model(&admin).Update("password_hash", hashedPassword)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "密码修改成功",
		"data":    nil,
	})
}

// ResetAdminPassword 重置管理员密码（需要已登录的管理员权限）
// POST /api/admin/auth/reset-password
func (h *AuthHandler) ResetAdminPassword(c *gin.Context) {
	// 安全修复：必须是已登录的管理员才能重置密码
	currentAdminID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "请先登录", "data": nil})
		return
	}

	var req struct {
		Username    string `json:"username" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if len(req.NewPassword) < 6 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "密码长度至少6位", "data": nil})
		return
	}

	db := database.GetDB()

	// 验证当前操作者是否存在
	var currentAdmin model.Admin
	if err := db.First(&currentAdmin, currentAdminID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "操作者不存在", "data": nil})
		return
	}

	var targetAdmin model.Admin
	if err := db.Where("username = ?", req.Username).First(&targetAdmin).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "管理员不存在", "data": nil})
		return
	}

	hashedPassword, err := service.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "密码加密失败", "data": nil})
		return
	}

	db.Model(&targetAdmin).Updates(map[string]interface{}{
		"password_hash":  hashedPassword,
		"login_fail_count": 0,
		"locked_until": nil,
	})

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "密码重置成功", "data": nil})
}

// UnfreezeAdmin 超级管理员手动解冻（MD 9.3.1）
// POST /api/admin/auth/unfreeze/:id
func (h *AuthHandler) UnfreezeAdmin(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var admin model.Admin
	if err := db.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "管理员不存在", "data": nil})
		return
	}

	db.Model(&admin).Updates(map[string]interface{}{
		"login_fail_count": 0,
		"locked_until":     nil,
	})

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "解冻成功", "data": nil})
}
