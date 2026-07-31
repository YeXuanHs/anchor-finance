package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type PaymentGatewayHandler struct {
	svc *service.PaymentGatewayService
	log *logger.Logger
}

func NewPaymentGatewayHandler(svc *service.PaymentGatewayService, log *logger.Logger) *PaymentGatewayHandler {
	return &PaymentGatewayHandler{svc: svc, log: log}
}

// List 返回支付方式列表
func (h *PaymentGatewayHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var isEnabled *bool
	if v := c.Query("is_enabled"); v != "" {
		b := v == "true" || v == "1"
		isEnabled = &b
	}

	items, total, err := h.svc.GetList(page, pageSize, isEnabled)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetDetail 返回单个支付方式详情
func (h *PaymentGatewayHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "支付方式不存在")
		return
	}
	response.Success(c, item)
}

// Create 创建支付方式
func (h *PaymentGatewayHandler) Create(c *gin.Context) {
	var req service.CreatePaymentGatewayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Create(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// Update 更新支付方式
func (h *PaymentGatewayHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req service.UpdatePaymentGatewayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Update(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// Delete 删除支付方式
func (h *PaymentGatewayHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// SetStatus 切换启用状态
func (h *PaymentGatewayHandler) SetStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.ToggleStatus(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "状态已切换")
}

// TestConnection 测试支付接口
func (h *PaymentGatewayHandler) TestConnection(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.TestConnection(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "接口测试通过"})
}

// GetEnabled 返回启用的支付方式（用户前台）
func (h *PaymentGatewayHandler) GetEnabled(c *gin.Context) {
	items, err := h.svc.GetEnabled()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// 转换为前端需要的格式
	var result []map[string]interface{}
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":    item.ID,
			"name":  item.Name,
			"title": item.Title,
			"code":  item.Code,
			"icon":  item.Icon,
		})
	}
	response.Success(c, result)
}

// GetSupportedInfo 返回支持的接口和支付类型
func (h *PaymentGatewayHandler) GetSupportedInfo(c *gin.Context) {
	info := h.svc.GetSupportedInfo()
	response.Success(c, info)
}

// 以下是路由别名（兼容旧路由）
func (h *PaymentGatewayHandler) AdminGetList(c *gin.Context)    { h.List(c) }
func (h *PaymentGatewayHandler) AdminGetDetail(c *gin.Context)  { h.GetDetail(c) }
func (h *PaymentGatewayHandler) AdminCreate(c *gin.Context)     { h.Create(c) }
func (h *PaymentGatewayHandler) AdminUpdate(c *gin.Context)     { h.Update(c) }
func (h *PaymentGatewayHandler) AdminDelete(c *gin.Context)     { h.Delete(c) }
func (h *PaymentGatewayHandler) AdminToggleStatus(c *gin.Context) { h.SetStatus(c) }
func (h *PaymentGatewayHandler) AdminUpdateSort(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req struct {
		SortOrder int `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.svc.UpdateSort(uint(id), req.SortOrder); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "排序已更新")
}
