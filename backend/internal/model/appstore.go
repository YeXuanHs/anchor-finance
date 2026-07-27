package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// App 应用商店应用
type App struct {
	gorm.Model
	Name        string         `gorm:"type:varchar(128);not null;index" json:"name"`
	Slug        string         `gorm:"type:varchar(128);uniqueIndex" json:"slug"`
	Description string         `gorm:"type:text" json:"description"`
	Icon        string         `gorm:"type:varchar(512)" json:"icon"`
	Banner      string         `gorm:"type:varchar(512)" json:"banner"`
	Version     string         `gorm:"type:varchar(32);not null" json:"version"`
	Author      string         `gorm:"type:varchar(128)" json:"author"`
	Website     string         `gorm:"type:varchar(512)" json:"website"`
	Price       datatypes.Decimal `gorm:"type:decimal(20,4);default:0;not null" json:"price"`
	Currency    string         `gorm:"type:varchar(8);default:'CNY';not null" json:"currency"`
	Category    string         `gorm:"type:varchar(64);index" json:"category"`
	Tags        datatypes.JSON `gorm:"type:jsonb" json:"tags"`
	Homepage    string         `gorm:"type:varchar(512)" json:"homepage"`
	Repository  string         `gorm:"type:varchar(512)" json:"repository"`
	MinVersion  string         `gorm:"type:varchar(32)" json:"min_version"`
	MaxVersion  string         `gorm:"type:varchar(32)" json:"max_version"`
	Status      int16          `gorm:"type:smallint;default:1;not null;index" json:"status"` // 1=上架 0=下架 2=审核中
	Downloads   int            `gorm:"default:0" json:"downloads"`
	Rating      float64        `gorm:"type:decimal(3,2);default:0" json:"rating"`
	RatingCount int            `gorm:"default:0" json:"rating_count"`
	SortOrder   int            `gorm:"default:0;index" json:"sort_order"`
	Featured    bool           `gorm:"default:false;index" json:"featured"`
	Verified    bool           `gorm:"default:false" json:"verified"`
	Screenshots datatypes.JSON `gorm:"type:jsonb" json:"screenshots"`
	Metadata    datatypes.JSON `gorm:"type:jsonb" json:"metadata"`
}

// AppInstall 应用安装记录
type AppInstall struct {
	gorm.Model
	UserID  uint   `gorm:"index;not null" json:"user_id"`
	AppID   uint   `gorm:"index;not null" json:"app_id"`
	App     App    `gorm:"foreignKey:AppID" json:"app,omitempty"`
	Version string `gorm:"type:varchar(32);not null" json:"version"`
	Status  int16  `gorm:"type:smallint;default:1;not null" json:"status"` // 1=已安装 2=待更新 3=已卸载
	Config  datatypes.JSON `gorm:"type:jsonb" json:"config"`
}

// AppReview 应用评价
type AppReview struct {
	gorm.Model
	UserID  uint   `gorm:"index;not null" json:"user_id"`
	AppID   uint   `gorm:"index;not null" json:"app_id"`
	App     App    `gorm:"foreignKey:AppID" json:"app,omitempty"`
	Rating  int8   `gorm:"type:smallint;not null" json:"rating"` // 1-5
	Content string `gorm:"type:text" json:"content"`
	Helpful int    `gorm:"default:0" json:"helpful"`
}
