package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type CustomFieldHandler struct {
	svc *service.CustomFieldService
	log *logger.Logger
}

func NewCustomFieldHandler(svc *service.CustomFieldService, log *logger.Logger) *CustomFieldHandler {
	return &CustomFieldHandler{svc: svc, log: log}
}

// ────────────────────────── Field Management ──────────────────────────

// GetFields returns fields filtered by group.
func (h *CustomFieldHandler) GetFields(c *gin.Context) {
	group := c.Query("group")
	fields, err := h.svc.GetFields(group)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, fields)
}

// GetFieldDetail returns a single field.
func (h *CustomFieldHandler) GetFieldDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid field id")
		return
	}
	field, err := h.svc.GetFieldByID(uint(id))
	if err != nil {
		response.NotFound(c, "field not found")
		return
	}
	response.Success(c, field)
}

// CreateField creates a new custom field.
func (h *CustomFieldHandler) CreateField(c *gin.Context) {
	var req service.CreateFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	field, err := h.svc.CreateField(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, field)
}

// UpdateField updates an existing field.
func (h *CustomFieldHandler) UpdateField(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid field id")
		return
	}
	var req service.UpdateFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	field, err := h.svc.UpdateField(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, field)
}

// DeleteField deletes a field and its values.
func (h *CustomFieldHandler) DeleteField(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid field id")
		return
	}
	if err := h.svc.DeleteField(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "field deleted")
}

// ReorderFields reorders fields by given IDs.
func (h *CustomFieldHandler) ReorderFields(c *gin.Context) {
	var req service.ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.ReorderFields(req.IDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "fields reordered")
}

// ────────────────────────── Field Values ──────────────────────────

// GetValues returns values for an owner.
func (h *CustomFieldHandler) GetValues(c *gin.Context) {
	ownerID, err := strconv.ParseUint(c.Query("owner_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid owner_id")
		return
	}
	ownerType := c.Query("owner_type")
	if ownerType == "" {
		response.BadRequest(c, "owner_type is required")
		return
	}
	values, err := h.svc.GetValues(uint(ownerID), ownerType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, values)
}

// SaveValues saves field values for an owner.
func (h *CustomFieldHandler) SaveValues(c *gin.Context) {
	ownerID, err := strconv.ParseUint(c.Query("owner_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid owner_id")
		return
	}
	ownerType := c.Query("owner_type")
	if ownerType == "" {
		response.BadRequest(c, "owner_type is required")
		return
	}
	var req service.SaveValuesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.SaveValues(uint(ownerID), ownerType, req.Values); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "values saved")
}

// DeleteValues deletes all values for an owner.
func (h *CustomFieldHandler) DeleteValues(c *gin.Context) {
	ownerID, err := strconv.ParseUint(c.Query("owner_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid owner_id")
		return
	}
	ownerType := c.Query("owner_type")
	if ownerType == "" {
		response.BadRequest(c, "owner_type is required")
		return
	}
	if err := h.svc.DeleteValues(uint(ownerID), ownerType); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "values deleted")
}

// GetSingleValue returns a single field value.
func (h *CustomFieldHandler) GetSingleValue(c *gin.Context) {
	fieldID, err := strconv.ParseUint(c.Query("field_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid field_id")
		return
	}
	ownerID, err := strconv.ParseUint(c.Query("owner_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid owner_id")
		return
	}
	ownerType := c.Query("owner_type")
	if ownerType == "" {
		response.BadRequest(c, "owner_type is required")
		return
	}
	val, err := h.svc.GetValue(uint(fieldID), uint(ownerID), ownerType)
	if err != nil {
		response.NotFound(c, "value not found")
		return
	}
	response.Success(c, val)
}

// ────────────────────────── Group Management ──────────────────────────

// GetGroups returns groups filtered by type.
func (h *CustomFieldHandler) GetGroups(c *gin.Context) {
	typ := c.Query("type")
	groups, err := h.svc.GetGroups(typ)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, groups)
}

// CreateGroup creates a new field group.
func (h *CustomFieldHandler) CreateGroup(c *gin.Context) {
	var req service.CreateCustomFieldGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	group, err := h.svc.CreateGroup(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, group)
}

// UpdateGroup updates an existing group.
func (h *CustomFieldHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}
	var req service.UpdateCustomFieldGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	group, err := h.svc.UpdateGroup(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, group)
}

// DeleteGroup deletes a field group.
func (h *CustomFieldHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}
	if err := h.svc.DeleteGroup(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "group deleted")
}

// ────────────────────────── Validation ──────────────────────────

// ValidateFields validates values against field definitions.
func (h *CustomFieldHandler) ValidateFields(c *gin.Context) {
	var req service.ValidateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	errs, err := h.svc.ValidateFields(req.Group, req.Values)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	if len(errs) > 0 {
		response.Error(c, 422, 422, "validation failed")
		return
	}
	response.SuccessMsg(c, "validation passed")
}

// GetCartCustomFields returns cart-specific fields.
func (h *CustomFieldHandler) GetCartCustomFields(c *gin.Context) {
	fields, err := h.svc.GetCartCustomFields()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, fields)
}

// GetProductCustomFields returns product fields with values.
func (h *CustomFieldHandler) GetProductCustomFields(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product_id")
		return
	}
	fields, values, err := h.svc.GetProductCustomFields(uint(productID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"fields": fields, "values": values})
}

// GetClientCustomFields returns client fields.
func (h *CustomFieldHandler) GetClientCustomFields(c *gin.Context) {
	fields, err := h.svc.GetClientCustomFields()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, fields)
}

// GetHostCustomFields returns host fields with values.
func (h *CustomFieldHandler) GetHostCustomFields(c *gin.Context) {
	hostID, err := strconv.ParseUint(c.Param("host_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid host_id")
		return
	}
	fields, values, err := h.svc.GetHostCustomFields(uint(hostID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"fields": fields, "values": values})
}

// ────────────────────────── Bulk Operations ──────────────────────────

// CopyFields copies fields between groups.
func (h *CustomFieldHandler) CopyFields(c *gin.Context) {
	var req service.CopyFieldsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.CopyFields(req.FromGroup, req.ToGroup); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "fields copied")
}

// ImportFields imports fields from JSON.
func (h *CustomFieldHandler) ImportFields(c *gin.Context) {
	var req service.ImportFieldsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	count, err := h.svc.ImportFields(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"imported": count})
}

// ExportFields exports fields of a group.
func (h *CustomFieldHandler) ExportFields(c *gin.Context) {
	group := c.Query("group")
	fields, err := h.svc.ExportFields(group)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, fields)
}
