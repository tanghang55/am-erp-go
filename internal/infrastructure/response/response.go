package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedData 分页数据结构
type PaginatedData struct {
	Data     interface{} `json:"data"`
	Total    int64       `json:"total"`
	Page     int         `json:"page,omitempty"`
	PageSize int         `json:"page_size,omitempty"`
}

// Success 成功响应（单条数据或无数据）
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessMessage 成功响应（自定义消息）
func SuccessMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（SuccessMessage 的别名）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	SuccessMessage(c, message, data)
}

// Paginated 分页响应
func Paginated(c *gin.Context, list interface{}, total int64) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PaginatedData{
			Data:  list,
			Total: total,
		},
	})
}

// SuccessPage 分页响应（包含页码，可选page_size）
func SuccessPage(c *gin.Context, list interface{}, total int64, page int, pageSize ...int) {
	data := PaginatedData{
		Data:  list,
		Total: total,
		Page:  page,
	}
	if len(pageSize) > 0 {
		data.PageSize = pageSize[0]
	}

	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, httpCode int, message string) {
	c.JSON(httpCode, Response{
		Code:    httpCode,
		Message: message,
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// NotFound 404错误
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// InternalError 500错误
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}

// ServerError 500错误（InternalError 别名）
func ServerError(c *gin.Context, message string) {
	InternalError(c, message)
}
