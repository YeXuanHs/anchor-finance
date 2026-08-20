package certification

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// WechatConfig 微信实名认证配置（腾讯云人脸核身）
type WechatConfig struct {
	SecretID  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
}

// WechatPlugin 微信实名认证插件
type WechatPlugin struct {
	config *WechatConfig
}

// NewWechatPlugin 创建微信实名认证插件
func NewWechatPlugin(configJSON string) (*WechatPlugin, error) {
	var config WechatConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("invalid wechat certification config: %w", err)
	}
	if config.SecretID == "" || config.SecretKey == "" {
		return nil, fmt.Errorf("missing required fields: secret_id, secret_key")
	}
	return &WechatPlugin{config: &config}, nil
}

func (p *WechatPlugin) Name() string  { return "Wechat" }
func (p *WechatPlugin) Title() string { return "微信实名认证" }

// Personal 个人实名认证
func (p *WechatPlugin) Personal(name, card string) (*CertResult, error) {
	return p.getDetectAuth(name, card)
}

// Company 企业实名认证
func (p *WechatPlugin) Company(name, card string) (*CertResult, error) {
	return p.getDetectAuth(name, card)
}

// GetStatus 查询认证状态
func (p *WechatPlugin) GetStatus(certifyId string) (*CertQueryResult, error) {
	return p.getDetectAuthResult(certifyId)
}

// getDetectAuth 发起人脸核身
func (p *WechatPlugin) getDetectAuth(name, idCard string) (*CertResult, error) {
	// 构建请求参数
	params := map[string]string{
		"Action":           "GetDetectAuth",
		"Version":          "2018-03-01",
		"Region":           "ap-guangzhou",
		"Name":             name,
		"IdCard":           idCard,
		"RuleId":           "1",
		"Nonce":            fmt.Sprintf("%d", time.Now().UnixNano()),
		"Timestamp":        fmt.Sprintf("%d", time.Now().Unix()),
	}

	// 生成签名
	params["Signature"] = p.sign(params)

	// 调用API
	result, err := p.doRequest(params)
	if err != nil {
		return &CertResult{
			Status: CertStatusFailed,
			Msg:    fmt.Sprintf("认证接口调用失败: %s", err.Error()),
		}, nil
	}

	// 解析响应
	response, ok := result["Response"].(map[string]interface{})
	if !ok {
		return &CertResult{
			Status: CertStatusFailed,
			Msg:    "无效的接口响应",
		}, nil
	}

	if errMsg, ok := response["Error"].(map[string]interface{}); ok {
		return &CertResult{
			Status: CertStatusFailed,
			Msg:    fmt.Sprintf("%v", errMsg["Message"]),
		}, nil
	}

	bizToken, _ := response["BizToken"].(string)
	url, _ := response["Url"].(string)

	if bizToken == "" || url == "" {
		return &CertResult{
			Status: CertStatusFailed,
			Msg:    "接口返回数据不完整",
		}, nil
	}

	return &CertResult{
		Status:    CertStatusSubmitted,
		CertifyID: bizToken,
		URL:       url,
		Msg:       "请使用微信扫码完成实名认证",
	}, nil
}

// getDetectAuthResult 查询核身结果
func (p *WechatPlugin) getDetectAuthResult(bizToken string) (*CertQueryResult, error) {
	params := map[string]string{
		"Action":    "GetDetectAuthResult",
		"Version":   "2018-03-01",
		"Region":    "ap-guangzhou",
		"BizToken":  bizToken,
		"Nonce":     fmt.Sprintf("%d", time.Now().UnixNano()),
		"Timestamp": fmt.Sprintf("%d", time.Now().Unix()),
	}

	params["Signature"] = p.sign(params)

	result, err := p.doRequest(params)
	if err != nil {
		return &CertQueryResult{
			Status: CertStatusFailed,
			Msg:    fmt.Sprintf("查询失败: %s", err.Error()),
		}, nil
	}

	response, ok := result["Response"].(map[string]interface{})
	if !ok {
		return &CertQueryResult{
			Status: CertStatusFailed,
			Msg:    "无效的接口响应",
		}, nil
	}

	// Status: 0未开始 1认证中 2认证成功 3认证失败
	status, _ := response["Status"].(float64)
	certStatus := CertStatusPending
	switch int(status) {
	case 2:
		certStatus = CertStatusPassed
	case 3:
		certStatus = CertStatusFailed
	}

	return &CertQueryResult{
		Status: certStatus,
		Msg:    fmt.Sprintf("认证状态: %d", int(status)),
	}, nil
}

// sign 签名（HMAC-SHA256）
func (p *WechatPlugin) sign(params map[string]string) string {
	// 排序参数
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", k, params[k]))
	}
	str := strings.Join(parts, "&")

	// HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(p.config.SecretKey))
	mac.Write([]byte(str))
	return hex.EncodeToString(mac.Sum(nil))
}

// doRequest 调用腾讯云API
func (p *WechatPlugin) doRequest(params map[string]string) (map[string]interface{}, error) {
	// 构建请求URL
	apiURL := "https://faceid.tencentcloudapi.com/"

	// 使用GET请求
	query := url.Values{}
	for k, v := range params {
		query.Set(k, v)
	}
	fullURL := apiURL + "?" + query.Encode()

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("invalid response: %s", string(body))
	}

	return result, nil
}
