package handler

import (
	"strconv"

	"github.com/anchor-finance/backend/internal/model"
	"github.com/anchor-finance/backend/pkg/logger"
	"github.com/anchor-finance/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ClientCareHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewClientCareHandler(db *gorm.DB, log *logger.Logger) *ClientCareHandler {
	return &ClientCareHandler{db: db, log: log}
}

// GetRules returns all client care rules (admin).
func (h *ClientCareHandler) GetRules(c *gin.Context) {
	var rules []model.ClientCareRule
	if err := h.db.Order("id DESC").Find(&rules).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, rules)
}

type CreateRuleRequest struct {
	Name       string     `json:"name" binding:"required,max=100"`
	Type       string     `json:"type" binding:"required,max=50"`
	Condition  model.JSON `json:"condition"`
	TemplateID uint       `json:"template_id"`
	Channel    string     `json:"channel" binding:"max=20"`
	DaysBefore int        `json:"days_before"`
}

// CreateRule creates a new client care rule (admin).
func (h *ClientCareHandler) CreateRule(c *gin.Context) {
	var req CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	rule := model.ClientCareRule{
		Name:       req.Name,
		Type:       req.Type,
		Condition:  req.Condition,
		TemplateID: req.TemplateID,
		Channel:    req.Channel,
		DaysBefore: req.DaysBefore,
		IsActive:   true,
	}

	if err := h.db.Create(&rule).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, rule)
}

type UpdateRuleRequest struct {
	Name       *string     `json:"name"`
	Type       *string     `json:"type"`
	Condition  *model.JSON `json:"condition"`
	TemplateID *uint       `json:"template_id"`
	Channel    *string     `json:"channel"`
	DaysBefore *int        `json:"days_before"`
	IsActive   *bool       `json:"is_active"`
}

// UpdateRule updates a client care rule (admin).
func (h *ClientCareHandler) UpdateRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	var req UpdateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Condition != nil {
		updates["condition"] = *req.Condition
	}
	if req.TemplateID != nil {
		updates["template_id"] = *req.TemplateID
	}
	if req.Channel != nil {
		updates["channel"] = *req.Channel
	}
	if req.DaysBefore != nil {
		updates["days_before"] = *req.DaysBefore
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := h.db.Model(&model.ClientCareRule{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}
	response.SuccessMsg(c, "rule updated")
}

// DeleteRule deletes a client care rule (admin).
func (h *ClientCareHandler) DeleteRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	if err := h.db.Delete(&model.ClientCareRule{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "rule deleted")
}

// GetLogs returns client care logs with pagination (admin).
func (h *ClientCareHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var ruleID *uint
	if rid := c.Query("rule_id"); rid != "" {
		v, _ := strconv.ParseUint(rid, 10, 64)
		uid := uint(v)
		ruleID = &uid
	}
	var userID *uint
	if uid := c.Query("user_id"); uid != "" {
		v, _ := strconv.ParseUint(uid, 10, 64)
		u := uint(v)
		userID = &u
	}

	var logs []model.ClientCareLog
	var total int64

	query := h.db.Model(&model.ClientCareLog{})
	if ruleID != nil {
		query = query.Where("rule_id = ?", *ruleID)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}
