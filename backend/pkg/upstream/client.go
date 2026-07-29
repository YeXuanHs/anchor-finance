package upstream

import (
	"fmt"
	"time"

	"anchorfinance/internal/model"
)

// RemoteProduct represents a product fetched from an upstream provider.
type RemoteProduct struct {
	RemoteID      string                 `json:"remote_id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Price         float64                `json:"price"`
	Currency      string                 `json:"currency"`
	BillingCycle  string                 `json:"billing_cycle"`
	Type          string                 `json:"type"`
	Stock         int                    `json:"stock"`
	ConfigOptions map[string]interface{} `json:"config_options,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// ConnectionResult holds the result of a connection test.
type ConnectionResult struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Latency int64  `json:"latency_ms"`
}

// Client defines the interface every upstream provider must implement.
type Client interface {
	// TestConnection sends a lightweight request to verify credentials and reachability.
	TestConnection() (*ConnectionResult, error)
	// FetchProducts retrieves the product catalogue from the upstream.
	FetchProducts() ([]RemoteProduct, error)
}

// NewClient returns the appropriate Client implementation for the given provider.
func NewClient(provider *model.UpstreamProvider) (Client, error) {
	switch provider.Type {
	case "zjmf", "zjmfv3":
		return newZJMFClient(provider), nil
	case "whmcs":
		return newWHMCSClient(provider), nil
	case "v10":
		return newV10Client(provider), nil
	case "custom":
		return newCustomClient(provider), nil
	default:
		return nil, fmt.Errorf("unsupported upstream type: %s", provider.Type)
	}
}

// httpClient is a shared helper with a configured timeout.
func newHTTPTimeout() time.Duration {
	return 15 * time.Second
}
