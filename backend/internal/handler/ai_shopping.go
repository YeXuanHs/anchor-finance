package handler

import (
	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type AIShoppingHandler struct {
	svc *service.AIShoppingService
	log *logger.Logger
}

func NewAIShoppingHandler(svc *service.AIShoppingService, log *logger.Logger) *AIShoppingHandler {
	return &AIShoppingHandler{svc: svc, log: log}
}

// ─── 配置管理（后台） ───

// GetConfig 获取 AI 购物助手配置
func (h *AIShoppingHandler) GetConfig(c *gin.Context) {
	config := h.svc.GetConfig()
	response.Success(c, config)
}

// SaveConfig 保存 AI 购物助手配置
func (h *AIShoppingHandler) SaveConfig(c *gin.Context) {
	var config model.AIShoppingAssistantConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.SaveConfig(&config); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.SuccessMsg(c, "保存成功")
}

// GetCatalogConfig 获取商品目录配置
func (h *AIShoppingHandler) GetCatalogConfig(c *gin.Context) {
	config := h.svc.GetCatalogConfig()
	response.Success(c, config)
}

// SaveCatalogConfig 保存商品目录配置
func (h *AIShoppingHandler) SaveCatalogConfig(c *gin.Context) {
	var config model.ProductCatalogConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.SaveCatalogConfig(&config); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.SuccessMsg(c, "保存成功")
}

// ─── 前台对话 ───

// StartSession 开始购物助手会话
func (h *AIShoppingHandler) StartSession(c *gin.Context) {
	userID := c.GetUint("user_id")
	sessionID, err := h.svc.StartSession(userID)
	if err != nil {
		response.ServerError(c, "创建会话失败")
		return
	}
	response.Success(c, gin.H{"session_id": sessionID})
}

// SendMessage 发送消息
func (h *AIShoppingHandler) SendMessage(c *gin.Context) {
	sessionID := c.Param("session_id")
	var req struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "消息不能为空")
		return
	}
	reply, recommendations, err := h.svc.SendMessage(sessionID, req.Message)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"reply":           reply,
		"recommendations": recommendations,
	})
}

// CloseSession 关闭会话
func (h *AIShoppingHandler) CloseSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	if err := h.svc.CloseSession(sessionID); err != nil {
		response.ServerError(c, "关闭会话失败")
		return
	}
	response.SuccessMsg(c, "会话已关闭")
}

// GetSessionMessages 获取会话消息
func (h *AIShoppingHandler) GetSessionMessages(c *gin.Context) {
	sessionID := c.Param("session_id")
	messages, err := h.svc.GetSessionMessages(sessionID)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, messages)
}
