package service

import (
	"errors"
	"time"

	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// UpgradeOrder 升降级订单
type UpgradeOrder struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	OrderNo         string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"order_no"`
	UserID          uint           `gorm:"index;not null" json:"user_id"`
	UserProductID   uint           `gorm:"index;not null" json:"user_product_id"`
	TargetProductID uint           `gorm:"index;not null" json:"target_product_id"`
	Type            int8           `gorm:"type:smallint;not null" json:"type"` // 1升级 2降级
	BillingCycle    string         `gorm:"type:varchar(20)" json:"billing_cycle"`
	Amount          float64        `gorm:"type:decimal(10,2);not null" json:"amount"` // 差价
	Status          int8           `gorm:"type:smallint;default:1" json:"status"`     // 1待支付 2已支付 3已开通 4已取消
	PaymentMethod   string         `gorm:"type:varchar(50)" json:"payment_method"`
	PaidAt          *time.Time     `json:"paid_at"`
	CompletedAt     *time.Time     `json:"completed_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// CreateUpgradeRequest 创建升降级请求
type CreateUpgradeRequest struct {
	UserProductID   uint   `json:"user_product_id" binding:"required"`
	TargetProductID uint   `json:"target_product_id" binding:"required"`
	BillingCycle    string `json:"billing_cycle"`
}

type UpgradeService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewUpgradeService(db *gorm.DB, log *logger.Logger) *UpgradeService {
	return &UpgradeService{db: db, log: log}
}

// GetAvailableUpgrades 获取可用的升降级选项
func (s *UpgradeService) GetAvailableUpgrades(userID, userProductID uint) ([]Product, error) {
	var up UserProduct
	if err := s.db.First(&up, userProductID).Error; err != nil {
		return nil, errors.New("user product not found")
	}
	if up.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	var currentProduct Product
	if err := s.db.First(&currentProduct, up.ProductID).Error; err != nil {
		return nil, errors.New("current product not found")
	}

	// 查找同分组的其他产品作为升降级选项
	var products []Product
	err := s.db.Where("id != ? AND group_id = ? AND status = 1", currentProduct.ID, currentProduct.GroupID).
		Order("price ASC").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

// CreateUpgrade 计算差价并创建升降级订单
func (s *UpgradeService) CreateUpgrade(userID uint, req CreateUpgradeRequest) (*UpgradeOrder, error) {
	var up UserProduct
	if err := s.db.First(&up, req.UserProductID).Error; err != nil {
		return nil, errors.New("user product not found")
	}
	if up.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	if up.Status != "Active" {
		return nil, errors.New("user product is not active")
	}

	var currentProduct Product
	if err := s.db.First(&currentProduct, up.ProductID).Error; err != nil {
		return nil, errors.New("current product not found")
	}

	var targetProduct Product
	if err := s.db.First(&targetProduct, req.TargetProductID).Error; err != nil {
		return nil, errors.New("target product not found")
	}
	if targetProduct.Status != 1 {
		return nil, errors.New("target product is disabled")
	}

	// 判断升降级类型
	var upgradeType int8
	if targetProduct.Price > currentProduct.Price {
		upgradeType = 1 // 升级
	} else if targetProduct.Price < currentProduct.Price {
		upgradeType = 2 // 降级
	} else {
		return nil, errors.New("target product has the same price, no upgrade needed")
	}

	// 计算差价（按剩余天数比例）
	var expireAt time.Time
	if up.NextDueDate != nil {
		expireAt = *up.NextDueDate
	} else {
		expireAt = time.Now().AddDate(1, 0, 0) // 默认1年后
	}
	amount := s.calcPriceDiff(currentProduct.Price, targetProduct.Price, expireAt)

	// 降级差价应为负数时按0处理
	if amount < 0 {
		amount = 0
	}

	order := &UpgradeOrder{
		OrderNo:         util.GenerateOrderNo(),
		UserID:          userID,
		UserProductID:   req.UserProductID,
		TargetProductID: req.TargetProductID,
		Type:            upgradeType,
		BillingCycle:    req.BillingCycle,
		Amount:          amount,
		Status:          1,
	}

	if err := s.db.Create(order).Error; err != nil {
		return nil, err
	}

	s.log.Infof("upgrade order created: %s (user=%d, type=%d, amount=%.2f)",
		order.OrderNo, userID, upgradeType, amount)
	return order, nil
}

// GetByID 根据ID获取升降级订单
func (s *UpgradeService) GetByID(id uint) (*UpgradeOrder, error) {
	var order UpgradeOrder
	if err := s.db.First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

// GetUserUpgrades 获取用户的升降级订单（分页）
func (s *UpgradeService) GetUserUpgrades(userID uint, page, pageSize int) ([]UpgradeOrder, int64, error) {
	var orders []UpgradeOrder
	var total int64

	query := s.db.Model(&UpgradeOrder{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// PayAndApply 支付并应用升降级
func (s *UpgradeService) PayAndApply(orderID uint) (*UpgradeOrder, error) {
	var order UpgradeOrder
	if err := s.db.First(&order, orderID).Error; err != nil {
		return nil, err
	}
	if order.Status != 1 {
		return nil, errors.New("order is not pending payment")
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		// 标记已支付
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"status": 2,
			"paid_at": now,
		}).Error; err != nil {
			return err
		}

		// 更新用户产品
		var targetProduct Product
		if err := tx.First(&targetProduct, order.TargetProductID).Error; err != nil {
			return err
		}

		updates := map[string]interface{}{
			"product_id": order.TargetProductID,
			"name":       targetProduct.Name,
		}
		if order.BillingCycle != "" {
			updates["billing_cycle"] = order.BillingCycle
		}

		if err := tx.Model(&UserProduct{}).Where("id = ?", order.UserProductID).
			Updates(updates).Error; err != nil {
			return err
		}

		// 标记已开通
		completedAt := time.Now()
		if err := tx.Model(&order).Updates(map[string]interface{}{
			"status":       3,
			"completed_at": completedAt,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.log.Infof("upgrade applied: %s", order.OrderNo)
	return s.GetByID(orderID)
}

// calcPriceDiff 按剩余天数比例计算差价
func (s *UpgradeService) calcPriceDiff(currentPrice, targetPrice float64, expireAt time.Time) float64 {
	diff := targetPrice - currentPrice
	if diff <= 0 {
		return 0
	}

	now := time.Now()
	if expireAt.After(now) {
		remaining := expireAt.Sub(now).Hours() / 24
		totalDays := 30.0 // 默认按月计算
		if remaining > 0 && remaining < totalDays {
			ratio := remaining / totalDays
			return diff * ratio
		}
	}

	return diff
}

// Cancel 取消升降级订单
func (s *UpgradeService) Cancel(orderID uint) error {
	var order UpgradeOrder
	if err := s.db.First(&order, orderID).Error; err != nil {
		return err
	}
	if order.Status != 1 {
		return errors.New("only pending orders can be cancelled")
	}
	return s.db.Model(&order).Update("status", 4).Error
}

// GetList 获取所有升降级订单 (admin, 分页)
func (s *UpgradeService) GetList(page, pageSize int, status *int8, userID *uint) ([]UpgradeOrder, int64, error) {
	var orders []UpgradeOrder
	var total int64

	query := s.db.Model(&UpgradeOrder{})
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
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}
