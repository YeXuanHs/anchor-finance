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

// zjmfClient implements Client for 智简魔方 (ZJMF) compatible panels.
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

// zjmfAPIResponse is the common envelope returned by ZJMF-style APIs.
type zjmfAPIResponse struct {
	Result string          `json:"result"`
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
}

// zjmfSign computes the ZJMF API signature: md5(sorted_params + api_key).
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

// doRequest builds a signed ZJMF API request and returns the parsed response.
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

// TestConnection calls a lightweight ZJMF API endpoint to verify credentials.
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

// FetchProducts retrieves the product list from a ZJMF panel.
func (c *zjmfClient) FetchProducts() ([]RemoteProduct, error) {
	resp, err := c.doRequest("getproducts", nil)
	if err != nil {
		return nil, err
	}
	if resp.Result != "success" {
		return nil, fmt.Errorf("zjmf api error: %s", resp.Msg)
	}

	var raw []struct {
		ID          int     `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Currency    string  `json:"currency"`
		Billing     string  `json:"billingcycle"`
		Type        string  `json:"type"`
		Stock       int     `json:"stock"`
	}
	if err := json.Unmarshal(resp.Data, &raw); err != nil {
		return nil, fmt.Errorf("zjmf parse products: %w", err)
	}

	products := make([]RemoteProduct, 0, len(raw))
	for _, p := range raw {
		currency := p.Currency
		if currency == "" {
			currency = "CNY"
		}
		products = append(products, RemoteProduct{
			RemoteID:     fmt.Sprintf("%d", p.ID),
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			Currency:     currency,
			BillingCycle: p.Billing,
			Type:         p.Type,
			Stock:        p.Stock,
		})
	}
	return products, nil
}
