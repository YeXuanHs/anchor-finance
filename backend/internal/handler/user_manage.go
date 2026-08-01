package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserManageHandler handles admin user management HTTP requests.
type UserManageHandler struct {
	svc *service.UserManageService
	log *logger.Logger
}

// NewUserManageHandler creates a new UserManageHandler.
func NewUserManageHandler(svc *service.UserManageService, log *logger.Logger) *UserManageHandler {
	return &UserManageHandler{svc: svc, log: log}
}

// ==================== Search & Filter ====================

// Search searches clients by keyword.
// GET /manage/users
func (h *UserManageHandler) Search(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	users, total, err := h.svc.SearchClients(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, users, total, page, pageSize)
}

// Filter filters clients by status/group/date.
// GET /manage/users/filter
func (h *UserManageHandler) Filter(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var f service.ClientFilter
	c.ShouldBindQuery(&f)

	users, total, err := h.svc.FilterClients(page, pageSize, f)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, users, total, page, pageSize)
}

// GetSummary returns client statistics summary.
// GET /manage/users/summary
func (h *UserManageHandler) GetSummary(c *gin.Context) {
	summary, err := h.svc.GetClientSummary()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, summary)
}

// ==================== Client Lifecycle ====================

// Create creates a new client account.
// POST /manage/users
func (h *UserManageHandler) Create(c *gin.Context) {
	var req service.RegisterClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.svc.CreateClient(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, user)
}

// Close suspends a client account.
// POST /manage/users/:id/close
func (h *UserManageHandler) Close(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	if err := h.svc.CloseClient(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "client closed")
}

// Delete deletes a client account (soft/hard).
// DELETE /manage/users/:id
func (h *UserManageHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	soft := c.DefaultQuery("soft", "true") != "false"
	if err := h.svc.DeleteClient(uint(id), soft); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "client deleted")
}

// Ban disables a client account.
// POST /manage/users/:id/ban
func (h *UserManageHandler) Ban(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	var req struct {
		Reason string `json:"reason"`
	}
	c.ShouldBindJSON(&req)

	if err := h.svc.BanClient(uint(id), req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "client banned")
}

// Unban re-enables a client account.
// POST /manage/users/:id/unban
func (h *UserManageHandler) Unban(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	if err := h.svc.UnbanClient(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "client unbanned")
}

// CancelBan cancels a ban request.
// POST /manage/users/:id/cancel-ban
func (h *UserManageHandler) CancelBan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	if err := h.svc.CancelBan(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "ban cancelled")
}

// ==================== Client Profile ====================

// GetProfile returns full client profile with stats.
// GET /manage/users/:id/profile
func (h *UserManageHandler) GetProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	profile, err := h.svc.GetClientProfile(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, profile)
}

// UpdateProfile updates client profile fields.
// PUT /manage/users/:id/profile
func (h *UserManageHandler) UpdateProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateClientProfile(uint(id), data); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "profile updated")
}

// GetHosts returns a client's hosts.
// GET /manage/users/:id/hosts
func (h *UserManageHandler) GetHosts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	hosts, total, err := h.svc.GetClientHosts(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, hosts, total, page, pageSize)
}

// GetInvoices returns a client's invoices.
// GET /manage/users/:id/invoices
func (h *UserManageHandler) GetInvoices(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	invoices, total, err := h.svc.GetClientInvoices(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, page, pageSize)
}

// GetOrders returns a client's orders.
// GET /manage/users/:id/orders
func (h *UserManageHandler) GetOrders(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.svc.GetClientOrders(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, orders, total, page, pageSize)
}

// GetTickets returns a client's tickets.
// GET /manage/users/:id/tickets
func (h *UserManageHandler) GetTickets(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	tickets, total, err := h.svc.GetClientTickets(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, tickets, total, page, pageSize)
}

// GetLogs returns client activity logs.
// GET /manage/users/:id/logs
func (h *UserManageHandler) GetLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.svc.GetClientLogs(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// ==================== Client Notes ====================

// AddNote adds an admin note for a client.
// POST /manage/users/:id/notes
func (h *UserManageHandler) AddNote(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	var req struct {
		Content   string `json:"content" binding:"required"`
		IsPrivate bool   `json:"is_private"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := c.GetUint("user_id")
	note, err := h.svc.AddNote(uint(clientID), adminID, req.Content, req.IsPrivate)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, note)
}

// GetNotes returns all admin notes for a client.
// GET /manage/users/:id/notes
func (h *UserManageHandler) GetNotes(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	notes, err := h.svc.GetNotes(uint(clientID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, notes)
}

// DeleteNote deletes an admin note.
// DELETE /manage/notes/:id
func (h *UserManageHandler) DeleteNote(c *gin.Context) {
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid note id")
		return
	}

	if err := h.svc.DeleteNote(uint(noteID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "note deleted")
}

// ==================== Client Authorization ====================

// Authorize sets client permissions.
// POST /manage/users/:id/authorize
func (h *UserManageHandler) Authorize(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AuthorizeClient(uint(clientID), req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "client authorized")
}

// GetAuth returns auth settings for a client.
// GET /manage/users/:id/auth
func (h *UserManageHandler) GetAuth(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	auth, err := h.svc.GetClientAuth(uint(clientID))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, auth)
}

// ==================== Client Group ====================

// AssignGroup assigns a client to a group.
// POST /manage/users/:id/group
func (h *UserManageHandler) AssignGroup(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	var req struct {
		GroupID uint `json:"group_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AssignGroup(uint(clientID), req.GroupID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "group assigned")
}

// RemoveFromGroup removes a client from a group.
// DELETE /manage/users/:id/group/:group_id
func (h *UserManageHandler) RemoveFromGroup(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	groupID, err := strconv.ParseUint(c.Param("group_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	if err := h.svc.RemoveFromGroup(uint(clientID), uint(groupID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "removed from group")
}

// ==================== Certification ====================

// GetCertificationStatus returns real-name verification status.
// GET /manage/users/:id/certification
func (h *UserManageHandler) GetCertificationStatus(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	cert, err := h.svc.GetCertificationStatus(uint(clientID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cert)
}

// ReviewCertification reviews a certification submission.
// POST /manage/users/:id/certification/review
func (h *UserManageHandler) ReviewCertification(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	var req struct {
		Approved bool   `json:"approved"`
		Reason   string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	reviewerID := c.GetUint("user_id")
	if err := h.svc.ReviewCertification(uint(clientID), req.Approved, req.Reason, reviewerID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	msg := "certification approved"
	if !req.Approved {
		msg = "certification rejected"
	}
	response.SuccessMsg(c, msg)
}

// ==================== Cancel Requests ====================

// GetCancelRequests lists cancellation requests.
// GET /manage/cancel-requests
func (h *UserManageHandler) GetCancelRequests(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	requests, total, err := h.svc.GetCancelRequests(page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, requests, total, page, pageSize)
}

// ProcessCancelRequest approves or rejects a cancellation request.
// POST /manage/cancel-requests/:id
func (h *UserManageHandler) ProcessCancelRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid request id")
		return
	}

	var req struct {
		Approved bool   `json:"approved"`
		Remark   string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.svc.ProcessCancelRequest(uint(requestID), req.Approved, adminID, req.Remark); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	msg := "cancel request approved"
	if !req.Approved {
		msg = "cancel request rejected"
	}
	response.SuccessMsg(c, msg)
}

// ==================== Balance ====================

// AdjustBalance adds or deducts balance for a user.
// POST /manage/users/:id/balance
func (h *UserManageHandler) AdjustBalance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		Amount      float64 `json:"amount" binding:"required"`
		Description string  `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AdjustBalance(uint(id), req.Amount, req.Description); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "balance adjusted")
}

// ResetPassword sets a new password for a user.
// POST /manage/users/:id/reset-password
func (h *UserManageHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.ResetPassword(uint(id), req.NewPassword); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "password reset")
}

// GetStatus returns user status information.
// GET /manage/users/:id/status
func (h *UserManageHandler) GetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	info, err := h.svc.GetUserStatus(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, info)
}

// GetOperationLogs returns operation logs for a user.
// GET /manage/users/:id/operation-logs
func (h *UserManageHandler) GetOperationLogs(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.svc.GetOperationLogs(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// ==================== 删除原因管理 ====================

// DelReason deletes a cancel reason.
// DELETE /manage/cancel-reasons/:id
func (h *UserManageHandler) DelReason(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid reason id")
		return
	}

	if err := h.svc.DeleteCancelReason(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cancel reason deleted")
}

// GetCancelReasons returns all cancel reasons.
// GET /manage/cancel-reasons
func (h *UserManageHandler) GetCancelReasons(c *gin.Context) {
	reasons, err := h.svc.GetCancelReasons()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, reasons)
}

// AddCancelReason adds a new cancel reason.
// POST /manage/cancel-reasons
func (h *UserManageHandler) AddCancelReason(c *gin.Context) {
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AddCancelReason(req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cancel reason added")
}

// ==================== 充值发票 ====================

// AddRechargeInvoice creates a recharge invoice for a user.
// POST /manage/users/:id/recharge-invoice
func (h *UserManageHandler) AddRechargeInvoice(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
		Notes  string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	invoice, err := h.svc.CreateRechargeInvoice(uint(uid), req.Amount, req.Notes)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, invoice)
}

// ==================== 用户发票 ====================

// AddUserInvoice creates a blank invoice for a user.
// POST /manage/users/:id/invoice
func (h *UserManageHandler) AddUserInvoice(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	invoice, err := h.svc.CreateUserInvoice(uint(uid))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, invoice)
}

// ==================== 认证下载 ====================

// CertifiDownload downloads a certification image.
// GET /manage/certification/:id/download
func (h *UserManageHandler) CertifiDownload(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid certification id")
		return
	}

	fileType := c.Query("type")
	if fileType == "" {
		response.BadRequest(c, "type is required")
		return
	}

	filePath, fileName, err := h.svc.GetCertificationFile(uint(id), fileType)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	c.FileAttachment(filePath, fileName)
}

// ==================== 个人认证下载 ====================

// CertifiPersonDownload downloads a personal certification image.
// GET /manage/certification/person/:client_id/download
func (h *UserManageHandler) CertifiPersonDownload(c *gin.Context) {
	clientID, err := strconv.ParseUint(c.Param("client_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid client id")
		return
	}

	fileType := c.Query("type")
	if fileType == "" {
		response.BadRequest(c, "type is required")
		return
	}

	filePath, fileName, err := h.svc.GetPersonalCertificationFile(uint(clientID), fileType)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	c.FileAttachment(filePath, fileName)
}

// ==================== 认证历史日志 ====================

// GetCertifyHistoryLog returns certification history logs for a client.
// GET /manage/certification/:id/history
func (h *UserManageHandler) GetCertifyHistoryLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid certification id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.svc.GetCertifyHistoryLog(uint(id), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// ==================== 主机绑定销售 ====================

// HostBindSale binds a salesperson to a user.
// POST /manage/users/:id/bind-sale
func (h *UserManageHandler) HostBindSale(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		SaleID uint `json:"sale_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.BindSale(uint(uid), req.SaleID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "sale bound successfully")
}

// ==================== 关联用户列表 ====================

// RelationUserList returns a list of users for relation selection.
// GET /manage/users/relation-list
func (h *UserManageHandler) RelationUserList(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	users, total, err := h.svc.GetRelationUserList(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, users, total, page, pageSize)
}

// ==================== 用户产品发票 ====================

// UserProductInvoice returns invoices related to a user's products.
// GET /manage/users/:id/product-invoices
func (h *UserManageHandler) UserProductInvoice(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	order := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "desc")

	invoices, total, err := h.svc.GetUserProductInvoices(uint(uid), page, pageSize, order, sort)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, invoices, total, page, pageSize)
}

// ==================== 新增缺失方法 ====================

// HostByUid returns hosts for a specific user.
// GET /manage/users/:uid/hosts-by-uid
func (h *UserManageHandler) HostByUid(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("uid"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	hosts, err := h.svc.GetHostsByUID(uint(uid))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, hosts)
}

// CerifyList returns certification list.
// GET /manage/certification/list
func (h *UserManageHandler) CerifyList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	items, total, err := h.svc.GetCertifyList(page, pageSize, status, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// CerifyLogList returns certification log list (v2).
// GET /manage/certification/log-list
func (h *UserManageHandler) CerifyLogList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	uid := c.Query("uid")
	logType := c.Query("type")

	logs, total, err := h.svc.GetCertifyLogList(page, pageSize, uid, logType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// CertifiPersonDetail returns personal certification details.
// GET /manage/certification/person/:id
func (h *UserManageHandler) CertifiPersonDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	detail, err := h.svc.GetCertifiPersonDetail(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, detail)
}

// CertifiPersonModify modifies personal certification.
// PUT /manage/certification/person/:id
func (h *UserManageHandler) CertifiPersonModify(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.ModifyCertifiPerson(uint(id), data); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "personal certification modified")
}

// CertifiCompanyDetail returns company certification details.
// GET /manage/certification/company/:id
func (h *UserManageHandler) CertifiCompanyDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	detail, err := h.svc.GetCertifiCompanyDetail(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, detail)
}

// CertifiCompanyModify modifies company certification.
// PUT /manage/certification/company/:id
func (h *UserManageHandler) CertifiCompanyModify(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.ModifyCertifiCompany(uint(id), data); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "company certification modified")
}

// UserProductaccounts returns product accounts for a user.
// GET /manage/users/:id/product-accounts
func (h *UserManageHandler) UserProductaccounts(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	accounts, total, err := h.svc.GetUserProductAccounts(uint(uid), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, accounts, total, page, pageSize)
}

// LoginByUser logs in as a specific user (admin impersonation).
// POST /manage/users/:id/login-as
func (h *UserManageHandler) LoginByUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	token, err := h.svc.GenerateUserToken(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{"token": token})
}

// DeleteCancelRequest deletes a cancel request.
// DELETE /manage/cancel-requests/:id
func (h *UserManageHandler) DeleteCancelRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	if err := h.svc.DeleteCancelRequest(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cancel request deleted")
}

// AddRecordLog adds an operation record log.
// POST /manage/users/:id/record-log
func (h *UserManageHandler) AddRecordLog(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.svc.AddRecordLog(uint(uid), adminID, req.Content); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "record log added")
}

// AddRemarkLog adds a remark log for a user.
// POST /manage/users/:id/remark-log
func (h *UserManageHandler) AddRemarkLog(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.svc.AddRemarkLog(uint(uid), adminID, req.Content); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "remark log added")
}
