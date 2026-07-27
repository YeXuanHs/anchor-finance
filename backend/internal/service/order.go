package service

import (
	"errors"
	"time"

	"github.com/anchor-finance/backend/internal/util"
	"github.com/anchor-finance/backend/pkg/logger"

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
	db    *gorm.DB
	log   *logger.Logger
	invSvc *InvoiceService
}

func NewOrderService(db *gorm.DB, log *logger.Logger, invSvc *InvoiceService) *OrderService {
	return &OrderService{db: db, log: log, invSvc: invSvc}
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
