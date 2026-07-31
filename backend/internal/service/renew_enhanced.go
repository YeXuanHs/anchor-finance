package service

import (
	"fmt"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// RenewEnhancedService 续费增强服务
type RenewEnhancedService struct {
	db *gorm.DB
}

// NewRenewEnhancedService 创建续费增强服务
func NewRenewEnhancedService(db *gorm.DB) *RenewEnhancedService {
	return &RenewEnhancedService{db: db}
}

// RenewInvoice 续费发票
type RenewInvoice struct {
	ID        uint      `gorm:"primaryKey"`
	HostID    uint      `gorm:"index"`
	UserID    uint      `gorm:"index"`
	InvoiceID uint      `gorm:"index"`
	Cycle     string    `gorm:"size:32"`
	Amount    float64
	Status    string    `gorm:"size:32"` // pending/paid/cancelled
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RenewParams 续费参数
type RenewParams struct {
	HostID      uint    `json:"host_id"`
	Cycle       string  `json:"cycle"`
	PaymentMethod string `json:"payment_method"`
	UseBalance  bool    `json:"use_balance"`
	Amount      float64 `json:"amount"`
}

// SetOtherParams 设置其他参数
func (s *RenewEnhancedService) SetOtherParams(hostID uint, params map[string]interface{}) error {
	var host struct {
		Config datatypes.JSON
	}
	if err := s.db.Table("hosts").Where("id = ?", hostID).Select("config").First(&host).Error; err != nil {
		return err
	}

	config := make(map[string]interface{})
	if host.Config != nil {
		host.Config.Unmarshal(&config)
	}

	for k, v := range params {
		config[k] = v
	}

	configJSON, _ := datatypes.MarshalJSON(config)
	return s.db.Table("hosts").Where("id = ?", hostID).Update("config", configJSON).Error
}

// GetOtherParams 获取其他参数
func (s *RenewEnhancedService) GetOtherParams(hostID uint) (map[string]interface{}, error) {
	var host struct {
		Config datatypes.JSON
	}
	if err := s.db.Table("hosts").Where("id = ?", hostID).Select("config").First(&host).Error; err != nil {
		return nil, err
	}

	config := make(map[string]interface{})
	if host.Config != nil {
		host.Config.Unmarshal(&config)
	}

	return config, nil
}

// SetPayType 设置支付方式
func (s *RenewEnhancedService) SetPayType(hostID uint, payType string) error {
	return s.SetOtherParams(hostID, map[string]interface{}{
		"renew_pay_type": payType,
	})
}

// GetPayType 获取支付方式
func (s *RenewEnhancedService) GetPayType(hostID uint) (string, error) {
	params, err := s.GetOtherParams(hostID)
	if err != nil {
		return "", err
	}

	if payType, ok := params["renew_pay_type"].(string); ok {
		return payType, nil
	}

	return "balance", nil
}

// DeleteRenewInvoice 删除续费发票
func (s *RenewEnhancedService) DeleteRenewInvoice(invoiceID uint) error {
	return s.db.Where("id = ? AND status = ?", invoiceID, "pending").Delete(&RenewInvoice{}).Error
}

// DeleteHostUnpaidUpgradeInvoice 删除主机未支付升级发票
func (s *RenewEnhancedService) DeleteHostUnpaidUpgradeInvoice(hostID uint) error {
	return s.db.Where("host_id = ? AND status = ? AND type = ?", hostID, "pending", "upgrade").
		Delete(&struct{ ID uint }{}).Error
}

// BatchRenew 批量续费（增强版）
func (s *RenewEnhancedService) BatchRenew(params []RenewParams, userID uint) ([]map[string]interface{}, error) {
	results := make([]map[string]interface{}, len(params))

	for i, p := range params {
		result := map[string]interface{}{
			"host_id": p.HostID,
			"status":  "pending",
		}

		// 计算续费价格
		price, err := s.CalculatedPrice(p.HostID, p.Cycle)
		if err != nil {
			result["status"] = "failed"
			result["error"] = err.Error()
			results[i] = result
			continue
		}

		result["price"] = price
		result["cycle"] = p.Cycle

		// 创建续费发票
		invoice := &RenewInvoice{
			HostID: p.HostID,
			UserID: userID,
			Cycle:  p.Cycle,
			Amount: price,
			Status: "pending",
		}
		s.db.Create(invoice)
		result["invoice_id"] = invoice.ID

		// 如果使用余额支付
		if p.UseBalance {
			err = s.payWithBalance(userID, price)
			if err != nil {
				result["status"] = "pending_payment"
				result["error"] = "Insufficient balance"
			} else {
				invoice.Status = "paid"
				s.db.Save(invoice)
				result["status"] = "paid"

				// 延长服务期限
				s.extendHost(p.HostID, p.Cycle)
			}
		}

		results[i] = result
	}

	return results, nil
}

// Renew 单个续费
func (s *RenewEnhancedService) Renew(params RenewParams, userID uint) (map[string]interface{}, error) {
	results, err := s.BatchRenew([]RenewParams{params}, userID)
	if err != nil {
		return nil, err
	}
	return results[0], nil
}

// CalculatedPrice 计算续费价格
func (s *RenewEnhancedService) CalculatedPrice(hostID uint, cycle string) (float64, error) {
	var host struct {
		ProductID uint
		Price     float64
		Cycle     string
	}
	if err := s.db.Table("hosts").Where("id = ?", hostID).Select("product_id, price, cycle").First(&host).Error; err != nil {
		return 0, err
	}

	// 查找产品定价
	var pricing struct {
		Price float64
	}
	err := s.db.Table("pricings").
		Where("product_id = ? AND cycle = ? AND enabled = ?", host.ProductID, cycle, true).
		Select("price").First(&pricing).Error

	if err != nil {
		// 使用默认价格
		return host.Price, nil
	}

	return pricing.Price, nil
}

// RenewHandle 续费处理
func (s *RenewEnhancedService) RenewHandle(hostID uint, cycle string, userID uint) (map[string]interface{}, error) {
	price, err := s.CalculatedPrice(hostID, cycle)
	if err != nil {
		return nil, err
	}

	// 创建续费订单
	order := map[string]interface{}{
		"user_id": userID,
		"type":    "renewal",
		"host_id": hostID,
		"cycle":   cycle,
		"amount":  price,
		"status":  "pending",
	}

	logger.Info("Creating renewal order", "host_id", hostID, "cycle", cycle, "price", price)

	return order, nil
}

// UnchangePrice 不变价格续费
func (s *RenewEnhancedService) UnchangePrice(hostID uint, cycle string) (float64, error) {
	var host struct {
		Price float64
	}
	if err := s.db.Table("hosts").Where("id = ?", hostID).Select("price").First(&host).Error; err != nil {
		return 0, err
	}

	return host.Price, nil
}

// GetRenewalPage 获取续费页面数据
func (s *RenewEnhancedService) GetRenewalPage(hostID uint) (map[string]interface{}, error) {
	var host struct {
		ID        uint
		ProductID uint
		Hostname  string
		Price     float64
		Cycle     string
		ExpiredAt time.Time
		Status    int
	}
	if err := s.db.Table("hosts").Where("id = ?", hostID).First(&host).Error; err != nil {
		return nil, err
	}

	// 获取可用续费周期
	var pricings []struct {
		Cycle string
		Price float64
	}
	s.db.Table("pricings").
		Where("product_id = ? AND enabled = ?", host.ProductID, true).
		Order("sort_order").
		Find(&pricings)

	// 获取支付方式
	var payTypes []map[string]interface{}
	s.db.Table("payment_gateways").Where("enabled = ?", true).Find(&payTypes)

	result := map[string]interface{}{
		"host": map[string]interface{}{
			"id":         host.ID,
			"hostname":   host.Hostname,
			"price":      host.Price,
			"cycle":      host.Cycle,
			"expired_at": host.ExpiredAt,
			"status":     host.Status,
		},
		"cycles":     pricings,
		"pay_types":  payTypes,
		"auto_renew": false,
	}

	return result, nil
}

// SetAutoRenew 设置自动续费
func (s *RenewEnhancedService) SetAutoRenew(hostID uint, enabled bool) error {
	return s.SetOtherParams(hostID, map[string]interface{}{
		"auto_renew": enabled,
	})
}

// payWithBalance 使用余额支付
func (s *RenewEnhancedService) payWithBalance(userID uint, amount float64) error {
	var user struct {
		Balance float64
	}
	if err := s.db.Table("users").Where("id = ?", userID).Select("balance").First(&user).Error; err != nil {
		return err
	}

	if user.Balance < amount {
		return fmt.Errorf("insufficient balance: %.2f < %.2f", user.Balance, amount)
	}

	// 扣除余额
	return s.db.Table("users").Where("id = ?", userID).
		Update("balance", gorm.Expr("balance - ?", amount)).Error
}

// extendHost 延长主机服务期限
func (s *RenewEnhancedService) extendHost(hostID uint, cycle string) {
	var host struct {
		ExpiredAt time.Time
	}
	s.db.Table("hosts").Where("id = ?", hostID).Select("expired_at").First(&host)

	// 计算新到期时间
	newExpired := host.ExpiredAt
	if newExpired.Before(time.Now()) {
		newExpired = time.Now()
	}

	switch cycle {
	case "monthly":
		newExpired = newExpired.AddDate(0, 1, 0)
	case "quarterly":
		newExpired = newExpired.AddDate(0, 3, 0)
	case "semi_annual":
		newExpired = newExpired.AddDate(0, 6, 0)
	case "yearly":
		newExpired = newExpired.AddDate(1, 0, 0)
	case "biennial":
		newExpired = newExpired.AddDate(2, 0, 0)
	case "triennial":
		newExpired = newExpired.AddDate(3, 0, 0)
	}

	s.db.Table("hosts").Where("id = ?", hostID).Update("expired_at", newExpired)
}
