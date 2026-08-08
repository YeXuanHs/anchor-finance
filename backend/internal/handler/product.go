package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

// GetProducts 获取产品列表
func (h *ProductHandler) GetProducts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"list": []interface{}{}, "total": 0})
}

// GetProduct 获取单个产品
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{"id": id})
}

// CreateProduct 创建产品
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "产品创建成功"})
}

// UpdateProduct 更新产品
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "产品更新成功"})
}

// DeleteProduct 删除产品
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "产品删除成功"})
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
