package admin

import (
	"encoding/json"
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetOSOptions 获取操作系统选项列表（从settings表读取，后台可配置）
// GET /api/admin/os-options
func GetOSOptions(c *gin.Context) {
	db := database.GetDB()

	// 从settings表读取OS选项配置
	var setting model.Setting
	if err := db.Where("`key` = ?", "os_options").First(&setting).Error; err == nil && setting.Value != "" {
		var options []gin.H
		if json.Unmarshal([]byte(setting.Value), &options) == nil {
			c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": options})
			return
		}
	}

	// 默认OS选项（首次使用时自动写入settings表）
	options := []gin.H{
		{"id": "centos7", "name": "CentOS 7", "type": "linux"},
		{"id": "centos8", "name": "CentOS 8", "type": "linux"},
		{"id": "ubuntu20", "name": "Ubuntu 20.04", "type": "linux"},
		{"id": "ubuntu22", "name": "Ubuntu 22.04", "type": "linux"},
		{"id": "debian11", "name": "Debian 11", "type": "linux"},
		{"id": "debian12", "name": "Debian 12", "type": "linux"},
		{"id": "windows2019", "name": "Windows Server 2019", "type": "windows"},
		{"id": "windows2022", "name": "Windows Server 2022", "type": "windows"},
	}

	// 自动写入settings表供后台修改
	jsonBytes, _ := json.Marshal(options)
	db.Where("`key` = ?", "os_options").Assign(model.Setting{Key: "os_options", Value: string(jsonBytes), Group: "dcim"}).FirstOrCreate(&model.Setting{})

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": options})
}
