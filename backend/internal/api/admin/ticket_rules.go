package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetTicketRuleList 获取工单传递规则列表
// GET /api/admin/ticket-rules
func GetTicketRuleList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	var total int64
	db.Model(&model.TicketRule{}).Count(&total)

	var rules []model.TicketRule
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("priority ASC").Find(&rules)

	if rules == nil {
		rules = []model.TicketRule{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      rules,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateTicketRule 创建工单传递规则
// POST /api/admin/ticket-rules
func CreateTicketRule(c *gin.Context) {
	var req struct {
		Name           string `json:"name" binding:"required"`
		ConditionType  string `json:"condition_type" binding:"required"`
		ConditionValue string `json:"condition_value" binding:"required"`
		ActionType     string `json:"action_type" binding:"required"`
		ActionValue    string `json:"action_value" binding:"required"`
		Priority       int    `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	rule := model.TicketRule{
		Name:           req.Name,
		ConditionType:  req.ConditionType,
		ConditionValue: req.ConditionValue,
		ActionType:     req.ActionType,
		ActionValue:    req.ActionValue,
		Priority:       req.Priority,
		Status:         "active",
	}

	if err := db.Create(&rule).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": rule.ID,
		},
	})
}

// UpdateTicketRule 更新工单传递规则
// PUT /api/admin/ticket-rules/:id
func UpdateTicketRule(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name           string `json:"name"`
		ConditionType  string `json:"condition_type"`
		ConditionValue string `json:"condition_value"`
		ActionType     string `json:"action_type"`
		ActionValue    string `json:"action_value"`
		Priority       int    `json:"priority"`
		Status         string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var rule model.TicketRule
	if err := db.First(&rule, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "规则不存在", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.ConditionType != "" {
		updates["condition_type"] = req.ConditionType
	}
	if req.ConditionValue != "" {
		updates["condition_value"] = req.ConditionValue
	}
	if req.ActionType != "" {
		updates["action_type"] = req.ActionType
	}
	if req.ActionValue != "" {
		updates["action_value"] = req.ActionValue
	}
	if req.Priority > 0 {
		updates["priority"] = req.Priority
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&rule).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteTicketRule 删除工单传递规则
// DELETE /api/admin/ticket-rules/:id
func DeleteTicketRule(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var rule model.TicketRule
	if err := db.First(&rule, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "规则不存在", "data": nil})
		return
	}

	db.Delete(&rule)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
