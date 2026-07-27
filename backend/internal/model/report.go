package model

import (
	"time"

	"gorm.io/gorm"
)

// DailyReport 每日统计报表
type DailyReport struct {
	gorm.Model
	Date           string  `gorm:"type:varchar(10);uniqueIndex" json:"date"` // 2024-01-01
	NewUsers       int     `gorm:"default:0" json:"new_users"`
	NewOrders      int     `gorm:"default:0" json:"new_orders"`
	PaidOrders     int     `gorm:"default:0" json:"paid_orders"`
	Revenue        float64 `gorm:"type:decimal(12,2);default:0" json:"revenue"`
	Refunds        float64 `gorm:"type:decimal(12,2);default:0" json:"refunds"`
	NewTickets     int     `gorm:"default:0" json:"new_tickets"`
	ActiveProducts int     `gorm:"default:0" json:"active_products"`
}

// MonthlyReport 月度统计
type MonthlyReport struct {
	gorm.Model
	Month      string  `gorm:"type:varchar(7);uniqueIndex" json:"month"` // 2024-01
	NewUsers   int     `gorm:"default:0" json:"new_users"`
	TotalUsers int     `gorm:"default:0" json:"total_users"`
	NewOrders  int     `gorm:"default:0" json:"new_orders"`
	PaidOrders int     `gorm:"default:0" json:"paid_orders"`
	Revenue    float64 `gorm:"type:decimal(12,2);default:0" json:"revenue"`
	Expenses   float64 `gorm:"type:decimal(12,2);default:0" json:"expenses"`
	Profit     float64 `gorm:"type:decimal(12,2);default:0" json:"profit"`
}

// Invoice 账单（报表统计用）
type InvoiceReport struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Status    int16     `gorm:"type:smallint" json:"status"`
	Amount    float64   `gorm:"type:decimal(20,4)" json:"amount"`
	PaidAt    *time.Time `json:"paid_at"`
	CreatedAt time.Time `json:"created_at"`
}
