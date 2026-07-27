package handler

import (
	"net/http"

	"anchorfinance/internal/service"
	"github.com/gin-gonic/gin"
)

// CaptchaHandler handles captcha-related HTTP requests.
type CaptchaHandler struct {
	captchaService *service.CaptchaService
}

// NewCaptchaHandler creates a new CaptchaHandler.
func NewCaptchaHandler(captchaService *service.CaptchaService) *CaptchaHandler {
	return &CaptchaHandler{captchaService: captchaService}
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

	// TODO: Integrate with SMS gateway to send the code.
	// The code is returned here for development/testing purposes only.
	// In production, remove the code from the response.
	_ = code // In production, send via SMS gateway

	c.JSON(http.StatusOK, gin.H{
		"message": "verification code sent",
		// "code": code, // Uncomment for development only
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

	// TODO: Integrate with email service to send the code.
	// The code is returned here for development/testing purposes only.
	// In production, remove the code from the response.
	_ = code // In production, send via email service

	c.JSON(http.StatusOK, gin.H{
		"message": "verification code sent",
		// "code": code, // Uncomment for development only
	})
}
