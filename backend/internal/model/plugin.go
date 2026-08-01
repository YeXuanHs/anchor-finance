package model

import "gorm.io/gorm"

// 插件类型常量
const (
	PluginTypeMail          = "mail"           // 邮件插件
	PluginTypeSMS           = "sms"            // 短信插件
	PluginTypeCertification = "certification"  // 实名认证
	PluginTypeGateway       = "gateway"        // 支付网关
	PluginTypeOAuth         = "oauth"          // OAuth登录
	PluginTypeServer        = "server"         // 服务器模块（自动开通）
	PluginTypeAddon         = "addon"          // 扩展插件
)

// 插件类型标签
var PluginTypeLabels = map[string]string{
	PluginTypeMail:          "邮件",
	PluginTypeSMS:           "短信",
	PluginTypeCertification: "实名认证",
	PluginTypeGateway:       "支付网关",
	PluginTypeOAuth:         "OAuth登录",
	PluginTypeServer:        "服务器模块",
	PluginTypeAddon:         "扩展插件",
}

// Plugin 插件模型
type Plugin struct {
	gorm.Model
	Name        string `gorm:"type:varchar(64);uniqueIndex;not null" json:"name"`        // 插件标识名
	Title       string `gorm:"type:varchar(64);not null" json:"title"`                   // 显示名称
	Type        string `gorm:"type:varchar(32);not null;index" json:"type"`               // 插件类型
	Description string `gorm:"type:varchar(512)" json:"description"`                     // 描述
	Author      string `gorm:"type:varchar(64)" json:"author"`                           // 作者
	Version     string `gorm:"type:varchar(32)" json:"version"`                          // 版本号
	HelpURL     string `gorm:"type:varchar(256)" json:"help_url"`                        // 帮助文档URL
	Config      string `gorm:"type:json" json:"config"`                                 // 配置JSON
	IsSystem    bool   `gorm:"default:false" json:"is_system"`                           // 是否系统内置
	IsEnabled   bool   `gorm:"default:false;index" json:"is_enabled"`                    // 是否启用
	SortOrder   int    `gorm:"default:0" json:"sort_order"`                              // 排序
	Module      string `gorm:"type:varchar(32)" json:"module"`                           // 模块目录名
	// 服务器模块专用字段
	MaxAccounts int    `gorm:"default:0" json:"max_accounts"`                            // 最大账户数
	ServerGroup string `gorm:"type:varchar(64)" json:"server_group"`                     // 服务器分组
}

// TableName 表名
func (Plugin) TableName() string {
	return "plugins"
}

// ServerModule 服务器模块（自动开通接口）
type ServerModule struct {
	gorm.Model
	Name         string `gorm:"type:varchar(64);uniqueIndex;not null" json:"name"`       // 模块标识名
	Title        string `gorm:"type:varchar(64);not null" json:"title"`                  // 显示名称
	Description  string `gorm:"type:varchar(512)" json:"description"`                    // 描述
	Module       string `gorm:"type:varchar(64);not null;index" json:"module"`           // 模块类型（对应plugin name）
	Hostname     string `gorm:"type:varchar(128)" json:"hostname"`                       // 主机名/IP
	Port         int    `gorm:"default" json:"port"`                                     // 端口
	Username     string `gorm:"type:varchar(64)" json:"username"`                        // 用户名
	Password     string `gorm:"type:varchar(256)" json:"password"`                       // 密码（加密存储）
	AccessHash   string `gorm:"type:text" json:"access_hash"`                            // 访问哈希
	Secure       bool   `gorm:"default:false" json:"secure"`                             // 是否使用SSL
	MaxAccounts  int    `gorm:"default:100" json:"max_accounts"`                         // 最大账户数
	CurrentCount int    `gorm:"default:0" json:"current_count"`                          // 当前账户数
	GroupId      uint   `gorm:"index" json:"group_id"`                                   // 分组ID
	GroupName    string `gorm:"type:varchar(64)" json:"group_name"`                      // 分组名称
	IsEnabled    bool   `gorm:"default:true;index" json:"is_enabled"`                    // 是否启用
	Status       string `gorm:"type:varchar(32);default:'active'" json:"status"`         // 状态
	Config       string `gorm:"type:json" json:"config"`                                // 模块配置
	SortOrder    int    `gorm:"default:0" json:"sort_order"`                             // 排序
}

// TableName 表名
func (ServerModule) TableName() string {
	return "server_modules"
}

// ServerGroup 服务器分组
type ServerGroup struct {
	gorm.Model
	Name        string `gorm:"type:varchar(64);uniqueIndex;not null" json:"name"`        // 分组名称
	Description string `gorm:"type:varchar(256)" json:"description"`                     // 描述
	ModuleType  string `gorm:"type:varchar(64)" json:"module_type"`                      // 模块类型
	IsEnabled   bool   `gorm:"default:true" json:"is_enabled"`                           // 是否启用
	SortOrder   int    `gorm:"default:0" json:"sort_order"`                              // 排序
}

// TableName 表名
func (ServerGroup) TableName() string {
	return "server_groups"
}
