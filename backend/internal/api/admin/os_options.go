package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOSOptions 获取操作系统选项列表
// GET /api/admin/os-options
func GetOSOptions(c *gin.Context) {
	// 返回常用的操作系统选项
	options := []gin.H{
		{"id": "centos7", "name": "CentOS 7", "type": "linux"},
		{"id": "centos8", "name": "CentOS 8", "type": "linux"},
		{"id": "centos9", "name": "CentOS 9", "type": "linux"},
		{"id": "ubuntu20", "name": "Ubuntu 20.04", "type": "linux"},
		{"id": "ubuntu22", "name": "Ubuntu 22.04", "type": "linux"},
		{"id": "ubuntu24", "name": "Ubuntu 24.04", "type": "linux"},
		{"id": "debian11", "name": "Debian 11", "type": "linux"},
		{"id": "debian12", "name": "Debian 12", "type": "linux"},
		{"id": "almalinux8", "name": "AlmaLinux 8", "type": "linux"},
		{"id": "almalinux9", "name": "AlmaLinux 9", "type": "linux"},
		{"id": "rocky8", "name": "Rocky Linux 8", "type": "linux"},
		{"id": "rocky9", "name": "Rocky Linux 9", "type": "linux"},
		{"id": "windows2019", "name": "Windows Server 2019", "type": "windows"},
		{"id": "windows2022", "name": "Windows Server 2022", "type": "windows"},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    options,
	})
}
