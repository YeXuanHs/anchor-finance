package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type Order struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	OrderNo    string         `gorm:"uniqueIndex;size:64;not null" json:"order_no"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	ProductID  uint           `gorm:"not null" json:"product_id"`
	Product    Product        `gorm:"foreignKey:ProductID" json:"product"`
	Quantity   int            `gorm:"default:1" json:"quantity"`
	TotalPrice float64        `gorm:"type:decimal(12,2);not null" json:"total_price"`
	Period     int            `gorm:"not null" json:"period"`
	PeriodUnit string         `gorm:"size:16;default:day" json:"period_unit"`
	Status     int            `gorm:"default:0;comment:0=pending 1=paid 2=cancelled 3=expired" json:"status"`
	Remark     string         `gorm:"size:256" json:"remark"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderService struct {
	db      *gorm.DB
	log     *logger.Logger
	invSvc  *InvoiceService
	provSvc *ProvisionService
}

func NewOrderService(db *gorm.DB, log *logger.Logger, invSvc *InvoiceService, provSvc *ProvisionService) *OrderService {
	return &OrderService{db: db, log: log, invSvc: invSvc, provSvc: provSvc}
}

type CreateOrderRequest struct {
	ProductID uint   `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"omitempty,gte=1"`
	Remark    string `json:"remark"`
}

// Create generates an order with auto order_no and calculates total.
func (s *OrderService) Create(userID uint, req CreateOrderRequest) (*Order, error) {
	var product Product
	if err := s.db.First(&product, req.ProductID).Error; err != nil {
		return nil, errors.New("product not found")
	}

	if product.Status != 1 {
		return nil, errors.New("product is disabled")
	}

	qty := req.Quantity
	if qty < 1 {
		qty = 1
	}
	if product.Stock >= 0 && product.Stock < qty {
		return nil, errors.New("insufficient stock")
	}

	order := &Order{
		OrderNo:    util.GenerateOrderNo(),
		UserID:     userID,
		ProductID:  product.ID,
		Quantity:   qty,
		TotalPrice: product.Price * float64(qty),
		Period:     product.Period,
		PeriodUnit: product.PeriodUnit,
		Status:     0,
		Remark:     req.Remark,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		if product.Stock >= 0 {
			if err := tx.Model(&product).Update("stock", gorm.Expr("stock - ?", qty)).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	s.log.Infof("order created: %s (user=%d, product=%d)", order.OrderNo, userID, product.ID)
	return order, nil
}

// GetByID fetches an order by ID.
func (s *OrderService) GetByID(id uint) (*Order, error) {
	var order Order
	if err := s.db.Preload("Product").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetUserOrders returns paginated orders for a user.
func (s *OrderService) GetUserOrders(userID, page, pageSize int) ([]Order, int64, error) {
	var orders []Order
	var total int64

	query := s.db.Model(&Order{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Preload("Product").Offset(offset).Limit(limit).
		Order("id DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// Pay marks order paid and creates invoice + user_product.
func (s *OrderService) Pay(orderID uint) (*Order, error) {
	var order Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return nil, err
	}
	if order.Status != 0 {
		return nil, errors.New("order is not pending")
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Mark paid
		if err := tx.Model(&order).Update("status", 1).Error; err != nil {
			return err
		}

		// Create invoice
		if _, err := s.invSvc.CreateWithTx(tx, order.UserID, order.ID, order.OrderNo, order.TotalPrice); err != nil {
			return err
		}

		// Create user_product
		now := time.Now()
		expire := calcExpire(now, order.Period, order.PeriodUnit)
		up := &UserProduct{
			UserID:    order.UserID,
			ProductID: order.ProductID,
			OrderID:   order.ID,
			OrderNo:   order.OrderNo,
			StartAt:   now,
			ExpireAt:  expire,
			Status:    1,
		}
		if err := tx.Create(up).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.log.Infof("order paid: %s", order.OrderNo)

	// Trigger auto-provisioning asynchronously
	if s.provSvc != nil {
		go func() {
			if err := s.provSvc.ProvisionOrder(order.ID); err != nil {
				s.log.Errorf("auto-provision for order %d failed: %v", order.ID, err)
			}
		}()
	}

	return s.GetByID(orderID)
}

// PayWithMethod pays an order with a specific payment method.
func (s *OrderService) PayWithMethod(userID, orderID uint, paymentMethod string) (*Order, error) {
	var order Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return nil, err
	}
	if order.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	if order.Status != 0 {
		return nil, errors.New("order is not pending")
	}

	// Default to balance payment
	if paymentMethod == "" {
		paymentMethod = "balance"
	}

	// Balance payment: deduct from user balance
	if paymentMethod == "balance" {
		var user struct {
			Balance float64
		}
		if err := s.db.Table("users").Where("id = ?", userID).Select("balance").First(&user).Error; err != nil {
			return nil, errors.New("user not found")
		}
		if user.Balance < order.TotalPrice {
			return nil, errors.New("insufficient balance")
		}
		// Deduct balance
		if err := s.db.Exec("UPDATE users SET balance = balance - ? WHERE id = ? AND balance >= ?",
			order.TotalPrice, userID, order.TotalPrice).Error; err != nil {
			return nil, errors.New("balance deduction failed")
		}
		// Record transaction
		s.db.Exec(`INSERT INTO balance_logs (user_id, type, amount, before_balance, after_balance, description, created_at)
			VALUES (?, 'payment', ?, ?, ?, ?, ?)`,
			userID, order.TotalPrice, user.Balance, user.Balance-order.TotalPrice,
			"Order payment: "+order.OrderNo, time.Now().Unix())
	}

	// Mark paid and create records
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"status":         1,
			"payment_method": paymentMethod,
			"paid_at":        time.Now(),
		}).Error; err != nil {
			return err
		}
		if _, err := s.invSvc.CreateWithTx(tx, order.UserID, order.ID, order.OrderNo, order.TotalPrice); err != nil {
			return err
		}
		now := time.Now()
		expire := calcExpire(now, order.Period, order.PeriodUnit)
		up := &UserProduct{
			UserID:    order.UserID,
			ProductID: order.ProductID,
			OrderID:   order.ID,
			OrderNo:   order.OrderNo,
			StartAt:   now,
			ExpireAt:  expire,
			Status:    1,
		}
		return tx.Create(up).Error
	})
	if err != nil {
		return nil, err
	}

	s.log.Infof("order paid via %s: %s", paymentMethod, order.OrderNo)
	if s.provSvc != nil {
		go func() {
			if err := s.provSvc.ProvisionOrder(order.ID); err != nil {
				s.log.Errorf("auto-provision for order %d failed: %v", order.ID, err)
			}
		}()
	}
	return s.GetByID(orderID)
}

// Cancel sets order status to cancelled.
func (s *OrderService) Cancel(orderID uint) error {
	var order Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return err
	}
	if order.Status != 0 {
		return errors.New("only pending orders can be cancelled")
	}
	return s.db.Model(&order).Update("status", 2).Error
}

// GetList returns all orders with pagination (admin).
func (s *OrderService) GetList(page, pageSize int, status *int, userID *uint) ([]Order, int64, error) {
	var orders []Order
	var total int64

	query := s.db.Model(&Order{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Preload("Product").Offset(offset).Limit(limit).
		Order("id DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

func calcExpire(start time.Time, period int, unit string) time.Time {
	switch unit {
	case "month":
		return start.AddDate(0, period, 0)
	case "year":
		return start.AddDate(period, 0, 0)
	default: // day
		return start.AddDate(0, 0, period)
	}
}

// ---------------------------------------------------------------------------
// Admin order management methods
// ---------------------------------------------------------------------------

// AdminCreateOrderItem is a single product item in admin order creation.
type AdminCreateOrderItem struct {
	ProductID    uint   `json:"product_id" binding:"required"`
	BillingCycle string `json:"billing_cycle"`
	Quantity     int    `json:"quantity"`
	PriceOverride float64 `json:"price_override"` // 管理员自定义价格，0表示使用产品默认价格
}

// AdminCreateOrderRequest is the payload for admin creating orders.
type AdminCreateOrderRequest struct {
	UserID       uint                   `json:"user_id" binding:"required"`
	Items        []AdminCreateOrderItem `json:"items" binding:"required,min=1"`
	PromoCode    string                 `json:"promo_code"`
	PaymentMethod string                `json:"payment_method"`
	Status       int16                  `json:"status"` // 0=Pending, 1=Active, etc.
	GenerateInvoice bool               `json:"generate_invoice"`
	AdminNotes   string                 `json:"admin_notes"`
}

// AdminCreateOrder allows admin to manually create an order for a user.
func (s *OrderService) AdminCreateOrder(adminID uint, req AdminCreateOrderRequest) (*Order, error) {
	// Verify user exists
	var user model.User
	if err := s.db.First(&user, req.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var totalAmount float64
	var promo *model.PromoCode

	// Validate promo code if provided
	if req.PromoCode != "" {
		var p model.PromoCode
		if err := s.db.Where("code = ? AND status = 1", req.PromoCode).First(&p).Error; err != nil {
			return nil, errors.New("invalid promo code")
		}
		now := time.Now().Unix()
		if p.StartTime > 0 && now < p.StartTime {
			return nil, errors.New("promo code not yet active")
		}
		if p.ExpirationTime > 0 && now > p.ExpirationTime {
			return nil, errors.New("promo code has expired")
		}
		if p.MaxTimes > 0 && p.UsedCount >= p.MaxTimes {
			return nil, errors.New("promo code usage limit reached")
		}
		promo = &p
	}

	// Validate products and calculate total
	type itemResult struct {
		product model.Product
		item    AdminCreateOrderItem
		price   float64
	}
	var items []itemResult

	for _, item := range req.Items {
		var product model.Product
		if err := s.db.First(&product, item.ProductID).Error; err != nil {
			return nil, fmt.Errorf("product not found: %d", item.ProductID)
		}
		if product.Status != 1 {
			return nil, fmt.Errorf("product is disabled: %s", product.Name)
		}

		qty := item.Quantity
		if qty < 1 {
			qty = 1
		}

		// Check stock
		if product.StockControl && product.Stock >= 0 && product.Stock < qty {
			return nil, fmt.Errorf("insufficient stock for %s: available %d", product.Name, product.Stock)
		}

		// Calculate price
		basePrice := product.Price
		price := basePrice
		cycle := item.BillingCycle
		if cycle == "" {
			cycle = product.BillingCycle
		}

		// Admin price override
		if item.PriceOverride > 0 {
			price = item.PriceOverride
		} else if cycle != "" && cycle != product.BillingCycle {
			switch cycle {
			case "monthly":
				price = basePrice
			case "quarterly":
				price = basePrice * 3 * 0.95
			case "semi-annually":
				price = basePrice * 6 * 0.90
			case "annually":
				price = basePrice * 12 * 0.85
			case "biennially":
				price = basePrice * 24 * 0.80
			case "triennially":
				price = basePrice * 36 * 0.75
			default:
				price = basePrice
			}
		}

		subtotal := price * float64(qty)
		totalAmount += subtotal

		items = append(items, itemResult{product: product, item: item, price: price})
	}

	// Apply promo code discount
	var discountAmount float64
	var promoID *uint
	if promo != nil {
		switch promo.Type {
		case "percent":
			discountAmount = totalAmount * promo.Value / 100
		case "fixed":
			discountAmount = promo.Value
			if discountAmount > totalAmount {
				discountAmount = totalAmount
			}
		case "override":
			discountAmount = totalAmount - promo.Value
			if discountAmount < 0 {
				discountAmount = 0
			}
		case "free":
			discountAmount = totalAmount
		}
		pid := promo.ID
		promoID = &pid
	}

	finalTotal := totalAmount - discountAmount
	if finalTotal < 0 {
		finalTotal = 0
	}

	// Create order in transaction
	var order *Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Create order record
		order = &Order{
			OrderNo:     util.GenerateOrderNo(),
			UserID:      req.UserID,
			ProductID:   items[0].product.ID, // primary product
			PromoCodeID: promoID,
			Quantity:    req.Items[0].Quantity,
			TotalPrice:  finalTotal,
			Period:      items[0].product.BillingCycle,
			Status:      int(req.Status),
			Remark:      req.AdminNotes,
		}
		if order.Status == 0 {
			order.Status = 0 // Pending
		}
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Update stock for each product
		for _, ir := range items {
			qty := ir.item.Quantity
			if qty < 1 {
				qty = 1
			}
			if ir.product.StockControl && ir.product.Stock >= 0 {
				if err := tx.Model(&ir.product).Update("stock", gorm.Expr("stock - ?", qty)).Error; err != nil {
					return err
				}
			}
		}

		// Record coupon usage
		// Record promo code usage
		if promo != nil {
			if err := tx.Model(promo).Update("used_count", gorm.Expr("used_count + 1")).Error; err != nil {
				return err
			}
			usageLog := model.PromoCodeLog{
				PromoID: promo.ID,
				UserID:  req.UserID,
				OrderID: order.ID,
				Amount:  discountAmount,
			}
			if err := tx.Create(&usageLog).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.log.Infof("admin created order: %s (user=%d, admin=%d)", order.OrderNo, req.UserID, adminID)
	return order, nil
}

// CheckOrderResult contains the result of order validation.
type CheckOrderResult struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

// CheckOrder validates whether an order can be executed.
func (s *OrderService) CheckOrder(orderID uint) (*CheckOrderResult, error) {
	var order Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return nil, errors.New("order not found")
	}

	result := &CheckOrderResult{Valid: true}

	// Check user exists
	var user model.User
	if err := s.db.First(&user, order.UserID).Error; err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, "associated user not found")
	}

	// Check product exists
	var product model.Product
	if err := s.db.First(&product, order.ProductID).Error; err != nil {
		result.Valid = false
		result.Errors = append(result.Errors, "associated product not found")
	} else {
		if product.Status != 1 {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("product '%s' is disabled", product.Name))
		}
		if product.StockControl && product.Stock >= 0 && product.Stock < order.Quantity {
			result.Valid = false
			result.Errors = append(result.Errors, fmt.Sprintf("insufficient stock for '%s': available %d, needed %d", product.Name, product.Stock, order.Quantity))
		}
	}

	// Check order status
	if order.Status != 0 {
		result.Warnings = append(result.Warnings, fmt.Sprintf("order status is %d (not pending)", order.Status))
	}

	return result, nil
}

// BatchAction represents the batch operation type.
type BatchAction string

const (
	BatchConfirm BatchAction = "confirm"
	BatchCancel  BatchAction = "cancel"
	BatchDelete  BatchAction = "delete"
)

// OrderBatchUpdateRequest is the payload for batch operations.
type OrderBatchUpdateRequest struct {
	OrderIDs []uint      `json:"order_ids" binding:"required,min=1"`
	Action   BatchAction `json:"action" binding:"required,oneof=confirm cancel delete"`
}

// BatchUpdate performs batch operations on orders.
func (s *OrderService) BatchUpdate(adminID uint, req OrderBatchUpdateRequest) (int, error) {
	var processed int

	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, orderID := range req.OrderIDs {
			var order Order
			if err := tx.First(&order, orderID).Error; err != nil {
				continue // skip invalid IDs
			}

			switch req.Action {
			case BatchConfirm:
				if order.Status == 0 {
					if err := tx.Model(&order).Update("status", 1).Error; err != nil {
						return err
					}
					processed++
				}
			case BatchCancel:
				if order.Status == 0 {
					now := time.Now()
					if err := tx.Model(&order).Updates(map[string]interface{}{
						"status":        4,
						"cancelled_at":  &now,
						"cancel_reason": "admin batch cancel",
					}).Error; err != nil {
						return err
					}
					processed++
				}
			case BatchDelete:
				if err := tx.Delete(&order).Error; err != nil {
					return err
				}
				processed++
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	s.log.Infof("admin batch %s: %d orders processed (admin=%d)", req.Action, processed, adminID)
	return processed, nil
}

// Delete soft-deletes an order.
func (s *OrderService) Delete(orderID uint) error {
	var order Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return errors.New("order not found")
	}
	return s.db.Delete(&order).Error
}

// AddNote adds an admin note to an order.
func (s *OrderService) AddNote(orderID, adminID uint, content string, isPrivate bool) (*model.OrderNote, error) {
	var order Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return nil, errors.New("order not found")
	}

	note := &model.OrderNote{
		OrderID:   orderID,
		AdminID:   adminID,
		Content:   content,
		IsPrivate: isPrivate,
	}
	if err := s.db.Create(note).Error; err != nil {
		return nil, err
	}

	// Also update the order's AdminNotes field for quick access
	s.db.Model(&order).Update("admin_notes", content)

	return note, nil
}

// GetNotes returns all notes for an order.
func (s *OrderService) GetNotes(orderID uint) ([]model.OrderNote, error) {
	var notes []model.OrderNote
	if err := s.db.Where("order_id = ?", orderID).Order("id DESC").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// GetSaleOrders returns orders related to sales with pagination.
func (s *OrderService) GetSaleOrders(page, pageSize int, saleID *uint, status *string) ([]Order, int64, error) {
	var orders []Order
	var total int64

	query := s.db.Model(&Order{})
	if status != nil && *status != "" {
		query = query.Where("status = ?", *status)
	}
	// If saleID is provided, filter by user's sale_id (requires joining users)
	if saleID != nil {
		query = query.Joins("JOIN users ON users.id = orders.user_id").Where("users.sale_id = ?", *saleID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Preload("Product").Offset(offset).Limit(limit).
		Order("id DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// ActivateOrder manually activates an order (opens service without payment).
func (s *OrderService) ActivateOrder(orderID uint) (*Order, error) {
	var order Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return nil, errors.New("order not found")
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// Mark order as paid/active
		now := time.Now()
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"status":  1,
			"paid_at": &now,
		}).Error; err != nil {
			return err
		}

		// Create user_product
		expire := calcExpire(now, order.Period, order.PeriodUnit)
		up := &UserProduct{
			UserID:    order.UserID,
			ProductID: order.ProductID,
			OrderID:   order.ID,
			OrderNo:   order.OrderNo,
			StartAt:   now,
			ExpireAt:  expire,
			Status:    1,
		}
		if err := tx.Create(up).Error; err != nil {
			return err
		}

		// Create invoice if needed
		if s.invSvc != nil {
			if _, err := s.invSvc.CreateWithTx(tx, order.UserID, order.ID, order.OrderNo, order.TotalPrice); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.log.Infof("order activated: %s", order.OrderNo)

	// Trigger auto-provisioning
	if s.provSvc != nil {
		go func() {
			if err := s.provSvc.ProvisionOrder(order.ID); err != nil {
				s.log.Errorf("auto-provision for order %d failed: %v", order.ID, err)
			}
		}()
	}

	return s.GetByID(orderID)
}

// ChangeStatus changes the status of an order.
func (s *OrderService) ChangeStatus(orderID uint, status int16) error {
	var order Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return errors.New("order not found")
	}

	updates := map[string]interface{}{
		"status": status,
	}

	// Set timestamps based on status
	now := time.Now()
	switch status {
	case 1: // Paid
		updates["paid_at"] = &now
	case 3: // Completed
		updates["completed_at"] = &now
	case 4: // Cancelled
		updates["cancelled_at"] = &now
	}

	return s.db.Model(&order).Updates(updates).Error
}

// MultiProductItem represents a product in a multi-product price calculation.
type MultiProductItem struct {
	ProductID    uint   `json:"product_id" binding:"required"`
	BillingCycle string `json:"billing_cycle"`
	Quantity     int    `json:"quantity"`
}

// MultiTotalResult is the result of multi-product price calculation.
type MultiTotalResult struct {
	Items    []MultiTotalItem `json:"items"`
	Subtotal float64          `json:"subtotal"`
	Discount float64          `json:"discount"`
	Total    float64          `json:"total"`
}

// MultiTotalItem is the price breakdown for a single product.
type MultiTotalItem struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Cycle       string  `json:"billing_cycle"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
}

// GetMultiTotal calculates the total price for multiple products.
func (s *OrderService) GetMultiTotal(items []MultiProductItem, couponCode string) (*MultiTotalResult, error) {
	result := &MultiTotalResult{}

	for _, item := range items {
		var product model.Product
		if err := s.db.First(&product, item.ProductID).Error; err != nil {
			return nil, fmt.Errorf("product not found: %d", item.ProductID)
		}
		if product.Status != 1 {
			return nil, fmt.Errorf("product is disabled: %s", product.Name)
		}

		qty := item.Quantity
		if qty < 1 {
			qty = 1
		}

		basePrice := product.Price
		price := basePrice
		cycle := item.BillingCycle
		if cycle == "" {
			cycle = product.BillingCycle
		}

		if cycle != "" && cycle != product.BillingCycle {
			switch cycle {
			case "monthly":
				price = basePrice
			case "quarterly":
				price = basePrice * 3 * 0.95
			case "semi-annually":
				price = basePrice * 6 * 0.90
			case "annually":
				price = basePrice * 12 * 0.85
			case "biennially":
				price = basePrice * 24 * 0.80
			case "triennially":
				price = basePrice * 36 * 0.75
			default:
				price = basePrice
			}
		}

		subtotal := price * float64(qty)
		result.Items = append(result.Items, MultiTotalItem{
			ProductID:   product.ID,
			ProductName: product.Name,
			Cycle:       cycle,
			Quantity:    qty,
			UnitPrice:   price,
			Subtotal:    subtotal,
		})
		result.Subtotal += subtotal
	}

	// Apply coupon
	if couponCode != "" {
		var promo model.PromoCode
		if err := s.db.Where("code = ? AND status = 1", couponCode).First(&promo).Error; err == nil {
			now := time.Now().Unix()
			valid := true
			if promo.StartTime > 0 && now < promo.StartTime {
				valid = false
			}
			if promo.ExpirationTime > 0 && now > promo.ExpirationTime {
				valid = false
			}
			if promo.MaxTimes > 0 && promo.UsedCount >= promo.MaxTimes {
				valid = false
			}

			if valid {
				switch promo.Type {
				case "percent":
					result.Discount = result.Subtotal * promo.Value / 100
				case "fixed":
					result.Discount = promo.Value
					if result.Discount > result.Subtotal {
						result.Discount = result.Subtotal
					}
				case "override":
					result.Discount = result.Subtotal - promo.Value
					if result.Discount < 0 {
						result.Discount = 0
					}
				case "free":
					result.Discount = result.Subtotal
				}
			}
		}
	}

	result.Total = result.Subtotal - result.Discount
	if result.Total < 0 {
		result.Total = 0
	}

	return result, nil
}

// ApplyCustomPromo creates a custom promo code for an order.
func (s *OrderService) ApplyCustomPromo(req CreateCustomPromoRequest) (*model.PromoCode, error) {
	// Check if code already exists
	var existing model.PromoCode
	if err := s.db.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		return nil, errors.New("promo code already exists")
	}

	now := time.Now().Unix()
	promo := &model.PromoCode{
		Code:           req.Code,
		Type:           req.Type,
		Value:          req.Value,
		Status:         1,
		StartTime:      now,
		ExpirationTime: now + 365*24*3600, // default 1 year
	}

	if err := s.db.Create(promo).Error; err != nil {
		return nil, err
	}

	s.log.Infof("custom promo created: %s", req.Code)
	return promo, nil
}

// CreateCustomPromoRequest is the payload for creating a custom promo.
type CreateCustomPromoRequest struct {
	Code        string  `json:"code" binding:"required,max=10"`
	Name        string  `json:"name" binding:"required"`
	Type        string  `json:"type" binding:"required,oneof=percent fixed override free"`
	Value       float64 `json:"value"`
	Description string  `json:"description"`
}

// ==================== P0-5: CheckProduct ====================

// CheckProduct 校验产品试用/购买资格
func (s *OrderService) CheckProduct(productID, userID uint) (map[string]interface{}, error) {
	// 查找产品
	var product Product
	if err := s.db.First(&product, productID).Error; err != nil {
		return nil, errors.New("product not found")
	}

	// 查找用户
	var user struct {
		ID           uint
		Email        string
		PhoneNumber  string
		WechatID     string
		RealName     string
		IsVerified   int
	}
	if err := s.db.Table("users").Where("id = ?", userID).
		Select("id, email, phone_number, wechat_id, real_name, is_verified").Scan(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}

	// 检查产品试用条件（payontrial_condition）
	var errors_list []string
	config := product.Config
	if config != "" {
		// 解析试用条件配置
		var productConfig map[string]interface{}
		if err := json.Unmarshal([]byte(config), &productConfig); err == nil {
			if conditions, ok := productConfig["payontrial_condition"].([]interface{}); ok {
				for _, cond := range conditions {
					switch cond {
					case "realname":
						if user.IsVerified == 0 {
							errors_list = append(errors_list, "实名认证")
						}
					case "email":
						if user.Email == "" {
							errors_list = append(errors_list, "邮箱验证")
						}
					case "phone":
						if user.PhoneNumber == "" {
							errors_list = append(errors_list, "手机验证")
						}
					case "wechat":
						if user.WechatID == "" {
							errors_list = append(errors_list, "微信验证")
						}
					}
				}
			}
		}
	}

	// 检查库存
	if product.StockControl && product.Stock >= 0 && product.Stock < 1 {
		errors_list = append(errors_list, "库存不足")
	}

	// 检查产品状态
	if product.Status != 1 {
		errors_list = append(errors_list, "产品已下架")
	}

	result := map[string]interface{}{
		"valid":       len(errors_list) == 0,
		"errors":      errors_list,
		"product_id":  productID,
		"user_id":     userID,
	}
	return result, nil
}

// ==================== P3-17: GetOrderConfig ====================

// GetOrderConfig 获取下单配置（产品定价/配置选项）
func (s *OrderService) GetOrderConfig(productID uint, cycle string) (map[string]interface{}, error) {
	var product Product
	if err := s.db.First(&product, productID).Error; err != nil {
		return nil, errors.New("product not found")
	}

	basePrice := product.Price

	// 计算周期价格
	price := basePrice
	switch cycle {
	case "quarterly":
		price = basePrice * 3 * 0.95
	case "semi-annually":
		price = basePrice * 6 * 0.90
	case "annually":
		price = basePrice * 12 * 0.85
	case "biennially":
		price = basePrice * 24 * 0.80
	case "triennially":
		price = basePrice * 36 * 0.75
	}

	// 获取配置选项
	var options []map[string]interface{}
	s.db.Table("config_options").Where("product_id = ?", productID).Find(&options)

	// 获取自定义字段
	var customFields []map[string]interface{}
	s.db.Table("custom_fields").Where("`group` = 'product' AND relid = ?", productID).Find(&customFields)

	return map[string]interface{}{
		"product":       product,
		"price":         price,
		"billing_cycle": cycle,
		"options":       options,
		"custom_fields": customFields,
	}, nil
}

// SearchPageConfig returns filter configuration for the order search page.
func (s *OrderService) SearchPageConfig() map[string]interface{} {
	// Order statuses
	statuses := map[int]string{
		0: "Pending",
		1: "Paid",
		2: "Processing",
		3: "Completed",
		4: "Cancelled",
		5: "Refunded",
		6: "Partially Refunded",
		7: "Fraud",
		8: "Dispute",
	}

	// Payment methods (from payment_gateways table)
	var gateways []struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	}
	s.db.Table("payment_gateways").Where("status = 1").Select("name, title").Find(&gateways)

	return map[string]interface{}{
		"statuses":        statuses,
		"payment_gateways": gateways,
	}
}
