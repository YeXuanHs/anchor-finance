package service

import (
	"encoding/json"
	"errors"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"

	"anchorfinance/pkg/logger"
)

// SalePromotion 促销活动
type SalePromotion struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:256;not null" json:"name"`
	Code            string         `gorm:"uniqueIndex;size:64" json:"code"`
	Type            string         `gorm:"size:32;not null;index" json:"type"` // amount_off/percent_off/first_order/free_trial
	Condition       datatypes.JSON `gorm:"type:jsonb;not null" json:"condition"`
	Discount        datatypes.JSON `gorm:"type:jsonb;not null" json:"discount"`
	StartAt         time.Time      `gorm:"index;not null" json:"start_at"`
	EndAt           time.Time      `gorm:"index;not null" json:"end_at"`
	MaxUses         int            `gorm:"default:0" json:"max_uses"`
	UsedCount       int            `gorm:"default:0" json:"used_count"`
	MaxUsesPerUser  int            `gorm:"default:1" json:"max_uses_per_user"`
	Stackable       bool           `gorm:"default:false" json:"stackable"`
	AutoApply       bool           `gorm:"default:false" json:"auto_apply"`
	Status          int            `gorm:"default:1;index" json:"status"` // 1=启用 0=禁用
	Description     string         `gorm:"type:text" json:"description"`
	AdminNotes      string         `gorm:"type:text" json:"admin_notes"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type SalePromotionLog struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PromotionID    uint      `gorm:"index;not null" json:"promotion_id"`
	UserID         uint      `gorm:"index;not null" json:"user_id"`
	OrderID        *uint     `gorm:"index" json:"order_id"`
	DiscountAmount float64   `gorm:"type:decimal(12,2);not null" json:"discount_amount"`
	UsedAt         time.Time `gorm:"not null" json:"used_at"`
}

type SaleService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewSaleService(db *gorm.DB, log *logger.Logger) *SaleService {
	return &SaleService{db: db, log: log}
}

type CreateSalePromotionRequest struct {
	Name           string                 `json:"name" binding:"required,max=256"`
	Code           string                 `json:"code" binding:"omitempty,max=64"`
	Type           string                 `json:"type" binding:"required,oneof=amount_off percent_off first_order free_trial"`
	Condition      map[string]interface{} `json:"condition" binding:"required"`
	Discount       map[string]interface{} `json:"discount" binding:"required"`
	StartAt        time.Time              `json:"start_at" binding:"required"`
	EndAt          time.Time              `json:"end_at" binding:"required"`
	MaxUses        int                    `json:"max_uses"`
	MaxUsesPerUser int                    `json:"max_uses_per_user"`
	Stackable      bool                   `json:"stackable"`
	AutoApply      bool                   `json:"auto_apply"`
	Description    string                 `json:"description"`
}

type UpdateSalePromotionRequest struct {
	Name           *string                `json:"name"`
	Code           *string                `json:"code"`
	Type           *string                `json:"type"`
	Condition      map[string]interface{} `json:"condition"`
	Discount       map[string]interface{} `json:"discount"`
	StartAt        *time.Time             `json:"start_at"`
	EndAt          *time.Time             `json:"end_at"`
	MaxUses        *int                   `json:"max_uses"`
	MaxUsesPerUser *int                   `json:"max_uses_per_user"`
	Stackable      *bool                  `json:"stackable"`
	AutoApply      *bool                  `json:"auto_apply"`
	Description    *string                `json:"description"`
	AdminNotes     *string                `json:"admin_notes"`
}

// Create creates a new sale promotion.
func (s *SaleService) Create(req CreateSalePromotionRequest) (*SalePromotion, error) {
	if req.EndAt.Before(req.StartAt) {
		return nil, errors.New("end time must be after start time")
	}

	condJSON, err := json.Marshal(req.Condition)
	if err != nil {
		return nil, err
	}
	discJSON, err := json.Marshal(req.Discount)
	if err != nil {
		return nil, err
	}

	promo := &SalePromotion{
		Name:           req.Name,
		Code:           req.Code,
		Type:           req.Type,
		Condition:      datatypes.JSON(condJSON),
		Discount:       datatypes.JSON(discJSON),
		StartAt:        req.StartAt,
		EndAt:          req.EndAt,
		MaxUses:        req.MaxUses,
		MaxUsesPerUser: req.MaxUsesPerUser,
		Stackable:      req.Stackable,
		AutoApply:      req.AutoApply,
		Description:    req.Description,
		Status:         1,
	}

	if err := s.db.Create(promo).Error; err != nil {
		return nil, err
	}

	s.log.Infof("sale promotion created: %s (type=%s)", promo.Name, promo.Type)
	return promo, nil
}

// GetByID returns a single promotion by ID.
func (s *SaleService) GetByID(id uint) (*SalePromotion, error) {
	var promo SalePromotion
	if err := s.db.First(&promo, id).Error; err != nil {
		return nil, err
	}
	return &promo, nil
}

// GetByCode returns a promotion by code.
func (s *SaleService) GetByCode(code string) (*SalePromotion, error) {
	var promo SalePromotion
	if err := s.db.Where("code = ?", code).First(&promo).Error; err != nil {
		return nil, err
	}
	return &promo, nil
}

// Update modifies an existing promotion.
func (s *SaleService) Update(id uint, req UpdateSalePromotionRequest) (*SalePromotion, error) {
	promo, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Condition != nil {
		condJSON, _ := json.Marshal(req.Condition)
		updates["condition"] = datatypes.JSON(condJSON)
	}
	if req.Discount != nil {
		discJSON, _ := json.Marshal(req.Discount)
		updates["discount"] = datatypes.JSON(discJSON)
	}
	if req.StartAt != nil {
		updates["start_at"] = *req.StartAt
	}
	if req.EndAt != nil {
		updates["end_at"] = *req.EndAt
	}
	if req.MaxUses != nil {
		updates["max_uses"] = *req.MaxUses
	}
	if req.MaxUsesPerUser != nil {
		updates["max_uses_per_user"] = *req.MaxUsesPerUser
	}
	if req.Stackable != nil {
		updates["stackable"] = *req.Stackable
	}
	if req.AutoApply != nil {
		updates["auto_apply"] = *req.AutoApply
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.AdminNotes != nil {
		updates["admin_notes"] = *req.AdminNotes
	}

	if err := s.db.Model(promo).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

// Delete soft-deletes a promotion.
func (s *SaleService) Delete(id uint) error {
	return s.db.Delete(&SalePromotion{}, id).Error
}

// Enable sets promotion status to 1.
func (s *SaleService) Enable(id uint) error {
	return s.db.Model(&SalePromotion{}).Where("id = ?", id).Update("status", 1).Error
}

// Disable sets promotion status to 0.
func (s *SaleService) Disable(id uint) error {
	return s.db.Model(&SalePromotion{}).Where("id = ?", id).Update("status", 0).Error
}

// GetList returns all promotions with pagination.
func (s *SaleService) GetList(page, pageSize int, promoType *string, status *int) ([]SalePromotion, int64, error) {
	var promos []SalePromotion
	var total int64

	query := s.db.Model(&SalePromotion{})
	if promoType != nil {
		query = query.Where("type = ?", *promoType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&promos).Error; err != nil {
		return nil, 0, err
	}
	return promos, total, nil
}

// GetUsageStats returns usage statistics for a promotion.
func (s *SaleService) GetUsageStats(id uint) (map[string]interface{}, error) {
	promo, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	var totalDiscount float64
	var userCount int64

	s.db.Model(&SalePromotionLog{}).Where("promotion_id = ?", id).
		Select("COALESCE(SUM(discount_amount), 0)").Scan(&totalDiscount)
	s.db.Model(&SalePromotionLog{}).Where("promotion_id = ?", id).
		Distinct("user_id").Count(&userCount)

	return map[string]interface{}{
		"promotion_id":    promo.ID,
		"name":            promo.Name,
		"used_count":      promo.UsedCount,
		"max_uses":        promo.MaxUses,
		"total_discount":  totalDiscount,
		"unique_users":    userCount,
		"status":          promo.Status,
	}, nil
}

// GetUsageLogs returns paginated usage logs for a promotion.
func (s *SaleService) GetUsageLogs(promotionID uint, page, pageSize int) ([]SalePromotionLog, int64, error) {
	var logs []SalePromotionLog
	var total int64

	query := s.db.Model(&SalePromotionLog{}).Where("promotion_id = ?", promotionID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// ValidateAndApply validates a promotion code and records usage.
func (s *SaleService) ValidateAndApply(code string, userID uint, orderID uint, orderAmount float64) (float64, error) {
	promo, err := s.GetByCode(code)
	if err != nil {
		return 0, errors.New("invalid promotion code")
	}

	now := time.Now()
	if promo.Status != 1 {
		return 0, errors.New("promotion is disabled")
	}
	if now.Before(promo.StartAt) || now.After(promo.EndAt) {
		return 0, errors.New("promotion is not active")
	}
	if promo.MaxUses > 0 && promo.UsedCount >= promo.MaxUses {
		return 0, errors.New("promotion usage limit reached")
	}

	// Check per-user limit
	if promo.MaxUsesPerUser > 0 {
		var userCount int64
		s.db.Model(&SalePromotionLog{}).Where("promotion_id = ? AND user_id = ?", promo.ID, userID).Count(&userCount)
		if int(userCount) >= promo.MaxUsesPerUser {
			return 0, errors.New("user usage limit reached")
		}
	}

	// Parse discount
	var discountMap map[string]interface{}
	if err := json.Unmarshal(promo.Discount, &discountMap); err != nil {
		return 0, err
	}

	discountValue, _ := discountMap["value"].(float64)
	maxDiscount, _ := discountMap["max_discount"].(float64)

	var discountAmount float64
	switch promo.Type {
	case "amount_off":
		discountAmount = discountValue
	case "percent_off":
		discountAmount = orderAmount * discountValue / 100
		if maxDiscount > 0 && discountAmount > maxDiscount {
			discountAmount = maxDiscount
		}
	default:
		discountAmount = 0
	}

	if discountAmount > orderAmount {
		discountAmount = orderAmount
	}

	// Record usage
	log := &SalePromotionLog{
		PromotionID:    promo.ID,
		UserID:         userID,
		OrderID:        &orderID,
		DiscountAmount: discountAmount,
		UsedAt:         now,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(log).Error; err != nil {
			return err
		}
		return tx.Model(promo).Update("used_count", gorm.Expr("used_count + 1")).Error
	})

	if err != nil {
		return 0, err
	}

	s.log.Infof("promotion applied: %s (user=%d, discount=%.2f)", promo.Code, userID, discountAmount)
	return discountAmount, nil
}
