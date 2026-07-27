package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserManageService provides admin operations on user accounts.
type UserManageService struct {
	db *gorm.DB
}

// NewUserManageService creates a new UserManageService.
func NewUserManageService(db *gorm.DB) *UserManageService {
	return &UserManageService{db: db}
}

// SearchUsers searches users by keyword across multiple fields.
func (s *UserManageService) SearchUsers(page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := s.db.Model(&model.User{})
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ? OR nickname LIKE ?", q, q, q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// Ban disables a user account (status=0).
func (s *UserManageService) Ban(userID uint, reason string) error {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	if user.Status == 0 {
		return errors.New("user is already banned")
	}
	return s.db.Model(&model.User{}).Where("id = ?", userID).
		Update("status", 0).Error
}

// Unban re-enables a user account (status=1).
func (s *UserManageService) Unban(userID uint) error {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	if user.Status == 1 {
		return errors.New("user is already active")
	}
	return s.db.Model(&model.User{}).Where("id = ?", userID).
		Update("status", 1).Error
}

// AdjustBalance adds or deducts balance for a user and logs the transaction.
func (s *UserManageService) AdjustBalance(userID uint, amount float64, description string) error {
	if amount == 0 {
		return errors.New("amount must not be zero")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.Select("id, balance").First(&user, userID).Error; err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		if amount < 0 && user.Balance < -amount {
			return fmt.Errorf("insufficient balance: have %.4f, need %.4f", user.Balance, -amount)
		}

		result := tx.Model(&model.User{}).Where("id = ?", userID).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount))
		if result.Error != nil {
			return result.Error
		}

		// fetch new balance
		if err := tx.Select("balance").First(&user, userID).Error; err != nil {
			return err
		}

		log := &model.BalanceLog{
			UserID:      userID,
			Amount:      amount,
			Balance:     user.Balance,
			RelatedType: "admin",
			Description: description,
		}
		return tx.Create(log).Error
	})
}

// ResetPassword sets a new password for a user by admin.
func (s *UserManageService) ResetPassword(userID uint, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.Model(&model.User{}).Where("id = ?", userID).
		Update("password_hash", string(hashed)).Error
}

// GetOperationLogs returns the operation log entries for a user.
func (s *UserManageService) GetOperationLogs(userID uint, page, pageSize int) ([]model.BalanceLog, int64, error) {
	var logs []model.BalanceLog
	var total int64

	query := s.db.Model(&model.BalanceLog{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// GetUserStatus returns user status information.
func (s *UserManageService) GetUserStatus(userID uint) (map[string]interface{}, error) {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	var serviceCount int64
	s.db.Model(&model.ClientService{}).Where("user_id = ? AND status = ?", userID, 1).Count(&serviceCount)

	var remarkCount int64
	s.db.Model(&model.UserRemark{}).Where("user_id = ?", userID).Count(&remarkCount)

	return map[string]interface{}{
		"user":          user,
		"service_count": serviceCount,
		"remark_count":  remarkCount,
		"checked_at":    time.Now(),
	}, nil
}
