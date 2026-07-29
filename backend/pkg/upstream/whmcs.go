package upstream

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"anchorfinance/internal/model"
)

// whmcsClient implements Client for WHMCS-compatible billing systems.
type whmcsClient struct {
	baseURL    string
	identifier string
	secret     string
}

func newWHMCSClient(p *model.UpstreamProvider) *whmcsClient {
	cfg := map[string]interface{}{}
	if p.Config != nil {
		cfg = p.Config
	}
	identifier, _ := cfg["identifier"].(string)
	if identifier == "" {
		identifier = p.APIKey
	}
	secret, _ := cfg["secret"].(string)

	return &whmcsClient{
		baseURL:    strings.TrimRight(p.APIURL, "/"),
		identifier: identifier,
		secret:     secret,
	}
}

// whmcsResponse is the common WHMCS API response envelope.
type whmcsResponse struct {
	Result  string          `json:"result"`
	Message string          `json:"message"`
	Raw     json.RawMessage `json:"-"`
}

// doRequest sends a signed WHMCS API request.
func (c *whmcsClient) doRequest(action string, params map[string]string) (*whmcsResponse, error) {
	form := url.Values{}
	form.Set("action", action)
	form.Set("identifier", c.identifier)
	form.Set("secret", c.secret)
	form.Set("responsetype", "json")
	for k, v := range params {
		form.Set(k, v)
	}

	client := &http.Client{Timeout: newHTTPTimeout()}
	apiURL := c.baseURL + "/includes/api.php"
	resp, err := client.Post(apiURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("whmcs request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("whmcs read body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("whmcs http %d: %s", resp.StatusCode, string(body))
	}

	var apiResp whmcsResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("whmcs parse response: %w", err)
	}
	apiResp.Raw = json.RawMessage(body)
	return &apiResp, nil
}

// TestConnection calls GetSystemStatus to verify WHMCS credentials.
func (c *whmcsClient) TestConnection() (*ConnectionResult, error) {
	start := time.Now()

	resp, err := c.doRequest("GetSystemStatus", nil)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return &ConnectionResult{OK: false, Message: err.Error(), Latency: latency}, err
	}
	if resp.Result != "success" {
		msg := fmt.Sprintf("whmcs error: %s", resp.Message)
		return &ConnectionResult{OK: false, Message: msg, Latency: latency}, fmt.Errorf("%s", msg)
	}

	return &ConnectionResult{
		OK:      true,
		Message: "connection successful",
		Latency: latency,
	}, nil
}

// FetchProducts retrieves products from WHMCS via GetProducts API.
func (c *whmcsClient) FetchProducts() ([]RemoteProduct, error) {
	resp, err := c.doRequest("GetProducts", map[string]string{
		"pid":    "0",
		"gid":    "0",
	})
	if err != nil {
		return nil, err
	}
	if resp.Result != "success" {
		return nil, fmt.Errorf("whmcs error: %s", resp.Message)
	}

	var parsed struct {
		Products struct {
			Product []struct {
				ID            int     `json:"pid"`
				Name          string  `json:"name"`
				Description   string  `json:"description"`
				Price         float64 `json:"pricing"`
				Currency      string  `json:"currency"`
				BillingCycle  string  `json:"paytype"`
				Type          string  `json:"type"`
				Stock         string  `json:"stockcontrol"`
			} `json:"product"`
		} `json:"products"`
	}
	if err := json.Unmarshal(resp.Raw, &parsed); err != nil {
		return nil, fmt.Errorf("whmcs parse products: %w", err)
	}

	products := make([]RemoteProduct, 0, len(parsed.Products.Product))
	for _, p := range parsed.Products.Product {
		currency := p.Currency
		if currency == "" {
			currency = "USD"
		}
		products = append(products, RemoteProduct{
			RemoteID:     fmt.Sprintf("%d", p.ID),
			Name:         p.Name,
			Description:  p.Description,
			Price:        p.Price,
			Currency:     currency,
			BillingCycle: normalizeBillingCycle(p.BillingCycle),
			Type:         p.Type,
		})
	}
	return products, nil
}
