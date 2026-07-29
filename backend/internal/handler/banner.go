package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type BannerHandler struct {
	svc *service.BannerService
	log *logger.Logger
}

func NewBannerHandler(svc *service.BannerService, log *logger.Logger) *BannerHandler {
	return &BannerHandler{svc: svc, log: log}
}

// AdminGetList returns paginated banner list for admin.
func (h *BannerHandler) AdminGetList(c *gin.Context) {
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

// AdminGetDetail returns a single banner by ID.
func (h *BannerHandler) AdminGetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid banner id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "banner not found")
		return
	}
	response.Success(c, item)
}

// AdminCreate creates a new banner.
func (h *BannerHandler) AdminCreate(c *gin.Context) {
	var req service.CreateBannerRequest
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

// AdminUpdate updates an existing banner.
func (h *BannerHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid banner id")
		return
	}

	var req service.UpdateBannerRequest
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

// AdminDelete deletes a banner.
func (h *BannerHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid banner id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "banner deleted")
}

// AdminToggleStatus toggles a banner's enabled/disabled status.
func (h *BannerHandler) AdminToggleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid banner id")
		return
	}

	if err := h.svc.ToggleStatus(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "banner status toggled")
}

// AdminUpdateSort updates sort order for a banner.
func (h *BannerHandler) AdminUpdateSort(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid banner id")
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

// GetActive returns active banners for frontend display.
func (h *BannerHandler) GetActive(c *gin.Context) {
	position := c.DefaultQuery("position", "home")

	items, err := h.svc.GetActive(position)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// List is an alias for AdminGetList.
func (h *BannerHandler) List(c *gin.Context) { h.AdminGetList(c) }

// GetDetail is an alias for AdminGetDetail.
func (h *BannerHandler) GetDetail(c *gin.Context) { h.AdminGetDetail(c) }

// Create is an alias for AdminCreate.
func (h *BannerHandler) Create(c *gin.Context) { h.AdminCreate(c) }

// Update is an alias for AdminUpdate.
func (h *BannerHandler) Update(c *gin.Context) { h.AdminUpdate(c) }

// Delete is an alias for AdminDelete.
func (h *BannerHandler) Delete(c *gin.Context) { h.AdminDelete(c) }

// SetStatus is an alias for AdminToggleStatus.
func (h *BannerHandler) SetStatus(c *gin.Context) { h.AdminToggleStatus(c) }

// Click records a banner click and redirects to the link.
func (h *BannerHandler) Click(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid banner id")
		return
	}

	banner, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "banner not found")
		return
	}

	_ = h.svc.IncrementClick(uint(id))

	if banner.LinkURL != "" {
		c.Redirect(302, banner.LinkURL)
		return
	}
	response.Success(c, gin.H{"id": banner.ID})
}
