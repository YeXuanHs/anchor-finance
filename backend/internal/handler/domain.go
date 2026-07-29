package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DomainHandler 域名处理器
type DomainHandler struct {
	domainSvc *service.DomainService
	db        *gorm.DB
	log       *logger.Logger
}

// NewDomainHandler 创建域名处理器实例
func NewDomainHandler(domainSvc *service.DomainService, db *gorm.DB, log *logger.Logger) *DomainHandler {
	return &DomainHandler{domainSvc: domainSvc, db: db, log: log}
}

// ──────────────────────────── 用户接口 ────────────────────────────

// GetList 获取用户的域名列表
func (h *DomainHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	params := service.DomainListParams{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Keyword:  keyword,
	}

	domains, total, err := h.domainSvc.GetList(userID, params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, domains, total, page, pageSize)
}

// GetDetail 获取域名详情
func (h *DomainHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid domain id")
		return
	}

	userID := c.GetUint("user_id")
	domain, err := h.domainSvc.GetByIDAndUser(uint(id), userID)
	if err != nil {
		response.NotFound(c, "domain not found")
		return
	}
	response.Success(c, domain)
}

// Renew 续费域名
func (h *DomainHandler) Renew(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid domain id")
		return
	}

	var req struct {
		Years int `json:"years" binding:"required,min=1,max=10"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "years must be between 1 and 10")
		return
	}

	userID := c.GetUint("user_id")
	// 验证域名属于当前用户
	_, err = h.domainSvc.GetByIDAndUser(uint(id), userID)
	if err != nil {
		response.NotFound(c, "domain not found")
		return
	}

	domain, err := h.domainSvc.Renew(uint(id), req.Years)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, domain)
}

// GetDNSRecords 获取域名DNS记录
func (h *DomainHandler) GetDNSRecords(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid domain id")
		return
	}

	userID := c.GetUint("user_id")
	// 验证域名属于当前用户
	_, err = h.domainSvc.GetByIDAndUser(uint(id), userID)
	if err != nil {
		response.NotFound(c, "domain not found")
		return
	}

	records, err := h.domainSvc.GetDNSRecords(uint(id))
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, records)
}

// AddDNSRecord 添加DNS记录
func (h *DomainHandler) AddDNSRecord(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid domain id")
		return
	}

	userID := c.GetUint("user_id")
	// 验证域名属于当前用户
	_, err = h.domainSvc.GetByIDAndUser(uint(id), userID)
	if err != nil {
		response.NotFound(c, "domain not found")
		return
	}

	var req struct {
		Type     string `json:"type" binding:"required,oneof=A AAAA CNAME MX TXT NS SRV"`
		Name     string `json:"name" binding:"required"`
		Value    string `json:"value" binding:"required"`
		TTL      int    `json:"ttl"`
		Priority *int   `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	record := &model.DomainDNSRecord{
		DomainID: uint(id),
		Type:     req.Type,
		Name:     req.Name,
		Value:    req.Value,
		TTL:      req.TTL,
		Priority: req.Priority,
	}
	if record.TTL == 0 {
		record.TTL = 3600
	}

	if err := h.domainSvc.AddDNSRecord(record); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, record)
}

// UpdateDNSRecord 更新DNS记录
func (h *DomainHandler) UpdateDNSRecord(c *gin.Context) {
	domainID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid domain id")
		return
	}

	recordID, err := strconv.ParseUint(c.Param("record_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid record id")
		return
	}

	userID := c.GetUint("user_id")
	// 验证域名属于当前用户
	_, err = h.domainSvc.GetByIDAndUser(uint(domainID), userID)
	if err != nil {
		response.NotFound(c, "domain not found")
		return
	}

	// 验证记录属于该域名
	record, err := h.domainSvc.GetDNSRecordByID(uint(recordID))
	if err != nil {
		response.NotFound(c, "record not found")
		return
	}
	if record.DomainID != uint(domainID) {
		response.BadRequest(c, "record does not belong to this domain")
		return
	}

	var req struct {
		Type     string `json:"type" binding:"required,oneof=A AAAA CNAME MX TXT NS SRV"`
		Name     string `json:"name" binding:"required"`
		Value    string `json:"value" binding:"required"`
		TTL      int    `json:"ttl"`
		Priority *int   `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	record.Type = req.Type
	record.Name = req.Name
	record.Value = req.Value
	record.TTL = req.TTL
	record.Priority = req.Priority

	if err := h.domainSvc.UpdateDNSRecord(record); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, record)
}

// DeleteDNSRecord 删除DNS记录
func (h *DomainHandler) DeleteDNSRecord(c *gin.Context) {
	domainID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid domain id")
		return
	}

	recordID, err := strconv.ParseUint(c.Param("record_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid record id")
		return
	}

	userID := c.GetUint("user_id")
	// 验证域名属于当前用户
	_, err = h.domainSvc.GetByIDAndUser(uint(domainID), userID)
	if err != nil {
		response.NotFound(c, "domain not found")
		return
	}

	// 验证记录属于该域名
	record, err := h.domainSvc.GetDNSRecordByID(uint(recordID))
	if err != nil {
		response.NotFound(c, "record not found")
		return
	}
	if record.DomainID != uint(domainID) {
		response.BadRequest(c, "record does not belong to this domain")
		return
	}

	if err := h.domainSvc.DeleteDNSRecord(uint(recordID)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "deleted")
}

// CheckAvailability 检查域名是否可用
func (h *DomainHandler) CheckAvailability(c *gin.Context) {
	domainName := c.Query("domain")
	if domainName == "" {
		response.BadRequest(c, "domain name is required")
		return
	}

	available, err := h.domainSvc.CheckAvailability(domainName)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, gin.H{"domain": domainName, "available": available})
}

// InitiateTransfer 发起域名转移
func (h *DomainHandler) InitiateTransfer(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		DomainName  string  `json:"domain_name" binding:"required"`
		EPPCode     string  `json:"epp_code" binding:"required"`
		RegistrarID *uint   `json:"registrar_id"`
		Price       float64 `json:"price"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	transfer, err := h.domainSvc.InitiateTransfer(userID, req.DomainName, req.EPPCode, req.Price, req.RegistrarID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, transfer)
}

// GetTransfers 获取用户的域名转移列表
func (h *DomainHandler) GetTransfers(c *gin.Context) {
	userID := c.GetUint("user_id")

	transfers, err := h.domainSvc.GetTransfers(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, transfers)
}

// ──────────────────────────── 管理员接口 ────────────────────────────

// AdminGetList 管理员获取域名列表
func (h *DomainHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	params := service.DomainListParams{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Keyword:  keyword,
	}

	domains, total, err := h.domainSvc.AdminGetList(params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, domains, total, page, pageSize)
}

// AdminGetDetail 管理员获取域名详情
func (h *DomainHandler) AdminGetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid domain id")
		return
	}

	domain, err := h.domainSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "domain not found")
		return
	}
	response.Success(c, domain)
}

// AdminCreate 管理员创建域名
func (h *DomainHandler) AdminCreate(c *gin.Context) {
	var req struct {
		UserID           uint    `json:"user_id" binding:"required"`
		DomainName       string  `json:"domain_name" binding:"required"`
		RegistrarID      *uint   `json:"registrar_id"`
		RegistrationDate string  `json:"registration_date"`
		ExpiryDate       string  `json:"expiry_date"`
		NextDueDate      string  `json:"next_due_date"`
		Nameservers      string  `json:"nameservers"`
		Status           string  `json:"status"`
		AutoRenew        bool    `json:"auto_renew"`
		WhoisPrivacy     bool    `json:"whois_privacy"`
		TransferLock     bool    `json:"transfer_lock"`
		EPPCode          string  `json:"epp_code"`
		DNSManaged       bool    `json:"dns_managed"`
		Price            float64 `json:"price"`
		Currency         string  `json:"currency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	domain := &model.Domain{
		UserID:       req.UserID,
		DomainName:   req.DomainName,
		RegistrarID:  req.RegistrarID,
		Nameservers:  req.Nameservers,
		Status:       req.Status,
		AutoRenew:    req.AutoRenew,
		WhoisPrivacy: req.WhoisPrivacy,
		TransferLock: req.TransferLock,
		EPPCode:      req.EPPCode,
		DNSManaged:   req.DNSManaged,
		Price:        req.Price,
		Currency:     req.Currency,
	}

	if domain.Status == "" {
		domain.Status = "active"
	}
	if domain.Currency == "" {
		domain.Currency = "CNY"
	}

	if err := h.domainSvc.Create(domain); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, domain)
}

// AdminUpdate 管理员更新域名
func (h *DomainHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid domain id")
		return
	}

	domain, err := h.domainSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "domain not found")
		return
	}

	var req struct {
		DomainName   *string  `json:"domain_name"`
		RegistrarID  *uint    `json:"registrar_id"`
		Nameservers  *string  `json:"nameservers"`
		Status       *string  `json:"status"`
		AutoRenew    *bool    `json:"auto_renew"`
		WhoisPrivacy *bool    `json:"whois_privacy"`
		TransferLock *bool    `json:"transfer_lock"`
		EPPCode      *string  `json:"epp_code"`
		DNSManaged   *bool    `json:"dns_managed"`
		Price        *float64 `json:"price"`
		Currency     *string  `json:"currency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.DomainName != nil {
		domain.DomainName = *req.DomainName
	}
	if req.RegistrarID != nil {
		domain.RegistrarID = req.RegistrarID
	}
	if req.Nameservers != nil {
		domain.Nameservers = *req.Nameservers
	}
	if req.Status != nil {
		domain.Status = *req.Status
	}
	if req.AutoRenew != nil {
		domain.AutoRenew = *req.AutoRenew
	}
	if req.WhoisPrivacy != nil {
		domain.WhoisPrivacy = *req.WhoisPrivacy
	}
	if req.TransferLock != nil {
		domain.TransferLock = *req.TransferLock
	}
	if req.EPPCode != nil {
		domain.EPPCode = *req.EPPCode
	}
	if req.DNSManaged != nil {
		domain.DNSManaged = *req.DNSManaged
	}
	if req.Price != nil {
		domain.Price = *req.Price
	}
	if req.Currency != nil {
		domain.Currency = *req.Currency
	}

	if err := h.domainSvc.Update(domain); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, domain)
}

// AdminDelete 管理员删除域名
func (h *DomainHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid domain id")
		return
	}

	if err := h.domainSvc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "deleted")
}

// AdminGetTransfers 管理员获取域名转移列表
func (h *DomainHandler) AdminGetTransfers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")
	keyword := c.Query("keyword")

	params := service.DomainListParams{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Keyword:  keyword,
	}

	transfers, total, err := h.domainSvc.AdminGetTransfers(params)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, transfers, total, page, pageSize)
}

// AdminUpdateTransfer 管理员更新域名转移
func (h *DomainHandler) AdminUpdateTransfer(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid transfer id")
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=pending approved rejected completed failed"`
		Remark string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	adminID := c.GetUint("user_id")
	transfer, err := h.domainSvc.AdminUpdateTransfer(uint(id), req.Status, req.Remark, &adminID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, transfer)
}
