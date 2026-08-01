package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type TicketDeptHandler struct {
	deptSvc *service.TicketDepartmentService
	log     *logger.Logger
}

func NewTicketDeptHandler(deptSvc *service.TicketDepartmentService, log *logger.Logger) *TicketDeptHandler {
	return &TicketDeptHandler{deptSvc: deptSvc, log: log}
}

// Create creates a new ticket department (admin).
func (h *TicketDeptHandler) Create(c *gin.Context) {
	var req service.CreateTicketDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	dept, err := h.deptSvc.Create(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, dept)
}

// GetDetail returns a single department.
func (h *TicketDeptHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid department id")
		return
	}

	dept, err := h.deptSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "department not found")
		return
	}
	response.Success(c, dept)
}

// GetList returns all departments (admin).
func (h *TicketDeptHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}

	depts, total, err := h.deptSvc.GetList(page, pageSize, status, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, depts, total, page, pageSize)
}

// GetTree returns departments as a tree structure.
func (h *TicketDeptHandler) GetTree(c *gin.Context) {
	tree, err := h.deptSvc.GetTree()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, tree)
}

// Update modifies a department (admin).
func (h *TicketDeptHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid department id")
		return
	}

	var req service.UpdateTicketDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	dept, err := h.deptSvc.Update(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, dept)
}

// Delete soft-deletes a department (admin).
func (h *TicketDeptHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid department id")
		return
	}

	if err := h.deptSvc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "department deleted")
}

// Enable activates a department (admin).
func (h *TicketDeptHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid department id")
		return
	}

	if err := h.deptSvc.Enable(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "department enabled")
}

// Disable deactivates a department (admin).
func (h *TicketDeptHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid department id")
		return
	}

	if err := h.deptSvc.Disable(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "department disabled")
}

// AddMember adds a member to the department (admin).
func (h *TicketDeptHandler) AddMember(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid department id")
		return
	}

	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.deptSvc.AddMember(uint(id), req.UserID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "member added")
}

// RemoveMember removes a member from the department (admin).
func (h *TicketDeptHandler) RemoveMember(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid department id")
		return
	}
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.deptSvc.RemoveMember(uint(deptID), uint(userID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "member removed")
}

// SetManagers sets the managers for a department (admin).
func (h *TicketDeptHandler) SetManagers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid department id")
		return
	}

	var req struct {
		ManagerIDs []uint `json:"manager_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.deptSvc.SetManagers(uint(id), req.ManagerIDs); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "managers updated")
}

// GetMembers returns the members of a department (admin).
func (h *TicketDeptHandler) GetMembers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid department id")
		return
	}

	members, err := h.deptSvc.GetMembers(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, members)
}

// GetCustomParamType returns the available custom field types.
// GET /ticket-departments/custom-param-types
func (h *TicketDeptHandler) GetCustomParamType(c *gin.Context) {
	types := map[string]string{
		"dropdown":  "下拉",
		"password":  "密码",
		"text":      "文本框",
		"tickbox":   "选项框",
		"textarea":  "文本域",
	}
	response.Success(c, types)
}

// AddTicketCustomParam adds a custom field to a ticket department.
// POST /ticket-departments/:id/custom-params
func (h *TicketDeptHandler) AddTicketCustomParam(c *gin.Context) {
	deptID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid department id")
		return
	}

	var req struct {
		FieldName   string `json:"field_name" binding:"required"`
		FieldType   string `json:"field_type" binding:"required"`
		Description string `json:"description" binding:"required"`
		FieldOption string `json:"field_option"`
		RegExpr     string `json:"reg_expr"`
		AdminOnly   int8   `json:"admin_only"`
		Required    int8   `json:"required"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.FieldType == "dropdown" && req.FieldOption == "" {
		response.BadRequest(c, "dropdown type requires field_option")
		return
	}

	field, err := h.deptSvc.AddCustomParam(uint(deptID), req.FieldName, req.FieldType, req.Description, req.FieldOption, req.RegExpr, req.AdminOnly, req.Required, req.SortOrder)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, field)
}

// GetTicketParamVal returns a custom field's value.
// GET /ticket-departments/custom-params/:field_id
func (h *TicketDeptHandler) GetTicketParamVal(c *gin.Context) {
	fieldID, err := strconv.ParseUint(c.Param("field_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid field id")
		return
	}

	field, err := h.deptSvc.GetCustomParam(uint(fieldID))
	if err != nil {
		response.NotFound(c, "field not found")
		return
	}
	response.Success(c, field)
}

// EditTicketCustomParam updates a custom field.
// PUT /ticket-departments/custom-params/:field_id
func (h *TicketDeptHandler) EditTicketCustomParam(c *gin.Context) {
	fieldID, err := strconv.ParseUint(c.Param("field_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid field id")
		return
	}

	var req struct {
		FieldName   string `json:"field_name" binding:"required"`
		FieldType   string `json:"field_type" binding:"required"`
		Description string `json:"description" binding:"required"`
		DeptID      uint   `json:"dept_id" binding:"required"`
		FieldOption string `json:"field_option"`
		RegExpr     string `json:"reg_expr"`
		AdminOnly   int8   `json:"admin_only"`
		Required    int8   `json:"required"`
		SortOrder   int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.FieldType == "dropdown" && req.FieldOption == "" {
		response.BadRequest(c, "dropdown type requires field_option")
		return
	}

	if err := h.deptSvc.EditCustomParam(uint(fieldID), req.DeptID, req.FieldName, req.FieldType, req.Description, req.FieldOption, req.RegExpr, req.AdminOnly, req.Required, req.SortOrder); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "custom param updated")
}

// DelTicketCustomParam deletes a custom field.
// DELETE /ticket-departments/custom-params/:field_id
func (h *TicketDeptHandler) DelTicketCustomParam(c *gin.Context) {
	fieldID, err := strconv.ParseUint(c.Param("field_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid field id")
		return
	}

	if err := h.deptSvc.DeleteCustomParam(uint(fieldID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "custom param deleted")
}

// MoveUp moves a department up in sort order.
// POST /ticket-departments/:id/move-up
func (h *TicketDeptHandler) MoveUp(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid department id")
		return
	}

	if err := h.deptSvc.MoveUp(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "moved up")
}

// MoveDown moves a department down in sort order.
// POST /ticket-departments/:id/move-down
func (h *TicketDeptHandler) MoveDown(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid department id")
		return
	}

	if err := h.deptSvc.MoveDown(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "moved down")
}
