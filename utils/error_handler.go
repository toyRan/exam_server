package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

// ValidationError 用于返回验证错误的详细信息
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// HandleValidationError 处理验证错误
func HandleValidationError(err error) []ValidationError {
	var validationErrors []ValidationError
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, e := range ve {
			error := ValidationError{
				Field:   e.Field(),
				Message: getValidationErrorMsg(e),
			}
			validationErrors = append(validationErrors, error)
		}
	}
	return validationErrors
}

// getValidationErrorMsg 获取验证错误的具体信息
func getValidationErrorMsg(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "此字段不能为空"
	case "email":
		return "请输入有效的电子邮件地址"
	case "min":
		return "不能小于 " + e.Param()
	case "max":
		return "不能大于 " + e.Param()
	// 添加更多验证规则的错误信息...
	default:
		return "验证错误"
	}
}

// // DBError 数据库错误结构
// type DBError struct {
// 	Field   string `json:"field,omitempty"`
// 	Message string `json:"message"`
// 	Type    string `json:"type,omitempty"`
// }

// // HandleMySQLError 处理 MySQL 错误，返回结构化的错误信息
// func HandleMySQLError(err error) (string, int, interface{}) {
// 	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
// 		switch mysqlErr.Number {
// 		case 1062:
// 			field := extractDuplicateField(mysqlErr.Message)
// 			return "数据重复错误", http.StatusBadRequest, DBError{
// 				Field:   field,
// 				Message: fmt.Sprintf("字段 '%s' 的值已存在", field),
// 				Type:    "DUPLICATE_ERROR",
// 			}
// 		case 1452:
// 			return "外键约束错误", http.StatusBadRequest, DBError{
// 				Message: "关联记录不存在",
// 				Type:    "FOREIGN_KEY_ERROR",
// 			}
// 		default:
// 			return "数据库错误", http.StatusInternalServerError, DBError{
// 				Message: mysqlErr.Message,
// 				Type:    "DB_ERROR",
// 			}
// 		}
// 	}
// 	return "数据库错误", http.StatusInternalServerError, nil
// }

// HandleMySQLError 处理 MySQL 错误，直接返回用户友好的错误消息
func HandleMySQLError(err error) (message string, statusCode int) {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		log.Printf("err:%v", err)
		switch mysqlErr.Number {
		case 1062:
			field := extractDuplicateField(mysqlErr.Message)
			return fmt.Sprintf("'%s' 已存在，请修改后重试", field), http.StatusBadRequest
		case 1452:
			return "关联数据不存在", http.StatusBadRequest
		default:
			return "数据库操作失败，请稍后重试", http.StatusInternalServerError
		}
	}

	return "数据库操作失败，请稍后重试", http.StatusInternalServerError
}

// extractDuplicateField 从 MySQL 错误信息中提取重复的字段名
func extractDuplicateField(errMsg string) string {
	// MySQL 错误信息格式类似：
	// Duplicate entry 'xxx' for key 'table_name.field_name'
	// 或 Duplicate entry 'xxx' for key 'field_name'

	// 1. 先找到最后一个点号后的内容
	parts := strings.Split(errMsg, ".")
	lastPart := parts[len(parts)-1]

	// 2. 找到 key 后面的字段名（去掉可能的引号）
	keyParts := strings.Split(lastPart, "'")
	if len(keyParts) > 1 {
		return strings.TrimSpace(keyParts[len(keyParts)-2])
	}

	// 如果没有点号，直接找 key 后面的内容
	keyIndex := strings.Index(errMsg, "key '")
	if keyIndex != -1 {
		fieldName := errMsg[keyIndex+5:]
		fieldName = strings.TrimRight(fieldName, "'")
		return fieldName
	}

	return "未知字段"
}
