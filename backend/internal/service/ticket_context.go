package service

import (
	"encoding/json"
	"fmt"
	"html"
	"regexp"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// TicketContextService 工单上下文服务
// 为 AI 构建完整的上下文 JSON，包含工单、客户、产品、对话历史等信息
type TicketContextService struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewTicketContextService 创建工单上下文服务
func NewTicketContextService(db *gorm.DB, log *logger.Logger) *TicketContextService {
	return &TicketContextService{db: db, log: log}
}

// BuildContextJson 为 AI 构建完整的上下文 JSON
func (s *TicketContextService) BuildContextJson(ticketID, uid uint) string {
	var ticket model.Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		return "{}"
	}

	var dept model.Department
	s.db.Select("name").First(&dept, ticket.DepartmentID)

	relatedHost := s.resolveRelatedHost(ticket, uid)

	context := map[string]interface{}{
		"ticket": map[string]interface{}{
			"id":              ticket.ID,
			"number":          ticket.TicketNo,
			"title":           s.plainText(ticket.Subject),
			"department":      dept.Name,
			"priority":        s.priorityLabel(ticket.Priority),
			"status":          s.ticketStatusLabel(ticket.Status),
			"created_at":      ticket.CreatedAt.Format("2006-01-02 15:04:05"),
			"related_product": relatedHost,
		},
		"user":         s.buildUserInfo(uid),
		"products":     s.buildUserProducts(uid),
		"conversation": s.buildConversation(ticketID, ticket),
	}

	b, _ := json.Marshal(context)
	return string(b)
}

// BuildMessages 构建对话历史消息列表（用于 AI 调用）
func (s *TicketContextService) BuildMessages(ticketID uint) []AIChatMessage {
	var ticket model.Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		return nil
	}
	return s.buildConversation(ticketID, ticket)
}

// resolveRelatedHost 解析工单关联的主机
// 优先用工单的 RelID 字段，没有则从工单内容中提取 IP 匹配
func (s *TicketContextService) resolveRelatedHost(ticket model.Ticket, uid uint) map[string]interface{} {
	if ticket.RelID > 0 && ticket.RelType == "product" {
		var host model.Host
		if err := s.db.First(&host, ticket.RelID).Error; err == nil {
			return s.formatHost(host, uid, true)
		}
	}
	return s.guessRelatedHost(ticket, uid)
}

// guessRelatedHost 从客户活跃主机中猜测关联主机
func (s *TicketContextService) guessRelatedHost(ticket model.Ticket, uid uint) map[string]interface{} {
	var hosts []model.Host
	s.db.Where("owner_id = ? AND status = ?", uid, 1).
		Order("id DESC").
		Limit(5).
		Find(&hosts)

	if len(hosts) == 0 {
		return nil
	}

	haystack := s.plainText(ticket.Subject)
	var content string
	s.db.Model(&model.TicketReply{}).Where("ticket_id = ?", ticket.ID).Pluck("content", &content)
	haystack += " " + s.plainText(content)

	for _, host := range hosts {
		if host.IP != "" && strings.Contains(haystack, host.IP) {
			return s.formatHost(host, uid, true)
		}
	}

	if len(hosts) == 1 {
		return s.formatHost(hosts[0], uid, true)
	}

	return nil
}

// buildUserInfo 构建客户信息
func (s *TicketContextService) buildUserInfo(uid uint) map[string]interface{} {
	var user model.User
	if err := s.db.First(&user, uid).Error; err != nil {
		return map[string]interface{}{"id": uid}
	}

	phone := strings.TrimSpace(user.CountryCode + " " + user.Phone)
	lastLogin := ""
	if user.LastLoginAt != nil {
		lastLogin = user.LastLoginAt.Format("2006-01-02 15:04:05")
	}

	return map[string]interface{}{
		"id":            user.ID,
		"username":      user.Username,
		"email":         user.Email,
		"phone":         phone,
		"company":       "",
		"credit":        user.Balance.String(),
		"status":        s.clientStatusLabel(user.Status),
		"registered":    user.CreatedAt.Format("2006-01-02 15:04:05"),
		"last_login":    lastLogin,
		"last_login_ip": user.LastLoginIP,
	}
}

// buildUserProducts 构建客户产品列表（最多20个，仅基础字段）
func (s *TicketContextService) buildUserProducts(uid uint) []map[string]interface{} {
	var hosts []model.Host
	s.db.Where("owner_id = ? AND status != ?", uid, 5). // 5=Deleted
		Order("id DESC").
		Limit(20).
		Find(&hosts)

	var list []map[string]interface{}
	for _, host := range hosts {
		list = append(list, s.formatHost(host, uid, false))
	}
	return list
}

// formatHost 格式化主机信息
func (s *TicketContextService) formatHost(host model.Host, uid uint, withConnect bool) map[string]interface{} {
	var productName string
	if host.ProductID != nil {
		var prod model.Product
		if err := s.db.Select("name").First(&prod, *host.ProductID).Error; err == nil {
			productName = prod.Name
		}
	}

	dueDate := ""
	if host.ExpiredAt != nil {
		dueDate = host.ExpiredAt.Format("2006-01-02")
	}

	item := map[string]interface{}{
		"host_id":       host.ID,
		"product_name":  productName,
		"hostname":      host.Hostname,
		"status":        s.hostStatusLabel(host.Status),
		"main_ip":       host.IP,
		"billing_cycle": "",
		"due_date":      dueDate,
	}

	if !withConnect {
		return item
	}

	connect := s.resolveHostConnect(host)
	if connect != nil {
		item["connect_ip"] = connect["host"]
		item["connect_port"] = connect["port"]
		item["login_username"] = connect["username"]
		item["login_password"] = connect["password"]
		item["connect_protocol"] = connect["connection_type"]
	}

	specs := s.resolveHostSpecs(host)
	if len(specs) > 0 {
		item["specs"] = specs
	}

	return item
}

// resolveHostConnect 获取主机连接信息（轻量版，只读表字段）
func (s *TicketContextService) resolveHostConnect(host model.Host) map[string]interface{} {
	if host.IP == "" {
		return nil
	}

	connType := "ssh"
	if strings.Contains(strings.ToLower(host.OS), "win") {
		connType = "rdp"
	}

	port := 22
	if connType == "rdp" {
		port = 3389
	}

	username := "root"
	if connType == "rdp" {
		username = "administrator"
	}

	return map[string]interface{}{
		"host_id":         host.ID,
		"host":            host.IP,
		"port":            fmt.Sprintf("%d", port),
		"username":        username,
		"password":        "",
		"connection_type": connType,
		"product_name":    host.Hostname,
	}
}

// resolveHostSpecs 获取主机配置参数
func (s *TicketContextService) resolveHostSpecs(host model.Host) map[string]interface{} {
	specs := make(map[string]interface{})

	if host.OS != "" {
		specs["操作系统"] = host.OS
	}
	if host.CPU != "" {
		specs["CPU"] = host.CPU
	}
	if host.CPUCores > 0 {
		specs["CPU核心"] = host.CPUCores
	}
	if host.MemoryMB > 0 {
		specs["内存"] = fmt.Sprintf("%dMB", host.MemoryMB)
	}
	if host.DiskSizeGB > 0 {
		specs["磁盘"] = fmt.Sprintf("%dGB %s", host.DiskSizeGB, host.DiskType)
	}
	if host.BandwidthMbps > 0 {
		specs["带宽"] = fmt.Sprintf("%dMbps", host.BandwidthMbps)
	}
	if host.Location != "" {
		specs["机房"] = host.Location
	}

	// 解析 Config JSON 字段
	if len(host.Config) > 0 {
		var configMap map[string]interface{}
		if json.Unmarshal(host.Config, &configMap) == nil {
			for k, v := range configMap {
				label := normalizeConfigKey(k)
				if label != "" && v != nil && v != "" {
					if _, exists := specs[label]; !exists {
						specs[label] = stringifyValue(v)
					}
				}
			}
		}
	}

	return specs
}

// buildConversation 构建对话历史
func (s *TicketContextService) buildConversation(ticketID uint, ticket model.Ticket) []AIChatMessage {
	var messages []AIChatMessage

	// 工单初始内容作为用户消息
	initial := s.plainText(ticket.Subject)
	if initial != "" {
		// 查询工单内容（第一条回复或工单描述）
		var firstReply model.TicketReply
		if err := s.db.Where("ticket_id = ?", ticketID).Order("id ASC").First(&firstReply).Error; err == nil {
			content := s.plainText(firstReply.Content)
			if content != "" {
				messages = append(messages, AIChatMessage{Role: "user", Content: content})
			}
		}
	}

	// 获取 AI 管理员 ID
	aiAdminID := uint(0)
	var cfg model.AITicketConfig
	if err := s.db.Where("`key` = ?", "reply_admin_id").First(&cfg).Error; err == nil {
		fmt.Sscanf(cfg.Value, "%d", &aiAdminID)
	}
	if aiAdminID == 0 {
		aiAdminID = 1
	}

	// 查询所有回复
	var replies []model.TicketReply
	s.db.Where("ticket_id = ?", ticketID).
		Order("created_at ASC, id ASC").
		Find(&replies)

	for _, reply := range replies {
		text := s.plainText(reply.Content)
		if text == "" {
			continue
		}

		// 附件标记
		var attachments []model.Attachment
		s.db.Where("reply_id = ?", reply.ID).Find(&attachments)
		for _, att := range attachments {
			text += fmt.Sprintf("\n[附件: %s]", att.FileName)
		}

		if reply.AdminID != nil && *reply.AdminID > 0 {
			// 管理员回复
			if *reply.AdminID == aiAdminID {
				messages = append(messages, AIChatMessage{Role: "assistant", Content: text})
			} else {
				messages = append(messages, AIChatMessage{Role: "assistant", Content: "[人工客服] " + text})
			}
		} else if reply.UserID != nil && *reply.UserID > 0 {
			// 客户回复
			messages = append(messages, AIChatMessage{Role: "user", Content: text})
		} else {
			// 系统消息
			messages = append(messages, AIChatMessage{Role: "system", Content: text})
		}
	}

	return messages
}

// ─── 辅助方法 ───

var tagRegex = regexp.MustCompile(`<[^>]+>`)

func (s *TicketContextService) plainText(htmlStr string) string {
	text := html.UnescapeString(htmlStr)
	text = tagRegex.ReplaceAllString(text, " ")
	text = regexp.MustCompile(`\s+`).ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

func (s *TicketContextService) priorityLabel(priority int16) string {
	switch priority {
	case 0:
		return "低"
	case 1:
		return "中"
	case 2:
		return "高"
	case 3:
		return "紧急"
	default:
		return fmt.Sprintf("%d", priority)
	}
}

func (s *TicketContextService) ticketStatusLabel(status int16) string {
	switch status {
	case 0:
		return "待回复"
	case 1:
		return "已回复"
	case 2:
		return "已关闭"
	case 3:
		return "待处理"
	case 4:
		return "已解决"
	case 5:
		return "已取消"
	default:
		return fmt.Sprintf("%d", status)
	}
}

func (s *TicketContextService) clientStatusLabel(status int16) string {
	switch status {
	case 0:
		return "未激活"
	case 1:
		return "正常"
	case 2:
		return "已关闭"
	default:
		return fmt.Sprintf("%d", status)
	}
}

func (s *TicketContextService) hostStatusLabel(status int16) string {
	switch status {
	case 0:
		return "关机"
	case 1:
		return "运行中"
	case 2:
		return "故障"
	case 3:
		return "维护"
	case 4:
		return "创建中"
	default:
		return fmt.Sprintf("%d", status)
	}
}

// normalizeConfigKey 标准化配置键名
func normalizeConfigKey(key string) string {
	key = strings.TrimSpace(key)
	if key == "" {
		return ""
	}
	lower := strings.ToLower(key)
	switch lower {
	case "cpu":
		return "CPU"
	case "memory", "mem":
		return "内存"
	case "bw", "bandwidth":
		return "带宽"
	case "data_disk", "disk":
		return "数据盘"
	case "system_disk":
		return "系统盘"
	case "ip_num":
		return "IP数量"
	case "area":
		return "区域"
	case "os":
		return "操作系统"
	default:
		return key
	}
}

// stringifyValue 将值转为字符串
func stringifyValue(v interface{}) string {
	switch val := v.(type) {
	case bool:
		if val {
			return "是"
		}
		return "否"
	case []interface{}:
		b, _ := json.Marshal(val)
		return string(b)
	case map[string]interface{}:
		b, _ := json.Marshal(val)
		return string(b)
	case float64:
		if val == float64(int64(val)) {
			return fmt.Sprintf("%d", int64(val))
		}
		return fmt.Sprintf("%.2f", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// BuildTicketInfoForNotify 构建工单通知信息
func (s *TicketContextService) BuildTicketInfoForNotify(ticketID uint) map[string]interface{} {
	var ticket model.Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		return nil
	}

	var deptName string
	s.db.Model(&model.Department{}).Where("id = ?", ticket.DepartmentID).Pluck("name", &deptName)

	var content string
	var firstReply model.TicketReply
	if err := s.db.Where("ticket_id = ?", ticketID).Order("id ASC").First(&firstReply).Error; err == nil {
		content = s.plainText(firstReply.Content)
	}

	return map[string]interface{}{
		"ticket_id":  ticket.ID,
		"tid":        ticket.TicketNo,
		"title":      ticket.Subject,
		"department": deptName,
		"priority":   s.priorityLabel(ticket.Priority),
		"content":    content,
		"time":       time.Now().Format("2006-01-02 15:04:05"),
	}
}
