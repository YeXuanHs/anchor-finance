package handler

import (
	"strconv"

	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

// UploadHandler handles file upload HTTP requests.
type UploadHandler struct {
	uploadSvc *service.UploadService
	log       *logger.Logger
}

// NewUploadHandler creates a new UploadHandler.
func NewUploadHandler(uploadSvc *service.UploadService, log *logger.Logger) *UploadHandler {
	return &UploadHandler{uploadSvc: uploadSvc, log: log}
}

// Upload handles single file upload.
// POST /upload
func (h *UploadHandler) Upload(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	relType := c.PostForm("rel_type")
	if relType == "" {
		relType = "general"
	}
	relID, _ := strconv.ParseUint(c.PostForm("rel_id"), 10, 64)

	uploaded, err := h.uploadSvc.Upload(userID, file, relType, uint(relID))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":           uploaded.ID,
		"file_name":    uploaded.FileName,
		"original_name": uploaded.OriginalName,
		"file_size":    uploaded.FileSize,
		"mime_type":    uploaded.MimeType,
		"url":          h.uploadSvc.GetURL(uploaded),
	})
}

// UploadAvatar handles avatar upload.
// POST /upload/avatar
func (h *UploadHandler) UploadAvatar(c *gin.Context) {
	userID := c.GetUint("user_id")

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}

	uploaded, err := h.uploadSvc.Upload(userID, file, "avatar", 0)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"id":  uploaded.ID,
		"url": h.uploadSvc.GetURL(uploaded),
	})
}

// GetList returns paginated uploaded files for the user.
// GET /upload
func (h *UploadHandler) GetList(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	relType := c.Query("rel_type")

	files, total, err := h.uploadSvc.GetList(userID, page, pageSize, relType)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	type fileInfo struct {
		ID           uint   `json:"id"`
		FileName     string `json:"file_name"`
		OriginalName string `json:"original_name"`
		FileSize     int64  `json:"file_size"`
		MimeType     string `json:"mime_type"`
		RelType      string `json:"rel_type"`
		URL          string `json:"url"`
		CreatedAt    string `json:"created_at"`
	}

	result := make([]fileInfo, 0, len(files))
	for _, f := range files {
		result = append(result, fileInfo{
			ID:           f.ID,
			FileName:     f.FileName,
			OriginalName: f.OriginalName,
			FileSize:     f.FileSize,
			MimeType:     f.MimeType,
			RelType:      f.RelType,
			URL:          h.uploadSvc.GetURL(&f),
			CreatedAt:    f.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	response.SuccessPage(c, result, total, page, pageSize)
}

// GetDetail returns a single file's details.
// GET /upload/:id
func (h *UploadHandler) GetDetail(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	file, err := h.uploadSvc.GetByID(userID, uint(id))
	if err != nil {
		response.NotFound(c, "file not found")
		return
	}

	response.Success(c, gin.H{
		"id":           file.ID,
		"file_name":    file.FileName,
		"original_name": file.OriginalName,
		"file_size":    file.FileSize,
		"mime_type":    file.MimeType,
		"extension":    file.Extension,
		"rel_type":     file.RelType,
		"rel_id":       file.RelID,
		"url":          h.uploadSvc.GetURL(file),
		"created_at":   file.CreatedAt,
	})
}

// Delete deletes a file.
// DELETE /upload/:id
func (h *UploadHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid file id")
		return
	}

	if err := h.uploadSvc.Delete(userID, uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "file deleted")
}
