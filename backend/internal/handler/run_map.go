package handler

import (
	"strconv"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RunMapHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewRunMapHandler(db *gorm.DB, log *logger.Logger) *RunMapHandler {
	return &RunMapHandler{db: db, log: log}
}

// GetRunMapList returns a list of run map tasks.
func (h *RunMapHandler) GetRunMapList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var tasks []model.RunMap
	var total int64

	query := h.db.Model(&model.RunMap{})
	if keyword := c.Query("keyword"); keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR code LIKE ?", like, like)
	}
	if taskType := c.Query("type"); taskType != "" {
		query = query.Where("type = ?", taskType)
	}
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&tasks).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, tasks, total, page, pageSize)
}

// GetRunMap returns a single run map task.
func (h *RunMapHandler) GetRunMap(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	var task model.RunMap
	if err := h.db.First(&task, id).Error; err != nil {
		response.NotFound(c, "task not found")
		return
	}
	response.Success(c, task)
}

// RepeatTask repeats a failed task.
func (h *RunMapHandler) RepeatTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid task id")
		return
	}

	var task model.RunMap
	if err := h.db.First(&task, id).Error; err != nil {
		response.NotFound(c, "task not found")
		return
	}

	newTask := model.RunMap{
		Name:      task.Name + " (copy)",
		Code:      task.Code + "_copy_" + strconv.FormatInt(time.Now().UnixMilli(), 10),
		Type:      task.Type,
		Config:    task.Config,
		IsEnabled: true,
	}
	if err := h.db.Create(&newTask).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":      newTask.ID,
		"message": "task repeated",
	})
}

// GetTaskTypes returns available task types.
func (h *RunMapHandler) GetTaskTypes(c *gin.Context) {
	types := []gin.H{
		{"value": "script", "label": "脚本执行"},
		{"value": "api", "label": "API调用"},
		{"value": "webhook", "label": "Webhook"},
		{"value": "auto_provision", "label": "自动开通"},
		{"value": "auto_renew", "label": "自动续费"},
		{"value": "auto_suspend", "label": "自动暂停"},
		{"value": "auto_terminate", "label": "自动终止"},
		{"value": "domain_sync", "label": "域名同步"},
		{"value": "certificate_sync", "label": "证书同步"},
		{"value": "backup", "label": "备份任务"},
	}
	response.Success(c, gin.H{
		"types": types,
	})
}
