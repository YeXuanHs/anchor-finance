package handler

import (
	"strconv"

	"anchorfinance/internal/backup"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// BackupHandler 备份管理处理器
type BackupHandler struct {
	backupSvc *backup.Service
	log       *logger.Logger
}

// NewBackupHandler 创建备份处理器
func NewBackupHandler(backupSvc *backup.Service, log *logger.Logger) *BackupHandler {
	return &BackupHandler{
		backupSvc: backupSvc,
		log:       log,
	}
}

// CreateBackup 创建数据库备份
func (h *BackupHandler) CreateBackup(c *gin.Context) {
	var req struct {
		Host     string `json:"host" binding:"required"`
		Port     int    `json:"port" binding:"required"`
		User     string `json:"user" binding:"required"`
		Password string `json:"password" binding:"required"`
		Database string `json:"database" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.backupSvc.Backup(backup.BackupConfig{
		Host:     req.Host,
		Port:     req.Port,
		User:     req.User,
		Password: req.Password,
		Database: req.Database,
	})
	if err != nil {
		h.log.Errorf("数据库备份失败: %v", err)
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, result)
}

// ListBackups 列出所有备份
func (h *BackupHandler) ListBackups(c *gin.Context) {
	backups, err := h.backupSvc.ListBackups()
	if err != nil {
		h.log.Errorf("列出备份失败: %v", err)
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"backups": backups})
}

// DeleteBackup 删除备份
func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		response.BadRequest(c, "文件名不能为空")
		return
	}

	if err := h.backupSvc.DeleteBackup(filename); err != nil {
		h.log.Errorf("删除备份失败: %v", err)
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "备份已删除")
}

// CleanOldBackups 清理旧备份
func (h *BackupHandler) CleanOldBackups(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days < 7 {
		days = 7
	}

	deleted, err := h.backupSvc.CleanOldBackups(days)
	if err != nil {
		h.log.Errorf("清理旧备份失败: %v", err)
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"deleted": deleted,
		"days":    days,
	})
}

// RestoreBackup 恢复备份
func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	var req struct {
		Host     string `json:"host" binding:"required"`
		Port     int    `json:"port" binding:"required"`
		User     string `json:"user" binding:"required"`
		Password string `json:"password" binding:"required"`
		Database string `json:"database" binding:"required"`
		Filename string `json:"filename" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.backupSvc.Restore(backup.BackupConfig{
		Host:     req.Host,
		Port:     req.Port,
		User:     req.User,
		Password: req.Password,
		Database: req.Database,
	}, req.Filename); err != nil {
		h.log.Errorf("恢复备份失败: %v", err)
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "数据库恢复成功")
}
