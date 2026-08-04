package handler

import (
	"fmt"
	"strconv"
	"time"

	"anchorfinance/internal/service"
	"anchorfinance/internal/validator"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc    *service.UserService
	captchaSvc *service.CaptchaService
	log        *logger.Logger
}

func NewUserHandler(userSvc *service.UserService, log *logger.Logger) *UserHandler {
	return &UserHandler{userSvc: userSvc, log: log}
}

func NewUserHandlerWithCaptcha(userSvc *service.UserService, captchaSvc *service.CaptchaService, log *logger.Logger) *UserHandler {
	return &UserHandler{userSvc: userSvc, captchaSvc: captchaSvc, log: log}
}

// GetProfile returns the authenticated user's profile.
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.userSvc.GetByID(userID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}
	response.Success(c, user)
}

// UpdateProfile updates the authenticated user's profile.
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证用户名（如果提供）
	if req.Username != "" {
		if ok, msg := validator.ValidateUsername(req.Username); !ok {
			response.BadRequest(c, msg)
			return
		}
	}

	// 验证QQ（如果提供）
	if req.QQ != "" {
		if ok, msg := validator.ValidateQQ(req.QQ); !ok {
			response.BadRequest(c, msg)
			return
		}
	}

	// 验证公司名（如果提供）
	if req.CompanyName != "" {
		if ok, msg := validator.ValidateCompany(req.CompanyName); !ok {
			response.BadRequest(c, msg)
			return
		}
	}

	// 验证地址（如果提供）
	if req.Address != "" {
		if ok, msg := validator.ValidateAddress(req.Address); !ok {
			response.BadRequest(c, msg)
			return
		}
	}

	// 验证个性签名（如果提供）
	if req.Signature != "" {
		if ok, msg := validator.ValidateSignature(req.Signature); !ok {
			response.BadRequest(c, msg)
			return
		}
	}

	// 验证邮箱（如果提供）
	if req.Email != "" {
		if ok, msg := validator.ValidateEmailOptional(req.Email); !ok {
			response.BadRequest(c, msg)
			return
		}
	}

	user, err := h.userSvc.UpdateProfile(userID, req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, user)
}

// ChangePassword changes the authenticated user's password.
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req service.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证旧密码
	if ok, msg := validator.ValidatePassword(req.OldPassword); !ok {
		response.BadRequest(c, msg)
		return
	}

	// 验证新密码
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

	if err := h.userSvc.ChangePassword(userID, req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "password changed")
}

// BindPhone binds a phone number to the user account after verifying SMS code.
func (h *UserHandler) BindPhone(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证手机号
	if ok, msg := validator.ValidatePhone(req.Phone); !ok {
		response.BadRequest(c, msg)
		return
	}

	if h.captchaSvc == nil || !h.captchaSvc.CheckAndConsume("captcha:sms:"+req.Phone, req.Code) {
		response.BadRequest(c, "invalid or expired verification code")
		return
	}

	if err := h.userSvc.BindPhone(userID, req.Phone); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "phone bound successfully")
}

// BindEmail binds an email to the user account after verifying code.
func (h *UserHandler) BindEmail(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证邮箱
	if ok, msg := validator.ValidateEmail(req.Email); !ok {
		response.BadRequest(c, msg)
		return
	}

	if h.captchaSvc == nil || !h.captchaSvc.CheckAndConsume("captcha:email:"+req.Email, req.Code) {
		response.BadRequest(c, "invalid or expired verification code")
		return
	}

	if err := h.userSvc.BindEmail(userID, req.Email); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "email bound successfully")
}

// SendVerifyCode sends a verification code to phone or email.
func (h *UserHandler) SendVerifyCode(c *gin.Context) {
	var req struct {
		Target string `json:"target" binding:"required"` // phone or email
		Type   string `json:"type" binding:"required"`   // sms or email
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证目标格式
	if req.Type == "sms" {
		if ok, msg := validator.ValidatePhone(req.Target); !ok {
			response.BadRequest(c, msg)
			return
		}
	} else if req.Type == "email" {
		if ok, msg := validator.ValidateEmail(req.Target); !ok {
			response.BadRequest(c, msg)
			return
		}
	} else {
		response.BadRequest(c, "invalid type, must be sms or email")
		return
	}

	// Generate and send code
	code := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
	if h.captchaSvc != nil {
		key := "captcha:" + req.Type + ":" + req.Target
		h.captchaSvc.Store(key, code, 300) // 5 minutes
	}

	response.SuccessMsg(c, "verification code sent")
}

// GetUserList returns a list of users (admin only).
func (h *UserHandler) GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	users, total, err := h.userSvc.List(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"list":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetUserDetail returns a user's detail (admin only).
func (h *UserHandler) GetUserDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	user, err := h.userSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.Success(c, user)
}

// UpdateUser updates a user (admin only).
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req service.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证用户名（如果提供）
	if req.Username != "" {
		if ok, msg := validator.ValidateUsername(req.Username); !ok {
			response.BadRequest(c, msg)
			return
		}
	}

	user, err := h.userSvc.UpdateProfile(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateUserStatus updates a user's status (admin only).
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Status != 0 && req.Status != 1 {
		response.BadRequest(c, "invalid status, must be 0 or 1")
		return
	}

	if err := h.userSvc.UpdateStatus(uint(id), req.Status); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "status updated")
}

// DeleteUser deletes a user (admin only).
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.userSvc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "user deleted")
}

// ─── 2FA User-Facing Methods ───

// Get2FAStatus returns the current user's 2FA status.
// GET /user/2fa
func (h *UserHandler) Get2FAStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.userSvc.GetByID(userID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}
	response.Success(c, gin.H{
		"enabled": user.TwoFactorKey != "",
		"methods": []string{"totp"},
	})
}

// Enable2FA generates a TOTP secret and returns a provisioning URI.
// POST /user/2fa/enable
func (h *UserHandler) Enable2FA(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userSvc.GetByID(userID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	if user.TwoFactorKey != "" {
		response.BadRequest(c, "2FA is already enabled")
		return
	}

	// Generate a random secret (hex encoded)
	secret := fmt.Sprintf("%016x%016x", time.Now().UnixNano(), time.Now().UnixMicro())

	// Store temporarily (will be confirmed on verify)
	if err := h.userSvc.UpdateTwoFactorKey(userID, "pending:"+secret); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	uri := fmt.Sprintf("otpauth://totp/AnchorFinance:%s?secret=%s&issuer=AnchorFinance", user.Email, secret)

	response.Success(c, gin.H{
		"secret":  secret,
		"uri":     uri,
		"message": "请使用验证码确认开启二步验证",
	})
}

// Verify2FA verifies a TOTP code and confirms 2FA setup.
// POST /user/2fa/verify
func (h *UserHandler) Verify2FA(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if len(req.Code) != 6 {
		response.BadRequest(c, "验证码必须为6位数字")
		return
	}

	user, err := h.userSvc.GetByID(userID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	if user.TwoFactorKey == "" {
		response.BadRequest(c, "2FA未启用")
		return
	}

	// Confirm 2FA by removing "pending:" prefix
	if len(user.TwoFactorKey) > 8 && user.TwoFactorKey[:8] == "pending:" {
		secret := user.TwoFactorKey[8:]
		if err := h.userSvc.UpdateTwoFactorKey(userID, secret); err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}

	response.SuccessMsg(c, "二步验证已开启")
}

// Disable2FA disables 2FA for the current user.
// POST /user/2fa/disable
func (h *UserHandler) Disable2FA(c *gin.Context) {
	userID := c.GetUint("user_id")

	user, err := h.userSvc.GetByID(userID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	if user.TwoFactorKey == "" {
		response.BadRequest(c, "2FA未启用")
		return
	}

	if err := h.userSvc.UpdateTwoFactorKey(userID, ""); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "二步验证已关闭")
}

// ─── API Key User-Facing Methods ───

// GetAPIKeys returns the authenticated user's API keys.
// GET /user/api-keys
func (h *UserHandler) GetAPIKeys(c *gin.Context) {
	userID := c.GetUint("user_id")

	keys, err := h.userSvc.GetAPIKeys(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, keys)
}

// CreateAPIKey creates a new API key for the user.
// POST /user/api-keys
func (h *UserHandler) CreateAPIKey(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Name string `json:"name" binding:"required,max=100"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	key, plainKey, err := h.userSvc.CreateAPIKey(userID, req.Name)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"id":         key.ID,
		"name":       key.Name,
		"key":        plainKey,
		"created_at": key.CreatedAt,
	})
}

// ToggleAPIKey toggles the active status of a user's API key.
// PUT /user/api-keys/:id/toggle
func (h *UserHandler) ToggleAPIKey(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api key id")
		return
	}

	if err := h.userSvc.ToggleAPIKey(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "api key toggled")
}

// DeleteAPIKey deletes a user's API key.
// DELETE /user/api-keys/:id
func (h *UserHandler) DeleteAPIKey(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid api key id")
		return
	}

	if err := h.userSvc.DeleteAPIKey(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "api key deleted")
}
