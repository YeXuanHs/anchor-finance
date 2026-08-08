package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct{}

func NewTransactionHandler() *TransactionHandler {
	return &TransactionHandler{}
}

// GetTransactions 获取交易列表
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"list": []interface{}{}, "total": 0})
}

// RegisterRoutes 注册路由
func (h *TransactionHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/transactions", h.GetTransactions)
}
