package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ============================================================
// ManualDriver 手动管理驱动（不对接API）
// ============================================================

type ManualDriver struct {
	supplierID uint
}

func NewManualDriver(supplierID uint) *ManualDriver {
	return &ManualDriver{supplierID: supplierID}
}

func (d *ManualDriver) Key() string              { return "manual" }
func (d *ManualDriver) Name() string              { return "手动管理" }
func (d *ManualDriver) Capabilities() []string    { return nil }
func (d *ManualDriver) FetchProducts() ([]RemoteProduct, error) {
	return nil, fmt.Errorf("手动管理不支持自动拉取商品")
}
func (d *ManualDriver) FetchProductGroups() ([]RemoteGroup, error) {
	return nil, fmt.Errorf("手动管理不支持自动拉取分组")
}
func (d *ManualDriver) GetProductStructure() (*ProductStructure, error) {
	return nil, fmt.Errorf("手动管理不支持获取商品结构")
}
func (d *ManualDriver) SyncStatus(serviceID string) (*StatusResult, error) {
	return nil, fmt.Errorf("手动管理不支持状态同步")
}
func (d *ManualDriver) CreateService(params CreateServiceParams) (*ServiceResult, error) {
	return nil, fmt.Errorf("手动管理不支持自动开通")
}
func (d *ManualDriver) SuspendService(serviceID string) error   { return fmt.Errorf("手动管理不支持自动暂停") }
func (d *ManualDriver) UnsuspendService(serviceID string) error { return fmt.Errorf("手动管理不支持自动取消暂停") }
func (d *ManualDriver) TerminateService(serviceID string) error { return fmt.Errorf("手动管理不支持自动终止") }
func (d *ManualDriver) RenewService(serviceID string, cycle string) error {
	return fmt.Errorf("手动管理不支持自动续费")
}

// ============================================================
// ZjmfDriver zjmf API驱动（/v1/端点）
// ============================================================

type ZjmfDriver struct {
	supplierID uint
	apiURL     string // https://example.com
	apiKey     string // API密钥（账号或密钥）
	apiSecret  string // API密码
	client     *http.Client
	jwt        string // 缓存的JWT
}

func NewZjmfDriver(supplierID uint, apiURL, apiKey, apiSecret string) *ZjmfDriver {
	return &ZjmfDriver{
		supplierID: supplierID,
		apiURL:     apiURL,
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		client:     &http.Client{Timeout: 30 * time.Second},
	}
}

func (d *ZjmfDriver) Key() string           { return "zjmf" }
func (d *ZjmfDriver) Name() string           { return "zjmf接口" }
func (d *ZjmfDriver) Capabilities() []string { return []string{"provisioning", "renewal", "status_sync", "product_sync"} }

// login 获取JWT
func (d *ZjmfDriver) login() (string, error) {
	if d.jwt != "" {
		return d.jwt, nil
	}

	payload := map[string]string{
		"account":  d.apiKey,
		"password": d.apiSecret,
	}
	body, _ := json.Marshal(payload)

	resp, err := d.client.Post(d.apiURL+"/v1/login_api", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			JWT string `json:"jwt"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析登录响应失败")
	}
	if result.Code != 200 {
		return "", fmt.Errorf("登录失败: %s", result.Msg)
	}

	d.jwt = result.Data.JWT
	return d.jwt, nil
}

// request 发送认证请求（含JWT过期自动刷新）
func (d *ZjmfDriver) request(method, path string, params map[string]interface{}) (map[string]interface{}, error) {
	return d.doRequest(method, path, params, false)
}

func (d *ZjmfDriver) doRequest(method, path string, params map[string]interface{}, isRetry bool) (map[string]interface{}, error) {
	jwt, err := d.login()
	if err != nil {
		return nil, err
	}

	var req *http.Request
	if method == "GET" {
		req, _ = http.NewRequest(method, d.apiURL+path, nil)
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, fmt.Sprintf("%v", v))
		}
		req.URL.RawQuery = q.Encode()
	} else {
		body, _ := json.Marshal(params)
		req, _ = http.NewRequest(method, d.apiURL+path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Authorization", "Bearer "+jwt)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// JWT过期自动刷新（创欧shouldRetryWithFreshJwt机制）
	if resp.StatusCode == 401 && !isRetry {
		d.jwt = ""
		return d.doRequest(method, path, params, true)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	return result, nil
}

func (d *ZjmfDriver) FetchProducts() ([]RemoteProduct, error) {
	result, err := d.request("GET", "/v1/products", nil)
	if err != nil {
		return nil, err
	}

	var products []RemoteProduct
	if data, ok := result["data"].(map[string]interface{}); ok {
		if items, ok := data["product"].([]interface{}); ok {
			for _, item := range items {
				if m, ok := item.(map[string]interface{}); ok {
					p := RemoteProduct{
						ID:   fmt.Sprintf("%v", m["id"]),
						Name: fmt.Sprintf("%v", m["name"]),
					}
					if price, ok := m["price"].(float64); ok {
						p.Price = price
					}
					if gid, ok := m["gid"].(float64); ok {
						p.GroupID = fmt.Sprintf("%.0f", gid)
					}
					if gname, ok := m["group_name"].(string); ok {
						p.GroupName = gname
					}
					if stock, ok := m["stock"].(float64); ok {
						p.Stock = int(stock)
					}
					products = append(products, p)
				}
			}
		}
	}
	return products, nil
}

func (d *ZjmfDriver) FetchProductGroups() ([]RemoteGroup, error) {
	result, err := d.request("GET", "/v1/hosts/cates", nil)
	if err != nil {
		return nil, err
	}

	var groups []RemoteGroup
	if data, ok := result["data"].(map[string]interface{}); ok {
		if items, ok := data["cate"].([]interface{}); ok {
			for _, item := range items {
				if m, ok := item.(map[string]interface{}); ok {
					g := RemoteGroup{
						ID:   fmt.Sprintf("%v", m["id"]),
						Name: fmt.Sprintf("%v", m["name"]),
					}
					groups = append(groups, g)
				}
			}
		}
	}
	return groups, nil
}

func (d *ZjmfDriver) GetProductStructure() (*ProductStructure, error) {
	groups, err := d.FetchProductGroups()
	if err != nil {
		return nil, err
	}
	return &ProductStructure{Groups: groups}, nil
}

func (d *ZjmfDriver) SyncStatus(serviceID string) (*StatusResult, error) {
	result, err := d.request("GET", "/v1/hosts/"+serviceID+"/module/status", map[string]interface{}{"type": "host"})
	if err != nil {
		return nil, err
	}

	sr := &StatusResult{Status: "unknown"}
	if data, ok := result["data"].(map[string]interface{}); ok {
		if s, ok := data["status"].(string); ok {
			sr.Status = s
		}
		if ip, ok := data["ip"].(string); ok {
			sr.IPAddress = ip
		}
	}
	return sr, nil
}

func (d *ZjmfDriver) CreateService(params CreateServiceParams) (*ServiceResult, error) {
	// zjmf通过购物车下单流程创建服务
	// 1. 加入购物车
	cartData := map[string]interface{}{
		"product_id": params.ProductID,
		"qty":        1,
		"cycle":      params.Cycle,
	}
	_, err := d.request("POST", "/cart/add", cartData)
	if err != nil {
		return nil, fmt.Errorf("加入购物车失败: %w", err)
	}

	// 2. 结算
	_, err = d.request("POST", "/cart/checkout", nil)
	if err != nil {
		return nil, fmt.Errorf("结算失败: %w", err)
	}

	return &ServiceResult{}, nil
}

func (d *ZjmfDriver) SuspendService(serviceID string) error {
	_, err := d.request("POST", "/v1/hosts/"+serviceID+"/module/suspend", nil)
	return err
}

func (d *ZjmfDriver) UnsuspendService(serviceID string) error {
	_, err := d.request("POST", "/v1/hosts/"+serviceID+"/module/unsuspend", nil)
	return err
}

func (d *ZjmfDriver) TerminateService(serviceID string) error {
	_, err := d.request("POST", "/v1/hosts/"+serviceID+"/module/terminate", nil)
	return err
}

func (d *ZjmfDriver) RenewService(serviceID string, cycle string) error {
	_, err := d.request("POST", "/v1/hosts/"+serviceID+"/renew", map[string]interface{}{"billingcycle": cycle})
	return err
}

// ============================================================
// V10Driver v10 API驱动（/console/v1/端点）
// ============================================================

type V10Driver struct {
	supplierID uint
	apiURL     string
	apiKey     string
	apiSecret  string
	client     *http.Client
	jwt        string
}

func NewV10Driver(supplierID uint, apiURL, apiKey, apiSecret string) *V10Driver {
	return &V10Driver{
		supplierID: supplierID,
		apiURL:     apiURL,
		apiKey:     apiKey,
		apiSecret:  apiSecret,
		client:     &http.Client{Timeout: 30 * time.Second},
	}
}

func (d *V10Driver) Key() string           { return "v10" }
func (d *V10Driver) Name() string           { return "v10接口" }
func (d *V10Driver) Capabilities() []string { return []string{"provisioning", "renewal", "status_sync", "product_sync"} }

func (d *V10Driver) login() (string, error) {
	if d.jwt != "" {
		return d.jwt, nil
	}

	payload := map[string]string{
		"account":  d.apiKey,
		"password": d.apiSecret,
	}
	body, _ := json.Marshal(payload)

	resp, err := d.client.Post(d.apiURL+"/v1/login_api", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("登录请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			JWT string `json:"jwt"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("解析登录响应失败")
	}
	if result.Code != 200 {
		return "", fmt.Errorf("登录失败: %s", result.Msg)
	}

	d.jwt = result.Data.JWT
	return d.jwt, nil
}

// request 发送认证请求（含JWT过期自动刷新）
func (d *V10Driver) request(method, path string, params map[string]interface{}) (map[string]interface{}, error) {
	return d.doRequest(method, path, params, false)
}

func (d *V10Driver) doRequest(method, path string, params map[string]interface{}, isRetry bool) (map[string]interface{}, error) {
	jwt, err := d.login()
	if err != nil {
		return nil, err
	}

	var req *http.Request
	if method == "GET" {
		req, _ = http.NewRequest(method, d.apiURL+path, nil)
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, fmt.Sprintf("%v", v))
		}
		req.URL.RawQuery = q.Encode()
	} else {
		body, _ := json.Marshal(params)
		req, _ = http.NewRequest(method, d.apiURL+path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Authorization", "Bearer "+jwt)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// JWT过期自动刷新
	if resp.StatusCode == 401 && !isRetry {
		d.jwt = ""
		return d.doRequest(method, path, params, true)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	return result, nil
}

func (d *V10Driver) FetchProducts() ([]RemoteProduct, error) {
	result, err := d.request("GET", "/console/v1/products", nil)
	if err != nil {
		return nil, err
	}

	var products []RemoteProduct
	if data, ok := result["data"].(map[string]interface{}); ok {
		if items, ok := data["product"].([]interface{}); ok {
			for _, item := range items {
				if m, ok := item.(map[string]interface{}); ok {
					p := RemoteProduct{
						ID:   fmt.Sprintf("%v", m["id"]),
						Name: fmt.Sprintf("%v", m["name"]),
					}
					if price, ok := m["price"].(float64); ok {
						p.Price = price
					}
					if gid, ok := m["gid"].(float64); ok {
						p.GroupID = fmt.Sprintf("%.0f", gid)
					}
					if gname, ok := m["group_name"].(string); ok {
						p.GroupName = gname
					}
					if stock, ok := m["stock"].(float64); ok {
						p.Stock = int(stock)
					}
					products = append(products, p)
				}
			}
		}
	}
	return products, nil
}

func (d *V10Driver) FetchProductGroups() ([]RemoteGroup, error) {
	result, err := d.request("GET", "/console/v1/hosts/cates", nil)
	if err != nil {
		return nil, err
	}

	var groups []RemoteGroup
	if data, ok := result["data"].(map[string]interface{}); ok {
		if items, ok := data["cate"].([]interface{}); ok {
			for _, item := range items {
				if m, ok := item.(map[string]interface{}); ok {
					g := RemoteGroup{
						ID:   fmt.Sprintf("%v", m["id"]),
						Name: fmt.Sprintf("%v", m["name"]),
					}
					groups = append(groups, g)
				}
			}
		}
	}
	return groups, nil
}

func (d *V10Driver) GetProductStructure() (*ProductStructure, error) {
	groups, err := d.FetchProductGroups()
	if err != nil {
		return nil, err
	}
	return &ProductStructure{Groups: groups}, nil
}

func (d *V10Driver) SyncStatus(serviceID string) (*StatusResult, error) {
	result, err := d.request("GET", "/console/v1/mf_cloud/"+serviceID+"/status", nil)
	if err != nil {
		// 尝试dcim类型
		result, err = d.request("GET", "/console/v1/mf_dcim/"+serviceID+"/status", nil)
		if err != nil {
			return nil, err
		}
	}

	sr := &StatusResult{Status: "unknown"}
	if data, ok := result["data"].(map[string]interface{}); ok {
		if s, ok := data["status"].(string); ok {
			sr.Status = s
		}
		if ip, ok := data["ip"].(string); ok {
			sr.IPAddress = ip
		}
	}
	return sr, nil
}

func (d *V10Driver) CreateService(params CreateServiceParams) (*ServiceResult, error) {
	// v10通过购物车流程
	cartData := map[string]interface{}{
		"product_id": params.ProductID,
		"qty":        1,
		"cycle":      params.Cycle,
	}
	_, err := d.request("POST", "/console/v1/cart", cartData)
	if err != nil {
		return nil, fmt.Errorf("加入购物车失败: %w", err)
	}

	_, err = d.request("POST", "/console/v1/cart/settle", nil)
	if err != nil {
		return nil, fmt.Errorf("结算失败: %w", err)
	}

	_, err = d.request("POST", "/console/v1/pay", nil)
	if err != nil {
		return nil, fmt.Errorf("支付失败: %w", err)
	}

	return &ServiceResult{}, nil
}

func (d *V10Driver) SuspendService(serviceID string) error {
	_, err := d.request("POST", "/console/v1/host/"+serviceID+"/module/suspend", nil)
	return err
}

func (d *V10Driver) UnsuspendService(serviceID string) error {
	_, err := d.request("POST", "/console/v1/host/"+serviceID+"/module/unsuspend", nil)
	return err
}

func (d *V10Driver) TerminateService(serviceID string) error {
	_, err := d.request("POST", "/console/v1/refund", map[string]interface{}{"host_id": serviceID})
	return err
}

func (d *V10Driver) RenewService(serviceID string, cycle string) error {
	_, err := d.request("POST", "/console/v1/host/"+serviceID+"/renew", map[string]interface{}{"billingcycle": cycle})
	return err
}

// ============================================================
// AnchorDriver 锚点自有API驱动（X-API-Key认证）
// ============================================================

type AnchorDriver struct {
	supplierID uint
	apiURL     string // 我们自己的API地址
	apiKey     string // API密钥
	client     *http.Client
}

func NewAnchorDriver(supplierID uint, apiURL, apiKey, _ string) *AnchorDriver {
	return &AnchorDriver{
		supplierID: supplierID,
		apiURL:     apiURL,
		apiKey:     apiKey,
		client:     &http.Client{Timeout: 30 * time.Second},
	}
}

func (d *AnchorDriver) Key() string           { return "anchor" }
func (d *AnchorDriver) Name() string           { return "锚点接口" }
func (d *AnchorDriver) Capabilities() []string { return []string{"provisioning", "renewal", "status_sync", "product_sync"} }

func (d *AnchorDriver) request(method, path string, params map[string]interface{}) (map[string]interface{}, error) {
	var req *http.Request
	if method == "GET" {
		req, _ = http.NewRequest(method, d.apiURL+path, nil)
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, fmt.Sprintf("%v", v))
		}
		req.URL.RawQuery = q.Encode()
	} else {
		body, _ := json.Marshal(params)
		req, _ = http.NewRequest(method, d.apiURL+path, bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("X-API-Key", d.apiKey)

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)
	return result, nil
}

func (d *AnchorDriver) FetchProducts() ([]RemoteProduct, error) {
	result, err := d.request("GET", "/api/client/products", nil)
	if err != nil {
		return nil, err
	}

	var products []RemoteProduct
	if data, ok := result["data"].(map[string]interface{}); ok {
		if items, ok := data["list"].([]interface{}); ok {
			for _, item := range items {
				if m, ok := item.(map[string]interface{}); ok {
					p := RemoteProduct{
						ID:   fmt.Sprintf("%v", m["id"]),
						Name: fmt.Sprintf("%v", m["name"]),
					}
					if price, ok := m["price"].(float64); ok {
						p.Price = price
					}
					if gid, ok := m["group_id"].(float64); ok {
						p.GroupID = fmt.Sprintf("%.0f", gid)
					}
					if stock, ok := m["stock"].(float64); ok {
						p.Stock = int(stock)
					}
					products = append(products, p)
				}
			}
		}
	}
	return products, nil
}

func (d *AnchorDriver) FetchProductGroups() ([]RemoteGroup, error) {
	result, err := d.request("GET", "/api/client/products/categories", nil)
	if err != nil {
		return nil, err
	}

	var groups []RemoteGroup
	if data, ok := result["data"].([]interface{}); ok {
		for _, item := range data {
			if m, ok := item.(map[string]interface{}); ok {
				g := RemoteGroup{
					ID:   fmt.Sprintf("%v", m["id"]),
					Name: fmt.Sprintf("%v", m["name"]),
				}
				groups = append(groups, g)
			}
		}
	}
	return groups, nil
}

func (d *AnchorDriver) GetProductStructure() (*ProductStructure, error) {
	groups, err := d.FetchProductGroups()
	if err != nil {
		return nil, err
	}
	return &ProductStructure{Groups: groups}, nil
}

func (d *AnchorDriver) SyncStatus(serviceID string) (*StatusResult, error) {
	result, err := d.request("GET", "/api/client/services/"+serviceID, nil)
	if err != nil {
		return nil, err
	}

	sr := &StatusResult{Status: "unknown"}
	if data, ok := result["data"].(map[string]interface{}); ok {
		if s, ok := data["status"].(string); ok {
			sr.Status = s
		}
	}
	return sr, nil
}

func (d *AnchorDriver) CreateService(params CreateServiceParams) (*ServiceResult, error) {
	// 通过购物车流程
	_, err := d.request("POST", "/api/client/cart", map[string]interface{}{
		"product_id": params.ProductID,
		"quantity":   1,
		"cycle":      params.Cycle,
	})
	if err != nil {
		return nil, fmt.Errorf("加入购物车失败: %w", err)
	}

	_, err = d.request("POST", "/api/client/cart/checkout", nil)
	if err != nil {
		return nil, fmt.Errorf("结算失败: %w", err)
	}

	return &ServiceResult{}, nil
}

func (d *AnchorDriver) SuspendService(serviceID string) error {
	_, err := d.request("POST", "/api/admin/services/"+serviceID+"/suspend", nil)
	return err
}

func (d *AnchorDriver) UnsuspendService(serviceID string) error {
	_, err := d.request("POST", "/api/admin/services/"+serviceID+"/unsuspend", nil)
	return err
}

func (d *AnchorDriver) TerminateService(serviceID string) error {
	_, err := d.request("POST", "/api/admin/services/"+serviceID+"/terminate", nil)
	return err
}

func (d *AnchorDriver) RenewService(serviceID string, cycle string) error {
	_, err := d.request("POST", "/api/admin/services/"+serviceID+"/renewals", map[string]interface{}{"cycle": cycle})
	return err
}
