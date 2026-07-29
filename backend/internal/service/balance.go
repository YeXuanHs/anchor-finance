package service

import (
	"errors"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type BalanceService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewBalanceService(db *gorm.DB, log *logger.Logger) *BalanceService {
	return &BalanceService{db: db, log: log}
}

func (s *BalanceService) GetBalance(userID uint) (float64, error) {
	var balance float64
	err := s.db.Table("users").Where("id = ?", userID).Select("balance").Scan(&balance).Error
	return balance, err
}

func (s *BalanceService) Add(userID uint, amount float64, description string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("users").Where("id = ?", userID).Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}
		log := model.BalanceLog{
			UserID:      userID,
			Type:        "credit",
			Amount:      amount,
			Description: description,
		}
		return tx.Create(&log).Error
	})
}

func (s *BalanceService) Deduct(userID uint, amount float64, description string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Table("users").Where("id = ? AND balance >= ?", userID, amount).Update("balance", gorm.Expr("balance - ?", amount))
		if result.RowsAffected == 0 {
			return errors.New("insufficient balance")
		}
		if result.Error != nil {
			return result.Error
		}
		log := model.BalanceLog{
			UserID:      userID,
			Type:        "debit",
			Amount:      -amount,
			Description: description,
		}
		return tx.Create(&log).Error
	})
}

func (s *BalanceService) GetLog(userID uint, page, pageSize int) ([]model.BalanceLog, int64, error) {
	var logs []model.BalanceLog
	var total int64

	s.db.Model(&model.BalanceLog{}).Where("user_id = ?", userID).Count(&total)
	offset := (page - 1) * pageSize
	s.db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)
	return logs, total, nil
}
