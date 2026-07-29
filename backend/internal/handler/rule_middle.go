package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type RuleMiddleHandler struct {
	svc *service.RuleMiddleService
	log *logger.Logger
}

func NewRuleMiddleHandler(svc *service.RuleMiddleService, log *logger.Logger) *RuleMiddleHandler {
	return &RuleMiddleHandler{svc: svc, log: log}
}

// GetMenuList returns the rule middle menu list.
func (h *RuleMiddleHandler) GetMenuList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	menus, total, err := h.svc.GetMenuList(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, menus, total, page, pageSize)
}

// AddMenu adds a new rule middle menu.
func (h *RuleMiddleHandler) AddMenu(c *gin.Context) {
	var req struct {
		Name         string `json:"name" binding:"required"`
		AddRole      int    `json:"add_role"`
		AddRoleMenu  string `json:"add_role_menu"`
		CatRole      int    `json:"cat_role"`
		CatRoleMenu  string `json:"cat_role_menu"`
		DelRole      int    `json:"del_role"`
		DelRoleMenu  string `json:"del_role_menu"`
		EditRole     int    `json:"edit_role"`
		EditRoleMenu string `json:"edit_role_menu"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	menu := &model.RuleMiddle{
		Name:         req.Name,
		AddRole:      req.AddRole,
		AddRoleMenu:  req.AddRoleMenu,
		CatRole:      req.CatRole,
		CatRoleMenu:  req.CatRoleMenu,
		DelRole:      req.DelRole,
		DelRoleMenu:  req.DelRoleMenu,
		EditRole:     req.EditRole,
		EditRoleMenu: req.EditRoleMenu,
		Status:       1,
	}

	if err := h.svc.Create(menu); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, menu)
}

// UpdateMenu updates an existing rule middle menu.
func (h *RuleMiddleHandler) UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid menu id")
		return
	}

	var req struct {
		Name         *string `json:"name"`
		AddRole      *int    `json:"add_role"`
		AddRoleMenu  *string `json:"add_role_menu"`
		CatRole      *int    `json:"cat_role"`
		CatRoleMenu  *string `json:"cat_role_menu"`
		DelRole      *int    `json:"del_role"`
		DelRoleMenu  *string `json:"del_role_menu"`
		EditRole     *int    `json:"edit_role"`
		EditRoleMenu *string `json:"edit_role_menu"`
		Status       *int16  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.AddRole != nil {
		updates["add_role"] = *req.AddRole
	}
	if req.AddRoleMenu != nil {
		updates["add_role_menu"] = *req.AddRoleMenu
	}
	if req.CatRole != nil {
		updates["cat_role"] = *req.CatRole
	}
	if req.CatRoleMenu != nil {
		updates["cat_role_menu"] = *req.CatRoleMenu
	}
	if req.DelRole != nil {
		updates["del_role"] = *req.DelRole
	}
	if req.DelRoleMenu != nil {
		updates["del_role_menu"] = *req.DelRoleMenu
	}
	if req.EditRole != nil {
		updates["edit_role"] = *req.EditRole
	}
	if req.EditRoleMenu != nil {
		updates["edit_role_menu"] = *req.EditRoleMenu
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := h.svc.Update(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "menu updated")
}

// DeleteMenu deletes a rule middle menu.
func (h *RuleMiddleHandler) DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid menu id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "menu deleted")
}
