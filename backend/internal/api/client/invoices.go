package client

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var invoice model.Invoice
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "账单不存在",
			"data": nil,
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
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var invoice model.Invoice
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "账单不存在",
			"data": nil,
		})
		return
	}

	if invoice.Status != "unpaid" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "只有待支付的账单才能取消",
			"data": nil,
		})
		return
	}

	db.Model(&invoice).Update("status", "cancelled")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "取消成功",
		"data":    nil,
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
			"data": nil,
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
			"data": nil,
		})
		return
	}

	if invoice.Status != "unpaid" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "账单不是待支付状态",
			"data": nil,
		})
		return
	}

	// 查询用户余额
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "获取用户信息失败", "data": nil})
		return
	}

	if user.Balance < invoice.Amount {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "余额不足", "data": nil})
		return
	}

	// 安全修复：使用原子操作扣除余额（防并发竞态）
	// UPDATE users SET balance = balance - ? WHERE id = ? AND balance >= ?
	result := db.Model(&model.User{}).Where("id = ? AND balance >= ?", userID, invoice.Amount).
		Update("balance", gorm.Expr("balance - ?", invoice.Amount))
	if result.RowsAffected == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "余额不足（并发冲突）", "data": nil})
		return
	}

	// 更新账单状态
	db.Model(&invoice).Update("status", "paid")

	// 查询最新余额
	var updatedUser model.User
	db.First(&updatedUser, userID)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "支付成功",
		"data": gin.H{
			"balance": updatedUser.Balance,
		},
	})
}

// CombineInvoices 合并账单
// POST /api/client/invoices/combines
func CombineInvoices(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		InvoiceIDs []uint `json:"invoice_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	if len(req.InvoiceIDs) < 2 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "至少选择2个账单才能合并", "data": nil})
		return
	}

	db := database.GetDB()

	// 验证所有账单都属于该用户且都是未支付状态
	var totalAmount float64
	var firstInvoice model.Invoice
	for i, invID := range req.InvoiceIDs {
		var inv model.Invoice
		if err := db.Where("id = ? AND user_id = ?", invID, userID).First(&inv).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 404, "message": "账单不存在", "data": nil})
			return
		}
		if inv.Status != "unpaid" {
			c.JSON(http.StatusOK, gin.H{"code": 400, "message": "只能合并未支付的账单", "data": nil})
			return
		}
		totalAmount += inv.Amount
		if i == 0 {
			firstInvoice = inv
		}
	}

	// 保留第一个账单，把金额加到第一个账单上，删除其他账单
	db.Model(&firstInvoice).Update("amount", totalAmount)

	for _, invID := range req.InvoiceIDs[1:] {
		db.Delete(&model.Invoice{}, invID)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "合并成功",
		"data": gin.H{
			"invoice_id": firstInvoice.ID,
			"total_amount": totalAmount,
		},
	})
}

// FundInvoice 资金支付账单
// POST /api/client/invoices/:id/fund
func FundInvoice(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id := c.Param("id")

	invID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的账单ID", "data": nil})
		return
	}

	db := database.GetDB()
	var invoice model.Invoice
	if err := db.Where("id = ? AND user_id = ?", invID, userID).First(&invoice).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "账单不存在", "data": nil})
		return
	}

	if invoice.Status == "paid" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "账单已支付", "data": nil})
		return
	}

	// 通过PHP插件引擎处理支付
	results, err := pluginengine.TriggerHook("invoice_fund", map[string]interface{}{
		"invoice_id": invoice.ID,
		"user_id":    userID,
		"amount":     invoice.Amount,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 502, "message": "插件引擎离线", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "支付处理中",
		"data": gin.H{"results": results},
	})
}
