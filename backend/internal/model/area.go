package model

// Area 地区
type Area struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ParentID  uint   `gorm:"index;default:0" json:"parent_id"`
	Name      string `gorm:"type:varchar(64);not null" json:"name"`
	Code      string `gorm:"type:varchar(16);index" json:"code"`
	Level     int    `gorm:"default:1" json:"level"` // 1=country 2=province 3=city 4=district
	SortOrder int    `gorm:"default:0" json:"sort_order"`
}
