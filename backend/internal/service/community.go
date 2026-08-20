package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type CommunityService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewCommunityService(db *gorm.DB, log *logger.Logger) *CommunityService {
	return &CommunityService{db: db, log: log}
}

// GetPostList returns paginated community posts.
func (s *CommunityService) GetPostList(page, pageSize int, category string, keyword string) ([]model.CommunityPost, int64, error) {
	var posts []model.CommunityPost
	var total int64

	query := s.db.Model(&model.CommunityPost{}).Where("status >= 0")
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("is_top DESC, id DESC").Find(&posts).Error; err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

// GetPostByID returns a single post by ID.
func (s *CommunityService) GetPostByID(id uint) (*model.CommunityPost, error) {
	var post model.CommunityPost
	if err := s.db.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// CreatePost creates a new community post.
func (s *CommunityService) CreatePost(post *model.CommunityPost) error {
	return s.db.Create(post).Error
}

// UpdatePost updates a community post.
func (s *CommunityService) UpdatePost(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.CommunityPost{}).Where("id = ?", id).Updates(updates).Error
}

// DeletePost soft-deletes a community post.
func (s *CommunityService) DeletePost(id uint) error {
	return s.db.Model(&model.CommunityPost{}).Where("id = ?", id).Update("status", -1).Error
}

// GetComments returns paginated comments for a post.
func (s *CommunityService) GetComments(postID uint, page, pageSize int) ([]model.CommunityComment, int64, error) {
	var comments []model.CommunityComment
	var total int64

	query := s.db.Model(&model.CommunityComment{}).Where("post_id = ? AND status = 1", postID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id ASC").Find(&comments).Error; err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

// CreateComment creates a new comment.
func (s *CommunityService) CreateComment(comment *model.CommunityComment) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		return tx.Model(&model.CommunityPost{}).Where("id = ?", comment.PostID).
			UpdateColumn("reply_count", gorm.Expr("reply_count + 1")).Error
	})
}

// DeleteComment deletes a comment.
func (s *CommunityService) DeleteComment(id uint) error {
	return s.db.Model(&model.CommunityComment{}).Where("id = ?", id).Update("status", 0).Error
}
