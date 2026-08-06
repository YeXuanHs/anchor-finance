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

// GetCronTrend 定时任务执行趋势图表数据
// GET /admin/run-map/cron-trend?time_type=1
func (h *RunMapHandler) GetCronTrend(c *gin.Context) {
	timeType, _ := strconv.Atoi(c.DefaultQuery("time_type", "1"))

	var days int
	switch timeType {
	case 2:
		days = 15
	case 3:
		days = 30
	case 4:
		days = 90
	default:
		days = 7
	}

	type TrendItem struct {
		Date  string `json:"date"`
		Total int    `json:"total"`
		Fail  int    `json:"fail"`
	}

	var trends []TrendItem
	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC).Unix()
		endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC).Unix()

		var item TrendItem
		item.Date = date.Format("01-02")

		// 从 system_logs 或 run_map 表统计执行情况
		var total int64
		h.db.Table("run_maps").
			Where("created_at >= ? AND created_at <= ?", startOfDay, endOfDay).
			Count(&total)
		item.Total = int(total)

		var fail int64
		h.db.Table("run_maps").
			Where("created_at >= ? AND created_at <= ? AND run_count = 0", startOfDay, endOfDay).
			Count(&fail)
		item.Fail = int(fail)

		trends = append(trends, item)
	}

	response.Success(c, gin.H{
		"trends": trends,
	})
}

// GetCronHistory 定时任务执行历史列表
// GET /admin/run-map/cron-history?datetime=20260802
func (h *RunMapHandler) GetCronHistory(c *gin.Context) {
	datetime := c.DefaultQuery("datetime", time.Now().Format("20060102"))

	// 解析日期
	date, err := time.Parse("20060102", datetime)
	if err != nil {
		response.BadRequest(c, "invalid date format")
		return
	}
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC)

	// 定时任务类型定义
	type CronType struct {
		Type      int    `json:"type"`
		Name      string `json:"name"`
		ShowFail  int    `json:"show_fail"`
		Keywords  string `json:"keywords"`
	}
	cronTypes := []CronType{
		{Type: 1, Name: "自动暂停", ShowFail: 1, Keywords: "产品到期暂停"},
		{Type: 2, Name: "未实名暂停", ShowFail: 1, Keywords: "未实名客户产品暂停"},
		{Type: 3, Name: "信用额产品暂停", ShowFail: 1, Keywords: "信用额账单未支付"},
		{Type: 4, Name: "自动删除", ShowFail: 1, Keywords: "产品到期删除"},
		{Type: 7, Name: "DCIM流量重置", ShowFail: 1, Keywords: "DCIM流量"},
		{Type: 8, Name: "魔方云流量重置", ShowFail: 1, Keywords: "魔方云流量"},
		{Type: 9, Name: "关闭工单", ShowFail: 0, Keywords: ""},
		{Type: 10, Name: "信用额账单提醒", ShowFail: 0, Keywords: ""},
		{Type: 12, Name: "账单提醒", ShowFail: 0, Keywords: ""},
		{Type: 13, Name: "生成账单", ShowFail: 0, Keywords: ""},
	}

	type HistoryItem struct {
		Type     int    `json:"type"`
		Name     string `json:"name"`
		ShowFail int    `json:"show_fail"`
		Keywords string `json:"keywords"`
		All      int    `json:"all"`
		Fail     int    `json:"fail"`
	}

	var result []HistoryItem
	for _, ct := range cronTypes {
		item := HistoryItem{
			Type:     ct.Type,
			Name:     ct.Name,
			ShowFail: ct.ShowFail,
			Keywords: ct.Keywords,
		}

		// 统计全部执行次数
		if ct.Keywords != "" {
			var allCount int64
			h.db.Table("run_maps").
				Where("description LIKE ? AND created_at >= ? AND created_at <= ?",
					"%"+ct.Keywords+"%", startOfDay.Unix(), endOfDay.Unix()).
				Count(&allCount)
			item.All = int(allCount)
		}

		// 统计失败次数
		if ct.ShowFail == 1 && ct.Keywords != "" {
			var failCount int64
			h.db.Table("run_maps").
				Where("description LIKE ? AND created_at >= ? AND created_at <= ? AND run_count = 0",
					"%"+ct.Keywords+"%", startOfDay.Unix(), endOfDay.Unix()).
				Count(&failCount)
			item.Fail = int(failCount)
		}

		result = append(result, item)
	}

	response.Success(c, gin.H{
		"data": result,
	})
}
