package handler

import (
	"net/http"

	"anchorfinance/internal/api/middleware"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/auth"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/oauth"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// OAuthHandler handles OAuth-related HTTP requests.
type OAuthHandler struct {
	oauthSvc *service.OAuthService
	log      *logger.Logger
	jwtMgr   *auth.JWTManager
}

// NewOAuthHandler creates a new OAuth handler.
func NewOAuthHandler(oauthSvc *service.OAuthService, log *logger.Logger, jwtMgr *auth.JWTManager) *OAuthHandler {
	return &OAuthHandler{
		oauthSvc: oauthSvc,
		log:      log,
		jwtMgr:   jwtMgr,
	}
}

// GetProviders returns a list of available OAuth providers.
func (h *OAuthHandler) GetProviders(c *gin.Context) {
	all := oauth.GetAll()
	providers := make([]gin.H, 0, len(all))
	for name := range all {
		providers = append(providers, gin.H{
			"name": name,
		})
	}
	response.Success(c, gin.H{"providers": providers})
}

// Login redirects the user to the OAuth provider's authorization page.
func (h *OAuthHandler) Login(c *gin.Context) {
	providerName := c.Param("provider")

	// Generate state for CSRF protection
	state, err := h.oauthSvc.GenerateState("")
	if err != nil {
		h.log.Errorf("OAuth state generation failed: %v", err)
		response.ServerError(c, "failed to generate state")
		return
	}

	// Get provider and login URL
	provider, ok := oauth.Get(providerName)
	if !ok {
		response.BadRequest(c, "unsupported oauth provider")
		return
	}

	url := provider.GetLoginURL(state)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"url":   url,
			"state": state,
		},
	})
}

// Callback handles the OAuth callback after user authorization.
func (h *OAuthHandler) Callback(c *gin.Context) {
	providerName := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		response.BadRequest(c, "missing code parameter")
		return
	}

	// Validate state for CSRF protection
	_, err := h.oauthSvc.ValidateState(state)
	if err != nil {
		response.BadRequest(c, "invalid or expired state")
		return
	}

	// Handle callback using the service (which handles code exchange, user creation, etc.)
	result, err := h.oauthSvc.HandleCallback(providerName, code)
	if err != nil {
		h.log.Errorf("OAuth callback failed: %v", err)
		response.BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": result,
	})
}

// Redirect redirects the user to the OAuth provider's authorization page.
func (h *OAuthHandler) Redirect(c *gin.Context) {
	providerName := c.Param("provider")

	// Generate state
	state, err := h.oauthSvc.GenerateState("")
	if err != nil {
		h.log.Errorf("OAuth state generation failed: %v", err)
		response.ServerError(c, "failed to generate state")
		return
	}

	provider, ok := oauth.Get(providerName)
	if !ok {
		response.BadRequest(c, "unsupported oauth provider")
		return
	}

	url := provider.GetLoginURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Unbind unbinds an OAuth account from the user.
func (h *OAuthHandler) Unbind(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Provider string `json:"provider" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.oauthSvc.UnbindAccount(userID, req.Provider); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "解绑成功")
}

// GetBoundAccounts returns the user's bound OAuth accounts.
func (h *OAuthHandler) GetBoundAccounts(c *gin.Context) {
	userID := c.GetUint("user_id")

	accounts, err := h.oauthSvc.GetBoundAccounts(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"accounts": accounts})
}

// BindAccount binds an OAuth account to the user.
func (h *OAuthHandler) BindAccount(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Provider string `json:"provider" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.oauthSvc.BindAccount(userID, req.Provider, req.Code); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "绑定成功")
}

// InvalidateUserTokens invalidates all tokens for a user (password change).
func InvalidateUserTokens(userID uint) {
	middleware.InvalidateUserTokens(userID)
}
