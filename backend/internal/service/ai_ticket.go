package service

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
)

// AITicketService AI工单服务（参考mianyu_ai_ticket插件架构）
type AITicketService struct {
	aiSvc *AIService
}

// NewAITicketService 创建AI工单服务
func NewAITicketService() *AITicketService {
	return &AITicketService{
		aiSvc: NewAIService(),
	}
}

// IsEnabled AI工单是否启用
func (s *AITicketService) IsEnabled() bool {
	db := database.GetDB()
	var setting model.Setting
	if err := db.Where("`key` = ?", "ai_ticket_enabled").First(&setting).Error; err != nil {
		return false
	}
	return setting.Value == "1"
}

// GetConfig 获取AI工单配置
func (s *AITicketService) GetConfig() map[string]interface{} {
	db := database.GetDB()
	var settings []model.Setting
	db.Where("`group` = ?", "ai_ticket").Find(&settings)

	config := make(map[string]interface{})
	for _, setting := range settings {
		config[setting.Key] = setting.Value
	}

	// 填充默认值
	defaults := map[string]interface{}{
		"ai_ticket_enabled":    "0",
		"ai_ticket_dptids":     "",
		"ai_ticket_reply_admin": "1",
		"ai_ticket_system_prompt": "",
		"ai_ticket_queue_mode": "1",
		"ai_ticket_queue_batch": "20",
		"ai_ticket_transfer_keywords": "",
	}
	for k, v := range defaults {
		if _, ok := config[k]; !ok {
			config[k] = v
		}
	}
	return config
}

// isDepartmentEnabled 检查部门是否启用AI（参考插件isDepartmentEnabled）
func (s *AITicketService) isDepartmentEnabled(config map[string]interface{}, deptID uint) bool {
	raw, _ := config["ai_ticket_dptids"].(string)
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return true // 留空=全部部门
	}
	for _, id := range strings.Split(raw, ",") {
		id = strings.TrimSpace(id)
		if id == fmt.Sprintf("%d", deptID) {
			return true
		}
	}
	return false
}

// isTicketAiMode 检查工单是否处于AI模式（参考插件isTicketAiMode）
func (s *AITicketService) isTicketAiMode(ticketID uint) bool {
	db := database.GetDB()
	var mode model.AITicketMode
	if err := db.Where("ticket_id = ?", ticketID).First(&mode).Error; err != nil {
		return true // 默认AI模式（没有记录=AI模式）
	}
	return mode.Mode != "human"
}

// SetTicketMode 设置工单AI/人工模式
func (s *AITicketService) SetTicketMode(ticketID uint, mode string) error {
	db := database.GetDB()
	var existing model.AITicketMode
	if err := db.Where("ticket_id = ?", ticketID).First(&existing).Error; err != nil {
		// 不存在则创建
		return db.Create(&model.AITicketMode{
			TicketID:  ticketID,
			Mode:      mode,
			UpdatedAt: time.Now(),
		}).Error
	}
	return db.Model(&existing).Updates(map[string]interface{}{
		"mode":       mode,
		"updated_at": time.Now(),
	}).Error
}

// HandleTicketEvent 处理工单事件（参考插件handleTicketEvent）
func (s *AITicketService) HandleTicketEvent(eventType string, ticketID, userID, deptID uint) error {
	if !s.IsEnabled() {
		return nil
	}

	config := s.GetConfig()

	// 部门过滤
	if !s.isDepartmentEnabled(config, deptID) {
		return nil
	}

	// 检查工单是否处于人工模式
	if !s.isTicketAiMode(ticketID) {
		// 人工模式下入队监控任务
		return s.enqueueJob(ticketID, 0, "monitoring", userID, deptID)
	}

	// AI模式下入队正常处理任务
	return s.enqueueJob(ticketID, 0, eventType, userID, deptID)
}

// enqueueJob 入队AI任务
func (s *AITicketService) enqueueJob(ticketID, replyID uint, eventType string, userID, deptID uint) error {
	db := database.GetDB()

	// 幂等检查：同一工单+回复+事件，不重复入队
	var existing model.AITicketQueue
	if err := db.Where("ticket_id = ? AND reply_id = ? AND event_type = ? AND status IN ?",
		ticketID, replyID, eventType, []string{"pending", "processing"}).First(&existing).Error; err == nil {
		return nil // 已有任务
	}

	job := model.AITicketQueue{
		TicketID:  ticketID,
		ReplyID:   replyID,
		EventType: eventType,
		UserID:    userID,
		DeptID:    deptID,
		Status:    "pending",
	}
	return db.Create(&job).Error
}

// ProcessNext 处理下一个待处理任务（参考QueueService.processNext）
func (s *AITicketService) ProcessNext() (bool, error) {
	db := database.GetDB()

	// 恢复超时任务
	s.recoverStuckJobs()

	// 获取下一个待处理任务
	var job model.AITicketQueue
	if err := db.Where("status = ?", "pending").Order("created_at ASC, id ASC").First(&job).Error; err != nil {
		return false, nil // 无待处理任务
	}

	// 锁定任务
	db.Model(&job).Updates(map[string]interface{}{
		"status":     "processing",
		"attempts":   job.Attempts + 1,
		"updated_at": time.Now(),
	})

	// 执行任务
	err := s.runJob(job)
	if err != nil {
		s.markFailed(job.ID, err.Error())
		return false, err
	}

	s.markDone(job.ID)
	return true, nil
}

// runJob 执行单个AI任务（参考QueueService.runJob）
func (s *AITicketService) runJob(job model.AITicketQueue) error {
	db := database.GetDB()

	// 获取工单信息
	var ticket model.Ticket
	if err := db.First(&ticket, job.TicketID).Error; err != nil {
		return fmt.Errorf("工单不存在: %d", job.TicketID)
	}

	// 获取最新回复内容
	var lastReply model.TicketReply
	ticketContent := ticket.Subject
	db.Where("ticket_id = ?", ticket.ID).Order("id DESC").First(&lastReply)
	if lastReply.Content != "" {
		ticketContent = lastReply.Content
	}

	// 获取客户信息
	var user model.User
	db.First(&user, ticket.UserID)

	// 获取知识库上下文
	knowledgeContext := s.getKnowledgeContext(ticketContent)

	// 构建系统提示词
	config := s.GetConfig()
	systemPrompt, _ := config["ai_ticket_system_prompt"].(string)
	if systemPrompt == "" {
		systemPrompt = "你是一个专业的IDC客服AI。根据客户的工单内容，生成专业、友好的回复。"
	}
	if knowledgeContext != "" {
		systemPrompt += "\n\n参考知识库：\n" + knowledgeContext
	}

	// 检查转人工关键词
	transferKeywords, _ := config["ai_ticket_transfer_keywords"].(string)
	if transferKeywords != "" {
		for _, kw := range strings.Split(transferKeywords, " ") {
			if strings.Contains(ticketContent, strings.TrimSpace(kw)) {
				// 命中转人工关键词，切到人工模式
				s.SetTicketMode(job.TicketID, "human")
				s.logProcess(job.TicketID, "keyword_transfer", "transfer", 0, "命中转人工关键词: "+kw, "", "skipped")
				return nil
			}
		}
	}

	// 调用AI生成回复
	reply, err := s.aiSvc.TicketAutoReply(ticket.Subject, ticketContent, map[string]interface{}{
		"customer_info": map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		},
		"system_prompt": systemPrompt,
	})
	if err != nil {
		s.logProcess(job.TicketID, "ai_call", "error", 0, err.Error(), "", "failed")
		return err
	}

	// 检查是否保持沉默（监控模式）
	if job.EventType == "monitoring" && strings.TrimSpace(reply) == "[STAY_SILENT]" {
		s.logProcess(job.TicketID, "monitoring", "stay_silent", 0, "AI判断不接管", "", "skipped")
		return nil
	}

	// 发送AI回复（以管理员身份，MD 10.3：可配置reply_admin）
	replyAdminID := uint(1) // 默认使用ID=1的管理员
	if idStr, ok := config["ai_ticket_reply_admin"].(string); ok && idStr != "" {
		if parsedID, err := strconv.ParseUint(idStr, 10, 32); err == nil && parsedID > 0 {
			replyAdminID = uint(parsedID)
		}
	}

	// 获取AI模型名称用于日志
	modelName := ""
	if aiConfig := s.aiSvc.GetConfig(); aiConfig != nil {
		if m, ok := aiConfig["model"].(string); ok {
			modelName = m
		}
	}

	// 创建工单回复
	ticketReply := model.TicketReply{
		TicketID: ticket.ID,
		UserID:   replyAdminID,
		Content:  reply,
		IsAdmin:  true,
	}
	if err := db.Create(&ticketReply).Error; err != nil {
		return fmt.Errorf("保存回复失败: %w", err)
	}

	// 更新工单状态
	db.Model(&ticket).Updates(map[string]interface{}{
		"status":     "answered",
		"updated_at": time.Now(),
	})

	s.logProcess(job.TicketID, "ai_reply", "success", 1.0, "AI回复成功", modelName, "success")

	// AI自动关单（MD 9.3：客户诉求解决后自动关闭，默认开启）
	var autoCloseSetting string
	database.GetDB().Model(&model.Setting{}).Where("`key` = ?", "ai_ticket_auto_close").Pluck("value", &autoCloseSetting)
	if autoCloseSetting != "0" && autoCloseSetting != "false" {
		// AI回复后自动关闭工单
		database.GetDB().Model(&model.Ticket{}).Where("id = ?", job.TicketID).Update("status", "closed")
		log.Printf("[AITicket] Auto-closed ticket %d after AI reply", job.TicketID)
	}

	return nil
}

// getKnowledgeContext 从知识库获取相关上下文
func (s *AITicketService) getKnowledgeContext(content string) string {
	db := database.GetDB()
	var knowledges []model.AITicketKnowledge
	db.Where("status = 1").Order("sort ASC").Find(&knowledges)

	var matched []string
	for _, k := range knowledges {
		// 简单关键词匹配
		if k.Keywords != "" {
			for _, kw := range strings.Split(k.Keywords, ",") {
				kw = strings.TrimSpace(kw)
				if kw != "" && strings.Contains(content, kw) {
					matched = append(matched, fmt.Sprintf("Q: %s\nA: %s", k.Question, k.Answer))
					break
				}
			}
		}
	}

	if len(matched) == 0 {
		return ""
	}
	return strings.Join(matched, "\n\n")
}

// recoverStuckJobs 恢复超时任务（参考QueueService.recoverStuckJobs）
func (s *AITicketService) recoverStuckJobs() {
	db := database.GetDB()
	threshold := time.Now().Add(-90 * time.Second)
	db.Model(&model.AITicketQueue{}).
		Where("status = ? AND updated_at < ?", "processing", threshold).
		Updates(map[string]interface{}{
			"status":     "pending",
			"error_msg":  "处理超时已自动恢复",
			"updated_at": time.Now(),
		})
}

// markDone 标记任务完成
func (s *AITicketService) markDone(jobID uint) {
	db := database.GetDB()
	db.Model(&model.AITicketQueue{}).Where("id = ?", jobID).Updates(map[string]interface{}{
		"status":     "done",
		"error_msg":  "",
		"updated_at": time.Now(),
	})
}

// markFailed 标记任务失败（最多重试3次）
func (s *AITicketService) markFailed(jobID uint, message string) {
	db := database.GetDB()
	var job model.AITicketQueue
	db.First(&job, jobID)

	status := "pending"
	if job.Attempts >= 3 {
		status = "failed"
	}

	db.Model(&model.AITicketQueue{}).Where("id = ?", jobID).Updates(map[string]interface{}{
		"status":     status,
		"error_msg":  message,
		"updated_at": time.Now(),
	})
}

// logProcess 记录处理日志
func (s *AITicketService) logProcess(ticketID uint, event, decision string, confidence float64, message, modelUsed, status string) {
	db := database.GetDB()
	db.Create(&model.AITicketProcessLog{
		TicketID:   ticketID,
		Event:      event,
		Decision:   decision,
		Confidence: confidence,
		Status:     status,
		Message:    message,
		ModelUsed:  modelUsed,
	})
}

// ProcessQueue 处理队列（批量，参考ImmediateProcessor.run）
func (s *AITicketService) ProcessQueue(maxJobs int) (int, error) {
	processed := 0
	for processed < maxJobs {
		ok, err := s.ProcessNext()
		if err != nil {
			return processed, err
		}
		if !ok {
			break
		}
		processed++
	}
	return processed, nil
}

// GetQueueStats 获取队列统计
func (s *AITicketService) GetQueueStats() map[string]interface{} {
	db := database.GetDB()

	var pending, processing, done, failed int64
	db.Model(&model.AITicketQueue{}).Where("status = ?", "pending").Count(&pending)
	db.Model(&model.AITicketQueue{}).Where("status = ?", "processing").Count(&processing)
	db.Model(&model.AITicketQueue{}).Where("status = ?", "done").Count(&done)
	db.Model(&model.AITicketQueue{}).Where("status = ?", "failed").Count(&failed)

	return map[string]interface{}{
		"pending":    pending,
		"processing": processing,
		"done":       done,
		"failed":     failed,
		"total":      pending + processing + done + failed,
	}
}
