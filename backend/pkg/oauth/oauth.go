package oauth

// Provider defines the interface for OAuth providers.
type Provider interface {
	Name() string
	GetAuthURL(state string) string
	GetUserInfo(code string) (*UserInfo, error)
}

// UserInfo contains the user information returned by an OAuth provider.
type UserInfo struct {
	Provider string                 `json:"provider"`
	OpenID   string                 `json:"openid"`
	UnionID  string                 `json:"unionid"`
	Username string                 `json:"username"`
	Email    string                 `json:"email"`
	Avatar   string                 `json:"avatar"`
	RawData  map[string]interface{} `json:"raw_data"`
}

// Global provider registry.
var providers = make(map[string]Provider)

// Register adds an OAuth provider to the registry.
func Register(p Provider) {
	providers[p.Name()] = p
}

// Get retrieves an OAuth provider by name.
func Get(name string) (Provider, bool) {
	p, ok := providers[name]
	return p, ok
}

// GetAll returns all registered OAuth providers.
func GetAll() map[string]Provider {
	return providers
}
