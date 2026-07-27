package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	captchaTTL        = 5 * time.Minute
	smsCodeLength     = 6
	emailCodeLength   = 6
	rateLimitDuration = 60 * time.Second
)

// CaptchaService handles captcha generation and verification.
type CaptchaService struct {
	rdb *redis.Client
}

func NewCaptchaService(rdb *redis.Client) *CaptchaService {
	return &CaptchaService{rdb: rdb}
}

// GenerateImage generates an image captcha, stores answer in Redis, returns PNG bytes.
func (s *CaptchaService) GenerateImage(key string) ([]byte, error) {
	ctx := context.Background()
	code := generateNumericCode(6)

	if err := s.rdb.Set(ctx, "captcha:image:"+key, code, captchaTTL).Err(); err != nil {
		return nil, fmt.Errorf("failed to store captcha: %w", err)
	}

	// Generate a simple PNG image with the code
	img, err := generateCaptchaImage(code)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// GenerateSMS generates and stores a numeric SMS code.
func (s *CaptchaService) GenerateSMS(phone string) (string, error) {
	ctx := context.Background()

	// Rate limit check
	if exists, _ := s.rdb.Exists(ctx, "captcha:sms:rl:"+phone).Result(); exists > 0 {
		return "", fmt.Errorf("please wait before requesting another code")
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

// GenerateEmail generates and stores a numeric email code.
func (s *CaptchaService) GenerateEmail(email string) (string, error) {
	ctx := context.Background()

	if exists, _ := s.rdb.Exists(ctx, "captcha:email:rl:"+email).Result(); exists > 0 {
		return "", fmt.Errorf("please wait before requesting another code")
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

// Verify checks if the code matches (does NOT consume).
func (s *CaptchaService) Verify(key, code string) bool {
	ctx := context.Background()

	for _, prefix := range []string{"captcha:image:", "captcha:sms:", "captcha:email:"} {
		stored, err := s.rdb.Get(ctx, prefix+key).Result()
		if err == nil && stored == code {
			return true
		}
	}
	return false
}

// CheckAndConsume verifies and deletes the code (one-time use).
func (s *CaptchaService) CheckAndConsume(key, code string) bool {
	ctx := context.Background()

	for _, prefix := range []string{"captcha:image:", "captcha:sms:", "captcha:email:"} {
		redisKey := prefix + key
		stored, err := s.rdb.Get(ctx, redisKey).Result()
		if err == nil && stored == code {
			s.rdb.Del(ctx, redisKey)
			return true
		}
	}
	return false
}

// generateNumericCode returns a random numeric string of the given length.
func generateNumericCode(length int) string {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)
	n, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%0*d", length, n.Int64())
}

// generateCaptchaImage creates a simple PNG image with the given code.
func generateCaptchaImage(code string) ([]byte, error) {
	width, height := 200, 80
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{240, 240, 240, 255}}, image.Point{}, draw.Src)

	// Draw noise lines
	for i := 0; i < 5; i++ {
		x1, _ := rand.Int(rand.Reader, big.NewInt(int64(width)))
		y1, _ := rand.Int(rand.Reader, big.NewInt(int64(height)))
		x2, _ := rand.Int(rand.Reader, big.NewInt(int64(width)))
		y2, _ := rand.Int(rand.Reader, big.NewInt(int64(height)))
		c := color.RGBA{uint8(randInt(100, 200)), uint8(randInt(100, 200)), uint8(randInt(100, 200)), 255}
		drawLine(img, int(x1.Int64()), int(y1.Int64()), int(x2.Int64()), int(y2.Int64()), c)
	}

	// Draw characters
	for i, ch := range code {
		x := 20 + i*30
		y := 40 + randInt(-10, 10)
		c := color.RGBA{uint8(randInt(0, 100)), uint8(randInt(0, 100)), uint8(randInt(0, 100)), 255}
		drawChar(img, ch, x, y, c)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func randInt(min, max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(n.Int64()) + min
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.RGBA) {
	steps := max(abs(x2-x1), abs(y2-y1))
	if steps == 0 {
		return
	}
	for i := 0; i <= steps; i++ {
		x := x1 + (x2-x1)*i/steps
		y := y1 + (y2-y1)*i/steps
		if x >= 0 && x < img.Bounds().Dx() && y >= 0 && y < img.Bounds().Dy() {
			img.Set(x, y, c)
		}
	}
}

func drawChar(img *image.RGBA, ch rune, x, y int, c color.RGBA) {
	// Simple dot-based character rendering
	patterns := map[rune][][2]int{
		'0': {{0, 0}, {1, 0}, {2, 0}, {0, 1}, {2, 1}, {0, 2}, {2, 2}, {0, 3}, {2, 3}, {0, 4}, {1, 4}, {2, 4}},
		'1': {{1, 0}, {1, 1}, {1, 2}, {1, 3}, {1, 4}, {0, 4}},
		'2': {{0, 0}, {1, 0}, {2, 0}, {2, 1}, {0, 2}, {1, 2}, {2, 2}, {0, 3}, {0, 4}, {1, 4}, {2, 4}},
		'3': {{0, 0}, {1, 0}, {2, 0}, {2, 1}, {0, 2}, {1, 2}, {2, 2}, {2, 3}, {0, 4}, {1, 4}, {2, 4}},
		'4': {{0, 0}, {2, 0}, {0, 1}, {2, 1}, {0, 2}, {1, 2}, {2, 2}, {2, 3}, {2, 4}},
		'5': {{0, 0}, {1, 0}, {2, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 2}, {2, 3}, {0, 4}, {1, 4}, {2, 4}},
		'6': {{0, 0}, {1, 0}, {2, 0}, {0, 1}, {0, 2}, {1, 2}, {2, 2}, {0, 3}, {2, 3}, {0, 4}, {1, 4}, {2, 4}},
		'7': {{0, 0}, {1, 0}, {2, 0}, {2, 1}, {2, 2}, {2, 3}, {2, 4}},
		'8': {{0, 0}, {1, 0}, {2, 0}, {0, 1}, {2, 1}, {0, 2}, {1, 2}, {2, 2}, {0, 3}, {2, 3}, {0, 4}, {1, 4}, {2, 4}},
		'9': {{0, 0}, {1, 0}, {2, 0}, {0, 1}, {2, 1}, {0, 2}, {1, 2}, {2, 2}, {2, 3}, {0, 4}, {1, 4}, {2, 4}},
	}

	pattern, ok := patterns[ch]
	if !ok {
		pattern = patterns['0']
	}

	for _, p := range pattern {
		px := x + p[0]*5
		py := y + p[1]*5
		for dx := 0; dx < 5; dx++ {
			for dy := 0; dy < 5; dy++ {
				ix, iy := px+dx, py+dy
				if ix >= 0 && ix < img.Bounds().Dx() && iy >= 0 && iy < img.Bounds().Dy() {
					img.Set(ix, iy, c)
				}
			}
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
