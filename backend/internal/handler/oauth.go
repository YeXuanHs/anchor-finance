package handler

import (
	"net/http"
	"strings"

	"anchorfinance/internal/api/middleware"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/auth"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/oauth"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// OAuthHandler handles OAuth-related HTTP requests (including aggregate login).
type OAuthHandler struct {
	oauthSvc     *service.OAuthService
	aggregateSvc *service.AggregateLoginService
	userSvc      *service.UserService
	log          *logger.Logger
	jwtMgr       *auth.JWTManager
}

// NewOAuthHandler creates a new OAuth handler.
func NewOAuthHandler(
	oauthSvc *service.OAuthService,
	log *logger.Logger,
	jwtMgr *auth.JWTManager,
	aggregateSvc *service.AggregateLoginService,
	userSvc *service.UserService,
) *OAuthHandler {
	return &OAuthHandler{
		oauthSvc:     oauthSvc,
		log:          log,
		jwtMgr:       jwtMgr,
		aggregateSvc: aggregateSvc,
		userSvc:      userSvc,
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

// ==================== Aggregate Login (聚合登录) ====================

// GetAggregateProviders returns a list of active aggregate login providers.
// GET /auth/aggregate/providers
func (h *OAuthHandler) GetAggregateProviders(c *gin.Context) {
	if h.aggregateSvc == nil {
		response.ServerError(c, "aggregate login service not available")
		return
	}

	providers, err := h.aggregateSvc.GetProviders()
	if err != nil {
		h.log.Errorf("get aggregate providers failed: %v", err)
		response.ServerError(c, "failed to get providers")
		return
	}

	type providerInfo struct {
		Code string `json:"code"`
		Name string `json:"name"`
		Icon string `json:"icon,omitempty"`
	}

	result := make([]providerInfo, 0, len(providers))
	for _, p := range providers {
		result = append(result, providerInfo{
			Code: p.Code,
			Name: p.Name,
		})
	}

	response.Success(c, gin.H{"providers": result})
}

// AggregateLogin redirects the user to the aggregated login authorization page.
// GET /auth/aggregate/:code
func (h *OAuthHandler) AggregateLogin(c *gin.Context) {
	if h.aggregateSvc == nil {
		response.ServerError(c, "aggregate login service not available")
		return
	}

	code := c.Param("code")
	redirect := c.Query("redirect")

	state, err := h.oauthSvc.GenerateState(redirect)
	if err != nil {
		h.log.Errorf("generate aggregate login state failed: %v", err)
		response.ServerError(c, "failed to generate state")
		return
	}

	authURL, err := h.aggregateSvc.GetAuthURL(code, state)
	if err != nil {
		h.log.Errorf("get aggregate auth url failed: provider=%s err=%v", code, err)
		response.BadRequest(c, err.Error())
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// AggregateCallback handles the aggregated login callback with the authorization code.
// GET /auth/aggregate/:code/callback
func (h *OAuthHandler) AggregateCallback(c *gin.Context) {
	if h.aggregateSvc == nil {
		response.ServerError(c, "aggregate login service not available")
		return
	}

	providerCode := c.Param("code")
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		response.BadRequest(c, "missing authorization code")
		return
	}

	// Validate state to prevent CSRF
	redirect, err := h.oauthSvc.ValidateState(state)
	if err != nil {
		response.BadRequest(c, "invalid or expired state")
		return
	}

	// Process the aggregate login callback
	user, token, isNewUser, err := h.aggregateSvc.HandleCallback(providerCode, code, c.ClientIP(), c.GetHeader("User-Agent"))
	if err != nil {
		h.log.Errorf("aggregate login callback failed: provider=%s err=%v", providerCode, err)
		response.Error(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	_ = user

	// If there's a redirect URL, append tokens as query params
	if redirect != "" {
		sep := "?"
		if strings.Contains(redirect, "?") {
			sep = "&"
		}
		isNew := "false"
		if isNewUser {
			isNew = "true"
		}
		c.Redirect(http.StatusTemporaryRedirect, redirect+sep+
			"access_token="+token+
			"&is_new_user="+isNew)
		return
	}

	response.Success(c, gin.H{
		"access_token": token,
		"is_new_user":  isNewUser,
	})
}

// BindAggregateAccount binds an aggregate login account to the currently authenticated user.
// POST /auth/aggregate/bind
func (h *OAuthHandler) BindAggregateAccount(c *gin.Context) {
	if h.aggregateSvc == nil {
		response.ServerError(c, "aggregate login service not available")
		return
	}

	var req struct {
		Code      string `json:"code" binding:"required"`
		OAuthCode string `json:"oauth_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get(middleware.ContextKeyUserID)
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	if err := h.aggregateSvc.BindAccount(userID.(uint), req.Code, req.OAuthCode); err != nil {
		response.Error(c, http.StatusConflict, 409, err.Error())
		return
	}

	response.SuccessMsg(c, "aggregate account bound successfully")
}

// GetBoundAggregateAccounts returns the aggregate login accounts bound to the current user.
// GET /auth/aggregate/accounts
func (h *OAuthHandler) GetBoundAggregateAccounts(c *gin.Context) {
	if h.aggregateSvc == nil {
		response.ServerError(c, "aggregate login service not available")
		return
	}

	userID, exists := c.Get(middleware.ContextKeyUserID)
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	accounts, err := h.aggregateSvc.GetBoundAccounts(userID.(uint))
	if err != nil {
		response.ServerError(c, "failed to get bound accounts")
		return
	}

	type accountInfo struct {
		ID        uint   `json:"id"`
		Provider  string `json:"provider"`
		Username  string `json:"username"`
		Avatar    string `json:"avatar"`
		SocialUID string `json:"social_uid"`
	}

	result := make([]accountInfo, 0, len(accounts))
	for _, a := range accounts {
		result = append(result, accountInfo{
			ID:        a.ID,
			Provider:  a.Provider,
			Username:  a.Username,
			Avatar:    a.Avatar,
			SocialUID: a.OpenID,
		})
	}

	response.Success(c, gin.H{"accounts": result})
}

// UnbindAggregateAccount unbinds an aggregate login account from the current user.
// POST /auth/aggregate/unbind
func (h *OAuthHandler) UnbindAggregateAccount(c *gin.Context) {
	if h.aggregateSvc == nil {
		response.ServerError(c, "aggregate login service not available")
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID, exists := c.Get(middleware.ContextKeyUserID)
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	if err := h.aggregateSvc.UnbindAccount(userID.(uint), req.Code); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "aggregate account unbound successfully")
}
