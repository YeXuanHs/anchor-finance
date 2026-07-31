package handler

import (
	"net/http"
	"strconv"

	"anchorfinance/internal/service"

	"github.com/gin-gonic/gin"
)

// PDFHandler PDF处理器
type PDFHandler struct {
	pdfSvc *service.PDFService
}

// NewPDFHandler 创建PDF处理器
func NewPDFHandler(pdfSvc *service.PDFService) *PDFHandler {
	return &PDFHandler{pdfSvc: pdfSvc}
}

// GenerateContractPDF 生成合同PDF
func (h *PDFHandler) GenerateContractPDF(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	pdfBytes, err := h.pdfSvc.GenerateContractPDF(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=contract.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GenerateInvoicePDF 生成发票PDF
func (h *PDFHandler) GenerateInvoicePDF(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	pdfBytes, err := h.pdfSvc.GenerateInvoicePDF(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=invoice.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GenerateContractWithSeal 生成带电子章的合同PDF
func (h *PDFHandler) GenerateContractWithSeal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid contract ID"})
		return
	}

	sealPath := c.Query("seal_path")

	pdfBytes, err := h.pdfSvc.GenerateContractWithSeal(uint(id), sealPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=contract_with_seal.pdf")
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
