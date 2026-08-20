package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type SSLCertificateHandler struct {
	svc *service.SSLCertificateService
	log *logger.Logger
}

func NewSSLCertificateHandler(svc *service.SSLCertificateService, log *logger.Logger) *SSLCertificateHandler {
	return &SSLCertificateHandler{svc: svc, log: log}
}

// ═══════════════════ User-Facing ═══════════════════

// GetList returns paginated SSL certificates for the authenticated user.
func (h *SSLCertificateHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	certs, total, err := h.svc.GetList(userID, page, pageSize, keyword, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, certs, total, page, pageSize)
}

// GetDetail returns a single SSL certificate.
func (h *SSLCertificateHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid certificate id")
		return
	}
	userID := c.GetUint("user_id")

	cert, err := h.svc.GetByID(uint(id), userID)
	if err != nil {
		response.NotFound(c, "certificate not found")
		return
	}
	response.Success(c, cert)
}

// Order creates a new SSL certificate order.
func (h *SSLCertificateHandler) Order(c *gin.Context) {
	var req service.OrderCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID := c.GetUint("user_id")

	cert, err := h.svc.OrderCertificate(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, cert)
}

// GenerateCSR generates a Certificate Signing Request.
func (h *SSLCertificateHandler) GenerateCSR(c *gin.Context) {
	var req service.CSRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	result, err := h.svc.GenerateCSR(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// Validate triggers domain validation for a pending certificate.
func (h *SSLCertificateHandler) Validate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid certificate id")
		return
	}
	userID := c.GetUint("user_id")

	cert, err := h.svc.ValidateCertificate(uint(id), userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, cert)
}

// Install installs the issued certificate, key, and CA bundle.
func (h *SSLCertificateHandler) Install(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid certificate id")
		return
	}
	userID := c.GetUint("user_id")

	var req service.InstallCertificateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cert, err := h.svc.InstallCertificate(uint(id), userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, cert)
}

// Renew creates a renewal order for an existing certificate.
func (h *SSLCertificateHandler) Renew(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid certificate id")
		return
	}
	userID := c.GetUint("user_id")

	cert, err := h.svc.RenewCertificate(uint(id), userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, cert)
}

// Revoke marks a certificate as revoked.
func (h *SSLCertificateHandler) Revoke(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid certificate id")
		return
	}
	userID := c.GetUint("user_id")

	if err := h.svc.RevokeCertificate(uint(id), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "certificate revoked")
}

// GetCertificateTypes returns available SSL certificate types with pricing.
func (h *SSLCertificateHandler) GetCertificateTypes(c *gin.Context) {
	types, err := h.svc.GetCertificateTypes()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, types)
}

// ═══════════════════ Admin ═══════════════════

// AdminGetList returns all SSL certificates (admin).
func (h *SSLCertificateHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")

	var userID *uint
	if uid := c.Query("user_id"); uid != "" {
		v, _ := strconv.ParseUint(uid, 10, 64)
		u := uint(v)
		userID = &u
	}

	certs, total, err := h.svc.AdminGetList(page, pageSize, keyword, status, userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, certs, total, page, pageSize)
}

// AdminGetDetail returns a single SSL certificate (admin, no user filter).
func (h *SSLCertificateHandler) AdminGetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid certificate id")
		return
	}

	cert, err := h.svc.GetByID(uint(id), 0)
	if err != nil {
		response.NotFound(c, "certificate not found")
		return
	}
	response.Success(c, cert)
}

// AdminUpdate updates SSL certificate fields (admin).
func (h *SSLCertificateHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid certificate id")
		return
	}

	var req service.UpdateSSLCertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cert, err := h.svc.AdminUpdate(uint(id), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, cert)
}

// AdminDelete soft-deletes an SSL certificate (admin).
func (h *SSLCertificateHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid certificate id")
		return
	}

	if err := h.svc.AdminDelete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "certificate deleted")
}
