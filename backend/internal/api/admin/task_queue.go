package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetTaskQueueOverview 获取任务队列概览
// GET /api/admin/task-queue/overview
func GetTaskQueueOverview(c *gin.Context) {
	db := database.GetDB()

	var pendingCount int64
	db.Model(&model.TaskQueue{}).Where("status = ?", "pending").Count(&pendingCount)

	var runningCount int64
	db.Model(&model.TaskQueue{}).Where("status = ?", "running").Count(&runningCount)

	var completedCount int64
	db.Model(&model.TaskQueue{}).Where("status = ?", "completed").Count(&completedCount)

	var failedCount int64
	db.Model(&model.TaskQueue{}).Where("status = ?", "failed").Count(&failedCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"pending":   pendingCount,
			"running":   runningCount,
			"completed": completedCount,
			"failed":    failedCount,
			"total":     pendingCount + runningCount + completedCount + failedCount,
		},
	})
}

// GetTaskQueueList 获取任务队列列表
// GET /api/admin/task-queue
func GetTaskQueueList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	taskType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.TaskQueue{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if taskType != "" {
		query = query.Where("type = ?", taskType)
	}

	var total int64
	query.Count(&total)

	var tasks []model.TaskQueue
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&tasks)

	if tasks == nil {
		tasks = []model.TaskQueue{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      tasks,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// RetryTask 重试任务
// POST /api/admin/task-queue/:id/retry
func RetryTask(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var task model.TaskQueue
	if err := db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "任务不存在", "data": nil})
		return
	}

	if task.Status != "failed" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "只有失败的任务才能重试", "data": nil})
		return
	}

	db.Model(&task).Updates(map[string]interface{}{
		"status":      "pending",
		"error":       "",
		"retry_count": 0,
	})

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "重试成功", "data": nil})
}

// DeleteTask 删除任务
// DELETE /api/admin/task-queue/:id
func DeleteTask(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var task model.TaskQueue
	if err := db.First(&task, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "任务不存在", "data": nil})
		return
	}

	db.Delete(&task)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
