package model

import (
	"time"

	"gorm.io/gorm"
)

// CommunityPost 社区帖子
type CommunityPost struct {
	gorm.Model
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Title     string         `gorm:"type:varchar(256);not null" json:"title"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	Category  string         `gorm:"type:varchar(64);index" json:"category"`
	Tags      string         `gorm:"type:varchar(256)" json:"tags"`
	Status    int16          `gorm:"type:smallint;default:1;index" json:"status"` // 1=正常 0=隐藏 -1=删除
	ViewCount int            `gorm:"default:0" json:"view_count"`
	LikeCount int            `gorm:"default:0" json:"like_count"`
	IsTop     bool           `gorm:"default:false" json:"is_top"`
	IsHot     bool           `gorm:"default:false" json:"is_hot"`
	LastReplyAt *time.Time   `json:"last_reply_at"`
	ReplyCount int           `gorm:"default:0" json:"reply_count"`
}

// CommunityComment 社区评论
type CommunityComment struct {
	gorm.Model
	PostID    uint   `gorm:"index;not null" json:"post_id"`
	UserID    uint   `gorm:"index;not null" json:"user_id"`
	ParentID  *uint  `gorm:"index" json:"parent_id"` // 父评论ID，支持嵌套回复
	Content   string `gorm:"type:text;not null" json:"content"`
	Status    int16  `gorm:"type:smallint;default:1;index" json:"status"` // 1=正常 0=隐藏
	LikeCount int    `gorm:"default:0" json:"like_count"`
}
