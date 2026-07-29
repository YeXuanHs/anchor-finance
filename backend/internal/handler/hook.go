package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type HookHandler struct {
	svc *service.HookService
	log *logger.Logger
}

func NewHookHandler(svc *service.HookService, log *logger.Logger) *HookHandler {
	return &HookHandler{svc: svc, log: log}
}

// GetList returns paginated hooks.
func (h *HookHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	event := c.Query("event")

	var status *int16
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		sv := int16(v)
		status = &sv
	}

	hooks, total, err := h.svc.GetList(page, pageSize, event, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, hooks, total, page, pageSize)
}

// GetDetail returns a single hook.
func (h *HookHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid hook id")
		return
	}

	hook, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "hook not found")
		return
	}
	response.Success(c, hook)
}

// Create creates a new hook.
func (h *HookHandler) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Event       string `json:"event" binding:"required"`
		Type        string `json:"type" binding:"required,oneof=url script plugin"`
		URL         string `json:"url"`
		Script      string `json:"script"`
		Headers     string `json:"headers"`
		Params      string `json:"params"`
		Timeout     int    `json:"timeout"`
		RetryCount  int    `json:"retry_count"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Timeout == 0 {
		req.Timeout = 30
	}

	hook := &model.Hook{
		Name:        req.Name,
		Code:        req.Code,
		Event:       req.Event,
		Type:        req.Type,
		URL:         req.URL,
		Script:      req.Script,
		Status:      1,
		Timeout:     req.Timeout,
		RetryCount:  req.RetryCount,
		Description: req.Description,
	}

	if err := h.svc.Create(hook); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, hook)
}

// Update updates a hook.
func (h *HookHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid hook id")
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Event       *string `json:"event"`
		Type        *string `json:"type"`
		URL         *string `json:"url"`
		Script      *string `json:"script"`
		Headers     *string `json:"headers"`
		Params      *string `json:"params"`
		Timeout     *int    `json:"timeout"`
		RetryCount  *int    `json:"retry_count"`
		Description *string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Event != nil {
		updates["event"] = *req.Event
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.URL != nil {
		updates["url"] = *req.URL
	}
	if req.Script != nil {
		updates["script"] = *req.Script
	}
	if req.Headers != nil {
		updates["headers"] = *req.Headers
	}
	if req.Params != nil {
		updates["params"] = *req.Params
	}
	if req.Timeout != nil {
		updates["timeout"] = *req.Timeout
	}
	if req.RetryCount != nil {
		updates["retry_count"] = *req.RetryCount
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}

	if err := h.svc.Update(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "hook updated")
}

// Delete deletes a hook.
func (h *HookHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid hook id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "hook deleted")
}

// SetStatus sets the status of a hook.
func (h *HookHandler) SetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid hook id")
		return
	}

	var req struct {
		Status int16 `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.SetStatus(uint(id), req.Status); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "hook status updated")
}

// Trigger triggers a hook manually.
func (h *HookHandler) Trigger(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid hook id")
		return
	}

	hook, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "hook not found")
		return
	}

	if err := h.svc.Trigger(hook.Event, nil); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "hook triggered")
}

// GetLogs returns paginated hook execution logs.
func (h *HookHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var hookID uint
	if hid := c.Query("hook_id"); hid != "" {
		v, _ := strconv.ParseUint(hid, 10, 64)
		hookID = uint(v)
	}

	logs, total, err := h.svc.GetLogs(hookID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}
