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

	amount := product.Price
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

// MultiRenewTask represents a batch renewal task.
type MultiRenewTask struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	Name       string     `gorm:"type:varchar(128)" json:"name"`
	UserID     uint       `gorm:"index" json:"user_id"`
	Status     int16      `gorm:"default:0" json:"status"` // 0=pending 1=running 2=completed 3=failed
	Period     int        `json:"period"`
	PeriodUnit string     `gorm:"type:varchar(16)" json:"period_unit"`
	AutoPay    bool       `json:"auto_pay"`
	Note       string     `gorm:"type:text" json:"note"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// MultiRenewTaskService represents a service in a batch renewal task.
type MultiRenewTaskService struct {
	ID        uint `gorm:"primaryKey" json:"id"`
	TaskID    uint `gorm:"index" json:"task_id"`
	ServiceID uint `json:"service_id"`
}

// MultiRenewLog represents a log entry for batch renewal.
type MultiRenewLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TaskID    uint      `gorm:"index" json:"task_id"`
	ServiceID uint      `json:"service_id"`
	Status    string    `gorm:"type:varchar(32)" json:"status"`
	Message   string    `gorm:"type:text" json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// List returns paginated batch renewal tasks.
func (s *MultiRenewService) List(page, pageSize int, status *int16) ([]MultiRenewTask, int64, error) {
	var items []MultiRenewTask
	var total int64
	query := s.db.Model(&MultiRenewTask{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)
	return items, total, nil
}

// GetByID returns a batch renewal task by ID.
func (s *MultiRenewService) GetByID(id uint) (*MultiRenewTask, error) {
	var task MultiRenewTask
	if err := s.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// Create creates a new batch renewal task.
func (s *MultiRenewService) Create(name string, serviceIDs []uint, period int, periodUnit string, autoPay bool, note string) (*MultiRenewTask, error) {
	task := &MultiRenewTask{
		Name:       name,
		Period:     period,
		PeriodUnit: periodUnit,
		AutoPay:    autoPay,
		Note:       note,
		Status:     0,
	}
	if err := s.db.Create(task).Error; err != nil {
		return nil, err
	}

	// Add services to task
	for _, serviceID := range serviceIDs {
		ts := &MultiRenewTaskService{
			TaskID:    task.ID,
			ServiceID: serviceID,
		}
		s.db.Create(ts)
	}

	return task, nil
}

// Execute executes a batch renewal task.
func (s *MultiRenewService) Execute(taskID uint) (map[string]interface{}, error) {
	task, err := s.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	// Get services for this task
	var taskServices []MultiRenewTaskService
	s.db.Where("task_id = ?", taskID).Find(&taskServices)

	serviceIDs := make([]uint, len(taskServices))
	for i, ts := range taskServices {
		serviceIDs[i] = ts.ServiceID
	}

	// Use BatchRenew to execute
	summary, err := s.BatchRenew(task.UserID, serviceIDs, task.PeriodUnit)
	if err != nil {
		return nil, err
	}

	// Update task status
	s.db.Model(task).Update("status", 2) // completed

	return map[string]interface{}{
		"task_id": taskID,
		"summary": summary,
	}, nil
}

// Cancel cancels a batch renewal task.
func (s *MultiRenewService) Cancel(taskID uint) error {
	return s.db.Model(&MultiRenewTask{}).Where("id = ?", taskID).Update("status", 3).Error
}

// Delete deletes a batch renewal task.
func (s *MultiRenewService) Delete(taskID uint) error {
	s.db.Where("task_id = ?", taskID).Delete(&MultiRenewTaskService{})
	s.db.Where("task_id = ?", taskID).Delete(&MultiRenewLog{})
	return s.db.Delete(&MultiRenewTask{}, taskID).Error
}

// GetLogs returns logs for a batch renewal task.
func (s *MultiRenewService) GetLogs(taskID uint, page, pageSize int) ([]MultiRenewLog, int64, error) {
	var logs []MultiRenewLog
	var total int64
	query := s.db.Model(&MultiRenewLog{}).Where("task_id = ?", taskID)
	query.Count(&total)
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)
	return logs, total, nil
}

// GetStats returns batch renewal statistics.
func (s *MultiRenewService) GetStats() (map[string]interface{}, error) {
	var totalTasks int64
	s.db.Model(&MultiRenewTask{}).Count(&totalTasks)

	var completedTasks int64
	s.db.Model(&MultiRenewTask{}).Where("status = ?", 2).Count(&completedTasks)

	var pendingTasks int64
	s.db.Model(&MultiRenewTask{}).Where("status = ?", 0).Count(&pendingTasks)

	return map[string]interface{}{
		"total_tasks":     totalTasks,
		"completed_tasks": completedTasks,
		"pending_tasks":   pendingTasks,
	}, nil
}

// Preview previews the services that would be renewed.
func (s *MultiRenewService) Preview(serviceIDs []uint, period int, periodUnit string) (map[string]interface{}, error) {
	var items []map[string]interface{}
	var totalAmount float64

	for _, serviceID := range serviceIDs {
		var svc model.ClientService
		if err := s.db.First(&svc, serviceID).Error; err != nil {
			continue
		}

		var product model.Product
		s.db.First(&product, svc.ProductID)

		item := map[string]interface{}{
			"service_id": svc.ID,
			"name":       svc.Name,
			"product":    product.Name,
			"amount":     product.Price,
		}
		items = append(items, item)
		totalAmount += product.Price
	}

	return map[string]interface{}{
		"services":     items,
		"total_amount": totalAmount,
		"period":       period,
		"period_unit":  periodUnit,
	}, nil
}
