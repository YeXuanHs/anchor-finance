package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type TicketDeliverHandler struct {
	svc *service.TicketDeliverService
	log *logger.Logger
}

func NewTicketDeliverHandler(svc *service.TicketDeliverService, log *logger.Logger) *TicketDeliverHandler {
	return &TicketDeliverHandler{svc: svc, log: log}
}

// GetAddPage returns data for adding a ticket deliver rule.
func (h *TicketDeliverHandler) GetAddPage(c *gin.Context) {
	data, err := h.svc.GetAddPage()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, data)
}

// GetRules returns a list of ticket deliver rules.
func (h *TicketDeliverHandler) GetRules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	rules, total, err := h.svc.GetRules(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, rules, total, page, pageSize)
}

// GetRule returns a single ticket deliver rule.
func (h *TicketDeliverHandler) GetRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	rule, err := h.svc.GetRuleByID(uint(id))
	if err != nil {
		response.NotFound(c, "rule not found")
		return
	}
	response.Success(c, rule)
}

// CreateRule creates a new ticket deliver rule.
func (h *TicketDeliverHandler) CreateRule(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Departments string `json:"departments" binding:"required"`
		Products    string `json:"products"`
		Priority    int    `json:"priority"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rule := &model.TicketDeliverRule{
		Name:        req.Name,
		Priority:    req.Priority,
		Status:      1,
		Description: req.Description,
	}

	if err := h.svc.CreateRule(rule); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, rule)
}

// UpdateRule updates an existing ticket deliver rule.
func (h *TicketDeliverHandler) UpdateRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Departments *string `json:"departments"`
		Products    *string `json:"products"`
		Priority    *int    `json:"priority"`
		Description *string `json:"description"`
		Status      *int16  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Departments != nil {
		updates["departments"] = *req.Departments
	}
	if req.Products != nil {
		updates["products"] = *req.Products
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := h.svc.UpdateRule(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "rule updated")
}

// DeleteRule deletes a ticket deliver rule.
func (h *TicketDeliverHandler) DeleteRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	if err := h.svc.DeleteRule(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "rule deleted")
}

// GetLogs returns deliver logs.
func (h *TicketDeliverHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var ticketID uint
	if tid := c.Query("ticket_id"); tid != "" {
		v, _ := strconv.ParseUint(tid, 10, 64)
		ticketID = uint(v)
	}

	logs, total, err := h.svc.GetLogs(ticketID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}
