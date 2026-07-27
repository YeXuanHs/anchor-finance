package handler

import (
	"strconv"

	"github.com/anchor-finance/backend/internal/service"
	"github.com/anchor-finance/backend/pkg/logger"
	"github.com/anchor-finance/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type VoucherHandler struct {
	voucherSvc *service.VoucherService
	log        *logger.Logger
}

func NewVoucherHandler(voucherSvc *service.VoucherService, log *logger.Logger) *VoucherHandler {
	return &VoucherHandler{voucherSvc: voucherSvc, log: log}
}

// GetUserVouchers 获取用户可用的代金券
// GET /user/vouchers
func (h *VoucherHandler) GetUserVouchers(c *gin.Context) {
	userID := c.GetUint("user_id")
	vouchers, err := h.voucherSvc.GetUserVouchers(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, vouchers)
}

// ClaimVoucher 用户领取代金券
// POST /vouchers/:id/claim
func (h *VoucherHandler) ClaimVoucher(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid voucher id")
		return
	}

	if err := h.voucherSvc.ClaimVoucher(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "voucher claimed")
}

// ValidateVoucher 验证代金券是否可用于订单
// POST /vouchers/validate
func (h *VoucherHandler) ValidateVoucher(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.ValidateVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	discount, voucher, err := h.voucherSvc.ValidateVoucher(userID, req.Code, req.ProductID, req.Amount)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"discount": discount,
		"voucher":  voucher,
	})
}

// AdminGetList 获取所有代金券 (admin)
// GET /admin/vouchers
func (h *VoucherHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	vouchers, total, err := h.voucherSvc.GetAll(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, vouchers, total, page, pageSize)
}

// AdminCreate 创建代金券 (admin)
// POST /admin/vouchers
func (h *VoucherHandler) AdminCreate(c *gin.Context) {
	var req service.CreateVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	voucher, err := h.voucherSvc.Create(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, voucher)
}

// AdminUpdate 更新代金券 (admin)
// PUT /admin/vouchers/:id
func (h *VoucherHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid voucher id")
		return
	}

	var req service.UpdateVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	voucher, err := h.voucherSvc.Update(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, voucher)
}

// AdminDelete 删除代金券 (admin)
// DELETE /admin/vouchers/:id
func (h *VoucherHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid voucher id")
		return
	}

	if err := h.voucherSvc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "voucher deleted")
}
