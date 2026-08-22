package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetScheduleOverview 获取定时任务概览
// GET /api/admin/schedules/overview
func GetScheduleOverview(c *gin.Context) {
	db := database.GetDB()

	var activeCount int64
	db.Model(&model.ScheduleTask{}).Where("status = ?", "active").Count(&activeCount)

	var disabledCount int64
	db.Model(&model.ScheduleTask{}).Where("status = ?", "disabled").Count(&disabledCount)

	var totalCount int64
	db.Model(&model.ScheduleTask{}).Count(&totalCount)

	// 统计成功/失败运行
	var successCount int64
	db.Model(&model.ScheduleRun{}).Where("status = ?", "success").Count(&successCount)

	var failedCount int64
	db.Model(&model.ScheduleRun{}).Where("status = ?", "failed").Count(&failedCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"active":        activeCount,
			"disabled":      disabledCount,
			"total":         totalCount,
			"success_runs":  successCount,
			"failed_runs":   failedCount,
		},
	})
}

// GetScheduleRunList 获取定时任务运行记录
// GET /api/admin/schedule-runs
func GetScheduleRunList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	taskID := c.Query("task_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.ScheduleRun{})

	if taskID != "" {
		query = query.Where("task_id = ?", taskID)
	}

	var total int64
	query.Count(&total)

	var runs []model.ScheduleRun
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&runs)

	if runs == nil {
		runs = []model.ScheduleRun{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      runs,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetScheduleRunDetail 获取定时任务运行详情
// GET /api/admin/schedule-runs/:id
func GetScheduleRunDetail(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var run model.ScheduleRun
	if err := db.First(&run, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "运行记录不存在", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": run})
}
