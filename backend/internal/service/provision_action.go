package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"anchorfinance/internal/model"
)

// ProvisionOrder provisions all products in an order after payment.
func (s *ProvisionService) ProvisionOrder(orderID uint) error {
	var order model.Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		return fmt.Errorf("order not found: %w", err)
	}

	var userProducts []model.UserProduct
	if err := s.db.Where("order_id = ? AND provisioning_status = 0", orderID).Find(&userProducts).Error; err != nil {
		return fmt.Errorf("query user products: %w", err)
	}

	for _, up := range userProducts {
		if err := s.ProvisionProduct(up.ID, "create"); err != nil {
			s.log.Errorf("provision product %d for order %d: %v", up.ID, orderID, err)
		}
	}

	s.log.Infof("provisioning completed for order %d (%d items)", orderID, len(userProducts))
	return nil
}

// ProvisionProduct provisions a single product with the given action.
// Actions: create, suspend, terminate, unsuspend, rebuild, renew
func (s *ProvisionService) ProvisionProduct(userProductID uint, action string) error {
	var up model.UserProduct
	if err := s.db.Preload("Product").First(&up, userProductID).Error; err != nil {
		return fmt.Errorf("user product not found: %w", err)
	}

	// Find matching provision module for this product
	module, err := s.findModuleForProduct(up.ProductID)
	if err != nil {
		s.updateProvisionStatus(userProductID, 3, fmt.Sprintf("no provision module: %v", err))
		return fmt.Errorf("find provision module: %w", err)
	}

	// Mark as processing
	s.updateProvisionStatus(userProductID, 1, "")

	startTime := time.Now()
	logEntry := &model.ProvisionLog{
		ModuleID: module.ID,
		Action:   action,
		Status:   1,
	}

	reqData, _ := json.Marshal(map[string]interface{}{
		"user_product_id": userProductID,
		"action":          action,
		"product_id":      up.ProductID,
		"user_id":         up.UserID,
	})
	logEntry.Request = reqData

	var provisionErr error

	switch module.Type {
	case "upstream":
		provisionErr = s.callUpstreamAPI(module, &up, action)
	case "local":
		provisionErr = s.handleLocalProvision(module, &up, action)
	default:
		provisionErr = fmt.Errorf("unsupported module type: %s", module.Type)
	}

	duration := int(time.Since(startTime).Milliseconds())
	logEntry.Duration = duration

	if provisionErr != nil {
		logEntry.Status = 3
		logEntry.Error = provisionErr.Error()
		s.updateProvisionStatus(userProductID, 3, provisionErr.Error())
		s.IncrProvisionCount(module.ID, false)
	} else {
		logEntry.Status = 2
		s.updateProvisionStatus(userProductID, 2, "")
		s.IncrProvisionCount(module.ID, true)
		s.applyActionToProduct(&up, action)
	}

	respData, _ := json.Marshal(map[string]interface{}{
		"success":  provisionErr == nil,
		"duration": duration,
	})
	logEntry.Response = respData

	s.db.Create(logEntry)

	if provisionErr != nil {
		return fmt.Errorf("provision %s failed: %w", action, provisionErr)
	}

	s.log.Infof("provision %s success: user_product=%d module=%s", action, userProductID, module.Slug)
	return nil
}

// SuspendService suspends a user product.
func (s *ProvisionService) SuspendService(userProductID uint, reason string) error {
	var up model.UserProduct
	if err := s.db.First(&up, userProductID).Error; err != nil {
		return fmt.Errorf("user product not found: %w", err)
	}
	if up.Status != 1 {
		return fmt.Errorf("product is not active (status=%d)", up.Status)
	}

	if err := s.ProvisionProduct(userProductID, "suspend"); err != nil {
		return err
	}

	s.db.Model(&up).Updates(map[string]interface{}{
		"status":         2,
		"suspend_reason": reason,
	})
	return nil
}

// TerminateService terminates a user product.
func (s *ProvisionService) TerminateService(userProductID uint, reason string) error {
	var up model.UserProduct
	if err := s.db.First(&up, userProductID).Error; err != nil {
		return fmt.Errorf("user product not found: %w", err)
	}
	if up.Status == 4 {
		return fmt.Errorf("product is already terminated")
	}

	if err := s.ProvisionProduct(userProductID, "terminate"); err != nil {
		return err
	}

	now := time.Now()
	s.db.Model(&up).Updates(map[string]interface{}{
		"status":           4,
		"termination_date": &now,
	})
	return nil
}

// UnsuspendService reactivates a suspended user product.
func (s *ProvisionService) UnsuspendService(userProductID uint) error {
	var up model.UserProduct
	if err := s.db.First(&up, userProductID).Error; err != nil {
		return fmt.Errorf("user product not found: %w", err)
	}
	if up.Status != 2 {
		return fmt.Errorf("product is not suspended (status=%d)", up.Status)
	}

	if err := s.ProvisionProduct(userProductID, "unsuspend"); err != nil {
		return err
	}

	s.db.Model(&up).Updates(map[string]interface{}{
		"status":         1,
		"suspend_reason": "",
	})
	return nil
}

// findModuleForProduct finds the active provision module that supports the given product.
func (s *ProvisionService) findModuleForProduct(productID uint) (*model.ProvisionModule, error) {
	var product model.Product
	if err := s.db.First(&product, productID).Error; err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Try to find a module that supports this product type
	var modules []model.ProvisionModule
	if err := s.db.Where("active = true AND type = ?", product.Type).
		Order("priority DESC, weight DESC").
		Find(&modules).Error; err != nil {
		return nil, err
	}

	if len(modules) > 0 {
		return &modules[0], nil
	}

	// Fallback: find any active module
	var fallback model.ProvisionModule
	if err := s.db.Where("active = true").Order("priority DESC").First(&fallback).Error; err != nil {
		return nil, fmt.Errorf("no active provision module found")
	}
	return &fallback, nil
}

// callUpstreamAPI calls an upstream provider's API for provisioning.
func (s *ProvisionService) callUpstreamAPI(module *model.ProvisionModule, up *model.UserProduct, action string) error {
	if module.ServerURL == "" {
		return fmt.Errorf("module server_url is empty")
	}

	slug := strings.ToLower(module.Slug)
	switch {
	case strings.Contains(slug, "zjmf"):
		return s.callZJMFProvision(module, up, action)
	case strings.Contains(slug, "whmcs"):
		return s.callWHMCSProvision(module, up, action)
	case strings.Contains(slug, "v10"):
		return s.callV10Provision(module, up, action)
	default:
		return s.callCustomProvision(module, up, action)
	}
}

// callZJMFProvision makes a signed provisioning request to a ZJMF-compatible panel.
func (s *ProvisionService) callZJMFProvision(module *model.ProvisionModule, up *model.UserProduct, action string) error {
	apiAction := mapZJMFAction(action)
	params := map[string]string{
		"action":     apiAction,
		"product_id": fmt.Sprintf("%d", up.ProductID),
		"user_id":    fmt.Sprintf("%d", up.UserID),
	}
	if up.Domain != "" {
		params["domain"] = up.Domain
	}
	if up.Username != "" {
		params["username"] = up.Username
	}
	if up.Hostname != "" {
		params["hostname"] = up.Hostname
	}
	if action == "suspend" && up.SuspendReason != "" {
		params["reason"] = up.SuspendReason
	}
	if action == "create" && up.ConfigOptions != nil {
		cfgJSON, _ := json.Marshal(up.ConfigOptions)
		params["configoptions"] = string(cfgJSON)
	}

	// ZJMF API signature: md5(sorted param values + api_key)
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(params[k])
	}
	sb.WriteString(module.APIKey)
	params["sign"] = fmt.Sprintf("%x", md5.Sum([]byte(sb.String())))

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	baseURL := strings.TrimRight(module.ServerURL, "/")
	client := &http.Client{Timeout: time.Duration(module.Timeout) * time.Second}
	resp, err := client.Post(baseURL+"/api.php", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("zjmf provision request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("zjmf read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("zjmf http %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result string          `json:"result"`
		Msg    string          `json:"msg"`
		Data   json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("zjmf parse response: %w", err)
	}
	if apiResp.Result != "success" {
		return fmt.Errorf("zjmf error: %s", apiResp.Msg)
	}

	s.log.Infof("zjmf provision success: module=%s action=%s", module.Slug, action)
	return nil
}

// callWHMCSProvision makes a provisioning request to a WHMCS-compatible API.
func (s *ProvisionService) callWHMCSProvision(module *model.ProvisionModule, up *model.UserProduct, action string) error {
	whmcsAction := mapWHMCSAction(action)
	if whmcsAction == "" {
		return fmt.Errorf("whmcs does not support action: %s", action)
	}

	form := url.Values{}
	form.Set("action", whmcsAction)
	form.Set("identifier", module.APIKey)
	form.Set("secret", module.APISecret)
	form.Set("responsetype", "json")
	form.Set("serviceid", fmt.Sprintf("%d", up.ID))
	if up.Domain != "" {
		form.Set("domain", up.Domain)
	}
	if action == "suspend" && up.SuspendReason != "" {
		form.Set("suspendreason", up.SuspendReason)
	}
	if action == "create" && up.ConfigOptions != nil {
		cfgJSON, _ := json.Marshal(up.ConfigOptions)
		form.Set("customfields", string(cfgJSON))
	}

	baseURL := strings.TrimRight(module.ServerURL, "/")
	client := &http.Client{Timeout: time.Duration(module.Timeout) * time.Second}
	resp, err := client.Post(baseURL+"/includes/api.php", "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("whmcs provision request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("whmcs read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("whmcs http %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result  string `json:"result"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return fmt.Errorf("whmcs parse response: %w", err)
	}
	if apiResp.Result != "success" {
		return fmt.Errorf("whmcs error: %s", apiResp.Message)
	}

	s.log.Infof("whmcs provision success: module=%s action=%s", module.Slug, action)
	return nil
}

// callV10Provision makes a provisioning request to a V10 panel API.
func (s *ProvisionService) callV10Provision(module *model.ProvisionModule, up *model.UserProduct, action string) error {
	apiPath := mapV10Path(action)
	if apiPath == "" {
		return fmt.Errorf("v10 does not support action: %s", action)
	}

	payload := map[string]interface{}{
		"product_id": up.ProductID,
		"user_id":    up.UserID,
	}
	if up.Domain != "" {
		payload["domain"] = up.Domain
	}
	if up.Username != "" {
		payload["username"] = up.Username
	}
	if action == "suspend" && up.SuspendReason != "" {
		payload["reason"] = up.SuspendReason
	}

	bodyBytes, _ := json.Marshal(payload)

	baseURL := strings.TrimRight(module.ServerURL, "/")
	client := &http.Client{Timeout: time.Duration(module.Timeout) * time.Second}
	req, err := http.NewRequest("POST", baseURL+apiPath, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return fmt.Errorf("v10 build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+module.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("v10 provision request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("v10 read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("v10 http %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("v10 parse response: %w", err)
	}
	if apiResp.Code != 0 && apiResp.Code != 200 {
		return fmt.Errorf("v10 error: %s", apiResp.Message)
	}

	s.log.Infof("v10 provision success: module=%s action=%s", module.Slug, action)
	return nil
}

// callCustomProvision makes a provisioning request to a custom upstream API.
func (s *ProvisionService) callCustomProvision(module *model.ProvisionModule, up *model.UserProduct, action string) error {
	payload := map[string]interface{}{
		"action":     action,
		"product_id": up.ProductID,
		"user_id":    up.UserID,
		"domain":     up.Domain,
		"username":   up.Username,
		"hostname":   up.Hostname,
	}
	if action == "create" && up.ConfigOptions != nil {
		payload["config"] = up.ConfigOptions
	}
	if action == "suspend" && up.SuspendReason != "" {
		payload["reason"] = up.SuspendReason
	}

	bodyBytes, _ := json.Marshal(payload)

	baseURL := strings.TrimRight(module.ServerURL, "/")
	endpoint := fmt.Sprintf("/api/provision/%s", action)

	client := &http.Client{Timeout: time.Duration(module.Timeout) * time.Second}
	req, err := http.NewRequest("POST", baseURL+endpoint, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return fmt.Errorf("custom build request: %w", err)
	}
	if module.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+module.APIKey)
	}
	if module.APISecret != "" {
		req.Header.Set("X-API-Secret", module.APISecret)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("custom provision request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("custom read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("custom http %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp map[string]interface{}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("custom parse response: %w", err)
	}
	if success, ok := apiResp["success"].(bool); ok && !success {
		msg, _ := apiResp["message"].(string)
		if msg == "" {
			msg = "upstream returned failure"
		}
		return fmt.Errorf("custom error: %s", msg)
	}

	s.log.Infof("custom provision success: module=%s action=%s", module.Slug, action)
	return nil
}

// mapZJMFAction maps generic provision actions to ZJMF API action names.
func mapZJMFAction(action string) string {
	switch action {
	case "create":
		return "createproduct"
	case "suspend":
		return "suspendproduct"
	case "unsuspend":
		return "unsuspendproduct"
	case "terminate":
		return "terminateproduct"
	case "renew":
		return "renewproduct"
	default:
		return action
	}
}

// mapWHMCSAction maps generic provision actions to WHMCS API action names.
func mapWHMCSAction(action string) string {
	switch action {
	case "create":
		return "ModuleCreate"
	case "suspend":
		return "ModuleSuspend"
	case "unsuspend":
		return "ModuleUnsuspend"
	case "terminate":
		return "ModuleTerminate"
	case "renew":
		return "ModuleRenew"
	default:
		return ""
	}
}

// mapV10Path maps generic provision actions to V10 API paths.
func mapV10Path(action string) string {
	switch action {
	case "create":
		return "/api/service/create"
	case "suspend":
		return "/api/service/suspend"
	case "unsuspend":
		return "/api/service/unsuspend"
	case "terminate":
		return "/api/service/terminate"
	case "renew":
		return "/api/service/renew"
	default:
		return ""
	}
}

// handleLocalProvision handles provisioning locally (e.g., for hosting products).
func (s *ProvisionService) handleLocalProvision(module *model.ProvisionModule, up *model.UserProduct, action string) error {
	switch action {
	case "create":
		return s.localCreate(module, up)
	case "suspend":
		return s.localSuspend(up)
	case "unsuspend":
		return s.localUnsuspend(up)
	case "terminate":
		return s.localTerminate(up)
	case "renew":
		// Extract renew_days from module config, default to billing cycle
		days := 30
		if len(module.Config) > 0 {
			var configMap map[string]interface{}
			if err := json.Unmarshal(module.Config, &configMap); err == nil {
				if d, ok := configMap["renew_days"].(float64); ok && d > 0 {
					days = int(d)
				}
			}
		}
		return s.localRenew(up, days)
	default:
		return fmt.Errorf("unsupported local action: %s", action)
	}
}

func (s *ProvisionService) localCreate(module *model.ProvisionModule, up *model.UserProduct) error {
	s.log.Infof("local create: product=%d user=%d", up.ProductID, up.UserID)

	// Load the product to determine service type
	var product model.Product
	if err := s.db.First(&product, up.ProductID).Error; err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Generate unique credentials
	username := generateServiceUsername(up.UserID, up.ID)
	password := generateSecurePassword(16)

	// Build default hostname if not set
	hostname := up.Hostname
	if hostname == "" {
		hostname = fmt.Sprintf("svc-%d-%d.local", up.UserID, up.ID)
	}

	// Determine IP allocation from module config if available
	ip := up.IP
	if ip == "" && module.ServerIP != "" {
		ip = module.ServerIP
	}

	// Calculate expiry based on billing cycle
	var expiryDate time.Time
	switch product.BillingCycle {
	case "monthly":
		expiryDate = time.Now().AddDate(0, 1, 0)
	case "quarterly":
		expiryDate = time.Now().AddDate(0, 3, 0)
	case "semi-annually":
		expiryDate = time.Now().AddDate(0, 6, 0)
	case "annually":
		expiryDate = time.Now().AddDate(1, 0, 0)
	case "biennially":
		expiryDate = time.Now().AddDate(2, 0, 0)
	case "triennially":
		expiryDate = time.Now().AddDate(3, 0, 0)
	default:
		expiryDate = time.Now().AddDate(0, 1, 0) // default monthly
	}

	now := time.Now()
	nextDue := expiryDate

	updates := map[string]interface{}{
		"username":          username,
		"password":          password,
		"hostname":          hostname,
		"ip":                ip,
		"status":            1, // active
		"provisioning_status": 2, // success
		"registration_date": &now,
		"next_due_date":     &nextDue,
	}

	if err := s.db.Model(up).Updates(updates).Error; err != nil {
		return fmt.Errorf("update user_product: %w", err)
	}

	s.log.Infof("local create success: user_product=%d username=%s expiry=%s", up.ID, username, expiryDate.Format("2006-01-02"))
	return nil
}

func (s *ProvisionService) localSuspend(up *model.UserProduct) error {
	s.log.Infof("local suspend: user_product=%d", up.ID)

	now := time.Now()
	if err := s.db.Model(up).Updates(map[string]interface{}{
		"status":         2, // suspended
		"suspend_reason": up.SuspendReason,
		"suspended_at":   &now,
	}).Error; err != nil {
		return fmt.Errorf("suspend user_product: %w", err)
	}

	s.log.Infof("local suspend success: user_product=%d", up.ID)
	return nil
}

func (s *ProvisionService) localUnsuspend(up *model.UserProduct) error {
	s.log.Infof("local unsuspend: user_product=%d", up.ID)

	if err := s.db.Model(up).Updates(map[string]interface{}{
		"status":         1, // active
		"suspend_reason": "",
		"suspended_at":   nil,
	}).Error; err != nil {
		return fmt.Errorf("unsuspend user_product: %w", err)
	}

	s.log.Infof("local unsuspend success: user_product=%d", up.ID)
	return nil
}

func (s *ProvisionService) localTerminate(up *model.UserProduct) error {
	s.log.Infof("local terminate: user_product=%d", up.ID)

	now := time.Now()
	if err := s.db.Model(up).Updates(map[string]interface{}{
		"status":            4, // terminated
		"termination_date":  &now,
		"provisioning_status": 0, // reset
	}).Error; err != nil {
		return fmt.Errorf("terminate user_product: %w", err)
	}

	s.log.Infof("local terminate success: user_product=%d", up.ID)
	return nil
}

func (s *ProvisionService) localRenew(up *model.UserProduct, days int) error {
	s.log.Infof("local renew: user_product=%d days=%d", up.ID, days)

	if days <= 0 {
		days = 30 // default to monthly
	}

	// Extend from current expiry or now, whichever is later
	baseTime := time.Now()
	if up.NextDueDate != nil && up.NextDueDate.After(baseTime) {
		baseTime = *up.NextDueDate
	}
	newExpiry := baseTime.AddDate(0, 0, days)
	nextDue := newExpiry

	if err := s.db.Model(up).Updates(map[string]interface{}{
		"next_due_date": &nextDue,
		"status":        1, // ensure active
	}).Error; err != nil {
		return fmt.Errorf("renew user_product: %w", err)
	}

	s.log.Infof("local renew success: user_product=%d new_expiry=%s", up.ID, newExpiry.Format("2006-01-02"))
	return nil
}

// generateServiceUsername creates a unique service username from user ID and product ID.
func generateServiceUsername(userID, productID uint) string {
	return fmt.Sprintf("u%dp%d", userID, productID)
}

// generateSecurePassword generates a random alphanumeric password of the given length.
func generateSecurePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// updateProvisionStatus updates the provisioning_status of a user product.
func (s *ProvisionService) updateProvisionStatus(userProductID uint, status int16, errMsg string) {
	updates := map[string]interface{}{
		"provisioning_status": status,
	}
	if errMsg != "" {
		s.db.Model(&model.UserProduct{}).Where("id = ?", userProductID).Updates(updates)
		s.log.Warnf("provision status %d for product %d: %s", status, userProductID, errMsg)
	} else {
		s.db.Model(&model.UserProduct{}).Where("id = ?", userProductID).Updates(updates)
	}
}

// applyActionToProduct updates user product fields after successful provisioning.
func (s *ProvisionService) applyActionToProduct(up *model.UserProduct, action string) {
	switch action {
	case "create":
		now := time.Now()
		s.db.Model(up).Update("registration_date", &now)
	case "suspend":
		s.db.Model(up).Update("status", 2)
	case "unsuspend":
		s.db.Model(up).Updates(map[string]interface{}{
			"status":         1,
			"suspend_reason": "",
		})
	case "terminate":
		now := time.Now()
		s.db.Model(up).Updates(map[string]interface{}{
			"status":           4,
			"termination_date": &now,
		})
	}
}
