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

// IdcsmartaliConfig 阿里云实名认证配置
type IdcsmartaliConfig struct {
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
}

// IdcsmartaliPlugin 阿里云实名认证插件
type IdcsmartaliPlugin struct {
	config *IdcsmartaliConfig
}

// NewIdcsmartaliPlugin 创建阿里云实名认证插件
func NewIdcsmartaliPlugin(configJSON string) (*IdcsmartaliPlugin, error) {
	var config IdcsmartaliConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, fmt.Errorf("invalid aliyun certification config: %w", err)
	}
	if config.AccessKeyID == "" || config.AccessKeySecret == "" {
		return nil, fmt.Errorf("missing required fields: access_key_id, access_key_secret")
	}
	return &IdcsmartaliPlugin{config: &config}, nil
}

func (p *IdcsmartaliPlugin) Name() string  { return "Idcsmartali" }
func (p *IdcsmartaliPlugin) Title() string { return "阿里云实名认证" }

// Personal 个人实名认证
func (p *IdcsmartaliPlugin) Personal(name, card string) (*CertResult, error) {
	// 调用阿里云实人认证API
	params := map[string]string{
		"Action":         "InitiateFaceVerification",
		"Version":        "2019-03-07",
		"Format":         "JSON",
		"RegionId":       "cn-hangzhou",
		"Name":           name,
		"CertNo":         card,
		"CertType":       "IDENTITY_CARD",
		"OuterOrderNo":   fmt.Sprintf("AF%d", time.Now().UnixMilli()),
		"SceneCode":      "FACE",
		"Model":          "LIVENESS",
	}

	result, err := p.doRequest(params)
	if err != nil {
		return &CertResult{
			Status: CertStatusFailed,
			Msg:    fmt.Sprintf("认证接口调用失败: %s", err.Error()),
		}, nil
	}

	// 解析响应
	certifyId, _ := result["CertifyId"].(string)
	if certifyId == "" {
		return &CertResult{
			Status: CertStatusFailed,
			Msg:    "接口返回数据不完整",
		}, nil
	}

	return &CertResult{
		Status:    CertStatusSubmitted,
		CertifyID: certifyId,
		URL:       fmt.Sprintf("https://cloudauth.aliyun.com/verify/certify?CertifyId=%s", certifyId),
		Msg:       "请点击链接完成实名认证",
	}, nil
}

// Company 企业实名认证
func (p *IdcsmartaliPlugin) Company(name, card string) (*CertResult, error) {
	// 企业认证使用不同的接口
	params := map[string]string{
		"Action":         "InitiateFaceVerification",
		"Version":        "2019-03-07",
		"Format":         "JSON",
		"RegionId":       "cn-hangzhou",
		"Name":           name,
		"CertNo":         card,
		"CertType":       "IDENTITY_CARD",
		"OuterOrderNo":   fmt.Sprintf("AF%d", time.Now().UnixMilli()),
		"SceneCode":      "FACE",
		"Model":          "LIVENESS",
	}

	result, err := p.doRequest(params)
	if err != nil {
		return &CertResult{
			Status: CertStatusFailed,
			Msg:    fmt.Sprintf("认证接口调用失败: %s", err.Error()),
		}, nil
	}

	certifyId, _ := result["CertifyId"].(string)
	if certifyId == "" {
		return &CertResult{
			Status: CertStatusFailed,
			Msg:    "接口返回数据不完整",
		}, nil
	}

	return &CertResult{
		Status:    CertStatusSubmitted,
		CertifyID: certifyId,
		URL:       fmt.Sprintf("https://cloudauth.aliyun.com/verify/certify?CertifyId=%s", certifyId),
		Msg:       "请点击链接完成企业认证",
	}, nil
}

// GetStatus 查询认证状态
func (p *IdcsmartaliPlugin) GetStatus(certifyId string) (*CertQueryResult, error) {
	params := map[string]string{
		"Action":       "DescribeFaceVerifyResult",
		"Version":      "2019-03-07",
		"Format":       "JSON",
		"RegionId":     "cn-hangzhou",
		"CertifyId":    certifyId,
		"SceneCode":    "FACE",
	}

	result, err := p.doRequest(params)
	if err != nil {
		return &CertQueryResult{
			Status: CertStatusFailed,
			Msg:    fmt.Sprintf("查询失败: %s", err.Error()),
		}, nil
	}

	// 解析状态
	code, _ := result["Code"].(string)
	if code != "200" {
		return &CertQueryResult{
			Status: CertStatusPending,
			Msg:    fmt.Sprintf("认证进行中"),
		}, nil
	}

	verifyResult, _ := result["VerifyResult"].(map[string]interface{})
	status, _ := verifyResult["VerifyStatus"].(string)

	certStatus := CertStatusPending
	switch status {
	case "PASSED":
		certStatus = CertStatusPassed
	case "FAILED":
		certStatus = CertStatusFailed
	}

	return &CertQueryResult{
		Status: certStatus,
		Msg:    fmt.Sprintf("认证状态: %s", status),
	}, nil
}

// doRequest 调用阿里云API
func (p *IdcsmartaliPlugin) doRequest(params map[string]string) (map[string]interface{}, error) {
	// 添加公共参数
	params["AccessKeyId"] = p.config.AccessKeyID
	params["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	params["SignatureMethod"] = "HMAC-SHA1"
	params["SignatureVersion"] = "1.0"
	params["SignatureNonce"] = fmt.Sprintf("%d", time.Now().UnixNano())

	// 生成签名
	params["Signature"] = p.sign(params)

	// 构建请求URL
	apiURL := "https://cloudauth.aliyuncs.com/"
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

// sign 签名
func (p *IdcsmartaliPlugin) sign(params map[string]string) string {
	// 排序参数
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "Signature" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// 拼接
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(params[k])))
	}
	str := "GET&%2F&" + url.QueryEscape(strings.Join(parts, "&"))

	// HMAC-SHA1
	mac := hmac.New(sha256.New, []byte(p.config.AccessKeySecret+"&"))
	mac.Write([]byte(str))
	return hex.EncodeToString(mac.Sum(nil))
}
