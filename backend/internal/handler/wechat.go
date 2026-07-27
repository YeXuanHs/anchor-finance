package handler

import (
	"io"
	"net/http"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// WechatHandler handles WeChat-related HTTP requests.
type WechatHandler struct {
	wechatSvc *service.WechatService
	log       *logger.Logger
}

// NewWechatHandler creates a new WechatHandler.
func NewWechatHandler(wechatSvc *service.WechatService, log *logger.Logger) *WechatHandler {
	return &WechatHandler{wechatSvc: wechatSvc, log: log}
}

// GetAuthURL returns the WeChat OAuth authorization URL.
// GET /wechat/auth
func (h *WechatHandler) GetAuthURL(c *gin.Context) {
	redirectURI := c.Query("redirect_uri")
	state := c.Query("state")
	if state == "" {
		state = "wechat_login"
	}

	authURL := h.wechatSvc.GetAuthURL(redirectURI, state)
	response.Success(c, gin.H{"auth_url": authURL})
}

// Login handles WeChat OAuth login callback.
// GET /wechat/callback
func (h *WechatHandler) Login(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.BadRequest(c, "missing authorization code")
		return
	}

	result, err := h.wechatSvc.HandleLogin(code, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		h.log.Errorf("wechat login failed: %v", err)
		response.Error(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	response.Success(c, gin.H{
		"openid":     result.OpenID,
		"is_new_user": result.IsNewUser,
	})
}

// PayNotify handles WeChat payment callback.
// POST /wechat/pay/notify
func (h *WechatHandler) PayNotify(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Data(http.StatusOK, "application/xml", h.wechatSvc.GeneratePayNotifyResponse())
		return
	}

	_, err = h.wechatSvc.HandlePayNotify(body)
	if err != nil {
		h.log.Errorf("wechat pay notify failed: %v", err)
		c.Data(http.StatusOK, "application/xml", h.wechatSvc.GeneratePayNotifyResponse())
		return
	}

	c.Data(http.StatusOK, "application/xml", h.wechatSvc.GeneratePayNotifyResponse())
}

// SendTemplateMessage sends a WeChat template message.
// POST /wechat/message/send
func (h *WechatHandler) SendTemplateMessage(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		TemplateID string            `json:"template_id"`
		Data       map[string]string `json:"data" binding:"required"`
		URL        string            `json:"url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Get user's WeChat openid
	var oauthAccount struct {
		OpenID string `gorm:"column:open_id"`
	}
	if err := h.wechatSvc.GetDB().Table("oauth_accounts").
		Where("user_id = ? AND provider = 'wechat'", userID).
		Select("open_id").First(&oauthAccount).Error; err != nil {
		response.BadRequest(c, "wechat account not bound")
		return
	}

	if err := h.wechatSvc.SendTemplateMessage(oauthAccount.OpenID, req.TemplateID, req.Data, req.URL); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "template message sent")
}

// CreatePayOrder creates a WeChat pay order.
// POST /wechat/pay/create
func (h *WechatHandler) CreatePayOrder(c *gin.Context) {
	var req struct {
		OrderNo string `json:"order_no" binding:"required"`
		Amount  int    `json:"amount" binding:"required,gt=0"`
		OpenID  string `json:"openid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.wechatSvc.CreatePayOrder(req.OrderNo, req.Amount, c.ClientIP(), req.OpenID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}
