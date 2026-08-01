package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type TicketStatusHandler struct {
	svc *service.TicketStatusService
	log *logger.Logger
}

func NewTicketStatusHandler(svc *service.TicketStatusService, log *logger.Logger) *TicketStatusHandler {
	return &TicketStatusHandler{svc: svc, log: log}
}

// GetStatuses returns a list of ticket statuses.
func (h *TicketStatusHandler) GetStatuses(c *gin.Context) {
	statuses, err := h.svc.GetAllStatuses()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"statuses": statuses})
}

// AddStatus adds a new ticket status.
func (h *TicketStatusHandler) AddStatus(c *gin.Context) {
	var req struct {
		Title      string `json:"title" binding:"required"`
		Code       string `json:"code" binding:"required"`
		Color      string `json:"color"`
		ShowActive int8   `json:"show_active"`
		Order      int    `json:"order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	status := &model.TicketStatus{
		Title:      req.Title,
		Code:       req.Code,
		Color:      req.Color,
		ShowActive: req.ShowActive,
		Order:      req.Order,
		Status:     1,
	}

	if err := h.svc.AddStatus(status); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, status)
}

// UpdateStatus updates an existing ticket status.
func (h *TicketStatusHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid status id")
		return
	}

	var req struct {
		Title      string `json:"title"`
		Code       string `json:"code"`
		Color      string `json:"color"`
		ShowActive *int8  `json:"show_active"`
		Order      *int   `json:"order"`
		Status     *int16 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}
	if req.Color != "" {
		updates["color"] = req.Color
	}
	if req.ShowActive != nil {
		updates["show_active"] = *req.ShowActive
	}
	if req.Order != nil {
		updates["order"] = *req.Order
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := h.svc.UpdateStatus(uint(id), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "status updated")
}

// DeleteStatus deletes a ticket status.
func (h *TicketStatusHandler) DeleteStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid status id")
		return
	}

	if err := h.svc.DeleteStatus(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "status deleted")
}

// GetDefaultStatuses returns default ticket statuses.
func (h *TicketStatusHandler) GetDefaultStatuses(c *gin.Context) {
	defaults := h.svc.GetDefaultStatuses()
	response.Success(c, gin.H{"default": defaults})
}

// GetStatusDetail returns a single ticket status by ID.
// GET /admin/ticket-statuses/:id
func (h *TicketStatusHandler) GetStatusDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid status id")
		return
	}

	status, err := h.svc.GetStatusByID(uint(id))
	if err != nil {
		response.NotFound(c, "status not found")
		return
	}
	response.Success(c, status)
}
