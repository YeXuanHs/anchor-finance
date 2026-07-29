package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type AppStoreHandler struct {
	appSvc *service.AppStoreService
	log    *logger.Logger
}

func NewAppStoreHandler(appSvc *service.AppStoreService, log *logger.Logger) *AppStoreHandler {
	return &AppStoreHandler{appSvc: appSvc, log: log}
}

// ---------- Public ----------

// GetList returns paginated published apps.
// GET /apps
func (h *AppStoreHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	apps, total, err := h.appSvc.GetList(page, pageSize, category, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, apps, total, page, pageSize)
}

// GetDetail returns a single app by ID.
// GET /apps/:id
func (h *AppStoreHandler) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid app id")
		return
	}

	app, err := h.appSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "app not found")
		return
	}
	response.Success(c, app)
}

// Install installs an app for the current user.
// POST /apps/:id/install
func (h *AppStoreHandler) Install(c *gin.Context) {
	userID := c.GetUint("user_id")
	appID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid app id")
		return
	}

	install, err := h.appSvc.Install(userID, uint(appID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, install)
}

// Uninstall uninstalls an app for the current user.
// POST /apps/:id/uninstall
func (h *AppStoreHandler) Uninstall(c *gin.Context) {
	userID := c.GetUint("user_id")
	appID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid app id")
		return
	}

	if err := h.appSvc.Uninstall(userID, uint(appID)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "app uninstalled")
}

// Update updates an installed app to the latest version.
// POST /apps/:id/update
func (h *AppStoreHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	appID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid app id")
		return
	}

	install, err := h.appSvc.Update(userID, uint(appID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, install)
}

// GetInstalled returns all installed apps for the current user.
// GET /apps/installed
func (h *AppStoreHandler) GetInstalled(c *gin.Context) {
	userID := c.GetUint("user_id")

	installs, err := h.appSvc.GetInstalledApps(userID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, installs)
}

// CreateReview creates or updates a review for an app.
// POST /apps/:id/review
func (h *AppStoreHandler) CreateReview(c *gin.Context) {
	userID := c.GetUint("user_id")
	appID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid app id")
		return
	}

	var req service.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	review, err := h.appSvc.CreateReview(userID, uint(appID), req.Rating, req.Content)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, review)
}

// GetReviews returns paginated reviews for an app.
// GET /apps/:id/reviews
func (h *AppStoreHandler) GetReviews(c *gin.Context) {
	appID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid app id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	reviews, total, err := h.appSvc.GetReviews(uint(appID), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, reviews, total, page, pageSize)
}

// ---------- Admin ----------

// AdminGetList returns all apps including inactive (admin).
// GET /admin/apps
func (h *AppStoreHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	apps, total, err := h.appSvc.AdminGetList(page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, apps, total, page, pageSize)
}

// AdminGetDetail returns a single app by ID (admin).
// GET /admin/apps/:id
func (h *AppStoreHandler) AdminGetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid app id")
		return
	}

	app, err := h.appSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "app not found")
		return
	}
	response.Success(c, app)
}

// AdminCreate creates a new app (admin).
// POST /admin/apps
func (h *AppStoreHandler) AdminCreate(c *gin.Context) {
	var req service.CreateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	app, err := h.appSvc.Create(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, app)
}

// AdminUpdate updates an app (admin).
// PUT /admin/apps/:id
func (h *AppStoreHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid app id")
		return
	}

	var req service.UpdateAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	app, err := h.appSvc.Update(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, app)
}

// AdminList is an alias for AdminGetList.
func (h *AppStoreHandler) AdminList(c *gin.Context) { h.AdminGetList(c) }

// AdminDelete deletes an app (admin).
// DELETE /admin/apps/:id
func (h *AppStoreHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid app id")
		return
	}

	if err := h.appSvc.Delete(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "app deleted")
}
