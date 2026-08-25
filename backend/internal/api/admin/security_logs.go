package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetSecurityLogs 获取安全审计日志列表（MD 9.1 功能9）
// GET /api/admin/security-logs
// 支持按类型、时间、IP筛选，支持导出
func GetSecurityLogs(c *gin.Context) {
	db := database.GetDB()

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	attackType := c.Query("attack_type")
	ip := c.Query("ip")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	query := db.Model(&model.SecurityLog{})

	if attackType != "" {
		query = query.Where("attack_type = ?", attackType)
	}
	if ip != "" {
		query = query.Where("ip = ?", ip)
	}
	if startTimeStr != "" {
		if t, err := time.Parse("2006-01-02", startTimeStr); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if endTimeStr != "" {
		if t, err := time.Parse("2006-01-02", endTimeStr); err == nil {
			query = query.Where("created_at <= ?", t.Add(24*time.Hour-time.Second))
		}
	}

	var total int64
	query.Count(&total)

	var logs []model.SecurityLog
	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)

	if logs == nil {
		logs = []model.SecurityLog{}
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

// GetSecurityLogSummary 获取安全日志统计（按攻击类型分组）
// GET /api/admin/security-logs/summary
func GetSecurityLogSummary(c *gin.Context) {
	db := database.GetDB()

	// 最近24小时统计
	oneDayAgo := time.Now().Add(-24 * time.Hour)
	var results []struct {
		AttackType string `json:"attack_type"`
		Count      int64  `json:"count"`
	}
	db.Model(&model.SecurityLog{}).
		Select("attack_type, count(*) as count").
		Where("created_at >= ?", oneDayAgo).
		Group("attack_type").
		Scan(&results)

	// 最近7天统计
	sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)
	var weeklyTotal int64
	db.Model(&model.SecurityLog{}).Where("created_at >= ?", sevenDaysAgo).Count(&weeklyTotal)

	// 最近24小时总数
	var dailyTotal int64
	db.Model(&model.SecurityLog{}).Where("created_at >= ?", oneDayAgo).Count(&dailyTotal)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"daily_total": dailyTotal,
			"weekly_total": weeklyTotal,
			"by_type":    results,
		},
	})
}

// GetSecurityLogDetail 获取安全日志详情
// GET /api/admin/security-logs/:id
func GetSecurityLogDetail(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")

	var log model.SecurityLog
	if err := db.First(&log, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "日志不存在", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    log,
	})
}
