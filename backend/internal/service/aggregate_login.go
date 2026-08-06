package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/auth"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/oauth"
	"anchorfinance/pkg/oauth/juhe"

	"gorm.io/gorm"
)

// AggregateLoginService handles aggregate login (聚合登录) business logic.
type AggregateLoginService struct {
	db        *gorm.DB
	log       *logger.Logger
	userSvc   *UserService
	jwtMgr    *auth.JWTManager
	baseURL   string // backend base URL for constructing redirect URIs
}

// NewAggregateLoginService creates a new aggregate login service.
func NewAggregateLoginService(
	db *gorm.DB,
	log *logger.Logger,
	userSvc *UserService,
	jwtMgr *auth.JWTManager,
	baseURL string,
) *AggregateLoginService {
	return &AggregateLoginService{
		db:      db,
		log:     log,
		userSvc: userSvc,
		jwtMgr:  jwtMgr,
		baseURL: baseURL,
	}
}

// GetProviders returns all active aggregate login providers.
func (s *AggregateLoginService) GetProviders() ([]model.AggregateLoginProvider, error) {
	var providers []model.AggregateLoginProvider
	if err := s.db.Where("is_active = ?", true).Order("sort_order ASC").Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}

// GetAuthURL generates the authorization URL for a given provider code.
func (s *AggregateLoginService) GetAuthURL(providerCode, state string) (string, error) {
	provider, err := s.getProviderByCode(providerCode)
	if err != nil {
		return "", err
	}

	juheProvider := s.createJuheProvider(provider)
	return juheProvider.GetAuthURL(state), nil
}

// HandleCallback processes an aggregate login callback:
// validates code via aggregation API, creates/finds user, and generates JWT token.
func (s *AggregateLoginService) HandleCallback(providerCode, code, ip, userAgent string) (*User, string, bool, error) {
	provider, err := s.getProviderByCode(providerCode)
	if err != nil {
		return nil, "", false, err
	}

	juheProvider := s.createJuheProvider(provider)

	// Exchange code for user info via aggregation API
	userInfo, err := juheProvider.GetUserInfo(code)
	if err != nil {
		s.logAggregateLogin(provider.ID, 0, userInfo, ip, userAgent, 0, err.Error())
		return nil, "", false, err
	}

	// Find existing OAuth account by provider + social_uid
	var oauthAccount model.OAuthAccount
	oauthProviderName := juheProvider.Name()
	err = s.db.Where("provider = ? AND open_id = ?", oauthProviderName, userInfo.OpenID).First(&oauthAccount).Error

	var user *User
	var isNewUser bool

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// New user - create account and OAuth binding
		user, err = s.createOAuthUser(oauthProviderName, userInfo)
		if err != nil {
			s.logAggregateLogin(provider.ID, 0, userInfo, ip, userAgent, 0, err.Error())
			return nil, "", false, err
		}
		isNewUser = true
	} else if err != nil {
		s.logAggregateLogin(provider.ID, 0, userInfo, ip, userAgent, 0, err.Error())
		return nil, "", false, err
	} else {
		// Existing user - update OAuth account info
		s.updateOAuthAccount(&oauthAccount, userInfo)

		user, err = s.userSvc.GetByID(oauthAccount.UserID)
		if err != nil {
			return nil, "", false, errors.New("user account not found")
		}

		if user.Status != 1 {
			return nil, "", false, errors.New("account disabled")
		}

		// Update last login
		now := time.Now()
		s.db.Model(user).Update("last_login", &now)
		isNewUser = false
	}

	// Generate JWT token
	token, err := s.jwtMgr.GenerateToken(user.ID, user.Role == "admin")
	if err != nil {
		return nil, "", false, errors.New("failed to generate token")
	}

	// Log successful login
	s.logAggregateLogin(provider.ID, user.ID, userInfo, ip, userAgent, 1, "")

	s.log.Infof("aggregate login: provider=%s social_uid=%s user_id=%d new=%v",
		providerCode, userInfo.OpenID, user.ID, isNewUser)

	return user, token, isNewUser, nil
}

// BindAccount binds an aggregate login account to an existing user.
func (s *AggregateLoginService) BindAccount(userID uint, providerCode, code string) error {
	provider, err := s.getProviderByCode(providerCode)
	if err != nil {
		return err
	}

	// Verify user exists
	_, err = s.userSvc.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	juheProvider := s.createJuheProvider(provider)

	// Exchange code for user info
	userInfo, err := juheProvider.GetUserInfo(code)
	if err != nil {
		return err
	}

	oauthProviderName := juheProvider.Name()

	// Check if already bound to another user
	var existing model.OAuthAccount
	err = s.db.Where("provider = ? AND open_id = ?", oauthProviderName, userInfo.OpenID).First(&existing).Error
	if err == nil {
		if existing.UserID == userID {
			return errors.New("already bound to this account")
		}
		return errors.New("this social account is already bound to another user")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create OAuth binding
	rawDataJSON, _ := json.Marshal(userInfo.RawData)
	oauthAccount := model.OAuthAccount{
		UserID:   userID,
		Provider: oauthProviderName,
		OpenID:   userInfo.OpenID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Avatar:   userInfo.Avatar,
		RawData:  rawDataJSON,
	}
	if err := s.db.Create(&oauthAccount).Error; err != nil {
		return err
	}

	s.log.Infof("aggregate account bound: provider=%s social_uid=%s user_id=%d",
		providerCode, userInfo.OpenID, userID)
	return nil
}

// GetBoundAccounts returns all aggregate login accounts bound to a user.
func (s *AggregateLoginService) GetBoundAccounts(userID uint) ([]model.OAuthAccount, error) {
	var accounts []model.OAuthAccount
	if err := s.db.Where("user_id = ? AND provider LIKE ?", userID, "juhe_%").Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

// UnbindAccount removes an aggregate login binding from a user.
func (s *AggregateLoginService) UnbindAccount(userID uint, providerCode string) error {
	provider, err := s.getProviderByCode(providerCode)
	if err != nil {
		return err
	}

	oauthProviderName := "juhe_" + provider.Code

	result := s.db.Where("user_id = ? AND provider = ?", userID, oauthProviderName).Delete(&model.OAuthAccount{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("aggregate account not bound")
	}

	s.log.Infof("aggregate account unbound: provider=%s user_id=%d", providerCode, userID)
	return nil
}

// getProviderByCode fetches an active provider config from the database.
func (s *AggregateLoginService) getProviderByCode(code string) (*model.AggregateLoginProvider, error) {
	var provider model.AggregateLoginProvider
	err := s.db.Where("code = ? AND is_active = ?", code, true).First(&provider).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("aggregate login provider not found or disabled")
	}
	if err != nil {
		return nil, err
	}
	return &provider, nil
}

// createJuheProvider creates a juhe.Provider from the database config.
func (s *AggregateLoginService) createJuheProvider(p *model.AggregateLoginProvider) *juhe.Provider {
	redirectURL := s.baseURL + "/api/v1/auth/aggregate/" + p.Code + "/callback"
	return juhe.New(p.APIURL, p.AppID, p.AppKey, p.Code, redirectURL)
}

// createOAuthUser creates a new user from OAuth info and binds the account.
func (s *AggregateLoginService) createOAuthUser(providerName string, userInfo *oauth.UserInfo) (*User, error) {
	username := s.generateUsername(userInfo)

	user, err := s.userSvc.Register(RegisterRequest{
		Username: username,
		Password: generateRandomPassword(),
		Email:    userInfo.Email,
		Nickname: userInfo.Username,
	})
	if err != nil {
		s.log.Errorf("create aggregate user failed: %v", err)
		return nil, errors.New("failed to create user account")
	}

	// Update avatar if available
	if userInfo.Avatar != "" {
		s.db.Model(user).Update("avatar", userInfo.Avatar)
	}

	// Create OAuth account binding
	rawDataJSON, _ := json.Marshal(userInfo.RawData)
	oauthAccount := model.OAuthAccount{
		UserID:   user.ID,
		Provider: providerName,
		OpenID:   userInfo.OpenID,
		Username: userInfo.Username,
		Email:    userInfo.Email,
		Avatar:   userInfo.Avatar,
		RawData:  rawDataJSON,
	}
	if err := s.db.Create(&oauthAccount).Error; err != nil {
		s.log.Errorf("create aggregate oauth account failed: %v", err)
	}

	return user, nil
}

// updateOAuthAccount updates an existing OAuth account with fresh info.
func (s *AggregateLoginService) updateOAuthAccount(account *model.OAuthAccount, userInfo *oauth.UserInfo) {
	updates := map[string]interface{}{
		"username": userInfo.Username,
		"avatar":   userInfo.Avatar,
	}
	if userInfo.Email != "" {
		updates["email"] = userInfo.Email
	}
	if userInfo.RawData != nil {
		rawDataJSON, _ := json.Marshal(userInfo.RawData)
		updates["raw_data"] = rawDataJSON
	}
	s.db.Model(account).Updates(updates)
}

// generateUsername creates a unique username from OAuth info.
func (s *AggregateLoginService) generateUsername(userInfo *oauth.UserInfo) string {
	base := userInfo.Username
	if base == "" {
		base = userInfo.Provider + "_user"
	}

	var count int64
	s.db.Model(&User{}).Where("username = ?", base).Count(&count)
	if count == 0 {
		return base
	}

	b := make([]byte, 4)
	rand.Read(b)
	return base + "_" + hex.EncodeToString(b)
}

// logAggregateLogin records an aggregate login attempt in the log table.
func (s *AggregateLoginService) logAggregateLogin(
	providerID, userID uint,
	userInfo *oauth.UserInfo,
	ip, userAgent string,
	status int16,
	remark string,
) {
	loginLog := model.AggregateLoginLog{
		UserID:     userID,
		ProviderID: providerID,
		IP:         ip,
		UserAgent:  userAgent,
		Status:     status,
		Remark:     remark,
	}
	if userInfo != nil {
		loginLog.SocialUID = userInfo.OpenID
		loginLog.Nickname = userInfo.Username
		loginLog.Avatar = userInfo.Avatar
	}
	if err := s.db.Create(&loginLog).Error; err != nil {
		s.log.Errorf("create aggregate login log failed: %v", err)
	}
}
