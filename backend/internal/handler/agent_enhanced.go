package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// AgentEnhancedHandler handles enhanced agent HTTP requests.
type AgentEnhancedHandler struct {
	svc *service.AgentEnhancedService
	log *logger.Logger
}

// NewAgentEnhancedHandler creates a new AgentEnhancedHandler.
func NewAgentEnhancedHandler(svc *service.AgentEnhancedService, log *logger.Logger) *AgentEnhancedHandler {
	return &AgentEnhancedHandler{svc: svc, log: log}
}

// resolveAgentID looks up the agent by the authenticated user ID and returns agent.ID.
func (h *AgentEnhancedHandler) resolveAgentID(c *gin.Context) uint {
	uid := getUserID(c)
	if uid == 0 {
		return 0
	}
	agent, err := h.svc.GetAgentByUserID(uid)
	if err != nil {
		response.NotFound(c, "agent not found")
		return 0
	}
	return agent.ID
}

// ==================== Resource Management ====================

// GetResourceInfo returns resource info for the authenticated agent.
// GET /agent/resources
func (h *AgentEnhancedHandler) GetResourceInfo(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	info, err := h.svc.GetResourceInfo(agentID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, info)
}

// PostResourceInfo updates resource info for the authenticated agent.
// POST /agent/resources
func (h *AgentEnhancedHandler) PostResourceInfo(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.PostResourceInfo(agentID, data); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "resource info updated")
}

// GetProducts returns products available to the agent.
// GET /agent/products
func (h *AgentEnhancedHandler) GetProducts(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	products, err := h.svc.GetProducts(agentID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, products)
}

// GetHostLists returns host lists for the agent.
// GET /agent/hosts
func (h *AgentEnhancedHandler) GetHostLists(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	hosts, err := h.svc.GetHostLists(agentID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, hosts)
}

// ==================== Inspection ====================

// GetInspectionLists returns paginated inspection records.
// GET /agent/inspections
func (h *AgentEnhancedHandler) GetInspectionLists(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, total, err := h.svc.GetInspectionLists(agentID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// createInspectionRequest is the payload for CreateInspection.
type createInspectionRequest struct {
	HostID uint   `json:"host_id" binding:"required"`
	Type   string `json:"type" binding:"required,oneof=routine emergency"`
}

// CreateInspection creates a new inspection.
// POST /agent/inspections
func (h *AgentEnhancedHandler) CreateInspection(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	var req createInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	inspection, err := h.svc.CreateInspection(agentID, req.HostID, req.Type)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, inspection)
}

// uploadInspectionRequest is the payload for PostUpload.
type uploadInspectionRequest struct {
	Images []string `json:"images" binding:"required,min=1"`
}

// PostUpload uploads images for an inspection.
// POST /agent/inspections/:id/upload
func (h *AgentEnhancedHandler) PostUpload(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid inspection id")
		return
	}

	var req uploadInspectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.PostUpload(uint(id), req.Images); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "images uploaded")
}

// GetInspectionDetail returns inspection detail.
// GET /agent/inspections/:id
func (h *AgentEnhancedHandler) GetInspectionDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid inspection id")
		return
	}

	inspection, err := h.svc.GetInspectionDetail(uint(id))
	if err != nil {
		response.NotFound(c, "inspection not found")
		return
	}
	response.Success(c, inspection)
}

// GetInspectionIPs returns IPs for inspection selection.
// GET /agent/inspections/ips
func (h *AgentEnhancedHandler) GetInspectionIPs(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	ips, err := h.svc.GetInspectionIPs(agentID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, ips)
}

// ==================== Orders & Finance ====================

// GetOrderSearchPage searches orders with filters.
// GET /agent/orders/search
func (h *AgentEnhancedHandler) GetOrderSearchPage(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	filter := service.OrderFilter{
		Keyword: c.Query("keyword"),
		Status:  c.Query("status"),
		Type:    c.Query("type"),
		StartAt: c.Query("start_at"),
		EndAt:   c.Query("end_at"),
	}

	orders, total, err := h.svc.GetOrderSearchPage(agentID, page, pageSize, filter)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, orders, total, page, pageSize)
}

// GetOrders returns paginated orders.
// GET /agent/orders
func (h *AgentEnhancedHandler) GetOrders(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	orders, total, err := h.svc.GetOrders(agentID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, orders, total, page, pageSize)
}

// GetRenewSearchPage searches renewal orders.
// GET /agent/renews/search
func (h *AgentEnhancedHandler) GetRenewSearchPage(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	filter := service.OrderFilter{
		Keyword: c.Query("keyword"),
		Status:  c.Query("status"),
		StartAt: c.Query("start_at"),
		EndAt:   c.Query("end_at"),
	}

	renews, total, err := h.svc.GetRenewSearchPage(agentID, page, pageSize, filter)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, renews, total, page, pageSize)
}

// GetRenews returns paginated renewal orders.
// GET /agent/renews
func (h *AgentEnhancedHandler) GetRenews(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	renews, total, err := h.svc.GetRenews(agentID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, renews, total, page, pageSize)
}

// GetIncome returns income statistics.
// GET /agent/income
func (h *AgentEnhancedHandler) GetIncome(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	period := c.DefaultQuery("period", "month")

	income, err := h.svc.GetIncome(agentID, period)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, income)
}

// GetConsumption returns consumption statistics.
// GET /agent/consumption
func (h *AgentEnhancedHandler) GetConsumption(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	period := c.DefaultQuery("period", "month")

	consumption, err := h.svc.GetConsumption(agentID, period)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, consumption)
}

// GetAgentLogs returns agent activity logs.
// GET /agent/logs
func (h *AgentEnhancedHandler) GetAgentLogs(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := h.svc.GetAgentLogs(agentID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// ==================== After-sale ====================

// GetAfterSaleDetail returns after-sale detail.
// GET /agent/aftersale/:id
func (h *AgentEnhancedHandler) GetAfterSaleDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid aftersale id")
		return
	}

	item, err := h.svc.GetAfterSaleDetail(uint(id))
	if err != nil {
		response.NotFound(c, "aftersale not found")
		return
	}
	response.Success(c, item)
}

// GetRefundDetail returns refund detail.
// GET /agent/refunds/:id
func (h *AgentEnhancedHandler) GetRefundDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid refund id")
		return
	}

	item, err := h.svc.GetRefundDetail(uint(id))
	if err != nil {
		response.NotFound(c, "refund not found")
		return
	}
	response.Success(c, item)
}

// PostAfterSale processes an after-sale request.
// POST /agent/aftersale/:id
func (h *AgentEnhancedHandler) PostAfterSale(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid aftersale id")
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.PostAfterSale(uint(id), data); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "aftersale processed")
}

// PostUnAfterSale cancels an after-sale request.
// POST /agent/aftersale/:id/cancel
func (h *AgentEnhancedHandler) PostUnAfterSale(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid aftersale id")
		return
	}

	if err := h.svc.PostUnAfterSale(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "aftersale cancelled")
}

// PostRefund processes a refund.
// POST /agent/refunds/:id
func (h *AgentEnhancedHandler) PostRefund(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid refund id")
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.PostRefund(uint(id), data); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "refund processed")
}

// ==================== Tickets ====================

// GetTickets returns paginated tickets for the agent.
// GET /agent/tickets
func (h *AgentEnhancedHandler) GetTickets(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	tickets, total, err := h.svc.GetTickets(agentID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, tickets, total, page, pageSize)
}

// ==================== Evaluation ====================

// PostEvaluation submits an evaluation.
// POST /agent/evaluation
func (h *AgentEnhancedHandler) PostEvaluation(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.PostEvaluation(agentID, data); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "evaluation submitted")
}

// GetRunMapLists returns run map lists.
// GET /agent/run-maps
func (h *AgentEnhancedHandler) GetRunMapLists(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	items, err := h.svc.GetRunMapLists(agentID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// ==================== Token Management ====================

// GetToken returns the agent's API token.
// GET /agent/token
func (h *AgentEnhancedHandler) GetToken(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	token, err := h.svc.GetToken(agentID)
	if err != nil {
		response.NotFound(c, "no active token found")
		return
	}
	response.Success(c, token)
}

// setTokenRequest is the payload for SetToken.
type setTokenRequest struct {
	Name string `json:"name" binding:"required"`
}

// SetToken creates a new API token.
// POST /agent/token
func (h *AgentEnhancedHandler) SetToken(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	var req setTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	token, err := h.svc.SetToken(agentID, req.Name)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, token)
}

// CheckToken validates an API token (public endpoint).
// POST /agent/token/check
func (h *AgentEnhancedHandler) CheckToken(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	agent, err := h.svc.CheckToken(req.Token)
	if err != nil {
		response.Unauthorized(c, "invalid or expired token")
		return
	}
	response.Success(c, agent)
}

// ==================== Base Info ====================

// GetBaseInfo returns base info for the agent.
// GET /agent/base-info
func (h *AgentEnhancedHandler) GetBaseInfo(c *gin.Context) {
	agentID := h.resolveAgentID(c)
	if agentID == 0 {
		return
	}

	info, err := h.svc.GetBaseInfo(agentID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, info)
}
