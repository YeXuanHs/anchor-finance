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
	GroupID       string                 `json:"group_id,omitempty"`
	GroupName     string                 `json:"group_name,omitempty"`
	ConfigOptions map[string]interface{} `json:"config_options,omitempty"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// RemoteProductGroup represents a product group from upstream.
type RemoteProductGroup struct {
	GroupID      string `json:"group_id"`
	Name         string `json:"name"`
	ParentID     string `json:"parent_id,omitempty"`
	ProductCount int    `json:"product_count"`
}

// UpstreamProductsResult contains products with group structure.
type UpstreamProductsResult struct {
	Products []RemoteProduct      `json:"products"`
	Groups   []RemoteProductGroup `json:"groups"`
	Currency string               `json:"currency"`
}

// ConnectionResult holds the result of a connection test.
type ConnectionResult struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Latency int64  `json:"latency_ms"`
}

// Client defines the interface every upstream provider must implement.
type Client interface {
	TestConnection() (*ConnectionResult, error)
	FetchProducts() ([]RemoteProduct, error)
	FetchProductsWithGroups() (*UpstreamProductsResult, error)
	FetchProductsByGroup(groupID string) ([]RemoteProduct, error)
}

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

func newHTTPTimeout() time.Duration {
	return 15 * time.Second
}
