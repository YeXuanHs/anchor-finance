package handler

import (
	"anchorfinance/pkg/response"
	"net/http"
	"strconv"
	"time"

	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct{}

func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{}
}

// GetTransactions 获取交易列表
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	txType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	database := db.GetDB()
	if database == nil {
		response.SuccessPage(c, []interface{}{}, 0, page, pageSize)
		return
	}

	query := database.Table("transactions")
	if txType != "" {
		query = query.Where("type = ?", txType)
	}

	var total int64
	query.Count(&total)

	var transactions []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Type      string    `json:"type"`
		Amount    float64   `json:"amount"`
		Remark    string    `json:"remark"`
		CreatedAt time.Time `json:"created_at"`
	}

	query.Select("id, user_id, type, amount, remark, created_at").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&transactions)

	c.JSON(http.StatusOK, gin.H{
		"list":  transactions,
		"total": total,
		"page":  page,
	})
}

// RegisterRoutes 注册路由
func (h *TransactionHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/transactions", h.GetTransactions)
}

// ==================== Admin Router Methods ====================

// GetList returns a paginated transaction list.
func (h *TransactionHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	txType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": gin.H{"list": []interface{}{}, "total": 0}})
		return
	}

	query := database.Table("transactions")
	if txType != "" {
		query = query.Where("type = ?", txType)
	}

	var total int64
	query.Count(&total)

	var transactions []struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Type      string    `json:"type"`
		Amount    float64   `json:"amount"`
		Remark    string    `json:"remark"`
		CreatedAt time.Time `json:"created_at"`
	}

	query.Select("id, user_id, type, amount, remark, created_at").
		Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&transactions)

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"list":  transactions,
			"total": total,
		},
	})
}

// GetDetail returns transaction detail.
func (h *TransactionHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")

	database := db.GetDB()
	if database == nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "data": nil})
		return
	}

	var transaction struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		Type      string    `json:"type"`
		Amount    float64   `json:"amount"`
		Remark    string    `json:"remark"`
		CreatedAt time.Time `json:"created_at"`
	}

	err := database.Table("transactions").Where("id = ?", id).First(&transaction).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "交易不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "data": transaction})
}
