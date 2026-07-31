package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemConfig 系统配置模型（统一配置表）
type SystemConfig struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Group     string         `gorm:"index;size:50;not null" json:"group"`      // 配置分组
	Key       string         `gorm:"uniqueIndex;size:100;not null" json:"key"` // 配置键
	Value     string         `gorm:"type:text" json:"value"`                   // 配置值
	Type      string         `gorm:"size:20;default:string" json:"type"`       // 值类型: string, int, bool, json
	Name      string         `gorm:"size:100" json:"name"`                     // 配置名称
	Options   string         `gorm:"type:text" json:"options"`                 // 选项(JSON)
	SortOrder int            `gorm:"default:0" json:"sort_order"`
	Status    bool           `gorm:"default:true" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}

// ConfigGroup 配置分组常量
const (
	ConfigGroupGeneral     = "general"      // 基本设置
	ConfigGroupSecurity    = "security"     // 安全设置
	ConfigGroupLogin       = "login"        // 登录注册
	ConfigGroupCaptcha     = "captcha"      // 验证码
	ConfigGroupPayment     = "payment"      // 支付配置
	ConfigGroupEmail       = "email"        // 邮件配置
	ConfigGroupSMS         = "sms"          // 短信配置
	ConfigGroupNotification = "notification" // 通知配置
	ConfigGroupInvoice     = "invoice"      // 发票配置
	ConfigGroupAffiliate   = "affiliate"    // 代理配置
	ConfigGroupCredit      = "credit"       // 信用额度
	ConfigGroupRecharge    = "recharge"     // 充值配置
	ConfigGroupDisplay     = "display"      // 显示配置
	ConfigGroupTemplate    = "template"     // 模板配置
	ConfigGroupSEO         = "seo"          // SEO配置
	ConfigGroupMaintenance = "maintenance"  // 维护模式
	ConfigGroupOAuth       = "oauth"        // OAuth配置
	ConfigGroupUpstream    = "upstream"     // 上游对接
	ConfigGroupAdvanced    = "advanced"     // 高级设置
)

// DefaultSystemConfigs 默认系统配置
var DefaultSystemConfigs = []SystemConfig{
	// ========== 基本设置 ==========
	{Group: ConfigGroupGeneral, Key: "company_name", Value: "", Type: "string", Name: "公司名称", SortOrder: 1},
	{Group: ConfigGroupGeneral, Key: "company_email", Value: "", Type: "string", Name: "公司邮箱", SortOrder: 2},
	{Group: ConfigGroupGeneral, Key: "company_phone", Value: "", Type: "string", Name: "公司电话", SortOrder: 3},
	{Group: ConfigGroupGeneral, Key: "company_address", Value: "", Type: "string", Name: "公司地址", SortOrder: 4},
	{Group: ConfigGroupGeneral, Key: "company_profile", Value: "", Type: "text", Name: "公司简介", SortOrder: 5},
	{Group: ConfigGroupGeneral, Key: "domain", Value: "", Type: "string", Name: "网站域名", SortOrder: 6},
	{Group: ConfigGroupGeneral, Key: "system_url", Value: "", Type: "string", Name: "系统链接", SortOrder: 7},
	{Group: ConfigGroupGeneral, Key: "record_no", Value: "", Type: "string", Name: "备案号", SortOrder: 8},
	{Group: ConfigGroupGeneral, Key: "map", Value: "", Type: "string", Name: "坐标", SortOrder: 9},

	// Logo配置
	{Group: ConfigGroupGeneral, Key: "logo_url", Value: "", Type: "string", Name: "默认Logo", SortOrder: 10},
	{Group: ConfigGroupGeneral, Key: "logo_url_home", Value: "", Type: "string", Name: "前台Logo", SortOrder: 11},
	{Group: ConfigGroupGeneral, Key: "logo_url_bill", Value: "", Type: "string", Name: "账单Logo", SortOrder: 12},
	{Group: ConfigGroupGeneral, Key: "logo_url_admin", Value: "", Type: "string", Name: "后台Logo", SortOrder: 13},
	{Group: ConfigGroupGeneral, Key: "www_logo", Value: "", Type: "string", Name: "官网Logo", SortOrder: 14},
	{Group: ConfigGroupGeneral, Key: "favicon_url", Value: "", Type: "string", Name: "网站图标", SortOrder: 15},

	// 法律条款
	{Group: ConfigGroupGeneral, Key: "server_clause_url", Value: "", Type: "string", Name: "服务条款地址", SortOrder: 20},
	{Group: ConfigGroupGeneral, Key: "privacy_clause_url", Value: "", Type: "string", Name: "隐私条款地址", SortOrder: 21},
	{Group: ConfigGroupGeneral, Key: "cancellation_time", Value: "7", Type: "int", Name: "注销时间(天)", SortOrder: 22},

	// ========== 安全设置 ==========
	{Group: ConfigGroupSecurity, Key: "required_pwstrength", Value: "alpha_num", Type: "string", Name: "密码强度要求", SortOrder: 1,
		Options: `[{"label":"无要求","value":"none"},{"label":"字母数字","value":"alpha_num"},{"label":"大小写字母数字","value":"capital_alpha_num"},{"label":"大小写字母数字符号","value":"alpha_num_special"}]`},
	{Group: ConfigGroupSecurity, Key: "invalid_logins_banlength", Value: "30", Type: "int", Name: "登录失败封禁时长(分钟)", SortOrder: 2},
	{Group: ConfigGroupSecurity, Key: "login_error_max_num", Value: "5", Type: "int", Name: "登录错误次数限制", SortOrder: 3},
	{Group: ConfigGroupSecurity, Key: "login_error_switch", Value: "1", Type: "bool", Name: "登录错误限制开关", SortOrder: 4},
	{Group: ConfigGroupSecurity, Key: "home_ip_check", Value: "0", Type: "bool", Name: "前台登录IP检测", SortOrder: 5},
	{Group: ConfigGroupSecurity, Key: "admin_ip_check", Value: "0", Type: "bool", Name: "后台登录IP检测", SortOrder: 6},
	{Group: ConfigGroupSecurity, Key: "second_verify_action_home", Value: "", Type: "string", Name: "前台二次验证操作", SortOrder: 7},
	{Group: ConfigGroupSecurity, Key: "second_verify_action_admin", Value: "", Type: "string", Name: "后台二次验证操作", SortOrder: 8},
	{Group: ConfigGroupSecurity, Key: "second_verify_action_home_type", Value: "", Type: "string", Name: "前台二次验证类型", SortOrder: 9},
	{Group: ConfigGroupSecurity, Key: "admin_path", Value: "admin", Type: "string", Name: "后台路径", SortOrder: 10},

	// ========== 登录注册 ==========
	{Group: ConfigGroupLogin, Key: "allow_phone", Value: "1", Type: "bool", Name: "允许手机登录", SortOrder: 1},
	{Group: ConfigGroupLogin, Key: "allow_email", Value: "1", Type: "bool", Name: "允许邮箱登录", SortOrder: 2},
	{Group: ConfigGroupLogin, Key: "allow_id", Value: "0", Type: "bool", Name: "允许ID登录", SortOrder: 3},
	{Group: ConfigGroupLogin, Key: "allow_register_phone", Value: "1", Type: "bool", Name: "允许手机注册", SortOrder: 4},
	{Group: ConfigGroupLogin, Key: "allow_register_email", Value: "1", Type: "bool", Name: "允许邮箱注册", SortOrder: 5},
	{Group: ConfigGroupLogin, Key: "allow_register_wechat", Value: "0", Type: "bool", Name: "允许微信注册", SortOrder: 6},
	{Group: ConfigGroupLogin, Key: "allow_login_phone", Value: "1", Type: "bool", Name: "允许手机登录", SortOrder: 7},
	{Group: ConfigGroupLogin, Key: "allow_login_email", Value: "1", Type: "bool", Name: "允许邮箱登录", SortOrder: 8},
	{Group: ConfigGroupLogin, Key: "allow_login_wechat", Value: "0", Type: "bool", Name: "允许微信登录", SortOrder: 9},
	{Group: ConfigGroupLogin, Key: "allow_email_register_code", Value: "1", Type: "bool", Name: "邮箱注册发送验证码", SortOrder: 10},
	{Group: ConfigGroupLogin, Key: "tel_cc_input", Value: "1", Type: "bool", Name: "显示国际区号选择", SortOrder: 11},
	{Group: ConfigGroupLogin, Key: "clients_profoptional", Value: "username,companyname", Type: "string", Name: "用户资料可选字段", SortOrder: 12},
	{Group: ConfigGroupLogin, Key: "clients_profuneditable", Value: "", Type: "string", Name: "用户资料不可编辑字段", SortOrder: 13},
	{Group: ConfigGroupLogin, Key: "login_register_custom_require", Value: "[]", Type: "json", Name: "自定义注册字段", SortOrder: 14},
	{Group: ConfigGroupLogin, Key: "allow_custom_clients_id", Value: "0", Type: "bool", Name: "自定义用户ID", SortOrder: 15},
	{Group: ConfigGroupLogin, Key: "custom_clients_id_start", Value: "10000", Type: "int", Name: "自定义ID起始值", SortOrder: 16},
	{Group: ConfigGroupLogin, Key: "marketing_emails_opt_in", Value: "1", Type: "bool", Name: "营销邮件默认勾选", SortOrder: 17},

	// ========== 显示配置 ==========
	{Group: ConfigGroupDisplay, Key: "language", Value: "zh-cn", Type: "string", Name: "默认语言", SortOrder: 1},
	{Group: ConfigGroupDisplay, Key: "allow_user_language", Value: "1", Type: "bool", Name: "允许用户切换语言", SortOrder: 2},
	{Group: ConfigGroupDisplay, Key: "date_format", Value: "YYYY-MM-DD", Type: "string", Name: "后台日期格式", SortOrder: 3,
		Options: `[{"label":"YYYY-MM-DD","value":"YYYY-MM-DD"},{"label":"YYYY/MM/DD","value":"YYYY/MM/DD"},{"label":"MM/DD/YYYY","value":"MM/DD/YYYY"},{"label":"DD.MM.YYYY","value":"DD.MM.YYYY"},{"label":"DD-MM-YYYY","value":"DD-MM-YYYY"}]`},
	{Group: ConfigGroupDisplay, Key: "client_date_format", Value: "full", Type: "string", Name: "用户日期格式", SortOrder: 4,
		Options: `[{"label":"2024年1月1日","value":"full"},{"label":"1月1日 2024年","value":"shortmonth"}]`},
	{Group: ConfigGroupDisplay, Key: "default_country", Value: "CN", Type: "string", Name: "默认国家", SortOrder: 5},
	{Group: ConfigGroupDisplay, Key: "per_page_limit", Value: "20", Type: "int", Name: "每页显示条数", SortOrder: 6,
		Options: `[{"label":"10条","value":"10"},{"label":"20条","value":"20"},{"label":"50条","value":"50"},{"label":"100条","value":"100"}]`},
	{Group: ConfigGroupDisplay, Key: "show_cancel", Value: "1", Type: "bool", Name: "显示取消按钮", SortOrder: 7},
	{Group: ConfigGroupDisplay, Key: "nologin_send_ticket", Value: "0", Type: "bool", Name: "未登录可发工单", SortOrder: 8},
	{Group: ConfigGroupDisplay, Key: "evaluate_ticket", Value: "1", Type: "bool", Name: "工单评价功能", SortOrder: 9},
	{Group: ConfigGroupDisplay, Key: "ticket_reply_order", Value: "asc", Type: "string", Name: "工单回复排序", SortOrder: 10,
		Options: `[{"label":"正序","value":"asc"},{"label":"倒序","value":"desc"}]`},
	{Group: ConfigGroupDisplay, Key: "dl_incl_product", Value: "1", Type: "bool", Name: "下载包含产品", SortOrder: 11},

	// ========== 模板配置 ==========
	{Group: ConfigGroupTemplate, Key: "themes_templates", Value: "default", Type: "string", Name: "主题模板", SortOrder: 1},
	{Group: ConfigGroupTemplate, Key: "clientarea_default_themes", Value: "default", Type: "string", Name: "用户中心主题", SortOrder: 2},
	{Group: ConfigGroupTemplate, Key: "order_page_style", Value: "default", Type: "string", Name: "订单页面样式", SortOrder: 3},
	{Group: ConfigGroupTemplate, Key: "admin_default_theme", Value: "default", Type: "string", Name: "后台主题", SortOrder: 4},
	{Group: ConfigGroupTemplate, Key: "header", Value: "", Type: "text", Name: "自定义头部", SortOrder: 5},
	{Group: ConfigGroupTemplate, Key: "footer", Value: "", Type: "text", Name: "自定义底部", SortOrder: 6},
	{Group: ConfigGroupTemplate, Key: "login_header", Value: "", Type: "text", Name: "登录页头部", SortOrder: 7},
	{Group: ConfigGroupTemplate, Key: "login_footer", Value: "", Type: "text", Name: "登录页底部", SortOrder: 8},
	{Group: ConfigGroupTemplate, Key: "login_header_footer", Value: "0", Type: "bool", Name: "登录页显示头底部", SortOrder: 9},
	{Group: ConfigGroupTemplate, Key: "web_widgets", Value: "", Type: "text", Name: "网页挂件", SortOrder: 10},
	{Group: ConfigGroupTemplate, Key: "cart_product_description", Value: "", Type: "text", Name: "购物车产品说明", SortOrder: 11},
	{Group: ConfigGroupTemplate, Key: "custom_login_background_description", Value: "", Type: "text", Name: "登录页背景说明", SortOrder: 12},

	// ========== 维护模式 ==========
	{Group: ConfigGroupMaintenance, Key: "main_tenance_mode", Value: "0", Type: "bool", Name: "维护模式", SortOrder: 1},
	{Group: ConfigGroupMaintenance, Key: "main_tenance_mode_url", Value: "", Type: "string", Name: "维护模式重定向URL", SortOrder: 2},
	{Group: ConfigGroupMaintenance, Key: "main_tenance_mode_message", Value: "系统维护中，请稍后再访问", Type: "text", Name: "维护模式提示信息", SortOrder: 3},

	// ========== 代理配置 ==========
	{Group: ConfigGroupAffiliate, Key: "affiliate_enabled", Value: "0", Type: "bool", Name: "代理功能开关", SortOrder: 1},
	{Group: ConfigGroupAffiliate, Key: "affiliate_percent", Value: "10", Type: "float", Name: "代理佣金比例(%)", SortOrder: 2},
	{Group: ConfigGroupAffiliate, Key: "affiliate_bonusde_posit", Value: "0", Type: "float", Name: "代理奖金", SortOrder: 3},
	{Group: ConfigGroupAffiliate, Key: "affiliate_payout", Value: "100", Type: "float", Name: "最低提现金额", SortOrder: 4},
	{Group: ConfigGroupAffiliate, Key: "affiliate_cookie", Value: "30", Type: "float", Name: "Cookie有效期(天)", SortOrder: 5},
	{Group: ConfigGroupAffiliate, Key: "affiliate_withdraw", Value: "1", Type: "bool", Name: "允许提现", SortOrder: 6},
	{Group: ConfigGroupAffiliate, Key: "affiliate_is_authentication", Value: "0", Type: "bool", Name: "代理需要认证", SortOrder: 7},
	{Group: ConfigGroupAffiliate, Key: "affiliate_delay_commission", Value: "0", Type: "int", Name: "延迟佣金(天)", SortOrder: 8},
	{Group: ConfigGroupAffiliate, Key: "affiliate_is_reorder", Value: "1", Type: "bool", Name: "计算重复订单", SortOrder: 9},
	{Group: ConfigGroupAffiliate, Key: "affiliate_reorder", Value: "0", Type: "float", Name: "重复订单佣金", SortOrder: 10},
	{Group: ConfigGroupAffiliate, Key: "affiliate_is_renew", Value: "1", Type: "bool", Name: "计算续费佣金", SortOrder: 11},
	{Group: ConfigGroupAffiliate, Key: "affiliate_renew", Value: "0", Type: "float", Name: "续费佣金比例", SortOrder: 12},
	{Group: ConfigGroupAffiliate, Key: "aff_report", Value: "1", Type: "bool", Name: "代理报告", SortOrder: 13},

	// ========== 充值配置 ==========
	{Group: ConfigGroupRecharge, Key: "addfunds_enabled", Value: "1", Type: "bool", Name: "充值功能开关", SortOrder: 1},
	{Group: ConfigGroupRecharge, Key: "addfunds_minimum", Value: "1", Type: "float", Name: "单笔最小金额", SortOrder: 2},
	{Group: ConfigGroupRecharge, Key: "addfunds_maximum", Value: "10000", Type: "float", Name: "单笔最大金额", SortOrder: 3},
	{Group: ConfigGroupRecharge, Key: "addfunds_maximum_balance", Value: "50000", Type: "float", Name: "最高余额", SortOrder: 4},
	{Group: ConfigGroupRecharge, Key: "addfunds_require_order", Value: "0", Type: "bool", Name: "充值需要关联订单", SortOrder: 5},

	// ========== 信用额度 ==========
	{Group: ConfigGroupCredit, Key: "credit_limit", Value: "0", Type: "bool", Name: "信用额度开关", SortOrder: 1},
	{Group: ConfigGroupCredit, Key: "no_auto_apply_credit", Value: "0", Type: "bool", Name: "不自动使用信用额", SortOrder: 2},
	{Group: ConfigGroupCredit, Key: "credit_on_downgrade", Value: "0", Type: "bool", Name: "降级时退还信用额", SortOrder: 3},

	// ========== 发票配置 ==========
	{Group: ConfigGroupInvoice, Key: "in_circulation_create", Value: "0", Type: "bool", Name: "循环订单创建发票", SortOrder: 1},
	{Group: ConfigGroupInvoice, Key: "in_pdf", Value: "1", Type: "bool", Name: "PDF发票", SortOrder: 2},
	{Group: ConfigGroupInvoice, Key: "in_save_user_info", Value: "1", Type: "bool", Name: "保存用户发票信息", SortOrder: 3},
	{Group: ConfigGroupInvoice, Key: "in_batch_pay", Value: "1", Type: "bool", Name: "批量支付", SortOrder: 4},
	{Group: ConfigGroupInvoice, Key: "in_select_payment", Value: "1", Type: "bool", Name: "选择支付方式", SortOrder: 5},
	{Group: ConfigGroupInvoice, Key: "in_unpaid_tick", Value: "1", Type: "bool", Name: "显示未支付订单", SortOrder: 6},
	{Group: ConfigGroupInvoice, Key: "in_continuous_pay_num", Value: "0", Type: "int", Name: "连续支付数量(0不限)", SortOrder: 7},
	{Group: ConfigGroupInvoice, Key: "in_continuous_pay_num_type", Value: "MONTH", Type: "string", Name: "连续支付类型", SortOrder: 8,
		Options: `[{"label":"年","value":"YEAR"},{"label":"月","value":"MONTH"},{"label":"日","value":"DAY"},{"label":"次数","value":"NUMBER"}]`},
	{Group: ConfigGroupInvoice, Key: "in_overdue_fine", Value: "0", Type: "bool", Name: "逾期罚款", SortOrder: 9},
	{Group: ConfigGroupInvoice, Key: "in_overdue_fine_min", Value: "0", Type: "float", Name: "逾期罚款最小金额", SortOrder: 10},
	{Group: ConfigGroupInvoice, Key: "invoice_payto", Value: "", Type: "text", Name: "付款条文", SortOrder: 11},

	// ========== SEO配置 ==========
	{Group: ConfigGroupSEO, Key: "seo_keywords", Value: "", Type: "string", Name: "SEO关键词", SortOrder: 1},
	{Group: ConfigGroupSEO, Key: "seo_desc", Value: "", Type: "text", Name: "SEO描述", SortOrder: 2},

	// ========== 日志配置 ==========
	{Group: ConfigGroupAdvanced, Key: "sendmsgtimes", Value: "10", Type: "int", Name: "每天短信发送次数", SortOrder: 1},
	{Group: ConfigGroupAdvanced, Key: "sendmsgphone", Value: "5", Type: "int", Name: "每天短信发送手机个数", SortOrder: 2},
	{Group: ConfigGroupAdvanced, Key: "deletelogtime", Value: "90", Type: "int", Name: "删除日志天数", SortOrder: 3},
	{Group: ConfigGroupAdvanced, Key: "activity_limit", Value: "1000", Type: "int", Name: "活动日志限制", SortOrder: 4},
	{Group: ConfigGroupAdvanced, Key: "display_errors", Value: "0", Type: "bool", Name: "显示错误信息", SortOrder: 5},
	{Group: ConfigGroupAdvanced, Key: "sql_error_reporting", Value: "0", Type: "bool", Name: "SQL错误报告", SortOrder: 6},
	{Group: ConfigGroupAdvanced, Key: "hooks_debug_mode", Value: "0", Type: "bool", Name: "钩子调试模式", SortOrder: 7},
	{Group: ConfigGroupAdvanced, Key: "in_circulation_create", Value: "0", Type: "bool", Name: "循环订单创建", SortOrder: 8},

	// ========== OAuth配置 ==========
	{Group: ConfigGroupOAuth, Key: "wechat_login_appid", Value: "", Type: "string", Name: "微信登录AppID", SortOrder: 1},
	{Group: ConfigGroupOAuth, Key: "wechat_login_secret", Value: "", Type: "string", Name: "微信登录Secret", SortOrder: 2},

	// ========== 短信配置 ==========
	{Group: ConfigGroupSMS, Key: "sms_provider", Value: "aliyun", Type: "string", Name: "短信服务商", SortOrder: 1,
		Options: `[{"label":"阿里云","value":"aliyun"},{"label":"腾讯云","value":"tencent"},{"label":"互亿无线","value":"huyi"},{"label":"自定义","value":"custom"}]`},
	{Group: ConfigGroupSMS, Key: "accesskeyid", Value: "", Type: "string", Name: "AccessKeyID", SortOrder: 2},
	{Group: ConfigGroupSMS, Key: "accesskeysecret", Value: "", Type: "string", Name: "AccessKeySecret", SortOrder: 3},
	{Group: ConfigGroupSMS, Key: "code", Value: "", Type: "string", Name: "短信模板码", SortOrder: 4},
	{Group: ConfigGroupSMS, Key: "signature", Value: "", Type: "string", Name: "短信签名", SortOrder: 5},
	{Group: ConfigGroupSMS, Key: "submail_appid", Value: "", Type: "string", Name: "SubMail AppID", SortOrder: 6},
	{Group: ConfigGroupSMS, Key: "submail_appkey", Value: "", Type: "string", Name: "SubMail AppKey", SortOrder: 7},

	// ========== 支付宝认证 ==========
	{Group: ConfigGroupPayment, Key: "alipay_app_id", Value: "", Type: "string", Name: "支付宝AppID", SortOrder: 1},
	{Group: ConfigGroupPayment, Key: "alipay_private_key", Value: "", Type: "text", Name: "支付宝私钥", SortOrder: 2},
	{Group: ConfigGroupPayment, Key: "alipay_public_key", Value: "", Type: "text", Name: "支付宝公钥", SortOrder: 3},
	{Group: ConfigGroupPayment, Key: "alipay_biz_code", Value: "FACE", Type: "string", Name: "支付宝业务码", SortOrder: 4},

	// ========== Cookie配置 ==========
	{Group: ConfigGroupAdvanced, Key: "cookie_clientarea_nmae", Value: "", Type: "string", Name: "用户中心Cookie名称", SortOrder: 10},
}

// ConfigService 配置服务
type ConfigService struct {
	db *gorm.DB
}

func NewConfigService(db *gorm.DB) *ConfigService {
	return &ConfigService{db: db}
}

// InitDefaultConfigs 初始化默认配置
func (s *ConfigService) InitDefaultConfigs() error {
	for _, config := range DefaultSystemConfigs {
		var existing SystemConfig
		if err := s.db.Where("`key` = ?", config.Key).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := s.db.Create(&config).Error; err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Get 获取配置值
func (s *ConfigService) Get(key string) (string, error) {
	var config SystemConfig
	if err := s.db.Where("`key` = ? AND status = ?", key, true).First(&config).Error; err != nil {
		return "", err
	}
	return config.Value, nil
}

// GetWithDefault 获取配置值（带默认值）
func (s *ConfigService) GetWithDefault(key, defaultValue string) string {
	value, err := s.Get(key)
	if err != nil || value == "" {
		return defaultValue
	}
	return value
}

// GetBool 获取布尔配置
func (s *ConfigService) GetBool(key string) bool {
	value, _ := s.Get(key)
	return value == "1" || value == "true"
}

// GetInt 获取整数配置
func (s *ConfigService) GetInt(key string) int {
	var config SystemConfig
	if err := s.db.Where("`key` = ? AND status = ?", key, true).First(&config).Error; err != nil {
		return 0
	}
	var result int
	s.db.Raw("SELECT CAST(? AS INTEGER)", config.Value).Scan(&result)
	return result
}

// Set 设置配置值
func (s *ConfigService) Set(key, value string) error {
	return s.db.Model(&SystemConfig{}).Where("`key` = ?", key).Update("value", value).Error
}

// SetBatch 批量设置配置
func (s *ConfigService) SetBatch(configs map[string]string) error {
	tx := s.db.Begin()
	for key, value := range configs {
		if err := tx.Model(&SystemConfig{}).Where("`key` = ?", key).Update("value", value).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// GetByGroup 按分组获取配置
func (s *ConfigService) GetByGroup(group string) ([]SystemConfig, error) {
	var configs []SystemConfig
	if err := s.db.Where("`group` = ? AND status = ?", group, true).Order("sort_order").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetAll 获取所有配置
func (s *ConfigService) GetAll() ([]SystemConfig, error) {
	var configs []SystemConfig
	if err := s.db.Order("`group`, sort_order").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

// GetAllAsMap 获取所有配置为Map
func (s *ConfigService) GetAllAsMap() (map[string]string, error) {
	var configs []SystemConfig
	if err := s.db.Find(&configs).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for _, config := range configs {
		result[config.Key] = config.Value
	}
	return result, nil
}

// GetGroups 获取所有配置分组
func (s *ConfigService) GetGroups() []map[string]string {
	return []map[string]string{
		{"key": ConfigGroupGeneral, "name": "基本设置"},
		{"key": ConfigGroupSecurity, "name": "安全设置"},
		{"key": ConfigGroupLogin, "name": "登录注册"},
		{"key": ConfigGroupCaptcha, "name": "验证码"},
		{"key": ConfigGroupDisplay, "name": "显示配置"},
		{"key": ConfigGroupTemplate, "name": "模板配置"},
		{"key": ConfigGroupMaintenance, "name": "维护模式"},
		{"key": ConfigGroupAffiliate, "name": "代理配置"},
		{"key": ConfigGroupRecharge, "name": "充值配置"},
		{"key": ConfigGroupCredit, "name": "信用额度"},
		{"key": ConfigGroupInvoice, "name": "发票配置"},
		{"key": ConfigGroupSEO, "name": "SEO配置"},
		{"key": ConfigGroupOAuth, "name": "OAuth配置"},
		{"key": ConfigGroupSMS, "name": "短信配置"},
		{"key": ConfigGroupPayment, "name": "支付配置"},
		{"key": ConfigGroupEmail, "name": "邮件配置"},
		{"key": ConfigGroupNotification, "name": "通知配置"},
		{"key": ConfigGroupAdvanced, "name": "高级设置"},
	}
}

// ========== 便捷方法 ==========

// IsMaintenanceMode 是否维护模式
func (s *ConfigService) IsMaintenanceMode() bool {
	return s.GetBool("main_tenance_mode")
}

// IsPhoneLoginAllowed 是否允许手机登录
func (s *ConfigService) IsPhoneLoginAllowed() bool {
	return s.GetBool("allow_phone") && s.GetBool("allow_login_phone")
}

// IsEmailLoginAllowed 是否允许邮箱登录
func (s *ConfigService) IsEmailLoginAllowed() bool {
	return s.GetBool("allow_email") && s.GetBool("allow_login_email")
}

// IsPhoneRegisterAllowed 是否允许手机注册
func (s *ConfigService) IsPhoneRegisterAllowed() bool {
	return s.GetBool("allow_register_phone")
}

// IsEmailRegisterAllowed 是否允许邮箱注册
func (s *ConfigService) IsEmailRegisterAllowed() bool {
	return s.GetBool("allow_register_email")
}

// IsWechatLoginAllowed 是否允许微信登录
func (s *ConfigService) IsWechatLoginAllowed() bool {
	return s.GetBool("allow_login_wechat")
}

// IsAffiliateEnabled 是否启用代理
func (s *ConfigService) IsAffiliateEnabled() bool {
	return s.GetBool("affiliate_enabled")
}

// IsRechargeEnabled 是否启用充值
func (s *ConfigService) IsRechargeEnabled() bool {
	return s.GetBool("addfunds_enabled")
}

// IsCreditEnabled 是否启用信用额度
func (s *ConfigService) IsCreditEnabled() bool {
	return s.GetBool("credit_limit")
}

// GetLogo 获取Logo（根据场景）
func (s *ConfigService) GetLogo(scene string) string {
	switch scene {
	case "home", "web":
		logo := s.GetWithDefault("logo_url_home", "")
		if logo != "" {
			return logo
		}
	case "bill", "invoice":
		logo := s.GetWithDefault("logo_url_bill", "")
		if logo != "" {
			return logo
		}
	case "admin":
		logo := s.GetWithDefault("logo_url_admin", "")
		if logo != "" {
			return logo
		}
	case "www":
		logo := s.GetWithDefault("www_logo", "")
		if logo != "" {
			return logo
		}
	}
	return s.GetWithDefault("logo_url", "")
}

// GetPasswordStrength 获取密码强度要求
func (s *ConfigService) GetPasswordStrength() string {
	return s.GetWithDefault("required_pwstrength", "alpha_num")
}

// GetLoginErrorConfig 获取登录错误配置
func (s *ConfigService) GetLoginErrorConfig() (bool, int, int) {
	enabled := s.GetBool("login_error_switch")
	maxNum := s.GetInt("login_error_max_num")
	banLength := s.GetInt("invalid_logins_banlength")
	return enabled, maxNum, banLength
}

// GetMaintenanceConfig 获取维护模式配置
func (s *ConfigService) GetMaintenanceConfig() (bool, string, string) {
	enabled := s.IsMaintenanceMode()
	url := s.GetWithDefault("main_tenance_mode_url", "")
	message := s.GetWithDefault("main_tenance_mode_message", "系统维护中，请稍后再访问")
	return enabled, url, message
}

// GetLoginMethods 获取登录方式配置
func (s *ConfigService) GetLoginMethods() map[string]bool {
	return map[string]bool{
		"phone":   s.IsPhoneLoginAllowed(),
		"email":   s.IsEmailLoginAllowed(),
		"wechat":  s.IsWechatLoginAllowed(),
		"id":      s.GetBool("allow_id"),
	}
}

// GetRegisterMethods 获取注册方式配置
func (s *ConfigService) GetRegisterMethods() map[string]bool {
	return map[string]bool{
		"phone":   s.IsPhoneRegisterAllowed(),
		"email":   s.IsEmailRegisterAllowed(),
		"wechat":  s.GetBool("allow_register_wechat"),
	}
}

// GetPublicConfig 获取公开配置（前端需要）
func (s *ConfigService) GetPublicConfig() map[string]interface{} {
	return map[string]interface{}{
		// 公司信息
		"company_name":        s.GetWithDefault("company_name", ""),
		"company_email":       s.GetWithDefault("company_email", ""),
		"company_phone":       s.GetWithDefault("company_phone", ""),
		"company_address":     s.GetWithDefault("company_address", ""),
		"record_no":           s.GetWithDefault("record_no", ""),
		"system_url":          s.GetWithDefault("system_url", ""),

		// Logo
		"logo_url":            s.GetLogo("default"),
		"logo_url_home":       s.GetLogo("home"),
		"favicon_url":         s.GetWithDefault("favicon_url", ""),

		// 登录注册方式
		"login_methods":       s.GetLoginMethods(),
		"register_methods":    s.GetRegisterMethods(),

		// 功能开关
		"affiliate_enabled":   s.IsAffiliateEnabled(),
		"addfunds_enabled":    s.IsRechargeEnabled(),
		"credit_limit":        s.IsCreditEnabled(),
		"show_cancel":         s.GetBool("show_cancel"),
		"nologin_send_ticket": s.GetBool("nologin_send_ticket"),
		"evaluate_ticket":     s.GetBool("evaluate_ticket"),

		// 显示配置
		"language":            s.GetWithDefault("language", "zh-cn"),
		"allow_user_language": s.GetBool("allow_user_language"),
		"default_country":     s.GetWithDefault("default_country", "CN"),

		// 法律条款
		"server_clause_url":   s.GetWithDefault("server_clause_url", ""),
		"privacy_clause_url":  s.GetWithDefault("privacy_clause_url", ""),

		// SEO
		"seo_keywords":        s.GetWithDefault("seo_keywords", ""),
		"seo_desc":            s.GetWithDefault("seo_desc", ""),

		// 维护模式
		"maintenance_mode":    s.IsMaintenanceMode(),
	}
}
