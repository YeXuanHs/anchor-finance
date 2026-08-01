package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type CSChatHandler struct {
	svc *service.CSChatService
	log *logger.Logger
}

func NewCSChatHandler(svc *service.CSChatService, log *logger.Logger) *CSChatHandler {
	return &CSChatHandler{svc: svc, log: log}
}

// ─── 配置 ───

// GetAIConfig 获取AI配置
func (h *CSChatHandler) GetAIConfig(c *gin.Context) {
	response.Success(c, h.svc.GetAIConfig())
}

// SaveAIConfig 保存AI配置
func (h *CSChatHandler) SaveAIConfig(c *gin.Context) {
	var req struct {
		APIEndpoint string `json:"api_endpoint"`
		APIKey      string `json:"api_key"`
		Model       string `json:"model"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.SaveAIConfig(req.APIEndpoint, req.APIKey, req.Model); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.SuccessMsg(c, "保存成功")
}

// GetAppearanceConfig 获取外观配置
func (h *CSChatHandler) GetAppearanceConfig(c *gin.Context) {
	response.Success(c, h.svc.GetAppearanceConfig())
}

// SaveAppearanceConfig 保存外观配置
func (h *CSChatHandler) SaveAppearanceConfig(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	h.svc.SaveAppearanceConfig(data)
	response.SuccessMsg(c, "保存成功")
}

// GetWorkingHours 获取工作时间配置
func (h *CSChatHandler) GetWorkingHours(c *gin.Context) {
	response.Success(c, h.svc.GetWorkingHoursConfig())
}

// SaveWorkingHours 保存工作时间配置
func (h *CSChatHandler) SaveWorkingHours(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	h.svc.SaveWorkingHoursConfig(data)
	response.SuccessMsg(c, "保存成功")
}

// ─── 会话 ───

// ListSessions 获取会话列表
func (h *CSChatHandler) ListSessions(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	items, total, err := h.svc.ListSessions(status, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{"items": items, "total": total, "page": page})
}

// GetSession 获取会话详情
func (h *CSChatHandler) GetSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	session, err := h.svc.GetSession(uint(id))
	if err != nil {
		response.NotFound(c, "会话不存在")
		return
	}
	messages, _ := h.svc.GetSessionMessages(uint(id))
	response.Success(c, gin.H{"session": session, "messages": messages})
}

// SendReply 发送回复
func (h *CSChatHandler) SendReply(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "回复内容不能为空")
		return
	}
	if err := h.svc.SendReply(uint(id), req.Content); err != nil {
		response.ServerError(c, "发送失败")
		return
	}
	response.SuccessMsg(c, "发送成功")
}

// TransferToHuman 转人工
func (h *CSChatHandler) TransferToHuman(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	staffID := c.GetUint("user_id")
	if err := h.svc.TransferToHuman(uint(id), staffID); err != nil {
		response.ServerError(c, "操作失败")
		return
	}
	response.SuccessMsg(c, "已转人工")
}

// CloseSession 关闭会话
func (h *CSChatHandler) CloseSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.svc.CloseSession(uint(id)); err != nil {
		response.ServerError(c, "操作失败")
		return
	}
	response.SuccessMsg(c, "已关闭")
}

// RateSession 评价
func (h *CSChatHandler) RateSession(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var req struct {
		Rating  int    `json:"rating"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	h.svc.RateSession(uint(id), req.Rating, req.Comment)
	response.SuccessMsg(c, "评价成功")
}

// GetStats 获取统计
func (h *CSChatHandler) GetStats(c *gin.Context) {
	response.Success(c, h.svc.GetStats())
}
