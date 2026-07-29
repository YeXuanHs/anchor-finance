package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type ApiKey struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	UserID    uint       `gorm:"index;not null" json:"user_id"`
	Name      string     `gorm:"type:varchar(64);not null" json:"name"`
	Key       string     `gorm:"type:varchar(64);uniqueIndex;not null" json:"key"`
	Secret    string     `gorm:"type:varchar(128)" json:"-"`
	Permissions string   `gorm:"type:jsonb" json:"permissions"`
	RateLimit int        `gorm:"default:100" json:"rate_limit"`
	IsActive  bool       `gorm:"default:true;index" json:"is_active"`
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt *time.Time  `json:"expires_at"`
	CreatedAt time.Time   `json:"created_at"`
}

type ApiManageService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewApiManageService(db *gorm.DB, log *logger.Logger) *ApiManageService {
	return &ApiManageService{db: db, log: log}
}

func (s *ApiManageService) CreateKey(userID uint, name string, permissions string) (*ApiKey, error) {
	keyBytes := make([]byte, 32)
	rand.Read(keyBytes)
	key := hex.EncodeToString(keyBytes)
	secretBytes := make([]byte, 32)
	rand.Read(secretBytes)
	secret := hex.EncodeToString(secretBytes)
	apiKey := ApiKey{
		UserID:    userID,
		Name:      name,
		Key:       key,
		Secret:    secret,
		Permissions: permissions,
		IsActive:  true,
	}
	if err := s.db.Create(&apiKey).Error; err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (s *ApiManageService) GetByUserID(userID uint) ([]ApiKey, error) {
	var keys []ApiKey
	err := s.db.Where("user_id = ?", userID).Order("id DESC").Find(&keys).Error
	return keys, err
}

func (s *ApiManageService) GetByKey(key string) (*ApiKey, error) {
	var apiKey ApiKey
	if err := s.db.Where("`key` = ? AND is_active = ?", key, true).First(&apiKey).Error; err != nil {
		return nil, err
	}
	return &apiKey, nil
}

func (s *ApiManageService) ToggleStatus(id uint, userID uint) error {
	result := s.db.Model(&ApiKey{}).Where("id = ? AND user_id = ?", id, userID).Update("is_active", gorm.Expr("NOT is_active"))
	if result.RowsAffected == 0 {
		return errors.New("api key not found")
	}
	return result.Error
}

func (s *ApiManageService) Delete(id uint, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&ApiKey{})
	if result.RowsAffected == 0 {
		return errors.New("api key not found")
	}
	return result.Error
}
