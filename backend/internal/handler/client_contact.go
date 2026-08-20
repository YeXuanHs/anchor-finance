package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// ClientContactHandler handles client contact HTTP requests.
type ClientContactHandler struct {
	svc *service.ClientContactService
	log *logger.Logger
}

// NewClientContactHandler creates a new ClientContactHandler.
func NewClientContactHandler(svc *service.ClientContactService, log *logger.Logger) *ClientContactHandler {
	return &ClientContactHandler{svc: svc, log: log}
}

// List returns a paginated list of contacts (admin).
// GET /client-contacts
func (h *ClientContactHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var userID uint
	if uid := c.Query("user_id"); uid != "" {
		v, _ := strconv.ParseUint(uid, 10, 64)
		userID = uint(v)
	}

	items, total, err := h.svc.GetList(page, pageSize, userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// Get returns a single contact by ID.
// GET /client-contacts/:id
func (h *ClientContactHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contact id")
		return
	}

	contact, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "contact not found")
		return
	}
	response.Success(c, contact)
}

// Create adds a new contact.
// POST /client-contacts
func (h *ClientContactHandler) Create(c *gin.Context) {
	var req service.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	contact, err := h.svc.Create(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, contact)
}

// Update modifies an existing contact.
// PUT /client-contacts/:id
func (h *ClientContactHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contact id")
		return
	}

	var req service.UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Update(uint(id), req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "contact updated")
}

// Delete removes a contact.
// DELETE /client-contacts/:id
func (h *ClientContactHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contact id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "contact deleted")
}

// MyContacts returns contacts for the authenticated user.
// GET /user/contacts
func (h *ClientContactHandler) MyContacts(c *gin.Context) {
	userID := c.GetUint("user_id")

	contacts, err := h.svc.GetByUser(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, contacts)
}

// MyContactCreate adds a contact for the authenticated user.
// POST /user/contacts
func (h *ClientContactHandler) MyContactCreate(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	req.UserID = userID

	contact, err := h.svc.Create(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, contact)
}

// MyContactUpdate updates a contact belonging to the authenticated user.
// PUT /user/contacts/:id
func (h *ClientContactHandler) MyContactUpdate(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contact id")
		return
	}

	contact, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "contact not found")
		return
	}
	if contact.UserID != userID {
		response.NotFound(c, "contact not found")
		return
	}

	var req service.UpdateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.Update(uint(id), req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "contact updated")
}

// MyContactDelete deletes a contact belonging to the authenticated user.
// DELETE /user/contacts/:id
func (h *ClientContactHandler) MyContactDelete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contact id")
		return
	}

	contact, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "contact not found")
		return
	}
	if contact.UserID != userID {
		response.NotFound(c, "contact not found")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "contact deleted")
}
