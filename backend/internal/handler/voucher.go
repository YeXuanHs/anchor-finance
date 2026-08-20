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

// ==================== Admin 接口 ====================

// GetRate 获取费率配置
// GET /admin/voucher-rate
func (h *VoucherHandler) GetRate(c *gin.Context) {
	data, err := h.voucherSvc.GetRateConfig()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, data)
}

// PostRate 更新费率配置
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
	response.SuccessMsg(c, "请求成功")
}

// GetVoucherList 获取发票申请列表
// GET /admin/voucher-list
func (h *VoucherHandler) GetVoucherList(c *gin.Context) {
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

// GetVoucherDetail 获取发票申请详情
// GET /admin/voucher-detail/:id
func (h *VoucherHandler) GetVoucherDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	detail, err := h.voucherSvc.GetVoucherDetail(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, detail)
}

// PostVoucherStatus 更新发票申请状态
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
		response.BadRequest(c, "参数错误")
		return
	}
	if len(req.Notes) > 500 {
		response.BadRequest(c, "备注不超过500个字符")
		return
	}

	if err := h.voucherSvc.UpdateVoucherStatus(req.ID, req.Status, req.Notes); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "请求成功")
}

// ==================== 用户接口 ====================

// GetUserVoucherList 获取用户发票申请列表
// GET /user/vouchers
func (h *VoucherHandler) GetUserVoucherList(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	vouchers, total, err := h.voucherSvc.GetUserVoucherList(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, vouchers, total, page, pageSize)
}

// CreateUserVoucher 创建发票申请
// POST /user/vouchers
func (h *VoucherHandler) CreateUserVoucher(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.CreateVoucherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	voucher, err := h.voucherSvc.CreateUserVoucher(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, voucher)
}

// ==================== 发票抬头 ====================

// GetVoucherTypes 获取用户发票抬头列表
// GET /user/voucher-types
func (h *VoucherHandler) GetVoucherTypes(c *gin.Context) {
	userID := c.GetUint("user_id")

	types, err := h.voucherSvc.GetVoucherTypes(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, types)
}

// CreateVoucherType 创建发票抬头
// POST /user/voucher-types
func (h *VoucherHandler) CreateVoucherType(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.CreateVoucherTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	vt, err := h.voucherSvc.CreateVoucherType(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, vt)
}

// UpdateVoucherType 更新发票抬头
// PUT /user/voucher-types/:id
func (h *VoucherHandler) UpdateVoucherType(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	var req service.CreateVoucherTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	vt, err := h.voucherSvc.UpdateVoucherType(userID, uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, vt)
}

// DeleteVoucherType 删除发票抬头
// DELETE /user/voucher-types/:id
func (h *VoucherHandler) DeleteVoucherType(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid voucher type id")
		return
	}

	if err := h.voucherSvc.DeleteVoucherType(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// ==================== 收件地址 ====================

// GetVoucherPosts 获取用户收件地址列表
// GET /user/voucher-posts
func (h *VoucherHandler) GetVoucherPosts(c *gin.Context) {
	userID := c.GetUint("user_id")

	posts, err := h.voucherSvc.GetVoucherPosts(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, posts)
}

// CreateVoucherPost 创建收件地址
// POST /user/voucher-posts
func (h *VoucherHandler) CreateVoucherPost(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.CreateVoucherPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	vp, err := h.voucherSvc.CreateVoucherPost(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, vp)
}

// UpdateVoucherPost 更新收件地址
// PUT /user/voucher-posts/:id
func (h *VoucherHandler) UpdateVoucherPost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "ID错误")
		return
	}

	var req service.CreateVoucherPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	vp, err := h.voucherSvc.UpdateVoucherPost(userID, uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, vp)
}

// DeleteVoucherPost 删除收件地址
// DELETE /user/voucher-posts/:id
func (h *VoucherHandler) DeleteVoucherPost(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid voucher post id")
		return
	}

	if err := h.voucherSvc.DeleteVoucherPost(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}
