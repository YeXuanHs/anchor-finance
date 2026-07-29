package payment

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// PayPalGateway implements the Gateway interface for PayPal.
type PayPalGateway struct {
	ClientID     string
	ClientSecret string
	Mode         string // sandbox or live
	WebhookID    string
	NotifyURL    string
	ReturnURL    string
	CancelURL    string

	// OAuth token cache — tokens are valid for ~8 hours.
	tokenMu       sync.Mutex
	cachedToken   string
	tokenExpiry   time.Time
}

// PayPalConfig is the JSON schema stored in payment_gateways.config.
type PayPalConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Mode         string `json:"mode"` // sandbox or live
	WebhookID    string `json:"webhook_id"`
	NotifyURL    string `json:"notify_url"`
	ReturnURL    string `json:"return_url"`
	CancelURL    string `json:"cancel_url"`
}

// NewPayPalGateway creates a PayPal gateway from a JSON config string.
func NewPayPalGateway(configJSON string) (*PayPalGateway, error) {
	var cfg PayPalConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("paypal config parse error: %w", err)
	}
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, fmt.Errorf("paypal client_id and client_secret are required")
	}
	if cfg.Mode == "" {
		cfg.Mode = "sandbox"
	}
	return &PayPalGateway{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Mode:         cfg.Mode,
		WebhookID:    cfg.WebhookID,
		NotifyURL:    cfg.NotifyURL,
		ReturnURL:    cfg.ReturnURL,
		CancelURL:    cfg.CancelURL,
	}, nil
}

func (g *PayPalGateway) Name() string { return "paypal" }

func (g *PayPalGateway) apiBase() string {
	if g.Mode == "live" {
		return "https://api-m.paypal.com"
	}
	return "https://api-m.sandbox.paypal.com"
}

// getAccessToken obtains an OAuth2 access token from PayPal, using a cached
// token when available. Tokens are cached until 5 minutes before expiry.
func (g *PayPalGateway) getAccessToken(ctx context.Context) (string, error) {
	g.tokenMu.Lock()
	defer g.tokenMu.Unlock()

	if g.cachedToken != "" && time.Now().Before(g.tokenExpiry) {
		return g.cachedToken, nil
	}

	tokenURL := g.apiBase() + "/v1/oauth2/token"
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBufferString("grant_type=client_credentials"))
	if err != nil {
		return "", fmt.Errorf("paypal: create token request: %w", err)
	}

	credentials := base64.StdEncoding.EncodeToString([]byte(g.ClientID + ":" + g.ClientSecret))
	req.Header.Set("Authorization", "Basic "+credentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("paypal: token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("paypal: read token response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("paypal: token error status=%d body=%s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("paypal: parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("paypal: empty access_token in response")
	}

	// Cache with a 5-minute safety margin before actual expiry.
	expirySeconds := tokenResp.ExpiresIn
	if expirySeconds <= 0 {
		expirySeconds = 28800 // default 8 hours
	}
	g.cachedToken = tokenResp.AccessToken
	g.tokenExpiry = time.Now().Add(time.Duration(expirySeconds-300) * time.Second)

	return g.cachedToken, nil
}

// CreatePayment creates a PayPal order and returns the approval URL.
func (g *PayPalGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("paypal: amount must be positive")
	}

	accessToken, err := g.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	orderRequest := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{
			{
				"reference_id": param.OrderNo,
				"description":  param.Subject,
				"amount": map[string]interface{}{
					"currency_code": "USD",
					"value":         fmt.Sprintf("%.2f", param.Amount),
				},
			},
		},
		"application_context": map[string]interface{}{
			"return_url":  fmt.Sprintf("%s?order_id=%s", g.ReturnURL, param.OrderNo),
			"cancel_url":  fmt.Sprintf("%s?order_id=%s", g.CancelURL, param.OrderNo),
			"brand_name":  "AnchorFinance",
			"user_action": "PAY_NOW",
		},
	}

	reqBody, err := json.Marshal(orderRequest)
	if err != nil {
		return nil, fmt.Errorf("paypal: marshal order request: %w", err)
	}

	url := g.apiBase() + "/v2/checkout/orders"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("paypal: create order request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("paypal: order request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("paypal: read order response: %w", err)
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("paypal: order error status=%d body=%s", resp.StatusCode, string(body))
	}

	var orderResp struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Links  []struct {
			Href   string `json:"href"`
			Rel    string `json:"rel"`
			Method string `json:"method"`
		} `json:"links"`
	}
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, fmt.Errorf("paypal: parse order response: %w", err)
	}

	// Find the approval URL (rel = "approve")
	approvalURL := ""
	for _, link := range orderResp.Links {
		if link.Rel == "approve" {
			approvalURL = link.Href
			break
		}
	}

	if approvalURL == "" {
		approvalURL = fmt.Sprintf("%s/checkout/orders/%s", g.apiBase(), orderResp.ID)
	}

	return &PaymentResult{
		TradeNo: orderResp.ID,
		PayURL:  approvalURL,
		RawData: map[string]interface{}{
			"order_id":  orderResp.ID,
			"status":    orderResp.Status,
			"gateway":   "paypal",
			"order_no":  param.OrderNo,
			"amount":    param.Amount,
			"currency":  "USD",
		},
	}, nil
}

// QueryPayment queries the status of a PayPal order.
func (g *PayPalGateway) QueryPayment(ctx context.Context, tradeNo string) (*QueryResult, error) {
	if tradeNo == "" {
		return nil, fmt.Errorf("paypal: trade_no is required")
	}

	accessToken, err := g.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/v2/checkout/orders/%s", g.apiBase(), tradeNo)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("paypal: create query request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("paypal: query request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("paypal: read query response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("paypal: query error status=%d body=%s", resp.StatusCode, string(body))
	}

	var orderResp struct {
		ID             string `json:"id"`
		Status         string `json:"status"`
		PurchaseUnits  []struct {
			Payments struct {
				Captures []struct {
					ID     string `json:"id"`
					Status string `json:"status"`
					Amount struct {
						Value    string `json:"value"`
						Currency string `json:"currency_code"`
					} `json:"amount"`
				} `json:"captures"`
			} `json:"payments"`
		} `json:"purchase_units"`
	}
	if err := json.Unmarshal(body, &orderResp); err != nil {
		return nil, fmt.Errorf("paypal: parse query response: %w", err)
	}

	status := "pending"
	switch orderResp.Status {
	case "COMPLETED", "APPROVED":
		status = "success"
	case "VOIDED":
		status = "failed"
	case "CREATED":
		status = "pending"
	}

	var amount float64
	var paidAt *time.Time
	if len(orderResp.PurchaseUnits) > 0 && len(orderResp.PurchaseUnits[0].Payments.Captures) > 0 {
		capture := orderResp.PurchaseUnits[0].Payments.Captures[0]
		fmt.Sscanf(capture.Amount.Value, "%f", &amount)
		if capture.Status == "COMPLETED" {
			now := time.Now()
			paidAt = &now
		}
	}

	return &QueryResult{
		TradeNo: tradeNo,
		Status:  status,
		Amount:  amount,
		PaidAt:  paidAt,
	}, nil
}

// Refund processes a refund through PayPal.
func (g *PayPalGateway) Refund(ctx context.Context, param *RefundParam) (*RefundResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("paypal: refund amount must be positive")
	}

	accessToken, err := g.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	// PayPal requires the capture ID, not the order ID for refunds
	// The caller should pass the capture ID as TradeNo
	refundRequest := map[string]interface{}{
		"amount": map[string]interface{}{
			"value":    fmt.Sprintf("%.2f", param.Amount),
			"currency": "USD",
		},
	}

	reqBody, _ := json.Marshal(refundRequest)
	url := fmt.Sprintf("%s/v2/payments/captures/%s/refund", g.apiBase(), param.TradeNo)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("paypal: create refund request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("paypal: refund request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("paypal: read refund response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("paypal: refund error status=%d body=%s", resp.StatusCode, string(body))
	}

	var refundResp struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	json.Unmarshal(body, &refundResp)

	status := "processing"
	if refundResp.Status == "COMPLETED" {
		status = "success"
	}

	return &RefundResult{
		RefundNo: refundResp.ID,
		Status:   status,
	}, nil
}

// ParseNotify parses a PayPal webhook event without signature verification.
func (g *PayPalGateway) ParseNotify(ctx context.Context, data []byte) (*NotifyResult, error) {
	return g.parseEvent(data)
}

// VerifyAndParseNotify verifies the PayPal webhook signature and parses the event.
// headers should contain auth_algo, cert_url, transmission_id, transmission_sig,
// and transmission_time from the incoming HTTP request headers.
func (g *PayPalGateway) VerifyAndParseNotify(ctx context.Context, data []byte, headers map[string]string) (*NotifyResult, error) {
	if g.WebhookID == "" {
		return nil, fmt.Errorf("paypal: webhook_id not configured, cannot verify signature")
	}

	verifyReq := map[string]interface{}{
		"auth_algo":         headers["paypal-auth-algo"],
		"cert_url":          headers["paypal-cert-url"],
		"transmission_id":   headers["paypal-transmission-id"],
		"transmission_sig":  headers["paypal-transmission-sig"],
		"transmission_time": headers["paypal-transmission-time"],
		"webhook_id":        g.WebhookID,
		"webhook_event":     json.RawMessage(data),
	}

	verifyBody, err := json.Marshal(verifyReq)
	if err != nil {
		return nil, fmt.Errorf("paypal: marshal verify request: %w", err)
	}

	accessToken, err := g.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	verifyURL := g.apiBase() + "/v1/notifications/verify-webhook-signature"
	req, err := http.NewRequestWithContext(ctx, "POST", verifyURL, bytes.NewBuffer(verifyBody))
	if err != nil {
		return nil, fmt.Errorf("paypal: create verify request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("paypal: verify request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("paypal: read verify response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("paypal: verify error status=%d body=%s", resp.StatusCode, string(respBody))
	}

	var verifyResp struct {
		VerificationStatus string `json:"verification_status"`
	}
	if err := json.Unmarshal(respBody, &verifyResp); err != nil {
		return nil, fmt.Errorf("paypal: parse verify response: %w", err)
	}
	if verifyResp.VerificationStatus != "SUCCESS" {
		return nil, fmt.Errorf("paypal: webhook signature verification failed: %s", verifyResp.VerificationStatus)
	}

	return g.parseEvent(data)
}

// parseEvent is the shared logic for parsing a PayPal webhook event body.
func (g *PayPalGateway) parseEvent(data []byte) (*NotifyResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("paypal: empty notification data")
	}

	var event struct {
		EventType string `json:"event_type"`
		Resource  struct {
			ID            string `json:"id"`
			Status        string `json:"status"`
			CustomID      string `json:"custom_id"`
			Supplementary struct {
				Amounts struct {
					Total struct {
						Value string `json:"value"`
					} `json:"total"`
				} `json:"amounts"`
			} `json:"supplementary_data"`
		} `json:"resource"`
	}
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("paypal: parse webhook: %w", err)
	}

	status := "pending"
	switch event.EventType {
	case "CHECKOUT.ORDER.APPROVED", "PAYMENT.CAPTURE.COMPLETED":
		status = "success"
	case "CHECKOUT.ORDER.CANCELLED", "PAYMENT.CAPTURE.DENIED", "PAYMENT.CAPTURE.REFUNDED":
		status = "failed"
	}

	orderNo := event.Resource.CustomID
	if orderNo == "" {
		orderNo = event.Resource.ID
	}

	return &NotifyResult{
		TradeNo: event.Resource.ID,
		OrderNo: orderNo,
		Status:  status,
	}, nil
}
