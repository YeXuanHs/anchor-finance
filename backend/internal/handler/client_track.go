package handler

import (
	"net/http"
	"strconv"

	"anchorfinance/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ClientTrackHandler 客户跟踪管理
type ClientTrackHandler struct {
	db *gorm.DB
}

func NewClientTrackHandler(db *gorm.DB) *ClientTrackHandler {
	return &ClientTrackHandler{db: db}
}

// List 获取客户跟踪记录列表
// GET /admin/client-tracks
func (h *ClientTrackHandler) List(c *gin.Context) {
	uid := c.Query("uid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query := h.db.Model(&model.ClientTrackRecord{})
	if uid != "" {
		query = query.Where("uid = ?", uid)
	}

	var total int64
	query.Count(&total)

	var records []model.ClientTrackRecord
	if err := query.Preload("Admin").Preload("Remarks.Admin").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"list":      records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetDetail 获取跟踪记录详情
// GET /admin/client-tracks/:id
func (h *ClientTrackHandler) GetDetail(c *gin.Context) {
	id := c.Param("id")

	var record model.ClientTrackRecord
	if err := h.db.Preload("Admin").Preload("Remarks.Admin").First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "记录不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": record})
}

// Create 创建跟踪记录
// POST /admin/client-tracks
func (h *ClientTrackHandler) Create(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var req struct {
		UID uint   `json:"uid" binding:"required"`
		Des string `json:"des" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	record := model.ClientTrackRecord{
		UID:     req.UID,
		AdminID: adminID.(uint),
		Des:     req.Des,
	}

	if err := h.db.Create(&record).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	// 重新加载关联数据
	h.db.Preload("Admin").First(&record, record.ID)

	c.JSON(http.StatusOK, gin.H{"data": record, "message": "创建成功"})
}

// Update 更新跟踪记录
// PUT /admin/client-tracks/:id
func (h *ClientTrackHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Des string `json:"des" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if err := h.db.Model(&model.ClientTrackRecord{}).Where("id = ?", id).Update("des", req.Des).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// Delete 删除跟踪记录
// DELETE /admin/client-tracks/:id
func (h *ClientTrackHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	// 删除关联的备注
	h.db.Where("track_id = ?", id).Delete(&model.ClientTrackRemark{})

	if err := h.db.Delete(&model.ClientTrackRecord{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// AddRemark 添加备注
// POST /admin/client-tracks/:id/remarks
func (h *ClientTrackHandler) AddRemark(c *gin.Context) {
	trackID := c.Param("id")
	adminID, _ := c.Get("admin_id")

	var req struct {
		Remark string `json:"remark" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	tid, _ := strconv.ParseUint(trackID, 10, 32)
	remark := model.ClientTrackRemark{
		TrackID:  uint(tid),
		AdminID:  adminID.(uint),
		Remark:   req.Remark,
	}

	if err := h.db.Create(&remark).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "添加失败"})
		return
	}

	// 重新加载关联数据
	h.db.Preload("Admin").First(&remark, remark.ID)

	c.JSON(http.StatusOK, gin.H{"data": remark, "message": "添加成功"})
}

// DeleteRemark 删除备注
// DELETE /admin/client-tracks/remarks/:id
func (h *ClientTrackHandler) DeleteRemark(c *gin.Context) {
	id := c.Param("id")

	if err := h.db.Delete(&model.ClientTrackRemark{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetUserTracks 获取用户的跟踪记录（前台用户查看）
// GET /user/tracks
func (h *ClientTrackHandler) GetUserTracks(c *gin.Context) {
	uid, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query := h.db.Model(&model.ClientTrackRecord{}).Where("uid = ?", uid)
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	var total int64
	query.Count(&total)

	var records []model.ClientTrackRecord
	if err := query.Preload("Remarks").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"list":      records,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// UpdateTrackStatus 更新客户跟踪状态
// PUT /admin/client-tracks/status
func (h *ClientTrackHandler) UpdateTrackStatus(c *gin.Context) {
	var req struct {
		UID         uint `json:"uid" binding:"required"`
		TrackStatus int  `json:"track_status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 验证状态值 (1:待跟进, 2:跟进中, 3:已完成)
	if req.TrackStatus < 1 || req.TrackStatus > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "状态只能为1,2,3"})
		return
	}

	// 更新用户表中的跟踪状态
	if err := h.db.Table("users").Where("id = ?", req.UID).Update("track_status", req.TrackStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// GetClientNotes 获取客户备注
// GET /admin/client-tracks/notes/:uid
func (h *ClientTrackHandler) GetClientNotes(c *gin.Context) {
	uid := c.Param("uid")

	var user struct {
		Notes string `json:"notes"`
	}
	if err := h.db.Table("users").Select("notes").Where("id = ?", uid).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "客户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"notes": user.Notes}})
}

// UpdateClientNotes 更新客户备注
// PUT /admin/client-tracks/notes/:uid
func (h *ClientTrackHandler) UpdateClientNotes(c *gin.Context) {
	uid := c.Param("uid")

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	if len(req.Notes) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "备注不超过500个字符"})
		return
	}

	if err := h.db.Table("users").Where("id = ?", uid).Update("notes", req.Notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
