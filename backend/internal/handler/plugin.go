package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type PluginHandler struct {
	pluginSvc *service.PluginService
	log       *logger.Logger
}

func NewPluginHandler(pluginSvc *service.PluginService, log *logger.Logger) *PluginHandler {
	return &PluginHandler{pluginSvc: pluginSvc, log: log}
}

// GetList returns all plugins with optional status filter.
// GET /plugins
func (h *PluginHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var enabled *bool
	if v := c.Query("enabled"); v != "" {
		b := v == "true" || v == "1"
		enabled = &b
	}

	plugins, total, err := h.pluginSvc.GetList(page, pageSize, enabled)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, plugins, total, page, pageSize)
}

// GetDetail returns a single plugin by ID.
// GET /plugins/:id
func (h *PluginHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid plugin id")
		return
	}

	plugin, err := h.pluginSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "plugin not found")
		return
	}
	response.Success(c, plugin)
}

// Install installs a new plugin.
// POST /plugins/install
func (h *PluginHandler) Install(c *gin.Context) {
	var req service.InstallPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	plugin, err := h.pluginSvc.Install(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, plugin)
}

// Uninstall uninstalls a plugin.
// POST /plugins/:id/uninstall
func (h *PluginHandler) Uninstall(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid plugin id")
		return
	}

	if err := h.pluginSvc.Uninstall(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "plugin uninstalled")
}

// Enable enables a plugin.
// POST /plugins/:id/enable
func (h *PluginHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid plugin id")
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.pluginSvc.Enable(uint(id), adminID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "plugin enabled")
}

// Disable disables a plugin.
// POST /plugins/:id/disable
func (h *PluginHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid plugin id")
		return
	}

	adminID := c.GetUint("user_id")
	if err := h.pluginSvc.Disable(uint(id), adminID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "plugin disabled")
}

// UpdateConfig updates a plugin's configuration.
// PUT /plugins/:id/config
func (h *PluginHandler) UpdateConfig(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid plugin id")
		return
	}

	var req service.UpdatePluginConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	plugin, err := h.pluginSvc.UpdateConfig(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, plugin)
}

// GetLogs returns plugin operation logs.
// GET /plugins/logs
func (h *PluginHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var pluginID *uint
	if pid := c.Query("plugin_id"); pid != "" {
		v, _ := strconv.ParseUint(pid, 10, 64)
		id := uint(v)
		pluginID = &id
	}

	logs, total, err := h.pluginSvc.GetLogs(page, pageSize, pluginID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}
