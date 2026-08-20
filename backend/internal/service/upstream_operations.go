package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/upstream"
)

// UpstreamOperations extends UpstreamService with server control operations.
type UpstreamOperations struct {
	*UpstreamService
}

// NewUpstreamOperations wraps an UpstreamService to expose server control methods.
func NewUpstreamOperations(svc *UpstreamService) *UpstreamOperations {
	return &UpstreamOperations{UpstreamService: svc}
}

// ==================== Generic upstream API callers ====================

// callUpstreamAPI dispatches an API call to the correct provider type.
func (s *UpstreamOperations) callUpstreamAPI(provider *model.UpstreamProvider, action string, params map[string]interface{}) (map[string]interface{}, error) {
	ptype := strings.ToLower(provider.Type)
	switch ptype {
	case "zjmf", "zjmfv3":
		return s.callZJMFAPI(provider, action, params)
	case "whmcs":
		return s.callWHMCSAPI(provider, action, params)
	case "v10":
		return s.callV10API(provider, action, params)
	default:
		return s.callCustomAPI(provider, action, params)
	}
}

// callZJMFAPI sends a signed request to a ZJMF-compatible API.
func (s *UpstreamOperations) callZJMFAPI(provider *model.UpstreamProvider, action string, params map[string]interface{}) (map[string]interface{}, error) {
	apiURL := strings.TrimRight(provider.APIURL, "/") + "/api.php"

	// Build param string map for signing.
	pmap := map[string]string{"action": action}
	for k, v := range params {
		pmap[k] = fmt.Sprintf("%v", v)
	}
	pmap["sign"] = zjmfSign(pmap, provider.APIKey)

	form := url.Values{}
	for k, v := range pmap {
		form.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("zjmf request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("zjmf read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zjmf http %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Result string          `json:"result"`
		Msg    string          `json:"msg"`
		Data   json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("zjmf parse response: %w", err)
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

// callWHMCSAPI sends a request to a WHMCS-compatible API.
func (s *UpstreamOperations) callWHMCSAPI(provider *model.UpstreamProvider, action string, params map[string]interface{}) (map[string]interface{}, error) {
	cfg := map[string]interface{}{}
	if provider.Config != nil {
		cfg = provider.Config
	}
	identifier := provider.APIKey
	if id, ok := cfg["identifier"].(string); ok && id != "" {
		identifier = id
	}
	secret, _ := cfg["secret"].(string)

	form := url.Values{}
	form.Set("action", action)
	form.Set("identifier", identifier)
	form.Set("secret", secret)
	form.Set("responsetype", "json")
	for k, v := range params {
		form.Set(k, fmt.Sprintf("%v", v))
	}

	apiURL := strings.TrimRight(provider.APIURL, "/") + "/includes/api.php"
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("whmcs request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("whmcs read body: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("whmcs http %d: %s", resp.StatusCode, string(body))
	}

	var raw json.RawMessage
	var apiResp struct {
		Result  string `json:"result"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("whmcs parse response: %w", err)
	}
	raw = json.RawMessage(body)

	if apiResp.Result != "success" {
		return nil, fmt.Errorf("whmcs error: %s", apiResp.Message)
	}

	var data interface{}
	json.Unmarshal(raw, &data)
	return map[string]interface{}{
		"result":  apiResp.Result,
		"message": apiResp.Message,
		"data":    data,
	}, nil
}

// callV10API sends a request to a V10-compatible API.
func (s *UpstreamOperations) callV10API(provider *model.UpstreamProvider, action string, params map[string]interface{}) (map[string]interface{}, error) {
	token := provider.APIKey
	if t, ok := provider.Config["token"].(string); ok && t != "" {
		token = t
	}

	payload, _ := json.Marshal(params)
	apiURL := strings.TrimRight(provider.APIURL, "/") + action

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(payload)))
	if err != nil {
		return nil, fmt.Errorf("v10 build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("v10 request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("v10 read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("v10 http %d: %s", resp.StatusCode, string(body))
	}

	var apiResp struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("v10 parse response: %w", err)
	}
	if apiResp.Code != 0 && apiResp.Code != 200 {
		return nil, fmt.Errorf("v10 error: %s", apiResp.Message)
	}

	var data interface{}
	json.Unmarshal(apiResp.Data, &data)
	return map[string]interface{}{
		"code":    apiResp.Code,
		"message": apiResp.Message,
		"data":    data,
	}, nil
}

// callCustomAPI sends a request to a custom upstream API.
func (s *UpstreamOperations) callCustomAPI(provider *model.UpstreamProvider, action string, params map[string]interface{}) (map[string]interface{}, error) {
	payload := map[string]interface{}{
		"action": action,
	}
	for k, v := range params {
		payload[k] = v
	}
	bodyBytes, _ := json.Marshal(payload)

	apiURL := strings.TrimRight(provider.APIURL, "/") + fmt.Sprintf("/api/host/%s", action)
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(bodyBytes)))
	if err != nil {
		return nil, fmt.Errorf("custom build request: %w", err)
	}
	if provider.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+provider.APIKey)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("custom request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("custom read body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("custom http %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp map[string]interface{}
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("custom parse response: %w", err)
	}
	if success, ok := apiResp["success"].(bool); ok && !success {
		msg, _ := apiResp["message"].(string)
		return nil, fmt.Errorf("custom error: %s", msg)
	}
	return apiResp, nil
}

// ==================== Upstream action mapping ====================

// upstreamAction returns the API action name for a given operation and provider type.
func upstreamAction(ptype, operation string) string {
	ptype = strings.ToLower(ptype)
	switch ptype {
	case "zjmf", "zjmfv3":
		return zjmfUpstreamAction(operation)
	case "whmcs":
		return whmcsUpstreamAction(operation)
	case "v10":
		return v10UpstreamAction(operation)
	default:
		return operation
	}
}

func zjmfUpstreamAction(op string) string {
	switch op {
	case "on", "boot":
		return "on"
	case "off", "shutdown":
		return "off"
	case "reboot":
		return "reboot"
	case "status", "getStatus":
		return "getStatus"
	case "vnc":
		return "vnc"
	case "reinstall":
		return "reinstall"
	case "getOs":
		return "getOs"
	case "crackPassword":
		return "CrackPassword"
	case "moduleClientButton":
		return "moduleClientButton"
	case "moduleAdminButton":
		return "moduleAdminButton"
	case "modulePowerStatus":
		return "modulePowerStatus"
	case "ipmiStatus":
		return "ipmiStatus"
	case "ipmiOn":
		return "ipmiOn"
	case "ipmiOff":
		return "ipmiOff"
	case "ipmiReboot":
		return "ipmiReboot"
	case "ipmiVnc":
		return "ipmiVnc"
	case "dcimClientStatus":
		return "dcimClientStatus"
	case "dcimClientOn":
		return "dcimClientOn"
	case "dcimClientOff":
		return "dcimClientOff"
	case "dcimClientReboot":
		return "dcimClientReboot"
	case "dcimClientVnc":
		return "dcimClientVnc"
	case "dcimClientReinstall":
		return "dcimClientReinstall"
	case "dcimClientCrackPass":
		return "dcimClientCrackPass"
	case "dcimClientCancelReinstall":
		return "dcimClientCancelReinstall"
	case "dcimClientReinstallStatus":
		return "dcimClientReinstallStatus"
	case "dcimClientGetOs":
		return "dcimClientGetOs"
	default:
		return op
	}
}

func whmcsUpstreamAction(op string) string {
	switch op {
	case "on", "boot":
		return "ModuleStart"
	case "off", "shutdown":
		return "ModuleStop"
	case "reboot":
		return "ModuleRestart"
	case "status", "getStatus":
		return "ModuleStatus"
	case "reinstall":
		return "ModuleReinstall"
	case "moduleClientButton":
		return "ModuleClientButton"
	case "moduleAdminButton":
		return "ModuleAdminButton"
	case "modulePowerStatus":
		return "ModulePowerStatus"
	default:
		return ""
	}
}

func v10UpstreamAction(op string) string {
	switch op {
	case "on", "boot":
		return "/api/server/start"
	case "off", "shutdown":
		return "/api/server/stop"
	case "reboot":
		return "/api/server/restart"
	case "status", "getStatus":
		return "/api/server/status"
	case "vnc":
		return "/api/server/vnc"
	case "reinstall":
		return "/api/server/reinstall"
	case "getOs":
		return "/api/server/os-list"
	case "crackPassword":
		return "/api/server/crack-password"
	default:
		return ""
	}
}

// ==================== Internal helpers ====================

// zjmfSign computes the ZJMF API signature: md5(sorted_params_values + api_key).
func zjmfSign(params map[string]string, apiKey string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(params[k])
	}
	sb.WriteString(apiKey)

	return fmt.Sprintf("%x", md5.Sum([]byte(sb.String())))
}

// getProviderAndHostID validates the provider and returns the remote host identifier.
// It looks for remote_host_id in host metadata, falling back to hostID as string.
func (s *UpstreamOperations) getProviderAndHostID(providerID, hostID uint) (*model.UpstreamProvider, string, error) {
	provider, err := s.GetProviderByID(providerID)
	if err != nil {
		return nil, "", fmt.Errorf("provider not found: %w", err)
	}
	if !provider.IsActive {
		return nil, "", fmt.Errorf("provider %d is not active", providerID)
	}

	// Try to get remote host ID from host metadata
	var host model.Host
	if err := s.db.First(&host, hostID).Error; err == nil {
		if len(host.Metadata) > 0 {
			var meta map[string]interface{}
			if err := json.Unmarshal(host.Metadata, &meta); err == nil {
				if rid, ok := meta["upstream_host_id"]; ok {
					return provider, fmt.Sprintf("%v", rid), nil
				}
				if rid, ok := meta["remote_host_id"]; ok {
					return provider, fmt.Sprintf("%v", rid), nil
				}
			}
		}
	}

	return provider, fmt.Sprintf("%d", hostID), nil
}

// recordOperation logs an upstream operation to the host_operations table.
func (s *UpstreamOperations) recordOperation(hostID, operatorID uint, action, params, result, errMsg string, success bool) {
	status := int8(2)
	if !success {
		status = 3
	}
	now := time.Now()
	op := model.HostOperation{
		HostID:     hostID,
		OperatorID: operatorID,
		Action:     action,
		Params:     params,
		Status:     status,
		Result:     result,
		ErrorMsg:   errMsg,
		StartedAt:  now,
		FinishedAt: &now,
	}
	s.db.Create(&op)
}

// callAndRecord wraps an upstream API call, recording the operation result.
func (s *UpstreamOperations) callAndRecord(providerID, hostID, operatorID uint, operation, actionOverride string, params map[string]interface{}) (map[string]interface{}, error) {
	provider, remoteHostID, err := s.getProviderAndHostID(providerID, hostID)
	if err != nil {
		s.recordOperation(hostID, operatorID, operation, "", "", err.Error(), false)
		return nil, err
	}

	callParams := map[string]interface{}{
		"vps_id": remoteHostID,
	}
	for k, v := range params {
		callParams[k] = v
	}

	apiAction := actionOverride
	if apiAction == "" {
		apiAction = upstreamAction(provider.Type, operation)
	}
	if apiAction == "" {
		err := fmt.Errorf("operation %s not supported for provider type %s", operation, provider.Type)
		s.recordOperation(hostID, operatorID, operation, "", "", err.Error(), false)
		return nil, err
	}

	result, err := s.callUpstreamAPI(provider, apiAction, callParams)
	if err != nil {
		s.recordOperation(hostID, operatorID, operation, fmt.Sprintf("%v", params), "", err.Error(), false)
		return nil, fmt.Errorf("upstream %s failed: %w", operation, err)
	}

	resultJSON, _ := json.Marshal(result)
	s.recordOperation(hostID, operatorID, operation, fmt.Sprintf("%v", params), string(resultJSON), "", true)
	return result, nil
}

// ==================== Power Operations ====================

// UpstreamBoot boots a server via the upstream provider API.
func (s *UpstreamOperations) UpstreamBoot(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "on", "", nil)
}

// UpstreamShutdown shuts down a server via the upstream provider API.
func (s *UpstreamOperations) UpstreamShutdown(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "off", "", nil)
}

// UpstreamReboot reboots a server via the upstream provider API.
func (s *UpstreamOperations) UpstreamReboot(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "reboot", "", nil)
}

// UpstreamGetStatus gets a server's status from the upstream provider API.
func (s *UpstreamOperations) UpstreamGetStatus(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "getStatus", "", nil)
}

// ==================== Console Access ====================

// UpstreamVNC gets the VNC console URL from the upstream provider.
func (s *UpstreamOperations) UpstreamVNC(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "vnc", "", nil)
}

// UpstreamKVM gets the KVM console URL from the upstream provider.
func (s *UpstreamOperations) UpstreamKVM(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "vnc", "kvm", nil)
}

// UpstreamIPMIStatus gets the IPMI status from the upstream provider.
func (s *UpstreamOperations) UpstreamIPMIStatus(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "ipmiStatus", "", nil)
}

// UpstreamIPMIOn powers on a server via IPMI through the upstream provider.
func (s *UpstreamOperations) UpstreamIPMIOn(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "ipmiOn", "", nil)
}

// UpstreamIPMIOff powers off a server via IPMI through the upstream provider.
func (s *UpstreamOperations) UpstreamIPMIOff(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "ipmiOff", "", nil)
}

// UpstreamIPMIReboot reboots a server via IPMI through the upstream provider.
func (s *UpstreamOperations) UpstreamIPMIReboot(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "ipmiReboot", "", nil)
}

// UpstreamIPMIVNC gets the IPMI VNC console from the upstream provider.
func (s *UpstreamOperations) UpstreamIPMIVNC(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "ipmiVnc", "", nil)
}

// ==================== Reinstall through Upstream ====================

// UpstreamReinstall initiates OS reinstall via the upstream provider.
func (s *UpstreamOperations) UpstreamReinstall(providerID, hostID, operatorID uint, osName, password string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"os": osName,
	}
	if password != "" {
		params["password"] = password
	}
	return s.callAndRecord(providerID, hostID, operatorID, "reinstall", "", params)
}

// UpstreamGetReinstallStatus checks the reinstall progress from the upstream.
func (s *UpstreamOperations) UpstreamGetReinstallStatus(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "reinstallStatus", "", nil)
}

// UpstreamCancelReinstall cancels an in-progress reinstall via the upstream.
func (s *UpstreamOperations) UpstreamCancelReinstall(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "cancelReinstall", "", nil)
}

// UpstreamGetOSList gets the available OS list from the upstream provider.
func (s *UpstreamOperations) UpstreamGetOSList(providerID uint, operatorID uint) (map[string]interface{}, error) {
	provider, err := s.GetProviderByID(providerID)
	if err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}
	if !provider.IsActive {
		return nil, fmt.Errorf("provider %d is not active", providerID)
	}

	apiAction := upstreamAction(provider.Type, "getOs")
	if apiAction == "" {
		return nil, fmt.Errorf("getOs not supported for provider type %s", provider.Type)
	}

	result, err := s.callUpstreamAPI(provider, apiAction, nil)
	if err != nil {
		return nil, fmt.Errorf("upstream getOs failed: %w", err)
	}
	return result, nil
}

// UpstreamCrackPassword resets a server's password via the upstream provider.
func (s *UpstreamOperations) UpstreamCrackPassword(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "crackPassword", "", nil)
}

// ==================== DCIM Client Operations ====================

// DcimClientStatus gets the DCIM client status from the upstream.
func (s *UpstreamOperations) DcimClientStatus(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "dcimClientStatus", "", nil)
}

// DcimClientOn powers on via DCIM client through the upstream.
func (s *UpstreamOperations) DcimClientOn(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "dcimClientOn", "", nil)
}

// DcimClientOff powers off via DCIM client through the upstream.
func (s *UpstreamOperations) DcimClientOff(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "dcimClientOff", "", nil)
}

// DcimClientReboot reboots via DCIM client through the upstream.
func (s *UpstreamOperations) DcimClientReboot(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "dcimClientReboot", "", nil)
}

// DcimClientVNC gets the DCIM client VNC console from the upstream.
func (s *UpstreamOperations) DcimClientVNC(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "dcimClientVnc", "", nil)
}

// DcimClientReinstall reinstalls via DCIM client through the upstream.
func (s *UpstreamOperations) DcimClientReinstall(providerID, hostID, operatorID uint, osName string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"os": osName,
	}
	return s.callAndRecord(providerID, hostID, operatorID, "dcimClientReinstall", "", params)
}

// DcimClientCrackPass resets password via DCIM client through the upstream.
func (s *UpstreamOperations) DcimClientCrackPass(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "dcimClientCrackPass", "", nil)
}

// DcimClientCancelReinstall cancels reinstall via DCIM client through the upstream.
func (s *UpstreamOperations) DcimClientCancelReinstall(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "dcimClientCancelReinstall", "", nil)
}

// DcimClientReinstallStatus gets the DCIM client reinstall status from the upstream.
func (s *UpstreamOperations) DcimClientReinstallStatus(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "dcimClientReinstallStatus", "", nil)
}

// DcimClientGetOS gets the OS list from the DCIM client through the upstream.
func (s *UpstreamOperations) DcimClientGetOS(providerID uint, operatorID uint) (map[string]interface{}, error) {
	provider, err := s.GetProviderByID(providerID)
	if err != nil {
		return nil, fmt.Errorf("provider not found: %w", err)
	}
	if !provider.IsActive {
		return nil, fmt.Errorf("provider %d is not active", providerID)
	}

	apiAction := upstreamAction(provider.Type, "dcimClientGetOs")
	if apiAction == "" {
		return nil, fmt.Errorf("dcimClientGetOs not supported for provider type %s", provider.Type)
	}

	result, err := s.callUpstreamAPI(provider, apiAction, nil)
	if err != nil {
		return nil, fmt.Errorf("upstream dcimClientGetOs failed: %w", err)
	}
	return result, nil
}

// ==================== Module Buttons ====================

// ModuleClientButton executes a client-side module button via the upstream.
func (s *UpstreamOperations) ModuleClientButton(providerID, hostID, operatorID uint, button string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"button": button,
	}
	return s.callAndRecord(providerID, hostID, operatorID, "moduleClientButton", "", params)
}

// ModuleAdminButton executes an admin-side module button via the upstream.
func (s *UpstreamOperations) ModuleAdminButton(providerID, hostID, operatorID uint, button string) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"button": button,
	}
	return s.callAndRecord(providerID, hostID, operatorID, "moduleAdminButton", "", params)
}

// ModulePowerStatus gets the power status from the upstream module.
func (s *UpstreamOperations) ModulePowerStatus(providerID, hostID, operatorID uint) (map[string]interface{}, error) {
	return s.callAndRecord(providerID, hostID, operatorID, "modulePowerStatus", "", nil)
}

// ==================== Upstream Connection Test (reuse) ====================

// GetUpstreamClient returns the upstream Client interface for the given provider.
func (s *UpstreamOperations) GetUpstreamClient(providerID uint) (upstream.Client, error) {
	provider, err := s.GetProviderByID(providerID)
	if err != nil {
		return nil, err
	}
	return upstream.NewClient(provider)
}
