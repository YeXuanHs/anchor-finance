package service

import (
	"errors"
	"time"

	"github.com/anchor-finance/backend/internal/util"
	"github.com/anchor-finance/backend/pkg/logger"

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
	Content string `json:"content" binding:"required"`
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

// GetByID fetches a ticket with replies.
func (s *TicketService) GetByID(id uint) (*Ticket, []TicketReply, error) {
	var ticket Ticket
	if err := s.db.First(&ticket, id).Error; err != nil {
		return nil, nil, err
	}

	var replies []TicketReply
	if err := s.db.Where("ticket_id = ?", id).Order("id ASC").Find(&replies).Error; err != nil {
		return nil, nil, err
	}

	return &ticket, replies, nil
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
