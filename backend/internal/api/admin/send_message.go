package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetSendMessageSearchParams 获取发送消息搜索参数
// GET /api/admin/send-message/search-params
func GetSendMessageSearchParams(c *gin.Context) {
	// 返回搜索参数选项
	params := gin.H{
		"user_status": []gin.H{
			{"value": "active", "label": "正常"},
			{"value": "suspended", "label": "暂停"},
			{"value": "closed", "label": "关闭"},
		},
		"message_type": []gin.H{
			{"value": "email", "label": "邮件"},
			{"value": "sms", "label": "短信"},
			{"value": "system", "label": "站内信"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    params,
	})
}

// GetSendMethodList 获取发送方式列表
// GET /api/admin/send-message/send-methods
func GetSendMethodList(c *gin.Context) {
	methods := []gin.H{
		{"id": "email", "name": "邮件", "description": "通过邮件发送"},
		{"id": "sms", "name": "短信", "description": "通过短信发送"},
		{"id": "system", "name": "站内信", "description": "通过站内信发送"},
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    methods,
	})
}

// SearchSendMessageList 搜索发送消息列表
// GET /api/admin/send-message/search
func SearchSendMessageList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	messageType := c.Query("type")
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 根据类型查询不同的表
	db := database.GetDB()
	var total int64
	var results []interface{}

	switch messageType {
	case "email":
		// 查询邮件日志
		var logs []model.SystemLog
		query := db.Model(&model.SystemLog{}).Where("type = ?", "email")
		if keyword != "" {
			query = query.Where("content LIKE ?", "%"+keyword+"%")
		}
		query.Count(&total)
		offset := (page - 1) * pageSize
		query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)
		for _, l := range logs {
			results = append(results, l)
		}

	case "sms":
		// 查询短信日志
		var logs []model.SystemLog
		query := db.Model(&model.SystemLog{}).Where("type = ?", "sms")
		if keyword != "" {
			query = query.Where("content LIKE ?", "%"+keyword+"%")
		}
		query.Count(&total)
		offset := (page - 1) * pageSize
		query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)
		for _, l := range logs {
			results = append(results, l)
		}

	default:
		// 查询所有消息
		var logs []model.SystemLog
		query := db.Model(&model.SystemLog{}).Where("type IN ?", []string{"email", "sms", "system"})
		if keyword != "" {
			query = query.Where("content LIKE ?", "%"+keyword+"%")
		}
		query.Count(&total)
		offset := (page - 1) * pageSize
		query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs)
		for _, l := range logs {
			results = append(results, l)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      results,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
