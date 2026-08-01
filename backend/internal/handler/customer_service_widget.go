package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type CustomerServiceWidgetHandler struct {
	svc *service.CustomerServiceWidgetService
	log *logger.Logger
}

func NewCustomerServiceWidgetHandler(svc *service.CustomerServiceWidgetService, log *logger.Logger) *CustomerServiceWidgetHandler {
	return &CustomerServiceWidgetHandler{svc: svc, log: log}
}

// List 获取客服列表
func (h *CustomerServiceWidgetHandler) List(c *gin.Context) {
	showInactive := c.Query("show_inactive") == "true"
	items, err := h.svc.List(showInactive)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, items)
}

// Get 获取单个客服
func (h *CustomerServiceWidgetHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "未找到")
		return
	}
	response.Success(c, item)
}

// Create 创建客服入口
func (h *CustomerServiceWidgetHandler) Create(c *gin.Context) {
	var item model.CustomerServiceWidget
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.Create(&item); err != nil {
		response.ServerError(c, "创建失败")
		return
	}
	response.Success(c, item)
}

// Update 更新客服入口
func (h *CustomerServiceWidgetHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.Update(uint(id), updates); err != nil {
		response.ServerError(c, "更新失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// Delete 删除客服入口
func (h *CustomerServiceWidgetHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, "删除失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// GetSettings 获取全局设置
func (h *CustomerServiceWidgetHandler) GetSettings(c *gin.Context) {
	settings := h.svc.GetSettings()
	response.Success(c, settings)
}

// UpdateSettings 更新全局设置
func (h *CustomerServiceWidgetHandler) UpdateSettings(c *gin.Context) {
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.UpdateSettings(updates); err != nil {
		response.ServerError(c, "更新失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// GetPublic 前台获取客服列表
func (h *CustomerServiceWidgetHandler) GetPublic(c *gin.Context) {
	items, settings := h.svc.GetActiveForDisplay()
	response.Success(c, gin.H{
		"items":    items,
		"settings": settings,
	})
}
