package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/anchor-finance/backend/internal/model"
	"github.com/anchor-finance/backend/pkg/logger"
	"github.com/anchor-finance/backend/pkg/oauth"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// OAuthService handles OAuth business logic.
type OAuthService struct {
	db        *gorm.DB
	log       *logger.Logger
	userSvc   *UserService
	frontendURL string
}

// NewOAuthService creates a new OAuth service.
func NewOAuthService(db *gorm.DB, log *logger.Logger, userSvc *UserService, frontendURL string) *OAuthService {
	return &OAuthService{
		db:          db,
		log:         log,
		userSvc:     userSvc,
		frontendURL: frontendURL,
	}
}

// OAuthLoginResult contains the result of an OAuth login.
type OAuthLoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IsNewUser    bool   `json:"is_new_user"`
}

// GenerateState creates a random state string and stores it.
func (s *OAuthService) GenerateState(redirect string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	state := hex.EncodeToString(b)

	oauthState := model.OAuthState{
		State:     state,
		Redirect:  redirect,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	if err := s.db.Create(&oauthState).Error; err != nil {
		return "", err
	}

	return state, nil
}

// ValidateState validates and consumes an OAuth state token.
func (s *OAuthService) ValidateState(state string) (string, error) {
	var oauthState model.OAuthState
	err := s.db.Where("state = ? AND expires_at > ?", state, time.Now()).First(&oauthState).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", errors.New("invalid or expired state")
	}
	if err != nil {
		return "", err
	}

	// Delete used state
	s.db.Delete(&oauthState)

	return oauthState.Redirect, nil
}

// HandleCallback processes an OAuth callback: validates code, fetches user info,
// creates a new user if needed, and returns login tokens.
func (s *OAuthService) HandleCallback(providerName, code string) (*OAuthLoginResult, error) {
	provider, ok := oauth.Get(providerName)
	if !ok {
		return nil, errors.New("unsupported oauth provider")
	}

	// Get user info from OAuth provider
	userInfo, err := provider.GetUserInfo(code)
	if err != nil {
		s.log.Errorf("oauth %s get user info failed: %v", providerName, err)
		return nil, err
	}

	// Check if OAuth account already exists
	var oauthAccount model.OAuthAccount
	err = s.db.Where("provider = ? AND open_id = ?", providerName, userInfo.OpenID).First(&oauthAccount).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// New OAuth user - create account
		return s.createOAuthUser(providerName, userInfo)
	}
	if err != nil {
		return nil, err
	}

	// Existing user - update OAuth account info and return tokens
	s.updateOAuthAccount(&oauthAccount, userInfo)

	// Find the user
	user, err := s.userSvc.GetByID(oauthAccount.UserID)
	if err != nil {
		return nil, errors.New("user account not found")
	}

	if user.Status != 1 {
		return nil, errors.New("account disabled")
	}

	// Update last login
	now := time.Now()
	s.db.Model(user).Update("last_login", &now)

	s.log.Infof("oauth login: provider=%s openid=%s userid=%d", providerName, userInfo.OpenID, user.ID)

	return &OAuthLoginResult{
		IsNewUser: false,
	}, nil
}

// createOAuthUser creates a new user from OAuth info and binds the account.
func (s *OAuthService) createOAuthUser(providerName string, userInfo *oauth.UserInfo) (*OAuthLoginResult, error) {
	// Generate a unique username
	username := s.generateUsername(userInfo)

	// Create user with a random password (OAuth users don't use password login)
	user, err := s.userSvc.Register(RegisterRequest{
		Username: username,
		Password: generateRandomPassword(),
		Email:    userInfo.Email,
		Nickname: userInfo.Username,
	})
	if err != nil {
		s.log.Errorf("create oauth user failed: %v", err)
		return nil, errors.New("failed to create user account")
	}

	// Update avatar if available
	if userInfo.Avatar != "" {
		s.db.Model(user).Update("avatar", userInfo.Avatar)
	}

	// Create OAuth account binding
	rawDataJSON, _ := datatypes.MarshalJSON(userInfo.RawData)
	oauthAccount := model.OAuthAccount{
		UserID:   user.ID,
		Provider: providerName,
		OpenID:   userInfo.OpenID,
		UnionID:  userInfo.UnionID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Avatar:   userInfo.Avatar,
		RawData:  rawDataJSON,
	}
	if err := s.db.Create(&oauthAccount).Error; err != nil {
		s.log.Errorf("create oauth account failed: %v", err)
	}

	s.log.Infof("new oauth user created: provider=%s openid=%s userid=%d", providerName, userInfo.OpenID, user.ID)

	return &OAuthLoginResult{
		IsNewUser: true,
	}, nil
}

// BindAccount binds an OAuth account to an existing user.
func (s *OAuthService) BindAccount(userID uint, providerName, code string) error {
	provider, ok := oauth.Get(providerName)
	if !ok {
		return errors.New("unsupported oauth provider")
	}

	// Check if user exists
	_, err := s.userSvc.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Get user info from OAuth provider
	userInfo, err := provider.GetUserInfo(code)
	if err != nil {
		return err
	}

	// Check if this OAuth account is already bound to another user
	var existing model.OAuthAccount
	err = s.db.Where("provider = ? AND open_id = ?", providerName, userInfo.OpenID).First(&existing).Error
	if err == nil {
		if existing.UserID == userID {
			return errors.New("already bound to this account")
		}
		return errors.New("this oauth account is already bound to another user")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create OAuth binding
	rawDataJSON, _ := datatypes.MarshalJSON(userInfo.RawData)
	oauthAccount := model.OAuthAccount{
		UserID:   userID,
		Provider: providerName,
		OpenID:   userInfo.OpenID,
		UnionID:  userInfo.UnionID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Avatar:   userInfo.Avatar,
		RawData:  rawDataJSON,
	}
	if err := s.db.Create(&oauthAccount).Error; err != nil {
		return err
	}

	s.log.Infof("oauth account bound: provider=%s openid=%s userid=%d", providerName, userInfo.OpenID, userID)
	return nil
}

// GetBoundAccounts returns all OAuth accounts bound to a user.
func (s *OAuthService) GetBoundAccounts(userID uint) ([]model.OAuthAccount, error) {
	var accounts []model.OAuthAccount
	if err := s.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

// UnbindAccount removes an OAuth binding from a user.
func (s *OAuthService) UnbindAccount(userID uint, providerName string) error {
	result := s.db.Where("user_id = ? AND provider = ?", userID, providerName).Delete(&model.OAuthAccount{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("oauth account not bound")
	}

	s.log.Infof("oauth account unbound: provider=%s userid=%d", providerName, userID)
	return nil
}

// FindUserByOAuth finds a user by their OAuth account.
func (s *OAuthService) FindUserByOAuth(providerName, openID string) (*User, error) {
	var oauthAccount model.OAuthAccount
	err := s.db.Where("provider = ? AND open_id = ?", providerName, openID).First(&oauthAccount).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	user, err := s.userSvc.GetByID(oauthAccount.UserID)
	if err != nil {
		return nil, nil
	}
	return user, nil
}

func (s *OAuthService) updateOAuthAccount(account *model.OAuthAccount, userInfo *oauth.UserInfo) {
	updates := map[string]interface{}{
		"username": userInfo.Username,
		"avatar":   userInfo.Avatar,
	}
	if userInfo.Email != "" {
		updates["email"] = userInfo.Email
	}
	if userInfo.UnionID != "" {
		updates["union_id"] = userInfo.UnionID
	}
	if userInfo.RawData != nil {
		rawDataJSON, _ := datatypes.MarshalJSON(userInfo.RawData)
		updates["raw_data"] = rawDataJSON
	}
	s.db.Model(account).Updates(updates)
}

func (s *OAuthService) generateUsername(userInfo *oauth.UserInfo) string {
	base := userInfo.Username
	if base == "" {
		base = userInfo.Provider + "_user"
	}

	// Check if username exists
	var count int64
	s.db.Model(&User{}).Where("username = ?", base).Count(&count)
	if count == 0 {
		return base
	}

	// Append random suffix
	b := make([]byte, 4)
	rand.Read(b)
	return base + "_" + hex.EncodeToString(b)
}

func generateRandomPassword() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
