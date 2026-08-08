package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/anchorfinance/backend/internal/model"
	"github.com/anchorfinance/backend/internal/service"
)

type PluginHandler struct {
	pluginService *service.PluginService
}

func NewPluginHandler() *PluginHandler {
	return &PluginHandler{
		pluginService: service.NewPluginService(),
	}
}

// GetPlugins 获取插件列表
func (h *PluginHandler) GetPlugins(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "100"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	params := &model.PluginQueryParams{
		Page:     page,
		PageSize: pageSize,
		Category: category,
		Keyword:  keyword,
	}

	result, err := h.pluginService.GetPlugins(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetPlugin 获取单个插件
func (h *PluginHandler) GetPlugin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	plugin, err := h.pluginService.GetPluginByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "插件不存在"})
		return
	}

	c.JSON(http.StatusOK, plugin)
}

// TogglePlugin 切换插件状态
func (h *PluginHandler) TogglePlugin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req model.TogglePluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.pluginService.TogglePlugin(uint(id), req.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status := "已禁用"
	if req.Enabled {
		status = "已启用"
	}
	c.JSON(http.StatusOK, gin.H{"message": status})
}

// GetPluginConfig 获取插件配置
func (h *PluginHandler) GetPluginConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	config, err := h.pluginService.GetPluginConfig(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdatePluginConfig 更新插件配置
func (h *PluginHandler) UpdatePluginConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
		return
	}

	var req model.UpdatePluginConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.pluginService.UpdatePluginConfig(uint(id), req.Config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "配置保存成功"})
}

// RegisterRoutes 注册路由
func (h *PluginHandler) RegisterRoutes(r *gin.RouterGroup) {
	plugin := r.Group("/plugins")
	{
		plugin.GET("", h.GetPlugins)
		plugin.GET("/:id", h.GetPlugin)
		plugin.POST("/:id/toggle", h.TogglePlugin)
		plugin.GET("/:id/config", h.GetPluginConfig)
		plugin.POST("/:id/config", h.UpdatePluginConfig)
	}
}
