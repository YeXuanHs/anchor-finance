package client

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetBalanceLogs 获取余额日志（从充值记录+操作日志合并）
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

	db := database.GetDB()

	// 查询该用户的充值记录（真实数据）
	var total int64
	db.Model(&model.Recharge{}).Where("user_id = ?", userID).Count(&total)

	var recharges []model.Recharge
	offset := (page - 1) * pageSize
	db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&recharges)

	if recharges == nil {
		recharges = []model.Recharge{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      recharges,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetRechargeGateways 获取充值网关（从plugins表读payment域）
// GET /api/client/recharge/gateways
func GetRechargeGateways(c *gin.Context) {
	db := database.GetDB()
	var gateways []model.Plugin
	db.Where("domain = ? AND status = ?", "payment", "active").Order("name ASC").Find(&gateways)

	if gateways == nil {
		gateways = []model.Plugin{}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": gateways})
}

// CreateRecharge 创建充值订单（生成真实交易号+存库）
// POST /api/client/recharge
func CreateRecharge(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Amount  float64 `json:"amount" binding:"required"`
		Gateway string  `json:"gateway" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	// 0元购防护
	if req.Amount <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "充值金额必须大于0", "data": nil})
		return
	}

	// 生成交易号
	paymentNo := fmt.Sprintf("RCH%s%06d", time.Now().Format("20060102"), time.Now().UnixNano()%1000000)

	db := database.GetDB()
	recharge := model.Recharge{
		UserID:        userID.(uint),
		Amount:        req.Amount,
		Gateway:       req.Gateway,
		TransactionNo: paymentNo,
		Status:        "pending",
	}

	if err := db.Create(&recharge).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建充值订单失败", "data": nil})
		return
	}

	// 调用PHP插件引擎创建支付（获取支付URL）
	result, err := pluginengine.TriggerHook("create_payment", map[string]interface{}{
		"invoice_id":   recharge.ID,
		"amount":       req.Amount,
		"gateway":      req.Gateway,
		"payment_no":   paymentNo,
		"user_id":      userID,
		"type":         "recharge",
	})
	if err != nil {
		// 插件引擎离线，返回502
		c.JSON(http.StatusBadGateway, gin.H{"code": 502, "message": "支付服务暂时不可用", "data": nil})
		return
	}

	// 从PHP插件返回结果中获取支付URL
	var paymentURL string
	if len(result) > 0 {
		if url, ok := result[0].Result.(map[string]interface{}); ok {
			paymentURL, _ = url["pay_url"].(string)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "充值订单创建成功",
		"data": gin.H{
			"recharge_id": recharge.ID,
			"payment_no":  paymentNo,
			"amount":      req.Amount,
			"gateway":     req.Gateway,
			"pay_url":     paymentURL,
		},
	})
}

// GetUserCoupons 获取用户可用的优惠券（真实数据）
// GET /api/client/coupons
func GetUserCoupons(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	now := time.Now()

	// 查询所有有效且未过期的优惠券
	query := db.Model(&model.Coupon{}).
		Where("status = ? AND (start_date IS NULL OR start_date <= ?) AND (end_date IS NULL OR end_date >= ?)",
			"active", now, now)

	var total int64
	query.Count(&total)

	var coupons []model.Coupon
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&coupons)

	if coupons == nil {
		coupons = []model.Coupon{}
	}

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

// ClaimCoupon 领取优惠券（校验有效期+使用次数）
// POST /api/client/coupons/:id/claim
func ClaimCoupon(c *gin.Context) {
	couponID := c.Param("id")
	userID, _ := c.Get("user_id")

	db := database.GetDB()
	var coupon model.Coupon
	if err := db.First(&coupon, couponID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "优惠券不存在", "data": nil})
		return
	}

	// 校验状态
	if coupon.Status != "active" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "优惠券已失效", "data": nil})
		return
	}

	// 校验有效期
	now := time.Now()
	if coupon.StartDate != nil && now.Before(*coupon.StartDate) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "优惠券未开始", "data": nil})
		return
	}
	if coupon.EndDate != nil && now.After(*coupon.EndDate) {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "优惠券已过期", "data": nil})
		return
	}

	// 学创欧：检查是否已领取过（user_coupons表）
	var existingCount int64
	db.Model(&model.UserCoupon{}).Where("coupon_id = ? AND user_id = ?", coupon.ID, userID).Count(&existingCount)
	if existingCount > 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "你已经领取过这张优惠券", "data": nil})
		return
	}

	// 校验领取上限
	if coupon.UsageLimit > 0 {
		var claimedCount int64
		db.Model(&model.UserCoupon{}).Where("coupon_id = ?", coupon.ID).Count(&claimedCount)
		if claimedCount >= int64(coupon.UsageLimit) {
			c.JSON(http.StatusOK, gin.H{"code": 400, "message": "优惠券已被领完", "data": nil})
			return
		}
	}

	// 创建领取记录
	claimedAt := time.Now()
	userCoupon := model.UserCoupon{
		CouponID:    coupon.ID,
		UserID:      userID.(uint),
		ReceiveType: "claim",
		Status:      1,
		ClaimedAt:   &claimedAt,
	}
	if err := db.Create(&userCoupon).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "领取失败，可能已领取过", "data": nil})
		return
	}

	// 更新优惠券使用次数
	db.Model(&coupon).Update("used_count", gorm.Expr("used_count + ?", 1))

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "领取成功", "data": gin.H{"coupon_id": coupon.ID, "user_coupon_id": userCoupon.ID}})
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

	db := database.GetDB()
	var total int64
	db.Model(&model.UserNotification{}).Where("user_id = ?", userID).Count(&total)

	var notifications []model.UserNotification
	offset := (page - 1) * pageSize
	db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&notifications)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      notifications,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetNotificationUnreadCount 获取未读通知数量
// GET /api/client/notifications/unread-count
func GetNotificationUnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	db := database.GetDB()
	var count int64
	db.Model(&model.UserNotification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"count": count,
		},
	})
}

// MarkNotificationRead 标记通知已读
// PUT /api/client/notifications/:id/read-state
func MarkNotificationRead(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	db := database.GetDB()
	var notification model.UserNotification
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&notification).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "通知不存在",
			"data": nil,
		})
		return
	}

	db.Model(&notification).Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "标记成功",
	})
}

// MarkAllNotificationsRead 标记所有通知已读
// POST /api/client/notifications/mark-all-read
func MarkAllNotificationsRead(c *gin.Context) {
	userID, _ := c.Get("user_id")

	db := database.GetDB()
	db.Model(&model.UserNotification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "全部标记成功",
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
