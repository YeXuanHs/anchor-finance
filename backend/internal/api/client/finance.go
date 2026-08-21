package client

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetBalanceLogs 获取余额日志
// GET /api/client/balance-logs
func GetBalanceLogs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 暂时返回空列表，后续实现balance_logs表
	_ = userID
	_ = page
	_ = pageSize

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      []interface{}{},
			"total":     0,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetRechargeGateways 获取充值网关
// GET /api/client/recharge/gateways
func GetRechargeGateways(c *gin.Context) {
	// 返回可用的支付方式
	gateways := []gin.H{
		{"id": "alipay", "name": "支付宝", "icon": "alipay"},
		{"id": "wxpay", "name": "微信支付", "icon": "wechat"},
		{"id": "balance", "name": "余额支付", "icon": "wallet"},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    gateways,
	})
}

// CreateRecharge 创建充值订单
// POST /api/client/recharge
func CreateRecharge(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Amount  float64 `json:"amount" binding:"required"`
		Gateway string  `json:"gateway" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "充值金额必须大于0",
		})
		return
	}

	// TODO: 创建充值订单，调用支付网关
	_ = userID

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "充值订单创建成功",
		"data": gin.H{
			"payment_no": "PAY202401010001",
			"amount":     req.Amount,
			"gateway":    req.Gateway,
		},
	})
}

// GetUserCoupons 获取用户优惠券
// GET /api/client/coupons
func GetUserCoupons(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 暂时返回空列表，后续实现用户优惠券关联表
	_ = userID
	_ = page
	_ = pageSize

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      []interface{}{},
			"total":     0,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ClaimCoupon 领取优惠券
// POST /api/client/coupons/:id/claim
func ClaimCoupon(c *gin.Context) {
	userID, _ := c.Get("user_id")
	couponID := c.Param("id")

	// TODO: 检查优惠券是否存在、是否已领取、是否过期等
	_ = userID
	_ = couponID

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "领取成功",
	})
}

// GetUserNotifications 获取用户通知
// GET /api/client/notifications
func GetUserNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 暂时返回空列表，后续实现通知表
	_ = userID
	_ = page
	_ = pageSize

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      []interface{}{},
			"total":     0,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetNotices 获取公告列表
// GET /api/client/notices
func GetNotices(c *gin.Context) {
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
	db.Model(&model.News{}).Where("status = ?", "published").Count(&total)

	var notices []model.News
	offset := (page - 1) * pageSize
	db.Where("status = ?", "published").
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&notices)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      notices,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetHelpArticles 获取帮助文章
// GET /api/client/help-articles
func GetHelpArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	categoryID := c.Query("category_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.KnowledgeArticle{}).Where("status = ?", "published")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	var total int64
	query.Count(&total)

	var articles []model.KnowledgeArticle
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&articles)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      articles,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
