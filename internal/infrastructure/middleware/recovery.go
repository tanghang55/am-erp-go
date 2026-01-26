package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"am-erp-go/internal/infrastructure/response"

	"github.com/gin-gonic/gin"
)

// Recovery 全局异常恢复中间件
// 捕获panic并返回统一的错误响应
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印堆栈信息
				log.Printf("[Recovery] panic recovered: %v\n%s", err, debug.Stack())

				// 返回统一的错误响应
				c.AbortWithStatusJSON(http.StatusInternalServerError, response.Response{
					Code:    500,
					Message: "服务器内部错误",
					Success: false,
				})
			}
		}()
		c.Next()
	}
}

// ErrorHandler 错误处理中间件
// 用于处理handler中设置的错误
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 检查是否有错误
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("[ErrorHandler] error: %v", err)

			// 如果还没有写入响应，则写入错误响应
			if !c.Writer.Written() {
				response.ServerError(c, err.Error())
			}
		}
	}
}
