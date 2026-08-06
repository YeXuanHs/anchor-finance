package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserManageService provides admin operations on user accounts.
type UserManageService struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewUserManageService creates a new UserManageService.
func NewUserManageService(db *gorm.DB, log *logger.Logger) *UserManageService {
	return &UserManageService{db: db, log: log}
}

// ==================== Search & Filter ====================

// SearchClients searches users by keyword across multiple fields.
func (s *UserManageService) SearchClients(page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := s.db.Model(&model.User{})
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ? OR nickname LIKE ?", q, q, q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// FilterFilters defines the filter criteria for client queries.
type ClientFilter struct {
	Status    string `form:"status"`
	GroupID   uint   `form:"group_id"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Keyword   string `form:"keyword"`
}

// FilterClients filters clients by status, group, and date range.
func (s *UserManageService) FilterClients(page, pageSize int, f ClientFilter) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := s.db.Model(&model.User{})
	if f.Status != "" {
		query = query.Where("status = ?", f.Status)
	}
	if f.GroupID > 0 {
		query = query.Where("group_id = ?", f.GroupID)
	}
	if f.StartDate != "" {
		query = query.Where("created_at >= ?", f.StartDate)
	}
	if f.EndDate != "" {
		query = query.Where("created_at <= ?", f.EndDate)
	}
	if f.Keyword != "" {
		q := "%" + f.Keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ? OR nickname LIKE ?", q, q, q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// GetClientSummary returns aggregated client statistics.
func (s *UserManageService) GetClientSummary() (map[string]interface{}, error) {
	var total int64
	if err := s.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, err
	}

	var active int64
	s.db.Model(&model.User{}).Where("status = 1").Count(&active)

	var disabled int64
	s.db.Model(&model.User{}).Where("status = 0").Count(&disabled)

	var pending int64
	s.db.Model(&model.User{}).Where("status = 2").Count(&pending)

	var todayNew int64
	s.db.Model(&model.User{}).Where("DATE(created_at) = CURRENT_DATE").Count(&todayNew)

	var weekNew int64
	s.db.Model(&model.User{}).Where("created_at >= CURRENT_DATE - INTERVAL '7 days'").Count(&weekNew)

	return map[string]interface{}{
		"total":      total,
		"active":     active,
		"disabled":   disabled,
		"pending":    pending,
		"today_new":  todayNew,
		"week_new":   weekNew,
	}, nil
}

// ==================== Client Lifecycle ====================

// CreateClient creates a new client account.
func (s *UserManageService) CreateClient(req RegisterClientRequest) (*model.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		Phone:        req.Phone,
		PasswordHash: string(hashed),
		Nickname:     req.Nickname,
		Status:       1,
		GroupID:      1,
		Balance:      0,
	}

	if req.GroupID > 0 {
		user.GroupID = req.GroupID
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	s.log.Infof("admin created client: %s (id=%d)", user.Username, user.ID)
	return user, nil
}

// RegisterClientRequest is the payload for CreateClient.
type RegisterClientRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	Nickname string `json:"nickname"`
	GroupID  uint   `json:"group_id"`
}

// CloseClient suspends a client account (status=0).
func (s *UserManageService) CloseClient(clientID uint) error {
	var user model.User
	if err := s.db.First(&user, clientID).Error; err != nil {
		return fmt.Errorf("client not found: %w", err)
	}
	if user.Status == 0 {
		return errors.New("client is already closed")
	}
	return s.db.Model(&model.User{}).Where("id = ?", clientID).
		Update("status", 0).Error
}

// DeleteClient deletes a client. soft=true performs soft delete, soft=false performs hard delete.
func (s *UserManageService) DeleteClient(clientID uint, soft bool) error {
	var user model.User
	if err := s.db.First(&user, clientID).Error; err != nil {
		return fmt.Errorf("client not found: %w", err)
	}

	if soft {
		return s.db.Delete(&model.User{}, clientID).Error
	}

	return s.db.Unscoped().Delete(&model.User{}, clientID).Error
}

// BanClient disables a user account with a reason (status=0).
func (s *UserManageService) BanClient(clientID uint, reason string) error {
	var user model.User
	if err := s.db.First(&user, clientID).Error; err != nil {
		return fmt.Errorf("client not found: %w", err)
	}
	if user.Status == 0 {
		return errors.New("client is already banned")
	}

	updates := map[string]interface{}{
		"status": 0,
	}
	if reason != "" {
		updates["remark"] = reason
	}

	return s.db.Model(&model.User{}).Where("id = ?", clientID).
		Updates(updates).Error
}

// UnbanClient re-enables a user account (status=1).
func (s *UserManageService) UnbanClient(clientID uint) error {
	var user model.User
	if err := s.db.First(&user, clientID).Error; err != nil {
		return fmt.Errorf("client not found: %w", err)
	}
	if user.Status == 1 {
		return errors.New("client is already active")
	}
	return s.db.Model(&model.User{}).Where("id = ?", clientID).
		Update("status", 1).Error
}

// CancelBan cancels a ban request by re-enabling the client.
func (s *UserManageService) CancelBan(clientID uint) error {
	return s.UnbanClient(clientID)
}

// ==================== Client Profile ====================

// GetClientProfile returns a full client profile with statistics.
func (s *UserManageService) GetClientProfile(clientID uint) (map[string]interface{}, error) {
	var user model.User
	if err := s.db.First(&user, clientID).Error; err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	var serviceCount int64
	s.db.Model(&model.ClientService{}).Where("user_id = ?", clientID).Count(&serviceCount)

	var activeServiceCount int64
	s.db.Model(&model.ClientService{}).Where("user_id = ? AND status = ?", clientID, 1).Count(&activeServiceCount)

	var orderCount int64
	s.db.Model(&model.Order{}).Where("user_id = ?", clientID).Count(&orderCount)

	var invoiceCount int64
	s.db.Model(&model.Invoice{}).Where("user_id = ?", clientID).Count(&invoiceCount)

	var ticketCount int64
	s.db.Model(&model.Ticket{}).Where("user_id = ?", clientID).Count(&ticketCount)

	var totalSpent float64
	s.db.Model(&model.Invoice{}).Where("user_id = ? AND status = 1", clientID).
		Select("COALESCE(SUM(total), 0)").Scan(&totalSpent)

	var remarkCount int64
	s.db.Model(&model.UserRemark{}).Where("user_id = ?", clientID).Count(&remarkCount)

	var noteCount int64
	s.db.Model(&model.AdminNote{}).Where("user_id = ?", clientID).Count(&noteCount)

	return map[string]interface{}{
		"user":                  user,
		"service_count":         serviceCount,
		"active_service_count":  activeServiceCount,
		"order_count":           orderCount,
		"invoice_count":         invoiceCount,
		"ticket_count":          ticketCount,
		"total_spent":           totalSpent,
		"remark_count":          remarkCount,
		"note_count":            noteCount,
	}, nil
}

// UpdateClientProfile updates client profile fields.
func (s *UserManageService) UpdateClientProfile(clientID uint, data map[string]interface{}) error {
	var user model.User
	if err := s.db.First(&user, clientID).Error; err != nil {
		return fmt.Errorf("client not found: %w", err)
	}

	if len(data) == 0 {
		return errors.New("no data to update")
	}

	return s.db.Model(&model.User{}).Where("id = ?", clientID).Updates(data).Error
}

// GetClientHosts returns a client's host/server list.
func (s *UserManageService) GetClientHosts(clientID uint, page, pageSize int) ([]model.Host, int64, error) {
	var hosts []model.Host
	var total int64

	query := s.db.Model(&model.Host{}).Where("owner_id = ?", clientID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&hosts).Error; err != nil {
		return nil, 0, err
	}
	return hosts, total, nil
}

// GetClientInvoices returns a client's invoice list.
func (s *UserManageService) GetClientInvoices(clientID uint, page, pageSize int) ([]model.Invoice, int64, error) {
	var invoices []model.Invoice
	var total int64

	query := s.db.Model(&model.Invoice{}).Where("user_id = ?", clientID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	return invoices, total, nil
}

// GetClientOrders returns a client's order list.
func (s *UserManageService) GetClientOrders(clientID uint, page, pageSize int) ([]model.Order, int64, error) {
	var orders []model.Order
	var total int64

	query := s.db.Model(&model.Order{}).Where("user_id = ?", clientID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// GetClientTickets returns a client's ticket list.
func (s *UserManageService) GetClientTickets(clientID uint, page, pageSize int) ([]model.Ticket, int64, error) {
	var tickets []model.Ticket
	var total int64

	query := s.db.Model(&model.Ticket{}).Where("user_id = ?", clientID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&tickets).Error; err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}

// GetClientLogs returns client activity logs (balance logs).
func (s *UserManageService) GetClientLogs(clientID uint, page, pageSize int) ([]model.BalanceLog, int64, error) {
	var logs []model.BalanceLog
	var total int64

	query := s.db.Model(&model.BalanceLog{}).Where("user_id = ?", clientID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// ==================== Client Notes ====================

// AddNote adds an admin note for a client.
func (s *UserManageService) AddNote(clientID, adminID uint, content string, isPrivate bool) (*model.AdminNote, error) {
	var user model.User
	if err := s.db.First(&user, clientID).Error; err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	note := &model.AdminNote{
		UserID:    clientID,
		AdminID:   adminID,
		Content:   content,
		IsPrivate: isPrivate,
	}

	if err := s.db.Create(note).Error; err != nil {
		return nil, err
	}

	return note, nil
}

// GetNotes returns all admin notes for a client.
func (s *UserManageService) GetNotes(clientID uint) ([]model.AdminNote, error) {
	var notes []model.AdminNote
	if err := s.db.Preload("Admin").Where("user_id = ?", clientID).
		Order("created_at DESC").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// DeleteNote deletes an admin note by ID.
func (s *UserManageService) DeleteNote(noteID uint) error {
	result := s.db.Delete(&model.AdminNote{}, noteID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("note not found")
	}
	return nil
}

// ==================== Client Authorization ====================

// AuthorizeClient sets client permissions (stored in remark field as JSON-like info).
func (s *UserManageService) AuthorizeClient(clientID uint, permissions map[string]interface{}) error {
	var user model.User
	if err := s.db.First(&user, clientID).Error; err != nil {
		return fmt.Errorf("client not found: %w", err)
	}

	return s.db.Model(&model.User{}).Where("id = ?", clientID).
		Update("is_staff", true).Error
}

// GetClientAuth returns authorization settings for a client.
func (s *UserManageService) GetClientAuth(clientID uint) (map[string]interface{}, error) {
	var user model.User
	if err := s.db.First(&user, clientID).Error; err != nil {
		return nil, fmt.Errorf("client not found: %w", err)
	}

	return map[string]interface{}{
		"user_id":    user.ID,
		"is_admin":   user.IsAdmin,
		"is_staff":   user.IsStaff,
		"group_id":   user.GroupID,
	}, nil
}

// ==================== Client Group ====================

// AssignGroup assigns a client to a group.
func (s *UserManageService) AssignGroup(clientID, groupID uint) error {
	var user model.User
	if err := s.db.First(&user, clientID).Error; err != nil {
		return fmt.Errorf("client not found: %w", err)
	}

	var group model.ClientGroup
	if err := s.db.First(&group, groupID).Error; err != nil {
		return fmt.Errorf("group not found: %w", err)
	}

	// Update primary group
	if err := s.db.Model(&model.User{}).Where("id = ?", clientID).
		Update("group_id", groupID).Error; err != nil {
		return err
	}

	// Also add to group member table
	var count int64
	s.db.Model(&model.ClientGroupMember{}).
		Where("group_id = ? AND user_id = ?", groupID, clientID).Count(&count)
	if count == 0 {
		member := &model.ClientGroupMember{
			GroupID: groupID,
			UserID:  clientID,
		}
		return s.db.Create(member).Error
	}

	return nil
}

// RemoveFromGroup removes a client from a group.
func (s *UserManageService) RemoveFromGroup(clientID, groupID uint) error {
	result := s.db.Where("group_id = ? AND user_id = ?", groupID, clientID).
		Delete(&model.ClientGroupMember{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("client is not in this group")
	}
	return nil
}

// ==================== Certification ====================

// GetCertificationStatus returns the real-name verification status for a client.
func (s *UserManageService) GetCertificationStatus(clientID uint) (*model.Certification, error) {
	var cert model.Certification
	err := s.db.Where("user_id = ?", clientID).First(&cert).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// ReviewCertification approves or rejects a certification submission.
func (s *UserManageService) ReviewCertification(clientID uint, approved bool, reason string, reviewerID uint) error {
	var cert model.Certification
	if err := s.db.Where("user_id = ?", clientID).First(&cert).Error; err != nil {
		return fmt.Errorf("certification not found: %w", err)
	}

	if cert.Status != 1 {
		return errors.New("certification is not pending review")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"reviewed_by": reviewerID,
		"reviewed_at": now,
	}

	if approved {
		updates["status"] = 2
	} else {
		updates["status"] = 3
		updates["reject_reason"] = reason
	}

	return s.db.Model(&model.Certification{}).Where("id = ?", cert.ID).
		Updates(updates).Error
}

// ==================== Cancel Requests ====================

// GetCancelRequests returns all pending cancellation requests.
func (s *UserManageService) GetCancelRequests(page, pageSize int, status string) ([]model.CancelRequest, int64, error) {
	var requests []model.CancelRequest
	var total int64

	query := s.db.Model(&model.CancelRequest{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("User").Offset(offset).Limit(pageSize).
		Order("id DESC").Find(&requests).Error; err != nil {
		return nil, 0, err
	}
	return requests, total, nil
}

// ProcessCancelRequest approves or rejects a cancellation request.
func (s *UserManageService) ProcessCancelRequest(requestID uint, approved bool, adminID uint, remark string) error {
	var req model.CancelRequest
	if err := s.db.First(&req, requestID).Error; err != nil {
		return fmt.Errorf("cancel request not found: %w", err)
	}

	if req.Status != "pending" {
		return errors.New("request has already been processed")
	}

	updates := map[string]interface{}{
		"admin_id": adminID,
		"remark":   remark,
	}

	if approved {
		updates["status"] = "approved"
	} else {
		updates["status"] = "rejected"
	}

	return s.db.Model(&model.CancelRequest{}).Where("id = ?", requestID).
		Updates(updates).Error
}

// ==================== Legacy aliases (backward compat) ====================

// SearchUsers is an alias for SearchClients.
func (s *UserManageService) SearchUsers(page, pageSize int, keyword string) ([]model.User, int64, error) {
	return s.SearchClients(page, pageSize, keyword)
}

// Ban disables a user account (status=0). Alias for BanClient.
func (s *UserManageService) Ban(userID uint, reason string) error {
	return s.BanClient(userID, reason)
}

// Unban re-enables a user account (status=1). Alias for UnbanClient.
func (s *UserManageService) Unban(userID uint) error {
	return s.UnbanClient(userID)
}

// AdjustBalance adds or deducts balance for a user and logs the transaction.
func (s *UserManageService) AdjustBalance(userID uint, amount float64, description string) error {
	if amount == 0 {
		return errors.New("amount must not be zero")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var currentBalance float64
		if err := tx.Table("users").Where("id = ?", userID).Select("balance").Scan(&currentBalance).Error; err != nil {
			return fmt.Errorf("user not found: %w", err)
		}

		if amount < 0 && currentBalance < -amount {
			return fmt.Errorf("insufficient balance: have %.4f, need %.4f", currentBalance, -amount)
		}

		if err := tx.Table("users").Where("id = ?", userID).
			Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}

		var newBalance float64
		tx.Table("users").Where("id = ?", userID).Select("balance").Scan(&newBalance)

		balLog := &model.BalanceLog{
			UserID:      userID,
			Amount:      amount,
			Balance:     newBalance,
			RelatedType: "admin",
			Description: description,
		}
		return tx.Create(balLog).Error
	})
}

// ResetPassword sets a new password for a user by admin.
func (s *UserManageService) ResetPassword(userID uint, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.Model(&model.User{}).Where("id = ?", userID).
		Update("password_hash", string(hashed)).Error
}

// GetOperationLogs returns the operation log entries for a user.
func (s *UserManageService) GetOperationLogs(userID uint, page, pageSize int) ([]model.BalanceLog, int64, error) {
	return s.GetClientLogs(userID, page, pageSize)
}

// GetUserStatus returns user status information.
func (s *UserManageService) GetUserStatus(userID uint) (map[string]interface{}, error) {
	return s.GetClientProfile(userID)
}

// ==================== P1-7: GetBlackList ====================

// GetBlackList 获取黑名单用户列表
func (s *UserManageService) GetBlackList(page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := s.db.Model(&model.User{}).Where("status = ?", -1)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

// RemoveBlackList 从黑名单移除用户
func (s *UserManageService) RemoveBlackList(userID uint) error {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}
	if user.Status != -1 {
		return errors.New("user is not in blacklist")
	}
	return s.db.Model(&model.User{}).Where("id = ?", userID).Update("status", 1).Error
}

// ==================== P3-21: GetUserInvoices ====================

// GetUserInvoices 获取用户发票列表（独立页面）
func (s *UserManageService) GetUserInvoices(userID uint, page, pageSize int, status string) ([]map[string]interface{}, int64, error) {
	var invoices []map[string]interface{}
	var total int64

	query := s.db.Table("invoices").Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&invoices).Error; err != nil {
		return nil, 0, err
	}
	return invoices, total, nil
}

// ==================== Cancel Reasons ====================

// CancelReason represents a reason for account cancellation.
type CancelReason struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Reason    string    `gorm:"type:varchar(256);not null" json:"reason"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	Status    int8      `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// GetCancelReasons returns all active cancel reasons.
func (s *UserManageService) GetCancelReasons() ([]CancelReason, error) {
	var reasons []CancelReason
	if err := s.db.Where("status = 1").Order("sort_order ASC").Find(&reasons).Error; err != nil {
		return nil, err
	}
	return reasons, nil
}

// AddCancelReason adds a new cancel reason.
func (s *UserManageService) AddCancelReason(reason string, sortOrder int) error {
	r := &CancelReason{
		Reason:    reason,
		SortOrder: sortOrder,
		Status:    1,
	}
	return s.db.Create(r).Error
}

// DeleteCancelReason deletes a cancel reason by ID.
func (s *UserManageService) DeleteCancelReason(id uint) error {
	result := s.db.Delete(&CancelReason{}, id)
	if result.RowsAffected == 0 {
		return errors.New("reason not found")
	}
	return result.Error
}

// ==================== Invoice Creation ====================

// CreateRechargeInvoice creates a recharge invoice for a user.
func (s *UserManageService) CreateRechargeInvoice(userID uint, amount float64, description string) (*model.Invoice, error) {
	invoice := &model.Invoice{
		UserID:      userID,
		Total:       amount,
		Status:      0,
		Description: description,
	}
	if err := s.db.Create(invoice).Error; err != nil {
		return nil, err
	}
	return invoice, nil
}

// CreateUserInvoice creates a general invoice for a user.
func (s *UserManageService) CreateUserInvoice(userID uint, invoiceType string, items []map[string]interface{}) (*model.Invoice, error) {
	invoice := &model.Invoice{
		UserID: userID,
		Status: 0,
		Type:   invoiceType,
	}
	if err := s.db.Create(invoice).Error; err != nil {
		return nil, err
	}
	return invoice, nil
}

// ==================== Certification File ====================

// GetCertificationFile returns the certification file info for a client.
func (s *UserManageService) GetCertificationFile(clientID uint) (map[string]interface{}, error) {
	var cert model.Certification
	err := s.db.Where("user_id = ?", clientID).First(&cert).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"id":             cert.ID,
		"user_id":        cert.UserID,
		"type":           cert.Type,
		"real_name":      cert.RealName,
		"id_number":      cert.IDNumber,
		"front_image":    cert.FrontImage,
		"back_image":     cert.BackImage,
		"handheld_image": cert.HandheldImage,
		"status":         cert.Status,
		"reject_reason":  cert.RejectReason,
		"reviewed_by":    cert.ReviewedBy,
		"reviewed_at":    cert.ReviewedAt,
		"created_at":     cert.CreatedAt,
	}, nil
}
