package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type Invoice struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	InvoiceNo       string         `gorm:"uniqueIndex;size:64;not null" json:"invoice_no"`
	UserID          uint           `gorm:"index;not null" json:"user_id"`
	OrderID         uint           `json:"order_id"`
	OrderNo         string         `gorm:"size:64" json:"order_no"`
	Amount          float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Type            string         `gorm:"size:32;default:new;comment:new/renew/upgrade/combined" json:"type"`
	Status          int            `gorm:"default:0;comment:0=pending 1=paid 2=cancelled 3=overdue 4=refunded" json:"status"`
	PaidAt          *time.Time     `json:"paid_at"`
	DueDate         time.Time      `gorm:"index" json:"due_date"`
	Payment         string         `gorm:"size:64" json:"payment"`
	PaymentNotes    string         `gorm:"size:256" json:"payment_notes"`
	Credit          float64        `gorm:"type:decimal(12,2);default:0" json:"credit"`
	UseCreditLimit  int            `gorm:"default:0" json:"use_credit_limit"`
	Notes           string         `gorm:"type:text" json:"notes"`
	Remark          string         `gorm:"size:256" json:"remark"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type InvoiceService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewInvoiceService(db *gorm.DB, log *logger.Logger) *InvoiceService {
	return &InvoiceService{db: db, log: log}
}

// Create generates an invoice for a new order.
func (s *InvoiceService) Create(userID, orderID uint, orderNo string, amount float64) (*Invoice, error) {
	return s.CreateWithTx(s.db, userID, orderID, orderNo, amount)
}

// CreateWithTx creates an invoice within an existing transaction.
func (s *InvoiceService) CreateWithTx(tx *gorm.DB, userID, orderID uint, orderNo string, amount float64) (*Invoice, error) {
	invoice := &Invoice{
		InvoiceNo: util.GenerateInvoiceNo(),
		UserID:    userID,
		OrderID:   orderID,
		OrderNo:   orderNo,
		Amount:    amount,
		Type:      "new",
		Status:    0,
		DueDate:   time.Now().AddDate(0, 0, 7), // 7-day due
	}
	if err := tx.Create(invoice).Error; err != nil {
		return nil, err
	}
	return invoice, nil
}

// CreateRenew generates a renewal invoice.
func (s *InvoiceService) CreateRenew(userID uint, amount float64, remark string) (*Invoice, error) {
	invoice := &Invoice{
		InvoiceNo: util.GenerateInvoiceNo(),
		UserID:    userID,
		Amount:    amount,
		Type:      "renew",
		Status:    0,
		DueDate:   time.Now().AddDate(0, 0, 7),
		Remark:    remark,
	}
	if err := s.db.Create(invoice).Error; err != nil {
		return nil, err
	}
	s.log.Infof("renewal invoice created: %s (user=%d, amount=%.2f)", invoice.InvoiceNo, userID, amount)
	return invoice, nil
}

// GetByID fetches an invoice by ID.
func (s *InvoiceService) GetByID(id uint) (*Invoice, error) {
	var inv Invoice
	if err := s.db.First(&inv, id).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

// GetUserInvoices returns paginated invoices for a user.
func (s *InvoiceService) GetUserInvoices(userID uint, page, pageSize int) ([]Invoice, int64, error) {
	var invoices []Invoice
	var total int64

	query := s.db.Model(&Invoice{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	return invoices, total, nil
}

// PayWithBalance pays an invoice using the user's balance.
func (s *InvoiceService) PayWithBalance(invoiceID uint) error {
	inv, err := s.GetByID(invoiceID)
	if err != nil {
		return err
	}
	if inv.Status != 0 {
		return errors.New("invoice is not pending")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var user User
		if err := tx.First(&user, inv.UserID).Error; err != nil {
			return err
		}
		if user.Balance < inv.Amount {
			return errors.New("insufficient balance")
		}

		// Deduct balance
		if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", inv.Amount)).Error; err != nil {
			return err
		}

		// Mark paid
		now := time.Now()
		if err := tx.Model(inv).Updates(map[string]interface{}{
			"status":  1,
			"paid_at": &now,
		}).Error; err != nil {
			return err
		}

		s.log.Infof("invoice paid with balance: %s (user=%d)", inv.InvoiceNo, inv.UserID)
		return nil
	})
}

// MarkPaid marks an invoice as paid (external payment callback).
func (s *InvoiceService) MarkPaid(invoiceID uint) error {
	inv, err := s.GetByID(invoiceID)
	if err != nil {
		return err
	}
	if inv.Status != 0 {
		return errors.New("invoice is not pending")
	}

	now := time.Now()
	return s.db.Model(inv).Updates(map[string]interface{}{
		"status":  1,
		"paid_at": &now,
	}).Error
}

// Cancel marks an invoice as cancelled.
func (s *InvoiceService) Cancel(invoiceID uint) error {
	inv, err := s.GetByID(invoiceID)
	if err != nil {
		return err
	}
	if inv.Status != 0 {
		return errors.New("only pending invoices can be cancelled")
	}
	return s.db.Model(inv).Update("status", 2).Error
}

// GetOverdueInvoices finds all pending invoices past due date.
func (s *InvoiceService) GetOverdueInvoices() ([]Invoice, error) {
	var invoices []Invoice
	if err := s.db.Where("status = 0 AND due_date < ?", time.Now()).
		Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

// GetList returns all invoices with pagination (admin).
func (s *InvoiceService) GetList(page, pageSize int, status *int, userID *uint) ([]Invoice, int64, error) {
	var invoices []Invoice
	var total int64

	query := s.db.Model(&Invoice{})
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
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	return invoices, total, nil
}

// GetListEnhanced 发票列表（支持复杂筛选：金额区间/时间范围/行项描述/付款方式细分）
func (s *InvoiceService) GetListEnhanced(page, pageSize int, status *int, userID *uint,
	totalSmall, totalBig, createTimeStart, createTimeEnd,
	dueTimeStart, dueTimeEnd, paidTimeStart, paidTimeEnd,
	payment, lineItemDesc, invType, saleID, invoiceID,
	orderField, sort string) ([]Invoice, int64, error) {
	var invoices []Invoice
	var total int64

	query := s.db.Model(&Invoice{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	// 金额区间
	if totalSmall != "" {
		query = query.Where("amount >= ?", totalSmall)
	}
	if totalBig != "" {
		query = query.Where("amount <= ?", totalBig)
	}

	// 时间范围
	if createTimeStart != "" {
		query = query.Where("created_at >= ?", createTimeStart)
	}
	if createTimeEnd != "" {
		query = query.Where("created_at <= ?", createTimeEnd+" 23:59:59")
	}
	if dueTimeStart != "" {
		query = query.Where("due_date >= ?", dueTimeStart)
	}
	if dueTimeEnd != "" {
		query = query.Where("due_date <= ?", dueTimeEnd+" 23:59:59")
	}
	if paidTimeStart != "" {
		query = query.Where("paid_at >= ?", paidTimeStart)
	}
	if paidTimeEnd != "" {
		query = query.Where("paid_at <= ?", paidTimeEnd+" 23:59:59")
	}

	// 付款方式细分
	if payment != "" {
		if payment == "creditLimitPay" {
			query = query.Where("use_credit_limit = 1")
		} else if payment == "creditPay" {
			query = query.Where("credit > 0")
		} else {
			query = query.Where("payment = ?", payment)
		}
	}

	// 行项描述
	if lineItemDesc != "" {
		subQuery := s.db.Table("invoice_items").Select("invoice_id").Where("description LIKE ?", "%"+lineItemDesc+"%")
		query = query.Where("id IN (?)", subQuery)
	}

	// 发票类型
	if invType != "" {
		query = query.Where("type = ?", invType)
	}

	// 销售筛选
	if saleID != "" {
		query = query.Joins("JOIN users ON users.id = invoices.user_id").Where("users.sale_id = ?", saleID)
	}

	// 发票号搜索
	if invoiceID != "" {
		query = query.Where("id LIKE ? OR invoice_no LIKE ?", "%"+invoiceID+"%", "%"+invoiceID+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	orderClause := fmt.Sprintf("%s %s", orderField, sort)
	if err := query.Offset(offset).Limit(limit).Order(orderClause).Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	return invoices, total, nil
}

// UpdateInvoice updates invoice fields.
func (s *InvoiceService) UpdateInvoice(id uint, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	return s.db.Model(&Invoice{}).Where("id = ?", id).Updates(updates).Error
}

// InvoiceUser represents a user for invoice operations.
type InvoiceUser struct {
	ID      uint    `json:"id"`
	Balance float64 `json:"balance"`
	Credit  float64 `json:"credit"`
}

// GetUser returns user info for invoice operations.
func (s *InvoiceService) GetUser(userID uint) (*InvoiceUser, error) {
	var user InvoiceUser
	if err := s.db.Table("users").Where("id = ?", userID).
		Select("id, balance, credit").Scan(&user).Error; err != nil {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

// GetDB returns the database instance.
func (s *InvoiceService) GetDB() *gorm.DB {
	return s.db
}
