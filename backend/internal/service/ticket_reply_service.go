package service

import (
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// TicketReplyService 工单回复服务
type TicketReplyService struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewTicketReplyService 创建工单回复服务
func NewTicketReplyService(db *gorm.DB, log *logger.Logger) *TicketReplyService {
	return &TicketReplyService{db: db, log: log}
}

// PostAdminReply 以管理员身份回复工单
// 1. Markdown 转 HTML
// 2. 追加 AI 免责声明
// 3. 写入 ticket_replies 表（UserID=nil, AdminID=adminID）
// 4. 更新工单状态为"已回复"（status=1）
// 5. 更新 FirstReplyAt（如果是第一次回复）
// 6. 设置 IsRead=false（客户未读）
func (s *TicketReplyService) PostAdminReply(ticketID uint, content string, adminID uint) error {
	// 验证工单存在
	var ticket model.Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		return fmt.Errorf("工单不存在: %w", err)
	}

	// Markdown 转 HTML
	htmlContent := MarkdownToHTML(content)

	// 追加 AI 免责声明
	htmlContent = AppendAIDisclaimer(htmlContent)

	// 构建回复数据
	now := time.Now()
	reply := model.TicketReply{
		TicketID: ticketID,
		AdminID:  &adminID,
		Content:  htmlContent,
		IsRead:   false,
	}

	// 开始事务
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 插入回复
		if err := tx.Create(&reply).Error; err != nil {
			return fmt.Errorf("写入回复失败: %w", err)
		}

		// 更新工单状态
		updates := map[string]interface{}{
			"status": 1, // 已回复
		}

		// 设置首次回复时间
		if ticket.FirstReplyAt == nil {
			updates["first_reply_at"] = now
		}

		if err := tx.Model(&model.Ticket{}).Where("id = ?", ticketID).Updates(updates).Error; err != nil {
			return fmt.Errorf("更新工单状态失败: %w", err)
		}

		// 记录操作日志
		s.log.Infof("AI回复工单#User ID:%d - Ticket ID:%d - %s", ticket.UserID, ticketID, ticket.Subject)

		// 将之前的回复标记为已读（管理员已处理）
		tx.Model(&model.TicketReply{}).
			Where("ticket_id = ? AND is_read = ? AND admin_id IS NULL", ticketID, false).
			Update("is_read", true)

		return nil
	})
}

// PostSystemReply 以系统身份回复工单（内部备注）
func (s *TicketReplyService) PostSystemReply(ticketID uint, content string) error {
	var ticket model.Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		return fmt.Errorf("工单不存在: %w", err)
	}

	reply := model.TicketReply{
		TicketID:   ticketID,
		Content:    content,
		IsInternal: true,
		IsRead:     true,
	}

	return s.db.Create(&reply).Error
}

// getAdminDisplayName 获取管理员显示名称
func (s *TicketReplyService) getAdminDisplayName(adminID uint) string {
	var admin model.Admin
	if err := s.db.Select("username, real_name").First(&admin, adminID).Error; err != nil {
		return "admin"
	}
	if admin.RealName != "" {
		return admin.RealName
	}
	if admin.Username != "" {
		return admin.Username
	}
	return "admin"
}

// GetTicketInfo 获取工单基本信息（用于通知）
func (s *TicketReplyService) GetTicketInfo(ticketID uint) map[string]interface{} {
	var ticket model.Ticket
	if err := s.db.First(&ticket, ticketID).Error; err != nil {
		return nil
	}

	var deptName string
	s.db.Model(&model.Department{}).Where("id = ?", ticket.DepartmentID).Pluck("name", &deptName)

	priorityMap := map[int16]string{0: "低", 1: "中", 2: "高", 3: "紧急"}

	return map[string]interface{}{
		"ticket_id":  ticket.ID,
		"tid":        ticket.TicketNo,
		"title":      ticket.Subject,
		"department": deptName,
		"priority":   priorityMap[ticket.Priority],
		"uid":        ticket.UserID,
		"status":     ticket.Status,
	}
}
