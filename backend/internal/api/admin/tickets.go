package admin

import (
	"net/http"
	"strconv"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetTicketList 获取工单列表
// GET /api/admin/tickets
func GetTicketList(c *gin.Context) {
	// 1. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	status := c.Query("status")
	userID := c.Query("user_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// 2. 构建查询
	db := database.GetDB()
	query := db.Model(&model.Ticket{})

	// 关键词搜索（工单号、主题）
	if keyword != "" {
		query = query.Where("ticket_no LIKE ? OR subject LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 用户筛选
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// 3. 获取总数
	var total int64
	query.Count(&total)

	// 4. 分页查询
	var tickets []model.Ticket
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&tickets)

	// 5. 返回统一格式
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

// GetTicket 获取工单详情
// GET /api/admin/tickets/:id
func GetTicket(c *gin.Context) {
	// 1. 获取工单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的工单ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询工单
	db := database.GetDB()
	var ticket model.Ticket
	if err := db.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "工单不存在",
			"data":    nil,
		})
		return
	}

	// 3. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    ticket,
	})
}

// ReplyTicket 回复工单
// POST /api/admin/tickets/:id/reply
func ReplyTicket(c *gin.Context) {
	// 1. 获取工单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的工单ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析请求参数
	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 3. 查询工单
	db := database.GetDB()
	var ticket model.Ticket
	if err := db.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "工单不存在",
			"data":    nil,
		})
		return
	}

	// 4. 检查工单状态
	if ticket.Status == "closed" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "已关闭的工单不能回复",
			"data":    nil,
		})
		return
	}

	// 5. 更新工单状态为pending（等待用户回复）
	if err := db.Model(&ticket).Update("status", "pending").Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "回复工单失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 6. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "回复成功",
		"data":    nil,
	})
}

// CloseTicket 关闭工单
// POST /api/admin/tickets/:id/close
func CloseTicket(c *gin.Context) {
	// 1. 获取工单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的工单ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询工单
	db := database.GetDB()
	var ticket model.Ticket
	if err := db.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "工单不存在",
			"data":    nil,
		})
		return
	}

	// 3. 检查工单状态
	if ticket.Status == "closed" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "工单已关闭",
			"data":    nil,
		})
		return
	}

	// 4. 关闭工单
	if err := db.Model(&ticket).Update("status", "closed").Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "关闭工单失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "关闭成功",
		"data":    nil,
	})
}

// GetTicketDepartments 获取工单部门列表
// GET /api/admin/ticket-departments
func GetTicketDepartments(c *gin.Context) {
	db := database.GetDB()
	var departments []model.TicketDepartment
	db.Where("status = ?", "active").Order("sort_order ASC").Find(&departments)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    departments,
	})
}

// GetTicketStatuses 获取工单状态列表
// GET /api/admin/ticket-statuses
func GetTicketStatuses(c *gin.Context) {
	db := database.GetDB()
	var statuses []model.TicketStatus
	db.Order("sort_order ASC").Find(&statuses)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    statuses,
	})
}

// GetTicketSummary 获取工单统计
// GET /api/admin/tickets/summary
func GetTicketSummary(c *gin.Context) {
	db := database.GetDB()

	var openCount int64
	db.Model(&model.Ticket{}).Where("status = ?", "open").Count(&openCount)

	var pendingCount int64
	db.Model(&model.Ticket{}).Where("status = ?", "pending").Count(&pendingCount)

	var closedCount int64
	db.Model(&model.Ticket{}).Where("status = ?", "closed").Count(&closedCount)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"open":    openCount,
			"pending": pendingCount,
			"closed":  closedCount,
			"total":   openCount + pendingCount + closedCount,
		},
	})
}

// ReopenTicket 重新打开工单
// POST /api/admin/tickets/:id/reopen
func ReopenTicket(c *gin.Context) {
	// 1. 获取工单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的工单ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询工单
	db := database.GetDB()
	var ticket model.Ticket
	if err := db.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "工单不存在",
			"data":    nil,
		})
		return
	}

	// 3. 检查状态
	if ticket.Status != "closed" {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "只有已关闭的工单才能重新打开",
			"data":    nil,
		})
		return
	}

	// 4. 重新打开
	if err := db.Model(&ticket).Update("status", "open").Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "重新打开工单失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "工单已重新打开",
		"data":    nil,
	})
}

// AssignTicket 分配工单
// PUT /api/admin/tickets/:id/assignment
func AssignTicket(c *gin.Context) {
	// 1. 获取工单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的工单ID",
			"data":    nil,
		})
		return
	}

	// 2. 解析请求参数
	var req struct {
		AssignedTo uint `json:"assigned_to" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 3. 查询工单
	db := database.GetDB()
	var ticket model.Ticket
	if err := db.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "工单不存在",
			"data":    nil,
		})
		return
	}

	// 4. 分配工单
	if err := db.Model(&ticket).Update("assigned_to", req.AssignedTo).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    500,
			"message": "分配工单失败: " + err.Error(),
			"data":    nil,
		})
		return
	}

	// 5. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "分配成功",
		"data":    nil,
	})
}

// GetTicketReplies 获取工单回复
// GET /api/admin/tickets/:id/replies
func GetTicketReplies(c *gin.Context) {
	// 1. 获取工单ID
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    400,
			"message": "无效的工单ID",
			"data":    nil,
		})
		return
	}

	// 2. 查询工单回复
	db := database.GetDB()
	var replies []model.TicketReply
	db.Where("ticket_id = ?", id).Order("id ASC").Find(&replies)

	// 3. 返回统一格式
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    replies,
	})
}
