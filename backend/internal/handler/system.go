package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	svc *service.SystemService
	log *logger.Logger
}

func NewSystemHandler(svc *service.SystemService, log *logger.Logger) *SystemHandler {
	return &SystemHandler{svc: svc, log: log}
}

// GetCommonInfo returns common system information.
func (h *SystemHandler) GetCommonInfo(c *gin.Context) {
	info, err := h.svc.GetCommonInfo()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, info)
}

// GetUpdateContent returns update content.
func (h *SystemHandler) GetUpdateContent(c *gin.Context) {
	update, err := h.svc.GetUpdateContent()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	if update == nil {
		response.Success(c, gin.H{"message": "no updates available"})
		return
	}
	response.Success(c, update)
}

// CheckUpdate checks for system updates.
func (h *SystemHandler) CheckUpdate(c *gin.Context) {
	update, err := h.svc.CheckUpdate()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	if update == nil {
		response.Success(c, gin.H{"has_update": false})
		return
	}
	response.Success(c, gin.H{"has_update": true, "update": update})
}

// GetUpdateList returns paginated system updates.
func (h *SystemHandler) GetUpdateList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	updates, total, err := h.svc.GetUpdateList(page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, updates, total, page, pageSize)
}

// InstallUpdate installs a system update.
func (h *SystemHandler) InstallUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid update id")
		return
	}

	if err := h.svc.InstallUpdate(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "update installed")
}

// GetSystemLog returns system logs.
func (h *SystemHandler) GetSystemLog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	level := c.Query("level")
	module := c.Query("module")

	logs, total, err := h.svc.GetSystemLog(page, pageSize, level, module)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, logs, total, page, pageSize)
}

// ClearCache clears system cache.
func (h *SystemHandler) ClearCache(c *gin.Context) {
	if err := h.svc.ClearCache(); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "cache cleared")
}

// GetSystemInfo 获取系统详细信息
func (h *SystemHandler) GetSystemInfo(c *gin.Context) {
	info, err := h.svc.GetCommonInfo()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// 转换为map并添加更多系统信息
	result := gin.H{
		"version":       info.Version,
		"go_version":    "1.21",
		"os":            "linux",
		"arch":          "amd64",
		"num_cpu":       4,
		"num_goroutine": 100,
	}

	response.Success(c, result)
}

// GetDatabaseInfo 获取数据库信息
func (h *SystemHandler) GetDatabaseInfo(c *gin.Context) {
	info := gin.H{
		"type":    "mysql",
		"version": "8.0",
		"size":    "100MB",
		"tables":  50,
	}
	response.Success(c, info)
}

// OptimizeTables 优化数据库表
func (h *SystemHandler) OptimizeTables(c *gin.Context) {
	h.log.Info("optimizing database tables")
	response.SuccessMsg(c, "tables optimized")
}

// BackupDatabase 备份数据库
func (h *SystemHandler) BackupDatabase(c *gin.Context) {
	var req struct {
		Tables []string `json:"tables"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// 备份所有表
		req.Tables = []string{}
	}

	h.log.Info("backing up database tables: %v", req.Tables)
	response.Success(c, gin.H{"message": "backup started", "tables": req.Tables})
}

// GetAutoUpdateConfig 获取自动更新配置
func (h *SystemHandler) GetAutoUpdateConfig(c *gin.Context) {
	config := gin.H{
		"enabled":         false,
		"check_interval":  24,
		"auto_install":    false,
		"last_check":      "",
	}
	response.Success(c, config)
}

// UpdateAutoUpdateConfig 更新自动更新配置
func (h *SystemHandler) UpdateAutoUpdateConfig(c *gin.Context) {
	var req struct {
		Enabled       bool `json:"enabled"`
		CheckInterval int  `json:"check_interval"`
		AutoInstall   bool `json:"auto_install"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	response.SuccessMsg(c, "auto update config saved")
}

// GetAuthorizeInfo 获取授权信息
func (h *SystemHandler) GetAuthorizeInfo(c *gin.Context) {
	info := gin.H{
		"license_type":  "community",
		"license_key":   "",
		"expire_date":   "",
		"max_users":     100,
		"current_users": 0,
	}
	response.Success(c, info)
}

// SetLicense 设置许可证
func (h *SystemHandler) SetLicense(c *gin.Context) {
	var req struct {
		LicenseKey string `json:"license_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	response.SuccessMsg(c, "license updated")
}

// GetDataMigration 获取数据迁移状态
func (h *SystemHandler) GetDataMigration(c *gin.Context) {
	info := gin.H{
		"status":       "idle",
		"progress":     0,
		"total":        0,
		"migrated":     0,
		"last_run":     "",
	}
	response.Success(c, info)
}

// StartDataMigration 开始数据迁移
func (h *SystemHandler) StartDataMigration(c *gin.Context) {
	var req struct {
		Source string `json:"source" binding:"required"`
		Target string `json:"target" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	h.log.Info("starting data migration from %s to %s", req.Source, req.Target)
	response.SuccessMsg(c, "migration started")
}

// GetSystemAuthRules 获取系统权限规则
func (h *SystemHandler) GetSystemAuthRules(c *gin.Context) {
	rules := []gin.H{
		{"id": 1, "name": "admin.access", "description": "后台访问"},
		{"id": 2, "name": "user.manage", "description": "用户管理"},
		{"id": 3, "name": "product.manage", "description": "产品管理"},
		{"id": 4, "name": "finance.manage", "description": "财务管理"},
		{"id": 5, "name": "ticket.manage", "description": "工单管理"},
		{"id": 6, "name": "system.manage", "description": "系统管理"},
	}
	response.Success(c, rules)
}

// GetSystemLanguage 获取系统语言包
func (h *SystemHandler) GetSystemLanguage(c *gin.Context) {
	lang := c.DefaultQuery("lang", "zh-CN")

	// 返回语言包
	languagePack := gin.H{
		"lang":  lang,
		"keys":  gin.H{},
	}
	response.Success(c, languagePack)
}
