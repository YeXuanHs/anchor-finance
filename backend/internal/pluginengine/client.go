package pluginengine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PluginEngine PHP插件引擎HTTP客户端
// Go → HTTP → PHP引擎(9000) → 插件handler
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// Default 默认客户端
var Default = &Client{
	BaseURL: "http://127.0.0.1:9000",
	HTTP:    &http.Client{Timeout: 30 * time.Second},
}

// Result Hook执行结果
type Result struct {
	Handler string      `json:"handler"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// TriggerHook 触发Hook
func TriggerHook(hook string, params map[string]interface{}) ([]Result, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"hook":   hook,
		"params": params,
	})

	resp, err := Default.HTTP.Post(
		Default.BaseURL+"/internal/hook/trigger",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("插件引擎连接失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var r struct {
		Code int `json:"code"`
		Data struct {
			Results []Result `json:"results"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	return r.Data.Results, nil
}

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	_, err := TriggerHook("send_email", map[string]interface{}{
		"to": to, "subject": subject, "body": body,
	})
	return err
}

// SendSMS 发送短信
func SendSMS(phone, content string) error {
	_, err := TriggerHook("send_sms", map[string]interface{}{
		"phone": phone, "content": content,
	})
	return err
}

// HealthCheck 检查插件引擎是否在线
func HealthCheck() error {
	resp, err := Default.HTTP.Get(Default.BaseURL + "/health")
	if err != nil {
		return fmt.Errorf("插件引擎离线: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("插件引擎状态异常: %d", resp.StatusCode)
	}
	return nil
}

// HandlePaymentCallback 转发支付回调给PHP插件引擎处理
// PHP插件负责：验证网关签名、更新账单状态、记录交易
func HandlePaymentCallback(gateway string, params map[string]interface{}) (map[string]interface{}, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"gateway": gateway,
		"params":  params,
	})

	resp, err := Default.HTTP.Post(
		Default.BaseURL+"/internal/payment/callback/"+gateway,
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("插件引擎连接失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var r struct {
		Code int                    `json:"code"`
		Data map[string]interface{} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &r); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	if r.Code != 0 {
		return nil, fmt.Errorf("支付处理失败: code=%d", r.Code)
	}
	return r.Data, nil
}
