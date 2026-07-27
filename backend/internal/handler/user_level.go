package handler

import (
	"strconv"

	"github.com/anchor-finance/backend/internal/model"
	"github.com/anchor-finance/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserLevelHandler handles user level HTTP requests.
type UserLevelHandler struct {
	db *gorm.DB
}

// NewUserLevelHandler creates a new UserLevelHandler.
func NewUserLevelHandler(db *gorm.DB) *UserLevelHandler {
	return &UserLevelHandler{db: db}
}

// GetAll returns all user levels.
// GET /user-levels
func (h *UserLevelHandler) GetAll(c *gin.Context) {
	var levels []model.UserLevel
	if err := h.db.Order("priority DESC, id ASC").Find(&levels).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, levels)
}

// AdminGetList returns a paginated list of user levels (admin).
// GET /admin/user-levels
func (h *UserLevelHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var levels []model.UserLevel
	var total int64

	query := h.db.Model(&model.UserLevel{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("priority DESC, id ASC").Find(&levels)

	response.SuccessPage(c, levels, total, page, pageSize)
}

// AdminCreate creates a new user level (admin).
// POST /admin/user-levels
func (h *UserLevelHandler) AdminCreate(c *gin.Context) {
	var level model.UserLevel
	if err := c.ShouldBindJSON(&level); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.db.Create(&level).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, level)
}

// AdminUpdate updates a user level (admin).
// PUT /admin/user-levels/:id
func (h *UserLevelHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid level id")
		return
	}

	var level model.UserLevel
	if err := h.db.First(&level, id).Error; err != nil {
		response.NotFound(c, "level not found")
		return
	}

	var req model.UserLevel
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"min_amount":  req.MinAmount,
		"discount":    req.Discount,
		"priority":    req.Priority,
		"icon":        req.Icon,
		"description": req.Description,
		"is_default":  req.IsDefault,
	}

	if err := h.db.Model(&level).Updates(updates).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	h.db.First(&level, id)
	response.Success(c, level)
}

// AdminDelete deletes a user level (admin).
// DELETE /admin/user-levels/:id
func (h *UserLevelHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid level id")
		return
	}

	if err := h.db.Delete(&model.UserLevel{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "level deleted")
}
