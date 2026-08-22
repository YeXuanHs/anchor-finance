package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetSalesConfig 获取销售设置
// GET /api/admin/sales-config
func GetSalesConfig(c *gin.Context) {
	db := database.GetDB()
	var configs []model.SalesConfig
	db.Where("`group` = ?", "sales").Find(&configs)

	result := make(map[string]interface{})
	for _, config := range configs {
		result[config.Key] = gin.H{
			"value":   config.Value,
			"name":    config.Name,
			"type":    config.Type,
			"options": config.Options,
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// UpdateSalesConfig 更新销售设置
// PUT /api/admin/sales-config
func UpdateSalesConfig(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	for key, value := range req {
		db.Where("key = ?", key).Assign(map[string]interface{}{"value": value}).FirstOrCreate(&model.SalesConfig{})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// GetSalesGroupList 获取销售分组列表
// GET /api/admin/sales-groups
func GetSalesGroupList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	var total int64
	db.Model(&model.SalesGroup{}).Count(&total)

	var groups []model.SalesGroup
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("sort_order ASC").Find(&groups)

	if groups == nil {
		groups = []model.SalesGroup{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      groups,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateSalesGroup 创建销售分组
// POST /api/admin/sales-groups
func CreateSalesGroup(c *gin.Context) {
	var req struct {
		Name       string  `json:"name" binding:"required"`
		Commission float64 `json:"commission"`
		SortOrder  int     `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if req.Commission < 0 || req.Commission > 100 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "佣金比例需在0-100之间", "data": nil})
		return
	}

	db := database.GetDB()
	group := model.SalesGroup{
		Name:       req.Name,
		Commission: req.Commission,
		SortOrder:  req.SortOrder,
		Status:     "active",
	}

	if err := db.Create(&group).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败", "data": nil})
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

// UpdateSalesGroup 更新销售分组
// PUT /api/admin/sales-groups/:id
func UpdateSalesGroup(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name       string  `json:"name"`
		Commission float64 `json:"commission"`
		SortOrder  int     `json:"sort_order"`
		Status     string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var group model.SalesGroup
	if err := db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "分组不存在", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Commission >= 0 && req.Commission <= 100 {
		updates["commission"] = req.Commission
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&group).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteSalesGroup 删除销售分组
// DELETE /api/admin/sales-groups/:id
func DeleteSalesGroup(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var group model.SalesGroup
	if err := db.First(&group, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "分组不存在", "data": nil})
		return
	}

	db.Delete(&group)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
