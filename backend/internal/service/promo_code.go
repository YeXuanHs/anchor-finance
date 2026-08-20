package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type PromoCodeService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewPromoCodeService(db *gorm.DB, log *logger.Logger) *PromoCodeService {
	return &PromoCodeService{db: db, log: log}
}

// Validate checks if a promo code is valid for the given user/product/amount.
func (s *PromoCodeService) Validate(code string, userID uint, productID uint, amount float64) (float64, *model.PromoCode, error) {
	var promo model.PromoCode
	if err := s.db.Where("code = ? AND status = 1", code).First(&promo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, nil, errors.New("promo code not found")
		}
		return 0, nil, err
	}

	now := time.Now().Unix()
	if promo.StartTime > 0 && now < promo.StartTime {
		return 0, nil, errors.New("promo code is not yet active")
	}
	if promo.ExpirationTime > 0 && now > promo.ExpirationTime {
		return 0, nil, errors.New("promo code has expired")
	}

	if promo.MaxTimes > 0 && promo.UsedCount >= promo.MaxTimes {
		return 0, nil, errors.New("promo code usage limit reached")
	}

	// Check new/old client restriction
	if promo.OnlyNewClient || promo.OnlyOldClient {
		var orderCount int64
		s.db.Model(&model.Order{}).Where("user_id = ? AND status IN ?", userID, []int{1, 2, 3}).Count(&orderCount)
		isNewClient := orderCount == 0
		if promo.OnlyNewClient && !isNewClient {
			return 0, nil, errors.New("promo code is only available for new clients")
		}
		if promo.OnlyOldClient && isNewClient {
			return 0, nil, errors.New("promo code is only available for existing clients")
		}
	}

	// Check required product existence (Requires field = comma-separated product IDs)
	if promo.Requires != "" && promo.RequiresExist {
		requiredIDs := stringToUintSlice(promo.Requires)
		for _, reqID := range requiredIDs {
			var orderCount int64
			s.db.Model(&model.Order{}).Where("user_id = ? AND product_id = ? AND status IN ?", userID, reqID, []int{1, 2, 3}).Count(&orderCount)
			if orderCount == 0 {
				return 0, nil, errors.New("promo code requires an existing active product")
			}
		}
	}

	// Check product restriction (AppliesTo = comma-separated product IDs)
	if promo.AppliesTo != "" && productID > 0 {
		appliesTo := stringToUintSlice(promo.AppliesTo)
		if len(appliesTo) > 0 {
			found := false
			for _, pid := range appliesTo {
				if pid == productID {
					found = true
					break
				}
			}
			if !found {
				return 0, nil, errors.New("promo code is not applicable to this product")
			}
		}
	}

	// Check per-client usage
	if promo.OncePerClient {
		var userUsage int64
		s.db.Model(&model.PromoCodeLog{}).Where("promo_id = ? AND user_id = ?", promo.ID, userID).Count(&userUsage)
		if userUsage > 0 {
			return 0, nil, errors.New("promo code can only be used once per client")
		}
	}

	// Calculate discount
	var discount float64
	switch promo.Type {
	case "fixed":
		discount = promo.Value
	case "percent":
		discount = amount * promo.Value / 100
	case "override":
		discount = amount - promo.Value
		if discount < 0 {
			discount = 0
		}
	case "free":
		discount = amount
	default:
		return 0, nil, fmt.Errorf("unknown promo code type: %s", promo.Type)
	}

	if discount > amount {
		discount = amount
	}
	if discount < 0 {
		discount = 0
	}

	return discount, &promo, nil
}

// Apply records a promo code usage against an order.
func (s *PromoCodeService) Apply(promoID, userID, orderID uint, discount float64) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.PromoCode{}).Where("id = ?", promoID).
			UpdateColumn("used_count", gorm.Expr("used_count + 1")).Error; err != nil {
			return err
		}
		return tx.Create(&model.PromoCodeLog{
			PromoID: promoID,
			UserID:  userID,
			OrderID: orderID,
			Amount:  discount,
		}).Error
	})
}

// GetUserPromoCodes returns available promo codes for a user.
func (s *PromoCodeService) GetUserPromoCodes(userID uint) ([]model.PromoCode, error) {
	var promos []model.PromoCode
	now := time.Now().Unix()

	err := s.db.Where("status = 1 AND (start_time = 0 OR start_time <= ?) AND (expiration_time = 0 OR expiration_time >= ?)", now, now).
		Find(&promos).Error
	if err != nil {
		return nil, err
	}

	// Filter out exhausted promo codes and client-type restrictions
	var orderCount int64
	s.db.Model(&model.Order{}).Where("user_id = ? AND status IN ?", userID, []int{1, 2, 3}).Count(&orderCount)
	isNewClient := orderCount == 0

	var result []model.PromoCode
	for _, p := range promos {
		if p.MaxTimes > 0 && p.UsedCount >= p.MaxTimes {
			continue
		}
		if p.OnlyNewClient && !isNewClient {
			continue
		}
		if p.OnlyOldClient && isNewClient {
			continue
		}
		result = append(result, p)
	}
	return result, nil
}

// stringToUintSlice parses a comma-separated string into a uint slice.
func stringToUintSlice(s string) []uint {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]uint, 0, len(parts))
	for _, p := range parts {
		if v, err := strconv.ParseUint(strings.TrimSpace(p), 10, 32); err == nil {
			result = append(result, uint(v))
		}
	}
	return result
}
