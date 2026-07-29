package service

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// DomainService 域名服务
type DomainService struct {
	db         *gorm.DB
	log        *logger.Logger
	httpClient *http.Client
}

// NewDomainService 创建域名服务实例
func NewDomainService(db *gorm.DB, log *logger.Logger) *DomainService {
	return &DomainService{
		db:  db,
		log: log,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// DomainListParams 域名列表查询参数
type DomainListParams struct {
	Page     int
	PageSize int
	Status   string
	Keyword  string
}

// GetList 获取域名分页列表
func (s *DomainService) GetList(userID uint, params DomainListParams) ([]model.Domain, int64, error) {
	var domains []model.Domain
	var total int64

	query := s.db.Model(&model.Domain{}).Where("user_id = ?", userID)

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Keyword != "" {
		q := "%" + params.Keyword + "%"
		query = query.Where("domain_name LIKE ?", q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Order("created_at DESC").Find(&domains).Error; err != nil {
		return nil, 0, err
	}

	return domains, total, nil
}

// AdminGetList 管理员获取域名分页列表
func (s *DomainService) AdminGetList(params DomainListParams) ([]model.Domain, int64, error) {
	var domains []model.Domain
	var total int64

	query := s.db.Model(&model.Domain{})

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Keyword != "" {
		q := "%" + params.Keyword + "%"
		query = query.Where("domain_name LIKE ?", q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Preload("User").Order("created_at DESC").Find(&domains).Error; err != nil {
		return nil, 0, err
	}

	return domains, total, nil
}

// GetByID 根据ID获取域名
func (s *DomainService) GetByID(id uint) (*model.Domain, error) {
	var domain model.Domain
	if err := s.db.Preload("User").First(&domain, id).Error; err != nil {
		return nil, err
	}
	return &domain, nil
}

// GetByIDAndUser 根据ID和用户ID获取域名
func (s *DomainService) GetByIDAndUser(id, userID uint) (*model.Domain, error) {
	var domain model.Domain
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&domain).Error; err != nil {
		return nil, err
	}
	return &domain, nil
}

// GetByUserID 获取用户的所有域名
func (s *DomainService) GetByUserID(userID uint) ([]model.Domain, error) {
	var domains []model.Domain
	if err := s.db.Where("user_id = ?", userID).Find(&domains).Error; err != nil {
		return nil, err
	}
	return domains, nil
}

// Create 创建域名
func (s *DomainService) Create(domain *model.Domain) error {
	if err := s.db.Create(domain).Error; err != nil {
		s.log.Errorf("创建域名失败: %v", err)
		return err
	}
	return nil
}

// Update 更新域名
func (s *DomainService) Update(domain *model.Domain) error {
	if err := s.db.Save(domain).Error; err != nil {
		s.log.Errorf("更新域名失败: %v", err)
		return err
	}
	return nil
}

// Delete 删除域名（软删除）
func (s *DomainService) Delete(id uint) error {
	if err := s.db.Delete(&model.Domain{}, id).Error; err != nil {
		s.log.Errorf("删除域名失败: %v", err)
		return err
	}
	return nil
}

// Renew 续费域名
func (s *DomainService) Renew(id uint, years int) (*model.Domain, error) {
	if years <= 0 || years > 10 {
		return nil, errors.New("续费年限必须在1-10之间")
	}

	var domain model.Domain
	if err := s.db.First(&domain, id).Error; err != nil {
		return nil, err
	}

	if domain.Status != "active" && domain.Status != "expired" {
		return nil, errors.New("当前状态不允许续费")
	}

	now := time.Now()
	if domain.ExpiryDate == nil {
		expiry := now.AddDate(years, 0, 0)
		domain.ExpiryDate = &expiry
	} else {
		expiry := domain.ExpiryDate.AddDate(years, 0, 0)
		domain.ExpiryDate = &expiry
	}

	domain.NextDueDate = domain.ExpiryDate
	domain.Status = "active"

	if err := s.db.Save(&domain).Error; err != nil {
		s.log.Errorf("域名续费失败: %v", err)
		return nil, err
	}

	return &domain, nil
}

// GetDNSRecords 获取域名DNS记录
func (s *DomainService) GetDNSRecords(domainID uint) ([]model.DomainDNSRecord, error) {
	var records []model.DomainDNSRecord
	if err := s.db.Where("domain_id = ?", domainID).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// AddDNSRecord 添加DNS记录
func (s *DomainService) AddDNSRecord(record *model.DomainDNSRecord) error {
	var domain model.Domain
	if err := s.db.First(&domain, record.DomainID).Error; err != nil {
		return errors.New("域名不存在")
	}

	if !domain.DNSManaged {
		return errors.New("该域名未启用DNS管理")
	}

	if err := s.db.Create(record).Error; err != nil {
		s.log.Errorf("添加DNS记录失败: %v", err)
		return err
	}
	return nil
}

// UpdateDNSRecord 更新DNS记录
func (s *DomainService) UpdateDNSRecord(record *model.DomainDNSRecord) error {
	if err := s.db.Save(record).Error; err != nil {
		s.log.Errorf("更新DNS记录失败: %v", err)
		return err
	}
	return nil
}

// DeleteDNSRecord 删除DNS记录
func (s *DomainService) DeleteDNSRecord(id uint) error {
	if err := s.db.Delete(&model.DomainDNSRecord{}, id).Error; err != nil {
		s.log.Errorf("删除DNS记录失败: %v", err)
		return err
	}
	return nil
}

// GetDNSRecordByID 根据ID获取DNS记录
func (s *DomainService) GetDNSRecordByID(id uint) (*model.DomainDNSRecord, error) {
	var record model.DomainDNSRecord
	if err := s.db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// ==================== 域名可用性检查 ====================

// whoisServer WHOIS服务器映射
var whoisServers = map[string]string{
	"com":    "whois.verisign-grs.com",
	"net":    "whois.verisign-grs.com",
	"org":    "whois.pir.org",
	"info":   "whois.afilias.net",
	"biz":    "whois.neulevel.biz",
	"cn":     "whois.cnnic.cn",
	"com.cn": "whois.cnnic.cn",
	"net.cn": "whois.cnnic.cn",
	"org.cn": "whois.cnnic.cn",
	"io":     "whois.nic.io",
	"co":     "whois.nic.co",
	"me":     "whois.nic.me",
	"cc":     "whois.nic.cc",
	"tv":     "whois.nic.tv",
}

// getWhoisServer 获取域名对应的WHOIS服务器
func getWhoisServer(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return ""
	}

	// 先尝试完整后缀（如 com.cn）
	if len(parts) >= 3 {
		suffix := strings.Join(parts[len(parts)-2:], ".")
		if server, ok := whoisServers[suffix]; ok {
			return server
		}
	}

	// 再尝试顶级域名
	tld := parts[len(parts)-1]
	if server, ok := whoisServers[tld]; ok {
		return server
	}

	// 默认使用 whois.iana.org
	return "whois.iana.org"
}

// checkWhoisAvail 通过WHOIS查询检查域名是否可用
func (s *DomainService) checkWhoisAvail(domainName string) (bool, error) {
	whoisServer := getWhoisServer(domainName)
	if whoisServer == "" {
		return false, fmt.Errorf("no WHOIS server found for domain: %s", domainName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", whoisServer+":43")
	if err != nil {
		return false, fmt.Errorf("connect to WHOIS server %s: %w", whoisServer, err)
	}
	defer conn.Close()

	// 设置读写超时
	conn.SetDeadline(time.Now().Add(10 * time.Second))

	// 发送查询
	_, err = fmt.Fprintf(conn, "%s\r\n", domainName)
	if err != nil {
		return false, fmt.Errorf("send WHOIS query: %w", err)
	}

	// 读取响应
	var response strings.Builder
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Text()
		response.WriteString(line + "\n")
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("read WHOIS response: %w", err)
	}

	respText := response.String()

	// 判断域名是否可用
	// 常见的"未找到"响应模式
	notFoundPatterns := []string{
		"No match for",
		"NOT FOUND",
		"No Data Found",
		"Domain not found",
		"Status: free",
		"AVAILABLE",
		"No entries found",
		"Domain Status: free",
		"is free",
		"not been registered",
		"No match",
	}

	upperResp := strings.ToUpper(respText)
	for _, pattern := range notFoundPatterns {
		if strings.Contains(upperResp, strings.ToUpper(pattern)) {
			return true, nil
		}
	}

	// 如果找到注册信息，域名已被占用
	if strings.Contains(respText, "Domain Name:") || strings.Contains(respText, "domain:") {
		return false, nil
	}

	// 默认认为不可用（安全保守）
	return false, nil
}

// checkUpstreamDomainAvail 通过上游注册商API检查域名可用性
func (s *DomainService) checkUpstreamDomainAvail(provider *model.UpstreamProvider, domainName string) (bool, error) {
	if provider.APIURL == "" {
		return false, fmt.Errorf("upstream provider API URL not configured")
	}

	switch provider.Type {
	case "zjmf", "zjmfv3":
		return s.checkZJMFDomaiAvail(provider, domainName)
	case "whmcs":
		return s.checkWHMCSDomainAvail(provider, domainName)
	default:
		return false, fmt.Errorf("unsupported upstream provider type: %s", provider.Type)
	}
}

// checkZJMFDomaiAvail 通过zjmf_api检查域名可用性
func (s *DomainService) checkZJMFDomaiAvail(provider *model.UpstreamProvider, domainName string) (bool, error) {
	apiURL := strings.TrimRight(provider.APIURL, "/") + "/api.php"

	form := url.Values{}
	form.Set("action", "DomainCheck")
	form.Set("domain", domainName)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return false, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return false, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result string `json:"result"`
		Msg    string `json:"msg"`
		Data   struct {
			Available bool   `json:"available"`
			Domain    string `json:"domain"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return false, fmt.Errorf("parse response: %w", err)
	}

	if apiResp.Result != "success" {
		return false, fmt.Errorf("API error: %s", apiResp.Msg)
	}

	return apiResp.Data.Available, nil
}

// checkWHMCSDomainAvail 通过WHMCS API检查域名可用性
func (s *DomainService) checkWHMCSDomainAvail(provider *model.UpstreamProvider, domainName string) (bool, error) {
	apiURL := strings.TrimRight(provider.APIURL, "/") + "/includes/api.php"

	form := url.Values{}
	form.Set("action", "DomainCheck")
	form.Set("identifier", provider.APIKey)
	form.Set("secret", provider.Config["api_secret"].(string))
	form.Set("responsetype", "json")
	form.Set("domain", domainName)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(form.Encode()))
	if err != nil {
		return false, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return false, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result      string `json:"result"`
		Status      string `json:"status"`
		Unavailable bool   `json:"unavailable"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return false, fmt.Errorf("parse response: %w", err)
	}

	if apiResp.Result != "success" {
		return false, fmt.Errorf("API error: status=%s", apiResp.Status)
	}

	return !apiResp.Unavailable, nil
}

// CheckAvailability 检查域名是否可用
func (s *DomainService) CheckAvailability(domainName string) (bool, error) {
	// 1. 先检查本地数据库
	var count int64
	if err := s.db.Model(&model.Domain{}).Where("domain_name = ?", domainName).Count(&count).Error; err != nil {
		return false, err
	}
	if count > 0 {
		return false, nil
	}

	// 2. 检查是否有上游注册商配置
	var provider model.UpstreamProvider
	err := s.db.Where("type IN ? AND is_active = ?", []string{"zjmf", "zjmfv3", "whmcs"}, true).
		Order("id ASC").First(&provider).Error
	if err == nil {
		// 有上游注册商，调用其API
		available, upErr := s.checkUpstreamDomainAvail(&provider, domainName)
		if upErr != nil {
			s.log.Warnf("upstream domain check failed for %s: %v, falling back to WHOIS", domainName, upErr)
		} else {
			s.log.Infof("domain availability check via upstream: %s available=%v", domainName, available)
			return available, nil
		}
	}

	// 3. 没有上游或上游失败，使用WHOIS查询
	available, whoisErr := s.checkWhoisAvail(domainName)
	if whoisErr != nil {
		s.log.Errorf("WHOIS check failed for %s: %v", domainName, whoisErr)
		return false, fmt.Errorf("domain check failed: %w", whoisErr)
	}

	s.log.Infof("domain availability check via WHOIS: %s available=%v", domainName, available)
	return available, nil
}

// InitiateTransfer 发起域名转移
func (s *DomainService) InitiateTransfer(userID uint, domainName, eppCode string, price float64, registrarID *uint) (*model.DomainTransfer, error) {
	if domainName == "" {
		return nil, errors.New("域名不能为空")
	}
	if eppCode == "" {
		return nil, errors.New("EPP码不能为空")
	}

	var existing model.DomainTransfer
	err := s.db.Where("user_id = ? AND domain_name = ? AND status IN ?", userID, domainName, []string{"pending", "approved"}).First(&existing).Error
	if err == nil {
		return nil, errors.New("该域名已有进行中的转移申请")
	}

	transfer := &model.DomainTransfer{
		UserID:      userID,
		DomainName:  domainName,
		EPPCode:     eppCode,
		Status:      "pending",
		RegistrarID: registrarID,
		Price:       price,
	}

	if err := s.db.Create(transfer).Error; err != nil {
		s.log.Errorf("创建域名转移申请失败: %v", err)
		return nil, err
	}

	return transfer, nil
}

// GetTransfers 获取用户的域名转移列表
func (s *DomainService) GetTransfers(userID uint) ([]model.DomainTransfer, error) {
	var transfers []model.DomainTransfer
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&transfers).Error; err != nil {
		return nil, err
	}
	return transfers, nil
}

// AdminGetTransfers 管理员获取域名转移列表
func (s *DomainService) AdminGetTransfers(params DomainListParams) ([]model.DomainTransfer, int64, error) {
	var transfers []model.DomainTransfer
	var total int64

	query := s.db.Model(&model.DomainTransfer{})

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Keyword != "" {
		q := "%" + params.Keyword + "%"
		query = query.Where("domain_name LIKE ?", q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Offset(offset).Limit(params.PageSize).Preload("User").Order("created_at DESC").Find(&transfers).Error; err != nil {
		return nil, 0, err
	}

	return transfers, total, nil
}

// AdminUpdateTransfer 管理员更新域名转移状态
func (s *DomainService) AdminUpdateTransfer(id uint, status, remark string, adminID *uint) (*model.DomainTransfer, error) {
	var transfer model.DomainTransfer
	if err := s.db.First(&transfer, id).Error; err != nil {
		return nil, err
	}

	transfer.Status = status
	transfer.AdminID = adminID
	if remark != "" {
		transfer.Remark = remark
	}

	if status == "completed" || status == "rejected" || status == "failed" {
		now := time.Now()
		transfer.CompletedAt = &now
	}

	if err := s.db.Save(&transfer).Error; err != nil {
		s.log.Errorf("更新域名转移失败: %v", err)
		return nil, err
	}

	return &transfer, nil
}
