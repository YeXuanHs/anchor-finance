package admin

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/gin-gonic/gin"
)

// tableNameRegex 表名白名单校验（防SQL注入）
var tableNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

// GetDatabaseStatus 获取数据库状态
// GET /api/admin/database/status
func GetDatabaseStatus(c *gin.Context) {
	db := database.GetDB()

	var dbSize float64
	db.Raw("SELECT ROUND(SUM(data_length + index_length) / 1024 / 1024, 2) FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&dbSize)

	var tableCount int64
	db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE()").Scan(&tableCount)

	type TableInfo struct {
		Name string `json:"name"`
		Rows int64  `json:"rows"`
	}
	var tables []TableInfo
	db.Raw("SELECT table_name as name, table_rows as rows FROM information_schema.tables WHERE table_schema = DATABASE() ORDER BY table_rows DESC LIMIT 20").Scan(&tables)

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "success",
		"data": gin.H{
			"size_mb":     dbSize,
			"table_count": tableCount,
			"tables":      tables,
		},
	})
}

// BackupDatabase 备份数据库（H2/H3修复：使用环境变量文件传密码，写入受限目录）
// POST /api/admin/database/backups
func BackupDatabase(c *gin.Context) {
	db := database.GetDB()

	var tables []string
	db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()").Pluck("table_name", &tables)

	// H3修复：写入受限目录而非/tmp
	backupDir := "/var/backups/anchor"
	os.MkdirAll(backupDir, 0700)

	backupFile := fmt.Sprintf("%s/backup_%s.sql", backupDir, time.Now().Format("20060102_150405"))

	// H2修复：使用MYSQL_PWD环境变量传递密码（不在进程参数中暴露）
	cmd := exec.Command("mysqldump",
		"-h", os.Getenv("DB_HOST"),
		"-P", os.Getenv("DB_PORT"),
		"-u", os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
	)
	cmd.Env = append(os.Environ(), "MYSQL_PWD="+os.Getenv("DB_PASSWORD"))

	output, err := cmd.Output()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "备份失败", "data": nil})
		return
	}

	// 写入文件（权限600，仅owner可读）
	if err := os.WriteFile(backupFile, output, 0600); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "写入备份文件失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "备份成功",
		"data": gin.H{"file": backupFile, "tables": len(tables)},
	})
}

// OptimizeDatabase 优化数据库（H1修复：表名校验防SQL注入）
// POST /api/admin/database/optimizations
func OptimizeDatabase(c *gin.Context) {
	db := database.GetDB()

	var tables []string
	db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE()").Pluck("table_name", &tables)

	optimized := 0
	for _, table := range tables {
		// H1修复：校验表名只包含合法字符
		if !tableNameRegex.MatchString(table) {
			continue
		}
		db.Exec("OPTIMIZE TABLE `" + table + "`")
		optimized++
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0, "message": "优化完成",
		"data": gin.H{"optimized_tables": optimized},
	})
}
