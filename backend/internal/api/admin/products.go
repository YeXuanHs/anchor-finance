package admin

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
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
		Name          string  `json:"name" binding:"required"`
		GroupID       uint    `json:"group_id"`
		Type          string  `json:"type"`
		Description   string  `json:"description"`
		Price         float64 `json:"price"`
		BillingCycle  string  `json:"billing_cycle"`
		ConfigOptions string  `json:"config_options"` // JSON配置选项
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	// 0元购防护：价格必须大于0
	if req.Price <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "价格必须大于0", "data": nil})
		return
	}

	// 2. 创建产品
	db := database.GetDB()
	product := model.Product{
		Name:          req.Name,
		GroupID:       req.GroupID,
		Type:          req.Type,
		Description:   req.Description,
		Amount:        req.Price,
		Price:         req.Price,
		BillingCycle:  req.BillingCycle,
		ConfigOptions: req.ConfigOptions,
		Status:        "active",
	}

	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建产品失败",
			"data":    nil,
		})
		return
	}

	// 触发Hook: product_create
	pluginengine.TriggerHook("product_create", map[string]interface{}{
		"product_id": product.ID, "name": product.Name, "price": product.Price,
	})

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
		Name          string  `json:"name"`
		GroupID       uint    `json:"group_id"`
		Type          string  `json:"type"`
		Description   string  `json:"description"`
		Price         float64 `json:"price"`
		BillingCycle  string  `json:"billing_cycle"`
		Status        string  `json:"status"`
		ConfigOptions string  `json:"config_options"` // JSON配置选项
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
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
		updates["amount"] = req.Price
	}
	if req.BillingCycle != "" {
		updates["billing_cycle"] = req.BillingCycle
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.ConfigOptions != "" {
		updates["config_options"] = req.ConfigOptions
	}

	if err := db.Model(&product).Updates(updates).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "更新产品失败",
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
			"message": "删除产品失败",
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

// RestoreProduct 恢复产品
// POST /api/admin/products/:id/restorations
func RestoreProduct(c *gin.Context) {
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

	// 2. 查询产品（包括软删除的）
	db := database.GetDB()
	var product model.Product
	if err := db.Unscoped().First(&product, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "产品不存在",
			"data":    nil,
		})
		return
	}

	// 3. 检查是否已删除
	if product.DeletedAt.Time.IsZero() {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "产品未被删除",
			"data":    nil,
		})
		return
	}

	// 4. 恢复产品
	if err := db.Unscoped().Model(&product).Update("deleted_at", nil).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "恢复产品失败",
			"data":    nil,
		})
		return
	}

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "恢复成功",
		"data":    nil,
	})
}

// GetProductGroupChildren 获取产品分组子级
// GET /api/admin/product-groups/:id/children
func GetProductGroupChildren(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var groups []model.ProductGroup
	db.Where("parent_id = ? AND status = ?", id, "active").Order("sort_order ASC").Find(&groups)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    groups,
	})
}

// GetProductGroupTree 获取产品分组树
// GET /api/admin/product-groups/tree
func GetProductGroupTree(c *gin.Context) {
	db := database.GetDB()
	var groups []model.ProductGroup
	db.Where("status = ?", "active").Order("sort_order ASC").Find(&groups)

	// 构建树形结构
	tree := buildProductGroupTree(groups, 0)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    tree,
	})
}

// buildProductGroupTree 构建产品分组树
func buildProductGroupTree(groups []model.ProductGroup, parentID uint) []gin.H {
	var result []gin.H
	for _, group := range groups {
		if group.ParentID == parentID {
			children := buildProductGroupTree(groups, group.ID)
			item := gin.H{
				"id":          group.ID,
				"name":        group.Name,
				"description": group.Description,
				"sort_order":  group.SortOrder,
			}
			if len(children) > 0 {
				item["children"] = children
			}
			result = append(result, item)
		}
	}
	return result
}

// GetProductSummary 获取产品统计
// GET /api/admin/products/summary
func GetProductSummary(c *gin.Context) {
	db := database.GetDB()

	var activeCount int64
	db.Model(&model.Product{}).Where("status = ?", "active").Count(&activeCount)

	var hiddenCount int64
	db.Model(&model.Product{}).Where("status = ?", "hidden").Count(&hiddenCount)

	var totalCount int64
	db.Model(&model.Product{}).Count(&totalCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"active": activeCount,
			"hidden": hiddenCount,
			"total":  totalCount,
		},
	})
}

// GetProductOwners 获取产品所有者列表
// GET /api/admin/products/:id/owners
func GetProductOwners(c *gin.Context) {
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

	// 2. 查询购买了该产品的用户
	db := database.GetDB()
	var users []model.User
	db.Joins("JOIN orders ON orders.user_id = users.id").
		Where("orders.product_id = ? AND orders.status = ?", id, "paid").
		Group("users.id").
		Find(&users)

	// 3. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    users,
	})
}

// UpdateProductStatus 更新产品状态
// PATCH /api/admin/products/:id/status
func UpdateProductStatus(c *gin.Context) {
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
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	// 3. 验证状态值
	validStatuses := map[string]bool{
		"active":  true,
		"hidden":  true,
		"deleted": true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的状态值",
			"data":    nil,
		})
		return
	}

	// 4. 查询产品
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

	// 5. 更新状态
	if err := db.Model(&product).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "更新状态失败",
			"data":    nil,
		})
		return
	}

	// 6. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "状态更新成功",
		"data":    nil,
	})
}

// CreateProductGroup 创建产品分组
// POST /api/admin/product-groups
func CreateProductGroup(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		ParentID    uint   `json:"parent_id"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	group := model.ProductGroup{
		Name:        req.Name,
		ParentID:    req.ParentID,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		Status:      "active",
	}

	if err := db.Create(&group).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建分组失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": group.ID,
		},
	})
}

// UpdateProductGroup 更新产品分组
// PUT /api/admin/product-groups/:id
func UpdateProductGroup(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		SortOrder   int    `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var group model.ProductGroup
	if err := db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "分组不存在",
			"data":    nil,
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}

	db.Model(&group).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    nil,
	})
}

// DeleteProductGroup 删除产品分组
// DELETE /api/admin/product-groups/:id
func DeleteProductGroup(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var group model.ProductGroup
	if err := db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "分组不存在",
			"data":    nil,
		})
		return
	}

	db.Delete(&group)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"data":    nil,
	})
}

// SplitProductPreview 产品拆分预览
// POST /api/admin/products/split-previews
func SplitProductPreview(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Count     int  `json:"count" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var product model.Product
	if err := db.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "产品不存在", "data": nil})
		return
	}

	previews := make([]gin.H, req.Count)
	for i := 0; i < req.Count; i++ {
		previews[i] = gin.H{
			"name":   fmt.Sprintf("%s #%d", product.Name, i+1),
			"price":  product.Price,
			"amount": product.Amount,
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"previews": previews, "original": product}})
}

// SplitProduct 执行产品拆分
// POST /api/admin/products/splits
func SplitProduct(c *gin.Context) {
	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Count     int  `json:"count" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if req.Count < 2 || req.Count > 100 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "拆分数量必须在2-100之间", "data": nil})
		return
	}

	db := database.GetDB()
	var product model.Product
	if err := db.First(&product, req.ProductID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "产品不存在", "data": nil})
		return
	}

	created := 0
	for i := 0; i < req.Count; i++ {
		newProduct := model.Product{
			Name:        fmt.Sprintf("%s #%d", product.Name, i+1),
			Type:        product.Type,
			Price:       product.Price,
			Amount:      product.Amount,
			Description: product.Description,
			GroupID:     product.GroupID,
			Status:      product.Status,
		}
		if err := db.Create(&newProduct).Error; err == nil {
			created++
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "拆分成功", "data": gin.H{"created": created}})
}

// BatchUpdateProvisionHostname 批量更新开通主机名模板
// POST /api/admin/products/provision-hostname-batches
func BatchUpdateProvisionHostname(c *gin.Context) {
	var req struct {
		ProductIDs []uint `json:"product_ids" binding:"required"`
		Hostname   string `json:"hostname" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	updated := db.Model(&model.Product{}).Where("id IN ?", req.ProductIDs).Update("provision_hostname", req.Hostname).RowsAffected

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": gin.H{"updated": updated}})
}
