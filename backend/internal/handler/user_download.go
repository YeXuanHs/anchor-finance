package handler

import (
	"net/http"
	"strconv"

	"anchorfinance/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserDownloadHandler 用户专属下载管理
type UserDownloadHandler struct {
	db *gorm.DB
}

func NewUserDownloadHandler(db *gorm.DB) *UserDownloadHandler {
	return &UserDownloadHandler{db: db}
}

// AdminList 管理员获取用户下载列表
// GET /admin/user-downloads
func (h *UserDownloadHandler) AdminList(c *gin.Context) {
	uid := c.Query("uid")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	query := h.db.Model(&model.UserDownload{})
	if uid != "" {
		query = query.Where("uid = ?", uid)
	}

	var total int64
	query.Count(&total)

	var downloads []model.UserDownload
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&downloads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"list":      downloads,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// AdminCreate 管理员为用户上传文件
// POST /admin/user-downloads
func (h *UserDownloadHandler) AdminCreate(c *gin.Context) {
	adminID, _ := c.Get("admin_id")

	var req struct {
		UID     uint   `json:"uid" binding:"required"`
		Name    string `json:"name" binding:"required"`
		URL     string `json:"url" binding:"required"`
		DowName string `json:"downame"`
		Remarks string `json:"remarks"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	download := model.UserDownload{
		UID:     req.UID,
		Name:    req.Name,
		URL:     req.URL,
		DowName: req.DowName,
		Remarks: req.Remarks,
		AdminID: adminID.(uint),
	}
	if err := h.db.Create(&download).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": download, "message": "创建成功"})
}

// AdminDelete 管理员删除用户下载
// DELETE /admin/user-downloads/:id
func (h *UserDownloadHandler) AdminDelete(c *gin.Context) {
	id := c.Param("id")
	if err := h.db.Delete(&model.UserDownload{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetUserDownloads 用户获取自己的下载列表
// GET /user/downloads
func (h *UserDownloadHandler) GetUserDownloads(c *gin.Context) {
	uid, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	offset := (page - 1) * pageSize

	var total int64
	h.db.Model(&model.UserDownload{}).Where("uid = ?", uid).Count(&total)

	var downloads []model.UserDownload
	if err := h.db.Where("uid = ?", uid).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&downloads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"list":      downloads,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// ========== UserTaste Handler ==========

// UserTasteHandler 用户偏好设置
type UserTasteHandler struct {
	db *gorm.DB
}

func NewUserTasteHandler(db *gorm.DB) *UserTasteHandler {
	return &UserTasteHandler{db: db}
}

// Get 获取用户偏好
// GET /user/tastes
func (h *UserTasteHandler) Get(c *gin.Context) {
	uid, _ := c.Get("user_id")

	var taste model.UserTaste
	if err := h.db.Where("uid = ?", uid).First(&taste).Error; err != nil {
		// 返回默认值
		taste = model.UserTaste{
			UID:           uid.(uint),
			TicketRefresh: "never",
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": taste})
}

// Update 更新用户偏好
// PUT /user/tastes
func (h *UserTasteHandler) Update(c *gin.Context) {
	uid, _ := c.Get("user_id")

	var req struct {
		TicketRefresh string `json:"ticket_refresh"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	taste := model.UserTaste{
		UID:           uid.(uint),
		TicketRefresh: req.TicketRefresh,
	}

	// 使用 Upsert
	if err := h.db.Where("uid = ?", uid).Assign(map[string]interface{}{
		"ticket_refresh": req.TicketRefresh,
	}).FirstOrCreate(&taste).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": taste, "message": "更新成功"})
}
