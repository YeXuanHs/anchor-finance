package admin

import (
	"net/http"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetLogCleanupOverview 获取日志清理概览
// GET /api/admin/log-cleanups/overview
func GetLogCleanupOverview(c *gin.Context) {
	db := database.GetDB()

	// 统计各类型日志数量
	var systemLogCount int64
	db.Model(&model.SystemLog{}).Count(&systemLogCount)

	var operationLogCount int64
	db.Model(&model.OperationLog{}).Count(&operationLogCount)

	var loginLogCount int64
	db.Model(&model.LoginLog{}).Count(&loginLogCount)

	// 统计30天前的日志数量
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	var oldSystemLogCount int64
	db.Model(&model.SystemLog{}).Where("created_at < ?", thirtyDaysAgo).Count(&oldSystemLogCount)

	var oldOperationLogCount int64
	db.Model(&model.OperationLog{}).Where("created_at < ?", thirtyDaysAgo).Count(&oldOperationLogCount)

	var oldLoginLogCount int64
	db.Model(&model.LoginLog{}).Where("created_at < ?", thirtyDaysAgo).Count(&oldLoginLogCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"system_logs": gin.H{
				"total": systemLogCount,
				"old":   oldSystemLogCount,
			},
			"operation_logs": gin.H{
				"total": operationLogCount,
				"old":   oldOperationLogCount,
			},
			"login_logs": gin.H{
				"total": loginLogCount,
				"old":   oldLoginLogCount,
			},
		},
	})
}

// CleanupLogs 清理日志
// POST /api/admin/log-cleanups
func CleanupLogs(c *gin.Context) {
	var req struct {
		Type     string `json:"type" binding:"required"` // system, operation, login, all
		KeepDays int    `json:"keep_days"`               // 保留天数，默认30
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	if req.KeepDays <= 0 {
		req.KeepDays = 30
	}

	db := database.GetDB()
	cutoff := time.Now().AddDate(0, 0, -req.KeepDays)

	var deletedCount int64

	switch req.Type {
	case "system":
		result := db.Where("created_at < ?", cutoff).Delete(&model.SystemLog{})
		deletedCount = result.RowsAffected
	case "operation":
		result := db.Where("created_at < ?", cutoff).Delete(&model.OperationLog{})
		deletedCount = result.RowsAffected
	case "login":
		result := db.Where("created_at < ?", cutoff).Delete(&model.LoginLog{})
		deletedCount = result.RowsAffected
	case "all":
		r1 := db.Where("created_at < ?", cutoff).Delete(&model.SystemLog{})
		r2 := db.Where("created_at < ?", cutoff).Delete(&model.OperationLog{})
		r3 := db.Where("created_at < ?", cutoff).Delete(&model.LoginLog{})
		deletedCount = r1.RowsAffected + r2.RowsAffected + r3.RowsAffected
	default:
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的日志类型",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "清理成功",
		"data": gin.H{
			"deleted_count": deletedCount,
		},
	})
}
