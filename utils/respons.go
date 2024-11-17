package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SuccessResponse 用于返回成功的响应
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	response := gin.H{
		"code":    1,
		"message": message,
	}

	// 如果有数据，添加到响应中
	if data != nil {
		response["result"] = data
	}

	c.JSON(http.StatusOK, response)
}

// ErrorResponse 用于返回失败的响应
func ErrorResponse(c *gin.Context, message string, status int, errors ...interface{}) {
	response := gin.H{
		"code":    0,
		"message": message,
	}

	// 如果提供了错误信息，添加到响应中
	if len(errors) > 0 {
		response["errors"] = errors
	}

	c.JSON(status, response)
}
