package service

import (
	"errors"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// Contact 用户联系人
type Contact struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Name      string         `gorm:"type:varchar(64);not null" json:"name"`
	Phone     string         `gorm:"type:varchar(32)" json:"phone"`
	Email     string         `gorm:"type:varchar(128)" json:"email"`
	Company   string         `gorm:"type:varchar(128)" json:"company"`
	Address   string         `gorm:"type:varchar(512)" json:"address"`
	City      string         `gorm:"type:varchar(64)" json:"city"`
	Province  string         `gorm:"type:varchar(64)" json:"province"`
	Country   string         `gorm:"type:varchar(64);default:'CN'" json:"country"`
	ZipCode   string         `gorm:"type:varchar(16)" json:"zip_code"`
	IsDefault bool           `gorm:"default:false" json:"is_default"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ContactService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewContactService(db *gorm.DB, log *logger.Logger) *ContactService {
	return &ContactService{db: db, log: log}
}

type CreateUserContactRequest struct {
	Name     string `json:"name" binding:"required,max=64"`
	Phone    string `json:"phone" binding:"omitempty,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Company  string `json:"company" binding:"omitempty,max=128"`
	Address  string `json:"address" binding:"omitempty,max=512"`
	City     string `json:"city" binding:"omitempty,max=64"`
	Province string `json:"province" binding:"omitempty,max=64"`
	Country  string `json:"country" binding:"omitempty,max=64"`
	ZipCode  string `json:"zip_code" binding:"omitempty,max=16"`
}

type UpdateUserContactRequest struct {
	Name     string `json:"name" binding:"omitempty,max=64"`
	Phone    string `json:"phone" binding:"omitempty,max=32"`
	Email    string `json:"email" binding:"omitempty,email"`
	Company  string `json:"company" binding:"omitempty,max=128"`
	Address  string `json:"address" binding:"omitempty,max=512"`
	City     string `json:"city" binding:"omitempty,max=64"`
	Province string `json:"province" binding:"omitempty,max=64"`
	Country  string `json:"country" binding:"omitempty,max=64"`
	ZipCode  string `json:"zip_code" binding:"omitempty,max=16"`
}

// Create creates a new contact for a user.
func (s *ContactService) Create(userID uint, req CreateUserContactRequest) (*Contact, error) {
	country := req.Country
	if country == "" {
		country = "CN"
	}

	contact := &Contact{
		UserID:   userID,
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
		Company:  req.Company,
		Address:  req.Address,
		City:     req.City,
		Province: req.Province,
		Country:  country,
		ZipCode:  req.ZipCode,
	}

	if err := s.db.Create(contact).Error; err != nil {
		return nil, err
	}

	s.log.Infof("contact created: user=%d name=%s", userID, req.Name)
	return contact, nil
}

// GetByID returns a contact by ID, scoped to user.
func (s *ContactService) GetByID(userID, contactID uint) (*Contact, error) {
	var contact Contact
	if err := s.db.Where("id = ? AND user_id = ?", contactID, userID).First(&contact).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// GetList returns paginated contacts for a user.
func (s *ContactService) GetList(userID uint, page, pageSize int, keyword string) ([]Contact, int64, error) {
	var contacts []Contact
	var total int64

	query := s.db.Model(&Contact{}).Where("user_id = ?", userID)
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR phone LIKE ? OR email LIKE ?", q, q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("is_default DESC, id DESC").Find(&contacts).Error; err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

// Update updates a contact.
func (s *ContactService) Update(userID, contactID uint, req UpdateUserContactRequest) (*Contact, error) {
	contact, err := s.GetByID(userID, contactID)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Company != "" {
		updates["company"] = req.Company
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.City != "" {
		updates["city"] = req.City
	}
	if req.Province != "" {
		updates["province"] = req.Province
	}
	if req.Country != "" {
		updates["country"] = req.Country
	}
	if req.ZipCode != "" {
		updates["zip_code"] = req.ZipCode
	}

	if len(updates) > 0 {
		if err := s.db.Model(contact).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return s.GetByID(userID, contactID)
}

// Delete soft-deletes a contact.
func (s *ContactService) Delete(userID, contactID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", contactID, userID).Delete(&Contact{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("contact not found")
	}
	return nil
}

// SetDefault marks a contact as default and unsets others.
func (s *ContactService) SetDefault(userID, contactID uint) error {
	_, err := s.GetByID(userID, contactID)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Contact{}).Where("user_id = ?", userID).Update("is_default", false).Error; err != nil {
			return err
		}
		return tx.Model(&Contact{}).Where("id = ? AND user_id = ?", contactID, userID).Update("is_default", true).Error
	})
}

// GetDefault returns the user's default contact.
func (s *ContactService) GetDefault(userID uint) (*Contact, error) {
	var contact Contact
	if err := s.db.Where("user_id = ? AND is_default = true", userID).First(&contact).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &contact, nil
}
