package admin

import (
	"net/http"
	"strconv"

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
		c.JSON(http.StatusOK, gin.H{
			"code":    404,
			"message": "文件不存在",
		})
		return
	}

	// TODO: 删除实际文件

	db.Delete(&file)

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "删除成功",
	})
}

// GetMediaFileReferences 获取媒体文件引用
// GET /api/admin/media-files/:id/references
func GetMediaFileReferences(c *gin.Context) {
	// TODO: 查找引用该文件的其他资源
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    []interface{}{},
	})
}
