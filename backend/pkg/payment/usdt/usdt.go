package usdt

import (
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"anchorfinance/pkg/payment"
)

const (
	// TronGridAPI is the public TRON blockchain API for querying transactions.
	TronGridAPI = "https://api.trongrid.io"
	// TRC20ContractAddress is the USDT TRC-20 contract address on TRON mainnet.
	TRC20ContractAddress = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	// DefaultConfirmations is the number of block confirmations to wait for.
	DefaultConfirmations = 19
)

func init() {
	payment.Register(&Gateway{})
}

// AddressProvider generates or retrieves a unique TRC-20 deposit address for a given order.
// Implementations can use HD wallets, pre-generated address pools, or any other strategy.
type AddressProvider interface {
	// GetDepositAddress returns a TRON base58 address and its hex-encoded private key
	// for the given order number. The address must be unique per order to allow
	// reliable payment detection.
	GetDepositAddress(orderNo string) (address string, err error)
	// GetAddressPrivateKey returns the hex private key for a known address.
	GetAddressPrivateKey(address string) (privateKey string, err error)
}

// BlockchainScanner queries the TRON blockchain for incoming USDT transfers.
// The default implementation uses TronGrid's REST API.
type BlockchainScanner interface {
	// GetTRC20Transfers returns USDT transfers sent to the given address after the specified timestamp.
	GetTRC20Transfers(ctx context.Context, address string, since time.Time) ([]TRC20Transfer, error)
}

// TRC20Transfer represents a single TRC-20 token transfer event.
type TRC20Transfer struct {
	TxID        string
	From        string
	To          string
	Amount      float64 // Amount in USDT (not sun)
	BlockTime   time.Time
	Confirmed   bool
}

// Gateway implements the payment.Gateway interface for USDT (TRC-20) payments.
//
// This gateway generates a unique TRON deposit address per order, monitors the
// blockchain for incoming USDT transfers, and verifies the payment amount.
//
// Flow:
//  1. CreatePayment: generate a deposit address for the order
//  2. Client sends USDT to the deposit address
//  3. ParseNotify: verify the incoming transfer matches the expected amount
//  4. QueryPayment: scan the blockchain for transfers to the deposit address
type Gateway struct {
	AddressProvider    AddressProvider
	Scanner            BlockchainScanner
	APIKey             string // TronGrid API key for higher rate limits (optional)
	Confirmations      int    // Required block confirmations (default 19)
	Timeout            time.Duration
	mu                 sync.RWMutex
	orderAddresses     map[string]string // orderNo -> deposit address
	orderAmounts       map[string]float64 // orderNo -> expected amount
}

// Name returns the gateway identifier.
func (g *Gateway) Name() string {
	return "usdt"
}

// httpClient returns the configured HTTP client.
func (g *Gateway) httpClient() *http.Client {
	timeout := g.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	return &http.Client{Timeout: timeout}
}

// confirmations returns the configured confirmation count or the default.
func (g *Gateway) confirmations() int {
	if g.Confirmations > 0 {
		return g.Confirmations
	}
	return DefaultConfirmations
}

// getScanner returns the configured scanner or a default TronGrid scanner.
func (g *Gateway) getScanner() BlockchainScanner {
	if g.Scanner != nil {
		return g.Scanner
	}
	return &TronGridScanner{
		APIKey:     g.APIKey,
		HTTPClient: g.httpClient(),
	}
}

// CreatePayment generates a unique TRC-20 deposit address for the order.
func (g *Gateway) CreatePayment(ctx context.Context, param *payment.PaymentParam) (*payment.PaymentResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("usdt: amount must be positive")
	}
	if g.AddressProvider == nil {
		return nil, fmt.Errorf("usdt: address provider not configured")
	}

	address, err := g.AddressProvider.GetDepositAddress(param.OrderNo)
	if err != nil {
		return nil, fmt.Errorf("usdt: failed to get deposit address: %w", err)
	}
	if address == "" {
		return nil, fmt.Errorf("usdt: address provider returned empty address")
	}

	// Store the mapping for later verification
	g.mu.Lock()
	if g.orderAddresses == nil {
		g.orderAddresses = make(map[string]string)
	}
	if g.orderAmounts == nil {
		g.orderAmounts = make(map[string]float64)
	}
	g.orderAddresses[param.OrderNo] = address
	g.orderAmounts[param.OrderNo] = param.Amount
	g.mu.Unlock()

	return &payment.PaymentResult{
		TradeNo: param.OrderNo,
		PayURL:  fmt.Sprintf("tron:%s?amount=%.6f", address, param.Amount),
		RawData: map[string]interface{}{
			"gateway":    "usdt",
			"network":    "TRC20",
			"address":    address,
			"amount":     param.Amount,
			"order_no":   param.OrderNo,
			"contract":   TRC20ContractAddress,
			"created_at": time.Now().Format(time.RFC3339),
		},
	}, nil
}

// QueryPayment scans the blockchain for incoming USDT transfers to the deposit address.
func (g *Gateway) QueryPayment(ctx context.Context, tradeNo string) (*payment.QueryResult, error) {
	if tradeNo == "" {
		return nil, fmt.Errorf("usdt: trade_no is required")
	}

	g.mu.RLock()
	address := g.orderAddresses[tradeNo]
	expectedAmount := g.orderAmounts[tradeNo]
	g.mu.RUnlock()

	if address == "" {
		return &payment.QueryResult{
			TradeNo: tradeNo,
			Status:  "pending",
		}, nil
	}

	scanner := g.getScanner()
	transfers, err := scanner.GetTRC20Transfers(ctx, address, time.Time{})
	if err != nil {
		return nil, fmt.Errorf("usdt: blockchain scan failed: %w", err)
	}

	for _, t := range transfers {
		if t.Confirmed && t.Amount >= expectedAmount {
			now := t.BlockTime
			return &payment.QueryResult{
				TradeNo: tradeNo,
				Status:  "success",
				Amount:  t.Amount,
				PaidAt:  &now,
			}, nil
		}
	}

	return &payment.QueryResult{
		TradeNo: tradeNo,
		Status:  "pending",
	}, nil
}

// Refund is not supported for USDT payments (irreversible blockchain transactions).
func (g *Gateway) Refund(ctx context.Context, param *payment.RefundParam) (*payment.RefundResult, error) {
	return nil, fmt.Errorf("usdt: refund is not supported for cryptocurrency payments")
}

// ParseNotify parses a USDT payment notification.
//
// The notification data is expected to be a JSON-encoded NotifyPayload.
// This can be triggered by a blockchain monitoring service or webhook.
func (g *Gateway) ParseNotify(ctx context.Context, data []byte) (*payment.NotifyResult, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("usdt: empty notification data")
	}

	// Parse the notification payload
	payload, err := parseNotifyPayload(data)
	if err != nil {
		return nil, fmt.Errorf("usdt: invalid notification: %w", err)
	}

	// Verify the deposit address matches our records
	g.mu.RLock()
	expectedAddr := g.orderAddresses[payload.OrderNo]
	expectedAmount := g.orderAmounts[payload.OrderNo]
	g.mu.RUnlock()

	if expectedAddr == "" {
		return nil, fmt.Errorf("usdt: unknown order %s", payload.OrderNo)
	}
	if payload.ToAddress != expectedAddr {
		return nil, fmt.Errorf("usdt: address mismatch (expected %s, got %s)", expectedAddr, payload.ToAddress)
	}
	if payload.Amount < expectedAmount {
		return nil, fmt.Errorf("usdt: insufficient amount (expected %.6f, got %.6f)", expectedAmount, payload.Amount)
	}

	status := "pending"
	if payload.Confirmed {
		status = "success"
	}

	return &payment.NotifyResult{
		TradeNo: payload.OrderNo,
		OrderNo: payload.OrderNo,
		Amount:  payload.Amount,
		Status:  status,
		Sign:    payload.TxID,
	}, nil
}

// NotifyPayload is the expected structure for USDT payment notifications.
type NotifyPayload struct {
	OrderNo   string  `json:"order_no"`
	TxID      string  `json:"tx_id"`
	FromAddr  string  `json:"from_address"`
	ToAddress string  `json:"to_address"`
	Amount    float64 `json:"amount"`
	Confirmed bool    `json:"confirmed"`
	BlockTime int64   `json:"block_time"`
}

func parseNotifyPayload(data []byte) (*NotifyPayload, error) {
	var p NotifyPayload
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	if p.OrderNo == "" || p.TxID == "" {
		return nil, fmt.Errorf("missing required fields (order_no, tx_id)")
	}
	return &p, nil
}

// TronGridScanner implements BlockchainScanner using the TronGrid REST API.
type TronGridScanner struct {
	APIKey     string
	HTTPClient *http.Client
}

// GetTRC20Transfers queries TronGrid for USDT TRC-20 transfers to the given address.
func (s *TronGridScanner) GetTRC20Transfers(ctx context.Context, address string, since time.Time) ([]TRC20Transfer, error) {
	client := s.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	// Build the TronGrid API URL for TRC-20 transfers
	apiURL := fmt.Sprintf(
		"%s/v1/accounts/%s/transactions/trc20?contract_address=%s&limit=200&order_by=block_timestamp,desc",
		TronGridAPI, address, TRC20ContractAddress,
	)
	if !since.IsZero() {
		apiURL += fmt.Sprintf("&min_timestamp=%d", since.UnixMilli())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	if s.APIKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", s.APIKey)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("trongrid: HTTP %d", resp.StatusCode)
	}

	var result trongridResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("trongrid: failed to decode response: %w", err)
	}

	transfers := make([]TRC20Transfer, 0, len(result.Data))
	for _, tx := range result.Data {
		amount := parseFloat(tx.Value, tx.TokenInfo.Decimals)
		blockTime := time.UnixMilli(tx.BlockTimestamp)

		transfers = append(transfers, TRC20Transfer{
			TxID:      tx.TransactionID,
			From:      tx.From,
			To:        tx.To,
			Amount:    amount,
			BlockTime: blockTime,
			Confirmed: true, // TronGrid only returns confirmed transactions
		})
	}

	return transfers, nil
}

// GenerateDepositAddress creates a random TRON address for demo/testing purposes.
// In production, use a proper HD wallet or address pool implementation.
func GenerateDepositAddress() (address string, privateKeyHex string, err error) {
	privKey := make([]byte, 32)
	if _, err := rand.Read(privKey); err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}
	privateKeyHex = hex.EncodeToString(privKey)

	// Generate a deterministic-looking TRON address from the private key hash.
	// In production, derive the address using secp256k1 and TRON's address encoding.
	addrHash := md5Hash(privKey)
	address = "T" + strings.ToUpper(addrHash[:33])

	return address, privateKeyHex, nil
}

// SimpleAddressProvider is a basic AddressProvider that generates random addresses.
// Suitable for testing; in production, use an HD wallet-based provider.
type SimpleAddressProvider struct {
	mu        sync.Mutex
	addresses map[string]string // orderNo -> address
}

// NewSimpleAddressProvider creates a new SimpleAddressProvider.
func NewSimpleAddressProvider() *SimpleAddressProvider {
	return &SimpleAddressProvider{
		addresses: make(map[string]string),
	}
}

// GetDepositAddress returns a unique address for the order, generating one if needed.
func (p *SimpleAddressProvider) GetDepositAddress(orderNo string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if addr, ok := p.addresses[orderNo]; ok {
		return addr, nil
	}

	addr, _, err := GenerateDepositAddress()
	if err != nil {
		return "", err
	}
	p.addresses[orderNo] = addr
	return addr, nil
}

// GetAddressPrivateKey returns the private key for the given address.
func (p *SimpleAddressProvider) GetAddressPrivateKey(address string) (string, error) {
	return "", fmt.Errorf("simple address provider does not store private keys")
}

// trongridResponse represents the TronGrid API response for TRC-20 transfers.
type trongridResponse struct {
	Data []trongridTRC20Tx `json:"data"`
}

type trongridTRC20Tx struct {
	TransactionID string           `json:"transaction_id"`
	From          string           `json:"from"`
	To            string           `json:"to"`
	Value         string           `json:"value"`
	BlockTimestamp int64           `json:"block_timestamp"`
	TokenInfo     trongridTokenInfo `json:"token_info"`
}

type trongridTokenInfo struct {
	Address  string `json:"address"`
	Symbol   string `json:"symbol"`
	Decimals int    `json:"decimals"`
}

// parseFloat parses a string amount considering token decimals.
// USDT has 6 decimals on TRC-20.
func parseFloat(value string, decimals int) float64 {
	if value == "" {
		return 0
	}
	// Handle very large numbers by padding with leading zeros if needed
	for len(value) <= decimals {
		value = "0" + value
	}
	integerPart := value[:len(value)-decimals]
	fractionalPart := value[len(value)-decimals:]

	f, _ := strconv.ParseFloat(integerPart+"."+fractionalPart, 64)
	return f
}

// md5Hash returns the hex-encoded MD5 hash of the input.
func md5Hash(data []byte) string {
	h := md5.Sum(data)
	return hex.EncodeToString(h[:])
}
