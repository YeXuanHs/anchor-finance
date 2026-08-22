package client

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PaymentNotify 通用支付回调（转发给PHP插件引擎处理）
// POST /api/client/payment/notify/:gateway
func PaymentNotify(c *gin.Context) {
	gateway := c.Param("gateway")

	callbackParams := make(map[string]interface{})
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			callbackParams[key] = values[0]
		}
	}
	if err := c.Request.ParseForm(); err == nil {
		for key, values := range c.Request.PostForm {
			if len(values) > 0 {
				callbackParams[key] = values[0]
			}
		}
	}

	result, err := pluginengine.HandlePaymentCallback(gateway, callbackParams)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"code": 502, "message": "支付处理服务暂时不可用", "data": nil})
		return
	}

	status, _ := result["status"].(string)
	if status == "success" {
		invoiceID, _ := result["invoice_id"].(float64)
		transactionID, _ := result["transaction_id"].(string)
		if invoiceID > 0 {
			db := database.GetDB()
			var invoice model.Invoice
			if err := db.First(&invoice, uint(invoiceID)).Error; err == nil {
				now := time.Now()
				db.Model(&invoice).Updates(map[string]interface{}{
					"status":         "paid",
					"paid_at":        &now,
					"payment_method": gateway,
				})
				db.Create(&model.Payment{
					UserID:        invoice.UserID,
					InvoiceID:     invoice.ID,
					Amount:        invoice.Amount,
					Gateway:       gateway,
					TransactionNo: transactionID,
					Status:        "completed",
				})
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

// GetRechargeStatus 获取充值状态
// GET /api/client/recharge/:paymentNo/status
func GetRechargeStatus(c *gin.Context) {
	paymentNo := c.Param("paymentNo")
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var recharge model.Recharge
	if err := db.Where("transaction_no = ? AND user_id = ?", paymentNo, userID).First(&recharge).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "充值记录不存在", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{
			"status":  recharge.Status,
			"amount":  recharge.Amount,
			"paid_at": recharge.PaidAt,
		},
	})
}

// GetBalanceLogsSummary 余额日志汇总
// GET /api/client/balance-logs/summary
func GetBalanceLogsSummary(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var totalRecharge float64
	db.Model(&model.Recharge{}).Where("user_id = ? AND status = ?", userID, "completed").
		Select("COALESCE(SUM(amount), 0)").Scan(&totalRecharge)

	var user model.User
	db.First(&user, userID)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{
			"balance":        user.Balance,
			"total_recharge": totalRecharge,
		},
	})
}

// GetNotificationPreferences 获取通知偏好设置
// GET /api/client/auth/notification-preferences
func GetNotificationPreferences(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var user model.User
	db.First(&user, userID)

	defaults := map[string]bool{
		"email_enabled":  true,
		"sms_enabled":    true,
		"system_enabled": true,
	}

	var settings []model.Setting
	db.Where("`group` = ?", "notification").Find(&settings)
	for _, s := range settings {
		if s.Value == "1" {
			defaults[s.Key] = true
		} else if s.Value == "0" {
			defaults[s.Key] = false
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": defaults})
}

// UpdateNotificationPreferences 更新通知偏好设置
// PUT /api/client/auth/notification-preferences
func UpdateNotificationPreferences(c *gin.Context) {
	var req map[string]bool
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	for key, val := range req {
		value := "0"
		if val {
			value = "1"
		}
		db.Where("`key` = ? AND `group` = ?", key, "notification").
			Assign(model.Setting{Value: value}).
			FirstOrCreate(&model.Setting{})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": nil})
}

// GetOrderSummary 订单汇总
// GET /api/client/orders/summary
func GetOrderSummary(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var totalOrders int64
	db.Model(&model.Order{}).Where("user_id = ?", userID).Count(&totalOrders)

	var pendingOrders int64
	db.Model(&model.Order{}).Where("user_id = ? AND status = ?", userID, "pending").Count(&pendingOrders)

	var activeOrders int64
	db.Model(&model.Order{}).Where("user_id = ? AND status = ?", userID, "active").Count(&activeOrders)

	var totalSpent float64
	db.Model(&model.Order{}).Where("user_id = ? AND status IN ?", userID, []string{"active", "completed"}).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalSpent)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{
			"total_orders":   totalOrders,
			"pending_orders": pendingOrders,
			"active_orders":  activeOrders,
			"total_spent":    totalSpent,
		},
	})
}

// GetTicketServiceOptions 工单服务选项
// GET /api/client/tickets/service-options
func GetTicketServiceOptions(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var services []model.Service
	db.Where("user_id = ? AND status = ?", userID, "active").Find(&services)

	type ServiceOption struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	options := make([]ServiceOption, 0)
	for _, s := range services {
		options = append(options, ServiceOption{ID: s.ID, Name: s.Domain})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": options})
}

// allowedImageExtensions 允许上传的图片扩展名白名单
var allowedImageExtensions = map[string]bool{
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true, ".bmp": true,
}

// UploadTicketImages 上传工单图片
// POST /api/client/tickets/upload-images
func UploadTicketImages(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请选择文件", "data": nil})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "文件大小不能超过5MB", "data": nil})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedImageExtensions[ext] {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "不支持的文件类型，仅允许 jpg/png/gif/webp", "data": nil})
		return
	}

	safeFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	savePath := filepath.Join("uploads/tickets", safeFilename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "上传失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "上传成功",
		"data": gin.H{"url": "/" + savePath, "filename": safeFilename},
	})
}

// GetUserDirectReferrals 获取直接推荐用户
// GET /api/client/referral/direct-referrals
func GetUserDirectReferrals(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var referrals []model.Referral
	db.Where("referrer_id = ?", userID).Order("id DESC").Find(&referrals)

	type ReferralInfo struct {
		ID        uint      `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		Reward    float64   `json:"reward"`
	}

	result := make([]ReferralInfo, 0)
	for _, r := range referrals {
		var user model.User
		if err := db.First(&user, r.UserID).Error; err == nil {
			result = append(result, ReferralInfo{
				ID:        user.ID,
				Username:  user.Username,
				Email:     user.Email,
				CreatedAt: r.CreatedAt,
				Reward:    r.Reward,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": result})
}

// GetUserReferralAccountLogs 推介账户日志
// GET /api/client/referral/account-logs
func GetUserReferralAccountLogs(c *gin.Context) {
	userID, _ := c.Get("user_id")
	db := database.GetDB()

	var logs []model.OperationLog
	db.Where("user_id = ? AND action LIKE ?", userID, "referral%").Order("id DESC").Limit(50).Find(&logs)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": logs})
}
