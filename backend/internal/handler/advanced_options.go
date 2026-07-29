package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdvancedOptionsHandler struct {
	svc *service.AdvancedOptionsService
	log *logger.Logger
}

func NewAdvancedOptionsHandler(svc *service.AdvancedOptionsService, log *logger.Logger) *AdvancedOptionsHandler {
	return &AdvancedOptionsHandler{svc: svc, log: log}
}

// GetOptions returns advanced configuration options for a product.
func (h *AdvancedOptionsHandler) GetOptions(c *gin.Context) {
	pid, err := strconv.ParseUint(c.Query("pid"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	options, err := h.svc.GetOptions(uint(pid))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	links, err := h.svc.GetLinks(uint(pid))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"product_id": pid,
		"options":    options,
		"links":      links,
	})
}

// CreateOption creates a new advanced configuration option.
func (h *AdvancedOptionsHandler) CreateOption(c *gin.Context) {
	var req struct {
		ProductID   uint   `json:"product_id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Value       string `json:"value"`
		Required    bool   `json:"required"`
		SortOrder   int    `json:"sort_order"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	option := &model.AdvancedOption{
		ProductID:   req.ProductID,
		Name:        req.Name,
		Type:        req.Type,
		Value:       req.Value,
		Required:    req.Required,
		SortOrder:   req.SortOrder,
		Status:      1,
		Description: req.Description,
	}

	if err := h.svc.CreateOption(option); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, option)
}

// UpdateOption updates an existing advanced configuration option.
func (h *AdvancedOptionsHandler) UpdateOption(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid option id")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Value       string `json:"value"`
		Required    *bool  `json:"required"`
		SortOrder   *int   `json:"sort_order"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Value != "" {
		updates["value"] = req.Value
	}
	if req.Required != nil {
		updates["required"] = *req.Required
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if err := h.svc.UpdateOption(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "option updated")
}

// DeleteOption deletes an advanced configuration option.
func (h *AdvancedOptionsHandler) DeleteOption(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid option id")
		return
	}

	if err := h.svc.DeleteOption(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "option deleted")
}

// GetLinks returns configuration links for a product.
func (h *AdvancedOptionsHandler) GetLinks(c *gin.Context) {
	pid, err := strconv.ParseUint(c.Query("pid"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	links, err := h.svc.GetLinks(uint(pid))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"product_id": pid,
		"links":      links,
	})
}

// CreateLink creates a new configuration link.
func (h *AdvancedOptionsHandler) CreateLink(c *gin.Context) {
	var req struct {
		ProductID uint   `json:"product_id" binding:"required"`
		ConfigID  uint   `json:"config_id" binding:"required"`
		Relation  string `json:"relation" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	link := &model.AdvancedOptionLink{
		ProductID: req.ProductID,
		ConfigID:  req.ConfigID,
		Relation:  req.Relation,
		Status:    1,
	}

	if err := h.svc.CreateLink(link); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, link)
}

// UpdateLink updates an existing configuration link.
func (h *AdvancedOptionsHandler) UpdateLink(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid link id")
		return
	}

	var req struct {
		ConfigID *uint  `json:"config_id"`
		Relation string `json:"relation"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.ConfigID != nil {
		updates["config_id"] = *req.ConfigID
	}
	if req.Relation != "" {
		updates["relation"] = req.Relation
	}

	if err := h.svc.UpdateLink(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "link updated")
}

// DeleteLink deletes a configuration link.
func (h *AdvancedOptionsHandler) DeleteLink(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid link id")
		return
	}

	if err := h.svc.DeleteLink(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "link deleted")
}
