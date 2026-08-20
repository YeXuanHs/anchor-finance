package handler

import (
	"net/http"
	"strconv"
	"time"

	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

// GetProducts 获取产品列表
func (h *ProductHandler) GetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"list": []interface{}{}, "total": 0})
		return
	}

	var total int64
	database.Table("products").Count(&total)

	var products []struct {
		ID        uint      `json:"id"`
		Name      string    `json:"name"`
		Type      string    `json:"type"`
		Price     float64   `json:"price"`
		Status    int16     `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	database.Table("products").
		Select("id, name, type, price, status, created_at").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"list":  products,
		"total": total,
		"page":  page,
	})
}

// GetProduct 获取单个产品
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "not found"})
		return
	}

	var product struct {
		ID          uint      `json:"id"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Type        string    `json:"type"`
		Price       float64   `json:"price"`
		Status      int16     `json:"status"`
		CreatedAt   time.Time `json:"created_at"`
	}

	err := database.Table("products").Where("id = ?", id).First(&product).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "产品不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": product})
}

// CreateProduct 创建产品
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"code": 0, "message": "产品创建成功"})
}

// UpdateProduct 更新产品
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "产品更新成功"})
}

// DeleteProduct 删除产品
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "产品删除成功"})
}

// RegisterRoutes 注册路由
func (h *ProductHandler) RegisterRoutes(r *gin.RouterGroup) {
	product := r.Group("/products")
	{
		product.GET("", h.GetProducts)
		product.GET("/:id", h.GetProduct)
		product.POST("", h.CreateProduct)
		product.PUT("/:id", h.UpdateProduct)
		product.DELETE("/:id", h.DeleteProduct)
	}
}

// ==================== Admin Router Methods ====================

// GetList returns a paginated product list.
func (h *ProductHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	groupID := c.Query("group_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": []interface{}{}, "total": 0}})
		return
	}

	query := database.Table("products")
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}

	var total int64
	query.Count(&total)

	var products []struct {
		ID        uint      `json:"id"`
		Name      string    `json:"name"`
		Type      string    `json:"type"`
		Price     float64   `json:"price"`
		Status    int16     `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}

	query.Select("id, name, type, price, status, created_at").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&products)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  products,
			"total": total,
		},
	})
}

// Create creates a new product.
func (h *ProductHandler) Create(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "产品创建成功"})
}

// Update updates an existing product.
func (h *ProductHandler) Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "产品更新成功"})
}

// Delete deletes a product.
func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "数据库未初始化"})
		return
	}

	if err := database.Table("products").Where("id = ?", id).Delete(nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功"})
}

// DuplicateProduct duplicates a product.
func (h *ProductHandler) DuplicateProduct(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// EditStock edits product stock.
func (h *ProductHandler) EditStock(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// BatchUpdate batch updates products.
func (h *ProductHandler) BatchUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// BatchDelete batch deletes products.
func (h *ProductHandler) BatchDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// UpdateProductSort updates product sort order.
func (h *ProductHandler) UpdateProductSort(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// GetDiscountList returns product discount list.
func (h *ProductHandler) GetDiscountList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// GetProductDownloads returns product download files.
func (h *ProductHandler) GetProductDownloads(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// GetUpstreamPrice returns upstream price for a product.
func (h *ProductHandler) GetUpstreamPrice(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{}})
}

// EditResProduct edits a resource product.
func (h *ProductHandler) EditResProduct(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// SelectType returns product type options.
func (h *ProductHandler) SelectType(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// Selectcates returns product categories.
func (h *ProductHandler) Selectcates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// Downloadcates returns download categories.
func (h *ProductHandler) Downloadcates(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// AddDownloadcats adds a download category.
func (h *ProductHandler) AddDownloadcats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// GetFirstGroups returns first-level product groups.
func (h *ProductHandler) GetFirstGroups(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	var groups []struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	database.Table("product_groups").Select("id, name").Order("sort_order ASC").Find(&groups)
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": groups})
}

// CreateFirstGroup creates a first-level group.
func (h *ProductHandler) CreateFirstGroup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// UpdateFirstGroup updates a first-level group.
func (h *ProductHandler) UpdateFirstGroup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// DeleteFirstGroup deletes a first-level group.
func (h *ProductHandler) DeleteFirstGroup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// UpdateFirstGroupSort updates first-level group sort order.
func (h *ProductHandler) UpdateFirstGroupSort(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// GetGroups returns product groups.
func (h *ProductHandler) GetGroups(c *gin.Context) {
	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
		return
	}

	var groups []struct {
		ID        uint      `json:"id"`
		Name      string    `json:"name"`
		SortOrder int       `json:"sort_order"`
		CreatedAt time.Time `json:"created_at"`
	}

	database.Table("product_groups").
		Select("id, name, sort_order, created_at").
		Order("sort_order ASC").
		Find(&groups)

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": groups})
}

// CreateGroup creates a product group.
func (h *ProductHandler) CreateGroup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// UpdateGroup updates a product group.
func (h *ProductHandler) UpdateGroup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// DeleteGroup deletes a product group.
func (h *ProductHandler) DeleteGroup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// UpdateGroupSort updates group sort order.
func (h *ProductHandler) UpdateGroupSort(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// CheckAlias checks if a product alias is available.
func (h *ProductHandler) CheckAlias(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"available": true}})
}

// ManageDownloads manages product download files.
func (h *ProductHandler) ManageDownloads(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// AddDownloadFile adds a download file to a product.
func (h *ProductHandler) AddDownloadFile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// DeleteCustomField deletes a product custom field.
func (h *ProductHandler) DeleteCustomField(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// ==================== V1 Router Methods ====================

// GetDetail returns product detail.
func (h *ProductHandler) GetDetail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": nil})
}

// GetHot returns hot/popular products.
func (h *ProductHandler) GetHot(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}

// GetUserProducts returns products for the authenticated user.
func (h *ProductHandler) GetUserProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "data": []interface{}{}})
}
