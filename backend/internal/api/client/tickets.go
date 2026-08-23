package client

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/pluginengine"
	"github.com/gin-gonic/gin"
)

// GetUserTickets 获取用户工单列表
// GET /api/client/tickets
func GetUserTickets(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	query := db.Model(&model.Ticket{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var tickets []model.Ticket
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&tickets)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      tickets,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetUserTicket 获取用户工单详情
// GET /api/client/tickets/:id
func GetUserTicket(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的工单ID",
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var ticket model.Ticket
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "工单不存在",
			"data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    ticket,
	})
}

// CreateUserTicket 创建用户工单
// POST /api/client/tickets
func CreateUserTicket(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Subject    string `json:"subject" binding:"required"`
		Department string `json:"department"`
		Priority   string `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	// 生成工单号
	ticketNo := fmt.Sprintf("TK%s%06d", time.Now().Format("20060102"), time.Now().UnixNano()%1000000)

	if req.Priority == "" {
		req.Priority = "normal"
	}

	db := database.GetDB()
	ticket := model.Ticket{
		UserID:     userID.(uint),
		TicketNo:   ticketNo,
		Subject:    req.Subject,
		Status:     "open",
		Priority:   req.Priority,
		Department: req.Department,
	}

	if err := db.Create(&ticket).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "创建工单失败",
			"data": nil,
		})
		return
	}

	// 触发Hook: ticket_open
	pluginengine.TriggerHook("ticket_open", map[string]interface{}{
		"ticket_id": ticket.ID, "user_id": ticket.UserID, "subject": ticket.Subject,
	})

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "创建成功",
		"data": gin.H{
			"id":        ticket.ID,
			"ticket_no": ticket.TicketNo,
		},
	})
}

// ReplyUserTicket 用户回复工单
// POST /api/client/tickets/:id/reply
func ReplyUserTicket(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的工单ID",
			"data": nil,
		})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误",
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var ticket model.Ticket
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "工单不存在",
			"data": nil,
		})
		return
	}

	if ticket.Status == "closed" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "已关闭的工单不能回复",
			"data": nil,
		})
		return
	}

	// 保存回复内容到数据库
	reply := model.TicketReply{
		TicketID: uint(id),
		UserID:   userID.(uint),
		Content:  req.Content,
		IsAdmin:  false,
	}
	db.Create(&reply)

	// 更新工单状态
	db.Model(&ticket).Updates(map[string]interface{}{
		"status":     "open",
		"updated_at": time.Now(),
	})

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "回复成功",
		"data":    gin.H{"reply_id": reply.ID},
	})
}

// GetUserTicketReplies 获取用户工单回复
// GET /api/client/tickets/:id/replies
func GetUserTicketReplies(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的工单ID",
			"data": nil,
		})
		return
	}

	// 验证工单属于该用户
	db := database.GetDB()
	var ticket model.Ticket
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "工单不存在",
			"data": nil,
		})
		return
	}

	// 获取回复
	var replies []model.TicketReply
	db.Where("ticket_id = ?", id).Order("id ASC").Find(&replies)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    replies,
	})
}

// CloseUserTicket 用户关闭工单
// POST /api/client/tickets/:id/close
func CloseUserTicket(c *gin.Context) {
	userID, _ := c.Get("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的工单ID",
			"data": nil,
		})
		return
	}

	db := database.GetDB()
	var ticket model.Ticket
	if err := db.Where("id = ? AND user_id = ?", id, userID).First(&ticket).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "工单不存在",
			"data": nil,
		})
		return
	}

	if ticket.Status == "closed" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "工单已关闭",
			"data": nil,
		})
		return
	}

	db.Model(&ticket).Update("status", "closed")

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "关闭成功",
		"data":    nil,
	})
}

// TicketUpstreamReply 上游工单回复回调（供应商回调，需要API Key验证）
// POST /api/client/tickets/upstream/replies
func TicketUpstreamReply(c *gin.Context) {
	// 验证上游回调API Key（从Header获取）
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		apiKey = c.Query("api_key")
	}
	if apiKey == "" {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "缺少API Key", "data": nil})
		return
	}

	// 验证API Key是否有效（检查是否有匹配的供应商）
	db := database.GetDB()
	var supplier model.Supplier
	if err := db.Where("api_key = ? AND status = ?", apiKey, "active").First(&supplier).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 401, "message": "无效的API Key", "data": nil})
		return
	}

	var req struct {
		TicketID uint   `json:"ticket_id" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Author   string `json:"author"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "参数错误", "data": nil})
		return
	}

	var ticket model.Ticket
	if err := db.First(&ticket, req.TicketID).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "工单不存在", "data": nil})
		return
	}

	// 创建回复
	reply := model.TicketReply{
		TicketID: ticket.ID,
		Content:  req.Content,
		IsAdmin:  false, // 上游回复标记为非管理员
	}
	if err := db.Create(&reply).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "保存回复失败", "data": nil})
		return
	}

	// 更新工单状态
	db.Model(&ticket).Update("status", "answered")

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "回复成功", "data": gin.H{"reply_id": reply.ID}})
}
