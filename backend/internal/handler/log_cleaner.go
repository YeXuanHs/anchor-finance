package handler

import (
	"strconv"
	"time"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LogCleanerRule 日志清理规则模型
type LogCleanerRule struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:100;not null"`
	Type      string         `json:"type" gorm:"size:50;not null"`
	Value     string         `json:"value" gorm:"size:255"`
	Enabled   bool           `json:"enabled" gorm:"default:true"`
	LastRun   *time.Time     `json:"last_run"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

func (LogCleanerRule) TableName() string {
	return "log_cleaner_rules"
}

// LogCleanerHandler 日志清理处理器
type LogCleanerHandler struct {
	cleaner *service.LogCleaner
	db      *gorm.DB
	log     *logger.Logger
}

// NewLogCleanerHandler 创建日志清理处理器
func NewLogCleanerHandler(cleaner *service.LogCleaner, log *logger.Logger) *LogCleanerHandler {
	return &LogCleanerHandler{cleaner: cleaner, db: cleaner.GetDB(), log: log}
}

// GetStats 获取日志统计
func (h *LogCleanerHandler) GetStats(c *gin.Context) {
	stats := h.cleaner.GetLogStats()
	response.Success(c, stats)
}

// CleanByDays 按天数清理
func (h *LogCleanerHandler) CleanByDays(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days < 1 {
		days = 30
	}

	count, err := h.cleaner.CleanByDays(days)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned": count,
		"days":    days,
	})
}

// CleanByCount 按数量保留
func (h *LogCleanerHandler) CleanByCount(c *gin.Context) {
	keepCount, _ := strconv.Atoi(c.DefaultQuery("keep_count", "10000"))
	if keepCount < 100 {
		keepCount = 100
	}

	count, err := h.cleaner.CleanByCount(keepCount)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned":    count,
		"keep_count": keepCount,
	})
}

// CleanByModule 按模块清理
func (h *LogCleanerHandler) CleanByModule(c *gin.Context) {
	module := c.Query("module")
	if module == "" {
		response.BadRequest(c, "模块名不能为空")
		return
	}

	count, err := h.cleaner.CleanByModule(module)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned": count,
		"module":  module,
	})
}

// CleanByStatus 按状态清理
func (h *LogCleanerHandler) CleanByStatus(c *gin.Context) {
	status := c.Query("status")
	if status == "" {
		response.BadRequest(c, "状态不能为空")
		return
	}

	count, err := h.cleaner.CleanByStatus(status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned": count,
		"status":  status,
	})
}

// CleanExpired 清理过期日志
func (h *LogCleanerHandler) CleanExpired(c *gin.Context) {
	count, err := h.cleaner.CleanExpired()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned": count,
	})
}

// CleanAll 清理所有日志（危险操作）
func (h *LogCleanerHandler) CleanAll(c *gin.Context) {
	confirm := c.Query("confirm")
	if confirm != "yes" {
		response.BadRequest(c, "请确认清理操作，添加参数 confirm=yes")
		return
	}

	count, err := h.cleaner.CleanByDays(0)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned": count,
		"warning": "已清理所有日志",
	})
}

// GetRules 获取清理规则列表
func (h *LogCleanerHandler) GetRules(c *gin.Context) {
	var rules []LogCleanerRule
	h.cleaner.GetDB().Find(&rules)
	response.Success(c, rules)
}

// CreateRule 创建清理规则
func (h *LogCleanerHandler) CreateRule(c *gin.Context) {
	var rule LogCleanerRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	h.cleaner.GetDB().Create(&rule)
	response.Success(c, rule)
}

// UpdateRule 更新清理规则
func (h *LogCleanerHandler) UpdateRule(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var rule LogCleanerRule
	if err := h.cleaner.GetDB().First(&rule, id).Error; err != nil {
		response.NotFound(c, "规则不存在")
		return
	}
	if err := c.ShouldBindJSON(&rule); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	h.cleaner.GetDB().Save(&rule)
	response.Success(c, rule)
}

// DeleteRule 删除清理规则
func (h *LogCleanerHandler) DeleteRule(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.cleaner.GetDB().Delete(&LogCleanerRule{}, id)
	response.Success(c, nil)
}

// RunRule 执行清理规则
func (h *LogCleanerHandler) RunRule(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var rule LogCleanerRule
	if err := h.cleaner.GetDB().First(&rule, id).Error; err != nil {
		response.NotFound(c, "规则不存在")
		return
	}

	var count int64
	var err error

	switch rule.Type {
	case "days":
		days, _ := strconv.Atoi(rule.Value)
		count, err = h.cleaner.CleanByDays(days)
	case "count":
		keepCount, _ := strconv.Atoi(rule.Value)
		count, err = h.cleaner.CleanByCount(keepCount)
	case "module":
		count, err = h.cleaner.CleanByModule(rule.Value)
	case "status":
		count, err = h.cleaner.CleanByStatus(rule.Value)
	}

	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	now := time.Now()
	h.cleaner.GetDB().Model(&rule).Update("last_run", now)

	response.Success(c, gin.H{
		"cleaned": count,
		"rule_id": id,
	})
}

// ManualClean 手动清理
func (h *LogCleanerHandler) ManualClean(c *gin.Context) {
	var req struct {
		Days      int    `json:"days"`
		KeepCount int    `json:"keep_count"`
		Module    string `json:"module"`
		Status    string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var count int64
	var err error

	if req.Days > 0 {
		count, err = h.cleaner.CleanByDays(req.Days)
	} else if req.KeepCount > 0 {
		count, err = h.cleaner.CleanByCount(req.KeepCount)
	} else if req.Module != "" {
		count, err = h.cleaner.CleanByModule(req.Module)
	} else if req.Status != "" {
		count, err = h.cleaner.CleanByStatus(req.Status)
	} else {
		response.BadRequest(c, "请指定清理参数")
		return
	}

	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"cleaned": count,
	})
}
