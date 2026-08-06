package service

import (
	"errors"

	"anchorfinance/internal/model"

	"gorm.io/gorm"
)

// ClientContactService manages client contact operations.
type ClientContactService struct {
	db *gorm.DB
}

// NewClientContactService creates a new ClientContactService.
func NewClientContactService(db *gorm.DB) *ClientContactService {
	return &ClientContactService{db: db}
}

// CreateContactRequest is the payload for creating a contact.
type CreateContactRequest struct {
	UserID    uint   `json:"user_id" binding:"required"`
	Name      string `json:"name" binding:"required,max=100"`
	Email     string `json:"email" binding:"omitempty,email"`
	Phone     string `json:"phone" binding:"omitempty,max=32"`
	Position  string `json:"position" binding:"omitempty,max=100"`
	IsDefault bool   `json:"is_default"`
}

// UpdateContactRequest is the payload for updating a contact.
type UpdateContactRequest struct {
	Name      *string `json:"name"`
	Email     *string `json:"email"`
	Phone     *string `json:"phone"`
	Position  *string `json:"position"`
	IsDefault *bool   `json:"is_default"`
}

// Create adds a new contact for a user.
func (s *ClientContactService) Create(req CreateContactRequest) (*model.ClientContact, error) {
	contact := &model.ClientContact{
		UserID:    req.UserID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Position:  req.Position,
		IsDefault: req.IsDefault,
	}
	err := s.db.Transaction(func(tx *gorm.DB) error {
		if req.IsDefault {
			if err := tx.Model(&model.ClientContact{}).
				Where("user_id = ? AND is_default = ?", req.UserID, true).
				Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(contact).Error
	})
	return contact, err
}

// GetByID fetches a contact by ID.
func (s *ClientContactService) GetByID(id uint) (*model.ClientContact, error) {
	var contact model.ClientContact
	if err := s.db.First(&contact, id).Error; err != nil {
		return nil, err
	}
	return &contact, nil
}

// GetByUser returns all contacts for a user.
func (s *ClientContactService) GetByUser(userID uint) ([]model.ClientContact, error) {
	var contacts []model.ClientContact
	if err := s.db.Where("user_id = ?", userID).Order("is_default DESC, id ASC").Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}

// GetList returns a paginated contact list (admin, optionally filtered by user).
func (s *ClientContactService) GetList(page, pageSize int, userID uint) ([]model.ClientContact, int64, error) {
	var items []model.ClientContact
	var total int64

	query := s.db.Model(&model.ClientContact{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("User").Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// Update modifies an existing contact.
func (s *ClientContactService) Update(id uint, req UpdateContactRequest) error {
	contact, err := s.GetByID(id)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if req.IsDefault != nil && *req.IsDefault {
			if err := tx.Model(&model.ClientContact{}).
				Where("user_id = ? AND is_default = ? AND id != ?", contact.UserID, true, id).
				Update("is_default", false).Error; err != nil {
				return err
			}
		}

		updates := map[string]interface{}{}
		if req.Name != nil {
			updates["name"] = *req.Name
		}
		if req.Email != nil {
			updates["email"] = *req.Email
		}
		if req.Phone != nil {
			updates["phone"] = *req.Phone
		}
		if req.Position != nil {
			updates["position"] = *req.Position
		}
		if req.IsDefault != nil {
			updates["is_default"] = *req.IsDefault
		}
		if len(updates) == 0 {
			return nil
		}
		return tx.Model(&model.ClientContact{}).Where("id = ?", id).Updates(updates).Error
	})
}

// Delete removes a contact.
func (s *ClientContactService) Delete(id uint) error {
	result := s.db.Delete(&model.ClientContact{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("contact not found")
	}
	return nil
}
