package model

import (
	"time"

	"gorm.io/gorm"
)

// CPUModel CPU型号模型
type CPUModel struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Cores     int            `json:"cores"`
	Threads   int            `json:"threads"`
	Frequency string         `gorm:"size:20" json:"frequency"`
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CPUModel) TableName() string {
	return "cpu_models"
}

// InstanceSpec 实例规格模型
type InstanceSpec struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	CPU       int            `json:"cpu"`
	Memory    int            `json:"memory"` // GB
	Status    string         `gorm:"size:20;default:active" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (InstanceSpec) TableName() string {
	return "instance_specs"
}
