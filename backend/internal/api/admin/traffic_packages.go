package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetTrafficPackageList 获取流量包列表
// GET /api/admin/traffic-packages
func GetTrafficPackageList(c *gin.Context) {
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
	db.Model(&model.TrafficPackage{}).Count(&total)

	var packages []model.TrafficPackage
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("sort_order ASC").Find(&packages)

	if packages == nil {
		packages = []model.TrafficPackage{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      packages,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateTrafficPackage 创建流量包（防0元购：价格必须>0）
// POST /api/admin/traffic-packages
func CreateTrafficPackage(c *gin.Context) {
	var req struct {
		Name      string  `json:"name" binding:"required"`
		Volume    int64   `json:"volume" binding:"required"`
		Price     float64 `json:"price" binding:"required"`
		Unit      string  `json:"unit"`
		SortOrder int     `json:"sort_order"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	// 0元购防护：价格必须大于0，流量必须大于0
	if req.Price < 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "价格不能为负数", "data": nil})
		return
	}
	if req.Volume <= 0 {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "流量大小必须大于0", "data": nil})
		return
	}

	if req.Unit == "" {
		req.Unit = "GB"
	}

	db := database.GetDB()
	pkg := model.TrafficPackage{
		Name:      req.Name,
		Volume:    req.Volume,
		Price:     req.Price,
		Unit:      req.Unit,
		SortOrder: req.SortOrder,
		Status:    "active",
	}

	if err := db.Create(&pkg).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": pkg.ID,
		},
	})
}

// UpdateTrafficPackage 更新流量包
// PUT /api/admin/traffic-packages/:id
func UpdateTrafficPackage(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name      string  `json:"name"`
		Volume    int64   `json:"volume"`
		Price     float64 `json:"price"`
		Unit      string  `json:"unit"`
		SortOrder int     `json:"sort_order"`
		Status    string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	db := database.GetDB()
	var pkg model.TrafficPackage
	if err := db.First(&pkg, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "流量包不存在", "data": nil})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Volume > 0 {
		updates["volume"] = req.Volume
	}
	if req.Price >= 0 {
		updates["price"] = req.Price
	}
	if req.Unit != "" {
		updates["unit"] = req.Unit
	}
	if req.SortOrder > 0 {
		updates["sort_order"] = req.SortOrder
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	db.Model(&pkg).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "更新成功", "data": nil})
}

// DeleteTrafficPackage 删除流量包
// DELETE /api/admin/traffic-packages/:id
func DeleteTrafficPackage(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var pkg model.TrafficPackage
	if err := db.First(&pkg, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "流量包不存在", "data": nil})
		return
	}

	db.Delete(&pkg)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}

// GetTrafficLogList 获取流量使用记录
// GET /api/admin/traffic-logs
func GetTrafficLogList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	userID := c.Query("user_id")
	serviceID := c.Query("service_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.TrafficLog{})

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if serviceID != "" {
		query = query.Where("service_id = ?", serviceID)
	}

	var total int64
	query.Count(&total)

	var logs []model.TrafficLog
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)

	if logs == nil {
		logs = []model.TrafficLog{}
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
