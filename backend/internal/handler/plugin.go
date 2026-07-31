package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// PluginHandler 插件处理器
type PluginHandler struct {
	svc *service.PluginService
	log *logger.Logger
}

// NewPluginHandler 创建插件处理器
func NewPluginHandler(svc *service.PluginService, log *logger.Logger) *PluginHandler {
	return &PluginHandler{svc: svc, log: log}
}

// ==================== 插件管理 ====================

// List 获取插件列表
func (h *PluginHandler) List(c *gin.Context) {
	pluginType := c.Query("type")

	var isEnabled *bool
	if v := c.Query("is_enabled"); v != "" {
		b := v == "true" || v == "1"
		isEnabled = &b
	}

	items, err := h.svc.GetList(pluginType, isEnabled)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// Install 安装插件
func (h *PluginHandler) Install(c *gin.Context) {
	var req service.CreatePluginRequest
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

// Uninstall 卸载插件
func (h *PluginHandler) Uninstall(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "卸载成功")
}

// Enable 启用插件
func (h *PluginHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.SetEnabled(uint(id), true); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "已启用")
}

// Disable 禁用插件
func (h *PluginHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.SetEnabled(uint(id), false); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "已禁用")
}

// UpdateConfig 更新插件配置
func (h *PluginHandler) UpdateConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var req struct {
		Config string `json:"config"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateConfig(uint(id), req.Config); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "配置已更新")
}

// GetByType 根据类型获取启用的插件
func (h *PluginHandler) GetByType(c *gin.Context) {
	pluginType := c.Param("type")
	if pluginType == "" {
		response.BadRequest(c, "请指定插件类型")
		return
	}

	item, err := h.svc.GetEnabledByType(pluginType)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, item)
}

// GetTypes 获取支持的插件类型
func (h *PluginHandler) GetTypes(c *gin.Context) {
	response.Success(c, model.PluginTypeLabels)
}

// ==================== 服务器模块管理 ====================

// ListServerModules 获取服务器模块列表
func (h *PluginHandler) ListServerModules(c *gin.Context) {
	moduleType := c.Query("module")

	items, err := h.svc.GetServerModuleList(moduleType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// CreateServerModule 创建服务器模块
func (h *PluginHandler) CreateServerModule(c *gin.Context) {
	var req model.ServerModule
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.CreateServerModule(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// UpdateServerModule 更新服务器模块
func (h *PluginHandler) UpdateServerModule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateServerModule(uint(id), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// DeleteServerModule 删除服务器模块
func (h *PluginHandler) DeleteServerModule(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.DeleteServerModule(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// ==================== 服务器分组管理 ====================

// ListServerGroups 获取服务器分组列表
func (h *PluginHandler) ListServerGroups(c *gin.Context) {
	items, err := h.svc.GetServerGroupList()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// CreateServerGroup 创建服务器分组
func (h *PluginHandler) CreateServerGroup(c *gin.Context) {
	var req model.ServerGroup
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.CreateServerGroup(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// UpdateServerGroup 更新服务器分组
func (h *PluginHandler) UpdateServerGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateServerGroup(uint(id), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// DeleteServerGroup 删除服务器分组
func (h *PluginHandler) DeleteServerGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.DeleteServerGroup(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// ==================== OAuth提供商管理 ====================

// ListOAuthProviders 获取OAuth提供商列表
func (h *PluginHandler) ListOAuthProviders(c *gin.Context) {
	items, err := h.svc.GetOAuthProviderList()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}

// CreateOAuthProvider 创建OAuth提供商
func (h *PluginHandler) CreateOAuthProvider(c *gin.Context) {
	var req model.OAuthProvider
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.CreateOAuthProvider(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// UpdateOAuthProvider 更新OAuth提供商
func (h *PluginHandler) UpdateOAuthProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateOAuthProvider(uint(id), updates); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// DeleteOAuthProvider 删除OAuth提供商
func (h *PluginHandler) DeleteOAuthProvider(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}

	if err := h.svc.DeleteOAuthProvider(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// GetEnabledOAuthProviders 获取启用的OAuth提供商
func (h *PluginHandler) GetEnabledOAuthProviders(c *gin.Context) {
	items, err := h.svc.GetEnabledOAuthProviders()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, items)
}
