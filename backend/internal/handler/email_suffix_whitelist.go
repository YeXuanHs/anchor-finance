package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type EmailSuffixWhitelistHandler struct {
	svc *service.EmailSuffixWhitelistService
	log *logger.Logger
}

func NewEmailSuffixWhitelistHandler(svc *service.EmailSuffixWhitelistService, log *logger.Logger) *EmailSuffixWhitelistHandler {
	return &EmailSuffixWhitelistHandler{svc: svc, log: log}
}

// List 获取邮箱后缀白名单列表
// GET /api/admin/email-suffixes?show_inactive=true
func (h *EmailSuffixWhitelistHandler) List(c *gin.Context) {
	showInactive := c.Query("show_inactive") == "true"
	items, err := h.svc.List(showInactive)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, items)
}

// Add 添加邮箱后缀
// POST /api/admin/email-suffixes
func (h *EmailSuffixWhitelistHandler) Add(c *gin.Context) {
	var req struct {
		Suffix string `json:"suffix" binding:"required"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "后缀不能为空")
		return
	}
	if err := h.svc.Add(req.Suffix, req.Remark); err != nil {
		response.BadRequest(c, "添加失败，可能已存在: "+err.Error())
		return
	}
	response.SuccessMsg(c, "添加成功")
}

// Update 更新邮箱后缀
// PUT /api/admin/email-suffixes/:id
func (h *EmailSuffixWhitelistHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		IsActive *bool   `json:"is_active"`
		Remark   *string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请求参数错误")
		return
	}
	if err := h.svc.Update(uint(id), req.IsActive, req.Remark); err != nil {
		response.ServerError(c, "更新失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// Delete 删除邮箱后缀
// DELETE /api/admin/email-suffixes/:id
func (h *EmailSuffixWhitelistHandler) Delete(c *gin.Context) {
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

// BatchDelete 批量删除
// POST /api/admin/email-suffixes/batch-delete
func (h *EmailSuffixWhitelistHandler) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请选择要删除的项")
		return
	}
	if err := h.svc.BatchDelete(req.IDs); err != nil {
		response.ServerError(c, "批量删除失败")
		return
	}
	response.SuccessMsg(c, "批量删除成功")
}

// ImportDefaults 导入默认邮箱后缀
// POST /api/admin/email-suffixes/import-defaults
func (h *EmailSuffixWhitelistHandler) ImportDefaults(c *gin.Context) {
	h.svc.ImportDefaults()
	response.SuccessMsg(c, "默认邮箱后缀导入完成")
}
