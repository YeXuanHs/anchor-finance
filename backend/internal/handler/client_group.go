package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// ClientGroupHandler handles client group HTTP requests.
type ClientGroupHandler struct {
	svc *service.ClientGroupService
	log *logger.Logger
}

// NewClientGroupHandler creates a new ClientGroupHandler.
func NewClientGroupHandler(svc *service.ClientGroupService, log *logger.Logger) *ClientGroupHandler {
	return &ClientGroupHandler{svc: svc, log: log}
}

// List returns a paginated list of client groups.
// GET /client-groups
func (h *ClientGroupHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	groups, total, err := h.svc.GetList(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, groups, total, page, pageSize)
}

// Get returns a single client group by ID.
// GET /client-groups/:id
func (h *ClientGroupHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	group, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "group not found")
		return
	}
	response.Success(c, group)
}

// Create creates a new client group.
// POST /client-groups
func (h *ClientGroupHandler) Create(c *gin.Context) {
	var req service.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	group, err := h.svc.Create(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, group)
}

// Update modifies a client group.
// PUT /client-groups/:id
func (h *ClientGroupHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	var req service.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Update(uint(id), req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "group updated")
}

// Delete removes a client group.
// DELETE /client-groups/:id
func (h *ClientGroupHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "group deleted")
}

// GetMembers returns all members of a client group.
// GET /client-groups/:id/members
func (h *ClientGroupHandler) GetMembers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	members, err := h.svc.GetMembers(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, members)
}

// AddMemberRequest is the payload for adding a group member.
type AddMemberRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// AddMember adds a user to a client group.
// POST /client-groups/:id/members
func (h *ClientGroupHandler) AddMember(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	var req AddMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AddMember(uint(id), req.UserID); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "member added")
}

// RemoveMember removes a user from a client group.
// DELETE /client-groups/:id/members/:user_id
func (h *ClientGroupHandler) RemoveMember(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	if err := h.svc.RemoveMember(uint(id), uint(userID)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "member removed")
}

// SetMembersRequest is the payload for batch-setting group members.
type SetMembersRequest struct {
	UserIDs []uint `json:"user_ids" binding:"required"`
}

// SetMembers replaces all members of a group.
// PUT /client-groups/:id/members
func (h *ClientGroupHandler) SetMembers(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}

	var req SetMembersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.SetMembers(uint(id), req.UserIDs); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "members updated")
}
