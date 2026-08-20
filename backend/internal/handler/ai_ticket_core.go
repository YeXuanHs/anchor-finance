package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type AITicketCoreHandler struct {
	svc *service.AITicketCoreService
	log *logger.Logger
}

func NewAITicketCoreHandler(svc *service.AITicketCoreService, log *logger.Logger) *AITicketCoreHandler {
	return &AITicketCoreHandler{svc: svc, log: log}
}

// ─── 控制台配置 ───

// GetDashboard 获取控制台配置
func (h *AITicketCoreHandler) GetDashboard(c *gin.Context) {
	response.Success(c, h.svc.GetDashboardConfig())
}

// SaveDashboard 保存控制台配置
func (h *AITicketCoreHandler) SaveDashboard(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	h.svc.SaveDashboardConfig(data)
	response.SuccessMsg(c, "保存成功")
}

// ─── 知识库 ───

// ListKnowledge 获取知识库列表
func (h *AITicketCoreHandler) ListKnowledge(c *gin.Context) {
	keyword := c.Query("keyword")
	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}
	items, err := h.svc.ListKnowledge(status, keyword)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, items)
}

// CreateKnowledge 创建知识条目
func (h *AITicketCoreHandler) CreateKnowledge(c *gin.Context) {
	var item model.AITicketKnowledge
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.CreateKnowledge(&item); err != nil {
		response.ServerError(c, "创建失败")
		return
	}
	response.Success(c, item)
}

// UpdateKnowledge 更新知识条目
func (h *AITicketCoreHandler) UpdateKnowledge(c *gin.Context) {
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
	if err := h.svc.UpdateKnowledge(uint(id), updates); err != nil {
		response.ServerError(c, "更新失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// DeleteKnowledge 删除知识条目
func (h *AITicketCoreHandler) DeleteKnowledge(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.svc.DeleteKnowledge(uint(id)); err != nil {
		response.ServerError(c, "删除失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// ImportDefaultKnowledge 导入默认知识库
func (h *AITicketCoreHandler) ImportDefaultKnowledge(c *gin.Context) {
	count := h.svc.ImportDefaultKnowledge()
	response.SuccessMsg(c, "导入完成")
	c.JSON(200, gin.H{"code": 0, "message": "导入完成", "count": count})
}

// ─── 规则 ───

// ListRules 获取规则列表
func (h *AITicketCoreHandler) ListRules(c *gin.Context) {
	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}
	items, err := h.svc.ListRules(status)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, items)
}

// CreateRule 创建规则
func (h *AITicketCoreHandler) CreateRule(c *gin.Context) {
	var item model.AITicketRule
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.CreateRule(&item); err != nil {
		response.ServerError(c, "创建失败")
		return
	}
	response.Success(c, item)
}

// UpdateRule 更新规则
func (h *AITicketCoreHandler) UpdateRule(c *gin.Context) {
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
	if err := h.svc.UpdateRule(uint(id), updates); err != nil {
		response.ServerError(c, "更新失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// DeleteRule 删除规则
func (h *AITicketCoreHandler) DeleteRule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.svc.DeleteRule(uint(id)); err != nil {
		response.ServerError(c, "删除失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// ─── 队列 ───

// ListQueue 获取队列
func (h *AITicketCoreHandler) ListQueue(c *gin.Context) {
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	items, total, err := h.svc.ListQueue(status, page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{"items": items, "total": total, "page": page})
}

// GetQueueStats 获取队列统计
func (h *AITicketCoreHandler) GetQueueStats(c *gin.Context) {
	response.Success(c, h.svc.GetQueueStats())
}

// ─── 日志 ───

// ListProcessLogs 获取处理日志
func (h *AITicketCoreHandler) ListProcessLogs(c *gin.Context) {
	ticketID, _ := strconv.ParseUint(c.Query("ticket_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	items, total, err := h.svc.ListProcessLogs(uint(ticketID), page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{"items": items, "total": total, "page": page})
}

// ListNotifyLogs 获取通知日志
func (h *AITicketCoreHandler) ListNotifyLogs(c *gin.Context) {
	ticketID, _ := strconv.ParseUint(c.Query("ticket_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	items, total, err := h.svc.ListNotifyLogs(uint(ticketID), page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{"items": items, "total": total, "page": page})
}

// ─── 工单模式 ───

// GetTicketMode 获取工单AI模式
func (h *AITicketCoreHandler) GetTicketMode(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("ticket_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的工单ID")
		return
	}
	mode := h.svc.GetTicketMode(uint(ticketID))
	response.Success(c, gin.H{"mode": mode})
}

// SetTicketMode 设置工单AI模式
func (h *AITicketCoreHandler) SetTicketMode(c *gin.Context) {
	ticketID, err := strconv.ParseUint(c.Param("ticket_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的工单ID")
		return
	}
	var req struct {
		Mode string `json:"mode" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请选择模式")
		return
	}
	if err := h.svc.SetTicketMode(uint(ticketID), req.Mode); err != nil {
		response.ServerError(c, "设置失败")
		return
	}
	response.SuccessMsg(c, "设置成功")
}

// ─── 工具管理 ───

// ListTools 获取工具列表（按分类）
func (h *AITicketCoreHandler) ListTools(c *gin.Context) {
	categories := h.svc.ListTools()
	response.Success(c, categories)
}

// SetToolEnabled 启用/禁用工具
func (h *AITicketCoreHandler) SetToolEnabled(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.BadRequest(c, "缺少工具名称")
		return
	}

	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := h.svc.SetToolEnabled(name, req.Enabled); err != nil {
		response.ServerError(c, "设置失败")
		return
	}

	msg := "工具已启用"
	if !req.Enabled {
		msg = "工具已禁用"
	}
	response.SuccessMsg(c, msg)
}

// ListToolExecutionLogs 获取工具执行日志
func (h *AITicketCoreHandler) ListToolExecutionLogs(c *gin.Context) {
	ticketID, _ := strconv.ParseUint(c.Query("ticket_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 { page = 1 }
	if pageSize < 1 || pageSize > 100 { pageSize = 20 }
	items, total, err := h.svc.ListToolExecutionLogs(uint(ticketID), page, pageSize)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{"items": items, "total": total, "page": page})
}

// ─── 测试 ───

// TestAutoReply 测试自动回复
func (h *AITicketCoreHandler) TestAutoReply(c *gin.Context) {
	var req service.TicketAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	result := h.svc.ProcessTicket(req)
	response.Success(c, result)
}
