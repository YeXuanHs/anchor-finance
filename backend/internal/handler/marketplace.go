package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// MarketplaceHandler 交易市场处理器
type MarketplaceHandler struct {
	svc *service.MarketplaceService
	log *logger.Logger
}

func NewMarketplaceHandler(svc *service.MarketplaceService, log *logger.Logger) *MarketplaceHandler {
	return &MarketplaceHandler{svc: svc, log: log}
}

// ─── 挂售管理 ───

// CreateListing 创建挂售
func (h *MarketplaceHandler) CreateListing(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		HostID      uint    `json:"host_id" binding:"required"`
		SellPrice   float64 `json:"sell_price" binding:"required,gt=0"`
		Description string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	listing, err := h.svc.CreateListing(userID, req.HostID, req.SellPrice, req.Description)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, listing)
}

// UpdateListing 更新挂售
func (h *MarketplaceHandler) UpdateListing(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	var req struct {
		SellPrice   float64 `json:"sell_price" binding:"required,gt=0"`
		Description string  `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.svc.UpdateListing(userID, uint(id), req.SellPrice, req.Description); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "更新成功")
}

// RemoveListing 下架挂售
func (h *MarketplaceHandler) RemoveListing(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.svc.RemoveListing(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "下架成功")
}

// GetListing 获取挂售详情
func (h *MarketplaceHandler) GetListing(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	listing, err := h.svc.GetListing(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, listing)
}

// GetListings 获取挂售列表
func (h *MarketplaceHandler) GetListings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	listings, total, err := h.svc.GetListings(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	response.Success(c, gin.H{
		"list":  listings,
		"total": total,
		"page":  page,
	})
}

// GetUserListings 获取用户的挂售列表
func (h *MarketplaceHandler) GetUserListings(c *gin.Context) {
	userID := c.GetUint("user_id")

	listings, err := h.svc.GetUserListings(userID)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	response.Success(c, listings)
}

// ─── 订单管理 ───

// CreateOrder 创建订单
func (h *MarketplaceHandler) CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		ListingID     uint   `json:"listing_id" binding:"required"`
		PaymentMethod string `json:"payment_method" binding:"required"` // full=全额 fee_only=仅手续费
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	order, err := h.svc.CreateOrder(userID, req.ListingID, req.PaymentMethod)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, order)
}

// GetBuyerOrders 获取买家订单
func (h *MarketplaceHandler) GetBuyerOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.svc.GetBuyerOrders(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	response.Success(c, gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
	})
}

// GetSellerOrders 获取卖家订单
func (h *MarketplaceHandler) GetSellerOrders(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.svc.GetSellerOrders(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	response.Success(c, gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
	})
}

// CompleteOrder 完成订单
func (h *MarketplaceHandler) CompleteOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.svc.CompleteOrder(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "订单已完成")
}

// CancelOrder 取消订单
func (h *MarketplaceHandler) CancelOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.svc.CancelOrder(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "订单已取消")
}

// ─── 私聊功能 ───

// SendMessage 发送私聊消息
func (h *MarketplaceHandler) SendMessage(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		ReceiverID  uint   `json:"receiver_id" binding:"required"`
		ListingID   uint   `json:"listing_id" binding:"required"`
		Content     string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	message, err := h.svc.SendMessage(userID, req.ReceiverID, req.ListingID, req.Content)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, message)
}

// GetChatMessages 获取聊天消息
func (h *MarketplaceHandler) GetChatMessages(c *gin.Context) {
	userID := c.GetUint("user_id")
	listingID, _ := strconv.ParseUint(c.Param("listing_id"), 10, 32)
	otherUserID, _ := strconv.ParseUint(c.Param("user_id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	messages, total, err := h.svc.GetChatMessages(userID, uint(listingID), uint(otherUserID), page, pageSize)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	response.Success(c, gin.H{
		"list":  messages,
		"total": total,
		"page":  page,
	})
}

// GetChatSessions 获取聊天会话列表
func (h *MarketplaceHandler) GetChatSessions(c *gin.Context) {
	userID := c.GetUint("user_id")

	sessions, err := h.svc.GetChatSessions(userID)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	response.Success(c, sessions)
}

// GetUnreadCount 获取未读消息数
func (h *MarketplaceHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	count := h.svc.GetUnreadCount(userID)

	response.Success(c, gin.H{"unread_count": count})
}

// ─── 管理员接口 ───

// AdminGetConfig 管理员获取市场配置
func (h *MarketplaceHandler) AdminGetConfig(c *gin.Context) {
	config := h.svc.GetConfig()
	response.Success(c, config)
}

// AdminSaveConfig 管理员保存市场配置
func (h *MarketplaceHandler) AdminSaveConfig(c *gin.Context) {
	var req struct {
		Enabled           *bool    `json:"enabled"`
		FeeRate           *float64 `json:"fee_rate"`
		MinFee            *float64 `json:"min_fee"`
		MaxListingDays    *int     `json:"max_listing_days"`
		MinHoldDays       *int     `json:"min_hold_days"`
		RequireRealName   *bool    `json:"require_real_name"`
		AllowFeeOnly      *bool    `json:"allow_fee_only"`
		AutoTransfer      *bool    `json:"auto_transfer"`
		NotifyEmail       *bool    `json:"notify_email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	config := h.svc.GetConfig()
	if req.Enabled != nil {
		config.Enabled = *req.Enabled
	}
	if req.FeeRate != nil {
		config.FeeRate = *req.FeeRate
	}
	if req.MinFee != nil {
		config.MinFee = *req.MinFee
	}
	if req.MaxListingDays != nil {
		config.MaxListingDays = *req.MaxListingDays
	}
	if req.MinHoldDays != nil {
		config.MinHoldDays = *req.MinHoldDays
	}
	if req.RequireRealName != nil {
		config.RequireRealName = *req.RequireRealName
	}
	if req.AllowFeeOnly != nil {
		config.AllowFeeOnly = *req.AllowFeeOnly
	}
	if req.AutoTransfer != nil {
		config.AutoTransfer = *req.AutoTransfer
	}
	if req.NotifyEmail != nil {
		config.NotifyEmail = *req.NotifyEmail
	}

	if err := h.svc.SaveConfig(config); err != nil {
		response.ServerError(c, "保存失败")
		return
	}

	response.SuccessMsg(c, "保存成功")
}

// AdminGetAllListings 管理员获取所有挂售
func (h *MarketplaceHandler) AdminGetAllListings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	var listings []model.MarketplaceListing
	var total int64

	query := h.svc.GetDB().Model(&model.MarketplaceListing{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").
		Preload("User").Preload("Host").Find(&listings)

	response.Success(c, gin.H{
		"list":  listings,
		"total": total,
		"page":  page,
	})
}

// AdminGetAllOrders 管理员获取所有订单
func (h *MarketplaceHandler) AdminGetAllOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	var orders []model.MarketplaceOrder
	var total int64

	query := h.svc.GetDB().Model(&model.MarketplaceOrder{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").
		Preload("Listing").Preload("Buyer").Preload("Seller").Find(&orders)

	response.Success(c, gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
	})
}

// ─── 交易记录/收益/提现/日志 ───

// GetTransactions returns the current user's marketplace transactions.
func (h *MarketplaceHandler) GetTransactions(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.svc.GetTransactions(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	response.Success(c, gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
	})
}

// GetTransactionSummary returns summary stats for the user's transactions.
func (h *MarketplaceHandler) GetTransactionSummary(c *gin.Context) {
	userID := c.GetUint("user_id")

	summary, err := h.svc.GetTransactionSummary(userID)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	response.Success(c, summary)
}

// GetEarnings returns earnings data for the current user (as seller).
func (h *MarketplaceHandler) GetEarnings(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.svc.GetEarnings(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	response.Success(c, gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
	})
}

// GetWithdrawals returns withdrawal records for the current user.
func (h *MarketplaceHandler) GetWithdrawals(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := h.svc.GetWithdrawals(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	response.Success(c, gin.H{
		"list":  logs,
		"total": total,
		"page":  page,
	})
}

// CreateWithdrawal creates a withdrawal request.
func (h *MarketplaceHandler) CreateWithdrawal(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Amount float64 `json:"amount" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.svc.CreateWithdrawal(userID, req.Amount); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "提现申请已提交")
}

// GetLogs returns marketplace operation logs for the current user.
func (h *MarketplaceHandler) GetLogs(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.svc.GetLogs(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}

	response.Success(c, gin.H{
		"list":  orders,
		"total": total,
		"page":  page,
	})
}

// PayOrder pays a marketplace order.
func (h *MarketplaceHandler) PayOrder(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if err := h.svc.PayOrder(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessMsg(c, "支付成功")
}
