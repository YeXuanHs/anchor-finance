package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
	"github.com/gin-gonic/gin"
)

// GetCancelRequestList 获取取消请求列表
// GET /api/admin/cancel-requests
func GetCancelRequestList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.CancelRequest{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var requests []model.CancelRequest
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&requests)

	if requests == nil {
		requests = []model.CancelRequest{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      requests,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ApproveCancelRequest 批准取消请求
// POST /api/admin/cancel-requests/:id/approve
func ApproveCancelRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的请求ID", "data": nil})
		return
	}

	db := database.GetDB()
	var request model.CancelRequest
	if err := db.First(&request, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "请求不存在", "data": nil})
		return
	}

	if request.Status != "pending" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "只有待处理的请求才能批准", "data": nil})
		return
	}

	// 更新请求状态
	db.Model(&request).Update("status", "approved")

	// 暂停相关服务（走PHP插件引擎）
	if request.ServiceID > 0 {
		pluginengine.TriggerHook("suspend_service", map[string]interface{}{
			"service_id": request.ServiceID,
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "批准成功", "data": nil})
}

// RejectCancelRequest 拒绝取消请求
// POST /api/admin/cancel-requests/:id/reject
func RejectCancelRequest(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "无效的请求ID", "data": nil})
		return
	}

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)

	db := database.GetDB()
	var request model.CancelRequest
	if err := db.First(&request, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "请求不存在", "data": nil})
		return
	}

	if request.Status != "pending" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "只有待处理的请求才能拒绝", "data": nil})
		return
	}

	updates := map[string]interface{}{
		"status": "rejected",
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	db.Model(&request).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "拒绝成功", "data": nil})
}
