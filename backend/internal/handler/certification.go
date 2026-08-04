package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type CertificationHandler struct {
	certSvc *service.CertificationService
	log     *logger.Logger
}

func NewCertificationHandler(certSvc *service.CertificationService, log *logger.Logger) *CertificationHandler {
	return &CertificationHandler{certSvc: certSvc, log: log}
}

// Submit creates or updates a certification request.
func (h *CertificationHandler) Submit(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req service.SubmitCertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cert, err := h.certSvc.Submit(userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, cert)
}

// GetStatus returns the current user's certification status.
func (h *CertificationHandler) GetStatus(c *gin.Context) {
	userID := c.GetUint("user_id")

	cert, err := h.certSvc.GetByUserID(userID)
	if err != nil {
		response.NotFound(c, "no certification found")
		return
	}
	response.Success(c, cert)
}

// AdminGetList returns a paginated certification list.
func (h *CertificationHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status, _ := strconv.Atoi(c.DefaultQuery("status", "0"))

	certs, total, err := h.certSvc.GetList(page, pageSize, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, certs, total, page, pageSize)
}

// SubmitEnterprise creates or updates an enterprise certification request.
// POST /certification/enterprise
func (h *CertificationHandler) SubmitEnterprise(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		EnterpriseName  string `json:"enterprise_name" binding:"required"`
		RealName        string `json:"real_name" binding:"required"`
		IDCard          string `json:"id_card"`
		BusinessLicense string `json:"business_license"`
		FrontImage      string `json:"front_image"`
		BackImage       string `json:"back_image"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	certReq := service.SubmitCertificationRequest{
		Type:            "enterprise",
		RealName:        req.RealName,
		IDCard:          req.IDCard,
		EnterpriseName:  req.EnterpriseName,
		BusinessLicense: req.BusinessLicense,
		FrontImage:      req.FrontImage,
		BackImage:       req.BackImage,
	}

	cert, err := h.certSvc.Submit(userID, certReq)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, cert)
}

// AdminReview approves or rejects a certification.
func (h *CertificationHandler) AdminReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid certification id")
		return
	}

	reviewerID := c.GetUint("user_id")

	var req service.ReviewCertificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.certSvc.Review(uint(id), reviewerID, req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Status == 2 {
		response.SuccessMsg(c, "certification approved")
	} else {
		response.SuccessMsg(c, "certification rejected")
	}
}
