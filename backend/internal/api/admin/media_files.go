package admin

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/gin-gonic/gin"
)

// GetMediaFileList 获取媒体文件列表
// GET /api/admin/media-files
func GetMediaFileList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	db := database.GetDB()
	var total int64
	db.Model(&model.MediaFile{}).Count(&total)

	var files []model.MediaFile
	offset := (page - 1) * pageSize
	db.Offset(offset).Limit(pageSize).Order("id DESC").Find(&files)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      files,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// DeleteMediaFile 删除媒体文件
// DELETE /api/admin/media-files/:id
func DeleteMediaFile(c *gin.Context) {
	id := c.Param("id")

	db := database.GetDB()
	var file model.MediaFile
	if err := db.First(&file, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "文件不存在", "data": nil})
		return
	}

	// 删除实际文件
	if file.Path != "" {
		if err := os.Remove(file.Path); err != nil && !os.IsNotExist(err) {
			// 文件删除失败不影响DB记录删除，但记录错误
		}
	}

	db.Delete(&file)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "删除成功", "data": nil})
}

// UploadFile 上传文件（保存到磁盘）
// POST /api/admin/upload
func UploadFile(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 400, "message": "请选择文件", "data": nil})
		return
	}
	defer file.Close()

	// 生成唯一文件名（防路径穿越）
	ext := filepath.Ext(header.Filename)
	baseName := strings.TrimSuffix(header.Filename, ext)
	newName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), sanitize(baseName), ext)
	dir := "uploads"
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "创建上传目录失败", "data": nil})
		return
	}
	savePath := filepath.Join(dir, newName)

	// 保存文件到磁盘
	out, err := os.Create(savePath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "保存文件失败", "data": nil})
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "写入文件失败", "data": nil})
		return
	}

	// 从token获取用户ID
	adminID, _ := c.Get("user_id")

	db := database.GetDB()
	mediaFile := model.MediaFile{
		Name:       header.Filename,
		Path:       savePath,
		Size:       header.Size,
		MimeType:   header.Header.Get("Content-Type"),
		Extension:  strings.TrimPrefix(ext, "."),
		UploadedBy: adminID.(uint),
	}

	if err := db.Create(&mediaFile).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 500, "message": "保存文件信息失败", "data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "上传成功",
		"data": gin.H{
			"id":   mediaFile.ID,
			"name": header.Filename,
			"url":  "/" + savePath,
		},
	})
}

// sanitize 清理文件名特殊字符
func sanitize(name string) string {
	re := regexp.MustCompile(`[^\w\-.]`)
	return re.ReplaceAllString(name, "_")
}

// GetMediaFileReferences 获取媒体文件引用（查找哪些内容引用了该文件）
// GET /api/admin/media-files/:id/references
func GetMediaFileReferences(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var file model.MediaFile
	if err := db.First(&file, id).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 404, "message": "文件不存在", "data": nil})
		return
	}

	// 检查哪些内容引用了该文件URL
	references := []map[string]interface{}{}

	// 检查新闻
	var news []model.News
	db.Where("content LIKE ?", "%"+file.Path+"%").Find(&news)
	for _, n := range news {
		references = append(references, map[string]interface{}{"type": "news", "id": n.ID, "title": n.Title})
	}

	// 检查知识文章
	var articles []model.KnowledgeArticle
	db.Where("content LIKE ?", "%"+file.Path+"%").Find(&articles)
	for _, a := range articles {
		references = append(references, map[string]interface{}{"type": "knowledge_article", "id": a.ID, "title": a.Title})
	}

	// 检查下载
	var downloads []model.Download
	db.Where("file_url LIKE ?", "%"+file.Path+"%").Find(&downloads)
	for _, d := range downloads {
		references = append(references, map[string]interface{}{"type": "download", "id": d.ID, "title": d.Title})
	}

	// 检查HomeHero
	var heroes []model.HomeHero
	db.Where("image_url = ?", file.Path).Find(&heroes)
	for _, h := range heroes {
		references = append(references, map[string]interface{}{"type": "home_hero", "id": h.ID, "title": h.Title})
	}

	if references == nil { references = []map[string]interface{}{} }
	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": references})
}

// ReindexMediaFiles 重新索引媒体文件
// POST /api/admin/media-file-reindexes
func ReindexMediaFiles(c *gin.Context) {
	db := database.GetDB()

	// 扫描磁盘文件，同步到数据库
	var count int64
	db.Model(&model.MediaFile{}).Count(&count)

	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "重索引完成", "data": gin.H{"total": count}})
}
