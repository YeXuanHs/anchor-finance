package service

import (
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// BusinessListService 业务列表Pro服务
// 移植自 anchor_cloud_finance_pro 的 AdminBusinessController
type BusinessListService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewBusinessListService(db *gorm.DB, log *logger.Logger) *BusinessListService {
	return &BusinessListService{db: db, log: log}
}

// BusinessFilter 业务筛选条件
type BusinessFilter struct {
	Status         string `form:"status"`
	Keyword        string `form:"keyword"`
	ProductType    string `form:"product_type"`
	ProductID      uint   `form:"product_id"`
	BillingCycle   string `form:"billingcycle"`
	ServerID       int    `form:"server_id"` // -1=上游, 0=全部, >0=指定
	DueFilter      int    `form:"due_filter"` // 0=全部, 1=今日到期, 2=3天内, 3=7天内, 4=30天内, 5=已过期
	Payment        string `form:"payment"`
	DomainFilter   string `form:"domain_filter"`
	IPFilter       string `form:"ip_filter"`
	UsernameFilter string `form:"username_filter"`
	StartTimeFrom  string `form:"start_time_from"`
	StartTimeTo    string `form:"start_time_to"`
	SortField      string `form:"sort_field"`
	SortDir        string `form:"sort_dir"`
	Page           int    `form:"page"`
	PageSize       int    `form:"page_size"`
}

// BusinessItem 业务列表项
type BusinessItem struct {
	ID                uint    `json:"id"`
	UID               uint    `json:"uid"`
	Domain            string  `json:"domain"`
	DedicatedIP       string  `json:"dedicatedip"`
	DomainStatus      string  `json:"domainstatus"`
	NextDueDate       int64   `json:"nextduedate"`
	RegDate           int64   `json:"regdate"`
	Amount            float64 `json:"amount"`
	FirstPaymentAmount float64 `json:"firstpaymentamount"`
	BillingCycle      string  `json:"billingcycle"`
	ServerID          uint    `json:"serverid"`
	Payment           string  `json:"payment"`
	ClientName        string  `json:"client_name"`
	ClientEmail       string  `json:"client_email"`
	ProductName       string  `json:"product_name"`
	ProductType       string  `json:"product_type"`
	APIType           string  `json:"api_type"`
	ZjmfAPIID         uint    `json:"zjmf_api_id"`
	ServerName        string  `json:"server_name"`
	UpstreamName      string  `json:"upstream_name"`
	SaleName          string  `json:"sale_name"`
}

// StatusStats 状态统计
type StatusStats struct {
	All       int64 `json:"all"`
	Active    int64 `json:"Active"`
	Pending   int64 `json:"Pending"`
	Suspended int64 `json:"Suspended"`
	Cancelled int64 `json:"Cancelled"`
	Fraud     int64 `json:"Fraud"`
	Deleted   int64 `json:"Deleted"`
}

// GetDB 获取数据库连接
func (s *BusinessListService) GetDB() *gorm.DB {
	return s.db
}

// GetDB 获取数据库连接
func (s *BusinessListService) GetDB() *gorm.DB {
	return s.db
}

// GetList 获取业务列表（带高级筛选）
func (s *BusinessListService) GetList(filter BusinessFilter) ([]BusinessItem, int64, *StatusStats, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.SortField == "" {
		filter.SortField = "h.id"
	}
	if filter.SortDir == "" {
		filter.SortDir = "DESC"
	}

	// 验证排序字段
	allowedSort := map[string]bool{
		"h.id": true, "h.domainstatus": true, "h.nextduedate": true,
		"h.amount": true, "h.regdate": true, "c.username": true,
		"p.name": true, "h.domain": true,
	}
	if !allowedSort[filter.SortField] {
		filter.SortField = "h.id"
	}

	// 构建查询
	q := s.buildQuery(filter)

	// 获取总数
	var total int64
	q.Count(&total)

	// 获取列表
	var items []BusinessItem
	err := q.Order(fmt.Sprintf("%s %s", filter.SortField, filter.SortDir)).
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Scan(&items).Error
	if err != nil {
		return nil, 0, nil, err
	}

	// 获取状态统计
	stats := s.GetStatusStats(filter)

	return items, total, stats, nil
}

// GetStatusStats 获取状态统计（跟随筛选条件但不跟随状态筛选）
func (s *BusinessListService) GetStatusStats(filter BusinessFilter) *StatusStats {
	stats := &StatusStats{}

	// 临时清除状态筛选
	statusBackup := filter.Status
	filter.Status = ""

	q := s.buildQuery(filter)

	type StatusCount struct {
		DomainStatus string
		Cnt          int64
	}
	var counts []StatusCount
	q.Select("h.domainstatus, count(*) as cnt").Group("h.domainstatus").Scan(&counts)

	for _, sc := range counts {
		switch sc.DomainStatus {
		case "Active":
			stats.Active = sc.Cnt
		case "Pending":
			stats.Pending = sc.Cnt
		case "Suspended":
			stats.Suspended = sc.Cnt
		case "Cancelled":
			stats.Cancelled = sc.Cnt
		case "Fraud":
			stats.Fraud = sc.Cnt
		case "Deleted":
			stats.Deleted = sc.Cnt
		}
		stats.All += sc.Cnt
	}

	// 恢复状态筛选
	filter.Status = statusBackup
	return stats
}

// buildQuery 构建查询（复用）
func (s *BusinessListService) buildQuery(filter BusinessFilter) *gorm.DB {
	q := s.db.Table("host h").
		Select(`h.id, h.uid, h.domain, h.dedicatedip, h.domainstatus, h.nextduedate, h.regdate,
			h.amount, h.firstpaymentamount, h.billingcycle, h.serverid, h.payment,
			c.username as client_name, c.email as client_email,
			p.name as product_name, p.type as product_type, p.api_type, p.zjmf_api_id,
			s.name as server_name,
			u.user_nickname as sale_name`).
		Joins("LEFT JOIN clients c ON h.uid = c.id").
		Joins("LEFT JOIN products p ON h.productid = p.id").
		Joins("LEFT JOIN servers s ON h.serverid = s.id").
		Joins("LEFT JOIN user u ON c.sale_id = u.id")

	if filter.Status != "" {
		q = q.Where("h.domainstatus = ?", filter.Status)
	}
	if filter.Keyword != "" {
		q = q.Where("(c.username LIKE ? OR h.domain LIKE ? OR h.dedicatedip LIKE ? OR c.email LIKE ? OR h.id = ?)",
			"%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%", "%"+filter.Keyword+"%",
			filter.Keyword)
	}
	if filter.ProductType != "" {
		q = q.Where("p.type = ?", filter.ProductType)
	}
	if filter.ProductID > 0 {
		q = q.Where("h.productid = ?", filter.ProductID)
	}
	if filter.BillingCycle != "" {
		q = q.Where("h.billingcycle = ?", filter.BillingCycle)
	}
	if filter.Payment != "" {
		q = q.Where("h.payment LIKE ?", "%"+filter.Payment+"%")
	}
	if filter.DomainFilter != "" {
		q = q.Where("h.domain LIKE ?", "%"+filter.DomainFilter+"%")
	}
	if filter.IPFilter != "" {
		q = q.Where("h.dedicatedip LIKE ?", "%"+filter.IPFilter+"%")
	}
	if filter.UsernameFilter != "" {
		q = q.Where("c.username LIKE ?", "%"+filter.UsernameFilter+"%")
	}
	if filter.ServerID == -1 {
		q = q.Where("p.api_type = ? AND p.zjmf_api_id > 0", "zjmf_api")
	} else if filter.ServerID > 0 {
		q = q.Where("h.serverid = ?", filter.ServerID)
	}

	// 日期范围筛选
	if filter.StartTimeFrom != "" {
		q = q.Where("h.nextduedate >= ?", filter.StartTimeFrom)
	}
	if filter.StartTimeTo != "" {
		q = q.Where("h.nextduedate <= ?", filter.StartTimeTo)
	}

	// 到期筛选
	if filter.DueFilter > 0 {
		today := time.Now().Unix()
		switch filter.DueFilter {
		case 1: // 今日到期
			endOfDay := today + 86400
			q = q.Where("h.nextduedate >= ? AND h.nextduedate <= ?", today, endOfDay)
		case 2: // 3天内
			q = q.Where("h.nextduedate >= ? AND h.nextduedate <= ?", today, today+259200)
		case 3: // 7天内
			q = q.Where("h.nextduedate >= ? AND h.nextduedate <= ?", today, today+604800)
		case 4: // 30天内
			q = q.Where("h.nextduedate >= ? AND h.nextduedate <= ?", today, today+2592000)
		case 5: // 已过期
			q = q.Where("h.nextduedate > 0 AND h.nextduedate < ?", today)
		}
	}

	return q
}

// GetRow 获取单行业务数据
func (s *BusinessListService) GetRow(hostID uint) (*BusinessItem, error) {
	var item BusinessItem
	err := s.db.Table("host h").
		Select(`h.id, h.uid, h.domain, h.dedicatedip, h.domainstatus, h.nextduedate, h.regdate,
			h.amount, h.firstpaymentamount, h.billingcycle, h.serverid, h.payment,
			c.username as client_name, c.email as client_email,
			p.name as product_name, p.type as product_type, p.api_type, p.zjmf_api_id,
			s.name as server_name`).
		Joins("LEFT JOIN clients c ON h.uid = c.id").
		Joins("LEFT JOIN products p ON h.productid = p.id").
		Joins("LEFT JOIN servers s ON h.serverid = s.id").
		Where("h.id = ?", hostID).
		Scan(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetSnapshot 获取业务快照（详细信息）
func (s *BusinessListService) GetSnapshot(hostID uint) (map[string]interface{}, error) {
	var host struct {
		ID                uint    `gorm:"column:id"`
		UID               uint    `gorm:"column:uid"`
		Domain            string  `gorm:"column:domain"`
		DedicatedIP       string  `gorm:"column:dedicatedip"`
		DomainStatus      string  `gorm:"column:domainstatus"`
		NextDueDate       int64   `gorm:"column:nextduedate"`
		RegDate           int64   `gorm:"column:regdate"`
		Amount            float64 `gorm:"column:amount"`
		FirstPaymentAmount float64 `gorm:"column:firstpaymentamount"`
		BillingCycle      string  `gorm:"column:billingcycle"`
		Username          string  `gorm:"column:username"`
		OS                string  `gorm:"column:os"`
		OSURL             string  `gorm:"column:os_url"`
		Remark            string  `gorm:"column:remark"`
		Notes             string  `gorm:"column:notes"`
		ServerID          uint    `gorm:"column:serverid"`
		ProductID         uint    `gorm:"column:productid"`
	}

	if err := s.db.Table("host").Where("id = ?", hostID).First(&host).Error; err != nil {
		return nil, err
	}

	// 获取客户信息
	var client model.User
	s.db.Where("id = ?", host.UID).First(&client)

	// 获取产品信息
	var product model.Product
	s.db.Where("id = ?", host.ProductID).First(&product)

	// 获取服务器信息
	var serverName string
	if host.ServerID > 0 {
		s.db.Table("servers").Where("id = ?", host.ServerID).Select("name").Scan(&serverName)
	}

	cycleMap := map[string]string{
		"free": "免费", "monthly": "月付", "quarterly": "季付",
		"semiannually": "半年付", "annually": "年付",
		"biennially": "两年付", "triennially": "三年付", "onetime": "一次性",
	}
	statusMap := map[string]string{
		"Pending": "待开通", "Active": "已激活", "Suspended": "已暂停",
		"Cancelled": "被取消", "Fraud": "有欺诈", "Deleted": "被删除",
	}

	formatted := map[string]interface{}{
		"业务ID":   host.ID,
		"客户":     fmt.Sprintf("%s (%s)", client.Username, client.Email),
		"商品":     product.Name,
		"主机名":   host.Domain,
		"IP地址":   host.DedicatedIP,
		"状态":     statusMap[host.DomainStatus],
		"计费周期": cycleMap[host.BillingCycle],
		"续费金额": host.Amount,
		"首次付款": host.FirstPaymentAmount,
		"产品类型": product.Type,
		"用户名":   host.Username,
		"系统":     host.OS,
		"系统URL":  host.OSURL,
		"备注":     host.Remark,
		"管理员备注": host.Notes,
	}

	if host.RegDate > 0 {
		formatted["开通时间"] = time.Unix(host.RegDate, 0).Format("2006-01-02 15:04")
	}
	if host.NextDueDate > 0 {
		formatted["到期时间"] = time.Unix(host.NextDueDate, 0).Format("2006-01-02")
	}

	if serverName != "" {
		formatted["开通接口"] = serverName
	}

	return formatted, nil
}

// BatchSync 批量同步（返回任务ID）
func (s *BusinessListService) BatchSync(hostIDs []uint) (string, int, error) {
	taskID := fmt.Sprintf("sync_%d", time.Now().UnixNano())
	total := len(hostIDs)
	// 这里返回任务ID，实际执行由前端轮询驱动
	return taskID, total, nil
}

// LogActivity 记录业务操作日志
func (s *BusinessListService) LogActivity(module, action, targetType string, targetID uint, detail string, success bool, adminID uint) {
	result := int8(0)
	if success {
		result = 1
	}
	s.db.Create(&model.ACFPLog{
		ModuleKey:  module,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Detail:     detail,
		Result:     result,
		AdminID:    adminID,
	})
}
