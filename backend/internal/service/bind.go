package service

import (
	"errors"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// BindService 账号绑定服务
type BindService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewBindService(db *gorm.DB, log *logger.Logger) *BindService {
	return &BindService{db: db, log: log}
}

// BindOAuthRequest 绑定OAuth请求
type BindOAuthRequest struct {
	UserID      uint   `json:"user_id" binding:"required"`
	Provider    string `json:"provider" binding:"required"`
	ProviderUID string `json:"provider_uid" binding:"required"`
	AccessToken string `json:"access_token"`
	Nickname    string `json:"nickname"`
	Avatar      string `json:"avatar"`
}

// BindOAuth 绑定OAuth账号
func (s *BindService) BindOAuth(req BindOAuthRequest) (*model.OAuthBind, error) {
	var existing model.OAuthBind
	err := s.db.Where("user_id = ? AND provider = ?", req.UserID, req.Provider).First(&existing).Error
	if err == nil {
		return nil, errors.New("provider already bound to this user")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var dup model.OAuthBind
	err = s.db.Where("provider = ? AND provider_uid = ?", req.Provider, req.ProviderUID).First(&dup).Error
	if err == nil {
		return nil, errors.New("this oauth account is already bound to another user")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	bind := model.OAuthBind{
		UserID:      req.UserID,
		Provider:    req.Provider,
		ProviderUID: req.ProviderUID,
		AccessToken: req.AccessToken,
		Nickname:    req.Nickname,
		Avatar:      req.Avatar,
	}
	if err := s.db.Create(&bind).Error; err != nil {
		return nil, err
	}
	s.log.Infof("oauth bound: user=%d provider=%s", req.UserID, req.Provider)
	return &bind, nil
}

// UnbindOAuth 解绑OAuth账号
func (s *BindService) UnbindOAuth(userID uint, provider string) error {
	result := s.db.Where("user_id = ? AND provider = ?", userID, provider).Delete(&model.OAuthBind{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("binding not found")
	}
	s.log.Infof("oauth unbound: user=%d provider=%s", userID, provider)
	return nil
}

// ListBindings 列出用户的所有绑定
func (s *BindService) ListBindings(userID uint) ([]model.OAuthBind, error) {
	var binds []model.OAuthBind
	if err := s.db.Where("user_id = ?", userID).Find(&binds).Error; err != nil {
		return nil, err
	}
	return binds, nil
}

// GetBinding 获取指定绑定
func (s *BindService) GetBinding(userID uint, provider string) (*model.OAuthBind, error) {
	var bind model.OAuthBind
	if err := s.db.Where("user_id = ? AND provider = ?", userID, provider).First(&bind).Error; err != nil {
		return nil, err
	}
	return &bind, nil
}

// UpdateTokens 更新绑定的令牌
func (s *BindService) UpdateTokens(userID uint, provider, accessToken, refreshToken string, expiresAt *time.Time) error {
	updates := map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
	if expiresAt != nil {
		updates["expires_at"] = expiresAt
	}
	result := s.db.Model(&model.OAuthBind{}).
		Where("user_id = ? AND provider = ?", userID, provider).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("binding not found")
	}
	return nil
}
