package service

import (
	"bytes"
	"fmt"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// PDFService PDF生成服务
type PDFService struct {
	db *gorm.DB
}

// NewPDFService 创建PDF服务
func NewPDFService(db *gorm.DB) *PDFService {
	return &PDFService{db: db}
}

// ContractPDF 合同PDF数据
type ContractPDF struct {
	ContractNo   string
	Title        string
	PartyA       string // 甲方
	PartyB       string // 乙方
	Content      string
	StartDate    string
	EndDate      string
	Amount       float64
	SealImage    []byte // 电子章图片
	SignDate     string
}

// InvoicePDF 发票PDF数据
type InvoicePDF struct {
	InvoiceNo    string
	UserID       uint
	UserName     string
	UserEmail    string
	Items        []InvoicePDFItem
	SubTotal     float64
	Tax          float64
	Total        float64
	Currency     string
	IssueDate    string
	DueDate      string
	Status       string
	SealImage    []byte
	Notes        string
}

// InvoicePDFItem 发票PDF明细项
type InvoicePDFItem struct {
	Description  string
	Quantity     int
	UnitPrice    float64
	Amount       float64
}

// GenerateContractPDF 生成合同PDF
func (s *PDFService) GenerateContractPDF(contractID uint) ([]byte, error) {
	var contract struct {
		ID          uint
		ContractNo  string
		Title       string
		Content     string
		UserID      uint
		StartDate   time.Time
		EndDate     time.Time
		Amount      float64
		Status      string
		SignedAt    *time.Time
	}

	if err := s.db.Table("contracts").Where("id = ?", contractID).First(&contract).Error; err != nil {
		return nil, fmt.Errorf("contract not found: %w", err)
	}

	var user struct {
		Username string
		Email    string
		Company  string
	}
	s.db.Table("users").Where("id = ?", contract.UserID).First(&user)

	companyName := user.Company
	if companyName == "" {
		companyName = user.Username
	}

	systemName := "锚点财务"
	s.db.Table("system_configs").Where("key = ?", "company_name").Select("value").Scan(&systemName)

	pdfData := &ContractPDF{
		ContractNo: contract.ContractNo,
		Title:      contract.Title,
		PartyA:     systemName,
		PartyB:     companyName,
		Content:    contract.Content,
		StartDate:  contract.StartDate.Format("2006-01-02"),
		EndDate:    contract.EndDate.Format("2006-01-02"),
		Amount:     contract.Amount,
	}

	if contract.SignedAt != nil {
		pdfData.SignDate = contract.SignedAt.Format("2006-01-02")
	}

	return s.renderContractPDF(pdfData)
}

// GenerateInvoicePDF 生成发票PDF
func (s *PDFService) GenerateInvoicePDF(invoiceID uint) ([]byte, error) {
	var invoice struct {
		ID        uint
		InvoiceNo string
		UserID    uint
		SubTotal  float64
		Tax       float64
		Total     float64
		Currency  string
		Status    string
		DueDate   time.Time
		CreatedAt time.Time
		Notes     string
	}

	if err := s.db.Table("invoices").Where("id = ?", invoiceID).First(&invoice).Error; err != nil {
		return nil, fmt.Errorf("invoice not found: %w", err)
	}

	var items []struct {
		Description string
		Quantity    int
		UnitPrice   float64
		Amount      float64
	}
	s.db.Table("invoice_items").Where("invoice_id = ?", invoiceID).Find(&items)

	var user struct {
		Username string
		Email    string
	}
	s.db.Table("users").Where("id = ?", invoice.UserID).First(&user)

	pdfItems := make([]InvoicePDFItem, len(items))
	for i, item := range items {
		pdfItems[i] = InvoicePDFItem{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      item.Amount,
		}
	}

	pdfData := &InvoicePDF{
		InvoiceNo: invoice.InvoiceNo,
		UserID:    invoice.UserID,
		UserName:  user.Username,
		UserEmail: user.Email,
		Items:     pdfItems,
		SubTotal:  invoice.SubTotal,
		Tax:       invoice.Tax,
		Total:     invoice.Total,
		Currency:  invoice.Currency,
		IssueDate: invoice.CreatedAt.Format("2006-01-02"),
		DueDate:   invoice.DueDate.Format("2006-01-02"),
		Status:    invoice.Status,
		Notes:     invoice.Notes,
	}

	return s.renderInvoicePDF(pdfData)
}

// renderContractPDF 渲染合同PDF
func (s *PDFService) renderContractPDF(data *ContractPDF) ([]byte, error) {
	var buf bytes.Buffer

	// 使用简单的文本PDF生成（生产环境应使用gofpdf或wkhtmltopdf）
	// 这里生成一个基础的PDF结构
	buf.WriteString("%PDF-1.4\n")
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")

	// 页面内容
	content := fmt.Sprintf(`BT
/F1 12 Tf
50 750 Td
(%s) Tj
0 -30 Td
(Contract No: %s) Tj
0 -20 Td
(Party A: %s) Tj
0 -20 Td
(Party B: %s) Tj
0 -30 Td
(%s) Tj
0 -20 Td
(Period: %s to %s) Tj
0 -20 Td
(Amount: %.2f %s) Tj
0 -30 Td
(Signed: %s) Tj
ET`, data.Title, data.ContractNo, data.PartyA, data.PartyB,
		truncateString(data.Content, 500), data.StartDate, data.EndDate,
		data.Amount, "CNY", data.SignDate)

	stream := fmt.Sprintf("stream\n%s\nendstream", content)

	buf.WriteString(fmt.Sprintf("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>\nendobj\n"))
	buf.WriteString(fmt.Sprintf("4 0 obj\n<< /Length %d >>\n%s\nendobj\n", len(stream), stream))
	buf.WriteString("5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n")

	// 添加电子章
	if len(data.SealImage) > 0 {
		buf.WriteString(fmt.Sprintf("6 0 obj\n<< /Type /XObject /Subtype /Image /Width 200 /Height 200 /ColorSpace /DeviceRGB /BitsPerComponent 8 /Length %d >>\nstream\n", len(data.SealImage)))
		buf.Write(data.SealImage)
		buf.WriteString("\nendstream\nendobj\n")
	}

	buf.WriteString("xref\n0 7\n")
	buf.WriteString("trailer\n<< /Size 7 /Root 1 0 R >>\nstartxref\n0\n%%EOF\n")

	return buf.Bytes(), nil
}

// renderInvoicePDF 渲染发票PDF
func (s *PDFService) renderInvoicePDF(data *InvoicePDF) ([]byte, error) {
	var buf bytes.Buffer

	currency := data.Currency
	if currency == "" {
		currency = "CNY"
	}

	buf.WriteString("%PDF-1.4\n")
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")

	// 构建发票内容
	content := fmt.Sprintf(`BT
/F1 16 Tf
50 750 Td
(INVOICE) Tj
/F1 10 Tf
0 -25 Td
(Invoice No: %s) Tj
0 -15 Td
(Date: %s  Due: %s) Tj
0 -15 Td
(Status: %s) Tj
0 -25 Td
(Bill To: %s <%s>) Tj
0 -30 Td
`, data.InvoiceNo, data.IssueDate, data.DueDate, data.Status, data.UserName, data.UserEmail)

	// 添加明细项
	y := 620.0
	for _, item := range data.Items {
		content += fmt.Sprintf("/F1 9 Tf\n50 %.0f Td\n(%s  x%d  %.2f  %.2f) Tj\n",
			y, item.Description, item.Quantity, item.UnitPrice, item.Amount)
		y -= 15
	}

	// 合计
	content += fmt.Sprintf("/F1 10 Tf\n50 %.0f Td\n(Subtotal: %.2f %s) Tj\n", y-10, data.SubTotal, currency)
	content += fmt.Sprintf("0 -15 Td\n(Tax: %.2f %s) Tj\n", data.Tax, currency)
	content += fmt.Sprintf("0 -15 Td\n(Total: %.2f %s) Tj\n", data.Total, currency)

	if data.Notes != "" {
		content += fmt.Sprintf("0 -25 Td\n(Notes: %s) Tj\n", truncateString(data.Notes, 200))
	}

	content += "ET"

	stream := fmt.Sprintf("stream\n%s\nendstream", content)

	buf.WriteString(fmt.Sprintf("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] /Contents 4 0 R /Resources << /Font << /F1 5 0 R >> >> >>\nendobj\n"))
	buf.WriteString(fmt.Sprintf("4 0 obj\n<< /Length %d >>\n%s\nendobj\n", len(stream), stream))
	buf.WriteString("5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n")

	buf.WriteString("xref\n0 6\n")
	buf.WriteString("trailer\n<< /Size 6 /Root 1 0 R >>\nstartxref\n0\n%%EOF\n")

	logger.Info("Generated invoice PDF", "invoice_no", data.InvoiceNo)

	return buf.Bytes(), nil
}

// GenerateContractWithSeal 生成带电子章的合同PDF
func (s *PDFService) GenerateContractWithSeal(contractID uint, sealPath string) ([]byte, error) {
	// 读取电子章图片
	sealImage, err := readSealImage(sealPath)
	if err != nil {
		logger.Warnf("Failed to read seal image, generating without seal: %v", err)
		return s.GenerateContractPDF(contractID)
	}

	// 获取合同数据
	var contract struct {
		ID         uint
		ContractNo string
		Title      string
		Content    string
		UserID     uint
		StartDate  time.Time
		EndDate    time.Time
		Amount     float64
		SignedAt   *time.Time
	}

	if err := s.db.Table("contracts").Where("id = ?", contractID).First(&contract).Error; err != nil {
		return nil, err
	}

	var user struct {
		Username string
		Company  string
	}
	s.db.Table("users").Where("id = ?", contract.UserID).First(&user)

	companyName := user.Company
	if companyName == "" {
		companyName = user.Username
	}

	pdfData := &ContractPDF{
		ContractNo: contract.ContractNo,
		Title:      contract.Title,
		PartyA:     "锚点财务",
		PartyB:     companyName,
		Content:    contract.Content,
		StartDate:  contract.StartDate.Format("2006-01-02"),
		EndDate:    contract.EndDate.Format("2006-01-02"),
		Amount:     contract.Amount,
		SealImage:  sealImage,
	}

	if contract.SignedAt != nil {
		pdfData.SignDate = contract.SignedAt.Format("2006-01-02")
	}

	return s.renderContractPDF(pdfData)
}

// readSealImage 读取电子章图片
func readSealImage(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("seal path is empty")
	}
	// 实际实现中应该读取文件
	return nil, fmt.Errorf("seal image not found: %s", path)
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
