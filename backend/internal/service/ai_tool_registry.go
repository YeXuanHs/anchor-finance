package service

import (
	"encoding/json"
)

// AIToolCategory 工具分类
type AIToolCategory struct {
	Key     string    `json:"key"`
	Name    string    `json:"name"`
	Tools   []AITool  `json:"tools"`
}

// AITool 单个工具定义
type AITool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Category    string                 `json:"category"`
	RiskLevel   string                 `json:"risk_level"`
	Schema      map[string]interface{} `json:"schema"`
	Enabled     bool                   `json:"enabled"`
}

// AIToolRegistry 工具注册表
type AIToolRegistry struct {
	tools     map[string]*AITool
	categories []AIToolCategory
}

// NewAIToolRegistry 创建工具注册表
func NewAIToolRegistry() *AIToolRegistry {
	r := &AIToolRegistry{
		tools: make(map[string]*AITool),
	}
	r.initCategories()
	return r
}

func (r *AIToolRegistry) initCategories() {
	r.categories = []AIToolCategory{
		r.buildReadonlyCategory(),
		r.buildServerOpsCategory(),
		r.buildTicketOpsCategory(),
		r.buildFinanceCategory(),
		r.buildShoppingCategory(),
		r.buildGeneralCategory(),
	}
	for _, cat := range r.categories {
		for i := range cat.Tools {
			cat.Tools[i].Category = cat.Key
			r.tools[cat.Tools[i].Name] = &cat.Tools[i]
		}
	}
}

// GetAllCategories 获取所有分类
func (r *AIToolRegistry) GetAllCategories() []AIToolCategory {
	return r.categories
}

// GetTool 获取单个工具
func (r *AIToolRegistry) GetTool(name string) (*AITool, bool) {
	t, ok := r.tools[name]
	return t, ok
}

// GetOpenAITools 获取 OpenAI function calling 格式的工具列表
func (r *AIToolRegistry) GetOpenAITools(enabledTools []string) []map[string]interface{} {
	enabled := make(map[string]bool)
	for _, name := range enabledTools {
		enabled[name] = true
	}

	var tools []map[string]interface{}
	for _, t := range r.tools {
		if !enabled[t.Name] {
			continue
		}
		tools = append(tools, map[string]interface{}{
			"type":     "function",
			"function": t.Schema,
		})
	}
	return tools
}

// GetEnabledToolNames 获取所有工具名
func (r *AIToolRegistry) GetAllToolNames() []string {
	var names []string
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// ─── 分类定义 ───

func (r *AIToolRegistry) buildReadonlyCategory() AIToolCategory {
	return AIToolCategory{
		Key:  "readonly",
		Name: "只读查询",
		Tools: []AITool{
			{
				Name:        "query_client_info",
				Description: "查询客户基本信息（用户名、邮箱、手机、公司、余额、注册时间等）",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "query_client_info",
					"description": "查询客户基本信息",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"client_id": map[string]interface{}{"type": "integer", "description": "客户ID"},
						},
						"required": []string{"client_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "query_host_service",
				Description: "查询主机/服务详情（产品名、状态、IP、到期时间、连接信息、specs配置）",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "query_host_service",
					"description": "查询主机服务详情",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"host_id": map[string]interface{}{"type": "integer", "description": "主机ID"},
						},
						"required": []string{"host_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "query_client_hosts",
				Description: "查询客户的所有主机列表（最多30条）",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "query_client_hosts",
					"description": "查询客户的所有主机列表",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"client_id": map[string]interface{}{"type": "integer", "description": "客户ID"},
						},
						"required": []string{"client_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "query_order",
				Description: "查询订单详情（类型、状态、金额、支付方式、开通时长）",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "query_order",
					"description": "查询订单详情",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"order_id": map[string]interface{}{"type": "integer", "description": "订单ID"},
						},
						"required": []string{"order_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "query_product",
				Description: "查询产品详情（名称、描述、定价）",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "query_product",
					"description": "查询产品详情",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"product_id": map[string]interface{}{"type": "integer", "description": "产品ID"},
						},
						"required": []string{"product_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "search_products",
				Description: "按关键词搜索产品",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "search_products",
					"description": "按关键词搜索产品",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"keyword": map[string]interface{}{"type": "string", "description": "搜索关键词"},
						},
						"required": []string{"keyword"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "query_ticket_history",
				Description: "查询客户的历史工单列表",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "query_ticket_history",
					"description": "查询客户历史工单",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"client_id": map[string]interface{}{"type": "integer", "description": "客户ID"},
							"limit":     map[string]interface{}{"type": "integer", "description": "返回数量，默认10，最大30"},
						},
						"required": []string{"client_id"},
					},
				},
				Enabled: true,
			},
		},
	}
}

func (r *AIToolRegistry) buildServerOpsCategory() AIToolCategory {
	return AIToolCategory{
		Key:  "server_ops",
		Name: "服务器运维",
		Tools: []AITool{
			{
				Name:        "check_server_status",
				Description: "检查服务器运行状态和电源状态",
				RiskLevel:   "medium",
				Schema: map[string]interface{}{
					"name":        "check_server_status",
					"description": "检查服务器运行状态",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"host_id": map[string]interface{}{"type": "integer", "description": "主机ID"},
						},
						"required": []string{"host_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "ping_check",
				Description: "对服务器IP执行ping检测",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "ping_check",
					"description": "Ping检测服务器IP",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"ip": map[string]interface{}{"type": "string", "description": "服务器IP地址"},
						},
						"required": []string{"ip"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "query_server_os",
				Description: "查询服务器当前操作系统和可用OS列表",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "query_server_os",
					"description": "查询服务器操作系统信息",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"host_id": map[string]interface{}{"type": "integer", "description": "主机ID"},
						},
						"required": []string{"host_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "server_sync",
				Description: "同步服务器状态（从上游拉取最新状态）",
				RiskLevel:   "medium",
				Schema: map[string]interface{}{
					"name":        "server_sync",
					"description": "同步服务器状态",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"host_id": map[string]interface{}{"type": "integer", "description": "主机ID"},
						},
						"required": []string{"host_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "server_boot",
				Description: "开机（需要客户明确确认）",
				RiskLevel:   "high",
				Schema: map[string]interface{}{
					"name":        "server_boot",
					"description": "开机操作",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"host_id": map[string]interface{}{"type": "integer", "description": "主机ID"},
						},
						"required": []string{"host_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "server_shutdown",
				Description: "关机（需要客户明确确认）",
				RiskLevel:   "high",
				Schema: map[string]interface{}{
					"name":        "server_shutdown",
					"description": "关机操作",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"host_id": map[string]interface{}{"type": "integer", "description": "主机ID"},
						},
						"required": []string{"host_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "server_reboot",
				Description: "重启服务器（需要客户明确确认）",
				RiskLevel:   "high",
				Schema: map[string]interface{}{
					"name":        "server_reboot",
					"description": "重启操作",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"host_id": map[string]interface{}{"type": "integer", "description": "主机ID"},
						},
						"required": []string{"host_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "server_reinstall",
				Description: "重装系统（高危，必须客户明确确认并指定OS）",
				RiskLevel:   "critical",
				Schema: map[string]interface{}{
					"name":        "server_reinstall",
					"description": "重装操作系统",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"host_id": map[string]interface{}{"type": "integer", "description": "主机ID"},
							"os":      map[string]interface{}{"type": "string", "description": "目标操作系统"},
						},
						"required": []string{"host_id", "os"},
					},
				},
				Enabled: true,
			},
		},
	}
}

func (r *AIToolRegistry) buildTicketOpsCategory() AIToolCategory {
	return AIToolCategory{
		Key:  "ticket_ops",
		Name: "工单管理",
		Tools: []AITool{
			{
				Name:        "create_ticket",
				Description: "为客户创建新工单",
				RiskLevel:   "medium",
				Schema: map[string]interface{}{
					"name":        "create_ticket",
					"description": "创建新工单",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"client_id": map[string]interface{}{"type": "integer", "description": "客户ID"},
							"dptid":     map[string]interface{}{"type": "integer", "description": "部门ID"},
							"title":     map[string]interface{}{"type": "string", "description": "工单标题"},
							"content":   map[string]interface{}{"type": "string", "description": "工单内容"},
							"host_id":   map[string]interface{}{"type": "integer", "description": "关联主机ID"},
						},
						"required": []string{"client_id", "dptid", "title", "content"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "transfer_ticket_department",
				Description: "将工单分流到其他部门（创建新工单并关闭原工单）",
				RiskLevel:   "medium",
				Schema: map[string]interface{}{
					"name":        "transfer_ticket_department",
					"description": "工单分流到其他部门",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"ticket_id":    map[string]interface{}{"type": "integer", "description": "工单ID"},
							"target_dptid": map[string]interface{}{"type": "integer", "description": "目标部门ID"},
						},
						"required": []string{"ticket_id", "target_dptid"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "add_ticket_note",
				Description: "添加内部备注（仅管理员可见，客户不可见）",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "add_ticket_note",
					"description": "添加内部备注",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"ticket_id": map[string]interface{}{"type": "integer", "description": "工单ID"},
							"note":      map[string]interface{}{"type": "string", "description": "备注内容"},
						},
						"required": []string{"ticket_id", "note"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "ticket_transmit_upstream",
				Description: "透传工单到上游渠道",
				RiskLevel:   "high",
				Schema: map[string]interface{}{
					"name":        "ticket_transmit_upstream",
					"description": "透传工单到上游",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"ticket_id": map[string]interface{}{"type": "integer", "description": "工单ID"},
							"host_id":   map[string]interface{}{"type": "integer", "description": "主机ID"},
							"subject":   map[string]interface{}{"type": "string", "description": "上游工单标题"},
							"content":   map[string]interface{}{"type": "string", "description": "上游工单内容"},
						},
						"required": []string{"ticket_id", "subject", "content"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "sync_upstream_reply",
				Description: "同步上游工单回复（AI脱敏润色后回复客户）",
				RiskLevel:   "medium",
				Schema: map[string]interface{}{
					"name":        "sync_upstream_reply",
					"description": "同步上游回复",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"ticket_id":      map[string]interface{}{"type": "integer", "description": "工单ID"},
							"upstream_reply": map[string]interface{}{"type": "string", "description": "上游回复内容"},
						},
						"required": []string{"ticket_id", "upstream_reply"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "check_upstream_reply",
				Description: "查询上游工单是否有新回复",
				RiskLevel:   "medium",
				Schema: map[string]interface{}{
					"name":        "check_upstream_reply",
					"description": "查询上游工单新回复",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"ticket_id": map[string]interface{}{"type": "integer", "description": "工单ID"},
							"host_id":   map[string]interface{}{"type": "integer", "description": "主机ID"},
						},
						"required": []string{"ticket_id", "host_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "close_upstream_ticket",
				Description: "关闭上游工单",
				RiskLevel:   "high",
				Schema: map[string]interface{}{
					"name":        "close_upstream_ticket",
					"description": "关闭上游工单",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"ticket_id": map[string]interface{}{"type": "integer", "description": "工单ID"},
							"host_id":   map[string]interface{}{"type": "integer", "description": "主机ID"},
						},
						"required": []string{"ticket_id", "host_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "transfer_to_human",
				Description: "将工单转接人工（AI停止主动应答，转为静默监听）",
				RiskLevel:   "medium",
				Schema: map[string]interface{}{
					"name":        "transfer_to_human",
					"description": "转接人工",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"ticket_id": map[string]interface{}{"type": "integer", "description": "工单ID"},
							"reason":    map[string]interface{}{"type": "string", "description": "转接原因"},
						},
						"required": []string{"ticket_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "auto_close_ticket",
				Description: "自动关闭工单（客户诉求已解决时使用）",
				RiskLevel:   "medium",
				Schema: map[string]interface{}{
					"name":        "auto_close_ticket",
					"description": "自动关闭工单",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"ticket_id": map[string]interface{}{"type": "integer", "description": "工单ID"},
						},
						"required": []string{"ticket_id"},
					},
				},
				Enabled: true,
			},
		},
	}
}

func (r *AIToolRegistry) buildFinanceCategory() AIToolCategory {
	return AIToolCategory{
		Key:  "finance",
		Name: "财务退款",
		Tools: []AITool{
			{
				Name:        "refund_transmit_upstream",
				Description: "透传退款工单到上游（仅新开通订单+24小时内）",
				RiskLevel:   "critical",
				Schema: map[string]interface{}{
					"name":        "refund_transmit_upstream",
					"description": "透传退款申请到上游",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"ticket_id": map[string]interface{}{"type": "integer", "description": "工单ID"},
							"order_id":  map[string]interface{}{"type": "integer", "description": "订单ID"},
							"host_id":   map[string]interface{}{"type": "integer", "description": "主机ID"},
							"reason":    map[string]interface{}{"type": "string", "description": "退款原因"},
						},
						"required": []string{"ticket_id", "order_id", "reason"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "check_upstream_refund_result",
				Description: "查询上游退款结果，符合条件时自动执行下游退款",
				RiskLevel:   "critical",
				Schema: map[string]interface{}{
					"name":        "check_upstream_refund_result",
					"description": "查询上游退款结果",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"ticket_id": map[string]interface{}{"type": "integer", "description": "工单ID"},
							"order_id":  map[string]interface{}{"type": "integer", "description": "订单ID"},
						},
						"required": []string{"ticket_id", "order_id"},
					},
				},
				Enabled: true,
			},
		},
	}
}

func (r *AIToolRegistry) buildShoppingCategory() AIToolCategory {
	return AIToolCategory{
		Key:  "shopping",
		Name: "商品导购",
		Tools: []AITool{
			{
				Name:        "search_products",
				Description: "搜索商品（按关键词搜索商品名称和描述）",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "search_products",
					"description": "搜索商品",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"keyword": map[string]interface{}{"type": "string", "description": "搜索关键词（如：香港、CN2、轻量）"},
							"limit":   map[string]interface{}{"type": "integer", "description": "返回数量限制，默认10"},
						},
						"required": []string{"keyword"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "get_product_detail",
				Description: "获取商品详情（价格、配置、描述等）",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "get_product_detail",
					"description": "获取商品详情",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"product_id": map[string]interface{}{"type": "integer", "description": "商品ID"},
						},
						"required": []string{"product_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "list_product_groups",
				Description: "列出商品分组（所有可用的商品分类）",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "list_product_groups",
					"description": "列出商品分组",
					"parameters": map[string]interface{}{
						"type":       "object",
						"properties": map[string]interface{}{},
					},
				},
				Enabled: true,
			},
			{
				Name:        "get_group_products",
				Description: "获取分组下的商品列表",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "get_group_products",
					"description": "获取分组下的商品列表",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"group_id": map[string]interface{}{"type": "integer", "description": "商品分组ID"},
						},
						"required": []string{"group_id"},
					},
				},
				Enabled: true,
			},
			{
				Name:        "compare_products",
				Description: "对比多个商品的配置和价格",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "compare_products",
					"description": "对比多个商品",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"product_ids": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "integer"}, "description": "要对比的商品ID列表"},
						},
						"required": []string{"product_ids"},
					},
				},
				Enabled: true,
			},
		},
	}
}

func (r *AIToolRegistry) buildGeneralCategory() AIToolCategory {
	return AIToolCategory{
		Key:  "general",
		Name: "通用",
		Tools: []AITool{
			{
				Name:        "list_available_tools",
				Description: "查询当前可用的工具列表",
				RiskLevel:   "low",
				Schema: map[string]interface{}{
					"name":        "list_available_tools",
					"description": "查询当前可用工具列表",
					"parameters": map[string]interface{}{
						"type":       "object",
						"properties": map[string]interface{}{},
					},
				},
				Enabled: true,
			},
			{
				Name:        "http_request",
				Description: "发送HTTP请求（仅支持http/https）",
				RiskLevel:   "high",
				Schema: map[string]interface{}{
					"name":        "http_request",
					"description": "发送HTTP请求",
					"parameters": map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"method": map[string]interface{}{"type": "string", "description": "请求方法GET/POST", "enum": []string{"GET", "POST"}},
							"url":    map[string]interface{}{"type": "string", "description": "请求URL"},
							"body":   map[string]interface{}{"type": "string", "description": "POST请求体"},
						},
						"required": []string{"url"},
					},
				},
				Enabled: true,
			},
		},
	}
}

// ToolToJSON 工具序列化为 JSON 字符串
func ToolToJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// ─── 兼容 ai_ticket_core.go 的类型和函数 ───

const MaxToolRounds = 5

// ToolCategory 工具分类（兼容旧代码）
type ToolCategory = AIToolCategory

// ToolDefinition 工具定义（兼容旧代码）
type ToolDefinition = AITool

// AllToolCategories 获取所有工具分类
func AllToolCategories() []ToolCategory {
	r := NewAIToolRegistry()
	return r.GetAllCategories()
}

// AllToolNames 获取所有工具名
func AllToolNames() []string {
	r := NewAIToolRegistry()
	return r.GetAllToolNames()
}

// allToolsMap 返回所有工具名→定义映射
func allToolsMap() map[string]ToolDefinition {
	r := NewAIToolRegistry()
	m := make(map[string]ToolDefinition)
	for name, t := range r.tools {
		m[name] = *t
	}
	return m
}

// GetOpenAITools 获取 OpenAI function calling 格式的工具列表
func GetOpenAITools(enabledTools []string) []map[string]interface{} {
	r := NewAIToolRegistry()
	return r.GetOpenAITools(enabledTools)
}
