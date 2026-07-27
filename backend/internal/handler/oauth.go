package handler

import (
	"net/http"

	"github.com/anchor-finance/backend/internal/api/middleware"
	"github.com/anchor-finance/backend/internal/service"
	"github.com/anchor-finance/backend/pkg/logger"
	"github.com/anchor-finance/backend/pkg/oauth"
	"github.com/anchor-finance/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// OAuthHandler handles OAuth-related HTTP requests.
type OAuthHandler struct {
	oauthSvc *service.OAuthService
	log      *logger.Logger
	jwtKey   []byte
}

// NewOAuthHandler creates a new OAuth handler.
func NewOAuthHandler(oauthSvc *service.OAuthService, log *logger.Logger, jwtKey string) *OAuthHandler {
	return &OAuthHandler{
		oauthSvc: oauthSvc,
		log:      log,
		jwtKey:   []byte(jwtKey),
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

	provider, ok := oauth.Get(providerName)
	if !ok {
		response.BadRequest(c, "unsupported oauth provider")
		return
	}

	redirect := c.Query("redirect")

	state, err := h.oauthSvc.GenerateState(redirect)
	if err != nil {
		h.log.Errorf("generate oauth state failed: %v", err)
		response.ServerError(c, "failed to generate state")
		return
	}

	authURL := provider.GetAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Callback handles the OAuth provider's callback with the authorization code.
func (h *OAuthHandler) Callback(c *gin.Context) {
	providerName := c.Param("provider")

	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		response.BadRequest(c, "missing authorization code")
		return
	}

	// Validate state
	redirect, err := h.oauthSvc.ValidateState(state)
	if err != nil {
		response.BadRequest(c, "invalid or expired state")
		return
	}

	// Process OAuth callback
	result, err := h.oauthSvc.HandleCallback(providerName, code)
	if err != nil {
		h.log.Errorf("oauth callback failed: provider=%s err=%v", providerName, err)
		response.Error(c, http.StatusUnauthorized, 401, err.Error())
		return
	}

	// Find the user to generate tokens
	user, err := h.oauthSvc.FindUserByOAuth(providerName, "")
	if err != nil || user == nil {
		// If we can't find user by OAuth, the callback created one - find by the result
		response.ServerError(c, "failed to retrieve user")
		return
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := h.generateTokens(user.ID, user.IsAdmin)
	if err != nil {
		response.ServerError(c, "failed to generate tokens")
		return
	}

	// If there's a redirect URL, append tokens as query params
	if redirect != "" {
		sep := "?"
		if containsChar(redirect, '?') {
			sep = "&"
		}
		c.Redirect(http.StatusTemporaryRedirect, redirect+sep+
			"access_token="+accessToken+
			"&refresh_token="+refreshToken+
			"&is_new_user="+boolStr(result.IsNewUser))
		return
	}

	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"is_new_user":   result.IsNewUser,
	})
}

// BindAccount binds an OAuth account to the currently authenticated user.
func (h *OAuthHandler) BindAccount(c *gin.Context) {
	var req struct {
		Provider string `json:"provider" binding:"required"`
		Code     string `json:"code" binding:"required"`
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

	if err := h.oauthSvc.BindAccount(userID.(uint), req.Provider, req.Code); err != nil {
		response.Error(c, http.StatusConflict, 409, err.Error())
		return
	}

	response.SuccessMsg(c, "oauth account bound successfully")
}

// GetBoundAccounts returns the OAuth accounts bound to the current user.
func (h *OAuthHandler) GetBoundAccounts(c *gin.Context) {
	userID, exists := c.Get(middleware.ContextKeyUserID)
	if !exists {
		response.Unauthorized(c, "not authenticated")
		return
	}

	accounts, err := h.oauthSvc.GetBoundAccounts(userID.(uint))
	if err != nil {
		response.ServerError(c, "failed to get bound accounts")
		return
	}

	type accountInfo struct {
		ID       uint   `json:"id"`
		Provider string `json:"provider"`
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}

	result := make([]accountInfo, 0, len(accounts))
	for _, a := range accounts {
		result = append(result, accountInfo{
			ID:       a.ID,
			Provider: a.Provider,
			Username: a.Username,
			Avatar:   a.Avatar,
		})
	}

	response.Success(c, gin.H{"accounts": result})
}

// UnbindAccount unbinds an OAuth account from the current user.
func (h *OAuthHandler) UnbindAccount(c *gin.Context) {
	var req struct {
		Provider string `json:"provider" binding:"required"`
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

	if err := h.oauthSvc.UnbindAccount(userID.(uint), req.Provider); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "oauth account unbound successfully")
}

func (h *OAuthHandler) generateTokens(userID uint, isAdmin bool) (string, string, error) {
	claims := jwt.MapClaims{
		"user_id":   userID,
		"is_admin":  isAdmin,
	}

	// Generate access token (same pattern as AuthHandler)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessStr, err := accessToken.SignedString(h.jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshStr, err := refreshToken.SignedString(h.jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessStr, refreshStr, nil
}

func containsChar(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
