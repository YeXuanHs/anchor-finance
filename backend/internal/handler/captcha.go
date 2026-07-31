package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/email"
	"anchorfinance/pkg/sms"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CaptchaHandler handles captcha-related HTTP requests.
type CaptchaHandler struct {
	captchaService *service.CaptchaService
	geetestService *service.GeetestService
	smsSender      *sms.Sender
	emailSender    *email.Sender
}

// NewCaptchaHandler creates a new CaptchaHandler.
func NewCaptchaHandler(captchaService *service.CaptchaService, db *gorm.DB) *CaptchaHandler {
	return &CaptchaHandler{
		captchaService: captchaService,
		smsSender:      sms.NewSender(db),
		emailSender:    email.NewSender(db),
	}
}

// SetGeetestService 设置极验服务
func (h *CaptchaHandler) SetGeetestService(geetestService *service.GeetestService) {
	h.geetestService = geetestService
}

// smsRequest is the payload for SendSMS.
type smsRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// emailRequest is the payload for SendEmail.
type emailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// verifyImageRequest is the payload for VerifyImage.
type verifyImageRequest struct {
	Key       string `json:"key" binding:"required"`
	CaptchaID string `json:"captcha_id" binding:"required"`
	Digits    string `json:"digits" binding:"required"`
}

// geetestVerifyRequest is the payload for Geetest verification.
type geetestVerifyRequest struct {
	LotNumber     string `json:"lot_number" binding:"required"`
	CaptchaOutput string `json:"captcha_output" binding:"required"`
	PassToken     string `json:"pass_token" binding:"required"`
	GenTime       string `json:"gen_time" binding:"required"`
}

// GetImage generates and returns an image captcha.
// GET /captcha/image?key=xxx
func (h *CaptchaHandler) GetImage(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}

	captchaID, imgBytes, err := h.captchaService.GenerateImage(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate captcha"})
		return
	}

	// Return both the image and captcha ID
	c.Header("X-Captcha-ID", captchaID)
	c.Data(http.StatusOK, "image/png", imgBytes)
}

// GetImageJSON generates and returns an image captcha as JSON with base64.
// GET /captcha/image/json?key=xxx
func (h *CaptchaHandler) GetImageJSON(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}

	captchaID, imgBytes, err := h.captchaService.GenerateImage(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate captcha"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"captcha_id": captchaID,
		"image":      fmt.Sprintf("data:image/png;base64,%s", base64Encode(imgBytes)),
	})
}

// VerifyImage verifies an image captcha.
// POST /captcha/image/verify
func (h *CaptchaHandler) VerifyImage(c *gin.Context) {
	var req verifyImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.captchaService.VerifyImage(req.Key, req.CaptchaID, req.Digits) {
		c.JSON(http.StatusOK, gin.H{"message": "验证成功"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误或已过期"})
	}
}

// VerifyGeetest verifies Geetest 4.0 captcha.
// POST /captcha/geetest/verify
func (h *CaptchaHandler) VerifyGeetest(c *gin.Context) {
	if h.geetestService == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "极验服务未配置"})
		return
	}

	var req geetestVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := &service.GeetestVerifyParams{
		LotNumber:     req.LotNumber,
		CaptchaOutput: req.CaptchaOutput,
		PassToken:     req.PassToken,
		GenTime:       req.GenTime,
	}

	success, reason, err := h.geetestService.Verify(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "验证成功"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": reason})
	}
}

// GetGeetestConfig 获取极验前端配置
// GET /captcha/geetest/config
func (h *CaptchaHandler) GetGeetestConfig(c *gin.Context) {
	if h.geetestService == nil {
		c.JSON(http.StatusOK, gin.H{
			"enabled": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"enabled":    true,
		"captcha_id": h.geetestService.GetCaptchaID(),
	})
}

// SendSMS sends an SMS verification code.
// POST /captcha/sms
func (h *CaptchaHandler) SendSMS(c *gin.Context) {
	var req smsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := h.captchaService.GenerateSMS(req.Phone)
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	// Send SMS via configured provider
	if err := h.smsSender.Send(req.Phone, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to send SMS: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "验证码已发送",
	})
}

// VerifySMS verifies an SMS code.
// POST /captcha/sms/verify
func (h *CaptchaHandler) VerifySMS(c *gin.Context) {
	var req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.captchaService.VerifySMS(req.Phone, req.Code) {
		c.JSON(http.StatusOK, gin.H{"message": "验证成功"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误或已过期"})
	}
}

// SendEmail sends an email verification code.
// POST /captcha/email
func (h *CaptchaHandler) SendEmail(c *gin.Context) {
	var req emailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	code, err := h.captchaService.GenerateEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	// Send email with verification code
	subject := "验证码"
	htmlBody := fmt.Sprintf("<p>您的验证码是：<strong>%s</strong></p><p>请在5分钟内完成验证。</p>", code)

	if err := h.emailSender.Send(req.Email, subject, htmlBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to send email: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "验证码已发送",
	})
}

// VerifyEmail verifies an email code.
// POST /captcha/email/verify
func (h *CaptchaHandler) VerifyEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.captchaService.VerifyEmail(req.Email, req.Code) {
		c.JSON(http.StatusOK, gin.H{"message": "验证成功"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码错误或已过期"})
	}
}

// base64Encode encodes bytes to base64 string.
func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
