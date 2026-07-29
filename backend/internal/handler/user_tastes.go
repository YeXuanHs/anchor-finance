package handler

import (
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserTastesHandler struct {
	svc *service.UserTasteService
	log *logger.Logger
}

func NewUserTastesHandler(svc *service.UserTasteService, log *logger.Logger) *UserTastesHandler {
	return &UserTastesHandler{svc: svc, log: log}
}

// GetUserTastes returns user tastes for the current admin.
func (h *UserTastesHandler) GetUserTastes(c *gin.Context) {
	userID := c.GetUint("user_id")

	taste, err := h.svc.GetByUserID(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, taste)
}

// UpdateUserTastes updates user tastes for the current admin.
func (h *UserTastesHandler) UpdateUserTastes(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		Theme    string      `json:"theme"`
		Language string      `json:"language"`
		Layout   string      `json:"layout"`
		Settings interface{} `json:"settings"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Theme != "" {
		updates["theme"] = req.Theme
	}
	if req.Language != "" {
		updates["language"] = req.Language
	}
	if req.Layout != "" {
		updates["layout"] = req.Layout
	}
	if req.Settings != nil {
		updates["settings"] = req.Settings
	}

	if err := h.svc.Update(userID, updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "user tastes updated")
}
