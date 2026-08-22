package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetPromoCodeList 获取优惠码列表
// GET /api/admin/promo-codes
func GetPromoCodeList(c *gin.Context) {
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
	query := db.Model(&model.PromoCode{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("code LIKE ? OR name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	var promoCodes []model.PromoCode
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&promoCodes)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      promoCodes,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreatePromoCode 创建优惠码
// POST /api/admin/promo-codes
func CreatePromoCode(c *gin.Context) {
	var req struct {
		Code        string     `json:"code" binding:"required"`
		Name        string     `json:"name"`
		Type        string     `json:"type" binding:"required"` // percent, fixed
		Value       float64    `json:"value" binding:"required"`
		MinAmount   float64    `json:"min_amount"`
		MaxDiscount float64    `json:"max_discount"`
		UsageLimit  int        `json:"usage_limit"`
		StartDate   *time.Time `json:"start_date"`
		EndDate     *time.Time `json:"end_date"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 验证类型
	if req.Type != "percent" && req.Type != "fixed" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "类型必须是 percent 或 fixed", "data": nil})
		return
	}

	// 0元购防护：优惠码价值校验
	if req.Type == "percent" && (req.Value < 0 || req.Value > 100) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "百分比优惠码价值需在0-100之间", "data": nil})
		return
	}
	if req.Type == "fixed" && req.Value < 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "固定金额优惠码价值不能为负数", "data": nil})
		return
	}

	db := database.GetDB()

	// 检查优惠码是否已存在
	var count int64
	db.Model(&model.PromoCode{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "优惠码已存在",
			"data":    nil,
		})
		return
	}

	promoCode := model.PromoCode{
		Code:        req.Code,
		Name:        req.Name,
		Type:        req.Type,
		Value:       req.Value,
		MinAmount:   req.MinAmount,
		MaxDiscount: req.MaxDiscount,
		UsageLimit:  req.UsageLimit,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Status:      "active",
	}

	if err := db.Create(&promoCode).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建优惠码失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": promoCode.ID,
		},
	})
}

// UpdatePromoCode 更新优惠码
// PUT /api/admin/promo-codes/:id
func UpdatePromoCode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的优惠码ID",
			"data":    nil,
		})
		return
	}

	var req struct {
		Name        string     `json:"name"`
		Type        string     `json:"type"`
		Value       float64    `json:"value"`
		MinAmount   float64    `json:"min_amount"`
		MaxDiscount float64    `json:"max_discount"`
		UsageLimit  int        `json:"usage_limit"`
		StartDate   *time.Time `json:"start_date"`
		EndDate     *time.Time `json:"end_date"`
		Status      string     `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var promoCode model.PromoCode
	if err := db.First(&promoCode, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "优惠码不存在",
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
	if req.StartDate != nil {
		updates["start_date"] = req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = req.EndDate
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&promoCode).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    nil,
	})
}

// DeletePromoCode 删除优惠码
// DELETE /api/admin/promo-codes/:id
func DeletePromoCode(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的优惠码ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var promoCode model.PromoCode
	if err := db.First(&promoCode, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "优惠码不存在",
			"data":    nil,
		})
		return
	}

	db.Delete(&promoCode)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
		"data":    nil,
	})
}
