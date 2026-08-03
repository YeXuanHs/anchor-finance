package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// V10CartItem V10购物车条目
type V10CartItem struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	ProductID    uint           `gorm:"index;not null" json:"product_id"`
	Product      Product        `gorm:"foreignKey:ProductID" json:"product"`
	BillingCycle string         `gorm:"type:varchar(32);not null" json:"billing_cycle"`
	Quantity     int            `gorm:"default:1" json:"quantity"`
	Config       datatypes.JSON `gorm:"type:json" json:"config"`
	Domain       string         `gorm:"type:varchar(255)" json:"domain"`
	CouponCode   string         `gorm:"type:varchar(64)" json:"coupon_code"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func (V10CartItem) TableName() string {
	return "v10_cart_items"
}

type V10CartService struct {
	db    *gorm.DB
	log   *logger.Logger
	oSvc  *OrderService
	pSvc  *PromoCodeService
}

func NewV10CartService(db *gorm.DB, log *logger.Logger, oSvc *OrderService, pSvc *PromoCodeService) *V10CartService {
	return &V10CartService{db: db, log: log, oSvc: oSvc, pSvc: pSvc}
}

type V10AddToCartRequest struct {
	ProductID    uint           `json:"product_id" binding:"required"`
	BillingCycle string        `json:"billing_cycle" binding:"required"`
	Quantity     int           `json:"quantity"`
	Config       datatypes.JSON `json:"config"`
	Domain       string         `json:"domain"`
}

type V10UpdateCartRequest struct {
	Quantity int            `json:"quantity"`
	Config   datatypes.JSON `json:"config"`
	Domain   string         `json:"domain"`
}

type V10ApplyCouponRequest struct {
	CouponCode string `json:"coupon_code" binding:"required"`
}

type V10CartSummary struct {
	Items        []V10CartItem `json:"items"`
	ItemCount    int           `json:"item_count"`
	SubTotal     float64       `json:"sub_total"`
	Discount     float64       `json:"discount"`
	Total        float64       `json:"total"`
	CouponCode   string        `json:"coupon_code"`
	Currency     string        `json:"currency"`
}

// GetCart returns all items in the user's cart with summary.
func (s *V10CartService) GetCart(userID uint) (*V10CartSummary, error) {
	var items []V10CartItem
	if err := s.db.Preload("Product").Where("user_id = ?", userID).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}

	summary := &V10CartSummary{
		Items:     items,
		ItemCount: len(items),
		Currency:  "CNY",
	}

	var couponCode string
	for _, item := range items {
		price := float64(item.Product.Price) * float64(item.Quantity)
		summary.SubTotal += price
		if item.CouponCode != "" {
			couponCode = item.CouponCode
		}
	}

	// Apply coupon discount
	if couponCode != "" && s.pSvc != nil {
		summary.CouponCode = couponCode
		for _, item := range items {
			discount, _, err := s.pSvc.Validate(couponCode, userID, item.ProductID, float64(item.Product.Price)*float64(item.Quantity))
			if err == nil {
				summary.Discount += discount
			}
		}
	}

	summary.Total = summary.SubTotal - summary.Discount
	if summary.Total < 0 {
		summary.Total = 0
	}

	return summary, nil
}

// AddItem adds a product to the V10 cart.
func (s *V10CartService) AddItem(userID, productID uint, billingCycle string, config datatypes.JSON, quantity int, domain string) (*V10CartItem, error) {
	if quantity < 1 {
		quantity = 1
	}

	// Verify product
	var product Product
	if err := s.db.First(&product, productID).Error; err != nil {
		return nil, errors.New("product not found")
	}
	if product.Status != 1 {
		return nil, errors.New("product is disabled")
	}

	// Check for duplicate
	var existing V10CartItem
	err := s.db.Where("user_id = ? AND product_id = ? AND billing_cycle = ?",
		userID, productID, billingCycle).First(&existing).Error

	if err == nil {
		existing.Quantity += quantity
		if config != nil {
			existing.Config = config
		}
		if domain != "" {
			existing.Domain = domain
		}
		if err := s.db.Save(&existing).Error; err != nil {
			return nil, err
		}
		if err := s.db.Preload("Product").First(&existing, existing.ID).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item := &V10CartItem{
		UserID:       userID,
		ProductID:    productID,
		BillingCycle: billingCycle,
		Quantity:     quantity,
		Config:       config,
		Domain:       domain,
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	if err := s.db.Preload("Product").First(item, item.ID).Error; err != nil {
		return nil, err
	}

	s.log.Infof("v10 cart item added: user=%d product=%d", userID, productID)
	return item, nil
}

// UpdateItem updates a cart item.
func (s *V10CartService) UpdateItem(userID, itemID uint, req V10UpdateCartRequest) (*V10CartItem, error) {
	var item V10CartItem
	if err := s.db.Where("id = ? AND user_id = ?", itemID, userID).First(&item).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	if req.Quantity > 0 {
		item.Quantity = req.Quantity
	}
	if req.Config != nil {
		item.Config = req.Config
	}
	if req.Domain != "" {
		item.Domain = req.Domain
	}

	if err := s.db.Save(&item).Error; err != nil {
		return nil, err
	}
	if err := s.db.Preload("Product").First(&item, item.ID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// RemoveItem removes an item from the cart.
func (s *V10CartService) RemoveItem(userID, itemID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", itemID, userID).Delete(&V10CartItem{})
	if result.RowsAffected == 0 {
		return errors.New("cart item not found")
	}
	return result.Error
}

// ClearCart removes all items from the user's cart.
func (s *V10CartService) ClearCart(userID uint) error {
	return s.db.Where("user_id = ?", userID).Delete(&V10CartItem{}).Error
}

// ApplyCoupon applies a coupon code to all cart items.
func (s *V10CartService) ApplyCoupon(userID uint, couponCode string) error {
	if s.pSvc == nil {
		return errors.New("promo code service not available")
	}

	// Validate promo code exists
	var count int64
	s.db.Table("promo_codes").Where("code = ? AND status = 1", couponCode).Count(&count)
	if count == 0 {
		return errors.New("invalid promo code")
	}

	// Apply to all cart items
	if err := s.db.Model(&V10CartItem{}).Where("user_id = ?", userID).
		Update("coupon_code", couponCode).Error; err != nil {
		return err
	}

	s.log.Infof("promo code applied to cart: user=%d code=%s", userID, couponCode)
	return nil
}

// RemoveCoupon removes the coupon from all cart items.
func (s *V10CartService) RemoveCoupon(userID uint) error {
	return s.db.Model(&V10CartItem{}).Where("user_id = ?", userID).
		Update("coupon_code", "").Error
}

// Checkout creates orders from cart items.
func (s *V10CartService) Checkout(userID uint) ([]interface{}, error) {
	var items []V10CartItem
	if err := s.db.Preload("Product").Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}
	if len(items) == 0 {
		return nil, errors.New("cart is empty")
	}

	var orders []interface{}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			totalPrice := float64(item.Product.Price) * float64(item.Quantity)

			// Apply coupon if present
			if item.CouponCode != "" && s.pSvc != nil {
				discount, _, err := s.pSvc.Validate(item.CouponCode, userID, item.ProductID, totalPrice)
				if err == nil {
					totalPrice -= discount
					if totalPrice < 0 {
						totalPrice = 0
					}
				}
			}

			order := map[string]interface{}{
				"order_no":     util.GenerateOrderNo(),
				"user_id":      userID,
				"product_id":   item.ProductID,
				"quantity":     item.Quantity,
				"total_price":  totalPrice,
				"billing_cycle": item.BillingCycle,
				"domain":       item.Domain,
				"status":       0,
			}

			// Create order record
			var orderRecord struct {
				ID uint
			}
			if err := tx.Table("orders").Create(&order).Scan(&orderRecord).Error; err != nil {
				return err
			}

			orders = append(orders, order)
			s.log.Infof("v10 order created: %s (user=%d, product=%d)", order["order_no"], userID, item.ProductID)
		}

		// Clear cart
		if err := tx.Where("user_id = ?", userID).Delete(&V10CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetItemCount returns the number of items in the cart.
func (s *V10CartService) GetItemCount(userID uint) (int, error) {
	var count int64
	if err := s.db.Model(&V10CartItem{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return int(count), nil
}
