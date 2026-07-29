package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type LinkCauseHandler struct {
	svc *service.LinkCauseService
	log *logger.Logger
}

func NewLinkCauseHandler(svc *service.LinkCauseService, log *logger.Logger) *LinkCauseHandler {
	return &LinkCauseHandler{svc: svc, log: log}
}

// GetCauses returns a list of link causes.
func (h *LinkCauseHandler) GetCauses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	linkType := c.Query("type")
	keyword := c.Query("keyword")

	causes, total, err := h.svc.GetList(page, pageSize, linkType, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, causes, total, page, pageSize)
}

// GetCause returns a single link cause.
func (h *LinkCauseHandler) GetCause(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cause id")
		return
	}

	cause, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "cause not found")
		return
	}
	response.Success(c, cause)
}

// GetTree returns link causes as a tree structure.
func (h *LinkCauseHandler) GetTree(c *gin.Context) {
	linkType := c.Query("type")

	tree, err := h.svc.GetTree(linkType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, tree)
}

// CreateCause creates a new link cause.
func (h *LinkCauseHandler) CreateCause(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		ParentID *uint  `json:"parent_id"`
		LinkType string `json:"link_type" binding:"required"`
		Level    int    `json:"level"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cause := &model.LinkCause{
		Name:     req.Name,
		ParentID: req.ParentID,
		LinkType: req.LinkType,
		Level:    req.Level,
		Status:   1,
	}

	if err := h.svc.Create(cause); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cause)
}

// UpdateCause updates an existing link cause.
func (h *LinkCauseHandler) UpdateCause(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cause id")
		return
	}

	var req struct {
		Name     string `json:"name"`
		ParentID *uint  `json:"parent_id"`
		LinkType string `json:"link_type"`
		Level    *int   `json:"level"`
		Status   *int16 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.LinkType != "" {
		updates["link_type"] = req.LinkType
	}
	if req.Level != nil {
		updates["level"] = *req.Level
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := h.svc.Update(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cause updated")
}

// DeleteCause deletes a link cause.
func (h *LinkCauseHandler) DeleteCause(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid cause id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cause deleted")
}
