package handler

import (
	"strconv"

	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DDoSHandler handles DDoS management HTTP requests for users.
type DDoSHandler struct {
	db *gorm.DB
}

// NewDDoSHandler creates a new DDoSHandler.
func NewDDoSHandler(db *gorm.DB) *DDoSHandler {
	return &DDoSHandler{db: db}
}

// GetIPs returns the user's DDoS protected IPs.
// GET /user/ddos/ips
func (h *DDoSHandler) GetIPs(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var items []map[string]interface{}
	var total int64

	query := h.db.Table("ddos_ips").Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items)

	response.SuccessPage(c, items, total, page, pageSize)
}

// PostIP adds a new DDoS protected IP.
// POST /user/ddos/ips
func (h *DDoSHandler) PostIP(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		IP          string `json:"ip" binding:"required"`
		Remark      string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item := map[string]interface{}{
		"user_id": userID,
		"ip":      req.IP,
		"remark":  req.Remark,
		"status":  1,
	}
	if err := h.db.Table("ddos_ips").Create(&item).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, item)
}

// DeleteIP removes a DDoS protected IP.
// DELETE /user/ddos/ips/:id
func (h *DDoSHandler) DeleteIP(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ip id")
		return
	}

	result := h.db.Table("ddos_ips").Where("id = ? AND user_id = ?", id, userID).Delete(nil)
	if result.RowsAffected == 0 {
		response.NotFound(c, "ip not found")
		return
	}

	response.SuccessMsg(c, "ip deleted")
}

// GetIPRules returns DDoS rules for a specific IP.
// GET /user/ddos/ips/:id/rules
func (h *DDoSHandler) GetIPRules(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ip id")
		return
	}

	// Verify IP belongs to user
	var ipRecord map[string]interface{}
	if err := h.db.Table("ddos_ips").Where("id = ? AND user_id = ?", id, userID).First(&ipRecord).Error; err != nil {
		response.NotFound(c, "ip not found")
		return
	}

	var rules []map[string]interface{}
	h.db.Table("ddos_rules").Where("ddos_ip_id = ?", id).Order("id DESC").Find(&rules)

	response.Success(c, rules)
}

// PostIPToggle toggles DDoS protection for an IP.
// POST /user/ddos/ips/:id/toggle
func (h *DDoSHandler) PostIPToggle(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid ip id")
		return
	}

	// Verify IP belongs to user
	var ipRecord map[string]interface{}
	if err := h.db.Table("ddos_ips").Where("id = ? AND user_id = ?", id, userID).First(&ipRecord).Error; err != nil {
		response.NotFound(c, "ip not found")
		return
	}

	// Toggle status
	currentStatus, _ := ipRecord["status"].(int64)
	newStatus := 1
	if currentStatus == 1 {
		newStatus = 0
	}

	h.db.Table("ddos_ips").Where("id = ?", id).Update("status", newStatus)

	response.Success(c, gin.H{"id": id, "status": newStatus})
}

// DeleteRule removes a DDoS rule.
// DELETE /user/ddos/rules/:id
func (h *DDoSHandler) DeleteRule(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	// Verify rule belongs to user's IP
	result := h.db.Table("ddos_rules").
		Where("id = ? AND ddos_ip_id IN (SELECT id FROM ddos_ips WHERE user_id = ?)", id, userID).
		Delete(nil)
	if result.RowsAffected == 0 {
		response.NotFound(c, "rule not found")
		return
	}

	response.SuccessMsg(c, "rule deleted")
}

// PutRule updates a DDoS rule.
// PUT /user/ddos/rules/:id
func (h *DDoSHandler) PutRule(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid rule id")
		return
	}

	// Verify rule belongs to user's IP
	var rule map[string]interface{}
	if err := h.db.Table("ddos_rules").
		Where("id = ? AND ddos_ip_id IN (SELECT id FROM ddos_ips WHERE user_id = ?)", id, userID).
		First(&rule).Error; err != nil {
		response.NotFound(c, "rule not found")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Protocol    string `json:"protocol"`
		Port        string `json:"port"`
		Action      string `json:"action"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Protocol != "" {
		updates["protocol"] = req.Protocol
	}
	if req.Port != "" {
		updates["port"] = req.Port
	}
	if req.Action != "" {
		updates["action"] = req.Action
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if len(updates) > 0 {
		h.db.Table("ddos_rules").Where("id = ?", id).Updates(updates)
	}

	response.SuccessMsg(c, "rule updated")
}

// GetTraffic returns DDoS traffic statistics.
// GET /user/ddos/traffic
func (h *DDoSHandler) GetTraffic(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Get user's IPs
	var ips []map[string]interface{}
	h.db.Table("ddos_ips").Where("user_id = ?", userID).Find(&ips)

	// Aggregate traffic data
	trafficData := map[string]interface{}{
		"total_attack_count":   0,
		"total_attack_traffic": 0,
		"current_attack":       false,
		"protected_ips":        len(ips),
	}

	response.Success(c, trafficData)
}

// GetOverview returns DDoS protection overview.
// GET /user/ddos/overview
func (h *DDoSHandler) GetOverview(c *gin.Context) {
	userID := c.GetUint("user_id")

	var ipCount int64
	h.db.Table("ddos_ips").Where("user_id = ?", userID).Count(&ipCount)

	var ruleCount int64
	h.db.Table("ddos_rules").
		Where("ddos_ip_id IN (SELECT id FROM ddos_ips WHERE user_id = ?)", userID).
		Count(&ruleCount)

	overview := map[string]interface{}{
		"protected_ips":  ipCount,
		"total_rules":    ruleCount,
		"protection_on":  true,
	}

	response.Success(c, overview)
}
