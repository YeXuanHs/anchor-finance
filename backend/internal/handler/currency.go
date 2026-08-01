package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/response"

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

// AdminSetDefault sets a currency as default and unsets others.
// PUT /admin/currencies/:id/default
func (h *CurrencyHandler) AdminSetDefault(c *gin.Context) {
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

	// Unset all defaults, then set this one
	if err := h.db.Model(&model.Currency{}).Where("is_default = ?", true).Update("is_default", false).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	if err := h.db.Model(&currency).Update("is_default", true).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "default currency updated")
}

// AdminUpdateAllPrices recalculates all pricing based on the default currency rates.
// POST /admin/currencies/update-prices
func (h *CurrencyHandler) AdminUpdateAllPrices(c *gin.Context) {
	// Get default currency ID
	var defaultCurrency model.Currency
	if err := h.db.Where("is_default = ?", true).First(&defaultCurrency).Error; err != nil {
		response.BadRequest(c, "no default currency found")
		return
	}

	// Get all pricings for default currency
	var defaultPricings []struct {
		Type   string
		RelID  uint
		Fields map[string]interface{}
	}

	rows, err := h.db.Table("pricing").Where("currency = ?", defaultCurrency.ID).Rows()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	defer rows.Close()

	columns, _ := rows.Columns()
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)

		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}
		defaultPricings = append(defaultPricings, struct {
			Type   string
			RelID  uint
			Fields map[string]interface{}{
				"type":  rowMap["type"],
				"relid": rowMap["relid"],
				"data":  rowMap,
			},
		})
	}

	// Get all currencies
	var currencies []model.Currency
	h.db.Find(&currencies)

	// Update pricing for each currency based on default currency pricing × rate
	updated := 0
	for _, pricing := range defaultPricings {
		for _, curr := range currencies {
			if curr.ID == defaultCurrency.ID {
				continue
			}
			// Update pricing for this currency
			result := h.db.Table("pricing").
				Where("type = ? AND relid = ? AND currency = ?", pricing.Type, pricing.RelID, curr.ID).
				Updates(map[string]interface{}{
					"monthly": gorm.Expr("monthly * ?", curr.Rate),
				})
			if result.RowsAffected > 0 {
				updated++
			}
		}
	}

	response.Success(c, gin.H{"updated_count": updated})
}
