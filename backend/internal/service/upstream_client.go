package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// UpstreamClient 上游 API 统一调用客户端
// 支持三种上游系统：zjmf（模块函数调用）、v10（REST API）、anchorfinance（HTTP API）
type UpstreamClient struct {
	db     *gorm.DB
	log    *logger.Logger
	client *http.Client
}

// NewUpstreamClient 创建上游客户端
func NewUpstreamClient(db *gorm.DB, log *logger.Logger) *UpstreamClient {
	return &UpstreamClient{
		db:     db,
		log:    log,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// UpstreamResult 上游调用结果
type UpstreamResult struct {
	Data  map[string]interface{} `json:"data,omitempty"`
	Error string                 `json:"error,omitempty"`
}

// Call 统一调用入口，根据上游类型自动路由
// upstreamType: zjmf / v10 / anchorfinance
// action: open_support_ticket / get_ticket_replies / close_support_ticket / get_order_details
// params: 业务参数，其中 service_id 为必须（对应 host 表的 parent_id）
func (c *UpstreamClient) Call(upstreamType string, action string, params map[string]interface{}) *UpstreamResult {
	switch upstreamType {
	case "v10":
		return c.v10Call(action, params)
	case "anchorfinance":
		return c.anchorFinanceCall(action, params)
	default:
		return c.zjmfCall(action, params)
	}
}

// CallForHost 根据 host 记录自动判断上游类型并调用
// serviceId: 上游服务的 host ID（即 parent_id）
func (c *UpstreamClient) CallForHost(serviceId int64, action string, params map[string]interface{}) *UpstreamResult {
	if serviceId <= 0 {
		return &UpstreamResult{Error: "缺少上游服务ID"}
	}

	// 查询 host 记录
	var host struct {
		ID        int64  `gorm:"column:id"`
		ProductID int64  `gorm:"column:productid"`
		Server    string `gorm:"column:server"`
		ParentID  int64  `gorm:"column:parent_id"`
		DcimID    int64  `gorm:"column:dcimid"`
	}
	if err := c.db.Raw("SELECT id, productid, IFNULL(server,'') as server, IFNULL(parent_id,0) as parent_id, IFNULL(dcimid,0) as dcimid FROM host WHERE id = ?", serviceId).Scan(&host).Error; err != nil {
		return &UpstreamResult{Error: "上游服务不存在"}
	}

	// 判断上游系统类型：通过产品的 zjmf_api_id 查 zjmf_finance_api.type
	upstreamType := "zjmf"
	var zjmfApiID int64
	var product struct {
		ZjmfApiID int64 `gorm:"column:zjmf_api_id"`
	}
	if err := c.db.Raw("SELECT IFNULL(zjmf_api_id,0) as zjmf_api_id FROM products WHERE id = ?", host.ProductID).Scan(&product).Error; err == nil && product.ZjmfApiID > 0 {
		zjmfApiID = product.ZjmfApiID
		var api struct {
			Type string `gorm:"column:type"`
		}
		if err := c.db.Raw("SELECT IFNULL(type,'') as type FROM zjmf_finance_api WHERE id = ?", zjmfApiID).Scan(&api).Error; err == nil && api.Type == "v10" {
			upstreamType = "v10"
		}
	}

	// 注入 host 信息到 params
	params["_host"] = host
	params["_zjmf_api_id"] = zjmfApiID

	c.log.WithField("upstream_type", upstreamType).
		WithField("service_id", serviceId).
		WithField("action", action).
		Info("上游调用路由")

	return c.Call(upstreamType, action, params)
}

// ─────────────────────────────────────────────
// zjmf 模块函数调用
// ─────────────────────────────────────────────

func (c *UpstreamClient) zjmfCall(action string, params map[string]interface{}) *UpstreamResult {
	host, _ := params["_host"].(struct {
		ID        int64  `gorm:"column:id"`
		ProductID int64  `gorm:"column:productid"`
		Server    string `gorm:"column:server"`
		ParentID  int64  `gorm:"column:parent_id"`
		DcimID    int64  `gorm:"column:dcimid"`
	})

	// 获取上游服务信息
	var upstream struct {
		URL      string `gorm:"column:url"`
		APIKey   string `gorm:"column:api_key"`
		Username string `gorm:"column:username"`
	}

	// 尝试从 ticket_upstreams 表获取配置
	if host.Server != "" {
		c.db.Raw("SELECT IFNULL(url,'') as url, IFNULL(api_key,'') as api_key, IFNULL(username,'') as username FROM ticket_upstreams WHERE name = ? AND status = 1 LIMIT 1", host.Server).Scan(&upstream)
	}

	// 如果没有配置，尝试通过 zjmf_finance_api 获取
	if upstream.URL == "" {
		var api struct {
			URL      string `gorm:"column:url"`
			APIKey   string `gorm:"column:api_key"`
			Username string `gorm:"column:username"`
		}
		zjmfApiID, _ := params["_zjmf_api_id"].(int64)
		if zjmfApiID > 0 {
			if err := c.db.Raw("SELECT IFNULL(url,'') as url, IFNULL(api_key,'') as api_key, IFNULL(username,'') as username FROM zjmf_finance_api WHERE id = ?", zjmfApiID).Scan(&api).Error; err == nil {
				upstream = api
			}
		}
	}

	if upstream.URL == "" {
		return &UpstreamResult{Error: "未配置上游模块或无法获取上游地址"}
	}

	baseURL := strings.TrimRight(upstream.URL, "/")

	switch action {
	case "open_support_ticket":
		return c.zjmfCreateTicket(baseURL, upstream.APIKey, params)
	case "get_ticket_replies":
		return c.zjmfGetTicketReplies(baseURL, upstream.APIKey, params)
	case "close_support_ticket":
		return c.zjmfCloseTicket(baseURL, upstream.APIKey, params)
	case "get_order_details":
		return c.zjmfGetOrderDetails(baseURL, upstream.APIKey, params)
	default:
		return &UpstreamResult{Error: "不支持的操作: " + action}
	}
}

func (c *UpstreamClient) zjmfCreateTicket(baseURL, apiKey string, params map[string]interface{}) *UpstreamResult {
	subject := strVal(params["subject"])
	message := strVal(params["message"])
	priority := strVal(params["priority"])
	if priority == "" {
		priority = "medium"
	}
	serviceID := intVal(params["service_id"])

	payload := map[string]interface{}{
		"subject":  subject,
		"message":  message,
		"priority": priority,
	}
	if serviceID > 0 {
		payload["service_id"] = serviceID
	}

	url := baseURL + "/api/v1/ticket/create"
	result := c.httpPost(url, apiKey, payload)
	if result.Error != "" {
		return result
	}

	// 解析 zjmf 响应
	data := result.Data
	ticketID := ""
	if tid, ok := data["tid"].(string); ok && tid != "" {
		ticketID = tid
	} else if id, ok := data["id"].(string); ok && id != "" {
		ticketID = id
	} else if id, ok := data["ticketid"].(string); ok && id != "" {
		ticketID = id
	} else if id, ok := data["ticket_id"].(string); ok && id != "" {
		ticketID = id
	}

	return &UpstreamResult{
		Data: map[string]interface{}{
			"ticketid": ticketID,
			"status":   "submitted",
			"raw":      data,
		},
	}
}

func (c *UpstreamClient) zjmfGetTicketReplies(baseURL, apiKey string, params map[string]interface{}) *UpstreamResult {
	ticketID := strVal(params["upstream_ticket_id"])
	if ticketID == "" {
		return &UpstreamResult{Error: "缺少上游工单ID"}
	}

	// 尝试多个端点
	endpoints := []string{
		baseURL + "/api/v1/ticket/" + ticketID + "/replies",
		baseURL + "/api/v1/ticket/" + ticketID,
	}

	for _, url := range endpoints {
		result := c.httpGet(url, apiKey)
		if result.Error == "" {
			return result
		}
	}

	return &UpstreamResult{Error: "上游模块不支持查询工单回复"}
}

func (c *UpstreamClient) zjmfCloseTicket(baseURL, apiKey string, params map[string]interface{}) *UpstreamResult {
	ticketID := strVal(params["upstream_ticket_id"])
	if ticketID == "" {
		return &UpstreamResult{Error: "缺少上游工单ID"}
	}

	endpoints := []string{
		baseURL + "/api/v1/ticket/" + ticketID + "/close",
		baseURL + "/api/v1/ticket/" + ticketID + "/Close",
	}

	for _, url := range endpoints {
		result := c.httpPut(url, apiKey, nil)
		if result.Error == "" {
			return &UpstreamResult{
				Data: map[string]interface{}{
					"status": "closed",
					"raw":    result.Data,
				},
			}
		}
	}

	return &UpstreamResult{Error: "上游模块不支持关闭工单"}
}

func (c *UpstreamClient) zjmfGetOrderDetails(baseURL, apiKey string, params map[string]interface{}) *UpstreamResult {
	serviceID := intVal(params["service_id"])
	if serviceID <= 0 {
		return &UpstreamResult{Error: "缺少服务ID"}
	}

	url := fmt.Sprintf("%s/api/v1/service/%d", baseURL, serviceID)
	result := c.httpGet(url, apiKey)
	if result.Error != "" {
		return &UpstreamResult{
			Data: map[string]interface{}{
				"service_id": serviceID,
				"status":     "Unknown",
				"msg":        "查询失败: " + result.Error,
			},
		}
	}

	return result
}

// ─────────────────────────────────────────────
// V10 REST API 调用
// ─────────────────────────────────────────────

func (c *UpstreamClient) v10Call(action string, params map[string]interface{}) *UpstreamResult {
	zjmfApiID, _ := params["_zjmf_api_id"].(int64)
	if zjmfApiID <= 0 {
		return &UpstreamResult{Error: "V10上游API配置不存在"}
	}

	var api struct {
		URL    string `gorm:"column:url"`
		APIKey string `gorm:"column:api_key"`
	}
	if err := c.db.Raw("SELECT IFNULL(url,'') as url, IFNULL(api_key,'') as api_key FROM zjmf_finance_api WHERE id = ?", zjmfApiID).Scan(&api).Error; err != nil {
		return &UpstreamResult{Error: "V10上游API配置不存在"}
	}

	baseURL := strings.TrimRight(api.URL, "/")

	host, _ := params["_host"].(struct {
		ID        int64  `gorm:"column:id"`
		ProductID int64  `gorm:"column:productid"`
		Server    string `gorm:"column:server"`
		ParentID  int64  `gorm:"column:parent_id"`
		DcimID    int64  `gorm:"column:dcimid"`
	})

	switch action {
	case "open_support_ticket":
		return c.v10CreateTicket(baseURL, api.APIKey, params)
	case "get_ticket_replies":
		return c.v10GetTicketReplies(baseURL, api.APIKey, params)
	case "close_support_ticket":
		return c.v10CloseTicket(baseURL, api.APIKey, params)
	case "get_order_details":
		return c.v10GetOrderDetails(baseURL, api.APIKey, host.DcimID)
	default:
		return &UpstreamResult{Error: "V10不支持的操作: " + action}
	}
}

func (c *UpstreamClient) v10CreateTicket(baseURL, apiKey string, params map[string]interface{}) *UpstreamResult {
	subject := strVal(params["subject"])
	message := strVal(params["message"])
	priority := strVal(params["priority"])
	if priority == "" {
		priority = "medium"
	}
	if subject == "" || message == "" {
		return &UpstreamResult{Error: "工单标题和内容不能为空"}
	}

	payload := map[string]interface{}{
		"title":    subject,
		"content":  message,
		"priority": priority,
	}

	url := baseURL + "/console/v1/ticket"
	result := c.httpPost(url, apiKey, payload)
	if result.Error != "" {
		return &UpstreamResult{Error: "V10创建工单失败: " + result.Error}
	}

	data := result.Data
	ticketID := ""
	if id, ok := data["id"].(string); ok {
		ticketID = id
	} else if id, ok := data["ticket_id"].(string); ok {
		ticketID = id
	} else if id, ok := data["id"].(float64); ok {
		ticketID = fmt.Sprintf("%.0f", id)
	}

	return &UpstreamResult{
		Data: map[string]interface{}{
			"ticketid": ticketID,
			"status":   "submitted",
			"raw":      data,
		},
	}
}

func (c *UpstreamClient) v10GetTicketReplies(baseURL, apiKey string, params map[string]interface{}) *UpstreamResult {
	ticketID := strVal(params["upstream_ticket_id"])
	if ticketID == "" {
		return &UpstreamResult{Error: "缺少上游工单ID"}
	}

	url := baseURL + "/console/v1/ticket/" + ticketID
	result := c.httpGet(url, apiKey)
	if result.Error != "" {
		return &UpstreamResult{Error: "V10查询工单失败: " + result.Error}
	}
	return result
}

func (c *UpstreamClient) v10CloseTicket(baseURL, apiKey string, params map[string]interface{}) *UpstreamResult {
	ticketID := strVal(params["upstream_ticket_id"])
	if ticketID == "" {
		return &UpstreamResult{Error: "缺少上游工单ID"}
	}

	url := baseURL + "/console/v1/ticket/" + ticketID + "/close"
	result := c.httpPut(url, apiKey, nil)
	if result.Error != "" {
		return &UpstreamResult{Error: "V10关闭工单失败: " + result.Error}
	}

	return &UpstreamResult{
		Data: map[string]interface{}{
			"status": "closed",
			"raw":    result.Data,
		},
	}
}

func (c *UpstreamClient) v10GetOrderDetails(baseURL, apiKey string, dcimID int64) *UpstreamResult {
	if dcimID <= 0 {
		return &UpstreamResult{
			Data: map[string]interface{}{
				"status": "Unknown",
				"msg":    "无dcimid，无法查询V10上游",
			},
		}
	}

	url := fmt.Sprintf("%s/console/v1/host/%d", baseURL, dcimID)
	result := c.httpGet(url, apiKey)
	if result.Error != "" {
		return &UpstreamResult{Error: "V10查询主机失败: " + result.Error}
	}

	data := result.Data
	hd := data
	if hostData, ok := data["host"].(map[string]interface{}); ok {
		hd = hostData
	}

	return &UpstreamResult{
		Data: map[string]interface{}{
			"status":       strVal(hd["domainstatus"]),
			"domainstatus": strVal(hd["domainstatus"]),
			"hostname":     strVal(hd["domain"]),
			"ip":           strVal(hd["dedicatedip"]),
			"os":           strVal(hd["os"]),
			"raw":          data,
		},
	}
}

// ─────────────────────────────────────────────
// AnchorFinance HTTP API 调用
// ─────────────────────────────────────────────

func (c *UpstreamClient) anchorFinanceCall(action string, params map[string]interface{}) *UpstreamResult {
	host, _ := params["_host"].(struct {
		ID        int64  `gorm:"column:id"`
		ProductID int64  `gorm:"column:productid"`
		Server    string `gorm:"column:server"`
		ParentID  int64  `gorm:"column:parent_id"`
		DcimID    int64  `gorm:"column:dcimid"`
	})

	// 从 ticket_upstreams 表获取配置
	var upstream struct {
		URL    string `gorm:"column:url"`
		APIKey string `gorm:"column:api_key"`
	}
	if host.Server != "" {
		if err := c.db.Raw("SELECT IFNULL(url,'') as url, IFNULL(api_key,'') as api_key FROM ticket_upstreams WHERE name = ? AND type = 'anchorfinance' AND status = 1 LIMIT 1", host.Server).Scan(&upstream).Error; err != nil || upstream.URL == "" {
			return &UpstreamResult{Error: "未配置AnchorFinance上游"}
		}
	} else {
		return &UpstreamResult{Error: "未配置上游渠道"}
	}

	baseURL := strings.TrimRight(upstream.URL, "/")

	switch action {
	case "open_support_ticket":
		return c.afCreateTicket(baseURL, upstream.APIKey, params)
	case "get_ticket_replies":
		return c.afGetTicketReplies(baseURL, upstream.APIKey, params)
	case "close_support_ticket":
		return c.afCloseTicket(baseURL, upstream.APIKey, params)
	case "get_order_details":
		return c.afGetOrderDetails(baseURL, upstream.APIKey, params)
	default:
		return &UpstreamResult{Error: "AnchorFinance不支持的操作: " + action}
	}
}

func (c *UpstreamClient) afCreateTicket(baseURL, apiKey string, params map[string]interface{}) *UpstreamResult {
	subject := strVal(params["subject"])
	message := strVal(params["message"])
	if subject == "" || message == "" {
		return &UpstreamResult{Error: "工单标题和内容不能为空"}
	}

	payload := map[string]interface{}{
		"title":   subject,
		"content": message,
		"is_api":  true,
	}
	if deptID := intVal(params["department_id"]); deptID > 0 {
		payload["department_id"] = deptID
	}

	url := baseURL + "/api/v1/tickets/create"
	result := c.httpPost(url, apiKey, payload)
	if result.Error != "" {
		return &UpstreamResult{Error: "AnchorFinance创建工单失败: " + result.Error}
	}

	data := result.Data
	ticketID := strVal(data["ticket_id"])
	if ticketID == "" {
		ticketID = strVal(data["id"])
	}

	return &UpstreamResult{
		Data: map[string]interface{}{
			"ticketid": ticketID,
			"status":   "submitted",
			"raw":      data,
		},
	}
}

func (c *UpstreamClient) afGetTicketReplies(baseURL, apiKey string, params map[string]interface{}) *UpstreamResult {
	ticketID := strVal(params["upstream_ticket_id"])
	if ticketID == "" {
		return &UpstreamResult{Error: "缺少上游工单ID"}
	}

	url := baseURL + "/api/v1/tickets/" + ticketID + "/replies"
	result := c.httpGet(url, apiKey)
	if result.Error != "" {
		return &UpstreamResult{Error: "查询失败: " + result.Error}
	}
	return result
}

func (c *UpstreamClient) afCloseTicket(baseURL, apiKey string, params map[string]interface{}) *UpstreamResult {
	ticketID := strVal(params["upstream_ticket_id"])
	if ticketID == "" {
		return &UpstreamResult{Error: "缺少上游工单ID"}
	}

	url := baseURL + "/api/v1/tickets/" + ticketID + "/close"
	result := c.httpPut(url, apiKey, nil)
	if result.Error != "" {
		return &UpstreamResult{Error: "关闭失败: " + result.Error}
	}

	return &UpstreamResult{
		Data: map[string]interface{}{
			"status": "closed",
			"raw":    result.Data,
		},
	}
}

func (c *UpstreamClient) afGetOrderDetails(baseURL, apiKey string, params map[string]interface{}) *UpstreamResult {
	serviceID := intVal(params["service_id"])
	if serviceID <= 0 {
		return &UpstreamResult{Error: "缺少服务ID"}
	}

	url := fmt.Sprintf("%s/api/v1/hosts/%d", baseURL, serviceID)
	result := c.httpGet(url, apiKey)
	if result.Error != "" {
		return &UpstreamResult{Error: "查询失败: " + result.Error}
	}
	return result
}

// ─────────────────────────────────────────────
// HTTP 请求工具方法
// ─────────────────────────────────────────────

func (c *UpstreamClient) httpPost(url, apiKey string, payload map[string]interface{}) *UpstreamResult {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return &UpstreamResult{Error: "序列化参数失败: " + err.Error()}
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return &UpstreamResult{Error: "创建请求失败: " + err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}

	return c.doRequest(req)
}

func (c *UpstreamClient) httpGet(url, apiKey string) *UpstreamResult {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &UpstreamResult{Error: "创建请求失败: " + err.Error()}
	}
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	return c.doRequest(req)
}

func (c *UpstreamClient) httpPut(url, apiKey string, payload map[string]interface{}) *UpstreamResult {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return &UpstreamResult{Error: "序列化参数失败: " + err.Error()}
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest("PUT", url, body)
	if err != nil {
		return &UpstreamResult{Error: "创建请求失败: " + err.Error()}
	}
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+apiKey)
	}
	return c.doRequest(req)
}

func (c *UpstreamClient) doRequest(req *http.Request) *UpstreamResult {
	resp, err := c.client.Do(req)
	if err != nil {
		return &UpstreamResult{Error: "请求失败: " + err.Error()}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if resp.StatusCode >= 400 {
		// 尝试解析错误信息
		var errResp map[string]interface{}
		if json.Unmarshal(body, &errResp) == nil {
			if msg, ok := errResp["msg"].(string); ok {
				return &UpstreamResult{Error: msg}
			}
			if msg, ok := errResp["message"].(string); ok {
				return &UpstreamResult{Error: msg}
			}
			if msg, ok := errResp["error"].(string); ok {
				return &UpstreamResult{Error: msg}
			}
		}
		return &UpstreamResult{Error: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, truncateStr(bodyStr, 200))}
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		// 非 JSON 响应，包装为 data
		return &UpstreamResult{
			Data: map[string]interface{}{
				"raw": bodyStr,
			},
		}
	}

	// 检查业务层错误（某些 API 在 200 中返回错误）
	if status, ok := data["status"].(float64); ok && status != 200 {
		msg := strVal(data["msg"])
		if msg == "" {
			msg = strVal(data["message"])
		}
		if msg != "" {
			return &UpstreamResult{Error: msg}
		}
	}

	return &UpstreamResult{Data: data}
}

func truncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
