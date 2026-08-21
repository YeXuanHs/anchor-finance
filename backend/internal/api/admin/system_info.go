package admin

import (
	"net/http"
	"runtime"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/gin-gonic/gin"
)

// GetSystemInfo 获取系统信息
// GET /api/admin/system/info
func GetSystemInfo(c *gin.Context) {
	db := database.GetDB()

	// 获取数据库大小
	var dbSize float64
	db.Raw("SELECT ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&dbSize)

	// 获取表数量
	var tableCount int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&tableCount)

	// 获取Go运行时信息
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"database": gin.H{
				"size_mb":     dbSize,
				"table_count": tableCount,
			},
			"server": gin.H{
				"go_version":    runtime.Version(),
				"goroutine_count": runtime.NumGoroutine(),
				"memory_alloc_mb": float64(memStats.Alloc) / 1024 / 1024,
				"memory_sys_mb":   float64(memStats.Sys) / 1024 / 1024,
				"uptime":          time.Now().Format("2006-01-02 15:04:05"),
			},
		},
	})
}

// GetSystemModules 获取系统模块列表
// GET /api/admin/system/modules
func GetSystemModules(c *gin.Context) {
	modules := []gin.H{
		{"name": "认证模块", "version": "1.0.0", "status": "active"},
		{"name": "用户管理", "version": "1.0.0", "status": "active"},
		{"name": "订单管理", "version": "1.0.0", "status": "active"},
		{"name": "服务管理", "version": "1.0.0", "status": "active"},
		{"name": "账单管理", "version": "1.0.0", "status": "active"},
		{"name": "工单管理", "version": "1.0.0", "status": "active"},
		{"name": "产品管理", "version": "1.0.0", "status": "active"},
		{"name": "插件管理", "version": "1.0.0", "status": "active"},
		{"name": "设置管理", "version": "1.0.0", "status": "active"},
		{"name": "日志管理", "version": "1.0.0", "status": "active"},
		{"name": "内容管理", "version": "1.0.0", "status": "active"},
		{"name": "财务管理", "version": "1.0.0", "status": "active"},
		{"name": "供应商管理", "version": "1.0.0", "status": "active"},
		{"name": "优惠券管理", "version": "1.0.0", "status": "active"},
		{"name": "推介系统", "version": "1.0.0", "status": "active"},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    modules,
	})
}
