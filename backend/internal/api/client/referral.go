package client

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetUserReferralOverview 获取用户推介概览
// GET /api/client/referral/overview
func GetUserReferralOverview(c *gin.Context) {
	userID, _ := c.Get("user_id")

	db := database.GetDB()

	// 获取用户的推介码（简单实现：用户ID作为推介码）
	var user model.User
	db.First(&user, userID)

	// 统计推介人数
	var referralCount int64
	db.Model(&model.Referral{}).Where("referrer_id = ?", userID).Count(&referralCount)

	// 统计总奖励
	var totalReward float64
	db.Model(&model.Referral{}).Where("referrer_id = ? AND status = ?", userID, "completed").
		Select("COALESCE(SUM(reward), 0)").Scan(&totalReward)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"referral_code":   user.ID,
			"referral_count":  referralCount,
			"total_reward":    totalReward,
		},
	})
}

// GetUserReferralRewards 获取用户推介奖励
// GET /api/client/referral/rewards
func GetUserReferralRewards(c *gin.Context) {
	userID, _ := c.Get("user_id")
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
	db.Model(&model.Referral{}).Where("referrer_id = ?", userID).Count(&total)

	var referrals []model.Referral
	offset := (page - 1) * pageSize
	db.Where("referrer_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&referrals)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      referrals,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ApplyReferralWithdrawal 申请推介提现
// POST /api/client/referral/withdrawals
func ApplyReferralWithdrawal(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Amount float64 `json:"amount" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "提现金额必须大于0",
		})
		return
	}

	// 检查是否有足够的奖励
	db := database.GetDB()
	var totalReward float64
	db.Model(&model.Referral{}).Where("referrer_id = ? AND status = ?", userID, "completed").
		Select("COALESCE(SUM(reward), 0)").Scan(&totalReward)

	// 检查已提现金额
	var totalWithdrawn float64
	db.Model(&model.ReferralWithdrawal{}).Where("user_id = ? AND status IN ?", userID, []string{"pending", "approved", "paid"}).
		Select("COALESCE(SUM(amount), 0)").Scan(&totalWithdrawn)

	available := totalReward - totalWithdrawn
	if req.Amount > available {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "可提现金额不足",
		})
		return
	}

	withdrawal := model.ReferralWithdrawal{
		UserID: userID.(uint),
		Amount: req.Amount,
		Status: "pending",
	}

	if err := db.Create(&withdrawal).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "申请提现失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "申请成功，等待审核",
	})
}

// GetUserReferralWithdrawals 获取用户提现记录
// GET /api/client/referral/withdrawals
func GetUserReferralWithdrawals(c *gin.Context) {
	userID, _ := c.Get("user_id")
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
	db.Model(&model.ReferralWithdrawal{}).Where("user_id = ?", userID).Count(&total)

	var withdrawals []model.ReferralWithdrawal
	offset := (page - 1) * pageSize
	db.Where("user_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&withdrawals)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      withdrawals,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}
