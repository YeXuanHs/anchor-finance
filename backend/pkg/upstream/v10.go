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

// v10Client implements Client for V10-type panel APIs (token-based REST).
type v10Client struct {
	baseURL string
	token   string
}

func newV10Client(p *model.UpstreamProvider) *v10Client {
	token := p.APIKey
	if t, ok := p.Config["token"].(string); ok && t != "" {
		token = t
	}
	return &v10Client{
		baseURL: strings.TrimRight(p.APIURL, "/"),
		token:   token,
	}
}

// v10Response is the common V10 API response envelope.
type v10Response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// doRequest sends a GET request with Bearer token auth.
func (c *v10Client) doRequest(method, path string) (*v10Response, error) {
	client := &http.Client{Timeout: newHTTPTimeout()}
	reqURL := c.baseURL + path

	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("v10 build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("v10 request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("v10 read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("v10 http %d: %s", resp.StatusCode, string(body))
	}

	var apiResp v10Response
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("v10 parse response: %w", err)
	}
	return &apiResp, nil
}

// TestConnection calls the /api/status endpoint to verify the token.
func (c *v10Client) TestConnection() (*ConnectionResult, error) {
	start := time.Now()

	_, err := c.doRequest("GET", "/api/status")
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return &ConnectionResult{OK: false, Message: err.Error(), Latency: latency}, err
	}

	return &ConnectionResult{
		OK:      true,
		Message: "connection successful",
		Latency: latency,
	}, nil
}

// FetchProducts retrieves products from the V10 panel.
func (c *v10Client) FetchProducts() ([]RemoteProduct, error) {
	resp, err := c.doRequest("GET", "/api/products")
	if err != nil {
		return nil, err
	}
	if resp.Code != 0 && resp.Code != 200 {
		return nil, fmt.Errorf("v10 api error: %s", resp.Message)
	}

	var raw []struct {
		ID           int                    `json:"id"`
		Name         string                 `json:"name"`
		Description  string                 `json:"description"`
		Price        float64                `json:"price"`
		Currency     string                 `json:"currency"`
		BillingCycle string                 `json:"billing_cycle"`
		Type         string                 `json:"type"`
		Stock        int                    `json:"stock"`
		Options      map[string]interface{} `json:"options"`
	}
	if err := json.Unmarshal(resp.Data, &raw); err != nil {
		return nil, fmt.Errorf("v10 parse products: %w", err)
	}

	products := make([]RemoteProduct, 0, len(raw))
	for _, p := range raw {
		currency := p.Currency
		if currency == "" {
			currency = "CNY"
		}
		products = append(products, RemoteProduct{
			RemoteID:      fmt.Sprintf("%d", p.ID),
			Name:          p.Name,
			Description:   p.Description,
			Price:         p.Price,
			Currency:      currency,
			BillingCycle:  p.BillingCycle,
			Type:          p.Type,
			Stock:         p.Stock,
			ConfigOptions: p.Options,
		})
	}
	return products, nil
}

func (c *v10Client) FetchProductsWithGroups() (*UpstreamProductsResult, error) {
	products, err := c.FetchProducts()
	if err != nil {
		return nil, err
	}
	return &UpstreamProductsResult{
		Products: products,
		Groups:   []RemoteProductGroup{},
		Currency: "CNY",
	}, nil
}

func (c *v10Client) FetchProductsByGroup(groupID string) ([]RemoteProduct, error) {
	return c.FetchProducts()
}
