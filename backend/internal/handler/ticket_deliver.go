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

// ────────────────────────────────────────────────────────────
// 上游透传管理
// ────────────────────────────────────────────────────────────

// GetUpstreams returns all upstream configurations.
func (h *TicketDeliverHandler) GetUpstreams(c *gin.Context) {
	upstreams, err := h.svc.GetUpstreams()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"upstreams": upstreams})
}

// GetUpstream returns a single upstream configuration.
func (h *TicketDeliverHandler) GetUpstream(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid upstream id")
		return
	}

	upstream, err := h.svc.GetUpstreamByID(uint(id))
	if err != nil {
		response.NotFound(c, "upstream not found")
		return
	}
	response.Success(c, upstream)
}

// CreateUpstream creates a new upstream configuration.
func (h *TicketDeliverHandler) CreateUpstream(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		URL         string `json:"url" binding:"required"`
		APIKey      string `json:"api_key"`
		Username    string `json:"username"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证类型
	validTypes := map[string]bool{"zjmf": true, "v10": true, "anchorfinance": true}
	if !validTypes[req.Type] {
		response.BadRequest(c, "不支持的上游类型，支持: zjmf, v10, anchorfinance")
		return
	}

	upstream := &model.TicketUpstream{
		Name:        req.Name,
		Type:        req.Type,
		URL:         req.URL,
		APIKey:      req.APIKey,
		Username:    req.Username,
		Status:      1,
		Description: req.Description,
	}

	if err := h.svc.CreateUpstream(upstream); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, upstream)
}

// UpdateUpstream updates an upstream configuration.
func (h *TicketDeliverHandler) UpdateUpstream(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid upstream id")
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Type        *string `json:"type"`
		URL         *string `json:"url"`
		APIKey      *string `json:"api_key"`
		Username    *string `json:"username"`
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
	if req.Type != nil {
		validTypes := map[string]bool{"zjmf": true, "v10": true, "anchorfinance": true}
		if !validTypes[*req.Type] {
			response.BadRequest(c, "不支持的上游类型")
			return
		}
		updates["type"] = *req.Type
	}
	if req.URL != nil {
		updates["url"] = *req.URL
	}
	if req.APIKey != nil {
		updates["api_key"] = *req.APIKey
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := h.svc.UpdateUpstream(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "upstream updated")
}

// DeleteUpstream deletes an upstream configuration.
func (h *TicketDeliverHandler) DeleteUpstream(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid upstream id")
		return
	}

	if err := h.svc.DeleteUpstream(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "upstream deleted")
}

// TestUpstreamConnection tests connection to an upstream system.
func (h *TicketDeliverHandler) TestUpstreamConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid upstream id")
		return
	}

	success, message, err := h.svc.TestUpstreamConnection(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"success": success,
		"message": message,
	})
}

// GetSupportedTypes returns supported upstream types.
func (h *TicketDeliverHandler) GetSupportedTypes(c *gin.Context) {
	types := []map[string]interface{}{
		{
			"type":        "anchorfinance",
			"name":        "锚点财务",
			"description": "另一个锚点财务实例",
		},
		{
			"type":        "zjmf",
			"name":        "智简魔方",
			"description": "智简魔方财务系统",
		},
		{
			"type":        "v10",
			"name":        "V10",
			"description": "V10财务系统",
		},
	}
	response.Success(c, gin.H{"types": types})
}
