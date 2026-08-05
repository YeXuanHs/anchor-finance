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
	baseURL string
	apiKey  string
	config  map[string]interface{}
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

// FetchConfigOptions retrieves configurable options for a product from zjmf upstream.
func (c *zjmfClient) FetchConfigOptions(productID string) ([]RemoteConfigOption, error) {
	resp, err := c.doRequest("getproductconfigoptions", map[string]string{"pid": productID})
	if err != nil {
		return nil, err
	}
	if resp.Result != "success" {
		return nil, fmt.Errorf("zjmf api error: %s", resp.Msg)
	}

	var raw []struct {
		Name    string   `json:"name"`
		Type    string   `json:"type"`
		Options []string `json:"options"`
	}
	if err := json.Unmarshal(resp.Data, &raw); err != nil {
		return nil, fmt.Errorf("parse config options: %w", err)
	}

	opts := make([]RemoteConfigOption, 0, len(raw))
	for _, o := range raw {
		opts = append(opts, RemoteConfigOption{
			Name:    o.Name,
			Type:    o.Type,
			Options: o.Options,
		})
	}
	return opts, nil
}
