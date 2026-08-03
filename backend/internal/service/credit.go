package service

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type CreditService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCreditService(db *gorm.DB, log *logger.Logger) *CreditService {
	return &CreditService{db: db, log: log}
}

func (s *CreditService) GetByUserID(userID uint) (*model.CreditLimit, error) {
	var credit model.CreditLimit
	if err := s.db.Where("user_id = ?", userID).First(&credit).Error; err != nil {
		return nil, err
	}
	return &credit, nil
}

func (s *CreditService) SetLimit(userID uint, limit float64) error {
	var credit model.CreditLimit
	err := s.db.Where("user_id = ?", userID).First(&credit).Error
	if err == gorm.ErrRecordNotFound {
		credit = model.CreditLimit{UserID: userID, Limit: limit, Used: 0}
		return s.db.Create(&credit).Error
	}
	return s.db.Model(&credit).Update("limit", limit).Error
}

func (s *CreditService) UseCredit(userID uint, amount float64) error {
	result := s.db.Model(&model.CreditLimit{}).
		Where("user_id = ? AND `limit` - `used` >= ?", userID, amount).
		Update("used", gorm.Expr("used + ?", amount))
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (s *CreditService) ReleaseCredit(userID uint, amount float64) error {
	return s.db.Model(&model.CreditLimit{}).Where("user_id = ?", userID).
		Update("used", gorm.Expr("GREATEST(used - ?, 0)", amount)).Error
}

func (s *CreditService) GetLog(userID uint, page, pageSize int) ([]model.CreditLog, int64, error) {
	var logs []model.CreditLog
	var total int64

	query := s.db.Model(&model.CreditLog{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// ─────────────────────── Billing Cycle ───────────────────────

// GenerateMonthlyBills generates monthly credit bills for all users with credit usage.
// Called by cron on each billing cycle.
func (s *CreditService) GenerateMonthlyBills() (string, error) {
	now := time.Now()
	billMonth := now.Format("2006-01")

	var credits []model.CreditLimit
	if err := s.db.Where("used > 0").Find(&credits).Error; err != nil {
		return "", fmt.Errorf("query credit limits: %w", err)
	}

	created := 0
	for _, credit := range credits {
		// Skip if bill already exists for this month
		var existingBill model.CreditBill
		if err := s.db.Where("user_id = ? AND bill_month = ?", credit.UserID, billMonth).First(&existingBill).Error; err == nil {
			continue
		}

		billDay := credit.BillGenerationDay
		if billDay < 1 || billDay > 28 {
			billDay = 1
		}
		repayDays := credit.RepaymentPeriod
		if repayDays <= 0 {
			repayDays = 15
		}

		billingDate := time.Date(now.Year(), now.Month(), billDay, 0, 0, 0, 0, now.Location())
		dueDate := billingDate.AddDate(0, 0, repayDays)

		bill := &model.CreditBill{
			UserID:          credit.UserID,
			BillMonth:       billMonth,
			BillingDate:     billingDate,
			DueDate:         dueDate,
			TotalAmount:     credit.Used,
			PaidAmount:      0,
			RemainingAmount: credit.Used,
			LateFee:         0,
			Status:          "unpaid",
		}
		if err := s.db.Create(bill).Error; err != nil {
			s.log.Warnf("create bill for user %d: %v", credit.UserID, err)
			continue
		}

		// Create bill item for usage
		item := &model.CreditBillItem{
			BillID:      bill.ID,
			Type:        "usage",
			Description: fmt.Sprintf("信用额度使用 - %s", billMonth),
			Amount:      credit.Used,
		}
		s.db.Create(item)

		created++
		s.log.Infof("generated bill for user %d: month=%s amount=%.4f", credit.UserID, billMonth, credit.Used)
	}

	output := fmt.Sprintf("Checked %d users with credit usage, created %d bills for %s", len(credits), created, billMonth)
	return output, nil
}

// GetBills returns paginated bills for a user.
func (s *CreditService) GetBills(userID uint, page, pageSize int) ([]model.CreditBill, int64, error) {
	var bills []model.CreditBill
	var total int64

	query := s.db.Model(&model.CreditBill{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&bills).Error; err != nil {
		return nil, 0, err
	}
	return bills, total, nil
}

// GetBillByID returns a single bill by ID.
func (s *CreditService) GetBillByID(id uint) (*model.CreditBill, error) {
	var bill model.CreditBill
	if err := s.db.First(&bill, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bill not found")
		}
		return nil, err
	}
	return &bill, nil
}

// PayBill processes a bill payment (partial or full).
func (s *CreditService) PayBill(billID uint, amount float64, paymentMethod string) (*model.CreditBill, error) {
	var bill model.CreditBill
	if err := s.db.First(&bill, billID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bill not found")
		}
		return nil, err
	}

	if bill.Status == "paid" || bill.Status == "written_off" {
		return nil, errors.New("bill is already settled")
	}

	totalDue := bill.RemainingAmount + bill.LateFee
	if amount > totalDue {
		amount = totalDue
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		newPaid := bill.PaidAmount + amount
		newRemaining := bill.TotalAmount - newPaid

		// Apply payment to late fee first, then remaining
		lateFeePaid := amount
		if lateFeePaid > bill.LateFee {
			lateFeePaid = bill.LateFee
		}
		newLateFee := bill.LateFee - lateFeePaid

		newStatus := "partial"
		if newRemaining <= 0 && newLateFee <= 0 {
			newStatus = "paid"
			newRemaining = 0
		}

		updates := map[string]interface{}{
			"paid_amount":      newPaid,
			"remaining_amount": newRemaining,
			"late_fee":         newLateFee,
			"status":           newStatus,
		}
		if err := tx.Model(&bill).Updates(updates).Error; err != nil {
			return err
		}

		// Release credit if fully paid
		if newStatus == "paid" {
			if err := tx.Model(&model.CreditLimit{}).Where("user_id = ?", bill.UserID).
				Update("used", gorm.Expr("GREATEST(used - ?, 0)", bill.TotalAmount)).Error; err != nil {
				return err
			}

			// Also update available
			if err := tx.Model(&model.CreditLimit{}).Where("user_id = ?", bill.UserID).
				Update("available", gorm.Expr("available + ?", bill.TotalAmount)).Error; err != nil {
				return err
			}
		}

		// Create payment log
		creditLog := model.CreditLog{
			UserID:      bill.UserID,
			Type:        "repay",
			Amount:      amount,
			Balance:     0,
			RelatedID:   billID,
			RelatedType: "credit_bill",
			Remark:      fmt.Sprintf("Bill payment for %s via %s", bill.BillMonth, paymentMethod),
		}
		return tx.Create(&creditLog).Error
	})

	if err != nil {
		return nil, err
	}

	// Reload
	s.db.First(&bill, billID)
	return &bill, nil
}

// CalculateLateFee calculates and applies late fee for an overdue bill.
func (s *CreditService) CalculateLateFee(billID uint) (float64, error) {
	var bill model.CreditBill
	if err := s.db.First(&bill, billID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("bill not found")
		}
		return 0, err
	}

	if bill.Status == "paid" || bill.Status == "written_off" {
		return bill.LateFee, nil
	}

	now := time.Now()
	if now.Before(bill.DueDate) || now.Equal(bill.DueDate) {
		return 0, nil
	}

	// Calculate overdue days
	overdueDays := int(now.Sub(bill.DueDate).Hours() / 24)
	if overdueDays <= 0 {
		return 0, nil
	}

	// Late fee: 0.05% per day on remaining amount, capped at 10% of total
	dailyRate := 0.0005
	maxRate := 0.10
	lateFee := bill.RemainingAmount * dailyRate * float64(overdueDays)
	maxFee := bill.TotalAmount * maxRate
	if lateFee > maxFee {
		lateFee = maxFee
	}

	// Round to 4 decimal places
	lateFee = float64(int(lateFee*10000)) / 10000

	if lateFee > bill.LateFee {
		newLateFee := lateFee - bill.LateFee
		s.db.Model(&bill).Updates(map[string]interface{}{
			"late_fee": lateFee,
			"status":   "overdue",
		})
		s.log.Infof("applied late fee for bill %d: days=%d fee=%.4f", billID, overdueDays, newLateFee)
	}

	return lateFee, nil
}

// GetCreditConfig returns the user's credit configuration.
func (s *CreditService) GetCreditConfig(userID uint) (*model.CreditLimit, error) {
	var credit model.CreditLimit
	if err := s.db.Where("user_id = ?", userID).First(&credit).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return default config
			return &model.CreditLimit{
				UserID:            userID,
				BillGenerationDay: 1,
				RepaymentPeriod:   15,
			}, nil
		}
		return nil, err
	}
	return &credit, nil
}

// UpdateCreditConfig updates the user's credit billing configuration.
func (s *CreditService) UpdateCreditConfig(userID uint, billDay int, repayPeriod int) error {
	if billDay < 1 || billDay > 28 {
		return errors.New("bill generation day must be between 1 and 28")
	}
	if repayPeriod < 1 || repayPeriod > 60 {
		return errors.New("repayment period must be between 1 and 60 days")
	}

	var credit model.CreditLimit
	err := s.db.Where("user_id = ?", userID).First(&credit).Error
	if err == gorm.ErrRecordNotFound {
		credit = model.CreditLimit{
			UserID:            userID,
			BillGenerationDay: billDay,
			RepaymentPeriod:   repayPeriod,
		}
		return s.db.Create(&credit).Error
	}
	if err != nil {
		return err
	}

	return s.db.Model(&credit).Updates(map[string]interface{}{
		"bill_generation_day": billDay,
		"repayment_period":    repayPeriod,
	}).Error
}

// ApplyLateFees applies late fees to all overdue bills. Called by cron.
func (s *CreditService) ApplyLateFees() (string, error) {
	now := time.Now()

	var bills []model.CreditBill
	if err := s.db.Where("status IN ? AND due_date < ?", []string{"unpaid", "partial", "overdue"}, now).Find(&bills).Error; err != nil {
		return "", fmt.Errorf("query overdue bills: %w", err)
	}

	applied := 0
	for _, bill := range bills {
		if _, err := s.CalculateLateFee(bill.ID); err != nil {
			s.log.Warnf("calculate late fee for bill %d: %v", bill.ID, err)
			continue
		}
		applied++
	}

	output := fmt.Sprintf("Checked %d overdue bills, applied late fees to %d bills", len(bills), applied)
	return output, nil
}

// GetAdminBills returns all bills with filters (admin).
func (s *CreditService) GetAdminBills(page, pageSize int, userID uint, status string) ([]model.CreditBill, int64, error) {
	var bills []model.CreditBill
	var total int64

	query := s.db.Model(&model.CreditBill{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&bills).Error; err != nil {
		return nil, 0, err
	}
	return bills, total, nil
}

// WaiveLateFee waives the late fee for a bill (admin).
func (s *CreditService) WaiveLateFee(billID uint, adminID uint) (*model.CreditBill, error) {
	var bill model.CreditBill
	if err := s.db.First(&bill, billID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("bill not found")
		}
		return nil, err
	}

	if bill.LateFee <= 0 {
		return nil, errors.New("no late fee to waive")
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		waivedAmount := bill.LateFee

		newStatus := bill.Status
		if bill.RemainingAmount <= 0 {
			newStatus = "paid"
		} else if bill.Status == "overdue" {
			newStatus = "partial"
		}

		if err := tx.Model(&bill).Updates(map[string]interface{}{
			"late_fee": 0,
			"status":   newStatus,
		}).Error; err != nil {
			return err
		}

		// Log the waiver
		creditLog := model.CreditLog{
			UserID:      bill.UserID,
			Type:        "adjust",
			Amount:      waivedAmount,
			Balance:     0,
			RelatedID:   billID,
			RelatedType: "late_fee_waive",
			AdminID:     &adminID,
			Remark:      fmt.Sprintf("Late fee waived for bill %s", bill.BillMonth),
		}
		return tx.Create(&creditLog).Error
	})

	if err != nil {
		return nil, err
	}

	s.db.First(&bill, billID)
	return &bill, nil
}

// GetStatement returns a bill with all its items for detailed viewing.
func (s *CreditService) GetStatement(billID uint) (*model.CreditBill, []model.CreditBillItem, error) {
	var bill model.CreditBill
	if err := s.db.First(&bill, billID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("bill not found")
		}
		return nil, nil, err
	}

	var items []model.CreditBillItem
	if err := s.db.Where("bill_id = ?", billID).Order("id ASC").Find(&items).Error; err != nil {
		return nil, nil, err
	}

	return &bill, items, nil
}

// GetBillItems returns all items for a specific bill.
func (s *CreditService) GetBillItems(billID uint) ([]model.CreditBillItem, error) {
	var items []model.CreditBillItem
	if err := s.db.Where("bill_id = ?", billID).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// GetUsageSummary returns the credit usage summary for a user.
func (s *CreditService) GetUsageSummary(userID uint) (map[string]interface{}, error) {
	credit, err := s.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]interface{}{
				"limit":     0,
				"used":      0,
				"available": 0,
			}, nil
		}
		return nil, err
	}

	var unpaidBills int64
	s.db.Model(&model.CreditBill{}).Where("user_id = ? AND status IN ?", userID, []string{"unpaid", "partial", "overdue"}).Count(&unpaidBills)

	var overdueAmount float64
	s.db.Model(&model.CreditBill{}).Where("user_id = ? AND status = ?", userID, "overdue").
		Select("COALESCE(SUM(remaining_amount + late_fee), 0)").Scan(&overdueAmount)

	return map[string]interface{}{
		"limit":          credit.Limit,
		"used":           credit.Used,
		"available":      credit.Limit - credit.Used,
		"unpaid_bills":   unpaidBills,
		"overdue_amount": overdueAmount,
	}, nil
}

// ─────────────────────── Admin Client Management ───────────────────────

// CreditClientItem represents a credit-enabled user in the client list.
type CreditClientItem struct {
	ID                  uint    `json:"id"`
	UserID              uint    `json:"user_id"`
	Username            string  `json:"username"`
	Email               string  `json:"email"`
	Phone               string  `json:"phone"`
	CreditLimit         float64 `json:"credit_limit"`
	CreditUsed          float64 `json:"credit_used"`
	CreditAvailable     float64 `json:"credit_available"`
	AmountToBeSettled   float64 `json:"amount_to_be_settled"`
	CreditLimitUnpaid   float64 `json:"credit_limit_unpaid"`
	PaymentStatus       string  `json:"payment_status"`
	BillGenerationDay   int     `json:"bill_generation_day"`
	RepaymentPeriod     int     `json:"repayment_period"`
	CreatedAt           string  `json:"created_at"`
}

// GetClientList returns all credit-enabled users with search/filter/sort.
func (s *CreditService) GetClientList(page, pageSize int, keyword, paymentStatus string, order, sort string) ([]CreditClientItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Sub-query: get all credit-enabled user IDs
	creditQuery := s.db.Model(&model.CreditLimit{}).Select("user_id")

	if paymentStatus != "" {
		// Filter by payment status - get user IDs with matching status
		now := time.Now()
		billQuery := s.db.Model(&model.CreditBill{}).Select("DISTINCT user_id")
		switch paymentStatus {
		case "prepayment":
			billQuery = billQuery.Where("status = ? AND created_at >= ?", "paid", beginningOfMonth())
		case "paid":
			billQuery = billQuery.Where("status = ? AND remaining_amount <= 0", "paid")
		case "unpaid":
			billQuery = billQuery.Where("status IN ? AND due_date > ?", []string{"unpaid", "partial"}, now)
		case "overdue":
			billQuery = billQuery.Where("status IN ? AND due_date <= ?", []string{"unpaid", "partial", "overdue"}, now)
		}
		var userIDs []uint
		billQuery.Pluck("user_id", &userIDs)
		if len(userIDs) == 0 {
			return []CreditClientItem{}, 0, nil
		}
		creditQuery = creditQuery.Where("user_id IN ?", userIDs)
	}

	var total int64
	creditQuery.Count(&total)

	// Get paginated credit limits
	var credits []model.CreditLimit
	offset := (page - 1) * pageSize
	query := s.db.Model(&model.CreditLimit{})
	if paymentStatus != "" {
		var userIDs []uint
		creditQuery.Pluck("user_id", &userIDs)
		query = query.Where("user_id IN ?", userIDs)
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&credits).Error; err != nil {
		return nil, 0, err
	}

	var items []CreditClientItem
	for _, credit := range credits {
		var user model.User
		s.db.First(&user, credit.UserID)

		// Apply keyword filter
		if keyword != "" {
			match := false
			if containsIgnoreCase(user.Username, keyword) || containsIgnoreCase(user.Email, keyword) || containsIgnoreCase(user.Phone, keyword) {
				match = true
			}
			if !match {
				continue
			}
		}

		// Calculate amounts
		amountToBeSettled := s.getAmountToBeSettled(credit.UserID)
		unpaid := s.getUnpaidAmount(credit.UserID)
		used := amountToBeSettled + unpaid
		available := credit.Limit - used
		if available < 0 {
			available = 0
		}

		paymentStatusStr := s.getUserPaymentStatus(credit.UserID)

		items = append(items, CreditClientItem{
			ID:                credit.ID,
			UserID:            credit.UserID,
			Username:          user.Username,
			Email:             user.Email,
			Phone:             user.Phone,
			CreditLimit:       credit.Limit,
			CreditUsed:        used,
			CreditAvailable:   available,
			AmountToBeSettled: amountToBeSettled,
			CreditLimitUnpaid: unpaid,
			PaymentStatus:     paymentStatusStr,
			BillGenerationDay: credit.BillGenerationDay,
			RepaymentPeriod:   credit.RepaymentPeriod,
			CreatedAt:         credit.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Handle sorting for computed fields
	if order == "payment_status" || order == "amount_to_be_settled" || order == "credit_limit_unpaid" || order == "credit_available" {
		sortField := order
		sortDir := sort
		sortClientItems(items, sortField, sortDir)
	}

	return items, total, nil
}

func (s *CreditService) getAmountToBeSettled(userID uint) float64 {
	var total float64
	s.db.Model(&model.CreditBill{}).
		Where("user_id = ? AND status = ?", userID, "paid").
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&total)
	return total
}

func (s *CreditService) getUnpaidAmount(userID uint) float64 {
	var total float64
	s.db.Model(&model.CreditBill{}).
		Where("user_id = ? AND status IN ?", userID, []string{"unpaid", "partial", "overdue"}).
		Select("COALESCE(SUM(remaining_amount + late_fee), 0)").
		Scan(&total)
	return total
}

func (s *CreditService) getUserPaymentStatus(userID uint) string {
	now := time.Now()
	var bill model.CreditBill
	err := s.db.Where("user_id = ?", userID).Order("id DESC").First(&bill).Error
	if err != nil {
		return ""
	}
	switch bill.Status {
	case "paid":
		return "paid"
	case "partial":
		return "partial"
	case "overdue":
		return "overdue"
	default: // unpaid
		if now.After(bill.DueDate) {
			return "overdue"
		}
		return "unpaid"
	}
}

func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) &&
		strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

func sortClientItems(items []CreditClientItem, field, dir string) {
	sort.Slice(items, func(i, j int) bool {
		var less bool
		switch field {
		case "payment_status":
			less = items[i].PaymentStatus < items[j].PaymentStatus
		case "amount_to_be_settled":
			less = items[i].AmountToBeSettled < items[j].AmountToBeSettled
		case "credit_limit_unpaid":
			less = items[i].CreditLimitUnpaid < items[j].CreditLimitUnpaid
		case "credit_available":
			less = items[i].CreditAvailable < items[j].CreditAvailable
		default:
			return false
		}
		if dir == "desc" {
			return !less
		}
		return less
	})
}

func beginningOfMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}

// EnableCredit enables credit for a user with initial settings.
func (s *CreditService) EnableCredit(userID uint, limit float64, billDay, repayPeriod int) error {
	if billDay < 1 || billDay > 28 {
		billDay = 1
	}
	if repayPeriod < 1 || repayPeriod > 60 {
		repayPeriod = 15
	}

	var credit model.CreditLimit
	err := s.db.Where("user_id = ?", userID).First(&credit).Error
	if err == gorm.ErrRecordNotFound {
		credit = model.CreditLimit{
			UserID:            userID,
			Limit:             limit,
			Available:         limit,
			BillGenerationDay: billDay,
			RepaymentPeriod:   repayPeriod,
		}
		return s.db.Create(&credit).Error
	}
	if err != nil {
		return err
	}

	// Already exists, update
	return s.db.Model(&credit).Updates(map[string]interface{}{
		"limit":               limit,
		"available":           limit - credit.Used,
		"bill_generation_day": billDay,
		"repayment_period":    repayPeriod,
	}).Error
}

// DisableCredit disables credit for a user.
func (s *CreditService) DisableCredit(userID uint) error {
	result := s.db.Model(&model.CreditLimit{}).Where("user_id = ?", userID).Delete(&model.CreditLimit{})
	if result.RowsAffected == 0 {
		return errors.New("user credit not found")
	}
	return result.Error
}

// UpdateUserCreditSettings updates user credit settings (limit, repayment date, bill time).
func (s *CreditService) UpdateUserCreditSettings(userID uint, limit *float64, billDay *int, repayPeriod *int) error {
	var credit model.CreditLimit
	if err := s.db.Where("user_id = ?", userID).First(&credit).Error; err != nil {
		return errors.New("user credit not found")
	}

	updates := map[string]interface{}{}
	if limit != nil {
		updates["limit"] = *limit
		updates["available"] = *limit - credit.Used
	}
	if billDay != nil {
		if *billDay < 1 || *billDay > 28 {
			return errors.New("bill generation day must be between 1 and 28")
		}
		updates["bill_generation_day"] = *billDay
	}
	if repayPeriod != nil {
		if *repayPeriod < 1 || *repayPeriod > 60 {
			return errors.New("repayment period must be between 1 and 60 days")
		}
		updates["repayment_period"] = *repayPeriod
	}

	if len(updates) == 0 {
		return nil
	}

	return s.db.Model(&credit).Updates(updates).Error
}

// CreditInvoiceItem represents a credit invoice in the list.
type CreditInvoiceItem struct {
	ID              uint    `json:"id"`
	UserID          uint    `json:"user_id"`
	Username        string  `json:"username"`
	TotalAmount     float64 `json:"total_amount"`
	PaidAmount      float64 `json:"paid_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	LateFee         float64 `json:"late_fee"`
	Status          string  `json:"status"`
	PaymentStatus   string  `json:"payment_status"`
	DueDate         string  `json:"due_date"`
	PaidDate        string  `json:"paid_date"`
	CreatedAt       string  `json:"created_at"`
	BillMonth       string  `json:"bill_month"`
	DueDays         int     `json:"due_days"`
}

// GetUserCreditInvoices returns credit limit invoices for a specific user.
func (s *CreditService) GetUserCreditInvoices(userID uint, page, pageSize int, status string) ([]CreditInvoiceItem, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := s.db.Model(&model.CreditBill{}).Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var bills []model.CreditBill
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&bills).Error; err != nil {
		return nil, 0, err
	}

	var items []CreditInvoiceItem
	now := time.Now()
	for _, bill := range bills {
		paymentStatus := bill.Status
		dueDays := 0
		if bill.Status == "unpaid" && now.After(bill.DueDate) {
			paymentStatus = "overdue"
			dueDays = int(now.Sub(bill.DueDate).Hours() / 24)
		}

		paidDate := ""
		if bill.Status == "paid" {
			paidDate = bill.UpdatedAt.Format("2006-01-02 15:04:05")
		}

		var user model.User
		s.db.First(&user, bill.UserID)

		items = append(items, CreditInvoiceItem{
			ID:              bill.ID,
			UserID:          bill.UserID,
			Username:        user.Username,
			TotalAmount:     bill.TotalAmount,
			PaidAmount:      bill.PaidAmount,
			RemainingAmount: bill.RemainingAmount,
			LateFee:         bill.LateFee,
			Status:          bill.Status,
			PaymentStatus:   paymentStatus,
			DueDate:         bill.DueDate.Format("2006-01-02 15:04:05"),
			PaidDate:        paidDate,
			CreatedAt:       bill.CreatedAt.Format("2006-01-02 15:04:05"),
			BillMonth:       bill.BillMonth,
			DueDays:         dueDays,
		})
	}

	return items, total, nil
}

// GetCreditInvoiceSubItems returns sub-items under a credit invoice.
func (s *CreditService) GetCreditInvoiceSubItems(billID uint) ([]model.CreditBillItem, error) {
	var items []model.CreditBillItem
	if err := s.db.Where("bill_id = ?", billID).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ─────────────────────── Global Credit Config ───────────────────────

// GlobalCreditConfig represents the global credit limit configuration.
type GlobalCreditConfig struct {
	Enabled                  bool    `json:"enabled"`
	DefaultLimit             float64 `json:"default_limit"`
	BillGenerationDay        int     `json:"bill_generation_day"`
	RepaymentPeriod          int     `json:"repayment_period"`
	LateFeeEnabled           bool    `json:"late_fee_enabled"`
	LateFeeDailyPercent      float64 `json:"late_fee_daily_percent"`
}

// GetGlobalCreditConfig returns the global credit limit configuration.
func (s *CreditService) GetGlobalCreditConfig() (*GlobalCreditConfig, error) {
	config := &GlobalCreditConfig{
		Enabled:             true,
		DefaultLimit:        0,
		BillGenerationDay:   1,
		RepaymentPeriod:     15,
		LateFeeEnabled:      false,
		LateFeeDailyPercent: 0.05,
	}

	// Read from config_options table
	var options []model.ConfigOption
	if err := s.db.Where("`group` = ?", "credit_limit").Find(&options).Error; err != nil {
		return config, nil
	}

	for _, opt := range options {
		switch opt.Code {
		case "shd_credit_limit":
			config.Enabled = opt.Value == "1"
		case "shd_credit_limit_amount":
			config.DefaultLimit = parseFloat(opt.Value)
		case "shd_credit_limit_bill_generation_date":
			config.BillGenerationDay = parseInt(opt.Value)
		case "shd_credit_limit_bill_repayment_period":
			config.RepaymentPeriod = parseInt(opt.Value)
		case "shd_credit_limit_liquidated_damages":
			config.LateFeeEnabled = opt.Value == "1"
		case "shd_credit_limit_liquidated_damages_percent":
			config.LateFeeDailyPercent = parseFloat(opt.Value)
		}
	}

	return config, nil
}

// UpdateGlobalCreditConfig updates the global credit limit configuration.
func (s *CreditService) UpdateGlobalCreditConfig(config *GlobalCreditConfig) error {
	options := map[string]string{
		"shd_credit_limit":                          boolToIntStr(config.Enabled),
		"shd_credit_limit_amount":                   fmt.Sprintf("%.2f", config.DefaultLimit),
		"shd_credit_limit_bill_generation_date":     fmt.Sprintf("%d", config.BillGenerationDay),
		"shd_credit_limit_bill_repayment_period":    fmt.Sprintf("%d", config.RepaymentPeriod),
		"shd_credit_limit_liquidated_damages":       boolToIntStr(config.LateFeeEnabled),
		"shd_credit_limit_liquidated_damages_percent": fmt.Sprintf("%.2f", config.LateFeeDailyPercent),
	}

	for code, value := range options {
		var opt model.ConfigOption
		err := s.db.Where("code = ?", code).First(&opt).Error
		if err == gorm.ErrRecordNotFound {
			opt = model.ConfigOption{
				Group: "credit_limit",
				Code:  code,
				Name:  code,
				Type:  "text",
				Value: value,
			}
			s.db.Create(&opt)
		} else if err == nil {
			s.db.Model(&opt).Update("value", value)
		}
	}

	return nil
}

