package admin

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
	"github.com/gin-gonic/gin"
)

// GetMarketingPushList 获取营销推送列表
// GET /api/admin/marketing/pushes
func GetMarketingPushList(c *gin.Context) {
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
	query := db.Model(&model.MarketingPush{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var pushes []model.MarketingPush
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&pushes)

	if pushes == nil {
		pushes = []model.MarketingPush{}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      pushes,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// CreateMarketingPush 创建营销推送
// POST /api/admin/marketing/pushes
func CreateMarketingPush(c *gin.Context) {
	var req struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		Type       string `json:"type" binding:"required"`
		TargetType string `json:"target_type"`
		TargetIDs  []uint `json:"target_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.TargetType == "" {
		req.TargetType = "all"
	}

	targetIDsJSON := "[]"
	if len(req.TargetIDs) > 0 {
		jsonBytes, _ := json.Marshal(req.TargetIDs)
		targetIDsJSON = string(jsonBytes)
	}

	db := database.GetDB()
	push := model.MarketingPush{
		Title:      req.Title,
		Content:    req.Content,
		Type:       req.Type,
		TargetType: req.TargetType,
		TargetIDs:  targetIDsJSON,
		Status:     "draft",
	}

	if err := db.Create(&push).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "操作失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id": push.ID,
		},
	})
}

// SendMarketingPush 发送营销推送（走PHP插件引擎）
// POST /api/admin/marketing/pushes/:id/send
func SendMarketingPush(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var push model.MarketingPush
	if err := db.First(&push, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "推送不存在", "data": nil})
		return
	}

	if push.Status != "draft" {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "只有草稿状态的推送才能发送", "data": nil})
		return
	}

	// 查询目标用户
	var users []model.User
	if push.TargetType == "all" {
		db.Where("status = ?", "active").Find(&users)
	} else if push.TargetType == "user" {
		var targetIDs []uint
		json.Unmarshal([]byte(push.TargetIDs), &targetIDs)
		if len(targetIDs) > 0 {
			db.Where("id IN ? AND status = ?", targetIDs, "active").Find(&users)
		}
	}

	// 通过PHP插件引擎逐个发送
	var sentCount, failedCount int
	now := time.Now()

	for _, user := range users {
		var sendErr error
		switch push.Type {
		case "sms":
			if user.Phone != "" {
				sendErr = pluginengine.SendSMS(user.Phone, push.Content)
			}
		case "email":
			if user.Email != "" {
				sendErr = pluginengine.SendEmail(user.Email, push.Title, push.Content)
			}
		}

		if sendErr != nil {
			failedCount++
		} else {
			sentCount++
		}
	}

	db.Model(&push).Updates(map[string]interface{}{
		"status":       "sent",
		"sent_at":      &now,
		"sent_count":   sentCount,
		"failed_count": failedCount,
	})

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "发送完成",
		"data": gin.H{
			"sent_count":   sentCount,
			"failed_count": failedCount,
		},
	})
}

// DeleteMarketingPush 删除营销推送
// DELETE /api/admin/marketing/pushes/:id
func DeleteMarketingPush(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var push model.MarketingPush
	if err := db.First(&push, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "推送不存在", "data": nil})
		return
	}

	db.Delete(&push)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}
