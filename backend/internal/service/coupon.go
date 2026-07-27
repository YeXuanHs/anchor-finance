package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Coupon 优惠券 (service layer type)
type Coupon struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Code           string         `gorm:"uniqueIndex;size:64;not null" json:"code"`
	Name           string         `gorm:"size:128;not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	Type           string         `gorm:"size:16;not null;default:percentage" json:"type"` // percentage/fixed
	Value          float64        `gorm:"type:decimal(12,2);not null" json:"value"`
	MaxDiscount    float64        `gorm:"type:decimal(12,2)" json:"max_discount"`
	MinOrderAmount float64        `gorm:"type:decimal(12,2);default:0" json:"min_order_amount"`
	MaxUses        int            `gorm:"default:0" json:"max_uses"` // 0=unlimited
	UsedCount      int            `gorm:"default:0" json:"used_count"`
	MaxUsesPerUser int            `gorm:"default:1" json:"max_uses_per_user"`
	StartDate      *time.Time     `gorm:"index" json:"start_date"`
	EndDate        *time.Time     `gorm:"index" json:"end_date"`
	ProductIDs     datatypes.JSON `gorm:"type:jsonb" json:"product_ids"`
	GroupIDs       datatypes.JSON `gorm:"type:jsonb" json:"group_ids"`
	Status         int16          `gorm:"default:1;index" json:"status"` // 1=enabled 0=disabled
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// CouponUsageLog 优惠券使用记录
type CouponUsageLog struct {
	ID       uint    `gorm:"primaryKey" json:"id"`
	CouponID uint    `gorm:"index;not null" json:"coupon_id"`
	UserID   uint    `gorm:"index;not null" json:"user_id"`
	OrderID  uint    `gorm:"index" json:"order_id"`
	Discount float64 `gorm:"type:decimal(12,2);not null" json:"discount"`
}

type CouponService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCouponService(db *gorm.DB, log *logger.Logger) *CouponService {
	return &CouponService{db: db, log: log}
}

// Validate checks if a coupon code is valid for the given user/product/amount.
func (s *CouponService) Validate(code string, userID, productID uint, amount float64) (float64, *Coupon, error) {
	var coupon Coupon
	if err := s.db.Where("code = ? AND status = 1", code).First(&coupon).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil, errors.New("coupon not found")
		}
		return 0, nil, err
	}

	now := time.Now()
	if coupon.StartDate != nil && now.Before(*coupon.StartDate) {
		return 0, nil, errors.New("coupon is not yet active")
	}
	if coupon.EndDate != nil && now.After(*coupon.EndDate) {
		return 0, nil, errors.New("coupon has expired")
	}

	if coupon.MaxUses > 0 && coupon.UsedCount >= coupon.MaxUses {
		return 0, nil, errors.New("coupon usage limit reached")
	}

	// Check product restriction
	if coupon.ProductIDs != nil {
		var productIDs []uint
		if err := datatypesJSONToUintSlice(coupon.ProductIDs, &productIDs); err == nil && len(productIDs) > 0 {
			found := false
			for _, pid := range productIDs {
				if pid == productID {
					found = true
					break
				}
			}
			if !found {
				return 0, nil, errors.New("coupon is not applicable to this product")
			}
		}
	}

	// Check minimum order amount
	if coupon.MinOrderAmount > 0 && amount < coupon.MinOrderAmount {
		return 0, nil, fmt.Errorf("order amount %.2f is below minimum %.2f", amount, coupon.MinOrderAmount)
	}

	// Check per-user usage
	var userUsage int64
	s.db.Model(&CouponUsageLog{}).Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).Count(&userUsage)
	if coupon.MaxUsesPerUser > 0 && int(userUsage) >= coupon.MaxUsesPerUser {
		return 0, nil, errors.New("per-user coupon usage limit reached")
	}

	// Calculate discount
	var discount float64
	switch coupon.Type {
	case "fixed":
		discount = coupon.Value
	case "percentage":
		discount = amount * coupon.Value / 100
		if coupon.MaxDiscount > 0 && discount > coupon.MaxDiscount {
			discount = coupon.MaxDiscount
		}
	default:
		return 0, nil, fmt.Errorf("unknown coupon type: %s", coupon.Type)
	}

	if discount > amount {
		discount = amount
	}
	if discount < 0 {
		discount = 0
	}

	return discount, &coupon, nil
}

// Apply records a coupon usage against an order.
func (s *CouponService) Apply(couponID, userID, orderID uint, discount float64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Coupon{}).Where("id = ?", couponID).
			UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error; err != nil {
			return err
		}
		return tx.Create(&CouponUsageLog{
			CouponID: couponID,
			UserID:   userID,
			OrderID:  orderID,
			Discount: discount,
		}).Error
	})
}

// GetUserCoupons returns available coupons for a user.
func (s *CouponService) GetUserCoupons(userID uint) ([]Coupon, error) {
	var coupons []Coupon
	now := time.Now()

	err := s.db.Where("status = 1 AND (start_date IS NULL OR start_date <= ?) AND (end_date IS NULL OR end_date >= ?)", now, now).
		Find(&coupons).Error
	if err != nil {
		return nil, err
	}

	// Filter out exhausted coupons
	var result []Coupon
	for _, c := range coupons {
		if c.MaxUses > 0 && c.UsedCount >= c.MaxUses {
			continue
		}
		result = append(result, c)
	}
	return result, nil
}

// Create creates a new coupon (admin).
func (s *CouponService) Create(coupon *Coupon) error {
	var count int64
	s.db.Model(&Coupon{}).Where("code = ?", coupon.Code).Count(&count)
	if count > 0 {
		return errors.New("coupon code already exists")
	}
	return s.db.Create(coupon).Error
}

// Update updates a coupon (admin).
func (s *CouponService) Update(coupon *Coupon) error {
	return s.db.Save(coupon).Error
}

// Delete soft-deletes a coupon (admin).
func (s *CouponService) Delete(id uint) error {
	result := s.db.Delete(&Coupon{}, id)
	if result.RowsAffected == 0 {
		return errors.New("coupon not found")
	}
	return result.Error
}

// GetByID retrieves a coupon by ID.
func (s *CouponService) GetByID(id uint) (*Coupon, error) {
	var coupon Coupon
	if err := s.db.First(&coupon, id).Error; err != nil {
		return nil, err
	}
	return &coupon, nil
}

// GetAll returns all coupons with pagination (admin).
func (s *CouponService) GetAll(page, pageSize int) ([]Coupon, int64, error) {
	var coupons []Coupon
	var total int64

	s.db.Model(&Coupon{}).Count(&total)
	offset, limit := Paginate(page, pageSize)
	err := s.db.Offset(offset).Limit(limit).Order("id DESC").Find(&coupons).Error
	return coupons, total, err
}

// IncrementUsage increments the used count of a coupon.
func (s *CouponService) IncrementUsage(couponID uint) error {
	return s.db.Model(&Coupon{}).Where("id = ?", couponID).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}

// datatypesJSONToUintSlice is a helper to parse JSON array to uint slice.
func datatypesJSONToUintSlice(data datatypes.JSON, out *[]uint) error {
	if data == nil {
		return nil
	}
	return data.Scan(out)
}
