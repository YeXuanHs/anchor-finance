package model

import (
	"time"

	"gorm.io/gorm"
)

// ─── anchor_cloud_finance_pro 插件模型 ───

// ACFPFailNotifyEvent 失败通知去重事件
type ACFPFailNotifyEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EventKey  string    `gorm:"type:varchar(128);uniqueIndex;not null" json:"event_key"`
	CreatedAt time.Time `json:"created_at"`
}

// ACFPUpstreamCache 上游快照缓存
type ACFPUpstreamCache struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	HostID    uint      `gorm:"uniqueIndex" json:"host_id"`
	Data      string    `gorm:"type:longtext;comment:上游数据JSON快照" json:"data"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ACFPIPHistory IP变更历史
type ACFPIPHistory struct {
	gorm.Model
	HostID       uint   `gorm:"index" json:"host_id"`
	UserID       uint   `gorm:"index" json:"user_id"`
	OldIP        string `gorm:"type:varchar(64)" json:"old_ip"`
	NewIP        string `gorm:"type:varchar(64)" json:"new_ip"`
	OldAssigned  string `gorm:"type:text" json:"old_assigned"`
	NewAssigned  string `gorm:"type:text" json:"new_assigned"`
	TriggerEvent string `gorm:"type:varchar(32);comment:create/terminate/sync/delete" json:"trigger_event"`
}

// ACFPLimitedSale 限量发售配置
type ACFPLimitedSale struct {
	gorm.Model
	ProductID    uint   `gorm:"uniqueIndex" json:"product_id"`
	MaxQty       int    `gorm:"default:0;comment:最大销售数量，0=不限" json:"max_qty"`
	SoldQty      int    `gorm:"default:0;comment:已售数量" json:"sold_qty"`
	OffsetQty    int    `gorm:"default:0;comment:配额偏移" json:"offset_qty"`
	AutoHide     bool   `gorm:"default:true;comment:售罄自动隐藏" json:"auto_hide"`
	ShowCount    bool   `gorm:"default:true;comment:前台显示剩余数量" json:"show_count"`
	Status       int8   `gorm:"default:1" json:"status"`
}

// ACFPPriceLock 价格锁定配置
type ACFPPriceLock struct {
	gorm.Model
	ProductID   uint    `gorm:"index" json:"product_id"`
	BillingCycle string  `gorm:"type:varchar(32);comment:monthly/quarterly/semiannually/annually/biennially/triennially" json:"billing_cycle"`
	LockedPrice float64 `gorm:"type:decimal(10,2)" json:"locked_price"`
	Status      int8    `gorm:"default:1" json:"status"`
}

// ACFPLog 操作日志
type ACFPLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Module    string    `gorm:"type:varchar(50);index" json:"module"`
	Action    string    `gorm:"type:varchar(50)" json:"action"`
	Target    string    `gorm:"type:varchar(50);comment:host/product/order/user" json:"target"`
	TargetID  uint      `gorm:"default:0" json:"target_id"`
	Content   string    `gorm:"type:text" json:"content"`
	Status    int8      `gorm:"default:1;comment:1=成功 0=失败" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ACFPCronStatus 定时任务状态
type ACFPCronStatus struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CronName  string    `gorm:"type:varchar(100);uniqueIndex" json:"cron_name"`
	LastRun   time.Time `json:"last_run"`
	Duration  int64     `gorm:"comment:执行耗时毫秒" json:"duration"`
	Status    string    `gorm:"type:varchar(16);default:success" json:"status"`
	ErrorMsg  string    `gorm:"type:text" json:"error_msg"`
}

// ACFPCertProConfig 实名认证Pro配置
type ACFPCertProConfig struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	MinAge int  `gorm:"default:16;comment:最低年龄限制" json:"min_age"`
}

// ACFPCertMinor 未成年用户记录
type ACFPCertMinor struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex" json:"user_id"`
	IDCard    string    `gorm:"type:varchar(18)" json:"id_card"`
	Birthday  string    `gorm:"type:varchar(10)" json:"birthday"`
	Age       int       `json:"age"`
	Status    string    `gorm:"type:varchar(16);default:pending;comment:pending/rejected" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ACFPBatchTask 批量修改任务
type ACFPBatchTask struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TaskType  string    `gorm:"type:varchar(50);comment:price/name/status/hidden" json:"task_type"`
	Filter    string    `gorm:"type:longtext;comment:筛选条件JSON" json:"filter"`
	Changes   string    `gorm:"type:longtext;comment:修改内容JSON" json:"changes"`
	Status    string    `gorm:"type:varchar(16);default:pending;comment:pending/done/failed" json:"status"`
	Total     int       `gorm:"default:0" json:"total"`
	Success   int       `gorm:"default:0" json:"success"`
	Failed    int       `gorm:"default:0" json:"failed"`
	ErrorMsg  string    `gorm:"type:text" json:"error_msg"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
