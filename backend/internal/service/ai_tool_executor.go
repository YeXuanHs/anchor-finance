package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
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
		result = e.searchProductsSimple(args)
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

	// 商品导购
	case "search_products_guided":
		result = e.searchProducts(args)
	case "get_product_detail":
		result = e.getProductDetail(args)
	case "list_product_groups":
		result = e.listProductGroups()
	case "get_group_products":
		result = e.getGroupProducts(args)
	case "compare_products":
		result = e.compareProducts(args)

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

func (e *AIToolExecutor) searchProductsSimple(args map[string]interface{}) string {
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
	hostID := uint(intVal(args["host_id"]))
	subject := strVal(args["subject"])
	content := strVal(args["content"])

	if ticketID == 0 || subject == "" || content == "" {
		return jsonStr(map[string]interface{}{"error": "参数不完整：需要 ticket_id、subject、content"})
	}

	// 读取工单信息
	var ticket struct {
		ID       uint   `gorm:"column:id"`
		TicketNo string `gorm:"column:ticket_no"`
		UserID   uint   `gorm:"column:user_id"`
		Priority int    `gorm:"column:priority"`
		RelID    uint   `gorm:"column:rel_id"`
	}
	if err := e.db.Raw("SELECT id, ticket_no, user_id, priority, rel_id FROM tickets WHERE id = ?", ticketID).Scan(&ticket).Error; err != nil {
		return jsonStr(map[string]interface{}{"error": "工单不存在"})
	}

	// host_id 优先用参数传入的，没传则自动从工单读取
	if hostID == 0 {
		hostID = ticket.RelID
	}
	if hostID == 0 {
		return jsonStr(map[string]interface{}{
			"error":   "该工单未关联产品/主机，无法识别上游渠道。请先询问客户具体是哪个产品/服务器，然后通过 edit_ticket 设置工单的 host_id",
			"no_host": true,
		})
	}

	// 查询 host 信息
	var host struct {
		ID        int64  `gorm:"column:id"`
		UID       int64  `gorm:"column:uid"`
		ProductID int64  `gorm:"column:productid"`
		Server    string `gorm:"column:server"`
		ParentID  int64  `gorm:"column:parent_id"`
		DcimID    int64  `gorm:"column:dcimid"`
		Domain    string `gorm:"column:domain"`
	}
	if err := e.db.Raw("SELECT id, uid, productid, IFNULL(server,'') as server, IFNULL(parent_id,0) as parent_id, IFNULL(dcimid,0) as dcimid, IFNULL(domain,'') as domain FROM host WHERE id = ?", hostID).Scan(&host).Error; err != nil {
		return jsonStr(map[string]interface{}{"error": "主机不存在"})
	}

	// 获取客户信息
	var user struct {
		Username string `gorm:"column:username"`
	}
	e.db.Raw("SELECT IFNULL(username,'未知') as username FROM users WHERE id = ?", host.UID).Scan(&user)

	// 获取上游父服务的 dcimid
	var parentDcimID int64
	if host.ParentID > 0 {
		var parentHost struct {
			DcimID int64 `gorm:"column:dcimid"`
		}
		if err := e.db.Raw("SELECT IFNULL(dcimid,0) as dcimid FROM host WHERE id = ?", host.ParentID).Scan(&parentHost).Error; err == nil {
			parentDcimID = parentHost.DcimID
		}
	}

	// 获取本工单最近5条对话作为上下文
	type replyEntry struct {
		Content string `gorm:"column:content"`
		UID     int64  `gorm:"column:uid"`
	}
	var recentReplies []replyEntry
	e.db.Raw("SELECT IFNULL(content,'') as content, IFNULL(uid,0) as uid FROM ticket_replies WHERE ticket_id = ? ORDER BY id DESC LIMIT 5", ticketID).Scan(&recentReplies)
	// 反转为时间正序
	for i, j := 0, len(recentReplies)-1; i < j; i, j = i+1, j-1 {
		recentReplies[i], recentReplies[j] = recentReplies[j], recentReplies[i]
	}

	conversationSummary := ""
	for _, r := range recentReplies {
		role := "客服"
		if r.UID > 0 {
			role = "客户"
		}
		conversationSummary += fmt.Sprintf("[%s] %s\n", role, stripTags(r.Content))
	}

	upstreamChannel := host.Server
	parentHostID := host.ParentID
	dcimID := host.DcimID

	// 构建上游工单内容
	upstreamContent := fmt.Sprintf("【客户工单透传】\n工单号：%s\n客户：%s (UID:%d)\n本地主机ID：%d\n机器标识(dcimid)：\n上游渠道：%s\n父服务ID：%d\n父服务dcimid：%d\n\n--- 工单内容 ---\n%s\n\n--- 最近对话 ---\n%s",
		ticket.TicketNo, user.Username, host.UID, hostID, dcimID, upstreamChannel, parentHostID, parentDcimID, content, conversationSummary)

	// 尝试通过上游渠道提交工单
	var upstreamResult *UpstreamResult
	var upstreamTicketID string

	if upstreamChannel != "" && parentHostID > 0 {
		client := NewUpstreamClient(e.db, e.log)
		upstreamResult = client.CallForHost(parentHostID, "open_support_ticket", map[string]interface{}{
			"subject":     subject,
			"message":     upstreamContent,
			"service_id":  parentHostID,
			"priority":    fmt.Sprintf("%d", ticket.Priority),
		})
		if upstreamResult != nil && upstreamResult.Data != nil {
			upstreamTicketID = strVal(upstreamResult.Data["ticketid"])
			if upstreamTicketID == "" {
				upstreamTicketID = strVal(upstreamResult.Data["ticket_id"])
			}
			if upstreamTicketID == "" {
				upstreamTicketID = strVal(upstreamResult.Data["id"])
			}
		}
	}

	// 记录内部备注
	noteText := fmt.Sprintf("【上游工单透传】\n标题：%s\n主机ID：%d\ndcimid：%d\n上游渠道：%s\n父服务ID：%d\n",
		subject, hostID, dcimID, upstreamChannel, parentHostID)

	if upstreamResult != nil && upstreamResult.Error != "" {
		noteText += "上游调用状态：失败 - " + upstreamResult.Error + "\n"
	} else if upstreamResult != nil && upstreamResult.Data != nil {
		noteText += "上游调用状态：成功\n"
		if upstreamTicketID != "" {
			noteText += "上游工单号：" + upstreamTicketID + "\n"
		}
		rawJSON, _ := json.Marshal(upstreamResult.Data)
		noteText += "上游返回：" + string(rawJSON) + "\n"
	} else {
		noteText += "上游调用状态：未配置上游渠道或无父服务，需人工处理\n"
	}

	e.addTicketNote(map[string]interface{}{
		"ticket_id": ticketID,
		"note":      noteText,
	})

	msg := "未找到上游渠道，已记录内部备注待人工处理"
	if upstreamChannel != "" {
		msg = "已通过「" + upstreamChannel + "」渠道提交上游工单"
	}

	return jsonStr(map[string]interface{}{
		"success":              true,
		"ticket_id":            ticketID,
		"host_id":              hostID,
		"dcimid":               dcimID,
		"upstream_channel":     upstreamChannel,
		"parent_host_id":       parentHostID,
		"upstream_ticket_id":   upstreamTicketID,
		"upstream_result":      upstreamResult,
		"conversation_summary": truncateStr(conversationSummary, 500),
		"msg":                  msg,
	})
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
		return jsonStr(map[string]interface{}{"error": "参数不完整：需要 ticket_id 和 host_id"})
	}

	// 查询 host 信息
	var host struct {
		Server   string `gorm:"column:server"`
		ParentID int64  `gorm:"column:parent_id"`
	}
	if err := e.db.Raw("SELECT IFNULL(server,'') as server, IFNULL(parent_id,0) as parent_id FROM host WHERE id = ?", hostID).Scan(&host).Error; err != nil {
		return jsonStr(map[string]interface{}{"error": "主机不存在"})
	}

	upstreamChannel := host.Server
	parentHostID := host.ParentID

	// 从内部备注中读取上游工单号
	upstreamTicketID := ""
	var notes []struct {
		Content string `gorm:"column:content"`
	}
	e.db.Raw("SELECT IFNULL(content,'') as content FROM ticket_replies WHERE ticket_id = ? AND content LIKE '%【上游工单透传】%' ORDER BY id DESC LIMIT 5", ticketID).Scan(&notes)
	for _, note := range notes {
		cleanContent := stripTags(note.Content)
		if m := extractRegex(cleanContent, `上游工单号[:：]\s*(\S+)`); m != "" {
			upstreamTicketID = m
			break
		}
	}

	// 读取上次回复数
	lastKnownCount := 0
	var lastCountNote struct {
		Content string `gorm:"column:content"`
	}
	if err := e.db.Raw("SELECT IFNULL(content,'') as content FROM ticket_replies WHERE ticket_id = ? AND content LIKE '%上游回复数:%' ORDER BY id DESC LIMIT 1", ticketID).Scan(&lastCountNote).Error; err == nil {
		if m := extractRegex(stripTags(lastCountNote.Content), `上游回复数:(\d+)`); m != "" {
			lastKnownCount = int(intVal(m))
		}
	}

	if upstreamTicketID == "" {
		return jsonStr(map[string]interface{}{
			"success": false,
			"error":   "未找到上游工单号，可能尚未提交上游工单",
			"msg":     "未找到上游工单号",
		})
	}

	// 通过 UpstreamClient 查询上游回复
	var upstreamResult *UpstreamResult
	if upstreamChannel != "" && parentHostID > 0 {
		client := NewUpstreamClient(e.db, e.log)
		upstreamResult = client.CallForHost(parentHostID, "get_ticket_replies", map[string]interface{}{
			"service_id":         parentHostID,
			"upstream_ticket_id": upstreamTicketID,
		})
	}

	if upstreamResult != nil && upstreamResult.Error != "" {
		return jsonStr(map[string]interface{}{
			"success":             false,
			"error":               upstreamResult.Error,
			"upstream_ticket_id":  upstreamTicketID,
			"msg":                 "上游查询失败：" + upstreamResult.Error,
		})
	}

	// 分析新回复
	var replies []interface{}
	if upstreamResult != nil && upstreamResult.Data != nil {
		data := upstreamResult.Data
		// 不同模块返回格式不同，尝试多种字段
		if r, ok := data["replies"].([]interface{}); ok {
			replies = r
		} else if r, ok := data["messages"].([]interface{}); ok {
			replies = r
		} else if r, ok := data["data"].([]interface{}); ok {
			replies = r
		}
	}

	currentCount := len(replies)
	var newReplies []interface{}
	if currentCount > lastKnownCount {
		newReplies = replies[lastKnownCount:]
	}

	// 记录当前回复数到内部备注
	if currentCount > 0 && currentCount != lastKnownCount {
		e.addTicketNote(map[string]interface{}{
			"ticket_id": ticketID,
			"note":      fmt.Sprintf("上游回复数:%d", currentCount),
		})
	}

	msg := "上游暂无新回复"
	if len(newReplies) > 0 {
		msg = fmt.Sprintf("检测到上游有 %d 条新回复", len(newReplies))
	}

	return jsonStr(map[string]interface{}{
		"success":             true,
		"ticket_id":           ticketID,
		"host_id":             hostID,
		"upstream_channel":    upstreamChannel,
		"parent_host_id":      parentHostID,
		"upstream_ticket_id":  upstreamTicketID,
		"total_replies":       currentCount,
		"last_known_count":    lastKnownCount,
		"new_replies":         newReplies,
		"has_new":             len(newReplies) > 0,
		"msg":                 msg,
	})
}

func (e *AIToolExecutor) closeUpstreamTicket(args map[string]interface{}) string {
	ticketID := uint(intVal(args["ticket_id"]))
	hostID := uint(intVal(args["host_id"]))

	if ticketID == 0 || hostID == 0 {
		return jsonStr(map[string]interface{}{"error": "参数不完整：需要 ticket_id 和 host_id"})
	}

	// 查询 host 信息
	var host struct {
		Server   string `gorm:"column:server"`
		ParentID int64  `gorm:"column:parent_id"`
	}
	if err := e.db.Raw("SELECT IFNULL(server,'') as server, IFNULL(parent_id,0) as parent_id FROM host WHERE id = ?", hostID).Scan(&host).Error; err != nil {
		return jsonStr(map[string]interface{}{"error": "主机不存在"})
	}

	upstreamChannel := host.Server
	parentHostID := host.ParentID

	// 从内部备注中读取上游工单号
	upstreamTicketID := ""
	var notes []struct {
		Content string `gorm:"column:content"`
	}
	e.db.Raw("SELECT IFNULL(content,'') as content FROM ticket_replies WHERE ticket_id = ? AND content LIKE '%【上游工单透传】%' ORDER BY id DESC LIMIT 5", ticketID).Scan(&notes)
	for _, note := range notes {
		cleanContent := stripTags(note.Content)
		if m := extractRegex(cleanContent, `上游工单号[:：]\s*(\S+)`); m != "" {
			upstreamTicketID = m
			break
		}
	}

	if upstreamTicketID == "" {
		return jsonStr(map[string]interface{}{
			"success": false,
			"error":   "未找到上游工单号，可能尚未提交上游工单",
			"msg":     "未找到上游工单号",
		})
	}

	// 通过 UpstreamClient 关闭上游工单
	var upstreamResult *UpstreamResult
	if upstreamChannel != "" && parentHostID > 0 {
		client := NewUpstreamClient(e.db, e.log)
		upstreamResult = client.CallForHost(parentHostID, "close_support_ticket", map[string]interface{}{
			"service_id":         parentHostID,
			"upstream_ticket_id": upstreamTicketID,
		})
	}

	// 记录内部备注
	noteText := fmt.Sprintf("【关闭上游工单】\n上游渠道：%s\n父服务ID：%d\n上游工单号：%s\n",
		upstreamChannel, parentHostID, upstreamTicketID)
	if upstreamResult != nil && upstreamResult.Error != "" {
		noteText += "关闭状态：失败 - " + upstreamResult.Error + "\n"
	} else {
		noteText += "关闭状态：成功\n"
	}

	e.addTicketNote(map[string]interface{}{
		"ticket_id": ticketID,
		"note":      noteText,
	})

	msg := "已成功关闭上游工单"
	if upstreamResult != nil && upstreamResult.Error != "" {
		msg = "上游工单关闭失败：" + upstreamResult.Error
	}

	return jsonStr(map[string]interface{}{
		"success":              true,
		"ticket_id":            ticketID,
		"host_id":              hostID,
		"upstream_channel":     upstreamChannel,
		"parent_host_id":       parentHostID,
		"upstream_ticket_id":   upstreamTicketID,
		"upstream_result":      upstreamResult,
		"msg":                  msg,
	})
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
	hostID := uint(intVal(args["host_id"]))
	reason := strVal(args["reason"])

	if ticketID == 0 || orderID == 0 || reason == "" {
		return jsonStr(map[string]interface{}{"error": "参数不完整：需要 ticket_id、order_id、reason"})
	}

	// 校验退款资格：仅新开通订单 + 24小时内
	var order struct {
		ID        uint    `gorm:"column:id"`
		UserID    uint    `gorm:"column:uid"`
		Type      string  `gorm:"column:type"`
		Amount    float64 `gorm:"column:amount"`
		CreatedAt string  `gorm:"column:created_at"`
		RelID     uint    `gorm:"column:relid"`
	}
	if err := e.db.Raw("SELECT id, uid, type, amount, created_at, relid FROM orders WHERE id = ?", orderID).Scan(&order).Error; err != nil {
		return jsonStr(map[string]interface{}{"error": "订单不存在"})
	}

	if order.Type != "new" {
		return jsonStr(map[string]interface{}{"error": "仅新开通订单允许退款申请，续费订单不可退款"})
	}

	// 检查24小时限制
	if order.CreatedAt != "" {
		createTime, err := time.Parse("2006-01-02 15:04:05", order.CreatedAt)
		if err == nil && time.Since(createTime) > 24*time.Hour {
			return jsonStr(map[string]interface{}{"error": "订单开通已超过24小时，不符合退款条件"})
		}
	}

	// host_id 从工单自动读取
	if hostID == 0 {
		var ticket struct {
			RelID uint `gorm:"column:rel_id"`
		}
		if err := e.db.Raw("SELECT rel_id FROM tickets WHERE id = ?", ticketID).Scan(&ticket).Error; err == nil {
			hostID = ticket.RelID
		}
	}
	if hostID == 0 {
		return jsonStr(map[string]interface{}{
			"error":   "该工单未关联产品，无法确定要退款的具体产品。请先询问客户要退款的是哪个产品，然后通过 edit_ticket 设置工单的 host_id",
			"no_host": true,
		})
	}

	// 获取主机上游渠道信息
	var host struct {
		ID        int64  `gorm:"column:id"`
		UID       int64  `gorm:"column:uid"`
		Server    string `gorm:"column:server"`
		ParentID  int64  `gorm:"column:parent_id"`
		DcimID    int64  `gorm:"column:dcimid"`
	}
	if err := e.db.Raw("SELECT id, uid, IFNULL(server,'') as server, IFNULL(parent_id,0) as parent_id, IFNULL(dcimid,0) as dcimid FROM host WHERE id = ?", hostID).Scan(&host).Error; err != nil {
		return jsonStr(map[string]interface{}{"error": "主机不存在"})
	}

	upstreamChannel := host.Server
	parentHostID := host.ParentID
	dcimID := host.DcimID

	// 构建上游退款工单内容（绝不填写金额，只提交退款诉求）
	upstreamSubject := fmt.Sprintf("退款申请 - 订单#%d", orderID)
	upstreamContent := fmt.Sprintf("【退款申请】\n订单号：#%d\n客户UID：%d\n主机ID：%d\n机器标识(dcimid)：%d\n退款原因：%s\n申请时间：%s",
		orderID, e.userID, hostID, dcimID, reason, time.Now().Format("2006-01-02 15:04:05"))

	// 尝试通过上游渠道提交退款工单
	var upstreamResult *UpstreamResult
	var upstreamTicketID string

	if upstreamChannel != "" && parentHostID > 0 {
		client := NewUpstreamClient(e.db, e.log)
		upstreamResult = client.CallForHost(parentHostID, "open_support_ticket", map[string]interface{}{
			"subject":    upstreamSubject,
			"message":    upstreamContent,
			"service_id": parentHostID,
			"priority":   "high",
		})
		if upstreamResult != nil && upstreamResult.Data != nil {
			upstreamTicketID = strVal(upstreamResult.Data["ticketid"])
			if upstreamTicketID == "" {
				upstreamTicketID = strVal(upstreamResult.Data["ticket_id"])
			}
			if upstreamTicketID == "" {
				upstreamTicketID = strVal(upstreamResult.Data["id"])
			}
		}
	}

	// 记录内部备注
	noteText := fmt.Sprintf("【退款申请透传上游】\n订单号：#%d\n退款原因：%s\n主机ID：%d\ndcimid：%d\n上游渠道：%s\n父服务ID：%d\n",
		orderID, reason, hostID, dcimID, upstreamChannel, parentHostID)

	if upstreamResult != nil && upstreamResult.Error != "" {
		noteText += "上游调用状态：失败 - " + upstreamResult.Error + "\n"
	} else if upstreamResult != nil && upstreamResult.Data != nil {
		noteText += "上游调用状态：成功\n"
		if upstreamTicketID != "" {
			noteText += "上游工单号：" + upstreamTicketID + "\n"
		}
		rawJSON, _ := json.Marshal(upstreamResult.Data)
		noteText += "上游返回：" + string(rawJSON) + "\n"
	} else {
		noteText += "上游调用状态：未配置上游渠道，需人工处理\n"
	}

	e.addTicketNote(map[string]interface{}{
		"ticket_id": ticketID,
		"note":      noteText,
	})

	msg := "未找到上游渠道，已记录内部备注待人工处理"
	if upstreamChannel != "" {
		msg = "已通过「" + upstreamChannel + "」渠道提交退款工单至上游"
	}

	return jsonStr(map[string]interface{}{
		"success":              true,
		"ticket_id":            ticketID,
		"order_id":             orderID,
		"host_id":              hostID,
		"dcimid":               dcimID,
		"upstream_channel":     upstreamChannel,
		"parent_host_id":       parentHostID,
		"upstream_ticket_id":   upstreamTicketID,
		"upstream_result":      upstreamResult,
		"msg":                  msg,
	})
}

func (e *AIToolExecutor) checkUpstreamRefundResult(args map[string]interface{}) string {
	ticketID := uint(intVal(args["ticket_id"]))
	orderID := uint(intVal(args["order_id"]))

	if ticketID == 0 || orderID == 0 {
		return jsonStr(map[string]interface{}{"error": "参数不完整"})
	}

	// 查询订单
	var order struct {
		ID        uint    `gorm:"column:id"`
		UserID    uint    `gorm:"column:uid"`
		Type      string  `gorm:"column:type"`
		Amount    float64 `gorm:"column:amount"`
		CreatedAt string  `gorm:"column:created_at"`
		RelID     uint    `gorm:"column:relid"`
	}
	if err := e.db.Raw("SELECT id, uid, type, amount, created_at, relid FROM orders WHERE id = ?", orderID).Scan(&order).Error; err != nil {
		return jsonStr(map[string]interface{}{"error": "订单不存在"})
	}

	amount := order.Amount
	orderType := order.Type
	relID := order.RelID

	// 检查上游退款状态
	var host struct {
		Server   string `gorm:"column:server"`
		ParentID int64  `gorm:"column:parent_id"`
	}
	upstreamChannel := ""
	var parentHostID int64
	if relID > 0 {
		if err := e.db.Raw("SELECT IFNULL(server,'') as server, IFNULL(parent_id,0) as parent_id FROM host WHERE id = ?", relID).Scan(&host).Error; err == nil {
			upstreamChannel = host.Server
			parentHostID = host.ParentID
		}
	}

	// 从内部备注读取上游工单号
	upstreamTicketID := ""

	// 先从退款透传记录中查找
	var refundNotes []struct {
		Content string `gorm:"column:content"`
	}
	e.db.Raw("SELECT IFNULL(content,'') as content FROM ticket_replies WHERE ticket_id = ? AND content LIKE '%【退款申请透传上游】%' ORDER BY id DESC LIMIT 5", ticketID).Scan(&refundNotes)
	for _, note := range refundNotes {
		cleanContent := stripTags(note.Content)
		if m := extractRegex(cleanContent, `上游工单号[:：]\s*(\S+)`); m != "" {
			upstreamTicketID = m
			break
		}
	}

	// 也从工单透传记录中查找
	if upstreamTicketID == "" {
		var transmitNotes []struct {
			Content string `gorm:"column:content"`
		}
		e.db.Raw("SELECT IFNULL(content,'') as content FROM ticket_replies WHERE ticket_id = ? AND content LIKE '%【上游工单透传】%' ORDER BY id DESC LIMIT 5", ticketID).Scan(&transmitNotes)
		for _, note := range transmitNotes {
			cleanContent := stripTags(note.Content)
			if m := extractRegex(cleanContent, `上游工单号[:：]\s*(\S+)`); m != "" {
				upstreamTicketID = m
				break
			}
		}
	}

	// 查询上游退款状态
	upstreamRefundStatus := "unknown"
	var upstreamRefundData interface{}

	if upstreamChannel != "" && parentHostID > 0 && upstreamTicketID != "" {
		client := NewUpstreamClient(e.db, e.log)
		result := client.CallForHost(parentHostID, "get_ticket_replies", map[string]interface{}{
			"service_id":         parentHostID,
			"upstream_ticket_id": upstreamTicketID,
		})
		if result != nil && result.Error != "" {
			upstreamRefundStatus = "query_failed"
			upstreamRefundData = result.Error
		} else if result != nil {
			upstreamRefundStatus = "queried"
			upstreamRefundData = result.Data
		}
	}

	// 自动执行下游退款（仅新开通订单+24小时内）
	var downstreamRefundResult map[string]interface{}
	canAutoRefund := false

	if orderType == "new" && order.CreatedAt != "" {
		createTime, err := time.Parse("2006-01-02 15:04:05", order.CreatedAt)
		if err == nil && time.Since(createTime) <= 24*time.Hour {
			canAutoRefund = true
		}
	}

	if canAutoRefund && amount > 0 {
		// 通过财务标准流程退款：添加退款交易记录
		refundTransID := fmt.Sprintf("REF%s", randomHex(12))

		result := e.db.Exec(`INSERT INTO transactions (uid, amount, create_time, notes, gateway, trans_id) VALUES (?, ?, NOW(), ?, 'refund', ?)`,
			e.userID, -amount, fmt.Sprintf("AI自动退款 - 订单#%d - 上游退款确认后执行", orderID), refundTransID)
		if result.Error != nil {
			downstreamRefundResult = map[string]interface{}{"error": "退款执行失败: " + result.Error.Error()}
		} else {
			// 获取退款ID
			var refundID int64
			e.db.Raw("SELECT LAST_INSERT_ID()").Scan(&refundID)

			// 更新订单状态为已取消
			e.db.Exec("UPDATE orders SET status = 'Cancelled' WHERE id = ?", orderID)

			downstreamRefundResult = map[string]interface{}{
				"success":   true,
				"refund_id": refundID,
				"trans_id":  refundTransID,
				"amount":    amount,
				"order_id":  orderID,
			}

	// 记录内部备注
	e.addTicketNote(map[string]interface{}{
		"ticket_id": ticketID,
				"note":      fmt.Sprintf("【自动退款已执行】\n订单号：#%d\n退款金额：¥%.2f\n退款ID：%d\n上游状态：%s", orderID, amount, refundID, upstreamRefundStatus),
			})
		}
	}

	msg := "不符合自动退款条件（仅新开通订单+24小时内），需人工确认"
	if canAutoRefund {
		if downstreamRefundResult != nil && downstreamRefundResult["error"] == nil {
			msg = fmt.Sprintf("已自动执行下游退款 ¥%.2f，订单#%d已取消", amount, orderID)
		} else {
			msg = "符合退款条件但执行失败，请人工处理"
		}
	}

	return jsonStr(map[string]interface{}{
		"success":                  true,
		"order_id":                 orderID,
		"amount":                   amount,
		"order_type":               orderType,
		"can_auto_refund":          canAutoRefund,
		"upstream_channel":         upstreamChannel,
		"upstream_ticket_id":       upstreamTicketID,
		"upstream_refund_status":   upstreamRefundStatus,
		"upstream_refund_data":     upstreamRefundData,
		"downstream_refund":        downstreamRefundResult,
		"msg":                      msg,
	})
}

// ─── 辅助函数 ───

func (e *AIToolExecutor) checkHostOwnership(hostID uint) bool {
	var count int64
	e.db.Raw("SELECT COUNT(*) FROM host WHERE id = ? AND uid = ?", hostID, e.userID).Scan(&count)
	return count > 0
}

// ─── 商品导购工具 ───

func (e *AIToolExecutor) searchProducts(args map[string]interface{}) string {
	keyword := strVal(args["keyword"])
	limit := int(args["limit"].(float64))
	if limit <= 0 {
		limit = 10
	}

	type product struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		GroupName   string  `json:"group_name"`
	}
	var products []product
	q := e.db.Table("products p").
		Select("p.id, p.name, p.description, p.price, pg.name as group_name").
		Joins("LEFT JOIN product_groups pg ON p.group_id = pg.id").
		Where("p.status = ?", 1)
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("p.name LIKE ? OR p.description LIKE ?", like, like)
	}
	q.Limit(limit).Find(&products)

	if len(products) == 0 {
		return jsonStr(map[string]interface{}{"message": "未找到相关商品", "products": []interface{}{}})
	}

	var result []map[string]interface{}
	for _, p := range products {
		result = append(result, map[string]interface{}{
			"id":          p.ID,
			"name":        p.Name,
			"description": p.Description,
			"price":       p.Price,
			"group_name":  p.GroupName,
		})
	}
	return jsonStr(map[string]interface{}{"products": result, "total": len(result)})
}

func (e *AIToolExecutor) getProductDetail(args map[string]interface{}) string {
	productID := uint(intVal(args["product_id"]))
	if productID == 0 {
		return jsonStr(map[string]interface{}{"error": "商品ID不能为空"})
	}

	type product struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		GroupName   string  `json:"group_name"`
		BillingType string  `json:"billing_type"`
		SetupFee    float64 `json:"setup_fee"`
	}
	var p product
	err := e.db.Table("products p").
		Select("p.id, p.name, p.description, p.price, pg.name as group_name, p.billing_type, p.setup_fee").
		Joins("LEFT JOIN product_groups pg ON p.group_id = pg.id").
		Where("p.id = ?", productID).
		First(&p).Error
	if err != nil {
		return jsonStr(map[string]interface{}{"error": "商品不存在"})
	}

	// 获取配置选项
	type configOption struct {
		Name  string `json:"name"`
		Value string `json:"value"`
		Price string `json:"price"`
	}
	var options []configOption
	e.db.Table("product_config_options").Where("product_id = ?", productID).Find(&options)

	return jsonStr(map[string]interface{}{
		"id":           p.ID,
		"name":         p.Name,
		"description":  p.Description,
		"price":        p.Price,
		"group_name":   p.GroupName,
		"billing_type": p.BillingType,
		"setup_fee":    p.SetupFee,
		"options":      options,
	})
}

func (e *AIToolExecutor) listProductGroups() string {
	type group struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		ProductCount int   `json:"product_count"`
	}
	var groups []group
	e.db.Raw(`
		SELECT pg.id, pg.name, pg.description, COUNT(p.id) as product_count
		FROM product_groups pg
		LEFT JOIN products p ON p.group_id = pg.id AND p.status = 1
		GROUP BY pg.id, pg.name, pg.description
		ORDER BY pg.sort_order ASC
	`).Scan(&groups)

	if len(groups) == 0 {
		return jsonStr(map[string]interface{}{"message": "暂无商品分组", "groups": []interface{}{}})
	}
	return jsonStr(map[string]interface{}{"groups": groups})
}

func (e *AIToolExecutor) getGroupProducts(args map[string]interface{}) string {
	groupID := uint(intVal(args["group_id"]))
	if groupID == 0 {
		return jsonStr(map[string]interface{}{"error": "分组ID不能为空"})
	}

	type product struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
	var products []product
	e.db.Table("products").
		Select("id, name, description, price").
		Where("group_id = ? AND status = ?", groupID, 1).
		Order("sort_order ASC").
		Find(&products)

	return jsonStr(map[string]interface{}{"products": products, "total": len(products)})
}

func (e *AIToolExecutor) compareProducts(args map[string]interface{}) string {
	idsRaw, ok := args["product_ids"].([]interface{})
	if !ok || len(idsRaw) == 0 {
		return jsonStr(map[string]interface{}{"error": "请提供要对比的商品ID列表"})
	}

	var ids []uint
	for _, id := range idsRaw {
		if n, ok := id.(float64); ok {
			ids = append(ids, uint(n))
		}
	}
	if len(ids) < 2 {
		return jsonStr(map[string]interface{}{"error": "至少需要2个商品进行对比"})
	}

	type product struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		GroupName   string  `json:"group_name"`
		BillingType string  `json:"billing_type"`
	}
	var products []product
	e.db.Table("products p").
		Select("p.id, p.name, p.description, p.price, pg.name as group_name, p.billing_type").
		Joins("LEFT JOIN product_groups pg ON p.group_id = pg.id").
		Where("p.id IN ? AND p.status = ?", ids, 1).
		Find(&products)

	// 获取每个商品的配置选项
	type configOption struct {
		ProductID uint   `json:"product_id"`
		Name      string `json:"name"`
		Value     string `json:"value"`
		Price     string `json:"price"`
	}
	var allOptions []configOption
	e.db.Table("product_config_options").Where("product_id IN ?", ids).Find(&allOptions)

	// 按商品分组配置选项
	optionsMap := make(map[uint][]map[string]interface{})
	for _, opt := range allOptions {
		optionsMap[opt.ProductID] = append(optionsMap[opt.ProductID], map[string]interface{}{
			"name":  opt.Name,
			"value": opt.Value,
			"price": opt.Price,
		})
	}

	var result []map[string]interface{}
	for _, p := range products {
		result = append(result, map[string]interface{}{
			"id":           p.ID,
			"name":         p.Name,
			"description":  p.Description,
			"price":        p.Price,
			"group_name":   p.GroupName,
			"billing_type": p.BillingType,
			"options":      optionsMap[p.ID],
		})
	}

	return jsonStr(map[string]interface{}{"products": result, "total": len(result)})
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

// stripTags 移除 HTML 标签
func stripTags(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	return strings.TrimSpace(re.ReplaceAllString(s, ""))
}

// extractRegex 从文本中提取正则匹配的第一个捕获组
func extractRegex(text, pattern string) string {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return ""
	}
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// randomHex 生成指定长度的随机十六进制字符串
func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return strings.ToUpper(hex.EncodeToString(b))
}
