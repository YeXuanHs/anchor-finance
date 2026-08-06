package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type Ticket struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	TicketNo   string         `gorm:"uniqueIndex;size:64;not null" json:"ticket_no"`
	UserID     uint           `gorm:"index;not null" json:"user_id"`
	Subject    string         `gorm:"size:256;not null" json:"subject"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	Priority   string         `gorm:"size:16;default:normal;comment:low/normal/high/urgent" json:"priority"`
	Status     int            `gorm:"default:0;comment:0=open 1=replied 2=closed" json:"status"`
	AssigneeID *uint          `gorm:"index" json:"assignee_id"`
	ClosedAt   *time.Time     `json:"closed_at"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type TicketReply struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TicketID  uint      `gorm:"index;not null" json:"ticket_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	IsAdmin   bool      `gorm:"default:false" json:"is_admin"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type TicketAttachment struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	TicketID     uint   `gorm:"index;not null" json:"ticket_id"`
	ReplyID      *uint  `gorm:"index" json:"reply_id"`
	FileName     string `gorm:"type:varchar(255);not null" json:"file_name"`
	FilePath     string `gorm:"type:varchar(512);not null" json:"file_path"`
	FileSize     int64  `gorm:"not null" json:"file_size"`
	MimeType     string `gorm:"type:varchar(128)" json:"mime_type"`
	UploaderID   uint   `gorm:"index;not null" json:"uploader_id"`
	StorageDriver string `gorm:"type:varchar(32);default:'local'" json:"storage_driver"`
	StorageKey   string `gorm:"type:varchar(512)" json:"storage_key"`
	Hash         string `gorm:"type:varchar(64);index" json:"hash"`
}

func (TicketAttachment) TableName() string { return "attachments" }

type TicketTransferLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TicketID   uint      `gorm:"index;not null" json:"ticket_id"`
	FromDeptID *uint     `json:"from_dept_id"`
	ToDeptID   *uint     `json:"to_dept_id"`
	FromAgentID *uint    `json:"from_agent_id"`
	ToAgentID  *uint     `json:"to_agent_id"`
	OperatorID uint      `gorm:"not null" json:"operator_id"`
	Reason     string    `gorm:"type:text" json:"reason"`
	CreatedAt  time.Time `json:"created_at"`
}

type TicketService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewTicketService(db *gorm.DB, log *logger.Logger) *TicketService {
	return &TicketService{db: db, log: log}
}

type CreateTicketRequest struct {
	Subject  string `json:"subject" binding:"required,max=256"`
	Content  string `json:"content" binding:"required"`
	Priority string `json:"priority" binding:"omitempty,oneof=low normal high urgent"`
}

type ReplyTicketRequest struct {
	Content       string `json:"content" binding:"required"`
	AttachmentIDs []uint `json:"attachment_ids"`
}

// Create opens a new ticket.
func (s *TicketService) Create(userID uint, req CreateTicketRequest) (*Ticket, error) {
	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	ticket := &Ticket{
		TicketNo: util.GenerateTicketNo(),
		UserID:   userID,
		Subject:  req.Subject,
		Content:  req.Content,
		Priority: priority,
		Status:   0,
	}

	if err := s.db.Create(ticket).Error; err != nil {
		return nil, err
	}

	s.log.Infof("ticket created: %s (user=%d, subject=%s)", ticket.TicketNo, userID, req.Subject)
	return ticket, nil
}

// GetByID fetches a ticket with replies and attachments.
func (s *TicketService) GetByID(id uint) (*Ticket, []TicketReply, []TicketAttachment, error) {
	var ticket Ticket
	if err := s.db.First(&ticket, id).Error; err != nil {
		return nil, nil, nil, err
	}

	var replies []TicketReply
	if err := s.db.Where("ticket_id = ?", id).Order("id ASC").Find(&replies).Error; err != nil {
		return nil, nil, nil, err
	}

	// 获取工单附件（不属于任何回复的）
	var attachments []TicketAttachment
	if err := s.db.Where("ticket_id = ? AND reply_id IS NULL", id).Order("id ASC").Find(&attachments).Error; err != nil {
		return nil, nil, nil, err
	}

	// 获取各回复的附件
	replyAttachments := make(map[uint][]TicketAttachment)
	var allReplyAttachments []TicketAttachment
	replyIDs := make([]uint, 0, len(replies))
	for _, r := range replies {
		replyIDs = append(replyIDs, r.ID)
	}
	if len(replyIDs) > 0 {
		if err := s.db.Where("ticket_id = ? AND reply_id IN ?", id, replyIDs).Order("id ASC").Find(&allReplyAttachments).Error; err != nil {
			return nil, nil, nil, err
		}
		for _, att := range allReplyAttachments {
			if att.ReplyID != nil {
				replyAttachments[*att.ReplyID] = append(replyAttachments[*att.ReplyID], att)
			}
		}
	}

	// 构建带附件的回复列表
	type ReplyWithAttachments struct {
		TicketReply
		Attachments []TicketAttachment `json:"attachments"`
	}

	return &ticket, replies, attachments, nil
}

// GetUserTickets returns paginated tickets for a user.
func (s *TicketService) GetUserTickets(userID uint, page, pageSize int) ([]Ticket, int64, error) {
	var tickets []Ticket
	var total int64

	query := s.db.Model(&Ticket{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&tickets).Error; err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}

// Reply adds a reply to a ticket (user or admin).
func (s *TicketService) Reply(ticketID, userID uint, isAdmin bool, req ReplyTicketRequest) (*TicketReply, error) {
	var ticket Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		return nil, err
	}
	if ticket.Status == 2 {
		return nil, errors.New("ticket is closed")
	}

	reply := &TicketReply{
		TicketID: ticketID,
		UserID:   userID,
		IsAdmin:  isAdmin,
		Content:  req.Content,
	}

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(reply).Error; err != nil {
			return err
		}

		// 处理附件关联
		if len(req.AttachmentIDs) > 0 {
			if err := tx.Model(&model.Attachment{}).
				Where("id IN ? AND ticket_id = ?", req.AttachmentIDs, ticketID).
				Update("reply_id", reply.ID).Error; err != nil {
				return err
			}
		}

		newStatus := 1
		if isAdmin {
			newStatus = 1 // replied
		}
		return tx.Model(&ticket).Update("status", newStatus).Error
	})

	if err != nil {
		return nil, err
	}

	return reply, nil
}

// Close marks a ticket as closed.
func (s *TicketService) Close(ticketID uint) error {
	var ticket Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		return err
	}
	if ticket.Status == 2 {
		return errors.New("ticket already closed")
	}

	now := time.Now()
	return s.db.Model(&ticket).Updates(map[string]interface{}{
		"status":    2,
		"closed_at": &now,
	}).Error
}

// Assign sets an admin assignee for a ticket.
func (s *TicketService) Assign(ticketID, assigneeID uint) error {
	var ticket Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		return err
	}
	return s.db.Model(&ticket).Update("assignee_id", assigneeID).Error
}

// GetList returns all tickets with pagination (admin).
func (s *TicketService) GetList(page, pageSize int, status *int, keyword string) ([]Ticket, int64, error) {
	var tickets []Ticket
	var total int64

	query := s.db.Model(&Ticket{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("subject LIKE ? OR ticket_no LIKE ?", q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&tickets).Error; err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}

// UploadAttachment saves an attachment record for a ticket or reply.
func (s *TicketService) UploadAttachment(ticketID uint, replyID *uint, fileName, filePath string, fileSize int64, mimeType, hash string, uploaderID uint) (*TicketAttachment, error) {
	var ticket Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		return nil, errors.New("ticket not found")
	}

	att := &TicketAttachment{
		TicketID:   ticketID,
		ReplyID:    replyID,
		FileName:   fileName,
		FilePath:   filePath,
		FileSize:   fileSize,
		MimeType:   mimeType,
		UploaderID: uploaderID,
		Hash:       hash,
	}

	if err := s.db.Create(att).Error; err != nil {
		return nil, err
	}

	s.log.Infof("attachment uploaded: ticket=%d file=%s size=%d", ticketID, fileName, fileSize)
	return att, nil
}

// GetAttachments returns all attachments for a ticket.
func (s *TicketService) GetAttachments(ticketID uint) ([]TicketAttachment, error) {
	var attachments []TicketAttachment
	if err := s.db.Where("ticket_id = ?", ticketID).Order("id ASC").Find(&attachments).Error; err != nil {
		return nil, err
	}
	return attachments, nil
}

// DeleteAttachment deletes an attachment after verifying uploader ownership.
func (s *TicketService) DeleteAttachment(id, userID uint) error {
	var att TicketAttachment
	if err := s.db.First(&att, id).Error; err != nil {
		return errors.New("attachment not found")
	}
	if att.UploaderID != userID {
		return errors.New("permission denied")
	}
	return s.db.Delete(&att).Error
}

// MergeTickets merges multiple source tickets into a target ticket.
// All replies from source tickets are moved to the target; sources are marked as merged.
func (s *TicketService) MergeTickets(sourceIDs []uint, targetID uint, operatorID uint) error {
	if len(sourceIDs) == 0 {
		return errors.New("no source tickets provided")
	}

	var target Ticket
	if err := s.db.First(&target, targetID).Error; err != nil {
		return errors.New("target ticket not found")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, srcID := range sourceIDs {
			if srcID == targetID {
				continue
			}

			var src Ticket
			if err := tx.First(&src, srcID).Error; err != nil {
				return fmt.Errorf("source ticket %d not found", srcID)
			}

			// Move replies to target ticket
			if err := tx.Model(&TicketReply{}).Where("ticket_id = ?", srcID).Update("ticket_id", targetID).Error; err != nil {
				return fmt.Errorf("failed to move replies from ticket %d: %w", srcID, err)
			}

			// Move attachments to target ticket
			if err := tx.Model(&TicketAttachment{}).Where("ticket_id = ?", srcID).Update("ticket_id", targetID).Error; err != nil {
				return fmt.Errorf("failed to move attachments from ticket %d: %w", srcID, err)
			}

			// Mark source as merged (status=5 = cancelled, merged_into=targetID)
			if err := tx.Model(&src).Updates(map[string]interface{}{
				"status":      5,
				"merged_into": targetID,
			}).Error; err != nil {
				return fmt.Errorf("failed to mark ticket %d as merged: %w", srcID, err)
			}
		}

		s.log.Infof("tickets merged into %d by operator %d, sources: %v", targetID, operatorID, sourceIDs)
		return nil
	})
}

// TransferTicket transfers a ticket to another department/agent and logs the transfer.
func (s *TicketService) TransferTicket(ticketID uint, toDeptID *uint, toAgentID *uint, operatorID uint, reason string) error {
	var ticket Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		return errors.New("ticket not found")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{}
		if toDeptID != nil {
			updates["department_id"] = *toDeptID
		}
		if toAgentID != nil {
			updates["assignee_id"] = *toAgentID
		}
		if len(updates) > 0 {
			if err := tx.Model(&ticket).Updates(updates).Error; err != nil {
				return err
			}
		}

		log := &TicketTransferLog{
			TicketID:   ticketID,
			ToDeptID:   toDeptID,
			ToAgentID:  toAgentID,
			OperatorID: operatorID,
			Reason:     reason,
		}
		if err := tx.Create(log).Error; err != nil {
			return err
		}

		s.log.Infof("ticket %d transferred by operator %d", ticketID, operatorID)
		return nil
	})
}

// GetTransferLogs returns the transfer history for a ticket.
func (s *TicketService) GetTransferLogs(ticketID uint) ([]TicketTransferLog, error) {
	var logs []TicketTransferLog
	if err := s.db.Where("ticket_id = ?", ticketID).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

// TicketNote is the service-level struct for ticket notes.
type TicketNote struct {
	ID        uint      `json:"id"`
	TicketID  uint      `json:"ticket_id"`
	AdminID   uint      `json:"admin_id"`
	AdminName string    `json:"admin_name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// AddNote adds an admin note to a ticket.
func (s *TicketService) AddNote(ticketID, adminID uint, content string) (*TicketNote, error) {
	var ticket struct{ ID uint }
	if err := s.db.Table("tickets").Select("id").Where("id = ?", ticketID).First(&ticket).Error; err != nil {
		return nil, errors.New("ticket not found")
	}

	note := map[string]interface{}{
		"ticket_id":  ticketID,
		"admin_id":   adminID,
		"content":    content,
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}
	if err := s.db.Table("ticket_notes").Create(&note).Error; err != nil {
		return nil, err
	}

	var adminName string
	s.db.Table("admins").Select("username").Where("id = ?", adminID).Scan(&adminName)

	return &TicketNote{
		TicketID:  ticketID,
		AdminID:   adminID,
		AdminName: adminName,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}

// DeleteNote deletes a ticket note.
func (s *TicketService) DeleteNote(noteID uint) error {
	result := s.db.Table("ticket_notes").Where("id = ?", noteID).Delete(nil)
	if result.RowsAffected == 0 {
		return errors.New("note not found")
	}
	return result.Error
}

// DeleteReply deletes a ticket reply.
func (s *TicketService) DeleteReply(replyID uint) error {
	result := s.db.Table("ticket_replies").Where("id = ?", replyID).Delete(nil)
	if result.RowsAffected == 0 {
		return errors.New("reply not found")
	}
	return result.Error
}

// HostInfo represents a host associated with a ticket user.
type HostInfo struct {
	ID            uint    `json:"id"`
	UserID        uint    `json:"user_id"`
	ProductID     uint    `json:"product_id"`
	ProductName   string  `json:"product_name"`
	Domain        string  `json:"domain"`
	Status        string  `json:"status"`
	BillingCycle  string  `json:"billing_cycle"`
	Amount        float64 `json:"amount"`
	NextDueDate   string  `json:"next_due_date"`
	CreatedAt     string  `json:"created_at"`
}

// GetTicketDetailHost returns hosts associated with a ticket's user.
func (s *TicketService) GetTicketDetailHost(ticketID uint, page, pageSize int) ([]HostInfo, int64, error) {
	var ticket struct{ UserID uint }
	if err := s.db.Table("tickets").Select("user_id").Where("id = ?", ticketID).First(&ticket).Error; err != nil {
		return nil, 0, errors.New("ticket not found")
	}

	var total int64
	s.db.Table("hosts").Where("user_id = ?", ticket.UserID).Count(&total)

	var hosts []HostInfo
	offset := (page - 1) * pageSize
	s.db.Table("hosts h").
		Select("h.id, h.user_id, h.product_id, p.name as product_name, h.domain, h.status, h.billing_cycle, h.amount, h.next_due_date, h.created_at").
		Joins("LEFT JOIN products p ON p.id = h.product_id").
		Where("h.user_id = ?", ticket.UserID).
		Order("h.id DESC").
		Offset(offset).Limit(pageSize).
		Find(&hosts)

	return hosts, total, nil
}

// TicketStatistics represents ticket statistics.
type TicketStatistics struct {
	Total       int64            `json:"total"`
	Open        int64            `json:"open"`
	Replied     int64            `json:"replied"`
	Closed      int64            `json:"closed"`
	Pending     int64            `json:"pending"`
	TodayOpen   int64            `json:"today_open"`
	TodayClosed int64            `json:"today_closed"`
	WeekOpen    int64            `json:"week_open"`
	MonthOpen   int64            `json:"month_open"`
	ByDept      []DeptStat       `json:"by_dept"`
	ByPriority  []PriorityStat   `json:"by_priority"`
}

type DeptStat struct {
	DeptID   uint   `json:"dept_id"`
	DeptName string `json:"dept_name"`
	Count    int64  `json:"count"`
}

type PriorityStat struct {
	Priority string `json:"priority"`
	Count    int64  `json:"count"`
}

// GetTicketStatistics returns ticket statistics.
func (s *TicketService) GetTicketStatistics() (*TicketStatistics, error) {
	stats := &TicketStatistics{}

	// Total count
	s.db.Table("tickets").Count(&stats.Total)

	// Status counts
	s.db.Table("tickets").Where("status = 0").Count(&stats.Open)
	s.db.Table("tickets").Where("status = 1").Count(&stats.Replied)
	s.db.Table("tickets").Where("status = 2").Count(&stats.Closed)
	s.db.Table("tickets").Where("status = 3").Count(&stats.Pending)

	// Today counts
	today := time.Now().Format("2006-01-02")
	s.db.Table("tickets").Where("DATE(created_at) = ?", today).Count(&stats.TodayOpen)
	s.db.Table("tickets").Where("status = 2 AND DATE(closed_at) = ?", today).Count(&stats.TodayClosed)

	// Week and month counts
	weekAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	monthStart := time.Now().Format("2006-01") + "-01"
	s.db.Table("tickets").Where("created_at >= ?", weekAgo).Count(&stats.WeekOpen)
	s.db.Table("tickets").Where("created_at >= ?", monthStart).Count(&stats.MonthOpen)

	// By department
	s.db.Table("tickets t").
		Select("t.department_id as dept_id, d.name as dept_name, COUNT(*) as count").
		Joins("LEFT JOIN departments d ON d.id = t.department_id").
		Group("t.department_id, d.name").
		Find(&stats.ByDept)

	// By priority
	s.db.Table("tickets").
		Select("priority, COUNT(*) as count").
		Group("priority").
		Find(&stats.ByPriority)

	return stats, nil
}

// ==================== P1-8: GetAttachmentByID ====================

// GetAttachmentByID 根据ID获取附件信息
func (s *TicketService) GetAttachmentByID(id uint) (*TicketAttachment, error) {
	var att TicketAttachment
	if err := s.db.First(&att, id).Error; err != nil {
		return nil, err
	}
	return &att, nil
}

// ==================== P1-9: ReceiveTicket ====================

// ReceiveTicket 工单接单/领取
func (s *TicketService) ReceiveTicket(ticketID, adminID uint) error {
	var ticket Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		return errors.New("工单不存在")
	}
	if ticket.AssigneeID != nil && *ticket.AssigneeID > 0 {
		return errors.New("工单已被其他人领取")
	}

	return s.db.Model(&ticket).Updates(map[string]interface{}{
		"assignee_id": adminID,
		"status":      1, // replied/in-progress
	}).Error
}

// ==================== P2-15: GetListEnhanced ====================

// GetListEnhanced 工单列表（支持部门权限过滤和高级筛选）
func (s *TicketService) GetListEnhanced(page, pageSize int, status *int, keyword string,
	deptID, assigneeID, userID *uint, priority, startTime, endTime string) ([]Ticket, int64, error) {
	var tickets []Ticket
	var total int64

	query := s.db.Model(&Ticket{})

	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("subject LIKE ? OR ticket_no LIKE ?", like, like)
	}
	if deptID != nil {
		query = query.Where("department_id = ?", *deptID)
	}
	if assigneeID != nil {
		query = query.Where("assignee_id = ?", *assigneeID)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime+" 23:59:59")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&tickets).Error; err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}
