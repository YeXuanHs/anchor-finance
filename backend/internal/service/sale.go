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
	OnlyNewClient   bool           `gorm:"default:false" json:"only_new_client"`
	OnlyOldClient   bool           `gorm:"default:false" json:"only_old_client"`
	OncePerClient   bool           `gorm:"default:false" json:"once_per_client"`
	RequiresExist   *uint          `json:"requires_exist"`                   // product ID that must exist for user
	Recurring       bool           `gorm:"default:false" json:"recurring"`
	RecurFor        int            `gorm:"default:0" json:"recur_for"`       // number of billing cycles
	IsDiscount      bool           `gorm:"default:false" json:"is_discount"`
	UpgradeType     string         `gorm:"type:varchar(32)" json:"upgrade_type"` // percentage, fixed
	UpgradeValue    float64        `gorm:"type:decimal(20,4)" json:"upgrade_value"`
	UpgradeOptions  datatypes.JSON `gorm:"type:jsonb" json:"upgrade_options"`
	AppliesToGroups datatypes.JSON `gorm:"type:jsonb" json:"applies_to_groups"` // user group IDs
	Commissions     []SaleCommission `gorm:"foreignKey:PromotionID" json:"commissions,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// SaleCommission 代理佣金阶梯
type SaleCommission struct {
	gorm.Model
	PromotionID uint           `gorm:"index" json:"promotion_id"`
	MinAmount   float64        `gorm:"type:decimal(20,4)" json:"min_amount"`
	MaxAmount   float64        `gorm:"type:decimal(20,4)" json:"max_amount"`
	Rate        float64        `gorm:"type:decimal(10,4)" json:"rate"`        // percentage
	FixedAmount float64        `gorm:"type:decimal(20,4)" json:"fixed_amount"`
	Promotion   *SalePromotion `gorm:"foreignKey:PromotionID" json:"promotion,omitempty"`
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
	OnlyNewClient  bool                   `json:"only_new_client"`
	OnlyOldClient  bool                   `json:"only_old_client"`
	OncePerClient  bool                   `json:"once_per_client"`
	RequiresExist  *uint                  `json:"requires_exist"`
	Recurring      bool                   `json:"recurring"`
	RecurFor       int                    `json:"recur_for"`
	IsDiscount     bool                   `json:"is_discount"`
	UpgradeType    string                 `json:"upgrade_type"`
	UpgradeValue   float64                `json:"upgrade_value"`
	UpgradeOptions map[string]interface{} `json:"upgrade_options"`
	AppliesToGroups []uint               `json:"applies_to_groups"`
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
	OnlyNewClient  *bool                  `json:"only_new_client"`
	OnlyOldClient  *bool                  `json:"only_old_client"`
	OncePerClient  *bool                  `json:"once_per_client"`
	RequiresExist  *uint                  `json:"requires_exist"`
	Recurring      *bool                  `json:"recurring"`
	RecurFor       *int                   `json:"recur_for"`
	IsDiscount     *bool                  `json:"is_discount"`
	UpgradeType    *string                `json:"upgrade_type"`
	UpgradeValue   *float64               `json:"upgrade_value"`
	UpgradeOptions map[string]interface{} `json:"upgrade_options"`
	AppliesToGroups []uint               `json:"applies_to_groups"`
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
		OnlyNewClient:  req.OnlyNewClient,
		OnlyOldClient:  req.OnlyOldClient,
		OncePerClient:  req.OncePerClient,
		RequiresExist:  req.RequiresExist,
		Recurring:      req.Recurring,
		RecurFor:       req.RecurFor,
		IsDiscount:     req.IsDiscount,
		UpgradeType:    req.UpgradeType,
		UpgradeValue:   req.UpgradeValue,
	}

	if req.UpgradeOptions != nil {
		optJSON, _ := json.Marshal(req.UpgradeOptions)
		promo.UpgradeOptions = datatypes.JSON(optJSON)
	}
	if len(req.AppliesToGroups) > 0 {
		grpJSON, _ := json.Marshal(req.AppliesToGroups)
		promo.AppliesToGroups = datatypes.JSON(grpJSON)
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
	if req.OnlyNewClient != nil {
		updates["only_new_client"] = *req.OnlyNewClient
	}
	if req.OnlyOldClient != nil {
		updates["only_old_client"] = *req.OnlyOldClient
	}
	if req.OncePerClient != nil {
		updates["once_per_client"] = *req.OncePerClient
	}
	if req.RequiresExist != nil {
		updates["requires_exist"] = *req.RequiresExist
	}
	if req.Recurring != nil {
		updates["recurring"] = *req.Recurring
	}
	if req.RecurFor != nil {
		updates["recur_for"] = *req.RecurFor
	}
	if req.IsDiscount != nil {
		updates["is_discount"] = *req.IsDiscount
	}
	if req.UpgradeType != nil {
		updates["upgrade_type"] = *req.UpgradeType
	}
	if req.UpgradeValue != nil {
		updates["upgrade_value"] = *req.UpgradeValue
	}
	if req.UpgradeOptions != nil {
		optJSON, _ := json.Marshal(req.UpgradeOptions)
		updates["upgrade_options"] = datatypes.JSON(optJSON)
	}
	if req.AppliesToGroups != nil {
		grpJSON, _ := json.Marshal(req.AppliesToGroups)
		updates["applies_to_groups"] = datatypes.JSON(grpJSON)
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

// SetStatus sets the status of a sale promotion (1=enabled, 0=disabled).
func (s *SaleService) SetStatus(id uint, status int) error {
	result := s.db.Model(&SalePromotion{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("promotion not found")
	}
	s.log.Infof("sale promotion status updated: id=%d status=%d", id, status)
	return nil
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

// ─── Commission Ladder ───

// GetCommissionLadder returns all commission tiers for a promotion.
func (s *SaleService) GetCommissionLadder(promotionID uint) ([]SaleCommission, error) {
	var tiers []SaleCommission
	if err := s.db.Where("promotion_id = ?", promotionID).Order("min_amount ASC").Find(&tiers).Error; err != nil {
		return nil, err
	}
	return tiers, nil
}

// SetCommissionLadder replaces all commission tiers for a promotion.
func (s *SaleService) SetCommissionLadder(promotionID uint, tiers []SaleCommission) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("promotion_id = ?", promotionID).Delete(&SaleCommission{}).Error; err != nil {
			return err
		}
		for i := range tiers {
			tiers[i].PromotionID = promotionID
			tiers[i].ID = 0
			if err := tx.Create(&tiers[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// CalculateCommission calculates the commission for a given amount based on the promotion's ladder.
func (s *SaleService) CalculateCommission(promotionID uint, amount float64) (float64, error) {
	tiers, err := s.GetCommissionLadder(promotionID)
	if err != nil {
		return 0, err
	}
	if len(tiers) == 0 {
		return 0, nil
	}

	for _, tier := range tiers {
		if amount >= tier.MinAmount && (tier.MaxAmount == 0 || amount <= tier.MaxAmount) {
			commission := amount*tier.Rate/100 + tier.FixedAmount
			return commission, nil
		}
	}
	return 0, nil
}

// ValidateUsage checks if a user can use a promotion based on client type and product ownership.
func (s *SaleService) ValidateUsage(promotionID uint, userID uint) (bool, string, error) {
	promo, err := s.GetByID(promotionID)
	if err != nil {
		return false, "promotion not found", err
	}

	// Check client type restrictions
	if promo.OnlyNewClient || promo.OnlyOldClient {
		var orderCount int64
		if err := s.db.Table("orders").Where("user_id = ? AND deleted_at IS NULL", userID).Count(&orderCount).Error; err != nil {
			return false, "failed to check order history", err
		}
		isNew := orderCount == 0
		if promo.OnlyNewClient && !isNew {
			return false, "promotion is only for new clients", nil
		}
		if promo.OnlyOldClient && isNew {
			return false, "promotion is only for existing clients", nil
		}
	}

	// Check once per client
	if promo.OncePerClient {
		var usageCount int64
		s.db.Model(&SalePromotionLog{}).Where("promotion_id = ? AND user_id = ?", promotionID, userID).Count(&usageCount)
		if usageCount > 0 {
			return false, "promotion can only be used once per client", nil
		}
	}

	// Check requires exist (user must own a specific product)
	if promo.RequiresExist != nil && *promo.RequiresExist > 0 {
		var productCount int64
		if err := s.db.Table("user_products").Where("user_id = ? AND product_id = ? AND deleted_at IS NULL", userID, *promo.RequiresExist).Count(&productCount).Error; err != nil {
			return false, "failed to check product ownership", err
		}
		if productCount == 0 {
			return false, "promotion requires ownership of a specific product", nil
		}
	}

	return true, "", nil
}
