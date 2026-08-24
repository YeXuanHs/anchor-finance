package admin

import (
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/response"
	"github.com/gin-gonic/gin"
)

// ============================================================
// 配置分组管理
// ============================================================

// GetConfigGroupList 获取配置分组列表
// GET /api/admin/config-groups
func GetConfigGroupList(c *gin.Context) {
	db := database.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "15"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 15
	}

	query := db.Model(&model.ProductConfigGroup{})
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var list []model.ProductConfigGroup
	query.Order("`order` ASC, id ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&list)

	if list == nil {
		list = []model.ProductConfigGroup{}
	}

	response.SuccessPage(c, list, total, page, pageSize)
}

// CreateConfigGroup 创建配置分组
// POST /api/admin/config-groups
func CreateConfigGroup(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Order       int    `json:"order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	db := database.GetDB()
	group := model.ProductConfigGroup{
		Name:        req.Name,
		Description: req.Description,
		Order:       req.Order,
	}

	if err := db.Create(&group).Error; err != nil {
		response.ServerError(c, "创建分组失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"id": group.ID})
}

// UpdateConfigGroup 更新配置分组
// PUT /api/admin/config-groups/:id
func UpdateConfigGroup(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Order       *int   `json:"order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	db := database.GetDB()
	var group model.ProductConfigGroup
	if err := db.First(&group, id).Error; err != nil {
		response.NotFound(c, "分组不存在")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Order != nil {
		updates["order"] = *req.Order
	}

	if len(updates) > 0 {
		db.Model(&group).Updates(updates)
	}

	response.SuccessMsg(c, "更新成功")
}

// DeleteConfigGroup 删除配置分组
// DELETE /api/admin/config-groups/:id
func DeleteConfigGroup(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var group model.ProductConfigGroup
	if err := db.First(&group, id).Error; err != nil {
		response.NotFound(c, "分组不存在")
		return
	}

	// 检查是否有关联的产品
	var linkCount int64
	db.Model(&model.ProductConfigLink{}).Where("gid = ?", id).Count(&linkCount)
	if linkCount > 0 {
		response.BadRequest(c, "该分组已关联产品，无法删除")
		return
	}

	// 删除分组下的所有选项和子选项
	var optionIDs []uint
	db.Model(&model.ProductConfigOption{}).Where("gid = ?", id).Pluck("id", &optionIDs)
	if len(optionIDs) > 0 {
		db.Where("config_id IN ?", optionIDs).Delete(&model.ProductConfigOptionSub{})
		db.Where("config_id IN ?", optionIDs).Delete(&model.ProductConfigPricing{})
	}
	db.Where("gid = ?", id).Delete(&model.ProductConfigOption{})

	db.Delete(&group)

	response.SuccessMsg(c, "删除成功")
}

// ============================================================
// 配置选项管理
// ============================================================

// GetConfigOptionList 获取分组下的配置选项列表
// GET /api/admin/config-groups/:gid/options
func GetConfigOptionList(c *gin.Context) {
	gid := c.Param("gid")

	db := database.GetDB()

	// 验证分组存在
	var group model.ProductConfigGroup
	if err := db.First(&group, gid).Error; err != nil {
		response.NotFound(c, "分组不存在")
		return
	}

	var list []model.ProductConfigOption
	db.Where("gid = ?", gid).
		Order("`order` ASC, id ASC").
		Find(&list)

	if list == nil {
		list = []model.ProductConfigOption{}
	}

	response.Success(c, list)
}

// CreateConfigOption 创建配置选项
// POST /api/admin/config-options
func CreateConfigOption(c *gin.Context) {
	var req struct {
		GID        uint   `json:"gid" binding:"required"`
		OptionName string `json:"option_name" binding:"required"`
		OptionType int    `json:"option_type"`
		Order      int    `json:"order"`
		Hidden     bool   `json:"hidden"`
		Auto       bool   `json:"auto"`
		IsDiscount bool   `json:"is_discount"`
		IsRebate   bool   `json:"is_rebate"`
		QtyMinimum int    `json:"qty_minimum"`
		QtyMaximum int    `json:"qty_maximum"`
		QtyStage   int    `json:"qty_stage"`
		Unit       string `json:"unit"`
		Upgrade    bool   `json:"upgrade"`
		Notes      string `json:"notes"`
		Senior     bool   `json:"senior"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	db := database.GetDB()

	// 验证分组存在
	var group model.ProductConfigGroup
	if err := db.First(&group, req.GID).Error; err != nil {
		response.BadRequest(c, "分组不存在")
		return
	}

	option := model.ProductConfigOption{
		GID:        req.GID,
		OptionName: req.OptionName,
		OptionType: req.OptionType,
		Order:      req.Order,
		Hidden:     req.Hidden,
		Auto:       req.Auto,
		IsDiscount: req.IsDiscount,
		IsRebate:   req.IsRebate,
		QtyMinimum: req.QtyMinimum,
		QtyMaximum: req.QtyMaximum,
		QtyStage:   req.QtyStage,
		Unit:       req.Unit,
		Upgrade:    req.Upgrade,
		Notes:      req.Notes,
		Senior:     req.Senior,
	}

	if req.QtyStage == 0 {
		option.QtyStage = 1
	}

	if err := db.Create(&option).Error; err != nil {
		response.ServerError(c, "创建选项失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"id": option.ID})
}

// UpdateConfigOption 更新配置选项
// PUT /api/admin/config-options/:id
func UpdateConfigOption(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		OptionName string `json:"option_name"`
		OptionType *int   `json:"option_type"`
		Order      *int   `json:"order"`
		Hidden     *bool  `json:"hidden"`
		Auto       *bool  `json:"auto"`
		IsDiscount *bool  `json:"is_discount"`
		IsRebate   *bool  `json:"is_rebate"`
		QtyMinimum *int   `json:"qty_minimum"`
		QtyMaximum *int   `json:"qty_maximum"`
		QtyStage   *int   `json:"qty_stage"`
		Unit       string `json:"unit"`
		Upgrade    *bool  `json:"upgrade"`
		Notes      string `json:"notes"`
		Senior     *bool  `json:"senior"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	db := database.GetDB()
	var option model.ProductConfigOption
	if err := db.First(&option, id).Error; err != nil {
		response.NotFound(c, "选项不存在")
		return
	}

	updates := map[string]interface{}{}
	if req.OptionName != "" {
		updates["option_name"] = req.OptionName
	}
	if req.OptionType != nil {
		updates["option_type"] = *req.OptionType
	}
	if req.Order != nil {
		updates["order"] = *req.Order
	}
	if req.Hidden != nil {
		updates["hidden"] = *req.Hidden
	}
	if req.Auto != nil {
		updates["auto"] = *req.Auto
	}
	if req.IsDiscount != nil {
		updates["is_discount"] = *req.IsDiscount
	}
	if req.IsRebate != nil {
		updates["is_rebate"] = *req.IsRebate
	}
	if req.QtyMinimum != nil {
		updates["qty_minimum"] = *req.QtyMinimum
	}
	if req.QtyMaximum != nil {
		updates["qty_maximum"] = *req.QtyMaximum
	}
	if req.QtyStage != nil {
		updates["qty_stage"] = *req.QtyStage
	}
	if req.Unit != "" {
		updates["unit"] = req.Unit
	}
	if req.Upgrade != nil {
		updates["upgrade"] = *req.Upgrade
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}
	if req.Senior != nil {
		updates["senior"] = *req.Senior
	}

	if len(updates) > 0 {
		db.Model(&option).Updates(updates)
	}

	response.SuccessMsg(c, "更新成功")
}

// DeleteConfigOption 删除配置选项
// DELETE /api/admin/config-options/:id
func DeleteConfigOption(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var option model.ProductConfigOption
	if err := db.First(&option, id).Error; err != nil {
		response.NotFound(c, "选项不存在")
		return
	}

	// 删除子选项和相关定价
	var subIDs []uint
	db.Model(&model.ProductConfigOptionSub{}).Where("config_id = ?", id).Pluck("id", &subIDs)
	if len(subIDs) > 0 {
		db.Where("rel_id IN ? AND type = ?", subIDs, "config_option").Delete(&model.ProductConfigPricing{})
	}
	db.Where("config_id = ?", id).Delete(&model.ProductConfigOptionSub{})

	db.Delete(&option)

	response.SuccessMsg(c, "删除成功")
}

// ============================================================
// 配置子选项管理
// ============================================================

// GetConfigOptionSubList 获取选项下的子选项列表
// GET /api/admin/config-options/:oid/subs
func GetConfigOptionSubList(c *gin.Context) {
	oid := c.Param("oid")

	db := database.GetDB()

	// 验证选项存在
	var option model.ProductConfigOption
	if err := db.First(&option, oid).Error; err != nil {
		response.NotFound(c, "选项不存在")
		return
	}

	var list []model.ProductConfigOptionSub
	db.Where("config_id = ?", oid).
		Order("sort_order ASC, id ASC").
		Find(&list)

	if list == nil {
		list = []model.ProductConfigOptionSub{}
	}

	response.Success(c, list)
}

// CreateConfigOptionSub 创建子选项
// POST /api/admin/config-option-subs
func CreateConfigOptionSub(c *gin.Context) {
	var req struct {
		ConfigID   uint   `json:"config_id" binding:"required"`
		OptionName string `json:"option_name" binding:"required"`
		SortOrder  int    `json:"sort_order"`
		Hidden     bool   `json:"hidden"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	db := database.GetDB()

	// 验证选项存在
	var option model.ProductConfigOption
	if err := db.First(&option, req.ConfigID).Error; err != nil {
		response.BadRequest(c, "配置选项不存在")
		return
	}

	sub := model.ProductConfigOptionSub{
		ConfigID:   req.ConfigID,
		OptionName: req.OptionName,
		SortOrder:  req.SortOrder,
		Hidden:     req.Hidden,
	}

	if err := db.Create(&sub).Error; err != nil {
		response.ServerError(c, "创建子选项失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"id": sub.ID})
}

// UpdateConfigOptionSub 更新子选项
// PUT /api/admin/config-option-subs/:id
func UpdateConfigOptionSub(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		OptionName string `json:"option_name"`
		SortOrder  *int   `json:"sort_order"`
		Hidden     *bool  `json:"hidden"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	db := database.GetDB()
	var sub model.ProductConfigOptionSub
	if err := db.First(&sub, id).Error; err != nil {
		response.NotFound(c, "子选项不存在")
		return
	}

	updates := map[string]interface{}{}
	if req.OptionName != "" {
		updates["option_name"] = req.OptionName
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Hidden != nil {
		updates["hidden"] = *req.Hidden
	}

	if len(updates) > 0 {
		db.Model(&sub).Updates(updates)
	}

	response.SuccessMsg(c, "更新成功")
}

// DeleteConfigOptionSub 删除子选项
// DELETE /api/admin/config-option-subs/:id
func DeleteConfigOptionSub(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var sub model.ProductConfigOptionSub
	if err := db.First(&sub, id).Error; err != nil {
		response.NotFound(c, "子选项不存在")
		return
	}

	// 删除相关定价
	db.Where("rel_id = ? AND type = ?", id, "config_option").Delete(&model.ProductConfigPricing{})

	db.Delete(&sub)

	response.SuccessMsg(c, "删除成功")
}

// ============================================================
// 产品-配置分组关联管理
// ============================================================

// CreateProductConfigLink 关联产品到配置分组
// POST /api/admin/products/:pid/config-links
func CreateProductConfigLink(c *gin.Context) {
	pid := c.Param("pid")

	var req struct {
		GID uint `json:"gid" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	db := database.GetDB()

	// 验证产品存在
	var product model.Product
	if err := db.First(&product, pid).Error; err != nil {
		response.BadRequest(c, "产品不存在")
		return
	}

	// 验证分组存在
	var group model.ProductConfigGroup
	if err := db.First(&group, req.GID).Error; err != nil {
		response.BadRequest(c, "配置分组不存在")
		return
	}

	pidUint, _ := strconv.ParseUint(pid, 10, 64)

	// 检查是否已关联
	var existing model.ProductConfigLink
	if err := db.Where("pid = ? AND gid = ?", pidUint, req.GID).First(&existing).Error; err == nil {
		response.BadRequest(c, "该产品已关联此配置分组")
		return
	}

	link := model.ProductConfigLink{
		PID: uint(pidUint),
		GID: req.GID,
	}

	if err := db.Create(&link).Error; err != nil {
		response.ServerError(c, "关联失败: "+err.Error())
		return
	}

	response.Success(c, gin.H{"id": link.ID})
}

// DeleteProductConfigLink 取消产品与配置分组的关联
// DELETE /api/admin/products/:pid/config-links/:gid
func DeleteProductConfigLink(c *gin.Context) {
	pid := c.Param("pid")
	gid := c.Param("gid")

	db := database.GetDB()

	pidUint, _ := strconv.ParseUint(pid, 10, 64)
	gidUint, _ := strconv.ParseUint(gid, 10, 64)

	result := db.Where("pid = ? AND gid = ?", pidUint, gidUint).Delete(&model.ProductConfigLink{})
	if result.RowsAffected == 0 {
		response.NotFound(c, "关联关系不存在")
		return
	}

	response.SuccessMsg(c, "取消关联成功")
}

// ============================================================
// 子选项定价管理
// ============================================================

// GetConfigOptionSubPricing 获取子选项定价
// GET /api/admin/config-option-subs/:sid/pricing
func GetConfigOptionSubPricing(c *gin.Context) {
	sid := c.Param("sid")

	db := database.GetDB()

	// 验证子选项存在
	var sub model.ProductConfigOptionSub
	if err := db.First(&sub, sid).Error; err != nil {
		response.NotFound(c, "子选项不存在")
		return
	}

	var pricings []model.ProductConfigPricing
	db.Where("rel_id = ? AND type = ?", sid, "config_option").Find(&pricings)

	if pricings == nil {
		pricings = []model.ProductConfigPricing{}
	}

	response.Success(c, pricings)
}

// SetConfigOptionSubPricing 设置子选项定价
// POST /api/admin/config-option-subs/:sid/pricing
func SetConfigOptionSubPricing(c *gin.Context) {
	sid := c.Param("sid")

	var req struct {
		Currency          uint    `json:"currency" binding:"required"`
		Monthly           float64 `json:"monthly"`
		Quarterly         float64 `json:"quarterly"`
		Semiannually      float64 `json:"semiannually"`
		Annually          float64 `json:"annually"`
		Biennially        float64 `json:"biennially"`
		Triennially       float64 `json:"triennially"`
		MonthlySetup      float64 `json:"monthlysetupfee"`
		QuarterlySetup    float64 `json:"quarterlysetupfee"`
		SemiannuallySetup float64 `json:"semiannuallysetupfee"`
		AnnuallySetup     float64 `json:"annuallysetupfee"`
		BienniallySetup   float64 `json:"bienniallysetupfee"`
		TrienniallySetup  float64 `json:"trienniallysetupfee"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	db := database.GetDB()

	// 验证子选项存在
	var sub model.ProductConfigOptionSub
	if err := db.First(&sub, sid).Error; err != nil {
		response.NotFound(c, "子选项不存在")
		return
	}

	sidUint, _ := strconv.ParseUint(sid, 10, 64)

	// 查找是否已有该货币的定价记录
	var pricing model.ProductConfigPricing
	err := db.Where("rel_id = ? AND type = ? AND currency = ?", sidUint, "config_option", req.Currency).
		First(&pricing).Error

	if err == nil {
		// 更新已有定价
		db.Model(&pricing).Updates(map[string]interface{}{
			"monthly":             req.Monthly,
			"quarterly":           req.Quarterly,
			"semiannually":        req.Semiannually,
			"annually":            req.Annually,
			"biennially":          req.Biennially,
			"triennially":         req.Triennially,
			"monthlysetupfee":     req.MonthlySetup,
			"quarterlysetupfee":   req.QuarterlySetup,
			"semiannuallysetupfee": req.SemiannuallySetup,
			"annuallysetupfee":    req.AnnuallySetup,
			"bienniallysetupfee":  req.BienniallySetup,
			"trienniallysetupfee": req.TrienniallySetup,
		})
		response.Success(c, gin.H{"id": pricing.ID})
	} else {
		// 创建新定价
		pricing = model.ProductConfigPricing{
			RelID:             uint(sidUint),
			Type:              "config_option",
			Currency:          int(req.Currency),
			Monthly:           req.Monthly,
			Quarterly:         req.Quarterly,
			Semiannual:        req.Semiannually,
			Annually:          req.Annually,
			Biennially:        req.Biennially,
			Triennially:       req.Triennially,
			MonthlySetup:      req.MonthlySetup,
			QuarterlySetup:    req.QuarterlySetup,
			SemiannualSetup:   req.SemiannuallySetup,
			AnnuallySetup:     req.AnnuallySetup,
			BienniallySetup:   req.BienniallySetup,
			TrienniallySetup:  req.TrienniallySetup,
		}

		if err := db.Create(&pricing).Error; err != nil {
			response.ServerError(c, "设置定价失败: "+err.Error())
			return
		}

		response.Success(c, gin.H{"id": pricing.ID})
	}
}
