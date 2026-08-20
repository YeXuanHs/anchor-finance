package api

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ZjmfCompatHandler handles zjmf-compatible API requests at POST /api.php.
type ZjmfCompatHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewZjmfCompatHandler creates a new ZjmfCompatHandler.
func NewZjmfCompatHandler(db *gorm.DB, log *logger.Logger) *ZjmfCompatHandler {
	return &ZjmfCompatHandler{db: db, log: log}
}

// Handle processes all zjmf-compatible API requests.
func (h *ZjmfCompatHandler) Handle(c *gin.Context) {
	// Parse form parameters
	if err := c.Request.ParseForm(); err != nil {
		h.respError(c, "invalid request")
		return
	}

	action := c.PostForm("action")
	if action == "" {
		action = c.Query("action")
	}
	if action == "" {
		h.respError(c, "missing action")
		return
	}

	// Collect all params for signature verification
	params := make(map[string]string)
	for k, vs := range c.Request.PostForm {
		if len(vs) > 0 {
			params[k] = vs[0]
		}
	}
	for k, vs := range c.Request.URL.Query() {
		if _, exists := params[k]; !exists && len(vs) > 0 {
			params[k] = vs[0]
		}
	}

	// Verify signature
	apiKey := db.GetSystemSetting("zjmf_api_key")
	if apiKey != "" {
		sign := c.PostForm("sign")
		if sign == "" {
			sign = c.Query("sign")
		}
		expected := h.calcSign(params, apiKey)
		if sign != expected {
			h.respError(c, "invalid sign")
			return
		}
	}

	switch action {
	case "getsysteminfo":
		h.handleGetSystemInfo(c)
	case "getproducts":
		h.handleGetProducts(c)
	case "getclients":
		h.handleGetClients(c)
	case "getorders":
		h.handleGetOrders(c)
	case "getinvoices":
		h.handleGetInvoices(c)
	case "gettickets":
		h.handleGetTickets(c)
	default:
		h.respError(c, "unknown action: "+action)
	}
}

func (h *ZjmfCompatHandler) calcSign(params map[string]string, apiKey string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		if k == "sign" {
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf []byte
	for _, k := range keys {
		buf = append(buf, params[k]...)
	}
	buf = append(buf, apiKey...)
	return fmt.Sprintf("%x", md5.Sum(buf))
}

func (h *ZjmfCompatHandler) respSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"msg":    "",
		"data":   data,
	})
}

func (h *ZjmfCompatHandler) respError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"result": "error",
		"msg":    msg,
		"data":   nil,
	})
}

func (h *ZjmfCompatHandler) handleGetSystemInfo(c *gin.Context) {
	totalUsers := int64(0)
	h.db.Model(&model.User{}).Count(&totalUsers)

	totalProducts := int64(0)
	h.db.Model(&model.Product{}).Count(&totalProducts)

	totalOrders := int64(0)
	h.db.Model(&model.Order{}).Count(&totalOrders)

	totalInvoices := int64(0)
	h.db.Model(&model.Invoice{}).Count(&totalInvoices)

	h.respSuccess(c, gin.H{
		"company_name":  db.GetSystemSetting("company_name"),
		"version":       "1.0.0",
		"total_users":   totalUsers,
		"total_products": totalProducts,
		"total_orders":  totalOrders,
		"total_invoices": totalInvoices,
	})
}

func (h *ZjmfCompatHandler) handleGetProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var products []model.Product
	var total int64

	query := h.db.Model(&model.Product{})
	query.Count(&total)

	if err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&products).Error; err != nil {
		h.respError(c, err.Error())
		return
	}

	// Group products by group
	type productItem struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Price       string `json:"price"`
		Currency    string `json:"currency"`
		Type        string `json:"type"`
		GroupID     uint   `json:"gid"`
	}
	var items []productItem
	for _, p := range products {
		items = append(items, productItem{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       fmt.Sprintf("%.4f", p.Price),
			Currency:    p.Currency,
			Type:        p.Type,
			GroupID:     p.GroupID,
		})
	}

	h.respSuccess(c, gin.H{
		"list":  items,
		"total": total,
	})
}

func (h *ZjmfCompatHandler) handleGetClients(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var users []model.User
	var total int64

	query := h.db.Model(&model.User{})
	query.Count(&total)

	if err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		h.respError(c, err.Error())
		return
	}

	type clientItem struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Status   int16  `json:"status"`
	}
	var items []clientItem
	for _, u := range users {
		items = append(items, clientItem{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Phone:    u.Phone,
			Status:   u.Status,
		})
	}

	h.respSuccess(c, gin.H{
		"list":  items,
		"total": total,
	})
}

func (h *ZjmfCompatHandler) handleGetOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var orders []model.Order
	var total int64

	query := h.db.Model(&model.Order{})
	query.Count(&total)

	if err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&orders).Error; err != nil {
		h.respError(c, err.Error())
		return
	}

	type orderItem struct {
		ID       uint   `json:"id"`
		OrderNo  string `json:"order_no"`
		UserID   uint   `json:"user_id"`
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
		Status   int16  `json:"status"`
	}
	var items []orderItem
	for _, o := range orders {
		items = append(items, orderItem{
			ID:       o.ID,
			OrderNo:  o.OrderNo,
			UserID:   o.UserID,
			Amount:   fmt.Sprintf("%.4f", o.Total),
			Currency: o.Currency,
			Status:   o.Status,
		})
	}

	h.respSuccess(c, gin.H{
		"list":  items,
		"total": total,
	})
}

func (h *ZjmfCompatHandler) handleGetInvoices(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var invoices []model.Invoice
	var total int64

	query := h.db.Model(&model.Invoice{})
	query.Count(&total)

	if err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&invoices).Error; err != nil {
		h.respError(c, err.Error())
		return
	}

	type invoiceItem struct {
		ID        uint   `json:"id"`
		InvoiceNo string `json:"invoice_no"`
		UserID    uint   `json:"user_id"`
		Total     string `json:"total"`
		Currency  string `json:"currency"`
		Status    int16  `json:"status"`
	}
	var items []invoiceItem
	for _, inv := range invoices {
		items = append(items, invoiceItem{
			ID:        inv.ID,
			InvoiceNo: inv.InvoiceNo,
			UserID:    inv.UserID,
			Total:     fmt.Sprintf("%.4f", inv.Total),
			Currency:  inv.Currency,
			Status:    inv.Status,
		})
	}

	h.respSuccess(c, gin.H{
		"list":  items,
		"total": total,
	})
}

func (h *ZjmfCompatHandler) handleGetTickets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var tickets []model.Ticket
	var total int64

	query := h.db.Model(&model.Ticket{})
	query.Count(&total)

	if err := query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&tickets).Error; err != nil {
		h.respError(c, err.Error())
		return
	}

	type ticketItem struct {
		ID       uint   `json:"id"`
		TicketNo string `json:"ticket_no"`
		UserID   uint   `json:"user_id"`
		Subject  string `json:"subject"`
		Status   int16  `json:"status"`
		Priority int16  `json:"priority"`
	}
	var items []ticketItem
	for _, t := range tickets {
		items = append(items, ticketItem{
			ID:       t.ID,
			TicketNo: t.TicketNo,
			UserID:   t.UserID,
			Subject:  t.Subject,
			Status:   t.Status,
			Priority: t.Priority,
		})
	}

	h.respSuccess(c, gin.H{
		"list":  items,
		"total": total,
	})
}
