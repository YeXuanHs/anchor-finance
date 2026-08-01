package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

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

// GetRate 获取代金券费率配置
// GET /admin/voucher-rate
func (h *VoucherHandler) GetRate(c *gin.Context) {
	data, err := h.voucherSvc.GetRateConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, data)
}

// PostRate 更新代金券费率配置
// POST /admin/voucher-rate
func (h *VoucherHandler) PostRate(c *gin.Context) {
	var req struct {
		VoucherManager int     `json:"voucher_manager"`
		Rate           float64 `json:"rate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.voucherSvc.UpdateRateConfig(req.VoucherManager, req.Rate); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "rate config updated")
}

// GetExpressList 获取快递方式列表
// GET /admin/expresses
func (h *VoucherHandler) GetExpressList(c *gin.Context) {
	expresses, total, err := h.voucherSvc.GetExpressList()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"express": expresses, "total": total})
}

// GetExpress 获取单个快递方式
// GET /admin/expresses/:id
func (h *VoucherHandler) GetExpress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid express id")
		return
	}

	express, err := h.voucherSvc.GetExpressByID(uint(id))
	if err != nil {
		response.NotFound(c, "express not found")
		return
	}
	response.Success(c, gin.H{"express": express})
}

// PostExpress 创建/更新快递方式
// POST /admin/expresses
func (h *VoucherHandler) PostExpress(c *gin.Context) {
	var req struct {
		ID    *uint   `json:"id"`
		Name  string  `json:"name" binding:"required"`
		Price float64 `json:"price"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if len(req.Name) > 50 {
		response.BadRequest(c, "name too long (max 50 chars)")
		return
	}
	if req.Price < 0 {
		response.BadRequest(c, "price must be non-negative")
		return
	}

	if err := h.voucherSvc.SaveExpress(req.ID, req.Name, req.Price); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "express saved")
}

// DeleteExpress 删除快递方式
// DELETE /admin/expresses/:id
func (h *VoucherHandler) DeleteExpress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid express id")
		return
	}

	if err := h.voucherSvc.DeleteExpress(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "express deleted")
}

// AdminGetVoucherList 获取代金券列表(含发票信息)
// GET /admin/voucher-list
func (h *VoucherHandler) AdminGetVoucherList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	order := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "DESC")

	vouchers, total, err := h.voucherSvc.GetVoucherList(page, pageSize, status, order, sort)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, vouchers, total, page, pageSize)
}

// AdminGetVoucherDetail 获取代金券详情(含发票)
// GET /admin/voucher-detail/:id
func (h *VoucherHandler) AdminGetVoucherDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid voucher id")
		return
	}

	detail, err := h.voucherSvc.GetVoucherDetail(uint(id))
	if err != nil {
		response.NotFound(c, "voucher not found")
		return
	}
	response.Success(c, detail)
}

// PostVoucherStatus 更新代金券状态(审核/拒绝)
// POST /admin/voucher-status
func (h *VoucherHandler) PostVoucherStatus(c *gin.Context) {
	var req struct {
		ID     uint   `json:"id" binding:"required"`
		Status string `json:"status" binding:"required"`
		Notes  string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Status != "Reject" && req.Status != "Send" {
		response.BadRequest(c, "invalid status, must be Reject or Send")
		return
	}
	if len(req.Notes) > 500 {
		response.BadRequest(c, "notes too long (max 500 chars)")
		return
	}

	if err := h.voucherSvc.UpdateVoucherStatus(req.ID, req.Status, req.Notes); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "voucher status updated")
}
