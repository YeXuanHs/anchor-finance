package service

import (
	"context"
	"fmt"
	"time"

	"github.com/anchor-finance/backend/pkg/payment"
)

// CreatePaymentInput contains optional input for creating a payment.
type CreatePaymentInput struct {
	Amount  float64
	Subject string
}

// CreatePaymentOutput contains the result of creating a payment.
type CreatePaymentOutput struct {
	TradeNo   string
	PayURL    string
	QrcodeURL string
}

// TransactionStore defines the interface for persisting payment transactions.
type TransactionStore interface {
	CreateTransaction(tx *Transaction) error
	UpdateTransactionStatus(tradeNo, status string) error
	GetTransactionByTradeNo(tradeNo string) (*Transaction, error)
	GetTransactionByInvoiceID(invoiceID string) (*Transaction, error)
}

// InvoiceStore defines the interface for accessing invoices.
type InvoiceStore interface {
	GetInvoice(invoiceID string) (*Invoice, error)
	UpdateInvoiceStatus(invoiceID, status string) error
}

// Invoice represents an invoice record.
type Invoice struct {
	ID     string
	UserID string
	Amount float64
	Status string // pending, paid, refunded
}

// Transaction represents a payment transaction record.
type Transaction struct {
	ID        string
	TradeNo   string
	OrderNo   string
	UserID    string
	InvoiceID string
	Gateway   string
	Amount    float64
	Status    string // pending, success, failed, closed
	CreatedAt time.Time
	PaidAt    *time.Time
}

// PaymentService provides business logic for payment operations.
type PaymentService struct {
	txStore    TransactionStore
	invoiceStore InvoiceStore
}

// NewPaymentService creates a new PaymentService.
func NewPaymentService(txStore TransactionStore, invoiceStore InvoiceStore) *PaymentService {
	return &PaymentService{
		txStore:      txStore,
		invoiceStore: invoiceStore,
	}
}

// CreatePayment creates a new payment through the specified gateway.
func (s *PaymentService) CreatePayment(ctx context.Context, userID, invoiceID, gatewayName string, input *CreatePaymentInput) (*CreatePaymentOutput, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if invoiceID == "" {
		return nil, fmt.Errorf("invoice_id is required")
	}

	gw, ok := payment.Get(gatewayName)
	if !ok {
		return nil, fmt.Errorf("payment gateway %q not found", gatewayName)
	}

	// Fetch invoice to determine amount.
	var amount float64
	var subject string
	if s.invoiceStore != nil {
		invoice, err := s.invoiceStore.GetInvoice(invoiceID)
		if err != nil {
			return nil, fmt.Errorf("failed to get invoice: %w", err)
		}
		if invoice.Status != "pending" {
			return nil, fmt.Errorf("invoice is not pending (current status: %s)", invoice.Status)
		}
		amount = invoice.Amount
		subject = fmt.Sprintf("Payment for invoice %s", invoiceID)
	}

	// Override with explicit input if provided.
	if input != nil {
		if input.Amount > 0 {
			amount = input.Amount
		}
		if input.Subject != "" {
			subject = input.Subject
		}
	}

	if amount <= 0 {
		return nil, fmt.Errorf("payment amount must be positive")
	}

	orderNo := fmt.Sprintf("ORD_%s_%d", invoiceID, time.Now().UnixNano())

	// For balance payments, inject userID into context.
	if gatewayName == "balance" {
		ctx = context.WithValue(ctx, "userID", userID)
	}

	param := &payment.PaymentParam{
		OrderNo:  orderNo,
		Amount:   amount,
		Subject:  subject,
		ClientIP: "",
	}

	result, err := gw.CreatePayment(ctx, param)
	if err != nil {
		return nil, fmt.Errorf("gateway create payment failed: %w", err)
	}

	// Persist transaction record.
	if s.txStore != nil {
		tx := &Transaction{
			TradeNo:   result.TradeNo,
			OrderNo:   orderNo,
			UserID:    userID,
			InvoiceID: invoiceID,
			Gateway:   gatewayName,
			Amount:    amount,
			Status:    "pending",
			CreatedAt: time.Now(),
		}
		if err := s.txStore.CreateTransaction(tx); err != nil {
			return nil, fmt.Errorf("failed to save transaction: %w", err)
		}
	}

	return &CreatePaymentOutput{
		TradeNo:   result.TradeNo,
		PayURL:    result.PayURL,
		QrcodeURL: result.QrcodeURL,
	}, nil
}

// HandleNotify processes an async payment notification from a gateway.
func (s *PaymentService) HandleNotify(ctx context.Context, gatewayName string, data []byte) error {
	gw, ok := payment.Get(gatewayName)
	if !ok {
		return fmt.Errorf("payment gateway %q not found", gatewayName)
	}

	notify, err := gw.ParseNotify(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to parse notification: %w", err)
	}

	if notify.Status != "success" {
		return fmt.Errorf("payment not successful, status: %s", notify.Status)
	}

	// Update transaction status.
	if s.txStore != nil && notify.TradeNo != "" {
		if err := s.txStore.UpdateTransactionStatus(notify.TradeNo, "success"); err != nil {
			return fmt.Errorf("failed to update transaction: %w", err)
		}

		// Update associated invoice.
		if s.invoiceStore != nil {
			tx, err := s.txStore.GetTransactionByTradeNo(notify.TradeNo)
			if err != nil {
				return fmt.Errorf("failed to get transaction: %w", err)
			}
			if tx != nil && tx.InvoiceID != "" {
				if err := s.invoiceStore.UpdateInvoiceStatus(tx.InvoiceID, "paid"); err != nil {
					return fmt.Errorf("failed to update invoice: %w", err)
				}
			}
		}
	}

	return nil
}

// ProcessRefund processes a refund for a completed transaction.
func (s *PaymentService) ProcessRefund(ctx context.Context, transactionID string, amount float64) error {
	if transactionID == "" {
		return fmt.Errorf("transaction_id is required")
	}
	if amount <= 0 {
		return fmt.Errorf("refund amount must be positive")
	}

	if s.txStore == nil {
		return fmt.Errorf("transaction store not configured")
	}

	tx, err := s.txStore.GetTransactionByTradeNo(transactionID)
	if err != nil {
		return fmt.Errorf("failed to get transaction: %w", err)
	}
	if tx == nil {
		return fmt.Errorf("transaction not found")
	}
	if tx.Status != "success" {
		return fmt.Errorf("can only refund successful transactions (current status: %s)", tx.Status)
	}
	if amount > tx.Amount {
		return fmt.Errorf("refund amount %.2f exceeds original amount %.2f", amount, tx.Amount)
	}

	gw, ok := payment.Get(tx.Gateway)
	if !ok {
		return fmt.Errorf("payment gateway %q not found", tx.Gateway)
	}

	refundNo := fmt.Sprintf("REF_%s_%d", transactionID, time.Now().UnixNano())

	// For balance payments, inject userID into context.
	if tx.Gateway == "balance" {
		ctx = context.WithValue(ctx, "userID", tx.UserID)
	}

	refundParam := &payment.RefundParam{
		TradeNo:  transactionID,
		RefundNo: refundNo,
		Amount:   amount,
		Reason:   "merchant refund",
	}

	result, err := gw.Refund(ctx, refundParam)
	if err != nil {
		return fmt.Errorf("gateway refund failed: %w", err)
	}

	if result.Status == "success" {
		if err := s.txStore.UpdateTransactionStatus(transactionID, "refunded"); err != nil {
			return fmt.Errorf("failed to update transaction status: %w", err)
		}
		if s.invoiceStore != nil && tx.InvoiceID != "" {
			if err := s.invoiceStore.UpdateInvoiceStatus(tx.InvoiceID, "refunded"); err != nil {
				return fmt.Errorf("failed to update invoice status: %w", err)
			}
		}
	}

	return nil
}
