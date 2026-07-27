package handler

import (
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// OAuthBindHandler handles OAuth account binding HTTP requests.
type OAuthBindHandler struct {
	oauthBindSvc *service.OAuthBindService
	log          *logger.Logger
}

// NewOAuthBindHandler creates a new OAuthBindHandler.
func NewOAuthBindHandler(oauthBindSvc *service.OAuthBindService, log *logger.Logger) *OAuthBindHandler {
	return &OAuthBindHandler{oauthBindSvc: oauthBindSvc, log: log}
}

// GetProviders returns available OAuth providers.
// GET /oauth-bind/providers
func (h *OAuthBindHandler) GetProviders(c *gin.Context) {
	providers := h.oauthBindSvc.GetProviders()
	response.Success(c, gin.H{"providers": providers})
}

// Bind binds an OAuth account to the current user.
// POST /oauth-bind/bind
func (h *OAuthBindHandler) Bind(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.BindOAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.oauthBindSvc.Bind(userID, req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "oauth account bound successfully")
}

// Unbind unbinds an OAuth account from the current user.
// POST /oauth-bind/unbind
func (h *OAuthBindHandler) Unbind(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Provider string `json:"provider" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.oauthBindSvc.Unbind(userID, req.Provider); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "oauth account unbound successfully")
}

// GetBoundAccounts returns OAuth accounts bound to the current user.
// GET /oauth-bind/accounts
func (h *OAuthBindHandler) GetBoundAccounts(c *gin.Context) {
	userID := c.GetUint("user_id")

	accounts, err := h.oauthBindSvc.GetBoundAccounts(userID)
	if err != nil {
		response.ServerError(c, "failed to get bound accounts")
		return
	}
	response.Success(c, gin.H{"accounts": accounts})
}

// CheckBinding checks if a specific provider is bound.
// GET /oauth-bind/check/:provider
func (h *OAuthBindHandler) CheckBinding(c *gin.Context) {
	userID := c.GetUint("user_id")
	provider := c.Param("provider")

	isBound := h.oauthBindSvc.IsBound(userID, provider)
	response.Success(c, gin.H{
		"provider": provider,
		"bound":    isBound,
	})
}
