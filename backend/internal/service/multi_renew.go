package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type MultiRenewService struct {
	db        *gorm.DB
	log       *logger.Logger
	invSvc    *InvoiceService
	balSvc    *BalanceService
}

func NewMultiRenewService(db *gorm.DB, log *logger.Logger, invSvc *InvoiceService, balSvc *BalanceService) *MultiRenewService {
	return &MultiRenewService{db: db, log: log, invSvc: invSvc, balSvc: balSvc}
}

type RenewResult struct {
	ServiceID    uint    `json:"service_id"`
	OrderID      uint    `json:"order_id"`
	InvoiceID    uint    `json:"invoice_id"`
	Amount       float64 `json:"amount"`
	NewDueDate   *time.Time `json:"new_due_date"`
	PaymentOK    bool    `json:"payment_ok"`
	Error        string  `json:"error,omitempty"`
}

type BatchRenewSummary struct {
	TotalAmount float64       `json:"total_amount"`
	Renewed     int           `json:"renewed"`
	Failed      int           `json:"failed"`
	Results     []RenewResult `json:"results"`
}

func (s *MultiRenewService) BatchRenew(userID uint, serviceIDs []uint, cycle string) (*BatchRenewSummary, error) {
	summary := &BatchRenewSummary{
		Results: make([]RenewResult, 0, len(serviceIDs)),
	}

	for _, serviceID := range serviceIDs {
		result := s.renewSingleService(userID, serviceID, cycle)
		summary.Results = append(summary.Results, result)

		if result.Error != "" {
			summary.Failed++
		} else {
			summary.Renewed++
			summary.TotalAmount += result.Amount
		}
	}

	s.log.WithFields(map[string]interface{}{
		"user_id":    userID,
		"cycle":      cycle,
		"renewed":    summary.Renewed,
		"failed":     summary.Failed,
		"total_amt":  summary.TotalAmount,
	}).Info("batch renew completed")

	return summary, nil
}

// renewSingleService handles the renewal of a single service instance.
func (s *MultiRenewService) renewSingleService(userID, serviceID uint, cycle string) RenewResult {
	result := RenewResult{ServiceID: serviceID}

	// 1. Verify ownership and load the client service
	var svc model.ClientService
	if err := s.db.Where("id = ? AND user_id = ?", serviceID, userID).First(&svc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			result.Error = "service not found or access denied"
		} else {
			result.Error = fmt.Sprintf("load service: %v", err)
		}
		return result
	}

	// 2. Load product to get pricing
	var product model.Product
	if err := s.db.First(&product, svc.ProductID).Error; err != nil {
		result.Error = fmt.Sprintf("load product: %v", err)
		return result
	}

	amount := product.Price.InexactFloat64()
	if amount <= 0 {
		result.Error = "product price is zero"
		return result
	}
	result.Amount = amount

	// 3. Create renewal order
	order := &model.Order{
		OrderNo:      util.GenerateOrderNo(),
		UserID:       userID,
		ProductID:    svc.ProductID,
		UserProductID: svc.ID,
		Type:         "renew",
		Amount:       product.Price,
		Total:        product.Price,
		Currency:     product.Currency,
		BillingCycle: cycle,
		Quantity:     1,
		Description:  fmt.Sprintf("Renewal for service %s (%s)", svc.Name, cycle),
		Status:       0,
		PaymentStatus: 0,
	}
	if err := s.db.Create(order).Error; err != nil {
		result.Error = fmt.Sprintf("create order: %v", err)
		return result
	}
	result.OrderID = order.ID

	// 4. Create invoice for the renewal
	invoice, err := s.invSvc.Create(userID, order.ID, order.OrderNo, amount)
	if err != nil {
		result.Error = fmt.Sprintf("create invoice: %v", err)
		return result
	}
	result.InvoiceID = invoice.ID

	// 5. Try immediate balance payment
	if s.balSvc != nil {
		if err := s.balSvc.Deduct(userID, amount, fmt.Sprintf("Renewal order %s", order.OrderNo)); err != nil {
			s.log.Warnf("balance payment failed for service %d: %v", serviceID, err)
			result.PaymentOK = false
			result.Error = fmt.Sprintf("balance payment failed: %v", err)
			return result
		}

		// Mark order and invoice as paid
		now := time.Now()
		s.db.Model(order).Updates(map[string]interface{}{
			"status":         1,
			"payment_status": 1,
			"payment_method": "balance",
			"paid_at":        &now,
		})
		s.db.Model(invoice).Updates(map[string]interface{}{
			"status":   1,
			"paid_at":  &now,
		})

		result.PaymentOK = true
	}

	// 6. Extend the service's expiration date
	var newDueDate time.Time
	if svc.ExpiredAt != nil && svc.ExpiredAt.After(time.Now()) {
		newDueDate = s.extendDate(*svc.ExpiredAt, cycle)
	} else {
		newDueDate = s.extendDate(time.Now(), cycle)
	}

	s.db.Model(&svc).Updates(map[string]interface{}{
		"expired_at": &newDueDate,
		"status":     model.ClientServiceActive,
	})
	result.NewDueDate = &newDueDate

	s.log.Infof("service renewed: service=%d user=%d cycle=%s new_due=%s amount=%.2f",
		serviceID, userID, cycle, newDueDate.Format("2006-01-02"), amount)

	return result
}

// extendDate adds the billing cycle duration to a date.
func (s *MultiRenewService) extendDate(from time.Time, cycle string) time.Time {
	switch cycle {
	case "monthly":
		return from.AddDate(0, 1, 0)
	case "quarterly":
		return from.AddDate(0, 3, 0)
	case "semi-annually":
		return from.AddDate(0, 6, 0)
	case "annually":
		return from.AddDate(1, 0, 0)
	case "biennially":
		return from.AddDate(2, 0, 0)
	case "triennially":
		return from.AddDate(3, 0, 0)
	default:
		return from.AddDate(0, 1, 0)
	}
}
