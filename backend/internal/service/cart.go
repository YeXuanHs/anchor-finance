package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/anchor-finance/backend/internal/util"
	"github.com/anchor-finance/backend/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CartItem represents a shopping cart item with product info.
type CartItem struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	ProductID    uint           `gorm:"index;not null" json:"product_id"`
	Product      Product        `gorm:"foreignKey:ProductID" json:"product"`
	BillingCycle string         `gorm:"size:32;not null" json:"billing_cycle"`
	Quantity     int            `gorm:"default:1" json:"quantity"`
	Config       datatypes.JSON `gorm:"type:jsonb" json:"config"`
	Domain       string         `gorm:"size:255" json:"domain"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type CartService struct {
	db    *gorm.DB
	log   *logger.Logger
	oSvc  *OrderService
	cSvc  *CouponService
}

func NewCartService(db *gorm.DB, log *logger.Logger, oSvc *OrderService, cSvc *CouponService) *CartService {
	return &CartService{db: db, log: log, oSvc: oSvc, cSvc: cSvc}
}

// GetCart returns all cart items for a user with product info.
func (s *CartService) GetCart(userID uint) ([]CartItem, error) {
	var items []CartItem
	err := s.db.Preload("Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&items).Error
	return items, err
}

// AddItem adds a product to cart. If same product+cycle exists, increments quantity.
func (s *CartService) AddItem(userID, productID uint, billingCycle string, config datatypes.JSON, quantity int) (*CartItem, error) {
	if quantity < 1 {
		quantity = 1
	}

	// Verify product exists
	var product Product
	if err := s.db.First(&product, productID).Error; err != nil {
		return nil, errors.New("product not found")
	}
	if product.Status != 1 {
		return nil, errors.New("product is disabled")
	}

	// Check for duplicate (same product + billing cycle)
	var existing CartItem
	err := s.db.Where("user_id = ? AND product_id = ? AND billing_cycle = ?",
		userID, productID, billingCycle).First(&existing).Error

	if err == nil {
		// Increment quantity on existing item
		existing.Quantity += quantity
		if config != nil {
			existing.Config = config
		}
		if err := s.db.Save(&existing).Error; err != nil {
			return nil, err
		}
		if err := s.db.Preload("Product").First(&existing, existing.ID).Error; err != nil {
			return nil, err
		}
		s.log.Infof("cart item updated: user=%d product=%d qty=%d", userID, productID, existing.Quantity)
		return &existing, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Create new cart item
	item := &CartItem{
		UserID:       userID,
		ProductID:    productID,
		BillingCycle: billingCycle,
		Quantity:     quantity,
		Config:       config,
	}
	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	if err := s.db.Preload("Product").First(item, item.ID).Error; err != nil {
		return nil, err
	}
	s.log.Infof("cart item added: user=%d product=%d", userID, productID)
	return item, nil
}

// UpdateItem updates a cart item's quantity and config.
func (s *CartService) UpdateItem(id, userID uint, quantity int, config datatypes.JSON) (*CartItem, error) {
	var item CartItem
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	if quantity > 0 {
		item.Quantity = quantity
	}
	if config != nil {
		item.Config = config
	}

	if err := s.db.Save(&item).Error; err != nil {
		return nil, err
	}
	if err := s.db.Preload("Product").First(&item, item.ID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// RemoveItem removes a cart item.
func (s *CartService) RemoveItem(id, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&CartItem{})
	if result.RowsAffected == 0 {
		return errors.New("cart item not found")
	}
	return result.Error
}

// ClearCart removes all items from a user's cart.
func (s *CartService) ClearCart(userID uint) error {
	return s.db.Where("user_id = ?", userID).Delete(&CartItem{}).Error
}

// Checkout creates orders from cart items and applies coupon if provided.
func (s *CartService) Checkout(userID uint, couponCode string) ([]*Order, error) {
	var items []CartItem
	if err := s.db.Preload("Product").Where("user_id = ?", userID).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %w", err)
	}
	if len(items) == 0 {
		return nil, errors.New("cart is empty")
	}

	var orders []*Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			// Apply coupon discount if provided
			totalPrice := item.Product.Price * float64(item.Quantity)
			var couponID *uint
			if couponCode != "" && s.cSvc != nil {
				discount, coupon, err := s.cSvc.Validate(couponCode, userID, item.ProductID, totalPrice)
				if err == nil && coupon != nil {
					totalPrice -= discount
					if totalPrice < 0 {
						totalPrice = 0
					}
					cid := coupon.ID
					couponID = &cid
					// Record coupon usage
					_ = s.cSvc.Apply(coupon.ID, userID, 0, discount) // orderID will be set after creation
				}
			}

			order := &Order{
				OrderNo:    util.GenerateOrderNo(),
				UserID:     userID,
				ProductID:  item.ProductID,
				Quantity:   item.Quantity,
				TotalPrice: totalPrice,
				Period:     item.Product.Period,
				PeriodUnit: item.Product.PeriodUnit,
				Status:     0,
				Remark:     fmt.Sprintf("Cart checkout - %s", item.BillingCycle),
			}
			if err := tx.Create(order).Error; err != nil {
				return err
			}

			// Update stock if controlled
			if item.Product.Stock >= 0 {
				if err := tx.Model(&item.Product).Update("stock",
					gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
					return err
				}
			}

			orders = append(orders, order)
			s.log.Infof("order created from cart: %s (user=%d, product=%d)",
				order.OrderNo, userID, item.ProductID)
			_ = couponID // used for tracking
		}

		// Clear cart
		if err := tx.Where("user_id = ?", userID).Delete(&CartItem{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// GetDB returns the underlying gorm.DB for direct queries.
func (s *CartService) GetDB() *gorm.DB {
	return s.db
}
