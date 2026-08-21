package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetVerificationList 获取实名认证列表
// GET /api/admin/verifications
func GetVerificationList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	verType := c.Query("type")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Verification{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if verType != "" {
		query = query.Where("type = ?", verType)
	}

	var total int64
	query.Count(&total)

	var verifications []model.Verification
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&verifications)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      verifications,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetVerificationSummary 获取实名认证统计
// GET /api/admin/verifications/summary
func GetVerificationSummary(c *gin.Context) {
	db := database.GetDB()

	var pendingCount int64
	db.Model(&model.Verification{}).Where("status = ?", "pending").Count(&pendingCount)

	var approvedCount int64
	db.Model(&model.Verification{}).Where("status = ?", "approved").Count(&approvedCount)

	var rejectedCount int64
	db.Model(&model.Verification{}).Where("status = ?", "rejected").Count(&rejectedCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"pending":  pendingCount,
			"approved": approvedCount,
			"rejected": rejectedCount,
			"total":    pendingCount + approvedCount + rejectedCount,
		},
	})
}

// GetVerificationDetail 获取实名认证详情
// GET /api/admin/verifications/:id
func GetVerificationDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的认证ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var verification model.Verification
	if err := db.First(&verification, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "认证记录不存在",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    verification,
	})
}

// ApproveVerification 批准实名认证
// POST /api/admin/verifications/:id/approve
func ApproveVerification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的认证ID",
			"data":    nil,
		})
		return
	}

	db := database.GetDB()
	var verification model.Verification
	if err := db.First(&verification, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "认证记录不存在",
			"data":    nil,
		})
		return
	}

	if verification.Status != "pending" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "只有待审核的认证才能批准",
			"data":    nil,
		})
		return
	}

	// 更新认证状态
	db.Model(&verification).Update("status", "approved")

	// 更新用户实名状态
	db.Model(&model.User{}).Where("id = ?", verification.UserID).Update("is_verified", true)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "批准成功",
		"data":    nil,
	})
}

// RejectVerification 拒绝实名认证
// POST /api/admin/verifications/:id/reject
func RejectVerification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的认证ID",
			"data":    nil,
		})
		return
	}

	var req struct {
		Remark string `json:"remark"`
	}
	c.ShouldBindJSON(&req)

	db := database.GetDB()
	var verification model.Verification
	if err := db.First(&verification, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "认证记录不存在",
			"data":    nil,
		})
		return
	}

	if verification.Status != "pending" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "只有待审核的认证才能拒绝",
			"data":    nil,
		})
		return
	}

	// 更新认证状态
	updates := map[string]interface{}{
		"status": "rejected",
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}
	db.Model(&verification).Updates(updates)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "拒绝成功",
		"data":    nil,
	})
}
