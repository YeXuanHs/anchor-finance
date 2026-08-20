package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// InterflowHandler handles product interflow/association HTTP requests.
type InterflowHandler struct {
	interflowSvc *service.InterflowService
	log          *logger.Logger
}

// NewInterflowHandler creates a new InterflowHandler.
func NewInterflowHandler(interflowSvc *service.InterflowService, log *logger.Logger) *InterflowHandler {
	return &InterflowHandler{interflowSvc: interflowSvc, log: log}
}

// Create creates a new product interflow relation.
// POST /interflows
func (h *InterflowHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req service.CreateInterflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	interflow, err := h.interflowSvc.Create(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, interflow)
}

// GetList returns paginated interflow relations for the authenticated user.
// GET /interflows
func (h *InterflowHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	relationType := c.Query("relation_type")

	interflows, total, err := h.interflowSvc.GetList(userID, page, pageSize, relationType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, interflows, total, page, pageSize)
}

// GetDetail returns a single interflow relation.
// GET /interflows/:id
func (h *InterflowHandler) GetDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid interflow id")
		return
	}

	interflow, err := h.interflowSvc.GetByID(userID, uint(id))
	if err != nil {
		response.NotFound(c, "interflow not found")
		return
	}
	response.Success(c, interflow)
}

// GetByProduct returns all interflow relations for a specific product.
// GET /interflows/product/:product_id
func (h *InterflowHandler) GetByProduct(c *gin.Context) {
	userID := c.GetUint("user_id")
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	interflows, err := h.interflowSvc.GetByProduct(userID, uint(productID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, interflows)
}

// Update updates an interflow relation.
// PUT /interflows/:id
func (h *InterflowHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid interflow id")
		return
	}

	var req service.UpdateInterflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	interflow, err := h.interflowSvc.Update(userID, uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, interflow)
}

// Delete deletes an interflow relation.
// DELETE /interflows/:id
func (h *InterflowHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid interflow id")
		return
	}

	if err := h.interflowSvc.Delete(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "interflow deleted")
}

// ToggleStatus toggles the status of an interflow relation.
// POST /interflows/:id/toggle
func (h *InterflowHandler) ToggleStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid interflow id")
		return
	}

	if err := h.interflowSvc.ToggleStatus(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "interflow status toggled")
}
