package service

import (
	"errors"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// SystemMessageService manages user-facing system messages/notifications.
type SystemMessageService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewSystemMessageService(db *gorm.DB, log *logger.Logger) *SystemMessageService {
	return &SystemMessageService{db: db, log: log}
}

type SystemMessageInfo struct {
	ID        uint       `json:"id"`
	Type      string     `json:"type"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	Link      string     `json:"link"`
	IsRead    bool       `json:"is_read"`
	ReadAt    *time.Time `json:"read_at"`
	CreatedAt time.Time  `json:"created_at"`
}

// GetList returns paginated system messages for a user.
func (s *SystemMessageService) GetList(userID uint, page, pageSize int, msgType string, onlyUnread bool) ([]SystemMessageInfo, int64, error) {
	var messages []model.SystemMessage
	var total int64

	query := s.db.Model(&model.SystemMessage{}).Where("user_id = ?", userID)
	if msgType != "" {
		query = query.Where("type = ?", msgType)
	}
	if onlyUnread {
		query = query.Where("is_read = false")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	result := make([]SystemMessageInfo, 0, len(messages))
	for _, m := range messages {
		result = append(result, SystemMessageInfo{
			ID:        m.ID,
			Type:      m.Type,
			Title:     m.Title,
			Content:   m.Content,
			Link:      m.Link,
			IsRead:    m.IsRead,
			ReadAt:    m.ReadAt,
			CreatedAt: m.CreatedAt,
		})
	}
	return result, total, nil
}

// GetByID returns a single message by ID.
func (s *SystemMessageService) GetByID(userID, msgID uint) (*SystemMessageInfo, error) {
	var msg model.SystemMessage
	if err := s.db.Where("id = ? AND user_id = ?", msgID, userID).First(&msg).Error; err != nil {
		return nil, err
	}

	return &SystemMessageInfo{
		ID:        msg.ID,
		Type:      msg.Type,
		Title:     msg.Title,
		Content:   msg.Content,
		Link:      msg.Link,
		IsRead:    msg.IsRead,
		ReadAt:    msg.ReadAt,
		CreatedAt: msg.CreatedAt,
	}, nil
}

// MarkRead marks a single message as read.
func (s *SystemMessageService) MarkRead(userID, msgID uint) error {
	now := time.Now()
	result := s.db.Model(&model.SystemMessage{}).
		Where("id = ? AND user_id = ? AND is_read = false", msgID, userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("message not found or already read")
	}
	return nil
}

// MarkAllRead marks all messages for a user as read.
func (s *SystemMessageService) MarkAllRead(userID uint) (int64, error) {
	now := time.Now()
	result := s.db.Model(&model.SystemMessage{}).
		Where("user_id = ? AND is_read = false", userID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		})
	return result.RowsAffected, result.Error
}

// Delete deletes a message for a user.
func (s *SystemMessageService) Delete(userID, msgID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", msgID, userID).Delete(&model.SystemMessage{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("message not found")
	}
	return nil
}

// DeleteAll deletes all messages for a user.
func (s *SystemMessageService) DeleteAll(userID uint) (int64, error) {
	result := s.db.Where("user_id = ?", userID).Delete(&model.SystemMessage{})
	return result.RowsAffected, result.Error
}

// GetUnreadCount returns the number of unread messages for a user.
func (s *SystemMessageService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	if err := s.db.Model(&model.SystemMessage{}).Where("user_id = ? AND is_read = false", userID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetTypes returns distinct message types for a user.
func (s *SystemMessageService) GetTypes(userID uint) ([]string, error) {
	var types []string
	if err := s.db.Model(&model.SystemMessage{}).Where("user_id = ?", userID).Distinct("type").Pluck("type", &types).Error; err != nil {
		return nil, err
	}
	return types, nil
}
