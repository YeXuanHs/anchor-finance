package handler

import (
	"net/http"
	"time"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	userSvc *service.UserService
	log     *logger.Logger
	jwtKey  []byte
}

func NewAuthHandler(userSvc *service.UserService, log *logger.Logger, jwtKey string) *AuthHandler {
	return &AuthHandler{
		userSvc: userSvc,
		log:     log,
		jwtKey:  []byte(jwtKey),
	}
}

type loginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
}

type tokenClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// Login authenticates and returns JWT tokens.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.userSvc.Login(service.LoginRequest{
		Account:  req.Account,
		Password: req.Password,
	})
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	accessToken, refreshToken, err := h.generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		response.ServerError(c, "failed to generate tokens")
		return
	}

	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

// Register creates a new account and returns tokens.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.userSvc.Register(service.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Phone:    req.Phone,
		Nickname: req.Nickname,
	})
	if err != nil {
		response.Error(c, http.StatusConflict, 409, err.Error())
		return
	}

	accessToken, refreshToken, err := h.generateTokens(user.ID, user.Username, user.Role)
	if err != nil {
		response.ServerError(c, "failed to generate tokens")
		return
	}

	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          user,
	})
}

// RefreshToken issues a new access token from a valid refresh token.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	claims, err := h.parseToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "invalid refresh token")
		return
	}

	accessToken, refreshToken, err := h.generateTokens(claims.UserID, claims.Username, claims.Role)
	if err != nil {
		response.ServerError(c, "failed to generate tokens")
		return
	}

	response.Success(c, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Logout is a client-side logout placeholder (token blacklisting can be added).
func (h *AuthHandler) Logout(c *gin.Context) {
	response.SuccessMsg(c, "logged out")
}

func (h *AuthHandler) generateTokens(userID uint, username, role string) (string, string, error) {
	accessClaims := tokenClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessStr, err := accessToken.SignedString(h.jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshClaims := tokenClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshStr, err := refreshToken.SignedString(h.jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessStr, refreshStr, nil
}

func (h *AuthHandler) parseToken(tokenStr string) (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &tokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return h.jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrSignatureInvalid
}
