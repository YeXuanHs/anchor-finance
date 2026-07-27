package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard API response envelope.
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData wraps paginated results.
type PageData struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

func JSON(c *gin.Context, httpCode int, resp Response) {
	c.JSON(httpCode, resp)
}

func Success(c *gin.Context, data interface{}) {
	JSON(c, http.StatusOK, Response{Code: 0, Message: "success", Data: data})
}

func SuccessMsg(c *gin.Context, msg string) {
	JSON(c, http.StatusOK, Response{Code: 0, Message: msg})
}

func SuccessPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	JSON(c, http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PageData{
			List:     list,
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	})
}

func Error(c *gin.Context, httpCode int, code int, msg string) {
	JSON(c, httpCode, Response{Code: code, Message: msg})
}

func BadRequest(c *gin.Context, msg string) {
	Error(c, http.StatusBadRequest, 400, msg)
}

func Unauthorized(c *gin.Context, msg string) {
	Error(c, http.StatusUnauthorized, 401, msg)
}

func Forbidden(c *gin.Context, msg string) {
	Error(c, http.StatusForbidden, 403, msg)
}

func NotFound(c *gin.Context, msg string) {
	Error(c, http.StatusNotFound, 404, msg)
}

func ServerError(c *gin.Context, msg string) {
	Error(c, http.StatusInternalServerError, 500, msg)
}
