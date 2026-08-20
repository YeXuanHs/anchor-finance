package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SendMessageBatchHandler struct {
	svc *service.SendMessageBatchService
	log *logger.Logger
	db  *gorm.DB
}

func NewSendMessageBatchHandler(svc *service.SendMessageBatchService, log *logger.Logger, db *gorm.DB) *SendMessageBatchHandler {
	return &SendMessageBatchHandler{svc: svc, log: log, db: db}
}

// GetSearchParams returns search parameters for batch messaging.
func (h *SendMessageBatchHandler) GetSearchParams(c *gin.Context) {
	params, err := h.svc.GetSearchParams()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, params)
}

// GetBatches returns paginated batch send tasks.
func (h *SendMessageBatchHandler) GetBatches(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	sendMethod := c.Query("send_method")

	var status *int8
	if s := c.Query("status"); s != "" {
		v, _ := strconv.ParseInt(s, 10, 8)
		st := int8(v)
		status = &st
	}

	batches, total, err := h.svc.GetBatches(page, pageSize, sendMethod, status)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, batches, total, page, pageSize)
}

// SendBatch sends messages in batch.
func (h *SendMessageBatchHandler) SendBatch(c *gin.Context) {
	var req struct {
		SendMethod string `json:"send_method" binding:"required"`
		Subject    string `json:"subject"`
		Content    string `json:"content" binding:"required"`
		UserIDs    []uint `json:"user_ids"`
		GroupIDs   []uint `json:"group_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	batch := &model.SendMessageBatch{
		SendMethod: req.SendMethod,
		Subject:    req.Subject,
		Content:    req.Content,
		Status:     0,
		CreatedBy:  c.GetUint("user_id"),
	}

	if err := h.svc.SendBatch(batch); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, batch)
}

// GetProgress returns the progress of a batch send operation.
func (h *SendMessageBatchHandler) GetProgress(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid batch id")
		return
	}

	batch, err := h.svc.GetProgress(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, batch)
}

// GetRecords returns batch send records.
func (h *SendMessageBatchHandler) GetRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var batchID *uint
	if bid := c.Query("batch_id"); bid != "" {
		v, _ := strconv.ParseUint(bid, 10, 64)
		id := uint(v)
		batchID = &id
	}

	records, total, err := h.svc.GetRecords(page, pageSize, batchID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, records, total, page, pageSize)
}

// GetSendMethod 获取发送方式配置（邮件、短信、站内信的详细配置）
// GET /admin/send-message-batch/send-method
func (h *SendMessageBatchHandler) GetSendMethod(c *gin.Context) {
	sendMethods := []gin.H{
		{"name": "email", "value": "邮件"},
		{"name": "mobile", "value": "手机"},
		{"name": "system", "value": "站内信"},
	}

	response.Success(c, gin.H{
		"send_method": sendMethods,
	})
}

// GetSMSTemplateList 短信模板列表
// GET /admin/send-message-batch/sms-templates
func (h *SendMessageBatchHandler) GetSMSTemplateList(c *gin.Context) {
	type SMSTemplate struct {
		ID         uint   `json:"id"`
		Title      string `json:"title"`
		TemplateID string `json:"template_id"`
		RangeType  int    `json:"range_type"`
		Content    string `json:"content"`
		Status     int    `json:"status"`
		Operator   string `json:"sms_operator"`
	}

	var templates []SMSTemplate
	h.db.Table("message_templates").
		Select("id, title, template_id, range_type, content, status, sms_operator").
		Where("status = ? AND range_type = ?", 2, 2).
		Scan(&templates)

	response.Success(c, gin.H{
		"templates": templates,
	})
}

// GetEmailTemplateList 邮件模板列表
// GET /admin/send-message-batch/email-templates
func (h *SendMessageBatchHandler) GetEmailTemplateList(c *gin.Context) {
	keyword := c.Query("keyword")

	type EmailTmpl struct {
		ID       uint   `json:"id"`
		Type     string `json:"type"`
		Name     string `json:"name"`
		Status   int16  `json:"status"`
	}

	query := h.db.Table("email_templates").
		Select("id, type, name, status").
		Where("language = '' OR language IS NULL")

	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	var templates []EmailTmpl
	query.Scan(&templates)

	response.Success(c, gin.H{
		"email_list": templates,
	})
}

// GetEmailTemplateParams 获取邮件模板参数
// GET /admin/send-message-batch/email-template-params
func (h *SendMessageBatchHandler) GetEmailTemplateParams(c *gin.Context) {
	sendType := c.DefaultQuery("send_type", "clients")

	type TemplateArg struct {
		Label string `json:"label"`
		Name  string `json:"name"`
		List  []gin.H `json:"list"`
	}

	baseArgs := []TemplateArg{
		{
			Label: "args_base",
			Name:  "其他",
			List: []gin.H{
				{"label": "{SYSTEM_COMPANYNAME}", "name": "公司名称"},
				{"label": "{COMPANY_DOMAIN}", "name": "公司域名"},
				{"label": "{TEMPLATE_DATE}", "name": "模板日期"},
				{"label": "{TEMPLATE_TIME}", "name": "模板时间"},
				{"label": "{CODE}", "name": "验证码"},
				{"label": "{SEND_TIME}", "name": "发送时间"},
				{"label": "{SYSTEM_URL}", "name": "系统URL"},
				{"label": "{SYSTEM_WEB_URL}", "name": "网站URL"},
			},
		},
		{
			Label: "args_clients",
			Name:  "客户相关",
			List: []gin.H{
				{"label": "{CLIENT_ID}", "name": "客户ID"},
				{"label": "{USERNAME}", "name": "用户名"},
				{"label": "{ACCOUNT_EMAIL}", "name": "邮箱"},
				{"label": "{CLIENT_SIGNUP_DATE}", "name": "注册时间"},
				{"label": "{CLIENT_STATUS}", "name": "客户状态"},
				{"label": "{CLIENT_GROUP_NAME}", "name": "客户组"},
				{"label": "{CLIENT_PHONENUMBER}", "name": "手机号"},
			},
		},
	}

	// 如果发送类型包含产品，追加产品变量
	if sendType == "clients_and_host" {
		baseArgs = append(baseArgs, TemplateArg{
			Label: "args_product",
			Name:  "产品/服务相关",
			List: []gin.H{
				{"label": "{PRODUCT_NAME}", "name": "产品名称"},
				{"label": "{HOSTNAME}", "name": "主机名"},
				{"label": "{PRODUCT_MAINIP}", "name": "主IP"},
				{"label": "{PRODUCT_FIRST_TIME}", "name": "开通时间"},
				{"label": "{PRODUCT_END_TIME}", "name": "到期时间"},
				{"label": "{PRODUCT_BINLLY_CYCLE}", "name": "计费周期"},
				{"label": "{PRODUCT_USER}", "name": "产品用户名"},
				{"label": "{ORDER_ID}", "name": "订单ID"},
				{"label": "{ORDER_TOTAL_FEE}", "name": "订单金额"},
				{"label": "{INVOICE_PAID_TIME}", "name": "支付时间"},
			},
		})
	}

	response.Success(c, gin.H{
		"data": baseArgs,
	})
}

// SearchBatches 搜索发送列表（支持多种筛选条件）
// POST /admin/send-message-batch/search
func (h *SendMessageBatchHandler) SearchBatches(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var req struct {
		ClientIDs   []uint `json:"client_ids"`
		ClientStatus []int `json:"client_status"`
		SaleIDs     []uint `json:"sale_ids"`
		ProductIDs  []uint `json:"product_ids"`
		Language    []string `json:"language"`
		Country     []string `json:"country"`
		SendType    string `json:"send_type"` // clients / clients_and_host
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.SendType == "" {
		req.SendType = "clients"
	}

	type ClientInfo struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Phone    string `json:"mobile"`
		Email    string `json:"email"`
	}

	query := h.db.Table("users").
		Select("id, username, phone as mobile, email").
		Where("status = 1")

	if len(req.ClientIDs) > 0 {
		query = query.Where("id IN ?", req.ClientIDs)
	}
	if len(req.ClientStatus) > 0 {
		query = query.Where("status IN ?", req.ClientStatus)
	}
	if len(req.Language) > 0 {
		query = query.Where("language IN ?", req.Language)
	}

	var clients []ClientInfo
	query.Scan(&clients)

	// 如果是 clients_and_host 类型，还需要关联产品
	type ClientWithHost struct {
		ClientInfo
		HostID     uint   `json:"host_id"`
		HostDomain string `json:"host_domain"`
		ProductName string `json:"product_name"`
		ProductID   uint   `json:"productid"`
	}

	var results []ClientWithHost
	if req.SendType == "clients_and_host" {
		for _, client := range clients {
			var hosts []ClientWithHost
			hostQuery := h.db.Table("hosts h").
				Select("h.id as host_id, h.owner_id as client_id, h.hostname as host_domain, p.name as product_name, h.product_id as productid").
				Joins("LEFT JOIN products p ON p.id = h.product_id").
				Where("h.owner_id = ?", client.ID)

			if len(req.ProductIDs) > 0 {
				hostQuery = hostQuery.Where("h.product_id IN ?", req.ProductIDs)
			}
			hostQuery.Scan(&hosts)

			for _, host := range hosts {
				host.ClientInfo = client
				results = append(results, host)
			}
		}
		if len(results) == 0 {
			response.BadRequest(c, "没有满足条件的发送对象")
			return
		}
	} else {
		for _, client := range clients {
			results = append(results, ClientWithHost{
				ClientInfo: client,
			})
		}
	}

	// 分页
	total := int64(len(results))
	offset := (page - 1) * pageSize
	if offset >= len(results) {
		results = []ClientWithHost{}
	} else {
		end := offset + pageSize
		if end > len(results) {
			end = len(results)
		}
		results = results[offset:end]
	}

	response.SuccessPage(c, results, total, page, pageSize)
}
