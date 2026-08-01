package handler

import (
	"fmt"
	"strconv"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

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

// ==================== 新增缺失方法 ====================

// Setting returns contract module settings (admin).
// GET /admin/contracts/setting
func (h *ContractHandler) Setting(c *gin.Context) {
	configKeys := []string{
		"contract_open", "contract_institutions", "contract_address",
		"contract_username", "contract_email", "contract_phonenumber",
		"contract_consignee_address", "contract_postcode", "contract_postcode_fee",
		"contract_number_custom", "contract_number", "contract_number_prefix",
		"contract_pdf_logo", "contract_company_logo",
	}

	configs := make(map[string]interface{})
	for _, key := range configKeys {
		var config model.SystemConfig
		if err := h.db.Where("`key` = ?", key).First(&config).Error; err == nil {
			configs[key] = config.Value
		}
	}

	// Get default currency
	var currency struct {
		Code   string `json:"code"`
		Prefix string `json:"prefix"`
		Suffix string `json:"suffix"`
	}
	h.db.Table("currencies").Where("`default` = 1").Select("code, prefix, suffix").Scan(&currency)
	configs["currency"] = currency

	response.Success(c, configs)
}

// SettingPost updates contract module settings (admin).
// POST /admin/contracts/setting
func (h *ContractHandler) SettingPost(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Validate contract_number length
	if num, ok := req["contract_number"]; ok {
		numStr := fmt.Sprintf("%v", num)
		if len(numStr) < 8 || len(numStr) > 25 {
			response.BadRequest(c, "contract number length must be 8-25")
			return
		}
	}

	configKeys := map[string]bool{
		"contract_open": true, "contract_institutions": true, "contract_address": true,
		"contract_username": true, "contract_email": true, "contract_phonenumber": true,
		"contract_consignee_address": true, "contract_postcode": true, "contract_postcode_fee": true,
		"contract_number_custom": true, "contract_number": true, "contract_number_prefix": true,
		"contract_pdf_logo": true, "contract_company_logo": true,
	}

	for key, value := range req {
		if !configKeys[key] {
			continue
		}
		// Upsert config
		var count int64
		h.db.Model(&model.SystemConfig{}).Where("`key` = ?", key).Count(&count)
		if count > 0 {
			h.db.Model(&model.SystemConfig{}).Where("`key` = ?", key).Update("value", fmt.Sprintf("%v", value))
		} else {
			h.db.Create(&model.SystemConfig{Key: key, Value: fmt.Sprintf("%v", value)})
		}
	}

	response.SuccessMsg(c, "contract settings updated")
}

// Tpl returns contract template list (admin).
// GET /admin/contracts/tpl
func (h *ContractHandler) Tpl(c *gin.Context) {
	var templates []model.ContractTemplate
	h.db.Order("id DESC").Find(&templates)

	// Contract argument placeholders
	args := []map[string]interface{}{
		{"title": "客户信息", "data": []map[string]string{
			{"name": "客户用户名", "arg": "{$client_username}"},
			{"name": "客户公司名", "arg": "{$client_company_name}"},
			{"name": "客户电话", "arg": "{$client_telephone}"},
			{"name": "客户邮箱", "arg": "{$client_email}"},
			{"name": "客户地址", "arg": "{$client_address}"},
		}},
		{"title": "产品信息", "data": []map[string]string{
			{"name": "产品名称", "arg": "{$client_product_name}"},
			{"name": "产品域名", "arg": "{$client_product_domain}"},
			{"name": "产品状态", "arg": "{$client_product_status}"},
			{"name": "开始日期", "arg": "{$client_product_startdate}"},
			{"name": "结束日期", "arg": "{$client_product_enddate}"},
			{"name": "计费周期", "arg": "{$client_product_billingcycle}"},
			{"name": "产品价格", "arg": "{$client_product_price}"},
		}},
	}

	response.Success(c, gin.H{
		"tpl": templates,
		"arg": args,
	})
}

// DeleteTpl deletes a contract template (admin).
// DELETE /admin/contracts/tpl/:id
func (h *ContractHandler) DeleteTpl(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid template id")
		return
	}

	if err := h.db.Delete(&model.ContractTemplate{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "template deleted")
}

// ContractPage returns contract list page for admin with filters.
// GET /admin/contracts/page
func (h *ContractHandler) ContractPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	order := c.DefaultQuery("order", "id")
	sort := c.DefaultQuery("sort", "DESC")

	query := h.db.Model(&model.Contract{})

	if uid := c.Query("uid"); uid != "" {
		query = query.Where("user_id = ?", uid)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if keyword := c.Query("keyword"); keyword != "" {
		query = query.Where("title LIKE ? OR contract_no LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if startTime := c.Query("start_time"); startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime := c.Query("end_time"); endTime != "" {
		query = query.Where("created_at <= ?", endTime)
	}

	var total int64
	query.Count(&total)

	var contracts []model.Contract
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order(order + " " + sort).Find(&contracts)

	// Status names
	statusNames := map[int8]string{
		1: "未签署",
		2: "待签署",
		3: "已签署",
		4: "已作废",
	}

	type EnrichedContract struct {
		model.Contract
		StatusName string `json:"status_name"`
		UserName   string `json:"user_name"`
	}

	enriched := make([]EnrichedContract, len(contracts))
	for i, contract := range contracts {
		enriched[i] = EnrichedContract{
			Contract:   contract,
			StatusName: statusNames[contract.Status],
		}
		var user model.User
		if err := h.db.Select("username").First(&user, contract.UserID).Error; err == nil {
			enriched[i].UserName = user.Username
		}
	}

	response.SuccessPage(c, enriched, total, page, pageSize)
}

// ContractPagePost creates or updates a contract via POST (admin).
// POST /admin/contracts/page
func (h *ContractHandler) ContractPagePost(c *gin.Context) {
	var req struct {
		ID         *uint   `json:"id"`
		ContractNo string  `json:"contract_no" binding:"required"`
		UserID     uint    `json:"user_id" binding:"required"`
		Title      string  `json:"title" binding:"required"`
		Content    string  `json:"content"`
		Type       string  `json:"type"`
		Amount     float64 `json:"amount"`
		StartDate  *string `json:"start_date"`
		EndDate    *string `json:"end_date"`
		Status     *int8   `json:"status"`
		TplID      *uint   `json:"tpl_id"`
		NoteType   *int    `json:"note_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// If template specified, use template content
	if req.TplID != nil && *req.TplID > 0 {
		var tpl model.ContractTemplate
		if err := h.db.First(&tpl, *req.TplID).Error; err == nil {
			req.Content = tpl.Content
		}
	}

	if req.ID != nil && *req.ID > 0 {
		// Update
		updates := map[string]interface{}{
			"contract_no": req.ContractNo,
			"user_id":     req.UserID,
			"title":       req.Title,
			"content":     req.Content,
			"type":        req.Type,
			"amount":      req.Amount,
		}
		if req.Status != nil {
			updates["status"] = *req.Status
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

		if err := h.db.Model(&model.Contract{}).Where("id = ?", *req.ID).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}

		response.SuccessMsg(c, "contract updated")
	} else {
		// Create
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
}
