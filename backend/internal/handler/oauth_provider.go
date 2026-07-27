package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type OAuthProviderHandler struct {
	svc *service.OAuthProviderService
	log *logger.Logger
}

func NewOAuthProviderHandler(svc *service.OAuthProviderService, log *logger.Logger) *OAuthProviderHandler {
	return &OAuthProviderHandler{svc: svc, log: log}
}

// AdminGetList returns paginated OAuth provider list.
func (h *OAuthProviderHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "-1"))

	items, total, err := h.svc.GetList(page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// AdminGetDetail returns a single OAuth provider by ID.
func (h *OAuthProviderHandler) AdminGetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "oauth provider not found")
		return
	}
	response.Success(c, item)
}

// AdminCreate creates a new OAuth provider.
func (h *OAuthProviderHandler) AdminCreate(c *gin.Context) {
	var req service.CreateOAuthProviderRequest
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

// AdminUpdate updates an OAuth provider.
func (h *OAuthProviderHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	var req service.UpdateOAuthProviderRequest
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

// AdminDelete deletes an OAuth provider.
func (h *OAuthProviderHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "oauth provider deleted")
}

// AdminToggleStatus toggles a provider's status.
func (h *OAuthProviderHandler) AdminToggleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	if err := h.svc.ToggleStatus(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "oauth provider status toggled")
}

// AdminToggleEnabled toggles a provider's enabled flag.
func (h *OAuthProviderHandler) AdminToggleEnabled(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	if err := h.svc.ToggleEnabled(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "oauth provider enabled toggled")
}

// AdminUpdateSort updates sort order for an OAuth provider.
func (h *OAuthProviderHandler) AdminUpdateSort(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid provider id")
		return
	}

	var req struct {
		SortOrder int `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateSort(uint(id), req.SortOrder); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "sort order updated")
}

// GetEnabled returns enabled OAuth providers for frontend.
func (h *OAuthProviderHandler) GetEnabled(c *gin.Context) {
	items, err := h.svc.GetEnabled()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}
