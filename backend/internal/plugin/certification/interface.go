package certification

// CertStatus 认证状态常量
const (
	CertStatusPending  = 0 // 待认证
	CertStatusPassed   = 1 // 认证通过
	CertStatusFailed   = 2 // 认证失败
	CertStatusExpired  = 3 // 认证过期
	CertStatusSubmitted = 4 // 已提交资料
)

// CertResult 认证结果
type CertResult struct {
	Status    int    `json:"status"`     // 认证状态
	CertifyID string `json:"certify_id"` // 认证证书ID
	URL       string `json:"url"`        // 认证URL（扫码链接等）
	Msg       string `json:"msg"`        // 提示信息
}

// CertQueryResult 认证查询结果
type CertQueryResult struct {
	Status int    `json:"status"` // 认证状态: 1通过 2失败 4已提交
	Msg    string `json:"msg"`    // 提示信息
}

// Certification 实名认证接口
type Certification interface {
	// Personal 个人实名认证
	// name: 姓名, card: 证件号码
	// 返回认证URL或错误
	Personal(name, card string) (*CertResult, error)

	// Company 企业实名认证
	// name: 企业名称, card: 统一社会信用代码
	// 返回认证URL或错误
	Company(name, card string) (*CertResult, error)

	// GetStatus 查询认证状态
	// certifyId: 认证证书ID
	// 返回认证状态和提示信息
	GetStatus(certifyId string) (*CertQueryResult, error)

	// Name 插件名称
	Name() string

	// Title 插件显示标题
	Title() string
}
