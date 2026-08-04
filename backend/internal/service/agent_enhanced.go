package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// AgentEnhancedService provides enhanced agent operations.
type AgentEnhancedService struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewAgentEnhancedService creates a new AgentEnhancedService.
func NewAgentEnhancedService(db *gorm.DB, log *logger.Logger) *AgentEnhancedService {
	return &AgentEnhancedService{db: db, log: log}
}

// GetAgentByUserID returns an agent by user ID.
func (s *AgentEnhancedService) GetAgentByUserID(userID uint) (*model.Agent, error) {
	var agent model.Agent
	if err := s.db.Where("user_id = ?", userID).First(&agent).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// ==================== Resource Management ====================

// GetResourceInfo returns aggregated resource info for an agent.
func (s *AgentEnhancedService) GetResourceInfo(agentID uint) (map[string]interface{}, error) {
	var agent model.Agent
	if err := s.db.First(&agent, agentID).Error; err != nil {
		return nil, err
	}

	var totalOrders int64
	s.db.Model(&model.Order{}).Where("user_id = ?", agent.UserID).Count(&totalOrders)

	var activeHosts int64
	s.db.Model(&model.Host{}).Where("owner_id = ? AND status = 1", agent.UserID).Count(&activeHosts)

	var totalProducts int64
	s.db.Model(&model.Product{}).Where("status = 1").Count(&totalProducts)

	return map[string]interface{}{
		"agent_id":       agent.ID,
		"balance":        agent.Balance,
		"total_earned":   agent.TotalEarned,
		"total_orders":   totalOrders,
		"active_hosts":   activeHosts,
		"total_products": totalProducts,
	}, nil
}

// PostResourceInfo updates resource-related info for an agent.
func (s *AgentEnhancedService) PostResourceInfo(agentID uint, data map[string]interface{}) error {
	var agent model.Agent
	if err := s.db.First(&agent, agentID).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{}
	if v, ok := data["commission_rate"].(float64); ok {
		updates["commission_rate"] = v
	}
	if v, ok := data["status"].(float64); ok {
		updates["status"] = int8(v)
	}
	if len(updates) == 0 {
		return nil
	}
	return s.db.Model(&agent).Updates(updates).Error
}

// GetProducts returns products available to an agent.
func (s *AgentEnhancedService) GetProducts(agentID uint) ([]model.Product, error) {
	var products []model.Product
	if err := s.db.Where("status = 1").Order("sort_order ASC, id DESC").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

// GetHostLists returns hosts associated with an agent's downline users.
func (s *AgentEnhancedService) GetHostLists(agentID uint) ([]model.Host, error) {
	var agent model.Agent
	if err := s.db.First(&agent, agentID).Error; err != nil {
		return nil, err
	}

	var hosts []model.Host
	if err := s.db.Where("owner_id = ?", agent.UserID).
		Preload("Product").
		Order("id DESC").
		Find(&hosts).Error; err != nil {
		return nil, err
	}
	return hosts, nil
}

// ==================== Inspection ====================

// GetInspectionLists returns paginated inspection records.
func (s *AgentEnhancedService) GetInspectionLists(agentID uint, page, pageSize int) ([]model.AgentInspection, int64, error) {
	var items []model.AgentInspection
	var total int64

	query := s.db.Model(&model.AgentInspection{}).Where("agent_id = ?", agentID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Preload("Host").Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// CreateInspection creates a new inspection record.
func (s *AgentEnhancedService) CreateInspection(agentID, hostID uint, inspectionType string) (*model.AgentInspection, error) {
	inspection := &model.AgentInspection{
		AgentID: agentID,
		HostID:  hostID,
		Type:    inspectionType,
		Status:  "pending",
	}
	if err := s.db.Create(inspection).Error; err != nil {
		return nil, err
	}
	return inspection, nil
}

// PostUpload uploads images for an inspection.
func (s *AgentEnhancedService) PostUpload(inspectionID uint, imageUrls []string) error {
	var inspection model.AgentInspection
	if err := s.db.First(&inspection, inspectionID).Error; err != nil {
		return err
	}

	var existing []string
	if inspection.Images != nil {
		_ = json.Unmarshal(inspection.Images, &existing)
	}
	existing = append(existing, imageUrls...)

	b, err := json.Marshal(existing)
	if err != nil {
		return err
	}
	return s.db.Model(&inspection).Update("images", datatypes.JSON(b)).Error
}

// GetInspectionDetail returns a single inspection record.
func (s *AgentEnhancedService) GetInspectionDetail(inspectionID uint) (*model.AgentInspection, error) {
	var inspection model.AgentInspection
	if err := s.db.Preload("Host").First(&inspection, inspectionID).Error; err != nil {
		return nil, err
	}
	return &inspection, nil
}

// GetInspectionIPs returns unique IPs from hosts for inspection selection.
func (s *AgentEnhancedService) GetInspectionIPs(agentID uint) ([]map[string]interface{}, error) {
	var agent model.Agent
	if err := s.db.First(&agent, agentID).Error; err != nil {
		return nil, err
	}

	var hosts []model.Host
	if err := s.db.Where("owner_id = ?", agent.UserID).Find(&hosts).Error; err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for _, h := range hosts {
		if h.IP != "" {
			results = append(results, map[string]interface{}{
				"host_id":  h.ID,
				"hostname": h.Hostname,
				"ip":       h.IP,
			})
		}
	}
	return results, nil
}

// ==================== Orders & Finance ====================

// OrderFilter holds filters for order search.
type OrderFilter struct {
	Keyword string
	Status  string
	Type    string
	StartAt string
	EndAt   string
}

// GetOrderSearchPage searches orders with filters.
func (s *AgentEnhancedService) GetOrderSearchPage(agentID uint, page, pageSize int, filter OrderFilter) ([]model.Order, int64, error) {
	var agent model.Agent
	if err := s.db.First(&agent, agentID).Error; err != nil {
		return nil, 0, err
	}

	query := s.db.Model(&model.Order{}).Where("user_id = ?", agent.UserID)
	if filter.Keyword != "" {
		q := "%" + filter.Keyword + "%"
		query = query.Where("order_no LIKE ? OR description LIKE ?", q, q)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.StartAt != "" {
		query = query.Where("created_at >= ?", filter.StartAt)
	}
	if filter.EndAt != "" {
		query = query.Where("created_at <= ?", filter.EndAt)
	}

	var total int64
	query.Count(&total)

	var orders []model.Order
	offset := (page - 1) * pageSize
	if err := query.Preload("Product").Offset(offset).Limit(pageSize).Order("id DESC").Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// GetOrders returns paginated orders for an agent.
func (s *AgentEnhancedService) GetOrders(agentID uint, page, pageSize int) ([]model.Order, int64, error) {
	var agent model.Agent
	if err := s.db.First(&agent, agentID).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	s.db.Model(&model.Order{}).Where("user_id = ?", agent.UserID).Count(&total)

	var orders []model.Order
	offset := (page - 1) * pageSize
	if err := s.db.Where("user_id = ?", agent.UserID).
		Preload("Product").
		Offset(offset).Limit(pageSize).
		Order("id DESC").
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// GetRenewSearchPage searches renewal orders with filters.
func (s *AgentEnhancedService) GetRenewSearchPage(agentID uint, page, pageSize int, filter OrderFilter) ([]model.Order, int64, error) {
	filter.Type = "renew"
	return s.GetOrderSearchPage(agentID, page, pageSize, filter)
}

// GetRenews returns paginated renewal orders for an agent.
func (s *AgentEnhancedService) GetRenews(agentID uint, page, pageSize int) ([]model.Order, int64, error) {
	var agent model.Agent
	if err := s.db.First(&agent, agentID).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	s.db.Model(&model.Order{}).Where("user_id = ? AND type = 'renew'", agent.UserID).Count(&total)

	var orders []model.Order
	offset := (page - 1) * pageSize
	if err := s.db.Where("user_id = ? AND type = 'renew'", agent.UserID).
		Preload("Product").
		Offset(offset).Limit(pageSize).
		Order("id DESC").
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}
	return orders, total, nil
}

// GetIncome returns income statistics for an agent within a period.
func (s *AgentEnhancedService) GetIncome(agentID uint, period string) (map[string]interface{}, error) {
	var agent model.Agent
	if err := s.db.First(&agent, agentID).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	var startTime time.Time
	switch period {
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	case "year":
		startTime = now.AddDate(-1, 0, 0)
	default:
		startTime = now.AddDate(0, -1, 0)
	}

	var totalIncome float64
	s.db.Model(&model.AgentCommission{}).
		Where("agent_id = ? AND status = 2 AND created_at >= ?", agentID, startTime).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	var totalOrders int64
	s.db.Model(&model.Order{}).
		Where("user_id = ? AND paid_at >= ?", agent.UserID, startTime).
		Count(&totalOrders)

	var commissionCount int64
	s.db.Model(&model.AgentCommission{}).
		Where("agent_id = ? AND created_at >= ?", agentID, startTime).
		Count(&commissionCount)

	return map[string]interface{}{
		"period":           period,
		"start_time":       startTime,
		"total_income":     totalIncome,
		"total_orders":     totalOrders,
		"commission_count": commissionCount,
		"balance":          agent.Balance,
		"total_earned":     agent.TotalEarned,
	}, nil
}

// GetConsumption returns consumption statistics for an agent.
func (s *AgentEnhancedService) GetConsumption(agentID uint, period string) (map[string]interface{}, error) {
	now := time.Now()
	var startTime time.Time
	switch period {
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	case "year":
		startTime = now.AddDate(-1, 0, 0)
	default:
		startTime = now.AddDate(0, -1, 0)
	}

	var totalAmount float64
	s.db.Model(&model.AgentConsumption{}).
		Where("agent_id = ? AND created_at >= ?", agentID, startTime).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalAmount)

	var totalCommission float64
	s.db.Model(&model.AgentConsumption{}).
		Where("agent_id = ? AND created_at >= ?", agentID, startTime).
		Select("COALESCE(SUM(commission), 0)").
		Scan(&totalCommission)

	var totalCount int64
	s.db.Model(&model.AgentConsumption{}).
		Where("agent_id = ? AND created_at >= ?", agentID, startTime).
		Count(&totalCount)

	return map[string]interface{}{
		"period":          period,
		"start_time":      startTime,
		"total_amount":    totalAmount,
		"total_commission": totalCommission,
		"total_count":     totalCount,
	}, nil
}

// GetAgentLogs returns paginated agent activity logs.
func (s *AgentEnhancedService) GetAgentLogs(agentID uint, page, pageSize int) ([]model.AgentLog, int64, error) {
	var items []model.AgentLog
	var total int64

	query := s.db.Model(&model.AgentLog{}).Where("agent_id = ?", agentID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// ==================== After-sale ====================

// GetAfterSaleDetail returns after-sale detail.
func (s *AgentEnhancedService) GetAfterSaleDetail(id uint) (*model.AgentAfterSale, error) {
	var item model.AgentAfterSale
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// GetRefundDetail returns refund detail (alias for after-sale with type=refund).
func (s *AgentEnhancedService) GetRefundDetail(id uint) (*model.AgentAfterSale, error) {
	var item model.AgentAfterSale
	if err := s.db.Where("id = ? AND type = 'refund'", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// PostAfterSale processes an after-sale request.
func (s *AgentEnhancedService) PostAfterSale(id uint, data map[string]interface{}) error {
	var item model.AgentAfterSale
	if err := s.db.First(&item, id).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{}
	if v, ok := data["status"].(string); ok {
		updates["status"] = v
	}
	if v, ok := data["reason"].(string); ok {
		updates["reason"] = v
	}
	if v, ok := data["amount"].(float64); ok {
		updates["amount"] = v
	}
	if len(updates) == 0 {
		return nil
	}
	return s.db.Model(&item).Updates(updates).Error
}

// PostUnAfterSale cancels an after-sale request.
func (s *AgentEnhancedService) PostUnAfterSale(id uint) error {
	return s.db.Model(&model.AgentAfterSale{}).Where("id = ?", id).
		Update("status", "cancelled").Error
}

// PostRefund processes a refund.
func (s *AgentEnhancedService) PostRefund(id uint, data map[string]interface{}) error {
	var item model.AgentAfterSale
	if err := s.db.Where("id = ? AND type = 'refund'", id).First(&item).Error; err != nil {
		return err
	}

	updates := map[string]interface{}{
		"status":       "approved",
		"processed_at": time.Now(),
	}
	if v, ok := data["amount"].(float64); ok {
		updates["amount"] = v
	}
	return s.db.Model(&item).Updates(updates).Error
}

// ==================== Tickets ====================

// GetTickets returns paginated tickets for an agent's user.
func (s *AgentEnhancedService) GetTickets(agentID uint, page, pageSize int) ([]model.Ticket, int64, error) {
	var agent model.Agent
	if err := s.db.First(&agent, agentID).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	s.db.Model(&model.Ticket{}).Where("user_id = ?", agent.UserID).Count(&total)

	var tickets []model.Ticket
	offset := (page - 1) * pageSize
	if err := s.db.Where("user_id = ?", agent.UserID).
		Preload("Department").
		Offset(offset).Limit(pageSize).
		Order("id DESC").
		Find(&tickets).Error; err != nil {
		return nil, 0, err
	}
	return tickets, total, nil
}

// ==================== Evaluation ====================

// PostEvaluation submits an evaluation for an order.
func (s *AgentEnhancedService) PostEvaluation(agentID uint, data map[string]interface{}) error {
	orderID, _ := data["order_id"].(float64)
	rating, _ := data["rating"].(float64)
	comment, _ := data["comment"].(string)

	if orderID == 0 || rating == 0 {
		return gorm.ErrInvalidData
	}

	log := &model.AgentLog{
		AgentID:  agentID,
		Action:   "evaluation",
		Target:   "order",
		TargetID: uint(orderID),
		Detail:   comment,
	}
	return s.db.Create(log).Error
}

// GetRunMapLists returns run map entries.
func (s *AgentEnhancedService) GetRunMapLists(agentID uint) ([]model.RunMap, error) {
	var items []model.RunMap
	if err := s.db.Where("is_enabled = true").Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// ==================== Token Management ====================

// GetToken returns the active API token for an agent.
func (s *AgentEnhancedService) GetToken(agentID uint) (*model.AgentToken, error) {
	var token model.AgentToken
	if err := s.db.Where("agent_id = ? AND enabled = true", agentID).Order("id DESC").First(&token).Error; err != nil {
		return nil, err
	}
	return &token, nil
}

// SetToken creates or updates an API token for an agent.
func (s *AgentEnhancedService) SetToken(agentID uint, name string) (*model.AgentToken, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	tokenStr := hex.EncodeToString(tokenBytes)

	token := &model.AgentToken{
		AgentID:   agentID,
		Token:     tokenStr,
		Name:      name,
		ExpiresAt: time.Now().AddDate(1, 0, 0), // 1 year
		Enabled:   true,
	}
	if err := s.db.Create(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}

// CheckToken validates an API token and returns the associated agent.
func (s *AgentEnhancedService) CheckToken(tokenStr string) (*model.Agent, error) {
	var token model.AgentToken
	if err := s.db.Where("token = ? AND enabled = true", tokenStr).First(&token).Error; err != nil {
		return nil, err
	}

	if token.ExpiresAt.Before(time.Now()) {
		return nil, gorm.ErrRecordNotFound
	}

	// Update last used time
	now := time.Now()
	s.db.Model(&token).Update("last_used_at", &now)

	var agent model.Agent
	if err := s.db.First(&agent, token.AgentID).Error; err != nil {
		return nil, err
	}
	return &agent, nil
}

// ==================== Base Info ====================

// GetBaseInfo returns base info for an agent.
func (s *AgentEnhancedService) GetBaseInfo(agentID uint) (map[string]interface{}, error) {
	var agent model.Agent
	if err := s.db.Preload("User").Preload("Parent").First(&agent, agentID).Error; err != nil {
		return nil, err
	}

	var downlineCount int64
	s.db.Model(&model.Agent{}).Where("parent_id = ?", agentID).Count(&downlineCount)

	var orderCount int64
	s.db.Model(&model.Order{}).Where("user_id = ?", agent.UserID).Count(&orderCount)

	return map[string]interface{}{
		"agent":          agent,
		"downline_count": downlineCount,
		"order_count":    orderCount,
	}, nil
}

// ==================== Resource Pool ====================

// ResourcePool 资源池
type ResourcePool struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:128;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Type        string    `gorm:"size:32;default:shared;comment:shared/exclusive" json:"type"`
	TotalCPU    int       `json:"total_cpu"`
	TotalMemory int       `json:"total_memory"` // MB
	TotalDisk   int       `json:"total_disk"`   // GB
	UsedCPU     int       `json:"used_cpu"`
	UsedMemory  int       `json:"used_memory"`
	UsedDisk    int       `json:"used_disk"`
	AgentIDs    string    `gorm:"type:text" json:"agent_ids"` // JSON array of agent IDs
	Status      string    `gorm:"size:32;default:active;comment:active/disabled" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetResourcePools returns paginated resource pools.
func (s *AgentEnhancedService) GetResourcePools(page, pageSize int) ([]ResourcePool, int64, error) {
	var items []ResourcePool
	var total int64

	query := s.db.Model(&ResourcePool{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// CreateResourcePool creates a new resource pool.
func (s *AgentEnhancedService) CreateResourcePool(pool *ResourcePool) error {
	return s.db.Create(pool).Error
}

// UpdateResourcePool updates a resource pool.
func (s *AgentEnhancedService) UpdateResourcePool(id uint, updates map[string]interface{}) error {
	return s.db.Model(&ResourcePool{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteResourcePool deletes a resource pool.
func (s *AgentEnhancedService) DeleteResourcePool(id uint) error {
	return s.db.Delete(&ResourcePool{}, id).Error
}

// ==================== Agent Level ====================

// AgentLevel 代理商等级
type AgentLevel struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:64;not null" json:"name"`
	Level          int       `gorm:"uniqueIndex;not null" json:"level"`
	CommissionRate float64   `gorm:"type:decimal(5,2);default:0" json:"commission_rate"`
	MinSales       float64   `gorm:"type:decimal(12,2);default:0" json:"min_sales"`
	MaxSubAgents   int       `gorm:"default:0" json:"max_sub_agents"`
	Description    string    `gorm:"type:text" json:"description"`
	Discount       float64   `gorm:"type:decimal(5,2);default:100" json:"discount"`
	Status         string    `gorm:"size:32;default:active;comment:active/disabled" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// GetAgentLevels returns all agent levels.
func (s *AgentEnhancedService) GetAgentLevels() ([]AgentLevel, error) {
	var items []AgentLevel
	if err := s.db.Order("level ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// CreateAgentLevel creates a new agent level.
func (s *AgentEnhancedService) CreateAgentLevel(level *AgentLevel) error {
	return s.db.Create(level).Error
}

// UpdateAgentLevel updates an agent level.
func (s *AgentEnhancedService) UpdateAgentLevel(id uint, updates map[string]interface{}) error {
	return s.db.Model(&AgentLevel{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteAgentLevel deletes an agent level.
func (s *AgentEnhancedService) DeleteAgentLevel(id uint) error {
	return s.db.Delete(&AgentLevel{}, id).Error
}

// ==================== Performance Report ====================

// GetPerformanceReport returns performance statistics for agents.
func (s *AgentEnhancedService) GetPerformanceReport(page, pageSize int, period string) ([]map[string]interface{}, int64, error) {
	var agents []model.Agent
	var total int64

	query := s.db.Model(&model.Agent{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("total_earned DESC").Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	now := time.Now()
	var startTime time.Time
	switch period {
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	case "year":
		startTime = now.AddDate(-1, 0, 0)
	default:
		startTime = now.AddDate(0, -1, 0)
	}

	var results []map[string]interface{}
	for _, agent := range agents {
		var periodSales float64
		s.db.Model(&model.Order{}).
			Where("user_id = ? AND paid_at >= ?", agent.UserID, startTime).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&periodSales)

		var periodCommission float64
		s.db.Model(&model.AgentCommission{}).
			Where("agent_id = ? AND status = 2 AND created_at >= ?", agent.ID, startTime).
			Select("COALESCE(SUM(amount), 0)").
			Scan(&periodCommission)

		var orderCount int64
		s.db.Model(&model.Order{}).
			Where("user_id = ? AND paid_at >= ?", agent.UserID, startTime).
			Count(&orderCount)

		results = append(results, map[string]interface{}{
			"agent_id":          agent.ID,
			"level":             agent.Level,
			"total_earned":      agent.TotalEarned,
			"period_sales":      periodSales,
			"period_commission": periodCommission,
			"period_orders":     orderCount,
		})
	}

	return results, total, nil
}

// GetPerformanceChart returns chart data for agent performance.
func (s *AgentEnhancedService) GetPerformanceChart(period string) (map[string]interface{}, error) {
	now := time.Now()
	var startTime time.Time
	var format string

	switch period {
	case "week":
		startTime = now.AddDate(0, 0, -7)
		format = "2006-01-02"
	case "month":
		startTime = now.AddDate(0, -1, 0)
		format = "2006-01-02"
	case "year":
		startTime = now.AddDate(-1, 0, 0)
		format = "2006-01"
	default:
		startTime = now.AddDate(0, -1, 0)
		format = "2006-01-02"
	}

	type dailyData struct {
		Date     string  `json:"date"`
		Sales    float64 `json:"sales"`
		Orders   int64   `json:"orders"`
		NewAgents int64  `json:"new_agents"`
	}

	var results []dailyData

	// Aggregate orders by date
	rows, err := s.db.Model(&model.Order{}).
		Where("paid_at >= ?", startTime).
		Select("DATE_FORMAT(paid_at, ? as date, COALESCE(SUM(amount), 0) as sales, COUNT(*) as orders", format).
		Group("date").
		Order("date ASC").
		Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dateMap := make(map[string]*dailyData)
	for rows.Next() {
		var d dailyData
		if err := rows.Scan(&d.Date, &d.Sales, &d.Orders); err != nil {
			continue
		}
		dateMap[d.Date] = &d
		results = append(results, d)
	}

	// Aggregate new agents by date
	agentRows, err := s.db.Model(&model.Agent{}).
		Where("created_at >= ?", startTime).
		Select("DATE_FORMAT(created_at, ? as date, COUNT(*) as cnt", format).
		Group("date").
		Order("date ASC").
		Rows()
	if err == nil {
		defer agentRows.Close()
		for agentRows.Next() {
			var date string
			var cnt int64
			if err := agentRows.Scan(&date, &cnt); err != nil {
				continue
			}
			if existing, ok := dateMap[date]; ok {
				existing.NewAgents = cnt
			}
		}
	}

	return map[string]interface{}{
		"period": period,
		"data":   results,
	}, nil
}

// ==================== Commission Settlement ====================

// CommissionSettlement 佣金结算记录
type CommissionSettlement struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AgentID   uint      `gorm:"index;not null" json:"agent_id"`
	Amount    float64   `gorm:"type:decimal(12,2);not null" json:"amount"`
	Status    string    `gorm:"size:32;default:pending;comment:pending/paid/rejected" json:"status"`
	Period    string    `gorm:"size:32" json:"period"` // e.g. "2024-01"
	Remark    string    `gorm:"type:text" json:"remark"`
	PaidAt    *time.Time `json:"paid_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CommissionRule 佣金规则
type CommissionRule struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:128;not null" json:"name"`
	Type           string    `gorm:"size:32;default:percentage;comment:percentage/fixed" json:"type"`
	Rate           float64   `gorm:"type:decimal(8,2)" json:"rate"`
	MinAmount      float64   `gorm:"type:decimal(12,2);default:0" json:"min_amount"`
	MaxAmount      float64   `gorm:"type:decimal(12,2);default:0" json:"max_amount"`
	ProductIDs     string    `gorm:"type:text" json:"product_ids"` // JSON array, empty = all
	LevelIDs       string    `gorm:"type:text" json:"level_ids"`   // JSON array, empty = all levels
	Status         string    `gorm:"size:32;default:active" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// GetCommissionSettlements returns paginated commission settlements.
func (s *AgentEnhancedService) GetCommissionSettlements(page, pageSize int) ([]CommissionSettlement, int64, error) {
	var items []CommissionSettlement
	var total int64

	query := s.db.Model(&CommissionSettlement{})
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// SettleCommission settles pending commissions for an agent.
func (s *AgentEnhancedService) SettleCommission(agentID uint, period string) (*CommissionSettlement, error) {
	var agent model.Agent
	if err := s.db.First(&agent, agentID).Error; err != nil {
		return nil, err
	}

	// Calculate total pending commissions
	var totalAmount float64
	query := s.db.Model(&model.AgentCommission{}).
		Where("agent_id = ? AND status = 1", agentID)
	if period != "" {
		query = query.Where("DATE_FORMAT(created_at, '%Y-%m') = ?", period)
	}
	query.Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount)

	if totalAmount <= 0 {
		return nil, fmt.Errorf("no pending commissions to settle")
	}

	settlement := &CommissionSettlement{
		AgentID: agentID,
		Amount:  totalAmount,
		Status:  "paid",
		Period:  period,
		Remark:  fmt.Sprintf("佣金结算 - %s", period),
	}
	now := time.Now()
	settlement.PaidAt = &now

	err := s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(settlement).Error; err != nil {
			return err
		}
		// Mark commissions as confirmed
		updateQuery := tx.Model(&model.AgentCommission{}).
			Where("agent_id = ? AND status = 1", agentID)
		if period != "" {
			updateQuery = updateQuery.Where("DATE_FORMAT(created_at, '%Y-%m') = ?", period)
		}
		return updateQuery.Update("status", 2).Error
	})

	if err != nil {
		return nil, err
	}

	s.log.Infof("commission settled for agent %d: %.2f (period=%s)", agentID, totalAmount, period)
	return settlement, nil
}

// GetCommissionRules returns all commission rules.
func (s *AgentEnhancedService) GetCommissionRules() ([]CommissionRule, error) {
	var items []CommissionRule
	if err := s.db.Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// UpdateCommissionRules updates commission rules (batch replace).
func (s *AgentEnhancedService) UpdateCommissionRules(rules []CommissionRule) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Clear existing rules
		if err := tx.Where("1 = 1").Delete(&CommissionRule{}).Error; err != nil {
			return err
		}
		// Insert new rules
		if len(rules) > 0 {
			return tx.CreateInBatches(rules, 100).Error
		}
		return nil
	})
}
