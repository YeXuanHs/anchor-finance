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

// AdminGetByID returns a single affiliate by ID with user info (admin).
func (s *AffiliateService) AdminGetByID(id uint) (*model.Affiliate, error) {
	var aff model.Affiliate
	if err := s.db.Preload("User").First(&aff, id).Error; err != nil {
		return nil, err
	}
	return &aff, nil
}

// AdminUpdate updates affiliate fields (admin).
func (s *AffiliateService) AdminUpdate(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&model.Affiliate{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		return errors.New("affiliate not found")
	}
	return result.Error
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

// GetGatewayList returns available payment gateways for affiliate.
func (s *AffiliateService) GetGatewayList() ([]map[string]interface{}, error) {
	var gateways []model.PaymentGateway
	if err := s.db.Where("is_active = ?", true).Find(&gateways).Error; err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for _, gw := range gateways {
		result = append(result, map[string]interface{}{
			"id":   gw.ID,
			"name": gw.Name,
			"code": gw.Code,
		})
	}
	return result, nil
}

// ProductAffiRequest represents affiliate product settings.
type ProductAffiRequest struct {
	ID                 uint    `json:"id"`
	PID                uint    `json:"pid" binding:"required"`
	AffiliateEnabled   int     `json:"affiliate_enabled"`
	AffiliateBates     float64 `json:"affiliate_bates"`
	AffiliateType      int     `json:"affiliate_type"`
	AffiliateIsReorder int     `json:"affiliate_is_reorder"`
	AffiliateReorder   int     `json:"affiliate_reorder"`
	AffiliateReorderType int   `json:"affiliate_reorder_type"`
	AffiliateIsRenew   int     `json:"affiliate_is_renew"`
	AffiliateRenew     int     `json:"affiliate_renew"`
	AffiliateRenewType int     `json:"affiliate_renew_type"`
}

// ProductAffiliateSetting model for affiliate product settings.
type ProductAffiliateSetting struct {
	ID                   uint    `gorm:"primaryKey" json:"id"`
	PID                  uint    `gorm:"index;not null" json:"pid"`
	AffiliateEnabled     int     `gorm:"default:0" json:"affiliate_enabled"`
	AffiliateBates       float64 `gorm:"type:decimal(10,2);default:0" json:"affiliate_bates"`
	AffiliateType        int     `gorm:"default:0" json:"affiliate_type"` // 0=percentage 1=fixed
	AffiliateIsReorder   int     `gorm:"default:0" json:"affiliate_is_reorder"`
	AffiliateReorder     int     `gorm:"default:0" json:"affiliate_reorder"`
	AffiliateReorderType int     `gorm:"default:0" json:"affiliate_reorder_type"`
	AffiliateIsRenew     int     `gorm:"default:0" json:"affiliate_is_renew"`
	AffiliateRenew       int     `gorm:"default:0" json:"affiliate_renew"`
	AffiliateRenewType   int     `gorm:"default:0" json:"affiliate_renew_type"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

// GetProductAffiSetting returns affiliate settings for a product.
func (s *AffiliateService) GetProductAffiSetting(pid uint) (*ProductAffiliateSetting, error) {
	var setting ProductAffiliateSetting
	if err := s.db.Where("pid = ?", pid).First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &setting, nil
}

// SaveProductAffiSetting creates or updates affiliate settings for a product.
func (s *AffiliateService) SaveProductAffiSetting(req *ProductAffiRequest) error {
	if req.ID > 0 {
		// Update existing
		return s.db.Model(&ProductAffiliateSetting{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
			"affiliate_enabled":      req.AffiliateEnabled,
			"affiliate_bates":        req.AffiliateBates,
			"affiliate_type":         req.AffiliateType,
			"affiliate_is_reorder":   req.AffiliateIsReorder,
			"affiliate_reorder":      req.AffiliateReorder,
			"affiliate_reorder_type": req.AffiliateReorderType,
			"affiliate_is_renew":     req.AffiliateIsRenew,
			"affiliate_renew":        req.AffiliateRenew,
			"affiliate_renew_type":   req.AffiliateRenewType,
		}).Error
	}

	// Create new
	setting := &ProductAffiliateSetting{
		PID:                  req.PID,
		AffiliateEnabled:     req.AffiliateEnabled,
		AffiliateBates:       req.AffiliateBates,
		AffiliateType:        req.AffiliateType,
		AffiliateIsReorder:   req.AffiliateIsReorder,
		AffiliateReorder:     req.AffiliateReorder,
		AffiliateReorderType: req.AffiliateReorderType,
		AffiliateIsRenew:     req.AffiliateIsRenew,
		AffiliateRenew:       req.AffiliateRenew,
		AffiliateRenewType:   req.AffiliateRenewType,
	}
	return s.db.Create(setting).Error
}

// GetUserBuyRecords returns affiliate purchase records for a user.
func (s *AffiliateService) GetUserBuyRecords(uid string, page, pageSize int) ([]map[string]interface{}, int64, error) {
	var total int64
	query := s.db.Table("invoices i").
		Joins("LEFT JOIN users u ON i.uid = u.id").
		Where("i.uid = ? AND i.status = ?", uid, "Paid")

	query.Count(&total)

	var records []map[string]interface{}
	offset := (page - 1) * pageSize
	if err := query.Select("i.id, i.type, i.subtotal, i.paid_time, u.username, u.id as uid").
		Offset(offset).Limit(pageSize).
		Order("i.id DESC").
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// ==================== 新增缺失方法 ====================

// UserAffiSettingRequest 用户推介计划设置请求
type UserAffiSettingRequest struct {
	ID                  uint    `json:"id"`
	UID                 uint    `json:"uid" binding:"required"`
	AffiliateEnabled    int     `json:"affiliate_enabled"`
	AffiliateBates      float64 `json:"affiliate_bates"`
	AffiliateType       int     `json:"affiliate_type"`
	AffiliateIsReorder  int     `json:"affiliate_is_reorder"`
	AffiliateReorder    int     `json:"affiliate_reorder"`
	AffiliateReorderType int    `json:"affiliate_reorder_type"`
	AffiliateIsRenew    int     `json:"affiliate_is_renew"`
	AffiliateRenew      int     `json:"affiliate_renew"`
	AffiliateRenewType  int     `json:"affiliate_renew_type"`
}

// GetUserAffiPage returns affiliate settings page data for a user.
func (s *AffiliateService) GetUserAffiPage(uid string) (map[string]interface{}, map[string]interface{}, error) {
	var setting map[string]interface{}
	if err := s.db.Table("affiliates_user_setting").Where("uid = ?", uid).Find(&setting).Error; err != nil {
		return nil, nil, err
	}

	var affData map[string]interface{}
	if err := s.db.Table("affiliates").Where("uid = ?", uid).Find(&affData).Error; err != nil {
		return nil, nil, err
	}

	return setting, affData, nil
}

// UpdateUserAffiBalance updates affiliate balance for a user.
func (s *AffiliateService) UpdateUserAffiBalance(uid uint, withdrawn, balance float64) error {
	updates := map[string]interface{}{}
	if withdrawn > 0 {
		updates["withdrawn"] = withdrawn
	}
	if balance > 0 {
		updates["balance"] = balance
	}
	if len(updates) == 0 {
		return nil
	}
	return s.db.Table("affiliates").Where("uid = ?", uid).Updates(updates).Error
}

// SaveUserAffiSetting creates or updates affiliate settings for a user.
func (s *AffiliateService) SaveUserAffiSetting(req *UserAffiSettingRequest) error {
	if req.ID == 0 {
		// Create new
		data := map[string]interface{}{
			"uid":                   req.UID,
			"affiliate_enabled":     req.AffiliateEnabled,
			"affiliate_bates":       req.AffiliateBates,
			"affiliate_type":        req.AffiliateType,
			"affiliate_is_reorder":  req.AffiliateIsReorder,
			"affiliate_reorder":     req.AffiliateReorder,
			"affiliate_reorder_type": req.AffiliateReorderType,
			"affiliate_is_renew":    req.AffiliateIsRenew,
			"affiliate_renew":       req.AffiliateRenew,
			"affiliate_renew_type":  req.AffiliateRenewType,
			"create_time":           time.Now().Unix(),
		}
		return s.db.Table("affiliates_user_setting").Create(data).Error
	}

	// Update existing
	data := map[string]interface{}{
		"affiliate_enabled":      req.AffiliateEnabled,
		"affiliate_bates":        req.AffiliateBates,
		"affiliate_type":         req.AffiliateType,
		"affiliate_is_reorder":   req.AffiliateIsReorder,
		"affiliate_reorder":      req.AffiliateReorder,
		"affiliate_reorder_type": req.AffiliateReorderType,
		"affiliate_is_renew":     req.AffiliateIsRenew,
		"affiliate_renew":        req.AffiliateRenew,
		"affiliate_renew_type":   req.AffiliateRenewType,
	}
	return s.db.Table("affiliates_user_setting").Where("id = ?", req.ID).Updates(data).Error
}

// GetUserAffiList returns referred users list for an affiliate.
func (s *AffiliateService) GetUserAffiList(uid string, page, pageSize int, keyword string) ([]map[string]interface{}, int64, error) {
	// Get affiliate ID
	var aff map[string]interface{}
	if err := s.db.Table("affiliates").Select("id").Where("uid = ?", uid).Find(&aff).Error; err != nil {
		return nil, 0, err
	}
	if aff == nil {
		return []map[string]interface{}{}, 0, nil
	}
	affID := aff["id"]

	query := s.db.Table("affiliates_user au").
		Joins("LEFT JOIN users u ON au.uid = u.id").
		Where("au.affid = ?", affID)

	if keyword != "" {
		query = query.Where("u.username LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var users []map[string]interface{}
	offset := (page - 1) * pageSize
	if err := query.Select("u.id, u.username, u.email, u.created_at, u.last_login_at").
		Offset(offset).Limit(pageSize).
		Order("au.id DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserAffiRecord returns withdrawal records for a specific user.
func (s *AffiliateService) GetUserAffiRecord(uid string, page, pageSize int, status string) ([]map[string]interface{}, int64, error) {
	query := s.db.Table("affiliates_withdraw aw").
		Joins("LEFT JOIN users u ON aw.uid = u.id").
		Where("aw.uid = ?", uid)

	if status != "" {
		query = query.Where("aw.status = ?", status)
	}

	var total int64
	query.Count(&total)

	var records []map[string]interface{}
	offset := (page - 1) * pageSize
	if err := query.Select("aw.id, aw.num, aw.type, aw.create_time, aw.status, aw.reason, u.username").
		Offset(offset).Limit(pageSize).
		Order("aw.id DESC").
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// GetAffiliateUserIDs returns user IDs referred by an affiliate.
func (s *AffiliateService) GetAffiliateUserIDs(uid string) ([]uint, error) {
	var aff map[string]interface{}
	if err := s.db.Table("affiliates").Select("id").Where("uid = ?", uid).Find(&aff).Error; err != nil {
		return nil, err
	}
	if aff == nil {
		return []uint{}, nil
	}
	affID := aff["id"]

	var uids []uint
	if err := s.db.Table("affiliates_user").Where("affid = ?", affID).Pluck("uid", &uids).Error; err != nil {
		return nil, err
	}

	return uids, nil
}

// GetAffiWithdrawRecords returns affiliate withdrawal records.
func (s *AffiliateService) GetAffiWithdrawRecords(page, pageSize int, keyword, status string) ([]map[string]interface{}, int64, error) {
	query := s.db.Table("affiliates_withdraw aw").
		Joins("LEFT JOIN users u ON aw.uid = u.id").
		Joins("LEFT JOIN currencies cu ON u.currency = cu.id")

	if keyword != "" {
		query = query.Where("u.username LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("aw.status = ?", status)
	}

	var total int64
	query.Count(&total)

	var records []map[string]interface{}
	offset := (page - 1) * pageSize
	if err := query.Select("aw.id, aw.num, aw.type, aw.create_time, aw.status, aw.reason, u.username, cu.suffix").
		Offset(offset).Limit(pageSize).
		Order("aw.id DESC").
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// ProcessAffiWithdrawSH processes affiliate withdrawal approval/rejection.
func (s *AffiliateService) ProcessAffiWithdrawSH(withdrawID uint, typ, status int, payment, transID, reason string, adminID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var withdraw model.AffiliateWithdraw
		if err := tx.First(&withdraw, withdrawID).Error; err != nil {
			return fmt.Errorf("withdrawal not found")
		}

		if withdraw.Status != 1 {
			return fmt.Errorf("withdrawal is not pending")
		}

		if status == 2 {
			// Approve
			if typ == 1 {
				// Add to user balance
				if err := tx.Table("users").Where("id = ?", withdraw.UserID).
					Update("credit", gorm.Expr("credit + ?", withdraw.Amount)).Error; err != nil {
					return err
				}
			} else if typ == 3 {
				// Create account record
				account := map[string]interface{}{
					"uid":          withdraw.UserID,
					"gateway":      payment,
					"create_time":  time.Now().Unix(),
					"pay_time":     time.Now().Unix(),
					"amount_out":   withdraw.Amount,
					"trans_id":     transID,
					"description":  "Affiliate withdrawal",
				}
				if err := tx.Table("accounts").Create(account).Error; err != nil {
					return err
				}
			}

			// Update withdrawal status
			if err := tx.Model(&withdraw).Updates(map[string]interface{}{
				"status":      2,
				"type":        typ,
				"admin_id":    adminID,
				"update_time": time.Now().Unix(),
			}).Error; err != nil {
				return err
			}

			// Update affiliate balance
			var aff model.Affiliate
			if err := tx.Where("uid = ?", withdraw.UserID).First(&aff).Error; err == nil {
				tx.Model(&aff).Updates(map[string]interface{}{
					"withdraw_ing": gorm.Expr("withdraw_ing - ?", withdraw.Amount),
					"withdrawn":    gorm.Expr("withdrawn + ?", withdraw.Amount),
					"updated_time": time.Now().Unix(),
				})
			}
		} else if status == 3 {
			// Reject
			if err := tx.Model(&withdraw).Updates(map[string]interface{}{
				"status":      3,
				"reason":      reason,
				"admin_id":    adminID,
				"update_time": time.Now().Unix(),
			}).Error; err != nil {
				return err
			}

			// Return balance to affiliate
			var aff model.Affiliate
			if err := tx.Where("uid = ?", withdraw.UserID).First(&aff).Error; err == nil {
				tx.Model(&aff).Updates(map[string]interface{}{
					"balance":      gorm.Expr("balance + ?", withdraw.Amount),
					"withdraw_ing": gorm.Expr("withdraw_ing - ?", withdraw.Amount),
					"updated_time": time.Now().Unix(),
				})
			}
		}

		return nil
	})
}

// GetCommissionRecords returns commission records with enrichment.
func (s *AffiliateService) GetCommissionRecords(uid string, page, pageSize int) ([]map[string]interface{}, int64, error) {
	var total int64
	query := s.db.Table("affiliate_records ar").
		Joins("LEFT JOIN users u ON ar.user_id = u.id").
		Where("ar.user_id = ?", uid)

	query.Count(&total)

	var records []map[string]interface{}
	offset := (page - 1) * pageSize
	if err := query.Select("ar.*").
		Offset(offset).Limit(pageSize).
		Order("ar.id DESC").
		Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}
