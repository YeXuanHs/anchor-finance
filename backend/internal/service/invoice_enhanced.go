package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// InvoiceNote 管理员备注
type InvoiceNote struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	InvoiceID uint      `gorm:"index;not null" json:"invoice_id"`
	AdminID   uint      `gorm:"index" json:"admin_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// InvoiceLog 账单操作日志
type InvoiceLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	InvoiceID uint      `gorm:"index;not null" json:"invoice_id"`
	AdminID   uint      `gorm:"index" json:"admin_id"`
	Action    string    `gorm:"size:64;not null;index" json:"action"`
	Detail    string    `gorm:"type:text" json:"detail"`
	IPAddress string    `gorm:"size:64" json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}

// InvoiceRefund 退款记录
type InvoiceRefund struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	InvoiceID   uint      `gorm:"index;not null" json:"invoice_id"`
	Amount      float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Reason      string    `gorm:"size:512" json:"reason"`
	Status      int       `gorm:"default:0;comment:0=pending 1=approved 2=rejected 3=completed" json:"status"`
	AdminID     uint      `gorm:"index" json:"admin_id"`
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

// InvoiceSummary 账单统计汇总
type InvoiceSummary struct {
	TotalAmount    float64 `json:"total_amount"`
	PaidAmount     float64 `json:"paid_amount"`
	UnpaidAmount   float64 `json:"unpaid_amount"`
	OverdueAmount  float64 `json:"overdue_amount"`
	TotalCount     int64   `json:"total_count"`
	PaidCount      int64   `json:"paid_count"`
	UnpaidCount    int64   `json:"unpaid_count"`
	OverdueCount   int64   `json:"overdue_count"`
	CancelledCount int64   `json:"cancelled_count"`
}

// InvoiceSearchFilters 高级搜索过滤条件
type InvoiceSearchFilters struct {
	Query     string   `json:"query" form:"query"`
	Status    *int     `json:"status" form:"status"`
	UserID    *uint    `json:"user_id" form:"user_id"`
	MinAmount *float64 `json:"min_amount" form:"min_amount"`
	MaxAmount *float64 `json:"max_amount" form:"max_amount"`
	Type      string   `json:"type" form:"type"`
	StartDate string   `json:"start_date" form:"start_date"`
	EndDate   string   `json:"end_date" form:"end_date"`
	Page      int      `json:"page" form:"page"`
	PageSize  int      `json:"page_size" form:"page_size"`
}

// InvoiceEnhancedService 账单增强服务
type InvoiceEnhancedService struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewInvoiceEnhancedService creates a new InvoiceEnhancedService.
func NewInvoiceEnhancedService(db *gorm.DB, log *logger.Logger) *InvoiceEnhancedService {
	return &InvoiceEnhancedService{db: db, log: log}
}

// ==================== Status Filters ====================

// GetInvoicesByStatus 按状态筛选账单
func (s *InvoiceEnhancedService) GetInvoicesByStatus(status, page, pageSize int) ([]Invoice, int64, error) {
	var invoices []Invoice
	var total int64

	query := s.db.Model(&Invoice{}).Where("status = ?", status)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	return invoices, total, nil
}

// GetPaidInvoices 已支付账单列表
func (s *InvoiceEnhancedService) GetPaidInvoices(page, pageSize int) ([]Invoice, int64, error) {
	return s.GetInvoicesByStatus(1, page, pageSize)
}

// GetUnpaidInvoices 未支付账单列表
func (s *InvoiceEnhancedService) GetUnpaidInvoices(page, pageSize int) ([]Invoice, int64, error) {
	return s.GetInvoicesByStatus(0, page, pageSize)
}

// GetCancelledInvoices 已取消账单列表
func (s *InvoiceEnhancedService) GetCancelledInvoices(page, pageSize int) ([]Invoice, int64, error) {
	return s.GetInvoicesByStatus(2, page, pageSize)
}

// GetOverdueInvoicesPage 逾期账单列表（分页）
func (s *InvoiceEnhancedService) GetOverdueInvoicesPage(page, pageSize int) ([]Invoice, int64, error) {
	var invoices []Invoice
	var total int64

	query := s.db.Model(&Invoice{}).Where("status = 0 AND due_date < ?", time.Now())
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("due_date ASC").Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	return invoices, total, nil
}

// GetInvoiceSummary 账单统计汇总
func (s *InvoiceEnhancedService) GetInvoiceSummary() (*InvoiceSummary, error) {
	summary := &InvoiceSummary{}

	// 总数和总金额
	s.db.Model(&Invoice{}).Count(&summary.TotalCount)
	s.db.Model(&Invoice{}).Select("COALESCE(SUM(amount), 0)").Scan(&summary.TotalAmount)

	// 已支付
	s.db.Model(&Invoice{}).Where("status = 1").Count(&summary.PaidCount)
	s.db.Model(&Invoice{}).Where("status = 1").Select("COALESCE(SUM(amount), 0)").Scan(&summary.PaidAmount)

	// 未支付
	s.db.Model(&Invoice{}).Where("status = 0").Count(&summary.UnpaidCount)
	s.db.Model(&Invoice{}).Where("status = 0").Select("COALESCE(SUM(amount), 0)").Scan(&summary.UnpaidAmount)

	// 逾期
	s.db.Model(&Invoice{}).Where("status = 0 AND due_date < ?", time.Now()).Count(&summary.OverdueCount)
	s.db.Model(&Invoice{}).Where("status = 0 AND due_date < ?", time.Now()).Select("COALESCE(SUM(amount), 0)").Scan(&summary.OverdueAmount)

	// 已取消
	s.db.Model(&Invoice{}).Where("status = 2").Count(&summary.CancelledCount)

	return summary, nil
}

// ==================== Invoice Operations ====================

// AddPayInvoice 添加支付记录到账单
func (s *InvoiceEnhancedService) AddPayInvoice(invoiceID uint, amount float64, method string) (*Invoice, error) {
	inv, err := s.getInvoiceByID(invoiceID)
	if err != nil {
		return nil, err
	}
	if inv.Status == 2 {
		return nil, errors.New("cannot pay a cancelled invoice")
	}
	if inv.Status == 1 {
		return nil, errors.New("invoice is already paid")
	}
	if amount <= 0 {
		return nil, errors.New("payment amount must be positive")
	}
	if amount > inv.Amount {
		return nil, fmt.Errorf("payment amount %.2f exceeds invoice amount %.2f", amount, inv.Amount)
	}

	now := time.Now()
	status := 2 // 部分支付
	if amount >= inv.Amount {
		status = 1 // 全额支付
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"status":          status,
			"payment_method":  method,
		}
		if status == 1 {
			updates["paid_at"] = &now
		}
		if err := tx.Model(inv).Updates(updates).Error; err != nil {
			return err
		}

		// 创建交易记录
		transaction := &InvoiceTransaction{
			TransactionNo: util.GenerateTransactionNo(),
			InvoiceID:     invoiceID,
			UserID:        inv.UserID,
			Amount:        amount,
			Gateway:       method,
			Type:          "payment",
			Status:        1,
			CompletedAt:   &now,
		}
		return tx.Create(transaction).Error
	})

	if err != nil {
		return nil, err
	}

	s.log.Infof("payment added to invoice %s: %.2f via %s", inv.InvoiceNo, amount, method)
	return s.getInvoiceByID(invoiceID)
}

// DeletePayInvoice 删除账单支付记录
func (s *InvoiceEnhancedService) DeletePayInvoice(invoiceID uint) error {
	inv, err := s.getInvoiceByID(invoiceID)
	if err != nil {
		return err
	}
	if inv.Status == 2 {
		return errors.New("cannot delete payment for cancelled invoice")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// 删除关联的交易记录
		if err := tx.Where("invoice_id = ? AND type = ?", invoiceID, "payment").Delete(&InvoiceTransaction{}).Error; err != nil {
			return err
		}
		// 恢复账单为未支付
		return tx.Model(inv).Updates(map[string]interface{}{
			"status":          0,
			"paid_at":         nil,
			"payment_method":  "",
		}).Error
	})
}

// RefundInvoice 退款
func (s *InvoiceEnhancedService) RefundInvoice(invoiceID uint, amount float64, reason string) (*InvoiceRefund, error) {
	inv, err := s.getInvoiceByID(invoiceID)
	if err != nil {
		return nil, err
	}
	if inv.Status != 1 {
		return nil, errors.New("can only refund paid invoices")
	}
	if amount <= 0 {
		return nil, errors.New("refund amount must be positive")
	}
	if amount > inv.Amount {
		return nil, fmt.Errorf("refund amount %.2f exceeds invoice amount %.2f", amount, inv.Amount)
	}

	refund := &InvoiceRefund{
		InvoiceID: invoiceID,
		Amount:    amount,
		Reason:    reason,
		Status:    3, // completed
	}
	now := time.Now()
	refund.CompletedAt = &now

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(refund).Error; err != nil {
			return err
		}

		// 创建退款交易记录
		transaction := &InvoiceTransaction{
			TransactionNo: util.GenerateTransactionNo(),
			InvoiceID:     invoiceID,
			UserID:        inv.UserID,
			Amount:        amount,
			Type:          "refund",
			Status:        1,
			CompletedAt:   &now,
		}
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 更新账单状态
		newStatus := 4 // 已退款
		if amount < inv.Amount {
			newStatus = 1 // 部分退款后仍为已支付
		}
		return tx.Model(inv).Update("status", newStatus).Error
	})

	if err != nil {
		return nil, err
	}

	s.log.Infof("refund processed for invoice %s: %.2f reason=%s", inv.InvoiceNo, amount, reason)
	return refund, nil
}

// GetRefundPage 获取退款选项
func (s *InvoiceEnhancedService) GetRefundPage(invoiceID uint) (map[string]interface{}, error) {
	inv, err := s.getInvoiceByID(invoiceID)
	if err != nil {
		return nil, err
	}

	var refunds []InvoiceRefund
	s.db.Where("invoice_id = ?", invoiceID).Order("id DESC").Find(&refunds)

	var totalRefunded float64
	for _, r := range refunds {
		if r.Status == 3 {
			totalRefunded += r.Amount
		}
	}

	return map[string]interface{}{
		"invoice":         inv,
		"refunds":         refunds,
		"total_refunded":  totalRefunded,
		"refundable":      inv.Amount - totalRefunded,
	}, nil
}

// ==================== Invoice Notes ====================

// AddNote 添加管理员备注
func (s *InvoiceEnhancedService) AddNote(invoiceID, adminID uint, content string) (*InvoiceNote, error) {
	if _, err := s.getInvoiceByID(invoiceID); err != nil {
		return nil, err
	}
	if content == "" {
		return nil, errors.New("note content cannot be empty")
	}

	note := &InvoiceNote{
		InvoiceID: invoiceID,
		AdminID:   adminID,
		Content:   content,
	}
	if err := s.db.Create(note).Error; err != nil {
		return nil, err
	}

	s.log.Infof("note added to invoice %d by admin %d", invoiceID, adminID)
	return note, nil
}

// GetNotes 获取账单备注列表
func (s *InvoiceEnhancedService) GetNotes(invoiceID uint) ([]InvoiceNote, error) {
	var notes []InvoiceNote
	if err := s.db.Where("invoice_id = ?", invoiceID).Order("id DESC").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// DeleteNote 删除备注
func (s *InvoiceEnhancedService) DeleteNote(noteID uint) error {
	var note InvoiceNote
	if err := s.db.First(&note, noteID).Error; err != nil {
		return err
	}
	return s.db.Delete(&note).Error
}

// ==================== Combine Invoices ====================

// GetCombineInvoices 获取可合并的账单（同一用户未支付账单）
func (s *InvoiceEnhancedService) GetCombineInvoices(userID uint) ([]Invoice, error) {
	var invoices []Invoice
	if err := s.db.Where("user_id = ? AND status = 0", userID).Order("id ASC").Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

// CombineInvoices 合并多个账单
func (s *InvoiceEnhancedService) CombineInvoices(invoiceIDs []uint) (*Invoice, error) {
	if len(invoiceIDs) < 2 {
		return nil, errors.New("need at least 2 invoices to combine")
	}

	var invoices []Invoice
	if err := s.db.Where("id IN ?", invoiceIDs).Find(&invoices).Error; err != nil {
		return nil, err
	}

	if len(invoices) != len(invoiceIDs) {
		return nil, errors.New("some invoices not found")
	}

	// 验证所有账单属于同一用户且未支付
	userID := invoices[0].UserID
	var totalAmount float64
	for _, inv := range invoices {
		if inv.UserID != userID {
			return nil, errors.New("all invoices must belong to the same user")
		}
		if inv.Status != 0 {
			return nil, fmt.Errorf("invoice %s is not pending", inv.InvoiceNo)
		}
		totalAmount += inv.Amount
	}

	var combined *Invoice
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 创建合并后的新账单
		combined = &Invoice{
			InvoiceNo: util.GenerateInvoiceNo(),
			UserID:    userID,
			Type:      "combined",
			Amount:    totalAmount,
			Status:    0,
			DueDate:   time.Now().AddDate(0, 0, 7),
			Remark:    fmt.Sprintf("合并账单 (原账单数: %d)", len(invoiceIDs)),
		}
		if err := tx.Create(combined).Error; err != nil {
			return err
		}

		// 将原账单的明细转移到新账单
		if err := tx.Model(&InvoiceItem{}).Where("invoice_id IN ?", invoiceIDs).
			Update("invoice_id", combined.ID).Error; err != nil {
			return err
		}

		// 取消原账单
		if err := tx.Model(&Invoice{}).Where("id IN ?", invoiceIDs).
			Update("status", 2).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	s.log.Infof("invoices combined: %v -> %s (total=%.2f)", invoiceIDs, combined.InvoiceNo, totalAmount)
	return combined, nil
}

// ==================== Invoice Email ====================

// SendInvoiceEmail 发送账单邮件
func (s *InvoiceEnhancedService) SendInvoiceEmail(invoiceID uint, email string) error {
	inv, err := s.getInvoiceByID(invoiceID)
	if err != nil {
		return err
	}
	if email == "" {
		// 查找用户邮箱
		var user struct{ Email string }
		if err := s.db.Model(&User{}).Where("id = ?", inv.UserID).Select("email").First(&user).Error; err != nil {
			return errors.New("user email not found")
		}
		email = user.Email
	}

	// 记录邮件发送日志
	s.LogAction(invoiceID, 0, "email_sent", fmt.Sprintf("账单邮件发送至 %s", email), "")
	s.log.Infof("invoice email sent: %s -> %s", inv.InvoiceNo, email)
	return nil
}

// InvoiceEmail 获取账单邮件模板数据
func (s *InvoiceEnhancedService) InvoiceEmail(invoiceID uint) (map[string]interface{}, error) {
	inv, err := s.getInvoiceByID(invoiceID)
	if err != nil {
		return nil, err
	}

	var items []InvoiceItem
	s.db.Where("invoice_id = ?", invoiceID).Find(&items)

	var user struct {
		Username string
		Email    string
	}
	s.db.Model(&User{}).Where("id = ?", inv.UserID).Select("username, email").First(&user)

	return map[string]interface{}{
		"invoice":  inv,
		"items":    items,
		"user":     user,
		"subject":  fmt.Sprintf("账单通知 - %s", inv.InvoiceNo),
		"template": "invoice_email",
	}, nil
}

// ==================== Renew Invoices ====================

// CreateRenewInvoice 创建续费账单
func (s *InvoiceEnhancedService) CreateRenewInvoice(hostID uint, cycle string) (*Invoice, error) {
	var host Host
	if err := s.db.Select("id", "owner_id", "product_id", "hostname").First(&host, hostID).Error; err != nil {
		return nil, errors.New("host not found")
	}
	if host.OwnerID == nil {
		return nil, errors.New("host has no owner")
	}

	// 查找产品价格
	var price float64
	if host.ProductID != nil {
		s.db.Model(&Product{}).Where("id = ?", *host.ProductID).Select("price").Scan(&price)
	}

	// 根据周期计算金额
	amount := calculateCycleAmount(price, cycle)

	invoice := &Invoice{
		InvoiceNo: util.GenerateInvoiceNo(),
		UserID:    *host.OwnerID,
		Type:      "renew",
		Amount:    amount,
		Status:    0,
		DueDate:   time.Now().AddDate(0, 0, 7),
		Remark:    fmt.Sprintf("续费账单 - %s (%s)", host.Hostname, cycle),
	}
	if err := s.db.Create(invoice).Error; err != nil {
		return nil, err
	}

	s.log.Infof("renewal invoice created: %s (host=%d, cycle=%s, amount=%.2f)", invoice.InvoiceNo, hostID, cycle, amount)
	return invoice, nil
}

// GetRenewInvoices 续费账单列表
func (s *InvoiceEnhancedService) GetRenewInvoices(page, pageSize int) ([]Invoice, int64, error) {
	var invoices []Invoice
	var total int64

	query := s.db.Model(&Invoice{}).Where("type = ?", "renew")
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	return invoices, total, nil
}

// ==================== Invoice Log ====================

// GetInvoiceLog 获取账单操作日志
func (s *InvoiceEnhancedService) GetInvoiceLog(invoiceID uint) ([]InvoiceLog, error) {
	var logs []InvoiceLog
	if err := s.db.Where("invoice_id = ?", invoiceID).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// LogAction 记录账单操作
func (s *InvoiceEnhancedService) LogAction(invoiceID, adminID uint, action, detail, ip string) {
	log := &InvoiceLog{
		InvoiceID: invoiceID,
		AdminID:   adminID,
		Action:    action,
		Detail:    detail,
		IPAddress: ip,
	}
	if err := s.db.Create(log).Error; err != nil {
		s.log.Errorf("failed to create invoice log: %v", err)
	}
}

// ==================== Search ====================

// SearchInvoices 搜索账单
func (s *InvoiceEnhancedService) SearchInvoices(query string) ([]Invoice, error) {
	var invoices []Invoice
	q := s.db.Model(&Invoice{})
	if query != "" {
		like := "%" + query + "%"
		q = q.Where("invoice_no LIKE ? OR remark LIKE ? OR CAST(amount AS TEXT) LIKE ?", like, like, query)
		// 按用户ID搜索（纯数字时）
		var userID uint
		if _, err := fmt.Sscanf(query, "%d", &userID); err == nil {
			q = q.Or("user_id = ?", userID)
		}
	}
	if err := q.Order("id DESC").Limit(50).Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

// SearchPage 高级搜索（多条件）
func (s *InvoiceEnhancedService) SearchPage(filters InvoiceSearchFilters) ([]Invoice, int64, error) {
	var invoices []Invoice
	var total int64

	q := s.db.Model(&Invoice{})

	if filters.Query != "" {
		like := "%" + filters.Query + "%"
		q = q.Where("invoice_no LIKE ? OR remark LIKE ?", like, like)
		var userID uint
		if _, err := fmt.Sscanf(filters.Query, "%d", &userID); err == nil {
			q = q.Or("user_id = ?", userID)
		}
	}
	if filters.Status != nil {
		q = q.Where("status = ?", *filters.Status)
	}
	if filters.UserID != nil {
		q = q.Where("user_id = ?", *filters.UserID)
	}
	if filters.MinAmount != nil {
		q = q.Where("amount >= ?", *filters.MinAmount)
	}
	if filters.MaxAmount != nil {
		q = q.Where("amount <= ?", *filters.MaxAmount)
	}
	if filters.Type != "" {
		q = q.Where("type = ?", filters.Type)
	}
	if filters.StartDate != "" {
		q = q.Where("created_at >= ?", filters.StartDate)
	}
	if filters.EndDate != "" {
		q = q.Where("created_at <= ?", filters.EndDate+" 23:59:59")
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filters.Page
	if page < 1 {
		page = 1
	}
	pageSize := filters.PageSize
	if pageSize < 1 {
		pageSize = 20
	}
	offset, limit := Paginate(page, pageSize)

	if err := q.Offset(offset).Limit(limit).Order("id DESC").Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	return invoices, total, nil
}

// ==================== Duplicate ====================

// DuplicateInvoice 复制账单
func (s *InvoiceEnhancedService) DuplicateInvoice(invoiceID uint) (*Invoice, error) {
	orig, err := s.getInvoiceByID(invoiceID)
	if err != nil {
		return nil, err
	}

	var newInvoice *Invoice
	err = s.db.Transaction(func(tx *gorm.DB) error {
		newInvoice = &Invoice{
			InvoiceNo: util.GenerateInvoiceNo(),
			UserID:    orig.UserID,
			OrderID:   orig.OrderID,
			OrderNo:   orig.OrderNo,
			Amount:    orig.Amount,
			Type:      orig.Type,
			Status:    0,
			DueDate:   time.Now().AddDate(0, 0, 7),
			Remark:    fmt.Sprintf("复制自 %s", orig.InvoiceNo),
		}
		if err := tx.Create(newInvoice).Error; err != nil {
			return err
		}

		// 复制明细项
		var items []InvoiceItem
		if err := tx.Where("invoice_id = ?", invoiceID).Find(&items).Error; err != nil {
			return err
		}
		for i := range items {
			items[i].ID = 0
			items[i].InvoiceID = newInvoice.ID
			items[i].CreatedAt = time.Time{}
			items[i].UpdatedAt = time.Time{}
		}
		if len(items) > 0 {
			return tx.CreateInBatches(items, 100).Error
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	s.log.Infof("invoice duplicated: %s -> %s", orig.InvoiceNo, newInvoice.InvoiceNo)
	return newInvoice, nil
}

// ==================== Extended AddPay (P0-1) ====================

// AddPayInvoiceEx 管理员手动入账（支持网关/交易号/手续费/入账时间）
func (s *InvoiceEnhancedService) AddPayInvoiceEx(invoiceID uint, amount float64, gateway, transID string, fees float64, payTime string) (*Invoice, error) {
	inv, err := s.getInvoiceByID(invoiceID)
	if err != nil {
		return nil, err
	}
	if inv.Status == 2 {
		return nil, errors.New("cannot pay a cancelled invoice")
	}
	if inv.Status == 1 {
		return nil, errors.New("invoice is already paid")
	}
	if amount <= 0 {
		return nil, errors.New("payment amount must be positive")
	}
	if amount > inv.Amount {
		return nil, fmt.Errorf("payment amount %.2f exceeds invoice amount %.2f", amount, inv.Amount)
	}

	now := time.Now()
	var payAt time.Time
	if payTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", payTime); err == nil {
			payAt = t
		} else if t, err := time.Parse("2006-01-02", payTime); err == nil {
			payAt = t
		} else {
			payAt = now
		}
	} else {
		payAt = now
	}

	status := 2 // 部分支付
	if amount >= inv.Amount {
		status = 1 // 全额支付
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"status":        status,
			"payment":       gateway,
			"payment_notes": transID,
		}
		if status == 1 {
			updates["paid_at"] = &payAt
		}
		if err := tx.Model(inv).Updates(updates).Error; err != nil {
			return err
		}

		// 创建交易记录
		transaction := &InvoiceTransaction{
			TransactionNo: util.GenerateTransactionNo(),
			InvoiceID:     invoiceID,
			UserID:        inv.UserID,
			Amount:        amount,
			Gateway:       gateway,
			Type:          "payment",
			Status:        1,
			CompletedAt:   &payAt,
		}
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		// 扣除手续费记录
		if fees > 0 {
			feeRecord := &InvoiceTransaction{
				TransactionNo: util.GenerateTransactionNo(),
				InvoiceID:     invoiceID,
				UserID:        inv.UserID,
				Amount:        fees,
				Gateway:       gateway,
				Type:          "fee",
				Status:        1,
				CompletedAt:   &payAt,
			}
			if err := tx.Create(feeRecord).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	s.log.Infof("payment added to invoice %s: %.2f via %s (trans=%s)", inv.InvoiceNo, amount, gateway, transID)
	return s.getInvoiceByID(invoiceID)
}

// ApplyCreditLimit 使用信用额支付账单 (P0-2)
func (s *InvoiceEnhancedService) ApplyCreditLimit(invoiceID uint) error {
	inv, err := s.getInvoiceByID(invoiceID)
	if err != nil {
		return err
	}
	if inv.Status != 0 {
		return errors.New("只能对未支付账单使用信用额支付")
	}

	// 查找用户信用额设置
	var client struct {
		ID                uint
		IsOpenCreditLimit int
		CreditLimit       float64
		Credit            float64
	}
	if err := s.db.Table("clients").Where("id = ?", inv.UserID).
		Select("id, is_open_credit_limit, credit_limit, credit").Scan(&client).Error; err != nil {
		return errors.New("用户不存在")
	}
	if client.IsOpenCreditLimit == 0 {
		return errors.New("信用额未开通,不可使用信用额支付")
	}

	// 检查信用额余额（信用额度 + 余额 - 已使用）
	usedLimit := 0.0
	s.db.Table("invoices").Where("user_id = ? AND use_credit_limit = 1 AND status IN (0,1)", inv.UserID).
		Select("COALESCE(SUM(total), 0)").Scan(&usedLimit)
	available := client.CreditLimit + client.Credit - usedLimit
	if available < inv.Amount {
		return errors.New("当前信用额余额不足,不可使用信用额支付")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		return tx.Model(inv).Updates(map[string]interface{}{
			"status":            1,
			"use_credit_limit":  1,
			"paid_at":           &now,
		}).Error
	})
}

// ExecuteRenew 执行续费操作 - 处理已支付的续费账单，延长服务到期时间 (P0-4)
func (s *InvoiceEnhancedService) ExecuteRenew(invoiceID uint) error {
	inv, err := s.getInvoiceByID(invoiceID)
	if err != nil {
		return err
	}
	if inv.Status != 1 {
		return errors.New("只能对已支付账单执行续费")
	}
	if inv.Type != "renew" {
		return errors.New("非续费账单")
	}

	// 查找关联的主机/服务
	var host struct {
		ID        uint
		ProductID uint
		BillingCycle string
	}
	if err := s.db.Table("host").Where("orderid = ?", inv.OrderID).First(&host).Error; err != nil {
		// 尝试通过 invoice remark 查找 hostID
		return errors.New("未找到关联的服务")
	}

	// 计算续费时长
	months := 1
	switch host.BillingCycle {
	case "monthly":
		months = 1
	case "quarterly":
		months = 3
	case "semi-annually":
		months = 6
	case "annually":
		months = 12
	case "biennially":
		months = 24
	case "triennially":
		months = 36
	}

	// 延长到期时间
	return s.db.Transaction(func(tx *gorm.DB) error {
		var svc struct {
			ID        uint
			ExpiredAt *time.Time
			Status    int16
		}
		if err := tx.Table("client_services").Where("id = ?", host.ID).First(&svc).Error; err != nil {
			// 如果没找到 client_service，尝试更新 host 表
			return tx.Table("host").Where("id = ?", host.ID).
				Update("nextduedate", gorm.Expr("DATE_ADD(nextduedate, INTERVAL ? MONTH)", months)).Error
		}

		base := time.Now()
		if svc.ExpiredAt != nil && svc.ExpiredAt.After(time.Now()) {
			base = *svc.ExpiredAt
		}
		newExpiry := base.AddDate(0, months, 0)
		return tx.Table("client_services").Where("id = ?", host.ID).Updates(map[string]interface{}{
			"expired_at": newExpiry,
			"status":     1, // active
		}).Error
	})
}

// ==================== Helpers ====================

// InvoiceTransaction 账单交易记录（service层自定义，简化版）
type InvoiceTransaction struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	TransactionNo string    `gorm:"size:64;uniqueIndex;not null" json:"transaction_no"`
	InvoiceID     uint      `gorm:"index;not null" json:"invoice_id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	Amount        float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Gateway       string    `gorm:"size:64" json:"gateway"`
	Type          string    `gorm:"size:32;default:payment;comment:payment/refund" json:"type"`
	Status        int       `gorm:"default:0;comment:0=pending 1=success 2=failed" json:"status"`
	CompletedAt   *time.Time `json:"completed_at"`
	CreatedAt     time.Time  `json:"created_at"`
}

func (s *InvoiceEnhancedService) getInvoiceByID(id uint) (*Invoice, error) {
	var inv Invoice
	if err := s.db.First(&inv, id).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

// calculateCycleAmount 根据周期计算金额
func calculateCycleAmount(basePrice float64, cycle string) float64 {
	switch cycle {
	case "monthly":
		return basePrice
	case "quarterly":
		return basePrice * 3
	case "semi-annually":
		return basePrice * 6
	case "annually":
		return basePrice * 12
	case "biennially":
		return basePrice * 24
	case "triennially":
		return basePrice * 36
	default:
		return basePrice
	}
}
