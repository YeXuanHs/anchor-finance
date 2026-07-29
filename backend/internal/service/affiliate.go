package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type AffiliateService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAffiliateService(db *gorm.DB, log *logger.Logger) *AffiliateService {
	return &AffiliateService{db: db, log: log}
}

// TrackVisit records a referral visit for tracking.
func (s *AffiliateService) TrackVisit(referralCode, ip, userAgent, refererURL, landingURL string) error {
	var aff model.Affiliate
	if err := s.db.Where("code = ? AND is_active = true", referralCode).First(&aff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid referral code")
		}
		return err
	}

	visit := &model.AffiliateVisit{
		AffiliateID: aff.ID,
		IP:          ip,
		UserAgent:   userAgent,
		RefererURL:  refererURL,
		LandingURL:  landingURL,
	}
	return s.db.Create(visit).Error
}

// TrackSignup links a new user to the referrer via the referral code.
func (s *AffiliateService) TrackSignup(userID uint, referralCode string) error {
	var aff model.Affiliate
	if err := s.db.Where("code = ? AND is_active = true", referralCode).First(&aff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid referral code")
		}
		return err
	}

	// Prevent self-referral
	if aff.UserID == userID {
		return errors.New("cannot refer yourself")
	}

	// Mark the most recent unconverted visit from this IP as converted
	var visit model.AffiliateVisit
	if err := s.db.Where("affiliate_id = ? AND converted = false", aff.ID).
		Order("id DESC").First(&visit).Error; err == nil {
		s.db.Model(&visit).Updates(map[string]interface{}{
			"converted": true,
			"user_id":   userID,
		})
	}

	s.log.Infof("affiliate signup tracked: user=%d referrer=%d code=%s", userID, aff.UserID, referralCode)
	return nil
}

// TrackOrder calculates and records commission for an order made by a referred user.
func (s *AffiliateService) TrackOrder(referralCode string, orderID uint, orderAmount float64) (*model.AffiliateRecord, error) {
	var aff model.Affiliate
	if err := s.db.Where("code = ? AND is_active = true", referralCode).First(&aff).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid referral code")
		}
		return nil, err
	}

	commission := orderAmount * aff.CommissionRate / 100
	if commission <= 0 {
		return nil, errors.New("commission amount is zero")
	}

	record := &model.AffiliateRecord{
		AffiliateID: aff.ID,
		UserID:      aff.UserID,
		OrderID:     orderID,
		Amount:      commission,
		Status:      1, // pending confirmation
		Description: fmt.Sprintf("Commission for order #%d (amount: %.2f, rate: %.1f%%)", orderID, orderAmount, aff.CommissionRate),
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		return tx.Model(&model.Affiliate{}).
			Where("id = ?", aff.ID).
			Updates(map[string]interface{}{
				"balance":      gorm.Expr("balance + ?", commission),
				"total_earned": gorm.Expr("total_earned + ?", commission),
			}).Error
	})
	if err != nil {
		return nil, err
	}

	s.log.Infof("affiliate order tracked: affiliate=%d order=%d commission=%.2f", aff.ID, orderID, commission)
	return record, nil
}

// GetVisits returns paginated visit records for an affiliate (admin).
func (s *AffiliateService) GetVisits(affiliateID, page, pageSize int) ([]model.AffiliateVisit, int64, error) {
	var visits []model.AffiliateVisit
	var total int64

	query := s.db.Model(&model.AffiliateVisit{}).Where("affiliate_id = ?", affiliateID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&visits).Error; err != nil {
		return nil, 0, err
	}

	return visits, total, nil
}

// GetStats returns aggregate statistics for an affiliate.
func (s *AffiliateService) GetStats(affiliateID uint) (map[string]interface{}, error) {
	var aff model.Affiliate
	if err := s.db.First(&aff, affiliateID).Error; err != nil {
		return nil, errors.New("affiliate not found")
	}

	var visitCount int64
	s.db.Model(&model.AffiliateVisit{}).Where("affiliate_id = ?", affiliateID).Count(&visitCount)

	var signupCount int64
	s.db.Model(&model.AffiliateVisit{}).Where("affiliate_id = ? AND converted = true", affiliateID).Count(&signupCount)

	var orderCount int64
	s.db.Model(&model.AffiliateRecord{}).Where("affiliate_id = ?", affiliateID).Count(&orderCount)

	var confirmedEarnings float64
	s.db.Model(&model.AffiliateRecord{}).Where("affiliate_id = ? AND status = 2", affiliateID).
		Select("COALESCE(SUM(amount), 0)").Scan(&confirmedEarnings)

	return map[string]interface{}{
		"total_visits":        visitCount,
		"total_signups":       signupCount,
		"total_orders":        orderCount,
		"confirmed_earnings":  confirmedEarnings,
		"balance":             aff.Balance,
		"total_earned":        aff.TotalEarned,
		"total_withdrawn":     aff.TotalWithdrawn,
		"commission_rate":     aff.CommissionRate,
	}, nil
}

// UpdateCommissionRate updates the commission rate for an affiliate (admin).
func (s *AffiliateService) UpdateCommissionRate(affiliateID uint, rate float64) error {
	if rate < 0 || rate > 100 {
		return errors.New("commission rate must be between 0 and 100")
	}
	return s.db.Model(&model.Affiliate{}).Where("id = ?", affiliateID).
		Update("commission_rate", rate).Error
}

// GetByUserID returns the affiliate info for a user, creating one if it doesn't exist.
func (s *AffiliateService) GetByUserID(userID uint) (*model.Affiliate, error) {
	var aff model.Affiliate
	err := s.db.Where("user_id = ?", userID).First(&aff).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		code, err := s.GenerateCode()
		if err != nil {
			return nil, fmt.Errorf("generate affiliate code: %w", err)
		}
		aff = model.Affiliate{
			UserID:         userID,
			Code:           code,
			CommissionRate: 10,
			WithdrawMin:    100,
			IsActive:       true,
		}
		if err := s.db.Create(&aff).Error; err != nil {
			return nil, err
		}
		return &aff, nil
	}
	if err != nil {
		return nil, err
	}
	return &aff, nil
}

// GenerateCode creates a random 8-character alphanumeric referral code.
func (s *AffiliateService) GenerateCode() (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for i := 0; i < 10; i++ {
		code := make([]byte, 8)
		for j := range code {
			n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
			if err != nil {
				return "", err
			}
			code[j] = chars[n.Int64()]
		}
		str := string(code)
		var count int64
		s.db.Model(&model.Affiliate{}).Where("code = ?", str).Count(&count)
		if count == 0 {
			return str, nil
		}
	}
	return "", errors.New("failed to generate unique code after 10 attempts")
}

// AddRecord adds a commission record when a referred user makes a payment.
func (s *AffiliateService) AddRecord(affiliateID, userID, orderID uint, amount float64, description string) (*model.AffiliateRecord, error) {
	record := &model.AffiliateRecord{
		AffiliateID: affiliateID,
		UserID:      userID,
		OrderID:     orderID,
		Amount:      amount,
		Status:      1,
		Description: description,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(record).Error; err != nil {
			return err
		}
		return tx.Model(&model.Affiliate{}).
			Where("id = ?", affiliateID).
			Updates(map[string]interface{}{
				"balance":      gorm.Expr("balance + ?", amount),
				"total_earned": gorm.Expr("total_earned + ?", amount),
			}).Error
	})
	if err != nil {
		return nil, err
	}

	s.log.Infof("affiliate record added: affiliate=%d user=%d amount=%.2f", affiliateID, userID, amount)
	return record, nil
}

// GetRecords returns paginated commission records for an affiliate.
func (s *AffiliateService) GetRecords(affiliateID, page, pageSize int) ([]model.AffiliateRecord, int64, error) {
	var records []model.AffiliateRecord
	var total int64

	query := s.db.Model(&model.AffiliateRecord{}).Where("affiliate_id = ?", affiliateID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// ConfirmRecord confirms a pending commission record (admin).
func (s *AffiliateService) ConfirmRecord(recordID uint) error {
	return s.db.Model(&model.AffiliateRecord{}).
		Where("id = ? AND status = 1", recordID).
		Update("status", 2).Error
}

// ApplyWithdraw creates a withdrawal request.
func (s *AffiliateService) ApplyWithdraw(affiliateID uint, amount float64, method, account string) (*model.AffiliateWithdraw, error) {
	var aff model.Affiliate
	if err := s.db.First(&aff, affiliateID).Error; err != nil {
		return nil, errors.New("affiliate not found")
	}

	if amount < aff.WithdrawMin {
		return nil, fmt.Errorf("minimum withdrawal amount is %.2f", aff.WithdrawMin)
	}
	if aff.Balance < amount {
		return nil, errors.New("insufficient balance")
	}

	fee := amount * 0.02 // 2% 手续费
	actual := amount - fee

	withdraw := &model.AffiliateWithdraw{
		AffiliateID: affiliateID,
		Amount:      amount,
		Fee:         fee,
		Actual:      actual,
		Method:      method,
		Account:     account,
		Status:      1,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(withdraw).Error; err != nil {
			return err
		}
		return tx.Model(&model.Affiliate{}).
			Where("id = ?", affiliateID).
			Update("balance", gorm.Expr("balance - ?", amount)).Error
	})
	if err != nil {
		return nil, err
	}

	s.log.Infof("affiliate withdrawal applied: affiliate=%d amount=%.2f", affiliateID, amount)
	return withdraw, nil
}

// GetWithdraws returns paginated withdrawal records for an affiliate.
func (s *AffiliateService) GetWithdraws(affiliateID, page, pageSize int) ([]model.AffiliateWithdraw, int64, error) {
	var withdraws []model.AffiliateWithdraw
	var total int64

	query := s.db.Model(&model.AffiliateWithdraw{}).Where("affiliate_id = ?", affiliateID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&withdraws).Error; err != nil {
		return nil, 0, err
	}

	return withdraws, total, nil
}

// GetList returns a paginated affiliate list (admin).
func (s *AffiliateService) GetList(page, pageSize int) ([]model.Affiliate, int64, error) {
	var affiliates []model.Affiliate
	var total int64

	query := s.db.Model(&model.Affiliate{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Preload("User").Find(&affiliates).Error; err != nil {
		return nil, 0, err
	}

	return affiliates, total, nil
}

// ProcessWithdraw approves or rejects a withdrawal (admin).
func (s *AffiliateService) ProcessWithdraw(withdrawID uint, approve bool, adminNote string) error {
	now := time.Now()
	if approve {
		err := s.db.Model(&model.AffiliateWithdraw{}).
			Where("id = ? AND status = 1", withdrawID).
			Updates(map[string]interface{}{
				"status":       2,
				"admin_note":   adminNote,
				"processed_at": &now,
			}).Error
		if err != nil {
			return err
		}
		s.log.Infof("affiliate withdrawal approved: id=%d", withdrawID)
		return nil
	}

	// Reject: refund the balance
	var withdraw model.AffiliateWithdraw
	if err := s.db.First(&withdraw, withdrawID).Error; err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.AffiliateWithdraw{}).
			Where("id = ? AND status = 1", withdrawID).
			Updates(map[string]interface{}{
				"status":       3,
				"admin_note":   adminNote,
				"processed_at": &now,
			}).Error; err != nil {
			return err
		}
		return tx.Model(&model.Affiliate{}).
			Where("id = ?", withdraw.AffiliateID).
			Update("balance", gorm.Expr("balance + ?", withdraw.Amount)).Error
	})
}
