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
	MergedInto   *uint          `gorm:"index" json:"merged_into"` // 合并目标工单 ID
	Replies      []TicketReply  `gorm:"foreignKey:TicketID" json:"replies,omitempty"`
	Attachments  []Attachment   `gorm:"foreignKey:TicketID" json:"attachments,omitempty"`
	Tags         datatypes.JSON `gorm:"type:json" json:"tags"`
	Metadata     datatypes.JSON `gorm:"type:json" json:"metadata"`
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
	Metadata    datatypes.JSON `gorm:"type:json" json:"metadata"`
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

// TicketTransferLog 工单转移/部门转交记录
type TicketTransferLog struct {
	gorm.Model
	TicketID    uint    `gorm:"index;not null" json:"ticket_id"`
	FromDeptID  *uint   `json:"from_dept_id"`
	ToDeptID    *uint   `json:"to_dept_id"`
	FromAgentID *uint   `json:"from_agent_id"`
	ToAgentID   *uint   `json:"to_agent_id"`
	OperatorID  uint    `gorm:"not null" json:"operator_id"`
	Reason      string  `gorm:"type:text" json:"reason"`
	Ticket      *Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
}

// TicketNote 工单备注（管理员内部备注）
type TicketNote struct {
	gorm.Model
	TicketID    uint   `gorm:"index;not null" json:"ticket_id"`
	Ticket      Ticket `gorm:"foreignKey:TicketID" json:"ticket,omitempty"`
	AdminID     uint   `gorm:"index;not null" json:"admin_id"`
	Admin       Admin  `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Content     string `gorm:"type:text;not null" json:"content"`
	Attachment  string `gorm:"type:text" json:"attachment"` // 逗号分隔的附件路径
}

// TicketCustomField 工单自定义字段
type TicketCustomField struct {
	gorm.Model
	Type        string `gorm:"type:varchar(32);not null;index" json:"type"`       // ticket/order/product
	RelID       uint   `gorm:"index;not null" json:"rel_id"`                     // 关联ID
	FieldName   string `gorm:"type:varchar(128);not null" json:"field_name"`
	FieldType   string `gorm:"type:varchar(32);default:text" json:"field_type"`  // text/textarea/dropdown/password/tickbox
	Description string `gorm:"type:varchar(255)" json:"description"`
	FieldOption string `gorm:"type:text" json:"field_option"`                    // 下拉选项（逗号分隔）
	RegExpr     string `gorm:"type:varchar(255)" json:"reg_expr"`
	AdminOnly   int8   `gorm:"default:0" json:"admin_only"`
	Required    int8   `gorm:"default:0" json:"required"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
}

// TicketCustomFieldValue 工单自定义字段值
type TicketCustomFieldValue struct {
	gorm.Model
	FieldID uint   `gorm:"index;not null" json:"field_id"`
	RelID   uint   `gorm:"index;not null" json:"rel_id"`
	Value   string `gorm:"type:text" json:"value"`
}

// FlowPacket 流量包
type FlowPacket struct {
	gorm.Model
	ServerID uint    `gorm:"index;not null" json:"server_id"`
	Name     string  `gorm:"type:varchar(128);not null" json:"name"`
	Flow     float64 `gorm:"not null" json:"flow"`         // 流量大小（GB）
	Price    float64 `gorm:"not null" json:"price"`
	Status   int8    `gorm:"default:1" json:"status"`       // 1=启用 0=禁用
}
