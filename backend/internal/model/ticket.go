package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Department 工单部门
type Department struct {
	gorm.Model
	Name        string `gorm:"type:varchar(128);not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Slug        string `gorm:"type:varchar(128);uniqueIndex" json:"slug"`
	Email       string `gorm:"type:varchar(255)" json:"email"` // 转发邮箱
	ParentID    *uint  `gorm:"index" json:"parent_id"`
	Parent      *Department `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	Status      int16  `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=启用 0=禁用
	AutoReply   string `gorm:"type:text" json:"auto_reply"`
	Tickets     []Ticket `gorm:"foreignKey:DepartmentID" json:"tickets,omitempty"`
}

// Ticket 工单
type Ticket struct {
	gorm.Model
	TicketNo     string         `gorm:"type:varchar(32);uniqueIndex;not null" json:"ticket_no"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DepartmentID uint           `gorm:"index;not null" json:"department_id"`
	Department   Department     `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	AssignedTo   *uint          `gorm:"index" json:"assigned_to"`
	Assignee     *Admin         `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`
	Subject      string         `gorm:"type:varchar(256);not null" json:"subject"`
	Priority     int16          `gorm:"type:smallint;default:1;not null;index" json:"priority"` // 0=低 1=中 2=高 3=紧急
	Status       int16          `gorm:"type:smallint;default:0;not null;index" json:"status"` // 0=待回复 1=已回复 2=已关闭 3=待处理 4=已解决 5=已取消
	Source       string         `gorm:"type:varchar(32);default:'web'" json:"source"` // web/email/api/phone
	RelType      string         `gorm:"type:varchar(32)" json:"rel_type"` // order/product/user
	RelID        uint           `gorm:"index" json:"rel_id"`
	FirstReplyAt *time.Time `json:"first_reply_at"`
	ClosedAt     *time.Time `json:"closed_at"`
	Satisfaction int16          `gorm:"type:smallint" json:"satisfaction"` // 1-5 评分
	AdminNotes   string         `gorm:"type:text" json:"admin_notes"`
	Replies      []TicketReply  `gorm:"foreignKey:TicketID" json:"replies,omitempty"`
	Attachments  []Attachment   `gorm:"foreignKey:TicketID" json:"attachments,omitempty"`
	Tags         datatypes.JSON `gorm:"type:jsonb" json:"tags"`
	Metadata     datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
}

// TicketReply 工单回复
type TicketReply struct {
	gorm.Model
	TicketID    uint         `gorm:"index;not null" json:"ticket_id"`
	Ticket      Ticket       `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	UserID      *uint        `gorm:"index" json:"user_id"` // 用户回复
	User        *User        `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AdminID     *uint        `gorm:"index" json:"admin_id"` // 管理员回复
	Admin       *Admin       `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Content     string       `gorm:"type:text;not null" json:"content"`
	IsInternal  bool         `gorm:"default:false" json:"is_internal"` // 内部备注
	IsRead      bool         `gorm:"default:false" json:"is_read"`
	ReadAt      *time.Time `json:"read_at"`
	Attachments []Attachment `gorm:"foreignKey:ReplyID" json:"attachments,omitempty"`
	Metadata    datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
}

// Attachment 附件
type Attachment struct {
	gorm.Model
	TicketID    uint   `gorm:"index" json:"ticket_id"`
	Ticket      *Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	ReplyID     *uint  `gorm:"index" json:"reply_id"`
	Reply       *TicketReply `gorm:"foreignKey:ReplyID" json:"reply,omitempty"`
	UploaderID  uint   `gorm:"index;not null" json:"uploader_id"`
	FileName    string `gorm:"type:varchar(256);not null" json:"file_name"`
	FilePath    string `gorm:"type:varchar(512);not null" json:"file_path"`
	FileSize    int64  `gorm:"not null" json:"file_size"`
	MimeType    string `gorm:"type:varchar(128)" json:"mime_type"`
	StorageDriver string `gorm:"type:varchar(32);default:'local'" json:"storage_driver"` // local/s3/oss/cos
	StorageKey  string `gorm:"type:varchar(512)" json:"storage_key"`
	DownloadCount int  `gorm:"default:0" json:"download_count"`
	IsPublic    bool   `gorm:"default:false" json:"is_public"`
	Hash        string `gorm:"type:varchar(64);index" json:"hash"` // SHA256
}
