package handler

import (
	"strconv"
	"time"

	"anchorfinance/internal/security"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// FrontendUserHandler 前台用户信息处理器
// 专门处理前台用户信息展示，敏感数据自动脱敏
type FrontendUserHandler struct {
	userSvc *service.UserService
	log     *logger.Logger
}

// NewFrontendUserHandler 创建前台用户处理器
func NewFrontendUserHandler(userSvc *service.UserService, log *logger.Logger) *FrontendUserHandler {
	return &FrontendUserHandler{userSvc: userSvc, log: log}
}

// GetPublicProfile 获取用户公开信息（脱敏）
// 用于前台展示其他用户信息，如评论区、工单等场景
func (h *FrontendUserHandler) GetPublicProfile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	user, err := h.userSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	// 返回脱敏后的用户信息
	response.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"phone":      security.MaskPhone(user.Phone),
		"email":      security.MaskEmail(user.Email),
		"status":     user.Status,
		"created_at": user.CreatedAt,
	})
}

// GetMyLoginLogs 获取当前用户的登录日志（IP脱敏）
func (h *FrontendUserHandler) GetMyLoginLogs(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	logs, total, err := h.userSvc.GetLoginLogs(userID, page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// 脱敏 IP 后返回
	maskedLogs := make([]gin.H, len(logs))
	for i, l := range logs {
		maskedLogs[i] = gin.H{
			"id":         l.ID,
			"ip":         security.MaskIP(l.IP),
			"user_agent": l.UserAgent,
			"status":     l.Status,
			"reason":     l.Reason,
			"created_at": l.CreatedAt.Format(time.RFC3339),
		}
	}

	response.Success(c, gin.H{
		"list":      maskedLogs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}
