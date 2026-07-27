package handler

import (
	"strconv"
	"time"

	"github.com/anchor-finance/backend/internal/model"
	"github.com/anchor-finance/backend/pkg/logger"
	"github.com/anchor-finance/backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ContractHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewContractHandler(db *gorm.DB, log *logger.Logger) *ContractHandler {
	return &ContractHandler{db: db, log: log}
}

// ---------- User Endpoints ----------

// GetUserContracts returns paginated contracts for the authenticated user.
func (h *ContractHandler) GetUserContracts(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var contracts []model.Contract
	var total int64

	query := h.db.Model(&model.Contract{}).Where("user_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&contracts).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, contracts, total, page, pageSize)
}

// GetDetail returns a single contract for the authenticated user.
func (h *ContractHandler) GetDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contract id")
		return
	}

	var contract model.Contract
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&contract).Error; err != nil {
		response.NotFound(c, "contract not found")
		return
	}
	response.Success(c, contract)
}

// SignContract marks a contract as signed by the user.
func (h *ContractHandler) SignContract(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contract id")
		return
	}

	var contract model.Contract
	if err := h.db.Where("id = ? AND user_id = ?", id, userID).First(&contract).Error; err != nil {
		response.NotFound(c, "contract not found")
		return
	}

	if contract.Status != 2 {
		response.BadRequest(c, "contract is not in pending sign status")
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":    3,
		"signed_at": now,
	}
	if err := h.db.Model(&contract).Updates(updates).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "contract signed")
}

// ---------- Admin Endpoints ----------

// AdminGetList returns all contracts with pagination (admin).
func (h *ContractHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var status *int
	if s := c.Query("status"); s != "" {
		v, _ := strconv.Atoi(s)
		status = &v
	}
	var userID *uint
	if u := c.Query("user_id"); u != "" {
		v, _ := strconv.ParseUint(u, 10, 64)
		uid := uint(v)
		userID = &uid
	}

	var contracts []model.Contract
	var total int64

	query := h.db.Model(&model.Contract{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&contracts).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, contracts, total, page, pageSize)
}

// AdminGetDetail returns a single contract by ID (admin).
func (h *ContractHandler) AdminGetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contract id")
		return
	}

	var contract model.Contract
	if err := h.db.First(&contract, id).Error; err != nil {
		response.NotFound(c, "contract not found")
		return
	}
	response.Success(c, contract)
}

type CreateContractRequest struct {
	ContractNo string   `json:"contract_no" binding:"required,max=32"`
	UserID     uint     `json:"user_id" binding:"required"`
	Title      string   `json:"title" binding:"required,max=255"`
	Content    string   `json:"content"`
	Type       string   `json:"type" binding:"max=50"`
	Amount     float64  `json:"amount"`
	StartDate  *string  `json:"start_date"`
	EndDate    *string  `json:"end_date"`
}

// AdminCreate creates a new contract (admin).
func (h *ContractHandler) AdminCreate(c *gin.Context) {
	var req CreateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	contract := model.Contract{
		ContractNo: req.ContractNo,
		UserID:     req.UserID,
		Title:      req.Title,
		Content:    req.Content,
		Type:       req.Type,
		Status:     1,
		Amount:     req.Amount,
		AdminID:    c.GetUint("user_id"),
	}

	if req.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *req.StartDate); err == nil {
			contract.StartDate = &t
		}
	}
	if req.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			contract.EndDate = &t
		}
	}

	if err := h.db.Create(&contract).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, contract)
}

type UpdateContractRequest struct {
	Title      *string  `json:"title"`
	Content    *string  `json:"content"`
	Type       *string  `json:"type"`
	Status     *int8    `json:"status"`
	Amount     *float64 `json:"amount"`
	StartDate  *string  `json:"start_date"`
	EndDate    *string  `json:"end_date"`
}

// AdminUpdate updates a contract (admin).
func (h *ContractHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contract id")
		return
	}

	var req UpdateContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}
	if req.StartDate != nil {
		if t, err := time.Parse("2006-01-02", *req.StartDate); err == nil {
			updates["start_date"] = t
		}
	}
	if req.EndDate != nil {
		if t, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			updates["end_date"] = t
		}
	}

	if len(updates) > 0 {
		if err := h.db.Model(&model.Contract{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}
	response.SuccessMsg(c, "contract updated")
}

// AdminDelete soft-deletes a contract (admin).
func (h *ContractHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contract id")
		return
	}

	if err := h.db.Delete(&model.Contract{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "contract deleted")
}

// AdminSign marks a contract as signed by admin.
func (h *ContractHandler) AdminSign(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid contract id")
		return
	}

	var contract model.Contract
	if err := h.db.First(&contract, id).Error; err != nil {
		response.NotFound(c, "contract not found")
		return
	}

	if contract.Status != 1 && contract.Status != 2 {
		response.BadRequest(c, "contract cannot be signed in current status")
		return
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":    3,
		"signed_at": now,
	}
	if err := h.db.Model(&contract).Updates(updates).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "contract signed by admin")
}
