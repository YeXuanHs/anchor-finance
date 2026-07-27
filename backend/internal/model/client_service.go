package model

import "time"

// ClientService 客户服务实例
type ClientService struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `gorm:"index;not null" json:"user_id"`
	User        User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProductID   uint       `gorm:"index;not null" json:"product_id"`
	Product     Product    `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Name        string     `gorm:"type:varchar(256);not null" json:"name"`
	Status      int16      `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=活跃 2=暂停 3=待开通 4=已终止 5=已过期
	OpenedAt    *time.Time `json:"opened_at"`
	ExpiredAt   *time.Time `gorm:"index" json:"expired_at"`
	AutoRenew   bool       `gorm:"default:false" json:"auto_renew"`
	Remark      string     `gorm:"type:text" json:"remark"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

const (
	ClientServiceActive     int16 = 1
	ClientServiceSuspended  int16 = 2
	ClientServicePending    int16 = 3
	ClientServiceTerminated int16 = 4
	ClientServiceExpired    int16 = 5
)
