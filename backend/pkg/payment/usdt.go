package payment

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// USDTGateway implements the Gateway interface for USDT/TRON crypto payment.
// Supports TRC20 and ERC20 networks via TronGrid / Ethereum node APIs.
type USDTGateway struct {
	APIKey      string // TronGrid API key (or Infura/Alchemy for ERC20)
	APISecret   string
	NotifyURL   string
	Network     string // TRC20 or ERC20
	CallbackURL string
	// MasterWallet is the HD wallet root address used to derive per-order deposit addresses.
	// For TRC20, this is a TRON base58 address (T...).
	MasterWallet string
	// APINodeURL is the blockchain node/rpc endpoint (e.g. https://api.trongrid.io).
	APINodeURL string
	// USDTContract is the TRC20/ERC20 contract address for USDT.
	USDTContract string
	// orders stores pending payment orders keyed by order_no.
	orders map[string]*usdtOrder
}

// usdtOrder tracks a pending USDT payment.
type usdtOrder struct {
	OrderNo       string
	TradeNo       string
	Amount        float64
	DepositAddr   string
	ExpectedHash  string
	CreatedAt     time.Time
	ExpiresAt     time.Time
	Status        string // pending, confirmed, expired
	TxHash        string
}

// USDTConfig is the JSON schema stored in payment_gateways.config.
type USDTConfig struct {
	APIKey       string `json:"api_key"`
	APISecret    string `json:"api_secret"`
	NotifyURL    string `json:"notify_url"`
	Network      string `json:"network"`       // TRC20 or ERC20
	CallbackURL  string `json:"callback_url"`
	MasterWallet string `json:"master_wallet"`  // base58 TRON address or EVM address
	APINodeURL   string `json:"api_node_url"`   // e.g. https://api.trongrid.io
	USDTContract string `json:"usdt_contract"`  // TRC20 contract: TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t
}

// NewUSDTGateway creates a USDT gateway from a JSON config string.
func NewUSDTGateway(configJSON string) (*USDTGateway, error) {
	var cfg USDTConfig
	if err := json.Unmarshal([]byte(configJSON), &cfg); err != nil {
		return nil, fmt.Errorf("usdt config parse error: %w", err)
	}
	if cfg.MasterWallet == "" {
		return nil, fmt.Errorf("usdt master_wallet is required")
	}
	if cfg.Network == "" {
		cfg.Network = "TRC20"
	}
	if cfg.APINodeURL == "" {
		if cfg.Network == "TRC20" {
			cfg.APINodeURL = "https://api.trongrid.io"
		}
	}
	if cfg.USDTContract == "" && cfg.Network == "TRC20" {
		cfg.USDTContract = "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"
	}
	return &USDTGateway{
		APIKey:       cfg.APIKey,
		APISecret:    cfg.APISecret,
		NotifyURL:    cfg.NotifyURL,
		Network:      cfg.Network,
		CallbackURL:  cfg.CallbackURL,
		MasterWallet: cfg.MasterWallet,
		APINodeURL:   cfg.APINodeURL,
		USDTContract: cfg.USDTContract,
		orders:       make(map[string]*usdtOrder),
	}, nil
}

func (g *USDTGateway) Name() string { return "usdt" }

// generateDepositAddress derives a per-order deposit address from the master wallet
// using a deterministic hash of order_no + master_wallet.
func (g *USDTGateway) generateDepositAddress(orderNo string) string {
	// Deterministic derivation: H(master_wallet + order_no)
	// In production, use an HD wallet library (e.g., go-ethereum/accounts/hdwallet
	// or tron-go/wallet) to derive child keys from a master seed.
	h := sha256.Sum256([]byte(g.MasterWallet + ":" + orderNo))
	hash := hex.EncodeToString(h[:])

	if g.Network == "TRC20" {
		// TRON address = 41 + last 20 bytes of keccak256(pubkey)
		// Here we use a hash-based address placeholder; production should use proper key derivation.
		return "T" + strings.ToUpper(hash[:33])
	}
	// ERC20 (EVM): 0x + last 20 bytes
	return "0x" + hash[:40]
}

// CreatePayment creates a USDT payment invoice with a unique deposit address and amount.
func (g *USDTGateway) CreatePayment(ctx context.Context, param *PaymentParam) (*PaymentResult, error) {
	if param.Amount <= 0 {
		return nil, fmt.Errorf("usdt: amount must be positive")
	}

	tradeNo := param.TradeNo
	if tradeNo == "" {
		tradeNo = fmt.Sprintf("USDT_%s_%d", param.OrderNo, time.Now().UnixNano())
	}

	depositAddr := g.generateDepositAddress(param.OrderNo)

	order := &usdtOrder{
		OrderNo:     param.OrderNo,
		TradeNo:     tradeNo,
		Amount:      param.Amount,
		DepositAddr: depositAddr,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(30 * time.Minute),
		Status:      "pending",
	}
	g.orders[param.OrderNo] = order

	// Build a QR code URL via a public QR API
	qrData := fmt.Sprintf("%s:%s?amount=%.6f", strings.ToLower(g.Network), depositAddr, param.Amount)
	qrcodeURL := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=300x300&data=%s", qrData)

	rawData := map[string]interface{}{
		"network":       g.Network,
		"amount":        param.Amount,
		"deposit_addr":  depositAddr,
		"contract":      g.USDTContract,
		"expire_at":     order.ExpiresAt.Format(time.RFC3339),
		"qrcode_url":    qrcodeURL,
	}

	return &PaymentResult{
		TradeNo:   tradeNo,
		PayURL:    fmt.Sprintf("usdt://%s/%s", strings.ToLower(g.Network), depositAddr),
		QrcodeURL: qrcodeURL,
		RawData:   rawData,
	}, nil
}

// QueryPayment checks the blockchain for incoming USDT transfers to the deposit address.
func (g *USDTGateway) QueryPayment(ctx context.Context, tradeNo string) (*QueryResult, error) {
	if tradeNo == "" {
		return nil, fmt.Errorf("usdt: trade_no is required")
	}

	// Find the order by trade_no
	var order *usdtOrder
	for _, o := range g.orders {
		if o.TradeNo == tradeNo {
			order = o
			break
		}
	}
	if order == nil {
		return &QueryResult{
			TradeNo: tradeNo,
			Status:  "failed",
		}, nil
	}

	// Check if expired
	if time.Now().After(order.ExpiresAt) && order.Status == "pending" {
		order.Status = "expired"
	}

	if order.Status != "pending" {
		paidAt := time.Now()
		status := order.Status
		if status == "confirmed" {
			status = "success"
		}
		return &QueryResult{
			TradeNo: tradeNo,
			Status:  status,
			Amount:  order.Amount,
			PaidAt:  &paidAt,
		}, nil
	}

	// Query blockchain for TRC20 transfers to the deposit address
	if g.Network == "TRC20" && g.APINodeURL != "" {
		confirmed, txHash, err := g.queryTRC20Transfers(ctx, order)
		if err != nil {
			return &QueryResult{
				TradeNo: tradeNo,
				Status:  "pending",
			}, nil
		}
		if confirmed {
			order.Status = "confirmed"
			order.TxHash = txHash
			paidAt := time.Now()
			return &QueryResult{
				TradeNo: tradeNo,
				Status:  "success",
				Amount:  order.Amount,
				PaidAt:  &paidAt,
			}, nil
		}
	}

	return &QueryResult{
		TradeNo: tradeNo,
		Status:  "pending",
	}, nil
}

// queryTRC20Transfers calls TronGrid API to check for TRC20 transfers to the deposit address.
func (g *USDTGateway) queryTRC20Transfers(ctx context.Context, order *usdtOrder) (bool, string, error) {
	url := fmt.Sprintf(
		"%s/v1/accounts/%s/transactions/trc20?limit=20&contract_address=%s",
		g.APINodeURL, order.DepositAddr, g.USDTContract,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, "", err
	}
	if g.APIKey != "" {
		req.Header.Set("TRON-PRO-API-KEY", g.APIKey)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	var result struct {
		Data []struct {
			TransactionID string `json:"transaction_id"`
			TokenInfo     struct {
				Address string `json:"address"`
			} `json:"token_info"`
			From   string `json:"from"`
			To     string `json:"to"`
			Type   string `json:"type"`
			Value  string `json:"value"` // in smallest unit (6 decimals for USDT)
		} `json:"data"`
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, "", err
	}

	for _, tx := range result.Data {
		// USDT has 6 decimal places
		var value float64
		fmt.Sscanf(tx.Value, "%f", &value)
		value = value / 1_000_000

		if tx.To == order.DepositAddr && value >= order.Amount {
			return true, tx.TransactionID, nil
		}
	}

	return false, "", nil
}

// Refund is not supported for direct crypto payments.
func (g *USDTGateway) Refund(ctx context.Context, param *RefundParam) (*RefundResult, error) {
	return &RefundResult{
		RefundNo: param.RefundNo,
		Status:   "unsupported",
	}, fmt.Errorf("usdt: crypto payments require manual refund processing")
}

// ParseNotify parses a USDT payment callback notification.
func (g *USDTGateway) ParseNotify(ctx context.Context, data []byte) (*NotifyResult, error) {
	var notify struct {
		OrderID string  `json:"order_id"`
		Status  string  `json:"status"`
		Amount  float64 `json:"amount"`
		TxHash  string  `json:"tx_hash"`
		Network string  `json:"network"`
	}
	if err := json.Unmarshal(data, &notify); err != nil {
		return nil, fmt.Errorf("parse usdt notify: %w", err)
	}

	status := "pending"
	if notify.Status == "confirmed" || notify.Status == "success" {
		status = "success"
	} else if notify.Status == "expired" || notify.Status == "failed" {
		status = "failed"
	}

	return &NotifyResult{
		TradeNo: notify.OrderID,
		OrderNo: notify.OrderID,
		Amount:  notify.Amount,
		Status:  status,
	}, nil
}
