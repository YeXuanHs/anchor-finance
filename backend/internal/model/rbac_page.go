package model

import (
	"gorm.io/gorm"
)

// RbacPage RBAC页面权限
type RbacPage struct {
	gorm.Model
	Name     string `gorm:"type:varchar(128);not null" json:"name"`
	Path     string `gorm:"type:varchar(256);not null;index" json:"path"` // 页面路径
	Method   string `gorm:"type:varchar(16);default:GET" json:"method"` // GET/POST/PUT/DELETE
	ParentID *uint  `gorm:"index" json:"parent_id"`
	Parent   *RbacPage `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []RbacPage `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Module   string `gorm:"type:varchar(64);index" json:"module"` // 所属模块
	Status   int16  `gorm:"type:smallint;default:1;not null;index" json:"status"`
	Remark   string `gorm:"type:text" json:"remark"`
	SortOrder int   `gorm:"default:0" json:"sort_order"`
}

// RbacPageAuth 页面权限规则
type RbacPageAuth struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	PageID   uint   `gorm:"index;not null" json:"page_id"`
	RoleID   uint   `gorm:"index;not null" json:"role_id"`
	CanView  bool   `gorm:"default:false" json:"can_view"`
	CanEdit  bool   `gorm:"default:false" json:"can_edit"`
	CanDelete bool  `gorm:"default:false" json:"can_delete"`
}
