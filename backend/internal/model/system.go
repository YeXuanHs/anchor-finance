package model

import "time"

// SystemInfo 系统信息
type SystemInfo struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Version     string     `gorm:"type:varchar(32)" json:"version"`
	LicenseType int        `gorm:"default:0" json:"license_type"` // 0=免费 1=商业 2=试用
	LicenseKey  string     `gorm:"type:varchar(256)" json:"license_key"`
	LicenseExp  *time.Time `json:"license_exp"`
	LastCheckAt *time.Time `json:"last_check_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// SystemUpdate 系统更新
type SystemUpdate struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Version     string    `gorm:"type:varchar(32);not null" json:"version"`
	Content     string    `gorm:"type:text" json:"content"` // 更新内容
	Type        string    `gorm:"type:varchar(32)" json:"type"` // major/minor/patch/hotfix
	Status      int8      `gorm:"type:smallint;default:0" json:"status"` // 0=待安装 1=已安装 2=已跳过
	DownloadURL string    `gorm:"type:varchar(512)" json:"download_url"`
	FileSize    int64     `gorm:"default:0" json:"file_size"`
	ReleasedAt  time.Time `json:"released_at"`
	InstalledAt *time.Time `json:"installed_at"`
	CreatedAt   time.Time `json:"created_at"`
}
