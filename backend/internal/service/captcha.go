package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/dchest/captcha"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	captchaTTL        = 5 * time.Minute
	smsCodeLength     = 6
	emailCodeLength   = 6
	rateLimitDuration = 60 * time.Second
)

// CaptchaType 验证码类型
type CaptchaType string

const (
	CaptchaTypeImage  CaptchaType = "image"  // 图形验证码（防人机）
	CaptchaTypeSMS    CaptchaType = "sms"    // 短信验证码
	CaptchaTypeEmail  CaptchaType = "email"  // 邮箱验证码
)

// CaptchaService handles captcha generation and verification.
type CaptchaService struct {
	rdb    *redis.Client
	config *CaptchaConfigService
}

func NewCaptchaService(rdb *redis.Client, db *gorm.DB) *CaptchaService {
	return &CaptchaService{
		rdb:    rdb,
		config: NewCaptchaConfigService(db),
	}
}

// ============ 防人机图形验证码 ============

// GenerateImage generates an image captcha using dchest/captcha library.
// Returns the captcha ID and PNG image bytes.
func (s *CaptchaService) GenerateImage(key string) (string, []byte, error) {
	ctx := context.Background()

	// 获取验证码长度配置
	length := s.config.GetCaptchaLength()
	if length < 4 {
		length = 4
	}

	// 生成验证码
	captchaID := captcha.NewLen(length)

	// Store mapping from our key to captcha ID
	if err := s.rdb.Set(ctx, "captcha:image:"+key, captchaID, captchaTTL).Err(); err != nil {
		return "", nil, fmt.Errorf("failed to store captcha mapping: %w", err)
	}

	// Encode captcha to PNG
	var buf bytes.Buffer
	if err := captcha.WriteImage(&buf, captchaID, 200, 80); err != nil {
		return "", nil, fmt.Errorf("failed to generate captcha image: %w", err)
	}

	return captchaID, buf.Bytes(), nil
}

// VerifyImage verifies an image captcha against the stored captcha ID.
func (s *CaptchaService) VerifyImage(key, captchaID, digits string) bool {
	ctx := context.Background()

	// Check if the captcha ID matches the stored one
	storedID, err := s.rdb.Get(ctx, "captcha:image:"+key).Result()
	if err != nil || storedID != captchaID {
		return false
	}

	// Verify using dchest/captcha
	if !captcha.VerifyString(captchaID, digits) {
		return false
	}

	// Consume the captcha (one-time use)
	s.rdb.Del(ctx, "captcha:image:"+key)
	return true
}

// ShouldShowCaptcha 检查某个场景是否应该显示验证码
func (s *CaptchaService) ShouldShowCaptcha(scene string) bool {
	return s.config.ShouldShowCaptcha(scene)
}

// GetSceneConfig 获取所有场景配置
func (s *CaptchaService) GetSceneConfig() map[string]bool {
	return s.config.GetSceneConfig()
}

// ============ 短信验证码 ============

// GenerateSMS generates and stores a numeric SMS code.
func (s *CaptchaService) GenerateSMS(phone string) (string, error) {
	ctx := context.Background()

	// Rate limit check
	if exists, _ := s.rdb.Exists(ctx, "captcha:sms:rl:"+phone).Result(); exists > 0 {
		return "", fmt.Errorf("请稍后再请求验证码")
	}

	code := generateNumericCode(smsCodeLength)

	if err := s.rdb.Set(ctx, "captcha:sms:"+phone, code, captchaTTL).Err(); err != nil {
		return "", fmt.Errorf("failed to store SMS code: %w", err)
	}
	if err := s.rdb.Set(ctx, "captcha:sms:rl:"+phone, "1", rateLimitDuration).Err(); err != nil {
		return "", fmt.Errorf("failed to set rate limit: %w", err)
	}

	return code, nil
}

// VerifySMS verifies an SMS code and consumes it (one-time use).
func (s *CaptchaService) VerifySMS(phone, code string) bool {
	ctx := context.Background()

	stored, err := s.rdb.Get(ctx, "captcha:sms:"+phone).Result()
	if err != nil || stored != code {
		return false
	}

	// Consume the code
	s.rdb.Del(ctx, "captcha:sms:"+phone)
	return true
}

// ============ 邮箱验证码 ============

// GenerateEmail generates and stores a numeric email code.
func (s *CaptchaService) GenerateEmail(email string) (string, error) {
	ctx := context.Background()

	if exists, _ := s.rdb.Exists(ctx, "captcha:email:rl:"+email).Result(); exists > 0 {
		return "", fmt.Errorf("请稍后再请求验证码")
	}

	code := generateNumericCode(emailCodeLength)

	if err := s.rdb.Set(ctx, "captcha:email:"+email, code, captchaTTL).Err(); err != nil {
		return "", fmt.Errorf("failed to store email code: %w", err)
	}
	if err := s.rdb.Set(ctx, "captcha:email:rl:"+email, "1", rateLimitDuration).Err(); err != nil {
		return "", fmt.Errorf("failed to set rate limit: %w", err)
	}

	return code, nil
}

// VerifyEmail verifies an email code and consumes it (one-time use).
func (s *CaptchaService) VerifyEmail(email, code string) bool {
	ctx := context.Background()

	stored, err := s.rdb.Get(ctx, "captcha:email:"+email).Result()
	if err != nil || stored != code {
		return false
	}

	// Consume the code
	s.rdb.Del(ctx, "captcha:email:"+email)
	return true
}

// ============ 工具函数 ============

// generateNumericCode returns a random numeric string of the given length.
func generateNumericCode(length int) string {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)
	n, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%0*d", length, n.Int64())
}

// ============ 配置管理 ============

// GetCaptchaConfigService 获取配置服务
func (s *CaptchaService) GetCaptchaConfigService() *CaptchaConfigService {
	return s.config
}

// InitDefaultConfigs 初始化默认配置
func (s *CaptchaService) InitDefaultConfigs() error {
	return s.config.InitDefaultConfigs()
}

// CheckAndConsume checks if a code matches the stored value and consumes it.
func (s *CaptchaService) CheckAndConsume(key, code string) bool {
	ctx := context.Background()
	stored, err := s.rdb.Get(ctx, key).Result()
	if err != nil || stored != code {
		return false
	}
	s.rdb.Del(ctx, key)
	return true
}
