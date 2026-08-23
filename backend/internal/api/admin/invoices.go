package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
	"github.com/gin-gonic/gin"
)

// GetInvoiceList 获取账单列表
// GET /api/admin/invoices
func GetInvoiceList(c *gin.Context) {
	// 1. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	userID := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 2. 构建查询
	db := database.GetDB()
	query := db.Model(&model.Invoice{})

	// 关键词搜索（账单号）
	if keyword != "" {
		query = query.Where("invoice_no LIKE ?", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 用户筛选
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 3. 获取总数
	var total int64
	query.Count(&total)

	// 4. 分页查询
	var invoices []model.Invoice
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&invoices)

	// 5. 返回统一格式
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

// GetInvoice 获取账单详情
// GET /api/admin/invoices/:id
func GetInvoice(c *gin.Context) {
	// 1. 获取账单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的账单ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询账单
	db := database.GetDB()
	var invoice model.Invoice
	if err := db.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "账单不存在",
			"data":    nil,
		})
		return
	}

	// 3. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    invoice,
	})
}

// CancelInvoice 取消账单
// POST /api/admin/invoices/:id/cancel
func CancelInvoice(c *gin.Context) {
	// 1. 获取账单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的账单ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询账单
	db := database.GetDB()
	var invoice model.Invoice
	if err := db.First(&invoice, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "账单不存在",
			"data":    nil,
		})
		return
	}

	// 3. 检查状态
	if invoice.Status == "cancelled" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "账单已取消",
			"data":    nil,
		})
		return
	}

	if invoice.Status == "paid" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "已支付的账单不能取消",
			"data":    nil,
		})
		return
	}

	// 4. 取消账单
	if err := db.Model(&invoice).Update("status", "cancelled").Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "取消账单失败",
			"data":    nil,
		})
		return
	}

	// 触发Hook: invoice_mark_cancelled
	pluginengine.TriggerHook("invoice_mark_cancelled", map[string]interface{}{
		"invoice_id": invoice.ID, "user_id": invoice.UserID,
	})

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "取消成功",
		"data":    nil,
	})
}

// GetTransactionList 获取交易流水
// GET /api/admin/transactions
func GetTransactionList(c *gin.Context) {
	// 1. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	userID := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 2. 构建查询（从支付记录表读取真实交易数据）
	db := database.GetDB()
	query := db.Model(&model.Payment{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	var total int64
	query.Count(&total)

	var transactions []model.Payment
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&transactions)

	if transactions == nil {
		transactions = []model.Payment{}
	}

	// 3. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      transactions,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
