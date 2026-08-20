package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// ContactsHandler handles contact management HTTP requests.
type ContactsHandler struct {
	contactSvc *service.ContactService
	log        *logger.Logger
}

// NewContactsHandler creates a new ContactsHandler.
func NewContactsHandler(contactSvc *service.ContactService, log *logger.Logger) *ContactsHandler {
	return &ContactsHandler{contactSvc: contactSvc, log: log}
}

// Create creates a new contact.
// POST /contacts
func (h *ContactsHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req service.CreateUserContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	contact, err := h.contactSvc.Create(userID, req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, contact)
}

// GetList returns paginated contacts for the authenticated user.
// GET /contacts
func (h *ContactsHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	contacts, total, err := h.contactSvc.GetList(userID, page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, contacts, total, page, pageSize)
}

// GetDetail returns a single contact.
// GET /contacts/:id
func (h *ContactsHandler) GetDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contact id")
		return
	}

	contact, err := h.contactSvc.GetByID(userID, uint(id))
	if err != nil {
		response.NotFound(c, "contact not found")
		return
	}
	response.Success(c, contact)
}

// Update updates a contact.
// PUT /contacts/:id
func (h *ContactsHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contact id")
		return
	}

	var req service.UpdateUserContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	contact, err := h.contactSvc.Update(userID, uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, contact)
}

// Delete deletes a contact.
// DELETE /contacts/:id
func (h *ContactsHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contact id")
		return
	}

	if err := h.contactSvc.Delete(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "contact deleted")
}

// SetDefault marks a contact as default.
// POST /contacts/:id/default
func (h *ContactsHandler) SetDefault(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contact id")
		return
	}

	if err := h.contactSvc.SetDefault(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "default contact set")
}

// GetDefault returns the user's default contact.
// GET /contacts/default
func (h *ContactsHandler) GetDefault(c *gin.Context) {
	userID := c.GetUint("user_id")

	contact, err := h.contactSvc.GetDefault(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	if contact == nil {
		response.NotFound(c, "no default contact set")
		return
	}
	response.Success(c, contact)
}
