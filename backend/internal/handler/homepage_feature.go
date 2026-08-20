package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type HomepageFeatureHandler struct {
	svc *service.HomepageFeatureService
	log *logger.Logger
}

func NewHomepageFeatureHandler(svc *service.HomepageFeatureService, log *logger.Logger) *HomepageFeatureHandler {
	return &HomepageFeatureHandler{svc: svc, log: log}
}

// List returns paginated feature list for admin.
func (h *HomepageFeatureHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	position := c.Query("position")
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))

	items, total, err := h.svc.GetList(page, pageSize, position, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetDetail returns a single feature by ID.
func (h *HomepageFeatureHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid feature id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "feature not found")
		return
	}
	response.Success(c, item)
}

// Create creates a new feature.
func (h *HomepageFeatureHandler) Create(c *gin.Context) {
	var req service.CreateFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Create(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// Update updates an existing feature.
func (h *HomepageFeatureHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid feature id")
		return
	}

	var req service.UpdateFeatureRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Update(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// Delete deletes a feature.
func (h *HomepageFeatureHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid feature id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "feature deleted")
}

// SetStatus toggles a feature's enabled/disabled status.
func (h *HomepageFeatureHandler) SetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid feature id")
		return
	}

	if err := h.svc.ToggleStatus(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "feature status toggled")
}

// GetActive returns active features for frontend display.
func (h *HomepageFeatureHandler) GetActive(c *gin.Context) {
	position := c.DefaultQuery("position", "home")

	items, err := h.svc.GetActive(position)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}
