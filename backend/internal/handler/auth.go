package handler

import (
	"net/http"
	"time"

	"anchorfinance/internal/api/middleware"
	"anchorfinance/internal/service"
	"anchorfinance/internal/validator"
	"anchorfinance/pkg/auth"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userSvc      *service.UserService
	captchaSvc   *service.CaptchaService
	emailSuffixSvc *service.EmailSuffixWhitelistService
	log          *logger.Logger
	jwtMgr       *auth.JWTManager
}

func NewAuthHandler(userSvc *service.UserService, log *logger.Logger, jwtMgr *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userSvc: userSvc,
		log:     log,
		jwtMgr:  jwtMgr,
	}
}

func NewAuthHandlerWithCaptcha(userSvc *service.UserService, captchaSvc *service.CaptchaService, log *logger.Logger, jwtMgr *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userSvc:    userSvc,
		captchaSvc: captchaSvc,
		log:        log,
		jwtMgr:     jwtMgr,
	}
}

// SetEmailSuffixService 设置邮箱后缀白名单服务
func (h *AuthHandler) SetEmailSuffixService(svc *service.EmailSuffixWhitelistService) {
	h.emailSuffixSvc = svc
}

type loginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
}

// Login authenticates and returns JWT tokens.
func (h *AuthHandler) Login(c *gin.Context) {
	ip := c.ClientIP()

	// 多层限流检查（移植自 zjmf）
	ok, msg := middleware.CheckLoginLimit(ip)
	if !ok {
		response.TooManyRequests(c, msg)
		return
	}

	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证账号不能为空
	if req.Account == "" {
		response.BadRequest(c, "账号不能为空")
		return
	}

	// 验证密码不能为空
	if req.Password == "" {
		response.BadRequest(c, "密码不能为空")
		return
	}

	// 验证密码长度（移植自 zjmf：6-255位）
	if len(req.Password) < 6 {
		response.BadRequest(c, "密码长度至少6位")
		return
	}

	user, err := h.userSvc.Login(service.LoginRequest{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		// 记录失败（移植自 zjmf login_inc）
		middleware.RecordLoginAttempt(ip)
		h.log.Warnf("登录失败: ip=%s account=%s err=%v", ip, req.Account, err)
		response.Unauthorized(c, "账号或密码错误")
		return
	}

	// 登录成功，清除失败记录
	middleware.ClearLoginAttempts(ip)

	// 判断是否管理员
	isAdmin := user.Role == "admin"

	accessToken, refreshToken, err := h.generateTokens(user.ID, isAdmin, ip)
	if err != nil {
		response.ServerError(c, "failed to generate tokens")
		return
	}

	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

// Register creates a new account and returns tokens.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证用户名（移植自 zjmf：4-20位）
	if ok, msg := validator.ValidateUsername(req.Username); !ok {
		response.BadRequest(c, msg)
		return
	}

	// 验证密码（移植自 zjmf：6-32位）
	if ok, msg := validator.ValidatePassword(req.Password); !ok {
		response.BadRequest(c, msg)
		return
	}

	// 验证邮箱（如果提供）
	if req.Email != "" {
		if ok, msg := validator.ValidateEmailOptional(req.Email); !ok {
			response.BadRequest(c, msg)
			return
		}
		// 邮箱后缀白名单校验
		if h.emailSuffixSvc != nil && !h.emailSuffixSvc.IsAllowed(req.Email) {
			response.BadRequest(c, "该邮箱后缀不允许注册，请使用其他邮箱")
			return
		}
	}

	// 验证手机号（如果提供）
	if req.Phone != "" {
		if ok, msg := validator.ValidatePhone(req.Phone); !ok {
			response.BadRequest(c, msg)
			return
		}
	}

	// 验证昵称（如果提供）
	if req.Nickname != "" {
		if ok, msg := validator.ValidateNickname(req.Nickname); !ok {
			response.BadRequest(c, msg)
			return
		}
	}

	user, err := h.userSvc.Register(service.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	accessToken, refreshToken, err := h.generateTokens(user.ID, user.Role == "admin", c.ClientIP())
	if err != nil {
		response.ServerError(c, "failed to generate tokens")
		return
	}

	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

// RefreshToken issues new tokens from a valid refresh token.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	claims, err := h.jwtMgr.ValidateToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "invalid refresh token")
		return
	}

	// 检查 Token 是否在密码修改之后签发（移植自 zjmf）
	if !middleware.IsTokenValid(claims.UserID, claims.IssuedAt.Unix()) {
		response.Unauthorized(c, "密码已修改，请重新登录")
		return
	}

	accessToken, refreshToken, err := h.generateTokens(claims.UserID, claims.IsAdmin, c.ClientIP())
	if err != nil {
		response.ServerError(c, "failed to generate tokens")
		return
	}

	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// SMSLogin handles login via SMS verification code.
func (h *AuthHandler) SMSLogin(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证手机号（移植自 zjmf：4-11位）
	if ok, msg := validator.ValidatePhone(req.Phone); !ok {
		response.BadRequest(c, msg)
		return
	}

	// Verify SMS code
	if h.captchaSvc == nil || !h.captchaSvc.VerifySMS(req.Phone, req.Code) {
		response.BadRequest(c, "验证码无效或已过期")
		return
	}

	// Find user by phone
	user, err := h.userSvc.GetByPhone(req.Phone)
	if err != nil {
		response.BadRequest(c, "手机号未注册")
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := h.generateTokens(user.ID, user.Role == "admin", c.ClientIP())
	if err != nil {
		response.ServerError(c, "failed to generate tokens")
		return
	}

	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

// generateTokens creates access and refresh JWT tokens.
func (h *AuthHandler) generateTokens(userID uint, isAdmin bool, ip string) (string, string, error) {
	// 使用 JWTManager 生成 token（带 IP 绑定）
	accessToken, err := h.jwtMgr.GenerateTokenWithIP(userID, isAdmin, ip)
	if err != nil {
		return "", "", err
	}

	// Refresh token 也使用 JWTManager，但有效期更长
	// 这里简化处理，实际可以创建专门的 refresh token
	refreshToken, err := h.jwtMgr.GenerateTokenWithIP(userID, isAdmin, ip)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ChangePassword handles password change.
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		OldPassword     string `json:"old_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证旧密码（移植自 zjmf：6-32位）
	if ok, msg := validator.ValidatePassword(req.OldPassword); !ok {
		response.BadRequest(c, msg)
		return
	}

	// 验证新密码（移植自 zjmf：6-32位）
	if ok, msg := validator.ValidatePassword(req.NewPassword); !ok {
		response.BadRequest(c, msg)
		return
	}

	// 验证确认密码
	if ok, msg := validator.ValidatePasswordMatch(req.NewPassword, req.ConfirmPassword); !ok {
		response.BadRequest(c, msg)
		return
	}

	// 验证新密码不能和旧密码相同
	if ok, msg := validator.ValidatePasswordNotSame(req.OldPassword, req.NewPassword); !ok {
		response.BadRequest(c, msg)
		return
	}

	if err := h.userSvc.ChangePassword(userID, service.ChangePasswordRequest{
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 密码修改后使所有 Token 失效（移植自 zjmf client_user_update_pass_）
	middleware.InvalidateUserTokens(userID)

	response.SuccessMsg(c, "密码修改成功，请重新登录")
}

// LoginRateLimitStatus returns the current login rate limit status for an IP.
func (h *AuthHandler) LoginRateLimitStatus(c *gin.Context) {
	ip := c.ClientIP()
	ok, msg := middleware.CheckLoginLimit(ip)
	if !ok {
		response.Success(c, gin.H{
			"locked": true,
			"message": msg,
		})
	} else {
		response.Success(c, gin.H{
			"locked": false,
		})
	}
}

// AccessTokenLogin handles login via access_token (from zjmf loginAccessToken).
// POST /auth/access-token
func (h *AuthHandler) AccessTokenLogin(c *gin.Context) {
	var req struct {
		AccessToken string `json:"access_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Validate the access token and get user
	user, err := h.userSvc.GetUserByAccessToken(req.AccessToken)
	if err != nil {
		response.Unauthorized(c, "access_token无效或已过期")
		return
	}

	isAdmin := user.Role == "admin"
	accessToken, refreshToken, err := h.generateTokens(user.ID, isAdmin, c.ClientIP())
	if err != nil {
		response.ServerError(c, "failed to generate tokens")
		return
	}

	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

// cleanup cleans up expired rate limit entries periodically.
func (h *AuthHandler) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		// 清理过期的限流记录
		middleware.CleanupExpiredEntries()
	}
}
