package handler

import (
	"net/http"
	"strconv"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AgentHandler handles agent (affiliate) HTTP requests.
type AgentHandler struct {
	db *gorm.DB
}

// NewAgentHandler creates a new AgentHandler.
func NewAgentHandler(db *gorm.DB) *AgentHandler {
	return &AgentHandler{db: db}
}

// GetInfo returns the authenticated user's agent info.
// GET /agent
func (h *AgentHandler) GetInfo(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var agent model.Agent
	if err := h.db.Where("user_id = ?", userID).First(&agent).Error; err != nil {
		response.NotFound(c, "agent not found")
		return
	}

	response.Success(c, agent)
}

// GetDownlines returns the downline users of the authenticated agent.
// GET /agent/downlines
func (h *AgentHandler) GetDownlines(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var agent model.Agent
	if err := h.db.Where("user_id = ?", userID).First(&agent).Error; err != nil {
		response.NotFound(c, "agent not found")
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

	var agents []model.Agent
	var total int64

	query := h.db.Model(&model.Agent{}).Where("parent_id = ?", agent.ID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("User").Offset(offset).Limit(pageSize).Order("id DESC").Find(&agents)

	response.SuccessPage(c, agents, total, page, pageSize)
}

// GetCommissions returns the authenticated agent's commission records.
// GET /agent/commissions
func (h *AgentHandler) GetCommissions(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var agent model.Agent
	if err := h.db.Where("user_id = ?", userID).First(&agent).Error; err != nil {
		response.NotFound(c, "agent not found")
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

	var commissions []model.AgentCommission
	var total int64

	query := h.db.Model(&model.AgentCommission{}).Where("agent_id = ?", agent.ID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&commissions)

	response.SuccessPage(c, commissions, total, page, pageSize)
}

// withdrawRequest is the payload for ApplyWithdraw.
type withdrawRequest struct {
	Amount float64 `json:"amount" binding:"required,gt=0"`
}

// ApplyWithdraw requests a commission withdrawal.
// POST /agent/withdraw
func (h *AgentHandler) ApplyWithdraw(c *gin.Context) {
	userID := getUserID(c)
	if userID == 0 {
		return
	}

	var req withdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var agent model.Agent
	if err := h.db.Where("user_id = ?", userID).First(&agent).Error; err != nil {
		response.NotFound(c, "agent not found")
		return
	}

	if agent.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient commission balance"})
		return
	}

	if err := h.db.Model(&agent).UpdateColumn("balance", gorm.Expr("balance - ?", req.Amount)).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "withdrawal request submitted")
}

// AdminGetList returns a paginated list of agents (admin).
// GET /admin/agents
func (h *AgentHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var agents []model.Agent
	var total int64

	query := h.db.Model(&model.Agent{})
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Joins("LEFT JOIN users ON users.id = agents.user_id").
			Where("agents.agent_no LIKE ? OR users.username LIKE ? OR users.email LIKE ?", q, q, q)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Preload("User").Preload("Parent").Offset(offset).Limit(pageSize).Order("id DESC").Find(&agents)

	response.SuccessPage(c, agents, total, page, pageSize)
}

// adminCreateAgentRequest is the payload for AdminCreate.
type adminCreateAgentRequest struct {
	UserID         uint    `json:"user_id" binding:"required"`
	ParentID       *uint   `json:"parent_id"`
	CommissionRate float64 `json:"commission_rate" binding:"omitempty,gte=0,lte=100"`
}

// AdminCreate creates a new agent (admin).
// POST /admin/agents
func (h *AgentHandler) AdminCreate(c *gin.Context) {
	var req adminCreateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Check user exists
	var user model.User
	if err := h.db.First(&user, req.UserID).Error; err != nil {
		response.BadRequest(c, "user not found")
		return
	}

	// Check not already an agent
	var existing model.Agent
	if err := h.db.Where("user_id = ?", req.UserID).First(&existing).Error; err == nil {
		response.BadRequest(c, "user is already an agent")
		return
	}

	level := 1
	if req.ParentID != nil {
		var parent model.Agent
		if err := h.db.First(&parent, *req.ParentID).Error; err != nil {
			response.BadRequest(c, "parent agent not found")
			return
		}
		level = parent.Level + 1
	}

	commissionRate := 10.0
	if req.CommissionRate > 0 {
		commissionRate = req.CommissionRate
	}

	agent := model.Agent{
		UserID:         req.UserID,
		AgentNo:        generateAgentNo(),
		ParentID:       req.ParentID,
		Level:          level,
		CommissionRate: commissionRate,
		Status:         1,
	}

	if err := h.db.Create(&agent).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, agent)
}

// adminUpdateAgentRequest is the payload for AdminUpdate.
type adminUpdateAgentRequest struct {
	CommissionRate *float64 `json:"commission_rate" binding:"omitempty,gte=0,lte=100"`
	Status         *int8    `json:"status" binding:"omitempty,oneof=1 2"`
	ParentID       *uint    `json:"parent_id"`
}

// AdminUpdate updates an agent (admin).
// PUT /admin/agents/:id
func (h *AgentHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid agent id")
		return
	}

	var req adminUpdateAgentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var agent model.Agent
	if err := h.db.First(&agent, id).Error; err != nil {
		response.NotFound(c, "agent not found")
		return
	}

	updates := map[string]interface{}{}
	if req.CommissionRate != nil {
		updates["commission_rate"] = *req.CommissionRate
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}

	if len(updates) > 0 {
		if err := h.db.Model(&agent).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}

	h.db.Preload("User").Preload("Parent").First(&agent, id)
	response.Success(c, agent)
}

// AdminConfirmCommission confirms or rejects a commission record (admin).
// POST /admin/agents/commissions/:id/confirm
func (h *AgentHandler) AdminConfirmCommission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid commission id")
		return
	}

	var req struct {
		Status int8 `json:"status" binding:"required,oneof=2 3"` // 2=confirmed 3=rejected
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var commission model.AgentCommission
	if err := h.db.First(&commission, id).Error; err != nil {
		response.NotFound(c, "commission not found")
		return
	}

	if commission.Status != 1 {
		response.BadRequest(c, "commission already processed")
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&commission).Update("status", req.Status).Error; err != nil {
			return err
		}
		// Credit agent balance on confirmation
		if req.Status == 2 {
			if err := tx.Model(&model.Agent{}).Where("id = ?", commission.AgentID).
				Updates(map[string]interface{}{
					"balance":      gorm.Expr("balance + ?", commission.Amount),
					"total_earned": gorm.Expr("total_earned + ?", commission.Amount),
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "commission confirmed")
}

// generateAgentNo generates a unique agent number.
func generateAgentNo() string {
	// Simple timestamp-based generator; replace with a proper sequence in production.
	return "AG" + strconv.FormatInt(timeNow().UnixMicro(), 36)
}

func timeNow() *time.Time {
	t := time.Now()
	return &t
}
