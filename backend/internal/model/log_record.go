package model

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// LogRecord 操作日志记录
type LogRecord struct {
	gorm.Model
	AdminID    uint           `gorm:"index;not null" json:"admin_id"`
	Admin      Admin          `gorm:"foreignKey:AdminID" json:"admin,omitempty"`
	Action     string         `gorm:"type:varchar(64);not null;index" json:"action"`      // 操作类型: create/update/delete/login/export/config/query
	Module     string         `gorm:"type:varchar(64);not null;index" json:"module"`      // 所属模块: user/order/product/ticket/system/payment/menu/rule
	TargetID   uint           `gorm:"index" json:"target_id"`                             // 操作目标ID
	TargetType string         `gorm:"type:varchar(64)" json:"target_type"`                // 操作目标类型
	Title      string         `gorm:"type:varchar(256)" json:"title"`                     // 操作标题
	OldData    datatypes.JSON `gorm:"type:jsonb" json:"old_data"`                         // 变更前数据
	NewData    datatypes.JSON `gorm:"type:jsonb" json:"new_data"`                         // 变更后数据
	IPAddress  string         `gorm:"type:varchar(64);not null" json:"ip_address"`        // IP地址
	UserAgent  string         `gorm:"type:varchar(512)" json:"user_agent"`                // 浏览器UA
	RequestMethod string      `gorm:"type:varchar(10)" json:"request_method"`             // 请求方法
	RequestPath   string      `gorm:"type:varchar(512)" json:"request_path"`              // 请求路径
	Duration   int64          `gorm:"default:0" json:"duration"`                          // 耗时(ms)
	Remark     string         `gorm:"type:text" json:"remark"`                            // 备注
	Status     int16          `gorm:"type:smallint;default:1;not null" json:"status"`     // 1=成功 0=失败
	ErrorMsg   string         `gorm:"type:text" json:"error_msg"`                         // 错误信息
}

// TableName 指定表名
func (LogRecord) TableName() string {
	return "log_records"
}

// LogRecordStat 日志统计
type LogRecordStat struct {
	Date   string `json:"date"`
	Module string `json:"module"`
	Action string `json:"action"`
	Count  int64  `json:"count"`
}

// LogRecordExport 日志导出结构
type LogRecordExport struct {
	ID         uint      `json:"id"`
	AdminName  string    `json:"admin_name"`
	Action     string    `json:"action"`
	Module     string    `json:"module"`
	Title      string    `json:"title"`
	IPAddress  string    `json:"ip_address"`
	Status     int16     `json:"status"`
	Remark     string    `json:"remark"`
	CreatedAt  time.Time `json:"created_at"`
}
