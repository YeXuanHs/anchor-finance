package service

import (
	"errors"
	"math"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/security"
	"anchorfinance/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:128" json:"email"`
	Phone     string         `gorm:"uniqueIndex;size:32" json:"phone"`
	Password  string         `gorm:"size:128;not null" json:"-"`
	Nickname  string         `gorm:"size:64" json:"nickname"`
	Avatar    string         `gorm:"size:256" json:"avatar"`
	Status    int            `gorm:"default:1;comment:1=active 0=disabled" json:"status"`
	Role      string         `gorm:"size:32;default:user;comment:user/admin" json:"role"`
	Balance   float64        `gorm:"type:decimal(12,2);default:0" json:"balance"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewUserService(db *gorm.DB, log *logger.Logger) *UserService {
	return &UserService{db: db, log: log}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty"`
	Nickname string `json:"nickname" binding:"omitempty,max=64"`
}

type LoginRequest struct {
	Account  string `json:"account" binding:"required"` // username, email, or phone
	Password string `json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"omitempty,max=64"`
	Avatar   string `json:"avatar" binding:"omitempty,max=256"`
	Email    string `json:"email" binding:"omitempty,email"`
	Phone    string `json:"phone" binding:"omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
}

// Register creates a new user account.
func (s *UserService) Register(req RegisterRequest) (*User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: string(hashed),
		Nickname: req.Nickname,
		Status:   1,
		Role:     "user",
		Balance:  0,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	s.log.Infof("user registered: %s (id=%d)", user.Username, user.ID)
	return user, nil
}

// Login authenticates by username, email, or phone.
func (s *UserService) Login(req LoginRequest) (*User, error) {
	var user User
	err := s.db.Where("username = ? OR email = ? OR phone = ?", req.Account, req.Account, req.Account).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("invalid credentials")
	}
	if err != nil {
		return nil, err
	}

	if user.Status != 1 {
		return nil, errors.New("account disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	now := time.Now()
	s.db.Model(&user).Update("last_login", &now)

	s.log.Infof("user logged in: %s (id=%d)", user.Username, user.ID)
	return &user, nil
}

// GetByID fetches a user by primary key.
func (s *UserService) GetByID(id uint) (*User, error) {
	var user User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone fetches a user by phone number.
func (s *UserService) GetByPhone(phone string) (*User, error) {
	var user User
	if err := s.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail fetches a user by email.
func (s *UserService) GetByEmail(email string) (*User, error) {
	var user User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByAccessToken finds a user by their access token (from zjmf loginAccessToken).
func (s *UserService) GetUserByAccessToken(accessToken string) (*User, error) {
	// Look up the token in user_access_tokens table
	var token struct {
		UserID    uint   `gorm:"column:user_id"`
		ExpiresAt int64  `gorm:"column:expires_at"`
	}
	if err := s.db.Table("user_access_tokens").Where("token = ? AND status = 1", accessToken).First(&token).Error; err != nil {
		return nil, errors.New("access_token not found")
	}

	// Check if token is expired
	if token.ExpiresAt > 0 && token.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("access_token expired")
	}

	return s.GetByID(token.UserID)
}

// ResetPassword resets a user's password by account (phone or email).
func (s *UserService) ResetPassword(account, newPassword string) error {
	var user User
	if err := s.db.Where("phone = ? OR email = ?", account, account).First(&user).Error; err != nil {
		return errors.New("account not found")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.Model(&user).Update("password", string(hashed)).Error
}

// UpdateProfile updates nickname, avatar, email, phone.
func (s *UserService) UpdateProfile(userID uint, req UpdateProfileRequest) (*User, error) {
	user, err := s.GetByID(userID)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	if err := s.db.Model(user).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(userID)
}

// ChangePassword verifies old password and sets new one.
func (s *UserService) ChangePassword(userID uint, req ChangePasswordRequest) error {
	user, err := s.GetByID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return errors.New("old password incorrect")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.Model(user).Update("password", string(hashed)).Error
}

// UpdateTwoFactorKey updates the two_factor_key field for a user.
func (s *UserService) UpdateTwoFactorKey(userID uint, key string) error {
	return s.db.Model(&User{}).Where("id = ?", userID).Update("two_factor_key", key).Error
}

// ─── API密钥管理（对齐zjmf逻辑：每用户一个密钥，开关/查看/重置） ───

// GetAPISummary 返回用户的API状态、密钥和调用统计
func (s *UserService) GetAPISummary(userID uint, encryptor *security.Encryptor) (gin.H, error) {
	var user model.User
	if err := s.db.Select("id", "api_open", "api_password", "api_create_time").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	if user.APIOpen != 1 {
		return nil, errors.New("API功能未开启")
	}
	// 解密密钥
	decryptedKey := ""
	if user.APIPassword != "" && encryptor != nil {
		if plain, err := encryptor.Decrypt(user.APIPassword); err == nil {
			decryptedKey = plain
		}
	}
	// 主机数
	var hostCount int64
	s.db.Model(&model.Host{}).Where("owner_id = ?", userID).Count(&hostCount)
	// 活跃主机数（status=1 运行中）
	var activeCount int64
	s.db.Model(&model.Host{}).Where("owner_id = ? AND status = ?", userID, 1).Count(&activeCount)
	return gin.H{
		"api_password":    decryptedKey,
		"api_create_time": user.APICreateTime,
		"host_count":      hostCount,
		"active_count":    activeCount,
	}, nil
}

// ToggleAPIOpen 开关用户的API功能
func (s *UserService) ToggleAPIOpen(userID uint, open int8) error {
	updates := map[string]interface{}{"api_open": open}
	if open == 1 {
		now := time.Now()
		updates["api_create_time"] = &now
		// 首次开启时自动生成密钥
		var user model.User
		s.db.Select("api_password").Where("id = ?", userID).First(&user)
		if user.APIPassword == "" {
			encryptor := security.NewEncryptor(s.getJWTSecret())
			encrypted, err := encryptor.Encrypt(model.GenerateRandomPassword(12))
			if err == nil {
				updates["api_password"] = encrypted
			}
		}
	}
	return s.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

// ResetAPIKey 重置用户的API密钥（自动生成新密钥并保存）
func (s *UserService) ResetAPIKey(userID uint, encryptor *security.Encryptor) (string, error) {
	newKey := model.GenerateRandomPassword(12)
	encrypted, err := encryptor.Encrypt(newKey)
	if err != nil {
		return "", err
	}
	if err := s.db.Model(&model.User{}).Where("id = ?", userID).Update("api_password", encrypted).Error; err != nil {
		return "", err
	}
	return newKey, nil
}

// getJWTSecret 从数据库获取JWT密钥（用于AES加密API密码）
func (s *UserService) getJWTSecret() string {
	var config model.SystemConfig
	if err := s.db.Where("`key` = ?", "jwt_secret").First(&config).Error; err != nil {
		return "default-secret-key-change-me"
	}
	return config.Value
}

// APIRequestLog API请求日志模型（用于统计）
// 如果model包中已有此模型则不需要重复定义

// BindPhone binds a phone number to the user account.
func (s *UserService) BindPhone(userID uint, phone string) error {
	// Check if phone is already bound to another user
	var count int64
	s.db.Model(&User{}).Where("phone = ? AND id != ?", phone, userID).Count(&count)
	if count > 0 {
		return errors.New("phone number already bound to another account")
	}
	return s.db.Model(&User{}).Where("id = ?", userID).Update("phone", phone).Error
}

// BindEmail binds an email to the user account.
func (s *UserService) BindEmail(userID uint, email string) error {
	// Check if email is already bound to another user
	var count int64
	s.db.Model(&User{}).Where("email = ? AND id != ?", email, userID).Count(&count)
	if count > 0 {
		return errors.New("email already bound to another account")
	}
	return s.db.Model(&User{}).Where("id = ?", userID).Update("email", email).Error
}

// GetList returns a paginated user list (admin).
func (s *UserService) GetList(page, pageSize int, keyword string) ([]User, int64, error) {
	var users []User
	var total int64

	query := s.db.Model(&User{})
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR email LIKE ? OR phone LIKE ? OR nickname LIKE ?", q, q, q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// Paginate helper (shared).
func Paginate(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	pageSize = int(math.Min(float64(pageSize), 100))
	return (page - 1) * pageSize, pageSize
}

// GetLoginLogs returns paginated login logs for a user.
func (s *UserService) GetLoginLogs(userID uint, page, pageSize int) ([]LoginLog, int64, error) {
	var logs []LoginLog
	var total int64

	query := s.db.Model(&LoginLog{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
