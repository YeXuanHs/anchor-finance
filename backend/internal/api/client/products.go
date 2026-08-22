package client

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetClientProducts 获取可购买产品列表
// GET /api/client/products
func GetClientProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	groupID := c.Query("group_id")
	keyword := c.Query("keyword")

	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }

	db := database.GetDB()
	query := db.Model(&model.Product{}).Where("status = ?", "active")

	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var products []model.Product
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("sort_order ASC").Find(&products)

	if products == nil { products = []model.Product{} }

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{"list": products, "total": total, "page": page, "page_size": pageSize},
	})
}

// GetClientProductCategories 获取产品分组（分类浏览）
// GET /api/client/products/categories
func GetClientProductCategories(c *gin.Context) {
	db := database.GetDB()
	var groups []model.ProductGroup
	db.Where("status = ?", "active").Order("sort_order ASC").Find(&groups)

	if groups == nil { groups = []model.ProductGroup{} }

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": groups})
}

// GetClientProductDetail 获取产品详情
// GET /api/client/products/:id
func GetClientProductDetail(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var product model.Product
	if err := db.Where("id = ? AND status = ?", id, "active").First(&product).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "产品不存在", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": product})
}
