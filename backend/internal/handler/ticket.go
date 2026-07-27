package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketSvc *service.TicketService
	log       *logger.Logger
}

func NewTicketHandler(ticketSvc *service.TicketService, log *logger.Logger) *TicketHandler {
	return &TicketHandler{ticketSvc: ticketSvc, log: log}
}

// Create opens a new ticket.
func (h *TicketHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req service.CreateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ticket, err := h.ticketSvc.Create(userID, req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, ticket)
}

// GetDetail returns a single ticket with replies.
func (h *TicketHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ticket id")
		return
	}

	ticket, replies, err := h.ticketSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "ticket not found")
		return
	}
	response.Success(c, gin.H{
		"ticket":  ticket,
		"replies": replies,
	})
}

// GetUserTickets returns paginated tickets for the authenticated user.
func (h *TicketHandler) GetUserTickets(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	tickets, total, err := h.ticketSvc.GetUserTickets(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, tickets, total, page, pageSize)
}

// Reply adds a user reply to a ticket.
func (h *TicketHandler) Reply(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ticket id")
		return
	}

	userID := c.GetUint("user_id")
	var req service.ReplyTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	reply, err := h.ticketSvc.Reply(uint(id), userID, false, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, reply)
}

// Close closes a ticket.
func (h *TicketHandler) Close(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ticket id")
		return
	}

	if err := h.ticketSvc.Close(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "ticket closed")
}

// GetList returns all tickets (admin).
func (h *TicketHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}

	tickets, total, err := h.ticketSvc.GetList(page, pageSize, status, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, tickets, total, page, pageSize)
}

// Assign assigns a ticket to an admin (admin).
func (h *TicketHandler) Assign(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ticket id")
		return
	}

	var req struct {
		AssigneeID uint `json:"assignee_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.ticketSvc.Assign(uint(id), req.AssigneeID); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "ticket assigned")
}
