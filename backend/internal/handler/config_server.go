package handler

import (
	"strconv"
	"strings"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type ConfigServerHandler struct {
	svc *service.ConfigServerService
	log *logger.Logger
}

func NewConfigServerHandler(svc *service.ConfigServerService, log *logger.Logger) *ConfigServerHandler {
	return &ConfigServerHandler{svc: svc, log: log}
}

// ---------- Server Config ----------

// GetList returns paginated server configs.
func (h *ConfigServerHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	serverType := c.Query("type")
	keyword := c.Query("keyword")

	items, total, err := h.svc.GetList(page, pageSize, serverType, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// GetDetail returns a single server config by ID.
func (h *ConfigServerHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server config id")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "server config not found")
		return
	}
	response.Success(c, item)
}

// Create creates a server config.
func (h *ConfigServerHandler) Create(c *gin.Context) {
	var req service.CreateServerConfigRequest
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

// Update updates a server config.
func (h *ConfigServerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server config id")
		return
	}

	var req service.UpdateServerConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.Update(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// Delete deletes a server config.
func (h *ConfigServerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server config id")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// BatchUpdateStatus batch-updates status for server configs.
func (h *ConfigServerHandler) BatchUpdateStatus(c *gin.Context) {
	var req struct {
		IDs    []uint `json:"ids" binding:"required,min=1"`
		Status int16  `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.BatchUpdateStatus(req.IDs, req.Status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "请求成功")
}

// BatchDelete batch-deletes server configs.
func (h *ConfigServerHandler) BatchDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.BatchDelete(req.IDs); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// UpdateSort updates sort order for a server config.
func (h *ConfigServerHandler) UpdateSort(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid server config id")
		return
	}

	var req struct {
		SortOrder int `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.UpdateSort(uint(id), req.SortOrder); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "请求成功")
}

// ---------- Server Template ----------

// GetTemplateList returns paginated server templates.
func (h *ConfigServerHandler) GetTemplateList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	templateType := c.Query("type")

	items, total, err := h.svc.GetTemplateList(page, pageSize, templateType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, items, total, page, pageSize)
}

// CreateTemplate creates a server template.
func (h *ConfigServerHandler) CreateTemplate(c *gin.Context) {
	var req service.CreateServerTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.CreateTemplate(req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// UpdateTemplate updates a server template.
func (h *ConfigServerHandler) UpdateTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	var req service.UpdateServerTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	item, err := h.svc.UpdateTemplate(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

// DeleteTemplate deletes a server template.
func (h *ConfigServerHandler) DeleteTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	if err := h.svc.DeleteTemplate(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// helper: convert comma-separated string to uint slice
func parseIDs(s string) []uint {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	ids := make([]uint, 0, len(parts))
	for _, p := range parts {
		if v, err := strconv.ParseUint(strings.TrimSpace(p), 10, 64); err == nil {
			ids = append(ids, uint(v))
		}
	}
	return ids
}

// ─── ConfigServers Admin Methods (from zjmf ConfigServersController) ───

// ServerList returns paginated servers with filters.
func (h *ConfigServerHandler) ServerList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	gid, _ := strconv.ParseUint(c.Query("gid"), 10, 64)
	search := c.Query("search")

	items, total, err := h.svc.AdminServerList(page, pageSize, uint(gid), search)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"data": items, "count": total})
}

// AddServers returns data for the add server page.
func (h *ConfigServerHandler) AddServers(c *gin.Context) {
	data := h.svc.AdminAddServersData()
	response.Success(c, data)
}

// GetModulesGroup returns server groups filtered by module type.
func (h *ConfigServerHandler) GetModulesGroup(c *gin.Context) {
	moduleType := c.Query("modules")
	groups := h.svc.AdminGetModulesGroup(moduleType)
	response.Success(c, gin.H{"groups": groups})
}

// AddServersPost creates a new server.
func (h *ConfigServerHandler) AddServersPost(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		IPAddress   string `json:"ip_address"`
		Hostname    string `json:"hostname"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		AccessHash  string `json:"accesshash"`
		Type        string `json:"type" binding:"required"`
		Port        int    `json:"port"`
		MaxAccounts int    `json:"max_accounts"`
		GID         int    `json:"gid"`
		Secure      int    `json:"secure"`
		Disabled    int    `json:"disabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	serverID, err := h.svc.AdminAddServersPost(req.Name, req.IPAddress, req.Hostname, req.Username, req.Password, req.AccessHash, req.Type, req.Port, req.MaxAccounts, req.GID, req.Secure, req.Disabled)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"server_id": serverID})
}

// EditServers returns server detail for editing.
func (h *ConfigServerHandler) EditServers(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.BadRequest(c, "invalid server id")
		return
	}

	server, err := h.svc.AdminEditServers(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	data := h.svc.AdminAddServersData()
	data["server"] = server
	response.Success(c, data)
}

// EditServersPost updates a server.
func (h *ConfigServerHandler) EditServersPost(c *gin.Context) {
	var req struct {
		ID          uint   `json:"id" binding:"required"`
		Name        string `json:"name" binding:"required"`
		IPAddress   string `json:"ip_address"`
		Hostname    string `json:"hostname"`
		Username    string `json:"username"`
		Password    string `json:"password"`
		AccessHash  string `json:"accesshash"`
		Type        string `json:"type" binding:"required"`
		Port        int    `json:"port"`
		MaxAccounts int    `json:"max_accounts"`
		GID         int    `json:"gid"`
		Secure      int    `json:"secure"`
		Disabled    int    `json:"disabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AdminEditServersPost(req.ID, req.Name, req.IPAddress, req.Hostname, req.Username, req.Password, req.AccessHash, req.Type, req.Port, req.MaxAccounts, req.GID, req.Secure, req.Disabled); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "请求成功")
}

// DeleteServers deletes a server.
func (h *ConfigServerHandler) DeleteServers(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		id, _ = strconv.ParseUint(c.Query("id"), 10, 64)
	}
	if id == 0 {
		response.BadRequest(c, "invalid server id")
		return
	}

	if err := h.svc.AdminDeleteServers(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// GroupsList returns paginated server groups.
func (h *ConfigServerHandler) GroupsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	items, total, err := h.svc.AdminServerGroupsList(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"data": items, "count": total})
}

// CreateGroups returns data for creating a server group.
func (h *ConfigServerHandler) CreateGroups(c *gin.Context) {
	data := h.svc.AdminCreateGroupPage()
	response.Success(c, data)
}

// CreateGroupsPost creates a new server group.
func (h *ConfigServerHandler) CreateGroupsPost(c *gin.Context) {
	var req struct {
		GroupName string `json:"group_name" binding:"required"`
		Mode      int    `json:"mode"`
		SID       []uint `json:"sid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AdminCreateGroupPost(req.GroupName, req.Mode, req.SID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server group created")
}

// EditServerGroups returns data for editing a server group.
func (h *ConfigServerHandler) EditServerGroups(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.BadRequest(c, "invalid group id")
		return
	}

	data, err := h.svc.AdminEditServerGroup(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, data)
}

// EditServerGroupsPost updates a server group.
func (h *ConfigServerHandler) EditServerGroupsPost(c *gin.Context) {
	var req struct {
		ID        uint   `json:"id" binding:"required"`
		GroupName string `json:"group_name" binding:"required"`
		Mode      int    `json:"mode"`
		SID       []uint `json:"sid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.svc.AdminEditServerGroupPost(req.ID, req.GroupName, req.Mode, req.SID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server group updated")
}

// DeleteServerGroups deletes a server group.
func (h *ConfigServerHandler) DeleteServerGroups(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		id, _ = strconv.ParseUint(c.Query("id"), 10, 64)
	}
	if id == 0 {
		response.BadRequest(c, "invalid group id")
		return
	}

	if err := h.svc.AdminDeleteServerGroup(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "server group deleted")
}

// TestLink tests connection to a server.
func (h *ConfigServerHandler) TestLink(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if id == 0 {
		response.BadRequest(c, "invalid server id")
		return
	}

	result, err := h.svc.AdminTestLink(uint(id))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}
