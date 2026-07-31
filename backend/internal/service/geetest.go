package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

// GeetestAPIServer 极验官方API服务器地址
const GeetestAPIServer = "https://gcaptcha4.geetest.com"

// GeetestService 极验验证码服务
type GeetestService struct {
	captchaID  string
	captchaKey string
	log        *zap.Logger
}

// GeetestVerifyResult 极验验证结果
type GeetestVerifyResult struct {
	Status      string                 `json:"status"`
	Result      string                 `json:"result"`
	Reason      string                 `json:"reason"`
	CaptchaArgs map[string]interface{} `json:"captcha_args"`
	Code        string                 `json:"code"`
	Msg         string                 `json:"msg"`
}

// GeetestVerifyParams 极验验证参数
type GeetestVerifyParams struct {
	LotNumber     string `json:"lot_number" form:"lot_number"`
	CaptchaOutput string `json:"captcha_output" form:"captcha_output"`
	PassToken     string `json:"pass_token" form:"pass_token"`
	GenTime       string `json:"gen_time" form:"gen_time"`
}

// NewGeetestService 创建极验服务
func NewGeetestService(captchaID, captchaKey string, log *zap.Logger) *GeetestService {
	return &GeetestService{
		captchaID:  captchaID,
		captchaKey: captchaKey,
		log:        log,
	}
}

// GetCaptchaID 获取 Captcha ID（前端用）
func (s *GeetestService) GetCaptchaID() string {
	return s.captchaID
}

// GenerateSignToken 生成签名令牌
// 使用 HMAC-SHA256，lot_number 作为 message，captcha_key 作为 key
func (s *GeetestService) GenerateSignToken(lotNumber string) string {
	h := hmac.New(sha256.New, []byte(s.captchaKey))
	h.Write([]byte(lotNumber))
	return hex.EncodeToString(h.Sum(nil))
}

// Verify 验证极验参数
func (s *GeetestService) Verify(params *GeetestVerifyParams) (bool, string, error) {
	if s.captchaID == "" || s.captchaKey == "" {
		return false, "极验未配置", fmt.Errorf("geetest not configured")
	}

	// 生成签名
	signToken := s.GenerateSignToken(params.LotNumber)

	// 构建请求参数
	formData := url.Values{}
	formData.Set("lot_number", params.LotNumber)
	formData.Set("captcha_output", params.CaptchaOutput)
	formData.Set("pass_token", params.PassToken)
	formData.Set("gen_time", params.GenTime)
	formData.Set("sign_token", signToken)

	// 构建请求URL（使用官方固定地址）
	requestURL := fmt.Sprintf("%s/validate?captcha_id=%s", GeetestAPIServer, s.captchaID)

	// 发送请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Post(requestURL, "application/x-www-form-urlencoded", strings.NewReader(formData.Encode()))
	if err != nil {
		s.log.Error("极验验证请求失败", zap.Error(err))
		// 请求失败时默认通过验证（容灾策略）
		return true, "geetest api request failed, pass by default", nil
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("极验验证读取响应失败", zap.Error(err))
		return false, "读取响应失败", err
	}

	// 解析响应
	var result GeetestVerifyResult
	if err := json.Unmarshal(body, &result); err != nil {
		s.log.Error("极验验证解析响应失败", zap.Error(err))
		return false, "解析响应失败", err
	}

	s.log.Info("极验验证结果",
		zap.String("status", result.Status),
		zap.String("result", result.Result),
		zap.String("reason", result.Reason),
	)

	// 判断验证结果
	if result.Status == "success" && result.Result == "success" {
		return true, "验证成功", nil
	}

	reason := result.Reason
	if reason == "" {
		reason = result.Msg
	}
	if reason == "" {
		reason = "验证失败"
	}

	return false, reason, nil
}

// VerifyFromRequest 从 HTTP 请求中验证极验参数
func (s *GeetestService) VerifyFromRequest(r *http.Request) (bool, string, error) {
	params := &GeetestVerifyParams{
		LotNumber:     r.FormValue("lot_number"),
		CaptchaOutput: r.FormValue("captcha_output"),
		PassToken:     r.FormValue("pass_token"),
		GenTime:       r.FormValue("gen_time"),
	}

	// 检查参数是否完整
	if params.LotNumber == "" || params.CaptchaOutput == "" || params.PassToken == "" || params.GenTime == "" {
		return false, "验证参数不完整", fmt.Errorf("missing geetest params")
	}

	return s.Verify(params)
}

// GetFrontendConfig 获取前端配置
func (s *GeetestService) GetFrontendConfig() map[string]interface{} {
	return map[string]interface{}{
		"captcha_id": s.captchaID,
	}
}
