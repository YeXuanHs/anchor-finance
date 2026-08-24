package admin

import (
	"net/http"
	"runtime"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
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

// GetSystemModules 获取系统模块列表（检测真实状态）
// GET /api/admin/system/modules
func GetSystemModules(c *gin.Context) {
	db := database.GetDB()

	// 检测各模块真实状态
	checkModule := func(name string, model interface{}) gin.H {
		status := "active"
		count := int64(0)
		if err := db.Model(model).Count(&count).Error; err != nil {
			status = "error"
		}
		return gin.H{"name": name, "version": "1.0.0", "status": status, "count": count}
	}

	modules := []gin.H{
		checkModule("用户管理", &model.User{}),
		checkModule("订单管理", &model.Order{}),
		checkModule("服务管理", &model.Service{}),
		checkModule("账单管理", &model.Invoice{}),
		checkModule("工单管理", &model.Ticket{}),
		checkModule("产品管理", &model.Product{}),
		checkModule("插件管理", &model.Plugin{}),
		checkModule("供应商管理", &model.Supplier{}),
		checkModule("内容管理", &model.News{}),
		checkModule("财务管理", &model.Payment{}),
		checkModule("日志管理", &model.OperationLog{}),
		{"name": "认证模块", "version": "1.0.0", "status": "active", "count": 0},
		{"name": "设置管理", "version": "1.0.0", "status": "active", "count": 0},
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": modules})
}
