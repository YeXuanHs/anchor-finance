package service

import (
	"errors"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type OAuthProviderService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewOAuthProviderService(db *gorm.DB, log *logger.Logger) *OAuthProviderService {
	return &OAuthProviderService{db: db, log: log}
}

type CreateOAuthProviderRequest struct {
	Name         string `json:"name" binding:"required,max=32"`
	DisplayName  string `json:"display_name" binding:"required,max=64"`
	Icon         string `json:"icon" binding:"max=512"`
	ClientID     string `json:"client_id" binding:"required,max=256"`
	ClientSecret string `json:"client_secret" binding:"required,max=256"`
	RedirectURL  string `json:"redirect_url" binding:"max=512"`
	AuthURL      string `json:"auth_url" binding:"max=512"`
	TokenURL     string `json:"token_url" binding:"max=512"`
	UserInfoURL  string `json:"user_info_url" binding:"max=512"`
	Scopes       string `json:"scopes" binding:"max=512"`
	Extra        string `json:"extra"`
	IsEnabled    bool   `json:"is_enabled"`
	SortOrder    int    `json:"sort_order"`
	Status       int16  `json:"status"`
}

type UpdateOAuthProviderRequest struct {
	DisplayName  *string `json:"display_name" binding:"omitempty,max=64"`
	Icon         *string `json:"icon" binding:"omitempty,max=512"`
	ClientID     *string `json:"client_id" binding:"omitempty,max=256"`
	ClientSecret *string `json:"client_secret" binding:"omitempty,max=256"`
	RedirectURL  *string `json:"redirect_url" binding:"omitempty,max=512"`
	AuthURL      *string `json:"auth_url" binding:"omitempty,max=512"`
	TokenURL     *string `json:"token_url" binding:"omitempty,max=512"`
	UserInfoURL  *string `json:"user_info_url" binding:"omitempty,max=512"`
	Scopes       *string `json:"scopes" binding:"omitempty,max=512"`
	Extra        *string `json:"extra"`
	IsEnabled    *bool   `json:"is_enabled"`
	SortOrder    *int    `json:"sort_order"`
	Status       *int16  `json:"status"`
}

func (s *OAuthProviderService) Create(req CreateOAuthProviderRequest) (*model.OAuthProvider, error) {
	var count int64
	s.db.Model(&model.OAuthProvider{}).Where("name = ?", req.Name).Count(&count)
	if count > 0 {
		return nil, errors.New("oauth provider name already exists")
	}

	provider := model.OAuthProvider{
		Name:         req.Name,
		DisplayName:  req.DisplayName,
		Icon:         req.Icon,
		ClientID:     req.ClientID,
		ClientSecret: req.ClientSecret,
		RedirectURL:  req.RedirectURL,
		AuthURL:      req.AuthURL,
		TokenURL:     req.TokenURL,
		UserInfoURL:  req.UserInfoURL,
		Scopes:       req.Scopes,
		IsEnabled:    req.IsEnabled,
		SortOrder:    req.SortOrder,
		Status:       req.Status,
	}
	if provider.Status == 0 {
		provider.Status = 1
	}
	if req.Extra != "" {
		provider.Extra = datatypes.JSON(req.Extra)
	}

	if err := s.db.Create(&provider).Error; err != nil {
		return nil, err
	}
	s.log.Infof("oauth provider created: id=%d name=%s", provider.ID, provider.Name)
	return &provider, nil
}

func (s *OAuthProviderService) Update(id uint, req UpdateOAuthProviderRequest) (*model.OAuthProvider, error) {
	var provider model.OAuthProvider
	if err := s.db.First(&provider, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.ClientID != nil {
		updates["client_id"] = *req.ClientID
	}
	if req.ClientSecret != nil {
		updates["client_secret"] = *req.ClientSecret
	}
	if req.RedirectURL != nil {
		updates["redirect_url"] = *req.RedirectURL
	}
	if req.AuthURL != nil {
		updates["auth_url"] = *req.AuthURL
	}
	if req.TokenURL != nil {
		updates["token_url"] = *req.TokenURL
	}
	if req.UserInfoURL != nil {
		updates["user_info_url"] = *req.UserInfoURL
	}
	if req.Scopes != nil {
		updates["scopes"] = *req.Scopes
	}
	if req.Extra != nil {
		updates["extra"] = datatypes.JSON(*req.Extra)
	}
	if req.IsEnabled != nil {
		updates["is_enabled"] = *req.IsEnabled
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := s.db.Model(&provider).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&provider, id).Error; err != nil {
		return nil, err
	}
	s.log.Infof("oauth provider updated: id=%d", id)
	return &provider, nil
}

func (s *OAuthProviderService) Delete(id uint) error {
	var provider model.OAuthProvider
	if err := s.db.First(&provider, id).Error; err != nil {
		return err
	}
	if provider.UserCount > 0 {
		return errors.New("cannot delete provider with bound users")
	}

	result := s.db.Delete(&model.OAuthProvider{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("oauth provider not found")
	}
	s.log.Infof("oauth provider deleted: id=%d", id)
	return nil
}

func (s *OAuthProviderService) GetByID(id uint) (*model.OAuthProvider, error) {
	var provider model.OAuthProvider
	if err := s.db.First(&provider, id).Error; err != nil {
		return nil, err
	}
	return &provider, nil
}

func (s *OAuthProviderService) GetList(page, pageSize int, status int) ([]model.OAuthProvider, int64, error) {
	var items []model.OAuthProvider
	var total int64

	query := s.db.Model(&model.OAuthProvider{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := util.Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("sort_order ASC, id ASC").
		Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

// GetEnabled returns enabled OAuth providers for frontend.
func (s *OAuthProviderService) GetEnabled() ([]model.OAuthProvider, error) {
	var items []model.OAuthProvider
	if err := s.db.Where("status = 1 AND is_enabled = true").
		Order("sort_order ASC, id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *OAuthProviderService) ToggleStatus(id uint) error {
	var provider model.OAuthProvider
	if err := s.db.First(&provider, id).Error; err != nil {
		return err
	}
	newStatus := int16(1)
	if provider.Status == 1 {
		newStatus = 0
	}
	return s.db.Model(&provider).Update("status", newStatus).Error
}

func (s *OAuthProviderService) ToggleEnabled(id uint) error {
	var provider model.OAuthProvider
	if err := s.db.First(&provider, id).Error; err != nil {
		return err
	}
	return s.db.Model(&provider).Update("is_enabled", !provider.IsEnabled).Error
}

func (s *OAuthProviderService) UpdateSort(id uint, sortOrder int) error {
	result := s.db.Model(&model.OAuthProvider{}).Where("id = ?", id).Update("sort_order", sortOrder)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("oauth provider not found")
	}
	return nil
}
