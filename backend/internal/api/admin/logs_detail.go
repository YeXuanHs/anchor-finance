package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// getPagedLogs 通用分页查询SystemLog
func getPagedLogs(c *gin.Context, logType string) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.SystemLog{}).Where("type = ?", logType)

	var total int64
	query.Count(&total)

	var logs []model.SystemLog
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)

	if logs == nil {
		logs = []model.SystemLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetSMSLogs 获取短信日志
// GET /api/admin/logs/sms
func GetSMSLogs(c *gin.Context) {
	getPagedLogs(c, "sms")
}

// GetEmailLogs 获取邮件日志
// GET /api/admin/logs/email
func GetEmailLogs(c *gin.Context) {
	getPagedLogs(c, "email")
}

// GetAPILogs 获取API日志
// GET /api/admin/logs/api
func GetAPILogs(c *gin.Context) {
	getPagedLogs(c, "api")
}

// GetCronLogs 获取定时任务日志
// GET /api/admin/logs/cron
func GetCronLogs(c *gin.Context) {
	getPagedLogs(c, "cron")
}

// GetAdminLoginLogs 获取管理员登录日志
// GET /api/admin/logs/admin-login
func GetAdminLoginLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.LoginLog{}).Where("is_admin = ?", true)

	var total int64
	query.Count(&total)

	var logs []model.LoginLog
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)

	if logs == nil {
		logs = []model.LoginLog{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      logs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetNotificationLogs 获取站内信日志
// GET /api/admin/logs/notification
func GetNotificationLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	var total int64
	db.Model(&model.UserNotification{}).Count(&total)

	var notifications []model.UserNotification
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&notifications)

	if notifications == nil {
		notifications = []model.UserNotification{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      notifications,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
