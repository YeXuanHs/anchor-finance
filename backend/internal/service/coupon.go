package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CouponService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCouponService(db *gorm.DB, log *logger.Logger) *CouponService {
	return &CouponService{db: db, log: log}
}

// Validate checks if a coupon code is valid for the given user/product/amount/cycle.
func (s *CouponService) Validate(code string, userID, productID uint, amount float64, cycle string, isNewClient bool) (float64, *model.Coupon, error) {
	var coupon model.Coupon
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

	// Check new/old client restriction
	if coupon.OnlyNewClient && !isNewClient {
		return 0, nil, errors.New("coupon is only available for new clients")
	}
	if coupon.OnlyOldClient && isNewClient {
		return 0, nil, errors.New("coupon is only available for existing clients")
	}

	// Check required product existence
	if coupon.RequiresExist != nil && *coupon.RequiresExist > 0 {
		var orderCount int64
		s.db.Model(&model.Order{}).Where("user_id = ? AND product_id = ? AND status IN ?", userID, *coupon.RequiresExist, []string{"active", "paid"}).Count(&orderCount)
		if orderCount == 0 {
			return 0, nil, errors.New("coupon requires an existing active product")
		}
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

	// Check billing cycle restriction
	if coupon.Cycles != nil && cycle != "" {
		var cycles []string
		if err := coupon.Cycles.Scan(&cycles); err == nil && len(cycles) > 0 {
			found := false
			for _, c := range cycles {
				if c == cycle {
					found = true
					break
				}
			}
			if !found {
				return 0, nil, errors.New("coupon is not applicable to this billing cycle")
			}
		}
	}

	// Check minimum order amount
	if coupon.MinOrderAmount > 0 && amount < float64(coupon.MinOrderAmount) {
		return 0, nil, fmt.Errorf("order amount %.2f is below minimum %.2f", amount, float64(coupon.MinOrderAmount))
	}

	// Check per-user usage
	var userUsage int64
	s.db.Model(&model.CouponUsageLog{}).Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).Count(&userUsage)
	if coupon.OncePerClient && userUsage > 0 {
		return 0, nil, errors.New("coupon can only be used once per client")
	}
	if coupon.MaxUsesPerUser > 0 && int(userUsage) >= coupon.MaxUsesPerUser {
		return 0, nil, errors.New("per-user coupon usage limit reached")
	}

	// Calculate discount
	var discount float64
	switch coupon.Type {
	case "fixed":
		discount = float64(coupon.Value)
	case "percentage":
		discount = amount * float64(coupon.Value) / 100
		if coupon.MaxDiscount > 0 && discount > float64(coupon.MaxDiscount) {
			discount = float64(coupon.MaxDiscount)
		}
	case "override":
		// Override: the value is the final price, discount is the difference
		discount = amount - float64(coupon.Value)
		if discount < 0 {
			discount = 0
		}
	case "free":
		discount = amount
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
		if err := tx.Model(&model.Coupon{}).Where("id = ?", couponID).
			UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error; err != nil {
			return err
		}
		return tx.Create(&model.CouponUsageLog{
			CouponID: couponID,
			UserID:   userID,
			OrderID:  orderID,
			Discount: discount,
		}).Error
	})
}

// GetUserCoupons returns available coupons for a user.
func (s *CouponService) GetUserCoupons(userID uint, isNewClient bool) ([]model.Coupon, error) {
	var coupons []model.Coupon
	now := time.Now()

	err := s.db.Where("status = 1 AND (start_date IS NULL OR start_date <= ?) AND (end_date IS NULL OR end_date >= ?)", now, now).
		Find(&coupons).Error
	if err != nil {
		return nil, err
	}

	// Filter out exhausted coupons and client-type restrictions
	var result []model.Coupon
	for _, c := range coupons {
		if c.MaxUses > 0 && c.UsedCount >= c.MaxUses {
			continue
		}
		if c.OnlyNewClient && !isNewClient {
			continue
		}
		if c.OnlyOldClient && isNewClient {
			continue
		}
		result = append(result, c)
	}
	return result, nil
}

// Create creates a new coupon (admin).
func (s *CouponService) Create(coupon *model.Coupon) error {
	var count int64
	s.db.Model(&model.Coupon{}).Where("code = ?", coupon.Code).Count(&count)
	if count > 0 {
		return errors.New("coupon code already exists")
	}
	return s.db.Create(coupon).Error
}

// Update updates a coupon (admin).
func (s *CouponService) Update(coupon *model.Coupon) error {
	return s.db.Save(coupon).Error
}

// Delete soft-deletes a coupon (admin).
func (s *CouponService) Delete(id uint) error {
	result := s.db.Delete(&model.Coupon{}, id)
	if result.RowsAffected == 0 {
		return errors.New("coupon not found")
	}
	return result.Error
}

// GetByID retrieves a coupon by ID.
func (s *CouponService) GetByID(id uint) (*model.Coupon, error) {
	var coupon model.Coupon
	if err := s.db.First(&coupon, id).Error; err != nil {
		return nil, err
	}
	return &coupon, nil
}

// GetAll returns all coupons with pagination (admin).
func (s *CouponService) GetAll(page, pageSize int) ([]model.Coupon, int64, error) {
	var coupons []model.Coupon
	var total int64

	s.db.Model(&model.Coupon{}).Count(&total)
	offset, limit := Paginate(page, pageSize)
	err := s.db.Offset(offset).Limit(limit).Order("id DESC").Find(&coupons).Error
	return coupons, total, err
}

// IncrementUsage increments the used count of a coupon.
func (s *CouponService) IncrementUsage(couponID uint) error {
	return s.db.Model(&model.Coupon{}).Where("id = ?", couponID).
		UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error
}

// datatypesJSONToUintSlice is a helper to parse JSON array to uint slice.
func datatypesJSONToUintSlice(data datatypes.JSON, out *[]uint) error {
	if data == nil {
		return nil
	}
	return data.Scan(out)
}
