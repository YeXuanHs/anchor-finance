package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type AITicketHandler struct {
	svc *service.AITicketService
	log *logger.Logger
}

func NewAITicketHandler(svc *service.AITicketService, log *logger.Logger) *AITicketHandler {
	return &AITicketHandler{svc: svc, log: log}
}

// ─── AI 配置管理 ───

// ListAIConfigs 获取 AI 配置列表
func (h *AITicketHandler) ListAIConfigs(c *gin.Context) {
	configs, err := h.svc.ListAIConfigs()
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, configs)
}

// GetAIConfig 获取单个 AI 配置
func (h *AITicketHandler) GetAIConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	config, err := h.svc.GetAIConfig(uint(id))
	if err != nil {
		response.NotFound(c, "未找到")
		return
	}
	response.Success(c, config)
}

// SaveAIConfig 保存 AI 配置（创建或更新）
func (h *AITicketHandler) SaveAIConfig(c *gin.Context) {
	var config model.AIConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.SaveAIConfig(&config); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.Success(c, config)
}

// DeleteAIConfig 删除 AI 配置
func (h *AITicketHandler) DeleteAIConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.svc.DeleteAIConfig(uint(id)); err != nil {
		response.ServerError(c, "删除失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// ─── 自动回复配置 ───

// GetAutoReplyConfig 获取自动回复配置
func (h *AITicketHandler) GetAutoReplyConfig(c *gin.Context) {
	config := h.svc.GetAutoReplyConfig()
	response.Success(c, config)
}

// SaveAutoReplyConfig 保存自动回复配置
func (h *AITicketHandler) SaveAutoReplyConfig(c *gin.Context) {
	var config model.AITicketAutoReplyConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.SaveAutoReplyConfig(&config); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.SuccessMsg(c, "保存成功")
}

// ─── 自动回复日志 ───

// GetAutoReplyLogs 获取自动回复日志
func (h *AITicketHandler) GetAutoReplyLogs(c *gin.Context) {
	ticketID, _ := strconv.ParseUint(c.Query("ticket_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := h.svc.GetAutoReplyLogs(uint(ticketID), page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{
		"items": logs,
		"total": total,
		"page":  page,
	})
}

// MarkReplyAccepted 标记回复是否被接受
func (h *AITicketHandler) MarkReplyAccepted(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req struct {
		Accepted bool `json:"accepted"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.MarkReplyAccepted(uint(id), req.Accepted); err != nil {
		response.ServerError(c, "操作失败")
		return
	}
	response.SuccessMsg(c, "操作成功")
}

// TestAutoReply 测试自动回复
func (h *AITicketHandler) TestAutoReply(c *gin.Context) {
	var req service.AutoReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	result := h.svc.GenerateAutoReply(req)
	response.Success(c, result)
}
