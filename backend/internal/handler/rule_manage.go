package handler

import (
	"encoding/json"
	"strconv"

	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RuleManageHandler manages admin auth rules (navigation menus).
type RuleManageHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewRuleManageHandler creates a new RuleManageHandler.
func NewRuleManageHandler(db *gorm.DB, log *logger.Logger) *RuleManageHandler {
	return &RuleManageHandler{db: db, log: log}
}

// AuthRule represents an admin navigation rule.
type AuthRule struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Status    int        `gorm:"default:1" json:"status"`
	App       string     `gorm:"type:varchar(64);default:admin" json:"app"`
	Type      string     `gorm:"type:varchar(64);default:admin_url" json:"type"`
	Name      string     `gorm:"type:varchar(128)" json:"name"`
	Param     string     `gorm:"type:varchar(256)" json:"param"`
	Title     string     `gorm:"type:varchar(128);not null" json:"title"`
	Condition string     `gorm:"type:varchar(256)" json:"condition"`
	PID       uint       `gorm:"default:0;index" json:"pid"`
	URL       string     `gorm:"type:varchar(256)" json:"url"`
	IsDisplay int        `gorm:"default:0" json:"is_display"` // 0=接口, 1=前台页面
	Order     int        `gorm:"default:0" json:"order"`
	Son       []AuthRule `gorm:"-" json:"son,omitempty"`
}

// TableName returns the table name for AuthRule.
func (AuthRule) TableName() string {
	return "auth_rule"
}

// GetMenuList returns the auth rule menu list as a tree.
// GET /admin/rule-manage/menus
func (h *RuleManageHandler) GetMenuList(c *gin.Context) {
	var rules []AuthRule
	if err := h.db.Order("`order` ASC, id ASC").Find(&rules).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// Add cn_name for display
	for i := range rules {
		if rules[i].IsDisplay == 1 {
			rules[i].Name = rules[i].Title + " (前台页面)"
		} else {
			rules[i].Name = rules[i].Title + " (接口)"
		}
	}

	tree := buildAuthRuleTree(rules, 0)
	response.Success(c, tree)
}

// AddMenu adds a new auth rule menu item.
// POST /admin/rule-manage/menus
func (h *RuleManageHandler) AddMenu(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Name      string `json:"name"`
		URL       string `json:"url"`
		IsDisplay int    `json:"is_display"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Title == "" {
		response.BadRequest(c, "title不能为空")
		return
	}

	rule := AuthRule{
		Status:    1,
		App:       "admin",
		Type:      "admin_url",
		Name:      req.Name,
		Title:     req.Title,
		URL:       req.URL,
		IsDisplay: req.IsDisplay,
		PID:       0,
		Order:     1,
	}

	if err := h.db.Create(&rule).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, rule)
}

// EditMenu edits an existing auth rule menu item.
// PUT /admin/rule-manage/menus/:id
func (h *RuleManageHandler) EditMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid menu id")
		return
	}

	var req struct {
		Title     *string `json:"title"`
		Name      *string `json:"name"`
		URL       *string `json:"url"`
		IsDisplay *int    `json:"is_display"`
		Status    *int    `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		if *req.Title == "" {
			response.BadRequest(c, "title不能为空")
			return
		}
		updates["title"] = *req.Title
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.URL != nil {
		updates["url"] = *req.URL
	}
	if req.IsDisplay != nil {
		updates["is_display"] = *req.IsDisplay
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) > 0 {
		if err := h.db.Model(&AuthRule{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}

	response.SuccessMsg(c, "更新成功")
}

// SaveMenuList saves the entire auth rule menu list (tree structure).
// POST /admin/rule-manage/menus/save
func (h *RuleManageHandler) SaveMenuList(c *gin.Context) {
	var req struct {
		List []AuthRule `json:"list" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Flatten the tree and assign order
	flatList := flattenAuthRuleTree(req.List, 0, 1)

	// Use transaction to replace all rules
	tx := h.db.Begin()
	if tx.Error != nil {
		response.ServerError(c, tx.Error.Error())
		return
	}

	// Delete all existing rules
	if err := tx.Where("id > 0").Delete(&AuthRule{}).Error; err != nil {
		tx.Rollback()
		response.ServerError(c, err.Error())
		return
	}

	// Insert all rules
	if len(flatList) > 0 {
		if err := tx.Create(&flatList).Error; err != nil {
			tx.Rollback()
			response.ServerError(c, err.Error())
			return
		}
	}

	tx.Commit()
	response.SuccessMsg(c, "保存成功")
}

// DeleteMenu deletes an auth rule menu item.
// DELETE /admin/rule-manage/menus/:id
func (h *RuleManageHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid menu id")
		return
	}

	// Delete the rule and its children
	if err := h.db.Where("id = ? OR pid = ?", id, id).Delete(&AuthRule{}).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "删除成功")
}

// buildAuthRuleTree builds a tree structure from flat auth rules.
func buildAuthRuleTree(rules []AuthRule, parentID uint) []AuthRule {
	var tree []AuthRule
	for _, rule := range rules {
		if rule.PID == parentID {
			children := buildAuthRuleTree(rules, rule.ID)
			rule.Son = children
			tree = append(tree, rule)
		}
	}
	return tree
}

// flattenAuthRuleTree flattens a tree structure into a slice with order and pid assigned.
func flattenAuthRuleTree(tree []AuthRule, parentID uint, startOrder int) []AuthRule {
	var result []AuthRule
	order := startOrder
	for _, node := range tree {
		rule := AuthRule{
			ID:        node.ID,
			Status:    node.Status,
			App:       node.App,
			Type:      node.Type,
			Name:      node.Name,
			Param:     node.Param,
			Title:     node.Title,
			Condition: node.Condition,
			PID:       parentID,
			URL:       node.URL,
			IsDisplay: node.IsDisplay,
			Order:     order,
		}
		result = append(result, rule)
		order++

		if len(node.Son) > 0 {
			children := flattenAuthRuleTree(node.Son, node.ID, 1)
			result = append(result, children...)
		}
	}
	return result
}

// unmarshalJSON is a helper to unmarshal JSON strings.
func unmarshalJSON(s string, v interface{}) error {
	if s == "" {
		return nil
	}
	return json.Unmarshal([]byte(s), v)
}
