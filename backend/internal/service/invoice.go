package service

import (
	"errors"
	"time"

	"github.com/anchor-finance/backend/internal/util"
	"github.com/anchor-finance/backend/pkg/logger"

	"gorm.io/gorm"
)

type Invoice struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	InvoiceNo  string         `gorm:"uniqueIndex;size:64;not null" json:"invoice_no"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	OrderID    uint           `json:"order_id"`
	OrderNo    string         `gorm:"size:64" json:"order_no"`
	Amount     float64        `gorm:"type:decimal(12,2);not null" json:"amount"`
	Type       string         `gorm:"size:32;default:new;comment:new/renew" json:"type"`
	Status     int            `gorm:"default:0;comment:0=pending 1=paid 2=cancelled 3=overdue" json:"status"`
	PaidAt     *time.Time     `json:"paid_at"`
	DueDate    time.Time      `gorm:"index" json:"due_date"`
	Remark     string         `gorm:"size:256" json:"remark"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
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
