package server

import "context"

// ServerStatus 服务器状态
type ServerStatus struct {
	Status string `json:"status"` // on/off/unknown
	Desc   string `json:"desc"`   // 状态描述
}

// AccountInfo 账户信息
type AccountInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	OS       string `json:"os"`
}

// ServerModule 服务器模块接口（自动开通）
type ServerModule interface {
	// Name 模块名称
	Name() string

	// Title 模块显示标题
	Title() string

	// TestConnection 测试连接
	TestConnection(ctx context.Context, params *ConnectionParams) error

	// CreateAccount 创建账户（开通服务）
	CreateAccount(ctx context.Context, params *CreateAccountParams) (*AccountInfo, error)

	// TerminateAccount 终止账户（删除服务）
	TerminateAccount(ctx context.Context, params *TerminateAccountParams) error

	// SuspendAccount 暂停账户
	SuspendAccount(ctx context.Context, params *SuspendAccountParams) error

	// UnsuspendAccount 取消暂停
	UnsuspendAccount(ctx context.Context, params *UnsuspendAccountParams) error

	// GetStatus 获取服务状态
	GetStatus(ctx context.Context, params *StatusParams) (*ServerStatus, error)

	// GetClientArea 获取客户端区域信息
	GetClientArea(ctx context.Context, params *ClientAreaParams) (map[string]interface{}, error)

	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, params *ChangePasswordParams) error

	// GetConfigOptions 获取配置选项
	GetConfigOptions() []ConfigOption
}

// ConnectionParams 连接参数
type ConnectionParams struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Config   map[string]interface{} `json:"config"`
}

// CreateAccountParams 创建账户参数
type CreateAccountParams struct {
	HostId      uint                   `json:"host_id"`
	Host        string                 `json:"host"`
	Port        int                    `json:"port"`
	Username    string                 `json:"username"`
	Password    string                 `json:"password"`
	ConfigOptions map[string]string    `json:"config_options"`
	Params      map[string]interface{} `json:"params"`
}

// TerminateAccountParams 终止账户参数
type TerminateAccountParams struct {
	HostId   uint   `json:"host_id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Config   map[string]interface{} `json:"config"`
}

// SuspendAccountParams 暂停账户参数
type SuspendAccountParams struct {
	HostId   uint   `json:"host_id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Reason   string `json:"reason"`
	Config   map[string]interface{} `json:"config"`
}

// UnsuspendAccountParams 取消暂停参数
type UnsuspendAccountParams struct {
	HostId   uint   `json:"host_id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Config   map[string]interface{} `json:"config"`
}

// StatusParams 状态查询参数
type StatusParams struct {
	HostId   uint   `json:"host_id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Config   map[string]interface{} `json:"config"`
}

// ClientAreaParams 客户端区域参数
type ClientAreaParams struct {
	HostId   uint   `json:"host_id"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Config   map[string]interface{} `json:"config"`
}

// ChangePasswordParams 修改密码参数
type ChangePasswordParams struct {
	HostId      uint   `json:"host_id"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	Config      map[string]interface{} `json:"config"`
}

// ConfigOption 配置选项
type ConfigOption struct {
	Type        string `json:"type"`        // text/password/select/checkbox
	Name        string `json:"name"`        // 显示名称
	Key         string `json:"key"`         // 配置键名
	Placeholder string `json:"placeholder"` // 占位符
	Description string `json:"description"` // 描述
	Options     []Option `json:"options"`   // 选项列表（type=select时使用）
	Required    bool   `json:"required"`    // 是否必填
	Default     string `json:"default"`     // 默认值
}

// Option 选项
type Option struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ModuleConstructor 模块构造函数
type ModuleConstructor func(config map[string]interface{}) (ServerModule, error)
