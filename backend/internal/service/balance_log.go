package service

import (
	"fmt"

	"anchorfinance/internal/model"
	"gorm.io/gorm"
)

// BalanceLogService manages user balance operations.
type BalanceLogService struct {
	db *gorm.DB
}

// NewBalanceLogService creates a new BalanceLogService.
func NewBalanceLogService(db *gorm.DB) *BalanceLogService {
	return &BalanceLogService{db: db}
}

// AddBalance credits the user's balance and logs the transaction.
func (s *BalanceLogService) AddBalance(userID uint, amount float64, relatedID uint, relatedType, description string) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Update user balance
		result := tx.Model(&model.User{}).Where("id = ?", userID).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return fmt.Errorf("user not found")
		}

		// Fetch new balance for logging
		var user model.User
		if err := tx.Select("balance").First(&user, userID).Error; err != nil {
			return err
		}

		// Create balance log entry
		log := &model.BalanceLog{
			UserID:      userID,
			Amount:      amount,
			Balance:     user.Balance,
			RelatedID:   relatedID,
			RelatedType: relatedType,
			Description: description,
		}
		return tx.Create(log).Error
	})
}

// DeductBalance debits the user's balance and logs the transaction.
func (s *BalanceLogService) DeductBalance(userID uint, amount float64, relatedID uint, relatedType, description string) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Check current balance
		var user model.User
		if err := tx.Select("id, balance").First(&user, userID).Error; err != nil {
			return fmt.Errorf("user not found: %w", err)
		}
		if user.Balance < amount {
			return fmt.Errorf("insufficient balance: have %.2f, need %.2f", user.Balance, amount)
		}

		// Deduct balance
		result := tx.Model(&model.User{}).Where("id = ?", userID).
			UpdateColumn("balance", gorm.Expr("balance - ?", amount))
		if result.Error != nil {
			return result.Error
		}

		// Fetch new balance for logging
		if err := tx.Select("balance").First(&user, userID).Error; err != nil {
			return err
		}

		// Create balance log entry (negative amount for deduction)
		log := &model.BalanceLog{
			UserID:      userID,
			Amount:      -amount,
			Balance:     user.Balance,
			RelatedID:   relatedID,
			RelatedType: relatedType,
			Description: description,
		}
		return tx.Create(log).Error
	})
}

// GetLogs returns paginated balance change logs for a user.
func (s *BalanceLogService) GetLogs(userID uint, page, pageSize int) ([]model.BalanceLog, int64, error) {
	var logs []model.BalanceLog
	var total int64

	query := s.db.Model(&model.BalanceLog{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error

	return logs, total, err
}

// GetBalance returns the current balance for a user.
func (s *BalanceLogService) GetBalance(userID uint) (float64, error) {
	var user model.User
	if err := s.db.Select("balance").First(&user, userID).Error; err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}
	return user.Balance, nil
}
