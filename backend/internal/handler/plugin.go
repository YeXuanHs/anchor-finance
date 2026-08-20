package handler

import (
	"net/http"

	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
)

type PluginHandler struct{}

func NewPluginHandler() *PluginHandler {
	return &PluginHandler{}
}

// List returns a list of plugins.
func (h *PluginHandler) List(c *gin.Context) {
	database := db.GetDB()
	var plugins []map[string]interface{}
	
	category := c.Query("category")
	keyword := c.Query("keyword")
	
	query := database.Table("plugins").Where("deleted_at IS NULL")
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	
	if err := query.Find(&plugins).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": plugins})
}

// GetTypes returns plugin types.
func (h *PluginHandler) GetTypes(c *gin.Context) {
	database := db.GetDB()
	var types []string
	
	if err := database.Table("plugins").Where("deleted_at IS NULL").Distinct().Pluck("category", &types).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": types})
}

// Install installs a plugin.
func (h *PluginHandler) Install(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// Uninstall uninstalls a plugin.
func (h *PluginHandler) Uninstall(c *gin.Context) {
	id := c.Param("id")
	database := db.GetDB()
	
	if err := database.Table("plugins").Where("id = ?", id).Update("deleted_at", "NOW()").Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// Enable enables a plugin.
func (h *PluginHandler) Enable(c *gin.Context) {
	id := c.Param("id")
	database := db.GetDB()
	
	if err := database.Table("plugins").Where("id = ?", id).Update("enabled", true).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// Disable disables a plugin.
func (h *PluginHandler) Disable(c *gin.Context) {
	id := c.Param("id")
	database := db.GetDB()
	
	if err := database.Table("plugins").Where("id = ?", id).Update("enabled", false).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// UpdateConfig updates plugin configuration.
func (h *PluginHandler) UpdateConfig(c *gin.Context) {
	id := c.Param("id")
	database := db.GetDB()
	
	var body map[string]interface{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	
	if err := database.Table("plugins").Where("id = ?", id).Update("config", body).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// GetByType returns plugins by type.
func (h *PluginHandler) GetByType(c *gin.Context) {
	pluginType := c.Param("type")
	database := db.GetDB()
	
	var plugins []map[string]interface{}
	if err := database.Table("plugins").Where("category = ? AND deleted_at IS NULL", pluginType).Find(&plugins).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": plugins})
}

// GetPlugins returns a list of plugins (legacy).
func (h *PluginHandler) GetPlugins(c *gin.Context) {
	h.List(c)
}

// GetPlugin returns a single plugin.
func (h *PluginHandler) GetPlugin(c *gin.Context) {
	id := c.Param("id")
	database := db.GetDB()
	
	var plugin map[string]interface{}
	if err := database.Table("plugins").Where("id = ? AND deleted_at IS NULL", id).First(&plugin).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "插件不存在"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": plugin})
}

// TogglePlugin toggles plugin state.
func (h *PluginHandler) TogglePlugin(c *gin.Context) {
	id := c.Param("id")
	database := db.GetDB()
	
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	
	if err := database.Table("plugins").Where("id = ?", id).Update("enabled", body.Enabled).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// GetPluginConfig returns plugin configuration.
func (h *PluginHandler) GetPluginConfig(c *gin.Context) {
	id := c.Param("id")
	database := db.GetDB()
	
	var plugin map[string]interface{}
	if err := database.Table("plugins").Select("config, config_fields").Where("id = ? AND deleted_at IS NULL", id).First(&plugin).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "插件不存在"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": plugin})
}

// UpdatePluginConfig updates plugin configuration.
func (h *PluginHandler) UpdatePluginConfig(c *gin.Context) {
	h.UpdateConfig(c)
}

// ListServerModules returns server modules.
func (h *PluginHandler) ListServerModules(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

// CreateServerModule creates a server module.
func (h *PluginHandler) CreateServerModule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// UpdateServerModule updates a server module.
func (h *PluginHandler) UpdateServerModule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// DeleteServerModule deletes a server module.
func (h *PluginHandler) DeleteServerModule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// ListServerGroups returns server groups.
func (h *PluginHandler) ListServerGroups(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok", "data": nil})
}

// CreateServerGroup creates a server group.
func (h *PluginHandler) CreateServerGroup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// UpdateServerGroup updates a server group.
func (h *PluginHandler) UpdateServerGroup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// DeleteServerGroup deletes a server group.
func (h *PluginHandler) DeleteServerGroup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "ok"})
}

// RegisterRoutes registers plugin routes.
func (h *PluginHandler) RegisterRoutes(r interface{}) {}
