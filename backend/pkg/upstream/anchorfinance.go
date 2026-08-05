package upstream

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"anchorfinance/internal/model"
)

type anchorFinanceClient struct {
	baseURL string
	token   string
	config  map[string]interface{}
}

func newAnchorFinanceClient(p *model.UpstreamProvider) *anchorFinanceClient {
	cfg := map[string]interface{}{}
	if p.Config != nil {
		cfg = p.Config
	}
	return &anchorFinanceClient{
		baseURL: strings.TrimRight(p.APIURL, "/"),
		token:   p.APIKey,
		config:  cfg,
	}
}

type afAPIResponse struct {
	Code int             `json:"code"`
	Data json.RawMessage `json:"data"`
}

type afPagedData struct {
	List  json.RawMessage `json:"list"`
	Total int64           `json:"total"`
}

func (c *anchorFinanceClient) doRequest(method, path string) ([]byte, int, error) {
	client := &http.Client{Timeout: newHTTPTimeout()}
	reqURL := c.baseURL + path

	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("af build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("af request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("af read body: %w", err)
	}
	return body, resp.StatusCode, nil
}

func (c *anchorFinanceClient) parseResponse(body []byte) (*afAPIResponse, error) {
	var apiResp afAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("af parse response: %w", err)
	}
	if apiResp.Code != 200 {
		return nil, fmt.Errorf("af api error: code=%d", apiResp.Code)
	}
	return &apiResp, nil
}

func (c *anchorFinanceClient) parsePagedData(data json.RawMessage) (*afPagedData, error) {
	var paged afPagedData
	if err := json.Unmarshal(data, &paged); err != nil {
		return nil, fmt.Errorf("af parse paged data: %w", err)
	}
	return &paged, nil
}

func (c *anchorFinanceClient) TestConnection() (*ConnectionResult, error) {
	start := time.Now()

	body, statusCode, err := c.doRequest("GET", "/api/v1/products?page=1&page_size=1")
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return &ConnectionResult{OK: false, Message: err.Error(), Latency: latency}, err
	}
	if statusCode >= 400 {
		msg := fmt.Sprintf("af api returned HTTP %d", statusCode)
		return &ConnectionResult{OK: false, Message: msg, Latency: latency}, fmt.Errorf("%s", msg)
	}

	if _, err := c.parseResponse(body); err != nil {
		return &ConnectionResult{OK: false, Message: err.Error(), Latency: latency}, err
	}

	return &ConnectionResult{
		OK:      true,
		Message: "connection successful",
		Latency: latency,
	}, nil
}

func (c *anchorFinanceClient) FetchProducts() ([]RemoteProduct, error) {
	body, statusCode, err := c.doRequest("GET", "/api/v1/products")
	if err != nil {
		return nil, err
	}
	if statusCode >= 400 {
		return nil, fmt.Errorf("af api returned HTTP %d: %s", statusCode, string(body))
	}

	apiResp, err := c.parseResponse(body)
	if err != nil {
		return nil, err
	}

	paged, err := c.parsePagedData(apiResp.Data)
	if err != nil {
		return nil, err
	}

	var raw []struct {
		ID          interface{} `json:"id"`
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Price       float64     `json:"price"`
		Currency    string      `json:"currency"`
		Billing     string      `json:"billing_cycle"`
		Type        string      `json:"type"`
		Stock       int         `json:"stock"`
		GroupID     interface{} `json:"group_id"`
		GroupName   string      `json:"group_name"`
	}
	if err := json.Unmarshal(paged.List, &raw); err != nil {
		return nil, fmt.Errorf("af parse products: %w", err)
	}

	products := make([]RemoteProduct, 0, len(raw))
	for _, p := range raw {
		currency := p.Currency
		if currency == "" {
			currency = "CNY"
		}
		groupID := fmt.Sprintf("%v", p.GroupID)
		if groupID == "<nil>" {
			groupID = ""
		}
		products = append(products, RemoteProduct{
			RemoteID:     fmt.Sprintf("%v", p.ID),
			Name:         p.Name,
			Description:  p.Description,
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

func (c *anchorFinanceClient) FetchProductsWithGroups() (*UpstreamProductsResult, error) {
	// Fetch groups
	groupsBody, _, err := c.doRequest("GET", "/api/v1/product-groups")
	if err != nil {
		return nil, err
	}

	groupsResp, err := c.parseResponse(groupsBody)
	if err != nil {
		return nil, err
	}

	groupsPaged, err := c.parsePagedData(groupsResp.Data)
	if err != nil {
		return nil, err
	}

	var rawGroups []struct {
		ID   interface{} `json:"id"`
		Name string      `json:"name"`
	}
	if err := json.Unmarshal(groupsPaged.List, &rawGroups); err != nil {
		return nil, fmt.Errorf("af parse groups: %w", err)
	}

	groupMap := make(map[string]string)
	var groups []RemoteProductGroup
	for _, g := range rawGroups {
		gid := fmt.Sprintf("%v", g.ID)
		groupMap[gid] = g.Name
		groups = append(groups, RemoteProductGroup{
			GroupID: gid,
			Name:    g.Name,
		})
	}

	// Fetch products
	products, err := c.FetchProducts()
	if err != nil {
		return nil, err
	}

	// Enrich group info
	for i := range products {
		if products[i].GroupName == "" && products[i].GroupID != "" {
			if name, ok := groupMap[products[i].GroupID]; ok {
				products[i].GroupName = name
			}
		}
	}

	// Update product counts per group
	countMap := make(map[string]int)
	for _, p := range products {
		if p.GroupID != "" {
			countMap[p.GroupID]++
		}
	}
	for i := range groups {
		groups[i].ProductCount = countMap[groups[i].GroupID]
	}

	return &UpstreamProductsResult{
		Products: products,
		Groups:   groups,
		Currency: "CNY",
	}, nil
}

func (c *anchorFinanceClient) FetchProductsByGroup(groupID string) ([]RemoteProduct, error) {
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

// FetchConfigOptions retrieves configurable options for a product from another AnchorFinance upstream.
func (c *anchorFinanceClient) FetchConfigOptions(productID string) ([]RemoteConfigOption, error) {
	body, statusCode, err := c.doRequest("GET", fmt.Sprintf("/api/v1/products/%s/config-options", productID))
	if err != nil {
		return nil, err
	}
	if statusCode >= 400 {
		return nil, fmt.Errorf("af api returned HTTP %d", statusCode)
	}

	apiResp, err := c.parseResponse(body)
	if err != nil {
		return nil, err
	}

	var raw []struct {
		Name    string   `json:"name"`
		Type    string   `json:"type"`
		Options []string `json:"options"`
	}
	if err := json.Unmarshal(apiResp.Data, &raw); err != nil {
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
