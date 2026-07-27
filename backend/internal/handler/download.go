package handler

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type DownloadHandler struct {
	downloadSvc *service.DownloadService
	log         *logger.Logger
}

func NewDownloadHandler(downloadSvc *service.DownloadService, log *logger.Logger) *DownloadHandler {
	return &DownloadHandler{downloadSvc: downloadSvc, log: log}
}

// ---------- Public ----------

// GetCategories returns all active download categories.
func (h *DownloadHandler) GetCategories(c *gin.Context) {
	cats, err := h.downloadSvc.GetCategories()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cats)
}

// GetFiles returns paginated published download files.
func (h *DownloadHandler) GetFiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var categoryID uint
	if cid := c.Query("category_id"); cid != "" {
		v, _ := strconv.ParseUint(cid, 10, 64)
		categoryID = uint(v)
	}

	files, total, err := h.downloadSvc.GetFiles(page, pageSize, categoryID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, files, total, page, pageSize)
}

// GetFile returns a single download file by ID.
func (h *DownloadHandler) GetFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	file, err := h.downloadSvc.GetFileByID(uint(id))
	if err != nil {
		response.NotFound(c, "file not found")
		return
	}
	response.Success(c, file)
}

// Download serves the actual file and increments download count.
func (h *DownloadHandler) Download(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	file, err := h.downloadSvc.GetFileByID(uint(id))
	if err != nil {
		response.NotFound(c, "file not found")
		return
	}

	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		response.NotFound(c, "file not found on disk")
		return
	}

	_ = h.downloadSvc.IncrementDownload(uint(id))

	filename := filepath.Base(file.FilePath)
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")
	c.File(file.FilePath)
}

// ---------- Admin ----------

// AdminGetCategories returns all download categories including inactive (admin).
func (h *DownloadHandler) AdminGetCategories(c *gin.Context) {
	cats, err := h.downloadSvc.GetCategories()
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cats)
}

// AdminCreateCategory creates a download category (admin).
func (h *DownloadHandler) AdminCreateCategory(c *gin.Context) {
	var req service.CreateDownloadCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cat, err := h.downloadSvc.CreateCategory(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cat)
}

// AdminUpdateCategory updates a download category (admin).
func (h *DownloadHandler) AdminUpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	var req service.UpdateDownloadCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cat, err := h.downloadSvc.UpdateCategory(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cat)
}

// AdminDeleteCategory deletes a download category (admin).
func (h *DownloadHandler) AdminDeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	if err := h.downloadSvc.DeleteCategory(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "category deleted")
}

// AdminGetFiles returns all files including unpublished (admin).
func (h *DownloadHandler) AdminGetFiles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var categoryID uint
	if cid := c.Query("category_id"); cid != "" {
		v, _ := strconv.ParseUint(cid, 10, 64)
		categoryID = uint(v)
	}

	files, total, err := h.downloadSvc.AdminGetFiles(page, pageSize, categoryID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, files, total, page, pageSize)
}

// AdminGetFile returns a single download file by ID (admin).
func (h *DownloadHandler) AdminGetFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	file, err := h.downloadSvc.AdminGetFile(uint(id))
	if err != nil {
		response.NotFound(c, "file not found")
		return
	}
	response.Success(c, file)
}

// AdminCreateFile creates a download file (admin).
func (h *DownloadHandler) AdminCreateFile(c *gin.Context) {
	var req service.CreateDownloadFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	file, err := h.downloadSvc.CreateFile(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, file)
}

// AdminUpdateFile updates a download file (admin).
func (h *DownloadHandler) AdminUpdateFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	var req service.UpdateDownloadFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	file, err := h.downloadSvc.UpdateFile(uint(id), req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, file)
}

// AdminDeleteFile deletes a download file (admin).
func (h *DownloadHandler) AdminDeleteFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	if err := h.downloadSvc.DeleteFile(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "file deleted")
}

// AdminUploadFile handles multipart file upload (admin).
func (h *DownloadHandler) AdminUploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	categoryID, _ := strconv.ParseUint(c.PostForm("category_id"), 10, 64)
	title := c.PostForm("title")
	if title == "" {
		title = fileHeader.Filename
	}

	saveDir := "./uploads/downloads"
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		response.ServerError(c, "failed to create upload directory")
		return
	}

	savePath := filepath.Join(saveDir, fileHeader.Filename)
	if err := c.SaveUploadedFile(fileHeader, savePath); err != nil {
		response.ServerError(c, "failed to save file")
		return
	}

	req := service.CreateDownloadFileRequest{
		CategoryID:  uint(categoryID),
		Title:       title,
		Description: c.PostForm("description"),
		FilePath:    savePath,
		FileSize:    fileHeader.Size,
		FileType:    filepath.Ext(fileHeader.Filename),
		IsPublished: c.PostForm("is_published") == "true",
	}

	dlFile, err := h.downloadSvc.CreateFile(req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "file uploaded",
		Data:    dlFile,
	})
}
