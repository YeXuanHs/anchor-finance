package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// GetUserService 用户查询服务
type GetUserService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewGetUserService(db *gorm.DB, log *logger.Logger) *GetUserService {
	return &GetUserService{db: db, log: log}
}

// UserProfile 用户档案
type UserProfile struct {
	model.User
	OrderCount   int64 `json:"order_count"`
	InvoiceCount int64 `json:"invoice_count"`
	ActiveHosts  int64 `json:"active_hosts"`
}

// UserStats 用户统计
type UserStats struct {
	TotalUsers    int64   `json:"total_users"`
	ActiveUsers   int64   `json:"active_users"`
	DisabledUsers int64   `json:"disabled_users"`
	PendingUsers  int64   `json:"pending_users"`
}

// GetByID 根据ID获取用户
func (s *GetUserService) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.Preload("Group").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (s *GetUserService) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := s.db.Preload("Group").Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Search 搜索用户
func (s *GetUserService) Search(keyword string, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := s.db.Model(&model.User{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ? OR phone LIKE ?", like, like, like, like)
	}

	query.Count(&total)
	offset, limit := util.Paginate(page, pageSize)
	query.Offset(offset).Limit(limit).Order("id DESC").Preload("Group").Find(&users)
	return users, total, nil
}

// GetProfile 获取用户档案（含统计）
func (s *GetUserService) GetProfile(id uint) (*UserProfile, error) {
	var user model.User
	if err := s.db.Preload("Group").First(&user, id).Error; err != nil {
		return nil, err
	}

	profile := &UserProfile{User: user}

	s.db.Model(&model.Order{}).Where("user_id = ?", id).Count(&profile.OrderCount)
	s.db.Model(&model.Invoice{}).Where("user_id = ? AND status = ?", id, 1).Count(&profile.InvoiceCount)
	s.db.Model(&model.Host{}).Where("owner_id = ? AND status = ?", id, 1).Count(&profile.ActiveHosts)

	return profile, nil
}

// GetStats 获取用户统计
func (s *GetUserService) GetStats() (*UserStats, error) {
	stats := &UserStats{}

	s.db.Model(&model.User{}).Count(&stats.TotalUsers)
	s.db.Model(&model.User{}).Where("status = ?", 1).Count(&stats.ActiveUsers)
	s.db.Model(&model.User{}).Where("status = ?", 0).Count(&stats.DisabledUsers)
	s.db.Model(&model.User{}).Where("status = ?", 2).Count(&stats.PendingUsers)

	return stats, nil
}
