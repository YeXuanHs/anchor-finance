package handler

import (
	"net/http"

	"anchorfinance/internal/api/middleware"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AggregateLoginHandler handles aggregate login (聚合登录) HTTP requests.
type AggregateLoginHandler struct {
	aggregateSvc *service.AggregateLoginService
	oauthSvc     *service.OAuthService
	userSvc      *service.UserService
	log          *logger.Logger
	jwtKey       []byte
}

// NewAggregateLoginHandler creates a new aggregate login handler.
func NewAggregateLoginHandler(
	aggregateSvc *service.AggregateLoginService,
	oauthSvc *service.OAuthService,
	userSvc *service.UserService,
	log *logger.Logger,
	jwtKey string,
) *AggregateLoginHandler {
	return &AggregateLoginHandler{
		aggregateSvc: aggregateSvc,
		oauthSvc:     oauthSvc,
		userSvc:      userSvc,
		log:          log,
		jwtKey:       []byte(jwtKey),
	}
}

// GetProviders returns a list of active aggregate login providers.
// GET /auth/aggregate/providers
func (h *AggregateLoginHandler) GetProviders(c *gin.Context) {
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

// Login redirects the user to the aggregated login authorization page.
// GET /auth/aggregate/:code
func (h *AggregateLoginHandler) Login(c *gin.Context) {
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

// Callback handles the aggregated login callback with the authorization code.
// GET /auth/aggregate/:code/callback
func (h *AggregateLoginHandler) Callback(c *gin.Context) {
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

	_ = user // user is available for further use if needed

	// If there's a redirect URL, append tokens as query params
	if redirect != "" {
		sep := "?"
		if containsChar(redirect, '?') {
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

// BindAccount binds an aggregate login account to the currently authenticated user.
// POST /auth/aggregate/bind
func (h *AggregateLoginHandler) BindAccount(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"` // provider code, e.g. qq, wx
		OAuthCode string `json:"oauth_code" binding:"required"` // authorization code from callback
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

// GetBoundAccounts returns the aggregate login accounts bound to the current user.
// GET /auth/aggregate/accounts
func (h *AggregateLoginHandler) GetBoundAccounts(c *gin.Context) {
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
		ID         uint   `json:"id"`
		Provider   string `json:"provider"`
		Username   string `json:"username"`
		Avatar     string `json:"avatar"`
		SocialUID  string `json:"social_uid"`
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

// UnbindAccount unbinds an aggregate login account from the current user.
// POST /auth/aggregate/unbind
func (h *AggregateLoginHandler) UnbindAccount(c *gin.Context) {
	var req struct {
		Code string `json:"code" binding:"required"` // provider code
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

func (h *AggregateLoginHandler) generateToken(userID uint, isAdmin bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"is_admin": isAdmin,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(h.jwtKey)
}
