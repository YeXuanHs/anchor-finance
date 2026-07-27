package service

import (
	"errors"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/oauth"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// OAuthBindService manages OAuth account binding/unbinding for the user frontend.
type OAuthBindService struct {
	db      *gorm.DB
	log     *logger.Logger
	userSvc *UserService
}

func NewOAuthBindService(db *gorm.DB, log *logger.Logger, userSvc *UserService) *OAuthBindService {
	return &OAuthBindService{db: db, log: log, userSvc: userSvc}
}

type BindOAuthRequest struct {
	Provider string `json:"provider" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

type OAuthAccountInfo struct {
	ID       uint   `json:"id"`
	Provider string `json:"provider"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
}

// GetProviders returns available OAuth providers.
func (s *OAuthBindService) GetProviders() []map[string]interface{} {
	all := oauth.GetAll()
	providers := make([]map[string]interface{}, 0, len(all))
	for name, p := range all {
		providers = append(providers, map[string]interface{}{
			"name": name,
			"auth_url": p.GetAuthURL(""),
		})
	}
	return providers
}

// Bind binds an OAuth account to the current user.
func (s *OAuthBindService) Bind(userID uint, req BindOAuthRequest) error {
	provider, ok := oauth.Get(req.Provider)
	if !ok {
		return errors.New("unsupported oauth provider")
	}

	// Verify user exists
	_, err := s.userSvc.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Exchange code for user info
	userInfo, err := provider.GetUserInfo(req.Code)
	if err != nil {
		s.log.Errorf("oauth %s get user info failed: %v", req.Provider, err)
		return errors.New("failed to get oauth user info")
	}

	// Check if this OAuth account is already bound
	var existing model.OAuthAccount
	err = s.db.Where("provider = ? AND open_id = ?", req.Provider, userInfo.OpenID).First(&existing).Error
	if err == nil {
		if existing.UserID == userID {
			return errors.New("already bound to your account")
		}
		return errors.New("this oauth account is already bound to another user")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Create binding
	rawDataJSON, _ := datatypes.MarshalJSON(userInfo.RawData)
	oauthAccount := model.OAuthAccount{
		UserID:   userID,
		Provider: req.Provider,
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

	s.log.Infof("oauth account bound: provider=%s openid=%s userid=%d", req.Provider, userInfo.OpenID, userID)
	return nil
}

// Unbind unbinds an OAuth account from the current user.
func (s *OAuthBindService) Unbind(userID uint, provider string) error {
	// Check if user has a password set (can't unbind last login method)
	var user User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Count remaining bindings
	var bindingCount int64
	s.db.Model(&model.OAuthAccount{}).Where("user_id = ?", userID).Count(&bindingCount)

	// If user has no password and only one binding left, prevent unbind
	if user.Password == "" && bindingCount <= 1 {
		return errors.New("cannot unbind: this is your only login method. Please set a password first")
	}

	result := s.db.Where("user_id = ? AND provider = ?", userID, provider).Delete(&model.OAuthAccount{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("oauth account not bound")
	}

	s.log.Infof("oauth account unbound: provider=%s userid=%d", provider, userID)
	return nil
}

// GetBoundAccounts returns all OAuth accounts bound to a user.
func (s *OAuthBindService) GetBoundAccounts(userID uint) ([]OAuthAccountInfo, error) {
	var accounts []model.OAuthAccount
	if err := s.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}

	result := make([]OAuthAccountInfo, 0, len(accounts))
	for _, a := range accounts {
		result = append(result, OAuthAccountInfo{
			ID:       a.ID,
			Provider: a.Provider,
			Username: a.Username,
			Avatar:   a.Avatar,
			Email:    a.Email,
		})
	}
	return result, nil
}

// IsBound checks if a specific provider is bound to the user.
func (s *OAuthBindService) IsBound(userID uint, provider string) bool {
	var count int64
	s.db.Model(&model.OAuthAccount{}).Where("user_id = ? AND provider = ?", userID, provider).Count(&count)
	return count > 0
}
