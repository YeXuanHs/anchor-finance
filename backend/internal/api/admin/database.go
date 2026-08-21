package admin

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/gin-gonic/gin"
)

// GetDatabaseStatus 获取数据库状态
// GET /api/admin/database/status
func GetDatabaseStatus(c *gin.Context) {
	db := database.GetDB()

	// 获取数据库大小
	var dbSize float64
	db.Raw("SELECT ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&dbSize)

	// 获取表数量
	var tableCount int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&tableCount)

	// 获取各表行数
	type TableInfo struct {
		Name string `json:"name"`
		Rows int64  `json:"rows"`
	}
	var tables []TableInfo
	db.Raw("SELECT table_name as name, table_rows as rows FROM information_schema.tables WHERE table_schema = DATABASE() ORDER BY table_rows DESC LIMIT 20").Scan(&tables)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"size_mb":     dbSize,
			"table_count": tableCount,
			"tables":      tables,
		},
	})
}

// OptimizeDatabase 优化数据库
// POST /api/admin/database/optimizations
func OptimizeDatabase(c *gin.Context) {
	db := database.GetDB()

	// 获取所有表
	var tables []string
	db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()").Pluck("table_name", &tables)

	// 优化每个表
	optimized := 0
	for _, table := range tables {
		db.Exec("OPTIMIZE TABLE " + table)
		optimized++
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "优化完成",
		"data": gin.H{
			"optimized_tables": optimized,
		},
	})
}
