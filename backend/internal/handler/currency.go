package handler

import (
	"strconv"

	"github.com/anchor-finance/backend/internal/model"
	"github.com/anchor-finance/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CurrencyHandler handles currency HTTP requests.
type CurrencyHandler struct {
	db *gorm.DB
}

// NewCurrencyHandler creates a new CurrencyHandler.
func NewCurrencyHandler(db *gorm.DB) *CurrencyHandler {
	return &CurrencyHandler{db: db}
}

// GetAll returns all active currencies.
// GET /currencies
func (h *CurrencyHandler) GetAll(c *gin.Context) {
	var currencies []model.Currency
	if err := h.db.Where("is_active = ?", true).Order("is_default DESC, id ASC").Find(&currencies).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, currencies)
}

// AdminGetList returns a paginated list of currencies (admin).
// GET /admin/currencies
func (h *CurrencyHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var currencies []model.Currency
	var total int64

	query := h.db.Model(&model.Currency{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("is_default DESC, id ASC").Find(&currencies)

	response.SuccessPage(c, currencies, total, page, pageSize)
}

// AdminCreate creates a new currency (admin).
// POST /admin/currencies
func (h *CurrencyHandler) AdminCreate(c *gin.Context) {
	var currency model.Currency
	if err := c.ShouldBindJSON(&currency); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.db.Create(&currency).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, currency)
}

// AdminUpdate updates a currency (admin).
// PUT /admin/currencies/:id
func (h *CurrencyHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid currency id")
		return
	}

	var currency model.Currency
	if err := h.db.First(&currency, id).Error; err != nil {
		response.NotFound(c, "currency not found")
		return
	}

	var req model.Currency
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{
		"code":       req.Code,
		"name":       req.Name,
		"symbol":     req.Symbol,
		"rate":       req.Rate,
		"is_default": req.IsDefault,
		"is_active":  req.IsActive,
		"precision":  req.Precision,
	}

	if err := h.db.Model(&currency).Updates(updates).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	h.db.First(&currency, id)
	response.Success(c, currency)
}

// AdminDelete deletes a currency (admin).
// DELETE /admin/currencies/:id
func (h *CurrencyHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid currency id")
		return
	}

	if err := h.db.Delete(&model.Currency{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "currency deleted")
}

// AdminUpdateRate updates only the exchange rate of a currency (admin).
// PUT /admin/currencies/:id/rate
func (h *CurrencyHandler) AdminUpdateRate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid currency id")
		return
	}

	var req struct {
		Rate float64 `json:"rate" binding:"required,gt=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var currency model.Currency
	if err := h.db.First(&currency, id).Error; err != nil {
		response.NotFound(c, "currency not found")
		return
	}

	if err := h.db.Model(&currency).Update("rate", req.Rate).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	h.db.First(&currency, id)
	response.Success(c, currency)
}
