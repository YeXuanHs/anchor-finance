package handler

import (
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type AIShoppingCoreHandler struct {
	svc *service.AIShoppingCoreService
	log *logger.Logger
}

func NewAIShoppingCoreHandler(svc *service.AIShoppingCoreService, log *logger.Logger) *AIShoppingCoreHandler {
	return &AIShoppingCoreHandler{svc: svc, log: log}
}

// ─── 后台配置 ───

// GetConfig 获取所有配置
func (h *AIShoppingCoreHandler) GetConfig(c *gin.Context) {
	response.Success(c, h.svc.GetAllConfig())
}

// SaveConfig 保存配置
func (h *AIShoppingCoreHandler) SaveConfig(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.SaveAllConfig(data); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.SuccessMsg(c, "保存成功")
}

// ─── 前台聊天 ───

// Chat 发送消息（支持页面上下文）
func (h *AIShoppingCoreHandler) Chat(c *gin.Context) {
	sessionID := c.Param("session_id")
	var req struct {
		Message     string `json:"message" binding:"required"`
		PageContext string `json:"page_context"` // 页面上下文
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "消息不能为空")
		return
	}
	userID := c.GetUint("user_id")
	reply, err := h.svc.Chat(sessionID, userID, req.Message, req.PageContext)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"reply": reply})
}

// GetChatHistory 获取聊天历史
func (h *AIShoppingCoreHandler) GetChatHistory(c *gin.Context) {
	sessionID := c.Param("session_id")
	messages := h.svc.GetChatHistory(sessionID)
	response.Success(c, messages)
}
