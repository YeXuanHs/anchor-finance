package payment

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// StripeGateway implements the Gateway interface for Stripe.
type StripeGateway struct {
	SecretKey     string
	WebhookSecret string
	NotifyURL     string
	ReturnURL     string
}

// StripeConfig is the JSON schema stored in payment_gateways.config.
type StripeConfig struct {
	SecretKey     string `json:"secret_key"`
	WebhookSecret string `json:"webhook_secret"`
	NotifyURL     string `json:"notify_url"`
	ReturnURL     string `json:"return_url"`
}

// NewStripeGateway creates a Stripe gateway from a JSON config string.
func NewStripeGateway(configJSON string) (*StripeGateway, error) {
	var cfg StripeConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("stripe config parse error: %w", err)
	}
	if cfg.SecretKey == "" {
		return nil, fmt.Errorf("stripe secret_key is required")
	}
	return &StripeGateway{
		SecretKey:     cfg.SecretKey,
		WebhookSecret: cfg.WebhookSecret,
		NotifyURL:     cfg.NotifyURL,
		ReturnURL:     cfg.ReturnURL,
	}, nil
}

func (g *StripeGateway) Name() string { return "stripe" }

// stripePost sends a POST request to the Stripe API using form-encoded data.
func (g *StripeGateway) stripePost(ctx context.Context, path string, params url.Values) ([]byte, int, error) {
	apiURL := "https://api.stripe.com" + path
	req, err := http.NewRequestWithContext(ctx, "POST", apiURL, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, 0, fmt.Errorf("stripe: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+g.SecretKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("stripe: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("stripe: read response: %w", err)
	}
	return body, resp.StatusCode, nil
}

// stripeGet sends a GET request to the Stripe API.
func (g *StripeGateway) stripeGet(ctx context.Context, path string) ([]byte, int, error) {
	apiURL := "https://api.stripe.com" + path
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("stripe: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+g.SecretKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("stripe: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("stripe: read response: %w", err)
	}
	return body, resp.StatusCode, nil
}

// CreatePayment creates a Stripe Checkout Session and returns the session URL.
func (g *StripeGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("stripe: amount must be positive")
	}

	amountInCents := int64(param.Amount * 100)
	successURL := g.ReturnURL
	if successURL == "" {
		successURL = "https://example.com/success"
	}

	params := url.Values{}
	params.Set("mode", "payment")
	params.Set("success_url", successURL+"?session_id={CHECKOUT_SESSION_ID}")
	if g.ReturnURL != "" {
		params.Set("cancel_url", g.ReturnURL+"?canceled=true")
	}
	params.Set("line_items[0][price_data][currency]", "usd")
	params.Set("line_items[0][price_data][product_data][name]", param.Subject)
	params.Set("line_items[0][price_data][unit_amount]", fmt.Sprintf("%d", amountInCents))
	params.Set("line_items[0][quantity]", "1")
	params.Set("metadata[order_no]", param.OrderNo)
	params.Set("client_reference_id", param.OrderNo)

	body, statusCode, err := g.stripePost(ctx, "/v1/checkout/sessions", params)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("stripe: create session error status=%d body=%s", statusCode, string(body))
	}

	var sessionResp struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &sessionResp); err != nil {
		return nil, fmt.Errorf("stripe: parse session response: %w", err)
	}

	if sessionResp.URL == "" {
		return nil, fmt.Errorf("stripe: session created but no URL returned")
	}

	return &PaymentResult{
		TradeNo: sessionResp.ID,
		PayURL:  sessionResp.URL,
		RawData: map[string]interface{}{
			"gateway":   "stripe",
			"session_id": sessionResp.ID,
			"order_no":  param.OrderNo,
			"amount":    param.Amount,
			"currency":  "usd",
		},
	}, nil
}

// QueryPayment queries the status of a Stripe PaymentIntent.
func (g *StripeGateway) QueryPayment(ctx context.Context, tradeNo string) (*QueryResult, error) {
	if tradeNo == "" {
		return nil, fmt.Errorf("stripe: payment_intent id is required")
	}

	body, statusCode, err := g.stripeGet(ctx, "/v1/payment_intents/"+tradeNo)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("stripe: query error status=%d body=%s", statusCode, string(body))
	}

	var piResp struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		Amount int64  `json:"amount"`
	}
	if err := json.Unmarshal(body, &piResp); err != nil {
		return nil, fmt.Errorf("stripe: parse query response: %w", err)
	}

	// Map Stripe PaymentIntent status to internal status.
	status := "pending"
	switch piResp.Status {
	case "succeeded":
		status = "success"
	case "canceled":
		status = "failed"
	case "requires_payment_method", "requires_confirmation", "requires_action":
		status = "pending"
	case "processing":
		status = "processing"
	}

	var paidAt *time.Time
	if status == "success" {
		now := time.Now()
		paidAt = &now
	}

	return &QueryResult{
		TradeNo: tradeNo,
		Status:  status,
		Amount:  float64(piResp.Amount) / 100,
		PaidAt:  paidAt,
	}, nil
}

// Refund processes a refund through Stripe.
func (g *StripeGateway) Refund(ctx context.Context, param *RefundParam) (*RefundResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("stripe: refund amount must be positive")
	}

	amountInCents := int64(param.Amount * 100)
	params := url.Values{}
	params.Set("amount", fmt.Sprintf("%d", amountInCents))
	params.Set("payment_intent", param.TradeNo)
	if param.Reason != "" {
		params.Set("reason", param.Reason)
	}

	body, statusCode, err := g.stripePost(ctx, "/v1/refunds", params)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("stripe: refund error status=%d body=%s", statusCode, string(body))
	}

	var refundResp struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}
	json.Unmarshal(body, &refundResp)

	status := "processing"
	switch refundResp.Status {
	case "succeeded":
		status = "success"
	case "failed":
		status = "failed"
	}

	return &RefundResult{
		RefundNo: refundResp.ID,
		Status:   status,
	}, nil
}

// ParseNotify parses a Stripe webhook event without signature verification.
func (g *StripeGateway) ParseNotify(ctx context.Context, data []byte) (*NotifyResult, error) {
	return g.parseEvent(data)
}

// VerifyAndParseNotify verifies the Stripe webhook signature and parses the event.
// headers must contain "Stripe-Signature" with the value from the HTTP request.
func (g *StripeGateway) VerifyAndParseNotify(ctx context.Context, data []byte, headers map[string]string) (*NotifyResult, error) {
	if g.WebhookSecret == "" {
		return nil, fmt.Errorf("stripe: webhook_secret not configured, cannot verify signature")
	}
	sigHeader := headers["Stripe-Signature"]
	if sigHeader == "" {
		return nil, fmt.Errorf("stripe: missing Stripe-Signature header")
	}

	// Parse the header: t=<timestamp>,v1=<signature>[,v1=...]
	var timestamp string
	var signatures []string
	for _, part := range strings.Split(sigHeader, ",") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "t":
			timestamp = kv[1]
		case "v1":
			signatures = append(signatures, kv[1])
		}
	}
	if timestamp == "" || len(signatures) == 0 {
		return nil, fmt.Errorf("stripe: malformed Stripe-Signature header")
	}

	// Reject events older than 5 minutes to prevent replay attacks.
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("stripe: invalid timestamp in signature: %w", err)
	}
	if abs(time.Now().Unix()-ts) > 300 {
		return nil, fmt.Errorf("stripe: webhook timestamp too old (>%ds)", 300)
	}

	// Compute expected signature: HMAC-SHA256(webhook_secret, "{timestamp}.{body}")
	mac := hmac.New(sha256.New, []byte(g.WebhookSecret))
	mac.Write([]byte(timestamp + "." + string(data)))
	expected := hex.EncodeToString(mac.Sum(nil))

	// Timing-safe comparison against all provided v1 signatures.
	matched := false
	for _, sig := range signatures {
		if subtle.ConstantTimeCompare([]byte(expected), []byte(sig)) == 1 {
			matched = true
			break
		}
	}
	if !matched {
		return nil, fmt.Errorf("stripe: webhook signature verification failed")
	}

	return g.parseEvent(data)
}

// parseEvent is the shared logic for parsing a Stripe webhook event body.
func (g *StripeGateway) parseEvent(data []byte) (*NotifyResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("stripe: empty notification data")
	}

	var event struct {
		Type string `json:"type"`
		Data struct {
			Object struct {
				ID            string            `json:"id"`
				Status        string            `json:"status"`
				Amount        int64             `json:"amount"`
				AmountTotal   int64             `json:"amount_total"`
				PaymentIntent string            `json:"payment_intent"`
				Metadata      map[string]string `json:"metadata"`
			} `json:"object"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("stripe: parse webhook: %w", err)
	}

	orderNo := event.Data.Object.Metadata["order_no"]
	if orderNo == "" {
		orderNo = event.Data.Object.ID
	}

	status := "pending"
	switch event.Type {
	case "checkout.session.completed", "payment_intent.succeeded":
		status = "success"
	case "payment_intent.payment_failed", "checkout.session.expired":
		status = "failed"
	case "charge.refunded":
		status = "refunded"
	}

	amount := float64(event.Data.Object.Amount)
	if amount == 0 {
		amount = float64(event.Data.Object.AmountTotal)
	}
	amount = amount / 100

	return &NotifyResult{
		TradeNo: event.Data.Object.ID,
		OrderNo: orderNo,
		Amount:  amount,
		Status:  status,
	}, nil
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
