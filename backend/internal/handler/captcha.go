package handler

import (
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

// smsRequest is the payload for SendSMS.
type smsRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// emailRequest is the payload for SendEmail.
type emailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// GetImage generates and returns an image captcha.
// GET /captcha/image?key=xxx
func (h *CaptchaHandler) GetImage(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
		return
	}

	imgBytes, err := h.captchaService.GenerateImage(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate captcha"})
		return
	}

	c.Data(http.StatusOK, "image/png", imgBytes)
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
		"message": "verification code sent",
	})
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
		"message": "verification code sent",
	})
}
