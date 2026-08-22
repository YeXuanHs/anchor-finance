package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetReferralOverview 获取推介概览
// GET /api/admin/referral/overview
func GetReferralOverview(c *gin.Context) {
	db := database.GetDB()

	var totalReferrals int64
	db.Model(&model.Referral{}).Count(&totalReferrals)

	var completedReferrals int64
	db.Model(&model.Referral{}).Where("status = ?", "completed").Count(&completedReferrals)

	var totalReward float64
	db.Model(&model.Referral{}).Where("status = ?", "completed").Select("COALESCE(SUM(reward), 0)").Scan(&totalReward)

	var pendingWithdrawals int64
	db.Model(&model.ReferralWithdrawal{}).Where("status = ?", "pending").Count(&pendingWithdrawals)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"total_referrals":     totalReferrals,
			"completed_referrals": completedReferrals,
			"total_reward":        totalReward,
			"pending_withdrawals": pendingWithdrawals,
		},
	})
}

// GetReferralRewards 获取推介奖励列表
// GET /api/admin/referral/rewards
func GetReferralRewards(c *gin.Context) {
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
	db.Model(&model.Referral{}).Count(&total)

	var referrals []model.Referral
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&referrals)

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

// GetReferralWithdrawals 获取推介提现列表
// GET /api/admin/referral-withdrawals
func GetReferralWithdrawals(c *gin.Context) {
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
	query := db.Model(&model.ReferralWithdrawal{})

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var withdrawals []model.ReferralWithdrawal
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&withdrawals)

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

// ApproveReferralWithdrawal 批准提现
// POST /api/admin/referral-withdrawals/:id/approve
func ApproveReferralWithdrawal(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的提现ID",,
			"data": nil})
		return
	}

	db := database.GetDB()
	var withdrawal model.ReferralWithdrawal
	if err := db.First(&withdrawal, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "提现记录不存在",,
			"data": nil})
		return
	}

	if withdrawal.Status != "pending" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "只有待审核的提现才能批准",,
			"data": nil})
		return
	}

	db.Model(&withdrawal).Update("status", "approved")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "批准成功",
	})
}

// RejectReferralWithdrawal 拒绝提现
// POST /api/admin/referral-withdrawals/:id/reject
func RejectReferralWithdrawal(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的提现ID",,
			"data": nil})
		return
	}

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)

	db := database.GetDB()
	var withdrawal model.ReferralWithdrawal
	if err := db.First(&withdrawal, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "提现记录不存在",,
			"data": nil})
		return
	}

	if withdrawal.Status != "pending" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "只有待审核的提现才能拒绝",,
			"data": nil})
		return
	}

	updates := map[string]interface{}{
		"status": "rejected",
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	db.Model(&withdrawal).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "拒绝成功",
	})
}
