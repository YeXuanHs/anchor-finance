package upstream

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
)

type zjmfClient struct {
	baseURL  string
	apiKey   string
	config   map[string]interface{}
	jwtToken string // cached JWT for authenticated API calls
}

func newZJMFClient(p *model.UpstreamProvider) *zjmfClient {
	cfg := map[string]interface{}{}
	if p.Config != nil {
		cfg = p.Config
	}
	return &zjmfClient{
		baseURL: strings.TrimRight(p.APIURL, "/"),
		apiKey:  p.APIKey,
		config:  cfg,
	}
}

type zjmfAPIResponse struct {
	Result string          `json:"result"`
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
}

func (c *zjmfClient) zjmfSign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(params[k])
	}
	sb.WriteString(c.apiKey)

	return fmt.Sprintf("%x", md5.Sum([]byte(sb.String())))
}

func (c *zjmfClient) doRequest(action string, extraParams map[string]string) (*zjmfAPIResponse, error) {
	params := map[string]string{
		"action": action,
	}
	for k, v := range extraParams {
		params[k] = v
	}
	params["sign"] = c.zjmfSign(params)

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	client := &http.Client{Timeout: newHTTPTimeout()}
	apiURL := c.baseURL + "/api.php"
	resp, err := client.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("zjmf request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("zjmf read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zjmf http %d: %s", resp.StatusCode, string(body))
	}

	var apiResp zjmfAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("zjmf parse response: %w", err)
	}
	return &apiResp, nil
}

func (c *zjmfClient) TestConnection() (*ConnectionResult, error) {
	start := time.Now()

	resp, err := c.doRequest("getsysteminfo", nil)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return &ConnectionResult{OK: false, Message: err.Error(), Latency: latency}, err
	}
	if resp.Result != "success" {
		msg := fmt.Sprintf("zjmf api error: %s", resp.Msg)
		return &ConnectionResult{OK: false, Message: msg, Latency: latency}, fmt.Errorf("%s", msg)
	}

	return &ConnectionResult{
		OK:      true,
		Message: "connection successful",
		Latency: latency,
	}, nil
}

func (c *zjmfClient) FetchProducts() ([]RemoteProduct, error) {
	resp, err := c.doRequest("getproducts", nil)
	if err != nil {
		return nil, err
	}
	if resp.Result != "success" {
		return nil, fmt.Errorf("zjmf api error: %s", resp.Msg)
	}

	var raw []struct {
		ID       int     `json:"id"`
		Name     string  `json:"name"`
		Desc     string  `json:"description"`
		Price    float64 `json:"price"`
		Currency string  `json:"currency"`
		Billing  string  `json:"billingcycle"`
		Type     string  `json:"type"`
		Stock    int     `json:"stock"`
		GID      int     `json:"gid"`
		GroupName string `json:"groupname"`
	}
	if err := json.Unmarshal(resp.Data, &raw); err != nil {
		// 尝试分组结构
		return c.parseGroupedProducts(resp.Data)
	}

	products := make([]RemoteProduct, 0, len(raw))
	for _, p := range raw {
		currency := p.Currency
		if currency == "" {
			currency = "CNY"
		}
		groupID := ""
		if p.GID > 0 {
			groupID = fmt.Sprintf("%d", p.GID)
		}
		products = append(products, RemoteProduct{
			RemoteID:     fmt.Sprintf("%d", p.ID),
			Name:         p.Name,
			Description:  p.Desc,
			Price:        p.Price,
			Currency:     currency,
			BillingCycle: p.Billing,
			Type:         p.Type,
			Stock:        p.Stock,
			GroupID:      groupID,
			GroupName:    p.GroupName,
		})
	}
	return products, nil
}

// parseGroupedProducts 解析分组结构的产品列表
func (c *zjmfClient) parseGroupedProducts(data json.RawMessage) ([]RemoteProduct, error) {
	var grouped []struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Products []struct {
			ID       int     `json:"id"`
			Name     string  `json:"name"`
			Desc     string  `json:"description"`
			Price    float64 `json:"price"`
			Currency string  `json:"currency"`
			Billing  string  `json:"billingcycle"`
			Type     string  `json:"type"`
			Stock    int     `json:"stock"`
		} `json:"products"`
	}
	if err := json.Unmarshal(data, &grouped); err != nil {
		return nil, fmt.Errorf("parse grouped products: %w", err)
	}

	var products []RemoteProduct
	for _, g := range grouped {
		for _, p := range g.Products {
			currency := p.Currency
			if currency == "" {
				currency = "CNY"
			}
			products = append(products, RemoteProduct{
				RemoteID:     fmt.Sprintf("%d", p.ID),
				Name:         p.Name,
				Description:  p.Desc,
				Price:        p.Price,
				Currency:     currency,
				BillingCycle: p.Billing,
				Type:         p.Type,
				Stock:        p.Stock,
				GroupID:      fmt.Sprintf("%d", g.ID),
				GroupName:    g.Name,
			})
		}
	}
	return products, nil
}

func (c *zjmfClient) FetchProductsWithGroups() (*UpstreamProductsResult, error) {
	resp, err := c.doRequest("getproducts", nil)
	if err != nil {
		return nil, err
	}
	if resp.Result != "success" {
		return nil, fmt.Errorf("zjmf api error: %s", resp.Msg)
	}

	// 先尝试分组结构
	var grouped []struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Products []struct {
			ID       int     `json:"id"`
			Name     string  `json:"name"`
			Desc     string  `json:"description"`
			Price    float64 `json:"price"`
			Currency string  `json:"currency"`
			Billing  string  `json:"billingcycle"`
			Type     string  `json:"type"`
			Stock    int     `json:"stock"`
		} `json:"products"`
	}

	if err := json.Unmarshal(resp.Data, &grouped); err == nil && len(grouped) > 0 && len(grouped[0].Products) > 0 {
		result := &UpstreamProductsResult{
			Groups:   make([]RemoteProductGroup, 0, len(grouped)),
			Products: make([]RemoteProduct, 0),
			Currency: "CNY",
		}
		for _, g := range grouped {
			result.Groups = append(result.Groups, RemoteProductGroup{
				GroupID:      fmt.Sprintf("%d", g.ID),
				Name:         g.Name,
				ProductCount: len(g.Products),
			})
			for _, p := range g.Products {
				currency := p.Currency
				if currency == "" {
					currency = "CNY"
				}
				result.Products = append(result.Products, RemoteProduct{
					RemoteID:     fmt.Sprintf("%d", p.ID),
					Name:         p.Name,
					Description:  p.Desc,
					Price:        p.Price,
					Currency:     currency,
					BillingCycle: p.Billing,
					Type:         p.Type,
					Stock:        p.Stock,
					GroupID:      fmt.Sprintf("%d", g.ID),
					GroupName:    g.Name,
				})
			}
		}
		return result, nil
	}

	// 平铺结构
	var flat []struct {
		ID        int     `json:"id"`
		Name      string  `json:"name"`
		Desc      string  `json:"description"`
		Price     float64 `json:"price"`
		Currency  string  `json:"currency"`
		Billing   string  `json:"billingcycle"`
		Type      string  `json:"type"`
		Stock     int     `json:"stock"`
		GID       int     `json:"gid"`
		GroupName string  `json:"groupname"`
	}
	if err := json.Unmarshal(resp.Data, &flat); err != nil {
		return nil, fmt.Errorf("parse products: %w", err)
	}

	groupMap := make(map[string]*RemoteProductGroup)
	result := &UpstreamProductsResult{
		Products: make([]RemoteProduct, 0, len(flat)),
		Groups:   make([]RemoteProductGroup, 0),
		Currency: "CNY",
	}

	for _, p := range flat {
		currency := p.Currency
		if currency == "" {
			currency = "CNY"
		}
		groupID := ""
		if p.GID > 0 {
			groupID = fmt.Sprintf("%d", p.GID)
		}
		if p.GroupName == "" {
			p.GroupName = "未分组"
		}

		if groupID != "" {
			if _, exists := groupMap[groupID]; !exists {
				groupMap[groupID] = &RemoteProductGroup{
					GroupID: groupID,
					Name:    p.GroupName,
				}
			}
			groupMap[groupID].ProductCount++
		}

		result.Products = append(result.Products, RemoteProduct{
			RemoteID:     fmt.Sprintf("%d", p.ID),
			Name:         p.Name,
			Description:  p.Desc,
			Price:        p.Price,
			Currency:     currency,
			BillingCycle: p.Billing,
			Type:         p.Type,
			Stock:        p.Stock,
			GroupID:      groupID,
			GroupName:    p.GroupName,
		})
		if currency != "" {
			result.Currency = currency
		}
	}

	for _, g := range groupMap {
		result.Groups = append(result.Groups, *g)
	}

	return result, nil
}

func (c *zjmfClient) FetchProductsByGroup(groupID string) ([]RemoteProduct, error) {
	all, err := c.FetchProducts()
	if err != nil {
		return nil, err
	}

	var filtered []RemoteProduct
	for _, p := range all {
		if p.GroupID == groupID {
			filtered = append(filtered, p)
		}
	}
	return filtered, nil
}

// zjmfLogin authenticates with the upstream zjmf and caches the JWT token.
func (c *zjmfClient) zjmfLogin() error {
	if c.jwtToken != "" {
		return nil
	}

	username, _ := c.config["username"].(string)
	password, _ := c.config["password"].(string)
	if username == "" || password == "" {
		return fmt.Errorf("zjmf JWT auth requires username and password in config")
	}

	loginURL := c.baseURL + "/zjmf_api_login"
	form := url.Values{}
	form.Set("username", username)
	form.Set("password", password)

	client := &http.Client{Timeout: newHTTPTimeout()}
	resp, err := client.Post(loginURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("zjmf login failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("zjmf login read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("zjmf login http %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Status int    `json:"status"`
		JWT    string `json:"jwt"`
		Data   struct {
			JWT string `json:"jwt"`
		} `json:"data"`
		Msg string `json:"msg"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("zjmf login parse: %w", err)
	}

	jwt := result.JWT
	if jwt == "" {
		jwt = result.Data.JWT
	}
	if jwt == "" {
		return fmt.Errorf("zjmf login no jwt returned: %s", result.Msg)
	}

	c.jwtToken = jwt
	return nil
}

// doJWTRequest sends a GET request with JWT authentication.
func (c *zjmfClient) doJWTRequest(path string, params map[string]string) ([]byte, error) {
	if err := c.zjmfLogin(); err != nil {
		return nil, err
	}

	reqURL := c.baseURL + "/" + strings.TrimLeft(path, "/")
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		reqURL += "?" + q.Encode()
	}

	client := &http.Client{Timeout: newHTTPTimeout()}
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("zjmf jwt build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.jwtToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("zjmf jwt request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("zjmf jwt read body: %w", err)
	}

	// If 405, token expired - retry once
	if resp.StatusCode == 405 {
		c.jwtToken = ""
		if err := c.zjmfLogin(); err != nil {
			return nil, err
		}
		req, _ = http.NewRequest("GET", reqURL, nil)
		req.Header.Set("Authorization", "Bearer "+c.jwtToken)
		resp, err = client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("zjmf jwt retry failed: %w", err)
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zjmf jwt http %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// FetchConfigOptions retrieves configurable options for a product from zjmf upstream.
// Uses JWT auth to call cart/get_product_config, same as zjmf's getZjmfUpstreamProductConfig.
func (c *zjmfClient) FetchConfigOptions(productID string) ([]RemoteConfigGroup, error) {
	body, err := c.doJWTRequest("cart/get_product_config", map[string]string{"pid": productID})
	if err != nil {
		return nil, err
	}

	var result struct {
		Status int `json:"status"`
		Data   struct {
			ConfigGroups []struct {
				ID      int    `json:"id"`
				Name    string `json:"name"`
				Options []struct {
					ID           int    `json:"id"`
					OptionName   string `json:"option_name"`
					OptionType   int    `json:"option_type"`
					QtyMinimum   int    `json:"qty_minimum"`
					QtyMaximum   int    `json:"qty_maximum"`
					UpstreamID   int    `json:"upstream_id"`
					Sub          []struct {
						ID         int    `json:"id"`
						OptionName string `json:"option_name"`
						SortOrder  int    `json:"sort_order"`
						UpstreamID int    `json:"upstream_id"`
						Pricing    []struct {
							Currency  int     `json:"currency"`
							Monthly   float64 `json:"monthly"`
							Quarterly float64 `json:"quarterly"`
							Annually  float64 `json:"annually"`
						} `json:"pricing"`
					} `json:"sub"`
				} `json:"options"`
			} `json:"config_groups"`
		} `json:"data"`
		Msg string `json:"msg"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse config options: %w", err)
	}
	if result.Status != 200 {
		return nil, fmt.Errorf("zjmf api error: %s", result.Msg)
	}

	groups := make([]RemoteConfigGroup, 0, len(result.Data.ConfigGroups))
	for _, g := range result.Data.ConfigGroups {
		group := RemoteConfigGroup{
			RemoteID: fmt.Sprintf("%d", g.ID),
			Name:     g.Name,
			Options:  make([]RemoteConfigOption, 0, len(g.Options)),
		}
		for _, opt := range g.Options {
			option := RemoteConfigOption{
				RemoteID:  fmt.Sprintf("%d", opt.ID),
				Name:      parseOptionName(opt.OptionName),
				Type:      opt.OptionType,
				GroupName: g.Name,
				GroupID:   fmt.Sprintf("%d", g.ID),
				Sub:       make([]RemoteConfigOptionSub, 0, len(opt.Sub)),
			}
			for _, sub := range opt.Sub {
				subOpt := RemoteConfigOptionSub{
					RemoteID: fmt.Sprintf("%d", sub.ID),
					Name:     parseSubOptionName(sub.OptionName),
				}
				// For OS type (option_type=5), extract OS and version
				if opt.OptionType == 5 {
					subOpt.OS, subOpt.Version = parseOSName(sub.OptionName)
				}
				// Use first pricing if available
				if len(sub.Pricing) > 0 {
					subOpt.Price = sub.Pricing[0].Monthly
				}
				option.Sub = append(option.Sub, subOpt)
			}
			group.Options = append(group.Options, option)
		}
		groups = append(groups, group)
	}
	return groups, nil
}

// parseOptionName extracts the display name from zjmf's pipe-separated format.
// Format: "key|DisplayName" -> "DisplayName", or just "name" if no pipe.
func parseOptionName(name string) string {
	parts := strings.SplitN(name, "|", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return name
}

// parseSubOptionName extracts the display name from zjmf's pipe-separated format.
func parseSubOptionName(name string) string {
	parts := strings.SplitN(name, "|", 2)
	if len(parts) == 2 {
		name = parts[1]
	}
	// Remove version part after ^
	if idx := strings.Index(name, "^"); idx >= 0 {
		return name[:idx]
	}
	return name
}

// parseOSName extracts OS name and version from zjmf's format.
// Format: "CentOS^7.6" or "Windows Server^2019" or "os_name|display^version"
func parseOSName(name string) (os, version string) {
	parts := strings.SplitN(name, "|", 2)
	if len(parts) == 2 {
		name = parts[1]
	}
	if idx := strings.Index(name, "^"); idx >= 0 {
		return name[:idx], name[idx+1:]
	}
	return name, ""
}
