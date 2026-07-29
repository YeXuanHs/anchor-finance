package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetUserHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewGetUserHandler(db *gorm.DB, log *logger.Logger) *GetUserHandler {
	return &GetUserHandler{db: db, log: log}
}

// GetUser returns a user by ID with sales restrictions.
func (h *GetUserHandler) GetUser(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var user model.User
	if err := h.db.Preload("Group").First(&user, uid).Error; err != nil {
		response.NotFound(c, "user not found")
		return
	}
	response.Success(c, user)
}

// GetUsers returns a list of users with sales restrictions.
func (h *GetUserHandler) GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var users []model.User
	var total int64

	query := h.db.Model(&model.User{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ? OR phone LIKE ?", like, like, like, like)
	}
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Preload("Group").Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, users, total, page, pageSize)
}

// CheckAccess checks if the current admin has access to the user.
func (h *GetUserHandler) CheckAccess(c *gin.Context) {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var user model.User
	if err := h.db.First(&user, uid).Error; err != nil {
		response.NotFound(c, "user not found")
		return
	}

	// Admin users always have access
	if user.IsAdmin {
		response.Success(c, gin.H{
			"user_id": uid,
			"access":  true,
		})
		return
	}

	// Check if user status is active
	access := user.Status == 1
	response.Success(c, gin.H{
		"user_id": uid,
		"access":  access,
		"status":  user.Status,
	})
}
