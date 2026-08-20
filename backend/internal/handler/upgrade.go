package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type UpgradeHandler struct {
	upgradeSvc *service.UpgradeService
	log        *logger.Logger
}

func NewUpgradeHandler(upgradeSvc *service.UpgradeService, log *logger.Logger) *UpgradeHandler {
	return &UpgradeHandler{upgradeSvc: upgradeSvc, log: log}
}

// GetAvailableUpgrades 获取可用的升降级选项
// GET /user/products/:id/upgrades
func (h *UpgradeHandler) GetAvailableUpgrades(c *gin.Context) {
	userID := c.GetUint("user_id")
	productID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid product id")
		return
	}

	products, err := h.upgradeSvc.GetAvailableUpgrades(userID, uint(productID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, products)
}

// CreateUpgrade 创建升降级订单
// POST /upgrades
func (h *UpgradeHandler) CreateUpgrade(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.CreateUpgradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	order, err := h.upgradeSvc.CreateUpgrade(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, order)
}

// GetAvailable is an alias for GetAvailableUpgrades.
func (h *UpgradeHandler) GetAvailable(c *gin.Context) {
	h.GetAvailableUpgrades(c)
}

// Submit is an alias for CreateUpgrade.
func (h *UpgradeHandler) Submit(c *gin.Context) {
	h.CreateUpgrade(c)
}

// GetUpgradeDetail 获取升降级订单详情
// GET /upgrades/:id
func (h *UpgradeHandler) GetUpgradeDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid upgrade id")
		return
	}

	order, err := h.upgradeSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "upgrade order not found")
		return
	}
	response.Success(c, order)
}

// GetUserUpgrades 获取用户的升降级订单列表
// GET /user/upgrades
func (h *UpgradeHandler) GetUserUpgrades(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.upgradeSvc.GetUserUpgrades(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, orders, total, page, pageSize)
}

// PayUpgrade 支付升降级订单
// POST /upgrades/:id/pay
func (h *UpgradeHandler) PayUpgrade(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid upgrade id")
		return
	}

	order, err := h.upgradeSvc.PayAndApply(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, order)
}
