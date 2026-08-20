package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type RuleHandler struct {
	ruleSvc *service.RuleService
	log     *logger.Logger
}

func NewRuleHandler(ruleSvc *service.RuleService, log *logger.Logger) *RuleHandler {
	return &RuleHandler{ruleSvc: ruleSvc, log: log}
}

// List returns paginated rules with filters.
func (h *RuleHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	ruleType := c.Query("type")
	keyword := c.Query("keyword")

	var status *int16
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		sv := int16(v)
		status = &sv
	}

	rules, total, err := h.ruleSvc.List(page, pageSize, ruleType, status, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, rules, total, page, pageSize)
}

// GetDetail returns a single rule by ID.
func (h *RuleHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	rule, err := h.ruleSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "rule not found")
		return
	}
	response.Success(c, rule)
}

// Create creates a new rule.
func (h *RuleHandler) Create(c *gin.Context) {
	var req service.CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rule, err := h.ruleSvc.Create(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, rule)
}

// Update updates a rule.
func (h *RuleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	var req service.UpdateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rule, err := h.ruleSvc.Update(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, rule)
}

// Delete deletes a rule.
func (h *RuleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	if err := h.ruleSvc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "rule deleted")
}

// Enable enables a rule.
func (h *RuleHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	if err := h.ruleSvc.Enable(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "rule enabled")
}

// Disable disables a rule.
func (h *RuleHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	if err := h.ruleSvc.Disable(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "rule disabled")
}

// Test tests a rule against provided data.
func (h *RuleHandler) Test(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	var req service.TestRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.ruleSvc.Test(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetRuleLogs returns execution logs for rules.
func (h *RuleHandler) GetRuleLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var ruleID *uint
	if rid := c.Query("rule_id"); rid != "" {
		v, _ := strconv.ParseUint(rid, 10, 64)
		id := uint(v)
		ruleID = &id
	}

	var success *bool
	if s := c.Query("success"); s != "" {
		b := s == "true"
		success = &b
	}

	logs, total, err := h.ruleSvc.GetRuleLogs(page, pageSize, ruleID, success)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}
