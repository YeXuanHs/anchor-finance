package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetProductList 获取产品列表
// GET /api/admin/products
func GetProductList(c *gin.Context) {
	// 1. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	groupID := c.Query("group_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 2. 构建查询
	db := database.GetDB()
	query := db.Model(&model.Product{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 分组筛选
	if groupID != "" {
		query = query.Where("group_id = ?", groupID)
	}

	// 3. 获取总数
	var total int64
	query.Count(&total)

	// 4. 分页查询
	var products []model.Product
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("sort_order ASC, id DESC").Find(&products)

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      products,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetProduct 获取产品详情
// GET /api/admin/products/:id
func GetProduct(c *gin.Context) {
	// 1. 获取产品ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的产品ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询产品
	db := database.GetDB()
	var product model.Product
	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "产品不存在",
			"data":    nil,
		})
		return
	}

	// 3. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    product,
	})
}

// CreateProduct 创建产品
// POST /api/admin/products
func CreateProduct(c *gin.Context) {
	// 1. 解析请求参数
	var req struct {
		Name         string  `json:"name" binding:"required"`
		GroupID      uint    `json:"group_id"`
		Type         string  `json:"type"`
		Description  string  `json:"description"`
		Price        float64 `json:"price"`
		BillingCycle string  `json:"billing_cycle"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 2. 创建产品
	db := database.GetDB()
	product := model.Product{
		Name:         req.Name,
		GroupID:      req.GroupID,
		Type:         req.Type,
		Description:  req.Description,
		Price:        req.Price,
		BillingCycle: req.BillingCycle,
		Status:       "active",
	}

	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建产品失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 3. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": product.ID,
		},
	})
}

// UpdateProduct 更新产品
// PUT /api/admin/products/:id
func UpdateProduct(c *gin.Context) {
	// 1. 获取产品ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的产品ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析请求参数
	var req struct {
		Name         string  `json:"name"`
		GroupID      uint    `json:"group_id"`
		Type         string  `json:"type"`
		Description  string  `json:"description"`
		Price        float64 `json:"price"`
		BillingCycle string  `json:"billing_cycle"`
		Status       string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 3. 查询产品
	db := database.GetDB()
	var product model.Product
	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "产品不存在",
			"data":    nil,
		})
		return
	}

	// 4. 更新产品
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.GroupID > 0 {
		updates["group_id"] = req.GroupID
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Price > 0 {
		updates["price"] = req.Price
	}
	if req.BillingCycle != "" {
		updates["billing_cycle"] = req.BillingCycle
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := db.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "更新产品失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    nil,
	})
}

// DeleteProduct 删除产品
// DELETE /api/admin/products/:id
func DeleteProduct(c *gin.Context) {
	// 1. 获取产品ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的产品ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询产品
	db := database.GetDB()
	var product model.Product
	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "产品不存在",
			"data":    nil,
		})
		return
	}

	// 3. 软删除产品
	if err := db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "删除产品失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 4. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"data":    nil,
	})
}

// GetProductGroups 获取产品分组列表
// GET /api/admin/product-groups
func GetProductGroups(c *gin.Context) {
	// 查询所有分组
	db := database.GetDB()
	var groups []model.ProductGroup
	db.Where("status = ?", "active").Order("sort_order ASC").Find(&groups)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    groups,
	})
}
