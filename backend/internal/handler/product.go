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
