package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type SSLCertificateService struct {
	db         *gorm.DB
	log        *logger.Logger
	httpClient *http.Client
}

func NewSSLCertificateService(db *gorm.DB, log *logger.Logger) *SSLCertificateService {
	return &SSLCertificateService{
		db:  db,
		log: log,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ────────────────────── Requests ──────────────────────

type OrderCertificateRequest struct {
	Domain          string `json:"domain" binding:"required,fqdn"`
	CertificateType string `json:"certificate_type" binding:"required,oneof=dv ov ev wildcard"`
	ProductID       *uint  `json:"product_id"`
	ValidationType  string `json:"validation_type" binding:"required,oneof=dns email http"`
	SANs            []string `json:"sans"`
}

type CSRRequest struct {
	Domain   string `json:"domain" binding:"required"`
	Country  string `json:"country" binding:"required,len=2"`
	State    string `json:"state"`
	City     string `json:"city"`
	Org      string `json:"org" binding:"required"`
	OrgUnit  string `json:"org_unit"`
	Email    string `json:"email" binding:"required,email"`
}

type CSRResponse struct {
	CSR        string `json:"csr"`
	PrivateKey string `json:"private_key"`
}

type InstallCertificateRequest struct {
	Certificate string `json:"certificate" binding:"required"`
	PrivateKey  string `json:"private_key" binding:"required"`
	CaBundle    string `json:"ca_bundle"`
}

type UpdateSSLCertRequest struct {
	Domain          string  `json:"domain"`
	CertificateType string  `json:"certificate_type"`
	Status          string  `json:"status"`
	Issuer          string  `json:"issuer"`
	SerialNumber    string  `json:"serial_number"`
	ExpiryDate      *string `json:"expiry_date"`
	AutoRenew       *bool   `json:"auto_renew"`
	Price           *float64 `json:"price"`
	AdminNotes      string  `json:"admin_notes"`
}

// ────────────────────── SSL Provider Types ──────────────────────

// sslProviderResponse SSL提供商API响应
type sslProviderResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// certPanelDCVResponse CertPanel DCV验证响应
type certPanelDCVResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		DCVMethod string `json:"dcv_method"`
		Record    string `json:"record"`
		Value     string `json:"value"`
		Token     string `json:"token"`
		URL       string `json:"url"`
	} `json:"data"`
}

// ────────────────────── Provider Actions ──────────────────────

// findSSLProvider 查找SSL证书的提供商配置
func (s *SSLCertificateService) findSSLProvider(cert *model.SSLCertificate) (*model.ProvisionModule, error) {
	// 如果证书关联了产品，查找产品的供应模块
	if cert.ProductID != nil {
		var product model.Product
		if err := s.db.First(&product, *cert.ProductID).Error; err == nil {
			var module model.ProvisionModule
			err := s.db.Where("active = true AND type = ?", "ssl").
				Order("priority DESC").
				First(&module).Error
			if err == nil {
				return &module, nil
			}
		}
	}

	// 查找默认的SSL供应模块
	var module model.ProvisionModule
	if err := s.db.Where("active = true AND type = ?", "ssl").
		Order("priority DESC").
		First(&module).Error; err != nil {
		return nil, fmt.Errorf("no active SSL provision module found")
	}
	return &module, nil
}

// callZJMFDcvApi 调用zjmf_api进行DCV验证
func (s *SSLCertificateService) callZJMFDcvApi(module *model.ProvisionModule, cert *model.SSLCertificate) (map[string]interface{}, error) {
	if module.ServerURL == "" {
		return nil, fmt.Errorf("zjmf API URL not configured")
	}

	apiURL := strings.TrimRight(module.ServerURL, "/") + "/api.php"

	form := url.Values{}
	form.Set("action", "SSLSubmitDCV")
	form.Set("cert_id", fmt.Sprintf("%d", cert.ID))
	form.Set("domain", cert.Domain)
	form.Set("validation_type", cert.ValidationType)
	if cert.CSR != "" {
		form.Set("csr", cert.CSR)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result string          `json:"result"`
		Msg    string          `json:"msg"`
		Data   json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if apiResp.Result != "success" {
		return nil, fmt.Errorf("zjmf error: %s", apiResp.Msg)
	}

	result := map[string]interface{}{
		"result": apiResp.Result,
		"msg":    apiResp.Msg,
	}
	if apiResp.Data != nil {
		var data interface{}
		json.Unmarshal(apiResp.Data, &data)
		result["data"] = data
	}
	return result, nil
}

// callCertPanelDcvApi 调用CertPanel API进行DCV验证
func (s *SSLCertificateService) callCertPanelDcvApi(module *model.ProvisionModule, cert *model.SSLCertificate) (map[string]interface{}, error) {
	if module.ServerURL == "" {
		return nil, fmt.Errorf("CertPanel API URL not configured")
	}

	apiURL := strings.TrimRight(module.ServerURL, "/") + "/api/v1/ssl/dcv"

	payload := map[string]interface{}{
		"domain":          cert.Domain,
		"validation_type": cert.ValidationType,
		"certificate_type": cert.CertificateType,
	}
	if cert.CSR != "" {
		payload["csr"] = cert.CSR
	}

	payloadBytes, _ := json.Marshal(payload)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if module.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+module.APIKey)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var apiResp certPanelDCVResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if apiResp.Code != 0 && apiResp.Code != 200 {
		return nil, fmt.Errorf("CertPanel error (code=%d): %s", apiResp.Code, apiResp.Message)
	}

	result := map[string]interface{}{
		"code":    apiResp.Code,
		"message": apiResp.Message,
		"data": map[string]interface{}{
			"dcv_method": apiResp.Data.DCVMethod,
			"record":     apiResp.Data.Record,
			"value":      apiResp.Data.Value,
			"token":      apiResp.Data.Token,
			"url":        apiResp.Data.URL,
		},
	}
	return result, nil
}

// callZJMFIssueCert 调用zjmf_api获取证书签发状态
func (s *SSLCertificateService) callZJMFIssueCert(module *model.ProvisionModule, cert *model.SSLCertificate) (map[string]interface{}, error) {
	if module.ServerURL == "" {
		return nil, fmt.Errorf("zjmf API URL not configured")
	}

	apiURL := strings.TrimRight(module.ServerURL, "/") + "/api.php"

	form := url.Values{}
	form.Set("action", "SSLGetCert")
	form.Set("cert_id", fmt.Sprintf("%d", cert.ID))
	form.Set("domain", cert.Domain)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result string          `json:"result"`
		Msg    string          `json:"msg"`
		Data   json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if apiResp.Result != "success" {
		return nil, fmt.Errorf("zjmf error: %s", apiResp.Msg)
	}

	result := map[string]interface{}{
		"result": apiResp.Result,
		"msg":    apiResp.Msg,
	}
	if apiResp.Data != nil {
		var data interface{}
		json.Unmarshal(apiResp.Data, &data)
		result["data"] = data
	}
	return result, nil
}

// ────────────────────── Certificate Types ──────────────────────

// GetCertificateTypes returns all enabled SSL certificate types with pricing.
func (s *SSLCertificateService) GetCertificateTypes() ([]model.SSLCertificateType, error) {
	var types []model.SSLCertificateType
	if err := s.db.Where("status = 1").Order("sort_order ASC").Find(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}

// ────────────────────── User CRUD ──────────────────────

// GetList returns paginated SSL certificates for a user.
func (s *SSLCertificateService) GetList(userID uint, page, pageSize int, keyword, status string) ([]model.SSLCertificate, int64, error) {
	var certs []model.SSLCertificate
	var total int64

	query := s.db.Model(&model.SSLCertificate{}).Where("user_id = ?", userID)
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("domain LIKE ?", q)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&certs).Error; err != nil {
		return nil, 0, err
	}
	return certs, total, nil
}

// GetByID returns a single certificate. Verifies ownership when userID > 0.
func (s *SSLCertificateService) GetByID(id, userID uint) (*model.SSLCertificate, error) {
	var cert model.SSLCertificate
	query := s.db
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if err := query.First(&cert, id).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

// ────────────────────── Order / Generate / Validate / Install / Renew / Revoke ──────────────────────

// OrderCertificate creates a new SSL certificate order.
func (s *SSLCertificateService) OrderCertificate(userID uint, req OrderCertificateRequest) (*model.SSLCertificate, error) {
	// Look up price from certificate type
	var certType model.SSLCertificateType
	if err := s.db.Where("code = ? AND status = 1", req.CertificateType).First(&certType).Error; err != nil {
		return nil, fmt.Errorf("certificate type %q not available", req.CertificateType)
	}

	sansJSON := "[]"
	if len(req.SANs) > 0 {
		b, _ := json.Marshal(req.SANs)
		sansJSON = string(b)
	}

	cert := &model.SSLCertificate{
		UserID:          userID,
		Domain:          req.Domain,
		ProductID:       req.ProductID,
		CertificateType: req.CertificateType,
		Status:          "pending",
		ValidationType:  req.ValidationType,
		SANs:            sansJSON,
		Price:           certType.Price,
		OrderID:         fmt.Sprintf("SSL-%d-%d", userID, time.Now().UnixMilli()),
	}

	if err := s.db.Create(cert).Error; err != nil {
		return nil, err
	}

	s.log.Infof("SSL certificate ordered: user=%d domain=%s type=%s", userID, req.Domain, req.CertificateType)
	return cert, nil
}

// GenerateCSR creates a CSR and private key for the given parameters.
func (s *SSLCertificateService) GenerateCSR(req CSRRequest) (*CSRResponse, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate private key: %w", err)
	}

	template := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   req.Domain,
			Organization: []string{req.Org},
			Country:      []string{req.Country},
		},
		EmailAddresses: []string{req.Email},
	}
	if req.OrgUnit != "" {
		template.Subject.OrganizationalUnit = []string{req.OrgUnit}
	}
	if req.State != "" {
		template.Subject.Province = []string{req.State}
	}
	if req.City != "" {
		template.Subject.Locality = []string{req.City}
	}

	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &template, key)
	if err != nil {
		return nil, fmt.Errorf("create CSR: %w", err)
	}

	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("marshal private key: %w", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	return &CSRResponse{CSR: string(csrPEM), PrivateKey: string(keyPEM)}, nil
}

// ValidateCertificate triggers domain validation for a pending certificate.
func (s *SSLCertificateService) ValidateCertificate(id, userID uint) (*model.SSLCertificate, error) {
	cert, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if cert.Status != "pending" {
		return nil, fmt.Errorf("certificate is %s, cannot validate", cert.Status)
	}

	s.log.Infof("SSL validation started: id=%d type=%s domain=%s", cert.ID, cert.ValidationType, cert.Domain)

	// 查找SSL提供商
	module, providerErr := s.findSSLProvider(cert)

	var validationData string

	if providerErr == nil && module != nil {
		// 有提供商配置，调用提供商API
		slug := strings.ToLower(module.Slug)
		switch {
		case strings.Contains(slug, "zjmf"):
			result, err := s.callZJMFDcvApi(module, cert)
			if err != nil {
				s.log.Warnf("zjmf DCV API failed for cert %d: %v", cert.ID, err)
			} else {
				validationData = s.buildValidationDataFromProvider(cert.ValidationType, result)
			}
		case strings.Contains(slug, "certpanel"):
			result, err := s.callCertPanelDcvApi(module, cert)
			if err != nil {
				s.log.Warnf("CertPanel DCV API failed for cert %d: %v", cert.ID, err)
			} else {
				validationData = s.buildValidationDataFromProvider(cert.ValidationType, result)
			}
		default:
			s.log.Infof("no specific DCV handler for provider %s, using local validation", module.Slug)
		}
	}

	// 如果提供商调用失败或没有提供商，使用本地验证数据
	if validationData == "" {
		switch cert.ValidationType {
		case "dns":
			validationData = `{"type":"dns","record":"_acme-challenge.` + cert.Domain + `","value":"pending-validation"}`
		case "email":
			validationData = `{"type":"email","sent_to":"admin@` + cert.Domain + `"}`
		case "http":
			validationData = `{"type":"http","path":"/.well-known/acme-challenge/pending","content":"pending-validation"}`
		default:
			return nil, fmt.Errorf("unsupported validation type: %s", cert.ValidationType)
		}
	}

	if err := s.db.Model(cert).Update("validation_data", validationData).Error; err != nil {
		return nil, err
	}
	return cert, nil
}

// buildValidationDataFromProvider 从提供商响应构建验证数据
func (s *SSLCertificateService) buildValidationDataFromProvider(validationType string, result map[string]interface{}) string {
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		return ""
	}

	switch validationType {
	case "dns":
		record, _ := data["record"].(string)
		value, _ := data["value"].(string)
		if record != "" && value != "" {
			return fmt.Sprintf(`{"type":"dns","record":"%s","value":"%s"}`, record, value)
		}
	case "email":
		sentTo, _ := data["sent_to"].(string)
		if sentTo != "" {
			return fmt.Sprintf(`{"type":"email","sent_to":"%s"}`, sentTo)
		}
	case "http":
		token, _ := data["token"].(string)
		urlPath, _ := data["url"].(string)
		if token != "" {
			return fmt.Sprintf(`{"type":"http","path":"%s","content":"%s"}`, urlPath, token)
		}
	}
	return ""
}

// SubmitDCV 提交DCV验证（供前端调用）
func (s *SSLCertificateService) SubmitDCV(id, userID uint) (*model.SSLCertificate, error) {
	cert, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if cert.Status != "pending" {
		return nil, fmt.Errorf("certificate is %s, cannot submit DCV", cert.Status)
	}

	// 查找SSL提供商
	module, providerErr := s.findSSLProvider(cert)
	if providerErr != nil {
		return nil, fmt.Errorf("no SSL provider configured: %w", providerErr)
	}

	slug := strings.ToLower(module.Slug)
	var result map[string]interface{}

	switch {
	case strings.Contains(slug, "zjmf"):
		result, err = s.callZJMFDcvApi(module, cert)
	case strings.Contains(slug, "certpanel"):
		result, err = s.callCertPanelDcvApi(module, cert)
	default:
		return nil, fmt.Errorf("unsupported SSL provider: %s", module.Slug)
	}

	if err != nil {
		return nil, fmt.Errorf("submit DCV failed: %w", err)
	}

	// 更新验证数据
	validationData := s.buildValidationDataFromProvider(cert.ValidationType, result)
	if validationData != "" {
		s.db.Model(cert).Update("validation_data", validationData)
	}

	// 轮询检查证书签发状态
	go s.pollCertificateIssuance(cert, module)

	s.log.Infof("SSL DCV submitted: id=%d domain=%s provider=%s", cert.ID, cert.Domain, module.Slug)
	return s.GetByID(id, 0)
}

// pollCertificateIssuance 轮询证书签发状态
func (s *SSLCertificateService) pollCertificateIssuance(cert *model.SSLCertificate, module *model.ProvisionModule) {
	maxAttempts := 30 // 最多轮询30次
	for i := 0; i < maxAttempts; i++ {
		time.Sleep(30 * time.Second)

		slug := strings.ToLower(module.Slug)
		var result map[string]interface{}
		var err error

		switch {
		case strings.Contains(slug, "zjmf"):
			result, err = s.callZJMFIssueCert(module, cert)
		default:
			// 其他提供商暂不支持自动轮询
			return
		}

		if err != nil {
			s.log.Warnf("poll cert issuance failed for cert %d: %v", cert.ID, err)
			continue
		}

		if data, ok := result["data"].(map[string]interface{}); ok {
			status, _ := data["status"].(string)
			if status == "issued" || status == "active" {
				// 证书已签发
				certificate, _ := data["certificate"].(string)
				caBundle, _ := data["ca_bundle"].(string)
				issuer, _ := data["issuer"].(string)

				now := time.Now()
				updates := map[string]interface{}{
					"status":      "issued",
					"issue_date":  &now,
				}
				if certificate != "" {
					updates["certificate"] = certificate
				}
				if caBundle != "" {
					updates["ca_bundle"] = caBundle
				}
				if issuer != "" {
					updates["issuer"] = issuer
				}

				s.db.Model(cert).Updates(updates)
				s.log.Infof("SSL certificate issued: id=%d domain=%s", cert.ID, cert.Domain)
				return
			}
		}
	}

	s.log.Warnf("SSL certificate issuance poll timeout: id=%d domain=%s", cert.ID, cert.Domain)
}

// InstallCertificate stores the issued certificate, key, and CA bundle.
func (s *SSLCertificateService) InstallCertificate(id, userID uint, req InstallCertificateRequest) (*model.SSLCertificate, error) {
	cert, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if cert.Status != "pending" && cert.Status != "issued" {
		return nil, fmt.Errorf("certificate is %s, cannot install", cert.Status)
	}

	now := time.Now()
	updates := map[string]interface{}{
		"certificate": req.Certificate,
		"private_key": req.PrivateKey,
		"ca_bundle":   req.CaBundle,
		"status":      "issued",
		"issue_date":  &now,
	}
	if err := s.db.Model(cert).Updates(updates).Error; err != nil {
		return nil, err
	}

	s.log.Infof("SSL certificate installed: id=%d domain=%s", cert.ID, cert.Domain)
	return s.GetByID(id, 0)
}

// RenewCertificate creates a renewal order for an existing certificate.
func (s *SSLCertificateService) RenewCertificate(id, userID uint) (*model.SSLCertificate, error) {
	cert, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}
	if cert.Status != "issued" && cert.Status != "expired" {
		return nil, fmt.Errorf("certificate is %s, cannot renew", cert.Status)
	}

	newCert := &model.SSLCertificate{
		UserID:          cert.UserID,
		Domain:          cert.Domain,
		ProductID:       cert.ProductID,
		CertificateType: cert.CertificateType,
		Status:          "pending",
		ValidationType:  cert.ValidationType,
		SANs:            cert.SANs,
		Price:           cert.Price,
		AutoRenew:       cert.AutoRenew,
		OrderID:         fmt.Sprintf("SSL-RENEW-%d-%d", cert.ID, time.Now().UnixMilli()),
	}

	if err := s.db.Create(newCert).Error; err != nil {
		return nil, err
	}

	s.log.Infof("SSL certificate renewed: old=%d new=%d domain=%s", cert.ID, newCert.ID, cert.Domain)
	return newCert, nil
}

// RevokeCertificate marks a certificate as revoked.
func (s *SSLCertificateService) RevokeCertificate(id, userID uint) error {
	cert, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}
	if cert.Status != "issued" {
		return fmt.Errorf("certificate is %s, cannot revoke", cert.Status)
	}

	if err := s.db.Model(cert).Update("status", "revoked").Error; err != nil {
		return err
	}

	s.log.Infof("SSL certificate revoked: id=%d domain=%s", cert.ID, cert.Domain)
	return nil
}

// ────────────────────── Admin ──────────────────────

// AdminGetList returns all SSL certificates with admin filters.
func (s *SSLCertificateService) AdminGetList(page, pageSize int, keyword, status string, userID *uint) ([]model.SSLCertificate, int64, error) {
	var certs []model.SSLCertificate
	var total int64

	query := s.db.Model(&model.SSLCertificate{})
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("domain LIKE ?", q)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&certs).Error; err != nil {
		return nil, 0, err
	}
	return certs, total, nil
}

// AdminUpdate updates certificate fields (admin operation).
func (s *SSLCertificateService) AdminUpdate(id uint, req UpdateSSLCertRequest) (*model.SSLCertificate, error) {
	var cert model.SSLCertificate
	if err := s.db.First(&cert, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("certificate not found")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Domain != "" {
		updates["domain"] = req.Domain
	}
	if req.CertificateType != "" {
		updates["certificate_type"] = req.CertificateType
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}
	if req.Issuer != "" {
		updates["issuer"] = req.Issuer
	}
	if req.SerialNumber != "" {
		updates["serial_number"] = req.SerialNumber
	}
	if req.ExpiryDate != nil {
		t, err := time.Parse("2006-01-02", *req.ExpiryDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expiry_date format, use YYYY-MM-DD")
		}
		updates["expiry_date"] = &t
	}
	if req.AutoRenew != nil {
		updates["auto_renew"] = *req.AutoRenew
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}

	if len(updates) == 0 {
		return &cert, nil
	}

	if err := s.db.Model(&cert).Updates(updates).Error; err != nil {
		return nil, err
	}

	return s.GetByID(id, 0)
}

// AdminDelete soft-deletes a certificate.
func (s *SSLCertificateService) AdminDelete(id uint) error {
	result := s.db.Delete(&model.SSLCertificate{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("certificate not found")
	}
	return nil
}

// ExpireCertificates marks certificates past their expiry date as expired (cron task).
func (s *SSLCertificateService) ExpireCertificates() (int64, error) {
	now := time.Now()
	result := s.db.Model(&model.SSLCertificate{}).
		Where("status = ? AND expiry_date IS NOT NULL AND expiry_date < ?", "issued", now).
		Update("status", "expired")
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected > 0 {
		s.log.Infof("expired %d SSL certificates", result.RowsAffected)
	}
	return result.RowsAffected, nil
}
