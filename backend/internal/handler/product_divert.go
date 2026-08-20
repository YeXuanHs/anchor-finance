package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// ProductDivertHandler handles product transfer HTTP requests.
type ProductDivertHandler struct {
	divertSvc *service.ProductDivertService
	log       *logger.Logger
}

// NewProductDivertHandler creates a new ProductDivertHandler.
func NewProductDivertHandler(divertSvc *service.ProductDivertService, log *logger.Logger) *ProductDivertHandler {
	return &ProductDivertHandler{divertSvc: divertSvc, log: log}
}

// Create creates a product transfer request.
// POST /product-transfers
func (h *ProductDivertHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req service.CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	transfer, err := h.divertSvc.Create(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, transfer)
}

// GetSent returns transfers sent by the authenticated user.
// GET /product-transfers/sent
func (h *ProductDivertHandler) GetSent(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	transfers, total, err := h.divertSvc.GetSent(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, transfers, total, page, pageSize)
}

// GetReceived returns transfers received by the authenticated user.
// GET /product-transfers/received
func (h *ProductDivertHandler) GetReceived(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	transfers, total, err := h.divertSvc.GetReceived(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, transfers, total, page, pageSize)
}

// GetDetail returns a single transfer request.
// GET /product-transfers/:id
func (h *ProductDivertHandler) GetDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid transfer id")
		return
	}

	transfer, err := h.divertSvc.GetByID(userID, uint(id))
	if err != nil {
		response.NotFound(c, "transfer not found")
		return
	}
	response.Success(c, transfer)
}

// Accept accepts a product transfer request.
// POST /product-transfers/:id/accept
func (h *ProductDivertHandler) Accept(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid transfer id")
		return
	}

	if err := h.divertSvc.Accept(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "transfer accepted")
}

// Reject rejects a product transfer request.
// POST /product-transfers/:id/reject
func (h *ProductDivertHandler) Reject(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid transfer id")
		return
	}

	if err := h.divertSvc.Reject(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "transfer rejected")
}

// Cancel cancels a product transfer (by sender).
// POST /product-transfers/:id/cancel
func (h *ProductDivertHandler) Cancel(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid transfer id")
		return
	}

	if err := h.divertSvc.Cancel(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "transfer cancelled")
}

// GetTransferCode 获取转移码
// GET /product-transfers/:id/code
func (h *ProductDivertHandler) GetTransferCode(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid transfer id")
		return
	}

	code, err := h.divertSvc.GetTransferCode(userID, uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"transfer_code": code})
}

// RegenerateCode 重新生成转移码
// POST /product-transfers/:id/regenerate-code
func (h *ProductDivertHandler) RegenerateCode(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid transfer id")
		return
	}

	code, err := h.divertSvc.RegenerateCode(userID, uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"transfer_code": code})
}

// AcceptByCode 通过转移码接受转移
// POST /product-transfers/accept-by-code
func (h *ProductDivertHandler) AcceptByCode(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		TransferCode string `json:"transfer_code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.divertSvc.AcceptByCode(userID, req.TransferCode); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "transfer accepted")
}

// ─── 管理员接口 ───

// AdminGetConfig 管理员获取转移配置
// GET /admin/product-transfer/config
func (h *ProductDivertHandler) AdminGetConfig(c *gin.Context) {
	config := h.divertSvc.GetConfig()
	response.Success(c, config)
}

// AdminSaveConfig 管理员保存转移配置
// PUT /admin/product-transfer/config
func (h *ProductDivertHandler) AdminSaveConfig(c *gin.Context) {
	var config service.TransferConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.divertSvc.SaveConfig(&config); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.SuccessMsg(c, "保存成功")
}

// AdminGetAllTransfers 管理员获取所有转移记录
// GET /admin/product-transfers
func (h *ProductDivertHandler) AdminGetAllTransfers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	var transfers []service.ProductTransfer
	var total int64

	query := h.divertSvc.GetDB().Model(&service.ProductTransfer{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").
		Preload("FromUser").Preload("ToUser").Find(&transfers)

	response.Success(c, gin.H{
		"list":  transfers,
		"total": total,
		"page":  page,
	})
}
