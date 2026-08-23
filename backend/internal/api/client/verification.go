package client

import (
	"net/http"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
	"github.com/gin-gonic/gin"
)

// GetVerificationStatus 获取实名认证状态
// GET /api/client/verification/status
func GetVerificationStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")

	db := database.GetDB()
	var verification model.Verification
	err := db.Where("user_id = ?", userID).Order("id DESC").First(&verification).Error

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "success",
			"data": gin.H{
				"status": "none",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"status":     verification.Status,
			"type":       verification.Type,
			"name":       verification.Name,
			"created_at": verification.CreatedAt,
		},
	})
}

// SubmitVerification 提交实名认证
// POST /api/client/verification/submit
func SubmitVerification(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Type    string `json:"type" binding:"required"` // person, company
		Name    string `json:"name" binding:"required"`
		IDCard  string `json:"id_card"`
		Company string `json:"company"`
		License string `json:"license"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	// 检查是否已有待审核的认证
	db := database.GetDB()
	var count int64
	db.Model(&model.Verification{}).Where("user_id = ? AND status = ?", userID, "pending").Count(&count)
	if count > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "已有待审核的认证，请等待审核",
			"data": nil,
		})
		return
	}

	// 触发Hook: before_certifi_submit（返回false可阻止提交）
	results, hookErr := pluginengine.TriggerHook("before_certifi_submit", map[string]interface{}{
		"user_id": userID, "type": req.Type, "name": req.Name,
	})
	if hookErr == nil && len(results) > 0 {
		if data, ok := results[0].Result.(map[string]interface{}); ok {
			if allowed, ok := data["allowed"].(bool); ok && !allowed {
				c.JSON(http.StatusOK, gin.H{"code": 403, "message": "认证提交被阻止", "data": nil})
				return
			}
		}
	}

	verification := model.Verification{
		UserID:  userID.(uint),
		Type:    req.Type,
		Name:    req.Name,
		IDCard:  req.IDCard,
		Company: req.Company,
		License: req.License,
		Status:  "pending",
	}

	if err := db.Create(&verification).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "提交认证失败",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "提交成功，等待审核",
	})
}
