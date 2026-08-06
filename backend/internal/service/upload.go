package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// UploadedFile 上传文件记录
type UploadedFile struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	FileName      string         `gorm:"type:varchar(256);not null" json:"file_name"`
	OriginalName  string         `gorm:"type:varchar(256);not null" json:"original_name"`
	FilePath      string         `gorm:"type:varchar(512);not null" json:"file_path"`
	FileSize      int64          `gorm:"not null" json:"file_size"`
	MimeType      string         `gorm:"type:varchar(128)" json:"mime_type"`
	Extension     string         `gorm:"type:varchar(16)" json:"extension"`
	Hash          string         `gorm:"type:varchar(64);index" json:"hash"`
	StorageDriver string         `gorm:"type:varchar(32);default:'local'" json:"storage_driver"`
	RelType       string         `gorm:"type:varchar(32);index" json:"rel_type"` // avatar/ticket/general
	RelID         uint           `gorm:"index" json:"rel_id"`
	DownloadCount int            `gorm:"default:0" json:"download_count"`
	IsPublic      bool           `gorm:"default:false" json:"is_public"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type UploadService struct {
	db        *gorm.DB
	log       *logger.Logger
	uploadDir string
	baseURL   string
}

func NewUploadService(db *gorm.DB, log *logger.Logger, uploadDir, baseURL string) *UploadService {
	if uploadDir == "" {
		uploadDir = "./uploads"
	}
	if baseURL == "" {
		baseURL = "/uploads"
	}
	return &UploadService{db: db, log: log, uploadDir: uploadDir, baseURL: baseURL}
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/gif":       true,
	"image/webp":      true,
	"application/pdf": true,
	"application/zip": true,
	"application/x-rar-compressed": true,
	"application/x-7z-compressed":  true,
	"text/plain":      true,
	"application/msword":                                    true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true,
	"application/vnd.ms-excel":                                                  true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
}

const maxFileSize = 20 * 1024 * 1024 // 20MB

// Upload handles file upload.
func (s *UploadService) Upload(userID uint, file *multipart.FileHeader, relType string, relID uint) (*UploadedFile, error) {
	if file.Size > maxFileSize {
		return nil, fmt.Errorf("file too large, max %dMB", maxFileSize/1024/1024)
	}

	// Validate mime type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	mimeType := file.Header.Get("Content-Type")
	if !allowedMimeTypes[mimeType] {
		return nil, errors.New("file type not allowed")
	}

	// Read file for hash
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	hasher := sha256.New()
	content, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	hasher.Write(content)
	hash := hex.EncodeToString(hasher.Sum(nil))

	// Generate storage path: uploads/2026/07/userID/filename
	now := time.Now()
	relDir := filepath.Join(now.Format("2006"), now.Format("01"))
	storageDir := filepath.Join(s.uploadDir, relDir, fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, err
	}

	// Generate unique filename
	fileName := fmt.Sprintf("%s%s",
		fmt.Sprintf("%d%s", now.UnixNano(), fmt.Sprintf("%04d", now.UnixNano()%10000)),
		ext)

	filePath := filepath.Join(storageDir, fileName)

	// Write file
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := dst.Write(content); err != nil {
		return nil, err
	}

	// Save record
	record := &UploadedFile{
		UserID:       userID,
		FileName:     fileName,
		OriginalName: file.Filename,
		FilePath:     filepath.Join(relDir, fmt.Sprintf("%d", userID), fileName),
		FileSize:     file.Size,
		MimeType:     mimeType,
		Extension:    ext,
		Hash:         hash,
		RelType:      relType,
		RelID:        relID,
		IsPublic:     relType == "avatar",
	}

	if err := s.db.Create(record).Error; err != nil {
		os.Remove(filePath)
		return nil, err
	}

	s.log.Infof("file uploaded: user=%d file=%s size=%d", userID, file.Filename, file.Size)
	return record, nil
}

// GetList returns paginated uploaded files for a user.
func (s *UploadService) GetList(userID uint, page, pageSize int, relType string) ([]UploadedFile, int64, error) {
	var files []UploadedFile
	var total int64

	query := s.db.Model(&UploadedFile{}).Where("user_id = ?", userID)
	if relType != "" {
		query = query.Where("rel_type = ?", relType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&files).Error; err != nil {
		return nil, 0, err
	}
	return files, total, nil
}

// GetByID returns a file record by ID.
func (s *UploadService) GetByID(userID, fileID uint) (*UploadedFile, error) {
	var file UploadedFile
	if err := s.db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

// Delete deletes a file and its record.
func (s *UploadService) Delete(userID, fileID uint) error {
	var file UploadedFile
	if err := s.db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		return errors.New("file not found")
	}

	// Delete physical file
	fullPath := filepath.Join(s.uploadDir, file.FilePath)
	os.Remove(fullPath)

	// Delete record
	return s.db.Delete(&file).Error
}

// GetURL returns the full URL for a file.
func (s *UploadService) GetURL(file *UploadedFile) string {
	return s.baseURL + "/" + file.FilePath
}

// UploadByType handles file upload with specific type validation.
func (s *UploadService) UploadByType(file *multipart.FileHeader, relType string) (*UploadedFile, error) {
	return s.Upload(0, file, relType, 0)
}
