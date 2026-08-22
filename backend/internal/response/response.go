package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

// SuccessMsg 成功响应（自定义消息）
func SuccessMsg(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": message,
		"data":    nil,
	})
}

// SuccessPage 分页成功响应
func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data": gin.H{
			"list":      list,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

// Forbidden 403错误
func Forbidden(c *gin.Context, message string) {
	Error(c, 403, message)
}

// NotFound 404错误
func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// ServerError 500错误
func ServerError(c *gin.Context, message string) {
	Error(c, 500, message)
}
