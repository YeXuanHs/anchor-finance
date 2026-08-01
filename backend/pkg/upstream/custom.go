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

// customClient implements Client for generic/custom upstream APIs.
// It supports configurable auth method, endpoint paths, and response mapping.
type customClient struct {
	baseURL      string
	apiKey       string
	authMethod   string // "bearer", "header", "query"
	authHeader   string
	testEndpoint string
	productPath  string
	resultField  string
}

func newCustomClient(p *model.UpstreamProvider) *customClient {
	cfg := map[string]interface{}{}
	if p.Config != nil {
		cfg = p.Config
	}

	authMethod, _ := cfg["auth_method"].(string)
	if authMethod == "" {
		authMethod = "bearer"
	}
	authHeader, _ := cfg["auth_header"].(string)
	if authHeader == "" {
		authHeader = "Authorization"
	}
	testEndpoint, _ := cfg["test_endpoint"].(string)
	if testEndpoint == "" {
		testEndpoint = "/api/test"
	}
	productPath, _ := cfg["product_endpoint"].(string)
	if productPath == "" {
		productPath = "/api/products"
	}
	resultField, _ := cfg["result_field"].(string)
	if resultField == "" {
		resultField = "data"
	}

	return &customClient{
		baseURL:      strings.TrimRight(p.APIURL, "/"),
		apiKey:       p.APIKey,
		authMethod:   authMethod,
		authHeader:   authHeader,
		testEndpoint: testEndpoint,
		productPath:  productPath,
		resultField:  resultField,
	}
}

// doRequest sends an HTTP request with the configured auth method.
func (c *customClient) doRequest(method, path string) ([]byte, int, error) {
	client := &http.Client{Timeout: newHTTPTimeout()}
	reqURL := c.baseURL + path

	req, err := http.NewRequest(method, reqURL, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("custom build request: %w", err)
	}

	switch c.authMethod {
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	case "header":
		req.Header.Set(c.authHeader, c.apiKey)
	case "query":
		q := req.URL.Query()
		q.Set("api_key", c.apiKey)
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("custom request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("custom read body: %w", err)
	}
	return body, resp.StatusCode, nil
}

// TestConnection sends a request to the configured test endpoint.
func (c *customClient) TestConnection() (*ConnectionResult, error) {
	start := time.Now()

	body, statusCode, err := c.doRequest("GET", c.testEndpoint)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return &ConnectionResult{OK: false, Message: err.Error(), Latency: latency}, err
	}

	if statusCode >= 400 {
		msg := fmt.Sprintf("custom api returned HTTP %d", statusCode)
		return &ConnectionResult{OK: false, Message: msg, Latency: latency}, fmt.Errorf("%s", msg)
	}

	return &ConnectionResult{
		OK:      true,
		Message: "connection successful",
		Latency: latency,
	}, nil
}

// FetchProducts retrieves products from the custom API endpoint.
func (c *customClient) FetchProducts() ([]RemoteProduct, error) {
	body, statusCode, err := c.doRequest("GET", c.productPath)
	if err != nil {
		return nil, err
	}
	if statusCode >= 400 {
		return nil, fmt.Errorf("custom api returned HTTP %d: %s", statusCode, string(body))
	}

	// Parse the response envelope and extract the products array from the configured field.
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("custom parse envelope: %w", err)
	}

	dataRaw, ok := envelope[c.resultField]
	if !ok {
		// Try treating the whole body as the data array directly.
		dataRaw = json.RawMessage(body)
	}

	var raw []struct {
		ID           interface{}            `json:"id"`
		Name         string                 `json:"name"`
		Description  string                 `json:"description"`
		Price        float64                `json:"price"`
		Currency     string                 `json:"currency"`
		BillingCycle string                 `json:"billing_cycle"`
		Type         string                 `json:"type"`
		Stock        int                    `json:"stock"`
		Options      map[string]interface{} `json:"options"`
	}
	if err := json.Unmarshal(dataRaw, &raw); err != nil {
		return nil, fmt.Errorf("custom parse products: %w", err)
	}

	products := make([]RemoteProduct, 0, len(raw))
	for _, p := range raw {
		var remoteID string
		switch v := p.ID.(type) {
		case float64:
			remoteID = fmt.Sprintf("%.0f", v)
		case string:
			remoteID = v
		default:
			remoteID = fmt.Sprintf("%v", v)
		}
		currency := p.Currency
		if currency == "" {
			currency = "CNY"
		}
		products = append(products, RemoteProduct{
			RemoteID:      remoteID,
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

func (c *customClient) FetchProductsWithGroups() (*UpstreamProductsResult, error) {
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

func (c *customClient) FetchProductsByGroup(groupID string) ([]RemoteProduct, error) {
	return c.FetchProducts()
}
