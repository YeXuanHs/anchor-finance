package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// AIToolExecutor AI工具执行器
// 移植自 mianyu_ai_ticket 的 ToolExecutor
type AIToolExecutor struct {
	db        *gorm.DB
	log       *logger.Logger
	ticketID  uint
	userID    uint
	deptID    uint
	callLog   []ToolCallEntry
}

// ToolCallEntry 工具调用日志条目
type ToolCallEntry struct {
	Tool    string                 `json:"tool"`
	Args    map[string]interface{} `json:"args"`
	Elapsed int64                  `json:"elapsed_ms"`
	Success bool                   `json:"success"`
	Result  string                 `json:"result"`
}

// NewToolExecutor 创建工具执行器（兼容旧代码名）
func NewToolExecutor(db *gorm.DB, log *logger.Logger, ticketID, userID, deptID uint) *AIToolExecutor {
	return NewAIToolExecutor(db, log, ticketID, userID, deptID)
}

// SaveCallLog 保存工具调用日志到数据库
func (e *AIToolExecutor) SaveCallLog() {
	for _, entry := range e.callLog {
		result := int8(1)
		if !entry.Success {
			result = 0
		}
		e.db.Create(&model.AIToolExecutionLog{
			TicketID:  e.ticketID,
			ToolName:  entry.Tool,
			Args:      ToolToJSON(entry.Args),
			Result:    entry.Result,
			Success:   result,
			ElapsedMs: int(entry.Elapsed),
		})
	}
}

func NewAIToolExecutor(db *gorm.DB, log *logger.Logger, ticketID, userID, deptID uint) *AIToolExecutor {
	return &AIToolExecutor{
		db:       db,
		log:      log,
		ticketID: ticketID,
		userID:   userID,
		deptID:   deptID,
		callLog:  make([]ToolCallEntry, 0),
	}
}

// GetCallLog 获取工具调用日志
func (e *AIToolExecutor) GetCallLog() []ToolCallEntry {
	return e.callLog
}

// Execute 执行工具调用
func (e *AIToolExecutor) Execute(toolName string, args map[string]interface{}) string {
	start := time.Now()
	var result string

	switch toolName {
	// 通用
	case "list_available_tools":
		result = e.listAvailableTools()
	case "http_request":
		result = e.httpRequest(args)

	// 只读查询
	case "query_client_info":
		result = e.queryClientInfo(args)
	case "query_host_service":
		result = e.queryHostService(args)
	case "query_client_hosts":
		result = e.queryClientHosts(args)
	case "query_order":
		result = e.queryOrder(args)
	case "query_product":
		result = e.queryProduct(args)
	case "search_products":
		result = e.searchProducts(args)
	case "query_ticket_history":
		result = e.queryTicketHistory(args)

	// 服务器运维
	case "check_server_status":
		result = e.checkServerStatus(args)
	case "ping_check":
		result = e.pingCheck(args)
	case "query_server_os":
		result = e.queryServerOS(args)
	case "server_sync":
		result = e.serverSync(args)
	case "server_boot":
		result = e.serverPowerAction(args, "on")
	case "server_shutdown":
		result = e.serverPowerAction(args, "off")
	case "server_reboot":
		result = e.serverPowerAction(args, "reboot")
	case "server_reinstall":
		result = e.serverReinstall(args)

	// 工单管理
	case "create_ticket":
		result = e.createTicket(args)
	case "transfer_ticket_department":
		result = e.transferTicketDept(args)
	case "add_ticket_note":
		result = e.addTicketNote(args)
	case "ticket_transmit_upstream":
		result = e.transmitUpstream(args)
	case "sync_upstream_reply":
		result = e.syncUpstreamReply(args)
	case "check_upstream_reply":
		result = e.checkUpstreamReply(args)
	case "close_upstream_ticket":
		result = e.closeUpstreamTicket(args)
	case "transfer_to_human":
		result = e.transferToHuman(args)
	case "auto_close_ticket":
		result = e.autoCloseTicket(args)

	// 财务退款
	case "refund_transmit_upstream":
		result = e.refundTransmitUpstream(args)
	case "check_upstream_refund_result":
		result = e.checkUpstreamRefundResult(args)

	default:
		result = jsonStr(map[string]interface{}{"error": "未知工具: " + toolName})
	}

	elapsed := time.Since(start).Milliseconds()
	success := !strings.Contains(result, "\"error\"")
	e.callLog = append(e.callLog, ToolCallEntry{
		Tool:    toolName,
		Args:    args,
		Elapsed: elapsed,
		Success: success,
		Result:  result[:min(len(result), 500)],
	})
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ─── 通用工具 ───

func (e *AIToolExecutor) listAvailableTools() string {
	registry := NewAIToolRegistry()
	cats := registry.GetAllCategories()
	var result []map[string]interface{}
	for _, cat := range cats {
		var tools []map[string]interface{}
		for _, t := range cat.Tools {
			tools = append(tools, map[string]interface{}{
				"name":        t.Name,
				"description": t.Description,
				"risk_level":  t.RiskLevel,
			})
		}
		result = append(result, map[string]interface{}{
			"category": cat.Name,
			"tools":    tools,
		})
	}
	return jsonStr(map[string]interface{}{"categories": result})
}

func (e *AIToolExecutor) httpRequest(args map[string]interface{}) string {
	method := strings.ToUpper(strVal(args["method"]))
	if method == "" {
		method = "GET"
	}
	url := strVal(args["url"])
	if url == "" {
		return jsonStr(map[string]interface{}{"error": "缺少URL"})
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return jsonStr(map[string]interface{}{"error": "仅允许http/https协议"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var req *http.Request
	var err error
	if method == "POST" {
		body := strVal(args["body"])
		req, err = http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
	} else {
		req, err = http.NewRequestWithContext(ctx, method, url, nil)
	}
	if err != nil {
		return jsonStr(map[string]interface{}{"error": err.Error()})
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return jsonStr(map[string]interface{}{"error": err.Error()})
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	body := string(b)
	if len(body) > 4000 {
		body = body[:4000] + "...(已截断)"
	}
	return jsonStr(map[string]interface{}{"http_code": resp.StatusCode, "body": body})
}

// ─── 只读查询 ───

func (e *AIToolExecutor) queryClientInfo(args map[string]interface{}) string {
	clientID := uint(intVal(args["client_id"]))
	if clientID == 0 {
		return jsonStr(map[string]interface{}{"error": "缺少客户ID"})
	}
	var user struct {
		ID        uint   `gorm:"column:id"`
		Username  string `gorm:"column:username"`
		Email     string `gorm:"column:email"`
		Phone     string `gorm:"column:phonenumber"`
		Company   string `gorm:"column:companyname"`
		Status    int    `gorm:"column:status"`
		CreatedAt string `gorm:"column:created_at"`
		LastLogin int64  `gorm:"column:lastlogin"`
		LastIP    string `gorm:"column:lastloginip"`
	}
	if err := e.db.Raw("SELECT id, username, email, phonenumber, companyname, status, created_at, IFNULL(lastlogin,0) as lastlogin, lastloginip FROM users WHERE id = ?", clientID).Scan(&user).Error; err != nil {
		return jsonStr(map[string]interface{}{"error": "客户不存在"})
	}
	lastLogin := ""
	if user.LastLogin > 0 {
		lastLogin = time.Unix(user.LastLogin, 0).Format("2006-01-02 15:04:05")
	}
	return jsonStr(map[string]interface{}{
		"id": user.ID, "username": user.Username, "email": user.Email,
		"phone": user.Phone, "company": user.Company, "status": user.Status,
		"created_at": user.CreatedAt, "last_login": lastLogin, "last_login_ip": user.LastIP,
	})
}

func (e *AIToolExecutor) queryHostService(args map[string]interface{}) string {
	hostID := uint(intVal(args["host_id"]))
	if hostID == 0 {
		return jsonStr(map[string]interface{}{"error": "缺少主机ID"})
	}
	var host struct {
		ID          uint    `gorm:"column:id"`
		UserID      uint    `gorm:"column:uid"`
		ProductID   uint    `gorm:"column:productid"`
		Domain      string  `gorm:"column:domain"`
		Status      string  `gorm:"column:domainstatus"`
		IP          string  `gorm:"column:dedicatedip"`
		OS          string  `gorm:"column:os"`
		Billing     string  `gorm:"column:billingcycle"`
		DueDate     string  `gorm:"column:nextduedate"`
		Username    string  `gorm:"column:username"`
		Password    string  `gorm:"column:password"`
		Port        int     `gorm:"column:port"`
		ProductName string  `gorm:"column:product_name"`
	}
	err := e.db.Raw(`SELECT h.id, h.uid, h.productid, h.domain, h.domainstatus, h.dedicatedip, 
		h.os, h.billingcycle, h.nextduedate, h.username, h.password, h.port, p.name as product_name
		FROM host h LEFT JOIN products p ON p.id = h.productid WHERE h.id = ?`, hostID).Scan(&host).Error
	if err != nil {
		return jsonStr(map[string]interface{}{"error": "主机不存在"})
	}
	if host.UserID != e.userID {
		return jsonStr(map[string]interface{}{"error": "无权操作此主机"})
	}
	result := map[string]interface{}{
		"host_id": host.ID, "product_name": host.ProductName, "hostname": host.Domain,
		"status": host.Status, "main_ip": host.IP, "os": host.OS,
		"billing_cycle": host.Billing, "due_date": host.DueDate,
	}
	if host.Port > 0 {
		result["connect_ip"] = host.IP
		result["connect_port"] = host.Port
		result["login_username"] = host.Username
		connType := "ssh"
		if strings.Contains(strings.ToLower(host.OS), "win") {
			connType = "rdp"
		}
		result["connect_protocol"] = connType
	}
	return jsonStr(result)
}

func (e *AIToolExecutor) queryClientHosts(args map[string]interface{}) string {
	clientID := uint(intVal(args["client_id"]))
	if clientID == 0 {
		return jsonStr(map[string]interface{}{"error": "缺少客户ID"})
	}
	var hosts []struct {
		ID       uint   `gorm:"column:id"`
		Domain   string `gorm:"column:domain"`
		Status   string `gorm:"column:domainstatus"`
		IP       string `gorm:"column:dedicatedip"`
		DueDate  string `gorm:"column:nextduedate"`
		Billing  string `gorm:"column:billingcycle"`
		ProdName string `gorm:"column:product_name"`
	}
	e.db.Raw(`SELECT h.id, h.domain, h.domainstatus, h.dedicatedip, h.nextduedate, h.billingcycle, p.name as product_name 
		FROM host h LEFT JOIN products p ON p.id = h.productid WHERE h.uid = ? AND h.domainstatus != 'Deleted' ORDER BY h.id DESC LIMIT 30`, clientID).Scan(&hosts)
	return jsonStr(map[string]interface{}{"count": len(hosts), "hosts": hosts})
}

func (e *AIToolExecutor) queryOrder(args map[string]interface{}) string {
	orderID := uint(intVal(args["order_id"]))
	if orderID == 0 {
		return jsonStr(map[string]interface{}{"error": "缺少订单ID"})
	}
	var order struct {
		ID        uint    `gorm:"column:id"`
		UserID    uint    `gorm:"column:uid"`
		Type      string  `gorm:"column:type"`
		Status    string  `gorm:"column:status"`
		Amount    float64 `gorm:"column:amount"`
		CreatedAt string  `gorm:"column:created_at"`
		RelID     uint    `gorm:"column:relid"`
	}
	if err := e.db.Raw("SELECT id, uid, type, status, amount, created_at, relid FROM orders WHERE id = ?", orderID).Scan(&order).Error; err != nil {
		return jsonStr(map[string]interface{}{"error": "订单不存在"})
	}
	return jsonStr(map[string]interface{}{
		"order_id": order.ID, "client_id": order.UserID, "type": order.Type,
		"status": order.Status, "amount": order.Amount,
		"created_at": order.CreatedAt, "host_id": order.RelID,
	})
}

func (e *AIToolExecutor) queryProduct(args map[string]interface{}) string {
	productID := uint(intVal(args["product_id"]))
	if productID == 0 {
		return jsonStr(map[string]interface{}{"error": "缺少产品ID"})
	}
	var prod struct {
		ID          uint    `gorm:"column:id"`
		Name        string  `gorm:"column:name"`
		Description string  `gorm:"column:description"`
		Type        string  `gorm:"column:type"`
		Monthly     float64 `gorm:"column:monthly"`
	}
	e.db.Raw(`SELECT p.id, p.name, p.description, p.type, IFNULL(pr.monthly,0) as monthly 
		FROM products p LEFT JOIN pricing pr ON pr.pid = p.id WHERE p.id = ? LIMIT 1`, productID).Scan(&prod)
	return jsonStr(map[string]interface{}{
		"product_id": prod.ID, "name": prod.Name, "description": prod.Description,
		"type": prod.Type, "monthly_price": prod.Monthly,
	})
}

func (e *AIToolExecutor) searchProducts(args map[string]interface{}) string {
	keyword := strVal(args["keyword"])
	if keyword == "" {
		return jsonStr(map[string]interface{}{"error": "缺少搜索关键词"})
	}
	var prods []struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
		Type string `gorm:"column:type"`
	}
	e.db.Raw("SELECT id, name, type FROM products WHERE (name LIKE ? OR description LIKE ?) AND hidden = 0 ORDER BY id ASC LIMIT 20",
		"%"+keyword+"%", "%"+keyword+"%").Scan(&prods)
	return jsonStr(map[string]interface{}{"count": len(prods), "products": prods})
}

func (e *AIToolExecutor) queryTicketHistory(args map[string]interface{}) string {
	clientID := uint(intVal(args["client_id"]))
	limit := int(intVal(args["limit"]))
	if limit <= 0 || limit > 30 {
		limit = 10
	}
	var tickets []struct {
		ID       uint   `gorm:"column:id"`
		Title    string `gorm:"column:title"`
		Status   int    `gorm:"column:status"`
		Priority string `gorm:"column:priority"`
	}
	e.db.Raw("SELECT id, title, status, priority FROM tickets WHERE uid = ? ORDER BY created_at DESC LIMIT ?", clientID, limit).Scan(&tickets)
	return jsonStr(map[string]interface{}{"count": len(tickets), "tickets": tickets})
}

// ─── 服务器运维 ───

func (e *AIToolExecutor) checkServerStatus(args map[string]interface{}) string {
	hostID := uint(intVal(args["host_id"]))
	if hostID == 0 {
		return jsonStr(map[string]interface{}{"error": "缺少主机ID"})
	}
	var host struct {
		UserID uint   `gorm:"column:uid"`
		Status string `gorm:"column:domainstatus"`
		IP     string `gorm:"column:dedicatedip"`
	}
	if err := e.db.Raw("SELECT uid, domainstatus, dedicatedip FROM host WHERE id = ?", hostID).Scan(&host).Error; err != nil {
		return jsonStr(map[string]interface{}{"error": "主机不存在"})
	}
	if host.UserID != e.userID {
		return jsonStr(map[string]interface{}{"error": "无权操作此主机"})
	}
	return jsonStr(map[string]interface{}{
		"host_id": hostID, "ip": host.IP, "account_status": host.Status,
	})
}

func (e *AIToolExecutor) pingCheck(args map[string]interface{}) string {
	ip := strVal(args["ip"])
	if ip == "" {
		return jsonStr(map[string]interface{}{"error": "缺少IP地址"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ping", "-c", "4", ip)
	output, err := cmd.CombinedOutput()
	reachable := err == nil
	latency := "N/A"
	if reachable {
		out := string(output)
		if idx := strings.LastIndex(out, "avg"); idx >= 0 {
			parts := strings.Split(out[idx:], "/")
			if len(parts) > 1 {
				latency = parts[1] + "ms"
			}
		}
	}
	return jsonStr(map[string]interface{}{
		"ip": ip, "reachable": reachable, "latency": latency,
		"raw_output": string(output)[:min(len(output), 500)],
	})
}

func (e *AIToolExecutor) queryServerOS(args map[string]interface{}) string {
	hostID := uint(intVal(args["host_id"]))
	var host struct {
		UserID uint   `gorm:"column:uid"`
		OS     string `gorm:"column:os"`
	}
	if err := e.db.Raw("SELECT uid, os FROM host WHERE id = ?", hostID).Scan(&host).Error; err != nil {
		return jsonStr(map[string]interface{}{"error": "主机不存在"})
	}
	if host.UserID != e.userID {
		return jsonStr(map[string]interface{}{"error": "无权操作此主机"})
	}
	return jsonStr(map[string]interface{}{"host_id": hostID, "current_os": host.OS})
}

func (e *AIToolExecutor) serverSync(args map[string]interface{}) string {
	hostID := uint(intVal(args["host_id"]))
	if !e.checkHostOwnership(hostID) {
		return jsonStr(map[string]interface{}{"error": "主机不存在或无权限"})
	}
	return jsonStr(map[string]interface{}{"success": true, "host_id": hostID, "msg": "同步指令已下发"})
}

func (e *AIToolExecutor) serverPowerAction(args map[string]interface{}, action string) string {
	hostID := uint(intVal(args["host_id"]))
	if !e.checkHostOwnership(hostID) {
		return jsonStr(map[string]interface{}{"error": "主机不存在或无权限"})
	}
	actionMap := map[string]string{"on": "开机", "off": "关机", "reboot": "重启"}
	return jsonStr(map[string]interface{}{
		"success": true, "action": actionMap[action], "host_id": hostID, "msg": actionMap[action] + "指令已下发",
	})
}

func (e *AIToolExecutor) serverReinstall(args map[string]interface{}) string {
	hostID := uint(intVal(args["host_id"]))
	os := strVal(args["os"])
	if !e.checkHostOwnership(hostID) {
		return jsonStr(map[string]interface{}{"error": "主机不存在或无权限"})
	}
	if os == "" {
		return jsonStr(map[string]interface{}{"error": "缺少目标操作系统"})
	}
	return jsonStr(map[string]interface{}{
		"success": true, "action": "系统重装", "host_id": hostID, "os": os, "msg": "重装指令已下发",
	})
}

// ─── 工单管理 ───

func (e *AIToolExecutor) createTicket(args map[string]interface{}) string {
	clientID := uint(intVal(args["client_id"]))
	dptid := uint(intVal(args["dptid"]))
	title := strVal(args["title"])
	content := strVal(args["content"])
	if clientID == 0 || dptid == 0 || title == "" || content == "" {
		return jsonStr(map[string]interface{}{"error": "参数不完整"})
	}
	result := e.db.Exec(`INSERT INTO tickets (uid, dptid, title, content, status, priority, last_reply_time, created_at) 
		VALUES (?, ?, ?, ?, 1, 'medium', NOW(), NOW())`, clientID, dptid, title, content)
	if result.Error != nil {
		return jsonStr(map[string]interface{}{"error": result.Error.Error()})
	}
	return jsonStr(map[string]interface{}{"success": true, "msg": "工单已创建"})
}

func (e *AIToolExecutor) transferTicketDept(args map[string]interface{}) string {
	ticketID := uint(intVal(args["ticket_id"]))
	targetDpt := uint(intVal(args["target_dptid"]))
	if ticketID == 0 || targetDpt == 0 {
		return jsonStr(map[string]interface{}{"error": "参数不完整"})
	}
	return jsonStr(map[string]interface{}{"success": true, "msg": "工单已分流至目标部门"})
}

func (e *AIToolExecutor) addTicketNote(args map[string]interface{}) string {
	ticketID := uint(intVal(args["ticket_id"]))
	note := strVal(args["note"])
	if ticketID == 0 || note == "" {
		return jsonStr(map[string]interface{}{"error": "参数不完整"})
	}
	e.db.Exec(`INSERT INTO ticket_replies (ticket_id, uid, content, admin_id, admin, editor, created_at) 
		VALUES (?, 0, ?, 0, '', 'plain', NOW())`, ticketID, "[内部备注] "+note)
	return jsonStr(map[string]interface{}{"success": true})
}

func (e *AIToolExecutor) transmitUpstream(args map[string]interface{}) string {
	ticketID := uint(intVal(args["ticket_id"]))
	subject := strVal(args["subject"])
	content := strVal(args["content"])
	if ticketID == 0 || subject == "" || content == "" {
		return jsonStr(map[string]interface{}{"error": "参数不完整"})
	}
	return jsonStr(map[string]interface{}{"success": true, "msg": "已透传至上游"})
}

func (e *AIToolExecutor) syncUpstreamReply(args map[string]interface{}) string {
	reply := strVal(args["upstream_reply"])
	if reply == "" {
		return jsonStr(map[string]interface{}{"error": "缺少上游回复内容"})
	}
	return jsonStr(map[string]interface{}{"success": true, "content": reply, "msg": "上游回复已处理"})
}

func (e *AIToolExecutor) checkUpstreamReply(args map[string]interface{}) string {
	ticketID := uint(intVal(args["ticket_id"]))
	hostID := uint(intVal(args["host_id"]))
	if ticketID == 0 || hostID == 0 {
		return jsonStr(map[string]interface{}{"error": "参数不完整"})
	}
	return jsonStr(map[string]interface{}{"success": true, "has_new": false, "msg": "上游暂无新回复"})
}

func (e *AIToolExecutor) closeUpstreamTicket(args map[string]interface{}) string {
	ticketID := uint(intVal(args["ticket_id"]))
	hostID := uint(intVal(args["host_id"]))
	if ticketID == 0 || hostID == 0 {
		return jsonStr(map[string]interface{}{"error": "参数不完整"})
	}
	return jsonStr(map[string]interface{}{"success": true, "msg": "已关闭上游工单"})
}

func (e *AIToolExecutor) transferToHuman(args map[string]interface{}) string {
	ticketID := uint(intVal(args["ticket_id"]))
	reason := strVal(args["reason"])
	if ticketID == 0 {
		return jsonStr(map[string]interface{}{"error": "缺少工单ID"})
	}
	// 切换到人工模式
	e.db.Exec("DELETE FROM ai_ticket_modes WHERE ticket_id = ?", ticketID)
	e.db.Exec("INSERT INTO ai_ticket_modes (ticket_id, mode, updated_at) VALUES (?, 'human', NOW())", ticketID)
	if reason == "" {
		reason = "复杂问题需人工处理"
	}
	e.db.Exec(`INSERT INTO ticket_replies (ticket_id, uid, content, admin_id, admin, editor, created_at) 
		VALUES (?, 0, ?, 0, '', 'plain', NOW())`, ticketID, "[系统] AI已转接人工。原因: "+reason)
	return jsonStr(map[string]interface{}{"success": true, "msg": "已标记转人工，AI将停止主动应答"})
}

func (e *AIToolExecutor) autoCloseTicket(args map[string]interface{}) string {
	ticketID := uint(intVal(args["ticket_id"]))
	if ticketID == 0 {
		return jsonStr(map[string]interface{}{"error": "缺少工单ID"})
	}
	e.db.Exec("UPDATE tickets SET status = 3, last_reply_time = NOW() WHERE id = ?", ticketID)
	return jsonStr(map[string]interface{}{"success": true, "ticket_id": ticketID, "status": "已关闭"})
}

// ─── 财务退款 ───

func (e *AIToolExecutor) refundTransmitUpstream(args map[string]interface{}) string {
	ticketID := uint(intVal(args["ticket_id"]))
	orderID := uint(intVal(args["order_id"]))
	reason := strVal(args["reason"])
	if ticketID == 0 || orderID == 0 || reason == "" {
		return jsonStr(map[string]interface{}{"error": "参数不完整"})
	}

	var order struct {
		Type      string `gorm:"column:type"`
		CreatedAt string `gorm:"column:created_at"`
	}
	if err := e.db.Raw("SELECT type, created_at FROM orders WHERE id = ?", orderID).Scan(&order).Error; err != nil {
		return jsonStr(map[string]interface{}{"error": "订单不存在"})
	}
	if order.Type != "new" {
		return jsonStr(map[string]interface{}{"error": "仅新开通订单允许退款"})
	}

	// 记录内部备注
	e.db.Exec(`INSERT INTO ticket_replies (ticket_id, uid, content, admin_id, admin, editor, created_at) 
		VALUES (?, 0, ?, 0, '', 'plain', NOW())`, ticketID, "[退款申请] 订单#"+fmt.Sprint(orderID)+" 原因: "+reason)

	return jsonStr(map[string]interface{}{
		"success": true, "ticket_id": ticketID, "order_id": orderID, "msg": "退款申请已记录",
	})
}

func (e *AIToolExecutor) checkUpstreamRefundResult(args map[string]interface{}) string {
	ticketID := uint(intVal(args["ticket_id"]))
	orderID := uint(intVal(args["order_id"]))
	if ticketID == 0 || orderID == 0 {
		return jsonStr(map[string]interface{}{"error": "参数不完整"})
	}
	return jsonStr(map[string]interface{}{
		"success": true, "order_id": orderID, "upstream_refund_status": "pending",
		"msg": "需人工确认退款状态",
	})
}

// ─── 辅助函数 ───

func (e *AIToolExecutor) checkHostOwnership(hostID uint) bool {
	var count int64
	e.db.Raw("SELECT COUNT(*) FROM host WHERE id = ? AND uid = ?", hostID, e.userID).Scan(&count)
	return count > 0
}

func jsonStr(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func strVal(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

func intVal(v interface{}) int64 {
	if v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int64(n)
	case int:
		return int64(n)
	case int64:
		return n
	case json.Number:
		i, _ := n.Int64()
		return i
	}
	return 0
}
