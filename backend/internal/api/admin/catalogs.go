package admin

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetCPUModelCatalog 获取CPU型号目录
// GET /api/admin/cpu-model-catalog
func GetCPUModelCatalog(c *gin.Context) {
	db := database.GetDB()
	var catalogs []model.CPUModel
	db.Where("status = ?", "active").Order("id ASC").Find(&catalogs)

	if catalogs == nil {
		catalogs = []model.CPUModel{}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": catalogs})
}

// GetInstanceSpecCatalog 获取实例规格目录
// GET /api/admin/instance-spec-catalog
func GetInstanceSpecCatalog(c *gin.Context) {
	db := database.GetDB()
	var specs []model.InstanceSpec
	db.Where("status = ?", "active").Order("id ASC").Find(&specs)

	if specs == nil {
		specs = []model.InstanceSpec{}
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": specs})
}
