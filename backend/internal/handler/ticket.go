package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketSvc *service.TicketService
	log       *logger.Logger
	uploadDir string
}

func NewTicketHandler(ticketSvc *service.TicketService, log *logger.Logger) *TicketHandler {
	return &TicketHandler{ticketSvc: ticketSvc, log: log, uploadDir: "uploads/tickets"}
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

// UploadAttachment handles file upload for a ticket or reply.
func (h *TicketHandler) UploadAttachment(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ticket id")
		return
	}

	userID := c.GetUint("user_id")

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	// Limit file size to 20MB
	if file.Size > 20*1024*1024 {
		response.BadRequest(c, "file size exceeds 20MB limit")
		return
	}

	// Ensure upload directory exists
	destDir := filepath.Join(h.uploadDir, fmt.Sprintf("%d", ticketID))
	if err := os.MkdirAll(destDir, 0755); err != nil {
		response.ServerError(c, "failed to create upload directory")
		return
	}

	// Generate unique filename to avoid collision
	ext := filepath.Ext(file.Filename)
	savedName := fmt.Sprintf("%d_%d%s", userID, ticketID, ext)
	savedPath := filepath.Join(destDir, savedName)

	if err := c.SaveUploadedFile(file, savedPath); err != nil {
		response.ServerError(c, "failed to save file")
		return
	}

	// Compute SHA256 hash
	var hash string
	f, err := os.Open(savedPath)
	if err == nil {
		h2 := sha256.New()
		io.Copy(h2, f)
		f.Close()
		hash = hex.EncodeToString(h2.Sum(nil))
	}

	// Parse optional reply_id
	var replyID *uint
	if rid := c.PostForm("reply_id"); rid != "" {
		v, _ := strconv.ParseUint(rid, 10, 64)
		if v > 0 {
			u := uint(v)
			replyID = &u
		}
	}

	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	att, err := h.ticketSvc.UploadAttachment(
		uint(ticketID), replyID, file.Filename, savedPath, file.Size, mimeType, hash, userID,
	)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, att)
}

// GetAttachments lists all attachments for a ticket.
func (h *TicketHandler) GetAttachments(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ticket id")
		return
	}

	attachments, err := h.ticketSvc.GetAttachments(uint(ticketID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, attachments)
}

// DeleteAttachment deletes an attachment (verifies ownership).
func (h *TicketHandler) DeleteAttachment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid attachment id")
		return
	}

	userID := c.GetUint("user_id")

	if err := h.ticketSvc.DeleteAttachment(uint(id), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "attachment deleted")
}

// MergeTickets merges multiple tickets into one (admin).
func (h *TicketHandler) MergeTickets(c *gin.Context) {
	adminID := c.GetUint("admin_id")

	var req struct {
		SourceIDs []uint `json:"source_ids" binding:"required,min=1"`
		TargetID  uint   `json:"target_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.ticketSvc.MergeTickets(req.SourceIDs, req.TargetID, adminID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "tickets merged")
}

// TransferTicket transfers a ticket to another department/agent (admin).
func (h *TicketHandler) TransferTicket(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ticket id")
		return
	}

	adminID := c.GetUint("admin_id")

	var req struct {
		DeptID  *uint  `json:"dept_id"`
		AgentID *uint  `json:"agent_id"`
		Reason  string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.DeptID == nil && req.AgentID == nil {
		response.BadRequest(c, "at least one of dept_id or agent_id is required")
		return
	}

	if err := h.ticketSvc.TransferTicket(uint(ticketID), req.DeptID, req.AgentID, adminID, req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "ticket transferred")
}

// GetTransferLogs returns the transfer history for a ticket (admin).
func (h *TicketHandler) GetTransferLogs(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ticket id")
		return
	}

	logs, err := h.ticketSvc.GetTransferLogs(uint(ticketID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, logs)
}
