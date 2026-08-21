package client

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetUserInvoices 获取用户账单列表
// GET /api/client/invoices
func GetUserInvoices(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Invoice{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var invoices []model.Invoice
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&invoices)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      invoices,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUserInvoice 获取用户账单详情
// GET /api/client/invoices/:id
func GetUserInvoice(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的账单ID",
		})
		return
	}

	db := database.GetDB()
	var invoice model.Invoice
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "账单不存在",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    invoice,
	})
}

// GetInvoiceSummary 获取账单统计
// GET /api/client/invoices/summary
func GetInvoiceSummary(c *gin.Context) {
	userID, _ := c.Get("user_id")

	db := database.GetDB()

	// 统计待支付金额
	var unpaidAmount float64
	db.Model(&model.Invoice{}).Where("user_id = ? AND status = ?", userID, "unpaid").
		Select("COALESCE(SUM(amount), 0)").Scan(&unpaidAmount)

	// 统计已支付金额
	var paidAmount float64
	db.Model(&model.Invoice{}).Where("user_id = ? AND status = ?", userID, "paid").
		Select("COALESCE(SUM(amount), 0)").Scan(&paidAmount)

	// 统计账单数量
	var unpaidCount int64
	db.Model(&model.Invoice{}).Where("user_id = ? AND status = ?", userID, "unpaid").Count(&unpaidCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"unpaid_amount": unpaidAmount,
			"paid_amount":   paidAmount,
			"unpaid_count":  unpaidCount,
		},
	})
}

// CancelUserInvoice 取消用户账单
// POST /api/client/invoices/:id/cancellations
func CancelUserInvoice(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的账单ID",
		})
		return
	}

	db := database.GetDB()
	var invoice model.Invoice
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "账单不存在",
		})
		return
	}

	if invoice.Status != "unpaid" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "只有待支付的账单才能取消",
		})
		return
	}

	db.Model(&invoice).Update("status", "cancelled")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "取消成功",
	})
}

// PayInvoiceByBalance 余额支付账单
// POST /api/client/invoices/:id/pay/balance
func PayInvoiceByBalance(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的账单ID",
		})
		return
	}

	db := database.GetDB()

	// 查询账单
	var invoice model.Invoice
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "账单不存在",
		})
		return
	}

	if invoice.Status != "unpaid" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "账单不是待支付状态",
		})
		return
	}

	// 查询用户余额
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "获取用户信息失败",
		})
		return
	}

	if user.Balance < invoice.Amount {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "余额不足",
		})
		return
	}

	// 扣除余额
	newBalance := user.Balance - invoice.Amount
	db.Model(&user).Update("balance", newBalance)

	// 更新账单状态
	db.Model(&invoice).Update("status", "paid")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "支付成功",
		"data": gin.H{
			"balance": newBalance,
		},
	})
}
