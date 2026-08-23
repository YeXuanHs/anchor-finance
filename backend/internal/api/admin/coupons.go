package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetCouponList 获取优惠券列表
// GET /api/admin/coupons
func GetCouponList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Coupon{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var coupons []model.Coupon
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&coupons)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      coupons,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetCouponSummary 获取优惠券统计
// GET /api/admin/coupons/summary
func GetCouponSummary(c *gin.Context) {
	db := database.GetDB()

	var activeCount int64
	db.Model(&model.Coupon{}).Where("status = ?", "active").Count(&activeCount)

	var disabledCount int64
	db.Model(&model.Coupon{}).Where("status = ?", "disabled").Count(&disabledCount)

	var totalCount int64
	db.Model(&model.Coupon{}).Count(&totalCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"active":   activeCount,
			"disabled": disabledCount,
			"total":    totalCount,
		},
	})
}

// CreateCoupon 创建优惠券
// POST /api/admin/coupons
func CreateCoupon(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Code        string  `json:"code" binding:"required"`
		Type        string  `json:"type" binding:"required"`
		Value       float64 `json:"value" binding:"required"`
		MinAmount   float64 `json:"min_amount"`
		MaxDiscount float64 `json:"max_discount"`
		UsageLimit  int     `json:"usage_limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	// 验证类型
	if req.Type != "percent" && req.Type != "fixed" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "类型必须是 percent 或 fixed", "data": nil})
		return
	}

	// 0元购防护：优惠券价值校验（percent 0-100，fixed > 0）
	if req.Type == "percent" && (req.Value < 0 || req.Value > 100) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "百分比优惠券价值需在0-100之间", "data": nil})
		return
	}
	if req.Type == "fixed" && req.Value < 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "固定金额优惠券价值不能为负数", "data": nil})
		return
	}

	db := database.GetDB()

	// 检查优惠券代码是否已存在
	var count int64
	db.Model(&model.Coupon{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "优惠券代码已存在",
			"data": nil,
		})
		return
	}

	coupon := model.Coupon{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		Value:       req.Value,
		MinAmount:   req.MinAmount,
		MaxDiscount: req.MaxDiscount,
		UsageLimit:  req.UsageLimit,
		Status:      "active",
	}

	if err := db.Create(&coupon).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建优惠券失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": coupon.ID,
		},
	})
}

// UpdateCoupon 更新优惠券
// PUT /api/admin/coupons/:id
func UpdateCoupon(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string  `json:"name"`
		Type        string  `json:"type"`
		Value       float64 `json:"value"`
		MinAmount   float64 `json:"min_amount"`
		MaxDiscount float64 `json:"max_discount"`
		UsageLimit  int     `json:"usage_limit"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var coupon model.Coupon
	if err := db.First(&coupon, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "优惠券不存在",
			"data": nil,
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
	if req.Value > 0 {
		updates["value"] = req.Value
	}
	if req.MinAmount > 0 {
		updates["min_amount"] = req.MinAmount
	}
	if req.MaxDiscount > 0 {
		updates["max_discount"] = req.MaxDiscount
	}
	if req.UsageLimit > 0 {
		updates["usage_limit"] = req.UsageLimit
	}

	db.Model(&coupon).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteCoupon 删除优惠券
// DELETE /api/admin/coupons/:id
func DeleteCoupon(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var coupon model.Coupon
	if err := db.First(&coupon, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "优惠券不存在",
			"data": nil,
		})
		return
	}

	db.Delete(&coupon)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// UpdateCouponStatus 更新优惠券状态
// PATCH /api/admin/coupons/:id/status
func UpdateCouponStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
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
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var coupon model.Coupon
	if err := db.First(&coupon, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "优惠券不存在",
			"data": nil,
		})
		return
	}

	db.Model(&coupon).Update("status", req.Status)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "状态更新成功",
	})
}

// GetCouponCampaignList 获取优惠券活动列表
// GET /api/admin/coupon-campaigns
func GetCouponCampaignList(c *gin.Context) {
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
	db.Model(&model.CouponCampaign{}).Count(&total)

	var campaigns []model.CouponCampaign
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&campaigns)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      campaigns,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetCouponCampaignSummary 获取优惠券活动统计
// GET /api/admin/coupon-campaigns/summary
func GetCouponCampaignSummary(c *gin.Context) {
	db := database.GetDB()

	var activeCount int64
	db.Model(&model.CouponCampaign{}).Where("status = ?", "active").Count(&activeCount)

	var totalCount int64
	db.Model(&model.CouponCampaign{}).Count(&totalCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"active": activeCount,
			"total":  totalCount,
		},
	})
}

// CreateCouponCampaign 创建优惠券活动
// POST /api/admin/coupon-campaigns
func CreateCouponCampaign(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		CouponID    uint   `json:"coupon_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	campaign := model.CouponCampaign{
		Name:        req.Name,
		Description: req.Description,
		CouponID:    req.CouponID,
		Status:      "active",
	}

	if err := db.Create(&campaign).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建活动失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": campaign.ID,
		},
	})
}

// UpdateCouponCampaign 更新优惠券活动
// PUT /api/admin/coupon-campaigns/:id
func UpdateCouponCampaign(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var campaign model.CouponCampaign
	if err := db.First(&campaign, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "活动不存在",
			"data": nil,
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

	db.Model(&campaign).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
	})
}

// DeleteCouponCampaign 删除优惠券活动
// DELETE /api/admin/coupon-campaigns/:id
func DeleteCouponCampaign(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var campaign model.CouponCampaign
	if err := db.First(&campaign, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "活动不存在",
			"data": nil,
		})
		return
	}

	db.Delete(&campaign)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// UpdateCouponCampaignStatus 更新优惠券活动状态
// PATCH /api/admin/coupon-campaigns/:id/status
func UpdateCouponCampaignStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var campaign model.CouponCampaign
	if err := db.First(&campaign, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "活动不存在",
			"data": nil,
		})
		return
	}

	db.Model(&campaign).Update("status", req.Status)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "状态更新成功",
	})
}
