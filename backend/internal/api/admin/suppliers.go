package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/util"
	"github.com/gin-gonic/gin"
)

// GetSupplierList 获取供应商列表
// GET /api/admin/suppliers
func GetSupplierList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Supplier{})

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var suppliers []model.Supplier
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&suppliers)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      suppliers,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetSupplierSummary 获取供应商统计
// GET /api/admin/suppliers/summary
func GetSupplierSummary(c *gin.Context) {
	db := database.GetDB()

	var activeCount int64
	db.Model(&model.Supplier{}).Where("status = ?", "active").Count(&activeCount)

	var disabledCount int64
	db.Model(&model.Supplier{}).Where("status = ?", "disabled").Count(&disabledCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"active":   activeCount,
			"disabled": disabledCount,
			"total":    activeCount + disabledCount,
		},
	})
}

// GetSupplierDetail 获取供应商详情
// GET /api/admin/suppliers/:id
func GetSupplierDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的供应商ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "供应商不存在",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    supplier,
	})
}

// CreateSupplier 创建供应商
// POST /api/admin/suppliers
func CreateSupplier(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		APIURL      string `json:"api_url"`
		APIKey      string `json:"api_key"`
		APISecret   string `json:"api_secret"`
		Description string `json:"description"`
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

	// 加密存储API密钥
	encryptedKey, _ := util.EncryptAES(req.APIKey)
	encryptedSecret, _ := util.EncryptAES(req.APISecret)

	supplier := model.Supplier{
		Name:        req.Name,
		Type:        req.Type,
		APIURL:      req.APIURL,
		APIKey:      encryptedKey,
		APISecret:   encryptedSecret,
		Description: req.Description,
		Status:      "active",
	}

	if err := db.Create(&supplier).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建供应商失败",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": supplier.ID,
		},
	})
}

// UpdateSupplier 更新供应商
// PUT /api/admin/suppliers/:id
func UpdateSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的供应商ID",
			"data":    nil,
		})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		APIURL      string `json:"api_url"`
		APIKey      string `json:"api_key"`
		APISecret   string `json:"api_secret"`
		Description string `json:"description"`
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
	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "供应商不存在",
			"data":    nil,
		})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.APIURL != "" {
		updates["api_url"] = req.APIURL
	}
	if req.APIKey != "" {
		encryptedKey, _ := util.EncryptAES(req.APIKey)
		updates["api_key"] = encryptedKey
	}
	if req.APISecret != "" {
		encryptedSecret, _ := util.EncryptAES(req.APISecret)
		updates["api_secret"] = encryptedSecret
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	db.Model(&supplier).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    nil,
	})
}

// DeleteSupplier 删除供应商
// DELETE /api/admin/suppliers/:id
func DeleteSupplier(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的供应商ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "供应商不存在",
			"data":    nil,
		})
		return
	}

	db.Delete(&supplier)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"data":    nil,
	})
}

// GetSupplierProviderTypes 获取供应商类型列表
// GET /api/admin/suppliers/provider-types
func GetSupplierProviderTypes(c *gin.Context) {
	// 返回支持的供应商类型
	types := []gin.H{
		{"id": "manual", "name": "手动管理", "description": "不对接API，手动管理"},
		{"id": "zjmf", "name": "zjmf接口", "description": "对接zjmf API"},
		{"id": "v10", "name": "v10接口", "description": "对接v10 API"},
		{"id": "anchor", "name": "锚点接口", "description": "对接锚点财务自有API"},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    types,
	})
}

// UpdateSupplierStatus 更新供应商状态
// PATCH /api/admin/suppliers/:id/status
func UpdateSupplierStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的供应商ID",
			"data":    nil,
		})
		return
	}

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

	// 验证状态值
	validStatuses := map[string]bool{
		"active":   true,
		"disabled": true,
	}
	if !validStatuses[req.Status] {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的状态值",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "供应商不存在",
			"data":    nil,
		})
		return
	}

	db.Model(&supplier).Update("status", req.Status)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "状态更新成功",
		"data":    nil,
	})
}

// GetSupplierBalance 获取供应商余额
// GET /api/admin/suppliers/:id/balance
func GetSupplierBalance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的供应商ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "供应商不存在",
			"data":    nil,
		})
		return
	}

	// 返回供应商余额（数据库记录；实时上游余额待供应商插件接入）
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"supplier_id": supplier.ID,
			"name":        supplier.Name,
			"balance":     supplier.Balance,
		},
	})
}

// GetSupplierProducts 获取供应商产品
// GET /api/admin/suppliers/:id/products
func GetSupplierProducts(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的供应商ID",
			"data":    nil,
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	db := database.GetDB()
	var total int64
	db.Model(&model.SupplierProduct{}).Where("supplier_id = ?", id).Count(&total)

	var products []model.SupplierProduct
	offset := (page - 1) * pageSize
	db.Where("supplier_id = ?", id).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&products)

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

// RevealSupplierSecret 查看供应商密钥（解密返回完整值）
// GET /api/admin/suppliers/:id/secrets/:key
func RevealSupplierSecret(c *gin.Context) {
	id := c.Param("id")
	key := c.Param("key") // api_key 或 api_secret

	db := database.GetDB()
	var supplier model.Supplier
	if err := db.First(&supplier, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "供应商不存在", "data": nil})
		return
	}

	var encrypted string
	switch key {
	case "api_key":
		encrypted = supplier.APIKey
	case "api_secret":
		encrypted = supplier.APISecret
	default:
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的密钥类型", "data": nil})
		return
	}

	if encrypted == "" {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"key": key, "value": ""}})
		return
	}

	decrypted, err := util.DecryptAES(encrypted)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "解密失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gin.H{"key": key, "value": decrypted}})
}

// ==================== 供应商分组映射（MD 7.2.5）====================

// GetSupplierGroupMappings 获取供应商分组映射
// GET /api/admin/suppliers/:id/group-mappings
func GetSupplierGroupMappings(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var mappings []model.SupplierGroupMapping
	db.Where("supplier_id = ?", id).Order("id ASC").Find(&mappings)
	if mappings == nil {
		mappings = []model.SupplierGroupMapping{}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": mappings})
}

// CreateSupplierGroupMapping 创建分组映射
// POST /api/admin/suppliers/:id/group-mappings
func CreateSupplierGroupMapping(c *gin.Context) {
	supplierID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		RemoteGroupID   string  `json:"remote_group_id" binding:"required"`
		RemoteGroupName string  `json:"remote_group_name"`
		LocalGroupID    uint    `json:"local_group_id" binding:"required"`
		ProfitRate      float64 `json:"profit_rate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if req.ProfitRate <= 0 {
		req.ProfitRate = 25
	}

	mapping := model.SupplierGroupMapping{
		SupplierID:      uint(supplierID),
		RemoteGroupID:   req.RemoteGroupID,
		RemoteGroupName: req.RemoteGroupName,
		LocalGroupID:    req.LocalGroupID,
		ProfitRate:      req.ProfitRate,
	}

	db := database.GetDB()
	if err := db.Create(&mapping).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "创建成功", "data": gin.H{"id": mapping.ID}})
}

// UpdateSupplierGroupMapping 更新分组映射
// PUT /api/admin/suppliers/:id/group-mappings/:mapping_id
func UpdateSupplierGroupMapping(c *gin.Context) {
	mappingID := c.Param("mapping_id")
	db := database.GetDB()

	var req struct {
		RemoteGroupID   string  `json:"remote_group_id"`
		RemoteGroupName string  `json:"remote_group_name"`
		LocalGroupID    uint    `json:"local_group_id"`
		ProfitRate      float64 `json:"profit_rate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.RemoteGroupID != "" {
		updates["remote_group_id"] = req.RemoteGroupID
	}
	if req.RemoteGroupName != "" {
		updates["remote_group_name"] = req.RemoteGroupName
	}
	if req.LocalGroupID > 0 {
		updates["local_group_id"] = req.LocalGroupID
	}
	if req.ProfitRate > 0 {
		updates["profit_rate"] = req.ProfitRate
	}

	db.Model(&model.SupplierGroupMapping{}).Where("id = ?", mappingID).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteSupplierGroupMapping 删除分组映射
// DELETE /api/admin/suppliers/:id/group-mappings/:mapping_id
func DeleteSupplierGroupMapping(c *gin.Context) {
	mappingID := c.Param("mapping_id")
	db := database.GetDB()

	db.Where("id = ?", mappingID).Delete(&model.SupplierGroupMapping{})
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
