package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"anchorfinance/internal/service"
	"anchorfinance/internal/validator"
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

	// 验证工单主题（移植自 zjmf：不能为空，最多200字符）
	if ok, msg := validator.ValidateTicketSubject(req.Subject); !ok {
		response.BadRequest(c, msg)
		return
	}

	// 验证工单内容（移植自 zjmf：不能为空）
	if ok, msg := validator.ValidateTicketContent(req.Content); !ok {
		response.BadRequest(c, msg)
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

	ticket, replies, attachments, err := h.ticketSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "ticket not found")
		return
	}

	// 获取每个回复的附件
	replyAttachments := make(map[uint][]service.TicketAttachment)
	for _, att := range attachments {
		if att.ReplyID != nil {
			replyAttachments[*att.ReplyID] = append(replyAttachments[*att.ReplyID], att)
		}
	}

	// 构建带附件的回复列表
	type ReplyWithAttachments struct {
		service.TicketReply
		Attachments []service.TicketAttachment `json:"attachments"`
	}

	repliesWithAttachments := make([]ReplyWithAttachments, 0, len(replies))
	for _, reply := range replies {
		atts := replyAttachments[reply.ID]
		if atts == nil {
			atts = []service.TicketAttachment{}
		}
		repliesWithAttachments = append(repliesWithAttachments, ReplyWithAttachments{
			TicketReply: reply,
			Attachments: atts,
		})
	}

	// 过滤出工单级别的附件（不属于回复的）
	ticketAttachments := make([]service.TicketAttachment, 0)
	for _, att := range attachments {
		if att.ReplyID == nil {
			ticketAttachments = append(ticketAttachments, att)
		}
	}

	response.Success(c, gin.H{
		"ticket":      ticket,
		"replies":     repliesWithAttachments,
		"attachments": ticketAttachments,
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

	// 验证回复内容（移植自 zjmf：不能为空）
	if ok, msg := validator.ValidateTicketContent(req.Content); !ok {
		response.BadRequest(c, msg)
		return
	}

	reply, err := h.ticketSvc.Reply(uint(id), userID, false, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, reply)
}

// AdminReply adds an admin reply to a ticket.
func (h *TicketHandler) AdminReply(c *gin.Context) {
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

	// 验证回复内容（移植自 zjmf：不能为空）
	if ok, msg := validator.ValidateTicketContent(req.Content); !ok {
		response.BadRequest(c, msg)
		return
	}

	reply, err := h.ticketSvc.Reply(uint(id), userID, true, req)
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

	// 文件类型白名单（禁止可执行文件和脚本）
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".txt": true, ".csv": true, ".zip": true, ".rar": true, ".7z": true,
		".mp4": true, ".mp3": true, ".wav": true,
	}
	dangerousExts := map[string]bool{
		".exe": true, ".bat": true, ".cmd": true, ".sh": true, ".ps1": true,
		".php": true, ".asp": true, ".aspx": true, ".jsp": true, ".py": true,
		".rb": true, ".pl": true, ".cgi": true, ".htaccess": true,
		".js": true, ".vbs": true, ".wsf": true, ".scr": true, ".com": true,
		".pif": true, ".msi": true, ".dll": true, ".so": true, ".dylib": true,
	}
	if dangerousExts[ext] {
		response.BadRequest(c, "该文件类型不允许上传")
		return
	}
	if !allowedExts[ext] {
		response.BadRequest(c, "不支持的文件类型: "+ext)
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

// AddNote adds an admin note to a ticket (admin).
// POST /tickets/:id/notes
func (h *TicketHandler) AddNote(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ticket id")
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if len(req.Content) > 10000 {
		response.BadRequest(c, "content too long (max 10000 chars)")
		return
	}

	adminID := c.GetUint("user_id")
	note, err := h.ticketSvc.AddNote(uint(ticketID), adminID, req.Content)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, note)
}

// DeleteNote deletes a ticket note (admin).
// DELETE /tickets/notes/:id
func (h *TicketHandler) DeleteNote(c *gin.Context) {
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid note id")
		return
	}

	if err := h.ticketSvc.DeleteNote(uint(noteID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "note deleted")
}

// DeleteReply deletes a ticket reply (admin).
// DELETE /tickets/replies/:id
func (h *TicketHandler) DeleteReply(c *gin.Context) {
	replyID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid reply id")
		return
	}

	if err := h.ticketSvc.DeleteReply(uint(replyID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "reply deleted")
}

// GetTicketDetailHost returns hosts associated with a ticket's user (admin).
// GET /tickets/:id/hosts
func (h *TicketHandler) GetTicketDetailHost(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ticket id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	hosts, total, err := h.ticketSvc.GetTicketDetailHost(uint(ticketID), page, pageSize)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessPage(c, hosts, total, page, pageSize)
}

// TicketStatistics returns ticket statistics (admin).
// GET /tickets/statistics
func (h *TicketHandler) TicketStatistics(c *gin.Context) {
	stats, err := h.ticketSvc.GetTicketStatistics()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, stats)
}

// ==================== P1-8: DownloadAttachment ====================

// Upload handles a generic file upload (not tied to a specific ticket).
func (h *TicketHandler) Upload(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	if file.Size > 20*1024*1024 {
		response.BadRequest(c, "file size exceeds 20MB limit")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".txt": true, ".csv": true, ".zip": true, ".rar": true, ".7z": true,
		".mp4": true, ".mp3": true, ".wav": true,
	}
	dangerousExts := map[string]bool{
		".exe": true, ".bat": true, ".cmd": true, ".sh": true, ".ps1": true,
		".php": true, ".asp": true, ".aspx": true, ".jsp": true, ".py": true,
		".rb": true, ".pl": true, ".cgi": true, ".htaccess": true,
		".js": true, ".vbs": true, ".wsf": true, ".scr": true, ".com": true,
		".pif": true, ".msi": true, ".dll": true, ".so": true, ".dylib": true,
	}
	if dangerousExts[ext] {
		response.BadRequest(c, "该文件类型不允许上传")
		return
	}
	if !allowedExts[ext] {
		response.BadRequest(c, "不支持的文件类型: "+ext)
		return
	}

	destDir := filepath.Join(h.uploadDir, "general")
	if err := os.MkdirAll(destDir, 0755); err != nil {
		response.ServerError(c, "failed to create upload directory")
		return
	}

	savedName := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixNano(), ext)
	savedPath := filepath.Join(destDir, savedName)

	if err := c.SaveUploadedFile(file, savedPath); err != nil {
		response.ServerError(c, "failed to save file")
		return
	}

	var hash string
	f, err := os.Open(savedPath)
	if err == nil {
		h2 := sha256.New()
		io.Copy(h2, f)
		f.Close()
		hash = hex.EncodeToString(h2.Sum(nil))
	}

	response.Success(c, gin.H{
		"filename":  file.Filename,
		"saved_name": savedName,
		"saved_path": savedPath,
		"size":      file.Size,
		"hash":      hash,
	})
}

// DownloadAttachment 下载工单附件
func (h *TicketHandler) DownloadAttachment(c *gin.Context) {
	attID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid attachment id")
		return
	}

	att, err := h.ticketSvc.GetAttachmentByID(uint(attID))
	if err != nil {
		response.NotFound(c, "attachment not found")
		return
	}

	c.FileAttachment(att.FilePath, att.FileName)
}

// ==================== P1-9: TicketReceive ====================

// TicketReceive 工单接单/领取
func (h *TicketHandler) TicketReceive(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ticket id")
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.ticketSvc.ReceiveTicket(uint(ticketID), adminID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "工单已领取")
}

// ==================== P2-15: Enhanced GetList ====================

// GetListEnhanced returns all tickets with department permission and advanced filters (admin).
func (h *TicketHandler) GetListEnhanced(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}
	var deptID *uint
	if d := c.Query("dept_id"); d != "" {
		v, _ := strconv.ParseUint(d, 10, 64)
		did := uint(v)
		deptID = &did
	}
	var assigneeID *uint
	if a := c.Query("assignee_id"); a != "" {
		v, _ := strconv.ParseUint(a, 10, 64)
		aid := uint(v)
		assigneeID = &aid
	}
	var userID *uint
	if u := c.Query("user_id"); u != "" {
		v, _ := strconv.ParseUint(u, 10, 64)
		uid := uint(v)
		userID = &uid
	}
	priority := c.Query("priority")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	tickets, total, err := h.ticketSvc.GetListEnhanced(page, pageSize, status, keyword, deptID, assigneeID, userID, priority, startTime, endTime)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, tickets, total, page, pageSize)
}
