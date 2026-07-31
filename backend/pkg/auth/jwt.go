package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenExpired = errors.New("token has expired")
	ErrTokenInvalid = errors.New("token is invalid")
)

// Claims defines the JWT claims structure.
type Claims struct {
	UserID  uint   `json:"user_id"`
	IsAdmin bool   `json:"is_admin"`
	IP      string `json:"ip,omitempty"` // IP 绑定（移植自 zjmf home_ip_check）
	jwt.RegisteredClaims
}

// JWTManager handles JWT token generation and validation.
type JWTManager struct {
	secret     []byte
	expireHour int
}

// NewJWTManager creates a new JWT manager.
func NewJWTManager(secret string, expireHour int) *JWTManager {
	if expireHour <= 0 {
		expireHour = 72
	}
	return &JWTManager{
		secret:     []byte(secret),
		expireHour: expireHour,
	}
}

// GenerateToken creates a new JWT token for a user.
func (m *JWTManager) GenerateToken(userID uint, isAdmin bool) (string, error) {
	return m.GenerateTokenWithIP(userID, isAdmin, "")
}

// GenerateTokenWithIP creates a new JWT token for a user with IP binding.
func (m *JWTManager) GenerateTokenWithIP(userID uint, isAdmin bool, ip string) (string, error) {
	claims := Claims{
		UserID:  userID,
		IsAdmin: isAdmin,
		IP:      ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(m.expireHour) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "anchor-finance",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ValidateToken parses and validates a JWT token string.
func (m *JWTManager) ValidateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenInvalid
		}
		return m.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}
