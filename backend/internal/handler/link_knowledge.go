package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type LinkKnowledgeHandler struct {
	svc *service.LinkKnowledgeService
	log *logger.Logger
}

func NewLinkKnowledgeHandler(svc *service.LinkKnowledgeService, log *logger.Logger) *LinkKnowledgeHandler {
	return &LinkKnowledgeHandler{svc: svc, log: log}
}

// GetKnowledges returns a list of link knowledges.
func (h *LinkKnowledgeHandler) GetKnowledges(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	knowledgeType := c.Query("type")
	category := c.Query("category")
	keyword := c.Query("keyword")

	items, total, err := h.svc.GetList(page, pageSize, knowledgeType, category, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetKnowledge returns a single link knowledge.
func (h *LinkKnowledgeHandler) GetKnowledge(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid knowledge id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "knowledge not found")
		return
	}
	response.Success(c, item)
}

// CreateKnowledge creates a new link knowledge.
func (h *LinkKnowledgeHandler) CreateKnowledge(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		Content   string `json:"content"`
		LinkCause uint   `json:"link_cause"`
		Type      string `json:"type"`
		Category  string `json:"category"`
		Status    *int16 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	status := int16(1)
	if req.Status != nil {
		status = *req.Status
	}

	item := &model.LinkKnowledge{
		Title:     req.Title,
		Content:   req.Content,
		LinkCause: req.LinkCause,
		Type:      req.Type,
		Category:  req.Category,
		Status:    status,
	}

	if err := h.svc.Create(item); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, item)
}

// UpdateKnowledge updates an existing link knowledge.
func (h *LinkKnowledgeHandler) UpdateKnowledge(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid knowledge id")
		return
	}

	var req struct {
		Title     string `json:"title"`
		Content   string `json:"content"`
		LinkCause *uint  `json:"link_cause"`
		Type      string `json:"type"`
		Category  string `json:"category"`
		Status    *int16 `json:"status"`
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
	if req.LinkCause != nil {
		updates["link_cause"] = *req.LinkCause
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := h.svc.Update(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "knowledge updated")
}

// DeleteKnowledge deletes a link knowledge.
func (h *LinkKnowledgeHandler) DeleteKnowledge(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid knowledge id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "knowledge deleted")
}
