package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type TicketPrereplyHandler struct {
	svc *service.TicketPrereplyService
	log *logger.Logger
}

func NewTicketPrereplyHandler(svc *service.TicketPrereplyService, log *logger.Logger) *TicketPrereplyHandler {
	return &TicketPrereplyHandler{svc: svc, log: log}
}

// GetReplyList returns a list of prereply categories with replies.
func (h *TicketPrereplyHandler) GetReplyList(c *gin.Context) {
	categories, err := h.svc.GetReplyList()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"prereply": categories})
}

// AddCategory adds a new prereply category.
func (h *TicketPrereplyHandler) AddCategory(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cat := &model.TicketPrereplyCategory{
		Name:   req.Name,
		Status: 1,
	}

	if err := h.svc.AddCategory(cat); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cat)
}

// UpdateCategory updates an existing prereply category.
func (h *TicketPrereplyHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	var req struct {
		Name      string `json:"name"`
		Status    *int16 `json:"status"`
		SortOrder *int   `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if err := h.svc.UpdateCategory(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "category updated")
}

// DeleteCategory deletes a prereply category.
func (h *TicketPrereplyHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	if err := h.svc.DeleteCategory(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "category deleted")
}

// AddReply adds a new prereply.
func (h *TicketPrereplyHandler) AddReply(c *gin.Context) {
	var req struct {
		CategoryID uint   `json:"category_id" binding:"required"`
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	reply := &model.TicketPrereply{
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Content:    req.Content,
		Status:     1,
	}

	if err := h.svc.AddReply(reply); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, reply)
}

// UpdateReply updates an existing prereply.
func (h *TicketPrereplyHandler) UpdateReply(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid reply id")
		return
	}

	var req struct {
		Title     string `json:"title"`
		Content   string `json:"content"`
		Status    *int16 `json:"status"`
		SortOrder *int   `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if err := h.svc.UpdateReply(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "reply updated")
}

// DeleteReply deletes a prereply.
func (h *TicketPrereplyHandler) DeleteReply(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid reply id")
		return
	}

	if err := h.svc.DeleteReply(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "reply deleted")
}

// GetCategoryDetail returns a single prereply category by ID.
// GET /admin/ticket-prereply-categories/:id
func (h *TicketPrereplyHandler) GetCategoryDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	category, err := h.svc.GetCategoryByID(uint(id))
	if err != nil {
		response.NotFound(c, "category not found")
		return
	}
	response.Success(c, category)
}

// GetPrereplyCategories returns all categories for dropdown selection.
// GET /admin/ticket-prereply-categories/options
func (h *TicketPrereplyHandler) GetPrereplyCategories(c *gin.Context) {
	categories, err := h.svc.GetAllCategories()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, categories)
}

// GetPrereplyDetail returns a single prereply by ID with its categories.
// GET /admin/ticket-prereplies/:id
func (h *TicketPrereplyHandler) GetPrereplyDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid reply id")
		return
	}

	reply, err := h.svc.GetReplyByID(uint(id))
	if err != nil {
		response.NotFound(c, "reply not found")
		return
	}

	categories, _ := h.svc.GetAllCategories()
	response.Success(c, gin.H{
		"list":       reply,
		"categories": categories,
	})
}

// SearchPrereply searches prereplies by title and content.
// POST /admin/ticket-prereplies/search
func (h *TicketPrereplyHandler) SearchPrereply(c *gin.Context) {
	var req struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	replies, err := h.svc.SearchReplies(req.Title, req.Content)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, replies)
}
