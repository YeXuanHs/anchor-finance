package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productSvc *service.ProductService
	log        *logger.Logger
}

func NewProductHandler(productSvc *service.ProductService, log *logger.Logger) *ProductHandler {
	return &ProductHandler{productSvc: productSvc, log: log}
}

// GetList returns active products.
func (h *ProductHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")

	products, total, err := h.productSvc.GetList(page, pageSize, category)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, products, total, page, pageSize)
}

// GetDetail returns a single product.
func (h *ProductHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	product, err := h.productSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "product not found")
		return
	}
	response.Success(c, product)
}

// GetHot returns the top selling products.
func (h *ProductHandler) GetHot(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "4"))
	if limit < 1 || limit > 20 {
		limit = 4
	}
	products, err := h.productSvc.GetHotProducts(limit)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, products)
}

// Create adds a new product (admin).
func (h *ProductHandler) Create(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	product, err := h.productSvc.Create(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, product)
}

// Update modifies a product (admin).
func (h *ProductHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	var req service.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	product, err := h.productSvc.Update(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, product)
}

// Delete soft-deletes a product (admin).
func (h *ProductHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	if err := h.productSvc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "product deleted")
}

// GetUserProducts returns the authenticated user's products.
func (h *ProductHandler) GetUserProducts(c *gin.Context) {
	userID := c.GetUint("user_id")
	products, err := h.productSvc.GetUserProducts(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, products)
}

// GetGroups returns product categories/groups (admin).
func (h *ProductHandler) GetGroups(c *gin.Context) {
	groups, err := h.productSvc.GetGroups()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, groups)
}

// CreateGroup creates a product group (admin).
func (h *ProductHandler) CreateGroup(c *gin.Context) {
	var req service.CreateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	group, err := h.productSvc.CreateGroup(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, group)
}

// UpdateGroup updates a product group (admin).
func (h *ProductHandler) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}
	var req service.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	group, err := h.productSvc.UpdateGroup(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, group)
}

// DeleteGroup deletes a product group (admin).
func (h *ProductHandler) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}
	if err := h.productSvc.DeleteGroup(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "group deleted")
}

// ==================== First Group (一级分组) ====================

// GetFirstGroups returns all first-level groups.
func (h *ProductHandler) GetFirstGroups(c *gin.Context) {
	groups, err := h.productSvc.GetFirstGroups()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, groups)
}

// CreateFirstGroup creates a new first-level group.
func (h *ProductHandler) CreateFirstGroup(c *gin.Context) {
	var req service.CreateFirstGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	group, err := h.productSvc.CreateFirstGroup(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, group)
}

// UpdateFirstGroup updates a first-level group.
func (h *ProductHandler) UpdateFirstGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}
	var req service.UpdateFirstGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.productSvc.UpdateFirstGroup(uint(id), req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "first group updated")
}

// DeleteFirstGroup deletes a first-level group.
func (h *ProductHandler) DeleteFirstGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid group id")
		return
	}
	if err := h.productSvc.DeleteFirstGroup(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "first group deleted")
}

// ==================== Sort Management ====================

// UpdateFirstGroupSort batch-updates first-level group sort orders.
func (h *ProductHandler) UpdateFirstGroupSort(c *gin.Context) {
	var req service.BatchSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.productSvc.UpdateFirstGroupSort(req.Items); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "first group sort updated")
}

// UpdateGroupSort batch-updates second-level group sort orders.
func (h *ProductHandler) UpdateGroupSort(c *gin.Context) {
	var req service.BatchSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.productSvc.UpdateGroupSort(req.Items); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "group sort updated")
}

// UpdateProductSort batch-updates product sort orders.
func (h *ProductHandler) UpdateProductSort(c *gin.Context) {
	var req service.BatchSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.productSvc.UpdateProductSort(req.Items); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "product sort updated")
}

// ==================== Duplicate Product ====================

// DuplicateProduct copies a product and its configurations.
func (h *ProductHandler) DuplicateProduct(c *gin.Context) {
	var req service.DuplicateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	product, err := h.productSvc.DuplicateProduct(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, product)
}

// ==================== Edit Stock ====================

// EditStock updates a product's stock quantity.
func (h *ProductHandler) EditStock(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}
	var req service.EditStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.productSvc.EditStock(uint(id), req.Stock); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "stock updated")
}

// ==================== Batch Operations ====================

// BatchUpdate updates multiple products at once.
func (h *ProductHandler) BatchUpdate(c *gin.Context) {
	var req service.BatchUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.productSvc.BatchUpdate(req.IDs, req.Fields); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "products updated")
}

// BatchDelete soft-deletes multiple products.
func (h *ProductHandler) BatchDelete(c *gin.Context) {
	var req service.BatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.productSvc.BatchDelete(req.IDs); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "products deleted")
}

// ==================== Check Alias ====================

// CheckAlias checks if a group alias is already in use.
func (h *ProductHandler) CheckAlias(c *gin.Context) {
	alias := c.Query("alias")
	if alias == "" {
		response.BadRequest(c, "alias is required")
		return
	}
	excludeID, _ := strconv.ParseUint(c.Query("exclude_id"), 10, 64)
	available, err := h.productSvc.CheckAlias(alias, uint(excludeID))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"available": available})
}

// ==================== Delete Custom Field ====================

// DeleteCustomField deletes a custom field by ID.
func (h *ProductHandler) DeleteCustomField(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid field id")
		return
	}
	if err := h.productSvc.DeleteCustomField(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "custom field deleted")
}

// ==================== Download Management ====================

// ManageDownloads adds or removes download file associations.
func (h *ProductHandler) ManageDownloads(c *gin.Context) {
	var req service.ManageDownloadsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := h.productSvc.ManageDownloads(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"result": result})
}

// GetProductDownloads returns download files linked to a product.
func (h *ProductHandler) GetProductDownloads(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}
	downloads, err := h.productSvc.GetProductDownloads(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, downloads)
}

// AddDownloadFile creates a new download file and links it to a product.
func (h *ProductHandler) AddDownloadFile(c *gin.Context) {
	var req service.AddDownloadFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.productSvc.AddDownloadFile(req); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "download file added")
}

// ==================== Discount List ====================

// GetDiscountList returns pricing/discount records for a product.
func (h *ProductHandler) GetDiscountList(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}
	list, err := h.productSvc.GetDiscountList(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, list)
}

// ==================== Get Upstream Price ====================

// GetUpstreamPrice fetches upstream pricing for a product.
func (h *ProductHandler) GetUpstreamPrice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}
	data, err := h.productSvc.GetUpstreamPrice(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}
