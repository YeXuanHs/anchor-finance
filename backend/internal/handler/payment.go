package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/payment"
)

// PaymentHandler handles HTTP requests for payment operations.
type PaymentHandler struct {
	paymentSvc *service.PaymentService
}

// NewPaymentHandler creates a new PaymentHandler.
func NewPaymentHandler(paymentSvc *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentSvc: paymentSvc}
}

// gatewayInfo represents a payment gateway in API responses.
type gatewayInfo struct {
	Name string `json:"name"`
}

// createPaymentRequest is the request body for CreatePayment.
type createPaymentRequest struct {
	UserID    string  `json:"user_id"`
	InvoiceID string  `json:"invoice_id"`
	Gateway   string  `json:"gateway"`
	Amount    float64 `json:"amount"`
	Subject   string  `json:"subject"`
}

// createPaymentResponse is the response body for CreatePayment.
type createPaymentResponse struct {
	TradeNo   string `json:"trade_no"`
	PayURL    string `json:"pay_url,omitempty"`
	QrcodeURL string `json:"qrcode_url,omitempty"`
}

// paymentResponse is a generic payment status response.
type paymentResponse struct {
	TradeNo string  `json:"trade_no"`
	Status  string  `json:"status"`
	Amount  float64 `json:"amount"`
}

// errorResponse is a JSON error response.
type errorResponse struct {
	Error string `json:"error"`
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error response.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

// GetGateways handles GET /payments/gateways
// Returns a list of all available payment gateways.
func (h *PaymentHandler) GetGateways(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	all := payment.GetAll()
	gateways := make([]gatewayInfo, 0, len(all))
	for name := range all {
		gateways = append(gateways, gatewayInfo{Name: name})
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"gateways": gateways,
	})
}

// CreatePayment handles POST /payments/create
// Creates a new payment through the specified gateway.
func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req createPaymentRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read request body")
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.UserID == "" {
		writeError(w, http.StatusBadRequest, "user_id is required")
		return
	}
	if req.InvoiceID == "" {
		writeError(w, http.StatusBadRequest, "invoice_id is required")
		return
	}
	if req.Gateway == "" {
		writeError(w, http.StatusBadRequest, "gateway is required")
		return
	}

	result, err := h.paymentSvc.CreatePayment(r.Context(), req.UserID, req.InvoiceID, req.Gateway, &service.CreatePaymentInput{
		Amount:  req.Amount,
		Subject: req.Subject,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, createPaymentResponse{
		TradeNo:   result.TradeNo,
		PayURL:    result.PayURL,
		QrcodeURL: result.QrcodeURL,
	})
}

// AlipayNotify handles POST /payments/alipay/notify
// Receives async payment notification from Alipay.
func (h *PaymentHandler) AlipayNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read notification body")
		return
	}
	defer r.Body.Close()

	// Convert body to map[string]string
	data := map[string]string{
		"body": string(body),
	}

	if err := h.paymentSvc.HandleNotify(r.Context(), "alipay", data); err != nil {
		// Alipay expects "failure" text on error.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("failure"))
		return
	}

	// Alipay expects "success" text on success.
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

// WechatNotify handles POST /payments/wechat/notify
// Receives async payment notification from WeChat Pay.
func (h *PaymentHandler) WechatNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "failed to read notification body")
		return
	}
	defer r.Body.Close()

	// Convert body to map[string]string
	data := map[string]string{
		"body": string(body),
	}

	if err := h.paymentSvc.HandleNotify(r.Context(), "wechat", data); err != nil {
		// WePay expects XML error response.
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[FAIL]]></return_msg></xml>`))
		return
	}

	// WePay expects XML success response.
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>`))
}

// ReturnURL handles GET /payments/return
// Called when the user is redirected back from the payment page.
func (h *PaymentHandler) ReturnURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	tradeNo := r.URL.Query().Get("trade_no")
	if tradeNo == "" {
		writeError(w, http.StatusBadRequest, "trade_no is required")
		return
	}

	// Try to get transaction details via the service's store
	// The actual status update happens via HandleNotify (async webhook)
	// This endpoint just confirms the user returned from the payment page
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"trade_no": tradeNo,
		"status":   "return",
		"message":  "payment return received, order status will be updated shortly",
	})
}
